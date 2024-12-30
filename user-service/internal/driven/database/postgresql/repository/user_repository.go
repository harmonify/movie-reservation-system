package repository

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserRepository(
	database *database.Database,
	pgErrTl database.PostgresqlErrorTranslator,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared_service.UserStorage {
	return &userRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *userRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserStorage {
	if tx == nil {
		return r
	}
	return NewUserRepository(
		r.database.WithTx(tx),
		r.pgErrTl,
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *userRepositoryImpl) SaveUser(ctx context.Context, createModel entity.SaveUser) (*entity.User, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userModel := (&model.User{}).FromSaveEntity(createModel)

	result := r.database.DB.
		WithContext(ctx).
		Create(userModel)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		switch e := (err).(type) {
		case *database.DuplicatedKeyError:
			if e.ColumnName == "username" {
				return nil, auth_service.ErrDuplicateUsername
			} else if e.ColumnName == "email" {
				return nil, auth_service.ErrDuplicateEmail
			} else if e.ColumnName == "phone_number" {
				return nil, auth_service.ErrDuplicatePhoneNumber
			}
		default:
			return nil, error_constant.ErrInternalServerError
		}
	}

	return userModel.ToEntity(), err
}

func (r *userRepositoryImpl) FindUser(ctx context.Context, findModel entity.FindUser) (*entity.User, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	userModel := model.User{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&userModel)
	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	return userModel.ToEntity(), err
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, findModel entity.FindUser, updateModel entity.UpdateUser) (*entity.User, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(updateModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	userModel := model.User{}
	result := r.database.DB.
		WithContext(ctx).
		Model(&userModel).
		Where(findMap).
		Clauses(clause.Returning{}).
		Updates(updateMap)

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *userRepositoryImpl) SoftDeleteUser(ctx context.Context, findModel entity.FindUser) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&model.User{})

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	return nil
}
