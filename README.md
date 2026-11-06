# PGE Stress Test

Projeto desenvolvido como parte do curso Pós Go Expert da Full Cycle. A aplicação é uma ferramenta de linha de comando escrita em Go para executar testes de carga HTTP. Ela envia uma quantidade definida de requisições para uma URL, limita o número de requisições simultâneas e apresenta métricas de duração, latência, throughput, erros de transporte e códigos de status HTTP.

O repositório também inclui um servidor HTTP simples para testes locais e uma configuração Docker Compose para executar cliente e servidor na mesma rede Docker.

## Executando localmente

```bash
go run ./cmd/stresstest/main.go \
  --url=http://localhost:8080 \
  --requests=1000 \
  --concurrency=10
```

As flags disponíveis também têm atalhos:

| Flag            | Descrição                          | Padrão      |
| --------------- | ---------------------------------- | ----------- |
| `--url`         | URL HTTP ou HTTPS a testar.        | obrigatório |
| `--requests`    | Número total de requisições.       | `10`        |
| `--concurrency` | Máximo de requisições simultâneas. | `1`         |

Exemplo contra um endpoint público:

```bash
go run ./cmd/stresstest/main.go \
  --url=https://www.google.com \
  --requests=100 \
  --concurrency=5
```

## Servidor local de exemplo

O projeto disponibiliza um servidor HTTP mínimo em `cmd/server` que responde na porta `8080`.

Em um terminal, inicie o servidor:

```bash
go run ./cmd/server/main.go
```

Em outro terminal, execute o stress test:

```bash
go run ./cmd/stresstest/main.go \
  --url=http://localhost:8080 \
  --requests=10000 \
  --concurrency=10
```

## Saída

Exemplo de resumo:

```text
════════════════════ Stress Test Summary ════════════════════
Duration:              523.07ms
Requests attempted:    13400
HTTP responses:        13400
Transport errors:      0

Throughput
  Attempted:           25618.21 req/s
  Responses:           25618.21 req/s
  Successful (2xx):    0.00 req/s

Latency
  Average:             10.70ms
  Minimum:             236.27µs
  Maximum:             109.11ms

HTTP status codes
  Status   Requests   Rate
  ------   --------   ------
  429      13400      100.00%
══════════════════════════════════════════════════════════════
```

### Métricas

| Métrica                  | Significado                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------- |
| **Duration**             | Tempo de relógio desde o início até o término de todas as requisições.                            |
| **Requests attempted**   | Total de tentativas realizadas.                                                                   |
| **HTTP responses**       | Tentativas que receberam uma resposta HTTP.                                                       |
| **Transport errors**     | Falhas antes de receber uma resposta HTTP, como conexão recusada, DNS ou timeout.                 |
| **Attempted throughput** | Tentativas por segundo.                                                                           |
| **Response throughput**  | Respostas HTTP por segundo.                                                                       |
| **Successful (2xx)**     | Respostas HTTP `2xx` por segundo.                                                                 |
| **Latency**              | Tempo individual de cada requisição, do início até o recebimento e descarte do corpo da resposta. |

> Códigos `4xx` e `5xx` são respostas HTTP, não erros de transporte. Por exemplo, uma resposta `429 Too Many Requests` aparece na tabela de status e conta em **HTTP responses**, mas não em **Transport errors**.

## Docker

### Build da imagem do stress test

```bash
docker build \
  -t pge-stress-test \
  --build-arg APP_PATH=./cmd/stresstest/main.go \
  .
```

### Executar contra uma URL pública

```bash
docker run --rm pge-stress-test \
  --url=https://www.google.com \
  --requests=100 \
  --concurrency=5
```

## Docker Compose

O arquivo `docker-compose.yaml` sobe:

| Serviço  | Papel                                                                           |
| -------- | ------------------------------------------------------------------------------- |
| `server` | Servidor HTTP de exemplo, publicado em `http://localhost:8080`.                 |
| `client` | Executa o stress test contra `http://server:8080` pela rede interna do Compose. |

Para criar as imagens e executar os serviços:

```bash
docker compose up --build
```

Na rede interna do Docker Compose, o cliente acessa o servidor usando:

```text
http://server:8080
```

```bash
docker compose down
```

## Makefile

O `Makefile` possui alguns atalhos para construir e executar a imagem:

```bash
make build
make run-google
make run-amazon
```
