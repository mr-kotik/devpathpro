package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/mr-kotik/devpathpro/pkg/config"
	"github.com/mr-kotik/devpathpro/pkg/tools"
	"github.com/mr-kotik/devpathpro/pkg/utils"
)

// MainMenu displays and handles the main menu
func MainMenu(cfg *config.Configuration) {
	reader := bufio.NewReader(os.Stdin)

	for {
		utils.ClearScreen()
		fmt.Println("\nDevPathPro - Development Environment Manager")
		utils.PrintDivider("=", 80)
		
		fmt.Println("\nMain Menu:")
		fmt.Println("1. Search and Configure Development Tools")
		fmt.Println("2. Verify Existing Configurations")
		fmt.Println("3. View Current Environment Settings")
		fmt.Println("4. Clear Screen")
		fmt.Println("5. Exit")
		
		fmt.Print("\nSelect an option (1-5): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			SearchToolsMenu(cfg)
		case "2":
			VerifyConfigMenu()
		case "3":
			ViewEnvironmentMenu()
		case "4":
			continue
		case "5":
			utils.ClearScreen()
			fmt.Println("\nThank you for using DevPathPro!")
			fmt.Println("Exiting program...")
			return
		default:
			fmt.Println("\nInvalid option. Press Enter to try again...")
			reader.ReadString('\n')
		}
	}
}

// SearchToolsMenu displays the tool search menu
func SearchToolsMenu(cfg *config.Configuration) {
	utils.ClearScreen()
	fmt.Println("\nAvailable Development Tools:")
	utils.PrintDivider("-", 80)

	// Group programs by category
	categories := make(map[string][]config.Program)
	var categoryOrder []string
	for _, prog := range cfg.Programs {
		if _, exists := categories[prog.Category]; !exists {
			categoryOrder = append(categoryOrder, prog.Category)
		}
		categories[prog.Category] = append(categories[prog.Category], prog)
	}

	// Sort categories
	sort.Strings(categoryOrder)

	// Display programs by category
	var numberedPrograms []config.Program
	currentNumber := 1
	
	for _, category := range categoryOrder {
		programs := categories[category]
		// Sort programs within category
		sort.Slice(programs, func(i, j int) bool {
			return programs[i].Name < programs[j].Name
		})

		fmt.Printf("\n%s:\n", category)
		for _, prog := range programs {
			fmt.Printf("[%3d] %-30s (%s)\n", currentNumber, prog.Name, prog.ExecutableName)
			numberedPrograms = append(numberedPrograms, prog)
			currentNumber++
		}
	}

	utils.PrintDivider("-", 80)
	fmt.Printf("\nTotal tools available: %d\n", len(numberedPrograms))
	fmt.Println("\nOptions:")
	fmt.Println("- Enter numbers separated by comma (e.g.: 1,3,5)")
	fmt.Println("- Enter range using hyphen (e.g.: 1-5)")
	fmt.Println("- Type 'all' to search for all tools")
	fmt.Println("- Type 'category:NAME' to select all tools in a category")
	fmt.Println("- Type 'back' to return to main menu")
	
	fmt.Print("\nYour choice: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.EqualFold(input, "back") {
		return
	}

	var selectedPrograms []config.Program
	if input == "all" {
		selectedPrograms = cfg.Programs
	} else if strings.HasPrefix(strings.ToLower(input), "category:") {
		category := strings.TrimPrefix(strings.ToLower(input), "category:")
		for _, prog := range cfg.Programs {
			if strings.EqualFold(prog.Category, category) {
				selectedPrograms = append(selectedPrograms, prog)
			}
		}
	} else {
		// Handle both individual numbers and ranges
		parts := strings.Split(input, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "-") {
				// Handle range
				rangeParts := strings.Split(part, "-")
				if len(rangeParts) == 2 {
					start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
					end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
					if err1 == nil && err2 == nil && start > 0 && end <= len(numberedPrograms) && start <= end {
						for i := start; i <= end; i++ {
							selectedPrograms = append(selectedPrograms, numberedPrograms[i-1])
						}
					}
				}
			} else {
				// Handle individual number
				if index, err := strconv.Atoi(part); err == nil && index > 0 && index <= len(numberedPrograms) {
					selectedPrograms = append(selectedPrograms, numberedPrograms[index-1])
				}
			}
		}
	}

	if len(selectedPrograms) == 0 {
		fmt.Println("\nNo valid tools selected. Press Enter to try again...")
		reader.ReadString('\n')
		return
	}

	// Process selected programs
	ProcessSelectedTools(selectedPrograms)
}

// VerifyConfigMenu displays the configuration verification menu
func VerifyConfigMenu() {
	fmt.Println("\nVerifying system configuration...")
	issues := config.VerifyConfigurations()

	if len(issues) == 0 {
		fmt.Println("✅ All checks passed successfully! No issues found.")
		return
	}

	fmt.Printf("\nFound %d issues:\n\n", len(issues))

	// Group issues by type
	issuesByType := make(map[string][]config.ConfigurationIssue)
	for _, issue := range issues {
		issuesByType[issue.Type] = append(issuesByType[issue.Type], issue)
	}

	// Display issues by group
	if pathIssues, ok := issuesByType["PATH"]; ok {
		fmt.Println("🔍 PATH Variable Issues:")
		for _, issue := range pathIssues {
			fmt.Printf("  • %s: %s\n", issue.Description, issue.Value)
		}
		fmt.Println()
	}

	if envIssues, ok := issuesByType["ENV"]; ok {
		fmt.Println("🔧 Environment Variable Issues:")
		for _, issue := range envIssues {
			fmt.Printf("  • %s: %s\n", issue.Description, issue.Value)
		}
		fmt.Println()
	}

	if programIssues, ok := issuesByType["PROGRAM"]; ok {
		fmt.Println("📦 Missing Programs:")
		for _, issue := range programIssues {
			fmt.Printf("  • %s (%s)\n", issue.Description, issue.Value)
		}
		fmt.Println()
	}

	// Ask user if they want to fix issues
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Would you like to attempt to fix these issues automatically? (y/n): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		fmt.Println("\nAttempting to fix detected issues...")
		if err := config.FixConfigurationIssues(issues); err != nil {
			fmt.Printf("❌ Error fixing issues: %v\n", err)
		} else {
			fmt.Println("✅ Fixes applied. Please restart the program for changes to take effect.")
		}
	}

	fmt.Println("\nPress Enter to return to main menu...")
	reader.ReadString('\n')
}

// ViewEnvironmentMenu displays current environment settings
func ViewEnvironmentMenu() {
	utils.ClearScreen()
	fmt.Println("\nCurrent Environment Settings:")
	utils.PrintDivider("-", 80)
	
	// Get and sort environment variables
	envVars := os.Environ()
	sort.Strings(envVars)
	
	for _, env := range envVars {
		fmt.Println(env)
	}
	
	fmt.Println("\nPress Enter to return to main menu...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// ProcessSelectedTools processes the selected tools
func ProcessSelectedTools(programs []config.Program) {
	utils.ClearScreen()
	fmt.Printf("\nSelected tools: ")
	for i, prog := range programs {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(prog.Name)
	}
	fmt.Println("\n")

	results := tools.ProcessTools(programs)
	
	// Display results
	configurationChanged := false
	notFoundPrograms := make([]config.Program, 0)

	for _, result := range results {
		fmt.Printf("\n=== %s ===\n", result.Program.Name)
		if !result.Found {
			fmt.Printf("❌ Not found in common locations\n")
			notFoundPrograms = append(notFoundPrograms, result.Program)
			continue
		}

		fmt.Printf("✅ Found in:\n")
		for _, path := range result.Paths {
			fmt.Printf("  - %s\n", path)
		}

		if result.Error != nil {
			fmt.Printf("⚠️ Configuration error: %v\n", result.Error)
		} else {
			fmt.Printf("✅ Successfully configured\n")
			configurationChanged = true
		}
	}

	// If some programs were not found, offer to search all drives
	if len(notFoundPrograms) > 0 {
		fmt.Printf("\n%d tools were not found in common locations.\n", len(notFoundPrograms))
		fmt.Println("Would you like to perform a deep search across all drives? This may take several minutes. (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		
		if answer == "y" || answer == "yes" {
			fmt.Println("\nStarting deep search. This may take a while...")
			deepResults := tools.ProcessToolsDeepSearch(notFoundPrograms)
			
			// Display deep search results
			for _, result := range deepResults {
				fmt.Printf("\n=== %s ===\n", result.Program.Name)
				if !result.Found {
					fmt.Printf("❌ Not found anywhere on this system\n")
					continue
				}

				fmt.Printf("✅ Found in:\n")
				for _, path := range result.Paths {
					fmt.Printf("  - %s\n", path)
				}

				if result.Error != nil {
					fmt.Printf("⚠️ Configuration error: %v\n", result.Error)
				} else {
					fmt.Printf("✅ Successfully configured\n")
					configurationChanged = true
				}
			}
		}
	}

	if configurationChanged {
		fmt.Println("\nConfiguration changes have been made.")
		fmt.Print("Would you like to restart your computer now to apply all changes? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer == "y" {
			fmt.Println("Restarting computer...")
			exec.Command("shutdown", "/r", "/t", "0").Run()
			return
		}
	}

	fmt.Println("\nPress Enter to return to main menu...")
	bufio.NewReader(os.Stdin).ReadString('\n')
} 