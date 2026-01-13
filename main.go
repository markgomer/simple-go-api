package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func insertNewUser(db *application, usr user) (int) {
	nextID := 0
	for id := range db.Data {
		if id >= nextID {
			nextID = id + 1
		}
	}
	db.Data[nextID] = usr
	return nextID
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

func handleInsertNewUser(dbJSON *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			// max bytes error
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(
					w,
					"Too large of a request",
					http.StatusRequestEntityTooLarge,
				)
			}
			fmt.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
		}
		var u user
		err = json.Unmarshal(data, &u)
		if err != nil {
			http.Error(w, "Invalid user format", http.StatusUnprocessableEntity)
			return
		}
		if u.FirstName == "" {
			http.Error(w, "First name missing", http.StatusBadRequest)
			return
		}
		if u.LastName == "" {
			http.Error(w, "Last name missing", http.StatusBadRequest)
			return
		}
		if u.Biography == "" {
			http.Error(w, "Biography missing", http.StatusBadRequest)
			return
		}
		insertNewUser(dbJSON, u)
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

func lazyCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	jsonfile, err := os.ReadFile("./mock.json")
	lazyCheck(err)

	var dbJSON application
	err = json.Unmarshal(jsonfile, &dbJSON.Data)
	lazyCheck(err)

	if dbJSON.Data == nil {
		dbJSON.Data = make(map[int]user)
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/users",
		handleGetUsers(&dbJSON),
	)
	mux.HandleFunc(
		"GET /api/users/{id}",
		handleGetUserByID(&dbJSON),
	)
	mux.HandleFunc(
		"POST /api/users",
		handleInsertNewUser(&dbJSON),
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
