package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCrVpcEndpointLinkedVpc_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_vpc_endpoint_linked_vpc.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrVpcEndpointLinkedVpcMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrVpcEndpointLinkedVpc")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCrVpcEndpointLinkedVpc%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrVpcEndpointLinkedVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"instance_id": acctest.CHECKSET,
						"vpc_id":      acctest.CHECKSET,
						"vswitch_id":  acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
