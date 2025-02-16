package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"devpathpro/pkg/config"
	"devpathpro/pkg/tools"
)

// DevPathProGUI represents the main GUI application
type DevPathProGUI struct {
	window       fyne.Window
	config       *config.Configuration
	tabContainer *container.AppTabs
}

// NewDevPathProGUI creates a new instance of the GUI application
func NewDevPathProGUI() *DevPathProGUI {
	a := app.New()
	window := a.NewWindow("DevPathPro")

	gui := &DevPathProGUI{
		window: window,
		config: &config.Configuration{
			Programs: config.GetDefaultPrograms(),
		},
	}

	gui.setupUI()
	return gui
}

// Run starts the GUI application
func (gui *DevPathProGUI) Run() {
	gui.window.Resize(fyne.NewSize(800, 600))
	gui.window.ShowAndRun()
}

// setupUI initializes the user interface
func (gui *DevPathProGUI) setupUI() {
	gui.setupTabs()
	gui.window.SetContent(gui.tabContainer)
}

// setupTabs creates and configures all tabs
func (gui *DevPathProGUI) setupTabs() {
	toolsTab := gui.createToolsTab()
	verifyTab := gui.createVerifyTab()
	settingsTab := gui.createSettingsTab()
	environmentTab := gui.createEnvironmentTab()

	gui.tabContainer = container.NewAppTabs(
		toolsTab,
		verifyTab,
		settingsTab,
		environmentTab,
	)
}

// createToolsTab creates the tools management tab
func (gui *DevPathProGUI) createToolsTab() *container.TabItem {
	// Создаем вертикальный контейнер для всех элементов
	mainContainer := container.NewVBox()

	// Группируем программы по категориям
	categories := make(map[string][]config.Program)
	for _, prog := range gui.config.Programs {
		categories[prog.Category] = append(categories[prog.Category], prog)
	}

	// Для каждой категории создаем группу
	for category, programs := range categories {
		// Создаем заголовок категории
		categoryLabel := widget.NewLabelWithStyle(category, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		mainContainer.Add(categoryLabel)

		// Создаем контейнер для программ в этой категории
		programsContainer := container.NewVBox()

		// Добавляем каждую программу
		for _, prog := range programs {
			program := prog // Локальная копия для замыкания

			// Создаем контейнер для программы
			progContainer := container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(program.Name),
				widget.NewButton("Configure", func() {
					gui.configureTool(program)
				}),
			)
			programsContainer.Add(progContainer)
		}

		mainContainer.Add(programsContainer)
		mainContainer.Add(widget.NewSeparator()) // Разделитель между категориями
	}

	// Добавляем кнопку поиска внизу
	searchBtn := widget.NewButtonWithIcon("Deep Search", theme.SearchIcon(), func() {
		gui.performDeepSearch()
	})

	// Создаем скролируемый контейнер
	scroll := container.NewVScroll(mainContainer)

	// Создаем основной контейнер с кнопкой внизу
	content := container.NewBorder(nil, searchBtn, nil, nil, scroll)

	return container.NewTabItem("Tools", content)
}

// createVerifyTab creates the configuration verification tab
func (gui *DevPathProGUI) createVerifyTab() *container.TabItem {
	// Create issues list
	issuesList := widget.NewTextGrid()

	// Create verify button
	verifyBtn := widget.NewButton("Verify Configuration", func() {
		issues := config.VerifyConfigurations()
		gui.displayIssues(issues, issuesList)
	})

	// Create fix button
	fixBtn := widget.NewButton("Fix Issues", func() {
		issues := config.VerifyConfigurations()
		if err := config.FixConfigurationIssues(issues); err != nil {
			gui.showError("Error fixing issues", err)
		} else {
			gui.showSuccess("Issues fixed successfully")
		}
	})

	buttons := container.NewHBox(verifyBtn, fixBtn)

	return container.NewTabItem("Verify",
		container.NewBorder(nil, buttons, nil, nil, issuesList))
}

// createSettingsTab creates the settings management tab
func (gui *DevPathProGUI) createSettingsTab() *container.TabItem {
	// Create settings form
	form := widget.NewForm()

	// Add settings fields here

	return container.NewTabItem("Settings", form)
}

// createEnvironmentTab creates the environment variables tab
func (gui *DevPathProGUI) createEnvironmentTab() *container.TabItem {
	// Создаем вертикальный контейнер для всех элементов
	mainContainer := container.NewVBox()

	// Добавляем заголовки
	headerContainer := container.NewHBox(
		widget.NewLabelWithStyle("Program", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Variable", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Value", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	mainContainer.Add(headerContainer)
	mainContainer.Add(widget.NewSeparator())

	// Создаем карту для хранения лейблов со значениями
	valueLabels := make(map[string]*widget.Label)

	// Добавляем строки для каждой программы
	for _, prog := range gui.config.Programs {
		valueLabel := widget.NewLabel("Searching...")
		valueLabels[prog.Name] = valueLabel

		rowContainer := container.NewHBox(
			widget.NewLabel(prog.Name),
			widget.NewLabel(prog.Name+"_HOME"),
			valueLabel,
		)
		mainContainer.Add(rowContainer)
	}

	// Функция для обновления значений
	updateValues := func() {
		for _, prog := range gui.config.Programs {
			label := valueLabels[prog.Name]
			label.SetText("Searching...")

			go func(p config.Program, l *widget.Label) {
				paths := tools.FindProgram(p)
				if len(paths) > 0 {
					l.SetText(filepath.Dir(filepath.Dir(paths[0])))
				} else {
					l.SetText("Not found")
				}
			}(prog, label)
		}
	}

	// Кнопка обновления
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		updateValues()
	})

	// Запускаем первоначальное обновление
	updateValues()

	// Создаем скролируемый контейнер
	scroll := container.NewVScroll(mainContainer)

	// Создаем основной контейнер с кнопкой вверху
	content := container.NewBorder(refreshBtn, nil, nil, nil, scroll)

	return container.NewTabItem("Environment", content)
}

// processSelectedTools processes the selected tools
func (gui *DevPathProGUI) processSelectedTools() {
	// TODO: Implement tool processing
	// This will use the existing tools.ProcessTools function
}

// displayIssues shows the configuration issues in the UI
func (gui *DevPathProGUI) displayIssues(issues []config.ConfigurationIssue, grid *widget.TextGrid) {
	var text string
	for _, issue := range issues {
		text += fmt.Sprintf("[%s] %s: %s\nSolution: %s\n\n",
			issue.Severity, issue.Type, issue.Description, issue.Solution)
	}
	grid.SetText(text)
}

// showError displays an error message
func (gui *DevPathProGUI) showError(title string, err error) {
	dialog := widget.NewModalPopUp(
		widget.NewLabel(fmt.Sprintf("Error: %v", err)),
		gui.window.Canvas(),
	)
	dialog.Show()
}

// showSuccess displays a success message
func (gui *DevPathProGUI) showSuccess(message string) {
	dialog := widget.NewModalPopUp(
		widget.NewLabel(message),
		gui.window.Canvas(),
	)
	dialog.Show()
}

func (gui *DevPathProGUI) configureTool(prog config.Program) {
	paths := tools.FindProgram(prog)
	if len(paths) == 0 {
		dialog.ShowInformation("Search",
			fmt.Sprintf("%s not found in standard locations. Would you like to perform a deep search?", prog.Name),
			gui.window)
		return
	}

	// Create path selection dialog
	var selectedPath string
	pathOptions := widget.NewRadioGroup(paths, func(value string) {
		selectedPath = value
	})

	configDialog := dialog.NewCustom(
		fmt.Sprintf("Configure %s", prog.Name),
		"Configure",
		container.NewVBox(
			widget.NewLabel("Select installation path:"),
			pathOptions,
		),
		gui.window,
	)

	configDialog.SetOnClosed(func() {
		if selectedPath == "" {
			return
		}

		// Get available configuration options
		options := tools.GetConfigOptions(prog)
		if len(options) == 0 {
			// If no special options, just configure
			if err := tools.ConfigureSelectedPath(prog, selectedPath); err != nil {
				dialog.ShowError(err, gui.window)
			} else {
				dialog.ShowInformation("Success",
					fmt.Sprintf("%s configured successfully", prog.Name),
					gui.window)
			}
			return
		}

		// Create options selection dialog
		var selectedVars []string
		optionsContainer := container.NewVBox()

		for _, opt := range options {
			check := widget.NewCheck(opt.Name, func(checked bool) {
				if checked {
					selectedVars = append(selectedVars, opt.Variables...)
				} else {
					// Remove variables from selected
					newVars := make([]string, 0)
					for _, v := range selectedVars {
						found := false
						for _, ov := range opt.Variables {
							if v == ov {
								found = true
								break
							}
						}
						if !found {
							newVars = append(newVars, v)
						}
					}
					selectedVars = newVars
				}
			})
			optionsContainer.Add(container.NewHBox(
				check,
				widget.NewLabel(opt.Description),
			))
		}

		optionsDialog := dialog.NewCustom(
			"Select Configuration Options",
			"Apply",
			optionsContainer,
			gui.window,
		)

		optionsDialog.SetOnClosed(func() {
			if err := tools.ConfigureSelectedPath(prog, selectedPath); err != nil {
				dialog.ShowError(err, gui.window)
			} else {
				dialog.ShowInformation("Success",
					fmt.Sprintf("%s configured successfully", prog.Name),
					gui.window)
			}
		})

		optionsDialog.Show()
	})

	configDialog.Show()
}

func (gui *DevPathProGUI) performDeepSearch() {
	progress := widget.NewProgressBarInfinite()
	progressDialog := dialog.NewCustom(
		"Deep Search",
		"Cancel",
		container.NewVBox(
			widget.NewLabel("Searching in all drives..."),
			progress,
		),
		gui.window,
	)

	go func() {
		results := tools.ProcessToolsDeepSearch(gui.config.Programs)
		gui.window.Content().Refresh()
		progressDialog.Hide()

		// Show results
		var text string
		for _, result := range results {
			if result.Found {
				text += fmt.Sprintf("✅ %s found:\n", result.Program.Name)
				for _, path := range result.Paths {
					text += fmt.Sprintf("  - %s\n", path)
				}
			} else {
				text += fmt.Sprintf("❌ %s not found\n", result.Program.Name)
			}
			if result.Error != nil {
				text += fmt.Sprintf("  Error: %v\n", result.Error)
			}
			text += "\n"
		}

		resultGrid := widget.NewTextGrid()
		resultGrid.SetText(text)
		dialog.ShowCustom(
			"Search Results",
			"Close",
			resultGrid,
			gui.window,
		)
	}()

	progressDialog.Show()
}
