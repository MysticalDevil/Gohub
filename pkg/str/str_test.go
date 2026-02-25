package str

import "testing"

func TestPluralSingular(t *testing.T) {
	if got := Plural("user"); got != "users" {
		t.Fatalf("unexpected plural: %s", got)
	}
	if got := Singular("users"); got != "user" {
		t.Fatalf("unexpected singular: %s", got)
	}
}

func TestCaseConversions(t *testing.T) {
	if got := Snake("TopicComment"); got != "topic_comment" {
		t.Fatalf("unexpected snake: %s", got)
	}
	if got := Camel("topic_comment"); got != "TopicComment" {
		t.Fatalf("unexpected camel: %s", got)
	}
	if got := LowerCamel("TopicComment"); got != "topicComment" {
		t.Fatalf("unexpected lower camel: %s", got)
	}
}
