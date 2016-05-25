package utilities

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type IConfigUtil interface {
	GetConfig(key string) string
}

type ConfigUtil struct {
	Viper *viper.Viper
}

func NewConfigUtil() *ConfigUtil {
	// read config
	var env string
	if os.Getenv("GO_ENV") != "" {
		env = os.Getenv("GO_ENV")
	}else {
		env = "development"
	}
	r := ConfigUtil{}
	r.Viper = viper.New()
	r.Viper.SetConfigName(env)
	r.Viper.AddConfigPath("./")
	r.Viper.AddConfigPath("./config/")
	r.Viper.AddConfigPath(".")
	r.Viper.SetConfigType("json")
	err := r.Viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return &r
}

func (r ConfigUtil) GetConfig(key string) string {
	return r.Viper.GetString(key)
}
