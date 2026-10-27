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

// Test VpcIpam IpamResourceDiscovery. >>> Resource test cases, automatically generated.
// Case test_ipam_resource_discovery 9871
func TestAccAliCloudVpcIpamIpamResourceDiscovery_basic9871(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_resource_discovery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamResourceDiscoveryMap9871)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpamResourceDiscovery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamResourceDiscoveryBasicDependence9871)
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
					"operating_region_list": []string{
						"cn-hangzhou"},
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"ipam_resource_discovery_description": "This is a custom IPAM resource discovery.",
					"ipam_resource_discovery_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operating_region_list.#":             "1",
						"resource_group_id":                   acctest.CHECKSET,
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
						"resource_group_id":                   acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcIpamIpamResourceDiscoveryMap9871 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
