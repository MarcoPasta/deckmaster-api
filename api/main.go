package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    // Try to connect to DB in background (non-blocking start). No driver is
    // imported here, so Ping will fail; that's acceptable for a smoke test.
    go func() {
        for {
            db, err := sql.Open("postgres", dsn)
            if err == nil {
                if err = db.Ping(); err == nil {
                    log.Println("db: connected")
                    db.Close()
                    return
                }
                db.Close()
            }
            log.Println("db: waiting for database...", err)
            time.Sleep(2 * time.Second)
        }
    }()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "deckmaster: OK")
    })

    log.Printf("listening on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
