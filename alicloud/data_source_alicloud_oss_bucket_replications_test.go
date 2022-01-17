package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOssBucketsDataSource_replication(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_bucket_replications.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-%d", rand),
		dataSourceOssBucketsConfigDependence_replication)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket_replication.default.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket_replication.default.bucket}-fake",
		}),
	}
	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#":                          "1",
			"names.#":                            "1",
			"buckets.0.name":                     fmt.Sprintf("tf-testacc-bucket-%d-default", rand),
			"buckets.0.acl":                      "private",
			"buckets.0.extranet_endpoint":        CHECKSET,
			"buckets.0.intranet_endpoint":        CHECKSET,
			"buckets.0.location":                 CHECKSET,
			"buckets.0.owner":                    CHECKSET,
			"buckets.0.storage_class":            "Standard",
			"buckets.0.redundancy_type":          "LRS",
			"buckets.0.creation_date":            CHECKSET,
			"buckets.0.cross_region_replication": CHECKSET,
			"buckets.0.transfer_acceleration":    CHECKSET,

			"buckets.0.replication_rule.#":                               "1",
			"buckets.0.replication_rule.0.prefix_set.#":                  "1",
			"buckets.0.replication_rule.0.prefix_set.0.prefixes.#":       "2",
			"buckets.0.replication_rule.0.prefix_set.0.prefixes.0":       "xx/",
			"buckets.0.replication_rule.0.prefix_set.0.prefixes.1":       "test/",
			"buckets.0.replication_rule.0.destination.#":                 "1",
			"buckets.0.replication_rule.0.destination.0.bucket":          fmt.Sprintf("tf-testacc-bucket-%d-target", rand),
			"buckets.0.replication_rule.0.destination.0.location":        "oss-cn-beijing",
			"buckets.0.replication_rule.0.destination.0.transfer_type":   "oss_acc",
			"buckets.0.replication_rule.0.action":                        "PUT",
			"buckets.0.replication_rule.0.historical_object_replication": "enabled",
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#": "0",
			"names.#":   "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
	}

	ossBucketsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceOssBucketsConfigDependence_replication(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket_replication" "target"{
	bucket = "${var.name}-target"
}

resource "alicloud_oss_bucket_replication" "default" {
	bucket = "${var.name}-default"

   replication_rule {
       prefix_set {
           prefixes = ["xx/", "test/"]
       }
       destination {
           bucket = "${alicloud_oss_bucket_replication.target.id}"
           location = "oss-cn-beijing"
           transfer_type = "oss_acc"
       }
       action = "PUT"
       historical_object_replication = "enabled"
   }
}

`, name)
}
