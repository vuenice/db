# ChatDB

ChatDB is a **single-binary** database viewer. Download or build one executable, run it, and open your browser — no server setup required.

- Supports **PostgreSQL** and **MySQL / MariaDB**
- Embedded web UI (Vue 3 SPA) — nothing extra to install for the frontend
- Saves connection credentials **encrypted at rest** in a local SQLite file

## Features

- **Auth**: register and log in with a local account; sessions use JWT
- **Connection registry**: save multiple DB connections with encrypted credentials
- **Browse catalog**: databases, tables, columns, indexes, and paginated row previews
- **Run SQL**: execute queries against read or write pools and cancel in-flight runs
- **Edit data**: update individual rows in the UI
- **DB operations**: truncate, delete, rename a database; import / export via **Workbench → Operations**
  - PostgreSQL: `pg_dump` / `psql` / `pg_restore` (requires `postgresql-client` on the same machine as the ChatDB binary)
  - MySQL: upload import; export returns a table listing
- **Bulk table operations**: drop / truncate / analyze / optimize / repair / check across multiple tables at once

## Quickstart

### Download a release binary

1. Download the binary for your OS from the Releases page.
2. Run it:

   **Windows (PowerShell)**
   ```powershell
   .\chatdb.exe
   ```

   **Linux / macOS**
   ```bash
   ./chatdb
   ```

3. Open `http://127.0.0.1:6366` in your browser.
4. Register an account and add your first database connection.

### Build from source

Prereqs: [Go](https://go.dev/dl/) and [Node.js](https://nodejs.org/) installed.

#### Windows (PowerShell)

Open PowerShell in the repo root and run the steps below. They replicate the Linux `make build` target without needing Make.

**Step 1 — build the frontend**

```powershell
Set-Location frontend
npm install
npm run build
Set-Location ..

# Copy the built SPA into the Go embed directory
Remove-Item -Recurse -Force backend\web\dist -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force backend\web\dist | Out-Null
Copy-Item -Recurse frontend\dist\* backend\web\dist\
```

**Step 2 — compile the Go binary**

```powershell
Set-Location backend
$env:CGO_ENABLED = "0"
go build -trimpath -ldflags="-s -w" -o ..\chatdb.exe .\cmd\chatdb
Set-Location ..
```

**Step 3 — run it**

```powershell
.\chatdb.exe
```

Open `http://127.0.0.1:6366`.

#### Linux / macOS

```bash
make build
./chatdb
```

## Configuration

On first start, ChatDB creates **`chatdb.config.json`** and **`chatdb.meta.sqlite`** in the OS user config directory:

| OS | Default path |
|----|--------------|
| Windows | `%APPDATA%\chatdb\` |
| Linux | `$XDG_CONFIG_HOME/chatdb/` (usually `~/.config/chatdb/`) |
| macOS | `~/Library/Application Support/chatdb/` |

Defaults: listen address `127.0.0.1:6366`, randomly generated `jwt_secret` and `app_key`. Edit the JSON file there to change the listen address or rotate secrets. You do not need to write a config file manually before first run.

## Security notes

- `app_key` must be **exactly 32 bytes** — it is the AES-256 key used to encrypt stored DB passwords.
- Keep `chatdb.config.json` and `chatdb.meta.sqlite` private (created with `0600` permissions where supported).
- Rotate `jwt_secret` to invalidate all existing sessions.

## Known limitations

- **One connection per user** is enforced in the current build.
- PostgreSQL full dump / restore requires `pg_dump`, `psql`, and `pg_restore` installed on the same machine that runs the ChatDB binary (e.g. the `postgresql-client` package). Large imports may time out if ChatDB sits behind a reverse proxy with short timeouts.
- Bulk operations (truncate, analyze, etc.) are dialect-sensitive; some edge cases may need SQL adjustments for strict Postgres / MySQL compatibility.

## License

Add your preferred license (MIT / Apache-2.0 / etc.) and a `LICENSE` file at the repo root.
