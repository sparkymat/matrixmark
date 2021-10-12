package config

import (
	"strings"

	"github.com/caarlos0/env/v6"
)

type Credentials struct {
	Username string
	Password string
}

type API interface {
	ShioriURL() string
	ShioriUsersMap() map[string]Credentials
}

func New() API {
	envConfig := envConfig{}
	if err := env.Parse(&envConfig); err != nil {
		panic(err)
	}

	shioriUsersMap := map[string]Credentials{}
	for _, userMapString := range strings.Split(envConfig.ShioriUsersMapString, ",") {
		mapWords := strings.Split(userMapString, ":")
		if len(mapWords) == 3 {
			shioriUsersMap[mapWords[0]] = Credentials{
				Username: mapWords[1],
				Password: mapWords[2],
			}
		}
	}

	if len(shioriUsersMap) == 0 {
		panic("shiori users not found or malformed")
	}

	return &service{
		envConfig:      envConfig,
		shioriUsersMap: shioriUsersMap,
	}
}

type service struct {
	envConfig      envConfig
	shioriUsersMap map[string]Credentials
}

type envConfig struct {
	ShioriURL            string `env:"SHIORI_URL,required"`
	ShioriUsersMapString string `env:"SHIORI_USERS,required"`
}

func (s *service) ShioriURL() string {
	return s.envConfig.ShioriURL
}

func (s *service) ShioriUsersMap() map[string]Credentials {
	return s.shioriUsersMap
}
