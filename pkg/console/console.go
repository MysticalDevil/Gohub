// Package console Command line helper method
package console

import (
	"fmt"
	"github.com/mgutz/ansi"
	"gohub/pkg/logger"
	"os"
)

// Success Print a success message, green output
func Success(msg string) {
	colorOut(msg, "green")
}

// Error Print an error message, red output
func Error(msg string) {
	colorOut(msg, "red")
}

// Warning Print a warning message, yellow output
func Warning(msg string) {
	colorOut(msg, "yellow")
}

// Exit Print an error message, an exit os.Exit(1)
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

// ExitIf Syntactic sugar, comes with err != nil judgment
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

// colorOut For internal use, set the highlight color
func colorOut(message, color string) {
	_, err := fmt.Fprintln(os.Stdout, ansi.Color(message, color))
	if err != nil {
		logger.LogIf(err)
		return
	}
}
