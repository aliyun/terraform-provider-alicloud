package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerDisksDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_disks.default"
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testacc-simpleapplicationserverdisk-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerDisksDependence)

	diskTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"disk_type": "System",
		}),
		fakeConfig: "",
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
		}),
		fakeConfig: "",
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "In_use",
		}),
		fakeConfig: "",
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"status":      "In_use",
			"disk_type":   "System",
		}),
		fakeConfig: "",
	}
	var existSimpleApplicationServerDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                CHECKSET,
			"names.#":              CHECKSET,
			"disks.#":              CHECKSET,
			"disks.0.id":           CHECKSET,
			"disks.0.device":       CHECKSET,
			"disks.0.disk_id":      CHECKSET,
			"disks.0.category":     CHECKSET,
			"disks.0.create_time":  CHECKSET,
			"disks.0.disk_name":    CHECKSET,
			"disks.0.disk_type":    CHECKSET,
			"disks.0.instance_id":  CHECKSET,
			"disks.0.payment_type": CHECKSET,
			"disks.0.size":         CHECKSET,
			"disks.0.status":       CHECKSET,
		}
	}

	var fakeSimpleApplicationServerDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"disks.#": "0",
		}
	}

	var SimpleApplicationServerDiskCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSimpleApplicationServerDiskMapFunc,
		fakeMapFunc:  fakeSimpleApplicationServerDiskMapFunc,
	}

	SimpleApplicationServerDiskCheckInfo.dataSourceTestCheck(t, rand, diskTypeConf, instanceIdConf, statusConf, allConf)
}

func dataSourceSimpleApplicationServerDisksDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_instances" "default" {}

data "alicloud_simple_application_server_images" "default" {}

data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server_instance" "default" {
  count          = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = "tf-testaccswas-disks"
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
}

locals {
  instance_id = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? data.alicloud_simple_application_server_instances.default.ids.0 : alicloud_simple_application_server_instance.default.0.id
}
`, name)
}
