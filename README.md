# Microservices com gRPC - Parte 3

Este projeto implementa uma arquitetura de microsserviços completa utilizando **gRPC** em **Go**, com três serviços independentes: **Order**, **Payment** e **Shipping**. Todos os serviços possuem persistência em **MySQL** e comunicação assíncrona via gRPC com mecanismos de resiliência (retry e timeout).

## Características

✅ **Arquitetura Hexagonal** em todos os microsserviços  
✅ **gRPC** para comunicação entre serviços  
✅ **Protocol Buffers** para serialização eficiente  
✅ **MySQL** para persistência de dados  
✅ **Docker & Docker Compose** para orquestração  
✅ **Retry automático** com backoff linear (5 tentativas)  
✅ **Timeout** de 2 segundos por requisição  
✅ **Validação de produtos** antes do processamento  
✅ **Cálculo de entrega** baseado na quantidade de itens  

## Pré-requisitos

### Sistema
- **Docker Desktop** (Windows/Mac) ou **Docker Engine** (Linux)
- **Docker Compose** v2.0 ou superior
- **Git** (para clonar o repositório)

### Para Desenvolvimento Local (Opcional)
- **Go** 1.24 ou superior
- **protoc** (compilador Protocol Buffers)
- **MySQL Client** (para acessar o banco manualmente)

## Instalação Rápida

### Passo 1: Clone o repositório

```bash
git clone <url-do-repositorio>
cd GRPC-Parte-3-main
```

### Passo 2: Execute com Docker Compose

```bash
docker-compose up --build
```

Isso irá:
- ✅ Construir e iniciar o container **MySQL 8.0** (porta 3308)
- ✅ Construir e iniciar o serviço **Payment** (porta 3001)
- ✅ Construir e iniciar o serviço **Order** (porta 3000)
- ✅ Construir e iniciar o serviço **Shipping** (porta 3002)
- ✅ Inicializar o banco de dados com tabelas e dados de exemplo

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
├── README.md                       # Esta documentação
│
├── microservices/
│   ├── order/                     # Serviço de Pedidos
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── cmd/main.go
│   │   ├── client/main.go         # Cliente de teste (português)
│   │   ├── config/
│   │   └── internal/
│   │       ├── adapter/
│   │       ├── adapters/
│   │       ├── application/
│   │       └── ports/
│   │
│   ├── payment/                   # Serviço de Pagamentos
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── cmd/main.go
│   │   └── internal/
│   │
│   └── shipping/                  # Serviço de Entrega
│       ├── Dockerfile
│       ├── go.mod
│       ├── cmd/main.go
│       └── internal/
│
└── microservices-proto/
    ├── order/order.proto          # Definição do serviço Order
    ├── payment/payment.proto      # Definição do serviço Payment
    ├── shipping/shipping.proto    # Definição do serviço Shipping
    └── golang/                    # Código gerado a partir dos .proto
```

## Arquitetura dos Microsserviços

### 1. Order Service (Porta 3000)

**Responsabilidade**: Orquestra todo o fluxo de pedidos

**Fluxo**:
1. Recebe requisição com itens do pedido
2. ✅ Valida se todos os produtos existem no banco
3. ✅ Salva o pedido no banco de dados
4. ✅ Processa pagamento no Payment Service
5. ✅ Se pagamento OK → Calcula entrega no Shipping Service
6. ✅ Retorna pedido com dias de entrega

**Comunicação**:
- HTTP/gRPC de entrada (cliente)
- gRPC saída → Payment Service
- gRPC saída → Shipping Service
- MySQL para persistência

**Tabelas**:
- `orders` - Armazena pedidos
- `products` - Armazena produtos para validação

### 2. Payment Service (Porta 3001)

**Responsabilidade**: Processa pagamentos

**Fluxo**:
1. Recebe requisição de pagamento
2. Valida dados do pedido
3. Processa pagamento (simulado)
4. Retorna confirmação

**Comunicação**:
- gRPC de entrada
- MySQL para persistência

**Tabelas**:
- `payments` - Armazena histórico de pagamentos

### 3. Shipping Service (Porta 3002)

**Responsabilidade**: Calcula prazos de entrega

**Fórmula de Cálculo**:
```
Dias de Entrega = 1 + (Quantidade Total - 1) / 5
```

Exemplos:
- 1 item = 1 dia
- 2-5 itens = 1 dia
- 6-10 itens = 2 dias
- 11-15 itens = 3 dias

**Comunicação**:
- gRPC de entrada
- Sem persistência (serviço stateless)

## Variáveis de Ambiente

### Order Service

```bash
PAYMENT_SERVICE_URL=payment:3001      # URL do serviço de pagamentos
SHIPPING_SERVICE_URL=shipping:3002    # URL do serviço de envios
DATA_SOURCE_URL=root:minhasenha@tcp(mysql:3306)/order_db  # Banco MySQL
APPLICATION_PORT=3000                 # Porta gRPC
ENV=development                        # Ambiente
```

### Payment Service

```bash
DATA_SOURCE_URL=root:minhasenha@tcp(mysql:3306)/payment  # Banco MySQL
APPLICATION_PORT=3001                 # Porta gRPC
ENV=development                        # Ambiente
```

### Shipping Service

```bash
APPLICATION_PORT=3002                 # Porta gRPC
ENV=development                        # Ambiente
```

## Banco de Dados

### Inicialização Automática

O arquivo `tmp_create_dbs.sql` é executado automaticamente ao iniciar o MySQL:

- **Cria bancos**: `order_db`, `payment`
- **Cria tabelas**: `orders`, `payments`, `products`
- **Insere dados**: Produtos de exemplo (prod1, prod2, ABC123, XYZ789)

### Estrutura das Tabelas

**products** (order_db):
```sql
CREATE TABLE products (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  product_code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**orders** (order_db):
```sql
CREATE TABLE orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at BIGINT NOT NULL
);
```

**payments** (payment):
```sql
CREATE TABLE payments (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Como Usar

### Opção 1: Docker Compose (Recomendado - Produção)

```bash
# Iniciar todos os serviços
docker-compose up --build

# Em outro terminal, testar
cd microservices/order
go run client/main.go

# Parar tudo
docker-compose down
```

### Opção 2: Execução Local (Desenvolvimento)

#### Terminal 1 - MySQL
```bash
docker-compose up mysql -d
```

#### Terminal 2 - Payment Service
```bash
export PAYMENT_SERVICE_URL=localhost:3001
export DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3308)/payment"
export APPLICATION_PORT=3001

cd microservices/payment
go run cmd/main.go
```

#### Terminal 3 - Shipping Service
```bash
export APPLICATION_PORT=3002
cd microservices/shipping
go run cmd/main.go
```

#### Terminal 4 - Order Service
```bash
export PAYMENT_SERVICE_URL=localhost:3001
export SHIPPING_SERVICE_URL=localhost:3002
export DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3308)/order_db"
export APPLICATION_PORT=3000

cd microservices/order
go run cmd/main.go
```

#### Terminal 5 - Cliente de Teste
```bash
cd microservices/order
go run client/main.go
```

### Opção 3: Acessar o Banco MySQL

Com o Docker Compose rodando:

```bash
docker-compose exec mysql mysql -u root -pminhasenha -D order_db
```

Consultas úteis:
```sql
-- Ver produtos
SELECT * FROM products;

-- Ver pedidos
SELECT * FROM orders;

-- Ver pagamentos
SELECT * FROM payments;
```

## Mecanismos de Resiliência

### Retry Automático
- **Máximo**: 5 tentativas
- **Backoff**: Linear de 1 segundo entre tentativas
- **Códigos**: Retorna em caso de `Unavailable` ou `ResourceExhausted`

### Timeout
- **Duração**: 2 segundos por requisição
- **Erro**: `DeadlineExceeded` se exceder

### Tratamento de Erros
```
Produto não existe → ErrProductNotFound
Pagamento falhou → Payment error
Entrega falhou → Shipping error
```

## Logs e Monitoramento

### Ver Logs de um Serviço

```bash
# Order Service
docker-compose logs ms-order -f

# Payment Service
docker-compose logs ms-payment -f

# Shipping Service
docker-compose logs ms-shipping -f

# MySQL
docker-compose logs ms-mysql -f
```

### Ver Todos os Logs
```bash
docker-compose logs -f
```

## Desenvolvimento

### Gerar Código Protobuf (Se Modificar .proto)

Você precisará do `protoc` instalado:

```bash
# Windows
choco install protoc

# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler
```

Depois instale os plugins do Go:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Gerar código:
```bash
# Para Order
protoc --go_out=microservices-proto/golang/order \
        --go-grpc_out=microservices-proto/golang/order \
        microservices-proto/order/order.proto

# Para Payment
protoc --go_out=microservices-proto/golang/payment \
        --go-grpc_out=microservices-proto/golang/payment \
        microservices-proto/payment/payment.proto

# Para Shipping
protoc --go_out=microservices-proto/golang/shipping \
        --go-grpc_out=microservices-proto/golang/shipping \
        microservices-proto/shipping/shipping.proto
```

### Compilar Localmente

```bash
# Order
cd microservices/order
go build -o order ./cmd/main.go

# Payment
cd microservices/payment
go build -o payment ./cmd/main.go

# Shipping
cd microservices/shipping
go build -o shipping ./cmd/main.go
```

### Executar Testes

```bash
# Testes do Order Service
cd microservices/order
go test ./...
```

## Troubleshooting

### Problema: "connection refused"

**Solução**: Aguarde o MySQL inicializar (leva ~10 segundos)

```bash
docker-compose logs mysql
# Procure por "ready for connections"
```

### Problema: "port already in use"

**Solução**: Mude as portas no `docker-compose.yml`:

```yaml
ports:
  - "3300:3000"  # Order em porta 3300
  - "3301:3001"  # Payment em porta 3301
```

### Problema: Build do Docker falha

**Solução**: Limpe o cache do Docker

```bash
docker-compose down
docker system prune -f
docker-compose up --build
```

### Problema: Banco não inicializa

**Solução**: Verifique o arquivo `tmp_create_dbs.sql`

```bash
docker-compose logs mysql | grep "ERROR"
```

## Performance e Escalabilidade

### Limite de Requisições Simultâneas
- Order Service: Sem limite (Go)
- Payment Service: Sem limite (Go)
- Shipping Service: Sem limite (Go)

### Timeout de Conexão
- 2 segundos por requisição gRPC
- Retry automático até 5 vezes

### Conexão com MySQL
- Connection pool automático do GORM
- Max open connections: 25 (padrão)
- Max idle connections: 5 (padrão)

## Segurança

⚠️ **Nota**: Este é um projeto educacional. Para produção:

- [ ] Use credenciais do banco seguras (secrets)
- [ ] Implemente autenticação gRPC (mTLS)
- [ ] Use variáveis de ambiente para senhas
- [ ] Configure HTTPS/TLS
- [ ] Implemente rate limiting
- [ ] Adicione validação de entrada robusta
- [ ] Use network policies do Docker

## Contribuindo

1. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
2. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
3. Push para a branch (`git push origin feature/AmazingFeature`)
4. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

## Suporte

Para problemas ou dúvidas:
1. Verifique os logs: `docker-compose logs`
2. Consulte a seção Troubleshooting
3. Abra uma issue no repositório

