package flibusta

type Sequence struct {
	Name string
	URLs []string
}

type Author struct {
	Name string
	URLs []string
}

type Book struct {
	Content  []byte
	FileName string
}
