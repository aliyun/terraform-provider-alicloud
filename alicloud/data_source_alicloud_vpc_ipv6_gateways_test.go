package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCIpv6GatewaysDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_ipv6_gateways.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.VpcIpv6GatewaySupportRegions)
	name := fmt.Sprintf("tf-testacc-vpcipv6gateway-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6GatewaysDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_ipv6_gateway.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"status": "Deleting",
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"vpc_id": "${alicloud_vpc_ipv6_gateway.default.vpc_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"vpc_id": "${alicloud_vpc_ipv6_gateway.default.vpc_id}-fake",
		}),
	}
	ipv6GatewayNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"ipv6_gateway_name": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"ipv6_gateway_name": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}",
			"ids":               []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"status":            "Available",
			"vpc_id":            "${alicloud_vpc_ipv6_gateway.default.vpc_id}",
			"ipv6_gateway_name": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}-fake",
			"ids":               []string{"${alicloud_vpc_ipv6_gateway.default.id}"},
			"status":            "Deleting",
			"vpc_id":            "${alicloud_vpc_ipv6_gateway.default.vpc_id}-fake",
			"ipv6_gateway_name": "${alicloud_vpc_ipv6_gateway.default.ipv6_gateway_name}-fake",
		}),
	}
	var existVpcIpv6GatewayMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           CHECKSET,
			"names.#":                         "1",
			"names.0":                         fmt.Sprintf("tf-testacc-vpcipv6gateway-%d", rand),
			"gateways.#":                      "1",
			"gateways.0.id":                   CHECKSET,
			"gateways.0.ipv6_gateway_id":      CHECKSET,
			"gateways.0.ipv6_gateway_name":    fmt.Sprintf("tf-testacc-vpcipv6gateway-%d", rand),
			"gateways.0.description":          fmt.Sprintf("tf-testacc-vpcipv6gateway-%d", rand),
			"gateways.0.status":               "Available",
			"gateways.0.spec":                 "Small",
			"gateways.0.vpc_id":               CHECKSET,
			"gateways.0.create_time":          CHECKSET,
			"gateways.0.instance_charge_type": "PayAsYouGo",
			"gateways.0.expired_time":         "",
			"gateways.0.business_status":      "Normal",
		}
	}

	var fakeVpcIpv6GatewayMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"gateways.#": "0",
		}
	}

	var VpcIpv6GatewayCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6GatewayMapFunc,
		fakeMapFunc:  fakeVpcIpv6GatewayMapFunc,
	}

	VpcIpv6GatewayCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, vpcIdConf, ipv6GatewayNameConf, allConf)
}

func dataSourceVpcIpv6GatewaysDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alicloud_vpc_ipv6_gateway" "default" {
  vpc_id            = alicloud_vpc.default.id
  ipv6_gateway_name = var.name
  description       = var.name
}`, name)
}
