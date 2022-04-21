package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSAESlbIntranet_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_load_balancer_intranet.default"
	ra := resourceAttrInit(resourceId, AlicloudSAESLBIntranetApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApplicationSlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAESlbIntranetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alicloud_sae_application.default.id}",
					"intranet_slb_id": "${alicloud_slb_load_balancer.default.id}",
					"intranet": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "80",
							"target_port": "8080",
						},
						{
							"protocol":    "TCP",
							"port":        "89",
							"target_port": "8989",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          CHECKSET,
						"intranet_slb_id": CHECKSET,
						"intranet.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alicloud_sae_application.default.id}",
					"intranet_slb_id": "${alicloud_slb_load_balancer.default.id}",
					"intranet": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "90",
							"target_port": "9090",
						},
						{
							"protocol":    "TCP",
							"port":        "99",
							"target_port": "9999",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          CHECKSET,
						"intranet_slb_id": CHECKSET,
						"intranet.#":      "2",
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

func AlicloudSAESlbIntranetBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_vpcs" "default"	{
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "intranet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = data.alicloud_vswitches.default.vswitches[0].id
}

resource "alicloud_sae_namespace" "default" {
	namespace_description = var.name
	namespace_id = "%s:%s"
	namespace_name = var.name
}
resource "alicloud_sae_application" "default" {
  app_description= var.name
  app_name=        var.name
  namespace_id=    alicloud_sae_namespace.default.namespace_id
  image_url=     "registry-vpc.%s.aliyuncs.com/lxepoo/apache-php5"
  package_type=    "Image"
  jdk=             "Open JDK 8"
  vswitch_id=      data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone = "Asia/Shanghai"
  replicas=        "5"
  cpu=             "500"
  memory =          "2048"
}


variable "name" {
  default = "%s"
}
`, defaultRegionToTest, name, defaultRegionToTest, name)
}

var AlicloudSAESLBIntranetApplicationMap0 = map[string]string{}
