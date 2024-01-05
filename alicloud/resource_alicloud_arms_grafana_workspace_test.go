package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Arms GrafanaWorkspace. >>> Resource test cases, automatically generated.
// Case 5667
func TestAccAliCloudArmsGrafanaWorkspace_basic5667(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5667)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5667)
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
					"grafana_workspace_name":    name,
					"grafana_version":           "9.0.x",
					"grafana_workspace_edition": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_workspace_name":    name,
						"grafana_version":           "9.0.x",
						"grafana_workspace_edition": "standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "8.2.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "8.2.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "grafana-rg-create-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "grafana-rg-create-test",
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
					"description": "grafana-rg-update-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "grafana-rg-update-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_workspace_name": name + "_update",
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
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name + "_update",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap5667 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence5667(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Case 5667  twin
func TestAccAliCloudArmsGrafanaWorkspace_basic5667_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5667)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5667)
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
					"grafana_version":           "9.0.x",
					"description":               "grafana-rg-update-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "9.0.x",
						"description":               "grafana-rg-update-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
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

// Test Arms GrafanaWorkspace. <<< Resource test cases, automatically generated.
