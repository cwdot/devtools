package cmd

import (
	"log"

	"github.com/pkg/errors"
)

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)

	}
	return v
}

func requireSingleArg(args []string, noArgs func() error) (bool, error) {
	if len(args) == 0 {
		return true, noArgs()
	} else if len(args) > 1 {
		return true, errors.New("too many arguments")
	}
	return false, nil
}
