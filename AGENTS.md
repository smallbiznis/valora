# **Billing-as-a-Service (Grafana-like Model)**

This document defines **agent responsibilities, boundaries, and policies** for the Billing-as-a-Service platform, inspired by Grafana OSS/Cloud design principles.

**Version:** 1.1

**Status:** Stable – Operational Clarifications

---

## 1. Philosophy

This system follows **strict separation of concerns**:

***Authentication** answers *who you are*

***Authorization** answers *what you can do*

***Policy & Governance** answer *what is recommended / restricted*

***Billing Engine** answers *how usage is measured and charged*

Like Grafana:

* Authentication is **non-blocking**
* Security policy is **contextual, not enforced at login**
* OSS and Cloud share **one codebase**, differentiated by **configuration and agents**

---

## 2. Deployment Modes

The platform operates in two explicit modes:

```text

OSS Mode     → Operator-driven, no public signup

Cloud Mode   → Multi-tenant SaaS, public signup

```

Mode is resolved at runtime:

```yaml

mode: oss | cloud

```

All agents **must respect deployment mode**.

---

## 3. Core Agents Overview

| Agent         | Responsibility       | Enforcement Level |

| ------------- | -------------------- | ----------------- |

| Auth Agent    | Identity & session   | Hard (technical)  |

| Signup Agent  | Tenant bootstrap     | Config-driven     |

| Policy Agent  | Security guidance    | Soft (advisory)   |

| Billing Agent | Usage → charge       | Hard (financial)  |

| Quota Agent   | Rate & limit control | Hard              |

| Audit Agent   | Compliance trail     | Hard              |

---

## 3.1 Agent Interaction Rules

Agents MUST remain loosely coupled and communicate only through

explicit, narrow interfaces.

The following interaction rules are **non-negotiable**:

* Auth Agent MUST NOT call Billing or Ledger
* Auth Agent MUST NOT query quota or usage state
* Policy Agent MUST NOT mutate system state
* Policy Agent MUST NOT block authentication
* Billing Agent MUST NOT inspect auth or policy annotations
* Billing Agent MUST NOT depend on UI or session context
* Quota Agent MAY block requests independently of billing
* Audit Agent observes events but NEVER influences control flow

Violating these rules introduces architectural risk and financial inconsistency.

---

## 4. Auth Agent

### Responsibility

* Authenticate users
* Issue session / token
* Manage login, logout, token rotation

### Explicit Non-Responsibilities

* ❌ Password policy enforcement
* ❌ Force password change
* ❌ Tenant lifecycle

### Behavior

* Login always returns a session **if credentials are valid**
* Session metadata may include:

```json

{

"user_id": "...",

"role": "OWNER",

"password_state": "default | rotated",

"auth_provider": "local | oauth | saml"

}

```

No login is blocked due to policy state.

---

## 5. Signup Agent

### OSS Mode

* Signup endpoints exist but are **disabled**
* Users are provisioned by:
* operator
* migration
* external IdP

### Cloud Mode

* Signup is **public and enabled**
* Signup creates:
* Tenant
* Owner user
* Trial subscription
* Default billing workspace

### Rules

* Signup creates **tenants**, not just users
* Signup cannot be called in OSS mode

---

## 6. Policy Agent

### Responsibility

* Provide security recommendations
* Annotate session with policy states

### Examples

* Default password detected
* MFA not enabled
* Billing contact not verified

### Enforcement Style

* ❌ No hard blocks at login
* ✅ Soft guidance via:
* UI banners
* warnings
* limited access to sensitive actions

Example:

```go

ifsession.PasswordState == "default" {

restrict("create_api_key")

}

```

---

## 7. Billing Agent

### Responsibility

* Meter usage
* Rate usage
* Generate charges
* Emit invoices

### Guarantees

* Idempotent processing
* Deterministic billing
* Immutable financial records

Billing **must never depend on UI or auth policy state**.

### Ledger Clarification

The Billing Agent MUST persist all financial effects

into an immutable Ledger subsystem.

Ledger characteristics:

* append-only
* deterministic
* source of truth for balances and invoices
* never recomputed from raw usage events

Billing correctness is defined by Ledger state,

not by usage logs or UI representations.

---

## 8. Quota Agent

### Responsibility

* Enforce usage limits
* Apply rate limits
* Protect system stability

### Characteristics

* Tenant-scoped
* Service-specific
* Hard enforcement

Quota violations:

* return explicit errors
* never silently fail

### Quota vs Billing Semantics

Quota enforcement occurs **BEFORE billing**.

* Quota Agent protects system stability and fairness
* Billing Agent records only what successfully passed enforcement

Exceeding quota:

* MAY block usage ingestion
* MUST return explicit errors
* MUST NOT retroactively affect billing or ledger state

---

## 9. Audit Agent

### Responsibility

* Append-only audit trail
* Security and billing events

### Rules

* No mutation of audit logs
* All timestamps UTC
* IDs use Snowflake / ULID

Audit is **mandatory in Cloud mode**.

---

## 10. OSS vs Cloud Capability Matrix

| Capability           | OSS      | Cloud       |

| -------------------- | -------- | ----------- |

| Public signup        | ❌        | ✅           |

| Multi-tenant         | ❌        | ✅           |

| Billing              | Optional | Mandatory   |

| Password enforcement | ❌        | Soft        |

| MFA                  | Optional | Recommended |

| Audit                | Optional | Mandatory   |

---

## 11. Design Constraints (Non-Negotiable)

* ❌ Do not block login for policy reasons
* ❌ Do not mix billing logic with auth
* ❌ Do not fork OSS vs Cloud codebases
* ✅ Use config & agents for behavior differences

---

## 11.1 Failure Semantics

The platform defines explicit failure behavior:

* Authentication failure → 401 / 403
* Authorization failure → 403
* Quota exceeded → 429 or 409 (explicit)
* Billing processing failure → async retry, never block ingestion
* Audit failure (Cloud mode) → request MUST fail

Silent failure is prohibited.

---

## 12. Summary

> **Authentication is a door, not a judge.**

> **Billing is strict. Security guidance is humane.**

This agent model ensures:

* OSS remains flexible and operator-friendly
* Cloud remains safe, auditable, and monetizable
* The platform scales without architectural regret

---
