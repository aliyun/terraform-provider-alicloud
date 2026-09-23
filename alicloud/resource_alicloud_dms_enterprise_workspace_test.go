// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DmsEnterprise Workspace. >>> Resource test cases, automatically generated.
// Case workspace_testcases_new_打包 11210
func TestAccAliCloudDmsEnterpriseWorkspace_basic11210(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudDmsEnterpriseWorkspaceMap11210)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DMSEnterpriseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsEnterpriseWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdmsenterprise%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDmsEnterpriseWorkspaceBasicDependence11210)
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
					"description":    name,
					"workspace_name": name,
					"vpc_id":         "${alicloud_vpc.vpc_create.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"workspace_name": name,
						"vpc_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "testworkspace-1755485906",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_name": name + "_update",
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

func TestAccAliCloudDmsEnterpriseWorkspace_basic11210_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_workspace.default"
	ra := resourceAttrInit(resourceId, AliCloudDmsEnterpriseWorkspaceMap11210)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DMSEnterpriseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsEnterpriseWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdmsenterprise%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDmsEnterpriseWorkspaceBasicDependence11210)
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
					"description":    name,
					"workspace_name": name,
					"vpc_id":         "${alicloud_vpc.vpc_create.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"workspace_name": name,
						"vpc_id":         CHECKSET,
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

var AliCloudDmsEnterpriseWorkspaceMap11210 = map[string]string{
	"region_id": CHECKSET,
}

func AliCloudDmsEnterpriseWorkspaceBasicDependence11210(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "vpc_create" {
  is_default  = false
  description = "test vpc"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "test_vpc_for_zhenyuan"
}


`, name)
}

// Test DmsEnterprise Workspace. <<< Resource test cases, automatically generated.
