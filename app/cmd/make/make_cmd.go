package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/console"
)

var CMD = &cobra.Command{
	Use:   "cmd",
	Short: "Create a command, should be snake_case, example: make cmd backup_database",
	Run:   runMakeCMD,
	Args:  cobra.ExactArgs(1), // Only one parameter is allowed and must be passed
}

func runMakeCMD(cmd *cobra.Command, args []string) {
	// format the model name and return a Model object
	model := makeModelFromString(args[0])

	// concatenate target file paths
	filePath := fmt.Sprintf("app/cmd/%s.go", model.PackageName)

	// create files from templates (with variable substitution done)
	createFileFromStub(filePath, "cmd", model)

	// friendly reminder
	console.Success("command name: " + model.PackageName)
	console.Success("command variable name: cmd.Cmd " + model.StructName)
	console.Warning("please edit main.go's app.Commands slice to register command")
}
