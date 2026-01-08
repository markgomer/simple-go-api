package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type application struct {
	Data map[string]user `json:"data"`
}

func findAll(db *application) {
	for id, entry := range db.Data {
		fmt.Printf("%s: %s %s\n",id,entry.FirstName,entry.LastName)
	}
}

func quickCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	jsonfile, err := os.ReadFile("./mock.json")
	quickCheck(err)

	var dbJson application
	err = json.Unmarshal(jsonfile, &dbJson.Data)
	quickCheck(err)

	findAll(&dbJson)
}
