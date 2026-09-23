package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
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

func TestAccAliCloudVPNGatewayVpnAttachment_basic0(t *testing.T) {
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
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"effect_immediately":  "false",
					"vpn_attachment_name": "${var.name}",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "des",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-testvpn2-1",
									"ike_pfs":      "group1",
									"remote_id":    "testbob2",
									"local_id":     "testalice2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "des",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "des",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-testvpn2-2",
									"ike_pfs":      "group1",
									"remote_id":    "testbob2",
									"local_id":     "testalice2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "des",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":                   "public",
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"effect_immediately":             "false",
						"vpn_attachment_name":            name,
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_attachment_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_attachment_name": name + "_update",
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
	ip_address = "42.${100 + tonumber(substr(var.name, -5, 2)) %% 100}.${100 + tonumber(substr(var.name, -3, 2)) %% 100}.${100 + tonumber(substr(var.name, -1, 1)) %% 10}"
	asn = "45014"
	description = "testAccVpnConnectionDesc"
}

resource "alicloud_vpn_customer_gateway" "defaultone" {
  name        = "${var.name}"
  ip_address  = "41.${100 + tonumber(substr(var.name, -5, 2)) %% 100}.${100 + tonumber(substr(var.name, -3, 2)) %% 100}.${100 + tonumber(substr(var.name, -1, 1)) %% 10}"
  asn = "45014"
  description = "${var.name}"
}

`, name)
}

func TestUnitAlicloudVPNGatewayVpnAttachmentBasicDependenceUsesUniqueIPs(t *testing.T) {
	config := AlicloudVPNGatewayVpnAttachmentBasicDependence0("tf-testaccvpnattachment12345")
	for _, fixedIP := range []string{"42.104.22.210", "41.104.22.229"} {
		if strings.Contains(config, fixedIP) {
			t.Fatalf("generated configuration must not contain fixed customer gateway IP %q", fixedIP)
		}
	}

	matches := regexp.MustCompile(`ip_address\s*=\s*"([^"]+)"`).FindAllStringSubmatch(config, -1)
	if len(matches) != 2 {
		t.Fatalf("expected two customer gateway IP expressions, got %d", len(matches))
	}
	if matches[0][1] == matches[1][1] {
		t.Fatalf("customer gateway IP expressions must differ, got %q", matches[0][1])
	}
	for _, match := range matches {
		if !strings.Contains(match[1], "tonumber(substr(var.name") {
			t.Fatalf("customer gateway IP must derive from var.name, got %q", match[1])
		}
	}
}

func TestAccAliCloudVPNGatewayVpnAttachment_basic1(t *testing.T) {
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
					"local_subnet":  "0.0.0.0/0",
					"remote_subnet": "0.0.0.0/0",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
							"role":                "master",
							"tunnel_index":        "1",
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
							"role":                "master",
							"tunnel_index":        "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"tunnel_options_specification.#": "2",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// lintignore: R001
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
	err = resourceAliCloudVpnGatewayVpnAttachmentCreate(dInit, rawClient)
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
		err := resourceAliCloudVpnGatewayVpnAttachmentCreate(dInit, rawClient)
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
	err = resourceAliCloudVpnGatewayVpnAttachmentUpdate(dExisted, rawClient)
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
		err := resourceAliCloudVpnGatewayVpnAttachmentUpdate(dExisted, rawClient)
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
		err := resourceAliCloudVpnGatewayVpnAttachmentRead(dExisted, rawClient)
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
	err = resourceAliCloudVpnGatewayVpnAttachmentDelete(dExisted, rawClient)
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
		err := resourceAliCloudVpnGatewayVpnAttachmentDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test VpnGateway VpnAttachment. >>> Resource test cases, automatically generated.
// Case 双隧道VpnAttachment测试用例-基础增删改查 10338
func TestAccAliCloudVpnGatewayVpnAttachment_basic10338(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMap10338)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependence10338)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-huhehaote"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "0.0.0.0/0",
					"enable_tunnels_bgp":  "true",
					"vpn_attachment_name": name,
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_index":         "1",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219001",
									"local_bgp_ip": "169.254.10.1",
									"tunnel_cidr":  "169.254.10.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_lifetime": "86100",
									"ike_mode":     "main",
									"ike_pfs":      "group2",
									"ike_version":  "ikev1",
									"local_id":     "1.1.1.1",
									"psk":          "12345678",
									"remote_id":    "2.2.2.2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes",
									"ipsec_lifetime": "86200",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_index":         "2",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219001",
									"local_bgp_ip": "169.254.20.1",
									"tunnel_cidr":  "169.254.20.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_lifetime": "86400",
									"ike_mode":     "main",
									"ike_pfs":      "group5",
									"ike_version":  "ikev2",
									"local_id":     "4.4.4.4",
									"psk":          "32333442",
									"remote_id":    "5.5.5.5",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "sha256",
									"ipsec_enc_alg":  "aes",
									"ipsec_lifetime": "86400",
									"ipsec_pfs":      "group5",
								},
							},
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.cgw1.id}",
						},
					},
					"remote_subnet":     "0.0.0.0/0",
					"network_type":      "public",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"enable_tunnels_bgp":             "true",
						"vpn_attachment_name":            name,
						"tunnel_options_specification.#": "2",
						"remote_subnet":                  "0.0.0.0/0",
						"network_type":                   "public",
						"resource_group_id":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet": "9.9.9.9/32",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_index":         "2",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"psk":          "tunnel2new",
									"ike_auth_alg": "sha384",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "86122",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "2.2.2.2",
									"remote_id":    "3.3.3.3",
								},
							},
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.cgw2.id}",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.42.1",
									"tunnel_cidr":  "169.254.42.0/30",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "sha512",
									"ipsec_enc_alg":  "aes192",
									"ipsec_lifetime": "86111",
									"ipsec_pfs":      "disabled",
								},
							},
						},
						{
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_index":         "1",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219002",
									"local_bgp_ip": "169.254.41.1",
									"tunnel_cidr":  "169.254.41.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha384",
									"ike_enc_alg":  "aes192",
									"ike_lifetime": "86022",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group2",
									"ike_version":  "ikev1",
									"local_id":     "5.5.5.5",
									"psk":          "123456789",
									"remote_id":    "4.4.4.4",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes192",
									"ipsec_lifetime": "86111",
									"ipsec_pfs":      "disabled",
								},
							},
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.cgw2.id}",
						},
					},
					"remote_subnet":     "9.0.0.0/8",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "9.9.9.9/32",
						"tunnel_options_specification.#": "2",
						"remote_subnet":                  "9.0.0.0/8",
						"resource_group_id":              CHECKSET,
					}),
				),
			},
			{
				// Step 2.5: only nested fields (psk, tunnel_bgp_config) change;
				// customer_gateway_id and tunnel_index stay the same as step 2.
				// With the old custom hash (tunnel_index + customer_gateway_id only)
				// this would produce no diff; the list schema must detect it.
				// The declaration order returns to tunnel_index 1, 2 here: every
				// intended change stays plan-visible (an explicit value equal to
				// the positional predecessor's cannot be presented by the legacy
				// list diff), and the state order realigns for the next step.
				Config: testAccConfig(map[string]interface{}{
					"local_subnet": "9.9.9.9/32",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_index":         "1",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219003",
									"local_bgp_ip": "169.254.51.1",
									"tunnel_cidr":  "169.254.51.0/30",
								},
							},
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha384",
									"ike_enc_alg":  "aes192",
									"ike_lifetime": "86022",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group2",
									"ike_version":  "ikev1",
									"local_id":     "5.5.5.5",
									"psk":          "987654321",
									"remote_id":    "4.4.4.4",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "md5",
									"ipsec_enc_alg":  "aes192",
									"ipsec_lifetime": "86111",
									"ipsec_pfs":      "disabled",
								},
							},
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.cgw2.id}",
						},
						{
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_index":         "2",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"psk":          "tunnel2newpsk",
									"ike_auth_alg": "sha384",
									"ike_enc_alg":  "aes256",
									"ike_lifetime": "86122",
									"ike_mode":     "aggressive",
									"ike_pfs":      "group14",
									"ike_version":  "ikev2",
									"local_id":     "2.2.2.2",
									"remote_id":    "3.3.3.3",
								},
							},
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.cgw2.id}",
							"tunnel_bgp_config": []map[string]interface{}{
								{
									"local_asn":    "1219003",
									"local_bgp_ip": "169.254.52.1",
									"tunnel_cidr":  "169.254.52.0/30",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "sha512",
									"ipsec_enc_alg":  "aes192",
									"ipsec_lifetime": "86111",
									"ipsec_pfs":      "disabled",
								},
							},
						},
					},
					"remote_subnet":     "9.0.0.0/8",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "9.9.9.9/32",
						"tunnel_options_specification.#": "2",
						"remote_subnet":                  "9.0.0.0/8",
						"resource_group_id":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":       "0.0.0.0/0",
					"enable_tunnels_bgp": "false",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_index":         "1",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha384",
									"ike_enc_alg":  "aes",
									"ike_lifetime": "55555",
									"ike_mode":     "main",
									"ike_pfs":      "group2",
									"ike_version":  "ikev2",
									"local_id":     "54.54.3.2",
									"psk":          "sadsaa",
									"remote_id":    "54.54.3.23",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "sha512",
									"ipsec_enc_alg":  "3des",
									"ipsec_lifetime": "44232",
									"ipsec_pfs":      "group5",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw2.id}",
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_index":         "2",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha1",
									"ike_enc_alg":  "aes",
									"ike_lifetime": "4432",
									"ike_mode":     "main",
									"ike_pfs":      "group14",
									"ike_version":  "ikev1",
									"local_id":     "54.54.3.2",
									"psk":          "wdascsax",
									"remote_id":    "54.54.3.29",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_auth_alg": "sha512",
									"ipsec_enc_alg":  "des",
									"ipsec_lifetime": "86400",
									"ipsec_pfs":      "group1",
								},
							},
						},
					},
					"remote_subnet":     "0.0.0.0/0",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"enable_tunnels_bgp":             "false",
						"tunnel_options_specification.#": "2",
						"remote_subnet":                  "0.0.0.0/0",
						"resource_group_id":              CHECKSET,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMap10338 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependence10338(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-huhehaote"
}

variable "az2" {
  default = "cn-huhehaote-b"
}

variable "az1" {
  default = "cn-huhehaote-a"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address = "2.2.2.2"
  asn        = "1219001"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address            = "43.43.3.22"
  asn                   = "44331"
  customer_gateway_name = "test_amp"
}


`, name)
}

// Case 单隧道VpnAttachment增删改查 10363
func TestAccAliCloudVpnGatewayVpnAttachment_basic10363(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMap10363)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependence10363)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-huhehaote"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "0.0.0.0/0",
					"vpn_attachment_name": name,
					"effect_immediately":  "true",
					"remote_subnet":       "0.0.0.0/0",
					"network_type":        "private",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "3des",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86100",
									"psk":          "122312421-1",
									"ike_pfs":      "group2",
									"remote_id":    "5.5.5.5",
									"local_id":     "32.32.32.32",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "3des",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86100",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "3des",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86100",
									"psk":          "122312421-2",
									"ike_pfs":      "group2",
									"remote_id":    "5.5.5.5",
									"local_id":     "32.32.32.32",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "3des",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86100",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"vpn_attachment_name":            name,
						"effect_immediately":             "true",
						"remote_subnet":                  "0.0.0.0/0",
						"network_type":                   "private",
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "2.0.0.0/8",
					"vpn_attachment_name": name + "_update",
					"effect_immediately":  "false",
					"remote_subnet":       "3.0.0.0/8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":        "2.0.0.0/8",
						"vpn_attachment_name": name + "_update",
						"effect_immediately":  "false",
						"remote_subnet":       "3.0.0.0/8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "0.0.0.0/0",
					"vpn_attachment_name": name + "_update",
					"effect_immediately":  "true",
					"remote_subnet":       "0.0.0.0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":        "0.0.0.0/0",
						"vpn_attachment_name": name + "_update",
						"effect_immediately":  "true",
						"remote_subnet":       "0.0.0.0/0",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMap10363 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependence10363(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-huhehaote"
}

variable "name2" {
  default = "test_amp2"
}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address = "54.54.54.21"
  asn        = "42311"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address = "3.12.22.33"
  asn        = "44492"
}


`, name)
}

// Case VpnAttachment测试用例-幸确-私网 5629
func TestAccAliCloudVpnGatewayVpnAttachment_basic5629(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMap5629)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependence5629)
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
					"local_subnet":        "0.0.0.0/0",
					"vpn_attachment_name": name,
					"effect_immediately":  "false",
					"remote_subnet":       "172.16.0.0/12",
					"network_type":        "private",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.用户网关1.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"ike_pfs":      "group2",
									"local_id":     "9.0.0.1",
									"psk":          "tf-5629-1",
									"remote_id":    "6.6.6.6",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group2",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.用户网关1.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"ike_pfs":      "group2",
									"local_id":     "9.0.0.1",
									"psk":          "tf-5629-2",
									"remote_id":    "6.6.6.6",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group2",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"vpn_attachment_name":            name,
						"effect_immediately":             "false",
						"remote_subnet":                  "172.16.0.0/12",
						"network_type":                   "private",
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":       "3.0.0.0/24",
					"effect_immediately": "true",
					"remote_subnet":      "2.0.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":       "3.0.0.0/24",
						"effect_immediately": "true",
						"remote_subnet":      "2.0.0.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":  "3.3.0.0/24",
					"remote_subnet": "3.3.2.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":  "3.3.0.0/24",
						"remote_subnet": "3.3.2.0/24",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMap5629 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependence5629(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "eu-central-1"
}

resource "alicloud_vpn_customer_gateway" "用户网关1" {
  ip_address            = "4.4.4.2"
  asn                   = "1219002"
  customer_gateway_name = "用户网关1-VpnAttachment"
  description           = "Xingque-Amp-test-vpn-attachement"
}


`, name)
}

// Case VpnAttachment测试用例-幸确 5358
func TestAccAliCloudVpnGatewayVpnAttachment_basic5358(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMap5358)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependence5358)
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
					"local_subnet":        "0.0.0.0/0",
					"vpn_attachment_name": name,
					"effect_immediately":  "false",
					"remote_subnet":       "172.16.0.0/12",
					"network_type":        "public",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.用户网关1.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"ike_pfs":      "group2",
									"local_id":     "9.0.0.1",
									"psk":          "tf-5358-1",
									"remote_id":    "6.6.6.6",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group2",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.用户网关1.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev1",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"ike_pfs":      "group2",
									"local_id":     "9.0.0.1",
									"psk":          "tf-5358-2",
									"remote_id":    "6.6.6.6",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group2",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":                   "0.0.0.0/0",
						"vpn_attachment_name":            name,
						"effect_immediately":             "false",
						"remote_subnet":                  "172.16.0.0/12",
						"network_type":                   "public",
						"resource_group_id":              CHECKSET,
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "192.168.0.0/24",
					"vpn_attachment_name": name + "_update",
					"remote_subnet":       "172.16.0.0/24",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":        "192.168.0.0/24",
						"vpn_attachment_name": name + "_update",
						"remote_subnet":       "172.16.0.0/24",
						"resource_group_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":        "192.168.0.0/25",
					"vpn_attachment_name": name + "_update",
					"remote_subnet":       "0.0.0.0/1",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":        "192.168.0.0/25",
						"vpn_attachment_name": name + "_update",
						"remote_subnet":       "0.0.0.0/1",
						"resource_group_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":      "192.168.0.0/24",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":      "192.168.0.0/24",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
						"resource_group_id":  CHECKSET,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMap5358 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependence5358(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "eu-central-1"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpn_customer_gateway" "用户网关1" {
  ip_address            = "4.4.4.1"
  asn                   = "1219002"
  customer_gateway_name = "用户网关1-VpnAttachment"
  description           = "Xingque-Amp-test-vpn-attachement"
}

resource "alicloud_vpn_customer_gateway" "用户网关2" {
  description           = "Xingque-Amp-test-vpn-attachement"
  ip_address            = "43.43.43.43"
  asn                   = "1219001"
  customer_gateway_name = "用户网关2-Vpnattachment"
}


`, name)
}

// Case 双隧道VpnAttachment指定主备角色 role
func TestAccAliCloudVpnGatewayVpnAttachment_role(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMapRole)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependenceRole)
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
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"enable_tunnels_bgp":  "false",
					"vpn_attachment_name": name,
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-testvpn-role-1",
									"ike_pfs":      "group2",
									"remote_id":    "role-remote-1",
									"local_id":     "role-local-1",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-testvpn-role-2",
									"ike_pfs":      "group2",
									"remote_id":    "role-remote-2",
									"local_id":     "role-local-2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":                   "public",
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"enable_tunnels_bgp":             "false",
						"vpn_attachment_name":            name,
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw2.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha1",
									"ike_enc_alg":  "aes256",
									"ike_version":  "ikev1",
									"ike_mode":     "aggressive",
									"ike_lifetime": "43200",
									"psk":          "tf-testvpn-role-upd-1",
									"ike_pfs":      "group14",
									"remote_id":    "role-remote-upd-1",
									"local_id":     "role-local-upd-1",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group14",
									"ipsec_enc_alg":  "aes256",
									"ipsec_auth_alg": "sha1",
									"ipsec_lifetime": "43200",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw2.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "false",
							"enable_nat_traversal": "false",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "sha1",
									"ike_enc_alg":  "aes256",
									"ike_version":  "ikev1",
									"ike_mode":     "aggressive",
									"ike_lifetime": "43200",
									"psk":          "tf-testvpn-role-upd-2",
									"ike_pfs":      "group14",
									"remote_id":    "role-remote-upd-2",
									"local_id":     "role-local-upd-2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group14",
									"ipsec_enc_alg":  "aes256",
									"ipsec_auth_alg": "sha1",
									"ipsec_lifetime": "43200",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMapRole = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependenceRole(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address            = "6.6.6.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65001"
  customer_gateway_name = "${var.name}-1"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address            = "6.6.7.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65002"
  customer_gateway_name = "${var.name}-2"
}

`, name)
}

// Case 单隧道VpnAttachment遗留字段覆盖
// 注意：该用例覆盖单隧道遗留字段（customer_gateway_id、ike_config、ipsec_config、
// bgp_config、health_check_config、enable_dpd、enable_nat_traversal）。
// 当前阿里云 API 已在所有区域强制要求双隧道，单隧道创建会返回
// VpnConnection.InvalidCreateTunnelOptions 错误，因此此用例在当前环境中
// 预期创建失败（ExpectError），测试通过表示错误符合预期。
func TestAccAliCloudVpnGatewayVpnAttachment_legacyFields(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMapLegacy)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentLegacyDependence)
	_ = testAccCheck
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":         "0.0.0.0/0",
					"remote_subnet":        "0.0.0.0/0",
					"vpn_attachment_name":  name,
					"effect_immediately":   "false",
					"network_type":         "public",
					"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
					"enable_dpd":           "true",
					"enable_nat_traversal": "true",
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg": "md5",
							"ike_enc_alg":  "aes",
							"ike_version":  "ikev1",
							"ike_mode":     "main",
							"ike_lifetime": "86400",
							"psk":          "tf-legacy-test-psk1",
							"ike_pfs":      "group2",
							"local_id":     "1.2.3.4",
							"remote_id":    "5.6.7.8",
						},
					},
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_auth_alg": "md5",
							"ipsec_enc_alg":  "aes",
							"ipsec_lifetime": "86400",
							"ipsec_pfs":      "group2",
						},
					},
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "true",
							"local_asn":    "65530",
							"tunnel_cidr":  "169.254.30.0/30",
							"local_bgp_ip": "169.254.30.1",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "true",
							"dip":      "10.0.0.2",
							"sip":      "192.168.1.1",
							"retry":    "3",
							"interval": "10",
							"policy":   "revoke_route",
						},
					},
				}),
				Check:       resource.ComposeTestCheckFunc(),
				ExpectError: regexp.MustCompile("VpnConnection.InvalidCreateTunnelOptions"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet":         "0.0.0.0/0",
					"remote_subnet":        "0.0.0.0/0",
					"vpn_attachment_name":  name,
					"effect_immediately":   "false",
					"network_type":         "public",
					"customer_gateway_id":  "${alicloud_vpn_customer_gateway.update.id}",
					"enable_dpd":           "false",
					"enable_nat_traversal": "false",
					"ike_config": []map[string]interface{}{
						{
							"ike_auth_alg": "sha1",
							"ike_enc_alg":  "3des",
							"ike_version":  "ikev2",
							"ike_mode":     "aggressive",
							"ike_lifetime": "43200",
							"psk":          "tf-legacy-test-psk2",
							"ike_pfs":      "group14",
							"local_id":     "10.10.10.1",
							"remote_id":    "10.10.10.2",
						},
					},
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_auth_alg": "sha1",
							"ipsec_enc_alg":  "3des",
							"ipsec_lifetime": "43200",
							"ipsec_pfs":      "group14",
						},
					},
					"bgp_config": []map[string]interface{}{
						{
							"enable":       "false",
							"local_asn":    "65531",
							"tunnel_cidr":  "169.254.31.0/30",
							"local_bgp_ip": "169.254.31.1",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"enable":   "false",
							"dip":      "10.0.0.3",
							"sip":      "192.168.1.2",
							"retry":    "5",
							"interval": "15",
							"policy":   "reserve_route",
						},
					},
				}),
				Check:       resource.ComposeTestCheckFunc(),
				ExpectError: regexp.MustCompile("VpnConnection.InvalidCreateTunnelOptions"),
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMapLegacy = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentLegacyDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpn_customer_gateway" "default" {
  ip_address            = "8.0.0.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65001"
  customer_gateway_name = "${var.name}-legacy-cgw"
}

resource "alicloud_vpn_customer_gateway" "update" {
  ip_address            = "8.0.1.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65002"
  customer_gateway_name = "${var.name}-legacy-cgw-update"
}
`, name)
}

// Case 双隧道VpnAttachment角色切换（master->slave）
func TestAccAliCloudVpnGatewayVpnAttachment_roleSwitch(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMapRoleSwitch)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentRoleSwitchDependence)
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
					"network_type":  "public",
					"local_subnet":  "0.0.0.0/0",
					"remote_subnet": "0.0.0.0/0",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-roleswitch-psk1",
									"ike_pfs":      "group2",
									"remote_id":    "rs-remote-1",
									"local_id":     "rs-local-1",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "master",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-roleswitch-psk2",
									"ike_pfs":      "group2",
									"remote_id":    "rs-remote-2",
									"local_id":     "rs-local-2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":                   "public",
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-roleswitch-psk1",
									"ike_pfs":      "group2",
									"remote_id":    "rs-remote-1",
									"local_id":     "rs-local-1",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.default.id}",
							"role":                 "slave",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-roleswitch-psk2",
									"ike_pfs":      "group2",
									"remote_id":    "rs-remote-2",
									"local_id":     "rs-local-2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpnGatewayVpnAttachmentMapRoleSwitch = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentRoleSwitchDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpn_customer_gateway" "default" {
  ip_address            = "9.0.0.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65001"
  customer_gateway_name = "${var.name}-roleswitch-cgw"
}
`, name)
}

// Test VpnGateway VpnAttachment. <<< Resource test cases, automatically generated.

// Case VpnAttachment tunnel_bandwidth 覆盖 Large/Standard 两值并回读 import
func TestAccAliCloudVpnGatewayVpnAttachment_tunnelBandwidth(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vpn_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayVpnAttachmentMapTunnelBandwidth)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayVpnAttachmentBasicDependenceTunnelBandwidth)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"tunnel_bandwidth":    "Large",
					"vpn_attachment_name": name,
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw1.id}",
							"role":                 "master",
							"tunnel_index":         "1",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-tunnelbw-psk-1",
									"ike_pfs":      "group2",
									"remote_id":    "tbw-remote-1",
									"local_id":     "tbw-local-1",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
						{
							"customer_gateway_id":  "${alicloud_vpn_customer_gateway.cgw2.id}",
							"role":                 "slave",
							"tunnel_index":         "2",
							"enable_dpd":           "true",
							"enable_nat_traversal": "true",
							"tunnel_ike_config": []map[string]interface{}{
								{
									"ike_auth_alg": "md5",
									"ike_enc_alg":  "aes",
									"ike_version":  "ikev2",
									"ike_mode":     "main",
									"ike_lifetime": "86400",
									"psk":          "tf-tunnelbw-psk-2",
									"ike_pfs":      "group2",
									"remote_id":    "tbw-remote-2",
									"local_id":     "tbw-local-2",
								},
							},
							"tunnel_ipsec_config": []map[string]interface{}{
								{
									"ipsec_pfs":      "group5",
									"ipsec_enc_alg":  "aes",
									"ipsec_auth_alg": "md5",
									"ipsec_lifetime": "86400",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":                   "public",
						"local_subnet":                   "0.0.0.0/0",
						"remote_subnet":                  "0.0.0.0/0",
						"tunnel_bandwidth":               "Large",
						"vpn_attachment_name":            name,
						"tunnel_options_specification.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tunnel_bandwidth": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_bandwidth":               "Standard",
						"tunnel_options_specification.#": "2",
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

var AlicloudVpnGatewayVpnAttachmentMapTunnelBandwidth = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayVpnAttachmentBasicDependenceTunnelBandwidth(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address            = "7.7.7.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65001"
  customer_gateway_name = "${var.name}-tbw-1"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address            = "7.7.8.${100 + tonumber(substr(var.name, -2, 2)) %% 50}"
  asn                   = "65002"
  customer_gateway_name = "${var.name}-tbw-2"
}

`, name)
}

func TestUnitVpnGatewayVpnAttachmentRejectsRemovedConfiguration(t *testing.T) {
	for _, removed := range []string{"all_tunnels", "one_tunnel", "tunnel_bgp_config", "tunnel_ike_config", "tunnel_ipsec_config"} {
		t.Run(removed, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentUnitConfig())
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			tunnels := config["tunnel_options_specification"].([]interface{})
			switch removed {
			case "all_tunnels":
				config["tunnel_options_specification"] = []interface{}{}
			case "one_tunnel":
				config["tunnel_options_specification"] = tunnels[:1]
			default:
				tunnels[0].(map[string]interface{})[removed] = []interface{}{}
			}
			if _, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil); err == nil {
				t.Fatal("expected a plan error: an omitted API parameter cannot remove existing tunnel configuration")
			}
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeRegistration(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	if r.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1 for the tunnel set-to-list migration", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 || r.StateUpgraders[0].Version != 0 || r.StateUpgraders[0].Upgrade == nil {
		t.Fatal("expected exactly one state upgrader from schema version 0")
	}
	if r.Schema["tunnel_options_specification"].Type != schema.TypeList {
		t.Fatal("current tunnel schema must be a list")
	}
	if !r.StateUpgraders[0].Type.Equals(resourceAliCloudVpnGatewayVpnAttachmentV0()) {
		t.Fatal("state upgrader must decode the frozen version 0 type")
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeV0Type(t *testing.T) {
	// Version 1 changes only the outer tunnel collection type. Derive the old
	// type from the real resource schema, not from the migration's snapshot.
	legacy := vpnAttachmentLegacyResourceForMigrationTest()
	if !resourceAliCloudVpnGatewayVpnAttachmentV0().Equals(legacy.CoreConfigSchema().ImpliedType()) {
		t.Fatal("version 0 state type differs from the original set-based resource schema")
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeV0JSON(t *testing.T) {
	for _, order := range [][]int{{1, 2}, {2, 1}} {
		t.Run(fmt.Sprint(order), func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			legacy := vpnAttachmentLegacyResourceForMigrationTest()
			raw := vpnAttachmentMigrationTestState(legacy.CoreConfigSchema().ImpliedType(), order)
			encoded, err := json.Marshal(raw)
			if err != nil {
				t.Fatal(err)
			}
			// Check the fixture against the real old schema. JSON decoding keeps
			// both possible serialized set orders for the upgrader to handle.
			if _, err := ctyjson.Unmarshal(encoded, legacy.CoreConfigSchema().ImpliedType()); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(encoded, &raw); err != nil {
				t.Fatal(err)
			}
			upgraded, err := r.StateUpgraders[0].Upgrade(raw, nil)
			if err != nil {
				t.Fatal(err)
			}
			want := vpnAttachmentMigrationTestState(legacy.CoreConfigSchema().ImpliedType(), []int{1, 2})
			if !reflect.DeepEqual(upgraded, want) {
				t.Fatal("migration must only order tunnels by index and retain every field, including id, timeouts and PSK")
			}
			if _, err := schema.JSONMapToStateValue(upgraded, r.CoreConfigSchema()); err != nil {
				t.Fatalf("upgraded state cannot be decoded by the list schema: %s", err)
			}
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeV0Flatmap(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	legacy := vpnAttachmentLegacyResourceForMigrationTest()
	legacyType := legacy.CoreConfigSchema().ImpliedType()
	raw := vpnAttachmentMigrationTestState(legacyType, []int{2, 1})
	legacyValue, err := schema.JSONMapToStateValue(raw, legacy.CoreConfigSchema())
	if err != nil {
		t.Fatal(err)
	}
	// Exercise the SDK's legacy flatmap-to-JSON conversion with the registered
	// frozen type, including SDK-only attributes such as timeouts.
	flatmap := terraform.NewInstanceStateShimmedFromValue(legacyValue, 0)
	decoded, err := schema.StateValueFromInstanceState(flatmap, r.StateUpgraders[0].Type)
	if err != nil {
		t.Fatal(err)
	}
	jsonState, err := schema.StateValueToJSONMap(decoded, r.StateUpgraders[0].Type)
	if err != nil {
		t.Fatal(err)
	}
	upgraded, err := r.StateUpgraders[0].Upgrade(jsonState, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := vpnAttachmentMigrationTestState(legacyType, []int{1, 2})
	if !reflect.DeepEqual(upgraded, want) {
		t.Fatal("flatmap upgrade must retain every version 0 state field")
	}

	// Also exercise automatic version dispatch through Resource.Refresh.
	// Replacing Read ensures this is an entirely offline state migration test.
	oldState, err := legacy.ShimInstanceStateFromValue(legacyValue)
	if err != nil {
		t.Fatal(err)
	}
	readCalled := false
	r.Read = func(d *schema.ResourceData, _ interface{}) error {
		readCalled = true
		tunnels := d.Get("tunnel_options_specification").([]interface{})
		if len(tunnels) != 2 || tunnels[0].(map[string]interface{})["tunnel_index"] != 1 || tunnels[1].(map[string]interface{})["tunnel_index"] != 2 {
			t.Fatal("Read must receive the upgraded, ordered tunnel list")
		}
		return nil
	}
	actual, err := r.Refresh(oldState, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !readCalled || actual.ID != raw["id"] || actual.Meta["schema_version"] != "1" {
		t.Fatal("Refresh must retain the ID and record schema version 1 after migration")
	}
	wantValue, err := schema.JSONMapToStateValue(want, r.CoreConfigSchema())
	if err != nil {
		t.Fatal(err)
	}
	wantState, err := r.ShimInstanceStateFromValue(wantValue)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual.Attributes, wantState.Attributes) {
		t.Fatal("automatic legacy state migration changed fields other than the tunnel collection order")
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeV0Empty(t *testing.T) {
	for name, raw := range map[string]map[string]interface{}{
		"nil state":     nil,
		"missing":       {"id": "vco-test"},
		"null tunnels":  {"id": "vco-test", "tunnel_options_specification": nil},
		"empty tunnels": {"id": "vco-test", "tunnel_options_specification": []interface{}{}},
	} {
		t.Run(name, func(t *testing.T) {
			var want map[string]interface{}
			if raw != nil {
				want = make(map[string]interface{}, len(raw))
				for key, value := range raw {
					want[key] = value
				}
			}
			actual, err := resourceAliCloudVpnGatewayVpnAttachmentStateUpgradeV0(raw, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, want) {
				t.Fatal("migration must preserve absent and empty tunnel collections")
			}
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentStateUpgradeV0Malformed(t *testing.T) {
	for name, tunnels := range map[string]interface{}{
		"not an array":  "invalid",
		"not an object": []interface{}{nil},
		"missing index": []interface{}{map[string]interface{}{"psk": "do-not-log"}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := resourceAliCloudVpnGatewayVpnAttachmentStateUpgradeV0(map[string]interface{}{"tunnel_options_specification": tunnels}, nil)
			if err == nil {
				t.Fatal("malformed tunnel state must fail without panicking")
			}
		})
	}
}

func vpnAttachmentLegacyResourceForMigrationTest() *schema.Resource {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	r.SchemaVersion = 0
	r.StateUpgraders = nil
	r.Schema["tunnel_options_specification"].Type = schema.TypeSet
	return r
}

func vpnAttachmentMigrationTestState(ty cty.Type, order []int) map[string]interface{} {
	state := vpnAttachmentMigrationTestValue(ty, "attachment").(map[string]interface{})
	state["id"] = "vco-migration-test"
	state["timeouts"] = map[string]interface{}{"create": "10m", "delete": "15m", "update": "20m"}
	tunnelType := ty.AttributeType("tunnel_options_specification").ElementType()
	tunnels := make([]interface{}, 0, len(order))
	for _, index := range order {
		tunnel := vpnAttachmentMigrationTestValue(tunnelType, fmt.Sprintf("tunnel-%d", index)).(map[string]interface{})
		tunnel["tunnel_index"] = float64(index)
		tunnels = append(tunnels, tunnel)
	}
	state["tunnel_options_specification"] = tunnels
	return state
}

// Populate every field from the real schema with nonzero, distinct values so
// migration cannot silently drop an unrelated or computed attribute.
func vpnAttachmentMigrationTestValue(ty cty.Type, path string) interface{} {
	switch {
	case ty.Equals(cty.String):
		return path
	case ty.Equals(cty.Number):
		return float64(64512)
	case ty.Equals(cty.Bool):
		return true
	case ty.IsObjectType():
		value := make(map[string]interface{})
		for name, attrType := range ty.AttributeTypes() {
			value[name] = vpnAttachmentMigrationTestValue(attrType, path+"."+name)
		}
		return value
	case ty.IsListType() || ty.IsSetType():
		return []interface{}{vpnAttachmentMigrationTestValue(ty.ElementType(), path+".0")}
	case ty.IsMapType():
		return map[string]interface{}{"test-key": vpnAttachmentMigrationTestValue(ty.ElementType(), path+".test-key")}
	default:
		panic("unsupported migration fixture type")
	}
}

// Exercise the actual SDK diff and Apply path, including the serialized RPC
// request: testing a hash alone cannot prove that a nested update reaches VPC.
func TestUnitVpnGatewayVpnAttachmentNestedDiffAndUpdate(t *testing.T) {
	tests := []struct {
		block, field, wireKey string
		value                 interface{}
	}{
		{"tunnel_ike_config", "psk", "TunnelIkeConfig.Psk", "updated-test-psk"},
		{"tunnel_ike_config", "ike_auth_alg", "TunnelIkeConfig.IkeAuthAlg", "sha256"},
		{"tunnel_ike_config", "ike_enc_alg", "TunnelIkeConfig.IkeEncAlg", "aes256"},
		{"tunnel_ike_config", "ike_version", "TunnelIkeConfig.IkeVersion", "ikev2"},
		{"tunnel_ike_config", "ike_mode", "TunnelIkeConfig.IkeMode", "aggressive"},
		{"tunnel_ike_config", "ike_lifetime", "TunnelIkeConfig.IkeLifetime", 43200},
		{"tunnel_ike_config", "ike_pfs", "TunnelIkeConfig.IkePfs", "group14"},
		{"tunnel_ike_config", "local_id", "TunnelIkeConfig.LocalId", "updated-local"},
		{"tunnel_ike_config", "remote_id", "TunnelIkeConfig.RemoteId", "updated-remote"},
		{"tunnel_bgp_config", "local_asn", "TunnelBgpConfig.LocalAsn", 65001},
		{"tunnel_bgp_config", "local_bgp_ip", "TunnelBgpConfig.LocalBgpIp", "169.254.1.2"},
		{"tunnel_bgp_config", "tunnel_cidr", "TunnelBgpConfig.TunnelCidr", "169.254.2.0/30"},
		{"tunnel_ipsec_config", "ipsec_auth_alg", "TunnelIpsecConfig.IpsecAuthAlg", "sha256"},
		{"tunnel_ipsec_config", "ipsec_enc_alg", "TunnelIpsecConfig.IpsecEncAlg", "aes256"},
		{"tunnel_ipsec_config", "ipsec_lifetime", "TunnelIpsecConfig.IpsecLifetime", 7200},
		{"tunnel_ipsec_config", "ipsec_pfs", "TunnelIpsecConfig.IpsecPfs", "group14"},
		{"", "enable_dpd", "EnableDpd", false},
		{"", "enable_nat_traversal", "EnableNatTraversal", false},
	}
	for _, tt := range tests {
		t.Run(tt.block+"/"+tt.field, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentUnitConfig())
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			tunnel := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			if tt.block != "" {
				tunnel = tunnel[tt.block].([]interface{})[0].(map[string]interface{})
			}
			tunnel[tt.field] = tt.value
			diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil || diff.Empty() {
				t.Fatalf("changing %s.%s produced no resource diff", tt.block, tt.field)
			}
			if diff.RequiresNew() {
				t.Fatal("nested tunnel change unexpectedly requires replacement")
			}
			found := false
			for key, change := range diff.Attributes {
				if strings.HasPrefix(key, "tunnel_options_specification.") && strings.HasSuffix(key, "."+tt.field) && change.Old != change.New {
					found = true
				}
			}
			if !found {
				t.Fatalf("diff omitted changed nested field %s", tt.field)
			}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
			state, err := r.Apply(old.State(), diff, client)
			if err != nil {
				t.Fatal(err)
			}
			calls := requests.modifyCalls()
			if len(calls) != 1 {
				t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
			}
			request := calls[0]
			if got := request.Get("TunnelOptionsSpecification.1.TunnelIndex"); got != "1" {
				t.Fatalf("changed tunnel index = %q, want 1", got)
			}
			if got := request.Get("TunnelOptionsSpecification.1." + tt.wireKey); got != fmt.Sprint(tt.value) {
				t.Errorf("RPC %s = %q, want %q", tt.wireKey, got, fmt.Sprint(tt.value))
			}
			if request.Has("TunnelOptionsSpecification.2.TunnelIndex") {
				t.Error("unchanged second tunnel was included in update")
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentOptionalComputedRead(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentUnitConfig()
	for _, value := range config["tunnel_options_specification"].([]interface{}) {
		tunnel := value.(map[string]interface{})
		for _, key := range []string{"tunnel_ike_config", "tunnel_ipsec_config", "tunnel_bgp_config", "enable_dpd", "enable_nat_traversal"} {
			delete(tunnel, key)
		}
	}
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	d.SetId("vco-unit-test")
	client, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentUnitConfig()))
	if err := r.Read(d, client); err != nil {
		t.Fatal(err)
	}
	if len(vpnAttachmentUnitStateTunnels(d)) != 2 {
		t.Fatal("read did not retain both tunnels")
	}
	vpnAttachmentUnitAssertNoDiff(t, r, d.State(), config, client)
}

func TestUnitVpnGatewayVpnAttachmentPartialComputedRead(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentUnitConfig()
	for _, value := range config["tunnel_options_specification"].([]interface{}) {
		tunnel := value.(map[string]interface{})
		ike := tunnel["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
		tunnel["tunnel_ike_config"] = []interface{}{map[string]interface{}{"psk": ike["psk"]}}
		tunnel["tunnel_bgp_config"] = []interface{}{map[string]interface{}{"local_asn": 65000}}
		tunnel["tunnel_ipsec_config"] = []interface{}{map[string]interface{}{"ipsec_enc_alg": "aes"}}
	}
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	d.SetId("vco-unit-test")
	client, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentUnitConfig()))
	if err := r.Read(d, client); err != nil {
		t.Fatal(err)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, d.State(), config, client)
}

func TestUnitVpnGatewayVpnAttachmentCreateTunnels(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentUnitConfig()
	diff, err := r.Diff(nil, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	state, err := r.Apply(nil, diff, client)
	if err != nil {
		t.Fatal(err)
	}
	requests.mu.Lock()
	calls := append([]url.Values(nil), requests.creates...)
	requests.mu.Unlock()
	if len(calls) != 1 {
		t.Fatalf("CreateVpnAttachment calls = %d, want 1", len(calls))
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Errorf("create unexpectedly followed by %d ModifyVpnAttachmentAttribute calls", got)
	}
	for i := 1; i <= 2; i++ {
		prefix := fmt.Sprintf("TunnelOptionsSpecification.%d.", i)
		if got := calls[0].Get(prefix + "TunnelIndex"); got != fmt.Sprint(i) {
			t.Errorf("created tunnel index = %q, want %d", got, i)
		}
		if got := calls[0].Get(prefix + "TunnelIkeConfig.Psk"); got != fmt.Sprintf("initial-test-psk-%d", i) {
			t.Errorf("created tunnel %d PSK = %q", i, got)
		}
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

func TestUnitVpnGatewayVpnAttachmentReadPSK(t *testing.T) {
	// "123456****" is the documented masked-suffix response form (kept as
	// returned); "literal*test*psk" exercises '*' as a legal PSK character.
	for _, value := range []string{"externally-rotated-test-psk", "literal*test*psk", "123456****"} {
		t.Run(value, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			config := vpnAttachmentUnitConfig()
			config["ike_config"] = []interface{}{map[string]interface{}{"psk": "old-top-level-test-psk"}}
			d := schema.TestResourceDataRaw(t, r.Schema, config)
			d.SetId("vco-unit-test")
			api := vpnAttachmentUnitAPI(vpnAttachmentUnitConfig())
			api["IkeConfig"] = map[string]interface{}{"Psk": value}
			for _, raw := range api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{}) {
				raw.(map[string]interface{})["TunnelIkeConfig"].(map[string]interface{})["Psk"] = value
			}
			client, _ := vpnAttachmentUnitClient(t, api)
			if err := r.Read(d, client); err != nil {
				t.Fatal(err)
			}
			if got := d.Get("ike_config.0.psk"); got != value {
				t.Errorf("top-level PSK = %q, want API value %q", got, value)
			}
			for _, raw := range vpnAttachmentUnitStateTunnels(d) {
				ike := raw.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
				if got := ike["psk"]; got != value {
					t.Errorf("tunnel PSK = %q, want API value %q", got, value)
				}
			}
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentReadStableOrder(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentUnitConfig()
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	d.SetId("vco-unit-test")
	api := vpnAttachmentUnitAPI(config)
	tunnels := api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	client, _ := vpnAttachmentUnitClient(t, api)
	if err := r.Read(d, client); err != nil {
		t.Fatal(err)
	}
	for i, value := range vpnAttachmentUnitStateTunnels(d) {
		if got := value.(map[string]interface{})["tunnel_index"]; got != i+1 {
			t.Errorf("state tunnel at position %d = %v, want index %d", i, got, i+1)
		}
	}
	vpnAttachmentUnitAssertNoDiff(t, r, d.State(), config, client)
}

func TestUnitVpnGatewayVpnAttachmentReorderWithoutModify(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentUnitConfig())
	old.SetId("vco-unit-test")
	config := vpnAttachmentUnitConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("configuration reorder must update positional state")
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentUnitConfig()))
	state, err := r.Apply(old.State(), diff, client)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Errorf("reordering unchanged tunnels invoked ModifyVpnAttachmentAttribute %d times", got)
	}
	if state.Attributes["tunnel_options_specification.0.tunnel_index"] != "2" || state.Attributes["tunnel_options_specification.1.tunnel_index"] != "1" {
		t.Error("Apply/Read did not preserve reordered configuration positions")
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

func TestUnitVpnGatewayVpnAttachmentEnableBgpOnly(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	oldConfig := vpnAttachmentUnitConfig()
	oldConfig["enable_tunnels_bgp"] = false
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vco-unit-test")
	config := vpnAttachmentUnitConfig()
	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	for key := range diff.Attributes {
		if strings.HasPrefix(key, "tunnel_options_specification.") {
			t.Fatalf("fixture should change only enable_tunnels_bgp, got %s", key)
		}
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	state, err := r.Apply(old.State(), diff, client)
	if err != nil {
		t.Fatal(err)
	}
	calls := requests.modifyCalls()
	if len(calls) != 1 {
		t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
	}
	request := calls[0]
	if got := request.Get("EnableTunnelsBgp"); got != "true" {
		t.Errorf("EnableTunnelsBgp = %q, want true", got)
	}
	for i := 1; i <= 2; i++ {
		prefix := fmt.Sprintf("TunnelOptionsSpecification.%d.", i)
		if got := request.Get(prefix + "TunnelIndex"); got != fmt.Sprint(i) {
			t.Errorf("tunnel index = %q, want %d", got, i)
		}
		if got := request.Get(prefix + "TunnelBgpConfig.LocalAsn"); got != "65000" {
			t.Errorf("tunnel %d BGP ASN = %q, want 65000", i, got)
		}
		for key := range request {
			if strings.HasPrefix(key, prefix+"TunnelIkeConfig.") || strings.HasPrefix(key, prefix+"TunnelIpsecConfig.") {
				t.Errorf("BGP-only change resent unchanged field %s", key)
			}
		}
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

func TestUnitVpnGatewayVpnAttachmentUpdateWaitsForCompletion(t *testing.T) {
	for _, completedState := range []string{"attached", "init", "active"} {
		t.Run(completedState, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			oldConfig := vpnAttachmentUnitConfig()
			old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			tunnel := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			tunnel["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = "updated-after-completion"
			tunnel["tunnel_bgp_config"].([]interface{})[0].(map[string]interface{})["local_asn"] = 65001
			diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			// ModifyVpnAttachmentAttribute is asynchronous. Its completion
			// indicator is top-level State, not IKE/IPsec negotiation Status.
			updating := vpnAttachmentUnitAPI(oldConfig)
			updating["State"] = "updating"
			completed := vpnAttachmentUnitAPI(config)
			completed["State"] = completedState
			client, requests := vpnAttachmentUnitClientResponses(t, []map[string]interface{}{updating, updating, completed})
			state, err := r.Apply(old.State(), diff, client)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(requests.modifyCalls()); got != 1 {
				t.Errorf("ModifyVpnAttachmentAttribute calls = %d, want 1", got)
			}
			requests.mu.Lock()
			describes := requests.describes
			requests.mu.Unlock()
			if describes < 3 {
				t.Errorf("Apply returned after %d Describe call(s), before configuration completed", describes)
			}
			if got := state.Attributes["tunnel_options_specification.0.tunnel_ike_config.0.psk"]; got != "updated-after-completion" {
				t.Errorf("Apply retained PSK from before asynchronous update: %q", got)
			}
			if got := state.Attributes["tunnel_options_specification.0.tunnel_bgp_config.0.local_asn"]; got != "65001" {
				t.Errorf("Apply retained BGP ASN from before asynchronous update: %q", got)
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentUpdateRejectsFailureStates(t *testing.T) {
	for _, failedState := range []string{"financialLocked", "deleted"} {
		t.Run(failedState, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			oldConfig := vpnAttachmentUnitConfig()
			old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			tunnel := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			tunnel["enable_dpd"] = false
			diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			failed := vpnAttachmentUnitAPI(oldConfig)
			failed["State"] = failedState
			client, requests := vpnAttachmentUnitClient(t, failed)
			if _, err := r.Apply(old.State(), diff, client); err == nil {
				t.Fatalf("Apply reported success for failed state %s", failedState)
			} else if !strings.Contains(err.Error(), failedState) {
				t.Errorf("error does not identify failed state %s: %s", failedState, err)
			}
			requests.mu.Lock()
			describes := requests.describes
			requests.mu.Unlock()
			if describes != 1 {
				t.Errorf("expected immediate failure after one Describe, got %d calls", describes)
			}
		})
	}
}

// The ACC basic10338 counterexample: a create fills the Computed-only leaves
// (tunnel_bgp_config bgp_status, peer_bgp_ip and peer_asn) with per-tunnel
// values, and the next apply reorders tunnels while writing changed values.
// The positional list diff presents those leaves as values inherited from the
// tunnel previously at the same position, differing from each tunnel's own;
// they are not user-writable, so they must not trigger the reorder ambiguity
// rejection, and the restore keeps each tunnel's own values.
func TestUnitVpnGatewayVpnAttachmentReorderWithComputedLeafState(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentDistinctConfig()
	diff, err := r.Diff(nil, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	createClient, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPIDistinctComputed(config))
	old, err := r.Apply(nil, diff, createClient)
	if err != nil {
		t.Fatal(err)
	}
	for position, want := range map[int]string{0: "169.254.1.2", 1: "169.254.2.2"} {
		if got := old.Attributes[fmt.Sprintf("tunnel_options_specification.%d.tunnel_bgp_config.0.peer_bgp_ip", position)]; got != want {
			t.Fatalf("fixture precondition: computed peer_bgp_ip at position %d = %q, want %q", position, got, want)
		}
	}

	updated := func() map[string]interface{} {
		config := vpnAttachmentDistinctConfig()
		config["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = "updated-tunnel-2-test-psk"
		return config
	}

	t.Run("reorder_and_modify_applies_once", func(t *testing.T) {
		reordered := updated()
		tunnels := reordered["tunnel_options_specification"].([]interface{})
		tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
		diff, err := r.Diff(old, terraform.NewResourceConfigRaw(reordered), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPIDistinctComputed(updated()))
		state, err := r.Apply(old, diff, client)
		if err != nil {
			t.Fatalf("reorder with explicit values must apply, not be rejected over computed-only leaves: %s", err)
		}
		calls := requests.modifyCalls()
		if len(calls) != 1 {
			t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
		}
		request := calls[0]
		prefix := "TunnelOptionsSpecification.1."
		if got := request.Get(prefix + "TunnelIndex"); got != "2" {
			t.Errorf("updated tunnel index = %q, want 2", got)
		}
		if got := request.Get(prefix + "TunnelIkeConfig.Psk"); got != "updated-tunnel-2-test-psk" {
			t.Errorf("TunnelIkeConfig.Psk = %q, want updated-tunnel-2-test-psk", got)
		}
		if got := request.Get(prefix + "TunnelIkeConfig.LocalId"); got != "tunnel-2-local" {
			t.Errorf("TunnelIkeConfig.LocalId = %q, want the tunnel's own tunnel-2-local", got)
		}
		for key := range request {
			if strings.HasPrefix(key, prefix+"TunnelBgpConfig.") {
				t.Errorf("computed peer values or unchanged BGP settings resent in update: %s", key)
			}
			if strings.HasPrefix(key, "TunnelOptionsSpecification.2.") {
				t.Errorf("unchanged tunnel included in update: %s", key)
			}
		}
		if got := state.Attributes["tunnel_options_specification.0.tunnel_bgp_config.0.peer_bgp_ip"]; got != "169.254.2.2" {
			t.Errorf("moved tunnel's own computed peer_bgp_ip = %q, want 169.254.2.2", got)
		}
		if got := state.Attributes["tunnel_options_specification.1.tunnel_bgp_config.0.peer_bgp_ip"]; got != "169.254.1.2" {
			t.Errorf("tunnel 1's own computed peer_bgp_ip = %q, want 169.254.1.2", got)
		}
		if got := state.Attributes["tunnel_options_specification.0.tunnel_ike_config.0.psk"]; got != "updated-tunnel-2-test-psk" {
			t.Errorf("moved tunnel psk = %q, want updated-tunnel-2-test-psk", got)
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, reordered, client)
	})

	t.Run("pure_reorder_stays_free_of_modify", func(t *testing.T) {
		reordered := vpnAttachmentDistinctConfig()
		tunnels := reordered["tunnel_options_specification"].([]interface{})
		tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
		diff, err := r.Diff(old, terraform.NewResourceConfigRaw(reordered), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPIDistinctComputed(config))
		state, err := r.Apply(old, diff, client)
		if err != nil {
			t.Fatal(err)
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Fatalf("pure reorder with computed leaf state invoked ModifyVpnAttachmentAttribute %d times", got)
		}
		if got := state.Attributes["tunnel_options_specification.0.tunnel_bgp_config.0.peer_bgp_ip"]; got != "169.254.2.2" {
			t.Errorf("moved tunnel's own computed peer_bgp_ip = %q, want 169.254.2.2", got)
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, reordered, client)
	})

	t.Run("writable_ambiguity_still_rejected", func(t *testing.T) {
		partial := updated()
		tunnels := partial["tunnel_options_specification"].([]interface{})
		tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
		tunnels[0].(map[string]interface{})["tunnel_ike_config"] = []interface{}{map[string]interface{}{"psk": "updated-tunnel-2-test-psk"}}
		diff, err := r.Diff(old, terraform.NewResourceConfigRaw(partial), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPIDistinctComputed(updated()))
		_, err = r.Apply(old, diff, client)
		if err == nil {
			t.Fatal("a partially written block inside a modified scope must still be rejected")
		}
		if !strings.Contains(err.Error(), "tunnel_options_specification.0.tunnel_ike_config.0.local_id") || !strings.Contains(err.Error(), "Split the operation into two applies") {
			t.Errorf("rejection does not name the writable inherited leaf: %s", err)
		}
		if strings.Contains(err.Error(), "peer_bgp_ip") {
			t.Errorf("computed-only leaf wrongly reported as a conflict: %s", err)
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Fatalf("rejected update still invoked ModifyVpnAttachmentAttribute %d times", got)
		}
	})
}

func TestUnitVpnGatewayVpnAttachmentReorderWithOmittedBlocks(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	for _, value := range tunnels {
		tunnel := value.(map[string]interface{})
		for _, key := range []string{"tunnel_ike_config", "tunnel_bgp_config", "tunnel_ipsec_config", "enable_dpd", "enable_nat_traversal"} {
			delete(tunnel, key)
		}
	}
	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("configuration reorder must update positional state")
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentDistinctConfig()))
	state, err := r.Apply(old.State(), diff, client)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Fatalf("reordering tunnels with omitted blocks invoked ModifyVpnAttachmentAttribute %d times with cross-associated values", got)
	}
	if state.Attributes["tunnel_options_specification.0.tunnel_index"] != "2" || state.Attributes["tunnel_options_specification.1.tunnel_index"] != "1" {
		t.Error("Apply/Read did not preserve reordered configuration positions")
	}
	if got := state.Attributes["tunnel_options_specification.0.tunnel_ike_config.0.psk"]; got != "tunnel-2-test-psk" {
		t.Errorf("reordered tunnel kept PSK %q, want its own tunnel-2-test-psk", got)
	}
	if got := state.Attributes["tunnel_options_specification.0.enable_dpd"]; got != "false" {
		t.Errorf("reordered tunnel kept enable_dpd %q, want its own false", got)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

func TestUnitVpnGatewayVpnAttachmentReorderWithExplicitChange(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
	old.SetId("vco-unit-test")
	updated := vpnAttachmentDistinctConfig()
	updated["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = "updated-tunnel-2-test-psk"
	swapped := func() map[string]interface{} {
		config := vpnAttachmentDistinctConfig()
		tunnels := config["tunnel_options_specification"].([]interface{})
		tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
		for _, value := range tunnels {
			tunnel := value.(map[string]interface{})
			for _, key := range []string{"tunnel_ike_config", "tunnel_bgp_config", "tunnel_ipsec_config", "enable_dpd", "enable_nat_traversal"} {
				delete(tunnel, key)
			}
		}
		return config
	}

	// Writing only part of a block inside an actively modified scope leaves the
	// unwritten leaves (local_id, remote_id) presented as no change, so their
	// planned values would be inherited from the tunnel previously at the same
	// position. Apply must reject the update instead of silently deferring the
	// explicit values to a second apply.
	t.Run("partial_block_write_rejected", func(t *testing.T) {
		partial := swapped()
		partial["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"] = []interface{}{map[string]interface{}{"psk": "updated-tunnel-2-test-psk"}}
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(partial), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(updated))
		_, err = r.Apply(old.State(), diff, client)
		if err == nil {
			t.Fatal("a partially written block inside a modified scope must be rejected, not silently deferred")
		}
		for _, want := range []string{"cannot be applied reliably", "tunnel_options_specification.0.tunnel_ike_config.0.local_id", "Split the operation into two applies"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q", err, want)
			}
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Fatalf("rejected update still invoked ModifyVpnAttachmentAttribute %d times", got)
		}
	})

	// Writing the complete block makes every differing leaf plan-visible, so a
	// single apply sends the tunnel's own values plus the rotation and nothing
	// else; the unchanged tunnel and the omitted blocks stay out of the update.
	t.Run("complete_block_write_applies_once", func(t *testing.T) {
		complete := swapped()
		complete["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"] = []interface{}{map[string]interface{}{
			"psk": "updated-tunnel-2-test-psk", "ike_auth_alg": "sha1", "ike_enc_alg": "aes", "ike_version": "ikev1",
			"ike_mode": "main", "ike_lifetime": 86400, "ike_pfs": "group2", "local_id": "tunnel-2-local", "remote_id": "tunnel-2-remote",
		}}
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(complete), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(updated))
		state, err := r.Apply(old.State(), diff, client)
		if err != nil {
			t.Fatal(err)
		}
		calls := requests.modifyCalls()
		if len(calls) != 1 {
			t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
		}
		request := calls[0]
		prefix := "TunnelOptionsSpecification.1."
		if got := request.Get(prefix + "TunnelIndex"); got != "2" {
			t.Errorf("updated tunnel index = %q, want 2", got)
		}
		if got := request.Get(prefix + "EnableDpd"); got != "false" {
			t.Errorf("updated tunnel sent EnableDpd = %q, want its own false", got)
		}
		if got := request.Get(prefix + "TunnelIkeConfig.Psk"); got != "updated-tunnel-2-test-psk" {
			t.Errorf("TunnelIkeConfig.Psk = %q, want updated-tunnel-2-test-psk", got)
		}
		if got := request.Get(prefix + "TunnelIkeConfig.LocalId"); got != "tunnel-2-local" {
			t.Errorf("TunnelIkeConfig.LocalId = %q, want the tunnel's own tunnel-2-local", got)
		}
		for key := range request {
			if strings.HasPrefix(key, prefix+"TunnelBgpConfig.") || strings.HasPrefix(key, prefix+"TunnelIpsecConfig.") {
				t.Errorf("PSK-only change resent unchanged block field %s", key)
			}
			if strings.HasPrefix(key, "TunnelOptionsSpecification.2.") {
				t.Errorf("unchanged tunnel included in update: %s", key)
			}
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, complete, client)
	})
}

func TestUnitVpnGatewayVpnAttachmentReorderOmittedAsymmetricBlocks(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	oldConfig := vpnAttachmentDistinctConfig()
	delete(oldConfig["tunnel_options_specification"].([]interface{})[0].(map[string]interface{}), "tunnel_bgp_config")
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vco-unit-test")
	api := vpnAttachmentUnitAPI(vpnAttachmentDistinctConfig())
	delete(api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{})[0].(map[string]interface{}), "TunnelBgpConfig")
	config := vpnAttachmentDistinctConfig()
	delete(config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{}), "tunnel_bgp_config")
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	for _, value := range tunnels {
		tunnel := value.(map[string]interface{})
		for _, key := range []string{"tunnel_ike_config", "tunnel_bgp_config", "tunnel_ipsec_config", "enable_dpd", "enable_nat_traversal"} {
			delete(tunnel, key)
		}
	}
	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("inherited absence from a reorder is not an explicit block clear: %s", err)
	}
	client, requests := vpnAttachmentUnitClient(t, api)
	state, err := r.Apply(old.State(), diff, client)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Fatalf("reordering asymmetric blocks invoked ModifyVpnAttachmentAttribute %d times", got)
	}
	if got := state.Attributes["tunnel_options_specification.0.tunnel_bgp_config.0.local_asn"]; got != "65002" {
		t.Errorf("tunnel 2 lost its own BGP configuration after reorder: local_asn = %q", got)
	}
	if got := state.Attributes["tunnel_options_specification.1.tunnel_bgp_config.#"]; got != "0" {
		t.Errorf("tunnel 1 gained a BGP block from the other tunnel: count = %q", got)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

// A reorder combined with a symmetric value swap asks each tunnel to take the
// other's current value. The positional list diff presents neither swap
// (planned equals the positional predecessor at both positions), so a single
// apply cannot tell the explicit swap from an inherited value. Apply must
// reject the combined change with guidance instead of silently succeeding and
// deferring the swap to a second apply; splitting the reorder from the value
// changes makes both steps converge in one apply each.
func TestUnitVpnGatewayVpnAttachmentReorderSymmetricSwap(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	oldConfig := vpnAttachmentDistinctConfig()
	oldConfig["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["ike_version"] = "ikev2"
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	ike2 := tunnels[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	ike2["ike_version"] = "ikev1"
	ike2["psk"] = "rotated-tunnel-2-psk"
	ike1 := tunnels[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	ike1["ike_version"] = "ikev2"
	ike1["psk"] = "rotated-tunnel-1-psk"

	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentDistinctConfig())
	if _, err := r.Apply(old.State(), diff, client); err == nil {
		t.Fatal("a reorder with an unpresentable explicit swap must be rejected, not silently deferred to a second apply")
	} else {
		for _, want := range []string{"cannot be applied reliably", "tunnel_options_specification.0.tunnel_ike_config.0.ike_version", "Split the operation into two applies"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q", err, want)
			}
		}
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Fatalf("rejected update still invoked ModifyVpnAttachmentAttribute %d times", got)
	}

	// Split step one: the reorder alone, with the configuration mirroring the
	// current state values. Positions move, every tunnel keeps its own values,
	// and no cloud modification is sent.
	reordered := vpnAttachmentDistinctConfig()
	reordered["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["ike_version"] = "ikev2"
	reorderTunnels := reordered["tunnel_options_specification"].([]interface{})
	reorderTunnels[0], reorderTunnels[1] = reorderTunnels[1], reorderTunnels[0]
	diff2, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(reordered), nil)
	if err != nil {
		t.Fatal(err)
	}
	client2, requests2 := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(reordered))
	state2, err := r.Apply(old.State(), diff2, client2)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(requests2.modifyCalls()); got != 0 {
		t.Fatalf("reorder-only apply invoked ModifyVpnAttachmentAttribute %d times", got)
	}
	if got := state2.Attributes["tunnel_options_specification.0.tunnel_ike_config.0.ike_version"]; got != "ikev2" {
		t.Errorf("reordered tunnel 2 ike_version = %q, want its own ikev2", got)
	}
	if got := state2.Attributes["tunnel_options_specification.1.tunnel_ike_config.0.ike_version"]; got != "ikev1" {
		t.Errorf("reordered tunnel 1 ike_version = %q, want its own ikev1", got)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state2, reordered, client2)

	// Split step two: with state order matching the configuration, the value
	// changes present as ordinary diffs and one apply sends both rotations.
	changed := vpnAttachmentDistinctConfig()
	changedTunnels := changed["tunnel_options_specification"].([]interface{})
	changedTunnels[0], changedTunnels[1] = changedTunnels[1], changedTunnels[0]
	changedIke2 := changedTunnels[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	changedIke2["ike_version"] = "ikev1"
	changedIke2["psk"] = "rotated-tunnel-2-psk"
	changedIke1 := changedTunnels[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	changedIke1["ike_version"] = "ikev2"
	changedIke1["psk"] = "rotated-tunnel-1-psk"
	diff3, err := r.Diff(state2, terraform.NewResourceConfigRaw(changed), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff3 == nil || diff3.Empty() {
		t.Fatal("value changes after the reorder produced no plan")
	}
	applied := vpnAttachmentDistinctConfig()
	appliedTunnels := applied["tunnel_options_specification"].([]interface{})
	appliedTunnels[0], appliedTunnels[1] = appliedTunnels[1], appliedTunnels[0]
	appliedIke2 := appliedTunnels[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	appliedIke2["ike_version"] = "ikev1"
	appliedIke2["psk"] = "rotated-tunnel-2-psk"
	appliedIke1 := appliedTunnels[1].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
	appliedIke1["ike_version"] = "ikev2"
	appliedIke1["psk"] = "rotated-tunnel-1-psk"
	client3, requests3 := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(applied))
	state3, err := r.Apply(state2, diff3, client3)
	if err != nil {
		t.Fatal(err)
	}
	calls := requests3.modifyCalls()
	if len(calls) != 1 {
		t.Fatalf("value-change apply invoked ModifyVpnAttachmentAttribute %d times, want 1", len(calls))
	}
	versions := map[string]string{}
	psks := map[string]string{}
	for i := 1; i <= 2; i++ {
		prefix := fmt.Sprintf("TunnelOptionsSpecification.%d.", i)
		versions[calls[0].Get(prefix+"TunnelIndex")] = calls[0].Get(prefix + "TunnelIkeConfig.IkeVersion")
		psks[calls[0].Get(prefix+"TunnelIndex")] = calls[0].Get(prefix + "TunnelIkeConfig.Psk")
	}
	if versions["1"] != "ikev2" || versions["2"] != "ikev1" {
		t.Errorf("value-change apply sent IkeVersion %v, want tunnel 1 ikev2 and tunnel 2 ikev1", versions)
	}
	if psks["1"] != "rotated-tunnel-1-psk" || psks["2"] != "rotated-tunnel-2-psk" {
		t.Errorf("value-change apply sent Psk %v, want the per-tunnel rotations", psks)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state3, changed, client3)
}

// An explicit write whose value coincides with the positional predecessor's is
// presented as no change even at the top level of a tunnel. Inside an actively
// modified field group that write is ambiguous and must be rejected with the
// affected attribute named, not silently inherited.
func TestUnitVpnGatewayVpnAttachmentReorderAmbiguousTopLevel(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	// Tunnel 2 moves into position 0 (the predecessor is tunnel 1, whose
	// enable_nat_traversal is true): enable_dpd=false is a visible change,
	// while the omitted enable_nat_traversal presents as no change against the
	// predecessor's true and would inherit it, differing from tunnel 2's own
	// false.
	tunnels[0].(map[string]interface{})["enable_dpd"] = false
	delete(tunnels[0].(map[string]interface{}), "enable_nat_traversal")
	diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentDistinctConfig())
	if _, err := r.Apply(old.State(), diff, client); err == nil {
		t.Fatal("an ambiguous top-level write inside a modified field group must be rejected, not silently inherited")
	} else {
		for _, want := range []string{"cannot be applied reliably", "tunnel_options_specification.0.enable_nat_traversal", "Split the operation into two applies"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q", err, want)
			}
		}
	}
	if got := len(requests.modifyCalls()); got != 0 {
		t.Fatalf("rejected update still invoked ModifyVpnAttachmentAttribute %d times", got)
	}
}

// A reordered tunnel that never had a tunnel_bgp_config must not inherit the
// positional predecessor's BGP leaves. A partial new block is an incomplete
// configuration and must be rejected explicitly; a complete new block applies
// with only the user's own values.
func TestUnitVpnGatewayVpnAttachmentReorderNewBlockNoInheritance(t *testing.T) {
	oldWithoutBgp := func() map[string]interface{} {
		config := vpnAttachmentDistinctConfig()
		delete(config["tunnel_options_specification"].([]interface{})[1].(map[string]interface{}), "tunnel_bgp_config")
		return config
	}
	reordered := func(bgp map[string]interface{}) map[string]interface{} {
		config := vpnAttachmentDistinctConfig()
		tunnels := config["tunnel_options_specification"].([]interface{})
		tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
		for _, value := range tunnels {
			tunnel := value.(map[string]interface{})
			for _, key := range []string{"tunnel_ike_config", "tunnel_bgp_config", "tunnel_ipsec_config", "enable_dpd", "enable_nat_traversal"} {
				delete(tunnel, key)
			}
		}
		if bgp != nil {
			tunnels[0].(map[string]interface{})["tunnel_bgp_config"] = []interface{}{bgp}
		}
		return config
	}

	t.Run("partial_block_rejected", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		old := schema.TestResourceDataRaw(t, r.Schema, oldWithoutBgp())
		old.SetId("vco-unit-test")
		config := reordered(map[string]interface{}{"local_asn": 65100})
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentDistinctConfig()))
		if _, err := r.Apply(old.State(), diff, client); err == nil {
			t.Fatal("applying a partial new tunnel_bgp_config must be rejected, not completed with inherited values")
		} else if !strings.Contains(err.Error(), "tunnel_cidr") {
			t.Errorf("error does not name the missing tunnel_cidr: %s", err)
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Errorf("incomplete BGP configuration reached ModifyVpnAttachmentAttribute %d times", got)
		}
	})

	t.Run("complete_block_applies_own_values", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		old := schema.TestResourceDataRaw(t, r.Schema, oldWithoutBgp())
		old.SetId("vco-unit-test")
		config := reordered(map[string]interface{}{
			"local_asn": 65100, "local_bgp_ip": "169.254.9.1", "tunnel_cidr": "169.254.9.0/30",
		})
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		applied := vpnAttachmentDistinctConfig()
		applied["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["tunnel_bgp_config"] = []interface{}{map[string]interface{}{
			"local_asn": 65100, "local_bgp_ip": "169.254.9.1", "tunnel_cidr": "169.254.9.0/30",
		}}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(applied))
		state, err := r.Apply(old.State(), diff, client)
		if err != nil {
			t.Fatal(err)
		}
		calls := requests.modifyCalls()
		if len(calls) != 1 {
			t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
		}
		request := calls[0]
		prefix := "TunnelOptionsSpecification.1."
		if got := request.Get(prefix + "TunnelIndex"); got != "2" {
			t.Errorf("updated tunnel index = %q, want 2", got)
		}
		for key, want := range map[string]string{"LocalAsn": "65100", "LocalBgpIp": "169.254.9.1", "TunnelCidr": "169.254.9.0/30"} {
			if got := request.Get(prefix + "TunnelBgpConfig." + key); got != want {
				t.Errorf("TunnelBgpConfig.%s = %q, want the user's own %q, not an inherited value", key, got, want)
			}
		}
		for key := range request {
			if strings.HasPrefix(key, prefix+"TunnelIkeConfig.") || strings.HasPrefix(key, prefix+"TunnelIpsecConfig.") {
				t.Errorf("new BGP block resent unchanged field %s", key)
			}
			if strings.HasPrefix(key, "TunnelOptionsSpecification.2.") {
				t.Errorf("unchanged tunnel included in update: %s", key)
			}
		}
		if got := state.Attributes["tunnel_options_specification.0.tunnel_bgp_config.0.tunnel_cidr"]; got != "169.254.9.0/30" {
			t.Errorf("state tunnel_cidr = %q, want the user's own 169.254.9.0/30", got)
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
	})
}

// Updating IKE parameters while the tunnel's pre-shared key is unknown (empty
// in state, e.g. after import) must be blocked: submitting the IKE block
// without Psk makes the API generate a random 16-character key.
func TestUnitVpnGatewayVpnAttachmentUpdateTunnelIkeWithoutPsk(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		oldConfig := vpnAttachmentUnitConfig()
		for _, value := range oldConfig["tunnel_options_specification"].([]interface{}) {
			value.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = ""
		}
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vco-unit-test")
		config := vpnAttachmentUnitConfig()
		for _, value := range config["tunnel_options_specification"].([]interface{}) {
			value.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = ""
		}
		config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["ike_lifetime"] = 43200
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
		if _, err := r.Apply(old.State(), diff, client); err == nil {
			t.Fatal("updating IKE config without a known PSK must be rejected, not submitted for a random key reset")
		} else if !strings.Contains(err.Error(), "pre-shared key") {
			t.Errorf("error does not explain the pre-shared key requirement: %s", err)
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Errorf("IKE update without Psk reached ModifyVpnAttachmentAttribute %d times", got)
		}
	})

	t.Run("explicit_psk_applies", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		oldConfig := vpnAttachmentUnitConfig()
		for _, value := range oldConfig["tunnel_options_specification"].([]interface{}) {
			value.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = ""
		}
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vco-unit-test")
		config := vpnAttachmentUnitConfig()
		for _, value := range config["tunnel_options_specification"].([]interface{}) {
			value.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = ""
		}
		ike := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
		ike["psk"] = "recovered-test-psk"
		ike["ike_lifetime"] = 43200
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
		state, err := r.Apply(old.State(), diff, client)
		if err != nil {
			t.Fatal(err)
		}
		calls := requests.modifyCalls()
		if len(calls) != 1 {
			t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
		}
		if got := calls[0].Get("TunnelOptionsSpecification.1.TunnelIkeConfig.Psk"); got != "recovered-test-psk" {
			t.Errorf("TunnelIkeConfig.Psk = %q, want the explicitly configured recovered-test-psk", got)
		}
		if got := calls[0].Get("TunnelOptionsSpecification.1.TunnelIkeConfig.IkeLifetime"); got != "43200" {
			t.Errorf("TunnelIkeConfig.IkeLifetime = %q, want 43200", got)
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
	})
}

// The single-tunnel ike_config path follows the same rule: an IKE update
// without a known pre-shared key must be blocked instead of triggering a
// random key reset.
func TestUnitVpnGatewayVpnAttachmentUpdateTopLevelIkeWithoutPsk(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		oldConfig := vpnAttachmentUnitConfig()
		oldConfig["ike_config"] = []interface{}{map[string]interface{}{"psk": "", "ike_lifetime": 86400}}
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vco-unit-test")
		config := vpnAttachmentUnitConfig()
		config["ike_config"] = []interface{}{map[string]interface{}{"psk": "", "ike_lifetime": 43200}}
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
		if _, err := r.Apply(old.State(), diff, client); err == nil {
			t.Fatal("updating ike_config without a known PSK must be rejected, not submitted for a random key reset")
		} else if !strings.Contains(err.Error(), "pre-shared key") {
			t.Errorf("error does not explain the pre-shared key requirement: %s", err)
		}
		if got := len(requests.modifyCalls()); got != 0 {
			t.Errorf("IKE update without Psk reached ModifyVpnAttachmentAttribute %d times", got)
		}
	})

	t.Run("explicit_psk_applies", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		oldConfig := vpnAttachmentUnitConfig()
		oldConfig["ike_config"] = []interface{}{map[string]interface{}{"psk": "", "ike_lifetime": 86400}}
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vco-unit-test")
		config := vpnAttachmentUnitConfig()
		config["ike_config"] = []interface{}{map[string]interface{}{"psk": "explicit-top-level-psk", "ike_lifetime": 43200}}
		diff, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		api := vpnAttachmentUnitAPI(config)
		api["IkeConfig"] = map[string]interface{}{"Psk": "explicit-top-level-psk", "IkeLifetime": 43200}
		client, requests := vpnAttachmentUnitClient(t, api)
		state, err := r.Apply(old.State(), diff, client)
		if err != nil {
			t.Fatal(err)
		}
		calls := requests.modifyCalls()
		if len(calls) != 1 {
			t.Fatalf("ModifyVpnAttachmentAttribute calls = %d, want 1", len(calls))
		}
		body := calls[0].Get("IkeConfig")
		if !strings.Contains(body, `"Psk":"explicit-top-level-psk"`) {
			t.Errorf("IkeConfig = %s, want the explicitly configured Psk", body)
		}
		if calls[0].Has("TunnelOptionsSpecification.1.TunnelIndex") {
			t.Error("unchanged tunnels were included in a top-level IKE update")
		}
		vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
	})
}

func TestUnitVpnGatewayVpnAttachmentReadMaskedPSK(t *testing.T) {
	t.Run("empty_keeps_prior", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		config := vpnAttachmentUnitConfig()
		config["ike_config"] = []interface{}{map[string]interface{}{"psk": "old-top-level-test-psk"}}
		d := schema.TestResourceDataRaw(t, r.Schema, config)
		d.SetId("vco-unit-test")
		api := vpnAttachmentUnitAPI(vpnAttachmentUnitConfig())
		api["IkeConfig"] = map[string]interface{}{"Psk": ""}
		for _, raw := range api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{}) {
			raw.(map[string]interface{})["TunnelIkeConfig"].(map[string]interface{})["Psk"] = ""
		}
		client, _ := vpnAttachmentUnitClient(t, api)
		if err := r.Read(d, client); err != nil {
			t.Fatal(err)
		}
		if got := d.Get("ike_config.0.psk"); got != "old-top-level-test-psk" {
			t.Errorf("top-level PSK = %q, want the prior configured value", got)
		}
		for position, raw := range vpnAttachmentUnitStateTunnels(d) {
			ike := raw.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
			want := fmt.Sprintf("initial-test-psk-%d", position+1)
			if got := ike["psk"]; got != want {
				t.Errorf("tunnel %d PSK = %q, want the prior configured value %q", position+1, got, want)
			}
		}
		vpnAttachmentUnitAssertNoDiff(t, r, d.State(), config, client)
	})

	// '*' is a legal PSK character and the API documents no all-asterisk mask
	// form, so an all-asterisk response is a real value. Adopting it keeps an
	// externally rotated key visible instead of masking it behind the prior.
	t.Run("all_asterisk_adopted", func(t *testing.T) {
		r := resourceAliCloudVpnGatewayVpnAttachment()
		config := vpnAttachmentUnitConfig()
		config["ike_config"] = []interface{}{map[string]interface{}{"psk": "old-top-level-test-psk"}}
		d := schema.TestResourceDataRaw(t, r.Schema, config)
		d.SetId("vco-unit-test")
		api := vpnAttachmentUnitAPI(vpnAttachmentUnitConfig())
		api["IkeConfig"] = map[string]interface{}{"Psk": "************"}
		for _, raw := range api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{}) {
			raw.(map[string]interface{})["TunnelIkeConfig"].(map[string]interface{})["Psk"] = "************"
		}
		client, _ := vpnAttachmentUnitClient(t, api)
		if err := r.Read(d, client); err != nil {
			t.Fatal(err)
		}
		if got := d.Get("ike_config.0.psk"); got != "************" {
			t.Errorf("top-level PSK = %q, want the adopted API value, not the prior", got)
		}
		for position, raw := range vpnAttachmentUnitStateTunnels(d) {
			ike := raw.(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
			if got := ike["psk"]; got != "************" {
				t.Errorf("tunnel %d PSK = %q, want the adopted API value, not the prior", position+1, got)
			}
		}
		diff, err := r.Diff(d.State(), terraform.NewResourceConfigRaw(config), client)
		if err != nil {
			t.Fatal(err)
		}
		if diff == nil || diff.Empty() {
			t.Error("an externally rotated all-asterisk PSK must surface as a plan diff, not stay hidden")
		}
	})
}

type vpnAttachmentUnitRequests struct {
	mu        sync.Mutex
	modifies  []url.Values
	creates   []url.Values
	describes int
}

func (r *vpnAttachmentUnitRequests) modifyCalls() []url.Values {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]url.Values(nil), r.modifies...)
}

func vpnAttachmentUnitClient(t *testing.T, api map[string]interface{}) (*connectivity.AliyunClient, *vpnAttachmentUnitRequests) {
	t.Helper()
	return vpnAttachmentUnitClientResponses(t, []map[string]interface{}{api})
}

func vpnAttachmentUnitClientResponses(t *testing.T, responses []map[string]interface{}) (*connectivity.AliyunClient, *vpnAttachmentUnitRequests) {
	t.Helper()
	requests := &vpnAttachmentUnitRequests{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse RPC request: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("Action") {
		case "DescribeVpnConnection":
			requests.mu.Lock()
			index := requests.describes
			requests.describes++
			if index >= len(responses) {
				index = len(responses) - 1
			}
			response := responses[index]
			requests.mu.Unlock()
			_ = json.NewEncoder(w).Encode(response)
		case "ModifyVpnAttachmentAttribute":
			requests.mu.Lock()
			requests.modifies = append(requests.modifies, r.PostForm)
			requests.mu.Unlock()
			_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
		case "CreateVpnAttachment":
			requests.mu.Lock()
			requests.creates = append(requests.creates, r.PostForm)
			requests.mu.Unlock()
			_, _ = w.Write([]byte(`{"RequestId":"mock","VpnConnectionId":"vco-unit-test"}`))
		default:
			t.Errorf("unexpected RPC action: %s", r.URL.Query().Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"UnexpectedAction","Message":"unit-test"}`))
		}
	}))
	t.Cleanup(server.Close)
	endpoints := &sync.Map{}
	endpoints.Store("vpc", strings.TrimPrefix(server.URL, "http://"))
	staticCredential, err := credential.NewCredential(new(credential.Config).
		SetType("access_key").SetAccessKeyId("unit-test-access-key").SetAccessKeySecret("unit-test-secret-key"))
	if err != nil {
		t.Fatal(err)
	}
	config := &connectivity.Config{
		AccessKey: "unit-test-access-key", SecretKey: "unit-test-secret-key", Credential: staticCredential,
		RegionId: "cn-hangzhou", AccountType: "ResourceOwner", Protocol: "HTTP", SkipRegionValidation: true,
		Endpoints: endpoints, SignVersion: &sync.Map{}, ClientReadTimeout: 5000, ClientConnectTimeout: 5000,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return client, requests
}

func vpnAttachmentUnitAssertNoDiff(t *testing.T, r *schema.Resource, state *terraform.InstanceState, config map[string]interface{}, meta interface{}) {
	t.Helper()
	diff, err := r.Diff(state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("unchanged configuration produced a diff after Read: %#v", diff.Attributes)
	}
}

func vpnAttachmentUnitStateTunnels(d *schema.ResourceData) []interface{} {
	if set, ok := d.Get("tunnel_options_specification").(*schema.Set); ok {
		return set.List()
	}
	return d.Get("tunnel_options_specification").([]interface{})
}

func vpnAttachmentUnitConfig() map[string]interface{} {
	return map[string]interface{}{
		"local_subnet": "10.0.0.0/24", "remote_subnet": "10.1.0.0/24", "enable_tunnels_bgp": true,
		"tunnel_options_specification": []interface{}{vpnAttachmentUnitTunnel(1), vpnAttachmentUnitTunnel(2)},
	}
}

func vpnAttachmentUnitTunnel(index int) map[string]interface{} {
	return map[string]interface{}{
		"tunnel_index": index, "customer_gateway_id": "cgw-unit-test", "enable_dpd": true, "enable_nat_traversal": true,
		"tunnel_ike_config": []interface{}{map[string]interface{}{
			"psk": fmt.Sprintf("initial-test-psk-%d", index), "ike_auth_alg": "sha1", "ike_enc_alg": "aes", "ike_version": "ikev1", "ike_mode": "main", "ike_lifetime": 86400, "ike_pfs": "group2", "local_id": "initial-local", "remote_id": "initial-remote",
		}},
		"tunnel_bgp_config": []interface{}{map[string]interface{}{
			"local_asn": 65000, "local_bgp_ip": "169.254.1.1", "tunnel_cidr": "169.254.1.0/30",
		}},
		"tunnel_ipsec_config": []interface{}{map[string]interface{}{
			"ipsec_auth_alg": "sha1", "ipsec_enc_alg": "aes", "ipsec_lifetime": 3600, "ipsec_pfs": "group2",
		}},
	}
}

// Distinct per-tunnel values make positional backfill across a reorder
// observable: inherited values always differ from the tunnel's own.
func vpnAttachmentDistinctConfig() map[string]interface{} {
	return map[string]interface{}{
		"local_subnet": "10.0.0.0/24", "remote_subnet": "10.1.0.0/24", "enable_tunnels_bgp": true,
		"tunnel_options_specification": []interface{}{vpnAttachmentDistinctTunnel(1), vpnAttachmentDistinctTunnel(2)},
	}
}

func vpnAttachmentDistinctTunnel(index int) map[string]interface{} {
	return map[string]interface{}{
		"tunnel_index": index, "customer_gateway_id": fmt.Sprintf("cgw-unit-%d", index),
		"enable_dpd": index == 1, "enable_nat_traversal": index == 1,
		"tunnel_ike_config": []interface{}{map[string]interface{}{
			"psk": fmt.Sprintf("tunnel-%d-test-psk", index), "ike_auth_alg": "sha1", "ike_enc_alg": "aes", "ike_version": "ikev1", "ike_mode": "main", "ike_lifetime": 86400, "ike_pfs": "group2",
			"local_id": fmt.Sprintf("tunnel-%d-local", index), "remote_id": fmt.Sprintf("tunnel-%d-remote", index),
		}},
		"tunnel_bgp_config": []interface{}{map[string]interface{}{
			"local_asn": 65000 + index, "local_bgp_ip": fmt.Sprintf("169.254.%d.1", index), "tunnel_cidr": fmt.Sprintf("169.254.%d.0/30", index),
		}},
		"tunnel_ipsec_config": []interface{}{map[string]interface{}{
			"ipsec_auth_alg": "sha1", "ipsec_enc_alg": []string{"", "aes", "aes256"}[index], "ipsec_lifetime": 3600, "ipsec_pfs": "group2",
		}},
	}
}

// Include service-owned defaults and fields omitted by the configuration.
func vpnAttachmentUnitAPI(config map[string]interface{}) map[string]interface{} {
	tunnels := make([]interface{}, 0, 2)
	for _, value := range config["tunnel_options_specification"].([]interface{}) {
		m := value.(map[string]interface{})
		index := m["tunnel_index"].(int)
		ike := m["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
		ipsec := m["tunnel_ipsec_config"].([]interface{})[0].(map[string]interface{})
		bgp := m["tunnel_bgp_config"].([]interface{})[0].(map[string]interface{})
		tunnels = append(tunnels, map[string]interface{}{
			"TunnelIndex": index, "CustomerGatewayId": m["customer_gateway_id"], "TunnelId": fmt.Sprintf("tun-unit-%d", index),
			"EnableDpd": m["enable_dpd"], "EnableNatTraversal": m["enable_nat_traversal"], "Role": "master", "State": "attached", "Status": "ipsec_sa_established",
			"InternetIp": fmt.Sprintf("192.0.2.%d", index), "ZoneNo": fmt.Sprintf("zone-%d", index),
			"TunnelIkeConfig": map[string]interface{}{
				"Psk": ike["psk"], "IkeAuthAlg": ike["ike_auth_alg"], "IkeEncAlg": ike["ike_enc_alg"], "IkeVersion": ike["ike_version"], "IkeMode": ike["ike_mode"],
				"IkeLifetime": ike["ike_lifetime"], "IkePfs": ike["ike_pfs"], "LocalId": ike["local_id"], "RemoteId": ike["remote_id"],
			},
			"TunnelIpsecConfig": map[string]interface{}{
				"IpsecAuthAlg": ipsec["ipsec_auth_alg"], "IpsecEncAlg": ipsec["ipsec_enc_alg"], "IpsecLifetime": ipsec["ipsec_lifetime"], "IpsecPfs": ipsec["ipsec_pfs"],
			},
			"TunnelBgpConfig": map[string]interface{}{
				"LocalAsn": bgp["local_asn"], "LocalBgpIp": bgp["local_bgp_ip"], "TunnelCidr": bgp["tunnel_cidr"], "PeerAsn": 65501, "PeerBgpIp": "169.254.1.2", "BgpStatus": "success",
			},
		})
	}
	return map[string]interface{}{
		"RequestId": "mock", "VpnConnectionId": "vco-unit-test", "LocalSubnet": config["local_subnet"], "RemoteSubnet": config["remote_subnet"],
		"EnableTunnelsBgp": config["enable_tunnels_bgp"], "NetworkType": "public", "State": "attached", "TunnelBandwidth": "standard",
		"TunnelOptionsSpecification": map[string]interface{}{"TunnelOptions": tunnels},
	}
}

// Give each tunnel its own Computed-only BGP status and peer values like a
// real API response: positional inheritance across a reorder then differs
// from each tunnel's own values, which is what exposed the guard's false
// conflicts on non-writable leaves.
func vpnAttachmentUnitAPIDistinctComputed(config map[string]interface{}) map[string]interface{} {
	api := vpnAttachmentUnitAPI(config)
	for _, raw := range api["TunnelOptionsSpecification"].(map[string]interface{})["TunnelOptions"].([]interface{}) {
		tunnel := raw.(map[string]interface{})
		index := tunnel["TunnelIndex"].(int)
		bgp := tunnel["TunnelBgpConfig"].(map[string]interface{})
		bgp["PeerAsn"] = fmt.Sprintf("655%02d", index)
		bgp["PeerBgpIp"] = fmt.Sprintf("169.254.%d.2", index)
		bgp["BgpStatus"] = []string{"", "success", "established"}[index]
	}
	return api
}

func TestAccAliCloudVPNGatewayVpnAttachment_tunnelOrder(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc%svpnattachmentorder%d", defaultRegionToTest, rand)
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
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"effect_immediately":  "false",
					"vpn_attachment_name": "${var.name}",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
							"role":                "master",
							"tunnel_index":        "1",
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultone.id}",
							"role":                "master",
							"tunnel_index":        "2",
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
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"effect_immediately":  "false",
					"vpn_attachment_name": "${var.name}",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultone.id}",
							"role":                "master",
							"tunnel_index":        "2",
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
							"role":                "master",
							"tunnel_index":        "1",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":        "public",
					"local_subnet":        "0.0.0.0/0",
					"remote_subnet":       "0.0.0.0/0",
					"effect_immediately":  "false",
					"vpn_attachment_name": "${var.name}",
					"tunnel_options_specification": []map[string]interface{}{
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.defaultone.id}",
							"role":                "master",
							"tunnel_index":        "2",
						},
						{
							"customer_gateway_id": "${alicloud_vpn_customer_gateway.default.id}",
							"role":                "master",
							"tunnel_index":        "1",
						},
					},
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
