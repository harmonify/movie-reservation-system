package config

type Config struct {
	Env string `mapstructure:"ENV"`

	AppName                   string `mapstructure:"APP_NAME"`
	AppDefaultCountryDialCode string `mapstructure:"APP_DEFAULT_COUNTRY_DIAL_CODE"`
	AppSecret                 string `mapstructure:"APP_SECRET"`
	AppJwtAudiences           string `mapstructure:"APP_JWT_AUDIENCES"`

	ServiceIdentifier       string `mapstructure:"SERVICE_IDENTIFIER"` // any identifier, used in RBAC
	ServicePort             string `mapstructure:"SERVICE_PORT"`
	ServiceBaseUrl          string `mapstructure:"SERVICE_BASE_URL"`
	ServiceBasePath         string `mapstructure:"SERVICE_BASE_PATH"`
	ServiceHttpReadTimeOut  string `mapstructure:"SERVICE_HTTP_READ_TIMEOUT"`
	ServiceHttpWriteTimeOut string `mapstructure:"SERVICE_HTTP_WRITE_TIMEOUT"`
	ServiceEnableCors       bool   `mapstructure:"SERVICE_ENABLE_CORS"`

	FrontEndUrl string `mapstructure:"FRONTEND_URL"`

	DbHost                string `mapstructure:"PG_HOST"`
	DbPort                int    `mapstructure:"PG_PORT"`
	DbUser                string `mapstructure:"PG_USER"`
	DbPassword            string `mapstructure:"PG_PASSWORD"`
	DbName                string `mapstructure:"PG_DATABASE"`
	DbMigration           bool   `mapstructure:"PG_AUTO_MIGRATION"`
	DbMaxIdleConn         int    `mapstructure:"PG_MAX_IDLE_CONN"`
	DbMaxOpenConn         int    `mapstructure:"PG_MAX_OPEN_CONN"`
	DbMaxLifetimeInMinute int    `mapstructure:"PG_MAX_LIFETIME_IN_MINUTE"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	GrpcPort                  string `mapstructure:"GRPC_PORT"`
	GrpcReservationServiceUrl string `mapstructure:"GRPC_RESERVATION_SERVICE_URL"`
	GrpcMovieServiceUrl       string `mapstructure:"GRPC_MOVIE_SERVICE_URL"`
	GrpcTheaterServiceUrl     string `mapstructure:"GRPC_THEATER_SERVICE_URL"`
	GrpcTicketServiceUrl      string `mapstructure:"GRPC_TICKET_SERVICE_URL"`
	GrpcUserServiceUrl        string `mapstructure:"GRPC_USER_SERVICE_URL"`

	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS"`
	KafkaVersion       string `mapstructure:"KAFKA_VERSION"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`

	LogType  string `mapstructure:"LOG_TYPE"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	LokiUrl  string `mapstructure:"LOKI_URL"`

	OtelHost     string `mapstructure:"OTEL_ENDPOINT"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE"`

	MailgunDefaultSender string `mapstructure:"MAILGUN_DEFAULT_SENDER"`
	MailgunDomain        string `mapstructure:"MAILGUN_DOMAIN"`
	MailgunApiKey        string `mapstructure:"MAILGUN_API_KEY"`

	TwilioAccountSid string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken  string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioServiceSID string `mapstructure:"TWILIO_SERVICE_SID"`
}
