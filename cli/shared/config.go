package shared

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env string `mapstructure:"ENV"`

	KafkaServerUrl string `mapstructure:"KAFKA_SERVER_URL"`
	SqlitePath     string `mapstructure:"SQLITE_PATH"`
}

type ConfigFile struct {
	Path string
}

func NewConfig(configFile *ConfigFile) (*Config, error) {
	cfg := &Config{}

	// Bind environment variables to Viper
	viper.AutomaticEnv()

	if configFile != nil && configFile.Path != "" {
		fmt.Printf(">> Config: Using %s as the configuration file\n", configFile.Path)
		viper.SetConfigFile(configFile.Path)
		viper.SetConfigType("env") // force the config file type to be env
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf(">> Config: Cannot read config. Error: %v\n", err)
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf(">> Config: Cannot parse config. Error: %v\n", err)
		return cfg, err
	}

	return cfg, nil
}
