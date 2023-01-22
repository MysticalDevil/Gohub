package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/console"
	"os"
)

var CmdMakeModel = &cobra.Command{
	Use:   "model",
	Short: "Create model file, example: make model user",
	Run:   runMakeModel,
	Args:  cobra.ExactArgs(1),
}

func runMakeModel(_ *cobra.Command, args []string) {
	// format the model name and return a Model object
	model := makeModelFromString(args[0])

	// make sure the models directory exists, e.g. `app/models/user`
	dir := fmt.Sprintf("app/models/%s/", model.PackageName)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		console.Exit(err.Error())
	}

	createFileFromStub(dir+model.PackageName+"_model.go", "model/model", model)
	createFileFromStub(dir+model.PackageName+"_util.go", "model/model_util", model)
	createFileFromStub(dir+model.PackageName+"_hooks.go", "model/model_hooks", model)
}
