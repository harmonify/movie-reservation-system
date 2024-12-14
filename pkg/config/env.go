package config

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Env         string `mapstructure:"ENV"`
	BaseUrl     string `mapstructure:"BASE_URL"`
	BasePath    string `mapstructure:"BASE_PATH"`
	FrontEndUrl string `mapstructure:"FRONTEND_URL"`
	EnableCors  bool   `mapstructure:"ENABLE_CORS"`
	AppSecret   string `mapstructure:"APP_SECRET"`
	AppPort     string `mapstructure:"SERVER_PORT"`
	GrpcPort    string `mapstructure:"GRPC_PORT"`

	HttpReadTimeOut  string `mapstructure:"HTTP_READ_TIMEOUT"`
	HttpWriteTimeOut string `mapstructure:"HTTP_WRITE_TIMEOUT"`

	DbHostMaster          string `mapstructure:"PG_HOST_MASTER"`
	DbPortMaster          int    `mapstructure:"PG_PORT_MASTER"`
	DbUserMaster          string `mapstructure:"PG_USER_MASTER"`
	DbPasswordMaster      string `mapstructure:"PG_PASSWORD_MASTER"`
	DbHostReplica         string `mapstructure:"PG_HOST_REPLICA"`
	DbPortReplica         int    `mapstructure:"PG_PORT_REPLICA"`
	DbUserReplica         string `mapstructure:"PG_USER_REPLICA"`
	DbPasswordReplica     string `mapstructure:"PG_PASSWORD_REPLICA"`
	DbName                string `mapstructure:"PG_DB"`
	DbMigration           bool   `mapstructure:"PG_AUTO_MIGRATION"`
	DbMaxIdleConn         int    `mapstructure:"PG_MAX_IDLE_CONN"`
	DbMaxOpenConn         int    `mapstructure:"PG_MAX_OPEN_CONN"`
	DbMaxLifetimeInMinute int    `mapstructure:"PG_MAX_LIFETIME_IN_MINUTE"`
	DbUseReplica          bool   `mapstructure:"PG_USE_REPLICA"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	LogType  string `mapstructure:"LOG_TYPE"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	LokiUrl  string `mapstructure:"LOKI_URL"`

	OtelHost     string `mapstructure:"OTEL_ENDPOINT"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE"`

	GrpcUserServiceUrl string `mapstructure:"GRPC_USER_SERVICE_URL"`
}
