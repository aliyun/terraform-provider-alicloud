package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA TransportLayerApplication. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaTransportLayerApplication_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_transport_layer_application.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaTransportLayerApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaTransportLayerApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%stla%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaTransportLayerApplicationBasicDependence0)
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
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"record_name": name + ".${data.alicloud_esa_sites.default.sites.0.site_name}",
					"rules": []map[string]interface{}{
						{
							"edge_port":                   "80",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "8080",
							"client_ip_pass_through_mode": "off",
							"source":                      "1.1.1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     CHECKSET,
						"record_name": CHECKSET,
						"rules.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_border_optimization": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_border_optimization": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_access_rule": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_access_rule": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"comment":                     name + "udp",
							"edge_port":                   "88",
							"source_type":                 "domain",
							"protocol":                    "UDP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv2",
							"source":                      name + "udp.${data.alicloud_esa_sites.default.sites.0.site_name}",
						},
						{
							"comment":                     name + "1",
							"edge_port":                   "82",
							"source_type":                 "ip",
							"protocol":                    "UDP",
							"source_port":                 "86",
							"client_ip_pass_through_mode": "PPv2",
							"source":                      "2.2.2.2",
						},
						{
							"comment":                     name + "2",
							"edge_port":                   "8080",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv1",
							"source":                      "6.6.6.6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"comment":                     name,
							"edge_port":                   "80",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "8080",
							"client_ip_pass_through_mode": "off",
							"source":                      "1.1.1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "1",
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

func TestAccAliCloudEsaTransportLayerApplication_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_transport_layer_application.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaTransportLayerApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaTransportLayerApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%stla%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaTransportLayerApplicationBasicDependence0)
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
					"site_id":                   "${data.alicloud_esa_sites.default.sites.0.id}",
					"record_name":               name + ".${data.alicloud_esa_sites.default.sites.0.site_name}",
					"cross_border_optimization": "on",
					"ip_access_rule":            "on",
					"ipv6":                      "on",
					"rules": []map[string]interface{}{
						{
							"comment":                     name + "udp",
							"edge_port":                   "88",
							"source_type":                 "domain",
							"protocol":                    "UDP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv2",
							"source":                      name + "udp.${data.alicloud_esa_sites.default.sites.0.site_name}",
						},
						{
							"comment":                     name + "1",
							"edge_port":                   "82",
							"source_type":                 "ip",
							"protocol":                    "UDP",
							"source_port":                 "86",
							"client_ip_pass_through_mode": "PPv2",
							"source":                      "2.2.2.2",
						},
						{
							"comment":                     name + "2",
							"edge_port":                   "8080",
							"source_type":                 "ip",
							"protocol":                    "TCP",
							"source_port":                 "80",
							"client_ip_pass_through_mode": "PPv1",
							"source":                      "6.6.6.6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                   CHECKSET,
						"record_name":               CHECKSET,
						"cross_border_optimization": "on",
						"ip_access_rule":            "on",
						"ipv6":                      "on",

						"rules.#": "3",
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

var AliCloudEsaTransportLayerApplicationMap0 = map[string]string{
	"application_id":            CHECKSET,
	"cross_border_optimization": CHECKSET,
	"ip_access_rule":            CHECKSET,
	"ipv6":                      CHECKSET,
	"status":                    CHECKSET,
}

func AliCloudEsaTransportLayerApplicationBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "tftestacc.com"
}
`, name)
}

// Test ESA TransportLayerApplication. <<< Resource test cases, automatically generated.
