package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

/*
NOTE: Entrypoint
*/
func main() {
	// "Load" "DB"
	dbJSON := LoadDB()

	mux := SetupHandlers(dbJSON)

	// setup server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}
	// Start server
	// careful that it locks the program
	fmt.Println("Server up!")
	err := server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	fmt.Println("Server down!")
}
