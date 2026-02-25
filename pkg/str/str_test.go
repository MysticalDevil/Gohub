package str

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluralSingular(t *testing.T) {
	require.Equal(t, "users", Plural("user"))
	require.Equal(t, "user", Singular("users"))
}

func TestCaseConversions(t *testing.T) {
	require.Equal(t, "topic_comment", Snake("TopicComment"))
	require.Equal(t, "TopicComment", Camel("topic_comment"))
	require.Equal(t, "topicComment", LowerCamel("TopicComment"))
}
