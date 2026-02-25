package validators

import (
	"mime/multipart"
	"testing"

	"github.com/go-playground/validator/v10"
)

func newTestValidator() *validator.Validate {
	v := validator.New()
	RegisterCustomValidations(v)
	return v
}

func TestValidateDigits(t *testing.T) {
	v := newTestValidator()
	if err := v.Var("123456", "digits=6"); err != nil {
		t.Fatalf("expected valid digits, got %v", err)
	}
	if err := v.Var("12345a", "digits=6"); err == nil {
		t.Fatalf("expected invalid digits")
	}
}

func TestValidateBetween(t *testing.T) {
	v := newTestValidator()
	if err := v.Var("abcd", "between=3_4"); err != nil {
		t.Fatalf("expected valid between, got %v", err)
	}
	if err := v.Var("ab", "between=3_4"); err == nil {
		t.Fatalf("expected invalid between")
	}
}

func TestValidateNumericBetween(t *testing.T) {
	v := newTestValidator()
	if err := v.Var("10", "numeric_between=2_100"); err != nil {
		t.Fatalf("expected valid numeric_between, got %v", err)
	}
	if err := v.Var("1", "numeric_between=2_100"); err == nil {
		t.Fatalf("expected invalid numeric_between")
	}
}

func TestValidateInNotIn(t *testing.T) {
	v := newTestValidator()
	if err := v.Var("asc", "in=asc_desc"); err != nil {
		t.Fatalf("expected valid in, got %v", err)
	}
	if err := v.Var("foo", "in=asc_desc"); err == nil {
		t.Fatalf("expected invalid in")
	}
	if err := v.Var("foo", "not_in=asc_desc"); err != nil {
		t.Fatalf("expected valid not_in, got %v", err)
	}
	if err := v.Var("asc", "not_in=asc_desc"); err == nil {
		t.Fatalf("expected invalid not_in")
	}
}

func TestValidateCnLength(t *testing.T) {
	v := newTestValidator()
	if err := v.Var("中文", "min_cn=2"); err != nil {
		t.Fatalf("expected valid min_cn, got %v", err)
	}
	if err := v.Var("中文", "max_cn=1"); err == nil {
		t.Fatalf("expected invalid max_cn")
	}
}

func TestValidateFileRuleValue(t *testing.T) {
	file := &multipart.FileHeader{
		Filename: "avatar.png",
		Size:     1024,
	}
	if !ValidateFileRuleValue("ext", file, "png_jpg") {
		t.Fatalf("expected valid ext")
	}
	if ValidateFileRuleValue("ext", file, "gif") {
		t.Fatalf("expected invalid ext")
	}
	if !ValidateFileRuleValue("size", file, "2048") {
		t.Fatalf("expected valid size")
	}
	if ValidateFileRuleValue("size", file, "10") {
		t.Fatalf("expected invalid size")
	}
}
