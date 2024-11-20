package registry

import "github.com/spf13/viper"

var (
	_configPath = "."
	_configType = "env"
	_configName = ".env"
)

func NewRegistry() (*viper.Viper, error) {
	reg := viper.New()
	reg.AutomaticEnv()

	reg.AddConfigPath(_configPath)
	reg.SetConfigType(_configType)
	reg.SetConfigName(_configName)

	if err := reg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return reg, nil
}

func SetConfigPath(v string) {
	_configPath = v
}

func SetConfigType(v string) {
	_configType = v
}

func SetConfigName(v string) {
	_configName = v
}
