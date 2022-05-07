package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenVbrHealthCheck_basic(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v cbn.VbrHealthCheck
	resourceId := "alicloud_cen_vbr_health_check.default"
	ra := resourceAttrInit(resourceId, CenVbrHealthCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenVbrHealthCheck")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenVbrHealthCheck%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenVbrHealthCheckBasicdependence)
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
					"cen_id":                 "${alicloud_cen_instance.default.id}",
					"health_check_source_ip": "192.168.1.2",
					"health_check_target_ip": "10.0.0.2",
					"vbr_instance_id":        "${alicloud_express_connect_virtual_border_router.default.id}",
					"vbr_instance_region_id": os.Getenv("ALICLOUD_REGION"),
					"health_check_interval":  "2",
					"healthy_threshold":      "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                 CHECKSET,
						"health_check_source_ip": "192.168.1.2",
						"health_check_target_ip": "10.0.0.2",
						"vbr_instance_id":        CHECKSET,
						"vbr_instance_region_id": os.Getenv("ALICLOUD_REGION"),
						"health_check_interval":  "2",
						"healthy_threshold":      "8",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_source_ip": "192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_source_ip": "192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_target_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_target_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_source_ip": "192.168.1.2",
					"health_check_target_ip": "10.0.0.2",
					"health_check_interval":  "2",
					"healthy_threshold":      "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_source_ip": "192.168.1.2",
						"health_check_target_ip": "10.0.0.2",
						"health_check_interval":  "2",
						"healthy_threshold":      "8",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCenVbrHealthCheck_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v cbn.VbrHealthCheck
	resourceId := "alicloud_cen_vbr_health_check.default"
	ra := resourceAttrInit(resourceId, CenVbrHealthCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenVbrHealthCheck")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tf-testAccCenVbrHealthCheck%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenVbrHealthCheckBasicdependence1)
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
					"cen_id":                 "${alicloud_cen_instance.default.id}",
					"health_check_source_ip": "192.168.1.2",
					"health_check_target_ip": "10.0.0.2",
					"vbr_instance_id":        "${alicloud_express_connect_virtual_border_router.default.id}",
					"vbr_instance_region_id": os.Getenv("ALICLOUD_REGION"),
					"health_check_interval":  "2",
					"healthy_threshold":      "8",
					"vbr_instance_owner_id":  "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                 CHECKSET,
						"health_check_source_ip": "192.168.1.2",
						"health_check_target_ip": "10.0.0.2",
						"vbr_instance_id":        CHECKSET,
						"vbr_instance_region_id": os.Getenv("ALICLOUD_REGION"),
						"health_check_interval":  "2",
						"healthy_threshold":      "8",
						"vbr_instance_owner_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vbr_instance_owner_id"},
			},
		},
	})
}

var CenVbrHealthCheckMap = map[string]string{
	"cen_id": CHECKSET,
}

func CenVbrHealthCheckBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}
resource "alicloud_cen_instance_attachment" "default" {
  instance_id = "${alicloud_cen_instance.default.id}"
  child_instance_id = alicloud_express_connect_virtual_border_router.default.id
  child_instance_type = "VBR"
  child_instance_region_id = "%s"
}
`, name, acctest.RandIntRange(1, 2999), defaultRegionToTest)
}

func CenVbrHealthCheckBasicdependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
   default = "%s"
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}
resource "alicloud_cen_instance_attachment" "default" {
  instance_id = "${alicloud_cen_instance.default.id}"
  child_instance_id = alicloud_express_connect_virtual_border_router.default.id
  child_instance_type = "VBR"
  child_instance_region_id = "%s"
}
data "alicloud_account" "default" {}
`, name, acctest.RandIntRange(1, 2999), os.Getenv("ALICLOUD_REGION"))
}
