# Developer Guide

## Running locally

1. Install Go 1.20+
2. Clone and fetch modules

```bash
git clone https://github.com/dionebr/preekeeper.git
cd preekeeper
go mod tidy
```

3. Run

```bash
go run . --help
```

## Testing

Add unit tests under `*_test.go`. Example suggestions:

- `techdetector` adapter unit tests (mock wrapper)
- worker logic smoke tests with `httptest` server

## Extending

- To add new detection rules, extend `internal/techdetector` or replace it with a different implementation.
- To change HTTP client behavior, update `NewFastHTTPClient` or `internal/fasthttpproxy.go`.

## Release process

1. Update `CHANGELOG.md` and version tag.
2. Create a git tag (annotated) for the version.
3. Build artifacts in `dist/`.
4. Create GitHub Release and attach artifacts.
