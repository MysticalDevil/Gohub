// Package make Command line make command
package make

import (
	"embed"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"gohub/pkg/console"
	"gohub/pkg/file"
	"gohub/pkg/str"
)

// Model parameter explanation
// A single word, the user command is passed as a parameter, taking the User model as an example:
// - user
// - User
// - users
// - Users
// organized data:
//
//	{
//		   "TableName": "users",
//		   "StructName": "User",
//		   "StructNamePlural": "Users",
//		   "VariableName": "user",
//		   "VariableNamePlural": "users",
//		   "PackageName": "user",
//	}
//
// -
// Two words or more, the user command to pass parameters, take the TopicComment model as an example:
// - topic_comment
// - topic_comments
// - TopicComment
// - TopicComments
// organized data:
//
//	{
//		   "TableName" :"topic_comments",
//		   "StructName": "TopicComment",
//		   "StructNamePlural": "TopicComments",
//		   "VariableName": "topicComment",
//		   "VariableNamePlural": "topicComments",
//		   "PackageName": "topic_comment"
//	}
type Model struct {
	TableName          string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNamePlural string
	PackageName        string
}

// stubsFS Convenient packaging of .stub files

//go:embed stubs
var stubsFS embed.FS

// Make Explain the cobra command
var Make = &cobra.Command{
	Use:   "make",
	Short: "Generate file nad code",
}

func init() {
	// Register make subcommands
	Make.AddCommand(
		CMD,
		CmdMakeApiController,
		CmdMakeFactory,
		CmdMakeMigration,
		CmdMakeModel,
		CmdMakePolicy,
		CmdMakeRequest,
		CmdMakeSeeder,
	)
}

// make ModelFromString Format user input
func makeModelFromString(name string) Model {
	model := Model{}

	model.StructName = str.Singular(strcase.ToCamel(name))
	model.StructNamePlural = str.Plural(model.StructName)
	model.VariableName = str.LowerCamel(model.StructName)
	model.VariableNamePlural = str.LowerCamel(model.StructNamePlural)
	model.TableName = str.Snake(model.StructNamePlural)
	model.PackageName = str.Snake(model.StructName)

	return model
}

// createFileFromStub Read the stub file and perform variable substitution
// The last option is optional. If you pass a parameter,
// you should pass in the map[string]string type as an additional variable search and replacement
func createFileFromStub(filePath, stubName string, model Model, variables ...any) {
	// implement the last parameter optional
	replaces := make(map[string]string)
	if len(variables) > 0 {
		replaces = variables[0].(map[string]string)
	}

	// target file already exists
	if file.Exists(filePath) {
		console.Exit(filePath + " already exists!")
	}

	// read the stub template file
	modelData, err := stubsFS.ReadFile("stubs/" + stubName + ".stub")
	if err != nil {
		console.Exit(err.Error())
	}
	modelStub := string(modelData)

	// add default substitution variable
	replaces["{{VariableName}}"] = model.VariableName
	replaces["{{VariableNamePlural}}"] = model.VariableNamePlural
	replaces["{{StructName}}"] = model.StructName
	replaces["{{StructNamePlural}}"] = model.StructNamePlural
	replaces["{{PackageName}}"] = model.PackageName
	replaces["{{TableName}}"] = model.TableName

	// variable substitution for template content
	for search, replace := range replaces {
		modelStub = strings.ReplaceAll(modelStub, search, replace)
	}

	// save to target file
	err = file.Put([]byte(modelStub), filePath)
	if err != nil {
		console.Exit(err.Error())
	}

	// prompt success
	console.Success(fmt.Sprintf("[%s] created.", filePath))
}
