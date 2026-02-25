package validators

import (
	"mime/multipart"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func newTestValidator() *validator.Validate {
	v := validator.New()
	RegisterCustomValidations(v)
	return v
}

func TestValidateDigits(t *testing.T) {
	v := newTestValidator()
	require.NoError(t, v.Var("123456", "digits=6"))
	require.Error(t, v.Var("12345a", "digits=6"))
}

func TestValidateBetween(t *testing.T) {
	v := newTestValidator()
	require.NoError(t, v.Var("abcd", "between=3_4"))
	require.Error(t, v.Var("ab", "between=3_4"))
}

func TestValidateNumericBetween(t *testing.T) {
	v := newTestValidator()
	require.NoError(t, v.Var("10", "numeric_between=2_100"))
	require.Error(t, v.Var("1", "numeric_between=2_100"))
}

func TestValidateInNotIn(t *testing.T) {
	v := newTestValidator()
	require.NoError(t, v.Var("asc", "in=asc_desc"))
	require.Error(t, v.Var("foo", "in=asc_desc"))
	require.NoError(t, v.Var("foo", "not_in=asc_desc"))
	require.Error(t, v.Var("asc", "not_in=asc_desc"))
}

func TestValidateCnLength(t *testing.T) {
	v := newTestValidator()
	require.NoError(t, v.Var("中文", "min_cn=2"))
	require.Error(t, v.Var("中文", "max_cn=1"))
}

func TestValidateFileRuleValue(t *testing.T) {
	file := &multipart.FileHeader{
		Filename: "avatar.png",
		Size:     1024,
	}
	require.True(t, ValidateFileRuleValue("ext", file, "png_jpg"))
	require.False(t, ValidateFileRuleValue("ext", file, "gif"))
	require.True(t, ValidateFileRuleValue("size", file, "2048"))
	require.False(t, ValidateFileRuleValue("size", file, "10"))
}
