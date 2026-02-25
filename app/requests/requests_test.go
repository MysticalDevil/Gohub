package requests

import (
	"mime/multipart"
	"testing"
)

type validateTestStruct struct {
	Name   string                `json:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func TestValidateSkipsEmptyNonRequired(t *testing.T) {
	data := &validateTestStruct{Name: ""}
	rules := MapData{
		"name": []string{"min=3"},
	}
	messages := MapData{
		"name": []string{"min:Name too short"},
	}
	errMap := validate(data, rules, messages)
	if len(errMap) != 0 {
		t.Fatalf("expected no errors, got %v", errMap)
	}
}

func TestValidateRequired(t *testing.T) {
	data := &validateTestStruct{Name: ""}
	rules := MapData{
		"name": []string{"required"},
	}
	messages := MapData{
		"name": []string{"required:Name required"},
	}
	errMap := validate(data, rules, messages)
	if len(errMap["name"]) == 0 {
		t.Fatalf("expected required error")
	}
	if errMap["name"][0] != "Name required" {
		t.Fatalf("unexpected message: %v", errMap["name"][0])
	}
}

func TestValidateMinRule(t *testing.T) {
	data := &validateTestStruct{Name: "ab"}
	rules := MapData{
		"name": []string{"min=3"},
	}
	messages := MapData{
		"name": []string{"min:Name too short"},
	}
	errMap := validate(data, rules, messages)
	if len(errMap["name"]) == 0 {
		t.Fatalf("expected min error")
	}
}

func TestValidateFileRequired(t *testing.T) {
	data := &validateTestStruct{Avatar: nil}
	rules := MapData{
		"file:avatar": []string{"required"},
	}
	messages := MapData{
		"avatar": []string{"required:Avatar required"},
	}
	errMap := validateFile(nil, data, rules, messages)
	if len(errMap["avatar"]) == 0 {
		t.Fatalf("expected avatar required error")
	}
}
