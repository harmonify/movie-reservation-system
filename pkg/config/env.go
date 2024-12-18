package config

type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	AppSecret   string `mapstructure:"APP_SECRET"`
	AppPort     string `mapstructure:"APP_PORT"`
	Env         string `mapstructure:"ENV"`
	BaseUrl     string `mapstructure:"BASE_URL"`
	BasePath    string `mapstructure:"BASE_PATH"`
	FrontEndUrl string `mapstructure:"FRONTEND_URL"`
	EnableCors  bool   `mapstructure:"ENABLE_CORS"`

	HttpReadTimeOut  string `mapstructure:"HTTP_READ_TIMEOUT"`
	HttpWriteTimeOut string `mapstructure:"HTTP_WRITE_TIMEOUT"`

	DbHostMaster          string `mapstructure:"PG_HOST_MASTER"`
	DbPortMaster          int    `mapstructure:"PG_PORT_MASTER"`
	DbUserMaster          string `mapstructure:"PG_USER_MASTER"`
	DbPasswordMaster      string `mapstructure:"PG_PASSWORD_MASTER"`
	DbName                string `mapstructure:"PG_DB"`
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

	KafkaProducerPort string `mapstructure:"KAFKA_PRODUCER_PORT"`

	LogType  string `mapstructure:"LOG_TYPE"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	LokiUrl  string `mapstructure:"LOKI_URL"`

	OtelHost     string `mapstructure:"OTEL_ENDPOINT"`
	OtelInsecure bool   `mapstructure:"OTEL_INSECURE"`
}
