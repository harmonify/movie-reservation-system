package config

type MovieServiceConfig struct {
	Env string `mapstructure:"ENV" validate:"required,oneof=dev test prod"`

	AppName   string `mapstructure:"APP_NAME" validate:"required"`
	AppSecret string `mapstructure:"APP_SECRET" validate:"required,base64"`

	AuthJwtIssuerIdentifier    string `mapstructure:"AUTH_JWT_ISSUER_IDENTIFIER" validate:"required"`
	AuthJwtAudienceIdentifiers string `mapstructure:"AUTH_JWT_AUDIENCE_IDENTIFIERS" validate:"required"`

	ServiceIdentifier       string `mapstructure:"SERVICE_IDENTIFIER" validate:"required"`
	ServiceHttpPort         string `mapstructure:"SERVICE_HTTP_PORT" validate:"required,numeric"`
	ServiceHttpBaseUrl      string `mapstructure:"SERVICE_HTTP_BASE_URL" validate:"required"`
	ServiceHttpBasePath     string `mapstructure:"SERVICE_HTTP_BASE_PATH" validate:"required"`
	ServiceHttpReadTimeOut  string `mapstructure:"SERVICE_HTTP_READ_TIMEOUT" validate:"required"`
	ServiceHttpWriteTimeOut string `mapstructure:"SERVICE_HTTP_WRITE_TIMEOUT" validate:"required"`
	ServiceHttpEnableCors   bool   `mapstructure:"SERVICE_HTTP_ENABLE_CORS" validate:"boolean"`

	FrontEndUrl string `mapstructure:"FRONTEND_URL"`

	MongoUri        string `mapstructure:"MONGO_URI" validate:"required,mongodb_connection_string"`
	MongoDbName     string `mapstructure:"MONGO_DB_NAME" validate:"required"`
	MongoReplicaSet string `mapstructure:"MONGO_REPLICA_SET" validate:"required"`

	RedisHost string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort string `mapstructure:"REDIS_PORT" validate:"required,numeric"`
	RedisPass string `mapstructure:"REDIS_PASS" validate:"required"`

	GrpcPort              int    `mapstructure:"GRPC_PORT" validate:"required,numeric"`
	GrpcAuthServiceUrl    string `mapstructure:"GRPC_AUTH_SERVICE_URL" validate:"required,url"`
	GrpcTheaterServiceUrl string `mapstructure:"GRPC_THEATER_SERVICE_URL" validate:"required,url"`

	LogType  string `mapstructure:"LOG_TYPE" validate:"required"`
	LogLevel string `mapstructure:"LOG_LEVEL" validate:"required"`
	LokiUrl  string `mapstructure:"LOKI_URL" validate:"required_if=LogType loki"`

	TracerType   string `mapstructure:"TRACER_TYPE" validate:"required,oneof=jaeger console nop"`
	OtelEndpoint string `mapstructure:"OTEL_ENDPOINT" validate:"required"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE" validate:"required,boolean"`
}
