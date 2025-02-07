package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAliCloudCddcDedicatedPropreHost_basic4362(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_propre_host.default"
	ra := resourceAttrInit(resourceId, AlicloudCddcDedicatedPropreHostMap4362)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedPropreHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccddcdedicatedproprehost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCddcDedicatedPropreHostBasicDependence4362)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_instance_name":       "testTf",
					"key_pair_name":           "${local.alicloud_key_pair_id}",
					"security_group_id":       "${local.alicloud_security_group_id}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids.0}",
					"ecs_zone_id":             "cn-hangzhou-i",
					"payment_type":            "Subscription",
					"ecs_deployment_set_id":   "${local.alicloud_ecs_deployment_set_id}",
					"vpc_id":                  "${local.vpc_id}",
					"ecs_host_name":           "testTf",
					"engine":                  "mysql",
					"dedicated_host_group_id": "${local.dedicated_host_group_id}",
					"ecs_class_list": []map[string]interface{}{
						{
							"sys_disk_capacity":             "40",
							"system_disk_performance_level": "PL1",
							"data_disk_performance_level":   "PL1",
							"disk_count":                    "1",
							"disk_capacity":                 "40",
							"disk_type":                     "cloud_essd",
							"sys_disk_type":                 "cloud_essd",
							"instance_type":                 "ecs.c6a.large",
						},
					},
					"image_id": "m-bp1d13fxs1ymbvw1dk5g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_instance_name":       "testTf",
						"key_pair_name":           CHECKSET,
						"security_group_id":       CHECKSET,
						"vswitch_id":              CHECKSET,
						"ecs_zone_id":             "cn-hangzhou-i",
						"payment_type":            "Subscription",
						"ecs_deployment_set_id":   CHECKSET,
						"vpc_id":                  CHECKSET,
						"ecs_host_name":           "testTf",
						"engine":                  "mysql",
						"dedicated_host_group_id": CHECKSET,
						"ecs_class_list.#":        "1",
						"image_id":                "m-bp1d13fxs1ymbvw1dk5g",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "ecs_unique_suffix", "os_password", "password_inherit", "period", "period_type"},
			},
		},
	})
}

var AlicloudCddcDedicatedPropreHostMap4362 = map[string]string{
	"dedicated_host_group_id": CHECKSET,
	"ecs_instance_id":         CHECKSET,
}

func AlicloudCddcDedicatedPropreHostBasicDependence4362(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6e" 
  network_type = "Vpc"
}

data "alicloud_images" "default" {
  name_regex  = "^aliyun_3_x64_20G_scc*"
  owners      = "system"
}

data "alicloud_instance_types" "essd" {
 	cpu_core_count    = 2
	memory_size       = 4
 	system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.groups.0.vpc_id : data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_security_groups" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_security_group" "default" {
  count = length(data.alicloud_security_groups.default.ids) > 0 ? 0 : 1
  vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
  name   = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_ecs_deployment_sets" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_ecs_deployment_set" "default" {
  count = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? 0 : 1
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "tf-testacc-cddc_dedicated_propre_host"
  description         = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_key_pairs" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_key_pair" "default" {
    count = length(data.alicloud_key_pairs.default.ids) > 0 ? 0 : 1
	key_pair_name = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  engine     = "MySQL"
  name_regex     = "^NO-DELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	count = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
	engine = "MySQL"
	vpc_id = local.vpc_id
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = "tf-testacc-cddc_dedicated_propre_host"
	open_permission = true
}
locals {
    vpc_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.groups.0.vpc_id : data.alicloud_vpcs.default.ids.0
    alicloud_security_group_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.ids.0 : concat(alicloud_security_group.default[*].id, [""])[0]
    alicloud_ecs_deployment_set_id = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? data.alicloud_ecs_deployment_sets.default.sets.0.deployment_set_id : concat(alicloud_ecs_deployment_set.default[*].id, [""])[0]
    alicloud_key_pair_id = length(data.alicloud_key_pairs.default.ids) > 0 ? data.alicloud_key_pairs.default.ids.0 : concat(alicloud_key_pair.default[*].id, [""])[0]
	dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : concat(alicloud_cddc_dedicated_host_group.default[*].id, [""])[0]
}

`, name)
}

// Case 4363
func SkipTestAccAliCloudCddcDedicatedPropreHost_basic4363(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_propre_host.default"
	ra := resourceAttrInit(resourceId, AlicloudCddcDedicatedPropreHostMap4363)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedPropreHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedproprehost%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCddcDedicatedPropreHostBasicDependence4363)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"period_type":                "Monthly",
					"auto_renew":                 "false",
					"period":                     "1",
					"ecs_unique_suffix":          "false",
					"password_inherit":           "false",
					"ecs_instance_name":          "testTf",
					"security_group_id":          "${local.alicloud_security_group_id}",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"ecs_zone_id":                "cn-hangzhou-i",
					"payment_type":               "Subscription",
					"ecs_deployment_set_id":      "${local.alicloud_ecs_deployment_set_id}",
					"vpc_id":                     "${local.vpc_id}",
					"ecs_host_name":              "testTf",
					"engine":                     "mysql",
					"dedicated_host_group_id":    "${local.dedicated_host_group_id}",
					"auto_pay":                   "true",
					"user_data_encoded":          "true",
					"user_data":                  "aGVsbG8gd29ybGQ=",
					"internet_charge_type":       "PayByBandwidth",
					"internet_max_bandwidth_out": "1",
					"resource_group_id":          "${local.resource_group_id}",
					"ecs_class_list": []map[string]interface{}{
						{
							"sys_disk_capacity":             "40",
							"system_disk_performance_level": "PL1",
							"data_disk_performance_level":   "PL1",
							"disk_count":                    "1",
							"disk_capacity":                 "40",
							"disk_type":                     "cloud_essd",
							"sys_disk_type":                 "cloud_essd",
							"instance_type":                 "ecs.c6a.large",
						},
					},
					"os_password": "YourPassword123!",
					"image_id":    "m-bp1d13fxs1ymbvw1dk5g",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period_type":             "Monthly",
						"auto_renew":              "false",
						"period":                  "1",
						"ecs_unique_suffix":       "false",
						"password_inherit":        "false",
						"ecs_instance_name":       "testTf",
						"security_group_id":       CHECKSET,
						"vswitch_id":              CHECKSET,
						"ecs_zone_id":             "cn-hangzhou-i",
						"payment_type":            "Subscription",
						"ecs_deployment_set_id":   CHECKSET,
						"vpc_id":                  CHECKSET,
						"ecs_host_name":           "testTf",
						"engine":                  "mysql",
						"dedicated_host_group_id": CHECKSET,
						"ecs_class_list.#":        "1",
						"os_password":             "YourPassword123!",
						"image_id":                "m-bp1d13fxs1ymbvw1dk5g",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "ecs_unique_suffix", "internet_charge_type", "internet_max_bandwidth_out", "os_password", "password_inherit", "period", "period_type", "user_data", "user_data_encoded"},
			},
		},
	})
}

var AlicloudCddcDedicatedPropreHostMap4363 = map[string]string{
	"dedicated_host_group_id": CHECKSET,
	"ecs_instance_id":         CHECKSET,
}

func AlicloudCddcDedicatedPropreHostBasicDependence4363(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6e" 
  network_type = "Vpc"
}

data "alicloud_images" "default" {
  name_regex  = "^aliyun_3_x64_20G_scc*"
  owners      = "system"
}

data "alicloud_instance_types" "essd" {
 	cpu_core_count    = 2
	memory_size       = 4
 	system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.groups.0.vpc_id : data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_security_groups" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_security_group" "default" {
  count = length(data.alicloud_security_groups.default.ids) > 0 ? 0 : 1
  vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
  name   = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_ecs_deployment_sets" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_ecs_deployment_set" "default" {
  count = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? 0 : 1
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "tf-testacc-cddc_dedicated_propre_host"
  description         = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_key_pairs" "default" {
  name_regex     = "tf-testacc-cddc_dedicated_propre_host"
}

resource "alicloud_key_pair" "default" {
    count = length(data.alicloud_key_pairs.default.ids) > 0 ? 0 : 1
	key_pair_name = "tf-testacc-cddc_dedicated_propre_host"
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  engine     = "MySQL"
  name_regex     = "^NO-DELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	count = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
	engine = "MySQL"
	vpc_id = local.vpc_id
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = "tf-testacc-cddc_dedicated_propre_host"
	open_permission = true
}
locals {
    vpc_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.groups.0.vpc_id : data.alicloud_vpcs.default.ids.0
    resource_group_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.groups.0.resource_group_id : data.alicloud_vpcs.default.vpcs.0.resource_group_id
    alicloud_security_group_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.ids.0 : concat(alicloud_security_group.default[*].id, [""])[0]
    alicloud_ecs_deployment_set_id = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? data.alicloud_ecs_deployment_sets.default.sets.0.deployment_set_id : concat(alicloud_ecs_deployment_set.default[*].id, [""])[0]
    alicloud_key_pair_id = length(data.alicloud_key_pairs.default.ids) > 0 ? data.alicloud_key_pairs.default.ids.0 : concat(alicloud_key_pair.default[*].id, [""])[0]
	dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : concat(alicloud_cddc_dedicated_host_group.default[*].id, [""])[0]
}

`, name)
}
