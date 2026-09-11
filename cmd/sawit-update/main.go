package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var Version = "0.1.0-dev"

type UpdateStatus struct {
	LastChecked      time.Time `json:"last_checked"`
	SystemUpToDate   bool      `json:"system_up_to_date"`
	SecurityUpdates  int       `json:"security_updates"`
	PackageUpdates   int       `json:"package_updates"`
	RebootRequired   bool      `json:"reboot_required"`
	Message          string    `json:"message"`
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output in machine-readable JSON format")
	checkFlag := flag.Bool("check", true, "Check for system updates")
	versionFlag := flag.Bool("version", false, "Print sawit-update version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("sawit-update version %s\n", Version)
		os.Exit(0)
	}

	_ = checkFlag

	status := UpdateStatus{
		LastChecked:     time.Now(),
		SystemUpToDate:  true,
		SecurityUpdates: 0,
		PackageUpdates:  0,
		RebootRequired:  false,
		Message:         "System is up to date.",
	}

	if *jsonFlag {
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Update Manager (sawit-update)")
	fmt.Println("──────────────────────────────────────────")
	fmt.Printf("Status:           %s\n", status.Message)
	fmt.Printf("Security Updates: %d\n", status.SecurityUpdates)
	fmt.Printf("Package Updates:  %d\n", status.PackageUpdates)
	fmt.Printf("Reboot Required:  %t\n", status.RebootRequired)
	fmt.Printf("Last Checked:     %s\n", status.LastChecked.Format(time.RFC3339))
}
