package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudNetworkInterfacesDataSourceBasic(t *testing.T) {

	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_network_interface_attachment.default.network_interface_id}_fake" ]`,
		}),
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"instance_id": `"${alicloud_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"instance_id": `"${alicloud_instance.default.id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d_fake"`, rand),
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"vpc_id": `"${alicloud_vpc.default.id}"`,
		}),
	}

	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"vswitch_id": `"${alicloud_vswitch.default.id}_fake"`,
		}),
	}

	privateIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"private_ip": `"192.168.0.2"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"private_ip": `"192.168.0.1"`,
		}),
	}

	securityGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":               `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"security_group_id": `"${alicloud_security_group.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":               `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"security_group_id": `"${alicloud_security_group.default.id}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":  `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"type": `"Secondary"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":  `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"type": `"Primary"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d"
						   }`, rand),
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d_fake"
						   }`, rand),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d"
						   }`, rand),
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"private_ip":        `"192.168.0.2"`,
			"security_group_id": `"${alicloud_security_group.default.id}"`,
			"type":              `"Secondary"`,
			"instance_id":       `"${alicloud_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_network_interface_attachment.default.network_interface_id}" ]`,
			"name_regex": fmt.Sprintf(`"tf-testAccNetworkInterfacesBasic%d"`, rand),
			"tags": fmt.Sprintf(`{
							 TF-VER = "0.11.3%d_fake"
						   }`, rand),
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"private_ip":        `"192.168.0.2"`,
			"security_group_id": `"${alicloud_security_group.default.id}"`,
			"type":              `"Primary"`,
			"instance_id":       `"${alicloud_instance.default.id}"`,
		}),
	}

	networkInterfacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, instanceIdConf, nameRegexConf, vpcIdConf, vswitchIdConf, privateIpConf,
		securityGroupIdConf, typeConf, tagsConf, allConf)
}

func testAccCheckAlicloudNetworkInterfacesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
 default = "tf-testAccNetworkInterfacesBasic"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    name = "${var.name}%d"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3%d"
	}
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
  eni_amount = 2
}

data "alicloud_images" "default" {
  	most_recent = true
	owners = "system"
}

resource "alicloud_instance" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    security_groups = ["${alicloud_security_group.default.id}"]
    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alicloud_images.default.images.0.image_id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface_attachment" "default" {
    instance_id = "${alicloud_instance.default.id}"
    network_interface_id = "${alicloud_network_interface.default.id}"
}

data "alicloud_network_interfaces" "default"  {
	%s
}`, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                          "1",
		"names.#":                        "1",
		"interfaces.#":                   "1",
		"interfaces.0.id":                CHECKSET,
		"interfaces.0.name":              fmt.Sprintf("tf-testAccNetworkInterfacesBasic%d", rand),
		"interfaces.0.status":            CHECKSET,
		"interfaces.0.vpc_id":            CHECKSET,
		"interfaces.0.zone_id":           CHECKSET,
		"interfaces.0.public_ip":         "",
		"interfaces.0.private_ip":        "192.168.0.2",
		"interfaces.0.private_ips.#":     "0",
		"interfaces.0.security_groups.#": "1",
		"interfaces.0.description":       "Basic test",
		"interfaces.0.instance_id":       CHECKSET,
		"interfaces.0.creation_time":     CHECKSET,
		"interfaces.0.tags.%":            "1",
		"interfaces.0.tags.TF-VER":       fmt.Sprintf("0.11.3%d", rand),
	}
}

var fakeNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#": "0",
		"names.#":      "0",
		"ids.#":        "0",
	}
}

var networkInterfacesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_network_interfaces.default",
	existMapFunc: existNetworkInterfacesMapFunc,
	fakeMapFunc:  fakeNetworkInterfacesMapFunc,
}
