# Detecção de Wildcard DNS

Preekeeper inclui uma detecção simples de wildcard DNS para evitar falsos positivos ao fuzzar subdomínios.

## Como funciona (resumo técnico)
1. Ao habilitar `--wildcard-detect` (padrão), o scanner resolve um subdomínio aleatório (por exemplo `zxy-<timestamp>.<host>`).
2. Se esse subdomínio resolve para um ou mais IPs, consideramos que o host tem wildcard DNS configurado e guardamos esses IPs em cache.
3. Para cada candidato de subdomínio real (por exemplo `admin.example.com`) resolvemos o nome e comparamos os IPs retornados com o conjunto cacheado.
4. Se houver intersecção, pulamos o candidato como provável wildcard (falso positivo).

## Limitações
- CDNs / Anycast podem tornar a comparação de IPs pouco confiável.
- DNS baseado em geolocalização pode retornar IPs diferentes dependendo de onde o processo roda.
- Alguns serviços respondem com diferent IPs para subdomínios inexistentes — pode gerar falso negativo/positivo.

## Recomendações
- Para alvos complexos, desative `--wildcard-detect` e use validação manual / requisições adicionais.
- Podemos melhorar detection: usar múltiplas queries, comparar padrões de TTL, ou usar probes HTTP com payloads específicos.

