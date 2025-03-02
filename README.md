![Logz Banner](./assets/top_banner.png)

---

**An advanced logging and metrics management tool with native support for Prometheus integration, dynamic notifications, and a powerful CLI.**

---

## **Table of Contents**
1. [About the Project](#about-the-project)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
    - [CLI](#cli)
    - [Configuration](#configuration)
5. [Prometheus Integration](#prometheus-integration)
6. [Roadmap](#roadmap)
7. [Contributing](#contributing)
8. [Contact](#contact)

---

## **About the Project**
Logz is a flexible and powerful solution for managing logs and metrics in modern systems. Built in **Go**, it provides extensive support for multiple notification methods such as **HTTP Webhooks**, **ZeroMQ**, and **DBus**, alongside seamless integration with **Prometheus** for advanced monitoring.

Logz is designed to be robust, highly configurable, and scalable, catering to developers, DevOps teams, and software architects who need a centralized approach to logging and metrics.

**Why Logz?**
- 💡 **Ease of Use**: Configure and manage logs effortlessly.
- 🌐 **Seamless Integration**: Easily integrates with Prometheus and other systems.
- 🔧 **Extensibility**: Add new notifiers and services as needed.

---

## **Features**
✨ **Dynamic Notifiers**:
- Support for multiple notifiers simultaneously.
- Centralized and flexible configuration via JSON or YAML.

📊 **Monitoring and Metrics**:
- Exposes Prometheus-compatible metrics.
- Dynamic management of metrics with persistence support.

💻 **Powerful CLI**:
- Straightforward commands to manage logs and services.
- Extensible for additional workflows.

🔒 **Resilient and Secure**:
- Validates against Prometheus naming conventions.
- Distinct modes for standalone and service execution.

---

## **Installation**
Requirements:
- **Go** version 1.19 or later.
- Prometheus (optional for advanced monitoring).

```bash
# Clone this repository
git clone https://github.com/faelmori/logz.git

# Navigate to the project directory
cd logz

# Compile the binary
go build -o logz

# (Optional) Add the binary to your PATH
export PATH=$PATH:$(pwd)
```

---

## **Usage**

### CLI
Here are some examples of commands you can execute with Logz’s CLI:

```bash
# Log at different levels
logz info --msg "Starting the application."
logz error --msg "Database connection failed."

# Start the detached service
logz start  

# Stop the detached service
logz stop  

# Watch logs in real-time
logz watch
```

### Configuration
Logz uses a JSON or YAML configuration file to centralize its setup. The file is automatically generated on first use or can be manually configured at:  
`~/.kubex/logz/config.json`.

**Example Configuration**:
```json
{
  "port": "2112",
  "bindAddress": "0.0.0.0",
  "logLevel": "info",
  "notifiers": {
    "webhook1": {
      "type": "http",
      "webhookURL": "https://example.com/webhook",
      "authToken": "your-token-here"
    }
  }
}
```

---

## **Prometheus Integration**
Once started, Logz exposes metrics at the endpoint:
```
http://localhost:2112/metrics
```

**Example Prometheus Configuration**:
```yaml
scrape_configs:
  - job_name: 'logz'
    static_configs:
      - targets: ['localhost:2112']
```

---

## **Roadmap**
🔜 **Upcoming Features**:
- Support for additional notifier types (e.g., Slack, Discord, and email).
- Integrated monitoring dashboard.
- Advanced configuration with automated validation.

---

## **Contributing**
Contributions are welcome! Feel free to open issues or submit pull requests. Check out the [Contributing Guide](CONTRIBUTING.md) for more details.

---

## **Contact**
💌 **Developer**:  
[Rafael Mori](mailto:faelmori@gmail.com)
💼 [Follow me on GitHub](https://github.com/faelmori)
I'm open to new work opportunities and collaborations. If you find this project interesting, don’t hesitate to reach out!
