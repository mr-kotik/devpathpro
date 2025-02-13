package utils

import (
	"os"
	"os/exec"
)

// ClearScreen clears the console screen
func ClearScreen() error {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// PrintDivider prints a divider line
func PrintDivider(char string, length int) {
	for i := 0; i < length; i++ {
		print(char)
	}
	println()
} 