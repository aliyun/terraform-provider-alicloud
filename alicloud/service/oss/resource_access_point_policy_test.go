package oss_test

import (
	"encoding/json"
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

// Apply-state keeps the configured policy document verbatim (the framework
// requires the applied value to match the plan), while refresh returns the
// API's canonical re-encoding; both are equivalent JSON, but the ID-only check
// compares raw strings. The policy is still verified by the per-step effect
// checks and by ImportStateVerify.

func TestAccAliCloudOssAccessPointPolicy_basic(t *testing.T) {
	resourceId := "alicloud_oss_access_point_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket and access point names only accept lowercase letters, digits
	// and hyphens, and the access point name is capped in length.
	name := fmt.Sprintf("tf-testacc%sossappol%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssAccessPointPolicy(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "access_point_name", name+"-ap"),
					resource.TestCheckResourceAttr(resourceId, "id", fmt.Sprintf("%s:%s-ap", name, name)),
					testAccCheckOssAccessPointPolicyEffect("Allow"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyFormatted(name, region, "Deny"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Deny"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Allow"),
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

func TestAccAliCloudOssAccessPointPolicy_twin(t *testing.T) {
	resourceId := "alicloud_oss_access_point_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket and access point names only accept lowercase letters, digits
	// and hyphens, and the access point name is capped in length.
	name := fmt.Sprintf("tf-testacc%sossappol%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssAccessPointPolicy(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// A create with a non-canonical document: the state must still
				// match the planned value after apply.
				Config: testAccAliCloudOSSAccessPointPolicyFormatted(name, region, "Deny"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "access_point_name", name+"-ap"),
					testAccCheckOssAccessPointPolicyEffect("Deny"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Allow"),
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

func testAccPreCheckOssAccessPointPolicy(t *testing.T) {
	acctest.PreCheck(t)
	// Access points and their policies are served from cn-hangzhou; the
	// flat-suite gate this case was migrated from converged to cn-hangzhou for
	// every configured region, so pin it here instead of skipping the case.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssAccessPointPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_access_point_policy" {
			continue
		}
		parts := strings.SplitN(rs.Primary.ID, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid resource id %s: expected bucket:access_point_name", rs.Primary.ID)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do("Oss", helper.XmlJsonParam("GET", "2019-05-17", "GetAccessPointPolicy", "/?accessPointPolicy"),
				map[string]*string{"x-oss-access-point-name": helper.StringPointer(parts[1])},
				nil, nil, map[string]*string{"bucket": helper.StringPointer(parts[0])}, true)
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
			// A gone bucket (post-test destroy removes the whole stack) or a gone
			// policy both count as deleted, as they did through the production
			// describe's NotFound mapping before.
			if errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicy", "NoSuchAccessPoint"}) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing access point policy %s failed: %s", rs.Primary.ID, err)
		}
		if response != nil {
			return fmt.Errorf("access point policy %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOssAccessPointPolicyEffect(effect string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith("alicloud_oss_access_point_policy.default", "policy", func(value string) error {
		var document map[string]interface{}
		if err := json.Unmarshal([]byte(value), &document); err != nil {
			return fmt.Errorf("policy is not valid JSON: %s", err)
		}
		statements, ok := document["Statement"].([]interface{})
		if !ok {
			return fmt.Errorf("policy has no Statement array: %s", value)
		}
		for _, statement := range statements {
			entry, ok := statement.(map[string]interface{})
			if !ok {
				continue
			}
			if entry["Effect"] == effect {
				return nil
			}
		}
		return fmt.Errorf("policy has no %q statement: %s", effect, value)
	})
}

func testAccAliCloudOSSAccessPointPolicyBasic(name, region, effect string) string {
	return testAccAliCloudOSSAccessPointPolicyConfig(name, region, testAccAliCloudOSSAccessPointPolicyDocument(effect))
}

// testAccAliCloudOSSAccessPointPolicyFormatted keeps the document semantically
// identical to the basic config but re-indents it, which the API does not
// preserve on read.
func testAccAliCloudOSSAccessPointPolicyFormatted(name, region, effect string) string {
	policy := fmt.Sprintf("replace(%s, \",\", \",\\n\")", testAccAliCloudOSSAccessPointPolicyDocument(effect))
	return testAccAliCloudOSSAccessPointPolicyConfig(name, region, policy)
}

func testAccAliCloudOSSAccessPointPolicyDocument(effect string) string {
	return fmt.Sprintf(`jsonencode({
    Version = "1"
    Statement = [
      {
        Action    = ["oss:PutObject", "oss:GetObject"]
        Effect    = "%s"
        Principal = [data.alicloud_caller_identity.current.account_id]
        Resource  = ["acs:oss:${var.region}:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.default.access_point_name}/object/*"]
      }
    ]
  })`, effect)
}

func testAccAliCloudOSSAccessPointPolicyConfig(name, region, policy string) string {
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

resource "alicloud_oss_access_point" "default" {
  access_point_name = "${var.name}-ap"
  bucket            = alicloud_oss_bucket.default.bucket
  # This case exercises the policy, not the access point's origin settings.
  # Pin the blocks the service populates so their server-side values cannot
  # surface as a plan difference in these steps.
  network_origin = "vpc"

  vpc_configuration {
    vpc_id = "vpc-abctest"
  }

  public_access_block_configuration {
    block_public_access = false
  }
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_access_point_policy" "default" {
  access_point_name = alicloud_oss_access_point.default.access_point_name
  bucket            = alicloud_oss_bucket.default.bucket
  policy            = %s
}
`, name, region, policy)
}
