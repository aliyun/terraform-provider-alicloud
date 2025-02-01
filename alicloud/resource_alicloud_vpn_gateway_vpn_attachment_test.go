package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
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

func init() {
	resource.AddTestSweepers(
		"alicloud_vpn_gateway_vpn_attachment",
		&resource.Sweeper{
			Name: "alicloud_vpn_gateway_vpn_attachment",
			F:    testSweepVpnGatewayVpnAttachment,
		})
}

func testSweepVpnGatewayVpnAttachment(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.VpnGatewayVpnAttachmentSupportRegions) {
		log.Printf("[INFO] Skipping Vpn Gateway Vpn Attachment unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeVpnConnections"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.VpnConnections.VpnConnection", response)

		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.VpnConnections.VpnConnection", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			name := fmt.Sprint(item["Name"])
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpn Gateway Vpn Attachment: %s", name)
				continue
			}
			action := "DeleteVpnAttachment"
			request := map[string]interface{}{
				"VpnConnectionId": item["VpnConnectionId"],
				"RegionId":        client.RegionId,
			}
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpn Gateway Vpn Attachment (%s): %s", name, err)
			}
			log.Printf("[INFO] Delete Vpn Gateway Vpn Attachment success: %s ", name)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudVPNGatewayVpnAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnAttachmentBasicDependence0)
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
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"effect_immediately":  "false",
					"ike_config": []map[string]string{
						{
							"ike_auth_alg": "md5",
							"ike_enc_alg":  "des",
							"ike_version":  "ikev2",
							"ike_mode":     "main",
							"ike_lifetime": "86400",
							"psk":          "tf-testvpn2",
							"ike_pfs":      "group1",
							"remote_id":    "testbob2",
							"local_id":     "testalice2",
						},
					},
					"ipsec_config": []map[string]string{
						{
							"ipsec_pfs":      "group5",
							"ipsec_enc_alg":  "des",
							"ipsec_auth_alg": "md5",
							"ipsec_lifetime": "86400",
						},
					},
					"bgp_config": []map[string]string{
						{
							"enable":       "true",
							"local_asn":    "45014",
							"tunnel_cidr":  "169.254.11.0/30",
							"local_bgp_ip": "169.254.11.1",
						},
					},
					"health_check_config": []map[string]string{
						{
							"enable":   "true",
							"sip":      "192.168.1.1",
							"dip":      "10.0.0.1",
							"interval": "10",
							"retry":    "10",
							"policy":   "revoke_route",
						},
					},
					"enable_dpd":           "true",
					"enable_nat_traversal": "true",
					"vpn_attachment_name":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customer_gateway_id":            CHECKSET,
						"network_type":                   "public",
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"effect_immediately":             "false",
						"ike_config.#":                   "1",
						"ike_config.0.ike_auth_alg":      "md5",
						"ike_config.0.ike_enc_alg":       "des",
						"ike_config.0.ike_version":       "ikev2",
						"ike_config.0.ike_mode":          "main",
						"ike_config.0.ike_lifetime":      "86400",
						"ike_config.0.psk":               "tf-testvpn2",
						"ike_config.0.ike_pfs":           "group1",
						"ike_config.0.remote_id":         "testbob2",
						"ike_config.0.local_id":          "testalice2",
						"ipsec_config.#":                 "1",
						"ipsec_config.0.ipsec_pfs":       "group5",
						"ipsec_config.0.ipsec_enc_alg":   "des",
						"ipsec_config.0.ipsec_auth_alg":  "md5",
						"ipsec_config.0.ipsec_lifetime":  "86400",
						"bgp_config.#":                   "1",
						"bgp_config.0.enable":            "true",
						"bgp_config.0.local_asn":         "45014",
						"bgp_config.0.local_bgp_ip":      "169.254.11.1",
						"bgp_config.0.tunnel_cidr":       "169.254.11.0/30",
						"health_check_config.#":          "1",
						"health_check_config.0.enable":   "true",
						"health_check_config.0.dip":      "10.0.0.1",
						"health_check_config.0.retry":    "10",
						"health_check_config.0.sip":      "192.168.1.1",
						"health_check_config.0.interval": "10",
						"health_check_config.0.policy":   "revoke_route",
						"enable_dpd":                     "true",
						"enable_nat_traversal":           "true",
						"vpn_attachment_name":            name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultone.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customer_gateway_id": CHECKSET,
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

var AlicloudVPNGatewayVpnAttachmentMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVPNGatewayVpnAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_vpn_customer_gateway" "default" {
	name = "${var.name}"
	ip_address = "42.104.22.210"
	asn = "45014"
	description = "testAccVpnConnectionDesc"
}

resource "alicloud_vpn_customer_gateway" "defaultone" {
  name        = "${var.name}"
  ip_address  = "41.104.22.229"
  asn = "45014"
  description = "${var.name}"
}

`, name)
}

func TestAccAlicloudVPNGatewayVpnAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVpnAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpnattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVpnAttachmentBasicDependence0)
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
					"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customer_gateway_id": CHECKSET,
						"local_subnet":        "0.0.0.0/0",
						"remote_subnet":       "0.0.0.0/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_attachment_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_attachment_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet": "192.168.1.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet": "192.168.1.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_subnet": "192.168.2.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_subnet": "192.168.2.0/24",
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
							"ike_auth_alg": "md5",
							"ike_enc_alg":  "des",
							"ike_version":  "ikev2",
							"ike_mode":     "main",
							"ike_lifetime": "86400",
							"psk":          "tf-testvpn2",
							"ike_pfs":      "group1",
							"remote_id":    "testbob2",
							"local_id":     "testalice2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ike_config.#":              "1",
						"ike_config.0.ike_auth_alg": "md5",
						"ike_config.0.ike_enc_alg":  "des",
						"ike_config.0.ike_version":  "ikev2",
						"ike_config.0.ike_mode":     "main",
						"ike_config.0.ike_lifetime": "86400",
						"ike_config.0.psk":          "tf-testvpn2",
						"ike_config.0.ike_pfs":      "group1",
						"ike_config.0.remote_id":    "testbob2",
						"ike_config.0.local_id":     "testalice2",
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
							"ipsec_lifetime": "86400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipsec_config.#":                "1",
						"ipsec_config.0.ipsec_pfs":      "group5",
						"ipsec_config.0.ipsec_enc_alg":  "des",
						"ipsec_config.0.ipsec_auth_alg": "md5",
						"ipsec_config.0.ipsec_lifetime": "86400",
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
						"bgp_config.#":              "1",
						"bgp_config.0.enable":       "true",
						"bgp_config.0.local_asn":    "45014",
						"bgp_config.0.local_bgp_ip": "169.254.11.1",
						"bgp_config.0.tunnel_cidr":  "169.254.11.0/30",
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
							"policy":   "revoke_route",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#":          "1",
						"health_check_config.0.enable":   "true",
						"health_check_config.0.dip":      "10.0.0.1",
						"health_check_config.0.retry":    "10",
						"health_check_config.0.sip":      "192.168.1.1",
						"health_check_config.0.interval": "10",
						"health_check_config.0.policy":   "revoke_route",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_dpd": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dpd": "false",
					}),
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitAccAlicloudVpnGatewayVpnAttachment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"customer_gateway_id": "CreateVpnGatewayVpnAttachmentValue",
		"network_type":        "CreateVpnGatewayVpnAttachmentValue",
		"local_subnet":        "CreateVpnGatewayVpnAttachmentValue",
		"remote_subnet":       "CreateVpnGatewayVpnAttachmentValue",
		"effect_immediately":  false,
		"ike_config": []map[string]interface{}{
			{
				"ike_auth_alg": "CreateVpnGatewayVpnAttachmentValue",
				"ike_enc_alg":  "CreateVpnGatewayVpnAttachmentValue",
				"ike_version":  "CreateVpnGatewayVpnAttachmentValue",
				"ike_mode":     "CreateVpnGatewayVpnAttachmentValue",
				"ike_lifetime": 86400,
				"psk":          "CreateVpnGatewayVpnAttachmentValue",
				"ike_pfs":      "CreateVpnGatewayVpnAttachmentValue",
				"remote_id":    "CreateVpnGatewayVpnAttachmentValue",
				"local_id":     "CreateVpnGatewayVpnAttachmentValue",
			},
		},
		"ipsec_config": []map[string]interface{}{
			{
				"ipsec_pfs":      "CreateVpnGatewayVpnAttachmentValue",
				"ipsec_enc_alg":  "CreateVpnGatewayVpnAttachmentValue",
				"ipsec_auth_alg": "CreateVpnGatewayVpnAttachmentValue",
				"ipsec_lifetime": 86400,
			},
		},
		"bgp_config": []map[string]interface{}{
			{
				"enable":       true,
				"local_asn":    45014,
				"tunnel_cidr":  "CreateVpnGatewayVpnAttachmentValue",
				"local_bgp_ip": "CreateVpnGatewayVpnAttachmentValue",
			},
		},
		"health_check_config": []map[string]interface{}{
			{
				"enable":   true,
				"sip":      "CreateVpnGatewayVpnAttachmentValue",
				"dip":      "CreateVpnGatewayVpnAttachmentValue",
				"interval": 10,
				"retry":    10,
				"policy":   "CreateVpnGatewayVpnAttachmentValue",
			},
		},
		"enable_dpd":           true,
		"enable_nat_traversal": true,
		"vpn_attachment_name":  "CreateVpnGatewayVpnAttachmentValue",
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
		"Name":              "CreateVpnGatewayVpnAttachmentValue",
		"AttachType":        "CreateVpnGatewayVpnAttachmentValue",
		"EffectImmediately": false,
		"RemoteSubnet":      "CreateVpnGatewayVpnAttachmentValue",
		"NetworkType":       "CreateVpnGatewayVpnAttachmentValue",
		"IpsecConfig": map[string]interface{}{
			"IpsecPfs":      "CreateVpnGatewayVpnAttachmentValue",
			"IpsecEncAlg":   "CreateVpnGatewayVpnAttachmentValue",
			"IpsecAuthAlg":  "CreateVpnGatewayVpnAttachmentValue",
			"IpsecLifetime": 86400,
		},
		"EnableNatTraversal": true,
		"AttachInstanceId":   "",
		"IkeConfig": map[string]interface{}{
			"IkeAuthAlg":  "CreateVpnGatewayVpnAttachmentValue",
			"LocalId":     "CreateVpnGatewayVpnAttachmentValue",
			"IkeEncAlg":   "CreateVpnGatewayVpnAttachmentValue",
			"IkeVersion":  "CreateVpnGatewayVpnAttachmentValue",
			"IkeMode":     "CreateVpnGatewayVpnAttachmentValue",
			"IkeLifetime": 86400,
			"RemoteId":    "CreateVpnGatewayVpnAttachmentValue",
			"Psk":         "CreateVpnGatewayVpnAttachmentValue",
			"IkePfs":      "CreateVpnGatewayVpnAttachmentValue",
		},
		"VpnBgpConfig": map[string]interface{}{
			"EnableBgp":  "true",
			"LocalAsn":   45014,
			"TunnelCidr": "CreateVpnGatewayVpnAttachmentValue",
			"PeerBgpIp":  "CreateVpnGatewayVpnAttachmentValue",
			"PeerAsn":    45014,
			"LocalBgpIp": "CreateVpnGatewayVpnAttachmentValue",
		},
		"LocalSubnet":       "CreateVpnGatewayVpnAttachmentValue",
		"CustomerGatewayId": "CreateVpnGatewayVpnAttachmentValue",
		"CreateTime":        1660027972000,
		"VcoHealthCheck": map[string]interface{}{
			"Policy":   "CreateVpnGatewayVpnAttachmentValue",
			"Enable":   "true",
			"Dip":      "CreateVpnGatewayVpnAttachmentValue",
			"Retry":    10,
			"Sip":      "CreateVpnGatewayVpnAttachmentValue",
			"Interval": 10,
		},
		"VpnGatewayId":    "CreateVpnGatewayVpnAttachmentValue",
		"State":           "init",
		"VpnConnectionId": "VpnGatewayVpnAttachmentId",
		"Spec":            "1000M",
		"EnableDpd":       true,
	}
	CreateMockResponse := map[string]interface{}{
		"VpnConnectionId": "VpnGatewayVpnAttachmentId",
		"Success":         true,
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_vpn_gateway_vpn_attachment", errorCode))
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
	err = resourceAlicloudVpnGatewayVpnAttachmentCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpnAttachment" {
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
		err := resourceAlicloudVpnGatewayVpnAttachmentCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpnGatewayVpnAttachmentUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"local_subnet":       "UpdateVpnGatewayVpnAttachmentValue",
		"remote_subnet":      "UpdateVpnGatewayVpnAttachmentValue",
		"effect_immediately": true,
		"ike_config": []map[string]interface{}{
			{
				"ike_auth_alg": "UpdateVpnGatewayVpnAttachmentValue",
				"ike_enc_alg":  "UpdateVpnGatewayVpnAttachmentValue",
				"ike_version":  "UpdateVpnGatewayVpnAttachmentValue",
				"ike_mode":     "UpdateVpnGatewayVpnAttachmentValue",
				"ike_lifetime": 86400,
				"psk":          "UpdateVpnGatewayVpnAttachmentValue",
				"ike_pfs":      "UpdateVpnGatewayVpnAttachmentValue",
				"remote_id":    "UpdateVpnGatewayVpnAttachmentValue",
				"local_id":     "UpdateVpnGatewayVpnAttachmentValue",
			},
		},
		"ipsec_config": []map[string]interface{}{
			{
				"ipsec_pfs":      "UpdateVpnGatewayVpnAttachmentValue",
				"ipsec_enc_alg":  "UpdateVpnGatewayVpnAttachmentValue",
				"ipsec_auth_alg": "UpdateVpnGatewayVpnAttachmentValue",
				"ipsec_lifetime": 86400,
			},
		},
		"bgp_config": []map[string]interface{}{
			{
				"enable":       true,
				"local_asn":    45014,
				"tunnel_cidr":  "UpdateVpnGatewayVpnAttachmentValue",
				"local_bgp_ip": "UpdateVpnGatewayVpnAttachmentValue",
			},
		},
		"health_check_config": []map[string]interface{}{
			{
				"enable":   true,
				"sip":      "UpdateVpnGatewayVpnAttachmentValue",
				"dip":      "UpdateVpnGatewayVpnAttachmentValue",
				"interval": 10,
				"retry":    10,
				"policy":   "UpdateVpnGatewayVpnAttachmentValue",
			},
		},
		"enable_dpd":           false,
		"enable_nat_traversal": false,
		"vpn_attachment_name":  "UpdateVpnGatewayVpnAttachmentValue",
	}
	diff, err := newInstanceDiff("alicloud_vpn_gateway_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"EffectImmediately": true,
		"RemoteSubnet":      "UpdateVpnGatewayVpnAttachmentValue",
		"Name":              "UpdateVpnGatewayVpnAttachmentValue",
		"IpsecConfig": map[string]interface{}{
			"IpsecPfs":      "UpdateVpnGatewayVpnAttachmentValue",
			"IpsecEncAlg":   "UpdateVpnGatewayVpnAttachmentValue",
			"IpsecAuthAlg":  "UpdateVpnGatewayVpnAttachmentValue",
			"IpsecLifetime": 86400,
		},
		"EnableNatTraversal": false,
		"IkeConfig": map[string]interface{}{
			"IkeAuthAlg":  "UpdateVpnGatewayVpnAttachmentValue",
			"LocalId":     "UpdateVpnGatewayVpnAttachmentValue",
			"IkeEncAlg":   "UpdateVpnGatewayVpnAttachmentValue",
			"IkeVersion":  "UpdateVpnGatewayVpnAttachmentValue",
			"IkeMode":     "UpdateVpnGatewayVpnAttachmentValue",
			"IkeLifetime": 86400,
			"RemoteId":    "UpdateVpnGatewayVpnAttachmentValue",
			"Psk":         "UpdateVpnGatewayVpnAttachmentValue",
			"IkePfs":      "UpdateVpnGatewayVpnAttachmentValue",
		},
		"VpnBgpConfig": map[string]interface{}{
			"EnableBgp":  "true",
			"LocalAsn":   45014,
			"TunnelCidr": "UpdateVpnGatewayVpnAttachmentValue",
			"PeerBgpIp":  "UpdateVpnGatewayVpnAttachmentValue",
			"PeerAsn":    45014,
			"LocalBgpIp": "UpdateVpnGatewayVpnAttachmentValue",
		},
		"LocalSubnet": "UpdateVpnGatewayVpnAttachmentValue",
		"VcoHealthCheck": map[string]interface{}{
			"Policy":   "UpdateVpnGatewayVpnAttachmentValue",
			"Enable":   "true",
			"Dip":      "UpdateVpnGatewayVpnAttachmentValue",
			"Retry":    10,
			"Sip":      "UpdateVpnGatewayVpnAttachmentValue",
			"Interval": 10,
		},
		"EnableDpd": false,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyVpnAttachmentAttribute" {
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
		err := resourceAlicloudVpnGatewayVpnAttachmentUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_vpn_gateway_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeVpnConnection" {
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
		err := resourceAlicloudVpnGatewayVpnAttachmentRead(dExisted, rawClient)
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
	err = resourceAlicloudVpnGatewayVpnAttachmentDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_vpn_gateway_vpn_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpn_gateway_vpn_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpnAttachment" {
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
		err := resourceAlicloudVpnGatewayVpnAttachmentDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
