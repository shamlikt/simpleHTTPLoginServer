package mySqlClient

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shamlikt/simpleHTTPLoginServer/pkg/utils"
)

const dbDriver string = "mysql"

type Client struct {
	Server   string
	Port     int
	UserName string
	Password string
	Dbname   string
	DbConn   *sql.DB
}

func (c Client) DbConnect() (db *sql.DB, err error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.UserName,
		c.Password, c.Server, c.Port, c.Dbname)
	db, err = sql.Open(dbDriver, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil

}

func (c Client) InsertUser(userInfo utils.UserInfo) (err error) {
	stmtIns, err := c.DbConn.Prepare("INSERT INTO users(name, password, email, age, salary)  VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	fmt.Println(userInfo.Password)
	hPass, err := utils.HashPassword(userInfo.Password)
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(userInfo.UserName, hPass, userInfo.Email, userInfo.Age, userInfo.Salary)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) ValidateUser(username string, password string) (isValid bool, err error) {
	userStmt, err := c.DbConn.Prepare("SELECT password FROM users WHERE name = ?")
	if err != nil {
		return false, err
	}
	defer userStmt.Close()

	var hPass string
	err = userStmt.QueryRow(username).Scan(&hPass)
	if err != nil {
		return false, err
	}

	if ok := utils.CompareHashAndPassword(hPass, password); ok {
		return true, nil
	} else {
		return false, nil
	}

}

func (c Client) GetUserData(username string, userInfo *utils.UserInfo) (err error) {
	userStmt, err := c.DbConn.Prepare("SELECT name, age, email, salary FROM users WHERE name = ?")
	if err != nil {
		return err
	}
	defer userStmt.Close()
	var name string
	var email string
	var age int
	var salary int

	err = userStmt.QueryRow(username).Scan(&name, &age, &email, &salary)
	if err != nil {
		return err
	}
	userInfo.UserName = name
	userInfo.Email = email
	userInfo.Age = age
	userInfo.Salary = salary
	return nil
}
