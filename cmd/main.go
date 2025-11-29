package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"flibustadl/internal/downloader"
	"flibustadl/internal/pkg/flibusta"
)

var (
	logLevel      = -4
	booksFilePath = "books/"
	unzipFiles    = true
	removeZipFile = true
	awaitInterval = time.Second
)

func main() {
	flag.IntVar(&logLevel, "loglevel", logLevel, "Log level")
	flag.StringVar(&booksFilePath, "books", booksFilePath, "Books file path")
	flag.BoolVar(&unzipFiles, "unzip", unzipFiles, "Unzip files")
	flag.BoolVar(&removeZipFile, "remove", removeZipFile, "Remove zip files")
	flag.DurationVar(&awaitInterval, "await", awaitInterval, "Await interval")

	flag.Parse()

	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.Level(logLevel),
		}),
	))

	target := flag.Arg(0)

	uri, err := url.Parse(target)
	if err != nil {
		log.Fatalf("failed to parse url: %v", err)
	}

	ctx := context.Background()

	dl := downloader.NewDownloader(flibusta.NewClient(uri.Scheme+"://"+uri.Host), &downloader.Config{
		BooksFilePath:       booksFilePath,
		ShouldUnzipFiles:    unzipFiles,
		ShouldRemoveZipFile: removeZipFile,
		AwaitInterval:       awaitInterval,
	})

	if err := dl.Download(ctx, target); err != nil {
		log.Fatalf("failed to download: %v", err)
	}
}
