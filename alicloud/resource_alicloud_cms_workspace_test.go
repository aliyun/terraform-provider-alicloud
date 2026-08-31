// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms Workspace. >>> Resource test cases, automatically generated.
// Case Umodel资源测试 8476
func TestAccAliCloudCmsWorkspace_basic8476(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsWorkspaceMap8476)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsWorkspaceBasicDependence8476)
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
					"sls_project":    "${alicloud_log_project.default.project_name}",
					"workspace_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project":    CHECKSET,
						"workspace_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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

func TestAccAliCloudCmsWorkspace_basic8476_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsWorkspaceMap8476)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsWorkspaceBasicDependence8476)
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
					"sls_project":    "${alicloud_log_project.default.project_name}",
					"workspace_name": name,
					"display_name":   name,
					"description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project":    CHECKSET,
						"workspace_name": name,
						"display_name":   name,
						"description":    name,
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

var AliCloudCmsWorkspaceMap8476 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudCmsWorkspaceBasicDependence8476(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

	resource "alicloud_log_project" "default" {
  		project_name = var.name
	}
`, name)
}

// Test Cms Workspace. <<< Resource test cases, automatically generated.
