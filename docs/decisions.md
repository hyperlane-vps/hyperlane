This document records intentional, permanent architectural decisions in Hyperlane. These decisions are not revisited casually.

If you are proposing a change that contradicts this document, it must go through a formal RFC and will likely be declined.

Decision: KVM + QEMU Only

Status: Locked

Hyperlane supports KVM with QEMU as the sole virtualization stack.

Rationale

Native Linux performance

Mature tooling and observability

Stable APIs

Broad hardware support

Supporting additional hypervisors increases complexity, fragments behavior, and harms reliability.

Decision: WireGuard for Networking

Status: Locked

All host-to-VM and inter-node networking is built on WireGuard.

Rationale

Simple mental model

Modern cryptography

Low operational overhead

Excellent performance

No alternative overlays or SDN layers will be added before v1.0.

Decision: IPv6-First Networking

Status: Locked

Hyperlane is IPv6-first internally.

Rationale

Eliminates unnecessary NAT

Simplifies addressing

Scales cleanly

IPv4 is supported only for external exposure and compatibility.

Decision: ZFS for Local Storage

Status: Locked (v0.x)

Hyperlane uses ZFS for local VM storage.

Rationale

Instant snapshots and clones

Data integrity guarantees

Predictable performance

Alternative storage backends may be considered after v0.3.

Decision: Immutable VM Images

Status: Locked

VM images are immutable, versioned artifacts.

Rationale

Reproducibility

Fast provisioning

Safe rollbacks

Mutable base images are explicitly rejected.

Decision: No Kubernetes Integration

Status: Locked

Hyperlane does not integrate with Kubernetes.

Rationale

Different abstraction level

Avoids conceptual overlap

Keeps scope focused

Users may run Kubernetes inside VMs if desired.

Decision: No Plugin System (Pre-v0.3)

Status: Locked

Hyperlane does not support plugins in early versions.

Rationale

Plugins freeze bad abstractions

Early extensibility increases maintenance burden

Extension points will be added only after core behavior stabilizes.

Decision: Single Control Plane Writer

Status: Locked

Only one control plane instance mutates desired state at any time.

Rationale

Deterministic behavior

Simplified recovery

Reduced failure modes

Horizontal scaling is introduced only when required.

Decision: Apache 2.0 License

Status: Locked

Hyperlane is licensed under Apache 2.0.

Rationale

Provider-friendly

Allows commercial adoption

Encourages ecosystem growth

No copyleft license will be introduced.

How to Propose a Change

Open a GitHub Discussion

Reference this document explicitly

Explain the problem, not the solution

Accept that "no" is a valid outcome


Hyperlane — Servers, without the ceremony.
