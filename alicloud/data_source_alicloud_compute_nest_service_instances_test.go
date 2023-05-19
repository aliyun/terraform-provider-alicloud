package alicloud

import (
	"fmt"

	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudComputeNestServiceInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ComputeNestSupportRegions)
	resourceId := "data.alicloud_compute_nest_service_instances.default"
	name := fmt.Sprintf("tf-testacc-ComputeNestServiceInstance%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceComputeNestServiceInstancesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_compute_nest_service_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Deployed",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "DeletedFailed",
		}),
	}
	filterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"name":  "ServiceInstanceId",
					"value": []string{"${alicloud_compute_nest_service_instance.default.id}"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"name":  "ServiceInstanceId",
					"value": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
				},
			},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ServiceInstance",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ServiceInstance_Update",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_compute_nest_service_instance.default.id}"},
			"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
			"status":     "Deployed",
			"filter": []map[string]interface{}{
				{
					"name":  "ServiceInstanceId",
					"value": []string{"${alicloud_compute_nest_service_instance.default.id}"},
				},
			},
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ServiceInstance",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
			"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}_fake",
			"status":     "DeletedFailed",
			"filter": []map[string]interface{}{
				{
					"name":  "ServiceInstanceId",
					"value": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
				},
			},
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ServiceInstance_Update",
			},
		}),
	}
	var existAlicloudComputeNestServiceInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"service_instances.#":    "1",
			"service_instances.0.id": CHECKSET,
			"service_instances.0.service_instance_id":                         CHECKSET,
			"service_instances.0.service_instance_name":                       CHECKSET,
			"service_instances.0.enable_instance_ops":                         CHECKSET,
			"service_instances.0.operation_start_time":                        CHECKSET,
			"service_instances.0.operation_end_time":                          CHECKSET,
			"service_instances.0.resources":                                   CHECKSET,
			"service_instances.0.source":                                      CHECKSET,
			"service_instances.0.service.#":                                   "1",
			"service_instances.0.service.0.service_id":                        CHECKSET,
			"service_instances.0.service.0.service_type":                      CHECKSET,
			"service_instances.0.service.0.deploy_type":                       CHECKSET,
			"service_instances.0.service.0.supplier_name":                     CHECKSET,
			"service_instances.0.service.0.supplier_url":                      CHECKSET,
			"service_instances.0.service.0.publish_time":                      CHECKSET,
			"service_instances.0.service.0.version":                           "1",
			"service_instances.0.service.0.version_name":                      CHECKSET,
			"service_instances.0.service.0.service_infos.#":                   "1",
			"service_instances.0.service.0.service_infos.0.name":              CHECKSET,
			"service_instances.0.service.0.service_infos.0.short_description": CHECKSET,
			"service_instances.0.service.0.service_infos.0.image":             CHECKSET,
			"service_instances.0.service.0.service_infos.0.locale":            CHECKSET,
			"service_instances.0.service.0.status":                            CHECKSET,
			"service_instances.0.tags.%":                                      "2",
			"service_instances.0.tags.Created":                                "TF",
			"service_instances.0.tags.For":                                    "ServiceInstance",
			"service_instances.0.status":                                      "Deployed",
		}
	}
	var fakeAlicloudComputeNestServiceInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "0",
			"names.#":             "0",
			"service_instances.#": "0",
		}
	}
	var alicloudComputeNestServiceInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_compute_nest_service_instances.default",
		existMapFunc: existAlicloudComputeNestServiceInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudComputeNestServiceInstancesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudComputeNestServiceInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, filterConf, tagsConf, allConf)
}

func dataSourceComputeNestServiceInstancesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_compute_nest_service_instance" "default" {
  		service_id            = "service-dd475e6e468348799f0f"
  		service_version       = "1"
  		service_instance_name = var.name
  		resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		payment_type          = "Permanent"
  		operation_metadata {
    		operation_start_time = "1681281179000"
    		operation_end_time   = "1681367579000"
    		resources            = "{\"Type\":\"ResourceIds\",\"ResourceIds\":{\"ALIYUN::ECS::INSTANCE\":[\"${alicloud_instance.default.id}\"]},\"RegionId\":\"cn-hangzhou\"}"
  		}
  		tags = {
    		Created = "TF"
    		For     = "ServiceInstance"
  		}
	}
`, name)
}
