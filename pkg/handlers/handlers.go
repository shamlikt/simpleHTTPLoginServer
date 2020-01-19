package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-sql-driver/mysql"
	"github.com/shamlikt/simpleHTTPLoginServer/pkg/mySqlClient"
	"github.com/shamlikt/simpleHTTPLoginServer/pkg/utils"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

type Result struct {
	Status  bool
	Message string
	Token   string
}

type SignUpResult struct {
	Status  bool
	Message string
}

type Env struct {
	JwtKey      []byte
	Mysqlclient *mySqlClient.Client
	TokenAuth   *jwtauth.JWTAuth
}

func (env *Env) SignUp(w http.ResponseWriter, r *http.Request) {
	var user utils.UserInfo
	var result SignUpResult

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rules := govalidator.MapData{
		"username": []string{"required", "max:20"},
		"password": []string{"required", "min:4"},
		"email":    []string{"required", "email"},
		"age":      []string{"required", "numeric_between:0,200"},
		"salary":   []string{"required", "min:0"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &user,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	w.Header().Set("Content-type", "applciation/json")
	if len(e) != 0 {
		respondWithError(w, 400, "Bad request")
		return
	}

	err := env.Mysqlclient.InsertUser(user)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				result.Status = false
				result.Message = "The user already registered"
				respondwithJSON(w, 403, result)
				return
			}
		}
		respondWithError(w, 500, "Internal Error")
		return
	}

	result.Status = true
	result.Message = "User Added successfully"
	respondwithJSON(w, 200, result)
	return

}

func (env *Env) LogIn(w http.ResponseWriter, r *http.Request) {
	var auth utils.AuthInfo
	var result Result

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rules := govalidator.MapData{
		"username": []string{"required"},
		"password": []string{"required"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &auth,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) != 0 {
		respondWithError(w, 400, "Bad request")
		return
	}

	username := auth.UserName
	password := auth.Password
	fmt.Println(username)
	fmt.Println(password)
	isValid, err := env.Mysqlclient.ValidateUser(username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 403, "User not found")
			return

		}
		respondWithError(w, 500, "internal error")
		return
	}

	if !isValid {
		respondWithError(w, 403, "Incorrect username or password")
		return
	}

	_, tokenStr, _ := env.TokenAuth.Encode(jwt.MapClaims{"username": username})
	if err != nil {
		respondWithError(w, 500, "internal error")
		return
	}

	result.Status = true
	result.Token = tokenStr
	result.Message = "User Authenticated successfully"
	respondwithJSON(w, 200, result)
	return
}

func (env *Env) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var result utils.UserInfoResult
	var userInfo utils.UserInfo

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	username, err := claims["username"].(string)

	if err != true {
		respondWithError(w, 400, "Bad request")
		return
	}

	dberr := env.Mysqlclient.GetUserData(username, &userInfo)
	if dberr != nil {
		respondWithError(w, 500, "Internal server error")
	}

	result.Status = true
	result.UserInfo = userInfo
	result.Message = "Got user information "
	respondwithJSON(w, 200, result)
	return

}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
