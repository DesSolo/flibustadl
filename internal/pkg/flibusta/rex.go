package flibusta

import "regexp"

var (
	rexFB2        = regexp.MustCompile(`<a href="(/b/\d+/fb2)">\(fb2\)</a>`)
	rexTitle      = regexp.MustCompile(`<h1 class="title">([^<]+)</h1>`)
	rexPagination = regexp.MustCompile(`<li class="pager-item"><a href=".*?page=(\d+)"`)
)
