package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSLaunchTemplatesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_launch_templates.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsLaunchTemplatesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ecs_launch_template.default.name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ecs_launch_template.default.name}-fake",
			"enable_details": "true",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ecs_launch_template.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ecs_launch_template.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_launch_template.default.id}"},
			"template_tags": map[string]interface{}{
				"tag1": "hello",
				"tag2": "world",
			},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_launch_template.default.id}"},
			"template_tags": map[string]interface{}{
				"tag1": "hello-fake",
				"tag2": "world-fake",
			},
			"enable_details": "true",
		}),
	}
	var existEcsLaunchTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                              "1",
			"ids.0":                                              CHECKSET,
			"names.#":                                            "1",
			"names.0":                                            name,
			"templates.#":                                        "1",
			"templates.0.auto_release_time":                      "",
			"templates.0.created_by":                             CHECKSET,
			"templates.0.data_disks.#":                           "2",
			"templates.0.default_version_number":                 CHECKSET,
			"templates.0.deployment_set_id":                      "",
			"templates.0.description":                            name,
			"templates.0.enable_vm_os_config":                    CHECKSET,
			"templates.0.host_name":                              name,
			"templates.0.image_id":                               CHECKSET,
			"templates.0.image_owner_alias":                      "",
			"templates.0.instance_charge_type":                   CHECKSET,
			"templates.0.instance_name":                          name,
			"templates.0.instance_type":                          CHECKSET,
			"templates.0.internet_charge_type":                   CHECKSET,
			"templates.0.internet_max_bandwidth_in":              CHECKSET,
			"templates.0.internet_max_bandwidth_out":             CHECKSET,
			"templates.0.io_optimized":                           "optimized",
			"templates.0.key_pair_name":                          name,
			"templates.0.latest_version_number":                  CHECKSET,
			"templates.0.id":                                     CHECKSET,
			"templates.0.launch_template_id":                     CHECKSET,
			"templates.0.network_interfaces.#":                   "1",
			"templates.0.network_interfaces.0.name":              "eth0",
			"templates.0.network_interfaces.0.description":       "hello1",
			"templates.0.network_interfaces.0.primary_ip":        "10.0.0.2",
			"templates.0.network_interfaces.0.security_group_id": "xxxx",
			"templates.0.network_interfaces.0.vswitch_id":        "xxxxxxx",
			"templates.0.network_type":                           "vpc",
			"templates.0.password_inherit":                       CHECKSET,
			"templates.0.period":                                 CHECKSET,
			"templates.0.private_ip_address":                     "",
			"templates.0.ram_role_name":                          name,
			"templates.0.spot_duration":                          "1",
			"templates.0.spot_price_limit":                       "5",
			"templates.0.spot_strategy":                          "SpotWithPriceLimit",
			"templates.0.system_disk.#":                          "1",
			"templates.0.user_data":                              CHECKSET,
			"templates.0.vswitch_id":                             CHECKSET,
			"templates.0.vpc_id":                                 CHECKSET,
			"templates.0.zone_id":                                "cn-hangzhou-i",
			"templates.0.template_tags.%":                        "2",
			"templates.0.template_tags.tag1":                     "hello",
			"templates.0.template_tags.tag2":                     "world",
		}
	}

	var fakeEcsLaunchTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"templates.#": "0",
		}
	}

	var EcsLaunchTemplatesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsLaunchTemplatesMapFunc,
		fakeMapFunc:  fakeEcsLaunchTemplatesMapFunc,
	}

	EcsLaunchTemplatesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf)
}

func dataSourceEcsLaunchTemplatesDependence(name string) string {
	return fmt.Sprintf(`
		data "alicloud_zones" "default" {
		  available_disk_category     = "cloud_efficiency"
		  available_resource_creation = "VSwitch"
		}
		data "alicloud_instance_types" "default" {
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		}
		data "alicloud_images" "default" {
		  name_regex  = "^ubuntu"
		  most_recent = true
		  owners      = "system"
		}
		data "alicloud_vpcs" "default" {
		  name_regex = "default-NODELETING"
		}
		data "alicloud_vswitches" "default" {
		 vpc_id = "${data.alicloud_vpcs.default.ids.0}"
		}
		resource "alicloud_security_group" "default" {
		  name   = "%[1]s"
		  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
		}

		resource "alicloud_ecs_launch_template" "default" {
			name  =                          "%[1]s"
			description =                    "%[1]s"
			image_id    =                    "${data.alicloud_images.default.images.0.id}"
			host_name   =                     "%[1]s"
			instance_charge_type  =          "PrePaid"
			instance_name         =           "%[1]s"
			instance_type         =           "${data.alicloud_instance_types.default.instance_types.0.id}"
			internet_charge_type  =          "PayByBandwidth"
			internet_max_bandwidth_in      =  "5"
			internet_max_bandwidth_out     =  "0"
			io_optimized                   =  "optimized"
			key_pair_name                  =  "%[1]s"
			ram_role_name                  =  "%[1]s"
			network_type                   =  "vpc"
			security_enhancement_strategy  =  "Active"
			spot_price_limit               =  "5"
			spot_strategy                  =  "SpotWithPriceLimit"
			security_group_id              =  "${alicloud_security_group.default.id}"
			system_disk {
					category =            "cloud_ssd"
					description =          "%[1]s"
					name =                 "%[1]s"
					size =                 "40"
					delete_with_instance = "false"
				}
			
			resource_group_id    =  "rg-zkdfjahg9zxncv0"
			user_data            =  "xxxxxxx"
			vswitch_id           =   "${data.alicloud_vswitches.default.vswitches.0.id}"
			vpc_id               =   "vpc-asdfnbg0as8dfk1nb2"
			zone_id              =   "cn-hangzhou-i"

			template_tags = {
				tag1 = "hello"
				tag2 = "world"
			}

			network_interfaces {
					name  =             "eth0"
					description =       "hello1"
					primary_ip  =      "10.0.0.2"
					security_group_id  = "xxxx"
					vswitch_id         = "xxxxxxx"
				}

			data_disks {
					name  =                 "disk1"
					description =          "test1"
					delete_with_instance = "true"
					category =            "cloud"
					encrypted =            "false"
					performance_level =   "PL0"
					size =                "20"
				}
			data_disks {
					name =                "disk2"
					description =          "test2"
					delete_with_instance = "true"
					category =             "cloud"
					encrypted =           "false"
					performance_level =    "PL0"
					size =                 "20"
				}
		}

`, name)
}
