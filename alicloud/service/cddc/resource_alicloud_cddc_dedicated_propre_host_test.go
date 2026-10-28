package cddc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cddc"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass ecs_class_list -- pre-existing ForceNew TypeList (present in baseline b5d337ad); the cddc migration moves the implementation file without changing the schema or member-order semantics; the convergent-reorder case is gated behind SkipTestAcc because the dedicated propre host is a Subscription (pre-paid) resource whose production Delete logs a warning and returns nil (no-op) and Update is a no-op, so a two-apply reorder would create orphaned hosts that cannot be torn down

// lintignore: AT001
func SkipTestAccAliCloudCddcDedicatedPropreHost_basic4362(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_propre_host.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCddcDedicatedPropreHostMap4362)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cddc.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDedicatedPropreHost")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccddcdedicatedproprehost%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCddcDedicatedPropreHostBasicDependence4362)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"key_pair_name":           acctest.CHECKSET,
						"security_group_id":       acctest.CHECKSET,
						"vswitch_id":              acctest.CHECKSET,
						"ecs_zone_id":             "cn-hangzhou-i",
						"payment_type":            "Subscription",
						"ecs_deployment_set_id":   acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
						"ecs_host_name":           "testTf",
						"engine":                  "mysql",
						"dedicated_host_group_id": acctest.CHECKSET,
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
	"dedicated_host_group_id": acctest.CHECKSET,
	"ecs_instance_id":         acctest.CHECKSET,
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
// lintignore: AT001
func SkipTestAccAliCloudCddcDedicatedPropreHost_basic4363(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_propre_host.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCddcDedicatedPropreHostMap4363)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cddc.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDedicatedPropreHost")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scddcdedicatedproprehost%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCddcDedicatedPropreHostBasicDependence4363)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"security_group_id":       acctest.CHECKSET,
						"vswitch_id":              acctest.CHECKSET,
						"ecs_zone_id":             "cn-hangzhou-i",
						"payment_type":            "Subscription",
						"ecs_deployment_set_id":   acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
						"ecs_host_name":           "testTf",
						"engine":                  "mysql",
						"dedicated_host_group_id": acctest.CHECKSET,
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
	"dedicated_host_group_id": acctest.CHECKSET,
	"ecs_instance_id":         acctest.CHECKSET,
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

// SkipTestAccAliCloudCddcDedicatedPropreHost_ecsClassListOrder covers the convergent
// reorder of the dedicated propre host's ecs_class_list TypeList. The schema marks
// ecs_class_list as ForceNew, so reordering members triggers replacement; the test
// applies the original order, then a convergent apply with the members reordered.
// See scripts/collection-order/README.md for the coverage contract.
//
// Skipped (SkipTestAcc) under the same disabled/safe boundary as basic4362/basic4363:
// the dedicated propre host is a Subscription (pre-paid) resource whose production
// Delete only logs a warning and returns nil (no-op) and whose Update is a no-op, so
// the two-apply convergent reorder would create replacement hosts that cannot be torn
// down (no PreCheckPrePaidResources guard and no CheckDestroy that deletes the
// resource). ecs_class_list remains a ForceNew TypeList with unchanged semantics; this
// stock bypass records the field's original reorder behavior without enabling a paid,
// non-deletable resource as a default-runnable TestAcc.
// lintignore: AT001
func SkipTestAccAliCloudCddcDedicatedPropreHost_ecsClassListOrder(t *testing.T) {
	resourceId := "alicloud_cddc_dedicated_propre_host.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccddcdedicatedproprehost%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCddcDedicatedPropreHostBasicDependence4362)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						{
							"sys_disk_capacity":             "40",
							"system_disk_performance_level": "PL1",
							"data_disk_performance_level":   "PL1",
							"disk_count":                    "1",
							"disk_capacity":                 "40",
							"disk_type":                     "cloud_essd",
							"sys_disk_type":                 "cloud_essd",
							"instance_type":                 "ecs.g6.large",
						},
					},
					"image_id": "m-bp1d13fxs1ymbvw1dk5g",
				}),
			},
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
							"instance_type":                 "ecs.g6.large",
						},
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
			},
		},
	})
}
