# Getting Started

## Requirements

- Go 1.20+ (Go 1.25 was used during development)
- Git

## Install

```bash
git clone https://github.com/dionebr/preekeeper.git
cd preekeeper
go mod tidy
```

## Build

```bash
# Build locally
go build -o preekeeper main.go

# Build cross-platform (example)
GOOS=linux GOARCH=amd64 go build -o dist/preekeeper-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/preekeeper-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o dist/preekeeper-darwin-amd64 .
```

## Quick Run

```bash
./preekeeper -u http://example.com -w wordlist.txt
```

## Notes

- Use `--tech` to enable the hidden/advanced technology detection engine.
- The project includes an internal wrapper for the detection engine to hide third-party package names from the main codebase.