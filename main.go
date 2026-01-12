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

type Response struct {
	Error    string `json:"error,omitempty"`
	Data any   `json:"postBody"`
}

func findAll(db *application) []string {
	allNames := make([]string, 0, len(db.Data))
	for id, entry := range db.Data {
		allNames = append(allNames, id + ": " + entry.FirstName+" "+entry.LastName)
	}
	sort.Strings(allNames)
	return allNames
}

func lazyCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func sendJSON(rw http.ResponseWriter, resp Response, status int) {
	rw.Header().Set("Content-Type", "application/json")
	// build Json
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Failed to marshal json data", "error", err)
		sendJSON(
			rw,
			Response{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}
	rw.WriteHeader(status)
	_, err = rw.Write(data)
	if err != nil {
		fmt.Println("Failed to write json data", "error", err)
		return
	}
}

func handleGetUsers(dbJSON *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allNames := findAll(dbJSON)
		sendJSON(w, Response{Data: allNames}, http.StatusOK)
	}
}

func main() {
	jsonfile, err := os.ReadFile("./mock.json")
	lazyCheck(err)

	var dbJSON application
	err = json.Unmarshal(jsonfile, &dbJSON.Data)
	lazyCheck(err)

	mux := http.NewServeMux()

	// endpoint in which we send response with list of everyone
	mux.HandleFunc(
		"GET /api/users",
		handleGetUsers(&dbJSON),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}
	// careful that it locks the program
	err = server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
