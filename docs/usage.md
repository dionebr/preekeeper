# Usage and CLI Reference

This file documents all available flags and examples.

## Flags (short)

- `-u, --url` (required) — target URL
- `-w, --wordlist` — wordlist path (default: wordlist.txt)
- `-t, --threads` — concurrent threads (default 20)
- `-T, --tech` — detect target technologies
- `-r, --recursive` — enable recursion
- `-d, --depth` — recursion depth
- `-m, --method` — HTTP method
- `-H, --headers` — custom headers
- `--proxy` — proxy URL
- `--timeout` — request timeout
- `--rate-limit` — requests per second
- `-s, --silent` — silent mode
- `-v, --verbose` — verbose
- `-o, --output` — output file

## Examples

Basic:
```bash
./preekeeper -u http://example.com -w wordlist.txt
```

Tech detection
```bash
./preekeeper -u http://example.com -w wordlist.txt -T
```

Recursive scan
```bash
./preekeeper -u http://example.com -w wordlist.txt -r -d 3
```

Proxy use
```bash
./preekeeper -u http://example.com -w wordlist.txt --proxy http://127.0.0.1:8080
```

Advanced filters
```bash
./preekeeper -u http://example.com -w wordlist.txt --mc 200,301,302 --fs 1024
```