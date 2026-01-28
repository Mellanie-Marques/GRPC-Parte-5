# Implementação do Microsserviço Shipping - Guia de Implementação

## 📝 O que foi criado:

✅ **Arquivos Protobuf:**
- `microservices-proto/shipping/shipping.proto` - Definição do serviço Shipping

✅ **Estrutura Shipping (Arquitetura Hexagonal):**
- `microservices/shipping/cmd/main.go` - Entrada da aplicação
- `microservices/shipping/internal/config/config.go` - Configurações
- `microservices/shipping/internal/application/core/domain/shipping.go` - Entidade de domínio
- `microservices/shipping/internal/application/core/api/api.go` - Lógica de negócio
- `microservices/shipping/internal/adapter/grpc/server.go` - Adapter gRPC
- `microservices/shipping/internal/ports/shipping.go` - Interface de porta
- `microservices/shipping/Dockerfile` - Container para Shipping
- `microservices/shipping/go.mod` - Dependências Go
- `microservices-proto/golang/shipping/go.mod` - Módulo protobuf

## 🔧 Próximos Passos - Manual:

### 1. Gerar arquivos protobuf
```bash
cd c:\Users\mella\Downloads\GRPC-Parte-3-main\GRPC-Parte-3-main
protoc -Imicroservices-proto/shipping --go_out=microservices-proto/golang/shipping --go-grpc_out=microservices-proto/golang/shipping microservices-proto/shipping/shipping.proto
```

### 2. Criar tabela de estoque (Product)
Adicionar à tabela inicial do banco:
```sql
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned AUTO_INCREMENT,
  `product_code` varchar(100) NOT NULL UNIQUE,
  `name` varchar(255),
  `quantity` int,
  PRIMARY KEY (`id`)
);

-- Inserir alguns produtos de teste
INSERT INTO products (product_code, name, quantity) VALUES 
('prod1', 'Produto 1', 100),
('prod2', 'Produto 2', 50),
('prod3', 'Produto 3', 75),
('prod4', 'Produto 4', 200);
```

### 3. Atualizar Order para:
- Validar produtos contra banco (adapter db)
- Chamar Shipping após sucesso de Payment

### 4. Adicionar Shipping ao docker-compose.yml
```yaml
  shipping:
    build:
      context: .
      dockerfile: microservices/shipping/Dockerfile
    container_name: ms-shipping
    depends_on:
      mysql:
        condition: service_healthy
    environment:
      APPLICATION_PORT: 3002
      ENV: development
    ports:
      - "3002:3002"
```

### 5. Atualizar Order para chamar Shipping
- Adicionar adapter de Shipping no Order
- Chamar após sucesso do pagamento

## 🎯 Funcionalidades do Shipping:
- Recebe OrderID e lista de itens
- Calcula prazo: 1 dia mínimo + 1 dia a cada 5 unidades
- Retorna número de dias de entrega

## 📦 Dependências necessárias:
```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

Esta é uma implementação complexa. Quer que eu continue com os próximos passos?
