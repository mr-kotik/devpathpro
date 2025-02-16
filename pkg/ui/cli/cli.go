package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"devpathpro/pkg/config"
	"devpathpro/pkg/tools"
)

type CLI struct {
	config *config.Configuration
}

func NewCLI(cfg *config.Configuration) *CLI {
	return &CLI{
		config: cfg,
	}
}

func (c *CLI) Run() {
	fmt.Println("\nDevPathPro - Environment Manager (CLI Mode)")
	fmt.Println("=========================================")

	for {
		fmt.Println("\nMain Menu:")
		fmt.Println("1. Search and Configure Tools")
		fmt.Println("2. Verify Configuration")
		fmt.Println("3. Show Environment Variables")
		fmt.Println("4. Exit")
		fmt.Print("\nSelect an option: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			c.searchAndConfigureTools()
		case 2:
			c.verifyConfiguration()
		case 3:
			c.showEnvironmentVariables()
		case 4:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func (c *CLI) searchAndConfigureTools() {
	fmt.Println("\nSearching for tools...")
	
	for _, prog := range c.config.Programs {
		fmt.Printf("\nChecking %s...\n", prog.Name)
		paths := tools.FindProgram(prog)
		
		if len(paths) == 0 {
			fmt.Printf("❌ %s not found in standard locations\n", prog.Name)
			continue
		}

		fmt.Printf("✅ %s found in:\n", prog.Name)
		for i, path := range paths {
			fmt.Printf("  %d. %s\n", i+1, path)
		}

		// Выбор пути установки
		selectedPath, err := tools.SelectPath(paths, prog.Name)
		if err != nil {
			fmt.Printf("Error selecting path: %v\n", err)
			continue
		}

		// Получение доступных опций конфигурации
		options := tools.GetConfigOptions(prog)
		if len(options) > 0 {
			fmt.Printf("\nConfiguration options for %s:\n", prog.Name)
			fmt.Println("0. All (recommended)")
			for i, opt := range options {
				fmt.Printf("%d. %s - %s\n", i+1, opt.Name, opt.Description)
			}
			
			fmt.Print("\nSelect options (comma-separated numbers, e.g., 1,3): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			var selectedVars []string
			if input != "" && input != "0" {
				numbers := strings.Split(input, ",")
				for _, num := range numbers {
					if idx, err := strconv.Atoi(strings.TrimSpace(num)); err == nil {
						if idx > 0 && idx <= len(options) {
							selectedVars = append(selectedVars, options[idx-1].Variables...)
						}
					}
				}
			}

			if err := tools.ConfigureSelectedPath(prog, selectedPath); err != nil {
				fmt.Printf("❌ Error configuring %s: %v\n", prog.Name, err)
			} else {
				fmt.Printf("✅ %s configured successfully\n", prog.Name)
			}
		} else {
			if err := tools.ConfigureSelectedPath(prog, selectedPath); err != nil {
				fmt.Printf("❌ Error configuring %s: %v\n", prog.Name, err)
			} else {
				fmt.Printf("✅ %s configured successfully\n", prog.Name)
			}
		}
	}
}

func (c *CLI) verifyConfiguration() {
	fmt.Println("\nVerifying configuration...")
	issues := config.VerifyConfigurations()
	
	if len(issues) == 0 {
		fmt.Println("✅ All checks passed successfully!")
		return
	}

	// Группируем проблемы по типу
	issuesByType := make(map[string][]config.ConfigurationIssue)
	for _, issue := range issues {
		issuesByType[issue.Type] = append(issuesByType[issue.Type], issue)
	}

	fmt.Printf("\nFound %d issues:\n", len(issues))

	// Отображаем проблемы по группам
	if pathIssues, ok := issuesByType["PATH"]; ok {
		fmt.Println("\n🔍 PATH Variable Issues:")
		for _, issue := range pathIssues {
			fmt.Printf("  • %s\n", issue.Description)
			fmt.Printf("    Solution: %s\n", issue.Solution)
		}
	}

	if envIssues, ok := issuesByType["ENV"]; ok {
		fmt.Println("\n🔧 Environment Variable Issues:")
		for _, issue := range envIssues {
			fmt.Printf("  • %s\n", issue.Description)
			fmt.Printf("    Solution: %s\n", issue.Solution)
		}
	}

	if programIssues, ok := issuesByType["PROGRAM"]; ok {
		fmt.Println("\n📦 Program Issues:")
		for _, issue := range programIssues {
			fmt.Printf("  • %s\n", issue.Description)
			fmt.Printf("    Solution: %s\n", issue.Solution)
		}
	}

	fmt.Print("\nWould you like to fix these issues? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) == "y" {
		fmt.Println("\nAttempting to fix issues...")
		if err := config.FixConfigurationIssues(issues); err != nil {
			fmt.Printf("❌ Error fixing issues: %v\n", err)
			fmt.Println("Some issues may require manual intervention.")
		} else {
			fmt.Println("✅ Issues fixed successfully!")
			fmt.Println("Please restart the program to apply all changes.")
		}
	}
}

func (c *CLI) showEnvironmentVariables() {
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("Program\t\tVariable\t\tValue\t\tStatus")
	fmt.Println("----------------------------------------------------------------")

	for _, prog := range c.config.Programs {
		paths := tools.FindProgram(prog)
		envVarName := prog.Name + "_HOME"
		envValue := os.Getenv(envVarName)
		
		status := "✅"
		value := envValue
		
		if envValue == "" {
			if len(paths) > 0 {
				value = filepath.Dir(paths[0]) // Используем найденный путь
				status = "⚠️" // Переменная не установлена, но программа найдена
			} else {
				value = "Not found"
				status = "❌" // Программа не найдена
			}
		}

		// Форматируем вывод с табуляцией
		fmt.Printf("%-12s\t%-16s\t%-24s\t%s\n",
			prog.Name,
			envVarName,
			value,
			status,
		)
	}

	fmt.Println("\nStatus Legend:")
	fmt.Println("✅ - Environment variable is set")
	fmt.Println("⚠️ - Program found but environment variable not set")
	fmt.Println("❌ - Program not found")
} 