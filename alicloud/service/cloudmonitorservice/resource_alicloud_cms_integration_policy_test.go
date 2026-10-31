// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cms IntegrationPolicy. >>> Resource test cases, automatically generated.
// Case IntegrationPolicy模型测试_副本1730689370617 8635
func TestAccAliCloudCmsIntegrationPolicy_basic8635(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_integration_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8635)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsIntegrationPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8635)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"workspace":               acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8635)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsIntegrationPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8635)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_type": "CS",
					"entity_group": []map[string]interface{}{
						{
							"cluster_id":           "${data.alicloud_cs_clusters.default.ids.0}",
							"cluster_entity_type":  "acs.ack.cluster",
							"entity_user_id":       "${data.alicloud_account.default.id}",
							"disable_policy_share": "true",
							"cluster_namespace":    "arms-prom",
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
						"workspace":               acctest.CHECKSET,
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
	"region_id": acctest.CHECKSET,
}

func AliCloudCmsIntegrationPolicyBasicDependence8635(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8626)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsIntegrationPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8626)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"workspace":               acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsIntegrationPolicyMap8626)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsIntegrationPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsIntegrationPolicyBasicDependence8626)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"workspace":               acctest.CHECKSET,
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
	"region_id": acctest.CHECKSET,
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
