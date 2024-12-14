package rest

import (
	"reflect"

	"github.com/gin-gonic/gin"
	globalInterface "github.com/harmonify/movie-reservation-system/pkg/interfaces"
	"github.com/harmonify/movie-reservation-system/pkg/middleware"
	commonRest "github.com/harmonify/movie-reservation-system/user-service/api/rest/common"
	"github.com/harmonify/movie-reservation-system/user-service/api/rest/interfaces"
	userRest "github.com/harmonify/movie-reservation-system/user-service/api/rest/user"
	"go.uber.org/fx"
)

var (
	provider = fx.Provide(
		commonRest.NewHealthCheckRestHandler,
		userRest.NewPayoutRestHandler,
		userRest.NewUserRestHandler,
		NewRest,
	)

	Module = fx.Options(
		middleware.Provider,
		provider,
	)
)

type Rest struct {
	HealthCheckRestHandler interfaces.IHealthCheckRestHandler
	PayoutRestHandler      interfaces.PayoutRestHandler
	UserRestHandler        interfaces.UserRestHandler
}

func NewRest(
	healthCheckRestHandler interfaces.IHealthCheckRestHandler,
	payoutRestHandler interfaces.PayoutRestHandler,
	userRestHandler interfaces.UserRestHandler,
) *Rest {
	return &Rest{
		HealthCheckRestHandler: healthCheckRestHandler,
		PayoutRestHandler:      payoutRestHandler,
		UserRestHandler:        userRestHandler,
	}
}

func (h *Rest) Initialize(g *gin.RouterGroup) {
	val := reflect.ValueOf(h).Elem()

	for i := 0; i < val.NumField(); i++ {
		if !val.Type().Field(i).IsExported() {
			continue
		}

		field := val.Field(i).Interface()

		if handler, ok := field.(globalInterface.RestHandler); ok {
			handler.Register(g)
		}
	}
}
