package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerSharedResourcesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `[format("%s:%s",alicloud_resource_manager_shared_resource.example.resource_id,alicloud_resource_manager_shared_resource.example.resource_type)]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `[format("%s:%s",alicloud_resource_manager_shared_resource.example.resource_id,alicloud_resource_manager_shared_resource.example.resource_type)]`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `["fake"]`,
			"status":            `"Associating"`,
		}),
	}

	var existResourceManagerSharedResourcesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"ids.#":                         "1",
			"resources.0.id":                CHECKSET,
			"resources.0.resource_type":     "VSwitch",
			"resources.0.resource_share_id": CHECKSET,
			"resources.0.resource_id":       CHECKSET,
			"resources.0.status":            "Associated",
		}
	}

	var fakeResourceManagerSharedResourcessMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": "0",
			"ids.#":       "0",
		}
	}

	var resourceManagerSharedResourcesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_shared_resources.example",
		existMapFunc: existResourceManagerSharedResourcesMapFunc,
		fakeMapFunc:  fakeResourceManagerSharedResourcessMapFunc,
	}

	resourceManagerSharedResourcesCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf)

}

func testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccResourceManagerSharedResources%d"
}
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "example" {
  name = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vswitch" "example" {
  availability_zone = data.alicloud_zones.example.ids.0
  cidr_block = "192.168.0.0/16"
  vpc_id = alicloud_vpc.example.id
}
resource "alicloud_resource_manager_resource_share" "example" {
	resource_share_name = var.name
}
resource "alicloud_resource_manager_shared_resource" "example" {
  resource_id       = alicloud_vswitch.example.id
  resource_share_id = alicloud_resource_manager_resource_share.example.id
  resource_type     = "VSwitch"
}
data "alicloud_resource_manager_shared_resources" "example"{
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
