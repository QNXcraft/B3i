# Agents.md

## Purpose
This file documents the automated and human agent roles used to develop, maintain, and release the BB10 / PlayBook App Manager repository (B3i). It explains responsibilities, workflows, security constraints, and how to extend or add agents.

## Roles and responsibilities

### Project Manager Agent
- Coordinates milestones and issues.
- Creates release checklists and updates CHANGELOG.md.
- Ensures tasks are prioritized and tracked.
- **Contact**: `@b3i/pm-agent`

### Lead Developer Agent
- Designs architecture and implements core Go modules.
- Reviews major design decisions and PRs.
- Ensures code quality and idiomatic Go.
- **Contact**: `@b3i/lead-dev-agent`

### Security Auditor Agent
- Reviews TLS handling, secrets, and UI security.
- Writes security guidance and tests for TLS behavior.
- Ensures no plaintext passwords are logged or stored.
- **Contact**: `@b3i/security-auditor`

### DevOps / Release Agent
- Maintains `Taskfile.yml` and CI workflows.
- Builds cross‑platform binaries and uploads release assets.
- Generates checksums for release artifacts.
- **Contact**: `@b3i/devops-agent`

### QA / Test Agent
- Writes unit and integration tests.
- Provides test HTTP servers that simulate BB10 Developer Mode endpoints.
- Ensures `go test ./...` passes.
- **Contact**: `@b3i/qa-agent`

### UX / Frontend Agent
- Produces minimal, secure static UI assets (HTML/CSS/JS).
- Ensures UI works offline and avoids remote dependencies.
- Implements drag‑and‑drop BAR upload flow.
- **Contact**: `@b3i/ux-agent`

### Docs & Onboarding Agent
- Writes README, CONTRIBUTING, and onboarding docs.
- Keeps documentation up to date with code changes.
- **Contact**: `@b3i/docs-agent`

### Release Manager (human or agent)
- Tags releases and verifies artifacts.
- Signs binaries if signing keys are available (human required for private keys).
- **Contact**: `@b3i/release-manager`

## Workflow
1. Issue created → Project Manager assigns role(s).
2. Lead Developer implements feature in a feature branch.
3. QA Agent writes tests and runs them locally.
4. PR opened; at least one human reviewer plus automated checks required.
5. Security Auditor reviews any TLS/secret changes.
6. DevOps Agent runs CI; on tag `v*` the Release Agent builds and uploads assets.
7. Release Manager verifies and publishes.

## Security rules for agents
- Never store plaintext passwords or tokens in the repository.
- Use `GITHUB_TOKEN` for release uploads; do not commit personal tokens.
- Default to TLS verification; allow `--insecure` only with explicit user consent and clear warnings.
- **Requirement**: Human approval for changes that weaken security (e.g., disabling TLS verification by default) or modify release signing.

## CI and release responsibilities
- `Taskfile.yml` maintained by DevOps Agent.
- `.github/workflows/release.yml` builds and uploads artifacts on tag push.
- Release artifacts must include checksums (`sha256`) and a release notes entry.

## Extending agents
- To add a new agent, create a new section here with responsibilities and acceptance criteria.
- Assign a GitHub team or username as the contact for the agent.

## Example commands
- Build all: `task build:all`
- Run tests: `task test` or `go test ./...`
- Start UI: `./dist/b3i-linux-amd64 serve --ui-port 8080`

## Contact and escalation
- For security issues: open an issue with label `security` and assign to `@b3i/security-auditor`.
- For release issues: ping `@b3i/release-manager`.
