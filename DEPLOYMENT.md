# Guia de Deployment - Microsserviços com gRPC

Este documento detalha os procedimentos para implantar e gerenciar os microsserviços em diferentes ambientes.

## 📋 Conteúdo

1. [Requisitos de Sistema](#requisitos-de-sistema)
2. [Arquivos Docker](#arquivos-docker)
3. [Deployment Local](#deployment-local)
4. [Deployment em Produção](#deployment-em-produção)
5. [Gerenciamento de Containers](#gerenciamento-de-containers)
6. [Monitoramento e Logs](#monitoramento-e-logs)
7. [Backup e Recuperação](#backup-e-recuperação)

## 🔧 Requisitos de Sistema

### Mínimo para Execução

| Componente | Versão | Tamanho |
|-----------|--------|--------|
| Docker | 20.10+ | N/A |
| Docker Compose | 2.0+ | N/A |
| Espaço em Disco | - | 2 GB |
| Memória RAM | - | 2 GB |
| CPU | - | 2 cores |

### Recomendado para Produção

| Componente | Versão | Tamanho |
|-----------|--------|--------|
| Docker | 24.0+ | N/A |
| Docker Compose | 2.20+ | N/A |
| Espaço em Disco | - | 10 GB |
| Memória RAM | - | 8 GB |
| CPU | - | 4 cores |

## 📦 Arquivos Docker

### Estrutura

```
.
├── docker-compose.yml           # Orquestração principal
├── microservices/
│   ├── order/Dockerfile        # Build para Order Service
│   ├── payment/Dockerfile      # Build para Payment Service
│   └── shipping/Dockerfile     # Build para Shipping Service
└── tmp_create_dbs.sql          # Inicialização do banco
```

### Dockerfile - Order Service

```dockerfile
# Estágio de build
FROM golang:1.24 AS builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/microservices/order
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go

# Estágio de execução
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/microservices/order/order .
EXPOSE 3000
CMD ["./order"]
```

### Dockerfile - Payment Service

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

### Dockerfile - Shipping Service

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

## 🚀 Deployment Local

### Passo 1: Preparar o Ambiente

```bash
# Clonar repositório
git clone <url-do-repositorio>
cd GRPC-Parte-3-main

# Verificar Docker
docker --version
docker-compose --version
```

### Passo 2: Iniciar os Serviços

```bash
# Build e execução completa
docker-compose up --build -d

# Ou apenas iniciar (se já foi buildado)
docker-compose up -d
```

### Passo 3: Verificar Status

```bash
# Ver status dos containers
docker-compose ps

# Ver logs do MySQL (aguardar inicialização)
docker-compose logs mysql | grep "ready for connections"

# Ver logs de cada serviço
docker-compose logs order
docker-compose logs payment
docker-compose logs shipping
```

### Passo 4: Testar Conectividade

```bash
# Verificar se Order está respondendo
curl localhost:3000

# Verificar se Payment está respondendo
curl localhost:3001

# Verificar se Shipping está respondendo
curl localhost:3002
```

### Passo 5: Executar Cliente de Teste

```bash
# Em outro terminal
cd microservices/order
go run client/main.go
```

### Passo 6: Parar os Serviços

```bash
# Parar sem remover volumes
docker-compose stop

# Parar e remover containers
docker-compose down

# Parar e remover tudo (incluindo volumes/dados)
docker-compose down -v
```

## 🏭 Deployment em Produção

### Variáveis de Ambiente para Produção

Crie um arquivo `.env.production`:

```env
# MySQL
MYSQL_ROOT_PASSWORD=MudePara_SenhaSegura123!
MYSQL_DATABASE_ORDER=order_db
MYSQL_DATABASE_PAYMENT=payment

# Order Service
ORDER_PORT=3000
ORDER_DB_HOST=mysql
ORDER_DB_PORT=3306
ORDER_DB_USER=root
ORDER_DB_PASSWORD=MudePara_SenhaSegura123!

# Payment Service
PAYMENT_PORT=3001
PAYMENT_DB_HOST=mysql
PAYMENT_DB_PORT=3306
PAYMENT_DB_USER=root
PAYMENT_DB_PASSWORD=MudePara_SenhaSegura123!

# Shipping Service
SHIPPING_PORT=3002

# Ambiente
ENV=production
```

### Modificar docker-compose.yml para Produção

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: ms-mysql
    restart: always  # Reinicia automaticamente
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    volumes:
      - mysql_data:/var/lib/mysql  # Persiste dados
      - ./tmp_create_dbs.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - microservices-net
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  payment:
    build:
      context: .
      dockerfile: microservices/payment/Dockerfile
    container_name: ms-payment
    restart: always
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      DB_DRIVER: mysql
      DATA_SOURCE_URL: root:${MYSQL_ROOT_PASSWORD}@tcp(mysql:3306)/payment
      APPLICATION_PORT: ${PAYMENT_PORT}
      ENV: ${ENV}
    networks:
      - microservices-net
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  order:
    build:
      context: .
      dockerfile: microservices/order/Dockerfile
    container_name: ms-order
    restart: always
    depends_on:
      mysql:
        condition: service_healthy
      payment:
        condition: service_started
      shipping:
        condition: service_started
    environment:
      PAYMENT_SERVICE_URL: payment:${PAYMENT_PORT}
      SHIPPING_SERVICE_URL: shipping:${SHIPPING_PORT}
      DB_DRIVER: mysql
      DATA_SOURCE_URL: root:${MYSQL_ROOT_PASSWORD}@tcp(mysql:3306)/order_db
      APPLICATION_PORT: ${ORDER_PORT}
      ENV: ${ENV}
    networks:
      - microservices-net
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  shipping:
    build:
      context: .
      dockerfile: microservices/shipping/Dockerfile
    container_name: ms-shipping
    restart: always
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      APPLICATION_PORT: ${SHIPPING_PORT}
      ENV: ${ENV}
    networks:
      - microservices-net
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  mysql_data:
    driver: local

networks:
  microservices-net:
    driver: bridge
```

### Iniciar em Produção

```bash
# Carregar variáveis de ambiente
export $(cat .env.production | xargs)

# Build e execução
docker-compose -f docker-compose.yml up --build -d

# Verificar status
docker-compose ps
```

## 🐳 Gerenciamento de Containers

### Comandos Comuns

```bash
# Verificar status
docker-compose ps
docker ps -a

# Ver logs
docker-compose logs -f order
docker-compose logs -f payment
docker-compose logs -f shipping
docker-compose logs -f mysql

# Parar serviço específico
docker-compose stop order

# Reiniciar serviço específico
docker-compose restart order

# Executar comando em um container
docker-compose exec mysql mysql -u root -p

# Ver uso de recursos
docker stats

# Remover volumes (dados persistidos)
docker-compose down -v
```

### Escalabilidade

Para escalar um serviço em múltiplas instâncias:

```bash
# Escalar Order Service para 3 instâncias
docker-compose up -d --scale order=3

# Balanceamento de carga necessário (Nginx/HAProxy)
```

## 📊 Monitoramento e Logs

### Configurar Logs

Logs são salvos em:
```
/var/lib/docker/containers/<container-id>/
```

### Acessar Logs

```bash
# Últimas 100 linhas
docker-compose logs --tail=100 order

# Seguir em tempo real
docker-compose logs -f order

# Com timestamp
docker-compose logs --timestamps order

# Desde uma data específica
docker-compose logs --since 2024-01-01 order
```

### Centralizar Logs (ELK Stack)

Para produção, considere usar ELK:

```yaml
# Adicionar ao docker-compose.yml
elasticsearch:
  image: docker.elastic.co/elasticsearch/elasticsearch:7.17.0
  environment:
    - discovery.type=single-node
  ports:
    - "9200:9200"

logstash:
  image: docker.elastic.co/logstash/logstash:7.17.0
  volumes:
    - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
  ports:
    - "5000:5000"

kibana:
  image: docker.elastic.co/kibana/kibana:7.17.0
  ports:
    - "5601:5601"
```

## 💾 Backup e Recuperação

### Backup de Dados

```bash
# Fazer dump do banco MySQL
docker-compose exec mysql mysqldump -u root -pminhasenha --all-databases > backup.sql

# Fazer backup de volumes
docker run --rm \
  -v mysql_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/mysql_backup.tar.gz /data
```

### Recuperar Dados

```bash
# Restaurar banco MySQL
docker-compose exec -T mysql mysql -u root -pminhasenha < backup.sql

# Restaurar volumes
docker run --rm \
  -v mysql_data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/mysql_backup.tar.gz -C /
```

## 🔐 Segurança para Produção

### Checklist de Segurança

- [ ] Alterar senha padrão do MySQL
- [ ] Usar `.env` para variáveis sensíveis
- [ ] Implementar mTLS entre serviços
- [ ] Configurar firewall/network policies
- [ ] Usar secrets do Docker/Kubernetes
- [ ] Ativar logging e monitoramento
- [ ] Implementar backup automático
- [ ] Usar healthchecks
- [ ] Configurar rate limiting
- [ ] Usar HTTPS em produção

### Exemplo com Secrets do Docker

```bash
# Criar secrets
echo "MudePara_SenhaSegura123!" | docker secret create db_password -

# Usar em docker-compose.yml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

secrets:
  db_password:
    external: true
```

## 🔍 Troubleshooting

### Container não inicia

```bash
# Ver erro
docker-compose logs service_name

# Solução comum: Remover dados corrompidos
docker-compose down -v
docker-compose up --build
```

### Erro de conexão entre serviços

```bash
# Verificar rede
docker network ls
docker network inspect grpc-parte-3-main_microservices-net

# Testar conectividade
docker-compose exec order ping payment
docker-compose exec order ping mysql
```

### Alto uso de memória

```bash
# Verificar uso
docker stats

# Limpar cache
docker system prune -f

# Limpar volumes não usados
docker volume prune
```

### Banco de dados não inicializa

```bash
# Verificar arquivo SQL
docker-compose logs mysql | grep ERROR

# Reiniciar com volumes zerados
docker-compose down -v
docker-compose up -d mysql
```

## 📈 Performance Tuning

### MySQL

```sql
-- Aumentar pool de conexões
SET GLOBAL max_connections = 1000;

-- Otimizar queries
CREATE INDEX idx_product_code ON products(product_code);
CREATE INDEX idx_order_customer ON orders(customer_id);
```

### Docker

```yaml
services:
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

## 📚 Referências

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Docker Best Practices](https://golang.org/doc/tutorial/database_access)
- [MySQL Docker Image](https://hub.docker.com/_/mysql)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance-best-practices/)

---

**Última Atualização**: Janeiro de 2026  
**Versão**: 1.0
