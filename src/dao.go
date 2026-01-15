package main

import (
	"encoding/json"
	"os"
)

func LoadDB() *application {
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

