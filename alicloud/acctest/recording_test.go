package acctest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Exercise the actual test-binary startup without a product TestMain, including
// randomness generated BEFORE PreCheck and requests made during test cleanup.
func TestSharedRecordingProcess(t *testing.T) {
	if os.Getenv("TF_RR_UNIT_CHILD") != "" {
		name := sdkacctest.RandString(20)
		token := helper.BuildClientToken("Create")
		scenario := os.Getenv("TF_RR_UNIT_CHILD")
		call := func() {
			r, _ := http.NewRequest("GET", os.Getenv("TF_RR_UNIT_URL")+"?Name="+name+"&ClientToken="+token, nil)
			resp, err := (&http.Client{Transport: recording.Transport(nil)}).Do(r)
			if err == nil {
				resp.Body.Close()
				if scenario == "unbound" {
					t.Fatal("unbound request reached HTTP")
				}
			} else if scenario == "unbound" {
				if !strings.Contains(err.Error(), "acctest.PreCheck") {
					t.Fatal(err)
				}
			} else if scenario != "swallow" && scenario != "unbound" {
				t.Fatal(err)
			}
		}
		if scenario == "unbound" {
			call()
			return
		}
		if scenario == "subtest" {
			t.Run("nested", func(t *testing.T) { PreCheck(t) })
			return
		}
		if scenario == "early-cleanup" {
			t.Cleanup(call)
		}
		PreCheck(t)
		PreCheck(t) // Must register finalization only once.
		call()
		if scenario != "unused" {
			t.Cleanup(call)
		}
		if scenario == "skip" {
			t.Skip("synthetic skip after HTTP")
		}
		if scenario == "fail" {
			t.Fatal("synthetic failure after HTTP")
		}
		return
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"State":"Ready"}`)
	}))
	url := server.URL
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	run := func(mode, dir, scenario string, extra ...string) (bool, string) {
		t.Helper()
		cmd := exec.Command(exe, append([]string{"-test.run=^TestSharedRecordingProcess$", "-test.v"}, extra...)...)
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "TF_ACC") && !strings.HasPrefix(e, "TF_RR_UNIT_") {
				cmd.Env = append(cmd.Env, e)
			}
		}
		cmd.Env = append(cmd.Env, "TF_ACC=1", "ALICLOUD_ACCESS_KEY=fake", "ALICLOUD_SECRET_KEY=fake", "ALICLOUD_REGION=cn-hangzhou", "TF_RR_UNIT_CHILD="+scenario, "TF_RR_UNIT_URL="+url, recording.ModeEnv+"="+mode, recording.PathEnv+"="+dir, "TF_ACC_RECORDING_RESULT_PATH="+filepath.Join(t.TempDir(), "result.json"))
		data, err := cmd.CombinedOutput()
		var result recording.Result
		for _, e := range cmd.Env {
			if strings.HasPrefix(e, "TF_ACC_RECORDING_RESULT_PATH=") {
				b, _ := os.ReadFile(strings.TrimPrefix(e, "TF_ACC_RECORDING_RESULT_PATH="))
				_ = json.Unmarshal(b, &result)
			}
		}
		if result.Success && (err != nil || result.InteractionCount != 2) {
			t.Fatalf("invalid success result: %+v\n%s", result, data)
		}
		if scenario == "ok" && len(extra) == 0 && (err != nil || !result.Success) {
			t.Fatalf("missing successful session: %+v\n%s", result, data)
		}
		return err == nil && result.Success, string(data)
	}
	dir := filepath.Join(t.TempDir(), "baseline")
	run("record", dir, "ok")
	before, _ := os.ReadFile(filepath.Join(dir, "cassette.json"))
	seed, _ := os.ReadFile(filepath.Join(dir, "seed"))
	for _, scenario := range []string{"skip", "fail"} {
		failedDir := filepath.Join(t.TempDir(), scenario)
		valid, out := run("record", failedDir, scenario)
		if valid || !strings.Contains(out, "FAIL") {
			t.Fatalf("failed recording passed: %s", out)
		}
		if _, err := os.Stat(filepath.Join(failedDir, "cassette.json")); !os.IsNotExist(err) {
			t.Fatal("failed recording published a cassette")
		}
	}
	server.Close() // Successful replay must not reach even a local HTTP server.
	run("replay", dir, "ok")
	for _, scenario := range []string{"unused", "swallow", "unbound", "skip", "fail", "subtest", "early-cleanup"} {
		t.Run(scenario, func(t *testing.T) {
			// A changed URL causes a swallowed mismatch; all other variants retain URL.
			previous := url
			if scenario == "swallow" {
				url += "/changed"
			}
			defer func() { url = previous }()
			valid, out := run("replay", dir, scenario)
			if valid {
				t.Fatalf("invalid case passed: %s", out)
			}
			if scenario != "unbound" && scenario != "early-cleanup" && !strings.Contains(out, "FAIL") {
				t.Fatalf("invalid lifecycle did not fail test: %s", out)
			}
		})
	}
	for _, args := range [][]string{{"-test.count=2"}, {"-test.shuffle=on"}, {"-sweep=cn-hangzhou"}, {"-test.run=Test"}} {
		valid, out := run("replay", dir, "ok", args...)
		if valid || !strings.Contains(out, "recording initialization failed") {
			t.Fatalf("invalid args accepted: %s", out)
		}
	}
	after, _ := os.ReadFile(filepath.Join(dir, "cassette.json"))
	afterSeed, _ := os.ReadFile(filepath.Join(dir, "seed"))
	if string(before) != string(after) || string(seed) != string(afterSeed) {
		t.Fatal("replay modified baseline")
	}
}
