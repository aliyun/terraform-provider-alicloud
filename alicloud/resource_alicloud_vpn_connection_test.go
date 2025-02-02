package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func equalSubnet(astr string, bstr string) bool {
	aArray := strings.Split(astr, ",")
	bArray := strings.Split(bstr, ",")
	if len(aArray) != len(bArray) {
		return false
	}

	for _, item := range aArray {
		if !strings.Contains(bstr, item) {
			return false
		}
	}
	return true
}

func testAccCheckVpnConnectionAttr(vpnConn *vpc.DescribeVpnConnectionResponse, localSubnet, remoteSubnet string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !equalSubnet(vpnConn.LocalSubnet, localSubnet) {
			return WrapError(Error("wrong local subnet, expect %s, get %s", localSubnet, vpnConn.LocalSubnet))
		}

		if !equalSubnet(vpnConn.RemoteSubnet, remoteSubnet) {
			return WrapError(Error("wrong remote subnet, expect %s, get %s", remoteSubnet, vpnConn.RemoteSubnet))
		}

		return nil
	}
}
func TestAccAliCloudVPNConnectionBasic(t *testing.T) {
	var v vpc.DescribeVpnConnectionResponse

	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, testAccVpnConnectionCheckMap)

	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccVpnConnectionBaisc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnConnectionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, connectivity.VPNSingleConnectRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id":      "${alicloud_vpn_gateway.default.id}",
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"local_subnet":        []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"name":                "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"local_subnet.#":  "2",
						"remote_subnet.#": "2",
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
					"local_subnet": []string{"172.16.2.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"local_subnet.#":  "1",
						"remote_subnet.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_subnet": []string{"10.4.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"local_subnet.#":  "1",
						"remote_subnet.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]string{
						{
							"ike_auth_alg":  "md5",
							"ike_enc_alg":   "des",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "tf-testvpn2",
							"ike_pfs":       "group1",
							"ike_remote_id": "testbob2",
							"ike_local_id":  "testalice2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ike_config.#":               "1",
						"ike_config.0.ike_auth_alg":  "md5",
						"ike_config.0.ike_enc_alg":   "des",
						"ike_config.0.ike_version":   "ikev2",
						"ike_config.0.ike_mode":      "main",
						"ike_config.0.ike_lifetime":  "86400",
						"ike_config.0.psk":           "tf-testvpn2",
						"ike_config.0.ike_pfs":       "group1",
						"ike_config.0.ike_remote_id": "testbob2",
						"ike_config.0.ike_local_id":  "testalice2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]string{
						{
							"ipsec_pfs":      "group5",
							"ipsec_enc_alg":  "des",
							"ipsec_auth_alg": "md5",
							"ipsec_lifetime": "8640",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipsec_config.#":                "1",
						"ipsec_config.0.ipsec_pfs":      "group5",
						"ipsec_config.0.ipsec_enc_alg":  "des",
						"ipsec_config.0.ipsec_auth_alg": "md5",
						"ipsec_config.0.ipsec_lifetime": "8640",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               "${var.name}",
					"local_subnet":       []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":      []string{"10.0.0.0/24", "10.0.1.0/24"},
					"effect_immediately": REMOVEKEY,
					"ike_config":         REMOVEKEY,
					"ipsec_config":       REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"effect_immediately": "false",
						"local_subnet.#":     "2",
						"remote_subnet.#":    "2",
					}),
				),
			},
		},
	})

}

var testAccVpnConnectionCheckMap = map[string]string{
	"vpn_gateway_id":                CHECKSET,
	"customer_gateway_id":           CHECKSET,
	"local_subnet.#":                "2",
	"remote_subnet.#":               "2",
	"effect_immediately":            "false",
	"ike_config.#":                  "1",
	"ike_config.0.ike_auth_alg":     "md5",
	"ike_config.0.ike_enc_alg":      "aes",
	"ike_config.0.ike_version":      "ikev2",
	"ike_config.0.ike_mode":         "main",
	"ike_config.0.ike_lifetime":     "86400",
	"ike_config.0.psk":              CHECKSET,
	"ike_config.0.ike_pfs":          "group2",
	"ike_config.0.ike_remote_id":    CHECKSET,
	"ike_config.0.ike_local_id":     CHECKSET,
	"ipsec_config.#":                "1",
	"ipsec_config.0.ipsec_pfs":      "group2",
	"ipsec_config.0.ipsec_enc_alg":  "aes",
	"ipsec_config.0.ipsec_auth_alg": "md5",
	"ipsec_config.0.ipsec_lifetime": "86400",
}

var resourceVpnConnectionConfigDependence = func(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

//data "alicloud_alb_zones" "default" {
//}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "me-east-1a"
	name = "${var.name}"
}

resource "alicloud_vswitch" "default1" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.17.0.0/22"
	availability_zone = "me-east-1a"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "default" {
	name = "${var.name}"
	ip_address = "42.104.22.210"
	description = "testAccVpnConnectionDesc"
}

resource "alicloud_vpn_customer_gateway" "default1" {
	name = "${var.name}"
	ip_address = "42.104.22.211"
	description = "testAccVpnConnectionDesc"
}

`, name)
}

func TestAccAliCloudVPNConnection_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnConnectionMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sonnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnConnectionBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, connectivity.VPNSingleConnectRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id":      "${alicloud_vpn_gateway.default.id}",
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"local_subnet":        []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"name":                "${var.name}",
					"effect_immediately":  "false",
					"ike_config": []map[string]string{
						{
							"ike_auth_alg":  "md5",
							"ike_enc_alg":   "des",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "tf-testvpn2",
							"ike_pfs":       "group1",
							"ike_remote_id": "testbob2",
							"ike_local_id":  "testalice2",
						},
					},
					"ipsec_config": []map[string]string{
						{
							"ipsec_pfs":      "group5",
							"ipsec_enc_alg":  "des",
							"ipsec_auth_alg": "md5",
							"ipsec_lifetime": "8640",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_id":      CHECKSET,
						"customer_gateway_id": CHECKSET,
						"name":                name,
						"local_subnet.#":      "2",
						"remote_subnet.#":     "2",
						"effect_immediately":  "false",
						"ike_config.#":        "1",
						"ipsec_config.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]string{
						{
							"enable":   "true",
							"dip":      "10.0.0.1",
							"sip":      "192.168.1.1",
							"interval": "10",
							"retry":    "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dpd": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_nat_traversal": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_nat_traversal": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_config": []map[string]string{
						{
							"enable":       "true",
							"local_asn":    "45014",
							"tunnel_cidr":  "169.254.11.0/30",
							"local_bgp_ip": "169.254.11.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd":           "false",
					"enable_nat_traversal": "false",
					"health_check_config": []map[string]string{
						{
							"enable": "false",
						},
					},
					"bgp_config": []map[string]string{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_config.#":          "1",
						"health_check_config.#": "1",
						"enable_dpd":            "false",
						"enable_nat_traversal":  "false",
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

var AlicloudVpnConnectionMap3 = map[string]string{}

func AlicloudVpnConnectionBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

//data "alicloud_alb_zones" "default" {
//}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "me-east-1a"
	name = "${var.name}"
}

resource "alicloud_vswitch" "default1" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.17.0.0/22"
	availability_zone = "me-east-1a"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    network_type =                 "public"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "default" {
	name = "${var.name}"
	ip_address = "42.104.22.211"
	description = "testAccVpnConnectionDesc"
    asn                   = "2224"
}

`, name)
}

// Test VPNGateway VpnConnection. >>> Resource test cases, automatically generated.
// Case VpnConnection测试用例---双Tunnel_关闭bgp 5678
func TestAccAliCloudVPNGatewayVpnConnection_basic5678(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap5678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence5678)
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
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"vpn_connection_name": name,
					"network_type":        "public",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "master",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_id": CHECKSET,
						//"local_subnet":        "3.0.0.0/24",
						//"remote_subnet":       "3.0.1.0/24",
						"vpn_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_tunnels_bgp": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type": "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"enable_tunnels_bgp": "false",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"enable_tunnels_bgp": "false",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_tunnels_bgp": "true",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "master",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
					},
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"vpn_connection_name": name + "_update",
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp":             "true",
						"tunnel_options_specification.#": "2",
						"vpn_gateway_id":                 CHECKSET,
						"vpn_connection_name":            name + "_update",
						//"local_subnet":                   "3.0.0.0/24",
						//"remote_subnet":                  "3.0.1.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route", "network_type"},
			},
		},
	})
}

var AlicloudVPNGatewayVpnConnectionMap5678 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudVPNGatewayVpnConnectionBasicDependence5678(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "spec" {
  default = "5"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 1)
  zone_id           = data.alicloud_vpn_gateway_zones.default.ids.0
  vswitch_name      = var.name
}

data "alicloud_vswitches" "default2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.1
}

resource "alicloud_vswitch" "vswitch2" {
  count             = length(data.alicloud_vswitches.default2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id           = data.alicloud_vpn_gateway_zones.default.ids.1
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id2 = length(data.alicloud_vswitches.default2.ids) > 0 ? data.alicloud_vswitches.default2.ids[0] : concat(alicloud_vswitch.vswitch2.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "HA-VPN" {
 vpn_type                     = "Normal"
 disaster_recovery_vswitch_id = local.vswitch_id2
 vpn_gateway_name             = var.name

 vswitch_id   = local.vswitch_id
 auto_pay     = true
 vpc_id       = data.alicloud_vpcs.default.ids.0
 network_type = "public"
 payment_type = "Subscription"
 enable_ipsec = true
 bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.5"
  asn                   = "2225"
  customer_gateway_name = var.name

}


`, name)
}

// Case VpnConnection测试用例---双Tunnel 3783
func TestAccAliCloudVPNGatewayVpnConnection_basic3783(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap3783)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence3783)
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
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"vpn_connection_name": name,
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "master",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_id": CHECKSET,
						//"local_subnet":        "3.0.0.0/24",
						//"remote_subnet":       "3.0.1.0/24",
						"vpn_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_tunnels_bgp": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_tunnels_bgp": "true",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.changeCustomerGateway.id}",
							"role":                 "master",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345679",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.changeCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345679",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
					},
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"vpn_connection_name": name + "_update",
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp":             "true",
						"tunnel_options_specification.#": "2",
						"vpn_gateway_id":                 CHECKSET,
						"vpn_connection_name":            name + "_update",
						//"local_subnet":                   "3.0.0.0/24",
						//"remote_subnet":                  "3.0.1.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.changeCustomerGateway.id}",
							"role":                 "master",
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.changeCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "36000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "1640",
									"ipsec_pfs":      "group5",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route"},
			},
		},
	})
}

var AlicloudVPNGatewayVpnConnectionMap3783 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudVPNGatewayVpnConnectionBasicDependence3783(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "spec" {
  default = "5"
}

data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 1)
  zone_id           = data.alicloud_vpn_gateway_zones.default.ids.0
  vswitch_name      = var.name
}

data "alicloud_vswitches" "default2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.1
}

resource "alicloud_vswitch" "vswitch2" {
  count             = length(data.alicloud_vswitches.default2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id           = data.alicloud_vpn_gateway_zones.default.ids.1
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id2 = length(data.alicloud_vswitches.default2.ids) > 0 ? data.alicloud_vswitches.default2.ids[0] : concat(alicloud_vswitch.vswitch2.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "HA-VPN" {
 vpn_type                     = "Normal"
 disaster_recovery_vswitch_id = local.vswitch_id2
 vpn_gateway_name             = var.name

 vswitch_id   = local.vswitch_id
 auto_pay     = true
 vpc_id       = data.alicloud_vpcs.default.ids.0
 network_type = "public"
 payment_type = "Subscription"
 enable_ipsec = true
 bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.5"
  asn                   = "2224"
  customer_gateway_name = var.name

}

resource "alicloud_vpn_customer_gateway" "changeCustomerGateway" {
  description           = "changeCustomerGateway"
  ip_address            = "2.2.2.6"
  asn                   = "2225"
  customer_gateway_name = var.name

}


`, name)
}

// Case VpnConnection测试用例---单隧道-2 5663
func TestAccAliCloudVPNGatewayVpnConnection_basic5663(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap5663)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence5663)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPNSingleConnectRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "1.1.1.1",
							"retry":    "3",
							"sip":      "3.3.3.3",
							"interval": "3",
						},
					},
					"auto_config_route":  "true",
					"effect_immediately": "true",
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "1219002",
							"local_bgp_ip": "169.254.10.1",
							"tunnel_cidr":  "169.254.10.0/30",
						},
					},
					"vpn_gateway_id": "${alicloud_vpn_gateway.default.id}",
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group2",
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_lifetime": "86400",
						},
					},
					"enable_nat_traversal": "true",
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha1",
							"ike_local_id":  "localid1",
							"ike_enc_alg":   "aes",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "12345678",
							"ike_remote_id": "remoteId2",
							"ike_pfs":       "group2",
						},
					},
					"enable_dpd":          "true",
					"vpn_connection_name": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"local_subnet":         "4.1.0.0/24",
						"customer_gateway_id": CHECKSET,
						"auto_config_route":   "true",
						"effect_immediately":  "true",
						//"remote_subnet":        "4.0.1.0/24",
						"vpn_gateway_id":       CHECKSET,
						"enable_nat_traversal": "true",
						"enable_dpd":           "true",
						"vpn_connection_name":  name,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "1.1.1.1",
							"retry":    "3",
							"sip":      "3.3.3.3",
							"interval": "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "1219002",
							"local_bgp_ip": "169.254.10.1",
							"tunnel_cidr":  "169.254.10.0/30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group2",
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_lifetime": "86400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_nat_traversal": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_nat_traversal": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha1",
							"ike_local_id":  "localid1",
							"ike_enc_alg":   "aes",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "12345678",
							"ike_remote_id": "remoteId2",
							"ike_pfs":       "group2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dpd": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet": []string{"4.1.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_subnet": []string{"4.0.1.0/24", "10.0.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_subnet.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"dip":      "11.111.11.111",
							"retry":    "4",
							"sip":      "32.32.32.232",
							"interval": "10",
							"enable":   "true",
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
				Config: testAccConfig(map[string]interface{}{
					"bgp_config": []map[string]interface{}{
						{
							"local_asn":    "1219004",
							"local_bgp_ip": "169.254.20.1",
							"tunnel_cidr":  "169.254.20.0/30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group14",
							"ipsec_enc_alg":  "aes192",
							"ipsec_auth_alg": "md5",
							"ipsec_lifetime": "8640",
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
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "md5",
							"ike_local_id":  "localid2",
							"ike_enc_alg":   "aes192",
							"ike_version":   "ikev2",
							"ike_mode":      "aggressive",
							"ike_lifetime":  "8640",
							"psk":           "2222222",
							"ike_remote_id": "remoteid2",
							"ike_pfs":       "group1",
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
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_config": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_enc_alg":  "aes256",
							"ipsec_auth_alg": "sha256",
							"ipsec_lifetime": "27000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_nat_traversal": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_nat_traversal": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha256",
							"ike_local_id":  "2.2.2.2",
							"ike_enc_alg":   "aes256",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "27000",
							"psk":           "999222",
							"ike_remote_id": "remoyteid1",
							"ike_pfs":       "group14",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dpd": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "disabled",
							"ipsec_enc_alg":  "des",
							"ipsec_auth_alg": "sha384",
							"ipsec_lifetime": "3600",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha384",
							"ike_local_id":  "localid3",
							"ike_enc_alg":   "des",
							"ike_version":   "ikev2",
							"ike_lifetime":  "3600",
							"psk":           "8888888",
							"ike_remote_id": "remoteid2",
							"ike_pfs":       "group1",
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
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "1.1.1.1",
							"retry":    "3",
							"sip":      "3.3.3.3",
							"interval": "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "1219002",
							"local_bgp_ip": "169.254.10.1",
							"tunnel_cidr":  "169.254.10.0/30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group2",
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_lifetime": "86400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha1",
							"ike_local_id":  "localid1",
							"ike_enc_alg":   "aes",
							"ike_version":   "ikev2",
							"ike_lifetime":  "86400",
							"psk":           "12345678",
							"ike_remote_id": "remoteId2",
							"ike_pfs":       "group2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dpd": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_connection_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "1.1.1.1",
							"retry":    "3",
							"sip":      "3.3.3.3",
							"interval": "3",
						},
					},
					"auto_config_route":  "true",
					"effect_immediately": "true",
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "1219002",
							"local_bgp_ip": "169.254.10.1",
							"tunnel_cidr":  "169.254.10.0/30",
						},
					},
					"vpn_gateway_id": "${alicloud_vpn_gateway.default.id}",
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group2",
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_lifetime": "86400",
						},
					},
					"enable_nat_traversal": "true",
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha1",
							"ike_local_id":  "localid1",
							"ike_enc_alg":   "aes",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "12345678",
							"ike_remote_id": "remoteId2",
							"ike_pfs":       "group2",
						},
					},
					"enable_dpd":          "true",
					"vpn_connection_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"local_subnet":         "4.1.0.0/24",
						"customer_gateway_id": CHECKSET,
						"auto_config_route":   "true",
						"effect_immediately":  "true",
						//"remote_subnet":        "4.0.1.0/24",
						"vpn_gateway_id":       CHECKSET,
						"enable_nat_traversal": "true",
						"enable_dpd":           "true",
						"vpn_connection_name":  name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route"},
			},
		},
	})
}

var AlicloudVPNGatewayVpnConnectionMap5663 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudVPNGatewayVpnConnectionBasicDependence5663(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "spec" {
  default = "20"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "me-east-1a"
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = "me-east-1a"
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
  vpn_type         = "Normal"
  vpn_gateway_name = var.name

  vswitch_id   = local.vswitch_id
  auto_pay     = true
  vpc_id       = data.alicloud_vpcs.default.ids.0
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "default" {
  description           = "defaultCustomerGateway"
  ip_address            = "4.3.2.10"
  asn                   = "1219002"
  customer_gateway_name = var.name
}


`, name)
}

// Case VpnConnection测试用例---双Tunnel_关闭bgp 5678  twin
func TestAccAliCloudVPNGatewayVpnConnection_basic5678_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap5678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence5678)
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
					"enable_tunnels_bgp": "true",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"enable_nat_traversal": "true",
							"enable_dpd":           "true",
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                 "master",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
					},
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"vpn_connection_name": name,
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp":             "true",
						"tunnel_options_specification.#": "2",
						"vpn_gateway_id":                 CHECKSET,
						"vpn_connection_name":            name,
						"local_subnet.#":                 CHECKSET,
						"remote_subnet.#":                CHECKSET,
						"tags.%":                         "2",
						"tags.Created":                   "TF",
						"tags.For":                       "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route"},
			},
		},
	})
}

// Case VpnConnection测试用例---双Tunnel 3783  twin
func TestAccAliCloudVPNGatewayVpnConnection_basic3783_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap3783)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence3783)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_tunnels_bgp": "true",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "master",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"tunnel_cidr":  "169.254.30.0/30",
									"local_bgp_ip": "169.254.30.1",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "3600",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "16400",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultCustomerGateway.id}",
							"role":                "slave",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.40.1",
									"tunnel_cidr":  "169.254.40.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "27000",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "localid_tunnel2",
									"psk":          "12345678",
									"remote_id":    "remote24",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes256",
									"ipsec_lifetime": "2700",
									"ipsec_pfs":      "group14",
								},
							},
						},
					},
					"vpn_gateway_id":      "${alicloud_vpn_gateway.HA-VPN.id}",
					"vpn_connection_name": name,
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_tunnels_bgp":             "true",
						"tunnel_options_specification.#": "2",
						"vpn_gateway_id":                 CHECKSET,
						"vpn_connection_name":            name,
						"tags.%":                         "2",
						"tags.Created":                   "TF",
						"tags.For":                       "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route"},
			},
		},
	})
}

// Case VpnConnection测试用例---单隧道-2 5663  twin
func TestAccAliCloudVPNGatewayVpnConnection_basic5663_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnConnectionMap5663)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVpnConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnConnectionBasicDependence5663)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPNSingleConnectRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        []string{"3.0.0.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "1.1.1.1",
							"retry":    "3",
							"sip":      "3.3.3.3",
							"interval": "3",
						},
					},
					"auto_config_route":  "true",
					"effect_immediately": "true",
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "1219002",
							"local_bgp_ip": "169.254.10.1",
							"tunnel_cidr":  "169.254.10.0/30",
						},
					},
					"vpn_gateway_id": "${alicloud_vpn_gateway.default.id}",
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_pfs":      "group2",
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_lifetime": "86400",
						},
					},
					"enable_nat_traversal": "true",
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg":  "sha1",
							"ike_local_id":  "localid1",
							"ike_enc_alg":   "aes",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "12345678",
							"ike_remote_id": "remoteId2",
							"ike_pfs":       "group2",
						},
					},
					"enable_dpd":          "true",
					"vpn_connection_name": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"local_subnet":         "4.1.0.0/24",
						"customer_gateway_id": CHECKSET,
						"auto_config_route":   "true",
						"effect_immediately":  "true",
						//"remote_subnet":        "4.0.1.0/24",
						"vpn_gateway_id":       CHECKSET,
						"enable_nat_traversal": "true",
						"enable_dpd":           "true",
						"vpn_connection_name":  name,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config_route"},
			},
		},
	})
}

// Test VPNGateway VpnConnection. <<< Resource test cases, automatically generated.
