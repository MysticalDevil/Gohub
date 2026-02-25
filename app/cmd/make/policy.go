package make

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gohub/pkg/logger"
)

var CmdMakePolicy = &cobra.Command{
	Use:   "policy",
	Short: "Create policy file, example: make policy user",
	Run:   runMakePolicy,
	Args:  cobra.ExactArgs(1),
}

func runMakePolicy(_ *cobra.Command, args []string) {
	model := makeModelFromString(args[0])

	err := os.MkdirAll("app/policies", os.ModePerm)
	if err != nil {
		logger.LogIf(err)
	}

	filePath := fmt.Sprintf("app/policies/%s_policy.go", model.PackageName)

	createFileFromStub(filePath, "policy", model)
}
