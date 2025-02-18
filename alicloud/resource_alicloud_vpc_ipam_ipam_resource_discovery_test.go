package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpcIpam IpamResourceDiscovery. >>> Resource test cases, automatically generated.
// Case test_ipam_resource_discovery 9871
func TestAccAliCloudVpcIpamIpamResourceDiscovery_basic9871(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_resource_discovery.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamResourceDiscoveryMap9871)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamResourceDiscovery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamResourceDiscoveryBasicDependence9871)
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
					"operating_region_list": []string{
						"cn-hangzhou"},
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"ipam_resource_discovery_description": "This is a custom IPAM resource discovery.",
					"ipam_resource_discovery_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operating_region_list.#":             "1",
						"resource_group_id":                   CHECKSET,
						"ipam_resource_discovery_description": "This is a custom IPAM resource discovery.",
						"ipam_resource_discovery_name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operating_region_list": []string{
						"cn-hangzhou", "ap-southeast-3", "ap-southeast-6"},
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"ipam_resource_discovery_description": "Description.",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operating_region_list.#":             "3",
						"resource_group_id":                   CHECKSET,
						"ipam_resource_discovery_description": "Description.",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operating_region_list": []string{
						"cn-hangzhou"},
					"ipam_resource_discovery_description": "This is my new rd.",
					"ipam_resource_discovery_name":        name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operating_region_list.#":             "1",
						"ipam_resource_discovery_description": "This is my new rd.",
						"ipam_resource_discovery_name":        name + "_update",
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

var AlicloudVpcIpamIpamResourceDiscoveryMap9871 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcIpamIpamResourceDiscoveryBasicDependence9871(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test VpcIpam IpamResourceDiscovery. <<< Resource test cases, automatically generated.
