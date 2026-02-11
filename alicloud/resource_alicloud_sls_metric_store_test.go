package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSlsMetricStore_basic7939(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_metric_store.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsMetricStoreMap7939)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsMetricStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsMetricStoreBasicDependence7939)
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
					"project_name":      "${alicloud_log_project.defaultj3iK32.project_name}",
					"metering_mode":     "ChargeByFunction",
					"mode":              "standard",
					"metric_type":       "prometheus",
					"metric_store_name": name,
					"ttl":               "7",
					"shard_count":       "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":      "sls-sdk-testp-metricstore",
						"metering_mode":     "ChargeByFunction",
						"mode":              "standard",
						"metric_type":       "prometheus",
						"metric_store_name": name,
						"ttl":               "7",
						"shard_count":       "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metering_mode"},
			},
		},
	})
}

var AlicloudSlsMetricStoreMap7939 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudSlsMetricStoreBasicDependence7939(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultj3iK32" {
  description  = "test"
  project_name = "sls-sdk-testp-metricstore"
}


`, name)
}
