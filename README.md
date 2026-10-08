# VueNiceDB

Need a lightweight, platform-independent PostgreSQL or MySQL DB viewer? CRUD and visualize your database like never before. Host it on your production server as a lightweight web server or run it as a native desktop application on your local machine. Built with **Golang** and **Vue.js**.

---

## 🚀 How to Run (2 Ways)

You can run VueNiceDB either as a **Web Server** (ideal for teams, VPS, Docker, or production servers) or as a **Desktop App** (ideal for local development). Both options offer a **quick pre-built method** (non-technical) and a **build-from-source method** (technical).

---

### 1. As a Web Server

#### Option A: Quick Start (Pre-built Binary — Non-Technical)
*No compiler or runtime required.*

1. **Download the binary**: Grab the latest `chatdb` (Linux/macOS) or `chatdb.exe` (Windows) binary from [GitHub Releases](https://github.com/vuenice/db/releases).
2. **Create a `.env` file** in the same directory:
   ```env
   PORT=6366
   AUTH_API_KEY=your_secret_key_here
   ```
3. **Run the server**:
   ```bash
   chmod +x ./chatdb
   ./chatdb
   ```
4. Open your browser at `http://localhost:6366`.

5. **(Optional) Run with Supervisor** (e.g. `/etc/supervisor/conf.d/vuenicedb.conf`):
   ```ini
   [program:vuenicedb]
   command=/path/to/chatdb
   directory=/path/to/db
   autostart=true
   autorestart=true
   user=www-data
   redirect_stderr=true
   stdout_logfile=/var/log/vuenicedb.log
   ```
   Reload supervisor to start:
   ```bash
   supervisorctl reread
   supervisorctl update
   supervisorctl start vuenicedb
   ```

#### Option B: Build from Source (Technical)
*Requires **Golang (1.23+)** and **Node.js (18+)** installed on your system.*

1. **Clone the repository**:
   ```bash
   git clone https://github.com/vuenice/db.git
   cd db
   ```
2. **Configure `.env`**:
   ```env
   PORT=6366
   AUTH_API_KEY=your_secret_key_here
   ```
3. **Build the web server and single binary**:
   ```bash
   make build
   ```
   *(This builds the Vue frontend into `backend/web/dist` and compiles a static Go binary `./chatdb`)*
4. **Start the web server**:
   ```bash
   ./chatdb
   ```

---

### 2. As an App

#### Option A: Pre-built Desktop App (Non-Technical)
*No development tools required.*

Download the installer or standalone executable for your operating system from [GitHub Releases](https://github.com/vuenice/db/releases):
- **Linux**: `.AppImage` or `.deb`
- **macOS**: `.dmg` or `.app`
- **Windows**: `.exe` installer

Run the downloaded file to launch VueNiceDB immediately.

#### Option B: Bundle from Source (Technical)
*Requires **Golang (1.23+)**, **Node.js (18+)**, and [Wails v3](https://v3.wails.io/) installed.*

1. **Clone the repository**:
   ```bash
   git clone https://github.com/vuenice/db.git
   cd db
   ```
2. **Install frontend dependencies**:
   ```bash
   cd frontend && npm install && npm run build && cd ..
   ```
3. **Bundle native executable for your OS**:
   ```bash
   task build:server
   # or using wails CLI directly:
   cd backend && wails3 build
   ```
   The bundled native application will be generated in `backend/bin/`.

---

## 📖 How to Use?

Here is a step-by-step walkthrough of all features and user flows inside VueNiceDB:

### Step 1. Register (Admin & First Connection)
When launching VueNiceDB for the first time, you are directed to the **Register** page to initialize your administrator account and primary database connection simultaneously:
- **Driver Selection**: Select either **PostgreSQL** (default port `5432`) or **MySQL / MariaDB** (default port `3306`).
- **Connection Label**: A distinct alias for this database connection (e.g. `Production DB` or `Analytics Store`).
- **Host & Port**: Database hostname (e.g. `127.0.0.1`, `localhost`, or remote IP) and port.
- **Database Name**: The default database to connect to.
- **SSL Mode**: Configure SSL encryption (`disable`, `require`, `verify-ca`, or `verify-full`).
- **Database Credentials**: Username and password. The password is encrypted with AES-256 derived from your `AUTH_API_KEY`.
- **Optional SSH Tunnel**: If your database is inside a private VPC or behind a firewall, enable **Use SSH Tunnel** and provide the SSH host, port (`22`), SSH user, and authentication method (Password or Private Key).
- Submitting the form automatically signs you in and opens the database workbench.

### Step 2. Login & Connection Selection
On subsequent visits, access the login page at `/login`:
- **Connection Label Selector**: VueNiceDB retrieves all registered database connections. Select the target database connection label from the dropdown.
- **Credentials**: Enter your username and password.
- Upon successful authentication, a secure JWT session token is generated and stored locally.
- *Tip*: If no connections or users exist yet, VueNiceDB automatically redirects you to the Register page.

### Step 3. Search & View Tables
Navigate the database structure and preview data through the main workbench:
- **Database & Schema Switcher**: Easily switch physical databases on the fly using the top-bar database selector, or select a specific schema (such as `public`).
- **Sidebar Table Catalog**: All tables and views are listed in the left sidebar alphabetically grouped (A-Z) with live search filtering.
- **Table Tabs**:
  - **Structure**: View all columns, data types, nullability flags, default values, and primary key indicators.
  - **Data**: Paginated data grid displaying table rows. Supports column sorting, row limit controls, and column shifting (move columns left/right).
  - **Indexes**: Displays all defined indexes and their SQL definitions.
- **Multi-Table Conditional Search**: Click the **Search** icon in the sidebar to open the advanced search workspace. Filter records across multiple tables using conditional operators (`=`, `!=`, `LIKE`, `>`, `<`).

### Step 4. Edit Row / Field Safely
Perform in-place data modifications without writing manual SQL:
1. In the **Data** tab, right-click on any row or cell to open the context menu.
2. Select:
   - **Update field**: Directly modify the clicked column value.
   - **Update record**: Open the full row editor dialog showing all columns and values.
3. **Live SQL Preview**: VueNiceDB displays a live preview of the exact SQL `UPDATE` statement that will be executed, with a strict `WHERE` clause pinned to the row's original snapshot to prevent accidental overwrites.
4. Confirm to execute the update against the write pool. The table data refreshes automatically.

### Step 5. Check Indexes & Constraints
Inspect table indexes and constraints to verify database performance:
1. Select any table from the sidebar.
2. Click the **Indexes** tab in the main view.
3. Review the index list, including primary keys, unique constraints, and secondary B-Tree/GIN indexes, along with their raw `CREATE INDEX` SQL definitions.

### Step 6. Connect with Model Context Protocol (MCP)
VueNiceDB includes built-in support for the **Model Context Protocol (MCP)**, allowing AI assistants like **Claude Desktop**, **Cursor**, and **Continue** to query and inspect your databases directly:
- Click the **"Connect MCP"** button in the sidebar to open the configuration modal.
- **Stdio Mode (Claude Desktop / Cursor)**:
  Add the following to your `claude_desktop_config.json`:
  ```json
  {
    "mcpServers": {
      "vuenicedb": {
        "command": "/path/to/chatdb",
        "args": ["mcp"],
        "env": {
          "AUTH_API_KEY": "your_secret_key_here"
        }
      }
    }
  }
  ```
- **HTTP Mode**:
  AI tools can also communicate over HTTP via `POST /api/connections/{id}/mcp` using JSON-RPC 2.0 and Bearer token authentication.
- **Available MCP Tools**:
  - `list_tables`: List all tables and views in the schema.
  - `describe_table`: Retrieve column details, data types, nullability, and indexes for a specific table.
  - `read_query`: Execute safe, read-only SQL queries (`SELECT`, `SHOW`, `EXPLAIN`).
  - `list_databases`: List all available databases.

### Step 7. Create Pages and Save Queries
Build custom SQL query pages and organize reusable queries:
- **Chat SQL Workbench**: Navigate to **Queries** -> **Chat SQL** to open the interactive SQL editor.
- **Pool Mode**: Toggle between **Read Pool** (safe read queries) and **Write Pool** (for DDL/DML changes).
- **Execution Metrics & CSV Export**: Run your query with `Ctrl + Enter` (or click **Run Query**). View execution duration in seconds, returned row counts, and export results directly to CSV.
- **Save Query Page**: Click **Save Page** in the editor toolbar, enter a custom title (e.g., `Monthly Active Users Report`), and persist the query.
- **Saved Queries Library**: Switch to **Queries** -> **Saved** to browse all saved query pages in an organized grid. Click any card to load the query into the editor.
- **+ New Query Page**: Click the **+ New Query Page** button in the Saved tab to start a new blank query workflow.
- **Execution History**: Open the **History** tab in the sidebar to review all past SQL executions with timestamps, duration, and rerun options.

---

## 🤝 How to Contribute?

We welcome contributions from developers! Below are the exact, tested steps to get VueNiceDB set up and running locally on Linux / macOS.

---

### Part 1: Setup

#### Prerequisites
- **Golang**: `1.23` or higher ([go.dev/dl](https://go.dev/dl/))
  - *Linux (Ubuntu/Debian)*: `sudo apt install golang` or extract Go tarball to `~/.local/go`
  - *macOS*: `brew install go`
- **Node.js**: `18.x` or higher with `npm` ([nodejs.org](https://nodejs.org/))
  - *Linux (Ubuntu/Debian)*: `sudo apt install nodejs npm` (or via `nvm`)
  - *macOS*: `brew install node`
- **Git**
- *(Optional)* **Docker & Docker Compose** (for running sample PostgreSQL test databases)

#### 1. Fork & Clone the Repository
```bash
git clone https://github.com/<your-username>/db.git
cd db
```

#### 2. Configure Environment (`.env`)
Create a `.env` file at the root of the project:
```bash
cat << 'EOF' > .env
PORT=6366
AUTH_API_KEY="aHR0cHM6Ly9hdXRoLWNvbmZpcm0tbmF2eS52ZXJjZWwuYXBwL2FwaQ=="
EOF
```
*(Or specify your preferred port and secret key)*

#### 3. Install Frontend Dependencies
```bash
cd frontend
npm install
cd ..
```

#### 4. Download Backend Dependencies
```bash
cd backend
go mod download
cd ..
```

---

### Part 2: Run

Choose how you want to run VueNiceDB locally:

#### Option A: Run in Development Mode (Live Hot-Reloading)
Start both the backend API server and the Vite frontend dev server with hot-reloading:

- **Using the convenience script**:
  ```bash
  chmod +x ./dev.sh
  ./dev.sh
  ```

- **Or manually in two terminals**:
  - **Terminal 1 (Backend API)**:
    ```bash
    cd backend
    go run ./cmd/chatdb
    ```
    *Starts on `http://127.0.0.1:6366`.*
  - **Terminal 2 (Frontend SPA)**:
    ```bash
    cd frontend
    npm run dev
    ```
    *Starts on `http://localhost:5173` with live reloading and proxying API calls to port 6366.*

- Open **`http://localhost:5173`** in your browser.

---

#### Option B: Build & Run Production Web Server (Single Binary)
Build the Vue SPA and embed it directly into a single self-contained Go executable:

1. **Build SPA and compile binary**:
   ```bash
   make build
   ```
   *Under the hood, this:*
   - Compiles the Vue frontend (`npm run build`).
   - Copies production assets into `backend/web/dist` for Go `//go:embed`.
   - Compiles a static Linux/macOS binary at `./chatdb` with `CGO_ENABLED=0`.

2. **Run the web server**:
   ```bash
   ./chatdb
   ```
3. Open **`http://127.0.0.1:6366`** in your browser.

---

#### Option C: (Optional) Run Desktop App on Linux / macOS
To test or build VueNiceDB as a native desktop application using Wails:

- **System prerequisites**:
  - *Linux (Ubuntu/Debian)*: `sudo apt install libgtk-4-dev libwebkitgtk-6.0-dev`
  - *macOS*: `xcode-select --install`
- **Install Wails v3 CLI**:
  ```bash
  go install github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```
- **Run desktop app in development**:
  ```bash
  cd backend
  go run -tags desktop ./cmd/chatdb
  ```
- **Bundle native executable**:
  ```bash
  cd backend
  wails3 build
  ```

---

#### Option D: (Optional) Run Test Databases with Docker
Start pre-configured PostgreSQL test databases for local querying:
```bash
docker compose up -d
```
- **App DB**: `localhost:5433` (Database: `chatdb_app`, User: `chatdb`, Password: `chatdb`)
- **Sample DB**: `localhost:5434` (Database: `sample`, User: `postgres`, Password: `postgres` with sample tables)

---

### Contributing Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature-name
   ```
2. Make your improvements and test them (`./dev.sh` or `make build`).
3. Commit your changes with clear messages:
   ```bash
   git commit -m "feat: describe your change"
   ```
4. Push to your fork:
   ```bash
   git push origin feature/my-feature-name
   ```
5. Open a **Pull Request** to the `main` branch on GitHub.

For bug reports or feature requests, feel free to open an issue on [GitHub Issues](https://github.com/vuenice/db/issues).

---

## 📬 Contact & Feedback

Contact for any feedback: [mr.yogesh.galav@gmail.com](mailto:mr.yogesh.galav@gmail.com)

---

## 👥 Contributors

- [Yogesh Galav](https://github.com/vuenice)
- [All Contributors](https://github.com/vuenice/db/graphs/contributors)

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
