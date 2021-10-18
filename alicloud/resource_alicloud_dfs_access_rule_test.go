package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDFSAccessRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_access_rule.default"
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDFSAccessRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsAccessRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsaccessrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDFSAccessRuleBasicDependence0)
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
					"network_segment": "192.0.2.0/24",
					"access_group_id": "${alicloud_dfs_access_group.default.id}",
					"description":     name,
					"rw_access_type":  "RDWR",
					"priority":        "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_segment": "192.0.2.0/24",
						"access_group_id": CHECKSET,
						"description":     name,
						"rw_access_type":  "RDWR",
						"priority":        "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rw_access_type": "RDONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_access_type": "RDONLY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    name,
					"rw_access_type": "RDWR",
					"priority":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"rw_access_type": "RDWR",
						"priority":       "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDFSAccessRuleMap0 = map[string]string{}

func AlicloudDFSAccessRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_dfs_access_group" "default" {
	network_type = "VPC"
	access_group_name = var.name
	description =  var.name
}

`, name)
}
