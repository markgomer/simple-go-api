package main

import (
	"go-api/src/core"
	"go-api/src/database"
	"log/slog"
	"net/http"
	"os"
	"time"
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
    db.PrettyPrintAll()

    slog.Info("Creating Handler")
    handler := core.NewHandler()
    slog.Info("Handler Created")

    server := http.Server{
        Addr: ":8080",
        Handler: handler,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout: time.Minute,
    }
    slog.Info("Server Created")

    slog.Info("Server Running")
    if err := server.ListenAndServe(); err != nil {
        slog.Error("Internal Server Error", "error", err)
        return err
    }

    return nil
}
