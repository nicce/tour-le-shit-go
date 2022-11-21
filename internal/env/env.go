package env

import "os"

type AppEnv struct {
	MembersMode    string
	Port           string
	ScoreboardMode string
	Db             Db
}

type Db struct {
	Username string
	Password string
	Name     string
}

func getEnvVariable(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic("missing env variable: " + key)
	}

	return v
}

func GetAppEnv() AppEnv {
	db := Db{
		Username: getEnvVariable("DATABASE_USER"),
		Password: getEnvVariable("DATABASE_PASSWORD"),
		Name:     getEnvVariable("DATABASE_NAME"),
	}

	return AppEnv{
		MembersMode:    getEnvVariable("MEMBERS_MODE"),
		Port:           getEnvVariable("PORT"),
		ScoreboardMode: getEnvVariable("SCOREBOARD_MODE"),
		Db:             db,
	}
}
