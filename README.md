# logz

[![Build Status](https://img.shields.io/travis/yourusername/logz.svg?style=flat)](https://travis-ci.org/yourusername/logz)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/logz)](https://goreportcard.com/report/github.com/yourusername/logz)
[![Coverage Status](https://coveralls.io/repos/github/yourusername/logz/badge.svg?branch=main)](https://coveralls.io/github/yourusername/logz?branch=main)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**logz** is a highly modular logging framework written in Go. It is designed to be used both as a native logger from within Go applications and as a standalone CLI tool to record logs, monitor metrics, and integrate with external notification systems.  
logz supports integrations with multiple platforms such as DBus, Prometheus, and secure ZeroMQ communication for centralized service management (e.g., with gkbxsrv).

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [CLI Commands](#cli-commands)
  - [Configuration](#configuration)
- [Integration](#integration)
  - [Logger API](#logger-api)
  - [Notifiers](#notifiers)
  - [Metrics and Prometheus](#metrics-and-prometheus)
- [Development](#development)
- [License](#license)

---

## Features

- **Core Logging**  
  - Structured `LogEntry` with chainable methods to set _Timestamp_, _Level_, _Message_, _Source_, _Caller_, _Tags_, _Metadata_, and other fields.
  - Support for multiple output formats (JSON, plaintext with ANSI color support).

- **Output and Reading**  
  - Flexible `LogWriter` interface for writing logs to stdout or files.
  - `LogReader` implementation that supports tailing log files in real time with interruption.

- **Notifiers**  
  - External notifications via HTTP and ZeroMQ (including secure ZeroMQ integration with JWT for gkbxsrv).
  - DBus-based integration for passive log delivery.
  - Centralized `NotifierManager` for dynamic management of notifiers through a CLI.

- **Metrics and Monitoring**  
  - Integrated metrics collection using the `PrometheusManager`, with in-memory caching and disk persistence.
  - Ability to update and export selected metrics to Prometheus based on a configurable whitelist.
  - Real-time monitoring dashboards for both metrics and logs through CLI commands.

- **Dynamic Configuration**  
  - Use [Viper](https://github.com/spf13/viper) to manage dynamic configuration of the service.
  - Automatic reload of configurations on change, with graceful service restart.

- **CLI Interface**  
  - Commands to start, stop, and monitor the service.
  - Dedicated commands for managing metrics and notifiers.
  - Interface for real-time log “tailing” and metrics dashboards.

---

## Installation

### Prerequisites

- [Go](https://golang.org) 1.18+
- Git

### Building

Clone the repository and build:

```bash
git clone https://github.com/yourusername/logz.git
cd logz
go build -o logz .
```

---

## Usage

### CLI Commands

You can use logz as a standalone CLI tool. Some key commands include:

- **Service Management**
    - `logz service start --port 9999` — Launch the logz HTTP service on port 9999.
    - `logz service stop` — Stop the running service.
    - `logz service status` — Display the status of the service.

- **Metrics Management**
    - `logz metrics add [name] [value]` — Add or update a metric.
    - `logz metrics remove [name]` — Remove a metric.
    - `logz metrics list` — List all registered metrics.
    - `logz metrics watch` — Display a real-time dashboard of metrics.

- **Logs**
    - `logz logs list` — Show static logs.
    - `logz logs watch` — Tail logs in real time.

- **Notifier Management**
    - `logz notifiers add --name discord --type external --url "https://discord.com/api/webhooks/..."` — Register a new notifier.
    - `logz notifiers remove --name discord` — Remove a notifier.
    - `logz notifiers enable --name external` — Enable a specific notifier.
    - `logz notifiers disable --name external` — Disable a specific notifier.
    - `logz notifiers reload` — Force reload of notifier configuration from file (if configured with Viper).

### Configuration

logz uses Viper for configuration. The configuration file can be in JSON, YAML, TOML, or INI format. It is searched in the following order:

1. The path defined in the environment variable `LOGZ_CONFIG_PATH`
2. The user's configuration directory
3. The user's home directory
4. The cache directory (fallback)

If no configuration file is found, a default configuration file (`service_config.json`) is generated automatically.

Key environment variables:
- `LOGZ_PROMETHEUS_ENABLED`: Enable/disable Prometheus integration.
- `LOGZ_NO_COLOR`: Disable color output.
- `LOGZ_PID_PATH`: Define the PID file location.
- `LOGZ_METRICS_FILE`: Define the metrics persistence file.
- `LOGZ_DBUS_ENABLED`: Enable DBus integration.
- `LOGZ_CONFIG_PATH`: Specify a custom path for the configuration file.
- Others for secure ZeroMQ (e.g., `LOGZ_ZMQSEC_ENABLED`, `LOGZ_ZMQSEC_ENDPOINT`, etc.)

---

## Integration

### Logger API

Import and use logz as a logging library integrated in your Go code:

```go
import "github.com/yourusername/logz/logger"

func main() {
    logz := logger.NewLogger(logger.INFO, "text", "stdout", "https://your.endpoint", "tcp://your.zmq", "")
    logz.SetMetadata("integration", "your-app")
    logz.Info("Application started", nil)
}
```

### Notifiers

Notifiers are dynamically managed by the NotifierManager. They support authentication via token (using SetAuthToken). For example, the external notifier can send logs via HTTP and ZeroMQ, while DBus and secure ZeroMQ (ZMQSecNotifier) are integrated for specific scenarios.

### Metrics and Prometheus

logz automatically updates metrics (like "logs_total" and "logs_total_<LEVEL>") for each log entry via the PrometheusManager. The metrics are stored in memory and persisted to disk. You can control which metrics are exported using the export whitelist configuration.

---

## Development

Contributions, bug reports, and feature requests are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
