package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Channeltoken  string
	Channelsecret string
	SpreadsheetID string
}

var Config config

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
