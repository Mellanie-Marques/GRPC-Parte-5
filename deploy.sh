#!/bin/bash

# Deploy Script para Microsserviços gRPC
# Uso: ./deploy.sh [develop|staging|production]
# Exemplo: ./deploy.sh production

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENVIRONMENT=${1:-development}
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Funções
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Validações
print_header "Microsserviços com gRPC - Deploy"

if [ ! -f "$SCRIPT_DIR/docker-compose.yml" ]; then
    print_error "docker-compose.yml não encontrado!"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    print_error "Docker não está instalado!"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose não está instalado!"
    exit 1
fi

print_success "Dependências validadas"

# Validar ambiente
case "$ENVIRONMENT" in
    development|staging|production)
        print_info "Ambiente: $ENVIRONMENT"
        ;;
    *)
        print_error "Ambiente inválido: $ENVIRONMENT"
        echo "Opções válidas: development, staging, production"
        exit 1
        ;;
esac

# Carregar arquivo de ambiente
ENV_FILE="$SCRIPT_DIR/.env.${ENVIRONMENT}"
if [ -f "$ENV_FILE" ]; then
    print_info "Carregando variáveis de $ENV_FILE"
    set -a
    source "$ENV_FILE"
    set +a
    print_success "Variáveis carregadas"
elif [ "$ENVIRONMENT" != "development" ]; then
    print_error "Arquivo de configuração não encontrado: $ENV_FILE"
    print_info "Crie o arquivo com: cp .env.example $ENV_FILE"
    exit 1
else
    print_warning "Usando configurações padrão para desenvolvimento"
fi

# Backup do banco de dados (se estiver rodando)
if docker ps | grep -q "ms-mysql"; then
    print_info "Backup do banco de dados..."
    mkdir -p "$SCRIPT_DIR/backups"
    docker-compose exec -T mysql mysqldump -u root -p"${MYSQL_ROOT_PASSWORD:-minhasenha}" \
        --all-databases > "$SCRIPT_DIR/backups/backup_${TIMESTAMP}.sql"
    print_success "Backup criado: backups/backup_${TIMESTAMP}.sql"
fi

# Selecionar arquivo docker-compose
if [ "$ENVIRONMENT" == "production" ]; then
    COMPOSE_FILE="docker-compose.prod.yml"
else
    COMPOSE_FILE="docker-compose.yml"
fi

if [ ! -f "$SCRIPT_DIR/$COMPOSE_FILE" ]; then
    print_warning "Arquivo $COMPOSE_FILE não encontrado, usando docker-compose.yml"
    COMPOSE_FILE="docker-compose.yml"
fi

print_info "Usando arquivo: $COMPOSE_FILE"

# Build
print_header "Build das Imagens"
print_info "Construindo imagens Docker..."
cd "$SCRIPT_DIR"
docker-compose -f "$COMPOSE_FILE" build
print_success "Build concluído"

# Parar containers antigos
print_header "Parando Containers"
if docker-compose -f "$COMPOSE_FILE" ps 2>/dev/null | grep -q "ms-"; then
    print_info "Parando containers antigos..."
    docker-compose -f "$COMPOSE_FILE" stop
    print_success "Containers parados"
else
    print_info "Nenhum container em execução"
fi

# Iniciar containers
print_header "Iniciando Serviços"
print_info "Iniciando containers..."
docker-compose -f "$COMPOSE_FILE" up -d

# Aguardar saúde dos serviços
print_info "Aguardando serviços iniciarem..."
sleep 10

# Verificar saúde
print_header "Verificação de Saúde"
for service in mysql payment order shipping; do
    if docker-compose -f "$COMPOSE_FILE" ps | grep -q "$service"; then
        print_success "Serviço $service está rodando"
    else
        print_error "Serviço $service NÃO está rodando!"
    fi
done

# Teste de conectividade
print_header "Teste de Conectividade"

# MySQL
if docker-compose -f "$COMPOSE_FILE" exec -T mysql mysqladmin -u root -p"${MYSQL_ROOT_PASSWORD:-minhasenha}" ping >/dev/null 2>&1; then
    print_success "MySQL respondendo"
else
    print_error "MySQL não está respondendo"
fi

# Mostrar status
print_header "Status Final"
docker-compose -f "$COMPOSE_FILE" ps

# Mostrar logs
print_info "Últimas linhas dos logs:"
echo ""
docker-compose -f "$COMPOSE_FILE" logs --tail=5

# Resumo
print_header "Deploy Concluído"
echo ""
echo "Ambiente: $ENVIRONMENT"
echo "Timestamp: $TIMESTAMP"
echo ""
echo "Serviços disponíveis:"
echo "  📋 Order Service:    http://localhost:3000 (gRPC)"
echo "  💳 Payment Service:  http://localhost:3001 (gRPC)"
echo "  🚚 Shipping Service: http://localhost:3002 (gRPC)"
echo "  🗄️  MySQL:            localhost:3308"
echo ""
echo "Comandos úteis:"
echo "  Ver logs:           docker-compose logs -f"
echo "  Ver status:         docker-compose ps"
echo "  Parar serviços:     docker-compose down"
echo "  Restaurar backup:   docker-compose exec -T mysql mysql -u root -p < backups/backup_XXX.sql"
echo ""

print_success "Deploy finalizado com sucesso!"

exit 0
