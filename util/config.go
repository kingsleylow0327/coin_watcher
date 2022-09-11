package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_NAME          string `mapstructure:"DB_NAME"`
	DB_PASSWORD      string `mapstructure:"DB_PASSWORD"`
	DB_ADDRESS       string `mapstructure:"DB_ADDRESS"`
	DB_SERVER        string `mapstructure:"DB_SERVER"`
	DB_PATH          string
	DISCORD_TOKEN    string `mapstructure:"DISCORD_TOKEN"`
	SENDER_CHANNEL   string `mapstructure:"SENDER_CHANNEL"`
	RECEIVER_CHANNEL string `mapstructure:"RECEIVER_CHANNEL"`
	CORNIX_ID        string `mapstructure:"CORNIX_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("Viper Error: %s", err))
		return
	}

	err = viper.Unmarshal(&config)
	config.DB_PATH = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.DB_NAME,
		config.DB_PASSWORD,
		config.DB_ADDRESS,
		config.DB_SERVER)
	return
}
