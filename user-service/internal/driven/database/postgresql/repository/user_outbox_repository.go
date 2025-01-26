package repository

import (
	"context"

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

type outboxRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewOutboxRepository(
	database *database.Database,
	pgErrTl database.PostgresqlErrorTranslator,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.OutboxStorage {
	return &outboxRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *outboxRepositoryImpl) WithTx(tx *database.Transaction) shared.OutboxStorage {
	if tx == nil {
		return r
	}
	return NewOutboxRepository(
		r.database.WithTx(tx),
		r.pgErrTl,
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *outboxRepositoryImpl) SaveOutbox(ctx context.Context, createModel entity.SaveUserOutbox) (*entity.UserOutbox, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	outboxModel := model.NewUserOutbox(createModel)

	result := r.database.DB.
		WithContext(ctx).
		Create(outboxModel)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	return outboxModel.ToEntity(), err
}
