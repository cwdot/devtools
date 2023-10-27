package util

import (
	"fmt"
	"strings"
)

func CreateCsvLinks(base string, csvLinks string) string {
	newLinks := make([]string, 0, 5)
	items := strings.Split(csvLinks, ",")
	base = strings.TrimRight(base, "/")
	for _, item := range items {
		item = strings.TrimLeft(item, "/")
		newLinks = append(newLinks, fmt.Sprintf("%s/%s", base, item))
	}
	return strings.Join(newLinks, " ")
}
