package cmd

import (
	"log"
	"os"
)

var homeDir string

func init() {
	homeDir = mustRet(os.UserHomeDir())
}

func mustRet[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return value
}
