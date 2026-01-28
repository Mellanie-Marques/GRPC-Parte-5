# CHECKLIST_DEPLOYMENT.md

# Checklist de Deployment

Use este documento como referência antes de fazer deploy em qualquer ambiente.

## 📋 Pré-Deployment (Dev/Staging/Prod)

- [ ] Clonar repositório
- [ ] Verificar versão do Docker (`docker --version`)
- [ ] Verificar versão do Docker Compose (`docker-compose --version`)
- [ ] Verificar espaço em disco (mínimo 2GB)
- [ ] Revisar `.env` com variáveis corretas
- [ ] Verificar que `docker-compose.yml` existe
- [ ] Verificar que `tmp_create_dbs.sql` existe
- [ ] Verificar que todos os `Dockerfile` existem

## 🔐 Segurança (Antes do Deploy)

### Desenvolvimento

- [ ] Usar senhas padrão (`minhasenha`) é aceitável
- [ ] Não expor portas públicas em localhost

### Staging

- [ ] Alterar `MYSQL_ROOT_PASSWORD`
- [ ] Usar variáveis de ambiente seguras
- [ ] Verificar firewall
- [ ] Revisar logs para erros sensíveis

### Produção

- [ ] ⚠️ OBRIGATÓRIO: Alterar `MYSQL_ROOT_PASSWORD`
- [ ] ⚠️ OBRIGATÓRIO: Usar variáveis de ambiente (não hardcoded)
- [ ] ⚠️ OBRIGATÓRIO: Implementar TLS/HTTPS
- [ ] ⚠️ OBRIGATÓRIO: Usar secrets do Docker/Kubernetes
- [ ] ⚠️ OBRIGATÓRIO: Configurar backup automático
- [ ] ⚠️ OBRIGATÓRIO: Revisar logs e monitoramento
- [ ] ⚠️ OBRIGATÓRIO: Implementar rate limiting
- [ ] ⚠️ OBRIGATÓRIO: Configurar health checks
- [ ] ⚠️ OBRIGATÓRIO: Testar recuperação de desastres
- [ ] ⚠️ OBRIGATÓRIO: Documentar plano de rollback

## 🚀 Build & Deploy

### 1. Preparação

```bash
- [ ] git clone <url>
- [ ] cd GRPC-Parte-3-main
- [ ] cp .env.example .env
- [ ] # Editar .env com valores corretos
```

### 2. Build Local (Opcional)

```bash
- [ ] go mod tidy (em cada serviço)
- [ ] go build (compilar localmente)
- [ ] go test (rodar testes)
```

### 3. Build Docker

```bash
- [ ] docker-compose build (criar imagens)
- [ ] docker images (verificar imagens criadas)
```

### 4. Deploy

```bash
- [ ] docker-compose up -d (iniciar)
- [ ] docker-compose ps (verificar status)
- [ ] sleep 10 (aguardar inicialização)
```

## ✅ Validação Pós-Deploy

### Containers

- [ ] `docker-compose ps` - Todos os 4 containers rodando
- [ ] Sem containers em estado `Exited`
- [ ] Sem containers com status `unhealthy`

### MySQL

- [ ] `docker-compose logs mysql | grep "ready for connections"`
- [ ] Banco de dados criados:
  ```bash
  docker-compose exec mysql mysql -u root -pSENHA -e "SHOW DATABASES;"
  ```
  - [ ] `information_schema`
  - [ ] `mysql`
  - [ ] `order_db`
  - [ ] `payment`

- [ ] Tabelas criadas:
  ```bash
  docker-compose exec mysql mysql -u root -pSENHA -D order_db -e "SHOW TABLES;"
  ```
  - [ ] `orders`
  - [ ] `products`

### Serviços gRPC

- [ ] Order (porta 3000) respondendo
- [ ] Payment (porta 3001) respondendo
- [ ] Shipping (porta 3002) respondendo

### Testes Funcionais

- [ ] Executar cliente de teste:
  ```bash
  cd microservices/order
  go run client/main.go
  ```
- [ ] Todos os 4 testes completados
- [ ] Nenhuma mensagem de erro crítico

### Logs

- [ ] Verificar logs do Order:
  ```bash
  docker-compose logs order | tail -20
  ```
- [ ] Verificar logs do Payment:
  ```bash
  docker-compose logs payment | tail -20
  ```
- [ ] Verificar logs do Shipping:
  ```bash
  docker-compose logs shipping | tail -20
  ```
- [ ] Sem erros (ERROR) nos logs

## 📊 Monitoramento Pós-Deploy

### Primeiras 24 horas

- [ ] Monitorar uso de memória:
  ```bash
  docker stats
  ```
- [ ] Monitorar uso de disco
- [ ] Verificar logs de erro
- [ ] Verificar timeout de conexões
- [ ] Testar funcionalidade básica

### Semana 1

- [ ] Executar carga mínima de testes
- [ ] Verificar performance
- [ ] Documentar problemas encontrados
- [ ] Verificar backup automático

### Mensal

- [ ] Revisar utilização de recursos
- [ ] Testar recuperação de backup
- [ ] Atualizar documentação
- [ ] Revisar logs de acesso

## 🔄 Backup & Recuperação

### Antes do Deploy

- [ ] Criar backup do banco anterior (se existe):
  ```bash
  docker-compose exec mysql mysqldump -u root -p --all-databases > backup_pre_deploy.sql
  ```

### Backup Pós-Deploy

- [ ] Confirmar que backup foi criado:
  ```bash
  ls -la backups/
  ```

### Teste de Recuperação

- [ ] Testar restauração de backup:
  ```bash
  docker-compose down -v
  docker-compose up -d mysql
  sleep 15
  docker-compose exec -T mysql mysql -u root -p < backup_pre_deploy.sql
  ```

## 🛑 Rollback Plan

### Caso de Falha

1. **Parar imediatamente**
   ```bash
   docker-compose down
   ```

2. **Restaurar dados**
   ```bash
   docker-compose up -d mysql
   docker-compose exec -T mysql mysql -u root -p < backup_pre_deploy.sql
   ```

3. **Reiniciar com versão anterior**
   ```bash
   git checkout <commit-anterior>
   docker-compose build
   docker-compose up -d
   ```

4. **Verificar status**
   ```bash
   docker-compose ps
   docker-compose logs
   ```

5. **Investigar causa**
   - Revisar logs
   - Verificar configuração
   - Testar componentes isoladamente

## 📈 Performance Baseline

Registrar depois do deploy bem-sucedido:

| Métrica | Valor | Data |
|---------|-------|------|
| Memory - MySQL | ___ MB | __/__/__ |
| Memory - Order | ___ MB | __/__/__ |
| Memory - Payment | ___ MB | __/__/__ |
| Memory - Shipping | ___ MB | __/__/__ |
| Disk Used | ___ GB | __/__/__ |
| Response Time (avg) | ___ ms | __/__/__ |
| Requests/sec | ___ | __/__/__ |

## 📞 Contato & Suporte

### Em caso de problema:

1. Verificar logs: `docker-compose logs`
2. Consultar [TROUBLESHOOTING.md](./DEPLOYMENT.md#troubleshooting)
3. Consultar [README.md](./README.md)
4. Abrir issue no repositório

### Informações úteis a fornecer:

- [ ] Ambiente (desenvolvimento/staging/produção)
- [ ] Saída de `docker-compose ps`
- [ ] Saída de `docker-compose logs` (últimas 50 linhas)
- [ ] Versão do Docker (`docker --version`)
- [ ] Sistema operacional
- [ ] Passos para reproduzir

## ✨ Checklist Completo

Antes de considerar o deploy bem-sucedido:

- [ ] Todos os containers rodando
- [ ] Testes funcionais passando
- [ ] Logs sem erros críticos
- [ ] Backup criado e testado
- [ ] Documentação atualizada
- [ ] Time notificado
- [ ] Monitoramento ativado

## 🎉 Deploy Concluído!

Se todos os itens acima foram marcados, o deploy foi bem-sucedido!

---

**Importante**: Mantenha esta checklist próximo durante cada deployment.

**Data do último deploy**: ________________  
**Responsável**: ________________________  
**Versão implantada**: ___________________  
**Ambiente**: ____________________________  
**Notas adicionais**: ______________________________________________________

---

**Última Atualização**: Janeiro de 2026  
**Versão**: 1.0
