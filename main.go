package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

// Program structure holds information about a development tool
type Program struct {
	name           string
	executableName string
	commonPaths    []string
	category       string
}

// Function to normalize path for comparison
func normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimRight(path, "/")
	return path
}

func main() {
	// Define available programs
	allPrograms := []Program{
		{
			name:           "CMake",
			executableName: "cmake.exe",
			commonPaths: []string{
				`C:\Program Files\CMake`,
				`C:\Program Files (x86)\CMake`,
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
				`C:\Program Files (x86)\Microsoft Visual Studio`,
				`C:\Program Files\Microsoft Visual Studio`,
				`C:\Windows\Microsoft.NET\Framework`,
				`C:\Windows\Microsoft.NET\Framework64`,
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
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		// Group programs by category
		categories := make(map[string][]Program)
		for _, prog := range allPrograms {
			categories[prog.category] = append(categories[prog.category], prog)
		}

		fmt.Println("\nPath Finder - Main Menu")
		fmt.Println(strings.Repeat("=", 80))
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
		fmt.Println("- Type 'exit' to quit the program")
		fmt.Print("\nYour choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "exit") {
			fmt.Println("\nExiting program...")
			return
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

		// Get current PATH
		currentPath := os.Getenv("PATH")
		pathList := strings.Split(currentPath, ";")
		
		// Create map for original paths
		pathMap := make(map[string]string)
		normalizedPathList := make([]string, 0, len(pathList))
		for _, path := range pathList {
			if path != "" {
				normalized := normalizePath(path)
				normalizedPathList = append(normalizedPathList, normalized)
				pathMap[normalized] = path
			}
		}

		// Process selected programs
		for _, prog := range selectedPrograms {
			fmt.Printf("\n=== Searching for %s ===\n", prog.name)
			paths := findProgram(prog)
			if len(paths) > 0 {
				var dirPaths []string
				execPathMap := make(map[string]string)
				inPathMap := make(map[string]bool)
				originalPathMap := make(map[string]string)
				
				for _, execPath := range paths {
					dirPath := filepath.Dir(execPath)
					dirPaths = append(dirPaths, dirPath)
					execPathMap[dirPath] = execPath
					
					// Check if path is in PATH
					normalizedDirPath := normalizePath(dirPath)
					for _, envPath := range normalizedPathList {
						if normalizedDirPath == envPath {
							inPathMap[dirPath] = true
							originalPathMap[dirPath] = pathMap[envPath]
							break
						}
					}
				}

				fmt.Printf("\nFound paths for %s:\n", prog.name)
				
				// Display existing paths first
				existingPaths := false
				for _, dirPath := range dirPaths {
					if inPathMap[dirPath] {
						if !existingPaths {
							fmt.Println("\nExisting paths in PATH:")
							fmt.Println(strings.Repeat("-", 80))
							existingPaths = true
						}
						fmt.Printf("Directory: %s\n", dirPath)
						fmt.Printf("  Current PATH: %s\n", originalPathMap[dirPath])
						fmt.Printf("  Executable: %s\n", execPathMap[dirPath])
						fmt.Println(strings.Repeat("-", 80))
					}
				}

				// Then display new paths
				newPaths := false
				fmt.Println("\nAvailable paths for addition:")
				fmt.Println(strings.Repeat("-", 80))
				for i, dirPath := range dirPaths {
					if !inPathMap[dirPath] {
						newPaths = true
						fmt.Printf("[%d] Directory: %s\n", i+1, dirPath)
						fmt.Printf("    Executable: %s\n", execPathMap[dirPath])
						fmt.Println(strings.Repeat("-", 80))
					}
				}

				if !newPaths {
					fmt.Println("\nAll found paths are already added to PATH")
					continue
				}
				
				fmt.Println("\nSelect action:")
				fmt.Println("  • Enter numbers separated by comma (e.g.: 1,2)")
				fmt.Println("  • Type 'all' to add all new paths")
				fmt.Println("  • Type 'skip' to skip")
				fmt.Print("\nYour choice: ")
				
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if input == "skip" {
					fmt.Println("Skipping addition of paths for", prog.name)
					continue
				}

				var pathsToAdd []string
				if input == "all" {
					// Добавляем только те пути, которых еще нет в PATH
					for _, dirPath := range dirPaths {
						if !inPathMap[dirPath] {
							pathsToAdd = append(pathsToAdd, dirPath)
						}
					}
				} else {
					// Разбираем введенные номера
					numbers := strings.Split(input, ",")
					for _, num := range numbers {
						num = strings.TrimSpace(num)
						if index, err := fmt.Sscanf(num, "%d"); err == nil && index > 0 && index <= len(dirPaths) {
							dirPath := dirPaths[index-1]
							if !inPathMap[dirPath] {
								pathsToAdd = append(pathsToAdd, dirPath)
							} else {
								fmt.Printf("Path %s is already in PATH, skipping...\n", dirPath)
							}
						}
					}
				}

				// Добавляем выбранные пути
				if len(pathsToAdd) > 0 {
					for _, path := range pathsToAdd {
						fmt.Printf("\nAdding path %s to PATH...\n", path)
						addToPath(path)
					}
				} else {
					fmt.Println("No new paths to add to PATH")
				}
			} else {
				fmt.Println(prog.name, "not found")
			}
			fmt.Println() // Пустая строка для разделения
		}

		fmt.Println("\nPress Enter to return to main menu...")
		reader.ReadString('\n')
	}
}

// Функция для проверки прав администратора
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// Функция для поиска программы в реестре Windows
func findInRegistry(progName string) []string {
	var results []string
	
	// Пути в реестре, где могут быть установлены программы
	regPaths := []string{
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, regPath := range regPaths {
		// Получаем список подключей через команду reg query
		cmd := exec.Command("reg", "query", regPath, "/s")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		// Разбираем вывод команды
		lines := strings.Split(string(output), "\n")
		var currentKey string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Если строка начинается с HKEY_, это новый ключ
			if strings.HasPrefix(line, "HKEY_") {
				currentKey = line
				continue
			}

			// Ищем значения Path и InstallLocation
			if strings.Contains(strings.ToLower(currentKey), strings.ToLower(progName)) {
				if strings.Contains(line, "Path") || strings.Contains(line, "InstallLocation") {
					parts := strings.SplitN(line, "REG_SZ", 2)
					if len(parts) == 2 {
						path := strings.TrimSpace(parts[1])
						// Убираем кавычки, если они есть
						path = strings.Trim(path, "\"")
						if path != "" {
							if _, err := os.Stat(path); err == nil {
								results = append(results, path)
							}
						}
					}
				}
			}
		}
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
					results = append(results, filePath) // Store full path to executable
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
	
	// List of directories to skip
	skipDirs := []string{
		"Windows", "Program Files", "Program Files (x86)",
		"ProgramData", "Users",
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

		// Skip system directories
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
			resultChan <- path // Return full path to executable
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

func addToPath(newPath string) {
	currentPath := os.Getenv("PATH")
	paths := strings.Split(currentPath, ";")
	
	normalizedNewPath := normalizePath(newPath)
	
	// Check if path already exists
	for _, path := range paths {
		if path != "" && normalizePath(path) == normalizedNewPath {
			fmt.Printf("Path %s is already in PATH\n", path)
			return
		}
	}
	
	// Add new path
	newPathValue := currentPath + ";" + newPath
	
	// Load Windows API
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		fmt.Printf("Error loading kernel32.dll: %v\n", err)
		return
	}
	proc, err := dll.FindProc("SetEnvironmentVariableW")
	if err != nil {
		fmt.Printf("Error finding SetEnvironmentVariableW: %v\n", err)
		return
	}
	
	pathPtr, err := syscall.UTF16PtrFromString("PATH")
	if err != nil {
		fmt.Printf("Error converting PATH: %v\n", err)
		return
	}
	valuePtr, err := syscall.UTF16PtrFromString(newPathValue)
	if err != nil {
		fmt.Printf("Error converting PATH value: %v\n", err)
		return
	}
	
	r1, _, err := proc.Call(uintptr(unsafe.Pointer(pathPtr)), uintptr(unsafe.Pointer(valuePtr)))
	if r1 == 0 {
		fmt.Printf("Error setting PATH: %v\n", err)
		return
	}
	
	fmt.Printf("Path %s successfully added to PATH\n", newPath)
} 