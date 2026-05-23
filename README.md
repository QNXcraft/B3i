# B3i (BlackBerry Bar installer)

B3i is a cross-platform (Linux, macOS, Windows) BB10 / PlayBook App Manager written in Go. It provides both a powerful CLI and a simple web UI for installing, uninstalling, and managing applications on your BlackBerry 10 or PlayBook devices.

## Features
- **Cross-platform**: Built with Go for Linux, macOS, and Windows.
- **Cobra CLI**: Command-line interface for automation and power users.
- **Gin Web UI**: Easy-to-use web interface with drag-and-drop BAR installation.
- **Embedded Assets**: No external dependencies or CDNs required.
- **Secure**: Proxying device calls through the local server to handle TLS safely.

## Installation
You can download the pre-built binaries from the [Releases](https://github.com/user/b3i/releases) page.

Alternatively, build from source using [Task](https://taskfile.dev):
```bash
task build
```

## Usage
### CLI
List installed apps:
```bash
./bin/b3i list --device <IP> --password <PWD>
```

Install a BAR file:
```bash
./bin/b3i install myapp.bar --device <IP> --password <PWD>
```

### Web UI
Start the web server:
```bash
./bin/b3i serve --device <IP> --password <PWD> --ui-port 8080
```
Then navigate to `http://localhost:8080` in your browser.

## Security
B3i communicates with BB10/PlayBook devices over HTTPS (port 1337). Since these devices typically use self-signed certificates in Developer Mode, B3i allows insecure TLS connections via the `--insecure` flag (defaulted to `true` for developer convenience).

**Warning**: Only use B3i on trusted networks when connecting to devices with Developer Mode enabled.

## Project agents and responsibilities
This project utilizes specialized AI agent roles for development and maintenance.
See [Agents.md](Agents.md) for a detailed description of roles, workflows, and security policies.

- **Lead Developer**: Architecture and core logic.
- **Security Auditor**: TLS and secrets management review.
- **DevOps**: CI/CD and Taskfile maintenance.
- **QA**: Testing and device simulation.
- **UX**: Web UI design and implementation.

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License
MIT
