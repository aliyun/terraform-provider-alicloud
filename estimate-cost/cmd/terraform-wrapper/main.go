// terraform-wrapper is a transparent wrapper around the real `terraform`
// binary:
//
//	terraform plan --estimate-cost   →  runs the real plan and appends a cost report
//	terraform <any other subcommand> →  syscall.Exec passthrough to the real terraform
//
// The real terraform binary is discovered by scanning $PATH and skipping the
// wrapper's own absolute path — no environment variable configuration needed.
// Cost estimation is delegated to a sibling `estimate-cost` binary; the
// mapping JSON files are embedded inside that binary. For development you
// can set TF_COST_MAPPINGS to a local mapping directory to override the
// embedded copy.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "plan" && hasFlag(args, "--estimate-cost") {
		runPlanWithCost(args)
		return
	}
	passthrough(args)
}

// passthrough exec's the real terraform with the original args, unchanged.
// We use syscall.Exec to replace the current process so signal handling,
// TTY behaviour and exit codes are inherited transparently.
func passthrough(args []string) {
	real := findRealTerraform()
	if real == "" {
		die("terraform-wrapper: real terraform binary not found in PATH")
	}
	if err := syscall.Exec(real, append([]string{real}, args...), os.Environ()); err != nil {
		die("exec %s: %v", real, err)
	}
}

// runPlanWithCost implements the `plan --estimate-cost` branch.
func runPlanWithCost(args []string) {
	real := findRealTerraform()
	if real == "" {
		die("terraform-wrapper: real terraform binary not found in PATH")
	}

	// Strip --estimate-cost and ensure -out= is present (use a temp file when
	// the user did not specify one).
	planArgs := []string{}
	hasOut := false
	outPath := ""
	for _, a := range args {
		if a == "--estimate-cost" {
			continue
		}
		if strings.HasPrefix(a, "-out=") {
			hasOut = true
			outPath = strings.TrimPrefix(a, "-out=")
		}
		planArgs = append(planArgs, a)
	}
	cleanup := func() {}
	// suppressTempPath is non-empty when we injected a temp -out= ourselves.
	// In that case we will filter terraform's "Saved the plan to: <tempPath>"
	// and "terraform apply <tempPath>" hint block out of stdout — that file is
	// about to be deleted, so the hint would mislead the user.
	suppressTempPath := ""
	if !hasOut {
		tmp, err := os.CreateTemp("", "wrap-tfplan-*.tfplan")
		if err != nil {
			die("failed to create temp plan file: %v", err)
		}
		outPath = tmp.Name()
		tmp.Close()
		cleanup = func() { os.Remove(outPath) }
		planArgs = append(planArgs, "-out="+outPath)
		suppressTempPath = outPath
	}
	defer cleanup()

	// Run the real `terraform plan` and stream stdout/stderr to the user.
	planCmd := exec.Command(real, planArgs...)
	planCmd.Stdin = os.Stdin
	planCmd.Stderr = os.Stderr

	// Stdout wiring: pass-through when the user supplied -out=, otherwise
	// pipe through filterApplyHint to rewrite the misleading apply-hint block.
	var pipeW io.Closer
	var filterDone chan struct{}
	if suppressTempPath == "" {
		planCmd.Stdout = os.Stdout
	} else {
		pr, pw := io.Pipe()
		planCmd.Stdout = pw
		pipeW = pw
		filterDone = make(chan struct{})
		go func() {
			defer close(filterDone)
			filterApplyHint(pr, os.Stdout, suppressTempPath)
		}()
	}

	runErr := planCmd.Run()
	if pipeW != nil {
		_ = pipeW.Close() // signal EOF to the filter goroutine
		<-filterDone      // wait for it to drain
	}
	if runErr != nil {
		os.Exit(exitCode(planCmd))
	}

	// Convert the saved plan to JSON via `terraform show -json`.
	showCmd := exec.Command(real, "show", "-json", outPath)
	planJSON, err := showCmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "\n[--estimate-cost] terraform show -json failed:", err)
		return // The plan itself succeeded; do not let a pricing failure derail it.
	}

	// Stash the JSON in a temp file and pass it to the estimate-cost engine.
	tmpJSON, err := os.CreateTemp("", "wrap-plan-*.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "\n[--estimate-cost] writing temp JSON failed:", err)
		return
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Write(planJSON)
	tmpJSON.Close()

	costPath := findEstimateCost()
	if costPath == "" {
		fmt.Fprintln(os.Stderr, "\n[--estimate-cost] estimate-cost binary not found (expected alongside the wrapper)")
		return
	}

	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Println("                       Cost Estimate")
	fmt.Println("─────────────────────────────────────────────────────────────")
	// estimate-cost ships with embedded mappings, so we normally do not need
	// to pass --mappings. For development, the user can set TF_COST_MAPPINGS
	// to a local directory; we forward it to the child.
	costArgs := []string{}
	if dir := os.Getenv("TF_COST_MAPPINGS"); dir != "" {
		costArgs = append(costArgs, "--mappings", dir)
	}
	costArgs = append(costArgs, tmpJSON.Name())
	costCmd := exec.Command(costPath, costArgs...)
	costCmd.Stdout = os.Stdout
	costCmd.Stderr = os.Stderr
	costCmd.Run()
}

// ansiRE matches CSI-style ANSI escape codes (e.g. \x1b[31m, \x1b[1;32m).
// terraform normally disables color when stdout is a pipe, but env vars
// like CLICOLOR_FORCE=1 or FORCE_COLOR=1 can override that. We strip them
// before pattern-matching so the filter is robust to either case.
var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	if !strings.ContainsRune(s, '\x1b') {
		return s
	}
	return ansiRE.ReplaceAllString(s, "")
}

// filterApplyHint reads terraform's stdout from r and writes it to w,
// removing the apply-hint block that points at the temporary plan file we
// injected (which is about to be deleted). The block looks like:
//
//	Saved the plan to: <tempPath>
//
//	To perform exactly these actions, run the following command to apply:
//	    terraform apply "<tempPath>"
//
// All lines from the first match up to and including the `terraform apply
// <tempPath>` line are removed; a friendly hint is printed in their place.
// All other terraform output is forwarded verbatim, line-by-line.
func filterApplyHint(r io.Reader, w io.Writer, tempPath string) {
	scanner := bufio.NewScanner(r)
	// Bump the scanner's max line size from the 64 KiB default to 1 MiB.
	// terraform plan output is usually small, but a single resource diff
	// containing a large blob would otherwise trigger ErrTooLong and we'd
	// drop output silently. 1 MiB is well beyond anything realistic.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	suppressing := false
	inserted := false
	for scanner.Scan() {
		line := scanner.Text()

		// Match against the ANSI-stripped form so coloured terraform output
		// (CLICOLOR_FORCE / FORCE_COLOR) still triggers correctly. We
		// always write the original line back unchanged so any colours
		// passing through are preserved.
		plain := stripANSI(line)

		// Detect the start of the apply-hint block. We trigger on the
		// "Saved the plan to:" prefix alone — terraform wraps long temp
		// paths to the next line when stdout is a pipe, so we cannot rely
		// on tempPath appearing on the same line.
		if !inserted && !suppressing && strings.HasPrefix(plain, "Saved the plan to:") {
			suppressing = true
			continue
		}

		if suppressing {
			// The block ends at the line `    terraform apply "<tempPath>"`.
			// Emit our friendly hint in its place, then resume pass-through.
			if strings.HasPrefix(strings.TrimSpace(plain), "terraform apply") {
				fmt.Fprintln(w, "[--estimate-cost] Plan was processed locally for cost estimation and not saved to disk.")
				fmt.Fprintln(w, "                  To apply, either:")
				fmt.Fprintln(w, "                    · re-run `terraform apply` to plan + apply interactively, or")
				fmt.Fprintln(w, "                    · re-run `terraform plan --estimate-cost -out=<path>` to save")
				fmt.Fprintln(w, "                      the plan, then `terraform apply <path>`.")
				suppressing = false
				inserted = true
			}
			continue
		}

		fmt.Fprintln(w, line)
	}
	// Surface scanner read errors instead of dropping them silently. An
	// over-long line (>1 MiB) would land here as bufio.ErrTooLong.
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "\n[--estimate-cost] warning: stdout read incomplete: %v\n", err)
	}
}

// findRealTerraform scans $PATH for a terraform binary that is not the
// wrapper itself.
func findRealTerraform() string {
	self, _ := os.Executable()
	self, _ = filepath.EvalSymlinks(self)
	for _, dir := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		if dir == "" {
			continue
		}
		c := filepath.Join(dir, "terraform")
		info, err := os.Stat(c)
		if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
			continue
		}
		resolved, _ := filepath.EvalSymlinks(c)
		if resolved != self {
			return c
		}
	}
	return ""
}

// findEstimateCost looks for an `estimate-cost` binary alongside the wrapper
// first, then falls back to scanning $PATH.
func findEstimateCost() string {
	self, err := os.Executable()
	if err == nil {
		cand := filepath.Join(filepath.Dir(self), "estimate-cost")
		if info, err := os.Stat(cand); err == nil && !info.IsDir() {
			return cand
		}
	}
	if p, err := exec.LookPath("estimate-cost"); err == nil {
		return p
	}
	return ""
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func exitCode(cmd *exec.Cmd) int {
	if cmd.ProcessState == nil {
		return 1
	}
	if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		return ws.ExitStatus()
	}
	return cmd.ProcessState.ExitCode()
}
