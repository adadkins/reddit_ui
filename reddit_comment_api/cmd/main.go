package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reddit_comment_api/pkg/reddit_comment_api"
	"time"

	"github.com/jmoiron/sqlx"
)

func main() {
	db := connectDB()
	a := reddit_comment_api.NewApp(db)

	err := a.Start()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func connectDB() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB_NAME"))
	fmt.Println(psqlInfo)
	db := &sqlx.DB{}
	err := errors.New("test")
	for i := 0; i < 10; i++ {
		db, err = sqlx.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		log.Println("PING...")
		err = db.Ping()
		if err != nil {
			log.Println("Ping unsuccessfull")
			log.Println(err)
			log.Println("Sleeping... then retrying")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	log.Println("PONG")
	return db
}
