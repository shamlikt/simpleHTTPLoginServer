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

	conf, err := config.ParseConfig(configFile)
	if err != nil {
		// log.Fatal("Error while parsing config")
		log.Fatal(err)

	}

	mysqlHost := conf.MysqlConfig.Host
	mysqlPort := conf.MysqlConfig.Port
	username := conf.MysqlConfig.UserName
	password := conf.MysqlConfig.Password
	dbName := conf.MysqlConfig.DataBase

	env.TokenAuth = jwtauth.New("HS256", []byte(conf.ServerConfig.JwtKey), nil)
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
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/data", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	log.Print("Listening on port: 9000")
	http.ListenAndServe(":9000", r)
}
