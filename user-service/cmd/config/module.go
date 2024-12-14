package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	Port          string `cfg:"port" cfgDefault:"9101"`
	PostgreSQLDSN string `cfg:"postgresql_dsn"`
}

var ConfigModule = fx.Module("config", fx.Provide(LoadConfig))
var Cfg = &config{}

func LoadConfig() *config {
	env := os.Getenv("ENV")

	if env == "test" {
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "..")
		err := os.Chdir(dir)
		if err != nil {
			panic(err)
		}

		viper.SetConfigFile("../.env")
		readConfig(Cfg)
		fmt.Println("🔍 Run on Testing Mode")
	} else {
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "..")
		err := os.Chdir(dir)
		if err != nil {
			panic(err)
		}

		viper.SetConfigFile("../.env")
		readConfig(Cfg)
		fmt.Println("🔍 Run on Testing Mode")
	}

	return Cfg
}

func readConfig(config *config) {
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("Cannot read configuration", err)
	}

	setDefaultEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Error("Cannot read configuration: ", err, ". will use default env")
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
