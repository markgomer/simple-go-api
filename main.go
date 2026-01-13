package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type application struct {
	Data map[int]user `json:"data"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data"`
}

func findAll(db *application) []string {
	allNames := make([]string, 0, len(db.Data))
	for id, entry := range db.Data {
		idStr := fmt.Sprintf("%d", id)
		allNames = append(allNames, entry.FirstName+" "+entry.LastName+"id: "+idStr)
	}
	sort.Strings(allNames)
	return allNames
}

func findByID(db *application, id int) (user, error) {
	userFound, exists := db.Data[id]
	if !exists {
		return user{}, fmt.Errorf("user with id %d not found", id)
	}
	return userFound, nil
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

func handleGetUserByID(dbJSON *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idQuery := r.PathValue("id")
		idInt, err := strconv.Atoi(idQuery)
		lazyCheck(err)

		userFound, err := findByID(dbJSON, idInt)
		if err != nil {
			sendJSON(w,Response{Error: err.Error()}, http.StatusNotFound)
			return
		}
		sendJSON(w, Response{Data: userFound}, http.StatusOK)
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
	mux.HandleFunc(
		"GET /api/users/{id}",
		handleGetUserByID(&dbJSON),
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
