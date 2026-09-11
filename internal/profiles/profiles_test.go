package profiles_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/profiles"
)

func TestGetProfiles(t *testing.T) {
	profs := profiles.GetProfiles()
	if len(profs) == 0 {
		t.Fatal("expected non-empty profile presets list")
	}
}

func TestGetProfileValid(t *testing.T) {
	p, err := profiles.GetProfile("container")
	if err != nil {
		t.Fatalf("unexpected error getting container profile: %v", err)
	}
	if p.Name != "container" {
		t.Errorf("expected container profile, got %s", p.Name)
	}
}

func TestGetProfileInvalid(t *testing.T) {
	_, err := profiles.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestApplyProfile(t *testing.T) {
	res, err := profiles.ApplyProfile("minimal")
	if err != nil {
		t.Fatalf("unexpected error applying minimal profile: %v", err)
	}
	if !res.Success {
		t.Errorf("expected success true, got false")
	}
}
