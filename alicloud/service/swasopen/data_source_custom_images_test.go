package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerCustomImagesDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_custom_images.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-swascustomimage-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerCustomImagesDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_simple_application_server_custom_image.default.id}-fake"},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}",
			"ids":        []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_simple_application_server_custom_image.default.custom_image_name}-fake",
			"ids":        []string{"${alicloud_simple_application_server_custom_image.default.id}"},
		}),
	}
	var existSimpleApplicationServerCustomImageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"ids.0":                      acctest.CHECKSET,
			"images.#":                   "1",
			"images.0.id":                acctest.CHECKSET,
			"images.0.custom_image_id":   acctest.CHECKSET,
			"images.0.custom_image_name": fmt.Sprintf("tf-testacc-swascustomimage-%d", rand),
			"images.0.description":       fmt.Sprintf("tf-testacc-swascustomimage-%d", rand),
			"images.0.platform":          acctest.CHECKSET,
		}
	}

	var fakeSimpleApplicationServerCustomImageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"images.#": "0",
		}
	}

	var SimpleApplicationServerCustomImageCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existSimpleApplicationServerCustomImageMapFunc,
		FakeMapFunc:  fakeSimpleApplicationServerCustomImageMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}

	SimpleApplicationServerCustomImageCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
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
