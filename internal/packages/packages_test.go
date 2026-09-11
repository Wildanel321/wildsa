package packages_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/packages"
)

func TestSearchPackages(t *testing.T) {
	results, err := packages.SearchPackages("nginx")
	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected search results for 'nginx'")
	}
	if results[0].Name == "" {
		t.Error("expected non-empty package name")
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	_, err := packages.SearchPackages("")
	if err == nil {
		t.Error("expected error for empty search query")
	}
}
