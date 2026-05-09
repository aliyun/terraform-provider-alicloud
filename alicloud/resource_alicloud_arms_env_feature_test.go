package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Arms EnvFeature. >>> Resource test cases, automatically generated.
// Case 4606
func TestAccAliCloudArmsEnvFeature_basic4606(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_feature.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvFeatureMap4606)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvFeature")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvfeature%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvFeatureBasicDependence4606)
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
					"env_feature_name": "metric-agent",
					"environment_id":   "${alicloud_arms_environment.env-feature.id}",
					"feature_version":  "1.1.17",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_feature_name": "metric-agent",
						"environment_id":   CHECKSET,
						"feature_version":  "1.1.17",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"feature_version": "1.1.17",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"feature_version": "1.1.17",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"feature_version": "1.1.18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"feature_version": "1.1.18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"feature_version": "1.1.17",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"feature_version": "1.1.17",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"feature_version": "1.1.18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"feature_version": "1.1.18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_feature_name": "metric-agent",
					"environment_id":   "${alicloud_arms_environment.env-feature.id}",
					"feature_version":  "1.1.17",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_feature_name": "metric-agent",
						"environment_id":   CHECKSET,
						"feature_version":  "1.1.17",
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

var AlicloudArmsEnvFeatureMap4606 = map[string]string{
	"status":    CHECKSET,
	"namespace": CHECKSET,
}

func AlicloudArmsEnvFeatureBasicDependence4606(name string) string {
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
  pod_cidr             = "10.125.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_arms_environment" "env-feature" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = local.cluster_id
  environment_sub_type = "ManagedKubernetes"
}


`, name)
}

// Case 4606  twin
func TestAccAliCloudArmsEnvFeature_basic4606_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_feature.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvFeatureMap4606)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvFeature")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvfeature%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvFeatureBasicDependence4606)
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
					"env_feature_name": "metric-agent",
					"environment_id":   "${alicloud_arms_environment.env-feature.id}",
					"feature_version":  "1.1.18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_feature_name": "metric-agent",
						"environment_id":   CHECKSET,
						"feature_version":  "1.1.18",
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

// Test Arms EnvFeature. <<< Resource test cases, automatically generated.
