package storage

import (
	"testing"
)

func TestGetStorageStatus(t *testing.T) {
	st, err := GetStorageStatus()
	if err != nil {
		t.Fatalf("GetStorageStatus failed: %v", err)
	}

	if st == nil {
		t.Fatal("expected StorageStatus to be non-nil")
	}

	if len(st.Mounts) == 0 {
		t.Error("expected at least root mount")
	}

	disks := GetDisks()
	if len(disks) == 0 {
		t.Error("expected at least mock or physical disks")
	}
}
