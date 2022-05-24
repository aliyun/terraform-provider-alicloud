package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEciImageCachesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacceci-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_eci_image_caches.default", name, dataSourceEciImageCachesConfigDependence)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_eci_image_cache.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_eci_image_cache.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_eci_image_cache.default.id}"},
			"name_regex": "${alicloud_eci_image_cache.default.image_cache_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_eci_image_cache.default.id}"},
			"name_regex": "${alicloud_eci_image_cache.default.image_cache_name}-fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_eci_image_cache.default.id}"},
			"status": "${alicloud_eci_image_cache.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_eci_image_cache.default.id}"},
			"status": "Failed",
		}),
	}

	var existEciImageCachesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"caches.#":                    "1",
			"caches.0.id":                 CHECKSET,
			"caches.0.container_group_id": CHECKSET,
			"caches.0.image_cache_id":     CHECKSET,
			"caches.0.image_cache_name":   name,
			"caches.0.images.#":           "1",
			"caches.0.images.0":           fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
			"caches.0.progress":           "100%",
			"caches.0.snapshot_id":        CHECKSET,
			"caches.0.status":             "Ready",
			"caches.0.events.#":           CHECKSET,
		}
	}

	var fakeEciImageCachesMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"caches.#": "0",
		}
	}

	var eciImageCachesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_eci_image_caches.default",
		existMapFunc: existEciImageCachesMapCheck,
		fakeMapFunc:  fakeEciImageCachesMapCheck,
	}

	eciImageCachesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, statusConf, idsConf)
}

func dataSourceEciImageCachesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf-eci-image-test"
  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
}
resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eci_image_cache" "default" {
  image_cache_name = var.name
  images            = ["registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine"]
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0
  eip_instance_id = alicloud_eip_address.default.id
}`, name, defaultRegionToTest)
}
