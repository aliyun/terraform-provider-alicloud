// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package dmsenterprise_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	dmsenterprise "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DmsEnterprise Workspace. >>> Resource test cases, automatically generated.
// Case workspace_testcases_new_打包 11210
func TestAccAliCloudDmsEnterpriseWorkspace_basic11210(t *testing.T) {
	tfacctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_workspace.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudDmsEnterpriseWorkspaceMap11210)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDMSEnterpriseServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDmsEnterpriseWorkspace")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdmsenterprise%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudDmsEnterpriseWorkspaceBasicDependence11210)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vpc_id":         tfacctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "testworkspace-1755485906",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": tfacctest.CHECKSET,
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
	tfacctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_workspace.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudDmsEnterpriseWorkspaceMap11210)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDMSEnterpriseServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDmsEnterpriseWorkspace")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdmsenterprise%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudDmsEnterpriseWorkspaceBasicDependence11210)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vpc_id":         tfacctest.CHECKSET,
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
	"region_id": tfacctest.CHECKSET,
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
