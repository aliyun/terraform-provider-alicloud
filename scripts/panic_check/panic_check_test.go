//nolint:all
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPanicCheckDetectsPanic(t *testing.T) {
	// Create a temporary directory for the test file
	tmpDir := t.TempDir()

	// Write a Go file that contains a panic() call
	panicFile := filepath.Join(tmpDir, "has_panic.go")
	panicCode := `package main

import "fmt"

func doSomething() {
	fmt.Println("hello")
	panic("something went wrong")
}

func main() {
	doSomething()
}
`
	if err := os.WriteFile(panicFile, []byte(panicCode), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Build the panic checker binary
	checkerBin := filepath.Join(tmpDir, "panic_check")
	buildCmd := exec.Command("go", "build", "-o", checkerBin, ".")
	buildCmd.Dir = "." // build from the panic_check directory
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build panic checker: %v\n%s", err, out)
	}

	// Run the checker against the test file
	runCmd := exec.Command(checkerBin, "-fileNames="+panicFile)
	output, err := runCmd.CombinedOutput()

	// Should exit with code 1 (panics found)
	if err == nil {
		t.Fatalf("expected panic checker to exit non-zero, but it exited 0.\nOutput: %s", output)
	}

	outStr := string(output)

	// Should report the panic in doSomething
	if !strings.Contains(outStr, "panic() call in function doSomething") {
		t.Errorf("expected output to contain 'panic() call in function doSomething', got:\n%s", outStr)
	}

	// Should report the file path
	if !strings.Contains(outStr, "has_panic.go") {
		t.Errorf("expected output to contain 'has_panic.go', got:\n%s", outStr)
	}

	// Should report FAIL
	if !strings.Contains(outStr, "FAIL") {
		t.Errorf("expected output to contain 'FAIL', got:\n%s", outStr)
	}
}

func TestPanicCheckCleanFile(t *testing.T) {
	// Create a temporary directory for the test file
	tmpDir := t.TempDir()

	// Write a Go file that does NOT contain a panic() call
	cleanFile := filepath.Join(tmpDir, "clean.go")
	cleanCode := `package main

import (
	"fmt"
	"log"
)

func doSomething() {
	fmt.Println("hello")
	log.Println("safe operation")
}

func main() {
	doSomething()
}
`
	if err := os.WriteFile(cleanFile, []byte(cleanCode), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Build the panic checker binary
	checkerBin := filepath.Join(tmpDir, "panic_check")
	buildCmd := exec.Command("go", "build", "-o", checkerBin, ".")
	buildCmd.Dir = "." // build from the panic_check directory
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build panic checker: %v\n%s", err, out)
	}

	// Run the checker against the clean file
	runCmd := exec.Command(checkerBin, "-fileNames="+cleanFile)
	output, err := runCmd.CombinedOutput()

	// Should exit with code 0 (no panics)
	if err != nil {
		t.Fatalf("expected panic checker to exit 0, but it failed.\nOutput: %s", output)
	}

	outStr := string(output)

	// Should report OK
	if !strings.Contains(outStr, "OK") {
		t.Errorf("expected output to contain 'OK', got:\n%s", outStr)
	}

	// Should report 0 panic calls
	if !strings.Contains(outStr, "0 panic calls found") {
		t.Errorf("expected output to contain '0 panic calls found', got:\n%s", outStr)
	}
}

func TestPanicCheckSkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a test file that contains a panic() call
	testFile := filepath.Join(tmpDir, "has_panic_test.go")
	testCode := `package main

func TestSomething(t *testing.T) {
	panic("this should be skipped")
}
`
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Write file list
	fileList := filepath.Join(tmpDir, "files.txt")
	if err := os.WriteFile(fileList, []byte(testFile+"\n"), 0644); err != nil {
		t.Fatalf("failed to write file list: %v", err)
	}

	// Build the panic checker binary
	checkerBin := filepath.Join(tmpDir, "panic_check")
	buildCmd := exec.Command("go", "build", "-o", checkerBin, ".")
	buildCmd.Dir = "."
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build panic checker: %v\n%s", err, out)
	}

	// Run without -includeTests — should skip the test file and exit 0
	runCmd := exec.Command(checkerBin, "-fileNames="+fileList)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected panic checker to exit 0 (test files skipped), but it failed.\nOutput: %s", output)
	}

	if !strings.Contains(string(output), "OK") {
		t.Errorf("expected OK output when test files are skipped, got:\n%s", output)
	}
}

func TestPanicCheckExclude(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a Go file with panic
	panicFile := filepath.Join(tmpDir, "common.go")
	panicCode := `package main

func helper() {
	panic("should be excluded")
}
`
	if err := os.WriteFile(panicFile, []byte(panicCode), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Write file list
	fileList := filepath.Join(tmpDir, "files.txt")
	if err := os.WriteFile(fileList, []byte(panicFile+"\n"), 0644); err != nil {
		t.Fatalf("failed to write file list: %v", err)
	}

	// Build the panic checker binary
	checkerBin := filepath.Join(tmpDir, "panic_check")
	buildCmd := exec.Command("go", "build", "-o", checkerBin, ".")
	buildCmd.Dir = "."
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build panic checker: %v\n%s", err, out)
	}

	// Run with -exclude common.go — should skip and exit 0
	runCmd := exec.Command(checkerBin, "-fileNames="+fileList, "-exclude=common.go")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected panic checker to exit 0 (excluded), but it failed.\nOutput: %s", output)
	}

	if !strings.Contains(string(output), "OK") {
		t.Errorf("expected OK output when file is excluded, got:\n%s", output)
	}
}
