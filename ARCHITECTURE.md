
# Valora OSS — Architecture

## Overview

Valora OSS is an open-source **billing engine** focused on **billing logic correctness**.

It is designed to extract billing concerns—usage metering, pricing, rating, subscriptions, and invoice generation—out of application code and into a **dedicated, deterministic engine**.

Valora OSS is intentionally **not** a payment processor.
It computes *what should be billed*, not *how money is collected*.

This document describes how Valora OSS is structured, where its boundaries are, and why those boundaries exist.

---

## What Valora OSS Is (and Is Not)

Valora OSS **is**:

- A billing computation engine
- A control plane for pricing and usage-based billing
- A deterministic system for producing invoices and billing states

Valora OSS **is not**:

- A payment gateway
- A merchant of record
- A financial settlement system
- An infrastructure or operations platform

This distinction is fundamental to every architectural decision in Valora.

---

## Architectural Intent

The architecture of Valora OSS is driven by three core intents:

1. **Billing logic must be isolated**Billing rules change faster than product logic. Hardcoding them into application flows creates long-term risk.
2. **Billing must be deterministic**Given the same inputs—usage, pricing configuration, subscription state, and time—Valora must always produce the same result.
3. **Trust boundaries must be explicit**
   High-risk domains (payments, credentials, infrastructure) are deliberately kept outside the engine.

---

## High-Level Structure

At a high level, Valora OSS sits between the application and external financial systems.

The application:

- Owns users and authentication
- Emits usage
- Executes payments

Valora OSS:

- Validates and aggregates usage
- Applies pricing and rating rules
- Manages subscription and invoice state

External systems:

- Execute payments
- Handle identity
- Perform accounting or reconciliation

Valora OSS never crosses into payment execution.

---

## Core Runtime Flow

A typical billing lifecycle in Valora OSS looks like this:

1. The application defines billing primitives:

   - Products
   - Meters
   - Prices (flat, usage-based, tiered, hybrid)
   - Subscriptions
2. The application sends usage events to Valora OSS.
3. Valora OSS:

   - Validates usage against meter definitions
   - Aggregates usage per billing period
   - Applies pricing and rating logic
   - Computes billable line items
4. At the billing boundary:

   - An invoice is generated
   - Totals and line items are finalized
   - Invoice state transitions are recorded
5. The application consumes the invoice output and performs payment execution externally.

At no point does Valora OSS store or process payment credentials.

---

## Internal Composition

Internally, Valora OSS is structured in layered responsibilities:

- **API layer**

  - Handles HTTP/gRPC transport
  - Enforces authentication and authorization
  - Performs request validation
  - Contains no billing logic
- **Application / service layer**

  - Orchestrates billing workflows
  - Coordinates domain operations
  - Enforces state transitions and invariants
- **Domain layer**

  - Pricing models and rating logic
  - Subscription lifecycle rules
  - Invoice state machines
  - Designed to be deterministic and transport-agnostic
- **Persistence layer**

  - Stores usage records, pricing configuration, subscription state, and invoices
  - Abstracted behind repositories
  - Assumes database security is handled by the deployment environment

This separation exists to keep billing logic testable, auditable, and transferable.

---

## Trust Boundaries

Valora OSS operates across clearly defined trust boundaries.

```mermaid
flowchart TB
    subgraph Client["Client / Application"]
        C1["Authenticates users"]
        C2["Emits usage events"]
        C3["Manages customers"]
    end

    subgraph Valora["Valora OSS"]
        V1["Metering"]
        V2["Pricing & Rating"]
        V3["Subscriptions"]
        V4["Invoices"]
        V5["Logical tenant isolation"]
    end

    subgraph DB["Database"]
        D1["Usage records"]
        D2["Pricing configuration"]
        D3["Invoice state"]
    end

    subgraph External["External Systems (Out of Scope)"]
        E1["Payment providers"]
        E2["Identity providers"]
        E3["Accounting / ERP systems"]
    end

    Client -->|Authenticated API requests| Valora
    Valora -->|Internal persistence| DB
    Valora -.->|References only| External
```
