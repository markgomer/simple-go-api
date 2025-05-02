package database

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/go-faker/faker/v4"
)

/* Types */
type id uuid.UUID

type user struct {
	FirstName string
	LastName  string
	biography string
}

type Application struct {
    data map[id]user
}

/* Methods */
func (u *user) InitRandomUser() *user {
    u.FirstName = faker.FirstName()
    u.LastName = faker.LastName()
    u.biography = faker.Sentence()
    return u
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

func (u user) PrettyPrint() {
    fmt.Printf("%s %s\nBiography: %s\n\n", u.FirstName, u.LastName, u.biography)
}

func (a Application) PrettyPrintAll() {
    for id, user := range a.data {
        fmt.Printf("%v:\n", id)
        user.PrettyPrint()
    }
}

