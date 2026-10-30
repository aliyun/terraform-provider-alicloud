package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Arms AddonRelease. >>> Resource test cases, automatically generated.
// Case 4607
func TestAccAliCloudArmsAddonRelease_basic4607(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_addon_release.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsAddonReleaseMap4607)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsAddonRelease")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsaddonrelease%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsAddonReleaseBasicDependence4607)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.env-addonrelease.id}",
					"addon_version":  "0.0.1",
					"addon_name":     "mysql",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": acctest.CHECKSET,
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
						"environment_id": acctest.CHECKSET,
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
	"addon_release_name": acctest.CHECKSET,
	"create_time":        acctest.CHECKSET,
}

func AlicloudArmsAddonReleaseBasicDependence4607(name string) string {
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
  pod_cidr             = "10.123.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_arms_environment" "env-addonrelease" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = local.cluster_id
  environment_sub_type = "ManagedKubernetes"
}


`, name)
}

// Case 4607  twin
func TestAccAliCloudArmsAddonRelease_basic4607_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_addon_release.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsAddonReleaseMap4607)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsAddonRelease")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsaddonrelease%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsAddonReleaseBasicDependence4607)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"environment_id":     acctest.CHECKSET,
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
