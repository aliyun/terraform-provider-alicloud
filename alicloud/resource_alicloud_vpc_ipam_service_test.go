package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpcIpam Service. >>> Resource test cases, automatically generated.
// Case ccc 8406
func TestAccAliCloudVpcIpamService_basic8406(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_service.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamServiceMap8406)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamServiceBasicDependence8406)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudVpcIpamServiceMap8406 = map[string]string{
	"enabled": CHECKSET,
}

func AlicloudVpcIpamServiceBasicDependence8406(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test VpcIpam Service. <<< Resource test cases, automatically generated.
