package flibusta

import "regexp"

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
