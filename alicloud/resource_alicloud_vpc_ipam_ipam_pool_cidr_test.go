package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpcIpam IpamPoolCidr. >>> Resource test cases, automatically generated.
// Case test_ipam_pool_cidr 8028
func TestAccAliCloudVpcIpamIpamPoolCidr_basic8028(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_cidr.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolCidrMap8028)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPoolCidr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampoolcidr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolCidrBasicDependence8028)
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
					"cidr":         "10.0.0.0/8",
					"ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr":         "10.0.0.0/8",
						"ipam_pool_id": CHECKSET,
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

var AlicloudVpcIpamIpamPoolCidrMap8028 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcIpamIpamPoolCidrBasicDependence8028(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv4"
}


`, name)
}

// Test VpcIpam IpamPoolCidr. <<< Resource test cases, automatically generated.
