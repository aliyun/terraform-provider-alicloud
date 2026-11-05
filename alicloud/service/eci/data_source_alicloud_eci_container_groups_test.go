package eci_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEciContainerGroupsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_container_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_container_group.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_container_group.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_container_group.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_container_group.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_container_group.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_container_group.default.id}"]`,
			"tags": `{
				"created" = "tf"
				"for" = "acceptance-test"
		}`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_container_group.default.id}"]`,
			"tags": `{
				"created" = "tf-fake"
				"for" = "acceptance-test-fake"
			}`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_container_group.default.container_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_container_group.default.container_group_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_container_group.default.id}"]`,
			"status": `"${alicloud_eci_container_group.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_container_group.default.id}"]`,
			"status": `"ScheduleFailed"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_container_group.default.id}"]`,
			"name_regex":        `"${alicloud_eci_container_group.default.container_group_name}"`,
			"resource_group_id": `"${alicloud_eci_container_group.default.resource_group_id}"`,
			"status":            `"${alicloud_eci_container_group.default.status}"`,
			"tags": `{
				"created" = "tf"
				"for" = "acceptance-test"
		}`,
		}),
		FakeConfig: testAccCheckAlicloudEciContainerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_container_group.default.id}_fake"]`,
			"name_regex":        `"${alicloud_eci_container_group.default.container_group_name}_fake"`,
			"resource_group_id": `"${alicloud_eci_container_group.default.resource_group_id}_fake"`,
			"status":            `"ScheduleFailed"`,
			"tags": `{
				"created" = "tf-fake"
				"for" = "acceptance-test-fake"
			}`,
		}),
	}
	var existAlicloudEciContainerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"names.#":                                      "1",
			"groups.#":                                     "1",
			"groups.0.id":                                  acctest.CHECKSET,
			"groups.0.container_group_id":                  acctest.CHECKSET,
			"groups.0.container_group_name":                acctest.CHECKSET,
			"groups.0.containers.#":                        "1",
			"groups.0.containers.0.image":                  fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
			"groups.0.containers.0.name":                   "nginx",
			"groups.0.containers.0.image_pull_policy":      "IfNotPresent",
			"groups.0.containers.0.volume_mounts.#":        "1",
			"groups.0.containers.0.ports.#":                "1",
			"groups.0.containers.0.environment_vars.#":     "1",
			"groups.0.cpu":                                 "2",
			"groups.0.dns_config.#":                        "0",
			"groups.0.host_aliases.#":                      "1",
			"groups.0.init_containers.#":                   "1",
			"groups.0.init_containers.0.name":              "init-busybox",
			"groups.0.init_containers.0.image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
			"groups.0.init_containers.0.image_pull_policy": "IfNotPresent",
			"groups.0.memory":                              "4",
			"groups.0.resource_group_id":                   acctest.CHECKSET,
			"groups.0.security_group_id":                   acctest.CHECKSET,
			"groups.0.status":                              "Running",
			"groups.0.vswitch_id":                          acctest.CHECKSET,
			"groups.0.vpc_id":                              acctest.CHECKSET,
			"groups.0.zone_id":                             acctest.CHECKSET,
			"groups.0.volumes.#":                           "1",
		}
	}
	var fakeAlicloudEciContainerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEciContainerGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_eci_container_groups.default",
		ExistMapFunc: existAlicloudEciContainerGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEciContainerGroupsDataSourceNameMapFunc,
	}
	alicloudEciContainerGroupsCheckInfo.DataSourceTestCheck(t, rand, idsConf, resourceGroupIdConf, tagsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEciContainerGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testacccontainergroup-%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "tf-eci-image-test"
  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
}

resource "alicloud_eci_container_group" "default" {
  container_group_name = var.name
  restart_policy       = "OnFailure"
  security_group_id    = alicloud_security_group.group.id
  vswitch_id           = data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0
  tags = {
	"created" = "tf"
	"for" = "acceptance-test"  
  }
  #################################
  # containers
  #################################
  containers {
    image             = "registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine"
    name              = "nginx"
    working_dir       = "/tmp/nginx"
    image_pull_policy = "IfNotPresent"
    commands          = ["/bin/sh", "-c", "sleep 9999"]
    volume_mounts {
      mount_path = "/tmp/test"
      read_only  = false
      name       = "empty1"
    }
    ports {
      port     = 80
      protocol = "TCP"
    }
    environment_vars {
      key   = "test"
      value = "nginx"
    }
  }
  host_aliases {
    ip        = "1.1.1.1"
    hostnames = ["hehe.com"]
  }

  #################################
  # init_containers
  #################################
  init_containers {
    name              = "init-busybox"
    image             = "registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30"
    image_pull_policy = "IfNotPresent"
    commands          = ["echo"]
    args              = ["hello initcontainer"]
  }

  #################################
  # volumes
  #################################
  volumes {
    name = "empty1"
    type = "EmptyDirVolume"
  }
}

data "alicloud_eci_container_groups" "default" {	
	enable_details = true
	%s	
}
`, rand, acctest.DefaultRegionToTest, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
