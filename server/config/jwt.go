package config

func loadJWTConfig(config *Config) {
	jwtSecretKey := getEnv("JWT_SECRET_KEY")

	config.JwtSecretKey = jwtSecretKey
}
