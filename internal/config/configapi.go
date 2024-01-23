package config

import (
	"os"
)

type ConfigApi struct {
	AGEAPI    string
	NATIONAPI string
	GENDERAPI string
}

func NewCfgApi() *ConfigApi {

	return &ConfigApi{
		AGEAPI:    os.Getenv("AGEAPI"),
		NATIONAPI: os.Getenv("NATIONAPI"),
		GENDERAPI: os.Getenv("GENDERAPI"),
	}
}
