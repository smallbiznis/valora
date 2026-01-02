# Valora OSS

**Valora OSS** is an open-source **deterministic billing engine** for modern SaaS and platform products.

![Release CI](https://github.com/smallbiznis/valora/actions/workflows/github-release.yml/badge.svg)
![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)
![Release](https://img.shields.io/github/v/release/smallbiznis/valora)
![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?logo=go)

Valora extracts billing concerns—usage metering, pricing, subscriptions, and invoicing—out of application code and into a **dedicated, self-hosted engine** with explicit boundaries.

> **Valora determines what should be billed.
> It does not execute payments.**

## Why Valora Exists

Valora was created as a deliberate response to a common gap in many organizations:
the desire for innovation without sufficient space to experiment, validate,
and evolve systems responsibly.

Rather than continuously giving more within constrained structures,
Valora represents a conscious decision to establish clear professional boundaries—
to build, test, and learn in an environment that values clarity,
determinism, and accountability.

This project is developed independently, outside of employment responsibilities,
and is not intended to replicate or compete with any internal company systems.

---

## What Valora Is (and Is Not)

### Valora **is**:

- A billing computation engine
- A control plane for pricing and usage-based billing
- A deterministic system for generating invoices and billing states
- Designed for self-hosted, multi-tenant SaaS architectures

### Valora **is not**:

- A payment gateway
- A merchant of record
- A settlement or reconciliation system
- An infrastructure, hosting, or SLA platform

Payment execution is intentionally **out of scope**.

---

## Core Capabilities

### Subscription Management

- Flat-rate and usage-based subscriptions
- Trialing, active, past_due, and canceled states
- Scheduled cancellation and renewal windows
- Organization-scoped isolation

### Usage Metering & Rating

- Metered billing via usage events
- Flexible pricing models (flat, tiered, hybrid)
- Deterministic aggregation per billing period
- Idempotent usage ingestion

### Invoicing

- Automated invoice generation
- Subscription and usage-based line items
- Explicit invoice state transitions
- Proration support where applicable

### Multi-Tenancy

- Logical tenant isolation
- Organization-scoped authorization
- Designed for SaaS and platform architectures

---

## Design Principles

Valora is built around a small set of strict principles:

- **Deterministic by design**Billing outputs are derived solely from persisted inputs and configuration.
- **Explicit boundaries**Billing logic is separated from payments, identity, and infrastructure concerns.
- **Self-hosted ownership**Teams retain full control over data, logic, and deployment.
- **Composable primitives**
  Pricing behavior is modeled explicitly, not hidden in application code.

---

## Architecture & Security

Valora OSS includes first-class documentation:

- **Architecture overview and trust boundaries** → `ARCHITECTURE.md`
- **Security scope and assumptions** → `SECURITY.md`
- **Lightweight threat model** → `THREAT_MODEL.md`

These documents describe:

- deterministic billing flows
- tenant isolation models
- security responsibilities
- explicit non-goals

---

## Deployment Model

Valora OSS is **self-hosted software**.

The adopting organization is responsible for:

- infrastructure security
- TLS termination and networking
- secrets management
- database operations and backups

Valora makes minimal assumptions about the runtime environment to remain portable.

---

## Who Valora Is For

Valora is designed for teams that:

- are scaling beyond hardcoded billing logic
- operate usage-based or hybrid pricing models
- build long-lived SaaS or platform systems
- value clarity, ownership, and auditability

---

## Roadmap (High-Level)

- Harden core billing lifecycle
- Expand pricing and proration models
- Improve usage ingestion and aggregation
- Optional hosted control plane (**Valora Cloud**, future)

---

## Releases

Releases are tag-driven and backed by a curated changelog.

- Release policy and tagging flow: `RELEASE.md`
- Release notes and history: `CHANGELOG.md`

---

## Observability

Grafana dashboard expects these metrics names:

- valora_scheduler_job_runs_total
- valora_scheduler_job_duration_seconds
- valora_scheduler_job_timeouts_total
- valora_scheduler_job_errors_total
- valora_scheduler_batch_processed_total
- valora_scheduler_batch_deferred_total
- valora_scheduler_runloop_lag_seconds

---

## License

Valora OSS is open-source.
See the `LICENSE` file for details.

---

## Documentation

📚 Documentation:

- Source: `./docs/docs/index.md`

> Valora aims to make billing **boring, predictable, and explainable**—
> so teams can focus on building their products.
