// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens NatGatewaySnatEntry. >>> Resource test cases, automatically generated.
// Case Snat_20241218 9626
// NOTE: ENS 订单系统对同一账号串行处理下单，NAT 网关与 EIP 并发创建会被拒绝（OrderFailed），
// 因此依赖链全串行：network -> vswitch -> nat_gateway -> eip -> attach -> snat_entry。
// standby_snat_ip 需要 FullCone 型 NAT/备用 EIP 池，enat.default（对称型）不支持
// （StartSnatStandbySnatIp 返回 SnatNotSupport），配置中以空串占位保持字段可见性。
func TestAccAliCloudEnsNatGatewaySnatEntry_basic9626(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway_snat_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewaySnatEntryMap9626)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEnsNatGatewaySnatEntryConfigStep1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name":   name,
						"source_cidr":       "10.0.0.0/8",
						"nat_gateway_id":    CHECKSET,
						"idle_timeout":      "50",
						"isp_affinity":      "false",
						"eip_affinity":      "false",
						"source_vswitch_id": CHECKSET,
						"source_network_id": CHECKSET,
						"standby_snat_ip":   "",
					}),
				),
			},
			{
				Config: AlicloudEnsNatGatewaySnatEntryConfigStep2(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name + "update",
						"snat_ip":         CHECKSET,
						"isp_affinity":    "true",
						"eip_affinity":    "true",
						"standby_snat_ip": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_vswitch_id", "source_network_id"},
			},
		},
	})
}

var AlicloudEnsNatGatewaySnatEntryMap9626 = map[string]string{
	"status":        CHECKSET,
	"creation_time": CHECKSET,
}

func AlicloudEnsNatGatewaySnatEntryDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
    default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultXqhlfk" {
        network_name = "${var.name}"
        cidr_block = "10.0.0.0/8"
        ens_region_id = "${var.ens_region_id}"
}

resource "alicloud_ens_vswitch" "defaultzkXvut" {
        cidr_block = "10.0.0.0/24"
        vswitch_name = "${var.name}"
        ens_region_id = "${alicloud_ens_network.defaultXqhlfk.ens_region_id}"
        network_id = "${alicloud_ens_network.defaultXqhlfk.id}"
}

resource "alicloud_ens_nat_gateway" "default2Kn0nu" {
        vswitch_id = "${alicloud_ens_vswitch.defaultzkXvut.id}"
        ens_region_id = "${alicloud_ens_vswitch.defaultzkXvut.ens_region_id}"
        network_id = "${alicloud_ens_vswitch.defaultzkXvut.network_id}"
        instance_type = "enat.default"
        nat_name = "${var.name}"
}

resource "alicloud_ens_eip" "defaultiUbwh0" {
        depends_on = [alicloud_ens_nat_gateway.default2Kn0nu]
        bandwidth = "5"
        payment_type = "PayAsYouGo"
        ens_region_id = "${alicloud_ens_vswitch.defaultzkXvut.ens_region_id}"
        eip_name = "${var.name}"
        internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultlI0M0t" {
        instance_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        allocation_id = "${alicloud_ens_eip.defaultiUbwh0.id}"
        instance_type = "Nat"
        standby = false
}

`, name)
}

func AlicloudEnsNatGatewaySnatEntryConfigStep1(name string) string {
	return AlicloudEnsNatGatewaySnatEntryDependence(name) + `
resource "alicloud_ens_nat_gateway_snat_entry" "default" {
        depends_on = [alicloud_ens_eip_instance_attachment.defaultlI0M0t]
        snat_entry_name = "${var.name}"
        source_cidr = "10.0.0.0/8"
        snat_ip = "${alicloud_ens_eip.defaultiUbwh0.ip_address}"
        standby_snat_ip = ""
        nat_gateway_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        idle_timeout = "50"
        isp_affinity = false
        eip_affinity = false
        source_vswitch_id = "${alicloud_ens_vswitch.defaultzkXvut.id}"
        source_network_id = "${alicloud_ens_network.defaultXqhlfk.id}"
}
`
}

func AlicloudEnsNatGatewaySnatEntryConfigStep2(name string) string {
	return AlicloudEnsNatGatewaySnatEntryDependence(name) + `
resource "alicloud_ens_eip" "eip3" {
        depends_on = [alicloud_ens_eip_instance_attachment.defaultlI0M0t]
        bandwidth = "5"
        payment_type = "PayAsYouGo"
        ens_region_id = "${var.ens_region_id}"
        eip_name = "${var.name}-update"
        internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultEip3" {
        instance_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        allocation_id = "${alicloud_ens_eip.eip3.id}"
        instance_type = "Nat"
        standby = false
}

resource "alicloud_ens_nat_gateway_snat_entry" "default" {
        depends_on = [alicloud_ens_eip_instance_attachment.defaultlI0M0t, alicloud_ens_eip_instance_attachment.defaultEip3]
        snat_entry_name = "${var.name}update"
        source_cidr = "10.0.0.0/8"
        snat_ip = "${alicloud_ens_eip.eip3.ip_address}"
        standby_snat_ip = ""
        nat_gateway_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        idle_timeout = "50"
        isp_affinity = true
        eip_affinity = true
        source_vswitch_id = "${alicloud_ens_vswitch.defaultzkXvut.id}"
        source_network_id = "${alicloud_ens_network.defaultXqhlfk.id}"
}
`
}

// Test Ens NatGatewaySnatEntry. <<< Resource test cases, automatically generated.
