# Project Plan: Mini Go Load Balancer (GitHub Showcase Project)

## Context
I'm learning Go and GCP and building a portfolio project to showcase in job searches. I want it **mini and finishable in limited time**, but polished enough that others could realistically clone and run it. Prioritize a working, well-documented core over broad scope. If anything below is ambiguous, ask me before proceeding.

## How I want to work (read this first — it overrides normal "just do it" helpfulness)
**The end goal is me learning Go.** Do not write the implementation code for me. This is a tutorial/mentorship session, not a code-generation session.

- **Teach step by step.** Break each feature into a small, digestible step. Explain the *concept* first (e.g., what a health check loop is, why we'd use a channel here) like a tutorial, before any code is involved.
- **TDD workflow, strictly.** For each feature:
  1. We discuss the concept/requirement briefly.
  2. You write (or guide me to write) a small test first that captures the desired behavior.
  3. I attempt to implement the feature myself to make the test pass.
  4. I report back what I tried/what happened; you review, give feedback, and only then help debug or explain further.
- **Don't skip ahead to full solutions.** If I'm stuck, give hints and nudges first (e.g., "think about what happens when two goroutines increment this counter at once") rather than the answer. Only give full code examples if I explicitly ask for advice or an example, or if I'm stuck after a couple of hint rounds.
- **Small steps over big dumps.** Prefer many small back-and-forth exchanges over one large response that implements a whole package. One test/feature at a time.
- **It's fine to be slow.** Time-efficiency of the *coding* is secondary to me actually understanding what I write. The "mini" scope exists so the learning stays manageable, not so the process gets rushed.
- Exception: scaffolding that isn't the learning target (e.g., docker-compose fake backends, boilerplate config file structure) can be provided directly rather than taught step by step — use judgment on what's "core Go/concurrency learning" vs. "supporting setup."

## Goal
Build a small, production-inspired HTTP load balancer in Go.

## Core Requirements (must ship)
- Two routing algorithms: round-robin and least-connections.
- Config loaded from a YAML file and/or environment variables (include a working `config.example.yaml` and `.env.example`).
- Active health checks: periodic `/healthz` polling, unhealthy backends excluded from rotation.
- Passive health checks: track forwarding failures per backend, trip a backend out of rotation after N consecutive failures (simple circuit breaker, not a library).
- Dynamic service discovery: backends are added/removed at runtime via a `BackendProvider` abstraction (see `architecture.md`) rather than a fixed static list. Ship a `StaticBackendProvider` for the local docker-compose demo, and a `GCPMIGBackendProvider` for the GCP deployment path (polling GCE MIG membership).
- Goroutines/channels for concurrency: health checks, request forwarding, metrics collection.
- Reverse proxy behavior: forward requests, preserve headers, timeouts, and a single one-shot retry on connection error (no general retry policy).
- `/metrics` endpoint — prefer Prometheus text format if low extra effort, otherwise plain JSON/text with backend health + request counts.
- Structured JSON logging to stdout (Cloud Logging-compatible by virtue of stdout ingestion — no GCP logging client needed), including per-request latency and which backend served it.
- Graceful shutdown.

## Demo / Usability (must ship — this is what makes it useful to others)
- `Dockerfile` for the load balancer.
- `docker-compose.yml` that runs the load balancer plus 2–3 fake backend services (simple Go echo servers or nginx) so someone can `docker-compose up` and see round-robin/least-connections working immediately via curl.
- README with quickstart instructions that get a stranger from `git clone` to a working demo in under a few minutes.

## GCP Integration (documented, not necessarily deployed)
- README section describing deployment: load balancer on a Compute Engine VM, backends in a Compute Engine **Managed Instance Group (MIG)**, discovered dynamically via `GCPMIGBackendProvider` (see `architecture.md`). MIG is chosen deliberately over Cloud Run: Cloud Run already load-balances across its own instances, so fronting it with a custom LB doesn't demonstrate anything — MIGs are the case where a hand-built load balancer is actually solving a real gap.
- Explain Cloud Logging compatibility (stdout JSON logs, no special client).
- Include k6 load-testing instructions runnable from Cloud Shell.
- Do NOT build GCP deployment automation or actually stand up infrastructure unless I explicitly ask — describing the steps clearly is sufficient for this pass. If time allows later, I may ask you to actually deploy it and capture real output.

## Testing
- Unit tests for routing algorithms, health check logic, and config loading — written **before** the implementation, per the TDD workflow above.
- A handful of integration-style tests using `httptest` (spin up fake backends in-process, verify proxy behavior). No docker-based test harness.
- Tests double as the teaching tool: a well-chosen test is how a concept gets introduced before I write code for it.

## Deliverables
1. Architecture: see `architecture.md` for the full design (layering, interfaces, concurrency model, GCP isolation). This plan does not repeat it — implementation should follow `architecture.md`, which takes precedence over this document if the two ever disagree on technical approach.
2. Mermaid or ASCII diagrams (request flow, health check loop) — supplementary to `architecture.md`, not a replacement.
3. Complete, idiomatic, documented Go codebase.
4. Dockerfile + docker-compose demo setup with fake backends.
5. README: quickstart, config reference, GCP deployment steps, k6 testing instructions.
6. Tests (unit + light integration).
7. A "Future Enhancements" section in the README listing things intentionally NOT built: TLS termination, weighted round-robin / IP hash, sticky sessions, rate limiting, admin UI, retry policies beyond one-shot.

## Explicit Non-Goals
- No production-grade features beyond what's listed above.
- No TLS termination.
- No general/configurable retry policy — one-shot retry only.
- No actual GCP infrastructure deployment unless separately requested.
- No third-party load balancing libraries — routing/health logic should be hand-written to demonstrate Go concurrency skills.

## Build order

1. Set up project scaffolding and Go module init
2. Define core interfaces (Backend, RoutingStrategy, BackendConnector,
BackendProvider, HealthChecker)
3. Write unit tests for config loading
4. Implement config loading (YAML + env vars)
5. Implement StaticBackendProvider (backend registry)
6. Implement BackendStats + StatsReporter/StatsReader (atomic, sync.Map-backed)
7. Write integration tests with httptest for reverse proxy behavior
8. Implement reverse proxy core (forwarding, header preservation, timeouts)
9. Implement one-shot retry on connection error
10. Write unit tests for routing algorithms
11. Implement round-robin routing strategy
12. Implement least-connections routing strategy
13. Write unit tests for health check logic
14. Implement active health checks (/healthz polling)
15. Implement passive health checks (consecutive-failure circuit breaker)
16. Implement /metrics endpoint (Prometheus or JSON/text)
17. Implement structured JSON logging (stdout, per-request latency + 18. backend served)
18. Implement graceful shutdown
19. Implement GCPMIGBackendProvider (dynamic GCP discovery)
20. Implement GCPManagedInstanceConnector
21. Create Dockerfile for the load balancer
22. Create docker-compose demo with fake backend services
23. Write README (quickstart, config reference)
24. Write README GCP deployment section
25. Write README k6 load-testing instructions
26. Write README "Future Enhancements" section

## How to proceed
1. For each issue: explain the concept → help me write a failing test (if applicable) → let me implement → review what I did → move to the next small piece. Do not jump ahead to the next feature until the current one works and I understand it.
2. Ask me before adding anything not listed above, even if it seems like a natural extension.