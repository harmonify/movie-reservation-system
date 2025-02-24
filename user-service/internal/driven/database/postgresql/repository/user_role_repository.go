package repository

import (
	"context"
	"errors"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"go.uber.org/zap"
)

type userRoleRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserRoleRepository(
	database *database.Database,
	pgErrTl database.PostgresqlErrorTranslator,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.UserRoleStorage {
	return &userRoleRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *userRoleRepositoryImpl) WithTx(tx *database.Transaction) shared.UserRoleStorage {
	if tx == nil {
		return r
	}
	return NewUserRoleRepository(
		r.database.WithTx(tx),
		r.pgErrTl,
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *userRoleRepositoryImpl) SearchUserRoles(ctx context.Context, searchModel entity.SearchUserRoles) ([]string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if searchModel.UserUUID == "" {
		err := errors.New("user uuid is required")
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	var userRoles []string
	err := r.database.DB.
		Raw(
			`SELECT roles.name
			FROM user_roles
			JOIN roles ON user_roles.role_id = roles.id
			WHERE user_roles.user_uuid = ?
			AND user_roles.deleted_at IS NULL
			AND roles.deleted_at IS NULL`,
			searchModel.UserUUID,
		).
		Scan(&userRoles).
		Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	return userRoles, nil
}

func (r *userRoleRepositoryImpl) SaveUserRoles(ctx context.Context, createModel entity.SaveUserRoles) ([]*entity.UserRole, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userRoleModels := make([]*model.UserRole, 0, len(createModel.RoleID))
	for _, roleID := range createModel.RoleID {
		userRoleModels = append(userRoleModels, model.NewUserRole(createModel.UserUUID, roleID))
	}

	result := r.database.DB.
		WithContext(ctx).
		Create(userRoleModels)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	userRoleEntities := make([]*entity.UserRole, 0, len(userRoleModels))
	for _, userRole := range userRoleModels {
		userRoleEntities = append(userRoleEntities, userRole.ToEntity())
	}

	return userRoleEntities, err
}

func (r *userRoleRepositoryImpl) SoftDeleteUserRoles(ctx context.Context, searchModel entity.SearchUserRoles) error {
	searchMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, searchModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(searchMap).
		Delete(&model.UserRole{})

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	return nil
}
