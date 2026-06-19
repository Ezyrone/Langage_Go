# URLWatch

Microservice Go de vérification d'URLs en masse. Envoie un lot d'URLs, les vérifie en parallèle (code HTTP, latence, disponibilité) et expose les résultats via une API REST.

## Build

```bash
go build ./...
```

## Run

```bash
go run ./cmd/urlwatch
```

Variables d'environnement :
- `ADDR` : adresse d'écoute (défaut `:8080`)
- `LOG_LEVEL` : niveau de log (`debug`, `info`, `warn`, `error` ; défaut `info`)

## Test

```bash
go test ./...
go test -race ./...
```

## Exemples curl

### Health check

```bash
curl http://localhost:8080/healthz
```

### Créer un lot de vérifications

```bash
curl -X POST http://localhost:8080/v1/checks \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://go.dev", "https://exemple.invalid", "https://github.com"],
    "options": { "concurrency": 4, "timeout_ms": 3000 }
  }'
```

Réponse (201 Created) :
```json
{
  "batch_id": "b_4f3c1a",
  "created_at": "2026-06-18T10:00:00Z",
  "summary": { "total": 3, "up": 2, "down": 1, "duration_ms": 812 },
  "results": [
    { "url": "https://go.dev", "status_code": 200, "ok": true, "latency_ms": 120 },
    { "url": "https://exemple.invalid", "ok": false, "error": "dns: no such host", "latency_ms": 2001 },
    { "url": "https://github.com", "status_code": 200, "ok": true, "latency_ms": 230 }
  ]
}
```

### Consulter un lot existant

```bash
curl http://localhost:8080/v1/checks/b_4f3c1a
```

### Erreur : lot introuvable

```bash
curl http://localhost:8080/v1/checks/nonexistent
# 404 : {"error":{"code":"batch_not_found","message":"aucun lot avec l'id nonexistent"}}
```
