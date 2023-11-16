package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication service on port %s\n", webPort)
	//TODO connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to DB")
	}

	// set up config

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// start the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// check that the connection is actually working
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Cannot connect to DB, retrying...")
			counts++
		} else {
			log.Println("Connected to DB")
			return connection
		}
		if counts > 10 {
			log.Println("Could not connect to DB after 10 tries, exiting...")
			log.Println(err)
			return nil
		}
		log.Println("Waiting 2 seconds before retrying...")
		time.Sleep(2 * time.Second)
		continue
	}
}
