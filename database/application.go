package database

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
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
func NewUser(firstName string, lastName string, bio string) *user {
    return &user{
        FirstName: firstName,
        LastName: lastName,
        biography: bio,
    }
}

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
        u := &user{}
        u.InitRandomUser()
        uid := uuid.New()
        app.data[id(uid)] = *u
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

func (a Application) FindById(query string) *user {
    var wanted *user
    parsedQuery, err := uuid.Parse(query)
    if err != nil {
        slog.Error("Invalid ID", "error", err)
        return nil
    }
    uid := id(parsedQuery)
    for id, user := range a.data {
        if id == uid {
            wanted = &user
        }
    }
    if wanted == nil {
        slog.Info("No user found with id: ", "info", query)
    }
    return wanted
}

func Insert(u user, a Application) {
}

func Update() { }

func Delete() { }

func (u user) ToString() string {
    userString := fmt.Sprintf(
        "%s %s\nBiography: %s\n",
        u.FirstName,
        u.LastName,
        u.biography,
    )
    return userString
}

func (a Application) ToString() string {
    var appString strings.Builder
    for id, user := range a.data {
        uid := uuid.UUID(id).String()
        fmt.Fprintf(&appString, "%s:\n%s", uid, user.ToString())
    }
    return appString.String()
}

