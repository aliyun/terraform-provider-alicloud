package nlb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudNlbServerGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_server_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_server_group.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_server_group.default.id}"]`,
			"resource_group_id": `"${alicloud_nlb_server_group.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_server_group.default.id}"]`,
			"resource_group_id": `"${alicloud_nlb_server_group.default.resource_group_id}_fake"`,
		}),
	}
	serverGroupNamesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_names": `["${alicloud_nlb_server_group.default.server_group_name}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_names": `["${alicloud_nlb_server_group.default.server_group_name}_fake"]`,
		}),
	}
	serverGroupTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_server_group.default.id}"]`,
			"server_group_type": `"${alicloud_nlb_server_group.default.server_group_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_nlb_server_group.default.id}"]`,
			"server_group_type": `"Ip"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_nlb_server_group.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_nlb_server_group.default.id}"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_server_group.default.server_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_server_group.default.server_group_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nlb_server_group.default.id}"]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nlb_server_group.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_nlb_server_group.default.id}"]`,
			"name_regex":         `"${alicloud_nlb_server_group.default.server_group_name}"`,
			"resource_group_id":  `"${alicloud_nlb_server_group.default.resource_group_id}"`,
			"server_group_names": `["${alicloud_nlb_server_group.default.server_group_name}"]`,
			"server_group_type":  `"${alicloud_nlb_server_group.default.server_group_type}"`,
			"status":             `"Available"`,
			"tags":               `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_nlb_server_group.default.id}_fake"]`,
			"name_regex":         `"${alicloud_nlb_server_group.default.server_group_name}_fake"`,
			"resource_group_id":  `"${alicloud_nlb_server_group.default.resource_group_id}_fake"`,
			"server_group_names": `["${alicloud_nlb_server_group.default.server_group_name}_fake"]`,
			"server_group_type":  `"Ip"`,
			"status":             `"Configuring"`,
			"tags":               `{Created = "TF_fake"}`,
		}),
	}
	var existAlicloudNlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                "1",
			"names.#":                                              "1",
			"groups.#":                                             "1",
			"groups.0.address_ip_version":                          "Ipv4",
			"groups.0.connection_drain":                            "true",
			"groups.0.connection_drain_timeout":                    "60",
			"groups.0.protocol":                                    "TCP",
			"groups.0.resource_group_id":                           acctest.CHECKSET,
			"groups.0.scheduler":                                   "Wrr",
			"groups.0.server_group_name":                           acctest.CHECKSET,
			"groups.0.server_group_type":                           "Instance",
			"groups.0.tags.%":                                      "1",
			"groups.0.tags.Created":                                "TF",
			"groups.0.vpc_id":                                      acctest.CHECKSET,
			"groups.0.id":                                          acctest.CHECKSET,
			"groups.0.related_load_balancer_ids.#":                 acctest.CHECKSET,
			"groups.0.server_count":                                acctest.CHECKSET,
			"groups.0.status":                                      acctest.CHECKSET,
			"groups.0.preserve_client_ip_enabled":                  "true",
			"groups.0.health_check.#":                              "1",
			"groups.0.health_check.0.health_check_enabled":         "true",
			"groups.0.health_check.0.health_check_type":            "TCP",
			"groups.0.health_check.0.health_check_connect_port":    "0",
			"groups.0.health_check.0.healthy_threshold":            "2",
			"groups.0.health_check.0.unhealthy_threshold":          "2",
			"groups.0.health_check.0.health_check_connect_timeout": "5",
			"groups.0.health_check.0.health_check_interval":        "10",
			"groups.0.health_check.0.http_check_method":            "GET",
			"groups.0.health_check.0.health_check_url":             "/test/index.html",
			"groups.0.health_check.0.health_check_domain":          "tf-testAcc.com",
			"groups.0.health_check.0.health_check_http_code.#":     "3",
		}
	}
	var fakeAlicloudNlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNlbServerGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nlb_server_groups.default",
		ExistMapFunc: existAlicloudNlbServerGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudNlbServerGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudNlbServerGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, serverGroupNamesConf, serverGroupTypeConf, tagsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudNlbServerGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccServerGroup-%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Instance"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
	health_check_url =           "/test/index.html"
	health_check_domain =       "tf-testAcc.com"
    health_check_enabled         = true
    health_check_type            = "TCP"
    health_check_connect_port    = 0
    healthy_threshold            = 2
    unhealthy_threshold          = 2
    health_check_connect_timeout = 5
    health_check_interval        = 10
    http_check_method            = "GET"
    health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
  }
  connection_drain           = true
  connection_drain_timeout   = 60
  preserve_client_ip_enabled = true
  tags = {
    Created = "TF"
  }
  address_ip_version = "Ipv4"
}

data "alicloud_nlb_server_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
