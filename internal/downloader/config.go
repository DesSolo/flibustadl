package downloader

import "time"

type Config struct {
	BooksFilePath       string
	ShouldUnzipFiles    bool
	ShouldRemoveZipFile bool
	AwaitInterval       time.Duration
}
