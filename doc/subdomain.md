# Subdomain fuzzing (design notes)

Este documento explica o comportamento das opções de subdomain no Preekeeper e recomenda melhores práticas.

## Modo simples (`--subdomain`)
- Cada palavra da wordlist é usada como label e combinada com o host alvo.
- Exemplo: label `admin` + host `example.com` => `http(s)://admin.example.com/`.

## Cartesian mode (`--subdomain-paths`)
- Para cada label na wordlist, testamos cada path também (produto cartesiano).
- Use com cuidado: um wordlist de 10k linhas se transformará em 100M combinações.
- Recomenda-se testar com wordlists pequenas ou usar filtros de rate-limit e delay.

## Both-schemes (`--http-https`)
- Tenta preferencialmente `https` e depois `http` por label quando ativado.
- Se o alvo tem apenas um esquema, a tentativa extra aumenta latência.

## Recomendações de uso
- Combine `--subdomain` com `--wildcard-detect` (padrão) para reduzir falsos positivos.
- Use `--rate-limit` e `--delay` para controlar taxa de requests.
- Se precisar de mais controle sobre combinações, considere gerar uma wordlist específica com o produto cartesiano fora do scanner e usar em modo padrão.

