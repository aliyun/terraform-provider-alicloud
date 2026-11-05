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
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// OSS validates the referenced Function Compute function and the assumed role
// at PutBucketEventNotification time, so the acceptance tests must provision a
// real FC service+function and a RAM role that OSS can assume, alongside the
// target bucket. FC v1 is served from cn-hangzhou in the test account.

func TestAccAliCloudOssBucketEventNotification_basic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_event_notification.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	name := fmt.Sprintf("tf-testacc%sossevn%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssBucketEventNotification(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketEventNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOssBucketEventNotificationBasic(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "id", name),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.events.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.events.0", "oss:ObjectCreated:PutObject"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.filter.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.filter.0.key.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.filter.0.key.0.prefix", "source"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.filter.0.key.0.suffix", ".png"),
					resource.TestCheckResourceAttrSet(resourceId, "function_compute_configuration.0.function_arn"),
					resource.TestCheckResourceAttrSet(resourceId, "function_compute_configuration.0.assume_role"),
					resource.TestCheckResourceAttrSet(resourceId, "function_compute_configuration.0.id"),
				),
			},
			{
				// Whole-replacement update: a second configuration is added and the
				// first one's events and filter change; PutBucketEventNotification
				// replaces the entire NotificationConfiguration.
				Config: testAccOssBucketEventNotificationUpdate(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.events.0", "oss:ObjectCreated:CompleteMultipartUpload"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.0.filter.0.key.0.prefix", "images"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.1.events.0", "oss:ObjectRemoved:Delete"),
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.1.filter.0.key.0.suffix", ".log"),
				),
			},
			{
				// Back to a single configuration, exercising the removal of items.
				Config: testAccOssBucketEventNotificationBasic(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "function_compute_configuration.#", "1"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"function_compute_configuration.0.id",
					"function_compute_configuration.1.id",
				},
			},
		},
	})
}

func testAccPreCheckOssBucketEventNotification(t *testing.T) {
	acctest.PreCheck(t)
	// FC v1 and the OSS event-notification endpoint are exercised from
	// cn-hangzhou in the test account; pin the region instead of skipping.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssBucketEventNotificationDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_event_notification" {
			continue
		}
		hostMap := map[string]*string{"bucket": helper.StringPointer(rs.Primary.ID)}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketEventNotification", "/?eventNotification"), nil, nil, nil, hostMap, true)
			if requestErr != nil {
				if helper.NeedRetry(requestErr) {
					wait()
					return retry.RetryableError(requestErr)
				}
				return retry.NonRetryableError(requestErr)
			}
			return nil
		})
		logs.AddDebug("/?eventNotification", response, rs.Primary.ID)
		if err != nil {
			// A missing configuration or a missing bucket both count as deleted.
			if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchNotificationConfiguration"}) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing bucket event notification %s failed: %s", rs.Primary.ID, err)
		}
		// An empty configuration (no FunctionComputeConfiguration) is the
		// deleted state after DeleteBucketEventNotification.
		if v, ok := response["NotificationConfiguration"].(map[string]interface{}); ok {
			if configs, ok := v["FunctionComputeConfiguration"].([]interface{}); ok && len(configs) > 0 {
				return fmt.Errorf("bucket event notification %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

// testAccOssBucketEventNotificationDependence provisions the cross-service
// prerequisites: the target bucket, an FC v1 service+function and a RAM role
// OSS can assume. The function ARN and role ARN are referenced by the
// event-notification resource.
func testAccOssBucketEventNotificationDependence(name, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_fc_service" "default" {
  name = var.name
}

resource "alicloud_fc_function" "default" {
  service    = alicloud_fc_service.default.name
  name       = var.name
  oss_bucket = alicloud_oss_bucket.default.id
  oss_key    = alicloud_oss_bucket_object.default.key
  runtime    = "python3.10"
  handler    = "hello.handler"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "fc/hello.zip"
  content = <<EOF
# -*- coding: utf-8 -*-
def handler(event, context):
    print("hello world")
    return 'hello world'
EOF
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "oss.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
EOF
  description = "Terraform AccTest"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunFCInvocationAccess"
  policy_type = "System"
}

locals {
  function_arn = "acs:fc:${var.region}:${data.alicloud_caller_identity.current.account_id}:services/${alicloud_fc_service.default.name}.LATEST/functions/${alicloud_fc_function.default.name}"
  assume_role  = alicloud_ram_role.default.arn
}
`, name, region)
}

func testAccOssBucketEventNotificationBasic(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_oss_bucket_event_notification" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  function_compute_configuration {
    events       = ["oss:ObjectCreated:PutObject"]
    function_arn = local.function_arn
    assume_role  = local.assume_role

    filter {
      key {
        prefix = "source"
        suffix = ".png"
      }
    }
  }

  depends_on = [alicloud_ram_role_policy_attachment.default]
}
`, testAccOssBucketEventNotificationDependence(name, region))
}

func testAccOssBucketEventNotificationUpdate(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_oss_bucket_event_notification" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  function_compute_configuration {
    events       = ["oss:ObjectCreated:CompleteMultipartUpload"]
    function_arn = local.function_arn
    assume_role  = local.assume_role

    filter {
      key {
        prefix = "images"
        suffix = ".png"
      }
    }
  }

  function_compute_configuration {
    events       = ["oss:ObjectRemoved:Delete"]
    function_arn = local.function_arn
    assume_role  = local.assume_role

    filter {
      key {
        suffix = ".log"
      }
    }
  }

  depends_on = [alicloud_ram_role_policy_attachment.default]
}
`, testAccOssBucketEventNotificationDependence(name, region))
}
