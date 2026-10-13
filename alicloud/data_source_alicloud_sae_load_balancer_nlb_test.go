package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSaeLoadBalancerNlbDataSource_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	resourceId := "alicloud_sae_load_balancer_nlb.default"
	dataSourceId := "data.alicloud_sae_load_balancer_nlb.default"
	ra := resourceAttrInit(resourceId, AlicloudSAELoadBalancerNlbMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApplicationNlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAELoadBalancerNlbBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alicloud_sae_application.default.id}",
					"nlb_id":       "${alicloud_nlb_load_balancer.default.id}",
					"address_type": "Internet",
					"listeners": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "80",
							"target_port": "8080",
							"cert_ids":    "cert-id-1",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default.id}",
							"zone_id":    "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
					},
				}) + fmt.Sprintf(`
data "alicloud_sae_load_balancer_nlb" "default" {
  app_id = alicloud_sae_load_balancer_nlb.default.app_id
}
`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          CHECKSET,
						"listeners.#":     "1",
						"zone_mappings.#": "1",
					}),
					resource.TestCheckResourceAttrSet(dataSourceId, "app_id"),
					resource.TestCheckResourceAttrSet(dataSourceId, "dns_name"),
				),
			},
		},
	})
}
