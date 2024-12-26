package rbac

import (
	"context"
	"path"
	"runtime"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"go.uber.org/fx"
)

type (
	CasbinParam struct {
		fx.In

		Database *database.Database
		Tracer   tracer.Tracer
		Logger   logger.Logger
		Config   *config.Config
	}

	CasbinResult struct {
		fx.Out

		RbacStorage shared_service.RbacStorage
	}

	casbinImpl struct {
		enforcer *casbin.Enforcer
		tracer   tracer.Tracer
		logger   logger.Logger
		config   *config.Config
	}
)

func NewCasbin(p CasbinParam) CasbinResult {
	// Initialize a Gorm adapter and use it in a Casbin enforcer
	a, _ := gormadapter.NewAdapterByDB(p.Database.DB)
	_, filename, _, _ := runtime.Caller(0)
	e, _ := casbin.NewEnforcer(
		path.Join(path.Dir(filename), "rbac_model.conf"),
		a,
	)

	// Load the policy from DB.
	e.LoadPolicy()

	return CasbinResult{
		RbacStorage: &casbinImpl{
			enforcer: e,
			tracer:   p.Tracer,
			logger:   p.Logger,
			config:   p.Config,
		},
	}
}

func (c *casbinImpl) WithTx(tx *database.Transaction) shared_service.RbacStorage {
	if tx == nil {
		return c
	}

	a, _ := gormadapter.NewAdapterByDB(tx.DB)
	c.enforcer.SetAdapter(a)

	return c
}

func (c *casbinImpl) CheckPermission(ctx context.Context, p shared_service.CheckPermissionParam) (bool, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	authorized, err := c.enforcer.Enforce(p.UUID, p.Resource, p.Action)
	return authorized, err
}

func (c *casbinImpl) GrantPermission(ctx context.Context, p shared_service.GrantPermissionParam) (bool, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	success, err := c.enforcer.AddPolicy(p.UUID, p.Resource, p.Action)
	return success, err
}

func (c *casbinImpl) BulkGrantPermission(ctx context.Context, p []shared_service.GrantPermissionParam) (bool, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var rules [][]string

	for _, rule := range p {
		rules = append(rules, []string{rule.UUID, rule.Resource, string(rule.Action)})
	}

	success, err := c.enforcer.AddPolicies(rules)
	return success, err
}

func (c *casbinImpl) RevokePermission(ctx context.Context, p shared_service.RevokePermissionParam) (bool, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	success, err := c.enforcer.RemovePolicy(p.UUID, p.Resource, p.Action)
	return success, err
}

func (c *casbinImpl) GrantRole(ctx context.Context, p shared_service.GrantRoleParam) (bool, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	success, err := c.enforcer.AddRoleForUserInDomain(p.UUID, string(p.Role), c.config.ServiceIdentifier)
	return success, err
}
