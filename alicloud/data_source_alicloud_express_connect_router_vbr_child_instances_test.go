package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterVbrChildInstancesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_vbr_child_instances.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervbrchildinstance%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterVbrChildInstancesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vbr_child_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vbr_child_instance.default.id}_fake"},
		}),
	}

	childInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"child_instance_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"child_instance_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_id}_fake",
		}),
	}

	childInstanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":              "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"child_instance_type": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":              "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"ids":                 []string{"${alicloud_express_connect_router_vbr_child_instance.default.id}_fake"},
			"child_instance_type": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_type}",
		}),
	}

	childInstanceRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                   "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"child_instance_region_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                   "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"child_instance_region_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_region_id}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"status": "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                   "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"ids":                      []string{"${alicloud_express_connect_router_vbr_child_instance.default.id}"},
			"child_instance_id":        "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_id}",
			"child_instance_type":      "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_type}",
			"child_instance_region_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_region_id}",
			"status":                   "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                   "${alicloud_express_connect_router_vbr_child_instance.default.ecr_id}",
			"ids":                      []string{"${alicloud_express_connect_router_vbr_child_instance.default.id}_fake"},
			"child_instance_id":        "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_id}_fake",
			"child_instance_type":      "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_type}",
			"child_instance_region_id": "${alicloud_express_connect_router_vbr_child_instance.default.child_instance_region_id}_fake",
			"status":                   "CREATING",
		}),
	}

	var existAliCloudExpressConnectRouterVbrChildInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"instances.#":                          "1",
			"instances.0.id":                       CHECKSET,
			"instances.0.ecr_id":                   CHECKSET,
			"instances.0.child_instance_id":        CHECKSET,
			"instances.0.child_instance_type":      CHECKSET,
			"instances.0.child_instance_owner_id":  CHECKSET,
			"instances.0.child_instance_region_id": CHECKSET,
			"instances.0.description":              CHECKSET,
			"instances.0.status":                   CHECKSET,
			"instances.0.create_time":              CHECKSET,
			"instances.0.modify_time":              CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterVbrChildInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var aliCloudExpressConnectRouterVbrChildInstancesInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_router_vbr_child_instances.default",
		existMapFunc: existAliCloudExpressConnectRouterVbrChildInstancesMapFunc,
		fakeMapFunc:  fakeAliCloudExpressConnectRouterVbrChildInstancesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterVbrChildInstancesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, childInstanceIdConf, childInstanceTypeConf, childInstanceRegionIdConf, statusConf, allConf)
}

func dataSourceExpressConnectRouterVbrChildInstancesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65532"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.default.connections.0.id
  vlan_id                = "1000"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}

resource "alicloud_express_connect_router_vbr_child_instance" "default" {
  ecr_id                   = alicloud_express_connect_router_express_connect_router.default.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.default.id
  child_instance_type      = "VBR"
  child_instance_owner_id  = data.alicloud_account.default.id
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
  description              = var.name
}
`, name)
}
