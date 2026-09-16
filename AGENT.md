# AGENT.md: Instructions for AI Coding Agents

Welcome, Agent. This document establishes technical constraints, architecture rules, project conventions, and non-negotiables for contributing code to the **SyncDrop** monorepo.

---

## 1. Project Philosophy & Core Constraints

* **Privacy & Local-First:** SyncDrop must never transmit telemetry, user analytics, or media to external third-party cloud services.
* **Minimal Dependencies:** Prioritize standard library utilities (both in Go and Dart) before adding external packages.
* **No Open Storage Ports:** Never implement code that assumes public WAN exposure of raw HTTP endpoints. WAN traffic is strictly routed through WireGuard.
* **Monorepo Discipline:** Maintain strict boundary separation between `server/` and `client/`. Never introduce cross-directory relative path imports between backend and frontend.

---

## 2. Directory Layout & Context

```text
syncdrop/
├── server/             # Go service daemon
│   ├── cmd/server/     # Entry point (main.go)
│   ├── internal/       # Core business logic (discovery, storage, auth, api)
│   └── scripts/        # DuckDNS updater and deployment scripts
├── client/             # Flutter mobile app
│   ├── lib/            # App code (clean/feature-first architecture)
│   └── test/           # Unit and widget tests
├── docs/               # PRD, specs, and architecture documents
├── AGENT.md            # Agent operational guide (this file)
└── README.md
