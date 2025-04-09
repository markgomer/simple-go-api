package core

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler() http.Handler {
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
    router.Get("/api/users/:id", handleFunc())
    router.Delete("/api/users/:id", handleFunc())
    router.Put("/api/users/:id", handleFunc())

    return router
}


func handleFunc() http.HandlerFunc {
    return func (rw http.ResponseWriter, req *http.Request) {
        return
    }
}
