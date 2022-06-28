package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaEndpointGroup_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGaEndpointGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaEndpointGroupBasicDependence)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
					"listener_id": "${alicloud_ga_listener.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_configurations.#": "1",
						"listener_id":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id", "endpoint_group_type"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
						{
							"endpoint": "${alicloud_eip_address.default.1.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "EndpointGroup_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "EndpointGroup_update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval_seconds": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval_seconds": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/healthcheckupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/healthcheckupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_port": `30`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_overrides.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_count": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_percentage": `30`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_percentage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
					"description":                   "EndpointGroup",
					"health_check_interval_seconds": `3`,
					"health_check_path":             "/healthcheck",
					"health_check_port":             `20`,
					"health_check_protocol":         "tcp",
					"name":                          name,
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "70",
						},
					},
					"threshold_count":    `3`,
					"traffic_percentage": `20`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#":     "1",
						"description":                   "EndpointGroup",
						"health_check_interval_seconds": "3",
						"health_check_path":             "/healthcheck",
						"health_check_port":             "20",
						"health_check_protocol":         "tcp",
						"name":                          name,
						"port_overrides.#":              "1",
						"threshold_count":               "3",
						"traffic_percentage":            "20",
					}),
				),
			},
		},
	})
}

var AlicloudGaEndpointGroupMap = map[string]string{
	"endpoint_group_type": "default",
	"status":              CHECKSET,
	"threshold_count":     "3",
}

func AlicloudGaEndpointGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default  = "%s"
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

resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="60"
    to_port="70"
  }
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity="SOURCE_IP"
  protocol="UDP"
  name=var.name
}

resource "alicloud_eip_address" "default" {
  count = 2
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name = var.name
}
`, name)
}
