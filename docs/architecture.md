# Architecture and Project Layout

This document explains the code layout and main components.

## Layout

```
/ (repo root)
  - main.go                 # application entry and TUI
  - go.mod                  # module and dependencies
  - README.md               # project readme (user-facing)
  - docs/                   # detailed docs
  - internal/               # internal helpers (proxy, techdetector)
      - fasthttpproxy.go    # local proxy dialer used by fasthttp
      - techdetector/       # wrapper hiding external detector
          - techdetector.go
  - dist/                   # built binaries
  - wordlist.txt            # default example wordlist
```

## Main components

- `TUI` (Bubble Tea) — handles user interface and input.
- `Scanner` — worker pool using fasthttp for fast HTTP requests.
- `RateLimiter` — simple token-based limiter for RPS control.
- `Proxy` — internal helper to support HTTP proxy for fasthttp.
- `Tech Detector` — hidden engine wrapper that provides technology fingerprints.

## Design notes
- Detection engine is abstracted behind `FingerprintEngine` to allow replacement/mock testing.
- The project uses `internal/` to keep non-public helpers and hide direct external package names from the public API.
- Cross-platform builds are supported with simple GOOS/GOARCH builds.