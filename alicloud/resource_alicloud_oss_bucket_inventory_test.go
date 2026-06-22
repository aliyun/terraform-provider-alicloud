package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketInventory. >>> Resource test cases, automatically generated.
// Case BucketInventory basic
func TestAccAliCloudOssBucketInventory_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_inventory.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketInventoryMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketInventory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketInventoryBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                   "${alicloud_oss_bucket.CreateBucket.id}",
					"inventory_id":             "report1",
					"is_enabled":               "true",
					"included_object_versions": "All",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Daily",
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{"Size", "LastModifiedDate"},
						},
					},
					"filter": []map[string]interface{}{
						{
							"prefix": "frontends/",
						},
					},
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.this.id}",
									"role_arn":   "${alicloud_ram_role.role.arn}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.DestBucket.id}",
									"prefix":     "inventory/",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                   CHECKSET,
						"inventory_id":             "report1",
						"is_enabled":               "true",
						"included_object_versions": "All",
						"schedule.#":               "1",
						"schedule.0.frequency":     "Daily",
						"destination.#":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_enabled":               "false",
					"included_object_versions": "Current",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Weekly",
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{"Size", "StorageClass", "ETag", "EncryptionStatus"},
						},
					},
					"filter": []map[string]interface{}{
						{
							"prefix": "logs/",
						},
					},
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.this.id}",
									"role_arn":   "${alicloud_ram_role.role.arn}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.DestBucket.id}",
									"prefix":     "report/",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled":               "false",
						"included_object_versions": "Current",
						"schedule.0.frequency":     "Weekly",
						"filter.0.prefix":          "logs/",
						"destination.0.oss_bucket_destination.0.prefix": "report/",
					}),
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

var AlicloudOssBucketInventoryMap = map[string]string{}

func AlicloudOssBucketInventoryBasicDependence(name string) string {
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
