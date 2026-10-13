package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Arms Environment. >>> Resource test cases, automatically generated.
// Case 4280
func TestAccAliCloudArmsEnvironment_basic4280(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4280)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "ECS",
					"environment_sub_type": "ECS",
					"environment_name":     name,
					"bind_resource_id":     "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "ECS",
						"environment_sub_type": "ECS",
						"environment_name":     name,
						"bind_resource_id":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "ECS",
					"environment_name":     name + "_update",
					"bind_resource_id":     "${data.alicloud_vpcs.default.ids.0}",
					"environment_sub_type": "ECS",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "ECS",
						"environment_name":     name + "_update",
						"bind_resource_id":     CHECKSET,
						"environment_sub_type": "ECS",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

var AlicloudArmsEnvironmentMap4280 = map[string]string{
	"environment_id": CHECKSET,
}

func AlicloudArmsEnvironmentBasicDependence4280(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}


`, name)
}

// Case 4697
func TestAccAliCloudArmsEnvironment_basic4697(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4697)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4697)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "CS",
					"bind_resource_id":     "${local.cluster_id}",
					"environment_sub_type": "ManagedKubernetes",
					"environment_name":     name,
					"managed_type":         "agent",
					"drop_metrics":         "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "CS",
						"environment_name":     name,
						"bind_resource_id":     CHECKSET,
						"environment_sub_type": "ManagedKubernetes",
						"managed_type":         "agent",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_name": name + "_update",
					"drop_metrics":     "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "CS",
					"environment_name":     name + "_update",
					"bind_resource_id":     "${local.cluster_id}",
					"environment_sub_type": "ManagedKubernetes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "CS",
						"environment_name":     name + "_update",
						"bind_resource_id":     CHECKSET,
						"environment_sub_type": "ManagedKubernetes",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang", "drop_metrics"},
			},
		},
	})
}

var AlicloudArmsEnvironmentMap4697 = map[string]string{
	"environment_id": CHECKSET,
}

func AlicloudArmsEnvironmentBasicDependence4697(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [data.alicloud_vswitches.default.ids.0]
  new_nat_gateway      = false
  pod_cidr             = "10.128.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}


`, name)
}

// Case 4543
func TestAccAliCloudArmsEnvironment_basic4543(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4543)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4543)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "Cloud",
					"environment_sub_type": "Cloud",
					"environment_name":     name,
					"bind_resource_id":     "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "Cloud",
						"environment_sub_type": "Cloud",
						"environment_name":     name,
						"bind_resource_id":     "cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "Cloud",
					"environment_name":     name + "_update",
					"environment_sub_type": "Cloud",
					"bind_resource_id":     "cn-hangzhou",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"aliyun_lang":          "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "Cloud",
						"environment_name":     name + "_update",
						"environment_sub_type": "Cloud",
						"bind_resource_id":     "cn-hangzhou",
						"resource_group_id":    CHECKSET,
						"aliyun_lang":          "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

var AlicloudArmsEnvironmentMap4543 = map[string]string{
	"environment_id": CHECKSET,
}

func AlicloudArmsEnvironmentBasicDependence4543(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}


`, name)
}

// Case 4280  twin
func TestAccAliCloudArmsEnvironment_basic4280_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4280)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "ECS",
					"environment_name":     name,
					"bind_resource_id":     "${data.alicloud_vpcs.default.ids.0}",
					"environment_sub_type": "ECS",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "ECS",
						"environment_name":     name,
						"bind_resource_id":     CHECKSET,
						"environment_sub_type": "ECS",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Case 4697  twin
func TestAccAliCloudArmsEnvironment_basic4697_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4697)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4697)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "CS",
					"environment_name":     name,
					"bind_resource_id":     "${local.cluster_id}",
					"environment_sub_type": "ManagedKubernetes",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "CS",
						"environment_name":     name,
						"bind_resource_id":     CHECKSET,
						"environment_sub_type": "ManagedKubernetes",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Case 4543  twin
func TestAccAliCloudArmsEnvironment_basic4543_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_environment.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvironmentMap4543)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvironment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvironmentBasicDependence4543)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_type":     "Cloud",
					"environment_name":     name,
					"environment_sub_type": "Cloud",
					"bind_resource_id":     "cn-hangzhou",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"aliyun_lang":          "zh",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_type":     "Cloud",
						"environment_name":     CHECKSET,
						"environment_sub_type": "Cloud",
						"bind_resource_id":     "cn-hangzhou",
						"resource_group_id":    CHECKSET,
						"aliyun_lang":          "zh",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Test Arms Environment. <<< Resource test cases, automatically generated.
