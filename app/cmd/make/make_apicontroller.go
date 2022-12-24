package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/console"
	"strings"
)

var ApiController = &cobra.Command{
	Use:   "apicontroller",
	Short: "Create api controller, example: make apicontroller v1/user",
	Run:   runMakeApiController,
	Args:  cobra.ExactArgs(1),
}

func runMakeApiController(cmd *cobra.Command, args []string) {
	// handle parameters, require an accompanying API version (v1 or v2)
	array := strings.Split(args[0], "/")
	if len(array) != 2 {
		console.Exit("api controller name format: v1/user")
	}

	// apiVersion is used to splice the target path
	apiVersion, name := array[0], array[1]
	model := makeModelFromString(name)

	// build target directory
	filePath := fmt.Sprintf("app/http/controllers/api/%s/%s_controller.go", apiVersion, model.TableName)

	createFileFromStub(filePath, "apicontroller", model)
}
