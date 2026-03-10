# Architecture Decision Records — finapp

> Backend Go para app de finanças pessoais com integração Pluggy.
> Data: 2026-03-09

---

## ADR-001: Chi como roteador HTTP

### Contexto
Precisávamos de um roteador HTTP para Go com suporte a middleware composável, parâmetros de rota e agrupamento de rotas protegidas por JWT.

### Decisão
Usar `go-chi/chi/v5`.

### Rationale
- **Compatível com `net/http`**: handlers seguem a assinatura padrão `http.HandlerFunc`, sem lock-in no framework. Handlers podem ser testados unitariamente sem instanciar o router.
- **Middleware por grupo**: o `r.Group(func(r chi.Router) { r.Use(auth) })` permite aplicar JWT somente nas rotas protegidas, sem afetar endpoints públicos (auth, webhook).
- **Leve e sem reflection**: Gin usa reflection para binding de parâmetros; Chi usa `chi.URLParam(r, "id")` — simples, explícito, performático.
- **Sem geração de código**: Echo e Gin têm helpers de validação embutidos que adicionam dependências implícitas. Chi mantém o projeto com controle total sobre validação.

### Consequências
- Parsing de JSON e validação são feitos manualmente nos handlers (mais verboso, mas mais explícito e testável).
- Suporte a middleware padrão do ecossistema `net/http` sem adaptadores.

---

## ADR-002: pgx v5 como driver PostgreSQL

### Contexto
Precisávamos de um driver PostgreSQL que suportasse tipos nativos do Postgres (UUID, JSONB, arrays, enums) sem wrappers.

### Decisão
Usar `jackc/pgx/v5` com `pgxpool`.

### Rationale
- **Protocolo nativo**: pgx implementa o wire protocol do PostgreSQL diretamente, sem a camada de abstração `database/sql`. Isso elimina overhead de conversão de tipo.
- **Tipos nativos**: `uuid.UUID`, `pgtype.JSONB`, `[]string` (arrays), enums mapeiam diretamente — sem precisar de `pq.Array()` ou wrappers.
- **pgxpool embutido**: pool de conexões com healthcheck, `MaxConns`, `MinConns`, `MaxConnLifetime` configuráveis — sem precisar de pgbouncer na camada da aplicação.
- **Melhor tratamento de erros**: `pgconn.PgError` expõe `Code` (código SQLSTATE) para detectar violações de constraint, conflitos e FK errors sem parsing de string.
- **`pgx.CollectRows`**: reduz boilerplate de scanning de rows.

### Consequências
- Não é compatível com `database/sql`, o que limita o uso de ORMs padrão (GORM, sqlx). Aceitável dado que usamos SQL puro.
- Curva de aprendizado ligeiramente maior que `lib/pq`, compensada pela expressividade e performance.

---

## ADR-003: Cache local dos dados Pluggy no PostgreSQL

### Contexto
A Pluggy retorna dados financeiros em tempo real, mas exige autenticação com TTL de 2h e tem rate limits. Realizar proxy direto de cada request do frontend para a Pluggy seria lento, frágil e caro.

### Decisão
Sincronizar dados da Pluggy (contas, transações) para tabelas locais no PostgreSQL e servir o frontend a partir do banco local.

### Rationale
- **Performance**: queries SQL locais são sub-milissegundo. Chamadas HTTP para a Pluggy têm latência de 200-800ms.
- **Aggregation**: relatórios, filtros, GROUP BY e ordenação são nativos do SQL. Fazer isso sobre a API da Pluggy exigiria buscar todos os dados e agregar em memória.
- **Enriquecimento do usuário**: categorias customizadas, tags e notas do usuário só fazem sentido em dado local — não existem no modelo de dados da Pluggy.
- **Resiliência**: se a Pluggy tiver downtime, o app continua funcionando com os dados já sincronizados. Um botão "Sincronizar agora" (`POST /pluggy/sync`) trata a expectativa de atualização do usuário.
- **Webhooks**: a Pluggy notifica via webhook quando um item é atualizado (`item/updated`). O backend recebe o webhook e dispara sincronização — dados ficam frescos poucos segundos após cada atualização bancária.

### Consequências
- Dados podem ter segundos/minutos de defasagem em relação à Pluggy (aceitável para finanças pessoais).
- Armazenamento extra no PostgreSQL. Mitigado por índices seletivos e sem armazenar dados brutos da Pluggy que não são usados.

---

## ADR-004: Autenticação JWT stateless

### Contexto
O app terá múltiplos usuários. Precisávamos de um mecanismo de autenticação que não introduza dependência de estado de sessão compartilhado.

### Decisão
JWT com access token de 15 minutos e refresh token de 7 dias. Assinado com HS256 e `JWT_SECRET` de no mínimo 32 chars.

### Rationale
- **Stateless**: o access token carrega `user_id` e `email` nas claims — nenhuma consulta ao banco para validar autenticidade. Middleware valida apenas a assinatura e a expiração.
- **Escalabilidade**: múltiplas instâncias do servidor compartilham apenas o `JWT_SECRET` — sem Redis ou banco de sessões.
- **Access token curto (15min)**: limita o blast radius em caso de vazamento. O refresh token (7 dias) permite sessões longas sem relogin.
- **Separação por `type` claim**: access tokens têm `type: "access"`, refresh tokens têm `type: "refresh"`. O middleware rejeita refresh tokens usados em rotas protegidas.

### Consequências
- Revogação de token individual não é possível sem um blocklist (Redis ou tabela). Aceitável no MVP; pode ser adicionada futuramente com tabela `revoked_tokens`.
- O `JWT_SECRET` é uma dependência crítica de segurança — deve ser gerenciado via secrets manager em produção.

---

## ADR-005: Arquitetura em camadas (Handler → Service → Repository)

### Contexto
O projeto tem múltiplos domínios (auth, pluggy, transações, orçamentos, metas, relatórios, simulações, projeções). Precisávamos de uma estrutura que separasse responsabilidades e permitisse testes unitários por camada.

### Decisão
Três camadas com interfaces explícitas em cada limite:

```
HTTP Request
     ↓
Handler (parse request, call service, write response — sem SQL, sem regras de negócio)
     ↓
Service (regras de negócio, cálculos, orquestração — sem HTTP, sem SQL)
     ↓
Repository (SQL + pgx — sem HTTP, sem regras de negócio)
     ↓
PostgreSQL
```

### Rationale
- **Handler**: só conhece `http.Request` e `http.ResponseWriter`. Delega tudo ao service.
- **Service**: detém regras como "não pode deletar categoria de sistema", "cálculo de juros compostos", "on_track de meta". Não sabe o que é HTTP nem SQL.
- **Repository**: só sabe executar queries. Recebe tipos do domínio, retorna tipos do domínio.
- **Interfaces em cada limite**: `repository.UserRepository`, `service.AuthService`, etc. Permitem mocks em testes sem DB real ou HTTP.
- Convenção Go: interfaces são definidas onde são consumidas (no package `service` ou `handler`), não onde são implementadas.

### Consequências
- Mais arquivos e indireção do que um handler que vai direto ao banco.
- Ganho: cada camada pode ser testada independentemente. Services de simulação e projeção são testáveis sem nenhum mock (lógica pura).

---

## ADR-006: golang-migrate com arquivos SQL numerados

### Contexto
O schema do banco evolui com o produto. Precisávamos de um mecanismo de versionamento de schema auditável, com suporte a rollback.

### Decisão
`golang-migrate/migrate/v4` com source `file://` e arquivos `.up.sql` / `.down.sql` numerados (`000001_`, `000002_`, ...).

### Rationale
- **SQL-first**: migrações legíveis por qualquer desenvolvedor ou DBA sem conhecer Go. O diff no git mostra exatamente o que muda no schema.
- **Auditável**: convenção `000001_create_users.up.sql` garante ordem determinística e identificação clara de cada migração no histórico do git.
- **Rollback**: cada migração tem um arquivo `.down.sql` com o DROP correspondente, permitindo reverter em emergências.
- **Integração no startup**: `main.go` chama `m.Up()` antes de iniciar o HTTP server — a aplicação nunca sobe com schema desatualizado.
- **Alternativa rejeitada — GORM AutoMigrate**: opaco, sem down migrations, acopla schema a struct tags, dificulta operações DDL avançadas (ENUMs, índices compostos, extensões).

### Consequências
- Migrações devem ser idempotentes quando possível (`CREATE TABLE IF NOT EXISTS`, `CREATE EXTENSION IF NOT EXISTS`).
- Em produção, migrações podem ser executadas separadamente do deploy (CI/CD pipeline com `migrate up` antes do rollout).

---

## Diagrama de Dependências

```
cmd/api/main.go
    ├── config
    ├── db (pgxpool)
    ├── repository/* (consome pgxpool, expõe interfaces)
    ├── pluggy/* (cliente HTTP Pluggy, token_manager)
    ├── service/* (consome repository interfaces + pluggy client)
    └── handler/* (consome service interfaces, expõe http.Handler)
           └── middleware/* (auth JWT, logging, recovery)
```

## Fluxo Pluggy

```
Frontend                Backend                     Pluggy API
   │                       │                             │
   │  POST /connect-token  │                             │
   │──────────────────────>│  POST /auth (CLIENT_ID)    │
   │                       │────────────────────────────>│
   │                       │<─── API Key (2h TTL) ───────│
   │                       │  POST /connect_token        │
   │                       │────────────────────────────>│
   │                       │<─── Connect Token (30min) ──│
   │<── { access_token } ──│                             │
   │                       │                             │
   │  [Pluggy Widget]      │                             │
   │──── connect token ───>│  Widget ──── POST /items ──>│
   │                       │                             │
   │                       │<─── Webhook: item/updated ──│
   │                       │  SyncItem()                 │
   │                       │──── GET /accounts ─────────>│
   │                       │──── GET /transactions ──────>│
   │                       │  Upsert local DB            │
   │  GET /transactions    │                             │
   │──────────────────────>│  SQL query (local DB)       │
   │<── transactions ──────│                             │
```
