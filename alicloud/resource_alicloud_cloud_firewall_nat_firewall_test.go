package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// Test CloudFirewall NatFirewall. >>> Resource test cases, automatically generated.
// Case 云防火墙资源测试-v1 6822
func TestAccAliCloudCloudFirewallNatFirewall_basic6822(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":         "normal",
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch":        "close",
					"vswitch_cidr":           "172.16.5.0/24",
					"vswitch_id":             "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
					"firewall_service_mode":  "PrimaryStandby",
					"firewall_service_zones": []interface{}{"cn-hangzhou-i"},
					"fw_vswitch_zone_id":     "cn-hangzhou-i",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "closed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "closed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_vswitch.defaultp4O7qi.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "fw_vswitch_zone_id", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

var AlicloudCloudFirewallNatFirewallMap6822 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCloudFirewallNatFirewallBasicDependence6822(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultikZ0gD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultp4O7qi" {
  vpc_id       = alicloud_vpc.defaultikZ0gD.id
  cidr_block   = "172.16.6.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_nat_gateway" "default2iRZpC" {
  eip_bind_mode    = "MULTI_BINDED"
  vpc_id           = alicloud_vpc.defaultikZ0gD.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultp4O7qi.id
  nat_type         = "Enhanced"
  network_type = "internet"
}

resource "alicloud_eip_address" "defaultyiRwgs" {
  address_name = var.name
}

resource "alicloud_eip_association" "defaults2MTuO" {
  instance_id   = alicloud_nat_gateway.default2iRZpC.id
  allocation_id = alicloud_eip_address.defaultyiRwgs.id
  mode          = "NAT"
  instance_type = "Nat"
}

resource "alicloud_snat_entry" "defaultAKE43g" {
  snat_ip           = alicloud_eip_address.defaultyiRwgs.ip_address
  snat_table_id     = alicloud_nat_gateway.default2iRZpC.snat_table_ids
  source_vswitch_id = alicloud_vswitch.defaultp4O7qi.id
}


`, name)
}

// Case 云防火墙资源测试-v1 6822  twin
func TestAccAliCloudCloudFirewallNatFirewall_basic6822_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "fw_vswitch_zone_id", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

// Case 云防火墙资源测试-v1 6822  raw
func TestAccAliCloudCloudFirewallNatFirewall_basic6822_raw(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              CHECKSET,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "fw_vswitch_zone_id", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

// Test CloudFirewall NatFirewall. <<< Resource test cases, automatically generated.

// TestUnitCloudFirewallNatFirewallSchema is a no-cloud regression for the three
// parameters mapped for property completeness (FirewallServiceMode /
// FirewallServiceZones / FwVswitchZoneId). It verifies the schema definitions
// without any cloud call or credentials.
func TestUnitCloudFirewallNatFirewallSchema(t *testing.T) {
	s := resourceAliCloudCloudFirewallNatFirewall().Schema

	mode, ok := s["firewall_service_mode"]
	assert.True(t, ok, "firewall_service_mode must be present in schema")
	assert.Equal(t, schema.TypeString, mode.Type)
	assert.True(t, mode.Optional)
	assert.NotNil(t, mode.ValidateFunc, "firewall_service_mode must carry an enum ValidateFunc")
	assert.True(t, mode.ForceNew, "firewall_service_mode is create-only (SwitchSecurityProxy does not accept it) so ForceNew")

	zones, ok := s["firewall_service_zones"]
	assert.True(t, ok, "firewall_service_zones must be present in schema")
	assert.Equal(t, schema.TypeList, zones.Type)
	assert.True(t, zones.Optional)
	zonesElem, ok := zones.Elem.(*schema.Schema)
	assert.True(t, ok, "firewall_service_zones elem must be *schema.Schema")
	assert.Equal(t, schema.TypeString, zonesElem.Type)
	assert.True(t, zones.ForceNew, "firewall_service_zones is create-only (SwitchSecurityProxy does not accept it) so ForceNew")

	fwVswitchZoneId, ok := s["fw_vswitch_zone_id"]
	assert.True(t, ok, "fw_vswitch_zone_id must be present in schema")
	assert.Equal(t, schema.TypeString, fwVswitchZoneId.Type)
	assert.True(t, fwVswitchZoneId.Optional)
	assert.True(t, fwVswitchZoneId.ForceNew, "fw_vswitch_zone_id is create-only (not read back) so ForceNew")
}

// TestUnitCloudFirewallNatFirewallCreateRequest is a no-cloud regression for
// the CreateSecurityProxy request marshalling. It feeds a populated
// ResourceData through expandCloudFirewallNatFirewallCreateRequest and asserts
// the three newly mapped parameters reach the request body with the correct
// names and shapes (FirewallServiceZones as a repeatList []string).
func TestUnitCloudFirewallNatFirewallCreateRequest(t *testing.T) {
	p := Provider().ResourcesMap
	d, err := schema.InternalMap(p["alicloud_cloud_firewall_nat_firewall"].Schema).Data(nil, nil)
	assert.Nil(t, err)

	// Literal d.Set keys keep tflint rule R001 (Set key must be a string
	// literal) green; a range loop with a variable key would fail it.
	assert.Nil(t, d.Set("nat_gateway_id", "ng-test-id"))
	assert.Nil(t, d.Set("nat_route_entry_list", []interface{}{
		map[string]interface{}{
			"nexthop_id":       "ng-test-id",
			"destination_cidr": "0.0.0.0/0",
			"nexthop_type":     "NatGateway",
			"route_table_id":   "rt-test-id",
		},
	}))
	assert.Nil(t, d.Set("proxy_name", "tf-test-proxy"))
	assert.Nil(t, d.Set("region_no", "cn-hangzhou"))
	assert.Nil(t, d.Set("vpc_id", "vpc-test-id"))
	assert.Nil(t, d.Set("firewall_switch", "close"))
	assert.Nil(t, d.Set("strict_mode", 0))
	assert.Nil(t, d.Set("vswitch_cidr", "172.16.5.0/24"))
	assert.Nil(t, d.Set("vswitch_id", "vsw-test-id"))
	assert.Nil(t, d.Set("firewall_service_mode", "PrimaryStandby"))
	assert.Nil(t, d.Set("firewall_service_zones", []interface{}{"cn-hangzhou-i", "cn-hangzhou-h"}))
	assert.Nil(t, d.Set("fw_vswitch_zone_id", "cn-hangzhou-i"))

	request := expandCloudFirewallNatFirewallCreateRequest(d)

	// Newly mapped parameters.
	assert.Equal(t, "PrimaryStandby", request["FirewallServiceMode"])
	assert.Equal(t, []string{"cn-hangzhou-i", "cn-hangzhou-h"}, request["FirewallServiceZones"])
	assert.Equal(t, "cn-hangzhou-i", request["FwVswitchZoneId"])

	// Existing parameters still marshalled (regression guard).
	assert.Equal(t, "ng-test-id", request["NatGatewayId"])
	assert.Equal(t, "vpc-test-id", request["VpcId"])
	assert.Equal(t, "tf-test-proxy", request["ProxyName"])
	assert.Equal(t, "cn-hangzhou", request["RegionNo"])
	assert.Equal(t, "close", request["FirewallSwitch"])
	assert.Equal(t, "vsw-test-id", request["VswitchId"])
	assert.Equal(t, "172.16.5.0/24", request["VswitchCidr"])
	natRouteEntryList, ok := request["NatRouteEntryList"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, natRouteEntryList, 1)
}

// TestUnitCloudFirewallNatFirewallCreateRequestOmitted asserts that when the
// three new optional parameters are not set, they are absent from the Create
// request body (d.GetOk returns false), preserving prior behaviour for
// existing configs.
func TestUnitCloudFirewallNatFirewallCreateRequestOmitted(t *testing.T) {
	p := Provider().ResourcesMap
	d, err := schema.InternalMap(p["alicloud_cloud_firewall_nat_firewall"].Schema).Data(nil, nil)
	assert.Nil(t, err)

	// Literal d.Set keys keep tflint rule R001 green (see the populated test
	// above). The three new parameters are intentionally left unset here.
	assert.Nil(t, d.Set("nat_gateway_id", "ng-test-id"))
	assert.Nil(t, d.Set("nat_route_entry_list", []interface{}{
		map[string]interface{}{
			"nexthop_id":       "ng-test-id",
			"destination_cidr": "0.0.0.0/0",
			"nexthop_type":     "NatGateway",
			"route_table_id":   "rt-test-id",
		},
	}))
	assert.Nil(t, d.Set("proxy_name", "tf-test-proxy"))
	assert.Nil(t, d.Set("region_no", "cn-hangzhou"))
	assert.Nil(t, d.Set("vpc_id", "vpc-test-id"))

	request := expandCloudFirewallNatFirewallCreateRequest(d)

	_, ok := request["FirewallServiceMode"]
	assert.False(t, ok, "FirewallServiceMode must be absent when unset")
	_, ok = request["FirewallServiceZones"]
	assert.False(t, ok, "FirewallServiceZones must be absent when unset")
	_, ok = request["FwVswitchZoneId"]
	assert.False(t, ok, "FwVswitchZoneId must be absent when unset")
}
