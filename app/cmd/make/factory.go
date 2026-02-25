package make

import (
	"fmt"
	"github.com/spf13/cobra"
)

var CmdMakeFactory = &cobra.Command{
	Use:   "factory",
	Short: "Create model's factory file, example: make factory user",
	Run:   runMakeFactory,
	Args:  cobra.ExactArgs(1),
}

func runMakeFactory(_ *cobra.Command, args []string) {
	model := makeModelFromString(args[0])

	filePath := fmt.Sprintf("database/factories/%s_factory.go", model.PackageName)

	createFileFromStub(filePath, "factory", model)
}
