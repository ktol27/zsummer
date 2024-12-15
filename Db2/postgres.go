package Db2

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

var DB *sql.DB

//DB_HOST=localhost;DB_NAME=02;DB_PASSWORD=982655;DB_PORT=5432;DB_USER=postgres

func InitDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	//if err != nil {
	//	main2.ErrorLogger.Error("Database connection error:", err)
	//	log.Fatal(err)
	//}

	for i := 0; i < 5; i++ {
		if err = DB.Ping(); err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}
