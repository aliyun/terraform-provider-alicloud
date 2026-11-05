package computenest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/computenest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudComputeNestServiceInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, ServiceInstanceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return computenest.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testaccComputeNestServiceInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ServiceInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// ComputeNest is Hangzhou-only (connectivity.ComputeNestSupportRegions);
			// the flat suite gate checkoutSupportedRegions converged to cn-hangzhou,
			// so pin it here to keep config and execution region consistent.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":             "service-0a15dbf13e4049daaf66",
					"service_version":        "1",
					"service_instance_name":  name,
					"parameters":             `{ \"ZoneId\": \"` + "${data.alicloud_zones.default.zones.0.id}" + `\", \"SystemDiskSize\": 40, \"InstanceChargeType\": \"PostPaid\", \"SecurityGroupId\": \"` + "${alicloud_security_group.default.id}" + `\", \"VSwitchId\": \"` + "${alicloud_vswitch.default.id}" + `\", \"Count\": 2, \"UserEnablePrometheus\": true, \"SystemDiskCategory\": \"cloud_efficiency\", \"InternetChargeType\": \"PayByTraffic\", \"InternetMaxBandwidthOut\": 0, \"VpcId\": \"` + "${alicloud_vpc.default.id}" + `\", \"RegionId\": \"cn-hangzhou\", \"DataDiskSize\": 100, \"DataDiskCategory\": \"cloud_efficiency\", \"InstanceType\": \"` + "${data.alicloud_instance_types.default.instance_types.0.id}" + `\", \"Password\": \"YourPassword123!\"}`,
					"enable_instance_ops":    "false",
					"template_name":          "模板1",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"specification_name":     "套餐一",
					"enable_user_prometheus": "true",
					"commodity": []map[string]interface{}{
						{
							"pay_period":      1,
							"pay_period_unit": "Month",
						},
					},
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
						"resource_group_id":      acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
				// commodity is create-only: resourceServiceInstanceCreate writes
				// Commodity into the CreateServiceInstance request, but
				// resourceServiceInstanceUpdate never handles d.HasChange("commodity")
				// and resourceServiceInstanceRead never calls d.Set("commodity", ...),
				// so the imported state cannot carry it. Listing both the block
				// name and the nested leaf names (commodity.pay_period /
				// commodity.pay_period_unit) is required because the coverage
				// checker removes schemaModifySet members by exact dotted-string
				// match; the bare block name does not exempt its nested fields.
				ImportStateVerifyIgnore: []string{"parameters", "commodity", "commodity.pay_period", "commodity.pay_period_unit"},
			},
		},
	})
}

func TestAccAliCloudComputeNestServiceInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, ServiceInstanceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return computenest.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testaccComputeNestServiceInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ServiceInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":     acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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

func TestAccAliCloudComputeNestServiceInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_compute_nest_service_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, ServiceInstanceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return computenest.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testaccComputeNestServiceInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ServiceInstanceBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":     acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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

var ServiceInstanceMap = map[string]string{
	"status": acctest.CHECKSET,
}

func ServiceInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.g6"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

func ServiceInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.g6"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
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
  		vswitch_id                 = alicloud_vswitch.default.id
	}
`, name)
}

func ServiceInstanceBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.g6"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_compute_nest_service_instance" "operated" {
  		service_id             = "service-0a15dbf13e4049daaf66"
  		service_version        = "1"
  		service_instance_name  = "${var.name}-operated"
  		parameters             = "{ \"ZoneId\": \"`+"${data.alicloud_zones.default.zones.0.id}"+`\", \"SystemDiskSize\": 40, \"InstanceChargeType\": \"PostPaid\", \"SecurityGroupId\": \"`+"${alicloud_security_group.default.id}"+`\", \"VSwitchId\": \"`+"${alicloud_vswitch.default.id}"+`\", \"Count\": 2, \"UserEnablePrometheus\": true, \"SystemDiskCategory\": \"cloud_efficiency\", \"InternetChargeType\": \"PayByTraffic\", \"InternetMaxBandwidthOut\": 0, \"VpcId\": \"`+"${alicloud_vpc.default.id}"+`\", \"RegionId\": \"cn-hangzhou\", \"DataDiskSize\": 100, \"DataDiskCategory\": \"cloud_efficiency\", \"InstanceType\": \"`+"${data.alicloud_instance_types.default.instance_types.0.id}"+`\", \"Password\": \"YourPassword123!\"}"
  		enable_instance_ops    = false
  		template_name          = "模板1"
  		resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		specification_name     = "套餐一"
  		enable_user_prometheus = true
	}
`, name)
}
