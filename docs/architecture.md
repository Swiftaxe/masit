# Load Balancer Architecture Plan

**Language:** Go
**Design goals:** Decoupling, Single Responsibility Principle, high-traffic support, dynamic backend scaling (GCP MIGs / K8s), testability

---

## 1. Layered Architecture

```
┌─────────────────────────────────────┐
│         Listener / Entry Point       │  accepts client connections
└───────────────┬───────────────────────┘
                │
┌───────────────▼───────────────────────┐
│         Routing Layer (Strategy)       │  "which backend?"
│   - RoundRobin / LeastConn / Random    │
│   - Reads BackendStats (passive)       │
│   - Knows nothing about protocols       │
└───────────────┬───────────────────────┘
                │ selects Backend
┌───────────────▼───────────────────────┐
│      Backend Registry / Pool           │  holds Backend + BackendStats
└───────────────┬───────────────────────┘
                │
┌───────────────▼───────────────────────┐
│    Connector / Adapter Layer           │  "how do I reach it?"
│   - TCPConnector, HTTPConnector,       │
│     GCPManagedInstanceConnector        │
│   - Reports events as they happen      │
└───────────────┬───────────────────────┘
                │
         Actual Backend (VM, pod, GCP MIG instance...)
```

Each layer has exactly one responsibility and depends only on interfaces, not concrete implementations of neighboring layers.

---

## 2. Core Interfaces

### Backend (data, not behavior)
```go
type Backend struct {
    ID       string
    Weight   int
    Metadata map[string]string // opaque to router
}
```

### Routing Strategy
```go
type RoutingStrategy interface {
    SelectBackend(pool []Backend, ctx RequestContext) Backend
}
```
Implementations: `RoundRobinStrategy`, `LeastConnectionsStrategy`, `WeightedRandomStrategy`

### Backend Connector
```go
type BackendConnector interface {
    Connect(backend Backend) (Connection, error)
    IsHealthy(backend Backend) bool
}
```
Implementations: `TCPConnector`, `HTTPConnector`, `GCPManagedInstanceConnector`

### Backend Provider (service discovery)
```go
type BackendProvider interface {
    ListBackends() ([]Backend, error)
}
```
Implementations: `StaticBackendProvider`, `GCPMIGBackendProvider`, `K8sServiceProvider`

### Health Checker
```go
type HealthChecker interface {
    Check(backend Backend, connector BackendConnector) HealthStatus
}
```

---

## 3. Stats Layer (Event-Driven, Atomic)

**BackendStats** — pure data, no behavior:
```go
type BackendStats struct {
    ActiveConnections int64
    TotalConnections  int64
    ErrorCount        int64
    AvgLatencyMs      float64
    LastUpdated       int64
}
```

Split into two narrow interfaces (Interface Segregation):

```go
// Write side — used only by Connectors
type StatsReporter interface {
    RecordConnectionOpened(backendId string)
    RecordConnectionClosed(backendId string)
    RecordLatency(backendId string, ms int64)
    RecordError(backendId string)
}

// Read side — used only by Routing Strategies
type StatsReader interface {
    GetStats(backendId string) BackendStats
}
```

A single `InMemoryStatsStore` implements both; each consumer only sees its narrowed interface.

### Concurrency model
- Backed by `sync.Map` (backendId → `*statsEntry`), chosen specifically because backend counts change dynamically under autoscaling — no stable numeric index can be assumed.
- All counters mutated via `atomic.AddInt64` / `atomic.LoadInt64` — no mutexes, no lock contention on the hot path.
- Latency stored as `sum` + `count` (both atomic), with average computed at read time — avoids read-modify-write races on a float.
- `GetStats` returns a "close enough, fast" snapshot, not a perfectly consistent one — acceptable for routing decisions.

### Backend lifecycle under dynamic scaling
- **Added:** stats lazily created via `LoadOrStore` on first event — no explicit registration step needed.
- **Removed:** `RemoveBackend(id)` deletes the entry from `sync.Map` to prevent unbounded growth as instances are recycled.

---

## 4. Design Principles Applied

| Principle | How it's applied |
|---|---|
| **SRP** | Routing, connecting, discovery, stats, health checks are all separate interfaces/types |
| **Decoupling** | Router never knows backend transport type; Connector never knows routing logic |
| **Interface Segregation** | `StatsReporter` vs `StatsReader` prevent writers from reading and vice versa |
| **Testability** | Routing strategies testable with mock `StatsReader`, no live network/GCP needed |
| **Dependency Injection** | Plain Go idiom — interfaces passed into constructors, no DI framework |

---

## 5. GCP-Specific Isolation

All GCP knowledge is fully contained in two implementations — nothing else in the system needs to change to support GCP:

- `GCPMIGBackendProvider` — discovers backends via GCP Managed Instance Group APIs
- `GCPManagedInstanceConnector` — resolves and connects to GCP instances

Both use `cloud.google.com/go/compute` and thread `context.Context` for cancellation/timeouts, per idiomatic Go + GCP SDK conventions.

---

## 6. Open / Next Steps

- [ ] Design `HealthChecker` (likely goroutine-per-backend + ticker, fan-in via channels)
- [ ] Define `BackendProvider` → `Registry` sync mechanism (polling vs. push/webhook from GCP)
- [ ] Implement `GCPMIGBackendProvider` and `GCPManagedInstanceConnector`
- [ ] Decide `context.Context` propagation conventions across all layers