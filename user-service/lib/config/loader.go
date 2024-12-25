package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigFile struct {
	Path string
}

func NewConfig(configFile *ConfigFile) (*Config, error) {
	cfg := &Config{}

	// Bind environment variables to Viper
	viper.AutomaticEnv()

	if configFile != nil && configFile.Path != "" {
		fmt.Println(fmt.Sprintf(">> Config: Using %s as the configuration file", configFile.Path))
		viper.SetConfigFile(configFile.Path)
		viper.SetConfigType("env") // force the config file type to be env
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf(">> Config: Cannot read config. Error: %s", err))
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println(fmt.Sprintf(">> Config: Cannot parse config. Error: %s", err))
		return cfg, err
	}

	return cfg, nil
}
