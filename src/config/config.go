package config

type DBEnv struct {
	DBName     string `envconfig:"NAME"`
	DBHost     string `envconfig:"HOST"`
	DBPort     int    `envconfig:"PORT"`
	DBUsername string `envconfig:"USERNAME"`
	DBPassword string `envconfig:"PASSWORD"`
}
