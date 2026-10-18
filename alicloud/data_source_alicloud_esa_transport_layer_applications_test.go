package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEsaTransportLayerApplicationsDataSource(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_transport_layer_application.default"
	dataSourceId := "data.alicloud_esa_transport_layer_applications.default"
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
					"keep_alive_protection":     "on",
					"static_ip":                 "on",
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
				}) + testAccAlicloudEsaTransportLayerApplicationsDataSourceConfig("exact", name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keep_alive_protection":     "on",
						"static_ip":                 "on",
						"cross_border_optimization": "on",
						"ip_access_rule":            "on",
						"ipv6":                      "on",
					}),
					resource.TestCheckResourceAttr(dataSourceId, "applications.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceId, "applications.0.id", resourceId, "id"),
					resource.TestCheckResourceAttrPair(dataSourceId, "applications.0.record_name", resourceId, "record_name"),
					resource.TestCheckResourceAttrPair(dataSourceId, "applications.0.site_id", resourceId, "site_id"),
					resource.TestCheckResourceAttr(dataSourceId, "applications.0.keep_alive_protection", "on"),
					resource.TestCheckResourceAttr(dataSourceId, "applications.0.static_ip", "on"),
					resource.TestCheckResourceAttr(dataSourceId, "applications.0.rules.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "applications.0.rules.0.edge_port", "80"),
				),
			},
		},
	})
}

func testAccAlicloudEsaTransportLayerApplicationsDataSourceConfig(matchType string, name string) string {
	return fmt.Sprintf(`
data "alicloud_esa_transport_layer_applications" "default" {
  site_id        = data.alicloud_esa_sites.default.sites.0.id
  match_type     = "%s"
  record_name    = "%s.${data.alicloud_esa_sites.default.sites.0.site_name}"
  enable_details = true
  ids            = ["${alicloud_esa_transport_layer_application.default.id}"]
}
`, matchType, name)
}
