package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ApiGateway TrafficControl. >>> Resource test cases.
func TestAccAliCloudApiGatewayTrafficControl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_traffic_control.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayTrafficControlMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayTrafficControl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_trafficcontrol%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayTrafficControlBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_control_name": name,
					"traffic_control_unit": "MINUTE",
					"api_default":          100,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_control_name": name,
						"traffic_control_unit": "MINUTE",
						"api_default":          "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_control_name": name + "u",
					"traffic_control_unit": "HOUR",
					"api_default":          200,
					"user_default":         50,
					"app_default":          30,
					"description":          "test traffic control",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_control_name": name + "u",
						"traffic_control_unit": "HOUR",
						"api_default":          "200",
						"user_default":         "50",
						"app_default":          "30",
						"description":          "test traffic control",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_control_name": name,
					"traffic_control_unit": "DAY",
					"api_default":          300,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_control_name": name,
						"traffic_control_unit": "DAY",
						"api_default":          "300",
					}),
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

var AlicloudApiGatewayTrafficControlMap0 = map[string]string{}

func AlicloudApiGatewayTrafficControlBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

// Test ApiGateway TrafficControl. <<< Resource test cases.
