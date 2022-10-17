package util

import (
	uuid "github.com/satori/go.uuid"
)

func Uuid() string {
	id := uuid.NewV4().String()
	return id
}
