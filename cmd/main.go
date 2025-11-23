package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"pr_service/internal/app"
	"pr_service/internal/repository"
	deps "pr_service/pkg/gen"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres:postgres@postgres:5432/appdb?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	userRepo := repository.NewRepository(db)

    service, _ := app.NewService(userRepo)

    srv, err := deps.NewServer(service)
    if err != nil {
        log.Fatal(err)
    }
    if err := http.ListenAndServe(":8080", srv); err != nil {
        log.Fatal(err)
    }
}