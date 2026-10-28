package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerDisksDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_disks.default"
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testacc-swasdisk-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerDisksDependence)

	diskTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"disk_type": "System",
		}),
		FakeConfig: "",
	}

	instanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_simple_application_server_instance.default.id}",
		}),
		FakeConfig: "",
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"status": "In_use",
		}),
		FakeConfig: "",
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_simple_application_server_instance.default.id}",
			"status":      "In_use",
			"disk_type":   "System",
		}),
		FakeConfig: "",
	}
	var existSimpleApplicationServerDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                acctest.CHECKSET,
			"names.#":              acctest.CHECKSET,
			"disks.#":              acctest.CHECKSET,
			"disks.0.id":           acctest.CHECKSET,
			"disks.0.device":       acctest.CHECKSET,
			"disks.0.disk_id":      acctest.CHECKSET,
			"disks.0.category":     acctest.CHECKSET,
			"disks.0.create_time":  acctest.CHECKSET,
			"disks.0.disk_name":    acctest.CHECKSET,
			"disks.0.disk_type":    acctest.CHECKSET,
			"disks.0.instance_id":  acctest.CHECKSET,
			"disks.0.payment_type": acctest.CHECKSET,
			"disks.0.size":         acctest.CHECKSET,
			"disks.0.status":       acctest.CHECKSET,
		}
	}

	var fakeSimpleApplicationServerDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"disks.#": "0",
		}
	}

	var SimpleApplicationServerDiskCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existSimpleApplicationServerDiskMapFunc,
		FakeMapFunc:  fakeSimpleApplicationServerDiskMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}

	SimpleApplicationServerDiskCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, diskTypeConf, instanceIdConf, statusConf, allConf)
}

func dataSourceSimpleApplicationServerDisksDependence(name string) string {
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
`, name)
}
