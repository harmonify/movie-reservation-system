package mongo_repository

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/config"
	mongo_pkg "github.com/harmonify/movie-reservation-system/pkg/database/mongo"
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
	*config.MovieSearchServiceConfig
	logger.Logger
	tracer.Tracer
	*mongo_pkg.MongoClient
	*util.Util
}

type movieMongoRepository struct {
	config          *config.MovieSearchServiceConfig
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
		movieCollection: p.Client.Database(p.MovieSearchServiceConfig.MongoDbName).Collection(movieCollectionName),
	}
}

// FacetResult is used to decode the output of our aggregation pipeline.
type FacetResult struct {
	Data       []*entity.SearchMovieResult `bson:"data"`
	TotalCount []struct {
		Count int64 `bson:"count"`
	} `bson:"totalCount"`
}

func (r *movieMongoRepository) SearchMovies(ctx context.Context, searchModel *entity.SearchMovieParam) (*shared.SearchMovieResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var lastSeenID bson.ObjectID
	var err error
	if searchModel.LastSeenID != "" {
		lastSeenID, err = bson.ObjectIDFromHex(searchModel.LastSeenID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert last seen ID: %w", err)
		}
	}

	attributeFilters := bson.A{}

	// Filter by movie IDs.
	if searchModel.MovieIDs != nil {
		if len(searchModel.MovieIDs) == 0 {
			return &shared.SearchMovieResult{
				Data: []*entity.SearchMovieResult{},
				Meta: &shared.SearchMovieMetadata{
					TotalCount:        0,
					HasNextPage:       false,
					LastSeenSortValue: nil,
					LastSeenID:        "",
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

	// Filter by last seen cursor. Used in STAGE 3.A.1.
	var lastSeenFilter bson.D
	if searchModel.LastSeenSortValue != nil && !lastSeenID.IsZero() {
		if searchModel.Keyword.Valid {
			lastSeenFilter = bson.D{{Key: "$or", Value: bson.A{
				bson.D{
					{Key: "score", Value: bson.D{
						{Key: "$lt", Value: searchModel.LastSeenSortValue},
					}},
				},
				bson.D{
					{Key: "score", Value: bson.D{
						{Key: "$eq", Value: searchModel.LastSeenSortValue},
					}},
					{Key: "_id", Value: bson.D{
						{Key: "$gt", Value: lastSeenID},
					}},
				},
			}}}
		} else {
			switch searchModel.SortBy {
			case entity.MovieSortByReleaseDateDesc:
				lastSeenFilter = bson.D{{Key: "$or", Value: bson.A{
					bson.D{
						{Key: "release_date", Value: bson.D{
							{Key: "$lt", Value: searchModel.LastSeenSortValue},
						}},
					},
					bson.D{
						{Key: "release_date", Value: bson.D{
							{Key: "$eq", Value: searchModel.LastSeenSortValue},
						}},
						{Key: "_id", Value: bson.D{
							{Key: "$gt", Value: lastSeenID},
						}},
					},
				}}}
			case entity.MovieSortByTitleDesc:
				lastSeenFilter = bson.D{{Key: "$or", Value: bson.A{
					bson.D{
						{Key: "title", Value: bson.D{
							{Key: "$lt", Value: searchModel.LastSeenSortValue},
						}},
					},
					bson.D{
						{Key: "title", Value: bson.D{
							{Key: "$eq", Value: searchModel.LastSeenSortValue},
						}},
						{Key: "_id", Value: bson.D{
							{Key: "$gt", Value: lastSeenID},
						}},
					},
				}}}
			default: // entity.MovieSortByTitleAsc
				lastSeenFilter = bson.D{{Key: "$or", Value: bson.A{
					bson.D{
						{Key: "title", Value: bson.D{
							{Key: "$gt", Value: searchModel.LastSeenSortValue},
						}},
					},
					bson.D{
						{Key: "title", Value: bson.D{
							{Key: "$eq", Value: searchModel.LastSeenSortValue},
						}},
						{Key: "_id", Value: bson.D{
							{Key: "$gt", Value: lastSeenID},
						}},
					},
				}}}
			}
		}
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
				// STAGE 3.A.1: MATCH by last seen cursor.
				bson.D{{Key: "$match", Value: lastSeenFilter}}, // Further filter the data result based on the last seen cursor.
				// STAGE 3.A.2: SORT by sort order.
				bson.D{{Key: "$sort", Value: sortValue}},
				// STAGE 3.A.3: LIMIT by limit + 1.
				// Limit the number of results to limit + 1 to determine if there's a next page.
				bson.D{{Key: "$limit", Value: searchModel.Limit + 1}},
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
	if err = cursor.All(ctx, &results); err != nil {
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

	// Determine if there's a next page.
	hasNextPage := false
	if len(movies) > int(searchModel.Limit) {
		hasNextPage = true
		// Remove the extra record.
		movies = movies[:int(searchModel.Limit)]
	}

	var lastSeenSortValue interface{}
	if len(movies) > 0 {
		lastSeenID = movies[len(movies)-1].MovieID
		if searchModel.Keyword.Valid {
			lastSeenSortValue = movies[len(movies)-1].Score
		} else {
			switch searchModel.SortBy {
			case entity.MovieSortByReleaseDateDesc:
				lastSeenSortValue = movies[len(movies)-1].ReleaseDate
			case entity.MovieSortByTitleDesc:
				lastSeenSortValue = movies[len(movies)-1].Title
			default: // entity.MovieSortByTitleAsc
				lastSeenSortValue = movies[len(movies)-1].Title
			}
		}
	}

	return &shared.SearchMovieResult{
		Data: movies,
		Meta: &shared.SearchMovieMetadata{
			TotalCount:        totalCount,
			HasNextPage:       hasNextPage,
			LastSeenSortValue: lastSeenSortValue,
			LastSeenID:        lastSeenID.Hex(),
		},
	}, nil
}
