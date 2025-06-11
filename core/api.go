package core

import (
	"fmt"
	"log/slog"

	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-api/src/database"
)

type PostBody struct {
    Firstname string `json:"firstname"`
    Lastname string `json:"lastname"`
    Bio string `json:"bio"`
}

type Response struct {
    Error string `json:"error,omitempty"`
    Data any `json:"data,omitempty"`
}

func NewHandler(db *database.Application) http.Handler {
    router := chi.NewMux()

    // call middlewares
    router.Use(middleware.Recoverer)
    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)

    /**
     * Set routes/endpoints
     *
     * Those functions have to be in the format:
     * func(rw http.ResponseWriter, req *http.Request) => void
    **/
    router.Post("/api/users", handleInsertUser(db))
    router.Get("/api/users", handleFunc())
    router.Get("/api/users/{id}", handleFindById(db))
    router.Delete("/api/users/{id}", handleFunc())
    router.Put("/api/users/{id}", handleFunc())
    router.NotFound(handleNotFound())

    return router
}


func handleFunc() http.HandlerFunc {
    return func (rw http.ResponseWriter, req *http.Request) {
        return
    }
}

/***
  curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"firstname":"John","lastname":"Doe","biography":"This is my bio."}' \
    http://localhost:8080/api/users
***/
func handleInsertUser(db *database.Application) http.HandlerFunc {
    return func (rw http.ResponseWriter, req *http.Request) {
        // We will fill this struct with the post body, which is in json
        var body PostBody

        // Here we fill the var body with the request(POST) body
        // if the user is trying to send invalid body, we kick his back
        if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
            sendJSON(
                rw,
                Response{Error: "That's invalid body, man! Send JSON!"},
                http.StatusUnprocessableEntity,
            )
            return 
        }
        newUser, err  := database.NewUser(
            body.Firstname,
            body.Lastname,
            body.Lastname,
        )
        // Check if json is complete
        if err != nil {
            fmt.Printf("ERROR=%s\n", err)
            sendJSON(
                rw,
                Response{
                    Error: err.Error(),
                },
                http.StatusBadRequest,
            )
            return
        }
        // TODO if user info is valid, then
        // TODO save user in database
        db.Insert(*newUser)
        // TODO answer with HTTP 201 (created)
        // TODO return new user's doc, including id

    }
}

func handleFindById(db *database.Application) http.HandlerFunc {
    return (
    func (rw http.ResponseWriter, req *http.Request) {
        urlParam := chi.URLParam(req, "id")       
        if db.FindById(urlParam) == nil {
            slog.Info("Id not found")
        }

        sendJSON(rw, Response{}, http.StatusCreated)
        return
    })
}

func handleNotFound() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        slog.Error("Route not found:", "URL path", r.URL.Path)
        return
    }
}

func sendJSON(rw http.ResponseWriter, resp Response, status int) {
    rw.Header().Set("Content-Type", "application/json")

    // build Json
    data, err := json.Marshal(resp)
    if err != nil {
        slog.Error("failed to marshal json data", "error", err)
        sendJSON(
            rw,
            Response{Error: "something went wrong"},
            http.StatusInternalServerError,
        )
        return
    }

    // write header and 
    rw.WriteHeader(status)
    if _, err := rw.Write(data); err != nil {
        slog.Error("failed to write json data", "error", err)
        return
    }
}
