# deckmaster

important notes

project domain: deckmaster.de

---

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

## Dev workflow (Dev Container & Compose override)

You can use a Dev Container (VS Code) or a Compose override for a development container
that has a shell and mounts the repository.

- Use the Dev Container: open the repo in VS Code and choose "Reopen in Container".
  The Dev Container uses `.devcontainer/devcontainer.json` and builds the `dev` service.

- Or use the Compose override locally to mount the workspace and start the database:

```bash
# Start the dev environment and database together
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build dev db

# Stop
docker compose -f docker-compose.yml -f docker-compose.dev.yml down
```

- The override file is `docker-compose.dev.yml` and ensures the `dev` service mounts
  the repo into `/workspace` (hot-reload your edits).

Notes:

- The Dev Container runs as the `dev` user (UID 1000) so files are owned correctly inside the container.
- `postCreateCommand` copies `.env.example` to `.env` if missing; still keep secrets out of the repo.
