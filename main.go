package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-api/src/core"
	"go-api/src/database"
)


func main() {
    slog.Info("Service initialized")
    if err := run(); err != nil {
        slog.Error("failed to execute code", "error", err)
        os.Exit(1)
    }
    slog.Info("All systems offline")
}


func run() error {
    db := database.InitWithRandom(1)
    fmt.Println(db.ToString())

    slog.Info("Creating Handler")
    handler := core.NewHandler(db)
    slog.Info("Handler Created")

    server := http.Server{
        Addr: ":8080",
        Handler: handler,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout: time.Minute,
    }
    slog.Info("Server created and running on http://localhost:8080")
    if err := server.ListenAndServe(); err != nil {
        slog.Error("Internal Server Error", "error", err)
        return err
    }

    return nil
}
