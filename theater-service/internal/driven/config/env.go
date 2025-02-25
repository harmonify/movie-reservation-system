package config

type TheaterServiceConfig struct {
	Env string `mapstructure:"ENV" validate:"required,oneof=dev test prod"`

	AppName                   string `mapstructure:"APP_NAME" validate:"required"`
	AppDefaultCountryDialCode string `mapstructure:"APP_DEFAULT_COUNTRY_DIAL_CODE"`
	AppDefaultSupportEmail    string `mapstructure:"APP_DEFAULT_SUPPORT_EMAIL"`
	AppSecret                 string `mapstructure:"APP_SECRET" validate:"required,base64"`

	AuthJwtIssuerIdentifier    string `mapstructure:"AUTH_JWT_ISSUER_IDENTIFIER" validate:"required"`
	AuthJwtAudienceIdentifiers string `mapstructure:"AUTH_JWT_AUDIENCE_IDENTIFIERS" validate:"required"`

	ServiceIdentifier       string `mapstructure:"SERVICE_IDENTIFIER" validate:"required"`
	ServiceHttpPort         string `mapstructure:"SERVICE_HTTP_PORT" validate:"required,numeric"`
	ServiceHttpBaseUrl      string `mapstructure:"SERVICE_HTTP_BASE_URL" validate:"required"`
	ServiceHttpBasePath     string `mapstructure:"SERVICE_HTTP_BASE_PATH" validate:"required"`
	ServiceHttpReadTimeOut  string `mapstructure:"SERVICE_HTTP_READ_TIMEOUT" validate:"required"`
	ServiceHttpWriteTimeOut string `mapstructure:"SERVICE_HTTP_WRITE_TIMEOUT" validate:"required"`
	ServiceHttpEnableCors   bool   `mapstructure:"SERVICE_HTTP_ENABLE_CORS" validate:"boolean"`

	FrontEndUrl string `mapstructure:"FRONTEND_URL" validate:"required,url"`

	DbType                string `mapstructure:"DB_TYPE" validate:"required,oneof=postgresql mysql"`
	DbHost                string `mapstructure:"DB_HOST"`
	DbPort                int    `mapstructure:"DB_PORT" validate:"min=1,max=65535"`
	DbUser                string `mapstructure:"DB_USER"`
	DbPassword            string `mapstructure:"DB_PASSWORD"`
	DbName                string `mapstructure:"DB_DATABASE"`
	DbMaxIdleConn         int    `mapstructure:"DB_MAX_IDLE_CONN"`
	DbMaxOpenConn         int    `mapstructure:"DB_MAX_OPEN_CONN"`
	DbMaxLifetimeInMinute int    `mapstructure:"DB_MAX_LIFETIME_IN_MINUTE"`

	RedisHost string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort string `mapstructure:"REDIS_PORT" validate:"required,numeric"`
	RedisPass string `mapstructure:"REDIS_PASS" validate:"required"`

	GrpcPort           int    `mapstructure:"GRPC_PORT" validate:"required,numeric,min=1024,max=65535"`
	GrpcAuthServiceUrl string `mapstructure:"GRPC_AUTH_SERVICE_URL" validate:"required,url"`

	KafkaBrokers               string `mapstructure:"KAFKA_BROKERS" validate:"required"`
	KafkaVersion               string `mapstructure:"KAFKA_VERSION" validate:"required"`
	KafkaConsumerGroup         string `mapstructure:"KAFKA_CONSUMER_GROUP" validate:"required"`
	KafkaTopicUserRegisteredV1 string `mapstructure:"KAFKA_TOPIC_USER_REGISTERED_V1" validate:"required"`

	LogType  string `mapstructure:"LOG_TYPE" validate:"required"`
	LogLevel string `mapstructure:"LOG_LEVEL" validate:"required"`
	LokiUrl  string `mapstructure:"LOKI_URL" validate:"required_if=LogType loki"`

	TracerType   string `mapstructure:"TRACER_TYPE" validate:"required,oneof=jaeger console nop"`
	OtelEndpoint string `mapstructure:"OTEL_ENDPOINT" validate:"required"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE" validate:"required,boolean"`
}
