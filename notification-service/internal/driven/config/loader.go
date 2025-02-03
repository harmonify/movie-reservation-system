package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewNotificationServiceConfig(configFile string) (*NotificationServiceConfig, error) {
	cfg := &NotificationServiceConfig{}

	// Bind environment variables to Viper
	viper.AutomaticEnv()

	if configFile != "" {
		fmt.Printf(">> Config: Using %s as the configuration file\n", configFile)
		viper.SetConfigFile(string(configFile))
		viper.SetConfigType("env") // force the config file type to be env
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(">> Config: Cannot read config. Error: %v\n", err)
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf(">> Config: Cannot parse config. Error: %v\n", err)
		return cfg, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		fmt.Printf(">> Config: Validation error. Error: %v\n", err)
		return cfg, err
	}

	return cfg, nil
}
