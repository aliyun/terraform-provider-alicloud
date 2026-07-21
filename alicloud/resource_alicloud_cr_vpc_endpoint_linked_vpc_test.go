package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCrVpcEndpointLinkedVpc_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_vpc_endpoint_linked_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudCrVpcEndpointLinkedVpcMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrVpcEndpointLinkedVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCrVpcEndpointLinkedVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrVpcEndpointLinkedVpcBasicDependence0)
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
					"instance_id":                      "${data.alicloud_cr_ee_instances.default.ids.0}",
					"vpc_id":                           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                       "${data.alicloud_vswitches.default.ids.0}",
					"module_name":                      "Registry",
					"enable_create_dns_record_in_pvzt": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"vpc_id":      CHECKSET,
						"vswitch_id":  CHECKSET,
						"module_name": "Registry",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_create_dns_record_in_pvzt"},
			},
		},
	})
}

var AlicloudCrVpcEndpointLinkedVpcMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudCrVpcEndpointLinkedVpcBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
`, name)
}
