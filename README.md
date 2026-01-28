# 🚀 Microsserviços com gRPC - Parte 5

Sistema de 3 microsserviços (Order, Payment, Shipping) conectados via gRPC e orquestrados com Docker.

## ⚡ Quick Start (5 minutos)

```bash
# Opção 1: Com Docker Compose
docker-compose up --build -d
cd microservices/order && go run client/main.go

# Opção 2: Com Makefile
make up
make test
make down

# Opção 3: Com script de deploy
./deploy.sh development
```

## 📋 Pré-requisitos

- **Docker** 20.10+ e **Docker Compose** 2.0+
- **Go** 1.24+ (apenas para desenvolvimento local)
- **MySQL** 8.0 (ou via Docker)

### Passo 3: Teste o Sistema

Em outro terminal:

```bash
cd microservices/order
go run client/main.go
```

O cliente executará 4 testes com diferentes cenários e exibirá os resultados.

### Passo 4: Parar os Serviços

```bash
docker-compose down
```

## Estrutura do Projeto

```
.
├── docker-compose.yml              # Orquestração de containers
├── tmp_create_dbs.sql             # Script de inicialização do banco

microservices/
├── order/           # Serviço Order (porta 3000)
├── payment/         # Serviço Payment (porta 3001)
└── shipping/        # Serviço Shipping (porta 3002)

microservices-proto/
├── order/           # Proto definitions
├── payment/         
└── shipping/        # Implementações gRPC geradas

docker-compose.yml   # Orquestração dev
docker-compose.prod.yml  # Orquestração produção
```

## 🔧 Variáveis de Ambiente

Copie `.env.example` para `.env`:

```bash
# Essencial
MYSQL_ROOT_PASSWORD=root123
MYSQL_DATABASE=microservices

# Serviços
ORDER_PORT=3000
PAYMENT_PORT=3001
SHIPPING_PORT=3002

# Timeouts
GRPC_TIMEOUT=2s
RETRY_ATTEMPTS=5
```

## 🛠️ Comandos Principais

| Comando | Ação |
|---------|------|
| `make up` | Inicia todos os serviços |
| `make down` | Para todos os serviços |
| `make logs` | Mostra logs |
| `make logs-follow` | Logs em tempo real |
| `make test` | Executa cliente de teste |
| `make mysql-cli` | Conecta ao MySQL |
| `make db-backup` | Faz backup do banco |
| `make db-restore` | Restaura backup |
| `make clean-all` | Remove tudo (imagens + volumes) |

## 🚀 Deployment

### Desenvolvimento (Local)
```bash
docker-compose up --build -d
```

### Produção (Com resource limits)
```bash
docker-compose -f docker-compose.prod.yml up -d
# ou
./deploy.sh production
```

### Staging
```bash
./deploy.sh staging
```

## 📊 Serviços

### Order (porta 3000)
- Criar pedido → valida com Payment → calcula entrega com Shipping
- Arquitetura hexagonal
- Repositório padrão

### Payment (porta 3001)
- Processa pagamentos
- Integrado com Order via gRPC
- Simula API de pagamento

### Shipping (porta 3002)
- Calcula prazo de entrega (fórmula simples)
- Integrado com Order via gRPC
- Recebe produto_id e retorna dias_entrega

## 🗄️ Banco de Dados

MySQL 8.0 é iniciado automaticamente com as tabelas:

```sql
-- Orders
CREATE TABLE orders (
  id UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  product_id INT NOT NULL,
  quantity INT NOT NULL,
  status VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Payments
CREATE TABLE payments (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  amount DECIMAL(10,2),
  status VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Shipments
CREATE TABLE shipments (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  estimated_days INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 📝 Desenvolvimento Local

```bash
# Build de um serviço
cd microservices/order
go build -o order ./cmd/main.go

# Ou rodar diretamente
go run ./cmd/main.go

# Testes
go test ./...
```

## 🔍 Troubleshooting

| Problema | Solução |
|----------|---------|
| Porta já em uso | Mude em `.env` ou `docker-compose.yml` |
| MySQL não conecta | Aguarde ~5s, confira senha em `.env` |
| gRPC timeout | Verifique se serviços estão up: `docker ps` |
| Imagens grandes | Use `make docker-prune` para limpar |
| Erro de build | Delete `go.mod` e `go.sum`, rode `go mod tidy` |

## 📦 Backup

```bash
# Backup automático
make db-backup

# Restaurar
make db-restore
```

Backups ficam em `backups/` com timestamp.

## ✅ Health Checks

Todos os serviços têm health checks automáticos:

```bash
# Ver status
docker-compose ps

# Ou
make status
```

## 🔐 Segurança em Produção

- Use `docker-compose.prod.yml` com resource limits
- Configure variáveis de ambiente em `.env` (não commite)
- Use secrets do Docker em produção
- Ative logging centralizado (ELK Stack opcional)

## 📚 Tecnologias

- **Go** 1.24
- **gRPC** + Protocol Buffers
- **MySQL** 8.0
- **Docker** & Docker Compose
- **Makefile** para automação

## 🎯 Próximos Passos

1. Leia [IMPLEMENTACAO_SHIPPING.md](./IMPLEMENTACAO_SHIPPING.md) para detalhes técnicos
2. Execute `make up` e teste com `make test`
3. Explore os serviços em `microservices/`

## 📄 Licença

MIT

