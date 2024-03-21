package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	GrpcServerAddress       string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

// read config from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // path to look for the config file in
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	// viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.AutomaticEnv()
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		return
	}

	err = viper.Unmarshal(&config)

	return config, err
}
