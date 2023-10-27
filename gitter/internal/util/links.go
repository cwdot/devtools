package util

import (
	"fmt"
	"strings"
)

func CreateCsvLinks(base string, csvLinks string) string {
	newLinks := make([]string, 0, 5)
	items := strings.Split(csvLinks, ",")
	for _, item := range items {
		newLinks = append(newLinks, fmt.Sprintf("%s/%s", base, item))
	}
	return strings.Join(newLinks, " ")
}
