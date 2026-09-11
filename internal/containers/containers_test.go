package containers_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/containers"
)

func TestListContainers(t *testing.T) {
	list, err := containers.ListContainers()
	if err != nil {
		t.Fatalf("unexpected error listing containers: %v", err)
	}
	if len(list) == 0 {
		t.Log("No active containers found or Docker engine offline")
	}
}

func TestListImages(t *testing.T) {
	images, err := containers.ListImages()
	if err != nil {
		t.Fatalf("unexpected error listing images: %v", err)
	}
	if len(images) == 0 {
		t.Log("No container images found")
	}
}

func TestEmptyContainerActions(t *testing.T) {
	if err := containers.StartContainer(""); err == nil {
		t.Error("expected error for empty container ID")
	}
	if err := containers.StopContainer(""); err == nil {
		t.Error("expected error for empty container ID")
	}
}
