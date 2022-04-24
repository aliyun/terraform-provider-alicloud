package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCTrafficMirrorSessionsDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_traffic_mirror_sessions.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorsession-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcTrafficMirrorSessionsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_traffic_mirror_session.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"name_regex": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"name_regex": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}-fake",
		}),
	}
	enabledConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"enabled": "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"enabled": "true",
		}),
	}
	priorityConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"priority": "11",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"priority": "2",
		}),
	}
	filterIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_filter_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_filter_id}-fake",
		}),
	}
	sessionNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                         []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_session_name": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                         []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_session_name": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}-fake",
		}),
	}
	sourceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_source_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_source_ids[0]}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_source_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_source_ids[0]}-fake",
		}),
	}
	targetIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_target_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_target_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"traffic_mirror_target_id": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_target_id}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"status": "Created",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"status": "Deleting",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}",
			"ids":                         []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"enabled":                     "false",
			"priority":                    "11",
			"traffic_mirror_filter_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_filter_id}",
			"traffic_mirror_session_name": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}",
			"traffic_mirror_source_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_source_ids[0]}",
			"traffic_mirror_target_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_target_id}",
			"status":                      "Created",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}-fake",
			"ids":                         []string{"${alicloud_vpc_traffic_mirror_session.default.id}"},
			"enabled":                     "true",
			"priority":                    "2",
			"traffic_mirror_filter_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_filter_id}-fake",
			"traffic_mirror_session_name": "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_session_name}-fake",
			"traffic_mirror_source_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_source_ids[0]}-fake",
			"traffic_mirror_target_id":    "${alicloud_vpc_traffic_mirror_session.default.traffic_mirror_target_id}-fake",
			"status":                      "Deleting",
		}),
	}
	var existVpcTrafficMirrorSessionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"sessions.#":                             "1",
			"sessions.0.traffic_mirror_session_id":   CHECKSET,
			"sessions.0.traffic_mirror_session_name": fmt.Sprintf("tf-testacc-vpctrafficmirrorsession-%d", rand),
			"sessions.0.traffic_mirror_session_description":     fmt.Sprintf("tf-testacc-vpctrafficmirrorsession-%d", rand),
			"sessions.0.status":                                 "Created",
			"sessions.0.traffic_mirror_target_type":             "NetworkInterface",
			"sessions.0.priority":                               "11",
			"sessions.0.traffic_mirror_target_id":               CHECKSET,
			"sessions.0.traffic_mirror_source_ids.#":            "1",
			"sessions.0.traffic_mirror_filter_id":               CHECKSET,
			"sessions.0.enabled":                                "false",
			"sessions.0.packet_length":                          "1500",
			"sessions.0.traffic_mirror_session_business_status": "Normal",
			"sessions.0.virtual_network_id":                     "10",
		}
	}

	var fakeVpcTrafficMirrorSessionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"sessions.#": "0",
		}
	}

	var vpcTrafficMirrorSessionCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcTrafficMirrorSessionMapFunc,
		fakeMapFunc:  fakeVpcTrafficMirrorSessionMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	vpcTrafficMirrorSessionCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, enabledConf, priorityConf, filterIdConf, sessionNameConf, sourceIdConf, targetIdConf, statusConf, allConf)
}

func dataSourceVpcTrafficMirrorSessionsDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_name        = var.name
  host_name            = var.name
  image_id             = data.alicloud_images.default.images.0.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  security_groups      = [alicloud_security_group.default.id]
  vswitch_id           = local.vswitch_id
  system_disk_category = "cloud_essd"
}

resource "alicloud_ecs_network_interface" "default" {
  count                  = 2
  network_interface_name = var.name
  vswitch_id             = local.vswitch_id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  count                = 2
  instance_id          = alicloud_instance.default.id
  network_interface_id = element(alicloud_ecs_network_interface.default.*.id, count.index)
}

resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_name        = var.name
  traffic_mirror_filter_description = var.name
}

resource "alicloud_vpc_traffic_mirror_session" "default" {
  priority                           = 11
  virtual_network_id                 = 10
  traffic_mirror_session_description = var.name
  traffic_mirror_session_name        = var.name
  traffic_mirror_target_id           = alicloud_ecs_network_interface_attachment.default[0].network_interface_id
  traffic_mirror_source_ids          = [alicloud_ecs_network_interface_attachment.default[1].network_interface_id]
  traffic_mirror_filter_id           = alicloud_vpc_traffic_mirror_filter.default.id
  traffic_mirror_target_type         = "NetworkInterface"
}`, name)
}
