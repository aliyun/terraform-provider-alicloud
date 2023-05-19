package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudComputeNestServiceInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ComputeNestSupportRegions)
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudComputeNestServiceInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ComputeNestService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeComputeNestServiceInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%sComputeNestServiceInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudComputeNestServiceInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":             "service-0a15dbf13e4049daaf66",
					"service_version":        "1",
					"service_instance_name":  name,
					"parameters":             `{ \"ZoneId\": \"` + "${data.alicloud_zones.default.zones.0.id}" + `\", \"SystemDiskSize\": 40, \"InstanceChargeType\": \"PostPaid\", \"SecurityGroupId\": \"` + "${alicloud_security_group.default.id}" + `\", \"VSwitchId\": \"` + "${data.alicloud_vswitches.default.ids.0}" + `\", \"Count\": 2, \"UserEnablePrometheus\": true, \"SystemDiskCategory\": \"cloud_efficiency\", \"InternetChargeType\": \"PayByTraffic\", \"InternetMaxBandwidthOut\": 0, \"VpcId\": \"` + "${data.alicloud_vpcs.default.ids.0}" + `\", \"RegionId\": \"` + defaultRegionToTest + `\", \"DataDiskSize\": 100, \"DataDiskCategory\": \"cloud_efficiency\", \"InstanceType\": \"` + "${data.alicloud_instance_types.default.instance_types.0.id}" + `\", \"Password\": \"YourPassword123!\"}`,
					"enable_instance_ops":    "false",
					"template_name":          "模板1",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"specification_name":     "套餐一",
					"enable_user_prometheus": "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":             "service-0a15dbf13e4049daaf66",
						"service_version":        "1",
						"service_instance_name":  name,
						"enable_instance_ops":    "false",
						"template_name":          "模板1",
						"resource_group_id":      CHECKSET,
						"specification_name":     "套餐一",
						"enable_user_prometheus": "true",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "ServiceInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "ServiceInstance_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "ServiceInstance_Update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

func TestAccAlicloudComputeNestServiceInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ComputeNestSupportRegions)
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudComputeNestServiceInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ComputeNestService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeComputeNestServiceInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%sComputeNestServiceInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudComputeNestServiceInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":            "service-dd475e6e468348799f0f",
					"service_version":       "1",
					"service_instance_name": name,
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"payment_type":          "Permanent",
					"operation_metadata": []map[string]interface{}{
						{
							"operation_start_time": "1681281179000",
							"operation_end_time":   "1681367579000",
							"resources":            `{\"Type\":\"ResourceIds\",\"ResourceIds\":{\"ALIYUN::ECS::INSTANCE\":[\"` + "${alicloud_instance.default.id}" + `\"]},\"RegionId\":\"cn-hangzhou\"}`,
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":            "service-dd475e6e468348799f0f",
						"service_version":       "1",
						"service_instance_name": name,
						"resource_group_id":     CHECKSET,
						"payment_type":          "Permanent",
						"operation_metadata.#":  "1",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "ServiceInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "ServiceInstance_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "ServiceInstance_Update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudComputeNestServiceInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ComputeNestSupportRegions)
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudComputeNestServiceInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ComputeNestService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeComputeNestServiceInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%sComputeNestServiceInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudComputeNestServiceInstanceBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":            "service-a5572d9808b949c1a9fd",
					"service_version":       "1",
					"service_instance_name": name,
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"payment_type":          "Permanent",
					"operation_metadata": []map[string]interface{}{
						{
							"operation_start_time":         "1681281179000",
							"operation_end_time":           "1681367579000",
							"operated_service_instance_id": "${alicloud_compute_nest_service_instance.operated.id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":            "service-a5572d9808b949c1a9fd",
						"service_version":       "1",
						"service_instance_name": name,
						"resource_group_id":     CHECKSET,
						"payment_type":          "Permanent",
						"operation_metadata.#":  "1",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "ServiceInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "ServiceInstance_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "ServiceInstance_Update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudComputeNestServiceInstanceMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudComputeNestServiceInstanceBasicDependence0(name string) string {
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
`, name)
}

func AlicloudComputeNestServiceInstanceBasicDependence1(name string) string {
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
`, name)
}

func AlicloudComputeNestServiceInstanceBasicDependence2(name string) string {
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

	resource "alicloud_compute_nest_service_instance" "operated" {
  		service_id             = "service-0a15dbf13e4049daaf66"
  		service_version        = "1"
  		service_instance_name  = "${var.name}-operated"
  		parameters             = "{ \"ZoneId\": \"`+"${data.alicloud_zones.default.zones.0.id}"+`\", \"SystemDiskSize\": 40, \"InstanceChargeType\": \"PostPaid\", \"SecurityGroupId\": \"`+"${alicloud_security_group.default.id}"+`\", \"VSwitchId\": \"`+"${data.alicloud_vswitches.default.ids.0}"+`\", \"Count\": 2, \"UserEnablePrometheus\": true, \"SystemDiskCategory\": \"cloud_efficiency\", \"InternetChargeType\": \"PayByTraffic\", \"InternetMaxBandwidthOut\": 0, \"VpcId\": \"`+"${data.alicloud_vpcs.default.ids.0}"+`\", \"RegionId\": \"`+defaultRegionToTest+`\", \"DataDiskSize\": 100, \"DataDiskCategory\": \"cloud_efficiency\", \"InstanceType\": \"`+"${data.alicloud_instance_types.default.instance_types.0.id}"+`\", \"Password\": \"YourPassword123!\"}"
  		enable_instance_ops    = false
  		template_name          = "模板1"
  		resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		specification_name     = "套餐一"
  		enable_user_prometheus = true
	}
`, name)
}
