package util

import (
	"github.com/spf13/viper"
)

type PcConfig struct {
	Url      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type MapConfig struct {
	To   string `mapstructure:"to"`
	From string `mapstructure:"from"`
}

type SkipConfig struct {
	Users  []string `mapstructure:"users"`
	Groups []string `mapstructure:"groups"`
}

type MappingConfig struct {
	Roles  []MapConfig `MapConfig,mapstructure:"roles"`
	Groups []MapConfig `MapConfig,mapstructure:"groups"`
}

type Config struct {
	From    PcConfig      `mapstructure:"from"`
	To      PcConfig      `mapstructure:"to"`
	Mapping MappingConfig `mapstructure:"mapping"`
	Skip    SkipConfig    `mapstructure:"skip"`
}

var vp *viper.Viper

func LoadConfig() (Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./util")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
