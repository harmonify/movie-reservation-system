package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database/postgresql"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"gorm.io/gorm/clause"
)

type UserKeyRepositoryParam struct {
	fx.In

	Database *database.Database
	Tracer   tracer.Tracer
	Logger   logger.Logger
	Util     *util.Util
}

type UserKeyRepositoryResult struct {
	fx.Out

	UserKeyStorage shared_service.UserKeyStorage
}

type userKeyRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserKeyRepository(p UserKeyRepositoryParam) UserKeyRepositoryResult {
	return UserKeyRepositoryResult{
		UserKeyStorage: &userKeyRepositoryImpl{
			database: p.Database,
			tracer:   p.Tracer,
			logger:   p.Logger,
			util:     p.Util,
		},
	}
}

func (r *userKeyRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserKeyStorage {
	if tx == nil {
		return r
	}

	return &userKeyRepositoryImpl{
		database: &database.Database{
			DB:     tx.DB,
			Logger: r.logger,
		},
		tracer: r.tracer,
		logger: r.logger,
		util:   r.util,
	}
}

func (r *userKeyRepositoryImpl) SaveUserKey(ctx context.Context, createModel entity.SaveUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userKeyModel := (&model.UserKey{}).FromSaveEntity(createModel)

	err := r.database.DB.WithContext(ctx).Create(&userKeyModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresql.UniqueViolation {
				if strings.Contains(err.Error(), "email") {
					return nil, auth_service.ErrDuplicateEmail
				} else if strings.Contains(err.Error(), "phone_number") {
					return nil, auth_service.ErrDuplicatePhoneNumber
				}
			}
		}
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) FindUserKey(ctx context.Context, findModel entity.FindUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	var userKeyModel model.UserKey
	err = r.database.DB.WithContext(ctx).Where(findMap).First(&userKeyModel).Error
	if err != nil {
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) UpdateUserKey(ctx context.Context, userUUID string, updateModel entity.UpdateUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(updateModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	parsedUUID, err := uuid.Parse(userUUID)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	var updatedUserKeyModel model.UserKey = model.UserKey{UserUUID: parsedUUID}
	err = r.database.DB.
		WithContext(ctx).
		Model(&updatedUserKeyModel).
		Clauses(clause.Returning{}).
		Updates(updateMap).
		Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
	}

	return updatedUserKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) SoftDeleteUserKey(ctx context.Context, userUUID string) error {
	return r.database.DB.
		WithContext(ctx).
		Delete(&model.UserKey{}, userUUID).
		Error
}
