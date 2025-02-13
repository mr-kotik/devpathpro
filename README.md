# DevPathPro

<div align="center">

<pre>
██████╗ ███████╗██╗   ██╗██████╗  █████╗ ████████╗██╗  ██╗██████╗ ██████╗  ██████╗ 
██╔══██╗██╔════╝██║   ██║██╔══██╗██╔══██╗╚══██╔══╝██║  ██║██╔══██╗██╔══██╗██╔═══██╗
██║  ██║█████╗  ██║   ██║██████╔╝███████║   ██║   ███████║██████╔╝██████╔╝██║   ██║
██║  ██║██╔══╝  ╚██╗ ██╔╝██╔═══╝ ██╔══██║   ██║   ██╔══██║██╔═══╝ ██╔══██╗██║   ██║
██████╔╝███████╗ ╚████╔╝ ██║     ██║  ██║   ██║   ██║  ██║██║     ██║  ██║╚██████╔╝
╚═════╝ ╚══════╝  ╚═══╝  ╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ 
</pre>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Windows](https://img.shields.io/badge/Windows-10%2B-blue.svg)](https://www.microsoft.com/windows)
[![GitHub release](https://img.shields.io/github/v/release/mr-kotik/DevPathPro?include_prereleases)](https://github.com/mr-kotik/devpathpro/releases)
[![GitHub issues](https://img.shields.io/github/issues/mr-kotik/DevPathPro)](https://github.com/mr-kotik/devpathpro/issues)

🛠️ A powerful Windows environment manager that automatically detects, verifies, and configures development tools. Supports over 40 tools including programming languages, databases, build systems, and infrastructure solutions.

`#windows` `#development` `#environment` `#path` `#devtools` `#golang`

[Installation](#-installation) • [Features](#-key-features) • [Documentation](#-features-in-detail) • [Contributing](#-contributing)

</div>

## 🚀 Overview

DevPathPro is an intelligent environment manager that streamlines the setup and maintenance of your development tools on Windows. It automatically detects installed tools, verifies their configurations, and ensures optimal environment settings.

## ✨ Key Features

- 🔍 **Automatic Tool Detection**: Finds installed development tools across your system
- ⚙️ **Smart Configuration**: Sets up and maintains environment variables
- 🛠️ **Path Management**: Intelligently manages system PATH entries
- 🔄 **Configuration Verification**: Checks and corrects existing settings
- 🎯 **Multi-Environment Support**: Handles multiple development environments
- 🎛️ **Selective Configuration**: Choose which settings to configure for each tool
- 🔒 **Security Settings**: Configures secure defaults for databases and services
- 💾 **Data Directory Management**: Sets up proper data and log directories
- 🌐 **Network Configuration**: Manages ports and connection settings
- 🔄 **Cache Management**: Configures cache locations and sizes

## 🎛️ Configuration Options

Each tool offers specific configuration options that can be selected during setup:

### Programming Languages

- **Python**:
  - Basic: HOME and PATH settings
  - Environment: Encoding, buffering settings
  - Development: Warnings, debug, optimization
  - Pip: Package manager configuration

- **Java**:
  - Basic: HOME and CLASSPATH
  - JVM: Memory and performance settings

- **Node.js**:
  - Basic: Runtime configuration
  - NPM: Package manager settings

- **Go**:
  - Basic: GOROOT, GOPATH, modules, proxy settings

- **Rust**:
  - Basic: Cargo, compiler, and documentation settings

### Build Systems

- **Maven**:
  - Basic: Repository, memory, and debug settings

- **Gradle**:
  - Basic: Home, cache, and daemon settings

### Databases

- **PostgreSQL**:
  - Basic: Data directory, ports, encoding, SSL

- **MySQL**:
  - Basic: Data directory, ports, configuration

- **MongoDB**:
  - Basic: Data directory, logs, port settings

- **Redis**:
  - Basic: Port, configuration, data directory

- **Elasticsearch**:
  - Basic: Memory, paths, ports configuration

### Infrastructure

- **Docker**:
  - Basic: BuildKit, experimental features, host settings

- **Kubernetes**:
  - Basic: Config, editor, Helm settings

## 🛠️ Supported Tools

<details>
<summary>💻 Programming Languages</summary>

- Python
- Node.js
- Java
- Go
- Rust
- Perl
- Scala
- Kotlin
- Swift
- Haskell
- Erlang
- Elixir
- .NET Core
- Ruby
</details>

<details>
<summary>🏗️ Build Systems</summary>

- CMake
- MSBuild
- Maven
- Gradle
- Make
- Ninja
</details>

<details>
<summary>🔧 Development Tools</summary>

- Git
- Visual Studio
- VS Code
- Docker
- Kubernetes
- Windows SDK
- WDK
- LLVM
- Jenkins
- SonarQube
- Grafana
</details>

<details>
<summary>📦 Package Managers</summary>

- vcpkg
- Conan
</details>

<details>
<summary>💾 Databases</summary>

- PostgreSQL
- MySQL
- MongoDB
- Redis
- Elasticsearch
- SQLite
- Oracle
- Cassandra
- Neo4j
- InfluxDB
</details>

<details>
<summary>☁️ Infrastructure</summary>

- Terraform
- Ansible
- Podman
- Helm
- Skaffold
</details>

## 📋 Requirements

- Windows 10 or later
- Administrator privileges
- PowerShell or Command Prompt
- Go 1.21 or later (for building from source)

## 🔨 Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/mr-kotik/devpathpro.git
   cd devpathpro
   ```

2. Build the project:
   ```bash
   go build -o DevPathPro.exe
   ```

3. (Optional) Run tests:
   ```bash
   go test ./...
   ```

## 📥 Installation

1. Download the latest release from the [releases page](https://github.com/mr-kotik/devpathpro/releases)
2. Run the executable as Administrator
3. Follow the initial setup wizard

## 🚦 Usage

1. Launch DevPathPro as Administrator
2. Choose from the main menu:
   - 🔍 Search and Configure Development Tools
   - ✔️ Verify Existing Configurations
   - 👀 View Current Environment Settings
3. When tools are found, select configuration options:
   - Choose "All" for recommended settings
   - Select specific options for custom configuration
4. Review and confirm the changes

## 🔧 Configuration Process

1. **Tool Detection**:
   - Scans common installation directories
   - Searches PATH for executables
   - Optional deep search across all drives

2. **Configuration Selection**:
   - Lists available configuration groups
   - Allows multiple option selection
   - Shows detailed description for each option

3. **Environment Setup**:
   - Sets required environment variables
   - Configures tool-specific settings
   - Manages PATH entries
   - Sets up data directories

4. **Verification**:
   - Checks applied configurations
   - Validates paths and permissions
   - Tests connectivity where applicable
   - Suggests corrections if needed

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 🙏 Acknowledgments

- Thanks to all contributors who have helped shape DevPathPro
- Special thanks to the open-source community for the amazing tools we support

---

<div align="center">

Made with ❤️ by [Alesta](https://github.com/mr-kotik)

</div> 