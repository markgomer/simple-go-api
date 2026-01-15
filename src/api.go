package main

import (
	"encoding/json"
	"os"
)

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type application struct {
	Data map[int]user `json:"data"`
}

func LoadDB() (*application) {
	var dbJSON application

	jsonfile, err := os.ReadFile("../mock.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonfile, &dbJSON.Data)
	if err != nil {
		panic(err)
	}
	if dbJSON.Data == nil {
		dbJSON.Data = make(map[int]user)
	}
	return &dbJSON
}
