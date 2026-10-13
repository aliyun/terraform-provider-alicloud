package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDataWorksRoutesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_routes.default"
	name := fmt.Sprintf("tf-testacc-dataworksroute%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksRoutesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_data_works_dw_resource_group.defaultVJvKvl.id}",
			"network_id":        "${alicloud_data_works_network.default.id}",
			"ids":               []string{"${alicloud_data_works_route.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_data_works_dw_resource_group.defaultVJvKvl.id}",
			"network_id":        "${alicloud_data_works_network.default.id}",
			"ids":               []string{"${alicloud_data_works_route.default.id}_fake"},
		}),
	}

	var existDataWorksRoutesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"routes.#":            "1",
			"routes.0.route_id":   CHECKSET,
			"routes.0.network_id": CHECKSET,
		}
	}

	var fakeDataWorksRoutesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"routes.#": "0",
			"ids.#":    "0",
		}
	}

	var DataWorksRoutesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksRoutesMapFunc,
		fakeMapFunc:  fakeDataWorksRoutesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		testAccPreCheck(t)
	}
	DataWorksRoutesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func dataSourceDataWorksRoutesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
}

data "alicloud_vpcs" "default2" {
	name_regex = "^default-NODELETING-2$"
}

data "alicloud_vswitches" "default2" {
	vpc_id  = data.alicloud_vpcs.default2.ids.0
}


resource "alicloud_data_works_dw_resource_group" "defaultVJvKvl" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = data.alicloud_vswitches.default.ids.0
  remark                = "OpenAPI测试用资源组"
  resource_group_name   = var.name
  default_vpc_id        = data.alicloud_vpcs.default.ids.0
  auto_renew = false
}

resource "alicloud_data_works_network" "default" {
  vpc_id               = data.alicloud_vpcs.default2.ids.0
  vswitch_id           = data.alicloud_vswitches.default2.ids.0
  dw_resource_group_id = alicloud_data_works_dw_resource_group.defaultVJvKvl.id
}

resource "alicloud_data_works_route" "default" {
  network_id       = alicloud_data_works_network.default.id
  destination_cidr = "192.168.30.0/24"
}

`, name)
}
