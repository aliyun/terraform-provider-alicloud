package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Sls MetricStore.
func TestAccAliCloudSlsMetricStore_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_metric_store.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsMetricStoreMap)
	serviceFunc := func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeSlsLogStore")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsms%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsMetricStoreBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SlsTestRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"project_name":          acctest.CHECKSET,
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
						"max_split_shard_count": acctest.REMOVEKEY,
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
	"project_name":          acctest.CHECKSET,
	"metric_store_name":     acctest.CHECKSET,
	"ttl":                   acctest.CHECKSET,
	"shard_count":           acctest.CHECKSET,
	"auto_split":            acctest.CHECKSET,
	"max_split_shard_count": acctest.CHECKSET,
	"append_meta":           acctest.CHECKSET,
	"hot_ttl":               acctest.CHECKSET,
	"mode":                  acctest.CHECKSET,
	"create_time":           acctest.CHECKSET,
	"last_modify_time":      acctest.CHECKSET,
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
