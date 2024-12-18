package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigFile struct {
	Path string
}

func NewConfig(configFile ConfigFile) *Config {
	// Bind environment variables to Viper
	viper.AutomaticEnv()

	if configFile.Path != "" {
		fmt.Println(fmt.Sprintf(">> Config: Using %s as the configuration file", configFile.Path))
		viper.SetConfigFile(configFile.Path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf(">> Config: Cannot read config. Erorr: %s", err))
	}

	cfg := &Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println(fmt.Sprintf(">> Config: Cannot parse config. Erorr: %s", err))
	}

	return cfg
}
