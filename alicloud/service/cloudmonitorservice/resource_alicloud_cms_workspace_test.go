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

// Test Cms Workspace. >>> Resource test cases, automatically generated.
// Case Umodel资源测试 8476
func TestAccAliCloudCmsWorkspace_basic8476(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_workspace.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsWorkspaceMap8476)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsWorkspaceBasicDependence8476)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_project":    "${alicloud_log_project.default.project_name}",
					"workspace_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project":    acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsWorkspaceMap8476)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsWorkspaceBasicDependence8476)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"sls_project":    acctest.CHECKSET,
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
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
