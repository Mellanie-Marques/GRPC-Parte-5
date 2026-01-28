.PHONY: help build up down logs test clean restart ps status lint fmt

# Variáveis
COMPOSE_FILE := docker-compose.yml
DOCKER_COMPOSE := docker-compose -f $(COMPOSE_FILE)
SERVICES := order payment shipping mysql

help: ## Mostrar esta mensagem de ajuda
	@echo "Microservices com gRPC - Makefile"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==========================================
# DOCKER & COMPOSE
# ==========================================

build: ## Build das imagens Docker
	@echo "📦 Building Docker images..."
	$(DOCKER_COMPOSE) build

up: ## Iniciar todos os serviços (com build)
	@echo "🚀 Starting services..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "✅ Services started successfully!"
	@echo "📊 Order:    http://localhost:3000 (gRPC)"
	@echo "💳 Payment:  http://localhost:3001 (gRPC)"
	@echo "🚚 Shipping: http://localhost:3002 (gRPC)"
	@echo "🗄️  MySQL:    localhost:3308"

up-no-build: ## Iniciar serviços sem fazer build
	@echo "🚀 Starting services..."
	$(DOCKER_COMPOSE) up -d

down: ## Parar todos os serviços
	@echo "🛑 Stopping services..."
	$(DOCKER_COMPOSE) down

restart: ## Reiniciar todos os serviços
	@echo "🔄 Restarting services..."
	$(DOCKER_COMPOSE) restart

ps: status ## Ver status dos containers (alias)

status: ## Ver status dos containers
	@echo "📊 Container Status:"
	@$(DOCKER_COMPOSE) ps

# ==========================================
# LOGS
# ==========================================

logs: ## Ver logs de todos os serviços (últimas 50 linhas)
	@$(DOCKER_COMPOSE) logs --tail=50

logs-follow: ## Ver logs em tempo real
	@$(DOCKER_COMPOSE) logs -f

logs-order: ## Ver logs do Order Service
	@$(DOCKER_COMPOSE) logs -f order

logs-payment: ## Ver logs do Payment Service
	@$(DOCKER_COMPOSE) logs -f payment

logs-shipping: ## Ver logs do Shipping Service
	@$(DOCKER_COMPOSE) logs -f shipping

logs-mysql: ## Ver logs do MySQL
	@$(DOCKER_COMPOSE) logs -f mysql

# ==========================================
# DATABASE
# ==========================================

mysql-cli: ## Acessar MySQL CLI
	@$(DOCKER_COMPOSE) exec mysql mysql -u root -pminhasenha

db-backup: ## Fazer backup do banco de dados
	@echo "💾 Backing up database..."
	@mkdir -p ./backups
	@$(DOCKER_COMPOSE) exec -T mysql mysqldump -u root -pminhasenha --all-databases > ./backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup completed: ./backups/backup_$(shell date +%Y%m%d_%H%M%S).sql"

db-restore: ## Restaurar banco de dados (use BACKUP_FILE=./backups/backup_XXX.sql)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "❌ Erro: BACKUP_FILE não definido"; \
		echo "Uso: make db-restore BACKUP_FILE=./backups/backup_XXX.sql"; \
	else \
		echo "🔄 Restoring database from $(BACKUP_FILE)..."; \
		$(DOCKER_COMPOSE) exec -T mysql mysql -u root -pminhasenha < $(BACKUP_FILE); \
		echo "✅ Database restored successfully"; \
	fi

db-init: ## Reinicializar banco de dados (remove volumes)
	@echo "⚠️  Removendo volumes do banco de dados..."
	@$(DOCKER_COMPOSE) down -v
	@echo "✅ Database reset. Execute 'make up' para criar novamente"

# ==========================================
# TESTES
# ==========================================

test: ## Executar cliente de teste
	@echo "🧪 Running test client..."
	@cd microservices/order && go run client/main.go

test-verbose: ## Executar testes com saída detalhada
	@cd microservices/order && go run client/main.go -v

unit-test: ## Executar unit tests dos serviços
	@echo "🧪 Running unit tests..."
	@cd microservices/order && go test -v ./...
	@cd ../payment && go test -v ./...
	@cd ../shipping && go test -v ./...

# ==========================================
# BUILD LOCAL
# ==========================================

build-order: ## Compilar Order Service localmente
	@echo "🔨 Building Order Service..."
	@cd microservices/order && go build -o order ./cmd/main.go
	@echo "✅ Order Service built: microservices/order/order"

build-payment: ## Compilar Payment Service localmente
	@echo "🔨 Building Payment Service..."
	@cd microservices/payment && go build -o payment ./cmd/main.go
	@echo "✅ Payment Service built: microservices/payment/payment"

build-shipping: ## Compilar Shipping Service localmente
	@echo "🔨 Building Shipping Service..."
	@cd microservices/shipping && go build -o shipping ./cmd/main.go
	@echo "✅ Shipping Service built: microservices/shipping/shipping"

build-all-local: build-order build-payment build-shipping ## Compilar todos os serviços localmente

# ==========================================
# GO MOD
# ==========================================

mod-tidy: ## Executar go mod tidy em todos os serviços
	@echo "🔄 Running go mod tidy..."
	@cd microservices/order && go mod tidy
	@cd ../payment && go mod tidy
	@cd ../shipping && go mod tidy
	@echo "✅ Dependencies updated"

# ==========================================
# LINTING & FORMATTING
# ==========================================

fmt: ## Formatar código Go
	@echo "💅 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Executar linter (requer golangci-lint instalado)
	@echo "🔍 Running linter..."
	@cd microservices/order && golangci-lint run
	@cd ../payment && golangci-lint run
	@cd ../shipping && golangci-lint run

# ==========================================
# LIMPEZA
# ==========================================

clean: ## Remover containers, imagens e volumes
	@echo "🧹 Cleaning up..."
	@$(DOCKER_COMPOSE) down -v
	@docker image prune -f
	@echo "✅ Cleanup completed"

clean-containers: ## Remover apenas containers
	@echo "🧹 Removing containers..."
	@$(DOCKER_COMPOSE) down

clean-volumes: ## Remover apenas volumes
	@echo "🧹 Removing volumes..."
	@$(DOCKER_COMPOSE) down -v

clean-images: ## Remover apenas imagens
	@echo "🧹 Removing images..."
	@docker image prune -f

clean-all: ## Limpeza completa (containers, volumes, imagens, cache)
	@echo "⚠️  Executando limpeza completa..."
	@$(DOCKER_COMPOSE) down -v
	@docker system prune -f
	@echo "✅ Full cleanup completed"

# ==========================================
# DOCKER SYSTEM
# ==========================================

docker-prune: ## Limpar cache e recursos não usados do Docker
	@echo "🧹 Pruning Docker system..."
	@docker system prune -f
	@echo "✅ Docker system pruned"

docker-stats: ## Ver uso de recursos dos containers
	@echo "📊 Container Resource Usage:"
	@docker stats --no-stream

# ==========================================
# PROTOBUF (Desenvolvimento)
# ==========================================

proto-generate: ## Gerar código Protobuf (requer protoc instalado)
	@echo "🔄 Generating Protobuf code..."
	@protoc --go_out=microservices-proto/golang/order \
		--go-grpc_out=microservices-proto/golang/order \
		microservices-proto/order/order.proto
	@protoc --go_out=microservices-proto/golang/payment \
		--go-grpc_out=microservices-proto/golang/payment \
		microservices-proto/payment/payment.proto
	@protoc --go_out=microservices-proto/golang/shipping \
		--go-grpc_out=microservices-proto/golang/shipping \
		microservices-proto/shipping/shipping.proto
	@echo "✅ Protobuf code generated"

# ==========================================
# DEVELOPMENT
# ==========================================

dev-order: ## Executar Order Service em modo desenvolvimento (localhost)
	@echo "🚀 Starting Order Service..."
	@export DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3308)/order_db"; \
	export PAYMENT_SERVICE_URL="localhost:3001"; \
	export SHIPPING_SERVICE_URL="localhost:3002"; \
	export APPLICATION_PORT="3000"; \
	cd microservices/order && go run cmd/main.go

dev-payment: ## Executar Payment Service em modo desenvolvimento
	@echo "🚀 Starting Payment Service..."
	@export DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3308)/payment"; \
	export APPLICATION_PORT="3001"; \
	cd microservices/payment && go run cmd/main.go

dev-shipping: ## Executar Shipping Service em modo desenvolvimento
	@echo "🚀 Starting Shipping Service..."
	@export APPLICATION_PORT="3002"; \
	cd microservices/shipping && go run cmd/main.go

# ==========================================
# DOCUMENTAÇÃO
# ==========================================

docs: ## Abrir documentação principal (README)
	@echo "📖 Opening README.md..."
	@cat README.md

docs-deployment: ## Abrir documentação de deployment
	@echo "📖 Opening DEPLOYMENT.md..."
	@cat DEPLOYMENT.md

# ==========================================
# QUICK START
# ==========================================

quick-start: clean up test ## Quick start: limpar, subir, testar

full-setup: mod-tidy build up test ## Setup completo: mod tidy, build, up, test

# ==========================================
# INFO
# ==========================================

info: ## Mostrar informações do sistema
	@echo "🖥️  System Information"
	@echo "====================="
	@echo "Docker Version:"
	@docker --version
	@echo "\nDocker Compose Version:"
	@docker-compose --version
	@echo "\nGo Version:"
	@go version
	@echo "\nServices Status:"
	@$(DOCKER_COMPOSE) ps || echo "❌ Services not running"

# =============================================
# PADRÃO
# =============================================

.DEFAULT_GOAL := help
