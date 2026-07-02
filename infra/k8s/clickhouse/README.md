# ClickHouse & Tabix Web UI Deployment

This directory contains the manifests to deploy **ClickHouse Server** and **Tabix Web client** (SQL GUI console) inside the `observability` namespace.

## Deployment Steps

1. **Deploy manifests:**
   ```bash
   kubectl apply -k .
   ```

2. **Verify installation:**
   ```bash
   kubectl get pods -n observability
   ```

---

## Accessing the User Interfaces

You have two choices to interact with ClickHouse locally:

### Option A: Built-in ClickHouse Playground (Play UI)
ClickHouse has a built-in lightweight SQL playground UI accessible directly on the HTTP server port:

1. Port-forward the ClickHouse HTTP port:
   ```bash
   kubectl port-forward svc/clickhouse -n observability 8123:8123
   ```
2. Open your browser and go to:
   [http://localhost:8123/play](http://localhost:8123/play)
3. Credentials:
   - **User:** `default`
   - **Password:** `lokiforce123`

### Option B: Tabix Dashboard (Advanced SQL Editor)
Tabix is a rich open-source web-based developer console and dashboard for ClickHouse:

1. Port-forward the Tabix port:
   ```bash
   kubectl port-forward svc/tabix -n observability 8080:80
   ```
2. Open your browser and go to:
   [http://localhost:8080](http://localhost:8080)
3. Connect Tabix to ClickHouse:
   - **ClickHouse host:** `http://localhost:8123` (Make sure the ClickHouse port-forward on `8123` is running!)
   - **Login:** `default`
   - **Password:** `lokiforce123`
