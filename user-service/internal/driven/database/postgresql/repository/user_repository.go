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
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"gorm.io/gorm/clause"
)

type UserRepositoryParam struct {
	fx.In

	Database                  *database.Database
	PostgresqlErrorTranslator database.PostgresqlErrorTranslator
	Tracer                    tracer.Tracer
	Logger                    logger.Logger
	Util                      *util.Util
}

type UserRepositoryResult struct {
	fx.Out

	UserStorage shared_service.UserStorage
}

type userRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserRepository(p UserRepositoryParam) UserRepositoryResult {
	return UserRepositoryResult{
		UserStorage: &userRepositoryImpl{
			database: p.Database,
			pgErrTl:  p.PostgresqlErrorTranslator,
			tracer:   p.Tracer,
			logger:   p.Logger,
			util:     p.Util,
		},
	}
}

func (r *userRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserStorage {
	if tx == nil {
		return r
	}

	return &userRepositoryImpl{
		database: &database.Database{
			DB:     tx.DB,
			Logger: r.logger,
		},
		tracer: r.tracer,
		logger: r.logger,
		util:   r.util,
	}
}

func (r *userRepositoryImpl) SaveUser(ctx context.Context, createModel entity.SaveUser) (*entity.User, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userModel := (&model.User{}).FromSaveEntity(createModel)

	err := r.database.DB.
		WithContext(ctx).
		Create(userModel).
		Error

	if err != nil {
		terr := r.pgErrTl.Translate(err)
		switch e := (terr).(type) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == database.UniqueViolation {
				if strings.Contains(err.Error(), "username") {
					return nil, auth_service.ErrDuplicateUsername
				} else if strings.Contains(err.Error(), "email") {
					return nil, auth_service.ErrDuplicateEmail
				} else if strings.Contains(err.Error(), "phone_number") {
					return nil, auth_service.ErrDuplicatePhoneNumber
				}
			}
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

	var userModel model.User
	err = r.database.DB.WithContext(ctx).Where(findMap).First(&userModel).Error
	if err != nil {
		return nil, err
	}

	return userModel.ToEntity(), err
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, userUUID string, updateModel entity.UpdateUser) (*entity.User, error) {
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

	var updatedUserModel model.User = model.User{UUID: parsedUUID}
	err = r.database.DB.
		WithContext(ctx).
		Model(&updatedUserModel).
		Clauses(clause.Returning{}).
		Updates(updateMap).
		Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
	}

	return updatedUserModel.ToEntity(), err
}

func (r *userRepositoryImpl) SoftDeleteUser(ctx context.Context, userUUID string) error {
	return r.database.DB.
		WithContext(ctx).
		Delete(&model.User{}, userUUID).
		Error
}
