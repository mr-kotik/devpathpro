package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"devpathpro/pkg/config"
)

// FindProgram searches for a program in the system
func FindProgram(prog config.Program) []string {
	var results []string
	var errors []error
	var mutex sync.Mutex

	// First check common paths
	for _, basePath := range prog.CommonPaths {
		// Expand environment variables in path
		basePath = os.ExpandEnv(basePath)
		
		// Check path existence
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
				
				if !info.IsDir() && strings.EqualFold(filepath.Base(filePath), prog.ExecutableName) {
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

	// Try using 'where' command
	cmd := exec.Command("where", prog.ExecutableName)
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

	return results
}

// GetAllDrives returns a list of available drives
func GetAllDrives() []string {
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

// SearchInDrive searches for a program in a specific drive
func SearchInDrive(drive, executableName string, resultChan chan<- string) {
	fmt.Printf("Searching on drive %s...\n", drive)
	
	// List of directories to skip
	skipDirs := []string{
		"Windows\\Temp", "Temp", "tmp", "cache", "Cache",
		"$Recycle.Bin", "$RECYCLE.BIN", "System Volume Information",
	}

	root := drive + ":\\"
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return filepath.SkipDir
			}
			return nil
		}

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

// SelectPath asks user to choose a path when multiple installations are found
func SelectPath(paths []string, programName string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("no paths provided")
	}
	
	if len(paths) == 1 {
		return paths[0], nil
	}

	fmt.Printf("\nMultiple installations of %s found:\n", programName)
	for i, path := range paths {
		fmt.Printf("[%d] %s\n", i+1, path)
	}
	fmt.Print("\nSelect path to use (enter number): ")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if num, err := strconv.Atoi(input); err == nil && num > 0 && num <= len(paths) {
			return paths[num-1], nil
		}
		fmt.Print("Invalid selection. Please try again: ")
	}
}

// ConfigureSelectedPath configures the program with the selected path
func ConfigureSelectedPath(prog config.Program, selectedPath string) error {
	// Set environment variables if specified
	if prog.EnvVar != "" {
		if err := os.Setenv(prog.EnvVar, filepath.Dir(selectedPath)); err != nil {
			return fmt.Errorf("failed to set %s: %v", prog.EnvVar, err)
		}
	}

	// Add to PATH if needed
	pathVar := os.Getenv("PATH")
	if !strings.Contains(pathVar, filepath.Dir(selectedPath)) {
		newPath := pathVar + ";" + filepath.Dir(selectedPath)
		if err := os.Setenv("PATH", newPath); err != nil {
			return fmt.Errorf("failed to update PATH: %v", err)
		}
	}

	return nil
} 