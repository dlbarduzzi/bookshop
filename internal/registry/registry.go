package registry

import "github.com/spf13/viper"

var (
	_configPath = "."
	_configName = ".env"
	_configType = "env"
)

func NewRegistry() (*viper.Viper, error) {
	reg := viper.New()
	reg.AutomaticEnv()

	reg.AddConfigPath(_configPath)
	reg.SetConfigName(_configName)
	reg.SetConfigType(_configType)

	if err := reg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return reg, nil
}

func SetConfigPath(p string) {
	_configPath = p
}

func SetConfigName(n string) {
	_configName = n
}

func SetConfigType(t string) {
	_configType = t
}
