package internal

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBuildAllPackages(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = filepath.Join(filepath.Dir(file), "..", "..", "..")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./... failed: %v\n%s", err, output)
	}
}
