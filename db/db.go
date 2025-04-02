package db

import "github.com/google/uuid"

type id uuid.UUID

type user struct {
	FirstName string
	LastName  string
	biography string
}

type application struct {
	data map[id]user
}
