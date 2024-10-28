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
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceWorkspaceMap7875)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceWorkspaceBasicDependence7875)
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
					"display_name":   "test_pop_584",
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    CHECKSET,
						"workspace_name": name,
						"display_name":   CHECKSET,
						"env_types.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "184",
					"display_name": "new_test_pop_559",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  CHECKSET,
						"display_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "986",
					"display_name": "test_pop_403",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  CHECKSET,
						"display_name": CHECKSET,
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

var AlicloudPaiWorkspaceWorkspaceMap7875 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceWorkspaceBasicDependence7875(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PaiWorkspace Workspace. <<< Resource test cases, automatically generated.
