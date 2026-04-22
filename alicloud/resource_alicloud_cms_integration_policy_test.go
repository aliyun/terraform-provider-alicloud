// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms IntegrationPolicy. >>> Resource test cases, automatically generated.
// Case IntegrationPolicy模型测试_副本1730689370617 8635
func TestAccAliCloudCmsIntegrationPolicy_basic8635(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8635)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8635)
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
					"policy_type": "CS",
					"entity_group": []map[string]interface{}{
						{
							"cluster_id": "${data.alicloud_cs_clusters.default.ids.0}",
						},
					},
					"integration_policy_name": name,
					"workspace":               "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":             "CS",
						"integration_policy_name": name,
						"workspace":               CHECKSET,
						"entity_group.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"integration_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"integration_policy_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudCmsIntegrationPolicy_basic8635_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8635)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8635)
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
					"policy_type": "CS",
					"entity_group": []map[string]interface{}{
						{
							"cluster_id":          "${data.alicloud_cs_clusters.default.ids.0}",
							"cluster_entity_type": "acs.ack.cluster",
						},
					},
					"integration_policy_name": name,
					"workspace":               "${alicloud_cms_workspace.default.id}",
					"force":                   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":             "CS",
						"integration_policy_name": name,
						"workspace":               CHECKSET,
						"entity_group.#":          "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AliCloudCmsIntegrationPolicyMap8635 = map[string]string{
	"region_id": CHECKSET,
}

func AliCloudCmsIntegrationPolicyBasicDependence8635(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_cs_clusters" "default" {
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

// Case IntegrationPolicy模型测试5 8626
func TestAccAliCloudCmsIntegrationPolicy_basic8626(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8626)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8626)
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
					"policy_type":             "ECS",
					"integration_policy_name": name,
					"workspace":               "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":             "ECS",
						"integration_policy_name": name,
						"workspace":               CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"integration_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"integration_policy_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudCmsIntegrationPolicy_basic8626_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_integration_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8626)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsIntegrationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8626)
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
					"policy_type":             "ECS",
					"integration_policy_name": name,
					"workspace":               "${alicloud_cms_workspace.default.id}",
					"force":                   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":             "ECS",
						"integration_policy_name": name,
						"workspace":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AliCloudCmsIntegrationPolicyMap8626 = map[string]string{
	"region_id": CHECKSET,
}

func AliCloudCmsIntegrationPolicyBasicDependence8626(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

// Test Cms IntegrationPolicy. <<< Resource test cases, automatically generated.
