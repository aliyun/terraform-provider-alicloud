package maxcompute_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudMaxComputeJobDataSource covers the alicloud_max_compute_job
// data source against a real job.
//
// A MaxCompute job is a transient runtime entity submitted through the
// MaxCompute data plane (ODPS SQL). The 2022-01-04 management API exposes only
// read/kill operations on jobs (GetJobInfo, ListJobInfos, GetRunningJobs,
// KillJobs, ...); there is no Create/Update/Delete API for a job (verified at
// triage), so there is no Terraform resource and no control-plane API to
// construct a job fixture in the test. A real job instance ID must therefore
// be supplied out of band.
//
// Set MAX_COMPUTE_JOB_INSTANCE_ID to the instance ID of a job that exists in
// the test region (cn-hangzhou) and is reachable by GetJobInfo. When the
// variable is unset the test is skipped rather than fabricating a fixture; this
// follows the provider's "borrow an existing ID" allowance with the documented
// reason that the entity cannot be self-constructed.
func TestAccAliCloudMaxComputeJobDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	instanceId := os.Getenv("MAX_COMPUTE_JOB_INSTANCE_ID")
	if instanceId == "" {
		t.Skipf("Skipping: MAX_COMPUTE_JOB_INSTANCE_ID is not set. A MaxCompute job has no Create API (verified), so it cannot be self-constructed; provide a real job instance ID from the test region (cn-hangzhou).")
	}

	resourceId := "data.alicloud_max_compute_job.default"
	config := fmt.Sprintf(`
data "alicloud_max_compute_job" "default" {
  instance_id = "%s"
}
`, instanceId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "id", instanceId),
					resource.TestCheckResourceAttr(resourceId, "instance_id", instanceId),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "job_type"),
					resource.TestCheckResourceAttrSet(resourceId, "project"),
					resource.TestCheckResourceAttrSet(resourceId, "region"),
					resource.TestCheckResourceAttrSet(resourceId, "tenant_id"),
				),
			},
		},
	})
}

// TestAccAliCloudMaxComputeJobDataSource_notFound verifies that reading a
// non-existent job instance surfaces an error instead of returning empty
// state. This case is deterministic — a fabricated instance ID never refers to
// a real job — so it runs without MAX_COMPUTE_JOB_INSTANCE_ID and exercises the
// GetJobInfo error path on every supported region.
func TestAccAliCloudMaxComputeJobDataSource_notFound(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	fakeInstanceId := "tf-acc-nonexistent-job-00000000"
	config := fmt.Sprintf(`
data "alicloud_max_compute_job" "default" {
  instance_id = "%s"
}
`, fakeInstanceId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?i)alicloud_max_compute_job`),
			},
		},
	})
}
