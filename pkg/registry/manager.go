package registry

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
)

const (
	envKey = `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	userEnvKey = `HKEY_CURRENT_USER\Environment`
)

// IsAdmin checks if the program has administrator privileges by verifying membership
// in the Windows Administrators group using Windows API
func IsAdmin() bool {
	var sid *windows.SID
	// Get SID for Windows Administrators group
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	// Check if current process is a member of Administrators group
	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

// AddToPath adds a new path to both system and user PATH environment variables.
// Requires administrator privileges to modify system PATH.
func AddToPath(newPath string) error {
	if !IsAdmin() {
		return fmt.Errorf("administrator privileges required")
	}

	// Add to system PATH
	if err := addToSystemPath(newPath); err != nil {
		return fmt.Errorf("error adding to system PATH: %v", err)
	}

	// Add to user PATH
	if err := addToUserPath(newPath); err != nil {
		return fmt.Errorf("error adding to user PATH: %v", err)
	}

	NotifyEnvironmentChange()
	return nil
}

// addToSystemPath adds a new path to the system PATH environment variable
func addToSystemPath(newPath string) error {
	return addToPathHelper(envKey, newPath)
}

// addToUserPath adds a new path to the user PATH environment variable
func addToUserPath(newPath string) error {
	return addToPathHelper(userEnvKey, newPath)
}

// addToPathHelper is a helper function to add paths to either system or user PATH.
// It checks if the path already exists and appends it if not.
func addToPathHelper(regKey string, newPath string) error {
	// Get current PATH
	cmd := exec.Command(`C:\Windows\System32\reg.exe`, "query", regKey, "/v", "Path")
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
	cmd = exec.Command(`C:\Windows\System32\reg.exe`, "add", regKey,
		"/v", "Path", "/t", "REG_EXPAND_SZ", "/d", currentPath, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error updating PATH: %v", err)
	}

	return nil
}

// SetEnvironmentVariable sets both system and user environment variables.
// Requires administrator privileges to modify system variables.
func SetEnvironmentVariable(name, value string) error {
	if !IsAdmin() {
		return fmt.Errorf("administrator privileges required")
	}

	// Set system environment variable
	if err := setSystemEnvironmentVariable(name, value); err != nil {
		return fmt.Errorf("error setting system environment variable: %v", err)
	}

	// Set user environment variable
	if err := setUserEnvironmentVariable(name, value); err != nil {
		return fmt.Errorf("error setting user environment variable: %v", err)
	}

	NotifyEnvironmentChange()
	return nil
}

// setSystemEnvironmentVariable sets a system environment variable in the Windows registry
func setSystemEnvironmentVariable(name, value string) error {
	cmd := exec.Command(`C:\Windows\System32\reg.exe`, "add", envKey,
		"/v", name, "/t", "REG_EXPAND_SZ", "/d", value, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting system environment variable: %v", err)
	}
	return nil
}

// setUserEnvironmentVariable sets a user environment variable in the Windows registry
func setUserEnvironmentVariable(name, value string) error {
	cmd := exec.Command(`C:\Windows\System32\reg.exe`, "add", userEnvKey,
		"/v", name, "/t", "REG_EXPAND_SZ", "/d", value, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting user environment variable: %v", err)
	}
	return nil
}

// NotifyEnvironmentChange broadcasts a message to all windows to notify them about
// environment variables changes using Windows API
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

// normalizePath normalizes a path for comparison by converting backslashes to forward slashes,
// removing trailing slashes, and converting to lowercase
func normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimRight(path, "/")
	return strings.ToLower(path)
} 