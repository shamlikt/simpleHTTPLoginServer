package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type mysqlConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
	DataBase string `yaml:"DataBase"`
}

type serverConfig struct {
	JwtKey  string `yaml:"MysqlConfig"`
	LogPath string `yaml:"LogPath"`
}

type systemConig struct {
	MysqlConfig  mysqlConfig  `yaml:"MysqlConfig"`
	ServerConfig serverConfig `yaml:"ServerConfig"`
}

type Config struct {
	Config systemConig `yaml:"Config"`
}

func ParseConfig(fileName string) (Config, error) {
	var config Config
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
