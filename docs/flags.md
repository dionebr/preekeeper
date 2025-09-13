# Flags e parâmetros do Preekeeper

Este documento descreve todas as flags suportadas pela ferramenta, incluindo as novas opções de subdomain e wildcard.

## Principais flags

- `-u, --url` (required): Target URL. Exemplo: `-u http://example.com`.
- `-w, --wordlist`: Wordlist file (default `wordlist.txt`).
- `-t, --threads`: Number of concurrent threads (default 20).
- `--delay`: Delay entre requests em ms (default 0).
- `--timeout`: Request timeout em segundos (default 10).
- `--retries`: Retries on failure (default 3).
- `--rate-limit`: Requests per second (0 = unlimited).

## HTTP / Output

- `-m, --method`: HTTP method (default GET).
- `-a, --user-agent`: User-Agent header.
- `-H, --headers`: Custom headers (can be used multiple times).
- `--cookies`: Cookies string.
- `--proxy`: Proxy URL (http://host:port).
- `-s, --silent`: Silent mode (no banner).
- `-v, --verbose`: Verbose logs (diagnostics go to stderr).
- `-o, --output`: Output file for results.

## Tecnologia

- `-T, --tech`: Detect target technologies (runs silently by default). Use `-v` to see diagnostics.

## Subdomain / Advanced

- `-S, --subdomain`: Fuzz subdomains using the wordlist (each entry becomes a label).
- `--subdomain-paths`: When used with `--subdomain`, combine subdomains and paths (cartesian product). Very costly.
- `--http-https`: When used with `--subdomain`, try both `https` and `http` per label (prefers https first).
- `--wildcard-detect`: (default true) Detect wildcard DNS by resolving a random label and skip results that match the wildcard IPs.

## Filtering

- `--mc`: Match status codes (comma-separated).
- `--fs`: Filter by response size.
- `--fl`: Filter by lines count.
- `--fr`: Filter by regex in response body.


### Observações
- Flags combinadas podem gerar comportamentos custosos (ex.: `--subdomain --subdomain-paths --http-https`). Use rate limiting to control.
- `--wildcard-detect` é ativado por padrão; desative se quiser tratar qualquer host que resolve como válido.


---

Para exemplos práticos veja `docs/usage.md`.
