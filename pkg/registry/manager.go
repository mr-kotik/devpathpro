package registry

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

const (
	envKey = `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
)

// IsAdmin checks if the program has administrator privileges
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("Administrator privileges required")
		fmt.Println("Please restart the program with administrator privileges")
		return false
	}
	return true
}

// AddToPath adds a new path to the system PATH
func AddToPath(newPath string) error {
	if !IsAdmin() {
		return fmt.Errorf("administrator privileges required")
	}

	// Get current PATH
	cmd := exec.Command(`C:\Windows\System32\reg.exe`, "query", envKey, "/v", "Path")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error reading PATH: %v", err)
	}

	// Parse output
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

	// Check if path already exists
	paths := strings.Split(currentPath, ";")
	normalizedNewPath := normalizePath(newPath)
	
	for _, path := range paths {
		if path != "" && normalizePath(path) == normalizedNewPath {
			return nil // Path already exists
		}
	}

	// Add new path
	if !strings.HasSuffix(currentPath, ";") {
		currentPath += ";"
	}
	currentPath += newPath

	// Update PATH in registry
	cmd = exec.Command(`C:\Windows\System32\reg.exe`, "add", envKey,
		"/v", "Path", "/t", "REG_EXPAND_SZ", "/d", currentPath, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error updating PATH: %v", err)
	}

	NotifyEnvironmentChange()
	return nil
}

// SetEnvironmentVariable sets a system environment variable
func SetEnvironmentVariable(name, value string) error {
	if !IsAdmin() {
		return fmt.Errorf("administrator privileges required")
	}

	cmd := exec.Command(`C:\Windows\System32\reg.exe`, "add", envKey,
		"/v", name, "/t", "REG_EXPAND_SZ", "/d", value, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting %s: %v", name, err)
	}

	NotifyEnvironmentChange()
	return nil
}

// NotifyEnvironmentChange notifies the system about environment changes
func NotifyEnvironmentChange() {
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

// normalizePath normalizes a path for comparison
func normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimRight(path, "/")
	return strings.ToLower(path)
} 