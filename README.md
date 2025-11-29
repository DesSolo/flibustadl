## flibustadl
Simple cli application for download books from flibusta

Simple usage:
```bash
flibustadl http://flibusta.is/a/1206
```

Features:
- Download author (a)
- Download sequence (sequence)
- Download series (s)
- Unzip files
- Await before downloading book
- Retries

Supported options:
```bash
  -await duration
        Await interval (default 1s)
  -books string
        Books file path (default "books/")
  -loglevel int
        Log level (default -4)
  -remove
        Remove zip files (default true)
  -unzip
        Unzip files (default true)
```