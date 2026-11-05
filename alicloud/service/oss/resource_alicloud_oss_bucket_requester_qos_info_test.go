package oss_test

import (
	"fmt"
	"os"
	"strings"
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

// The requester QoS configuration lives on the bucket and is addressed by the
// requester's RAM user id (composite resource id <bucket>:<user_id>). The ACC
// therefore builds a bucket and a RAM user, then exercises create/read/update/
// delete/import on the QoS info attached to that user.

func TestAccAliCloudOssBucketRequesterQoSInfo_basic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_requester_qos_info.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	name := fmt.Sprintf("tf-testacc%sossqos%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssBucketRequesterQoSInfo(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketRequesterQoSInfoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOSSBucketRequesterQoSInfoConfig(name, region, 1048576, 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "extranet_download_bandwidth", "1048576"),
					resource.TestCheckResourceAttr(resourceId, "extranet_qps", "100"),
					resource.TestCheckResourceAttrPair(resourceId, "bucket", "alicloud_oss_bucket.default", "bucket"),
					resource.TestCheckResourceAttrPair(resourceId, "qos_requester", "alicloud_ram_user.default", "id"),
				),
			},
			{
				// In-place update of the quota fields exercises the upsert Put path.
				Config: testAccAliCloudOSSBucketRequesterQoSInfoConfig(name, region, 2097152, 200),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "extranet_download_bandwidth", "2097152"),
					resource.TestCheckResourceAttr(resourceId, "extranet_qps", "200"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudOssBucketRequesterQoSInfo_twin(t *testing.T) {
	resourceId := "alicloud_oss_bucket_requester_qos_info.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	name := fmt.Sprintf("tf-testacc%sossqos2%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssBucketRequesterQoSInfo(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketRequesterQoSInfoDestroy,
		Steps: []resource.TestStep{
			{
				// Create directly with the updated quotas to exercise the
				// create-then-refresh path with a different baseline.
				Config: testAccAliCloudOSSBucketRequesterQoSInfoConfig(name, region, 3145728, 300),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "extranet_download_bandwidth", "3145728"),
					resource.TestCheckResourceAttr(resourceId, "extranet_qps", "300"),
					resource.TestCheckResourceAttrPair(resourceId, "qos_requester", "alicloud_ram_user.default", "id"),
				),
			},
			{
				Config: testAccAliCloudOSSBucketRequesterQoSInfoConfig(name, region, 1048576, 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "extranet_download_bandwidth", "1048576"),
					resource.TestCheckResourceAttr(resourceId, "extranet_qps", "100"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccPreCheckOssBucketRequesterQoSInfo(t *testing.T) {
	acctest.PreCheck(t)
	// Requester QoS info is served from cn-hangzhou; pin the region so the
	// suite is deterministic regardless of the configured default.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssBucketRequesterQoSInfoDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_requester_qos_info" {
			continue
		}
		parts := strings.SplitN(rs.Primary.ID, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid resource id %s: expected bucket:qos_requester", rs.Primary.ID)
		}
		bucket, qosRequester := parts[0], parts[1]
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do("Oss", helper.XmlJsonParam("GET", "2019-05-17", "GetBucketRequesterQoSInfo", "/?requesterQosInfo"),
				map[string]*string{"qosRequester": helper.StringPointer(qosRequester)},
				nil, nil, map[string]*string{"bucket": helper.StringPointer(bucket)}, true)
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
			// A gone requester QoS config or a gone bucket both count as deleted,
			// mirroring the production describe's NotFound mapping.
			if errs.IsExpectedErrors(err, []string{"NoSuchRequesterQoSInfo", "NoSuchBucket", "NoSuchRequesterQosConfiguration"}) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing requester qos info %s failed: %s", rs.Primary.ID, err)
		}
		if response != nil {
			return fmt.Errorf("requester qos info %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudOSSBucketRequesterQoSInfoConfig(name, region string, bandwidth, qps int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_ram_user" "default" {
  name         = var.name
  display_name = var.name
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_bucket_requester_qos_info" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  # The requester is the RAM user the QoS limits apply to; its resource id is
  # the UserId that the OSS API addresses as qosRequester.
  qos_requester = alicloud_ram_user.default.id

  extranet_download_bandwidth  = %d
  extranet_qps                 = %d
  extranet_upload_bandwidth    = %d
  intranet_download_bandwidth  = %d
  intranet_qps                 = %d
  intranet_upload_bandwidth    = %d
  total_download_bandwidth     = %d
  total_qps                    = %d
  total_upload_bandwidth       = %d
}
`, name, region, bandwidth, qps, bandwidth, bandwidth, qps, bandwidth, bandwidth, qps, bandwidth)
}
