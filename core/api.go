package core

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"go-api/src/database"
)

type PostBody struct {
    Firstname string `json:"firstname,omitempty"`
    Lastname string `json:"lastname,omitempty"`
    Bio string `json:"biography,omitempty"`
}

type Response struct {
    Error string `json:"error,omitempty"`
    UserID string `json:"userId,omitempty"`
    PostBody PostBody `json:"postBody"`
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
                Response{ Error: "That's invalid body, man! Send JSON!" },
                http.StatusUnprocessableEntity,
            )
            return 
        }

        // err checks if user has all the fields
        newUser, err  := database.NewUser(
            body.Firstname,
            body.Lastname,
            body.Bio,
        )
        if err != nil {
            response := Response{ Error: err.Error() }
            sendJSON(rw, response, http.StatusBadRequest)
            return
        }

        // err checks if the new id is unique
        _, newId, err := db.Insert(*newUser)
        if err != nil {
            response := Response{ Error: err.Error() }
            sendJSON(rw, response, http.StatusInternalServerError)
            return
        }

        // Return new user's doc, including id
        sendJSON(
            rw,
            Response{ UserID: uuid.UUID(newId).String(), PostBody: body, },
            http.StatusCreated,
        )
    }
}

func handleFindById(db *database.Application) http.HandlerFunc {
    return (
    func (rw http.ResponseWriter, req *http.Request) {
        urlParam := chi.URLParam(req, "id")       
        if db.FindById(urlParam) == nil {
            slog.Info("Id not found")
            sendJSON(rw, Response{}, http.StatusNotFound)
            return
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

    rw.WriteHeader(status)
    if _, err := rw.Write(data); err != nil {
        slog.Error("failed to write json data", "error", err)
        return
    }
}
