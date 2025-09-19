package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA TransportLayerApplication. >>> Resource test cases, automatically generated.
// Case transportLayerApplication_test
func TestAccAliCloudESATransportLayerApplicationtransportLayerApplication_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_transport_layer_application.default"
	ra := resourceAttrInit(resourceId, AliCloudESATransportLayerApplicationtransportLayerApplication_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaTransportLayerApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESATransportLayerApplication%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESATransportLayerApplicationtransportLayerApplication_testBasicDependence)
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
					"record_name":               "test-resource2.gositecdn.cn",
					"site_id":                   "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"ip_access_rule":            "off",
					"ipv6":                      "off",
					"cross_border_optimization": "off",
					"rules": []map[string]interface{}{

						{
							"comment":                     "test_transportLayerApplication",
							"edge_port":                   "80",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "8080",
							"client_ip_pass_through_mode": "off",
							"source":                      "1.2.3.4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_access_rule": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_border_optimization": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{

						{
							"comment":                     "test_transportLayerApplication",
							"edge_port":                   "80",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "8080",
							"client_ip_pass_through_mode": "off",
							"source":                      "1.2.3.4",
						},

						{
							"comment":                     "test_transportLayerApplication_modifies",
							"edge_port":                   "88",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv1",
							"source":                      "1.2.3.5",
						},

						{
							"comment":                     "test_transportLayerApplication_modifies-ip2",
							"edge_port":                   "8080",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv1",
							"source":                      "1.2.3.6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
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

var AliCloudESATransportLayerApplicationtransportLayerApplication_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESATransportLayerApplicationtransportLayerApplication_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test ESA TransportLayerApplication. <<< Resource test cases, automatically generated.
