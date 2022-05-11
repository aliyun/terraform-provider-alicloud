package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerCustomImagesDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_custom_images.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-swascustomimage-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerCustomImagesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_custom_image.default.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}",
			"ids":        []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}-fake",
			"ids":        []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
	}
	var existSimpleApplicationServerCustomImageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"ids.0":                      CHECKSET,
			"images.#":                   "1",
			"images.0.id":                CHECKSET,
			"images.0.custom_image_id":   CHECKSET,
			"images.0.custom_image_name": fmt.Sprintf("tf-testacc-swascustomimage-%d", rand),
			"images.0.description":       fmt.Sprintf("tf-testacc-swascustomimage-%d", rand),
			"images.0.platform":          CHECKSET,
		}
	}

	var fakeSimpleApplicationServerCustomImageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"images.#": "0",
		}
	}

	var SimpleApplicationServerCustomImageCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSimpleApplicationServerCustomImageMapFunc,
		fakeMapFunc:  fakeSimpleApplicationServerCustomImageMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}

	SimpleApplicationServerCustomImageCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceSimpleApplicationServerCustomImagesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_images" "default" {
	platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
	platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}
data "alicloud_simple_application_server_disks" "default" {
  disk_type = "System"
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
  snapshot_name = var.name
}

resource "alicloud_simple_application_server_custom_image" "default" {
  custom_image_name  = var.name
  description        = var.name
  system_snapshot_id = alicloud_simple_application_server_snapshot.default.id
  instance_id        = data.alicloud_simple_application_server_disks.default.disks.0.instance_id
}`, name)
}
