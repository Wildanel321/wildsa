package health_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/config"
	"github.com/sawitos/sawit/internal/health"
)

func TestEvaluateDefaultConfig(t *testing.T) {
	cfg := config.Default()
	report, err := health.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error evaluating health: %v", err)
	}

	if report.Overall == "" {
		t.Error("expected non-empty overall status")
	}
	if len(report.Checks) == 0 {
		t.Error("expected check items in health report")
	}
}

func TestEvaluateCriticalThresholds(t *testing.T) {
	cfg := config.Default()
	cfg.Health.CPUCriticalPercent = 0.0 // Force CPU check to trigger Critical

	report, err := health.Evaluate(cfg)
	if err != nil {
		t.Fatalf("unexpected error evaluating health: %v", err)
	}

	if report.Overall != health.StatusCritical {
		t.Errorf("expected overall status CRITICAL, got %s", report.Overall)
	}
}
