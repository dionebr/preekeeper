# Preekeeper — Uso e exemplos

Este documento explica o uso prático da ferramenta e exemplos que incluem as opções mais avançadas (subdomínios, cartesian product e detecção de wildcard DNS).

## Exemplos básicos

# Scan simples
./preekeeper -u http://example.com -w wordlist.txt

# Com threads e extensões
./preekeeper -u http://example.com -w wordlist.txt -t 50 -x .php,.html

# Fuzz com placeholder FUZZ
./preekeeper -u http://example.com/api/FUZZ -w api-endpoints.txt --mc 200,201

## Subdomain fuzzing (feroxbuster-like)

# Fuzz de subdomínios (cada palavra da wordlist vira um label)
./preekeeper -u http://example.com -w wordlist.txt --subdomain

# Tentar both schemes (https + http) por label
./preekeeper -u http://example.com -w wordlist.txt --subdomain --http-https

# Cartesian product: subdomínio x paths (muito custoso)
# Para cada label da wordlist, será testado cada path da mesma wordlist
./preekeeper -u http://example.com -w wordlist.txt --subdomain --subdomain-paths

# Desativar detecção de wildcard (caso queira assumir sempre que um subdomínio que resolve é válido)
./preekeeper -u http://example.com -w wordlist.txt --subdomain --wildcard-detect=false

## Detecção de tecnologias

# Detectar tecnologias do alvo (silencioso por padrão; use -v para ver logs de diagnóstico)
./preekeeper -u http://example.com --tech

# Exibir mensagens de diagnóstico da detecção
./preekeeper -u http://example.com --tech -v

## Recomendações práticas

- O modo `--subdomain-paths` (cartesian) pode gerar N^2 requests e encher a fila de jobs — use apenas em alvos controlados ou com wordlists pequenas.
- Quando usar `--http-https` espere o dobro das tentativas por label (latência aumentada).
- Se estiver usando wordlists grandes, prefira `--rate-limit` e `--delay` para evitar sobrecarregar a rede ou o alvo.

## Exemplos combinados

# Subdomain + cartesian + both schemes + rate limit
./preekeeper -u example.com -w wordlist.txt --subdomain --subdomain-paths --http-https --rate-limit 50 --delay 50


---

Para detalhes completos de cada flag veja `doc/flags.md`. Para explicação da detecção de wildcard e limitações veja `doc/wildcard.md`. Para comportamento e recomendações sobre subdomain cartesian veja `doc/subdomain.md`.
