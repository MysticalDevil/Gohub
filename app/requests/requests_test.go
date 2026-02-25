package requests

import (
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type validateTestStruct struct {
	Name   string                `json:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func newTestContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c
}

func TestValidateSkipsEmptyNonRequired(t *testing.T) {
	data := &validateTestStruct{Name: ""}
	rules := MapData{
		"name": []string{"min=3"},
	}
	messages := MapData{
		"name": []string{"min:Name too short"},
	}
	errMap := validate(newTestContext(), data, rules, messages)
	require.Empty(t, errMap)
}

func TestValidateRequired(t *testing.T) {
	data := &validateTestStruct{Name: ""}
	rules := MapData{
		"name": []string{"required"},
	}
	messages := MapData{
		"name": []string{"required:Name required"},
	}
	errMap := validate(newTestContext(), data, rules, messages)
	require.NotEmpty(t, errMap["name"])
	require.Equal(t, "Name required", errMap["name"][0])
}

func TestValidateMinRule(t *testing.T) {
	data := &validateTestStruct{Name: "ab"}
	rules := MapData{
		"name": []string{"min=3"},
	}
	messages := MapData{
		"name": []string{"min:Name too short"},
	}
	errMap := validate(newTestContext(), data, rules, messages)
	require.NotEmpty(t, errMap["name"])
}

func TestValidateFileRequired(t *testing.T) {
	data := &validateTestStruct{Avatar: nil}
	rules := MapData{
		"file:avatar": []string{"required"},
	}
	messages := MapData{
		"avatar": []string{"required:Avatar required"},
	}
	errMap := validateFile(newTestContext(), data, rules, messages)
	require.NotEmpty(t, errMap["avatar"])
}
