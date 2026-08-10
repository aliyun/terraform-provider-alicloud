package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEnsHaVip_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_ha_vip.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsHaVipMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsHaVip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccenshavip%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsHaVipBasicDependence)

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
					"vswitch_id":  "${alicloud_ens_vswitch.default.id}",
					"amount":      1,
					"description": "desc1",
					"ip_address":  "10.0.9.5",
					"ha_vip_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"amount":        "1",
						"description":   "desc1",
						"ip_address":    "10.0.9.5",
						"ha_vip_name":   name,
						"vswitch_id":    CHECKSET,
						"status":        CHECKSET,
						"create_time":   CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":  "${alicloud_ens_vswitch.default.id}",
					"amount":      1,
					"description": "desc1",
					"ip_address":  "10.0.9.5",
					"ha_vip_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":  "${alicloud_ens_vswitch.default.id}",
					"ip_address":  "10.0.9.5",
					"ha_vip_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name": name,
						"vswitch_id":  CHECKSET,
						"ip_address":  "10.0.9.5",
						"status":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount"},
			},
		},
	})
}

func TestAccAliCloudEnsHaVip_autoIpAddress(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_ha_vip.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsHaVipMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsHaVip")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccenshavip%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsHaVipBasicDependence)

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
					"vswitch_id":  "${alicloud_ens_vswitch.default.id}",
					"amount":      2,
					"ha_vip_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"amount":        "2",
						"ha_vip_name":   name,
						"vswitch_id":    CHECKSET,
						"status":        CHECKSET,
						"ip_address":    CHECKSET,
						"create_time":   CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudEnsHaVipMap = map[string]string{
	"status":        CHECKSET,
	"create_time":   CHECKSET,
	"ens_region_id": CHECKSET,
	"network_id":    CHECKSET,
	"ip_address":    CHECKSET,
	"vswitch_id":    CHECKSET,
	"ha_vip_name":   CHECKSET,
	"description":   CHECKSET,
}

func AliCloudEnsHaVipBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_network" "default" {
  network_name  = "%s"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "10.0.9.0/24"
  vswitch_name  = "%s"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}
`, name, name, name)
}
