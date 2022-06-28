package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaListener_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudGaListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaListenerBasicDependence)
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
					"accelerator_id": "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
					"description":    "create_description",
					"name":           "${var.name}",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"proxy_protocol": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": CHECKSET,
						"description":    "create_description",
						"name":           name,
						"port_ranges.#":  "1",
						"proxy_protocol": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id", "proxy_protocol"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "100",
							"to_port":   "110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "NONE",
					"protocol":        "UDP",
					"proxy_protocol":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "NONE",
						"protocol":        "UDP",
						"proxy_protocol":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_affinity": "SOURCE_IP",
					"description":     "create_description",
					"protocol":        "TCP",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "60",
							"to_port":   "70",
						},
					},
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_affinity": "SOURCE_IP",
						"description":     "create_description",
						"protocol":        "TCP",
						"port_ranges.#":   "1",
						"name":            name,
					}),
				),
			},
		},
	})
}

var AlicloudGaListenerMap = map[string]string{
	"client_affinity": "NONE",
	"protocol":        "TCP",
	"status":          CHECKSET,
}

func AlicloudGaListenerBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}
`, name)
}
