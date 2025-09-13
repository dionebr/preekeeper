# Preekeeper

**Version:** v1.1.0

**Preekeeper** is a professional, high-speed web directory scanner with an interactive TUI (Terminal User Interface) built in Go using Bubble Tea. It offers features comparable to gobuster and dirb, but with a modern visual interface and advanced capabilities.

## 📝 Changelog

### v1.1.0
- Feature: Technology detection (WappalyzerGo integration, CLI flag --tech)
- Improvement: Encapsulated tech detection in interface for testability
- Improvement: Robust error handling and user logs for tech detection
- Docs: README fully translated to English
- Refactor: Codebase ready for future releases and testing

## 🚀 Main Features

- **Interactive TUI Interface**: Modern interface with Bubble Tea framework
- **High Performance**: Concurrent scanning with FastHTTP
- **Custom Color Palette**: Professional visual scheme
- **Multiple Filters**: By size, lines, regex, and status codes
- **Technology Detection**: Automatic identification of frameworks, CMS, servers, and languages (Wappalyzer)
- **Recursive Scanning**: Configurable deep exploration
- **Rate Limiting**: Speed control to avoid overload
- **Proxy Support**: Compatible with HTTP proxies
- **Custom Headers**: Full support for HTTP headers
- **Multiple Extensions**: Automatic scanning with configurable extensions

## 📦 Installation

```bash
# Clone the repository
git clone <repository-url>
cd preekeeper-scanner

# Install dependencies
go mod tidy
# Requirement for technology detection:
go get github.com/projectdiscovery/wappalyzergo

# Build the project
go build -o preekeeper main.go
```

## 🎯 Basic Usage

```bash
# Basic example
./preekeeper -u http://example.com -w wordlist.txt

# Custom threads
./preekeeper -u http://example.com -w wordlist.txt -t 50

# Specific extensions
./preekeeper -u http://example.com -w wordlist.txt -x .php,.html,.js

# Recursive scanning
./preekeeper -u http://example.com -w wordlist.txt -r -d 3

# Specific status codes
./preekeeper -u http://example.com -w wordlist.txt --mc 200,301,302
```

## 🛠️ Full Parameters
| `-T, --tech` | - | Detect target technologies (Wappalyzer) | `--tech` |
# Detect target technologies
./preekeeper -u http://example.com --tech

### Required Parameters
| Flag | Description | Example |
|------|-------------|---------|
| `-u, --url` | Target URL (required) | `-u http://example.com` |

### Performance Parameters
| Flag | Default | Description | Example |
|------|---------|-------------|---------|
| `-t, --threads` | 20 | Number of concurrent threads | `-t 50` |
| `--delay` | 0 | Delay between requests (ms) | `--delay 100` |
| `--timeout` | 10 | Request timeout (seconds) | `--timeout 30` |
| `--retries` | 3 | Number of retries on failure | `--retries 5` |
| `--rate-limit` | 0 | Requests per second limit | `--rate-limit 100` |

### HTTP Parameters
| Flag | Default | Description | Example |
|------|---------|-------------|---------|
| `-w, --wordlist` | wordlist.txt | Wordlist file | `-w /path/to/wordlist.txt` |
| `-m, --method` | GET | HTTP method | `-m POST` |
| `-a, --user-agent` | Preekeeper/1.0 🐝 | Custom user agent | `-a "Custom Agent"` |
| `-H, --headers` | - | Custom headers | `-H "Authorization: Bearer token"` |
| `--cookies` | - | Cookies for requests | `--cookies "session=abc123"` |
| `--proxy` | - | Proxy URL | `--proxy http://127.0.0.1:8080` |

### Filtering Parameters
| Flag | Default | Description | Example |
|------|---------|-------------|---------|
| `--mc` | 200,204,301,302,307,403,401,500 | Status codes to display | `--mc 200,301,302` |
| `--fs` | - | Filter by response size | `--fs 1024,2048` |
| `--fl` | - | Filter by number of lines | `--fl 10,20` |
| `--fr` | - | Filter by regex in response | `--fr "error|404"` |

### Extension and Recursion Parameters
| Flag | Default | Description | Example |
|------|---------|-------------|---------|
| `-x, --extensions` | - | File extensions | `-x .php,.html,.js,.txt` |
| `-r, --recursive` | false | Enable recursive scanning | `-r` |
| `-d, --depth` | 2 | Maximum recursion depth | `-d 5` |

### Security Parameters
| Flag | Description | Example |
|------|-------------|---------|
| `--no-tls-validation` | Skip TLS validation | `--no-tls-validation` |

### Output Parameters
| Flag | Description | Example |
|------|-------------|---------|
| `-s, --silent` | Silent mode | `-s` |
| `-v, --verbose` | Verbose output | `-v` |
| `-o, --output` | Output file | `-o results.txt` |

## 🎨 TUI Interface

### Interactive Controls
- **`s`** - Start scan
- **`p`** - Pause/Resume scan
- **`r`** - Restart scan
- **`h`** - Show full help
- **`q`** - Quit application
- **`↑/k`** - Scroll up in results
- **`↓/j`** - Scroll down in results

### Status Filters
- **`1`** - Show only 2xx responses (Success)
- **`2`** - Show only 3xx responses (Redirect)
- **`3`** - Show only 4xx responses (Client Error)
- **`4`** - Show only 5xx responses (Server Error)
- **`5`** - Show all responses

### Color Scheme
- **Gray Brown** (#94866d) - 2xx (Success)
- **Wood Brown** (#746142) - 3xx (Redirect)
- **Dark Brown** (#42432e) - 4xx (Client Error)
- **Black Brown** (#1d1f10) - 5xx (Server Error)
- **Beige Gray** (#b4b2a7) - Other responses

## 📝 Advanced Examples
### Technology Detection
```bash
./preekeeper -u http://example.com --tech
```

### Scanning with Authentication
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  --cookies "session=abc123; csrf_token=xyz789"
```

### Scanning via Proxy
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  --proxy http://127.0.0.1:8080 \
  --no-tls-validation
```

### Deep Recursive Scanning
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  -r -d 5 -t 30 \
  -x .php,.html,.js,.css,.txt,.xml
```

### Scanning with Rate Limiting
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  --rate-limit 50 \
  --delay 200 \
  -t 10
```

### Scanning with Specific Filters
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  --mc 200,301,403 \
  --fs 1024,2048,4096 \
  --fr "(?i)(admin|login|password)"
```

### Scanning with Custom FUZZ
```bash
./preekeeper -u http://example.com/api/FUZZ -w api-endpoints.txt \
  -m POST \
  -H "Content-Type: application/json" \
  --mc 200,201,202
```

## 📊 Real-Time Metrics

The TUI interface displays real-time metrics:

- **Elapsed**: Scan elapsed time
- **Found**: Number of endpoints found
- **RPS**: Current requests per second
- **Processed**: Total requests processed
- **Current**: URL currently being tested
- **Recursion**: Recursive scanning status

## 🔧 Wordlist Configuration

Preekeeper includes a default wordlist with over 100 common terms, but you can use your own wordlists:

```bash
# Popular wordlists
./preekeeper -u http://example.com -w /usr/share/wordlists/dirb/common.txt
./preekeeper -u http://example.com -w /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt
./preekeeper -u http://example.com -w custom-endpoints.txt
```

## ⚡ Performance Optimizations

### For Maximum Speed
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  -t 100 \
  --rate-limit 0 \
  --delay 0 \
  --timeout 5
```

### For Sensitive Targets
```bash
./preekeeper -u http://example.com -w wordlist.txt \
  -t 5 \
  --rate-limit 10 \
  --delay 500 \
  --timeout 30
```

## 🚨 Limitations and Considerations

- **Rate Limiting**: Always use rate limiting in production environments
- **Thread Count**: More than 100 threads may cause network issues
- **Memory Usage**: Very large wordlists may consume memory
- **Target Stability**: Monitor target stability during intensive scans

## 🐛 Troubleshooting

### Common Issues

1. **"URL is required"**
   - Solution: Use the `-u` flag with a valid URL

2. **"Wordlist file not found"**
   - Solution: Check the wordlist path with `-w`

3. **Connection timeouts**
   - Solution: Increase timeout with `--timeout 30`

4. **Rate limiting by target**
   - Solution: Reduce threads `-t 5` and add delay `--delay 1000`

## 🤝 Contributing

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## 📜 License

This project is under the MIT license. See the `LICENSE` file for details.

## 👤 Author

**Dione Lima - Brazil**

- Tool: Preekeeper Scanner 🐝
- Framework: Bubble Tea + FastHTTP
- Performance: High-speed concurrent scanning

---

**Preekeeper Scanner** - Making directory scanning a visual and efficient experience! 🚀