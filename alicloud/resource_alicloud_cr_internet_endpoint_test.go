package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudCrInternetEndpoint_basic0 covers the full lifecycle:
// create with two ACL entries -> update to a single entry (add + remove) ->
// update to no entries (remove all) -> import.
//
// The config is expressed as an inline map[string]interface{} (built by
// resourceTestAccConfigFunc) rather than a hand-rolled fmt.Sprintf HCL string
// so that scripts/testing/testing_coverage_rate_check.go can parse the
// `entries` set and its nested `entry`/`comment` attributes for both set and
// modify coverage. The generator-emitted HCL is semantically identical to the
// previous hand-written blocks.
func TestAccAliCloudCrInternetEndpoint_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_internet_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudCrInternetEndpointMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrInternetEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrInternetEndpointBasicDependence0)
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
					"instance_id": "${local.cr_endpoint_instance_id}",
					"entries": []map[string]interface{}{
						{
							"entry":   "192.168.1.0/24",
							"comment": "entry-1",
						},
						{
							"entry":   "10.0.0.0/8",
							"comment": "entry-2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entries": []map[string]interface{}{
						{
							"entry":   "172.16.0.0/12",
							"comment": "entry-3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "1"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entries": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "0"),
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

var AlicloudCrInternetEndpointMap0 = map[string]string{}

func AlicloudCrInternetEndpointBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Economy"
  instance_name  = var.name
  image_scanner  = "DISABLE"
}

locals {
  cr_endpoint_instance_id = alicloud_cr_ee_instance.default.id
}
`, name)
}
