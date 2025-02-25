package repository

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	selectTheaterQuery = "theater.*, ST_X(location) AS longitude, ST_Y(location) AS latitude"
)

type theaterRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewTheaterRepository(
	database *database.Database,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.TheaterStorage {
	return &theaterRepositoryImpl{
		database: database,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *theaterRepositoryImpl) WithTx(tx *database.Transaction) shared.TheaterStorage {
	if tx == nil {
		return r
	}
	return NewTheaterRepository(
		r.database.WithTx(tx),
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *theaterRepositoryImpl) SaveTheater(ctx context.Context, create *entity.SaveTheater) (*entity.SaveTheaterResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	theater := entity.NewTheater(create)

	result := r.database.DB.
		WithContext(ctx).
		Create(&theater)

	err := result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &entity.SaveTheaterResult{
		TheaterID: theater.TheaterID,
	}, nil
}

func (r *theaterRepositoryImpl) UpdateTheater(ctx context.Context, find *entity.FindOneTheater, update *entity.UpdateTheater) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, update)
	if err != nil {
		return err
	}

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Model(&entity.Theater{}).
		Where(findMap).
		Updates(updateMap)

	err = result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	return nil
}

func (r *theaterRepositoryImpl) SoftDeleteTheater(ctx context.Context, find *entity.FindOneTheater) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&entity.Theater{})

	err = result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	return nil
}

func (r *theaterRepositoryImpl) FindOneTheater(ctx context.Context, find *entity.FindOneTheater) (*entity.Theater, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	theater := &entity.Theater{}
	result := r.database.DB.WithContext(ctx).Where(findMap).Select(selectTheaterQuery).First(&theater)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return theater, err
}

func (r *theaterRepositoryImpl) FindManyTheaters(ctx context.Context, find *entity.FindManyTheaters) (*entity.FindManyTheatersResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	theaters := []*entity.Theater{}
	err := r.database.DB.
		WithContext(ctx).
		Scopes(
			theaterKeywordFilter(find.Keyword.String),
			theaterGeoFilter(find.Location),
		).
		Select(selectTheaterQuery).
		Offset(buildOffset(find.Page, find.PageSize)).
		Limit(int(find.PageSize)).
		Find(&theaters).Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	// Count total results
	var totalResults int64

	err = r.database.DB.
		WithContext(ctx).
		Scopes(
			theaterKeywordFilter(find.Keyword.String),
			theaterGeoFilter(find.Location),
		).
		Model(&entity.Theater{}).
		Count(&totalResults).Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &entity.FindManyTheatersResult{
		Theaters: theaters,
		Metadata: &entity.FindManyTheatersMetadata{
			TotalResults: totalResults,
		},
	}, err
}

// theaterKeywordFilter applies a search filter if the keyword is provided.
// It searches (for example) in the name, address, and email fields.
func theaterKeywordFilter(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" {
			return db
		}
		likePattern := "%" + keyword + "%"
		return db.Where("name LIKE ? OR address LIKE ? OR email LIKE ? OR website LIKE ?", likePattern, likePattern, likePattern, likePattern)
	}
}

// theaterGeoFilter applies a geospatial filter if valid latitude, longitude and radius are provided.
func theaterGeoFilter(loc *entity.FindManyTheatersLocation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if loc == nil || loc.Latitude <= 0 || loc.Longitude <= 0 || loc.Radius <= 0 {
			return db
		}
		// Note: MySQL POINT() takes (longitude, latitude)
		return db.Where("ST_Distance_Sphere(location, POINT(?, ?)) <= ?", loc.Longitude, loc.Latitude, loc.Radius)
	}
}

func buildOffset(page, pageSize uint32) int {
	if page < 1 {
		page = 1
	}
	return int((page - 1) * pageSize)
}
