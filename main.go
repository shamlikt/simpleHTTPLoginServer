package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/shamlikt/simpleHTTPLoginServer/config"
	"github.com/shamlikt/simpleHTTPLoginServer/pkg/handlers"
	"github.com/shamlikt/simpleHTTPLoginServer/pkg/mySqlClient"
	"log"
	"net/http"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	var env handlers.Env
	var configFile string

	flag.StringVar(&configFile, "c", "", "conifg file without extension")
	flag.Parse()
	if configFile == "" {
		log.Fatal("No config file found, Please use --help command")
	}

	c, err := config.ParseConfig(configFile)
	if err != nil {
		// log.Fatal("Error while parsing config")
		log.Fatal(err)

	}

	mysqlHost := c.Config.MysqlConfig.Host
	mysqlPort := c.Config.MysqlConfig.Port
	username := c.Config.MysqlConfig.UserName
	password := c.Config.MysqlConfig.Password
	dbName := c.Config.MysqlConfig.DataBase
	fmt.Println(mysqlHost)
	fmt.Println(mysqlPort)

	env.TokenAuth = jwtauth.New("HS256", []byte(c.Config.ServerConfig.JwtKey), nil)
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

	log.Print("Mysql successfully conneted")
	env.Mysqlclient.DbConn = db

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/signup", env.SignUp)
		r.Post("/login", env.LogIn)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(env.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/data", env.GetUserInfo)

	})

	log.Print("Listening on port: 9000")
	http.ListenAndServe(":9000", r)
}
