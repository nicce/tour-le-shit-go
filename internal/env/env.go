package env

import "os"

type AppEnv struct {
	ScoreboardMode string
	Port           string
}

func getEnvVariable(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic("missing env variable: " + key)
	}

	return v
}

func GetAppEnv() AppEnv {
	return AppEnv{
		ScoreboardMode: getEnvVariable("SCOREBOARD_MODE"),
		Port:           getEnvVariable("PORT"),
	}
}
