package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() *Config {
	cfg := &Config{}
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..", ".env")

	viper.SetConfigFile(dir)
	readConfig(cfg)

	return cfg
}

func readConfig(config *Config) {
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Cannot read configuration", err)
	}

	setDefaultEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Cannot read configuration: ", err, ". will use default env")
	}
}

func setDefaultEnv() map[string]string {
	m := make(map[string]string)
	for _, s := range os.Environ() {
		a := strings.Split(s, "=")
		if viper.Get("TEMPLATE") == "true" {
			viper.Set(a[0], a[1])
		}
		m[a[0]] = a[1]
	}
	return m
}
