package database

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type id uuid.UUID

type Application struct {
	data map[id]user
}

func InitWithRandom(numberOfEntries int) *Application {
    app := &Application{}
    app.data = make(map[id]user)
    for range numberOfEntries {
        uid := id(uuid.New())
        u := &user{}
        u.InitRandomUser()
        app.data[uid] = *u
    }
    return app
}

func InitEmpty() *Application {
    app := &Application{}
    app.data = make(map[id]user)
    return app
}

func (a Application) FindAll() []user {
    var userSlice []user
    for _, user := range a.data {
        userSlice = append(userSlice, user)
    }
    return userSlice
}

func (a Application) FindById(query id) user {
    var wanted *user
    for id, user := range a.data {
        if id == query {
            wanted = &user
        }
    }
    if wanted == nil {
        slog.Info("No user found with id: ", "info", query)
    }
    return *wanted
}

func (a Application) PrettyPrintAll() {
    for id, user := range a.data {
        fmt.Printf("%v:\n", id)
        user.PrettyPrint()
    }
}
