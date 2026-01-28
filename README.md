# Hyperlane VPS

**Developer‑first, opinionated, open‑source VPS platform.**
Think *Vercel*, but for real servers.

Hyperlane provisions production‑grade virtual machines in seconds, with zero‑configuration defaults, instant SSH access, and a clean, modern control plane built for developers—not sysadmins.

---

## Why Hyperlane Exists

Traditional VPS platforms optimize for operators. Hyperlane optimizes for **developers**.

* No complex installers
* No YAML on day one
* No OpenStack‑style sprawl
* No guessing what state your server is in

Just:

```bash
vps up
```

And you’re live.

---

## Core Principles

### Opinionated by Default

Hyperlane intentionally removes choices early:

* **KVM + QEMU only**
* **WireGuard networking**
* **IPv6‑first**
* **ZFS local storage**
* **Immutable VM images**
* **Single‑binary control plane**

These decisions make installs predictable, performance consistent, and contributions coherent.

---

### Developer Experience First

Hyperlane treats servers like deployable artifacts.

```bash
vps init
vps up
vps ssh api-prod
vps scale api-prod --cpu 4 --ram 8G
vps snapshot api-prod
vps destroy api-prod
```

* No cloud‑init gymnastics
* No SSH key management
* No manual firewall rules

Everything works out of the box.

---

## Architecture Overview

```
┌───────────────┐
│   Web UI      │  React + WebSockets
└───────┬───────┘
        │
┌───────▼───────┐
│  API Gateway  │  REST + gRPC
└───────┬───────┘
        │
┌───────▼────────┐
│ Control Plane  │  Scheduler + State Engine
└───────┬────────┘
        │ mTLS
┌───────▼────────┐
│ Node Agent     │  Runs on each hypervisor
└────────────────┘
```

### Design Guarantees

* Event‑driven
* Idempotent operations
* Declarative desired state (internal)
* No shell‑based orchestration

---

## Ultra‑Fast Provisioning

Hyperlane provisions VMs in **under 5 seconds** by design:

* Pre‑warmed base images
* ZFS instant cloning
* Static kernels
* Minimal boot path
* SSH certificates (no keys)

VM creation is effectively:

> clone → boot → ready

---

## Networking Model

* WireGuard mesh between hosts and VMs
* No public IPs by default
* Identity‑based firewall rules
* Public exposure is explicit

```bash
vps expose api-prod --port 443
```

Internally:

* eBPF firewall
* Floating IP abstraction
* NAT only when required

---

## Images as First‑Class Citizens

Images are versioned, immutable artifacts.

```bash
vps image list
vps image create my-api --from api-prod
vps up --image my-api:v3
```

This enables:

* Preview environments
* Instant rollbacks
* Dev/prod parity

---

## GitOps (Optional)

GitOps is available—but never forced.

```yaml
vm:
  name: api-prod
  cpu: 4
  ram: 8G
  image: my-api:v3
```

* Git is the source of truth
* CLI and UI still work
* Drift is reconciled automatically

---

## Installation

Single‑node install (Linux host):

```bash
curl -fsSL get.hyperlane.sh | sh
hyperlane init
```

Requirements:

* Linux host
* KVM enabled
* ZFS installed

---

## Roadmap

### v0.1

* Single‑node
* CLI + API
* VM lifecycle
* ZFS snapshots
* SSH certificates
* Web UI

### v0.2

* Multi‑node scheduling
* Networking policies
* Metrics & logs

### v0.3

* GitOps reconciliation
* RBAC
* Multi‑tenant support
* Provider billing hooks

---

## Open Source

* **License:** Apache 2.0
* No open‑core bait‑and‑switch
* Provider‑friendly by default

Enterprise features live in separate repositories.

---

## Who This Is For

* Indie hosting providers
* Dev‑focused infrastructure teams
* SaaS platforms needing real servers
* Power users who want speed without chaos

---

## Repository Structure

Hyperlane is structured to keep responsibilities clean, testable, and contributor-friendly.

```
hyperlane/
├─ cmd/
│  ├─ hyperlane/        # Control plane binary
│  ├─ hyperlane-agent/  # Node agent (runs on hypervisors)
│  └─ vps/              # Developer CLI
│
├─ internal/
│  ├─ api/              # gRPC / REST definitions
│  ├─ auth/             # mTLS, SSH certs, identity
│  ├─ controlplane/     # Scheduler, reconciliation engine
│  ├─ images/           # Image lifecycle & caching
│  ├─ network/          # WireGuard, eBPF firewall logic
│  ├─ state/            # Desired vs actual state machine
│  ├─ storage/          # ZFS + snapshot abstraction
│  └─ telemetry/        # Metrics, logs, events
│
├─ pkg/                 # Public Go SDK (stable API)
├─ web/                 # Web UI (React + Vite)
├─ scripts/             # Dev & CI helpers only
├─ docs/                # Architecture, RFCs, guides
├─ CONTRIBUTING.md
├─ LICENSE
└─ README.md
```

### Structural Rules

* `internal/` may change freely
* `pkg/` is versioned and stable
* No business logic in `cmd/`
* Node agent and control plane are strictly separated

---

## CLI UX Specification (v0.1)

The CLI is the primary interface. It must be:

* Predictable
* Scriptable
* Fast
* Human-readable by default

### Global Rules

* Defaults over flags
* No required config files
* JSON output via `--json`
* Clear, opinionated error messages

### Core Commands

```bash
vps init                 # Initialize project + auth
vps up                   # Create or reconcile VM
vps list                 # List VMs
vps ssh <name>           # Instant SSH access
vps logs <name>          # Stream logs
vps scale <name>         # Resize CPU / RAM
vps snapshot <name>      # Create snapshot
vps destroy <name>       # Destroy VM
```

### Example

```bash
vps up api-prod --cpu 4 --ram 8G
```

Output:

```
✔ Image ready
✔ Network attached
✔ VM started (3.2s)
→ SSH: vps ssh api-prod
```

---

## VM State Machine (Core Invariant)

Every VM exists in exactly one state:

```
requested → provisioning → running → stopped → destroyed
                     ↘ failed
```

Rules:

* All operations are idempotent
* Failed states are inspectable and recoverable
* Control plane reconciles until desired == actual

No hidden or implicit transitions.

---

## Open Source

* **License:** Apache 2.0
* No open-core bait-and-switch
* Provider-friendly by default

Enterprise features live in separate repositories.

---

## Who This Is For

* Indie hosting providers
* Dev-focused infrastructure teams
* SaaS platforms needing real servers
* Power users who want speed without chaos

---

## Contributing

Hyperlane welcomes contributors who value:

* Strong opinions
* Clean abstractions
* Performance discipline
* Developer empathy

See `CONTRIBUTING.md` to get started.

---

## Control Plane RFC (v0.1)

This document defines the **non-negotiable behavior** of the Hyperlane control plane. All implementations must preserve these guarantees.

---

### Responsibilities

The control plane is responsible for:

* Translating user intent into desired state
* Scheduling VMs onto nodes
* Reconciling actual state to desired state
* Handling failures deterministically
* Exposing a stable API to the CLI and UI

It is **not** responsible for:

* Executing hypervisor commands directly
* Persisting large binaries or images
* Handling node-local decisions

---

### Core Loop (Reconciliation Engine)

The control plane operates on a continuous reconciliation loop:

```
observe → diff → plan → act → verify
```

* **observe**: read current state from agents
* **diff**: compare against desired state
* **plan**: compute minimal corrective actions
* **act**: dispatch commands to node agents
* **verify**: confirm convergence

All actions must be:

* Idempotent
* Retryable
* Order-independent

---

### State Model

The control plane persists only **desired state** and **observed state**.

```text
desired_vm_state
observed_vm_state
```

No transitional state is trusted unless confirmed by an agent.

---

### Scheduling (v0.1)

v0.1 uses a deterministic, single-pass scheduler.

Placement rules:

* Respect CPU and RAM capacity
* Prefer least-loaded node
* No live migration
* No overcommit

Future versions may add:

* Affinity rules
* Preemption
* Cost-based placement

---

### Failure Handling

Failures are explicit states, never silent.

Examples:

* Image pull failure → `failed`
* Network attach failure → `failed`
* Boot timeout → `failed`

Recovery requires one of:

* User retry
* Automated retry (bounded)
* Explicit destroy

---

### Node Agent Contract

Node agents:

* Execute commands from control plane
* Report state continuously
* Never make scheduling decisions
* Never mutate desired state

All communication:

* mTLS
* Versioned protocol
* Backward compatible for one minor version

---

### Persistence

* Embedded database (SQLite / BoltDB)
* Single writer
* Event-sourced log (append-only)

This guarantees:

* Crash recovery
* Deterministic replay
* Auditable history

---

### Invariants

The following must always hold true:

* One VM maps to exactly one node
* No VM runs without a desired state
* No agent mutates global state
* Control plane is restart-safe

Violations are bugs.

---

## Maintainer Philosophy

Hyperlane prioritizes:

* Clarity over flexibility
* Determinism over magic
* Developer trust over features

PRs that compromise these values will be rejected.

## Contributing

Hyperlane is a **serious infrastructure project**. Contributions are welcome, but not all contributions are accepted. This document exists to protect the project, its users, and its maintainers.

If you are looking for a flexible playground or a general-purpose virtualization framework, this project is not a good fit.

---

## Project Philosophy

Hyperlane prioritizes:

* **Clarity over flexibility**
* **Determinism over magic**
* **Strong opinions over endless configuration**
* **Developer experience over operator convenience**

Any contribution that weakens these principles will be declined, even if technically impressive.

---

## What We Are Building

Hyperlane is:

* A developer-first VPS platform
* Opinionated by default
* Designed for fast provisioning and predictable behavior
* Suitable for real hosting providers

Hyperlane is **not**:

* A general-purpose cloud
* A replacement for OpenStack
* A Kubernetes distribution
* A hypervisor abstraction layer

---

## Contribution Types We Welcome

### ✅ Bug Fixes

* Correctness issues
* Race conditions
* Resource leaks
* Security fixes

### ✅ Performance Improvements

* Faster provisioning paths
* Lower memory usage
* Reduced VM boot time

### ✅ Developer Experience

* CLI improvements
* Better error messages
* Documentation clarity

### ✅ Core Features (Only When Approved)

* Must align with the roadmap
* Must be discussed before implementation
* Must not expand scope unintentionally

---

## Contributions We Will Decline

These PRs will be closed without debate:

* Adding new hypervisors (Xen, ESXi, Hyper-V)
* Adding alternative storage backends before v0.3
* Adding Kubernetes integration
* Adding plugin systems before v0.3
* Replacing core technologies (WireGuard, ZFS)
* Large refactors without prior discussion

---

## Design First, Code Second

Before writing code for anything non-trivial:

1. Open a **GitHub Discussion** or **RFC document**
2. Clearly state the problem being solved
3. Explain why existing mechanisms are insufficient
4. Describe the minimal acceptable solution

PRs without prior alignment may be closed.

---

## Coding Standards

* Language: **Go** (primary)
* Clear, boring code is preferred
* No clever abstractions
* Explicit state transitions only
* No global mutable state

### Rules

* No business logic in `cmd/`
* Node agents may not mutate desired state
* Control plane must be restart-safe
* Errors must be actionable

---

## Testing Expectations

All contributions must include:

* Unit tests for new logic
* Deterministic test behavior
* Clear failure cases

If a feature cannot be tested reliably, it is likely not acceptable.

---

## Commit & PR Guidelines

### Commits

* Small, focused commits
* Clear intent in messages
* No "misc" or "cleanup" commits

### Pull Requests

* One concern per PR
* Describe **why**, not just **what**
* Link to related issue or discussion
* Expect review and requested changes

---

## Maintainer Authority

Maintainers:

* Decide scope
* Decide roadmap
* May reject PRs unilaterally
* Are not obligated to merge features

This is not hostility—it is how infrastructure stays reliable.

---

## Security Issues

Do **not** open public issues for security vulnerabilities.

Instead:

* Email: `security@hyperlane.dev` (placeholder)

Responsible disclosure is required.

---

## Final Note

Hyperlane values contributors who:

* Respect constraints
* Think in systems
* Prefer boring correctness
* Care about long-term maintainability

If that sounds like you, welcome.

---

**Hyperlane** — *Servers, without the ceremony.*
