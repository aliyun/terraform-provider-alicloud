package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Workspace. >>> Resource test cases, automatically generated.
// Case CC发布PAIWorkspace_副本1726115966862 7875
func TestAccAliCloudPaiWorkspaceWorkspace_basic7875(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceWorkspaceMap7875)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceWorkspaceBasicDependence7875)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "452",
					"workspace_name": name,
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    CHECKSET,
						"workspace_name": name,
						"env_types.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "new_test_pop_559",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "986",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudPaiWorkspaceWorkspace_basic7875_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceWorkspaceMap7875)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceWorkspaceBasicDependence7875)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "452",
					"workspace_name":    name,
					"display_name":      "new_test_pop_559",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       CHECKSET,
						"workspace_name":    name,
						"display_name":      CHECKSET,
						"resource_group_id": CHECKSET,
						"env_types.#":       "1",
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

var AliCloudPaiWorkspaceWorkspaceMap7875 = map[string]string{
	"status":            CHECKSET,
	"create_time":       CHECKSET,
	"resource_group_id": CHECKSET,
	"display_name":      CHECKSET,
}

func AliCloudPaiWorkspaceWorkspaceBasicDependence7875(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Test PaiWorkspace Workspace. <<< Resource test cases, automatically generated.
