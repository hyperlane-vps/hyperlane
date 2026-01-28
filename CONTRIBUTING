Contributing to Hyperlane

Thank you for your interest in contributing to Hyperlane. This guide explains how to get set up quickly and contribute effectively without guesswork.

Hyperlane is a production‑grade infrastructure project with strong opinions. Please read this document fully before opening issues or pull requests.

Prerequisites

Before contributing, you should be comfortable with:

Go (1.22+)

Linux systems

Virtualization concepts (KVM/QEMU)

Networking fundamentals

Reading and respecting architectural constraints

A local hypervisor is not required for most contributions.

Local Development Setup
1. Fork and Clone
git clone https://github.com/<your-username>/hyperlane.git
cd hyperlane
2. Install Tooling

Required:

Go 1.22+

GNU Make

Docker (for integration tests, optional)

Verify:

go version
Repository Orientation

Start here:

README.md – product vision and architecture

docs/ – design docs and RFCs

internal/state/ – VM state machine

internal/controlplane/ – reconciliation logic

cmd/vps/ – CLI entrypoint

If you do not understand the state machine, do not start coding yet.

Your First Contribution

Recommended starting points:

Documentation improvements

CLI error messages

Small, well‑scoped bug fixes

Look for issues labeled:

good-first-issue

docs

cli

Making Changes
Branching
git checkout -b feat/<short-description>
Coding Rules

Keep changes small and focused

Avoid refactors unless approved

Do not introduce new abstractions casually

Prefer explicit code over clever code

Testing

At minimum:

go test ./...

All tests must be deterministic. Flaky tests will block merges.

Submitting a Pull Request

Before opening a PR:

Ensure tests pass

Ensure the change aligns with project philosophy

Update documentation if behavior changes

Your PR should clearly explain:

Why the change is needed

What problem it solves

PRs that expand scope or bypass design discussion may be closed.

Communication

Use GitHub Discussions for design questions

Use Issues only for concrete problems

Be concise and technical

Final Reminder

Hyperlane values contributors who:

Respect constraints

Think in systems

Optimize for long‑term maintainability

If that sounds like you — welcome.

Hyperlane — Servers, without the ceremony.
