package tools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/mr-kotik/devpathpro/pkg/config"
	"github.com/mr-kotik/devpathpro/pkg/registry"
)

type ProcessResult struct {
	Program config.Program
	Found   bool
	Paths   []string
	Error   error
}

// ConfigOption represents a configuration option
type ConfigOption struct {
	Name        string
	Description string
	Variables   []string
}

// GetConfigOptions returns available configuration options for a program
func GetConfigOptions(prog config.Program) []ConfigOption {
	switch prog.Name {
	case "Python":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Python configuration (HOME and PATH)",
				Variables:   []string{"PYTHON_HOME", "PYTHONPATH"},
			},
			{
				Name:        "Environment",
				Description: "Python environment settings (encoding, buffering, etc.)",
				Variables:   []string{"PYTHONUNBUFFERED", "PYTHONDONTWRITEBYTECODE", "PYTHONIOENCODING", "PYTHONUTF8"},
			},
			{
				Name:        "Development",
				Description: "Development settings (warnings, debug, optimize)",
				Variables:   []string{"PYTHONWARNINGS", "PYTHONDEBUG", "PYTHONOPTIMIZE"},
			},
			{
				Name:        "Pip",
				Description: "Pip package manager settings",
				Variables:   []string{"PIP_CONFIG_FILE", "PIP_DEFAULT_TIMEOUT", "PIP_DISABLE_PIP_VERSION_CHECK"},
			},
		}
	case "Java", "OpenJDK":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Java configuration (HOME and CLASSPATH)",
				Variables:   []string{"JAVA_HOME", "CLASSPATH"},
			},
			{
				Name:        "JVM",
				Description: "JVM memory and performance settings",
				Variables:   []string{"_JAVA_OPTIONS"},
			},
		}
	case "Node.js":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Node.js configuration",
				Variables:   []string{"NODE_PATH"},
			},
			{
				Name:        "NPM",
				Description: "NPM package manager settings",
				Variables:   []string{"NPM_CONFIG_PREFIX", "NPM_CONFIG_CACHE", "NPM_CONFIG_TMP", "NPM_CONFIG_REGISTRY"},
			},
		}
	case "Go":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Go configuration",
				Variables:   []string{"GOROOT", "GOPATH", "GOBIN", "GO111MODULE", "GOCACHE", "GOTMPDIR", "GOPROXY", "GOSUMDB"},
			},
		}
	case "Rust":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Rust configuration",
				Variables:   []string{"RUST_HOME", "CARGO_HOME", "RUSTUP_HOME", "RUST_BACKTRACE", "RUSTC_WRAPPER", "CARGO_TARGET_DIR", "RUSTDOC_THEME"},
			},
		}
	case "Maven":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Maven configuration",
				Variables:   []string{"M2_HOME", "MAVEN_HOME", "MAVEN_OPTS", "MAVEN_CONFIG", "MAVEN_REPOSITORY", "MAVEN_DEBUG_OPTS"},
			},
		}
	case "Gradle":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Gradle configuration",
				Variables:   []string{"GRADLE_HOME", "GRADLE_USER_HOME", "GRADLE_OPTS", "GRADLE_CACHE", "GRADLE_DAEMON", "GRADLE_WORKERS"},
			},
		}
	case "Scala":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Scala configuration",
				Variables:   []string{"SCALA_HOME", "SCALA_OPTS", "SBT_OPTS", "SBT_HOME", "COURSIER_CACHE", "SCALA_CACHE"},
			},
		}
	case "Kotlin":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Kotlin configuration",
				Variables:   []string{"KOTLIN_HOME", "KOTLINC_OPTS", "KOTLIN_COMPILER_OPTS", "KOTLIN_DAEMON_OPTS", "KOTLIN_CACHE_DIR", "KOTLIN_COMPILER_CACHE"},
			},
		}
	case "Erlang":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Erlang configuration",
				Variables:   []string{"ERLANG_HOME", "ERL_LIBS", "ERL_CRASH_DUMP", "ERL_AFLAGS", "ERL_EPMD_PORT", "ERL_MAX_PORTS", "ERL_MAX_ETS_TABLES"},
			},
		}
	case "Elixir":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Elixir configuration",
				Variables:   []string{"ELIXIR_HOME", "MIX_HOME", "HEX_HOME", "MIX_ARCHIVES", "MIX_DEBUG", "MIX_ENV", "ELIXIR_EDITOR", "ELIXIR_ERL_OPTIONS"},
			},
		}
	case "Docker":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Docker configuration",
				Variables:   []string{"DOCKER_HOME", "DOCKER_CONFIG", "DOCKER_CLI_EXPERIMENTAL", "DOCKER_BUILDKIT", "COMPOSE_DOCKER_CLI_BUILD", "DOCKER_HOST"},
			},
		}
	case "Kubernetes":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Kubernetes configuration",
				Variables:   []string{"KUBECONFIG", "KUBE_EDITOR", "HELM_HOME", "HELM_REPOSITORY_CACHE", "HELM_REPOSITORY_CONFIG"},
			},
		}
	case "PostgreSQL":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic PostgreSQL configuration",
				Variables:   []string{"POSTGRES_HOME", "PGDATA", "PGHOST", "PGPORT", "PGLOCALEDIR", "PGLOG", "PGDATABASE", "PGUSER", "PGPASSWORD", "PGTZ", "PGCLIENTENCODING", "PGSSLMODE", "PGCONNECT_TIMEOUT", "PGPOOL_PORT", "PGBOUNCER_PORT", "PGADMIN_PORT"},
			},
		}
	case "MySQL":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic MySQL configuration",
				Variables:   []string{"MYSQL_HOME", "MYSQL_TCP_PORT", "MYSQL_UNIX_PORT", "MYSQL_DATA_DIR", "MYSQL_LOG_DIR", "MYSQL_CONFIG_FILE"},
			},
		}
	case "MongoDB":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic MongoDB configuration",
				Variables:   []string{"MONGODB_HOME", "MONGO_DATA_DIR", "MONGO_LOG_DIR", "MONGO_CONFIG", "MONGO_PORT"},
			},
		}
	case "Redis":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Redis configuration",
				Variables:   []string{"REDIS_HOME", "REDIS_PORT", "REDIS_CONFIG_FILE", "REDIS_DATA_DIR", "REDIS_LOG_FILE"},
			},
		}
	case "Elasticsearch":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Elasticsearch configuration",
				Variables:   []string{"ES_HOME", "ES_PATH_CONF", "ES_PATH_DATA", "ES_PATH_LOGS", "ES_JAVA_OPTS", "ES_PORT", "ES_TRANSPORT_PORT"},
			},
		}
	case "Oracle":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Oracle configuration",
				Variables:   []string{"ORACLE_HOME", "ORACLE_BASE", "ORACLE_SID", "NLS_LANG", "TNS_ADMIN", "ORACLE_TERM"},
			},
		}
	case "Cassandra":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Cassandra configuration",
				Variables:   []string{"CASSANDRA_HOME", "CASSANDRA_CONF", "CASSANDRA_DATA", "CASSANDRA_LOGS", "MAX_HEAP_SIZE", "HEAP_NEWSIZE", "CASSANDRA_PORT", "JMX_PORT"},
			},
		}
	case "Neo4j":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic Neo4j configuration",
				Variables:   []string{"NEO4J_HOME", "NEO4J_CONF", "NEO4J_DATA", "NEO4J_LOGS", "NEO4J_HEAP_MEMORY", "NEO4J_CACHE_MEMORY", "NEO4J_PAGE_CACHE", "NEO4J_HTTP_PORT", "NEO4J_BOLT_PORT", "NEO4J_HTTPS_PORT", "NEO4J_ACCEPT_LICENSE_AGREEMENT", "NEO4J_AUTH", "NEO4J_dbms_memory_pagecache_size", "NEO4J_dbms_memory_heap_initial_size", "NEO4J_dbms_memory_heap_max_size"},
			},
		}
	case "InfluxDB":
		return []ConfigOption{
			{
				Name:        "Basic",
				Description: "Basic InfluxDB configuration",
				Variables:   []string{"INFLUXDB_HOME", "INFLUXDB_CONFIG_PATH", "INFLUXDB_DATA_DIR", "INFLUXDB_META_DIR", "INFLUXDB_WAL_DIR", "INFLUXDB_HTTP_PORT", "INFLUXDB_RPC_PORT", "INFLUXDB_RETENTION", "INFLUXDB_MAX_SERIES_PER_DATABASE", "INFLUXDB_MAX_VALUES_PER_TAG", "INFLUXDB_CACHE_MAX_MEMORY_SIZE", "INFLUXDB_CACHE_SNAPSHOT_MEMORY_SIZE", "INFLUXDB_QUERY_TIMEOUT", "INFLUXDB_HTTP_AUTH_ENABLED", "INFLUXDB_ADMIN_USER", "INFLUXDB_ADMIN_PASSWORD"},
			},
		}
	}
	return nil
}

func showConfigMenu(prog config.Program) []string {
	options := GetConfigOptions(prog)
	if options == nil {
		// If no specific options defined, configure everything
		return nil
	}

	fmt.Printf("\nConfiguration options for %s:\n", prog.Name)
	fmt.Println("0. All (recommended)")
	for i, opt := range options {
		fmt.Printf("%d. %s - %s\n", i+1, opt.Name, opt.Description)
	}
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nSelect options (comma-separated numbers, e.g., 1,3): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" || input == "0" {
		return nil // Configure everything
	}

	var selectedVars []string
	numbers := strings.Split(input, ",")
	for _, num := range numbers {
		if idx, err := strconv.Atoi(strings.TrimSpace(num)); err == nil {
			if idx > 0 && idx <= len(options) {
				selectedVars = append(selectedVars, options[idx-1].Variables...)
			}
		}
	}

	return selectedVars
}

func ProcessTools(programs []config.Program) []ProcessResult {
	results := make([]ProcessResult, len(programs))

	for i, prog := range programs {
		result := ProcessResult{
			Program: prog,
		}

		paths := FindProgram(prog)
		if len(paths) > 0 {
			result.Found = true
			result.Paths = paths

			fmt.Printf("\nFound %s in:\n", prog.Name)
			for _, p := range paths {
				fmt.Printf("  - %s\n", p)
			}

			selectedVars := showConfigMenu(prog)
			if err := configureProgram(prog, paths[0], selectedVars); err != nil {
				result.Error = fmt.Errorf("error configuring %s: %v", prog.Name, err)
			}
		}

		results[i] = result
	}

	return results
}

// ProcessToolsDeepSearch performs a deep search for tools across all drives
func ProcessToolsDeepSearch(programs []config.Program) []ProcessResult {
	results := make([]ProcessResult, len(programs))

	for i, prog := range programs {
		result := ProcessResult{
			Program: prog,
		}

		// Get all available drives
		drives := GetAllDrives()
		var allPaths []string
		resultChan := make(chan string, 100)
		
		fmt.Printf("Searching for %s...\n", prog.Name)
		
		// Search each drive in parallel
		var wg sync.WaitGroup
		for _, drive := range drives {
			wg.Add(1)
			go func(d string) {
				defer wg.Done()
				SearchInDrive(d, prog.ExecutableName, resultChan)
			}(drive)
		}

		// Close result channel when all searches are complete
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Collect results
		for path := range resultChan {
			allPaths = append(allPaths, path)
		}

		if len(allPaths) > 0 {
			result.Found = true
			result.Paths = allPaths

			if err := configureProgram(prog, allPaths[0], nil); err != nil {
				result.Error = fmt.Errorf("error configuring %s: %v", prog.Name, err)
			}
		}

		results[i] = result
	}

	return results
}

func configureProgram(prog config.Program, path string, selectedVars []string) error {
	progDir := filepath.Dir(path)
	if err := registry.AddToPath(progDir); err != nil {
		return fmt.Errorf("error adding to PATH: %v", err)
	}

	switch prog.Name {
	case "Python":
		if err := configurePython(path, selectedVars); err != nil {
			return err
		}
	case "Java", "OpenJDK":
		if err := configureJava(path, selectedVars); err != nil {
			return err
		}
	case "Node.js":
		if err := configureNodejs(path, selectedVars); err != nil {
			return err
		}
	case "Go":
		if err := configureGo(path, selectedVars); err != nil {
			return err
		}
	case "Rust":
		if err := configureRust(path); err != nil {
			return err
		}
	case "Maven":
		if err := configureMaven(path); err != nil {
			return err
		}
	case "Gradle":
		if err := configureGradle(path); err != nil {
			return err
		}
	case "Scala":
		if err := configureScala(path); err != nil {
			return err
		}
	case "Kotlin":
		if err := configureKotlin(path); err != nil {
			return err
		}
	case "Erlang":
		if err := configureErlang(path); err != nil {
			return err
		}
	case "Elixir":
		if err := configureElixir(path); err != nil {
			return err
		}
	case "Docker":
		if err := configureDocker(path, selectedVars); err != nil {
			return err
		}
	case "Kubernetes":
		if err := configureKubernetes(path, selectedVars); err != nil {
			return err
		}
	case "PostgreSQL":
		if err := configurePostgreSQL(path); err != nil {
			return err
		}
	case "MySQL":
		if err := configureMySQL(path); err != nil {
			return err
		}
	case "MongoDB":
		if err := configureMongoDB(path); err != nil {
			return err
		}
	case "Redis":
		if err := configureRedis(path); err != nil {
			return err
		}
	case "Elasticsearch":
		if err := configureElasticsearch(path); err != nil {
			return err
		}
	case "Oracle":
		if err := configureOracle(path); err != nil {
			return err
		}
	case "Cassandra":
		if err := configureCassandra(path); err != nil {
			return err
		}
	case "Neo4j":
		if err := configureNeo4j(path); err != nil {
			return err
		}
	case "InfluxDB":
		if err := configureInfluxDB(path); err != nil {
			return err
		}
	}

	return nil
}

func configurePython(path string, selectedVars []string) error {
	pythonDir := filepath.Dir(path)
	scriptsDir := filepath.Join(pythonDir, "Scripts")

	// Define all possible configurations
	pythonConfig := map[string]string{
		"PYTHON_HOME": pythonDir,
		"PYTHONPATH": fmt.Sprintf("%s;%s", pythonDir, filepath.Join(pythonDir, "Lib", "site-packages")),
		"PYTHONUNBUFFERED": "1",
		"PYTHONDONTWRITEBYTECODE": "1",
		"PYTHONIOENCODING": "utf-8",
		"PYTHONUTF8": "1",
		"PYTHONWARNINGS": "default",
		"PYTHONDEBUG": "1",
		"PYTHONOPTIMIZE": "1",
		"PIP_CONFIG_FILE": filepath.Join(os.Getenv("APPDATA"), "pip", "pip.ini"),
		"PIP_DEFAULT_TIMEOUT": "100",
		"PIP_DISABLE_PIP_VERSION_CHECK": "1",
		"VIRTUAL_ENV_DISABLE_PROMPT": "1",
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range pythonConfig {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := pythonConfig[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	// Always add Scripts to PATH
	if err := registry.AddToPath(scriptsDir); err != nil {
		return fmt.Errorf("error adding Python Scripts to PATH: %v", err)
	}

	return nil
}

func configureJava(path string, selectedVars []string) error {
	javaDir := filepath.Dir(filepath.Dir(path))
	
	// Set JAVA_HOME
	if err := os.Setenv("JAVA_HOME", javaDir); err != nil {
		return fmt.Errorf("error setting JAVA_HOME: %v", err)
	}

	// Добавляем новые настройки
	// JDK инструменты
	toolsJar := filepath.Join(javaDir, "lib", "tools.jar")
	dtJar := filepath.Join(javaDir, "lib", "dt.jar")
	
	classpath := []string{toolsJar, dtJar}
	if existing := os.Getenv("CLASSPATH"); existing != "" {
		classpath = append(classpath, existing)
	}
	
	if err := os.Setenv("CLASSPATH", strings.Join(classpath, ";")); err != nil {
		return fmt.Errorf("error setting CLASSPATH: %v", err)
	}

	// Настройки JVM
	if err := os.Setenv("_JAVA_OPTIONS", "-Xmx2048m -Xms512m"); err != nil {
		return fmt.Errorf("error setting _JAVA_OPTIONS: %v", err)
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range map[string]string{
			"JAVA_HOME": javaDir,
			"CLASSPATH": strings.Join(classpath, ";"),
			"_JAVA_OPTIONS": "-Xmx2048m -Xms512m",
		} {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := map[string]string{
				"JAVA_HOME": javaDir,
				"CLASSPATH": strings.Join(classpath, ";"),
				"_JAVA_OPTIONS": "-Xmx2048m -Xms512m",
			}[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func configureNodejs(path string, selectedVars []string) error {
	nodeDir := filepath.Dir(path)
	
	// Define all possible configurations
	nodeConfig := map[string]string{
		"NODE_PATH": filepath.Join(nodeDir, "node_modules"),
		"NPM_CONFIG_PREFIX": filepath.Join(os.Getenv("APPDATA"), "npm"),
		"NPM_CONFIG_CACHE": filepath.Join(os.Getenv("APPDATA"), "npm-cache"),
		"NPM_CONFIG_TMP": filepath.Join(os.Getenv("TEMP"), "npm"),
		"NPM_CONFIG_REGISTRY": "https://registry.npmjs.org/",
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range nodeConfig {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := nodeConfig[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func configureGo(path string, selectedVars []string) error {
	goDir := filepath.Dir(filepath.Dir(path))
	
	// Define all possible configurations
	goConfig := map[string]string{
		"GOROOT": goDir,
		"GOPATH": filepath.Join(os.Getenv("USERPROFILE"), "go"),
		"GOBIN": filepath.Join(os.Getenv("USERPROFILE"), "go", "bin"),
		"GO111MODULE": "on",
		"GOCACHE": filepath.Join(os.Getenv("USERPROFILE"), "go", "cache"),
		"GOTMPDIR": filepath.Join(os.Getenv("TEMP"), "go-build"),
		"GOPROXY": "https://proxy.golang.org,direct",
		"GOSUMDB": "sum.golang.org",
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range goConfig {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := goConfig[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func configureRust(path string) error {
	rustDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Rust
	rustConfig := map[string]string{
		"RUST_HOME": rustDir,
		"CARGO_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".cargo"),
		"RUSTUP_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".rustup"),
		"RUST_BACKTRACE": "1",
		"RUSTC_WRAPPER": "sccache", // Если установлен sccache
		"CARGO_TARGET_DIR": filepath.Join(os.Getenv("USERPROFILE"), ".cargo", "target"),
		"RUSTDOC_THEME": "dark",
	}

	for key, value := range rustConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureMaven(path string) error {
	mavenDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Maven
	mavenConfig := map[string]string{
		"M2_HOME": mavenDir,
		"MAVEN_HOME": mavenDir,
		"MAVEN_OPTS": "-Xmx2048m -Xms1024m",
		"MAVEN_CONFIG": filepath.Join(os.Getenv("USERPROFILE"), ".m2"),
		"MAVEN_REPOSITORY": filepath.Join(os.Getenv("USERPROFILE"), ".m2", "repository"),
		"MAVEN_DEBUG_OPTS": "-Xdebug -Xnoagent -Djava.compiler=NONE -Xrunjdwp:transport=dt_socket,server=y,suspend=n,address=8000",
	}

	for key, value := range mavenConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureGradle(path string) error {
	gradleDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Gradle
	gradleConfig := map[string]string{
		"GRADLE_HOME": gradleDir,
		"GRADLE_USER_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".gradle"),
		"GRADLE_OPTS": "-Xmx2048m -Xms512m -XX:MaxPermSize=512m -XX:+HeapDumpOnOutOfMemoryError",
		"GRADLE_CACHE": filepath.Join(os.Getenv("USERPROFILE"), ".gradle", "caches"),
		"GRADLE_DAEMON": "true",
		"GRADLE_WORKERS": "4",
	}

	for key, value := range gradleConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureScala(path string) error {
	scalaDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Scala
	scalaConfig := map[string]string{
		"SCALA_HOME": scalaDir,
		"SCALA_OPTS": "-Xmx2048m -Xms1024m",
		"SBT_OPTS": "-Xmx2G -XX:+UseConcMarkSweepGC -XX:+CMSClassUnloadingEnabled",
		"SBT_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".sbt"),
		"COURSIER_CACHE": filepath.Join(os.Getenv("USERPROFILE"), ".coursier", "cache"),
		"SCALA_CACHE": filepath.Join(os.Getenv("USERPROFILE"), ".scala", "cache"),
	}

	for key, value := range scalaConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureKotlin(path string) error {
	kotlinDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Kotlin
	kotlinConfig := map[string]string{
		"KOTLIN_HOME": kotlinDir,
		"KOTLINC_OPTS": "-Xmx2G -Xms512M",
		"KOTLIN_COMPILER_OPTS": "-Xjvm-default=enable -Xopt-in=kotlin.RequiresOptIn",
		"KOTLIN_DAEMON_OPTS": "-Xmx2G -Xms512M",
		"KOTLIN_CACHE_DIR": filepath.Join(os.Getenv("USERPROFILE"), ".kotlin", "cache"),
		"KOTLIN_COMPILER_CACHE": filepath.Join(os.Getenv("USERPROFILE"), ".kotlin", "daemon"),
	}

	for key, value := range kotlinConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureErlang(path string) error {
	erlangDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Erlang
	erlangConfig := map[string]string{
		"ERLANG_HOME": erlangDir,
		"ERL_LIBS": filepath.Join(erlangDir, "lib"),
		"ERL_CRASH_DUMP": filepath.Join(os.Getenv("USERPROFILE"), ".erlang_crash.dump"),
		"ERL_AFLAGS": "-kernel shell_history enabled",
		"ERL_EPMD_PORT": "4369",
		"ERL_MAX_PORTS": "32768",
		"ERL_MAX_ETS_TABLES": "32768",
	}

	for key, value := range erlangConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureElixir(path string) error {
	elixirDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Elixir
	elixirConfig := map[string]string{
		"ELIXIR_HOME": elixirDir,
		"MIX_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".mix"),
		"HEX_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".hex"),
		"MIX_ARCHIVES": filepath.Join(os.Getenv("USERPROFILE"), ".mix", "archives"),
		"MIX_DEBUG": "1",
		"MIX_ENV": "dev",
		"ELIXIR_EDITOR": "code --wait",
		"ELIXIR_ERL_OPTIONS": "-kernel shell_history enabled",
	}

	for key, value := range elixirConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureDocker(path string, selectedVars []string) error {
	dockerDir := filepath.Dir(filepath.Dir(path))
	
	// Define all possible configurations
	dockerConfig := map[string]string{
		"DOCKER_HOME": dockerDir,
		"DOCKER_CONFIG": filepath.Join(os.Getenv("USERPROFILE"), ".docker"),
		"DOCKER_CLI_EXPERIMENTAL": "enabled",
		"DOCKER_BUILDKIT": "1",
		"COMPOSE_DOCKER_CLI_BUILD": "1",
		"DOCKER_HOST": "tcp://localhost:2375",
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range dockerConfig {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := dockerConfig[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func configureKubernetes(path string, selectedVars []string) error {
	// Define all possible configurations
	k8sConfig := map[string]string{
		"KUBECONFIG": filepath.Join(os.Getenv("USERPROFILE"), ".kube", "config"),
		"KUBE_EDITOR": "code --wait",
		"HELM_HOME": filepath.Join(os.Getenv("USERPROFILE"), ".helm"),
		"HELM_REPOSITORY_CACHE": filepath.Join(os.Getenv("USERPROFILE"), ".helm", "repository", "cache"),
		"HELM_REPOSITORY_CONFIG": filepath.Join(os.Getenv("USERPROFILE"), ".helm", "repository", "repositories.yaml"),
	}

	// If no specific variables selected, configure all
	if len(selectedVars) == 0 {
		for key, value := range k8sConfig {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("error setting %s: %v", key, err)
			}
		}
	} else {
		// Configure only selected variables
		for _, key := range selectedVars {
			if value, exists := k8sConfig[key]; exists {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func configurePostgreSQL(path string) error {
	pgDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки PostgreSQL
	pgConfig := map[string]string{
		"POSTGRES_HOME": pgDir,
		"PGDATA": filepath.Join(pgDir, "data"),
		"PGHOST": "localhost",
		"PGPORT": "5432",
		"PGLOCALEDIR": filepath.Join(pgDir, "share", "locale"),
		"PGLOG": filepath.Join(pgDir, "log", "postgresql.log"),
		"PGDATABASE": "postgres",
		"PGUSER": "postgres",
		"PGPASSWORD": "postgres",
		"PGTZ": "UTC",
		"PGCLIENTENCODING": "UTF8",
		"PGSSLMODE": "prefer",
		"PGCONNECT_TIMEOUT": "10",
		"PGPOOL_PORT": "9999",
		"PGBOUNCER_PORT": "6432",
		"PGADMIN_PORT": "5050",
	}

	for key, value := range pgConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureMySQL(path string) error {
	mysqlDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки
	mysqlConfig := map[string]string{
		"MYSQL_HOME": mysqlDir,
		"MYSQL_TCP_PORT": "3306",
		"MYSQL_UNIX_PORT": "3306",
		"MYSQL_DATA_DIR": filepath.Join(mysqlDir, "data"),
		"MYSQL_LOG_DIR": filepath.Join(mysqlDir, "log"),
		"MYSQL_CONFIG_FILE": filepath.Join(mysqlDir, "my.ini"),
	}

	for key, value := range mysqlConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureMongoDB(path string) error {
	mongoDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки
	mongoConfig := map[string]string{
		"MONGODB_HOME": mongoDir,
		"MONGO_DATA_DIR": filepath.Join(mongoDir, "data", "db"),
		"MONGO_LOG_DIR": filepath.Join(mongoDir, "log"),
		"MONGO_CONFIG": filepath.Join(mongoDir, "mongod.cfg"),
		"MONGO_PORT": "27017",
	}

	for key, value := range mongoConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureRedis(path string) error {
	redisDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки
	redisConfig := map[string]string{
		"REDIS_HOME": redisDir,
		"REDIS_PORT": "6379",
		"REDIS_CONFIG_FILE": filepath.Join(redisDir, "redis.windows.conf"),
		"REDIS_DATA_DIR": filepath.Join(redisDir, "data"),
		"REDIS_LOG_FILE": filepath.Join(redisDir, "log", "redis.log"),
	}

	for key, value := range redisConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureElasticsearch(path string) error {
	esDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки
	esConfig := map[string]string{
		"ES_HOME": esDir,
		"ES_PATH_CONF": filepath.Join(esDir, "config"),
		"ES_PATH_DATA": filepath.Join(esDir, "data"),
		"ES_PATH_LOGS": filepath.Join(esDir, "logs"),
		"ES_JAVA_OPTS": "-Xms1g -Xmx1g",
		"ES_PORT": "9200",
		"ES_TRANSPORT_PORT": "9300",
	}

	for key, value := range esConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureOracle(path string) error {
	oracleDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Oracle
	oracleConfig := map[string]string{
		"ORACLE_HOME": oracleDir,
		"ORACLE_BASE": filepath.Dir(oracleDir),
		"ORACLE_SID": "ORCL",
		"NLS_LANG": "AMERICAN_AMERICA.AL32UTF8",
		"TNS_ADMIN": filepath.Join(oracleDir, "network", "admin"),
		"ORACLE_TERM": "xterm",
	}

	for key, value := range oracleConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureCassandra(path string) error {
	cassandraDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Cassandra
	cassandraConfig := map[string]string{
		"CASSANDRA_HOME": cassandraDir,
		"CASSANDRA_CONF": filepath.Join(cassandraDir, "conf"),
		"CASSANDRA_DATA": filepath.Join(cassandraDir, "data"),
		"CASSANDRA_LOGS": filepath.Join(cassandraDir, "logs"),
		"MAX_HEAP_SIZE": "1G",
		"HEAP_NEWSIZE": "250M",
		"CASSANDRA_PORT": "9042",
		"JMX_PORT": "7199",
	}

	for key, value := range cassandraConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureNeo4j(path string) error {
	neo4jDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки Neo4j
	neo4jConfig := map[string]string{
		"NEO4J_HOME": neo4jDir,
		"NEO4J_CONF": filepath.Join(neo4jDir, "conf"),
		"NEO4J_DATA": filepath.Join(neo4jDir, "data"),
		"NEO4J_LOGS": filepath.Join(neo4jDir, "logs"),
		"NEO4J_HEAP_MEMORY": "4G",
		"NEO4J_CACHE_MEMORY": "2G",
		"NEO4J_PAGE_CACHE": "2G",
		"NEO4J_HTTP_PORT": "7474",
		"NEO4J_BOLT_PORT": "7687",
		"NEO4J_HTTPS_PORT": "7473",
		"NEO4J_ACCEPT_LICENSE_AGREEMENT": "yes",
		"NEO4J_AUTH": "neo4j/neo4j",
		"NEO4J_dbms_memory_pagecache_size": "2G",
		"NEO4J_dbms_memory_heap_initial_size": "2G",
		"NEO4J_dbms_memory_heap_max_size": "4G",
	}

	for key, value := range neo4jConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
}

func configureInfluxDB(path string) error {
	influxDir := filepath.Dir(filepath.Dir(path))
	
	// Основные настройки InfluxDB
	influxConfig := map[string]string{
		"INFLUXDB_HOME": influxDir,
		"INFLUXDB_CONFIG_PATH": filepath.Join(influxDir, "influxdb.conf"),
		"INFLUXDB_DATA_DIR": filepath.Join(influxDir, "data"),
		"INFLUXDB_META_DIR": filepath.Join(influxDir, "meta"),
		"INFLUXDB_WAL_DIR": filepath.Join(influxDir, "wal"),
		"INFLUXDB_HTTP_PORT": "8086",
		"INFLUXDB_RPC_PORT": "8088",
		"INFLUXDB_RETENTION": "52w",
		"INFLUXDB_MAX_SERIES_PER_DATABASE": "1000000",
		"INFLUXDB_MAX_VALUES_PER_TAG": "100000",
		"INFLUXDB_CACHE_MAX_MEMORY_SIZE": "1g",
		"INFLUXDB_CACHE_SNAPSHOT_MEMORY_SIZE": "256m",
		"INFLUXDB_QUERY_TIMEOUT": "60s",
		"INFLUXDB_HTTP_AUTH_ENABLED": "true",
		"INFLUXDB_ADMIN_USER": "admin",
		"INFLUXDB_ADMIN_PASSWORD": "admin",
	}

	for key, value := range influxConfig {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s: %v", key, err)
		}
	}

	return nil
} 