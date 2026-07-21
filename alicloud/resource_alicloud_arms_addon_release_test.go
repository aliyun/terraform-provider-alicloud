package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Arms AddonRelease. >>> Resource test cases, automatically generated.
// Case 4607
func TestAccAliCloudArmsAddonRelease_basic4607(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_addon_release.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsAddonReleaseMap4607)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsAddonRelease")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsaddonrelease%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsAddonReleaseBasicDependence4607)
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
					"environment_id": "${alicloud_arms_environment.env-addonrelease.id}",
					"addon_version":  "0.0.1",
					"addon_name":     "mysql",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": CHECKSET,
						"addon_version":  "0.0.1",
						"addon_name":     "mysql",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"values": "{\\\"host\\\":\\\"mysql-service.default\\\",\\\"port\\\":3306,\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"roots\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"values": "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"roots\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addon_version": "0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addon_version": "0.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addon_version": "0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addon_version": "0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"values": "{\\\"host\\\":\\\"mysql-service.default\\\",\\\"port\\\":3306,\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"rootroot\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"values": "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"rootroot\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addon_version": "0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addon_version": "0.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"values": "{\\\"host\\\":\\\"mysql-service.default\\\",\\\"port\\\":3306,\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"roots\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"values": "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"roots\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.env-addonrelease.id}",
					"addon_version":  "0.0.1",
					"addon_name":     "mysql",
					"values":         "{\\\"host\\\":\\\"mysql-service.default\\\",\\\"port\\\":3306,\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"roots\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": CHECKSET,
						"addon_version":  "0.0.1",
						"addon_name":     "mysql",
						"values":         "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"roots\"}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudArmsAddonReleaseMap4607 = map[string]string{
	"addon_release_name": CHECKSET,
	"create_time":        CHECKSET,
}

func AlicloudArmsAddonReleaseBasicDependence4607(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-test1-hz"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vswitch" {
  description  = "api-resource-test1-hz"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = var.name

  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
}


resource "alicloud_snapshot_policy" "default" {
  name            = var.name
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

data "alicloud_instance_types" "default" {
  availability_zone    = alicloud_vswitch.vswitch.zone_id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.sn1ne"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name               = var.name
  cluster_spec       = "ack.pro.small"
  version            = "1.24.6-aliyun.1"
  new_nat_gateway    = true
  node_cidr_mask     = 26
  proxy_mode         = "ipvs"
  service_cidr       = "172.23.0.0/16"
  pod_cidr           = "10.95.0.0/16"
  worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "env-addonrelease" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ACK"
}


`, name)
}

// Case 4607  twin
func TestAccAliCloudArmsAddonRelease_basic4607_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_addon_release.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsAddonReleaseMap4607)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsAddonRelease")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsaddonrelease%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsAddonReleaseBasicDependence4607)
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
					"addon_release_name": name,
					"environment_id":     "${alicloud_arms_environment.env-addonrelease.id}",
					"addon_version":      "0.0.1",
					"addon_name":         "mysql",
					"aliyun_lang":        "en",
					"values":             "{\\\"host\\\":\\\"mysql-service.default\\\",\\\"port\\\":3306,\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"roots\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addon_release_name": name,
						"environment_id":     CHECKSET,
						"addon_version":      "0.0.1",
						"addon_name":         "mysql",
						"aliyun_lang":        "en",
						"values":             "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"roots\"}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Arms AddonRelease. <<< Resource test cases, automatically generated.
