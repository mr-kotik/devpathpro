package config

// GetDefaultPrograms returns the default list of supported programs
func GetDefaultPrograms() []Program {
	return []Program{
		// Build Systems
		{
			Name:           "CMake",
			ExecutableName: "cmake.exe",
			CommonPaths: []string{
				`C:\Program Files\CMake\bin`,
				`C:\Program Files (x86)\CMake\bin`,
				`C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
				`C:\Program Files\Microsoft Visual Studio\2022\Professional\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
				`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\Common7\IDE\CommonExtensions\Microsoft\CMake\CMake\bin`,
			},
			Category: "Build Systems",
		},
		{
			Name:           "MSBuild",
			ExecutableName: "msbuild.exe",
			CommonPaths: []string{
				`C:\Program Files\Microsoft Visual Studio\2022\Community\MSBuild\Current\Bin`,
				`C:\Program Files\Microsoft Visual Studio\2022\Professional\MSBuild\Current\Bin`,
				`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\MSBuild\Current\Bin`,
				`C:\Windows\Microsoft.NET\Framework\v4.0.30319`,
				`C:\Windows\Microsoft.NET\Framework64\v4.0.30319`,
			},
			Category: "Build Systems",
		},
		{
			Name:           "Make",
			ExecutableName: "make.exe",
			CommonPaths: []string{
				`C:\Program Files\GnuWin32\bin`,
				`C:\Program Files (x86)\GnuWin32\bin`,
				`C:\MinGW\bin`,
				`C:\msys64\usr\bin`,
				`C:\msys64\mingw64\bin`,
				`C:\cygwin64\bin`,
			},
			Category: "Build Systems",
		},
		{
			Name:           "Ninja",
			ExecutableName: "ninja.exe",
			CommonPaths: []string{
				`C:\Program Files\Ninja`,
				`C:\Program Files (x86)\Ninja`,
			},
			Category: "Build Systems",
		},
		{
			Name:           "Maven",
			ExecutableName: "mvn.cmd",
			CommonPaths: []string{
				`C:\Program Files\Apache\maven\bin`,
				`C:\ProgramData\chocolatey\lib\maven\apache-maven-*\bin`,
			},
			Category: "Build Systems",
		},
		{
			Name:           "Gradle",
			ExecutableName: "gradle.bat",
			CommonPaths: []string{
				`C:\Program Files\Gradle\bin`,
				`C:\ProgramData\chocolatey\lib\gradle\tools\gradle-*\bin`,
			},
			Category: "Build Systems",
		},

		// Development Tools
		{
			Name:           "Git",
			ExecutableName: "git.exe",
			CommonPaths: []string{
				`C:\Program Files\Git\bin`,
				`C:\Program Files (x86)\Git\bin`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "Visual Studio",
			ExecutableName: "devenv.exe",
			CommonPaths: []string{
				`C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\IDE`,
				`C:\Program Files\Microsoft Visual Studio\2022\Professional\Common7\IDE`,
				`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\Common7\IDE`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "VS Code",
			ExecutableName: "code.exe",
			CommonPaths: []string{
				`C:\Program Files\Microsoft VS Code`,
				`C:\Users\%USERNAME%\AppData\Local\Programs\Microsoft VS Code`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "LLVM",
			ExecutableName: "clang.exe",
			CommonPaths: []string{
				`C:\Program Files\LLVM\bin`,
				`C:\Program Files (x86)\LLVM\bin`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "Windows SDK",
			ExecutableName: "rc.exe",
			CommonPaths: []string{
				`C:\Program Files (x86)\Windows Kits\10\bin\*\x64`,
				`C:\Program Files (x86)\Windows Kits\10\bin\*\x86`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "WDK",
			ExecutableName: "devcon.exe",
			CommonPaths: []string{
				`C:\Program Files (x86)\Windows Kits\10\Tools\*\x64`,
				`C:\Program Files (x86)\Windows Kits\10\Tools\*\x86`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "Jenkins",
			ExecutableName: "jenkins.exe",
			CommonPaths: []string{
				`C:\Program Files\Jenkins`,
				`C:\Program Files (x86)\Jenkins`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "SonarQube",
			ExecutableName: "sonar-scanner.bat",
			CommonPaths: []string{
				`C:\Program Files\SonarQube\bin`,
				`C:\Program Files (x86)\SonarQube\bin`,
				`C:\sonar-scanner\bin`,
			},
			Category: "Development Tools",
		},
		{
			Name:           "Grafana",
			ExecutableName: "grafana-server.exe",
			CommonPaths: []string{
				`C:\Program Files\GrafanaLabs\grafana\bin`,
				`C:\Program Files (x86)\GrafanaLabs\grafana\bin`,
			},
			Category: "Development Tools",
		},

		// Programming Languages
		{
			Name:           "Python",
			ExecutableName: "python.exe",
			CommonPaths: []string{
				`C:\Python3*`,
				`C:\Program Files\Python*`,
				`C:\Program Files (x86)\Python*`,
				`C:\Users\%USERNAME%\AppData\Local\Programs\Python\Python*`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Node.js",
			ExecutableName: "node.exe",
			CommonPaths: []string{
				`C:\Program Files\nodejs`,
				`C:\Program Files (x86)\nodejs`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Java",
			ExecutableName: "javac.exe",
			CommonPaths: []string{
				`C:\Program Files\Java\*\bin`,
				`C:\Program Files (x86)\Java\*\bin`,
				`C:\Program Files\Eclipse Foundation\*\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Go",
			ExecutableName: "go.exe",
			CommonPaths: []string{
				`C:\Program Files\Go\bin`,
				`C:\Go\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           ".NET Core",
			ExecutableName: "dotnet.exe",
			CommonPaths: []string{
				`C:\Program Files\dotnet`,
				`C:\Program Files (x86)\dotnet`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Ruby",
			ExecutableName: "ruby.exe",
			CommonPaths: []string{
				`C:\Ruby*\bin`,
				`C:\Program Files\Ruby*\bin`,
				`C:\Program Files (x86)\Ruby*\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Rust",
			ExecutableName: "rustc.exe",
			CommonPaths: []string{
				`C:\Users\%USERNAME%\.cargo\bin`,
				`C:\Program Files\Rust\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Perl",
			ExecutableName: "perl.exe",
			CommonPaths: []string{
				`C:\Perl*\bin`,
				`C:\Strawberry\perl\bin`,
				`C:\Program Files\Perl*\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Scala",
			ExecutableName: "scala.bat",
			CommonPaths: []string{
				`C:\Program Files (x86)\scala\bin`,
				`C:\Program Files\scala\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Kotlin",
			ExecutableName: "kotlin.bat",
			CommonPaths: []string{
				`C:\Program Files\Kotlin\bin`,
				`C:\Program Files (x86)\Kotlin\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Swift",
			ExecutableName: "swift.exe",
			CommonPaths: []string{
				`C:\Program Files\Swift\bin`,
				`C:\Program Files (x86)\Swift\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Haskell",
			ExecutableName: "ghc.exe",
			CommonPaths: []string{
				`C:\Program Files\Haskell\bin`,
				`C:\Program Files\GHC\*\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Erlang",
			ExecutableName: "erl.exe",
			CommonPaths: []string{
				`C:\Program Files\erl*\bin`,
				`C:\Program Files (x86)\erl*\bin`,
			},
			Category: "Programming Languages",
		},
		{
			Name:           "Elixir",
			ExecutableName: "elixir.bat",
			CommonPaths: []string{
				`C:\Program Files\Elixir\bin`,
				`C:\Program Files (x86)\Elixir\bin`,
			},
			Category: "Programming Languages",
		},

		// Package Managers
		{
			Name:           "vcpkg",
			ExecutableName: "vcpkg.exe",
			CommonPaths: []string{
				`C:\vcpkg`,
				`C:\dev\vcpkg`,
				`C:\Program Files\vcpkg`,
			},
			Category: "Package Managers",
		},
		{
			Name:           "Conan",
			ExecutableName: "conan.exe",
			CommonPaths: []string{
				`C:\Program Files\Conan`,
				`C:\Users\%USERNAME%\AppData\Local\Programs\Python\Python3*\Scripts`,
			},
			Category: "Package Managers",
		},

		// Databases
		{
			Name:           "PostgreSQL",
			ExecutableName: "psql.exe",
			CommonPaths: []string{
				`C:\Program Files\PostgreSQL\*\bin`,
				`C:\Program Files (x86)\PostgreSQL\*\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "MySQL",
			ExecutableName: "mysql.exe",
			CommonPaths: []string{
				`C:\Program Files\MySQL\MySQL Server *\bin`,
				`C:\Program Files (x86)\MySQL\MySQL Server *\bin`,
				`C:\Program Files\MariaDB *\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "MongoDB",
			ExecutableName: "mongod.exe",
			CommonPaths: []string{
				`C:\Program Files\MongoDB\Server\*\bin`,
				`C:\Program Files\MongoDB\*\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "Redis",
			ExecutableName: "redis-server.exe",
			CommonPaths: []string{
				`C:\Program Files\Redis`,
				`C:\Program Files (x86)\Redis`,
			},
			Category: "Databases",
		},
		{
			Name:           "Elasticsearch",
			ExecutableName: "elasticsearch.bat",
			CommonPaths: []string{
				`C:\Program Files\Elastic\Elasticsearch\*\bin`,
				`C:\Program Files (x86)\Elastic\Elasticsearch\*\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "SQLite",
			ExecutableName: "sqlite3.exe",
			CommonPaths: []string{
				`C:\Program Files\SQLite`,
				`C:\Program Files (x86)\SQLite`,
			},
			Category: "Databases",
		},
		{
			Name:           "Oracle",
			ExecutableName: "sqlplus.exe",
			CommonPaths: []string{
				`C:\Program Files\Oracle\*\bin`,
				`C:\Program Files (x86)\Oracle\*\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "Cassandra",
			ExecutableName: "cassandra.bat",
			CommonPaths: []string{
				`C:\Program Files\Apache\cassandra\bin`,
				`C:\Program Files (x86)\Apache\cassandra\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "Neo4j",
			ExecutableName: "neo4j.bat",
			CommonPaths: []string{
				`C:\Program Files\Neo4j*\bin`,
				`C:\Program Files (x86)\Neo4j*\bin`,
			},
			Category: "Databases",
		},
		{
			Name:           "InfluxDB",
			ExecutableName: "influxd.exe",
			CommonPaths: []string{
				`C:\Program Files\InfluxData\InfluxDB\bin`,
				`C:\Program Files (x86)\InfluxData\InfluxDB\bin`,
			},
			Category: "Databases",
		},

		// Infrastructure
		{
			Name:           "Docker",
			ExecutableName: "docker.exe",
			CommonPaths: []string{
				`C:\Program Files\Docker\Docker\resources\bin`,
				`C:\Program Files\Docker Toolbox`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Kubernetes",
			ExecutableName: "kubectl.exe",
			CommonPaths: []string{
				`C:\Program Files\Kubernetes\Minikube`,
				`C:\Program Files (x86)\Kubernetes`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Podman",
			ExecutableName: "podman.exe",
			CommonPaths: []string{
				`C:\Program Files\RedHat\Podman`,
				`C:\Program Files (x86)\RedHat\Podman`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Terraform",
			ExecutableName: "terraform.exe",
			CommonPaths: []string{
				`C:\Program Files\Terraform`,
				`C:\Program Files (x86)\Terraform`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Ansible",
			ExecutableName: "ansible.exe",
			CommonPaths: []string{
				`C:\Program Files\Ansible`,
				`C:\Program Files (x86)\Ansible`,
				`C:\Users\%USERNAME%\AppData\Local\Programs\Python\Python*\Scripts`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Helm",
			ExecutableName: "helm.exe",
			CommonPaths: []string{
				`C:\Program Files\Helm`,
				`C:\Program Files (x86)\Helm`,
			},
			Category: "Infrastructure",
		},
		{
			Name:           "Skaffold",
			ExecutableName: "skaffold.exe",
			CommonPaths: []string{
				`C:\Program Files\Skaffold`,
				`C:\Program Files (x86)\Skaffold`,
			},
			Category: "Infrastructure",
		},
	}
} 