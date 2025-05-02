package core

import (
	"go-api/src/database"
	"log/slog"

	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PostBody struct {
    URL string `json:"url"`
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
    router.Post("/api/users", handleFunc())
    router.Get("/api/users", handleFunc())
    router.Get("/api/users/{id}", handleFindById(db))
    router.Delete("/api/users/{id}", handleFunc())
    router.Put("/api/users/{id}", handleFunc())

    return router
}


func handleFunc() http.HandlerFunc {
    return func (rw http.ResponseWriter, req *http.Request) {
        return
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


func sendJSON(rw http.ResponseWriter, resp Response, status int) {
    rw.Header().Set("Content-Type", "application/json")

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
