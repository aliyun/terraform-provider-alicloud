// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketInventory. >>> Resource test cases, automatically generated.
// Case BucketInventory测试 6645
func TestAccAliCloudOssBucketInventory_basic6645(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_inventory.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketInventoryMap6645)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketInventory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketInventoryBasicDependence6645)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.this.id}",
									"role_arn":   "${alicloud_ram_role.role.arn}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.DestBucket.id}",
									"prefix":     "Pics/",
								},
							},
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{
								"Size", "LastModifiedDate", "ETag"},
						},
					},
					"bucket": "${alicloud_oss_bucket.CreateBucket.id}",
					"filter": []map[string]interface{}{
						{
							"prefix":           "Pics/",
							"lower_size_bound": "256",
							"upper_size_bound": "999999",
							"storage_class":    "Standard",
						},
					},
					"included_object_versions": "Current",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Daily",
						},
					},
					"inventory_id": "report01",
					"is_enabled":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                   CHECKSET,
						"included_object_versions": "Current",
						"inventory_id":             "report01",
						"is_enabled":               "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.this.id}",
									"role_arn":   "${alicloud_ram_role.role.arn}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.DestBucket.id}",
									"prefix":     "Pics/",
								},
							},
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{
								"Size", "LastModifiedDate", "ETag"},
						},
					},
					"bucket": "${alicloud_oss_bucket.CreateBucket.id}",
					"filter": []map[string]interface{}{
						{
							"prefix":           "Pics/",
							"lower_size_bound": "256",
							"upper_size_bound": "999999",
							"storage_class":    "Standard",
						},
					},
					"included_object_versions": "Current",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Daily",
						},
					},
					"inventory_id": "report01",
					"is_enabled":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
			},
		},
	})
}

var AlicloudOssBucketInventoryMap6645 = map[string]string{}

func AlicloudOssBucketInventoryBasicDependence6645(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-src"
}

resource "alicloud_oss_bucket" "DestBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-dst"
}

resource "alicloud_ram_role" "role" {
  name     = "${var.name}"
  document = <<EOF
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
}
`, name)
}

// Test Oss BucketInventory. <<< Resource test cases, automatically generated.
