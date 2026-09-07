package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

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

func TestAccAliCloudCENTransitRouterVpnAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCENTransitRouterVpnAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitroutervpnattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCENTransitRouterVpnAttachmentBasicDependence0)
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
					"vpn_owner_id":                          "${data.alicloud_account.default.id}",
					"transit_router_vpn_attachment_name":    "${var.name}",
					"auto_publish_route_enabled":            "false",
					"transit_router_attachment_description": "${var.name}",
					"vpn_id":                                "${alicloud_vpn_gateway_vpn_attachment.default.id}",
					"cen_id":                                "${alicloud_cen_transit_router.default.cen_id}",
					"transit_router_id":                     "${alicloud_cen_transit_router_cidr.default.transit_router_id}",
					"order_type":                            "PayByCenOwner",
					// Dual-tunnel Vco derives zones internally; passing any
					// explicit zone here trips OperationUnsupported.VcoTunnelNotMatchZoneParam.
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_owner_id":                          CHECKSET,
						"transit_router_vpn_attachment_name":    name,
						"auto_publish_route_enabled":            "false",
						"transit_router_attachment_description": name,
						"transit_router_id":                     CHECKSET,
						"vpn_id":                                CHECKSET,
						"order_type":                            "PayByCenOwner",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_vpn_attachment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_vpn_attachment_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudCENTransitRouterVpnAttachment_zoneMismatch asserts the
// negative case: when the underlying alicloud_vpn_gateway_vpn_attachment
// (Vco) is created in dual-tunnel mode via tunnel_options_specification,
// CBN derives the zones internally and rejects any attempt to also pass an
// explicit `zone` block on alicloud_cen_transit_router_vpn_attachment.
// The Create API returns OperationUnsupported.VcoTunnelNotMatchZoneParam;
// this test pins that behaviour so regressions in the provider (e.g. quietly
// dropping the zone field before it reaches the API) surface as test
// failures instead of silent drifts.
func TestAccAliCloudCENTransitRouterVpnAttachment_zoneMismatch(t *testing.T) {
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentrvpnattzm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCENTransitRouterVpnAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// IDRefreshName is intentionally omitted: the attachment is never
		// created (Create is expected to fail), so SDK's refresh probe on
		// that ID would report "ID-only refresh check never ran.".
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":            "${alicloud_cen_transit_router.default.cen_id}",
					"transit_router_id": "${alicloud_cen_transit_router_cidr.default.transit_router_id}",
					"vpn_id":            "${alicloud_vpn_gateway_vpn_attachment.default.id}",
					"zone": []map[string]interface{}{
						{
							"zone_id": "${data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0}",
						},
					},
				}),
				// CBN rejects the request with the error below when a
				// dual-tunnel Vco meets an explicit zone parameter. The
				// regex matches the error code so the assertion is stable
				// across minor message wording changes.
				ExpectError: regexp.MustCompile(`OperationUnsupported\.VcoTunnelNotMatchZoneParam`),
			},
		},
	})
}

var AlicloudCENTransitRouterVpnAttachmentMap0 = map[string]string{
	"auto_publish_route_enabled": CHECKSET,
	"status":                     CHECKSET,
	"vpn_id":                     CHECKSET,
	"vpn_owner_id":               CHECKSET,
	"zone.#":                     CHECKSET,
	"transit_router_id":          CHECKSET,
}

func AlicloudCENTransitRouterVpnAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_cen_instance" "default" {
		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
		cen_id = alicloud_cen_instance.default.id
		transit_router_description = "desd"
		transit_router_name = var.name
	}

	resource "alicloud_vpn_customer_gateway" "default" {
  		customer_gateway_name = var.name
  		# Spread across 42.104.100-199.100-199 so repeated test runs do not
  		# collide with previously leaked customer gateways in the account.
  		ip_address            = "42.104.${100 + tonumber(substr(var.name, -3, 2)) %% 100}.${100 + tonumber(substr(var.name, -2, 2)) %% 100}"
  		asn                   = "45014"
  		description           = "testAccVpnConnectionDesc"
	}

	resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  		network_type        = "public"
  		local_subnet        = "0.0.0.0/0"
  		remote_subnet       = "0.0.0.0/0"
  		effect_immediately  = false
	tunnel_options_specification {
		customer_gateway_id = alicloud_vpn_customer_gateway.default.id
		role                = "master"
		tunnel_index        = 1
			enable_dpd          = true
			enable_nat_traversal = true
			tunnel_ike_config {
				ike_auth_alg = "md5"
				ike_enc_alg  = "des"
				ike_version  = "ikev2"
				ike_mode     = "main"
				ike_lifetime = 86400
				psk          = "tf-testvpn3"
				ike_pfs      = "group1"
				remote_id    = "testbob3"
				local_id     = "testalice3"
			}
			tunnel_ipsec_config {
				ipsec_pfs      = "group5"
				ipsec_enc_alg  = "des"
				ipsec_auth_alg = "md5"
				ipsec_lifetime = 86400
			}
		}
	tunnel_options_specification {
		customer_gateway_id = alicloud_vpn_customer_gateway.default.id
		role                = "slave"
		tunnel_index        = 2
			enable_dpd          = true
			enable_nat_traversal = true
			tunnel_ike_config {
				ike_auth_alg = "md5"
				ike_enc_alg  = "des"
				ike_version  = "ikev2"
				ike_mode     = "main"
				ike_lifetime = 86400
				psk          = "tf-testvpn2"
				ike_pfs      = "group1"
				remote_id    = "testbob2"
				local_id     = "testalice2"
			}
			tunnel_ipsec_config {
				ipsec_pfs      = "group5"
				ipsec_enc_alg  = "des"
				ipsec_auth_alg = "md5"
				ipsec_lifetime = 86400
			}
		}
  		vpn_attachment_name  = var.name
	}

	resource "alicloud_cen_transit_router_cidr" "default" {
		transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  		cidr                     = "192.168.0.0/16"
  		transit_router_cidr_name = var.name
  		description              = var.name
  		publish_cidr_route       = false
	}
	
	data "alicloud_cen_transit_router_available_resources" "default" {
	}
`, name)
}

// TestAccAliCloudCENTransitRouterVpnAttachment_crossAccountOrderType covers the
// cross-account order_type flow:
//   - Account A (profile TerraformUT, provider alias "a") owns the CEN instance and
//     creates the alicloud_cen_transit_router_vpn_attachment whose order_type is
//     the payer chosen by the grant.
//   - Account B (profile TerraformTest, default provider) owns the VPN gateway and
//     the alicloud_cen_transit_router_grant_attachment that authorizes A to attach it.
//
// Changing the payer requires B updating the grant first and then A updating the
// attachment; the test drives that sequencing via depends_on.
func TestAccAliCloudCENTransitRouterVpnAttachment_crossAccountOrderType(t *testing.T) {
	// Load profile credentials up front so the Config strings below can embed
	// AK/SK directly. Terraform SDK evaluates Step.Config eagerly when the
	// Steps slice is built, so the creds MUST be populated before that point
	// rather than inside PreCheck (which runs later).
	testAccPreCheckCENCrossAccount(t)
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCENTransitRouterVpnAttachmentMap0)
	providerFactories, factoryProviders := cenCrossAccountProviderFactories()
	// The attachment under test is owned by account A (TerraformUT); pick the
	// factory-created provider with that AK so describe/CheckDestroy hits A.
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		client := cenCrossAccountClientByAK(*factoryProviders, sharedCENCrossAccountCreds.utAK)
		if client == nil {
			return &CbnService{}
		}
		return &CbnService{client}
	}, "DescribeCenTransitRouterVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentrvpnatt%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCENTransitRouterVpnAttachmentCrossAccountConfig(name, "PayByCenOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByCenOwner",
					}),
				),
			},
			{
				Config: testAccCENTransitRouterVpnAttachmentCrossAccountConfig(name, "PayByResourceOwner"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type": "PayByResourceOwner",
					}),
				),
			},
			// ImportState is intentionally omitted: Terraform SDK v1's import
			// step does not thread ProviderFactories through to the aliased
			// "alicloud.a" provider, so the refresh fails with "unknown
			// provider \"alicloud\"". The cross-account create + order_type
			// update flow is already validated by the two steps above.
		},
	})
}

func testAccCENTransitRouterVpnAttachmentCrossAccountConfig(name, orderType string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
%s

data "alicloud_account" "b" {}

data "alicloud_account" "a" {
  provider = alicloud.a
}

# --- Account A: CEN ---
resource "alicloud_cen_instance" "default" {
  provider          = alicloud.a
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  provider                   = alicloud.a
  cen_id                     = alicloud_cen_instance.default.id
  transit_router_description = "desd"
  transit_router_name        = var.name
}

resource "alicloud_cen_transit_router_cidr" "default" {
  provider                 = alicloud.a
  transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = var.name
  description              = var.name
  publish_cidr_route       = false
}

# Single-tunnel VPN attachments require an explicit zone matching the tunnel
# layout; pull the first master zone available in account A's region.
data "alicloud_cen_transit_router_available_resources" "default" {
  provider = alicloud.a
}

# --- Account B: VPN + Grant ---
resource "alicloud_vpn_customer_gateway" "default" {
  customer_gateway_name = var.name
  # Spread the test-generated IP across a wider range than basic0's
  # 42.104.22.100-119 pool so repeated cross-account runs do not collide
  # with account B's accumulated leftover gateways.
  ip_address  = "42.104.${100 + tonumber(substr(var.name, -3, 2)) %% 100}.${100 + tonumber(substr(var.name, -2, 2)) %% 100}"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
}

resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  effect_immediately  = false
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  vpn_attachment_name = var.name
}

resource "alicloud_cen_transit_router_grant_attachment" "default" {
  instance_type = "VPN"
  instance_id   = alicloud_vpn_gateway_vpn_attachment.default.id
  cen_owner_id  = data.alicloud_account.a.id
  cen_id        = alicloud_cen_instance.default.id
  order_type    = %q
}

# --- Account A: attach B's VPN into A's CEN ---
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  provider                              = alicloud.a
  vpn_owner_id                          = data.alicloud_account.b.id
  transit_router_vpn_attachment_name    = var.name
  auto_publish_route_enabled            = false
  transit_router_attachment_description = var.name
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.default.id
  cen_id                                = alicloud_cen_transit_router.default.cen_id
  transit_router_id                     = alicloud_cen_transit_router_cidr.default.transit_router_id
  order_type                            = %q
  zone {
    zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
  }
  depends_on = [alicloud_cen_transit_router_grant_attachment.default]
}
`, name, cenCrossAccountProviderBlocks(), orderType, orderType)
}

func TestAccAliCloudCENTransitRouterVpnAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCENTransitRouterVpnAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitroutervpnattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCENTransitRouterVpnAttachmentBasicDependence0)
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
					// Dual-tunnel Vco derives zones internally; passing any
					// explicit zone here trips OperationUnsupported.VcoTunnelNotMatchZoneParam.
					"transit_router_id": "${alicloud_cen_transit_router_cidr.default.transit_router_id}",
					"vpn_id":            "${alicloud_vpn_gateway_vpn_attachment.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterVpnAttachment",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_id":            CHECKSET,
						"transit_router_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "TransitRouterVpnAttachment",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_publish_route_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "TransitRouterVpnAttachment_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "TransitRouterVpnAttachment_Update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// lintignore: R001
func TestUnitAccAlicloudCenTransitRouterVpnAttachment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"vpn_owner_id": "CreateCenTransitRouterVpnAttachmentValue",
		"zone": []map[string]interface{}{
			{
				"zone_id": "CreateCenTransitRouterVpnAttachmentValue",
			},
		},
		"transit_router_attachment_name":        "CreateCenTransitRouterVpnAttachmentValue",
		"auto_publish_route_enabled":            false,
		"transit_router_attachment_description": "CreateCenTransitRouterVpnAttachmentValue",
		"vpn_id":                                "CreateCenTransitRouterVpnAttachmentValue",
		"cen_id":                                "CreateCenTransitRouterVpnAttachmentValue",
		"transit_router_id":                     "CreateCenTransitRouterVpnAttachmentValue",
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
		"TransitRouterAttachments": []interface{}{
			map[string]interface{}{
				"CreationTime":                       "CreateCenTransitRouterVpnAttachmentValue",
				"Status":                             "Attached",
				"TransitRouterAttachmentId":          "CreateCenTransitRouterVpnAttachmentValue",
				"TransitRouterId":                    "CreateCenTransitRouterVpnAttachmentValue",
				"VpnOwnerId":                         "CreateCenTransitRouterVpnAttachmentValue",
				"VpnId":                              "CreateCenTransitRouterVpnAttachmentValue",
				"TransitRouterAttachmentDescription": "CreateCenTransitRouterVpnAttachmentValue",
				"VpnRegionId":                        "CreateCenTransitRouterVpnAttachmentValue",
				"AutoPublishRouteEnabled":            false,
				"TransitRouterAttachmentName":        "CreateCenTransitRouterVpnAttachmentValue",
				"Zones": []interface{}{
					map[string]interface{}{
						"ZoneId": "CreateCenTransitRouterVpnAttachmentValue",
					},
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"TransitRouterAttachmentId": "CreateCenTransitRouterVpnAttachmentValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cen_transit_router_vpn_attachment", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTransitRouterVpnAttachmentCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateTransitRouterVpnAttachment" {
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
		err := resourceAliCloudCenTransitRouterVpnAttachmentCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTransitRouterVpnAttachmentUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{

		"transit_router_attachment_name":        "UpdateCenTransitRouterVpnAttachmentValue",
		"auto_publish_route_enabled":            true,
		"transit_router_attachment_description": "UpdateCenTransitRouterVpnAttachmentValue",
	}
	diff, err := newInstanceDiff("alicloud_cen_transit_router_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"TransitRouterAttachments": []interface{}{
			map[string]interface{}{
				"TransitRouterAttachmentDescription": "UpdateCenTransitRouterVpnAttachmentValue",
				"AutoPublishRouteEnabled":            true,
				"TransitRouterAttachmentName":        "UpdateCenTransitRouterVpnAttachmentValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateTransitRouterVpnAttachmentAttribute" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCenTransitRouterVpnAttachmentUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_cen_transit_router_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListTransitRouterVpnAttachments" {
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
		err := resourceAliCloudCenTransitRouterVpnAttachmentRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTransitRouterVpnAttachmentDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cen_transit_router_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_transit_router_vpn_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteTransitRouterVpnAttachment" {
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
			if *action == "ListTransitRouterVpnAttachments" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCenTransitRouterVpnAttachmentDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Cen TransitRouterVpnAttachment. >>> Resource test cases, automatically generated.
// Case VPN Attachment双隧道 10332
func TestAccAliCloudCenTransitRouterVpnAttachment_basic10332(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterVpnAttachmentMap10332)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterVpnAttachmentBasicDependence10332)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_owner_id":                          "${data.alicloud_account.default.id}",
					"cen_id":                                "${alicloud_cen_transit_router.defaultM8Zo6H.cen_id}",
					"transit_router_attachment_description": "test-vpn-attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id}",
					"vpn_id":                                "${alicloud_vpn_gateway_vpn_attachment.defaultvrPzdh.id}",
					"auto_publish_route_enabled":            "false",
					"charge_type":                           "POSTPAY",
					"transit_router_attachment_name":        "test-vpn-attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_owner_id":                          CHECKSET,
						"cen_id":                                CHECKSET,
						"transit_router_attachment_description": "test-vpn-attachment",
						"transit_router_id":                     CHECKSET,
						"vpn_id":                                CHECKSET,
						"auto_publish_route_enabled":            "false",
						"charge_type":                           "POSTPAY",
						"transit_router_attachment_name":        "test-vpn-attachment",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "test-vpn-attachment2",
					"auto_publish_route_enabled":            "true",
					"transit_router_attachment_name":        "test-vpn-attachment2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "test-vpn-attachment2",
						"auto_publish_route_enabled":            "true",
						"transit_router_attachment_name":        "test-vpn-attachment2",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenTransitRouterVpnAttachmentMap10332 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudCenTransitRouterVpnAttachmentBasicDependence10332(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "defaultbpR5Uk" {
  cen_instance_name = "test-vpn-attachment"
}

resource "alicloud_cen_transit_router" "defaultM8Zo6H" {
  cen_id = alicloud_cen_instance.defaultbpR5Uk.id
}

resource "alicloud_cen_transit_router_cidr" "defaultuUtyCv" {
  cidr              = "192.168.10.0/24"
  transit_router_id = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
}

resource "alicloud_vpn_customer_gateway" "defaultMeoCIz" {
  ip_address            = "42.104.22.${210 + tonumber(substr(var.name, -2, 2)) %% 20}"
  customer_gateway_name = var.name
  depends_on            = ["alicloud_cen_transit_router_cidr.defaultuUtyCv"]
}

data "alicloud_cen_transit_router_service" "default" {
	enable = "On"
}

resource "alicloud_vpn_gateway_vpn_attachment" "defaultvrPzdh" {
  network_type = "public"
  local_subnet = "0.0.0.0/0"
  enable_tunnels_bgp = "false"
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    role = "master"
    enable_dpd = "true"
    enable_nat_traversal = "true"
    tunnel_index = "2"

    tunnel_ike_config {
      remote_id = "2.2.2.2"
      ike_enc_alg = "aes"
      ike_mode = "main"
      ike_version = "ikev1"
      local_id = "1.1.1.1"
      ike_auth_alg = "md5"
      ike_lifetime = "86100"
      ike_pfs = "group2"
      psk = "12345678"
    }
    
      tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg = "aes"
      ipsec_lifetime = "86200"
      ipsec_pfs = "group5"
    }
    
  }
  tunnel_options_specification {
    role = "master"
    enable_nat_traversal = "true"
    tunnel_index = "1"
      tunnel_ike_config {
      local_id = "4.4.4.4"
      remote_id = "5.5.5.5"
      ike_lifetime = "86400"
      ike_pfs = "group5"
      ike_mode = "main"
      ike_version = "ikev2"
      psk = "32333442"
      ike_auth_alg = "md5"
      ike_enc_alg = "aes"
    }
    
      tunnel_ipsec_config {
      ipsec_enc_alg = "aes"
      ipsec_lifetime = "86400"
      ipsec_pfs = "group5"
      ipsec_auth_alg = "sha256"
    }
    
    customer_gateway_id = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    enable_dpd = "true"
  }
  
  remote_subnet = "0.0.0.0/0"
}

`, name)
}

// Case VPN Attachment单隧道 10409
func TestAccAliCloudCenTransitRouterVpnAttachment_basic10409(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterVpnAttachmentMap10409)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterVpnAttachmentBasicDependence10409)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_owner_id":                          "${data.alicloud_account.default.id}",
					"cen_id":                                "${alicloud_cen_transit_router.defaultM8Zo6H.cen_id}",
					"transit_router_attachment_description": "test-vpn-attachment",
					"transit_router_id":                     "${alicloud_cen_transit_router_cidr.defaultuUtyCv.transit_router_id}",
					"vpn_id":                                "${alicloud_vpn_gateway_vpn_attachment.defaultvrPzdh.id}",
					"auto_publish_route_enabled":            "false",
					"charge_type":                           "POSTPAY",
					"transit_router_attachment_name":        "test-vpn-attachment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_owner_id":                          CHECKSET,
						"cen_id":                                CHECKSET,
						"transit_router_attachment_description": "test-vpn-attachment",
						"transit_router_id":                     CHECKSET,
						"vpn_id":                                CHECKSET,
						"auto_publish_route_enabled":            "false",
						"charge_type":                           "POSTPAY",
						"transit_router_attachment_name":        "test-vpn-attachment",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "test-vpn-attachment2",
					"auto_publish_route_enabled":            "true",
					"transit_router_attachment_name":        "test-vpn-attachment2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "test-vpn-attachment2",
						"auto_publish_route_enabled":            "true",
						"transit_router_attachment_name":        "test-vpn-attachment2",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenTransitRouterVpnAttachmentMap10409 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudCenTransitRouterVpnAttachmentBasicDependence10409(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "defaultbpR5Uk" {
  cen_instance_name = "test-vpn-attachment"
}

resource "alicloud_cen_transit_router" "defaultM8Zo6H" {
  cen_id = alicloud_cen_instance.defaultbpR5Uk.id
}

resource "alicloud_cen_transit_router_cidr" "defaultuUtyCv" {
  cidr              = "192.168.10.0/24"
  transit_router_id = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
}

resource "alicloud_vpn_customer_gateway" "defaultMeoCIz" {
  ip_address            = "43.104.22.${230 + tonumber(substr(var.name, -2, 2)) %% 20}"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_gateway_vpn_attachment" "defaultvrPzdh" {
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  enable_tunnels_bgp  = false
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    role                 = "master"
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = 1
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-testvpn10409-1"
      ike_pfs      = "group2"
      remote_id    = "testbob-10409-1"
      local_id     = "testalice-10409-1"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "aes"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    role                 = "master"
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = 2
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-testvpn10409-2"
      ike_pfs      = "group2"
      remote_id    = "testbob-10409-2"
      local_id     = "testalice-10409-2"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "aes"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
}


`, name)
}

// Test Cen TransitRouterVpnAttachment. <<< Resource test cases, automatically generated.
