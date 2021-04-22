package models

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "arjuna.db.elephantsql.com"
	port     = 5432
	user     = "ydhobhyt"
	password = "FckQIjb9THM3s8rLbxhzNhAyV765Yj2-"
	dbname   = "ydhobhyt"
)

var (
	DB *sql.DB
)

func ConnectToDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		DB, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		// defer DB.Close()

		err = DB.Ping()
		if err != nil {
			panic(err)
		}

		fmt.Println("Successfully connected!")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
