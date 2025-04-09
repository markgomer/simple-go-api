package database

import (
	"fmt"

	"github.com/go-faker/faker/v4"
)

type user struct {
	FirstName string
	LastName  string
	biography string
}

func (u *user) InitRandomUser() *user {
    u.FirstName = faker.FirstName()
    u.LastName = faker.LastName()
    u.biography = faker.Sentence()
    return u
}

func (u user) PrettyPrint() {
    fmt.Printf("%s %s\nBiography: %s\n\n", u.FirstName, u.LastName, u.biography)
}
