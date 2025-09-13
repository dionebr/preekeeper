# Preekeeper

**Preekeeper** é um scanner de diretórios web profissional e de alta velocidade com interface TUI (Terminal User Interface) interativa, construído em Go usando Bubble Tea. Oferece funcionalidades comparáveis ao gobuster e dirb, mas com uma interface visual moderna e recursos avançados.

## 🚀 Características Principais

- **Interface TUI Interativa**: Interface moderna com Bubble Tea framework
- **Alta Performance**: Scanning concorrente com FastHTTP
- **Paleta de Cores Personalizada**: Esquema visual profissional
- **Múltiplos Filtros**: Por tamanho, linhas, regex e códigos de status
- **Descoberta de Tecnologias**: Identificação automática de frameworks, CMS, servidores e linguagens do alvo (Wappalyzer)
- **Scanning Recursivo**: Exploração em profundidade configurável
- **Rate Limiting**: Controle de velocidade para evitar sobrecarga
- **Suporte a Proxy**: Compatible com proxies HTTP
- **Headers Customizados**: Suporte completo a headers HTTP
- **Múltiplas Extensões**: Scanning automático com extensões configuráveis

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

## 🎯 Uso Básico

```bash
# Exemplo básico
./preekeeper -u http://example.com -w wordlist.txt

# Com threads personalizadas
./preekeeper -u http://example.com -w wordlist.txt -t 50

# Com extensões específicas
./preekeeper -u http://example.com -w wordlist.txt -x .php,.html,.js

# Scanning recursivo
./preekeeper -u http://example.com -w wordlist.txt -r -d 3

# Com códigos de status específicos
./preekeeper -u http://example.com -w wordlist.txt --mc 200,301,302
```

## 🛠️ Parâmetros Completos
| `-T, --tech` | - | Detectar tecnologias do alvo (Wappalyzer) | `--tech` |
# Detectar tecnologias do alvo
./preekeeper -u http://example.com --tech

### Parâmetros Obrigatórios
| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `-u, --url` | URL alvo (obrigatório) | `-u http://example.com` |

### Parâmetros de Performance
| Flag | Padrão | Descrição | Exemplo |
|------|--------|-----------|---------|
| `-t, --threads` | 20 | Número de threads concorrentes | `-t 50` |
| `--delay` | 0 | Delay entre requests (ms) | `--delay 100` |
| `--timeout` | 10 | Timeout de request (segundos) | `--timeout 30` |
| `--retries` | 3 | Tentativas em caso de falha | `--retries 5` |
| `--rate-limit` | 0 | Limite de requests por segundo | `--rate-limit 100` |

### Parâmetros HTTP
| Flag | Padrão | Descrição | Exemplo |
|------|--------|-----------|---------|
| `-w, --wordlist` | wordlist.txt | Arquivo de wordlist | `-w /path/to/wordlist.txt` |
| `-m, --method` | GET | Método HTTP | `-m POST` |
| `-a, --user-agent` | Preekeeper/1.0 🐝 | User agent personalizado | `-a \"Custom Agent\"` |
| `-H, --headers` | - | Headers personalizados | `-H \"Authorization: Bearer token\"` |
| `--cookies` | - | Cookies para requests | `--cookies \"session=abc123\"` |
| `--proxy` | - | URL do proxy | `--proxy http://127.0.0.1:8080` |

### Parâmetros de Filtragem
| Flag | Padrão | Descrição | Exemplo |
|------|--------|-----------|---------|
| `--mc` | 200,204,301,302,307,403,401,500 | Códigos de status para exibir | `--mc 200,301,302` |
| `--fs` | - | Filtrar por tamanho de resposta | `--fs 1024,2048` |
| `--fl` | - | Filtrar por número de linhas | `--fl 10,20` |
| `--fr` | - | Filtrar por regex na resposta | `--fr \"error\\|404\"` |

### Parâmetros de Extensão e Recursão
| Flag | Padrão | Descrição | Exemplo |
|------|--------|-----------|---------|
| `-x, --extensions` | - | Extensões de arquivo | `-x .php,.html,.js,.txt` |
| `-r, --recursive` | false | Habilitar scanning recursivo | `-r` |
| `-d, --depth` | 2 | Profundidade máxima recursiva | `-d 5` |

### Parâmetros de Segurança
| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `--no-tls-validation` | Ignorar validação TLS | `--no-tls-validation` |

### Parâmetros de Output
| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `-s, --silent` | Modo silencioso | `-s` |
| `-v, --verbose` | Output verboso | `-v` |
| `-o, --output` | Arquivo de saída | `-o results.txt` |

## 🎨 Interface TUI

### Controles Interativos
- **`s`** - Iniciar scan
- **`p`** - Pausar/Retomar scan
- **`r`** - Reiniciar scan
- **`h`** - Exibir ajuda completa
- **`q`** - Sair da aplicação
- **`↑/k`** - Scroll para cima nos resultados
- **`↓/j`** - Scroll para baixo nos resultados

### Filtros de Status
- **`1`** - Mostrar apenas respostas 2xx (Success)
- **`2`** - Mostrar apenas respostas 3xx (Redirect)
- **`3`** - Mostrar apenas respostas 4xx (Client Error)
- **`4`** - Mostrar apenas respostas 5xx (Server Error)
- **`5`** - Mostrar todas as respostas

### Esquema de Cores
- **Marrom Acinzentado** (#94866d) - Respostas 2xx (Success)
- **Marrom Madeira** (#746142) - Respostas 3xx (Redirect)
- **Marrom Escuro** (#42432e) - Respostas 4xx (Client Error)
- **Preto** (#1d1f10) - Respostas 5xx (Server Error)
- **Cinza-Bege** (#b4b2a7) - Outras respostas

## 📝 Exemplos Avançados
### Descoberta de Tecnologias
```bash
./preekeeper -u http://example.com --tech
```

### Scanning com Autenticação
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  -H \"Authorization: Bearer eyJhbGciOiJIUzI1NiIs...\" \\
  --cookies \"session=abc123; csrf_token=xyz789\"
```

### Scanning através de Proxy
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  --proxy http://127.0.0.1:8080 \\
  --no-tls-validation
```

### Scanning Recursivo Profundo
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  -r -d 5 -t 30 \\
  -x .php,.html,.js,.css,.txt,.xml
```

### Scanning com Rate Limiting
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  --rate-limit 50 \\
  --delay 200 \\
  -t 10
```

### Scanning com Filtros Específicos
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  --mc 200,301,403 \\
  --fs 1024,2048,4096 \\
  --fr \"(?i)(admin|login|password)\"
```

### Scanning com FUZZ personalizado
```bash
./preekeeper -u http://example.com/api/FUZZ -w api-endpoints.txt \\
  -m POST \\
  -H \"Content-Type: application/json\" \\
  --mc 200,201,202
```

## 📊 Métricas em Tempo Real

A interface TUI exibe métricas em tempo real:

- **Elapsed**: Tempo decorrido do scan
- **Found**: Número de endpoints encontrados
- **RPS**: Requests por segundo atual
- **Processed**: Total de requests processados
- **Current**: URL sendo testada atualmente
- **Recursion**: Status de scanning recursivo

## 🔧 Configuração de Wordlists

O Preekeeper inclui uma wordlist padrão com mais de 100 termos comuns, mas você pode usar suas próprias wordlists:

```bash
# Wordlists populares
./preekeeper -u http://example.com -w /usr/share/wordlists/dirb/common.txt
./preekeeper -u http://example.com -w /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt
./preekeeper -u http://example.com -w custom-endpoints.txt
```

## ⚡ Otimizações de Performance

### Para Máxima Velocidade
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  -t 100 \\
  --rate-limit 0 \\
  --delay 0 \\
  --timeout 5
```

### Para Targets Sensíveis
```bash
./preekeeper -u http://example.com -w wordlist.txt \\
  -t 5 \\
  --rate-limit 10 \\
  --delay 500 \\
  --timeout 30
```

## 🚨 Limitações e Considerações

- **Rate Limiting**: Use sempre rate limiting em ambientes de produção
- **Thread Count**: Mais de 100 threads pode causar problemas de rede
- **Memory Usage**: Wordlists muito grandes podem consumir memória
- **Target Stability**: Monitore a estabilidade do target durante scans intensivos

## 🐛 Troubleshooting

### Problemas Comuns

1. **\"URL is required\"**
   - Solução: Use a flag `-u` com uma URL válida

2. **\"Wordlist file not found\"**
   - Solução: Verifique o caminho da wordlist com `-w`

3. **Connection timeouts**
   - Solução: Aumente o timeout com `--timeout 30`

4. **Rate limiting pelo target**
   - Solução: Reduza threads `-t 5` e adicione delay `--delay 1000`

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📜 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para detalhes.

## 👤 Autor

**Dione Lima - Brazil**

- Ferramenta: Preekeeper Scanner 🐝
- Framework: Bubble Tea + FastHTTP
- Performance: Scanning concorrente de alta velocidade

---

**Preekeeper Scanner** - Transformando scanning de diretórios em uma experiência visual e eficiente! 🚀