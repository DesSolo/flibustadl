package flibusta

type Sequence struct {
	Name     string
	BookURLs []string
}

type Author struct {
	Name     string
	BookURLs []string
}

type Book struct {
	Content  []byte
	FileName string
}
