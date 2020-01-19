package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

	if len(e) != 0 {
		result.Status = false
		result.Message = e.Encode()
		resultJson, _ := json.Marshal(result)
		w.Write(resultJson)
		return
	}

	w.Header().Set("Content-type", "applciation/json")
	if len(e) != 0 {
		result.Status = false
		result.Message = e.Encode()
		resultJson, _ := json.Marshal(result)
		w.Write(resultJson)
		return
	}

	fmt.Println(env.Mysqlclient)
	err := env.Mysqlclient.InsertUser(user)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				result.Status = false
				result.Message = "The user already registered"
				resultJson, _ := json.Marshal(result)
				w.Write(resultJson)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result.Status = true
	result.Message = "User Added successfully"
	resultJson, _ := json.Marshal(result)
	w.Write(resultJson)
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
		result.Status = false
		result.Message = e.Encode()
		resultJson, _ := json.Marshal(result)
		w.Write(resultJson)
		return
	}

	username := auth.UserName
	password := auth.Password
	fmt.Println(username)
	fmt.Println(password)
	isValid, err := env.Mysqlclient.ValidateUser(username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			result.Status = false
			result.Message = "No User Found"
			resultJson, _ := json.Marshal(result)
			w.Write(resultJson)
			w.WriteHeader(http.StatusUnauthorized)
			return

		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isValid {
		result.Status = false
		result.Message = "Incorrect username or password"
		resultJson, _ := json.Marshal(result)
		w.Write(resultJson)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenStr, err := utils.ClaimsJWT(username, env.JwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result.Status = true
	result.Token = tokenStr
	result.Message = "User Authenticated successfully"
	resultJson, _ := json.Marshal(result)
	w.Write(resultJson)
	return
}

func (env *Env) GetUserInfo(w http.ResponseWriter, r *http.Request) {

	var result utils.UserInfoResult
	var userInfo utils.UserInfo

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Auth-token")
	w.Header().Set("Content-type", "applciation/json")

	if token == "" {
		result.Status = false
		result.Message = "Auth-token: Header not found"
		resultJson, _ := json.Marshal(result)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resultJson)
		return
	}

	claims := &utils.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return env.JwtKey, nil
	})

	if !tkn.Valid {
		result.Status = false
		result.Message = "Invalid token"
		resultJson, _ := json.Marshal(result)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resultJson)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			result.Status = false
			result.Message = "Invalid token"
			resultJson, _ := json.Marshal(result)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resultJson)
			return
		}
		result.Status = false
		result.Message = "Bad Request"
		resultJson, _ := json.Marshal(result)
		w.Write(resultJson)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Mysqlclient.GetUserData(claims.Username, &userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result.Status = true
	result.UserInfo = userInfo
	result.Message = "Got user information "
	resultJson, _ := json.Marshal(result)
	w.Write(resultJson)
	return

}
