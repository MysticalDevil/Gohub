package make

import (
	"fmt"
	"github.com/spf13/cobra"
)

var CmdMakeSeeder = &cobra.Command{
	Use:   "seeder",
	Short: "Create seeder file, example: make seeder user",
	Run:   runMakeSeeder,
	Args:  cobra.ExactArgs(1),
}

func runMakeSeeder(_ *cobra.Command, args []string) {
	model := makeModelFromString(args[0])

	filePath := fmt.Sprintf("database/seeders/%s_seeder.go", model.TableName)

	createFileFromStub(filePath, "seeder", model)
}
