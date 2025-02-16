package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"devpathpro/pkg/config"
	"devpathpro/pkg/tools"
	"devpathpro/pkg/registry"
)

type GUI struct {
	app    fyne.App
	window fyne.Window
	config *config.Configuration
}

func NewGUI() *GUI {
	return &GUI{
		app: app.New(),
		config: &config.Configuration{
			LogFile:  "devpathpro.log",
			Programs: config.GetDefaultPrograms(),
		},
	}
}

func (g *GUI) Run() {
	// Создаем временное окно для проверки прав
	tempWindow := g.app.NewWindow("DevPathPro")
	
	// Проверяем права администратора
	if !registry.IsAdmin() {
		dialog.ShowInformation("Administrator Rights Required", 
			"This program requires administrator privileges to modify system environment variables.\nPlease restart the program as administrator.", 
			tempWindow)
		tempWindow.Show()
		g.app.Run()
		return
	}
	
	// Закрываем временное окно
	tempWindow.Close()

	// Создаем основное окно с заданным размером
	g.window = g.app.NewWindow("DevPathPro")
	g.window.Resize(fyne.NewSize(300, 400))

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Tools", g.createToolsTab()),
		container.NewTabItem("Verify", g.createVerifyTab()),
		container.NewTabItem("Environment", g.createEnvironmentTab()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	// Создаем контейнер с отступами и фиксированным размером
	content := container.NewMax(container.NewPadded(tabs))

	// Устанавливаем контент
	g.window.SetContent(content)
	
	// Центрируем окно на экране
	g.window.CenterOnScreen()

	// Запускаем приложение
	g.window.Show()
	g.app.Run()
}

func (g *GUI) createToolsTab() fyne.CanvasObject {
	// Создаем вертикальный контейнер для всех элементов
	mainContainer := container.NewVBox()

	// Группируем программы по категориям
	categories := make(map[string][]config.Program)
	for _, prog := range g.config.Programs {
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

			// Создаем элементы управления
			check := widget.NewCheck("", nil)
			label := widget.NewLabel(program.Name)
			configBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
				g.configureTool(program)
			})

			// Создаем контейнер для программы
			progContainer := container.NewHBox(
				check,
				label,
				configBtn,
			)

			programsContainer.Add(progContainer)
		}

		mainContainer.Add(programsContainer)
		mainContainer.Add(widget.NewSeparator())
	}

	// Добавляем кнопку поиска внизу
	searchBtn := widget.NewButtonWithIcon("Deep Search", theme.SearchIcon(), func() {
		g.performDeepSearch()
	})

	// Создаем скролируемый контейнер
	scroll := container.NewVScroll(mainContainer)

	// Создаем основной контейнер с кнопкой внизу
	content := container.NewBorder(nil, searchBtn, nil, nil, scroll)

	return content
}

func (g *GUI) createVerifyTab() fyne.CanvasObject {
	// Создаем контейнер для результатов с прокруткой
	results := widget.NewTextGrid()
	resultsScroll := container.NewVScroll(results)

	// Создаем индикатор прогресса
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// Создаем канал для отмены операции
	done := make(chan bool)

	// Предварительно объявляем кнопки
	var verifyBtn *widget.Button
	var fixBtn *widget.Button
	var cancelBtn *widget.Button

	// Кнопка отмены
	cancelBtn = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		close(done)
		progress.Hide()
		verifyBtn.Enable()
		fixBtn.Enable()
		results.SetText("Operation cancelled")
		results.Refresh()
	})

	// Кнопка исправления
	fixBtn = widget.NewButtonWithIcon("Fix Issues", theme.MediaReplayIcon(), func() {
		// Блокируем кнопки во время исправления
		verifyBtn.Disable()
		fixBtn.Disable()

		// Показываем индикатор прогресса
		progress.Show()
		results.SetText("Fixing issues...")
		results.Refresh()

		// Запускаем исправление в отдельной горутине
		go func() {
			defer func() {
				// Разблокируем кнопки после завершения
				verifyBtn.Enable()
				fixBtn.Enable()
				progress.Hide()
				verifyBtn.Refresh()
				fixBtn.Refresh()
				progress.Refresh()
			}()

			select {
			case <-done:
				return
			default:
				issues := config.VerifyConfigurations()

				// Проверяем, не была ли операция отменена
				select {
				case <-done:
					return
				default:
					if err := config.FixConfigurationIssues(issues); err != nil {
						dialog.ShowError(fmt.Errorf("Fix error: %v", err), g.window)
					} else {
						dialog.ShowInformation("Success", "Issues fixed successfully!", g.window)
						g.showRestartDialog()
					}
				}
			}
		}()
	})

	// Кнопка проверки
	verifyBtn = widget.NewButtonWithIcon("Verify", theme.ConfirmIcon(), func() {
		// Блокируем кнопки во время проверки
		verifyBtn.Disable()
		fixBtn.Disable()

		// Показываем индикатор прогресса
		progress.Show()
		results.SetText("Checking configuration...")
		results.Refresh()

		// Запускаем проверку в отдельной горутине
		go func() {
			defer func() {
				// Разблокируем кнопки после завершения
				verifyBtn.Enable()
				fixBtn.Enable()
				progress.Hide()
				verifyBtn.Refresh()
				fixBtn.Refresh()
				progress.Refresh()
			}()

			select {
			case <-done:
				return
			default:
				issues := config.VerifyConfigurations()

				// Проверяем, не была ли операция отменена
				select {
				case <-done:
					return
				default:
					if len(issues) == 0 {
						results.SetText("✅ All checks passed successfully!")
					} else {
						text := fmt.Sprintf("Found %d issues:\n\n", len(issues))
						for _, issue := range issues {
							text += fmt.Sprintf("[%s] %s: %s\nSolution: %s\n\n",
								issue.Severity, issue.Type, issue.Description, issue.Solution)
						}
						results.SetText(text)
					}
					results.Refresh()
				}
			}
		}()
	})

	// Создаем контейнер с кнопками
	buttons := container.NewHBox(verifyBtn, fixBtn, cancelBtn)

	// Создаем основной контейнер
	content := container.NewBorder(
		buttons,  // top
		progress, // bottom
		nil,      // left
		nil,      // right
		resultsScroll,
	)

	return content
}

func (g *GUI) createEnvironmentTab() fyne.CanvasObject {
	// Создаем основной контейнер с отступами
	mainContainer := container.NewPadded(container.NewVBox())

	// Добавляем заголовки
	headerContainer := container.NewHBox(
		widget.NewLabelWithStyle("Program", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Variable", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Value", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	mainContainer.Add(headerContainer)
	mainContainer.Add(widget.NewSeparator())

	// Создаем контейнер для программ
	programsContainer := container.NewVBox()

	// Создаем карту для хранения лейблов со значениями
	valueLabels := make(map[string]*widget.Label)

	// Создаем индикатор прогресса
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// Функция для обновления значений
	updateValues := func() {
		progress.Show()
		for _, prog := range g.config.Programs {
			label := valueLabels[prog.Name]
			label.SetText("Searching...")
			label.Refresh()

			go func(p config.Program, l *widget.Label) {
				defer func() {
					// Проверяем, все ли метки обновлены
					allUpdated := true
					for _, label := range valueLabels {
						if label.Text == "Searching..." {
							allUpdated = false
							break
						}
					}
					if allUpdated {
						progress.Hide()
						progress.Refresh()
					}
				}()

				paths := tools.FindProgram(p)
				if len(paths) > 0 {
					// Берем директорию программы (один уровень вверх)
					path := filepath.Dir(paths[0])
					
					// Проверяем, является ли это bin директорией
					if strings.EqualFold(filepath.Base(path), "bin") {
						// Если это bin, берем родительскую директорию
						path = filepath.Dir(path)
					}
					
					// Получаем переменную окружения
					envValue := os.Getenv(p.Name + "_HOME")
					if envValue != "" {
						// Если переменная окружения установлена, показываем её
						l.SetText(envValue)
					} else {
						// Иначе показываем найденный путь
						l.SetText(path)
					}
				} else {
					l.SetText("Not found")
				}
				l.Refresh()
			}(prog, label)
		}
	}

	// Кнопка обновления
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		updateValues()
	})

	// Добавляем строки для каждой программы
	for _, prog := range g.config.Programs {
		valueLabel := widget.NewLabel("...")
		valueLabels[prog.Name] = valueLabel

		// Создаем контейнер для строки с фиксированными размерами
		rowContainer := container.NewHBox(
			container.NewHBox(widget.NewLabelWithStyle(prog.Name, fyne.TextAlignLeading, fyne.TextStyle{})),
			container.NewHBox(widget.NewLabelWithStyle(prog.Name+"_HOME", fyne.TextAlignLeading, fyne.TextStyle{})),
			container.NewHBox(valueLabel),
		)
		programsContainer.Add(rowContainer)
	}

	// Создаем прокручиваемый контейнер
	scroll := container.NewVScroll(programsContainer)

	// Запускаем первоначальное обновление
	updateValues()

	// Создаем контейнер для кнопок с отступами
	topContainer := container.NewPadded(
		container.NewHBox(
			refreshBtn,
			progress,
		),
	)

	// Создаем основной контейнер
	content := container.NewBorder(topContainer, nil, nil, nil, scroll)

	return content
}

// showRestartDialog показывает диалог с предложением перезагрузить компьютер
func (g *GUI) showRestartDialog() {
	confirmDialog := dialog.NewConfirm(
		"Restart Required",
		"Changes to environment variables require a system restart to take effect.\nWould you like to restart now?",
		func(restart bool) {
			if restart {
				if err := exec.Command("shutdown", "/r", "/t", "0").Run(); err != nil {
					dialog.ShowError(fmt.Errorf("Failed to restart: %v", err), g.window)
				}
			}
		},
		g.window,
	)
	confirmDialog.Show()
}

func (g *GUI) configureTool(prog config.Program) {
	paths := tools.FindProgram(prog)
	if len(paths) == 0 {
		dialog.ShowInformation("Search",
			fmt.Sprintf("%s not found in standard locations. Would you like to perform a deep search?", prog.Name),
			g.window)
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
		g.window,
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
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Success",
					fmt.Sprintf("%s configured successfully", prog.Name),
					g.window)
				g.showRestartDialog()
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
			g.window,
		)

		optionsDialog.SetOnClosed(func() {
			if err := tools.ConfigureSelectedPath(prog, selectedPath); err != nil {
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Success",
					fmt.Sprintf("%s configured successfully", prog.Name),
					g.window)
				g.showRestartDialog()
			}
		})

		optionsDialog.Show()
	})

	configDialog.Show()
}

func (g *GUI) performDeepSearch() {
	progress := widget.NewProgressBarInfinite()
	progressDialog := dialog.NewCustom(
		"Deep Search",
		"Cancel",
		container.NewVBox(
			widget.NewLabel("Searching in all drives..."),
			progress,
		),
		g.window,
	)

	go func() {
		results := tools.ProcessToolsDeepSearch(g.config.Programs)
		g.window.Content().Refresh()
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
			g.window,
		)
	}()

	progressDialog.Show()
}
