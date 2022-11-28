// Package str String helper methods
package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural Pluralize user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular SConvert to singular users -> user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake Convert to snake_case, such as: TopicComment -> topic_comment
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// Camel Convert to CamelCase, such as: topic_comment -> TopicComment
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// LowerCamel Convert to lowerCamelCase, such as: TopicComment -> topicComment
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}
