package config

type NotificationServiceConfig struct {
	Env string `mapstructure:"ENV" validate:"required,oneof=dev test prod"`

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
	ServiceHttpEnableCors   bool   `mapstructure:"SERVICE_HTTP_ENABLE_CORS" validate:"boolean"`

	FrontEndUrl string `mapstructure:"FRONTEND_URL" validate:"required,url"`

	GrpcPort int `mapstructure:"GRPC_PORT" validate:"required,numeric,min=1024,max=65535"`

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

	MailgunDefaultSender string `mapstructure:"MAILGUN_DEFAULT_SENDER" validate:"required"`
	MailgunDomain        string `mapstructure:"MAILGUN_DOMAIN" validate:"required"`
	MailgunApiKey        string `mapstructure:"MAILGUN_API_KEY" validate:"required"`

	TwilioAccountSid string `mapstructure:"TWILIO_ACCOUNT_SID" validate:"required"`
	TwilioAuthToken  string `mapstructure:"TWILIO_AUTH_TOKEN" validate:"required"`
	TwilioServiceSID string `mapstructure:"TWILIO_SERVICE_SID" validate:"required"`
}
