package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sawitos/sawit/internal/config"
	"github.com/sawitos/sawit/internal/health"
)

var Version = "0.1.0-dev"

func main() {
	jsonFlag := flag.Bool("json", false, "Output in machine-readable JSON format")
	checkAllFlag := flag.Bool("check-all", true, "Run all health inspection modules")
	configPath := flag.String("config", "/etc/sawit/sawitd.yaml", "Path to config file")
	versionFlag := flag.Bool("version", false, "Print sawit-health version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("sawit-health version %s\n", Version)
		os.Exit(0)
	}

	_ = checkAllFlag

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		cfg = config.Default()
	}

	report, err := health.Evaluate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error evaluating system health: %v\n", err)
		os.Exit(1)
	}

	if *jsonFlag {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Health Inspector (sawit-health)")
	fmt.Println("──────────────────────────────────────────")
	for _, check := range report.Checks {
		icon := "[✓]"
		if check.Status == health.StatusWarning {
			icon = "[!]"
		} else if check.Status == health.StatusCritical {
			icon = "[X]"
		}
		fmt.Printf("%-12s %s  %s\n", check.Name, icon, check.Message)
	}
	fmt.Println("──────────────────────────────────────────")
	fmt.Printf("Overall Health: %s\n", report.Overall)

	if report.Overall == health.StatusCritical {
		os.Exit(2)
	} else if report.Overall == health.StatusWarning {
		os.Exit(1)
	}
	os.Exit(0)
}
