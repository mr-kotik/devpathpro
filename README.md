# Path Finder

A utility for finding and managing development tools in Windows PATH environment variable.

## Features

- Interactive tool selection from multiple categories
- Smart search across multiple locations:
  - Common installation directories
  - Windows Registry
  - System PATH
  - All available drives
- Detailed information about found executables
- Easy PATH management:
  - View existing paths
  - Add new paths
  - Skip duplicate paths

## Supported Tools

### Build Systems
- CMake
- Make
- Ninja
- Maven
- Gradle

### Programming Languages
- Python
- Node.js
- Java JDK
- Rust

### Development Tools
- Git
- Visual Studio Code
- Docker
- Kubernetes

## Usage

1. Run the program:
   ```
   pathfinder.exe
   ```

2. Select tools to search for:
   - Enter numbers separated by comma (e.g.: 1,3,5)
   - Type 'all' to search for all tools
   - Type 'skip' to skip current tool

3. For each found tool, you can:
   - View existing paths in PATH
   - Add new paths to PATH
   - Skip adding paths

## System Requirements

- Windows operating system
- Administrative privileges (for modifying PATH)

## Notes

- The program normalizes paths for comparison but preserves original path format when adding to PATH
- Search performance is optimized by:
  - Parallel search across drives
  - Skipping system directories
  - Smart search order (common locations first) 