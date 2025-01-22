package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    _ "github.com/lib/pq"
)

type ServiceEndpoint struct {
    IPAddress   string
    Port        int
}

type ServiceEndpointList struct {
    lock        sync.RWMutex
    endpoints   []ServiceEndpoint
}

var serviceMapLock sync.RWMutex
var serviceMap = make(map[int]*ServiceEndpointList)

func main() {
    
    port := os.Getenv("PS_DISCOVERY_PORT")
    if port == "" {
        port = "4900"
    }

    dbHost := os.Getenv("PS_DATABASE_HOST")
    if dbHost == "" {
        dbHost = "localhost"
    }

    dbPort := os.Getenv("PS_DATABASE_PORT")
    if dbPort == "" {
        dbPort = "5432"
    }

    dbUser := os.Getenv("PS_DATABASE_USER")
    if dbUser == "" {
        dbUser = "dbuser"
    }

    dbPassword := os.Getenv("PS_DATABASE_PASSWORD")
    if dbPassword == "" {
        dbPassword = "dbpass"
    }

    dbName := os.Getenv("PS_DATABASE_NAME")
    if dbName == "" {
        dbName = "portfolio_simulator"
    }

    dbConnectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
    db, err := sql.Open("postgres", dbConnectionStr)
    if err != nil {
        log.Fatalf("Could not open db: %v\n", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Could not ping db: %v\n", err)
    }

    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Post("/service", createServiceHandler(db))
    r.Post("/register", registerServiceNodeHandler(db))

    log.Printf("Starting server on :%s...\n", port)
    if err := http.ListenAndServe(":" + port, r); err != nil {
        log.Fatal(err)
    }
}
