package database

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

/* Types */
type Id uuid.UUID


type Application struct {
    data map[Id]User
}


func InitEmpty() *Application {
    app := &Application{}
    app.data = make(map[Id]User)
    return app
}


func InitWithRandom(numberOfEntries int) *Application {
    app := InitEmpty()
    for range numberOfEntries {
        u := &User{}
        u.InitRandomUser()
        uid := uuid.New()
        app.data[Id(uid)] = *u
    }
    return app
}


func (a *Application) FindAll() []User {
    var userSlice []User
    for _, User := range a.data {
        userSlice = append(userSlice, User)
    }
    return userSlice
}


func (a *Application) FindById(query string) *User {
    var wanted *User
    parsedQuery, err := uuid.Parse(query)
    if err != nil {
        slog.Error("Invalid ID", "error", err)
        return nil
    }
    uid := Id(parsedQuery)
    for Id, User := range a.data {
        if Id == uid {
            wanted = &User
        }
    }
    if wanted == nil {
        slog.Info("No User found with Id: ", "info", query)
    }
    return wanted
}


func (a *Application) Insert(u User) (User, Id, error) {
    var (
        err error
        newUser User
    )
    newid := Id(uuid.New())

    if _, ok := a.data[newid]; ok {
        err = errors.New("User Id already exists")
        newUser = User{}
    } else {
        a.data[newid] = u
    }

    return newUser, newid, err
}


func (a *Application) Update() { }


func (a *Application) Delete() { }


func (a *Application) ToString() string {
    var appString strings.Builder
    for Id, User := range a.data {
        uid := uuid.UUID(Id).String()
        fmt.Fprintf(&appString, "%s:\n%s", uid, User.ToString())
    }
    return appString.String()
}

