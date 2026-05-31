package config

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    var err error
    dsn := "root:@tcp(127.0.0.1:3306)/ritel?charset=utf8mb4&parseTime=True&loc=Local"
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        panic(err.Error())
    }

    err = DB.Ping()
    if err != nil {
        panic(err.Error())
    }

    fmt.Println("Database connected")
}