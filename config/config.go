package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type mysqlConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	DataBase string
}

type serverConfig struct {
	JwtKey  string
	LogPath string
}

type SystemConig struct {
	MysqlConfig  mysqlConfig
	ServerConfig serverConfig
}

func ParseConfig(fileName string) (SystemConig, error) {
	var config SystemConig
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
