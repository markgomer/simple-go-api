package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"
)

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type application struct {
	Data map[string]user `json:"data"`
}

func findAll(db *application) []string {
	allNames := make([]string, 0, len(db.Data))
	for _, entry := range db.Data {
		allNames = append(allNames, entry.FirstName + " " + entry.LastName)
	}
	sort.Strings(allNames)
	return allNames
}

func lazyCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	jsonfile, err := os.ReadFile("./mock.json")
	lazyCheck(err)

	var dbJson application
	err = json.Unmarshal(jsonfile, &dbJson.Data)
	lazyCheck(err)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /api/users",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, findAll(&dbJson))
		},
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}
	// careful that it locks the program
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
