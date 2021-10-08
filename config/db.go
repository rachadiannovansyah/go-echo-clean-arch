package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// DBConfig ..
type DBConfig struct {
	DSN string
}

func LoadConfig() DBConfig {
	return DBConfig{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			viper.GetString(`database.host`),
			viper.GetString(`database.port`),
			viper.GetString(`database.user`),
			viper.GetString(`database.pass`),
			viper.GetString(`database.name`),
		),
	}
}
