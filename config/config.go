package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config ..
type Config struct {
	DB DBConfig
}

// NewConfig ..
func NewConfig() *Config {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	return &Config{
		DB: LoadConfig(),
	}
}
