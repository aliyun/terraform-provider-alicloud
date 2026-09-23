package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudVPNGatewayVcoRoute_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vco_route.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVcoRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVcoRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvcoroute%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVcoRouteBasicDependence0)
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
					"route_dest":        "192.168.10.0/24",
					"next_hop":          "${alicloud_cen_transit_router_vpn_attachment.default.vpn_id}",
					"vpn_connection_id": "${alicloud_cen_transit_router_vpn_attachment.default.vpn_id}",
					"weight":            "100",
					"overlay_mode":      "Ipsec",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_dest":        "192.168.10.0/24",
						"next_hop":          CHECKSET,
						"vpn_connection_id": CHECKSET,
						"weight":            "100",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overlay_mode"},
			},
		},
	})
}

var AlicloudVPNGatewayVcoRouteMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVPNGatewayVcoRouteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
}
resource "alicloud_cen_transit_router" "default" {
	cen_id = alicloud_cen_instance.default.id
	transit_router_description = "desd"
	transit_router_name = var.name
}
resource "alicloud_cen_transit_router_cidr" "default" {
  transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = var.name
  description              = var.name
  publish_cidr_route       = true
}

resource "alicloud_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "43.${100 + tonumber(substr(var.name, -5, 2)) %% 100}.${100 + tonumber(substr(var.name, -3, 2)) %% 100}.${100 + tonumber(substr(var.name, -1, 1)) %% 10}"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
}
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  network_type       = "public"
  local_subnet       = "0.0.0.0/0"
  remote_subnet      = "0.0.0.0/0"
  effect_immediately = false
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.default.id
    role                 = "master"
    tunnel_index         = 1
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "des"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-testvpn-master"
      ike_pfs      = "group1"
      remote_id    = "testbob-master"
      local_id     = "testalice-master"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "des"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.default.id
    role                 = "slave"
    tunnel_index         = 2
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "des"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-testvpn-slave"
      ike_pfs      = "group1"
      remote_id    = "testbob-slave"
      local_id     = "testalice-slave"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "des"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
  vpn_attachment_name = var.name
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
	auto_publish_route_enabled = false
	transit_router_attachment_description = var.name
	transit_router_attachment_name = var.name
	cen_id = alicloud_cen_transit_router.default.cen_id
	transit_router_id = alicloud_cen_transit_router_cidr.default.transit_router_id
	vpn_id = alicloud_vpn_gateway_vpn_attachment.default.id
}
`, name)
}

func TestUnitAlicloudVPNGatewayVcoRouteBasicDependenceUsesUniqueIP(t *testing.T) {
	config := AlicloudVPNGatewayVcoRouteBasicDependence0("tf-testaccvpngatewayvcoroute12345")
	if strings.Contains(config, "42.104.22.210") {
		t.Fatal("generated configuration must not contain the fixed customer gateway IP")
	}
	if !strings.Contains(config, "tonumber(substr(var.name") {
		t.Fatal("customer gateway IP must derive from var.name")
	}
	for _, legacy := range []string{"\n  customer_gateway_id =", "\n  ike_config {", "\n  ipsec_config {", "\n  bgp_config {"} {
		if strings.Contains(config, legacy) {
			t.Fatalf("generated configuration must not contain legacy single-tunnel setting %q", strings.TrimSpace(legacy))
		}
	}
	if got := strings.Count(config, "tunnel_options_specification {"); got != 2 {
		t.Fatalf("expected two tunnel options specifications, got %d", got)
	}
	for _, expected := range []string{"tunnel_index         = 1", "tunnel_index         = 2", `role                 = "master"`, `role                 = "slave"`} {
		if !strings.Contains(config, expected) {
			t.Fatalf("generated configuration must contain %q", expected)
		}
	}
	if strings.Contains(config, "zone {") {
		t.Fatal("transit router VPN attachment must not configure a zone for dual-tunnel VCO")
	}
}

func TestUnitWaitForVpnGatewayVcoRouteDeleted(t *testing.T) {
	exists := []bool{true, false, true, false, false, false}
	refreshCalls := 0
	refresh := func() (interface{}, string, error) {
		if exists[refreshCalls] {
			refreshCalls++
			return map[string]interface{}{"State": "published"}, "published", nil
		}
		refreshCalls++
		return nil, "", nil
	}

	err := waitForVpnGatewayVcoRouteDeleted(time.Second, time.Millisecond, refresh)
	assert.NoError(t, err)
	assert.Equal(t, len(exists), refreshCalls, "delete must require three consecutive not-found responses")
}

// lintignore: R001
func TestUnitAccAlicloudVpnGatewayVcoRoute(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_vpn_gateway_vco_route"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_vpn_gateway_vco_route"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"route_dest":        "CreateVpnGatewayVcoRouteValue",
		"next_hop":          "CreateVpnGatewayVcoRouteValue",
		"vpn_connection_id": "CreateVpnGatewayVcoRouteValue",
		"weight":            100,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"VcoRouteEntries": []interface{}{
			map[string]interface{}{
				"RouteDest":       "CreateVpnGatewayVcoRouteValue",
				"State":           "published",
				"NextHop":         "CreateVpnGatewayVcoRouteValue",
				"RouteEntryType":  "CreateVpnGatewayVcoRouteValue",
				"VpnConnectionId": "CreateVpnGatewayVcoRouteValue",
				"Weight":          100,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"VpnConnectionId": "CreateVpnGatewayVcoRouteValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_vpn_gateway_vco_route", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpnGatewayVcoRouteCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVcoRouteEntry" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpnGatewayVcoRouteCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpn_gateway_vco_route"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_vpn_gateway_vco_route", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpn_gateway_vco_route"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeVcoRouteEntries" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpnGatewayVcoRouteRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpnGatewayVcoRouteDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_vpn_gateway_vco_route", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpn_gateway_vco_route"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVcoRouteEntry" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpnGatewayVcoRouteDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
