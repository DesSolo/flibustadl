package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"path"
	"strconv"

	"flibustadl/internal/pkg/flibusta"
)

var (
	logLevel      = -4
	booksFilePath = "books/"
)

func main() {
	flag.IntVar(&logLevel, "loglevel", logLevel, "Log level")
	flag.StringVar(&booksFilePath, "books", booksFilePath, "Books file path")

	flag.Parse()

	uri, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to parse url: %v", err)
	}

	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.Level(logLevel),
		}),
	))

	ctx := context.Background()

	id, err := strconv.ParseUint(path.Base(uri.Path), 10, 64)
	if err != nil {
		log.Fatalf("failed to parse url: %v", err)
	}

	flClient := flibusta.NewClient(uri.Scheme + "://" + uri.Host)

	switch path.Dir(uri.Path) {
	case "/a":
		if err := downloadAuthor(ctx, flClient, id); err != nil {
			log.Fatalf("failed to download author: %v", err)
		}
	case "/sequence":
		if err := downloadSequence(ctx, flClient, id); err != nil {
			log.Fatalf("download sequence failed: %v", err)
		}

	default:
		log.Fatalf("unsupported path: %v", path.Dir(uri.Path))
	}
}

func downloadAuthor(ctx context.Context, cl *flibusta.Client, authorID uint64) error {
	author, err := cl.Author(ctx, authorID)
	if err != nil {
		return fmt.Errorf("failed to get author: %w", err)
	}

	slog.InfoContext(ctx, "processing author", "name", author.Name, "books", len(author.BookURLs))

	return downloadBooks(ctx, cl, author.Name, author.BookURLs)
}

func downloadSequence(ctx context.Context, client *flibusta.Client, sequenceID uint64) error {
	sequence, err := client.Sequence(ctx, sequenceID)
	if err != nil {
		return fmt.Errorf("client.Sequence: %w", err)
	}

	slog.InfoContext(ctx, "processing sequence", "name", sequence.Name, "books", len(sequence.BookURLs))

	return downloadBooks(ctx, client, sequence.Name, sequence.BookURLs)
}

func downloadBooks(ctx context.Context, client *flibusta.Client, root string, URLs []string) error {
	targetsPath := path.Join(booksFilePath, root)

	if err := os.MkdirAll(targetsPath, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	for _, uri := range URLs {
		slog.InfoContext(ctx, "downloading book", "uri", uri, "to", targetsPath)

		book, err := client.Download(ctx, uri)
		if err != nil {
			return fmt.Errorf("client.Download: %w", err)
		}

		file, err := os.Create(path.Join(targetsPath, book.FileName))
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}

		if _, err := file.Write(book.Content); err != nil {
			return fmt.Errorf("file.Write: %w", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("file.Close: %w", err)
		}
	}

	return nil
}
