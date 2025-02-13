package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Program structure holds information about a development tool
type Program struct {
	name           string
	executableName string
	commonPaths    []string
	category       string
}

// Define available programs
var allPrograms = []Program{
	{
		name:           "CMake",
		executableName: "cmake.exe",
		commonPaths: []string{
			`C:\Program Files\CMake\bin`,
			`C:\Program Files (x86)\CMake\bin`,
			`C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
			`C:\Program Files\Microsoft Visual Studio\2022\Professional\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
			`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
		},
		category: "Build Systems",
	},
	{
		name:           "Git",
		executableName: "git.exe",
		commonPaths: []string{
			`C:\Program Files\Git`,
			`C:\Program Files (x86)\Git`,
		},
		category: "Development Tools",
	},
	{
		name:           "MSBuild",
		executableName: "msbuild.exe",
		commonPaths: []string{
			`C:\Program Files\Microsoft Visual Studio\2022\Community\MSBuild\Current\Bin`,
			`C:\Program Files\Microsoft Visual Studio\2022\Professional\MSBuild\Current\Bin`,
			`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\MSBuild\Current\Bin`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\MSBuild\Current\Bin`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Professional\MSBuild\Current\Bin`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\MSBuild\Current\Bin`,
			`C:\Windows\Microsoft.NET\Framework\v4.0.30319`,
			`C:\Windows\Microsoft.NET\Framework64\v4.0.30319`,
		},
		category: "Build Systems",
	},
	{
		name:           "Python",
		executableName: "python.exe",
		commonPaths: []string{
			`C:\Python3`,
			`C:\Program Files\Python`,
			`C:\Program Files (x86)\Python`,
			`C:\Users\%USERNAME%\AppData\Local\Programs\Python`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Node.js",
		executableName: "node.exe",
		commonPaths: []string{
			`C:\Program Files\nodejs`,
			`C:\Program Files (x86)\nodejs`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Java",
		executableName: "javac.exe",
		commonPaths: []string{
			`C:\Program Files\Java`,
			`C:\Program Files (x86)\Java`,
			`C:\Program Files\Eclipse Foundation`,
		},
		category: "Programming Languages",
	},
	{
		name:           "VS Code",
		executableName: "code.exe",
		commonPaths: []string{
			`C:\Program Files\Microsoft VS Code`,
			`C:\Users\%USERNAME%\AppData\Local\Programs\Microsoft VS Code`,
		},
		category: "Development Tools",
	},
	{
		name:           "Rust",
		executableName: "rustc.exe",
		commonPaths: []string{
			`C:\Users\%USERNAME%\.cargo\bin`,
			`C:\Program Files\Rust`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Maven",
		executableName: "mvn.exe",
		commonPaths: []string{
			`C:\Program Files\Apache\maven`,
			`C:\Program Files (x86)\Apache\maven`,
		},
		category: "Build Systems",
	},
	{
		name:           "Gradle",
		executableName: "gradle.exe",
		commonPaths: []string{
			`C:\Program Files\Gradle`,
			`C:\Gradle`,
		},
		category: "Build Systems",
	},
	{
		name:           "Make",
		executableName: "make.exe",
		commonPaths: []string{
			`C:\Program Files\GnuWin32\bin`,
			`C:\Program Files (x86)\GnuWin32\bin`,
			`C:\MinGW\bin`,
			`C:\msys64\usr\bin`,           // MSYS2 path
			`C:\msys64\mingw64\bin`,       // MSYS2 MinGW-w64 path
			`C:\cygwin64\bin`,             // Cygwin path
			`C:\mingw-w64\mingw64\bin`,    // MinGW-w64 path
			`C:\Program Files\Git\usr\bin`, // Git for Windows path
		},
		category: "Build Systems",
	},
	{
		name:           "Ninja",
		executableName: "ninja.exe",
		commonPaths: []string{
			`C:\Program Files\Ninja`,
			`C:\Program Files (x86)\Ninja`,
		},
		category: "Build Systems",
	},
	{
		name:           "Docker",
		executableName: "docker.exe",
		commonPaths: []string{
			`C:\Program Files\Docker\Docker\resources\bin`,
			`C:\Program Files\Docker Toolbox`,
		},
		category: "Development Tools",
	},
	{
		name:           "Kubernetes",
		executableName: "kubectl.exe",
		commonPaths: []string{
			`C:\Program Files\Kubernetes\Minikube`,
			`C:\Program Files (x86)\Kubernetes`,
		},
		category: "Development Tools",
	},
	{
		name:           "Visual Studio",
		executableName: "devenv.exe",
		commonPaths: []string{
			`C:\Program Files\Microsoft Visual Studio\2022\Community`,
			`C:\Program Files\Microsoft Visual Studio\2022\Professional`,
			`C:\Program Files\Microsoft Visual Studio\2022\Enterprise`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Community`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Professional`,
			`C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise`,
		},
		category: "Development Tools",
	},
	{
		name:           "Windows SDK",
		executableName: "rc.exe",
		commonPaths: []string{
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64`,
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.22000.0\x64`,
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.19041.0\x64`,
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x86`,
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.22000.0\x86`,
			`C:\Program Files (x86)\Windows Kits\10\bin\10.0.19041.0\x86`,
		},
		category: "Development Tools",
	},
	{
		name:           "WDK",
		executableName: "devcon.exe",
		commonPaths: []string{
			`C:\Program Files (x86)\Windows Kits\10\Tools\x64`,
			`C:\Program Files (x86)\Windows Kits\10\Tools\x86`,
		},
		category: "Driver Development",
	},
	{
		name:           "LLVM",
		executableName: "clang.exe",
		commonPaths: []string{
			`C:\Program Files\LLVM`,
			`C:\Program Files (x86)\LLVM`,
		},
		category: "Development Tools",
	},
	{
		name:           "vcpkg",
		executableName: "vcpkg.exe",
		commonPaths: []string{
			`C:\vcpkg`,
			`C:\dev\vcpkg`,
			`C:\Program Files\vcpkg`,
		},
		category: "Package Manager",
	},
	{
		name:           "Conan",
		executableName: "conan.exe",
		commonPaths: []string{
			`C:\Program Files\Conan`,
			`C:\Users\%USERNAME%\AppData\Local\Programs\Python\Python3*\Scripts`,
		},
		category: "Package Manager",
	},
	{
		name:           "PostgreSQL",
		executableName: "psql.exe",
		commonPaths: []string{
			`C:\Program Files\PostgreSQL\*\bin`,
			`C:\Program Files (x86)\PostgreSQL\*\bin`,
		},
		category: "Databases",
	},
	{
		name:           "MySQL",
		executableName: "mysql.exe",
		commonPaths: []string{
			`C:\Program Files\MySQL\MySQL Server *\bin`,
			`C:\Program Files (x86)\MySQL\MySQL Server *\bin`,
			`C:\Program Files\MariaDB *\bin`,
		},
		category: "Databases",
	},
	{
		name:           "MongoDB",
		executableName: "mongod.exe",
		commonPaths: []string{
			`C:\Program Files\MongoDB\Server\*\bin`,
			`C:\Program Files\MongoDB\*\bin`,
		},
		category: "Databases",
	},
	{
		name:           "Redis",
		executableName: "redis-server.exe",
		commonPaths: []string{
			`C:\Program Files\Redis`,
			`C:\Program Files (x86)\Redis`,
		},
		category: "Databases",
	},
	{
		name:           "Elasticsearch",
		executableName: "elasticsearch.bat",
		commonPaths: []string{
			`C:\Program Files\Elastic\Elasticsearch\*\bin`,
			`C:\Program Files (x86)\Elastic\Elasticsearch\*\bin`,
		},
		category: "Databases",
	},
	{
		name:           ".NET Core",
		executableName: "dotnet.exe",
		commonPaths: []string{
			`C:\Program Files\dotnet`,
			`C:\Program Files (x86)\dotnet`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Ruby",
		executableName: "ruby.exe",
		commonPaths: []string{
			`C:\Ruby*\bin`,
			`C:\Program Files\Ruby*\bin`,
			`C:\Program Files (x86)\Ruby*\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Go",
		executableName: "go.exe",
		commonPaths: []string{
			`C:\Program Files\Go\bin`,
			`C:\Go\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Podman",
		executableName: "podman.exe",
		commonPaths: []string{
			`C:\Program Files\RedHat\Podman`,
			`C:\Program Files (x86)\RedHat\Podman`,
		},
		category: "Containerization",
	},
	{
		name:           "Helm",
		executableName: "helm.exe",
		commonPaths: []string{
			`C:\Program Files\Helm`,
			`C:\Program Files (x86)\Helm`,
			`C:\Users\%USERNAME%\.helm\bin`,
		},
		category: "Containerization",
	},
	{
		name:           "Skaffold",
		executableName: "skaffold.exe",
		commonPaths: []string{
			`C:\Program Files\Skaffold`,
			`C:\Program Files (x86)\Skaffold`,
			`C:\Users\%USERNAME%\.skaffold\bin`,
		},
		category: "Development Tools",
	},
	{
		name:           "SQLite",
		executableName: "sqlite3.exe",
		commonPaths: []string{
			`C:\Program Files\SQLite`,
			`C:\Program Files (x86)\SQLite`,
			`C:\sqlite`,
		},
		category: "Databases",
	},
	{
		name:           "Oracle",
		executableName: "sqlplus.exe",
		commonPaths: []string{
			`C:\Program Files\Oracle\*\bin`,
			`C:\Program Files (x86)\Oracle\*\bin`,
			`C:\oracle\*\bin`,
		},
		category: "Databases",
	},
	{
		name:           "Cassandra",
		executableName: "cassandra.bat",
		commonPaths: []string{
			`C:\Program Files\Apache\cassandra\bin`,
			`C:\Program Files (x86)\Apache\cassandra\bin`,
			`C:\cassandra\bin`,
		},
		category: "Databases",
	},
	{
		name:           "Neo4j",
		executableName: "neo4j.bat",
		commonPaths: []string{
			`C:\Program Files\Neo4j\*\bin`,
			`C:\Program Files (x86)\Neo4j\*\bin`,
			`C:\neo4j\bin`,
		},
		category: "Databases",
	},
	{
		name:           "InfluxDB",
		executableName: "influxd.exe",
		commonPaths: []string{
			`C:\Program Files\InfluxDB\*\bin`,
			`C:\Program Files (x86)\InfluxDB\*\bin`,
			`C:\InfluxDB\bin`,
		},
		category: "Databases",
	},
	{
		name:           "Perl",
		executableName: "perl.exe",
		commonPaths: []string{
			`C:\Perl64\bin`,
			`C:\Strawberry\perl\bin`,
			`C:\Program Files\Perl\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Scala",
		executableName: "scala.bat",
		commonPaths: []string{
			`C:\Program Files\Scala\bin`,
			`C:\Program Files (x86)\Scala\bin`,
			`C:\scala\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Kotlin",
		executableName: "kotlin.bat",
		commonPaths: []string{
			`C:\Program Files\Kotlin\bin`,
			`C:\Program Files (x86)\Kotlin\bin`,
			`C:\kotlin\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Swift",
		executableName: "swift.exe",
		commonPaths: []string{
			`C:\Program Files\Swift\bin`,
			`C:\Program Files (x86)\Swift\bin`,
			`C:\swift\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Haskell",
		executableName: "ghc.exe",
		commonPaths: []string{
			`C:\Program Files\Haskell\bin`,
			`C:\Program Files\GHC\*\bin`,
			`C:\ghc\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Erlang",
		executableName: "erl.exe",
		commonPaths: []string{
			`C:\Program Files\erl*\bin`,
			`C:\Program Files (x86)\erl*\bin`,
			`C:\erlang\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Elixir",
		executableName: "elixir.bat",
		commonPaths: []string{
			`C:\Program Files\Elixir\bin`,
			`C:\Program Files (x86)\Elixir\bin`,
			`C:\elixir\bin`,
		},
		category: "Programming Languages",
	},
	{
		name:           "Terraform",
		executableName: "terraform.exe",
		commonPaths: []string{
			`C:\Program Files\Terraform`,
			`C:\Program Files (x86)\Terraform`,
			`C:\terraform`,
		},
		category: "Infrastructure",
	},
	{
		name:           "Ansible",
		executableName: "ansible.exe",
		commonPaths: []string{
			`C:\Program Files\Ansible`,
			`C:\Program Files (x86)\Ansible`,
			`C:\Users\%USERNAME%\AppData\Roaming\Python\Scripts`,
		},
		category: "Infrastructure",
	},
	{
		name:           "Jenkins",
		executableName: "jenkins.exe",
		commonPaths: []string{
			`C:\Program Files\Jenkins`,
			`C:\Program Files (x86)\Jenkins`,
			`C:\Jenkins`,
		},
		category: "Development Tools",
	},
	{
		name:           "SonarQube",
		executableName: "sonar.bat",
		commonPaths: []string{
			`C:\Program Files\SonarQube\bin`,
			`C:\Program Files (x86)\SonarQube\bin`,
			`C:\sonarqube\bin`,
		},
		category: "Development Tools",
	},
	{
		name:           "Grafana",
		executableName: "grafana-server.exe",
		commonPaths: []string{
			`C:\Program Files\GrafanaLabs\grafana\bin`,
			`C:\Program Files (x86)\GrafanaLabs\grafana\bin`,
			`C:\grafana\bin`,
		},
		category: "Development Tools",
	},
}

// Function to normalize path for comparison
func normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimRight(path, "/")
	return path
}

// Функция для проверки прав администратора
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("WARNING: Program must be run as administrator")
		fmt.Println("Please restart the program with administrator privileges")
		return false
	}
	return true
}

// Function to verify existing environment variables
func verifyEnvironmentSettings(prog Program, path string) bool {
	needsUpdate := false
	
	switch prog.name {
	case "Python":
		userProfile := os.Getenv("USERPROFILE")
		expectedVars := map[string]string{
			"PYTHONHOME":      filepath.Dir(path),
			"PYTHONPATH":      filepath.Dir(path),
			"PYTHONUSERBASE":  filepath.Join(userProfile, ".local"),
			"PYTHONDONTWRITEBYTECODE": "1",
			"PIPENV_VENV_IN_PROJECT": "1",
			"PIP_CACHE_DIR":   filepath.Join(userProfile, ".pip", "cache"),
		}
		
		for name, expectedValue := range expectedVars {
			currentValue := os.Getenv(name)
			if currentValue != expectedValue {
				fmt.Printf("Environment variable %s needs update:\n", name)
				fmt.Printf("  Current: %s\n", currentValue)
				fmt.Printf("  Expected: %s\n", expectedValue)
				needsUpdate = true
			}
		}
		
	// Add similar checks for other tools
	// ... existing code for other tools ...
	}
	
	return needsUpdate
}

func main() {
	// Check administrator privileges at startup
	if !isAdmin() {
		fmt.Println("WARNING: Program must be run as administrator")
		fmt.Println("Please restart the program with administrator privileges")
		fmt.Println("\nPress Enter to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		// Group programs by category
		categories := make(map[string][]Program)
		for _, prog := range allPrograms {
			categories[prog.category] = append(categories[prog.category], prog)
		}

		fmt.Println("\nDevPathPro - Development Environment Manager")
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println("\nMain Menu:")
		fmt.Println("1. Search and Configure Development Tools")
		fmt.Println("2. Verify Existing Configurations")
		fmt.Println("3. View Current Environment Settings")
		fmt.Println("4. Exit")
		fmt.Print("\nSelect an option (1-4): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("\nAvailable tools to search:")
			fmt.Println(strings.Repeat("-", 80))
			
			// Display programs by category
			var numberedPrograms []Program
			currentNumber := 1
			for category, programs := range categories {
				fmt.Printf("\n%s:\n", category)
				for _, prog := range programs {
					fmt.Printf("[%d] %s\n", currentNumber, prog.name)
					numberedPrograms = append(numberedPrograms, prog)
					currentNumber++
				}
			}
			fmt.Println(strings.Repeat("-", 80))
			fmt.Println("\nSelect tools to search (options):")
			fmt.Println("- Enter numbers separated by comma (e.g.: 1,3,5)")
			fmt.Println("- Type 'all' to search for all tools")
			fmt.Println("- Type 'back' to return to main menu")
			fmt.Print("\nYour choice: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if strings.EqualFold(input, "back") {
				continue
			}

			var selectedPrograms []Program
			if input == "all" {
				selectedPrograms = allPrograms
			} else {
				numbers := strings.Split(input, ",")
				for _, num := range numbers {
					num = strings.TrimSpace(num)
					if index, err := strconv.Atoi(num); err == nil && index > 0 && index <= len(numberedPrograms) {
						selectedPrograms = append(selectedPrograms, numberedPrograms[index-1])
					}
				}
			}

			if len(selectedPrograms) == 0 {
				fmt.Println("No valid tools selected. Please try again.")
				continue
			}

			fmt.Printf("\nSelected tools: ")
			for i, prog := range selectedPrograms {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(prog.name)
			}
			fmt.Println("\n")

			// Process selected programs
			configurationChanged := false
			for _, prog := range selectedPrograms {
				fmt.Printf("\n=== Searching for %s ===\n", prog.name)
				paths := findProgram(prog)
				if len(paths) > 0 {
					for _, execPath := range paths {
						if verifyEnvironmentSettings(prog, execPath) {
							setupRegistryForProgram(prog, execPath)
							configurationChanged = true
						}
					}
				} else {
					fmt.Printf("%s not found\n", prog.name)
				}
			}

			if configurationChanged {
				fmt.Println("\nConfiguration changes have been made.")
				fmt.Print("Would you like to restart your computer now? (y/n): ")
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer == "y" {
					fmt.Println("Restarting computer...")
					exec.Command("shutdown", "/r", "/t", "0").Run()
					return
				}
			}

			fmt.Println("\nPress Enter to return to main menu...")
			reader.ReadString('\n')

		case "2":
			fmt.Println("\nVerifying existing configurations...")
			// Add verification logic here
			fmt.Println("\nPress Enter to return to main menu...")
			reader.ReadString('\n')

		case "3":
			fmt.Println("\nCurrent Environment Settings:")
			fmt.Println(strings.Repeat("-", 80))
			for _, env := range os.Environ() {
				fmt.Println(env)
			}
			fmt.Println("\nPress Enter to return to main menu...")
			reader.ReadString('\n')

		case "4":
			fmt.Println("\nExiting program...")
			return

		default:
			fmt.Println("\nInvalid option. Please try again.")
		}
	}
}

// Add logging setup
func init() {
	// Set up logging to file
	logFile, err := os.OpenFile("devpathpro.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Add error type for registry operations
type RegistryError struct {
	Operation string
	Path      string
	Err       error
}

func (e *RegistryError) Error() string {
	return fmt.Sprintf("Registry %s failed for %s: %v", e.Operation, e.Path, e.Err)
}

// Optimize registry search with concurrent operations
func findInRegistry(progName string) []string {
	var results []string
	var mutex sync.Mutex
	var wg sync.WaitGroup
	
	// Registry paths where programs might be installed
	regPaths := []string{
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	// Create buffered channel for errors
	errorChan := make(chan *RegistryError, len(regPaths))

	for _, regPath := range regPaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			
			// Get list of subkeys using reg query command with timeout
			cmd := exec.Command("reg", "query", path, "/s")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			
			// Add timeout for the command
			done := make(chan error, 1)
			go func() {
				output, err := cmd.Output()
				if err != nil {
					done <- err
					return
				}

				// Parse command output
				lines := strings.Split(string(output), "\n")
				var currentKey string
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}

					if strings.HasPrefix(line, "HKEY_") {
						currentKey = line
						continue
					}

					if strings.Contains(strings.ToLower(currentKey), strings.ToLower(progName)) {
						if strings.Contains(line, "Path") || strings.Contains(line, "InstallLocation") {
							parts := strings.SplitN(line, "REG_SZ", 2)
							if len(parts) == 2 {
								path := strings.TrimSpace(parts[1])
								path = strings.Trim(path, "\"")
								if path != "" {
									if _, err := os.Stat(path); err == nil {
										mutex.Lock()
										results = append(results, path)
										mutex.Unlock()
									}
								}
							}
						}
					}
				}
				done <- nil
			}()

			// Wait for command completion or timeout
			select {
			case err := <-done:
				if err != nil {
					errorChan <- &RegistryError{"query", path, err}
				}
			case <-time.After(5 * time.Second):
				cmd.Process.Kill()
				errorChan <- &RegistryError{"timeout", path, fmt.Errorf("operation timed out")}
			}
		}(regPath)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Log any errors that occurred
	for err := range errorChan {
		log.Printf("Registry search error: %v", err)
	}

	return results
}

// Функция для поиска программы с обработкой ошибок и повторными попытками
func findProgram(prog Program) []string {
	var results []string
	var errors []error
	var mutex sync.Mutex

	// First check common paths
	for _, basePath := range prog.commonPaths {
		// Expand environment variables in path
		basePath = os.ExpandEnv(basePath)
		
		// Проверяем существование пути
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue
		}

		var wg sync.WaitGroup
		errorsChan := make(chan error, 100)
		
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					if os.IsPermission(err) {
						return filepath.SkipDir
					}
					errorsChan <- fmt.Errorf("access error to %s: %v", filePath, err)
					return filepath.SkipDir
				}
				
				if !info.IsDir() && strings.EqualFold(filepath.Base(filePath), prog.executableName) {
					mutex.Lock()
					results = append(results, filePath)
					fmt.Printf("Found file: %s\n", filePath)
					mutex.Unlock()
				}
				return nil
			})
			if err != nil {
				errorsChan <- fmt.Errorf("error searching in %s: %v", path, err)
			}
		}(basePath)

		wg.Wait()
		close(errorsChan)
		
		for err := range errorsChan {
			errors = append(errors, err)
		}
	}

	// Try using 'where' command first
	cmd := exec.Command("where", prog.executableName)
	output, err := cmd.Output()

	if err == nil {
		paths := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, path := range paths {
			path = strings.TrimSpace(path)
			if path == "" {
				continue
			}
			found := false
			for _, existingPath := range results {
				if strings.EqualFold(existingPath, path) {
					found = true
					break
				}
			}
			if !found {
				results = append(results, path)
			}
		}
	}

	// If not found anywhere, search all drives
	if len(results) == 0 {
		fmt.Printf("Searching %s on all drives (this may take some time)...\n", prog.name)
		if len(errors) > 0 {
			fmt.Println("Warnings during search:")
			for _, err := range errors {
				fmt.Printf("- %v\n", err)
			}
		}

		// Get all available drives
		drives := getAllDrives()
		resultChan := make(chan string, 100)
		var wg sync.WaitGroup

		for _, drive := range drives {
			wg.Add(1)
			go func(d string) {
				defer wg.Done()
				searchInDrive(d, prog.executableName, resultChan)
			}(drive)
		}

		// Wait for all goroutines to finish
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Collect results
		for path := range resultChan {
			results = append(results, path)
		}
	}

	return results
}

func searchInDrive(drive, executableName string, resultChan chan<- string) {
	fmt.Printf("Searching on drive %s...\n", drive)
	
	// List of directories to skip (только временные и кэш-директории)
	skipDirs := []string{
		"Windows\\Temp", "Temp", "tmp", "cache", "Cache",
		"$Recycle.Bin", "$RECYCLE.BIN", "System Volume Information",
	}

	// Start from drive root
	root := drive + ":\\"
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip only temp and cache directories
		if info.IsDir() {
			baseName := filepath.Base(path)
			for _, skip := range skipDirs {
				if strings.EqualFold(baseName, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if strings.EqualFold(filepath.Base(path), executableName) {
			resultChan <- path
		}
		return nil
	})
}

func getAllDrives() []string {
	var drives []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		_, err := os.Stat(drivePath)
		if err == nil {
			drives = append(drives, string(drive))
		}
	}
	return drives
}

// Функция для настройки дополнительных параметров реестра
func setupRegistryForProgram(prog Program, path string) {
	if !isAdmin() {
		return
	}

	switch prog.name {
	case "Visual Studio":
		vsRoot := filepath.Dir(path)
		// Основные переменные Visual Studio
		envVars := map[string]string{
			"VSINSTALLDIR":    vsRoot,
			"VisualStudioDir": vsRoot,
			"VS_PATH":         vsRoot,
			"VSCMD_VER":      filepath.Base(filepath.Dir(vsRoot)), // версия VS (2019/2022)
		}

		// Добавляем пути для компонентов VS
		addToPath(filepath.Join(vsRoot, "Common7", "IDE"))
		addToPath(filepath.Join(vsRoot, "VC", "Tools", "MSVC", "*", "bin", "Hostx64", "x64"))
		addToPath(filepath.Join(vsRoot, "Common7", "Tools"))
		addToPath(filepath.Join(vsRoot, "Common7", "IDE", "CommonExtensions", "Microsoft", "TeamFoundation", "Team Explorer"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Windows SDK":
		sdkRoot := filepath.Dir(filepath.Dir(filepath.Dir(path)))
		version := filepath.Base(filepath.Dir(path))
		
		// Основные переменные SDK
		envVars := map[string]string{
			"WindowsSdkDir":           sdkRoot,
			"WindowsSdkVersion":       version,
			"WindowsSdkBinPath":       filepath.Join(sdkRoot, "bin", version),
			"WindowsSdkIncludePath":   filepath.Join(sdkRoot, "Include", version),
			"WindowsSdkLibPath":       filepath.Join(sdkRoot, "Lib", version),
			"WindowsSDK_ExecutablePath_x64": filepath.Join(sdkRoot, "bin", version, "x64"),
			"WindowsSDK_ExecutablePath_x86": filepath.Join(sdkRoot, "bin", version, "x86"),
		}

		// Добавляем пути для компонентов SDK
		addToPath(filepath.Join(sdkRoot, "bin", version, "x64"))
		addToPath(filepath.Join(sdkRoot, "bin", version, "x86"))
		addToPath(filepath.Join(sdkRoot, "Tools"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "WDK":
		wdkRoot := filepath.Dir(filepath.Dir(filepath.Dir(path)))
		
		// Основные переменные WDK
		envVars := map[string]string{
			"WINDDK_ROOT":     wdkRoot,
			"WDK_ROOT":        wdkRoot,
			"WDK_BIN_ROOT":    filepath.Join(wdkRoot, "bin"),
			"WDK_INC_ROOT":    filepath.Join(wdkRoot, "inc"),
			"WDK_LIB_ROOT":    filepath.Join(wdkRoot, "lib"),
		}

		// Добавляем пути для инструментов WDK
		addToPath(filepath.Join(wdkRoot, "Tools", "x64"))
		addToPath(filepath.Join(wdkRoot, "Tools", "x86"))
		addToPath(filepath.Join(wdkRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "LLVM":
		llvmRoot := filepath.Dir(path)
		
		// Основные переменные LLVM
		envVars := map[string]string{
			"LLVM_ROOT":       llvmRoot,
			"CLANG_ROOT":      llvmRoot,
			"LLVM_INCLUDE":    filepath.Join(llvmRoot, "include"),
			"LLVM_LIB":        filepath.Join(llvmRoot, "lib"),
		}

		// Добавляем путь для инструментов LLVM
		addToPath(filepath.Join(llvmRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "CMake":
		cmakeRoot := filepath.Dir(path)
		
		// Основные переменные CMake
		envVars := map[string]string{
			"CMAKE_ROOT":      cmakeRoot,
			"CMAKE_HOME":      cmakeRoot,
			"CMAKE_MODULE_PATH": filepath.Join(cmakeRoot, "share", "cmake-*", "Modules"),
		}

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Python":
		pythonRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Python
		envVars := map[string]string{
			"PYTHONHOME":      pythonRoot,
			"PYTHONPATH":      pythonRoot,
			"PYTHON_INCLUDE":  filepath.Join(pythonRoot, "include"),
			"PYTHON_LIB":      filepath.Join(pythonRoot, "libs"),
			"PYTHONUSERBASE":  filepath.Join(userProfile, ".local"),
			"PYTHONDONTWRITEBYTECODE": "1",  // Отключаем создание .pyc файлов
			"PIPENV_VENV_IN_PROJECT": "1",   // Создаем virtualenv в директории проекта
			"PIP_CACHE_DIR":   filepath.Join(userProfile, ".pip", "cache"),
			"VIRTUAL_ENV_DISABLE_PROMPT": "1", // Отключаем изменение PS1
		}

		// Добавляем пути для Python
		addToPath(filepath.Join(pythonRoot, "Scripts"))
		addToPath(pythonRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Make":
		makeRoot := filepath.Dir(path)
		
		// Определяем тип установки (MinGW, MSYS2, Cygwin)
		var installType string
		if strings.Contains(strings.ToLower(makeRoot), "mingw") {
			installType = "MinGW"
		} else if strings.Contains(strings.ToLower(makeRoot), "msys") {
			installType = "MSYS2"
		} else if strings.Contains(strings.ToLower(makeRoot), "cygwin") {
			installType = "Cygwin"
		}

		// Основные переменные Make и MinGW
		envVars := map[string]string{
			"MAKE_ROOT": makeRoot,
			"MAKE_HOME": makeRoot,
		}

		// Дополнительные настройки в зависимости от типа установки
		switch installType {
		case "MinGW":
			mingwRoot := filepath.Dir(makeRoot)
			envVars["MINGW_HOME"] = mingwRoot
			envVars["MINGW_ROOT"] = mingwRoot
			
			// Добавляем пути для MinGW
			addToPath(filepath.Join(mingwRoot, "bin"))
			addToPath(filepath.Join(mingwRoot, "lib"))
			addToPath(filepath.Join(mingwRoot, "include"))

		case "MSYS2":
			msysRoot := filepath.Dir(filepath.Dir(makeRoot))
			envVars["MSYS_HOME"] = msysRoot
			envVars["MSYSTEM"] = "MINGW64"
			
			// Добавляем пути для MSYS2
			addToPath(filepath.Join(msysRoot, "usr", "bin"))
			addToPath(filepath.Join(msysRoot, "mingw64", "bin"))

		case "Cygwin":
			cygwinRoot := filepath.Dir(filepath.Dir(makeRoot))
			envVars["CYGWIN_HOME"] = cygwinRoot
			
			// Добавляем пути для Cygwin
			addToPath(filepath.Join(cygwinRoot, "bin"))
			addToPath(filepath.Join(cygwinRoot, "usr", "bin"))
		}

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "vcpkg":
		vcpkgRoot := filepath.Dir(path)
		
		// Основные переменные vcpkg
		envVars := map[string]string{
			"VCPKG_ROOT":      vcpkgRoot,
			"VCPKG_DEFAULT_TRIPLET": "x64-windows",
		}

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Java":
		javaRoot := filepath.Dir(filepath.Dir(path))
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Java
		envVars := map[string]string{
			"JAVA_HOME":      javaRoot,
			"JDK_HOME":       javaRoot,
			"CLASSPATH":      ".",
			"MAVEN_OPTS":     "-Xmx2048m -XX:MaxPermSize=512m",
			"GRADLE_OPTS":    "-Xmx2048m -Dfile.encoding=UTF-8",
			"JAVA_TOOL_OPTIONS": "-Dfile.encoding=UTF8",
			"SPRING_PROFILES_ACTIVE": "local,development",
			"_JAVA_OPTIONS":  "-Djava.io.tmpdir=" + filepath.Join(userProfile, "temp", "java"),
			"JAVA_OPTS":      "-server -XX:+UseG1GC -XX:MaxGCPauseMillis=200",
		}

		// Добавляем пути для Java
		addToPath(filepath.Join(javaRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Node.js":
		nodeRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Node.js
		envVars := map[string]string{
			"NODE_HOME":      nodeRoot,
			"NODE_PATH":      filepath.Join(nodeRoot, "node_modules"),
			"NPM_CONFIG_PREFIX": filepath.Join(nodeRoot, "npm"),
			"NPM_CONFIG_CACHE": filepath.Join(userProfile, "AppData", "npm-cache"),
			"YARN_CACHE_FOLDER": filepath.Join(userProfile, "AppData", "Local", "Yarn"),
			"NODE_OPTIONS":    "--max-old-space-size=4096",
			"NODE_ENV":       "development",
			"NPM_CONFIG_INIT_AUTHOR_NAME": os.Getenv("USERNAME"),
			"NPM_CONFIG_INIT_LICENSE": "MIT",
			"NPM_CONFIG_SAVE_EXACT": "true",
			"NPM_CONFIG_USERCONFIG": filepath.Join(userProfile, ".npmrc"),
		}

		// Добавляем пути для npm и yarn
		addToPath(filepath.Join(nodeRoot, "node_modules", ".bin"))
		addToPath(filepath.Join(os.Getenv("APPDATA"), "npm"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Go":
		goRoot := filepath.Dir(filepath.Dir(path))
		userProfile := os.Getenv("USERPROFILE")
		goPath := filepath.Join(userProfile, "go")
		
		// Основные переменные Go
		envVars := map[string]string{
			"GOROOT":         goRoot,
			"GOPATH":         goPath,
			"GO111MODULE":    "on",
			"GOCACHE":        filepath.Join(userProfile, "AppData", "Local", "go-build"),
			"GOENV":          filepath.Join(userProfile, "AppData", "Roaming", "go", "env"),
			"GOTOOLDIR":      filepath.Join(goRoot, "pkg", "tool", "windows_amd64"),
			"GOPROXY":        "https://proxy.golang.org,direct",
			"GOSUMDB":        "sum.golang.org",
			"GOPRIVATE":      "",
			"GOMODCACHE":     filepath.Join(goPath, "pkg", "mod"),
			"GOTMPDIR":       filepath.Join(userProfile, "temp", "go"),
			"GOFLAGS":        "-mod=vendor",
		}

		// Добавляем пути для Go
		addToPath(filepath.Join(goRoot, "bin"))
		addToPath(filepath.Join(goPath, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Docker":
		dockerRoot := filepath.Dir(filepath.Dir(path))
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Docker
		envVars := map[string]string{
			"DOCKER_HOME":    dockerRoot,
			"DOCKER_CONFIG": filepath.Join(userProfile, ".docker"),
			"DOCKER_HOST":   "tcp://localhost:2375",
			"COMPOSE_CONVERT_WINDOWS_PATHS": "1",
			"DOCKER_BUILDKIT": "1",
			"COMPOSE_HTTP_TIMEOUT": "300",
			"DOCKER_CLI_EXPERIMENTAL": "enabled",
			"DOCKER_SCAN_SUGGEST": "false",
			"DOCKER_CONTENT_TRUST": "1",
			"DOCKER_REGISTRY": "docker.io",
			"DOCKER_DEFAULT_PLATFORM": "windows/amd64",
			"COMPOSE_PROJECT_NAME": "${PWD##*/}",
		}

		// Добавляем пути для docker-compose и дополнительных утилит
		addToPath(filepath.Join(dockerRoot, "cli-plugins"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Kubernetes":
		k8sRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Kubernetes
		envVars := map[string]string{
			"KUBECONFIG":      filepath.Join(userProfile, ".kube", "config"),
			"MINIKUBE_HOME":   filepath.Join(userProfile, ".minikube"),
			"MINIKUBE_WANTUPDATENOTIFICATION": "false",
			"KUBE_EDITOR":     "code --wait",
			"KIND_HOME":       filepath.Join(userProfile, ".kind"),
			"HELM_HOME":       filepath.Join(userProfile, ".helm"),
			"KUBE_CLUSTER_NAME": "dev-cluster",
			"KUBE_NAMESPACE":  "default",
			"KUBE_CONTEXT":   "minikube",
			"KUBE_PS1_SYMBOL_ENABLE": "true",
			"KUBE_PS1_SEPARATOR": "|",
			"KUBE_PS1_PREFIX": "⎈ ",
		}

		// Добавляем пути для kubectl и minikube
		addToPath(k8sRoot)
		addToPath(filepath.Join(userProfile, ".krew", "bin")) // Для плагинов kubectl

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "PostgreSQL":
		pgRoot := filepath.Dir(path)
		
		// Основные переменные PostgreSQL
		envVars := map[string]string{
			"PGDATA":          filepath.Join(pgRoot, "data"),
			"PGHOST":          "localhost",
			"PGPORT":          "5432",
			"PGLOCALEDIR":     filepath.Join(pgRoot, "share", "locale"),
			"PGDATABASE":      "postgres",
			"PGUSER":          "postgres",
			"PGPASSFILE":      filepath.Join(os.Getenv("APPDATA"), "postgresql", "pgpass.conf"),
			"PGSSLMODE":       "prefer",
			"PGCLIENTENCODING": "UTF8",
			"PGTZ":            "UTC",
			"PGOPTIONS":       "--client-min-messages=warning",
		}

		// Добавляем пути для PostgreSQL
		addToPath(filepath.Join(pgRoot, "bin"))
		addToPath(filepath.Join(pgRoot, "lib"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "MySQL":
		mysqlRoot := filepath.Dir(path)
		
		// Основные переменные MySQL
		envVars := map[string]string{
			"MYSQL_HOME":      mysqlRoot,
			"MYSQL_DATA":      filepath.Join(mysqlRoot, "data"),
			"MYSQL_PORT":      "3306",
			"MYSQL_HOST":      "localhost",
			"MYSQL_CHARSET":   "utf8mb4",
		}

		// Добавляем пути для MySQL
		addToPath(filepath.Join(mysqlRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "MongoDB":
		mongoRoot := filepath.Dir(path)
		
		// Основные переменные MongoDB
		envVars := map[string]string{
			"MONGO_HOME":      mongoRoot,
			"MONGO_DATA":      filepath.Join(mongoRoot, "data", "db"),
			"MONGO_LOG":       filepath.Join(mongoRoot, "log"),
			"MONGO_PORT":      "27017",
			"MONGO_HOST":      "localhost",
			"MONGOSH_EDITOR":  "code --wait",
			"MONGOSH_HISTORYFILE": filepath.Join(os.Getenv("APPDATA"), "mongodb", "mongosh_history"),
			"MONGOSH_TIMEOUT_MS": "30000",
			"MONGO_URL":       "mongodb://localhost:27017",
			"MONGO_INITDB_ROOT_USERNAME": "admin",
		}

		// Добавляем пути для MongoDB
		addToPath(filepath.Join(mongoRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Redis":
		redisRoot := filepath.Dir(path)
		
		// Основные переменные Redis
		envVars := map[string]string{
			"REDIS_HOME":     redisRoot,
			"REDIS_PORT":     "6379",
			"REDIS_HOST":     "localhost",
			"REDIS_CONFIG":   filepath.Join(redisRoot, "redis.windows.conf"),
			"REDIS_DATA":     filepath.Join(redisRoot, "data"),
			"REDIS_LOG":      filepath.Join(redisRoot, "logs", "redis.log"),
			"REDIS_MAX_MEMORY": "2gb",
			"REDIS_MAXCLIENTS": "10000",
			"REDIS_TIMEOUT":   "300",
			"REDIS_TLS_PORT":  "6380",
		}

		// Добавляем путь для Redis
		addToPath(redisRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Elasticsearch":
		esRoot := filepath.Dir(filepath.Dir(path))
		
		// Основные переменные Elasticsearch
		envVars := map[string]string{
			"ES_HOME":        esRoot,
			"ES_PATH_CONF":   filepath.Join(esRoot, "config"),
			"ES_JAVA_HOME":   os.Getenv("JAVA_HOME"),
			"ES_TMPDIR":      filepath.Join(esRoot, "tmp"),
			"ES_PATH_DATA":   filepath.Join(esRoot, "data"),
			"ES_PATH_LOGS":   filepath.Join(esRoot, "logs"),
			"ES_JAVA_OPTS":   "-Xms2g -Xmx2g",
			"ES_HEAP_SIZE":   "2g",
			"ES_NODE_NAME":   "node-1",
			"ES_CLUSTER_NAME": "dev-cluster",
			"ES_NETWORK_HOST": "localhost",
			"ES_HTTP_PORT":   "9200",
			"ES_DISCOVERY_TYPE": "single-node",
		}

		// Добавляем пути для Elasticsearch
		addToPath(filepath.Join(esRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case ".NET Core":
		dotnetRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные .NET Core
		envVars := map[string]string{
			"DOTNET_ROOT":     dotnetRoot,
			"DOTNET_CLI_HOME": userProfile,
			"DOTNET_NOLOGO":   "1", // Отключаем приветственное сообщение
			"DOTNET_CLI_TELEMETRY_OPTOUT": "1", // Отключаем телеметрию
			"NUGET_PACKAGES":  filepath.Join(userProfile, ".nuget", "packages"),
		}

		// Добавляем пути для .NET Core
		addToPath(dotnetRoot)
		addToPath(filepath.Join(userProfile, ".dotnet", "tools"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Ruby":
		rubyRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Ruby
		envVars := map[string]string{
			"RUBY_HOME":      rubyRoot,
			"GEM_HOME":       filepath.Join(userProfile, ".gem", "ruby"),
			"GEM_PATH":       filepath.Join(userProfile, ".gem", "ruby"),
			"BUNDLE_USER_HOME": filepath.Join(userProfile, ".bundle"),
			"BUNDLE_PATH":    filepath.Join(userProfile, ".bundle", "vendor"),
		}

		// Добавляем пути для Ruby и Gems
		addToPath(filepath.Join(rubyRoot, "bin"))
		addToPath(filepath.Join(userProfile, ".gem", "ruby", "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Rust":
		rustRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Rust
		envVars := map[string]string{
			"RUST_HOME":      rustRoot,
			"RUSTUP_HOME":    filepath.Join(userProfile, ".rustup"),
			"CARGO_HOME":     filepath.Join(userProfile, ".cargo"),
			"RUSTC_WRAPPER":  "sccache",
			"RUST_BACKTRACE": "1",
			"RUST_LOG":       "info",
			"CARGO_TARGET_DIR": filepath.Join(userProfile, ".cargo", "target"),
			"CARGO_INCREMENTAL": "1",
			"CARGO_NET_RETRY": "5",
			"RUSTFLAGS":      "-C target-cpu=native",
			"RUST_MIN_STACK": "8388608",
		}

		// Добавляем пути для Rust
		addToPath(filepath.Join(rustRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Podman":
		podmanRoot := filepath.Dir(path)
		
		// Основные переменные Podman
		envVars := map[string]string{
			"PODMAN_HOME":     podmanRoot,
			"CONTAINERS_CONF": filepath.Join(podmanRoot, "containers.conf"),
			"PODMAN_SOCKET":   `\\.\pipe\podman-machine-default`,
			"PODMAN_MACHINE_DEFAULT": "1",
		}

		// Добавляем пути для Podman
		addToPath(podmanRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Helm":
		helmRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Helm
		envVars := map[string]string{
			"HELM_HOME":     filepath.Join(userProfile, ".helm"),
			"HELM_CONFIG":   filepath.Join(userProfile, ".helm", "config"),
			"HELM_CACHE_HOME": filepath.Join(userProfile, ".helm", "cache"),
			"HELM_DATA_HOME": filepath.Join(userProfile, ".helm", "data"),
		}

		// Добавляем пути для Helm
		addToPath(helmRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Skaffold":
		skaffoldRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Skaffold
		envVars := map[string]string{
			"SKAFFOLD_HOME":   filepath.Join(userProfile, ".skaffold"),
			"SKAFFOLD_CONFIG": filepath.Join(userProfile, ".skaffold", "config"),
			"SKAFFOLD_UPDATE_CHECK": "false",
		}

		// Добавляем пути для Skaffold
		addToPath(skaffoldRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Perl":
		perlRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Perl
		envVars := map[string]string{
			"PERL_HOME":      perlRoot,
			"PERL5LIB":       filepath.Join(perlRoot, "lib"),
			"PERL_LOCAL_LIB_ROOT": filepath.Join(userProfile, "perl5"),
			"PERL_MB_OPT":    "--install_base " + filepath.Join(userProfile, "perl5"),
			"PERL_MM_OPT":    "INSTALL_BASE=" + filepath.Join(userProfile, "perl5"),
			"PERL_UNICODE":   "AS",
			"PERL_CPANM_HOME": filepath.Join(userProfile, ".cpanm"),
			"PERL_CRITIC_THEME": "stern",
		}

		// Добавляем пути для Perl
		addToPath(filepath.Join(perlRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Scala":
		scalaRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Scala
		envVars := map[string]string{
			"SCALA_HOME":     scalaRoot,
			"SBT_HOME":       filepath.Join(scalaRoot, "sbt"),
			"SBT_OPTS":       "-Xmx2G -XX:+UseG1GC -Xss2M",
			"COURSIER_CACHE": filepath.Join(userProfile, ".coursier"),
			"SCALA_CACHE":    filepath.Join(userProfile, ".scala"),
			"SCALA_VERSION":  "2.13.8",
			"SBT_CREDENTIALS": filepath.Join(userProfile, ".sbt", ".credentials"),
		}

		// Добавляем пути для Scala
		addToPath(filepath.Join(scalaRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Kotlin":
		kotlinRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Kotlin
		envVars := map[string]string{
			"KOTLIN_HOME":    kotlinRoot,
			"KOTLINC_HOME":   kotlinRoot,
			"KOTLIN_OPTS":    "-Xmx2G -XX:+UseG1GC",
			"KOTLIN_DAEMON_CLIENT_OPTIONS": "-Xmx2G",
			"KOTLIN_COMPILER_VERSION": "1.8.0",
			"KOTLIN_COMPILER_CACHE": filepath.Join(userProfile, ".kotlin", "cache"),
		}

		// Добавляем пути для Kotlin
		addToPath(filepath.Join(kotlinRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Swift":
		swiftRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Swift
		envVars := map[string]string{
			"SWIFT_HOME":     swiftRoot,
			"SWIFT_TOOLS_PATH": filepath.Join(swiftRoot, "usr", "bin"),
			"SWIFT_BUILD_PATH": filepath.Join(userProfile, ".swift", "build"),
			"SWIFT_PACKAGE_PATH": filepath.Join(userProfile, ".swift", "packages"),
			"SWIFT_CACHE_PATH": filepath.Join(userProfile, ".swift", "cache"),
			"SWIFT_USE_LIBDISPATCH": "1",
		}

		// Добавляем пути для Swift
		addToPath(filepath.Join(swiftRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Haskell":
		ghcRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Haskell
		envVars := map[string]string{
			"GHC_HOME":       ghcRoot,
			"CABAL_HOME":     filepath.Join(userProfile, "AppData", "Roaming", "cabal"),
			"STACK_ROOT":     filepath.Join(userProfile, "AppData", "Roaming", "stack"),
			"GHC_PACKAGE_PATH": filepath.Join(ghcRoot, "lib", "package.conf.d"),
			"HASKELL_DIST_DIR": filepath.Join(userProfile, ".cabal", "dist"),
			"STACK_WORK_DIR":  ".stack-work",
		}

		// Добавляем пути для Haskell
		addToPath(filepath.Join(ghcRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Erlang":
		erlRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Erlang
		envVars := map[string]string{
			"ERLANG_HOME":    erlRoot,
			"ERL_TOP":        erlRoot,
			"REBAR_CACHE_DIR": filepath.Join(userProfile, ".cache", "rebar3"),
			"ERL_LIBS":       filepath.Join(erlRoot, "lib"),
			"ERL_CRASH_DUMP": filepath.Join(userProfile, "erl_crash.dump"),
			"ERL_AFLAGS":     "-kernel shell_history enabled",
		}

		// Добавляем пути для Erlang
		addToPath(filepath.Join(erlRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Elixir":
		elixirRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Elixir
		envVars := map[string]string{
			"ELIXIR_HOME":    elixirRoot,
			"MIX_HOME":       filepath.Join(userProfile, ".mix"),
			"HEX_HOME":       filepath.Join(userProfile, ".hex"),
			"MIX_ARCHIVES":   filepath.Join(userProfile, ".mix", "archives"),
			"MIX_DEBUG":      "1",
			"MIX_ENV":        "dev",
		}

		// Добавляем пути для Elixir
		addToPath(filepath.Join(elixirRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Terraform":
		terraformRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Terraform
		envVars := map[string]string{
			"TERRAFORM_HOME": terraformRoot,
			"TF_CLI_CONFIG_FILE": filepath.Join(userProfile, ".terraformrc"),
			"TF_DATA_DIR":    filepath.Join(userProfile, ".terraform.d"),
			"TF_PLUGIN_CACHE_DIR": filepath.Join(userProfile, ".terraform.d", "plugin-cache"),
			"TF_IN_AUTOMATION": "true",
			"TF_LOG":         "INFO",
			"TF_LOG_PATH":    filepath.Join(userProfile, "terraform.log"),
			"TF_WORKSPACE":   "default",
			"TF_VAR_user_home": userProfile,
		}

		// Добавляем пути для Terraform
		addToPath(terraformRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Ansible":
		ansibleRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Ansible
		envVars := map[string]string{
			"ANSIBLE_HOME":   ansibleRoot,
			"ANSIBLE_CONFIG": filepath.Join(userProfile, ".ansible.cfg"),
			"ANSIBLE_INVENTORY": filepath.Join(userProfile, "ansible_hosts"),
			"ANSIBLE_LIBRARY": filepath.Join(ansibleRoot, "library"),
			"ANSIBLE_ROLES_PATH": filepath.Join(userProfile, ".ansible", "roles"),
			"ANSIBLE_COLLECTIONS_PATH": filepath.Join(userProfile, ".ansible", "collections"),
			"ANSIBLE_LOCAL_TEMP": filepath.Join(userProfile, ".ansible", "tmp"),
			"ANSIBLE_SSH_CONTROL_PATH": filepath.Join(userProfile, ".ansible", "cp", "%%h-%%p-%%r"),
		}

		// Добавляем пути для Ansible
		addToPath(ansibleRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Jenkins":
		jenkinsRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Jenkins
		envVars := map[string]string{
			"JENKINS_HOME":   jenkinsRoot,
			"JENKINS_URL":    "http://localhost:8080",
			"JENKINS_USER_ID": os.Getenv("USERNAME"),
			"JENKINS_UC":     "https://updates.jenkins.io",
			"JENKINS_UC_EXPERIMENTAL": "https://updates.jenkins.io/experimental",
			"JENKINS_INCREMENTALS_REPO_MIRROR": "https://repo.jenkins-ci.org/incrementals",
			"JENKINS_WAR":    filepath.Join(jenkinsRoot, "jenkins.war"),
			"JENKINS_LOG_DIR": filepath.Join(userProfile, ".jenkins", "logs"),
		}

		// Добавляем пути для Jenkins
		addToPath(jenkinsRoot)

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "SonarQube":
		sonarRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные SonarQube
		envVars := map[string]string{
			"SONAR_HOME":     sonarRoot,
			"SONAR_SCANNER_HOME": filepath.Join(sonarRoot, "sonar-scanner"),
			"SONAR_RUNNER_HOME": filepath.Join(sonarRoot, "sonar-runner"),
			"SONAR_USER_HOME": filepath.Join(userProfile, ".sonar"),
			"SONAR_HOST_URL": "http://localhost:9000",
			"SONAR_TOKEN":    "",  // Должен быть установлен пользователем
			"SONAR_SCANNER_OPTS": "-Xmx2048m",
			"SONAR_JAVA_PATH": filepath.Join(os.Getenv("JAVA_HOME"), "bin", "java.exe"),
		}

		// Добавляем пути для SonarQube
		addToPath(filepath.Join(sonarRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}

	case "Grafana":
		grafanaRoot := filepath.Dir(path)
		userProfile := os.Getenv("USERPROFILE")
		
		// Основные переменные Grafana
		envVars := map[string]string{
			"GF_HOME":        grafanaRoot,
			"GF_PATHS_DATA":  filepath.Join(grafanaRoot, "data"),
			"GF_PATHS_LOGS":  filepath.Join(grafanaRoot, "logs"),
			"GF_PATHS_PLUGINS": filepath.Join(grafanaRoot, "plugins"),
			"GF_PATHS_PROVISIONING": filepath.Join(grafanaRoot, "conf", "provisioning"),
			"GF_SERVER_HTTP_PORT": "3000",
			"GF_SECURITY_ADMIN_USER": "admin",
			"GF_INSTALL_PLUGINS": "grafana-clock-panel,grafana-simple-json-datasource",
			"GF_LOG_MODE":    "console file",
			"GF_LOG_LEVEL":   "info",
			"GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS": "*",
			"GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH": filepath.Join(userProfile, ".grafana", "dashboards", "home.json"),
		}

		// Добавляем пути для Grafana
		addToPath(filepath.Join(grafanaRoot, "bin"))

		// Устанавливаем переменные окружения
		for name, value := range envVars {
			cmd := exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				"/v", name, "/t", "REG_SZ", "/d", value, "/f")
			cmd.Run()
		}
	}

	// Уведомляем систему об изменениях
	notifyEnvironmentChange()
}

// Функция для уведомления системы об изменениях переменных окружения
func notifyEnvironmentChange() {
	dll, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		return
	}
	
	proc, err := dll.FindProc("SendMessageTimeoutW")
	if err != nil {
		return
	}

	msgPtr, _ := syscall.UTF16PtrFromString("Environment")
	proc.Call(
		uintptr(0xFFFF), // HWND_BROADCAST
		uintptr(0x001A), // WM_SETTINGCHANGE
		0,
		uintptr(unsafe.Pointer(msgPtr)),
		uintptr(0x2),    // SMTO_ABORTIFHUNG
		uintptr(5000),
		0,
	)
}

func addToPath(newPath string) {
	if !isAdmin() {
		fmt.Println("WARNING: Program must be run as administrator to modify PATH")
		return
	}

	// Получаем текущий PATH из реестра
	cmd := exec.Command("reg", "query", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, "/v", "Path")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error reading PATH from registry: %v\n", err)
		return
	}

	// Парсим вывод команды reg query
	lines := strings.Split(string(output), "\n")
	var currentPath string
	for _, line := range lines {
		if strings.Contains(line, "REG_EXPAND_SZ") || strings.Contains(line, "REG_SZ") {
			parts := strings.SplitN(line, "REG_", 2)
			if len(parts) == 2 {
				currentPath = strings.TrimSpace(strings.SplitN(parts[1], "    ", 2)[1])
			}
		}
	}

	// Проверяем, не существует ли уже путь
	paths := strings.Split(currentPath, ";")
	normalizedNewPath := normalizePath(newPath)
	
	for _, path := range paths {
		if path != "" && normalizePath(path) == normalizedNewPath {
			fmt.Printf("Path %s already exists in PATH\n", path)
			return
		}
	}

	// Добавляем новый путь
	newPathValue := currentPath
	if !strings.HasSuffix(newPathValue, ";") {
		newPathValue += ";"
	}
	newPathValue += newPath

	// Сохраняем обновленный PATH в реестр
	cmd = exec.Command("reg", "add", `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, "/v", "Path", "/t", "REG_EXPAND_SZ", "/d", newPathValue, "/f")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error saving PATH to registry: %v\n", err)
		return
	}

	// Уведомляем систему об изменении переменной окружения
	dll, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		fmt.Printf("Error loading user32.dll: %v\n", err)
		return
	}
	
	proc, err := dll.FindProc("SendMessageTimeoutW")
	if err != nil {
		fmt.Printf("Error finding SendMessageTimeoutW: %v\n", err)
		return
	}

	const HWND_BROADCAST = 0xFFFF
	const WM_SETTINGCHANGE = 0x001A
	msgPtr, _ := syscall.UTF16PtrFromString("Environment")
	
	ret, _, _ := proc.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(msgPtr)),
		uintptr(0x2), // SMTO_ABORTIFHUNG
		uintptr(5000),
		0,
	)

	if ret == 0 {
		fmt.Println("Warning: Failed to notify system about PATH change")
	}

	fmt.Printf("Path %s successfully added to system PATH\n", newPath)
	fmt.Println("Note: You may need to restart programs or terminal to apply changes")

	// Configure additional registry settings for the program
	for _, prog := range allPrograms {
		if strings.Contains(newPath, prog.executableName) {
			setupRegistryForProgram(prog, newPath)
			break
		}
	}
} 