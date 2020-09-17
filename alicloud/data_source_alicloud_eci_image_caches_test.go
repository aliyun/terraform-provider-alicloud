package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEciImageCachesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_eci_image_caches.default", strconv.FormatInt(int64(rand), 10), dataSourceEciImageCachesConfigDependence)
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
			"caches.0.image_cache_name":   "test-image-cache",
			"caches.0.images.#":           "1",
			"caches.0.images.0":           "registry.cn-beijing.aliyuncs.com/sceneplatform/sae-image-demo:latest",
			"caches.0.progress":           "100%",
			"caches.0.snapshot_id":        CHECKSET,
			"caches.0.status":             "Ready",
			"caches.0.events.#":           "3",
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

	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
		testAccPreCheckWithNoDefaultVswitch(t)
	}

	eciImageCachesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, statusConf, idsConf)
}

func dataSourceEciImageCachesConfigDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = true
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
resource "alicloud_security_group" "default" {
  name        = "test-image-cache"
  description = "tf-eci-image-test"
  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
}
resource "alicloud_eip" "default" {
  name = "test-image-cache"
}

resource "alicloud_eci_image_cache" "default" {
  image_cache_name = "test-image-cache"
  images            = ["registry.cn-beijing.aliyuncs.com/sceneplatform/sae-image-demo:latest"]
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0
  eip_instance_id = alicloud_eip.default.id
}`)
}
