package flibusta

import (
	"regexp"
	"strconv"
)

func rexGroup(rex *regexp.Regexp, src []byte, group int) string {
	items := rex.FindSubmatch(src)
	if len(items) < group {
		return ""
	}

	return string(items[group])
}

func rexGroupAll(rex *regexp.Regexp, src []byte, group int) []string {
	var items []string

	for _, item := range rex.FindAllSubmatch(src, -1) {
		if len(item) < group {
			continue
		}

		items = append(items, string(item[group]))
	}

	return items
}

func totalPages(content []byte) []int {
	var items []int

	for _, item := range rexPagination.FindAllSubmatch(content, -1) {
		if len(item) < 1 {
			continue
		}

		page, err := strconv.Atoi(string(item[1]))
		if err != nil {
			continue
		}

		items = append(items, page)
	}

	return items
}
