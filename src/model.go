package main

import (
	"errors"
	"fmt"
	"sort"
)

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type application struct {
	Data map[int]user `json:"data"`
}

func FindAll(db *application) []string {
	allNames := make([]string, 0, len(db.Data))
	for id, entry := range db.Data {
		idStr := fmt.Sprintf("%d", id)
		allNames = append(allNames, entry.FirstName+" "+entry.LastName+" id: "+idStr)
	}
	sort.Strings(allNames)
	return allNames
}

func FindByID(db *application, id int) (user, error) {
	userFound, exists := db.Data[id]
	if !exists {
		return user{}, fmt.Errorf("user with id %d not found", id)
	}
	return userFound, nil
}

func InsertNewUser(db *application, usr user) (int, error) {
	// find next available id
	nextID := 0
	for id := range db.Data {
		if id >= nextID {
			nextID = id + 1
		}
	}

	if usr.FirstName == "" {
		return 0, errors.New("first name missing")
	}
	if usr.LastName == "" {
		return 0, errors.New("last name missing")
	}
	if usr.Biography == "" {
		return 0, errors.New("biography missing")
	}

	db.Data[nextID] = usr
	return nextID, nil
}

func UpdateUser(db *application, id int, updatedUser user) (user, error) {
	_, err := FindByID(db, id)
	if err != nil {
		return user{}, fmt.Errorf("user with id %d not found", id)
	}
	if updatedUser.FirstName == "" {
		return user{}, errors.New("first name missing")
	}
	if updatedUser.LastName == "" {
		return user{}, errors.New("last name missing")
	}
	if updatedUser.Biography == "" {
		return user{}, errors.New("biography missing")
	}
	db.Data[id] = updatedUser
	return updatedUser, nil
}

func DeleteUser(db *application, idToDelete int) error {
	_, err := FindByID(db, idToDelete)
	if err != nil {
		return err
	}

	delete(db.Data, idToDelete)
	return nil
}
