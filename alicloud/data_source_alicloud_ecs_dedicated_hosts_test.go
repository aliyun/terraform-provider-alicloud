package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudECSDedicatedHostsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "data.alicloud_ecs_dedicated_hosts.default"
	name := fmt.Sprintf("tf_testAccEcsDedicatedHostsDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsDedicatedHostsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_instance.default.dedicated_host_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_instance.default.dedicated_host_id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_instance.default.dedicated_host_id}"},
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_instance.default.dedicated_host_id}"},
			"name_regex": name + "fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_instance.default.dedicated_host_id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_instance.default.dedicated_host_id}"},
			"status": "UnderAssessment",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_instance.default.dedicated_host_id}"},
			"dedicated_host_type": "ddh.c5",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_instance.default.dedicated_host_id}"},
			"dedicated_host_type": "ddh.g5",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_instance.default.dedicated_host_id}"},
			"tags": map[string]string{
				"Create": "TF",
				"For":    "ddh-test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_instance.default.dedicated_host_id}"},
			"tags": map[string]string{
				"Create": "ddh-test",
				"For":    "TF",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_instance.default.dedicated_host_id}"},
			"dedicated_host_type": "ddh.c5",
			"status":              "Available",
			"name_regex":          name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alicloud_instance.default.dedicated_host_id}"},
			"dedicated_host_type": "ddh.c5",
			"name_regex":          name + "fake",
			"status":              "UnderAssessment",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"ids.0":                                 CHECKSET,
			"names.#":                               "1",
			"names.0":                               CHECKSET,
			"hosts.0.action_on_maintenance":         "Migrate",
			"hosts.0.auto_placement":                "on",
			"hosts.0.auto_release_time":             "",
			"hosts.0.id":                            CHECKSET,
			"hosts.0.dedicated_host_id":             CHECKSET,
			"hosts.0.dedicated_host_name":           CHECKSET,
			"hosts.0.dedicated_host_type":           CHECKSET,
			"hosts.0.description":                   "From_Terraform",
			"hosts.0.expired_time":                  CHECKSET,
			"hosts.0.gpu_spec":                      "",
			"hosts.0.machine_id":                    CHECKSET,
			"hosts.0.payment_type":                  "PostPaid",
			"hosts.0.physical_gpus":                 CHECKSET,
			"hosts.0.resource_group_id":             "",
			"hosts.0.sale_cycle":                    "",
			"hosts.0.sockets":                       CHECKSET,
			"hosts.0.status":                        "Available",
			"hosts.0.supported_instance_types_list": NOSET,
			"hosts.0.zone_id":                       CHECKSET,
			"hosts.0.instances.#":                   "1",
			"hosts.0.instances.0.instance_id":       CHECKSET,
			"hosts.0.instances.0.instance_type":     CHECKSET,
			"hosts.0.instances.0.socket_id":         "",
			"hosts.0.instances.0.instance_owner_id": CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"hosts.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, typeConf, statusConf, tagsConf, allConf)
}

func dataSourceEcsDedicatedHostsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		instance_type_family = "ecs.c5"
  		availability_zone    = alicloud_ecs_dedicated_host.default.zone_id
  		image_id             = data.alicloud_images.default.images.0.id
  		system_disk_category = "cloud_efficiency"
	}

	resource "alicloud_vpc" "default" {
  		cidr_block = "192.168.0.0/16"
  		vpc_name   = var.name
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  		zone_id      = alicloud_ecs_dedicated_host.default.zone_id
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_ecs_dedicated_host" "default" {
  		dedicated_host_type   = "ddh.c5"
  		description           = "From_Terraform"
  		dedicated_host_name   = var.name
  		action_on_maintenance = "Migrate"
  		tags = {
    		Create = "TF"
    		For    = "ddh-test",
  		}
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
  		dedicated_host_id          = alicloud_ecs_dedicated_host.default.id
	}
	`, name)
}
