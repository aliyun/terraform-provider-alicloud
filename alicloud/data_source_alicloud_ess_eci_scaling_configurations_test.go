package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEssEciScalingconfigurationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_eci_scaling_configuration.default.id}"]`,
			"name_regex":       `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_eci_scaling_configuration.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
	}

	var existEssEciScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"configurations.#":                  "1",
			"configurations.0.name":             fmt.Sprintf("tf-testAccDataSourceEssConfigurations-%d", rand),
			"configurations.0.scaling_group_id": CHECKSET,
			"configurations.0.host_name":        "hostname",
			"configurations.0.spot_strategy":    "SpotWithPriceLimit",
		}
	}

	var fakeEssEciScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"configurations.#": "0",
			"ids.#":            "0",
			"names.#":          "0",
		}
	}

	var essScalingconfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_eci_scaling_configurations.default",
		existMapFunc: existEssEciScalingconfigurationsMapFunc,
		fakeMapFunc:  fakeEssEciScalingconfigurationsMapFunc,
	}

	essScalingconfigurationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudEssEciScalingconfigurationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssConfigurations-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	group_type = "ECI"
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_ess_eci_scaling_configuration" "default" {
  scaling_group_id           = "${alicloud_ess_scaling_group.default.id}"
  scaling_configuration_name = "${var.name}"
  security_group_id          = "${alicloud_security_group.default.id}"
  restart_policy       = "Always"
  container_group_name = "${var.name}"
  resource_group_id    = "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}"
  description          = "newDesc"
  dns_policy           = "Default"
  spot_price_limit     = 1.2
  auto_create_eip      = true
  memory               = 2
  eip_bandwidth        = 5
  ram_role_name        = "ramRoleName"
  ingress_bandwidth    = 1024000
  host_name            = "hostname"
  spot_strategy        = "SpotWithPriceLimit"
  cpu                  = 1
  tags = {
    "name" : "tf-test2"
  }	
  containers {
    ports {
      protocol = "newProtocol"
      port     = 2
    }
    environment_vars {
      key   = "newKey"
      value = "newValue"
    }
    working_dir       = "newWorkingDir"
    args              = ["arg2"]
    cpu               = 2
    gpu               = 2
    memory            = 2
    name              = "newName"
    image             = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
    image_pull_policy = "Always"
    volume_mounts {
      mount_path = "newPath"
      name       = "newName"
      read_only  = false
    }
    commands = ["cmd2"]
  }

 volumes {
    config_file_volume_config_file_to_paths {
      content = "content2"
      path    = "path2"
    }
    disk_volume_disk_id   = "disk_volume_disk_id2"
    disk_volume_fs_type   = "disk_volume_fs_type2"
    disk_volume_disk_size = 2
    flex_volume_driver    = "flex_volume_driver2"
    flex_volume_fs_type   = "flex_volume_fs_type2"
    flex_volume_options   = "flex_volume_options2"
    nfs_volume_path       = "nfs_volume_path2"
    nfs_volume_read_only  = false
    nfs_volume_server     = "nfs_volume_server2"
    name                  = "name2"
    type                  = "type2"
  }

  host_aliases {
    hostnames = ["hostnames2"]
    ip        = "ip2"
  }

  image_registry_credentials {
    password = "password"
    server   = "server"
    username = "username"
  }

  init_containers {
    ports {
      protocol = "protocol"
      port     = 1
    }
    environment_vars {
      key   = "key"
      value = "value"
    }
    working_dir       = "workingDir"
    args              = ["arg"]
    cpu               = 1
    gpu               = 1
    memory            = 1
    name              = "name"
    image             = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
    image_pull_policy = "Always"
    volume_mounts {
      mount_path = "path"
      name       = "name"
      read_only  = true
    }
    commands = ["cmd"]
  }

  acr_registry_infos {
    domains       = ["test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"]
    instance_id   = "cri-47rme9691uiowvfv"
    region_id     = "cn-hangzhou"
    instance_name = "zzz"
  }

}


data "alicloud_ess_eci_scaling_configurations" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
