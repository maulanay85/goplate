package config

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    int    `mapstructure:"APP_PORT"`
	Env     string `mapstructure:"APP_ENV"`
}
