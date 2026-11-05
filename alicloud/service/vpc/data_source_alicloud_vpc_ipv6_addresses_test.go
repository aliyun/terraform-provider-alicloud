package vpc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVPCIpv6AddressesDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_ipv6_addresses.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-vpcipv6address-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6AddressesDependence)

	associatedInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"associated_instance_id": "${data.alicloud_instances.default.instances.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"associated_instance_id": "${data.alicloud_instances.default.instances.0.id}_fake",
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}_fake",
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${data.alicloud_vpcs.default.ids.0}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Pending",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":             "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
			"associated_instance_id": "${data.alicloud_instances.default.instances.0.id}",
			"status":                 "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}_fake",
			"vswitch_id":             "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}_fake",
			"associated_instance_id": "${data.alicloud_instances.default.instances.0.id}_fake",
			"status":                 "Pending",
		}),
	}
	var existVpcIpv6AddressMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                acctest.CHECKSET,
			"ids.0":                                acctest.CHECKSET,
			"addresses.#":                          acctest.CHECKSET,
			"addresses.0.id":                       acctest.CHECKSET,
			"addresses.0.status":                   "Available",
			"addresses.0.associated_instance_id":   acctest.CHECKSET,
			"addresses.0.associated_instance_type": acctest.CHECKSET,
			"addresses.0.ipv6_address":             acctest.CHECKSET,
			"addresses.0.ipv6_address_id":          acctest.CHECKSET,
			"addresses.0.ipv6_address_name":        "",
			"addresses.0.ipv6_gateway_id":          acctest.CHECKSET,
			"addresses.0.network_type":             acctest.CHECKSET,
			"addresses.0.create_time":              acctest.CHECKSET,
			"addresses.0.vswitch_id":               acctest.CHECKSET,
			"addresses.0.vpc_id":                   acctest.CHECKSET,
		}
	}

	var fakeVpcIpv6AddressMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"addresses.#": "0",
		}
	}

	var VpcIpv6AddressCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6AddressMapFunc,
		fakeMapFunc:  fakeVpcIpv6AddressMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ECS_WITH_IPV6_ADDRESS")
	}

	VpcIpv6AddressCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, associatedInstanceIdConf, vswitchIdConf, vpcIdConf, statusConf, allConf)
}

func dataSourceVpcIpv6AddressesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alicloud_vpcs" "default" {
  name_regex = "no-deleteing-ipv6-address"
}
`, name)
}
