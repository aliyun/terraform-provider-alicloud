package ecp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEcpInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccInstance-%d", rand)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecp_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecp_instance.default.id}_fake"]`,
		}),
	}
	zoneIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_ecp_instance.default.id}"]`,
			"zone_id": `"${data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_ecp_instance.default.id}_fake"]`,
			"zone_id": `"${data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id}"`,
		}),
	}
	imageIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecp_instance.default.id}"]`,
			"image_id": `"${alicloud_ecp_instance.default.image_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecp_instance.default.id}"]`,
			"image_id": `"${alicloud_ecp_instance.default.image_id}_fake"`,
		}),
	}

	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"instance_name": `"${alicloud_ecp_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"instance_name": `"${alicloud_ecp_instance.default.instance_name}_fake"`,
		}),
	}
	instanceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"instance_type": `"${local.instance_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}_fake"]`,
			"instance_type": `"${local.instance_type}"`,
		}),
	}
	keyPairNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"key_pair_name": `"${alicloud_ecp_instance.default.key_pair_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"key_pair_name": `"${alicloud_ecp_instance.default.key_pair_name}_fake"`,
		}),
	}
	resolutionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_instance.default.id}"]`,
			"resolution": `"${alicloud_ecp_instance.default.resolution}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_instance.default.id}"]`,
			"resolution": `"${alicloud_ecp_instance.default.resolution}_fake"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecp_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecp_instance.default.instance_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecp_instance.default.id}"]`,
			"status": `"${alicloud_ecp_instance.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecp_instance.default.id}"]`,
			"status": `"Stopping"`,
		}),
	}
	chargeTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecp_instance.default.id}"]`,
			"payment_type": `"PayAsYouGo"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecp_instance.default.id}"]`,
			"payment_type": `"Subscription"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}"]`,
			"image_id":      `"${alicloud_ecp_instance.default.image_id}"`,
			"instance_name": `"${alicloud_ecp_instance.default.instance_name}"`,
			"payment_type":  `"PayAsYouGo"`,
			"instance_type": `"${local.instance_type}"`,
			"key_pair_name": `"${alicloud_ecp_instance.default.key_pair_name}"`,
			"name_regex":    `"${alicloud_ecp_instance.default.instance_name}"`,
			"resolution":    `"${alicloud_ecp_instance.default.resolution}"`,
			"status":        `"${alicloud_ecp_instance.default.status}"`,
			"zone_id":       `"${data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcpInstancesDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecp_instance.default.id}_fake"]`,
			"image_id":      `"${alicloud_ecp_instance.default.image_id}_fake"`,
			"instance_name": `"${alicloud_ecp_instance.default.instance_name}_fake"`,
			"instance_type": `"${local.instance_type}"`,
			"payment_type":  `"Subscription"`,
			"key_pair_name": `"${alicloud_ecp_instance.default.key_pair_name}_fake"`,
			"name_regex":    `"${alicloud_ecp_instance.default.instance_name}_fake"`,
			"resolution":    `"${alicloud_ecp_instance.default.resolution}_fake"`,
			"status":        `"Stopping"`,
			"zone_id":       `"${data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id}"`,
		}),
	}
	var existAlicloudEcpInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.description":       name,
			"instances.0.image_id":          `android_9_0_0_release_2851157_20211201.vhd`,
			"instances.0.instance_name":     name,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.key_pair_name":     name,
			"instances.0.resolution":        acctest.CHECKSET,
			"instances.0.security_group_id": acctest.CHECKSET,
			"instances.0.vswitch_id":        acctest.CHECKSET,
			"instances.0.payment_type":      acctest.CHECKSET,
			"instances.0.zone_id":           acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
		}
	}
	var fakeAlicloudEcpInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcpInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecp_instances.default",
		ExistMapFunc: existAlicloudEcpInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcpInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudEcpInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, chargeTypeConf, zoneIdConf, imageIdConf, instanceNameConf, keyPairNameConf, resolutionConf, nameRegexConf, statusConf, instanceTypeConf, allConf)
}

func testAccCheckAlicloudEcpInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccInstance-%d"
}

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

resource "alicloud_ecp_instance" "default" {
  instance_name     = var.name
  description       = var.name
  key_pair_name     = "${alicloud_ecp_key_pair.default.key_pair_name}"
  security_group_id = "${alicloud_security_group.group.id}"
  vswitch_id        = "${data.alicloud_vswitches.default.ids.0}"
  image_id          = "android_9_0_0_release_2851157_20211201.vhd"
  instance_type     = "${local.instance_type}"
  vnc_password      = "Cp1234"
  force             = "true"
  payment_type      = "PayAsYouGo"
}

data "alicloud_ecp_instances" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
