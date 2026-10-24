package alicloud

// The alicloud_polardb_cluster sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 and plugin-testing
// helper/resource packages cannot both register -sweep in one binary. The
// resource implementation itself lives in alicloud/service/polardb.

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	polardbsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_cluster", &resource.Sweeper{
		Name: "alicloud_polardb_cluster",
		F:    testSweepPolarDBClusters,
	})
}

func testSweepPolarDBClusters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []polardb.DBCluster
	req := polardb.CreateDescribeDBClustersRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusters(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Polardb Instances: %s", err)
		}
		resp, _ := raw.(*polardb.DescribeDBClustersResponse)
		if resp == nil || len(resp.Items.DBCluster) < 1 {
			break
		}
		insts = append(insts, resp.Items.DBCluster...)

		if len(resp.Items.DBCluster) < PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.DBClusterDescription
		id := v.DBClusterId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			if skip {
				if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
					skip = !need
				}
			}

			if skip {
				log.Printf("[INFO] Skipping Polardb Instance: %s (%s)", name, id)
				continue
			}
		}

		log.Printf("[INFO] Deleting Polardb Instance: %s (%s)", name, id)

		if strings.EqualFold(v.PayType, string(polardbsvc.Prepaid)) {
			// A subscription cluster is released through the billing refund API;
			// DeleteDBCluster only accepts pay-as-you-go clusters.
			service := polardbsvc.NewService(client)
			if _, err := service.RefundPrePaidCluster(id, 10*time.Minute); err != nil {
				log.Printf("[ERROR] Failed to refund Polardb subscription cluster (%s (%s)): %s", name, id, err)
				continue
			}
			stateConf := helper.BuildStateConf([]string{"Creating", "Running", "Deleting"}, []string{}, 10*time.Minute, 30*time.Second, service.ClusterStateRefreshFunc(id, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				log.Printf("[ERROR] Failed to wait for Polardb subscription cluster release (%s (%s)): %s", name, id, err)
			}
			continue
		}

		req := polardb.CreateDeleteDBClusterRequest()
		req.DBClusterId = id
		_, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DeleteDBCluster(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Polardb Instance (%s (%s)): %s", name, id, err)
		}
	}

	return nil
}

// TestUnitPolarDBClusterSweeperEntry verifies the cleanup entry without
// deleting any cloud resource: make sweep must load the root package, and the
// root binary's sweeper registry must select alicloud_polardb_cluster and fail
// on missing credentials before any API call.
func TestUnitPolarDBClusterSweeperEntry(t *testing.T) {
	t.Run("make-sweep-load-path", func(t *testing.T) {
		makePath, err := exec.LookPath("make")
		if err != nil {
			t.Fatal("make required for sweeper entry test:", err)
		}
		source, err := os.ReadFile("../GNUmakefile")
		if err != nil {
			t.Fatal(err)
		}
		temp := t.TempDir()
		bin := filepath.Join(temp, "bin")
		if err := os.Mkdir(bin, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(temp, "GNUmakefile"), source, 0600); err != nil {
			t.Fatal(err)
		}
		shim := `#!/bin/sh
set -eu
if [ "$1" != test ]; then echo 'unexpected fake go command' >&2; exit 81; fi
printf '%s\n' "$@" > "$SWEEP_ARGS_FILE"
printf '%s\n' "${TF_ACC:-}" > "$SWEEP_ENV_FILE"
`
		if err := os.WriteFile(filepath.Join(bin, "go"), []byte(shim), 0700); err != nil {
			t.Fatal(err)
		}
		argsFile := filepath.Join(temp, "args")
		envFile := filepath.Join(temp, "env")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, makePath, "--no-print-directory", "-f", "GNUmakefile", "sweep", "REGION=cn-hangzhou", "RESOURCE=alicloud_polardb_cluster")
		cmd.Dir = temp
		// Explicit minimal environment: never inherit real credentials, MAKEFLAGS
		// or a real Go executable ahead of the shim.
		cmd.Env = []string{"PATH=" + bin + string(os.PathListSeparator) + os.Getenv("PATH"), "HOME=" + temp, "SWEEP_ARGS_FILE=" + argsFile, "SWEEP_ENV_FILE=" + envFile}
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("make sweep failed: %v: %s", err, output)
		}
		raw, err := os.ReadFile(argsFile)
		if err != nil {
			t.Fatal(err)
		}
		got := strings.Split(strings.TrimSpace(string(raw)), "\n")
		want := []string{"test", "./alicloud", "-v", "-timeout=10m", "-sweep=cn-hangzhou", "-sweep-run=alicloud_polardb_cluster"}
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Fatalf("make sweep loads %q; want root package %q", got, want)
		}
		env, err := os.ReadFile(envFile)
		if err != nil {
			t.Fatal(err)
		}
		if string(env) != "1\n" {
			t.Fatalf("unexpected sweep env %q", env)
		}
	})

	probe := func(t *testing.T, sweepRun string) (string, error) {
		t.Helper()
		exe, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, exe, "-sweep=cn-hangzhou", "-sweep-run="+sweepRun)
		// No cloud credentials: the sweeper must fail in sharedClientForRegion
		// before issuing any API call, so this probe never touches the cloud.
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir()}
		output, err := cmd.CombinedOutput()
		return string(output), err
	}

	t.Run("registered-in-make-sweep-runner", func(t *testing.T) {
		output, err := probe(t, "alicloud_polardb_cluster")
		if err == nil {
			t.Fatalf("sweeper runner probe unexpectedly succeeded; output: %s", output)
		}
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 1 {
			t.Fatalf("sweeper runner probe returned %v; output: %s", err, output)
		}
		for _, want := range []string{
			"Running Sweeper (alicloud_polardb_cluster) in region (cn-hangzhou)",
			"empty ALICLOUD_ACCESS_KEY",
		} {
			if !strings.Contains(output, want) {
				t.Fatalf("sweeper runner output missing %q; output: %s", want, output)
			}
		}
	})

	t.Run("unregistered-name-is-not-selected", func(t *testing.T) {
		output, err := probe(t, "alicloud_polardb_cluster_probe_unregistered")
		if err != nil {
			t.Fatalf("unknown sweeper filter must run nothing and exit 0: %v: %s", err, output)
		}
		if strings.Contains(output, "Running Sweeper (alicloud_polardb_cluster)") {
			t.Fatalf("unknown sweeper filter still ran the cluster sweeper: %s", output)
		}
	})

	t.Run("single-registration", func(t *testing.T) {
		source, err := os.ReadFile("service/polardb/resource_cluster_test.go")
		if err != nil {
			t.Fatal(err)
		}
		for _, marker := range []string{
			`AddTestSweepers("alicloud_polardb_cluster"`,
			"func testSweepPolarDBClusters(",
		} {
			if strings.Contains(string(source), marker) {
				t.Fatalf("product package test file still owns the sweeper marker %q; keep the entry only in the root test binary", marker)
			}
		}
	})
}
