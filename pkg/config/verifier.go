package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigurationIssue represents a configuration problem
type ConfigurationIssue struct {
	Type        string // PATH, ENV, PROGRAM, PERMISSION, SECURITY
	Severity    string // HIGH, MEDIUM, LOW
	Description string
	Value       string
	Solution    string
}

// VerifyConfigurations checks the system's PATH and environment variables for issues
func VerifyConfigurations() []ConfigurationIssue {
	var issues []ConfigurationIssue

	// Check PATH
	pathIssues := verifyPath()
	issues = append(issues, pathIssues...)

	// Check environment variables
	envIssues := verifyEnvironmentVariables()
	issues = append(issues, envIssues...)

	// Check program paths and versions
	programIssues := verifyProgramPaths()
	issues = append(issues, programIssues...)

	// Check security settings
	securityIssues := verifySecuritySettings()
	issues = append(issues, securityIssues...)

	// Check permissions
	permissionIssues := verifyPermissions()
	issues = append(issues, permissionIssues...)

	return issues
}

func verifyPath() []ConfigurationIssue {
	var issues []ConfigurationIssue
	path := os.Getenv("PATH")
	paths := strings.Split(path, ";")

	// Check for duplicate paths
	seen := make(map[string]bool)
	for _, p := range paths {
		if p == "" {
			continue
		}
		
		normalized := strings.ToLower(filepath.Clean(p))
		
		// Check for duplicate paths
		if seen[normalized] {
			issues = append(issues, ConfigurationIssue{
				Type:        "PATH",
				Severity:    "LOW",
				Description: "Duplicate PATH entry found",
				Value:       p,
				Solution:    "Remove duplicate entry from PATH",
			})
		}
		seen[normalized] = true

		// Check if path exists
		if _, err := os.Stat(p); os.IsNotExist(err) {
			issues = append(issues, ConfigurationIssue{
				Type:        "PATH",
				Severity:    "MEDIUM",
				Description: "PATH entry does not exist",
				Value:       p,
				Solution:    "Remove non-existent path or create directory",
			})
		}

		// Check path length
		if len(p) > 260 {
			issues = append(issues, ConfigurationIssue{
				Type:        "PATH",
				Severity:    "HIGH",
				Description: "PATH entry exceeds Windows path length limit",
				Value:       p,
				Solution:    "Shorten path or use subst to create drive letter mapping",
			})
		}
	}

	return issues
}

func verifyEnvironmentVariables() []ConfigurationIssue {
	var issues []ConfigurationIssue
	
	// Расширенный список переменных окружения для проверки
	varsToCheck := map[string]struct {
		desc     string
		required bool
		verify   func(string) bool
	}{
		"JAVA_HOME": {"Java Development Kit", true, nil},
		"PYTHON_HOME": {"Python", true, nil},
		"GOROOT": {"Go Programming Language", true, nil},
		"GOPATH": {"Go Workspace", true, nil},
		"NODE_PATH": {"Node.js modules", false, nil},
		"MAVEN_HOME": {"Apache Maven", false, nil},
		"GRADLE_HOME": {"Gradle", false, nil},
		"DOCKER_HOME": {"Docker", false, nil},
		"KUBECONFIG": {"Kubernetes", false, nil},
		"RUST_HOME": {"Rust", false, nil},
		"CARGO_HOME": {"Cargo (Rust package manager)", false, nil},
		"POSTGRES_HOME": {"PostgreSQL", false, nil},
		"MYSQL_HOME": {"MySQL", false, nil},
		"MONGODB_HOME": {"MongoDB", false, nil},
		"REDIS_HOME": {"Redis", false, nil},
		"ES_HOME": {"Elasticsearch", false, nil},
		"NEO4J_HOME": {"Neo4j", false, nil},
		"INFLUXDB_HOME": {"InfluxDB", false, nil},
	}

	for env, info := range varsToCheck {
		value := os.Getenv(env)
		
		// Проверка наличия обязательных переменных
		if info.required && value == "" {
			issues = append(issues, ConfigurationIssue{
				Type:        "ENV",
				Severity:    "HIGH",
				Description: fmt.Sprintf("Missing required %s environment variable", info.desc),
				Value:       env,
				Solution:    fmt.Sprintf("Set %s environment variable to the installation directory", env),
			})
			continue
		}

		// Если переменная установлена, проверяем путь
		if value != "" {
			// Проверка существования пути
			if _, err := os.Stat(value); os.IsNotExist(err) {
				issues = append(issues, ConfigurationIssue{
					Type:        "ENV",
					Severity:    "MEDIUM",
					Description: fmt.Sprintf("%s path does not exist", info.desc),
					Value:       fmt.Sprintf("%s=%s", env, value),
					Solution:    "Update path to correct installation directory",
				})
			}

			// Проверка дополнительных условий
			if info.verify != nil && !info.verify(value) {
				issues = append(issues, ConfigurationIssue{
					Type:        "ENV",
					Severity:    "LOW",
					Description: fmt.Sprintf("%s value validation failed", info.desc),
					Value:       fmt.Sprintf("%s=%s", env, value),
					Solution:    "Check value format and requirements",
				})
			}
		}
	}

	return issues
}

func verifyProgramPaths() []ConfigurationIssue {
	var issues []ConfigurationIssue
	programs := GetDefaultPrograms()

	for _, prog := range programs {
		found := false
		var foundPath string

		for _, path := range prog.CommonPaths {
			// Replace environment variables
			path = os.ExpandEnv(path)
			
			// Handle wildcard paths
			if strings.Contains(path, "*") {
				matches, err := filepath.Glob(path)
				if err == nil && len(matches) > 0 {
					for _, match := range matches {
						execPath := filepath.Join(match, prog.ExecutableName)
						if _, err := os.Stat(execPath); err == nil {
							found = true
							foundPath = execPath
							break
						}
					}
				}
			} else {
				execPath := filepath.Join(path, prog.ExecutableName)
				if _, err := os.Stat(execPath); err == nil {
					found = true
					foundPath = execPath
					break
				}
			}
		}

		if !found {
			issues = append(issues, ConfigurationIssue{
				Type:        "PROGRAM",
				Severity:    "MEDIUM",
				Description: fmt.Sprintf("%s not found in common installation paths", prog.Name),
				Value:       prog.ExecutableName,
				Solution:    fmt.Sprintf("Install %s or update PATH if already installed", prog.Name),
			})
		} else {
			// Проверка прав доступа к исполняемому файлу
			if info, err := os.Stat(foundPath); err == nil {
				if info.Mode().Perm()&0111 == 0 {
					issues = append(issues, ConfigurationIssue{
						Type:        "PERMISSION",
						Severity:    "HIGH",
						Description: fmt.Sprintf("%s executable permissions are incorrect", prog.Name),
						Value:       foundPath,
						Solution:    "Update file permissions to allow execution",
					})
				}
			}
		}
	}

	return issues
}

func verifySecuritySettings() []ConfigurationIssue {
	var issues []ConfigurationIssue

	// Проверка настроек безопасности для баз данных
	dbSecurityChecks := map[string]struct {
		envVars []string
		desc    string
	}{
		"PostgreSQL": {
			envVars: []string{"PGPASSWORD", "PGUSER"},
			desc:    "PostgreSQL credentials in environment",
		},
		"MySQL": {
			envVars: []string{"MYSQL_ROOT_PASSWORD", "MYSQL_USER"},
			desc:    "MySQL credentials in environment",
		},
		"MongoDB": {
			envVars: []string{"MONGO_INITDB_ROOT_PASSWORD", "MONGO_INITDB_ROOT_USERNAME"},
			desc:    "MongoDB credentials in environment",
		},
	}

	for dbName, check := range dbSecurityChecks {
		for _, envVar := range check.envVars {
			if value := os.Getenv(envVar); value != "" {
				issues = append(issues, ConfigurationIssue{
					Type:        "SECURITY",
					Severity:    "HIGH",
					Description: fmt.Sprintf("%s: %s", dbName, check.desc),
					Value:       envVar,
					Solution:    "Use configuration files instead of environment variables for credentials",
				})
			}
		}
	}

	return issues
}

func verifyPermissions() []ConfigurationIssue {
	var issues []ConfigurationIssue

	// Проверка прав доступа к важным директориям
	dirsToCheck := []struct {
		path        string
		description string
		required    bool
	}{
		{os.Getenv("GOPATH"), "Go workspace", true},
		{os.Getenv("MAVEN_REPOSITORY"), "Maven repository", false},
		{os.Getenv("GRADLE_USER_HOME"), "Gradle home", false},
		{os.Getenv("DOCKER_CONFIG"), "Docker configuration", true},
		{filepath.Join(os.Getenv("USERPROFILE"), ".kube"), "Kubernetes configuration", false},
	}

	for _, dir := range dirsToCheck {
		if dir.path == "" {
			if dir.required {
				issues = append(issues, ConfigurationIssue{
					Type:        "PERMISSION",
					Severity:    "HIGH",
					Description: fmt.Sprintf("Required directory path not set: %s", dir.description),
					Value:       dir.path,
					Solution:    "Set correct path and ensure proper permissions",
				})
			}
			continue
		}

		if info, err := os.Stat(dir.path); err == nil {
			// Проверка прав на запись
			if info.Mode().Perm()&0200 == 0 {
				issues = append(issues, ConfigurationIssue{
					Type:        "PERMISSION",
					Severity:    "HIGH",
					Description: fmt.Sprintf("No write permission: %s", dir.description),
					Value:       dir.path,
					Solution:    "Grant write permissions to the current user",
				})
			}
		}
	}

	return issues
}

// FixConfigurationIssues attempts to fix identified configuration issues
func FixConfigurationIssues(issues []ConfigurationIssue) error {
	for _, issue := range issues {
		switch issue.Type {
		case "PATH":
			if strings.Contains(issue.Description, "Duplicate") || strings.Contains(issue.Description, "does not exist") {
				newPath := removePath(os.Getenv("PATH"), issue.Value)
				if err := os.Setenv("PATH", newPath); err != nil {
					return fmt.Errorf("failed to update PATH: %v", err)
				}
			}

		case "ENV":
			if strings.Contains(issue.Description, "Missing") {
				if value := findProgramPath(issue.Value); value != "" {
					if err := os.Setenv(issue.Value, value); err != nil {
						return fmt.Errorf("failed to set %s: %v", issue.Value, err)
					}
				}
			}

		case "PERMISSION":
			if strings.Contains(issue.Description, "No write permission") {
				if err := os.Chmod(issue.Value, 0755); err != nil {
					return fmt.Errorf("failed to update permissions for %s: %v", issue.Value, err)
				}
			}

		case "SECURITY":
			// Для проблем безопасности только выводим предупреждение
			fmt.Printf("Security warning: %s\nRecommended solution: %s\n", issue.Description, issue.Solution)
		}
	}
	return nil
}

func removePath(path, valueToRemove string) string {
	paths := strings.Split(path, ";")
	var newPaths []string
	for _, p := range paths {
		if strings.ToLower(filepath.Clean(p)) != strings.ToLower(filepath.Clean(valueToRemove)) {
			newPaths = append(newPaths, p)
		}
	}
	return strings.Join(newPaths, ";")
}

func findProgramPath(envVar string) string {
	commonPaths := map[string][]string{
		"JAVA_HOME": {
			`C:\Program Files\Java\*`,
			`C:\Program Files (x86)\Java\*`,
			`C:\Program Files\Eclipse Foundation\*`,
		},
		"PYTHON_HOME": {
			`C:\Python3*`,
			`C:\Program Files\Python*`,
			`C:\Program Files (x86)\Python*`,
			`C:\Users\%USERNAME%\AppData\Local\Programs\Python\Python*`,
		},
		"GOROOT": {
			`C:\Go`,
			`C:\Program Files\Go`,
		},
		"NODE_PATH": {
			`C:\Program Files\nodejs`,
			`C:\Program Files (x86)\nodejs`,
		},
		"DOCKER_HOME": {
			`C:\Program Files\Docker`,
			`C:\Program Files\Docker\Docker`,
		},
		"RUST_HOME": {
			`C:\Users\%USERNAME%\.cargo`,
			`C:\Program Files\Rust`,
		},
	}

	if paths, ok := commonPaths[envVar]; ok {
		for _, pathPattern := range paths {
			// Expand environment variables
			pathPattern = os.ExpandEnv(pathPattern)
			
			// Handle wildcards
			if strings.Contains(pathPattern, "*") {
				if matches, err := filepath.Glob(pathPattern); err == nil {
					for _, match := range matches {
						if _, err := os.Stat(match); err == nil {
							return match
						}
					}
				}
			} else {
				if _, err := os.Stat(pathPattern); err == nil {
					return pathPattern
				}
			}
		}
	}
	return ""
} 