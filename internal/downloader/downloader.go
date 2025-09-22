package downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path"
	"strconv"

	"flibustadl/internal/pkg/flibusta"
)

type Downloader struct {
	client *flibusta.Client
	config *Config
}

func NewDownloader(client *flibusta.Client, config *Config) *Downloader {
	return &Downloader{
		client: client,
		config: config,
	}
}

func (d *Downloader) Download(ctx context.Context, uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("url.Parse: %w", err)
	}

	id, err := strconv.ParseUint(path.Base(u.Path), 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseUint: %w", err)
	}

	switch path.Dir(u.Path) {
	case "/a":
		if err := d.downloadAuthor(ctx, id); err != nil {
			return fmt.Errorf("downloadAuthor: %w", err)
		}
	case "/sequence":
		if err := d.downloadSequence(ctx, id); err != nil {
			return fmt.Errorf("downloadSequence: %w", err)
		}
	case "/s":
		if err := d.downloadSeries(ctx, id); err != nil {
			return fmt.Errorf("downloadSeries: %w", err)
		}
	default:
		return fmt.Errorf("unknown path: %s", u.Path)
	}

	return nil
}

func (d *Downloader) downloadAuthor(ctx context.Context, authorID uint64) error {
	author, err := d.client.Author(ctx, authorID)
	if err != nil {
		return fmt.Errorf("failed to get author: %w", err)
	}

	slog.InfoContext(ctx, "processing author", "name", author.Name, "books", len(author.BookURLs))

	return d.downloadBooks(ctx, author.Name, author.BookURLs)
}

func (d *Downloader) downloadSequence(ctx context.Context, sequenceID uint64) error {
	sequence, err := d.client.Sequence(ctx, sequenceID)
	if err != nil {
		return fmt.Errorf("client.Sequence: %w", err)
	}

	slog.InfoContext(ctx, "processing sequence", "name", sequence.Name, "books", len(sequence.BookURLs))

	return d.downloadBooks(ctx, sequence.Name, sequence.BookURLs)
}

func (d *Downloader) downloadSeries(ctx context.Context, seriesID uint64) error {
	series, err := d.client.Series(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("client.Series: %w", err)
	}

	slog.InfoContext(ctx, "processing series", "name", series.Name, "books", len(series.BookURLs))
	return d.downloadBooks(ctx, series.Name, series.BookURLs)
}

func (d *Downloader) downloadBooks(ctx context.Context, root string, URLs []string) error {
	targetsPath := path.Join(d.config.BooksFilePath, root)

	if err := os.MkdirAll(targetsPath, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	for _, uri := range URLs {
		slog.InfoContext(ctx, "downloading book", "uri", uri, "to", targetsPath)

		book, err := d.client.Download(ctx, uri)
		if err != nil {
			return fmt.Errorf("client.Download: %w", err)
		}

		bookFilePath := path.Join(targetsPath, book.FileName)

		file, err := os.Create(bookFilePath)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}

		if _, err := file.Write(book.Content); err != nil {
			return fmt.Errorf("file.Write: %w", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("file.Close: %w", err)
		}

		if !d.config.ShouldUnzipFiles {
			continue
		}

		slog.DebugContext(ctx, "unzipping book", "file", bookFilePath)
		if err := unzip(targetsPath, bookFilePath); err != nil {
			return fmt.Errorf("unzip: %w", err)
		}

		if d.config.ShouldRemoveZipFile {
			slog.DebugContext(ctx, "removing zip file", "file", bookFilePath)
			if err := os.Remove(bookFilePath); err != nil {
				return fmt.Errorf("os.Remove: %w", err)
			}
		}
	}

	return nil
}

func unzip(root, filePath string) error {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("zip.OpenReader: %w", err)
	}

	defer reader.Close()

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		zipFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("f.Open: %w", err)
		}

		fsFile, err := os.OpenFile(path.Join(root, f.Name), os.O_WRONLY|os.O_CREATE, f.Mode())
		if err != nil {
			return fmt.Errorf("os.OpenFile: %w", err)
		}

		if _, err := io.Copy(fsFile, zipFile); err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}

		if err := zipFile.Close(); err != nil {
			return fmt.Errorf("file.Close: %w", err)
		}

		if err := fsFile.Close(); err != nil {
			return fmt.Errorf("fsFile.Close: %w", err)
		}
	}

	return nil
}
