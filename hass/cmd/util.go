package cmd

import "log"

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)

	}
	return v
}
