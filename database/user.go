package database

import (
	"errors"
	"fmt"

	"github.com/go-faker/faker/v4"
)


type User struct {
    firstName string
    lastName  string
    biography string
}



/* Methods */
func NewUser(firstName string, lastName string, bio string) (*User, error) {
    var err error
    if firstName == "" || bio == "" || lastName == "" {
        err = errors.New("Why dontcha tell me all about ya, ya dirty dawg?")
    }
    newUser := &User{
        firstName: firstName,
        lastName: lastName,
        biography: bio,
    }
    return newUser, err
}


func (u *User) InitRandomUser() *User {
    u.firstName = faker.FirstName()
    u.lastName = faker.LastName()
    u.biography = faker.Sentence()
    return u
}


func (u *User) ToString() string {
    userString := fmt.Sprintf(
        "%s %s\nBiography: %s\n",
        u.firstName,
        u.lastName,
        u.biography,
    )
    return userString
}
