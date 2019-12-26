package main

import (
	"flag"
	"fmt"
	"github.com/simpleHTTPLoginServer/config"
	"github.com/simpleHTTPLoginServer/pkg/handlers"
	"github.com/simpleHTTPLoginServer/pkg/mySqlClient"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	var env handlers.Env
	var configFile string

	flag.StringVar(&configFile, "c", "noConfig", "conifg file without extension")
	flag.Parse()
	if configFile == "noConfig" {
		log.Fatal("No config file found, Please use --help command")
	}

	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	var conf config.SystemConig

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode configuration, %v", err)
	}

	mysqlHost := conf.MysqlConfig.Host
	mysqlPort := conf.MysqlConfig.Port
	username := conf.MysqlConfig.UserName
	password := conf.MysqlConfig.Password //ideally loads from system ENV
	dbName := conf.MysqlConfig.DataBase

	env.JwtKey = []byte(conf.ServerConfig.JwtKey) //ideally loads from system ENV

	env.Mysqlclient = &mySqlClient.Client{mysqlHost,
		mysqlPort,
		username,
		password,
		dbName,
		nil,
	}

	db, err := env.Mysqlclient.DbConnect()
	if err != nil {
		log.Fatalf("unable to connect with db, %v", err)
	}

	log.Print("My sql successfully conneted")
	env.Mysqlclient.DbConn = db
	http.HandleFunc("/user/signup", env.SignUp)
	http.HandleFunc("/user/login", env.LogIn)
	http.HandleFunc("/user", env.GetUserInfo)

	log.Print("Listening on port: 9000")
	http.ListenAndServe(":9000", nil)
}
