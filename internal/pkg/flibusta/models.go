package flibusta

import "io"

type Sequence struct {
	Name     string
	BookURLs []string
}

type Author struct {
	Name     string
	BookURLs []string
}

type Series struct {
	Name     string
	BookURLs []string
}

type Book struct {
	Content  io.ReadCloser
	FileName string
}
