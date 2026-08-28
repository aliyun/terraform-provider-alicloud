package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls MetricStore.
func TestAccAliCloudSlsMetricStore_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_metric_store.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsMetricStoreMap)
	serviceFunc := func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeSlsLogStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsms%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsMetricStoreBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SlsTestRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name":          "${alicloud_log_project.default.project_name}",
					"metric_store_name":     name,
					"ttl":                   "30",
					"shard_count":           "2",
					"auto_split":            "true",
					"max_split_shard_count": "64",
					"append_meta":           "false",
					"hot_ttl":               "7",
					"mode":                  "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":          CHECKSET,
						"metric_store_name":     name,
						"ttl":                   "30",
						"shard_count":           "2",
						"auto_split":            "true",
						"max_split_shard_count": "64",
						"append_meta":           "false",
						"hot_ttl":               "7",
						"mode":                  "standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl":                   "180",
					"max_split_shard_count": "128",
					"append_meta":           "true",
					"hot_ttl":               "14",
					"infrequent_access_ttl": "90",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  "${data.alicloud_regions.default.regions.0.id}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl":                   "180",
						"max_split_shard_count": "128",
						"append_meta":           "true",
						"hot_ttl":               "14",
						"infrequent_access_ttl": "90",
						"encrypt_conf.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_split": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_split":            "false",
						"max_split_shard_count": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"encrypt_conf",
				},
			},
		},
	})
}

var AlicloudSlsMetricStoreMap = map[string]string{
	"project_name":          CHECKSET,
	"metric_store_name":     CHECKSET,
	"ttl":                   CHECKSET,
	"shard_count":           CHECKSET,
	"auto_split":            CHECKSET,
	"max_split_shard_count": CHECKSET,
	"append_meta":           CHECKSET,
	"hot_ttl":               CHECKSET,
	"mode":                  CHECKSET,
	"create_time":           CHECKSET,
	"last_modify_time":      CHECKSET,
}

func AlicloudSlsMetricStoreBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_account" "default" {}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "terraform test project for metric store"
}

resource "alicloud_kms_key" "key" {
  description            = "${var.name}-kms"
  pending_window_in_days = "7"
  status                 = "Enabled"
}
`, name)
}
