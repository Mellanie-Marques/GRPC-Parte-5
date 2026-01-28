# DOCKERFILE_BEST_PRACTICES.md

# Dockerfile Best Practices

Documentação sobre as melhores práticas utilizadas nos Dockerfiles deste projeto.

## 📋 Índice

1. [Princípios Gerais](#princípios-gerais)
2. [Multi-Stage Builds](#multi-stage-builds)
3. [Seleção de Imagem Base](#seleção-de-imagem-base)
4. [Otimizações](#otimizações)
5. [Segurança](#segurança)
6. [Boas Práticas](#boas-práticas)

## 🎯 Princípios Gerais

### 1. Tamanho Mínimo

```dockerfile
❌ ERRADO (800MB+)
FROM golang:1.24
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/main.go
CMD ["./app"]

✅ CORRETO (30MB)
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/app .
CMD ["./app"]
```

### 2. Segurança

```dockerfile
❌ ERRADO (usuário root)
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y golang
CMD ["./app"]

✅ CORRETO (sem privilégios)
FROM alpine:latest
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser
USER appuser
CMD ["./app"]
```

### 3. Performance (Cache)

```dockerfile
❌ ERRADO (invalida cache frequentemente)
COPY . .
RUN go mod download
RUN go build -o app ./cmd/main.go

✅ CORRETO (aproveita cache)
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/main.go
```

## 🔄 Multi-Stage Builds

### Por que usar?

1. **Separação de Preocupações**: Build separado do runtime
2. **Redução de Tamanho**: Não inclui ferramentas de build na imagem final
3. **Melhor Performance**: Cache de layers mais eficiente

### Estrutura

```dockerfile
# ====== ESTÁGIO 1: BUILD ======
FROM golang:1.24 AS builder
# Instruções para compilar

# ====== ESTÁGIO 2: RUNTIME ======
FROM alpine:latest
# Apenas o necessário para executar
```

### Exemplo Order Service

```dockerfile
# Estágio 1: Build
FROM golang:1.24 AS builder

# Definir working directory no container
WORKDIR /usr/src/app

# Copiar arquivos do projeto
COPY . .

# Navegar para o diretório do serviço
WORKDIR /usr/src/app/microservices/order

# Compilar Go estaticamente
# CGO_ENABLED=0: Sem dependências C
# GOOS=linux: Compilar para Linux
# -a: Force rebuild de todos os packages
# -installsuffix cgo: Sufixo customizado para evitar conflitos
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -o order ./cmd/main.go

# Estágio 2: Runtime
FROM alpine:latest

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

# Working directory
WORKDIR /root/

# Copiar binário do estágio builder
COPY --from=builder /usr/src/app/microservices/order/order .

# Expor porta
EXPOSE 3000

# Comando padrão
CMD ["./order"]
```

### Comparação de Tamanho

| Abordagem | Tamanho | Observação |
|-----------|---------|-----------|
| Sem multi-stage | 800MB+ | Inclui compilador Go, git, etc |
| Com multi-stage | ~30MB | Apenas binário + Alpine |
| Alpine vs Ubuntu | 3.5MB vs 77MB | Alpine é 22x menor |

## 📦 Seleção de Imagem Base

### Alpine Linux (RECOMENDADO)

```dockerfile
FROM alpine:latest

# Vantagens:
# ✅ Pequeno (3.5MB)
# ✅ Seguro (menos superfície de ataque)
# ✅ Rápido (startup rápido)
# ✅ Muitos pacotes disponíveis

# Desvantagens:
# ❌ Usa musl libc (pode ter incompatibilidades com glibc)
# ❌ Menos comum que Ubuntu/Debian
```

### Ubuntu/Debian

```dockerfile
FROM ubuntu:20.04

# Vantagens:
# ✅ Familiaridade
# ✅ Muitos pacotes
# ✅ Compatibilidade glibc
# ✅ Comunidade grande

# Desvantagens:
# ❌ Muito maior (77MB+)
# ❌ Mais superfície de ataque
# ❌ Startup mais lento
```

### Distroless (AVANÇADO)

```dockerfile
FROM gcr.io/distroless/base:latest

# Vantagens:
# ✅ Mínimo absoluto
# ✅ Máxima segurança
# ✅ Pequeno tamanho

# Desvantagens:
# ❌ Sem shell (difícil debugar)
# ❌ Sem ferramentas padrão
# ❌ Curva de aprendizado
```

### Recomendação

Para este projeto: **Alpine + Go estático**
- ✅ Tamanho: ~30MB
- ✅ Segurança: Excelente
- ✅ Performance: Excelente
- ✅ Compatibilidade: Boa (Go estático não depende de glibc)

## ⚡ Otimizações

### 1. Build Cache

```dockerfile
# Ordem importa para cache!

# ❌ Inválida cache frequentemente
COPY . .
RUN go mod download
RUN go build -o app ./cmd/main.go

# ✅ Aproveita cache
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/main.go
```

### 2. Consolidar RUN Commands

```dockerfile
# ❌ Múltiplas camadas
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl
RUN apk --no-cache add openssl

# ✅ Uma única camada
RUN apk --no-cache add ca-certificates curl openssl
```

### 3. Usar .dockerignore

```
# .dockerignore
node_modules/
npm-debug.log
.git
.gitignore
README.md
.env
__pycache__
venv/
.vscode/
.idea/
*.log
```

### 4. Build Args

```dockerfile
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION} AS builder

# Uso: docker build --build-arg GO_VERSION=1.25
```

## 🔐 Segurança

### 1. Não Rodar como Root

```dockerfile
❌ ERRADO (default é root)
FROM alpine:latest
CMD ["./app"]

✅ CORRETO (usuário dedicado)
FROM alpine:latest
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser
USER appuser
CMD ["./app"]
```

### 2. Não Usar Latest para Versões Base

```dockerfile
❌ ERRADO (pode mudar)
FROM alpine:latest

✅ CORRETO (versão específica)
FROM alpine:3.19
```

### 3. Certificados CA

```dockerfile
# Essencial para HTTPS
RUN apk --no-cache add ca-certificates
```

### 4. Read-Only Root Filesystem

```dockerfile
# Em produção, considere:
# docker run --read-only myimage
```

### 5. Scan de Vulnerabilidades

```bash
# Usar Trivy
trivy image myimage:latest

# Ou Docker Scout
docker scout cves myimage:latest
```

## ✅ Boas Práticas

### 1. Labels para Metadados

```dockerfile
LABEL maintainer="seu-email@example.com"
LABEL version="1.0"
LABEL description="Order Service - gRPC"
```

### 2. Health Check

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1
```

### 3. Documentar Portas

```dockerfile
EXPOSE 3000/tcp
```

### 4. .dockerignore

```dockerfile
# Sempre criar um .dockerignore para não enviar arquivos desnecessários
echo "node_modules\n.git\n*.log" > .dockerignore
```

### 5. Usar ENTRYPOINT + CMD

```dockerfile
# Mais flexível que apenas CMD
ENTRYPOINT ["/app"]
CMD ["serve"]

# docker run myimage              → /app serve
# docker run myimage --help       → /app --help
```

### 6. Limpar Cache APK

```dockerfile
RUN apk --no-cache add package
# O --no-cache evita armazenar índices de pacotes
```

## 🎓 Exemplo Completo

```dockerfile
# ====== ESTÁGIO 1: BUILD ======
FROM golang:1.24 AS builder

# Metadados de build
LABEL stage=builder

# Working directory
WORKDIR /usr/src/app

# Copiar módulos Go primeiro (aproveita cache)
COPY microservices/order/go.* ./microservices/order/

# Download de dependências
WORKDIR /usr/src/app/microservices/order
RUN go mod download

# Copiar código
COPY . /usr/src/app/
WORKDIR /usr/src/app/microservices/order

# Compilar com flags otimizados
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o order ./cmd/main.go

# ====== ESTÁGIO 2: RUNTIME ======
FROM alpine:3.19

# Metadados
LABEL maintainer="seu-email@example.com"
LABEL version="1.0"
LABEL description="Order Service"

# Instalar apenas o necessário
RUN apk --no-cache add ca-certificates

# Criar usuário não-root
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Working directory
WORKDIR /home/appuser

# Copiar binário do builder
COPY --from=builder /usr/src/app/microservices/order/order /usr/local/bin/order

# Mudar propriedade
RUN chown -R appuser:appuser /home/appuser

# Mudar para usuário não-root
USER appuser

# Porta
EXPOSE 3000/tcp

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget -q -O- http://localhost:3000/health || exit 1

# Entrypoint
ENTRYPOINT ["order"]
```

## 🔗 Referências

- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
- [Multi-stage Builds](https://docs.docker.com/develop/develop-images/multistage-build/)
- [Alpine Linux](https://alpinelinux.org/)
- [Go Docker](https://golang.org/doc/tutorial/database_access)

---

**Última Atualização**: Janeiro de 2026  
**Versão**: 1.0
