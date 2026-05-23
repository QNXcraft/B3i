# Contributing to B3i

Thank you for your interest in contributing to B3i!

## Development Workflow
This project follows the workflow described in [Agents.md](Agents.md). Both human contributors and AI agents collaborate using the following process:

1. **Issue Creation**: Open an issue to describe the feature or bug.
2. **Branching**: Create a feature branch for your changes.
3. **Testing**: Ensure all tests pass by running `task test`. Add new tests for any new features.
4. **Pull Request**: Open a PR and request review.
5. **Security Review**: Any changes affecting TLS, authentication, or secrets require a review from the Security Auditor role (or a human equivalent).

## Prerequisites
- Go 1.21 or later
- [Task](https://taskfile.dev) for running build tasks

## Building
```bash
task build
```

## Running Tests
```bash
task test
```

## AI Agents
We use AI agents to assist with various roles. If you are an agent, please refer to [Agents.md](Agents.md) for your specific responsibilities and workflow instructions.
