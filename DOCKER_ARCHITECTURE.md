# DOCKER_ARCHITECTURE.md

# Arquitetura de Containers Docker

Este documento descreve a arquitetura de containerização dos microsserviços.

## 📋 Sumário

1. [Estrutura de Containers](#estrutura-de-containers)
2. [Imagens Docker](#imagens-docker)
3. [Volume Persistence](#volume-persistence)
4. [Networking](#networking)
5. [Health Checks](#health-checks)
6. [Resource Limits](#resource-limits)
7. [Logging](#logging)

## 🐳 Estrutura de Containers

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Network                   │
│                   (microservices-net)                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ ms-order     │  │ ms-payment   │  │ ms-shipping  │     │
│  │ (Port 3000)  │  │ (Port 3001)  │  │ (Port 3002)  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│        │                  │                  │             │
│        └──────────────────┼──────────────────┘             │
│                           │                                │
│                    ┌──────────────┐                        │
│                    │   ms-mysql   │                        │
│                    │ (Port 3308)  │                        │
│                    └──────────────┘                        │
│                           │                                │
│                    ┌──────────────┐                        │
│                    │ mysql_data   │                        │
│                    │  (Volume)    │                        │
│                    └──────────────┘                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 📦 Imagens Docker

### 1. Order Service

**Dockerfile**: `microservices/order/Dockerfile`

```dockerfile
# Multi-stage build para reduzir tamanho
FROM golang:1.24 AS builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/microservices/order
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/microservices/order/order .
EXPOSE 3000
CMD ["./order"]
```

**Características**:
- ✅ Multi-stage build para otimização
- ✅ Imagem base Alpine (3.5MB)
- ✅ Go compilado como binário estático
- ✅ Sem dependências de runtime desnecessárias

**Tamanho estimado**: ~30 MB

### 2. Payment Service

**Dockerfile**: `microservices/payment/Dockerfile`

```dockerfile
FROM golang:1.24 AS builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/microservices/payment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/microservices/payment/payment .
EXPOSE 3001
CMD ["./payment"]
```

**Tamanho estimado**: ~30 MB

### 3. Shipping Service

**Dockerfile**: `microservices/shipping/Dockerfile`

```dockerfile
FROM golang:1.24 AS builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/microservices/shipping
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipping ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/microservices/shipping/shipping .
EXPOSE 3002
CMD ["./shipping"]
```

**Tamanho estimado**: ~30 MB

### 4. MySQL

**Imagem**: `mysql:8.0` (oficial)

- ✅ Imagem oficial da Docker Library
- ✅ Suporte completo a MySQL 8.0
- ✅ Inicialização automática de scripts SQL

**Tamanho**: ~445 MB

## 💾 Volume Persistence

### MySQL Data Volume

```yaml
volumes:
  mysql_data:
    driver: local
```

**Função**: Persistir dados do banco de dados entre restarts

**Localização padrão**:
```
/var/lib/docker/volumes/grpc-parte-3-main_mysql_data/_data
```

**Backup do volume**:
```bash
# Criar backup
docker run --rm \
  -v mysql_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/mysql_backup.tar.gz /data

# Restaurar backup
docker run --rm \
  -v mysql_data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/mysql_backup.tar.gz -C /
```

## 🌐 Networking

### Rede Docker Compose

```yaml
networks:
  microservices-net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/16
```

### Resolução de Nomes (DNS)

Dentro da rede Docker, os serviços podem se comunicar usando o nome do container:

```
order   → payment:3001  ✅
order   → shipping:3002 ✅
mysql   → localhost:3306 (interno)
```

### Portas Expostas

| Serviço | Porta Interna | Porta Host | Protocolo |
|---------|---------------|-----------|-----------|
| Order | 3000 | 3000 | gRPC |
| Payment | 3001 | 3001 | gRPC |
| Shipping | 3002 | 3002 | gRPC |
| MySQL | 3306 | 3308 | TCP |

## 🏥 Health Checks

### Order Service

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

Estados:
- ✅ **healthy**: Serviço respondendo normalmente
- ⚠️ **starting**: Aguardando inicialização
- ❌ **unhealthy**: Serviço fora

### Payment Service

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:3001/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### Shipping Service

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:3002/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### MySQL

```yaml
healthcheck:
  test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
  interval: 10s
  timeout: 5s
  retries: 5
```

## 📊 Resource Limits

### Configuração padrão (docker-compose.yml)

Sem limite de recursos (desenvolvimento).

### Configuração produção (docker-compose.prod.yml)

```yaml
services:
  mysql:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G

  order:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

**Entender Limites**:
- `limits.cpus`: Máximo de CPU que pode usar
- `limits.memory`: Máximo de memória que pode usar
- `reservations.cpus`: CPU garantida
- `reservations.memory`: Memória garantida

## 📝 Logging

### Configuração de Logs

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

**Opções**:
- `max-size`: Tamanho máximo por arquivo de log
- `max-file`: Número máximo de arquivos de log

### Acessar Logs

```bash
# Ver últimas linhas
docker-compose logs --tail=50 order

# Seguir em tempo real
docker-compose logs -f order

# Com timestamp
docker-compose logs --timestamps order
```

### Localização dos Logs

```
/var/lib/docker/containers/<container-id>/<container-id>-json.log
```

### Centralizar Logs

Para produção, considere usar:
- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Splunk**
- **DataDog**
- **New Relic**

## 🔄 Container Lifecycle

### Startup

```
docker-compose up
  ↓
[+] Running 4/4 (containers iniciando)
  ↓
[MySQL] healthcheck: waiting
  ↓
[MySQL] healthy → ready
  ↓
[Payment] inicia
  ↓
[Order] inicia (depende de Payment saudável)
  ↓
[Shipping] inicia
  ↓
Todos rodando ✅
```

### Restart Policy

```yaml
restart: always  # Reinicia sempre que cai
```

Opções:
- `no`: Não reinicia
- `always`: Sempre reinicia
- `on-failure`: Apenas em falha
- `unless-stopped`: Até ser parado manualmente

## 🔧 Troubleshooting Docker

### Container não inicia

```bash
# Ver erro completo
docker-compose logs mysql

# Reiniciar do zero
docker-compose down -v
docker-compose up
```

### Alto uso de memória

```bash
# Ver uso
docker stats

# Aumentar limite
# Editar docker-compose.yml e aumentar memory limits
docker-compose restart
```

### Volume corrompido

```bash
# Remover e recriar
docker-compose down -v
docker-compose up -d mysql
```

## 📈 Otimizações

### Tamanho de Imagem

- ✅ Multi-stage builds: Reduz de 800MB para 30MB
- ✅ Alpine Linux: Imagem base pequena (3.5MB)
- ✅ Go estático: Sem dependências de runtime

### Performance

- ✅ Caching de layers Docker
- ✅ Connection pooling MySQL
- ✅ gRPC com HTTP/2

### Segurança

- ✅ Imagens Alpine (menor superfície de ataque)
- ✅ Binários estáticos (sem vulnerabilidades de runtime)
- ✅ Network isolation

## 📚 Referências

- [Docker Documentation](https://docs.docker.com/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Compose Specification](https://github.com/compose-spec/compose-spec)
- [Alpine Linux](https://alpinelinux.org/)
- [Go Docker](https://golang.org/doc/tutorial/database_access)

---

**Última Atualização**: Janeiro de 2026  
**Versão**: 1.0
