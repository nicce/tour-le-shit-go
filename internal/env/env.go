package env

import "os"

type AppEnv struct {
	ScoreboardMode string
}

func getEnvVariable(key, d string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return d
	}

	return v
}

func GetAppEnv() AppEnv {
	return AppEnv{
		ScoreboardMode: getEnvVariable("SCOREBOARD_MODE", "MOCK"),
	}
}
