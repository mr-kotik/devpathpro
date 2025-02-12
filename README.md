# 🔍 Path Finder

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows-lightgrey.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

*A utility for finding and managing PATH environment variables for Windows development tools*
</div>

## 📋 Table of Contents
- [Features](#-features)
- [Supported Tools](#-supported-tools)
- [Installation](#-installation)
- [Usage](#-usage)
- [System Requirements](#-system-requirements)
- [Notes](#-notes)

## ✨ Features

- 🔄 Interactive tool selection by categories
- 🔍 Smart search across multiple locations:
  - Standard installation paths
  - Windows Registry
  - Current PATH variables
  - All available drives
- 📝 Detailed information about found executables
- ⚡ Parallel search for high performance
- 🛠️ Easy PATH management:
  - View existing paths
  - Add new paths
  - Skip duplicates

## 🔧 Supported Tools

### 🏗️ Build Systems
- CMake
- Make (GnuWin32, MSYS2, Cygwin, MinGW)
- Ninja
- Maven
- Gradle
- MSBuild

### 💻 Programming Languages
- Python
- Node.js
- Java JDK
- Rust

### 🛠️ Development Tools
- Git
- Visual Studio Code
- Docker
- Kubernetes

## 📥 Installation

1. Download the latest version of `pathfinder.exe` from [releases](https://github.com/mr-kotik/DevPathVariableRestorer/releases)
2. Place the file in a convenient directory
3. Run the program as administrator

## 🚀 Usage

1. Launch the program:
   ```
   pathfinder.exe
   ```

2. In the main menu:
   - Select tools to search for by entering their numbers separated by commas (e.g.: 1,3,5)
   - Type 'all' to search for all tools
   - Type 'exit' to quit the program

3. For each found tool:
   - View existing paths in PATH
   - Add new paths to PATH
   - Skip the current tool

4. After completion:
   - Press Enter to return to the main menu
   - Select other tools or exit

## 💻 System Requirements

- Operating System: Windows
- Administrator privileges (for modifying PATH)
- Minimum 50 MB free disk space

## 📝 Notes

- The program normalizes paths for comparison but preserves the original format when adding to PATH
- Search optimization:
  - Parallel search across drives
  - Skipping system directories
  - Smart search order (standard installation locations checked first)
- Security:
  - Duplicate checking before addition
  - Original path format preservation
  - Protection against invalid paths

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE) file for more information.

## 🤝 Contributing

We welcome your contributions to the project! If you have suggestions for improvements or found a bug:

1. Create an Issue
2. Submit a Pull Request
3. Contact the developers

---
<div align="center">
Made with ❤️ for developers
</div> 