package database


import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)


/* Types */
type id uuid.UUID


type Application struct {
    data map[id]User
}


func InitEmpty() *Application {
    app := &Application{}
    app.data = make(map[id]User)
    return app
}


func InitWithRandom(numberOfEntries int) *Application {
    app := InitEmpty()
    for range numberOfEntries {
        u := &User{}
        u.InitRandomUser()
        uid := uuid.New()
        app.data[id(uid)] = *u
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
    uid := id(parsedQuery)
    for id, User := range a.data {
        if id == uid {
            wanted = &User
        }
    }
    if wanted == nil {
        slog.Info("No User found with id: ", "info", query)
    }
    return wanted
}


func (a *Application) Insert(u User) (string, User) {
    newid := id(uuid.New())
    slog.Debug("id(newid)=", "debug", newid)
    a.data[newid] = u
    stringID := uuid.UUID(newid).String()
    slog.Debug(stringID)
    return stringID, u
}


func (a *Application) Update() { }


func (a *Application) Delete() { }


func (a *Application) ToString() string {
    var appString strings.Builder
    for id, User := range a.data {
        uid := uuid.UUID(id).String()
        fmt.Fprintf(&appString, "%s:\n%s", uid, User.ToString())
    }
    return appString.String()
}

