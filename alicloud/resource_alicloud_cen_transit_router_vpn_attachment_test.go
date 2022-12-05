package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCENTransitRouterVpnAttachment_basic0(t *testing.T) {
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
					"vpn_owner_id": "${data.alicloud_account.default.id}",
					"zone": []map[string]interface{}{
						{
							"zone_id": "${data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0}",
						},
					},
					"transit_router_attachment_name":        "${var.name}",
					"auto_publish_route_enabled":            "false",
					"transit_router_attachment_description": "${var.name}",
					"vpn_id":                                "${alicloud_vpn_gateway_vpn_attachment.default.id}",
					"cen_id":                                "${alicloud_cen_transit_router.default.cen_id}",
					"transit_router_id":                     "${alicloud_cen_transit_router_cidr.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_owner_id":                          CHECKSET,
						"zone.#":                                "1",
						"transit_router_attachment_name":        name,
						"auto_publish_route_enabled":            "false",
						"transit_router_attachment_description": name,
						"transit_router_id":                     CHECKSET,
						"vpn_id":                                CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id"},
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
  		name        = "${var.name}"
  		ip_address  = "42.104.22.212"
  		asn         = "45014"
  		description = "testAccVpnConnectionDesc"
	}

	resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  		customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  		network_type        = "public"
  		local_subnet        = "0.0.0.0/0"
  		remote_subnet       = "0.0.0.0/0"
  		effect_immediately  = false
  		ike_config {
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

  		ipsec_config {
    		ipsec_pfs      = "group5"
    		ipsec_enc_alg  = "des"
			ipsec_auth_alg = "md5"
    		ipsec_lifetime = 86400
  		}
		bgp_config {
			enable       = true
    		local_asn    = 45014
    		tunnel_cidr  = "169.254.11.0/30"
    		local_bgp_ip = "169.254.11.1"
  		}
  		health_check_config {
    		enable   = true
    		sip      = "192.168.1.1"
    		dip      = "10.0.0.1"
    		interval = 10
    		retry    = 10
    		policy   = "revoke_route"
  		}
  		enable_dpd           = true
  		enable_nat_traversal = true
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

func TestAccAlicloudCENTransitRouterVpnAttachment_basic1(t *testing.T) {
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
					"zone": []map[string]interface{}{
						{
							"zone_id": "${data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0}",
						},
					},
					"transit_router_id": "${alicloud_cen_transit_router_cidr.default.transit_router_id}",
					"vpn_id":            "${alicloud_vpn_gateway_vpn_attachment.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterVpnAttachment",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone.#":            "1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id"},
			},
		},
	})
}

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
	err = resourceAlicloudCenTransitRouterVpnAttachmentCreate(dInit, rawClient)
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
		err := resourceAlicloudCenTransitRouterVpnAttachmentCreate(dInit, rawClient)
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
	err = resourceAlicloudCenTransitRouterVpnAttachmentUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenTransitRouterVpnAttachmentUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenTransitRouterVpnAttachmentRead(dExisted, rawClient)
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
	err = resourceAlicloudCenTransitRouterVpnAttachmentDelete(dExisted, rawClient)
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
		err := resourceAlicloudCenTransitRouterVpnAttachmentDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
