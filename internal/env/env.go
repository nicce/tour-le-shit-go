package env

import "os"

type AppEnv struct {
	MembersMode    string
	Port           string
	ScoreboardMode string
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
		MembersMode:    getEnvVariable("MEMBERS_MODE"),
		Port:           getEnvVariable("PORT"),
		ScoreboardMode: getEnvVariable("SCOREBOARD_MODE"),
	}
}
