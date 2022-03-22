package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBServerGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	serverGroupIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_ids": `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_ids": `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_server_group.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_server_group.default.id}_fake"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_server_group.default.server_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"vpc_id": `"${alicloud_alb_server_group.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}_fake"]`,
			"vpc_id": `"${alicloud_alb_server_group.default.vpc_id}"`,
		}),
	}
	serverGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}"]`,
			"name_regex":        `"${alicloud_alb_server_group.default.server_group_name}"`,
			"vpc_id":            `"${alicloud_alb_server_group.default.vpc_id}"`,
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}"`,
			"status":            `"Available"`,
			"tags":              `{Created = "TF"}`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}"`,
			"server_group_ids":  `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}_fake"]`,
			"name_regex":        `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
			"vpc_id":            `"${alicloud_alb_server_group.default.vpc_id}"`,
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
			"status":            `"Configuring"`,
			"tags":              `{Created = "TF_fake"}`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}_fake"`,
			"server_group_ids":  `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	var existDataAlicloudAlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                    "1",
			"groups.#":                                                 "1",
			"groups.0.server_group_name":                               fmt.Sprintf("tf-testAccAlbServerGroup-%d", rand),
			"groups.0.protocol":                                        "HTTP",
			"groups.0.scheduler":                                       CHECKSET,
			"groups.0.id":                                              CHECKSET,
			"groups.0.server_group_id":                                 CHECKSET,
			"groups.0.status":                                          CHECKSET,
			"groups.0.vpc_id":                                          CHECKSET,
			"groups.0.servers.#":                                       "1",
			"groups.0.servers.0.server_type":                           "Ecs",
			"groups.0.servers.0.description":                           CHECKSET,
			"groups.0.servers.0.port":                                  "80",
			"groups.0.servers.0.server_id":                             CHECKSET,
			"groups.0.servers.0.server_ip":                             CHECKSET,
			"groups.0.servers.0.weight":                                "10",
			"groups.0.servers.0.status":                                CHECKSET,
			"groups.0.sticky_session_config.#":                         "1",
			"groups.0.sticky_session_config.0.cookie":                  CHECKSET,
			"groups.0.sticky_session_config.0.cookie_timeout":          CHECKSET,
			"groups.0.sticky_session_config.0.sticky_session_enabled":  "true",
			"groups.0.sticky_session_config.0.sticky_session_type":     CHECKSET,
			"groups.0.health_check_config.#":                           "1",
			"groups.0.health_check_config.0.health_check_enabled":      "true",
			"groups.0.health_check_config.0.health_check_connect_port": CHECKSET,
			"groups.0.health_check_config.0.health_check_codes.#":      "3",
			"groups.0.health_check_config.0.health_check_host":         CHECKSET,
			"groups.0.health_check_config.0.health_check_http_version": CHECKSET,
			"groups.0.health_check_config.0.health_check_interval":     CHECKSET,
			"groups.0.health_check_config.0.health_check_method":       CHECKSET,
			"groups.0.health_check_config.0.health_check_path":         CHECKSET,
			"groups.0.health_check_config.0.health_check_protocol":     CHECKSET,
			"groups.0.health_check_config.0.health_check_timeout":      CHECKSET,
			"groups.0.health_check_config.0.healthy_threshold":         CHECKSET,
			"groups.0.health_check_config.0.unhealthy_threshold":       CHECKSET,
			"groups.0.tags.%":                                          "1",
			"groups.0.tags.Created":                                    "TF",
		}
	}
	var fakeDataAlicloudAlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudAlbServerGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_server_groups.default",
		existMapFunc: existDataAlicloudAlbServerGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbServerGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbServerGroupCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, statusConf, serverGroupNameConf, tagsConf, serverGroupIdsConf, resourceGroupIdConf, allConf)
}
func testAccCheckAlicloudAlbServerGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAlbServerGroup-%d"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = local.vswitch_id
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  server_group_name = var.name
  health_check_config {
    health_check_connect_port = "46325"
    health_check_enabled      = true
    health_check_host         = "tf-testAcc.com"
    health_check_codes        = ["http_2xx", "http_3xx", "http_4xx"]
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_path         = "/tf-testAcc"
    health_check_protocol     = "HTTP"
    health_check_timeout      = 5
    healthy_threshold         = 3
    unhealthy_threshold       = 3
  }
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-testAcc"
    sticky_session_type    = "Server"
  }
  tags = {
    Created = "TF"
  }
  servers {
    description = var.name
    port        = 80
    server_id   = alicloud_instance.default.id
    server_ip   = alicloud_instance.default.private_ip
    server_type = "Ecs"
    weight      = 10
  }
}

data "alicloud_alb_server_groups" "default" {	
    enable_details = true
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
