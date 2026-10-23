package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCrVpcEndpointLinkedVpc_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_vpc_endpoint_linked_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudCrVpcEndpointLinkedVpcMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrVpcEndpointLinkedVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-vpc-ep-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrVpcEndpointLinkedVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                      "${alicloud_cr_ee_instance.default.id}",
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

	resource "alicloud_cr_ee_instance" "default" {
		payment_type   = "Subscription"
		period         = 1
		renewal_status = "ManualRenewal"
		instance_type  = "Advanced"
		instance_name  = var.name
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
`, name)
}
