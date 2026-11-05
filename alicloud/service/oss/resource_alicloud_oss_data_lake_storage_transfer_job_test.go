package oss_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudOssDataLakeStorageTransferJob exercises the full lifecycle of
// alicloud_oss_data_lake_storage_transfer_job: create with a Rule, update the
// Rule and the enable operational flag, then import. The DataLakeStorageTransferJob
// APIs are OSS data-plane subresource operations reached through
// /?x-oss-datalake-storage-transfer-job; six of the seven CRUD APIs are
// Private, so this case also confirms the inferred Private request encoding
// (job-id query parameter, PUT body wrapper, Start/Stop POST) against the
// real service.
func TestAccAliCloudOssDataLakeStorageTransferJob_basic(t *testing.T) {
	resourceId := "alicloud_oss_data_lake_storage_transfer_job.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket names only accept lowercase letters, digits and hyphens.
	name := fmt.Sprintf("tf-testacc%sossdlst%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssDataLakeStorageTransferJob(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssDataLakeStorageTransferJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssDataLakeStorageTransferJobConfig(name, region, "logs/transfer", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttrSet(resourceId, "job_id"),
					resource.TestCheckResourceAttr(resourceId, "log_base_dir", "logs/transfer"),
					resource.TestCheckResourceAttr(resourceId, "need_verify", "false"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				// Update the Rule payload (reconfigure the job).
				Config: testAccAliCloudOssDataLakeStorageTransferJobConfig(name, region, "logs/transfer/v2", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "log_base_dir", "logs/transfer/v2"),
					resource.TestCheckResourceAttr(resourceId, "need_verify", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable"},
			},
		},
	})
}

// testAccPreCheckOssDataLakeStorageTransferJob pins the region: the
// DataLakeStorageTransferJob plane is served from cn-hangzhou and the
// OSS-HDFS data lake transfer feature must be ticket-enabled for the test
// account. The feature gate is a known release blocker tracked with the OSS
// product owner; absence produces an AccessDenied / feature-not-open error
// rather than a provider defect.
func testAccPreCheckOssDataLakeStorageTransferJob(t *testing.T) {
	acctest.PreCheck(t)
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

// testAccCheckOssDataLakeStorageTransferJobDestroy verifies that every job
// created by the test has been deleted, calling GetDataLakeStorageTransferJob
// through the provider meta exactly as the production Read does.
func testAccCheckOssDataLakeStorageTransferJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_data_lake_storage_transfer_job" {
			continue
		}
		parts, err := helper.ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return fmt.Errorf("invalid resource id %s: expected bucket:job_id", rs.Primary.ID)
		}
		bucket, jobID := parts[0], parts[1]
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		query := map[string]*string{
			"x-oss-datalake-storage-transfer-job-id": helper.StringPointer(jobID),
		}
		hostMap := map[string]*string{"bucket": helper.StringPointer(bucket)}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetDataLakeStorageTransferJob", "/?x-oss-datalake-storage-transfer-job"), query, nil, nil, hostMap, true)
			if requestErr != nil {
				if helper.NeedRetry(requestErr) {
					wait()
					return retry.RetryableError(requestErr)
				}
				return retry.NonRetryableError(requestErr)
			}
			return nil
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NoSuchDataLakeStorageTransferJob", "NoSuchBucket"}) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing data lake storage transfer job %s failed: %s", rs.Primary.ID, err)
		}
		if response != nil {
			return fmt.Errorf("data lake storage transfer job %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// testAccAliCloudOssDataLakeStorageTransferJobConfig builds a self-contained
// stack: a bucket, the RAM executor role OSS-HDFS assumes to run the transfer,
// and the data lake storage transfer job itself. logDir and needVerify vary
// across steps to exercise Rule reconfiguration; enable toggles Start/Stop.
func testAccAliCloudOssDataLakeStorageTransferJobConfig(name, region, logDir string, needVerify bool) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

provider "alicloud" {
  region = "%[2]s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_ram_role" "executor" {
  name        = "${var.name}-role"
  document    = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["oss.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
EOF
  force = true
}

resource "alicloud_oss_data_lake_storage_transfer_job" "default" {
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  executor_role_id  = alicloud_ram_role.executor.arn
  log_base_dir      = "%[3]s"
  need_verify       = %[4]t
  prefix_filter_includes = ["logs/", "data/"]
  enable            = %[5]t
}
`, name, region, logDir, needVerify, needVerify)
}
