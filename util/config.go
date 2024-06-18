package util

import "github.com/spf13/viper"

type Config struct {
	ClientId     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	ModuleSecret string `mapstructure:"MODULE_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("connection")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
