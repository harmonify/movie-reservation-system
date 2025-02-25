package mongo_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	mongo_pkg "github.com/harmonify/movie-reservation-system/pkg/database/mongo"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MovieMongoRepositoryParam struct {
	fx.In
	*config.MovieServiceConfig
	logger.Logger
	tracer.Tracer
	*mongo_pkg.MongoClient
	*util.Util
}

type movieMongoRepository struct {
	config          *config.MovieServiceConfig
	logger          logger.Logger
	tracer          tracer.Tracer
	util            *util.Util
	client          *mongo_pkg.MongoClient
	movieCollection *mongo.Collection
}

var movieCollectionName = "movies"

func NewMovieMongoRepository(p MovieMongoRepositoryParam) shared.MovieStorage {
	return &movieMongoRepository{
		logger:          p.Logger,
		tracer:          p.Tracer,
		util:            p.Util,
		client:          p.MongoClient,
		movieCollection: p.Client.Database(p.MovieServiceConfig.MongoDbName).Collection(movieCollectionName),
	}
}

// FacetResult is used to decode the output of our aggregation pipeline.
type FacetResult struct {
	Data       []*entity.SearchMovieResult `bson:"data"`
	TotalCount []struct {
		Count int64 `bson:"count"`
	} `bson:"totalCount"`
}

func (r *movieMongoRepository) SearchMovies(ctx context.Context, searchModel *entity.SearchMovie) (*shared.SearchMovieResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	attributeFilters := bson.A{}

	// Filter by movie IDs.
	if searchModel.MovieIDs != nil {
		if len(searchModel.MovieIDs) == 0 {
			return &shared.SearchMovieResult{
				Data: []*entity.SearchMovieResult{},
				Meta: &shared.SearchMovieMetadata{
					TotalCount:  0,
					HasNextPage: false,
				},
			}, nil
		} else {
			attributeFilters = append(attributeFilters, bson.D{{Key: "movie_id", Value: bson.D{
				{Key: "$in", Value: searchModel.MovieIDs},
			}}})
		}
	}

	// Filter by genres.
	if searchModel.Genre.Valid {
		attributeFilters = append(attributeFilters, bson.D{{Key: "genres", Value: bson.D{
			{Key: "$regex", Value: searchModel.Genre.String},
			{Key: "$options", Value: "i"},
		}}})
	}

	// Filter by keyword.
	if searchModel.Keyword.Valid {
		attributeFilters = append(attributeFilters, bson.D{{Key: "$text", Value: bson.D{
			{Key: "$search", Value: searchModel.Keyword.String},
		}}})
	}

	if searchModel.ReleaseDateFrom.Valid {
		attributeFilters = append(attributeFilters, bson.D{{Key: "release_date", Value: bson.D{
			{Key: "$gte", Value: searchModel.ReleaseDateFrom.Time},
		}}})
	}

	if searchModel.ReleaseDateTo.Valid {
		attributeFilters = append(attributeFilters, bson.D{{Key: "release_date", Value: bson.D{
			{Key: "$lte", Value: searchModel.ReleaseDateTo.Time},
		}}})
	}

	// BUILD PIPELINE
	pipeline := mongo.Pipeline{}

	// STAGE 1: MATCH by movie attribute filters.
	// If there are no filters, add an empty filter.
	if len(attributeFilters) == 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{}}})
	} else {
		// Apply the filter.
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "$and", Value: attributeFilters},
		}}})
	}

	// STAGE 2: ADD FIELDS (OPTIONAL)
	// Append the text score if a keyword is provided.
	if searchModel.Keyword.Valid {
		pipeline = append(
			pipeline,
			// Include the text score in the documents if a keyword is provided.
			//   { $addFields: { score: { $meta: "textScore" } } },
			bson.D{{Key: "$addFields", Value: bson.D{
				{Key: "score", Value: bson.D{
					{Key: "$meta", Value: "textScore"},
				}},
			}}},
		)
	}

	// Sort the results based on the sort order. Used in STAGE 3.A.2.
	var sortValue bson.D
	if searchModel.Keyword.Valid {
		sortValue = bson.D{
			{Key: "score", Value: bson.D{
				{Key: "$meta", Value: "textScore"},
			}},
		}
	} else {
		switch searchModel.SortBy {
		case entity.MovieSortByReleaseDateDesc:
			sortValue = bson.D{
				{Key: "release_date", Value: -1},
				{Key: "_id", Value: -1},
			}
		case entity.MovieSortByTitleDesc:
			sortValue = bson.D{
				{Key: "title", Value: -1},
				{Key: "_id", Value: -1},
			}
		default: // entity.MovieSortByTitleAsc
			sortValue = bson.D{
				{Key: "title", Value: 1},
				{Key: "_id", Value: 1},
			}
		}
	}

	pipeline = append(
		pipeline,
		// STAGE 3: FACET to get different parts of the result.
		bson.D{{Key: "$facet", Value: bson.D{
			// STAGE 3.A: FACET for paginated results.
			{Key: "data", Value: bson.A{
				// STAGE 3.A.1: SORT
				bson.D{{Key: "$sort", Value: sortValue}},
				// STAGE 3.A.2: OFFSET
				bson.D{{Key: "$skip", Value: buildOffset(searchModel.Page, searchModel.PageSize)}},
				// STAGE 3.A.3: LIMIT
				bson.D{{Key: "$limit", Value: searchModel.PageSize}},
			}},
			// STAGE 3.B: FACET for total count.
			// The "totalCount" facet for the total number of matching documents
			{Key: "totalCount", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
	)

	r.logger.WithCtx(ctx).Debug("search movies pipeline", zap.Any("pipeline", pipeline))

	cursor, err := r.movieCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to find movies: %w", err)
	}
	defer cursor.Close(ctx)

	// The aggregation result will be a single document containing our two facets.
	var results []FacetResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode movies: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	facetResult := results[0]

	movies := facetResult.Data

	var totalCount int64
	if len(facetResult.TotalCount) > 0 {
		totalCount = facetResult.TotalCount[0].Count
	} else {
		totalCount = 0
	}

	hasNextPage := totalCount > searchModel.Page*searchModel.PageSize

	return &shared.SearchMovieResult{
		Data: movies,
		Meta: &shared.SearchMovieMetadata{
			TotalCount:  totalCount,
			HasNextPage: hasNextPage,
		},
	}, nil
}

func (r *movieMongoRepository) GetMovieByID(ctx context.Context, movieId string) (*entity.Movie, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	movieOid, err := bson.ObjectIDFromHex(movieId)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to convert movie id to object id", zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	res := r.movieCollection.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: movieOid,
		},
	})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, error_pkg.NotFoundError
		}
		r.logger.WithCtx(ctx).Error("failed to find movie", zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	var movie entity.Movie
	if err := res.Decode(&movie); err != nil {
		r.logger.WithCtx(ctx).Error("failed to decode movie", zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	return &movie, nil
}

func (r *movieMongoRepository) SaveMovie(ctx context.Context, saveModel *entity.SaveMovie) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	saveBson, err := r.util.MongoStructUtil.ConvertStructToBsonDocument(ctx, saveModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to convert struct to bson", zap.Error(err))
		return "", error_pkg.InternalServerError
	}

	res, err := r.movieCollection.InsertOne(ctx, saveBson)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to insert movie", zap.Error(err))
		return "", error_pkg.InternalServerError
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		r.logger.WithCtx(ctx).Error("failed to convert inserted id to object id", zap.Any("inserted_id", res.InsertedID))
		return "", error_pkg.InternalServerError
	}

	return id.Hex(), nil
}

func (r *movieMongoRepository) UpdateMovie(ctx context.Context, movieId string, updateModel *entity.UpdateMovie) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// convert to object id
	movieOid, err := bson.ObjectIDFromHex(movieId)
	if err != nil {
		return fmt.Errorf("failed to convert movie id to object id: %w", err)
	}

	updateBson, err := r.util.MongoStructUtil.ConvertStructToBsonDocument(ctx, updateModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to convert struct to bson", zap.Error(err))
		return error_pkg.InternalServerError
	}

	updateResult, err := r.movieCollection.UpdateOne(
		ctx,
		bson.D{
			{
				Key:   "_id",
				Value: movieOid,
			},
		},
		updateBson,
	)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to update movie", zap.Error(err))
		return error_pkg.InternalServerError
	}

	if updateResult.MatchedCount == 0 {
		return error_pkg.NotFoundError
	}

	return nil
}

func (r *movieMongoRepository) SoftDeleteMovie(ctx context.Context, movieId string) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := r.movieCollection.UpdateOne(
		ctx,
		bson.D{
			{
				Key:   "_id",
				Value: movieId,
			},
		},
		bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "deleted_at",
						Value: time.Now(),
					},
				},
			},
		},
	)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to soft delete movie", zap.Error(err))
		return error_pkg.InternalServerError
	}

	if res.MatchedCount == 0 {
		return error_pkg.NotFoundError
	}

	return nil
}

func buildOffset(page, pageSize int64) int64 {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}
