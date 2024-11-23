package cmd

import (
	"log"
)

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)

	}
	return v
}

func requireSingleArg(args []string, wrongArgs func(message string) error) error {
	if len(args) == 0 {
		return wrongArgs("missing argument")
	} else if len(args) > 1 {
		return wrongArgs("too many arguments")
	}
	return nil
}
