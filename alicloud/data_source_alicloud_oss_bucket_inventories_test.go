package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudOssBucketInventoriesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_oss_bucket_inventories.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfaccossinv%d", rand),
		dataSourceOssBucketInventoriesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket": "${alicloud_oss_bucket_inventory.default.bucket}",
			"ids":    []string{"${alicloud_oss_bucket_inventory.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket": "${alicloud_oss_bucket_inventory.default.bucket}",
			"ids":    []string{"${alicloud_oss_bucket_inventory.default.id}-fake"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                "1",
			"ids.0":                                                CHECKSET,
			"inventories.#":                                        "1",
			"inventories.0.id":                                     CHECKSET,
			"inventories.0.inventory_id":                           fmt.Sprintf("tfaccossinv%d", rand),
			"inventories.0.is_enabled":                             "true",
			"inventories.0.included_object_versions":               "All",
			"inventories.0.schedule.#":                             "1",
			"inventories.0.schedule.0.frequency":                   "Daily",
			"inventories.0.optional_fields.#":                      "1",
			"inventories.0.optional_fields.0.field.#":              "2",
			"inventories.0.destination.#":                          "1",
			"inventories.0.destination.0.oss_bucket_destination.#": "1",
			"inventories.0.destination.0.oss_bucket_destination.0.format": "CSV",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"inventories.#": "0",
		}
	}

	var checkInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	checkInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceOssBucketInventoriesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_oss_bucket" "source" {
  bucket        = "${var.name}-src"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "dest" {
  bucket        = "${var.name}-dst"
  storage_class = "Standard"
}

resource "alicloud_ram_role" "default" {
  name        = "${var.name}"
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
  description = "this is a test for bucket inventory."
  force       = true
}

resource "alicloud_ram_policy" "default" {
  name     = "${var.name}"
  document = <<EOF
{
  "Statement": [
    {
      "Action": ["oss:PutObject", "oss:AbortMultipartUpload"],
      "Effect": "Allow",
      "Resource": ["acs:oss:*:*:${alicloud_oss_bucket.dest.id}/*"]
    }
  ],
  "Version": "1"
}
EOF
  force    = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  policy_name = "${alicloud_ram_policy.default.name}"
  policy_type = "${alicloud_ram_policy.default.type}"
  role_name   = "${alicloud_ram_role.default.name}"
}

resource "alicloud_oss_bucket_inventory" "default" {
  bucket       = "${alicloud_oss_bucket.source.id}"
  inventory_id = "${var.name}"
  is_enabled   = true

  included_object_versions = "All"

  schedule {
    frequency = "Daily"
  }

  optional_fields {
    field = ["Size", "LastModifiedDate"]
  }

  destination {
    oss_bucket_destination {
      format     = "CSV"
      account_id = "${data.alicloud_account.current.id}"
      role_arn   = "acs:ram::${data.alicloud_account.current.id}:role/${alicloud_ram_role_policy_attachment.default.role_name}"
      bucket     = "acs:oss:::${alicloud_oss_bucket.dest.id}"
      prefix     = "inv/"
    }
  }
}

`, name)
}
