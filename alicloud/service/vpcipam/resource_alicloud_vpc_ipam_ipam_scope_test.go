package vpcipam_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcipam"
)

// Test VpcIpam IpamScope. >>> Resource test cases, automatically generated.
// Case test_ipam_scope 8006
func TestAccAliCloudVpcIpamIpamScope_basic8006(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_scope.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamScopeMap8006)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpamScope")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipamscope%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamScopeBasicDependence8006)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_name":        name,
					"ipam_id":                "${alicloud_vpc_ipam_ipam.defaultIpam.id}",
					"ipam_scope_description": "This is a ipam scope.",
					"ipam_scope_type":        "private",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_name":        name,
						"ipam_id":                acctest.CHECKSET,
						"ipam_scope_description": "This is a ipam scope.",
						"ipam_scope_type":        "private",
						"resource_group_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_name":        name + "_update",
					"ipam_scope_description": "This is my ipam scope.",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_name":        name + "_update",
						"ipam_scope_description": "This is my ipam scope.",
						"resource_group_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region_id"},
			},
		},
	})
}

var AlicloudVpcIpamIpamScopeMap8006 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudVpcIpamIpamScopeBasicDependence8006(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}


`, name)
}

// Test VpcIpam IpamScope. <<< Resource test cases, automatically generated.
