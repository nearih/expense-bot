package config

import (
	"github.com/spf13/viper"
)

type Line struct {
	Channeltoken  string
	Channelsecret string
}

type RootConfig struct {
	Line          Line
	SpreadsheetID string
	SheetRange    string
	Port          int
}

var Config RootConfig

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
