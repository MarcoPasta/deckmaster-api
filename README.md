# deckmaster

## Development: Docker Compose & .env

- **Copy example env:** Create a local `.env` from the safe example before running compose:

```bash
cp .env.example .env

# Edit `.env` to set any secrets or custom ports
```

- **Build & run (foreground):**

```bash
docker compose up --build
```

- **Build & run (detached):**

```bash
docker compose up --build -d
```

- **Stop & remove containers:**

```bash
docker compose down
```

- **Notes:**
  - The repository includes `.env.example` (safe to commit). Keep your real `.env` private — `.gitignore` already excludes `.env`.
  - `docker-compose.yml` reads the project `.env` for variable substitution and passes runtime variables to the `app` service using `env_file: .env`.
  - Confirm environment inside containers if needed:

```bash
docker compose exec app env | grep DB_
docker compose exec db env | grep POSTGRES
```
