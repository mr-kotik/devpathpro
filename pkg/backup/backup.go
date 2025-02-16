package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type EnvironmentBackup struct {
	Timestamp time.Time          `json:"timestamp"`
	Variables map[string]string `json:"variables"`
}

// CreateBackup creates a backup of both registry and environment variables
func CreateBackup() error {
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// Backup registry
	regFile := filepath.Join(backupDir, fmt.Sprintf("registry_%s.reg", timestamp))
	cmd := exec.Command("reg", "export", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment", regFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to backup registry: %v", err)
	}

	// Backup environment variables
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				envVars[env[:i]] = env[i+1:]
				break
			}
		}
	}

	backup := EnvironmentBackup{
		Timestamp: time.Now(),
		Variables: envVars,
	}

	envFile := filepath.Join(backupDir, fmt.Sprintf("env_%s.json", timestamp))
	f, err := os.Create(envFile)
	if err != nil {
		return fmt.Errorf("failed to create environment backup file: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backup); err != nil {
		return fmt.Errorf("failed to write environment backup: %v", err)
	}

	return nil
}

// RestoreBackup restores environment from a specific backup
func RestoreBackup(timestamp string) error {
	backupDir := "backups"
	
	// Restore registry
	regFile := filepath.Join(backupDir, fmt.Sprintf("registry_%s.reg", timestamp))
	cmd := exec.Command("reg", "import", regFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore registry: %v", err)
	}

	// Restore environment variables
	envFile := filepath.Join(backupDir, fmt.Sprintf("env_%s.json", timestamp))
	f, err := os.Open(envFile)
	if err != nil {
		return fmt.Errorf("failed to open environment backup file: %v", err)
	}
	defer f.Close()

	var backup EnvironmentBackup
	if err := json.NewDecoder(f).Decode(&backup); err != nil {
		return fmt.Errorf("failed to read environment backup: %v", err)
	}

	// Apply environment variables
	for key, value := range backup.Variables {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to restore environment variable %s: %v", key, err)
		}
	}

	return nil
}

// ListBackups returns a list of available backups
func ListBackups() ([]string, error) {
	backupDir := "backups"
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var timestamps []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".reg" {
			// Extract timestamp from filename (registry_2006-01-02_15-04-05.reg)
			timestamp := entry.Name()[9 : len(entry.Name())-4]
			timestamps = append(timestamps, timestamp)
		}
	}

	return timestamps, nil
} 