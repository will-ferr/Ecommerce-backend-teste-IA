# Smart Choice - E-commerce Backend

Uma API de e-commerce robusta e escalável construída com Go, Gin Gonic e PostgreSQL, seguindo princípios SOLID e Clean Architecture.

## 🚀 Tecnologias

- **Framework**: Gin Gonic
- **ORM**: GORM
- **Banco de Dados**: PostgreSQL
- **Autenticação**: JWT + 2FA
- **Monitoramento**: Prometheus + OpenTelemetry
- **Container**: Docker + Docker Compose

## 📁 Estrutura do Projeto

```
├── config/          # Configurações e variáveis de ambiente
├── controllers/     # Handlers HTTP
├── database/        # Conexão e migrações do banco
├── logger/          # Configuração de logging
├── middlewares/     # Middlewares (autenticação, CORS, rate limiting)
├── models/          # Models do GORM
├── repository/      # Camada de acesso a dados
├── routes/          # Definição de rotas
├── services/        # Lógica de negócio
├── tracing/         # Configuração de OpenTelemetry
├── utils/           # Utilitários
├── Dockerfile
├── docker-compose.yml
└── prometheus.yml
```

## 🛡️ Funcionalidades

### Autenticação e Segurança
- JWT tokens com expiração configurável
- Suporte a Two-Factor Authentication (2FA)
- Rate limiting para prevenção de ataques
- CORS configurado
- Middleware de logging para auditoria

### Gestão de Produtos
- CRUD completo de produtos
- Filtros avançados (nome, preço, estoque)
- Paginação eficiente
- Alertas automáticos de estoque baixo via GORM Hooks

### Sistema de Cupons
- Validação de cupons (validade, uso máximo, valor mínimo)
- Controle de utilização

### Dashboard e Métricas
- Vendas diárias e mensais
- Contagem de novos usuários
- Status dos pedidos
- Métricas Prometheus

### Webhooks de Pagamento
- Endpoint seguro para webhooks
- Validação de assinatura HMAC
- Transações ACID para atualização de status

### SEO Backend
- Meta tags dinâmicas para produtos
- Open Graph tags
- URLs canônicas

## 🚀 Setup Rápido

### Pré-requisitos
- Docker e Docker Compose
- Go 1.25.5+ (para desenvolvimento local)

### Executando com Docker

1. Clone o repositório:
```bash
git clone <repository-url>
cd Smart-choice01
```

2. Configure as variáveis de ambiente:
```bash
cp .env.example .env
# Edite .env com suas configurações
```

3. Inicie os serviços:
```bash
docker-compose up -d
```

4. Acesse a API:
- API: http://localhost:8080
- Prometheus: http://localhost:9090
- PostgreSQL: localhost:5432

### Desenvolvimento Local

1. Instale dependências:
```bash
go mod download
```

2. Configure o banco PostgreSQL:
```bash
# Crie o banco de dados
createdb smart_choice
```

3. Execute a aplicação:
```bash
go run main.go
```

## 📚 Endpoints da API

### Autenticação
- `POST /auth/register` - Registro de usuário
- `POST /auth/login` - Login
- `POST /auth/2fa/generate` - Gerar 2FA
- `POST /auth/2fa/validate` - Validar 2FA

### Produtos
- `GET /api/products` - Listar produtos (com filtros)
- `GET /api/products/:id` - Obter produto
- `POST /api/products` - Criar produto (admin)
- `PUT /api/products/:id` - Atualizar produto (admin)
- `DELETE /api/products/:id` - Deletar produto (admin)

### Cupons
- `POST /api/coupons/validate` - Validar cupom

### Dashboard
- `GET /api/dashboard/metrics` - Métricas administrativas

### Webhooks
- `POST /webhooks/payment` - Webhook de pagamento

### SEO
- `GET /seo/product/:id` - Meta tags de produto
- `GET /seo/category/:category` - Meta tags de categoria
- `GET /seo/home` - Meta tags da home

### Sistema
- `GET /health` - Health check
- `GET /metrics` - Métricas Prometheus

## 🔧 Configuração

### Variáveis de Ambiente
```bash
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=smart_choice
DB_PORT=5432

# JWT
JWT_SECRET=your_super_secret_key

# Webhook
WEBHOOK_SECRET=your_webhook_secret

# Gin
GIN_MODE=release
```

## 📊 Monitoramento

### Prometheus
A aplicação expõe métricas em `/metrics`. O Prometheus está configurado para coletar:
- HTTP requests
- Latência
- Taxa de erros
- Uso de memória

### Logging
- Logs estruturados com zerolog
- Níveis: trace, debug, info, warn, error
- Logs administrativos para auditoria

## 🔒 Segurança

- Senhas hasheadas com bcrypt
- Tokens JWT com expiração
- Validação de entrada em todos os endpoints
- Rate limiting configurável
- CORS restrito
- Validação de webhook com HMAC

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar com coverage
go test -cover ./...

# Executar testes de benchmark
go test -bench=. ./...
```

## 📈 Performance

- Conexão pool com PostgreSQL
- Índices otimizados
- Paginação eficiente
- Cache configurável
- Middleware de compressão

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma feature branch
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a MIT License.

## 🆘 Suporte

Para dúvidas e suporte:
- Abra uma issue no GitHub
- Contato: [email]

---

**Smart Choice** - A escolha inteligente para seu e-commerce!
