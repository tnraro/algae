package util

import (
	"os"
	"testing"
)

func SetupDataDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "algae_*")
	if err != nil {
		t.Fatalf("Failed to setup DATA_DIR: %v", err)
	}
	defer os.RemoveAll(dir)
	os.Setenv("DATA_DIR", dir)
	return dir
}

func AssertEq[T comparable](t *testing.T, a T, b T) {
	if a != b {
		t.Fatalf("%v is not %v", a, b)
	}
}
