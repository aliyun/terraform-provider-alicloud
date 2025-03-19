package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEfloNode_basic10171(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap10171)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence10171)
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-hangzhou-b",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "B1",
					"server_arch":       "bmserver",
					"computing_server":  "efg2.C48cA3sen",
					"stage_num":         "36",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-hangzhou-b",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "B1",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48cA3sen",
						"stage_num":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "pricing_cycle", "product_code", "product_form", "product_type", "renew_period", "renewal_status", "server_arch", "stage_num", "subscription_type", "zone"},
			},
		},
	})
}

var AlicloudEfloNodeMap10171 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence10171(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}
