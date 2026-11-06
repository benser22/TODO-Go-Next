# TODO — API Go + Next.js + Postgres

TODO clásico full-stack: API REST en Go (stdlib + pgx + sqlc), PostgreSQL en Docker
y front en Next.js (App Router, TypeScript, Tailwind).

```
backend/    API Go: cmd/todo-api + internal/{config,domain,service,handler,store/postgres,db}
frontend/   UI Next.js: app/ + components/todo-app.tsx + lib/api.ts
docker-compose.yml   Postgres 17 + API (dockerizada) + web (dockerizada, standalone)
```

## Arranque con un solo comando

```powershell
docker compose up --build
```

- UI: http://localhost:3000
- API: http://localhost:8080 (healthz, readyz)
- Postgres en el host: localhost:5433 (credenciales `todo/todo`, BD `todos`)

El contenedor de la API aplica las migraciones con goose automáticamente antes
de arrancar. Desarrollo sin Docker: `go run ./cmd/todo-api` en `backend/` y
`npm run dev` en `frontend/` (usa el `.env` local).

## API

| Método | Ruta | Descripción | Códigos |
|---|---|---|---|
| GET | `/api/v1/todos` | listar | 200 |
| POST | `/api/v1/todos` | crear `{ "title": "..." }` | 201 / 400 |
| GET | `/api/v1/todos/{id}` | obtener | 200 / 400 / 404 |
| PATCH | `/api/v1/todos/{id}` | editar parcial `{ "title"?, "completed"? }` | 200 / 400 / 404 |
| DELETE | `/api/v1/todos/{id}` | borrar | 204 / 400 / 404 |

Health: `GET /healthz` (liveness), `GET /readyz` (readiness, hace ping a la BD).
Errores con formato estable:

```json
{ "error": { "code": "NOT_FOUND", "message": "todo not found" } }
```

## Herramientas (fijadas como go tools, sin instalación global)

```powershell
go tool sqlc generate   # regenera internal/db desde db/queries
go tool goose ...       # migraciones (up / down / status)
go vet ./...
```

## Configuración (backend/.env)

```
ADDR=:8080
DATABASE_URL=postgres://todo:todo@localhost:5433/todos?sslmode=disable
FRONTEND_URL=http://localhost:3000   # origen CORS permitido
```

## Capas (dependencias hacia adentro)

```
handler (HTTP)  →  service (lógica)  →  domain (entidades y reglas)
                        ↑ consume interfaz TodoStore (definida en service)
store/postgres  ────────┘ implementación con sqlc + pgx
```

- Los errores de dominio (`ErrNotFound`, `ErrInvalidTitle`, `ErrInvalidID`) se
  mapean a códigos HTTP solo en `handler`.
- Queries en SQL puro en `db/queries`; sqlc genera Go tipado en `internal/db`.
