package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)


type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data"`
}

/*
NOTE: Handlers
*/
func handleGetUsers(dbJSON *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allNames := FindAll(dbJSON)
		sendJSON(w, Response{Data: allNames}, http.StatusOK)
	}
}

func handleGetUserByID(dbJSON *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idQuery := r.PathValue("id")
		idToUpdate, err := strconv.Atoi(idQuery)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusBadRequest)
		}

		userFound, err := FindByID(dbJSON, idToUpdate)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusNotFound)
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
		newID, err := InsertNewUser(dbJSON, u)
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

		_, err = FindByID(dbJSON, idToUpdate)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusNotFound)
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

		updatedUser, err = UpdateUser(dbJSON, idToUpdate, updatedUser)
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

func handleDeleteUser(db *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idQuery := r.PathValue("id")
		idToDelete, err := strconv.Atoi(idQuery)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		err = DeleteUser(db, idToDelete)
		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusNotFound)
			return
		}
		sendJSON(
			w,
			Response{
				Data: fmt.Sprintf("Deleted user id: %d", idToDelete),
			},
			http.StatusOK,
		)
	}
}

/*
  NOTE: helper
*/

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

func SetupHandlers(db *application) (*http.ServeMux) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /api/users",
		handleGetUsers(db),
	)
	mux.HandleFunc(
		"GET /api/users/{id}",
		handleGetUserByID(db),
	)
	mux.HandleFunc(
		"POST /api/users",
		handleInsertNewUser(db),
	)
	mux.HandleFunc(
		"PUT /api/users/{id}",
		handleUpdateUser(db),
	)
	mux.HandleFunc(
		"DELETE /api/users/{id}",
		handleDeleteUser(db),
	)
	return mux
}
