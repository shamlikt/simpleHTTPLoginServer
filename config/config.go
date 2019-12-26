package config

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
