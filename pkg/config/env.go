package config

type Config struct {
	Env string `mapstructure:"ENV" validate:"required,oneof=development staging production"`

	AppName                   string `mapstructure:"APP_NAME" validate:"required"`
	AppDefaultCountryDialCode string `mapstructure:"APP_DEFAULT_COUNTRY_DIAL_CODE"`
	AppSecret                 string `mapstructure:"APP_SECRET" validate:"required"`
	AppJwtAudiences           string `mapstructure:"APP_JWT_AUDIENCES" validate:"required"`

	ServiceIdentifier       string `mapstructure:"SERVICE_IDENTIFIER" validate:"required"`
	ServiceHttpPort         string `mapstructure:"SERVICE_HTTP_PORT" validate:"required,numeric"`
	ServiceHttpBaseUrl      string `mapstructure:"SERVICE_HTTP_BASE_URL" validate:"required"`
	ServiceHttpBasePath     string `mapstructure:"SERVICE_HTTP_BASE_PATH" validate:"required"`
	ServiceHttpReadTimeOut  string `mapstructure:"SERVICE_HTTP_READ_TIMEOUT" validate:"required"`
	ServiceHttpWriteTimeOut string `mapstructure:"SERVICE_HTTP_WRITE_TIMEOUT" validate:"required"`
	ServiceHttpEnableCors   bool   `mapstructure:"SERVICE_HTTP_ENABLE_CORS" validate:"required,boolean"`

	FrontEndUrl string `mapstructure:"FRONTEND_URL"`

	DbHost                string `mapstructure:"PG_HOST"`
	DbPort                int    `mapstructure:"PG_PORT" validate:"min=1,max=65535"`
	DbUser                string `mapstructure:"PG_USER"`
	DbPassword            string `mapstructure:"PG_PASSWORD"`
	DbName                string `mapstructure:"PG_DATABASE"`
	DbMigration           bool   `mapstructure:"PG_AUTO_MIGRATION" validate:"boolean"`
	DbMaxIdleConn         int    `mapstructure:"PG_MAX_IDLE_CONN"`
	DbMaxOpenConn         int    `mapstructure:"PG_MAX_OPEN_CONN"`
	DbMaxLifetimeInMinute int    `mapstructure:"PG_MAX_LIFETIME_IN_MINUTE"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	GrpcPort                       string `mapstructure:"GRPC_PORT" validate:"numeric"`
	NotificationServiceGrpcAddress string `mapstructure:"NOTIFICATION_SERVICE_GRPC_ADDRESS" validate:"required,url"`

	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS" validate:"required"`
	KafkaVersion       string `mapstructure:"KAFKA_VERSION" validate:"required"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP" validate:"required"`

	LogType  string `mapstructure:"LOG_TYPE" validate:"required"`
	LogLevel string `mapstructure:"LOG_LEVEL" validate:"required"`
	LokiUrl  string `mapstructure:"LOKI_URL" validate:"required_if=LogType loki"`

	OtelEndpoint string `mapstructure:"OTEL_ENDPOINT" validate:"required"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE" validate:"required,boolean"`

	MailgunDefaultSender string `mapstructure:"MAILGUN_DEFAULT_SENDER"`
	MailgunDomain        string `mapstructure:"MAILGUN_DOMAIN"`
	MailgunApiKey        string `mapstructure:"MAILGUN_API_KEY"`

	TwilioAccountSid string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken  string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioServiceSID string `mapstructure:"TWILIO_SERVICE_SID"`
}
