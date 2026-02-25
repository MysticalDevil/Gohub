package paginator

import (
	"os"
	"testing"

	appconfig "gohub/config"
	pkgconfig "gohub/pkg/config"
)

func initTestConfig(t *testing.T) {
	t.Helper()

	if err := os.WriteFile(".env", []byte("APP_ENV=local\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(".env")
	})

	appconfig.Initialize()
	pkgconfig.InitConfig("")
}

func TestFormatBaseURL(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{}
	if got := p.formatBaseURL("/topics"); got != "/topics?page=" {
		t.Fatalf("unexpected url: %s", got)
	}
	if got := p.formatBaseURL("/topics?order=asc"); got != "/topics?order=asc&page=" {
		t.Fatalf("unexpected url: %s", got)
	}
}

func TestGetPageLink(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{
		BaseURL: "/topics?page=",
		Sort:    "created_at",
		Order:   "desc",
		PerPage: 20,
	}
	if got := p.getPageLink(2); got != "/topics?page=2&sort=created_at&order=desc&per_page=20" {
		t.Fatalf("unexpected page link: %s", got)
	}
}

func TestPrevNextLinks(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{
		BaseURL:   "/topics?page=",
		Sort:      "id",
		Order:     "asc",
		PerPage:   10,
		TotalPage: 3,
		Page:      2,
	}
	if got := p.getPrevPageURL(); got != "/topics?page=1&sort=id&order=asc&per_page=10" {
		t.Fatalf("unexpected prev link: %s", got)
	}
	if got := p.getNextPageURL(); got != "/topics?page=3&sort=id&order=asc&per_page=10" {
		t.Fatalf("unexpected next link: %s", got)
	}
}

func TestPrevNextLinksBounds(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{BaseURL: "/topics?page=", Sort: "id", Order: "asc", PerPage: 10, TotalPage: 1, Page: 1}
	if got := p.getPrevPageURL(); got != "" {
		t.Fatalf("expected empty prev link")
	}
	if got := p.getNextPageURL(); got != "" {
		t.Fatalf("expected empty next link")
	}
}
