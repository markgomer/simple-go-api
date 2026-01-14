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

func insertNewUser(db *application, usr user) (int, error) {
	// find next available id
	nextID := 0
	for id := range db.Data {
		if id >= nextID {
			nextID = id + 1
		}
	}

	if usr.FirstName == "" {
		return 0, errors.New("first name missing")
	}
	if usr.LastName == "" {
		return 0, errors.New("last name missing")
	}
	if usr.Biography == "" {
		return 0, errors.New("biography missing")
	}

	db.Data[nextID] = usr
	return nextID, nil
}

func updateUser(db *application, id int, updatedUser user) (user, error) {
	_, err := findByID(db, id)
	if err != nil {
		return user{}, fmt.Errorf("user with id %d not found", id)
	}
	if updatedUser.FirstName == "" {
		return user{}, errors.New("first name missing")
	}
	if updatedUser.LastName == "" {
		return user{}, errors.New("last name missing")
	}
	if updatedUser.Biography == "" {
		return user{}, errors.New("biography missing")
	}
	db.Data[id] = updatedUser
	return updatedUser, nil
}

/*
	NOTE: Handlers
*/

func handleGetUsers(dbJSON *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allNames := findAll(dbJSON)
		sendJSON(w, Response{Data: allNames}, http.StatusOK)
	}
}

func handleGetUserByID(dbJSON *application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idQuery := r.PathValue("id")
		idToUpdate, err := strconv.Atoi(idQuery)
		lazyCheck(err)

		userFound, err := findByID(dbJSON, idToUpdate)
		if err != nil {
			sendJSON(w,Response{Error: err.Error()}, http.StatusNotFound)
			return
		}
		sendJSON(w, Response{Data: userFound}, http.StatusOK)
	}
}

func handleInsertNewUser(dbJSON *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(
					w,
					"Too large of a request",
					http.StatusRequestEntityTooLarge,
				)
				return
			}
			fmt.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			return
		}
		var u user
		err = json.Unmarshal(data, &u)
		if err != nil {
			fmt.Println("Failed to unmarshal json data", "error", err)
			sendJSON(
				w,
				Response{Error: "Internal Server Error"},
				http.StatusInternalServerError,
			)
			return
		}
		newID, err := insertNewUser(dbJSON, u)
		if err != nil {
			fmt.Println(err)
			sendJSON(w, Response{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		sendJSON(w, Response{Data: fmt.Sprintf("New id = %d", newID)}, http.StatusOK)
	}
}

func handleUpdateUser(dbJSON *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idQuery := r.PathValue("id")
		idToUpdate, err := strconv.Atoi(idQuery)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusBadRequest)
			return
		}

		_, err = findByID(dbJSON, idToUpdate)
		if err != nil {
			sendJSON(w,Response{Error: err.Error()}, http.StatusNotFound)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
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
		var updatedUser user
		err = json.Unmarshal(bodyData, &updatedUser)
		if err != nil {
			fmt.Println("Failed to unmarshal json data", "error", err)
			sendJSON(
				w,
				Response{Error: "Internal Server Error"},
				http.StatusInternalServerError,
			)
			return
		}

		updatedUser, err = updateUser(dbJSON, idToUpdate, updatedUser)
		if err != nil {
			fmt.Println(err)
			sendJSON(w, Response{Error: err.Error()}, http.StatusBadRequest)
			return
		}

		sendJSON(
			w,
			Response{
				Data: fmt.Sprintf("User updated = %v", updatedUser),
			},
			http.StatusOK,
		)
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
	mux.HandleFunc(
		"PUT /api/users/{id}",
		handleUpdateUser(&dbJSON),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}
	// careful that it locks the program
	fmt.Println("Server up!")
	err = server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	fmt.Println("Server down!")
}
