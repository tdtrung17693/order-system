package config

func loadAppConfig(config *Config) {
	config.AppPort = getEnv("APP_PORT")
	config.AppEnv = AppEnv(getEnvWithDefault("APP_ENV", AppEnvDevelopment))
	config.AppURL = getEnvWithDefault("APP_URL", "*")
}
