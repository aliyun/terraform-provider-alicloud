package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"net"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	p := Provider().ResourcesMap
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				// Nested changes must be visible even when the tunnel index
				// and customer gateway are unchanged.
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	var originalID string
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					vpnAttachmentCheckRolePhase(resourceId, &originalID, 0),
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
									"ipsec_pfs":      "group14",
									"ipsec_enc_alg":  "aes256",
									"ipsec_auth_alg": "sha1",
									"ipsec_lifetime": "43200",
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
				Check: resource.ComposeTestCheckFunc(vpnAttachmentCheckRolePhase(resourceId, &originalID, 1), func(*terraform.State) error { t.Log("SAFE_IPSEC_ONLY_UPDATE_VERIFIED"); return nil }),
			},
			{Config: testAccConfig(nil), PlanOnly: true, ExpectNonEmptyPlan: false},
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
					vpnAttachmentCheckRolePhase(resourceId, &originalID, 2),
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
			{Config: testAccConfig(nil), PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: testAccConfig(nil), Check: vpnAttachmentCheckRolePhase(resourceId, &originalID, 2)},
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Legacy role configuration is ignored; state retains the service-reported role.
func TestAccAliCloudVpnGatewayVpnAttachment_roleSwitch(t *testing.T) {
	var v map[string]interface{}
	var originalID string
	checkRole := func(state *terraform.State) error {
		entry := state.RootModule().Resources["alicloud_vpn_gateway_vpn_attachment.default"]
		if originalID == "" {
			originalID = entry.Primary.ID
		}
		if originalID != entry.Primary.ID {
			return fmt.Errorf("SAFE_UNEXPECTED_RECREATION")
		}
		service := VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		obj, err := service.DescribeVpnGatewayVpnAttachment(entry.Primary.ID)
		if err != nil {
			return fmt.Errorf("SAFE_VPN_LIVE_READ_FAILED")
		}
		outer, ok := obj["TunnelOptionsSpecification"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("SAFE_TUNNEL_CONFIGURATION_MISMATCH")
		}
		tunnels, ok := outer["TunnelOptions"].([]interface{})
		if !ok || len(tunnels) != 2 {
			return fmt.Errorf("SAFE_TUNNEL_CONFIGURATION_MISMATCH")
		}
		seen := map[int]bool{}
		for _, raw := range tunnels {
			tunnel := raw.(map[string]interface{})
			index := formatInt(tunnel["TunnelIndex"])
			if index < 1 || index > 2 || seen[index] || vpnAttachmentStateAttrByIndex(entry.Primary, index, "role") != fmt.Sprint(tunnel["Role"]) {
				return fmt.Errorf("SAFE_TUNNEL_CONFIGURATION_MISMATCH")
			}
			seen[index] = true
		}
		return nil
	}
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					checkRole,
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
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{Config: testAccConfig(nil), Check: checkRole},
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
			if _, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil); err == nil {
				t.Fatal("expected a plan error: an omitted API parameter cannot remove existing tunnel configuration")
			}
		})
	}
}

// Changing a ForceNew field plans a replacement of the whole attachment, so
// the in-place guards against removing tunnels or clearing nested blocks do
// not apply: the removed configuration is simply not part of the new
// resource. A plan without a ForceNew change keeps the guards, covered by
// TestUnitVpnGatewayVpnAttachmentRejectsRemovedConfiguration.
func TestUnitVpnGatewayVpnAttachmentRebuildAllowsRemovedConfiguration(t *testing.T) {
	cases := []struct {
		name        string
		forceNewKey string
		forceNewVal string
		rewrite     func(config map[string]interface{})
	}{
		{"tunnel_bandwidth_change_allows_removed_tunnel", "tunnel_bandwidth", "Large", func(config map[string]interface{}) {
			config["tunnel_options_specification"] = config["tunnel_options_specification"].([]interface{})[:1]
		}},
		{"tunnel_bandwidth_change_allows_cleared_bgp", "tunnel_bandwidth", "Large", func(config map[string]interface{}) {
			config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_bgp_config"] = []interface{}{}
		}},
		{"network_type_change_allows_removed_tunnel", "network_type", "private", func(config map[string]interface{}) {
			config["tunnel_options_specification"] = config["tunnel_options_specification"].([]interface{})[:1]
		}},
		{"network_type_change_allows_cleared_ike", "network_type", "private", func(config map[string]interface{}) {
			config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"] = []interface{}{}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentUnitConfig())
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			config[tc.forceNewKey] = tc.forceNewVal
			tc.rewrite(config)
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("changing ForceNew field %s plans a rebuild; removing the configuration must not be rejected by the in-place guards: %s", tc.forceNewKey, err)
			}
			if diff == nil || !diff.RequiresNew() {
				t.Fatalf("changing ForceNew field %s must plan a replacement, not an in-place update", tc.forceNewKey)
			}
		})
	}
}

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
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
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
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
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
	config["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["role"] = "slave"
	diff, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	state, diags := r.Apply(context.Background(), nil, diff, client)
	if diags.HasError() {
		t.Fatalf("%v", diags)
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
		if calls[0].Has(prefix + "Role") {
			t.Fatal("read-only role sent to CreateVpnAttachment")
		}
		if vpnAttachmentStateAttrByIndex(state, i, "role") != "master" {
			t.Fatal("configured role replaced service-reported state")
		}
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

func TestUnitVpnGatewayVpnAttachmentEnableBgpOnly(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	oldConfig := vpnAttachmentUnitConfig()
	oldConfig["enable_tunnels_bgp"] = false
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vco-unit-test")
	config := vpnAttachmentUnitConfig()
	diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	for key := range diff.Attributes {
		if strings.HasPrefix(key, "tunnel_options_specification.") {
			t.Fatalf("fixture should change only enable_tunnels_bgp, got %s", key)
		}
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	state, diags := r.Apply(context.Background(), old.State(), diff, client)
	if diags.HasError() {
		t.Fatalf("%v", diags)
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
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
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
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
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
			if got := vpnAttachmentStateAttrByIndex(state, 1, "tunnel_ike_config.0.psk"); got != "updated-after-completion" {
				t.Errorf("Apply retained PSK from before asynchronous update: %q", got)
			}
			if got := vpnAttachmentStateAttrByIndex(state, 1, "tunnel_bgp_config.0.local_asn"); got != "65001" {
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
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			failed := vpnAttachmentUnitAPI(oldConfig)
			failed["State"] = failedState
			client, requests := vpnAttachmentUnitClient(t, failed)
			if _, diags := r.Apply(context.Background(), old.State(), diff, client); !diags.HasError() {
				t.Fatalf("Apply reported success for failed state %s", failedState)
			} else if !strings.Contains(fmt.Sprintf("%v", diags), failedState) {
				t.Errorf("error does not identify failed state %s: %v", failedState, diags)
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
		diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
		if _, diags := r.Apply(context.Background(), old.State(), diff, client); !diags.HasError() {
			t.Fatal("updating IKE config without a known PSK must be rejected, not submitted for a random key reset")
		} else if !strings.Contains(fmt.Sprintf("%v", diags), "pre-shared key") {
			t.Errorf("error does not explain the pre-shared key requirement: %v", diags)
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
		diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
		state, diags := r.Apply(context.Background(), old.State(), diff, client)
		if diags.HasError() {
			t.Fatalf("%v", diags)
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
	for _, priorPsk := range []string{"", "prior-top-level-psk"} {
		t.Run("blocked_empty_psk_from_"+map[bool]string{true: "empty_state", false: "known_state"}[priorPsk == ""], func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			oldConfig := vpnAttachmentUnitConfig()
			oldConfig["ike_config"] = []interface{}{map[string]interface{}{"psk": priorPsk, "ike_lifetime": 86400}}
			old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
			old.SetId("vco-unit-test")
			config := vpnAttachmentUnitConfig()
			config["ike_config"] = []interface{}{map[string]interface{}{"psk": "", "ike_lifetime": 43200}}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), client)
			if err == nil {
				if _, diags := r.Apply(context.Background(), old.State(), diff, client); diags.HasError() {
					err = fmt.Errorf("%v", diags)
				}
			}
			if err == nil {
				t.Fatal("updating ike_config without a known PSK must be rejected, not submitted for a random key reset")
			} else if !strings.Contains(err.Error(), "pre-shared key") {
				t.Errorf("error does not explain the pre-shared key requirement: %s", err)
			}
			if got := len(requests.modifyCalls()); got != 0 {
				t.Errorf("IKE update without Psk reached ModifyVpnAttachmentAttribute %d times", got)
			}
		})

	}
	for _, mode := range []string{"explicit_psk_applies", "omitted_psk_inherits", "unknown_psk_resolved"} {
		t.Run(mode, func(t *testing.T) {
			omitPsk := mode == "omitted_psk_inherits"
			r := resourceAliCloudVpnGatewayVpnAttachment()
			oldConfig := vpnAttachmentUnitConfig()
			priorPsk := ""
			if mode != "explicit_psk_applies" {
				priorPsk = "explicit-top-level-psk"
			}
			oldConfig["ike_config"] = []interface{}{map[string]interface{}{"psk": priorPsk, "ike_lifetime": 86400}}
			old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
			old.SetId("vco-unit-test")

			if mode == "unknown_psk_resolved" {
				pending := vpnAttachmentUnitConfig()
				pending["ike_config"] = []interface{}{map[string]interface{}{"psk": "74D93920-ED26-11E3-AC10-0800200C9A66", "ike_lifetime": 43200}}
				pendingDiff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(pending), nil)
				if err != nil {
					t.Fatal(err)
				}
				if attr := pendingDiff.Attributes["ike_config.0.psk"]; attr == nil || !attr.NewComputed {
					t.Fatal("unknown PSK must remain computed at its own field")
				}
			}
			config := vpnAttachmentUnitConfig()
			config["ike_config"] = []interface{}{map[string]interface{}{"psk": "explicit-top-level-psk", "ike_lifetime": 43200}}
			if omitPsk {
				delete(config["ike_config"].([]interface{})[0].(map[string]interface{}), "psk")
			}
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			api := vpnAttachmentUnitAPI(config)
			api["IkeConfig"] = map[string]interface{}{"Psk": "explicit-top-level-psk", "IkeLifetime": 43200}
			client, requests := vpnAttachmentUnitClient(t, api)
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
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

	for _, replacement := range []bool{false, true} {
		t.Run(map[bool]string{false: "create_empty_psk_allowed", true: "replacement_empty_psk_allowed"}[replacement], func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			config := vpnAttachmentUnitConfig()
			config["ike_config"] = []interface{}{map[string]interface{}{"psk": "", "ike_lifetime": 43200}}
			var prior *terraform.InstanceState
			if replacement {
				oldConfig := vpnAttachmentUnitConfig()
				oldConfig["ike_config"] = []interface{}{map[string]interface{}{"psk": "prior-top-level-psk", "ike_lifetime": 86400}}
				old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
				old.SetId("vco-unit-test")
				prior = old.State()
				config["network_type"] = "private"
			}
			if _, err := r.Diff(context.Background(), prior, terraform.NewResourceConfigRaw(config), nil); err != nil {
				t.Fatal(err)
			}
		})
	}
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
		diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), client)
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
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("unchanged configuration produced a diff after Read: %#v", diff.Attributes)
	}
}

func vpnAttachmentUnitStateTunnels(d *schema.ResourceData) []interface{} {
	if set, ok := d.Get("tunnel_options_specification").(*schema.Set); ok {
		return vpnAttachmentSortedTunnels(set.List())
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

// Give each tunnel distinct service-owned values to exercise hash exclusion.
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
	const vpn = "alicloud_vpn_gateway_vpn_attachment.default"
	name := fmt.Sprintf("tfaccbackfill%d", acctest.RandIntRange(100000, 999999))
	before := map[int]string{}
	originalID := ""
	unassigned := map[int]bool{}
	after := map[int]string{}
	afterIP := map[int]string{}
	live := func(s *terraform.State, bound bool) error {
		entry, ok := s.RootModule().Resources[vpn]
		if !ok {
			return fmt.Errorf("SAFE_VPN_STATE_MISSING")
		}
		service := VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		deadline := time.Now().Add(3 * time.Minute)
		for {
			obj, err := service.DescribeVpnGatewayVpnAttachment(entry.Primary.ID)
			if err != nil {
				return fmt.Errorf("SAFE_VPN_LIVE_READ_FAILED")
			}
			outer, ok := obj["TunnelOptionsSpecification"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("SAFE_TUNNEL_DATA_MISSING")
			}
			tunnels, ok := outer["TunnelOptions"].([]interface{})
			if !ok || len(tunnels) != 2 {
				return fmt.Errorf("SAFE_TUNNEL_COUNT_MISMATCH")
			}
			ready := true
			seen := map[int]bool{}
			for _, raw := range tunnels {
				tunnel, ok := raw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("SAFE_TUNNEL_DATA_MISSING")
				}
				index := formatInt(tunnel["TunnelIndex"])
				if (index != 1 && index != 2) || seen[index] {
					return fmt.Errorf("SAFE_TUNNEL_INDEX_MISMATCH")
				}
				seen[index] = true
				cgw, ok := s.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", index)]
				if !ok || tunnel["CustomerGatewayId"] != cgw.Primary.ID {
					return fmt.Errorf("SAFE_TUNNEL_GATEWAY_MISMATCH")
				}
				ike, _ := tunnel["TunnelIkeConfig"].(map[string]interface{})
				local, _ := ike["LocalId"].(string)
				ip, _ := tunnel["InternetIp"].(string)
				if !bound {
					before[index] = local
					unassigned[index] = !vpnAttachmentCenReproAssignedIP(ip) && !vpnAttachmentCenReproAssignedIP(local)
					if unassigned[index] {
						t.Logf("SAFE_TUNNEL%d_INITIAL_IP_UNASSIGNED", index)
					}
				} else {
					after[index] = local
					afterIP[index] = ip
					ready = ready && vpnAttachmentCenReproAssignedIP(ip) && vpnAttachmentCenReproAssignedIP(local)
				}
			}
			if !bound {
				originalID = entry.Primary.ID
				t.Log("SAFE_UNBOUND_SNAPSHOT_VERIFIED")
				return nil
			}
			if ready {
				changed := 0
				for _, index := range []int{1, 2} {
					if unassigned[index] && before[index] != after[index] {
						changed++
						t.Logf("SAFE_TUNNEL%d_BACKFILL_CHANGED", index)
					}
				}
				if changed == 2 {
					t.Log("SAFE_BOUND_GATEWAYS_AND_IPS_VERIFIED")
					return nil
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("SAFE_BACKFILL_NOT_OBSERVED")
			}
			time.Sleep(5 * time.Second)
		}
	}
	verifyState := func(s *terraform.State) error {
		entry, ok := s.RootModule().Resources[vpn]
		if !ok {
			return fmt.Errorf("SAFE_VPN_STATE_MISSING")
		}
		if entry.Primary.ID != originalID {
			return fmt.Errorf("SAFE_UNEXPECTED_RECREATION")
		}
		for _, index := range []int{1, 2} {
			if vpnAttachmentStateAttrByIndex(entry.Primary, index, "tunnel_ike_config.0.psk") != fmt.Sprintf("tf-acc-backfill-%d", index) {
				return fmt.Errorf("SAFE_PSK_RETENTION_MISMATCH")
			}
			if vpnAttachmentStateAttrByIndex(entry.Primary, index, "internet_ip") != afterIP[index] {
				return fmt.Errorf("SAFE_REFRESH_INTERNET_IP_MISMATCH")
			}
			if vpnAttachmentStateAttrByIndex(entry.Primary, index, "tunnel_ike_config.0.local_id") != after[index] {
				return fmt.Errorf("SAFE_REFRESH_LOCAL_ID_MISMATCH")
			}
			if vpnAttachmentStateAttrByIndex(entry.Primary, index, "customer_gateway_id") != s.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", index)].Primary.ID {
				return fmt.Errorf("SAFE_TUNNEL_GATEWAY_MISMATCH")
			}
		}
		t.Log("SAFE_REPEAT_APPLY_STATE_VERIFIED")
		return nil
	}
	checkDestroy := func(s *terraform.State) error {
		client := func() *connectivity.AliyunClient { return testAccProvider.Meta().(*connectivity.AliyunClient) }
		checks := []struct {
			id, method string
			service    func() interface{}
		}{
			{vpn, "DescribeVpnGatewayVpnAttachment", func() interface{} { return &VPNGatewayServiceV2{client()} }},
			{"alicloud_cen_transit_router_vpn_attachment.default", "DescribeCenTransitRouterVpnAttachment", func() interface{} { return &CbnService{client()} }},
			{"alicloud_cen_transit_router_cidr.default", "DescribeCenTransitRouterCidr", func() interface{} { return &CbnService{client()} }},
			{"alicloud_cen_transit_router.default", "DescribeCenTransitRouter", func() interface{} { return &CenServiceV2{client()} }},
			{"alicloud_cen_instance.default", "DescribeCenInstance", func() interface{} { return &CbnService{client()} }},
		}
		for _, check := range checks {
			var object map[string]interface{}
			rc := resourceCheckInitWithDescribeMethod(check.id, &object, check.service, check.method)
			if err := rc.checkResourceDestroy()(s); err != nil {
				return fmt.Errorf("SAFE_DEPENDENCY_DESTROY_FAILED")
			}
		}
		if err := testAccCheckVpnCustomerGatewayDestroy(s); err != nil {
			return fmt.Errorf("SAFE_CUSTOMER_GATEWAY_DESTROY_FAILED")
		}
		t.Log("SAFE_ALL_RESOURCES_DESTROY_VERIFIED")
		return nil
	}
	unbound := vpnAttachmentCenBackfillConfig(name, false, false)
	reversed := vpnAttachmentCenBackfillConfig(name, false, true)
	bound := vpnAttachmentCenBackfillConfig(name, true, true)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory, CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{Config: unbound, Check: func(s *terraform.State) error { return live(s, false) }},
			{Config: reversed, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: reversed, Check: func(s *terraform.State) error {
				e := s.RootModule().Resources[vpn]
				if e == nil || e.Primary.ID != originalID {
					return fmt.Errorf("SAFE_UNEXPECTED_RECREATION")
				}
				for _, index := range []int{1, 2} {
					if vpnAttachmentStateAttrByIndex(e.Primary, index, "tunnel_ike_config.0.psk") != fmt.Sprintf("tf-acc-backfill-%d", index) || vpnAttachmentStateAttrByIndex(e.Primary, index, "customer_gateway_id") != s.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", index)].Primary.ID {
						return fmt.Errorf("SAFE_TUNNEL_CONFIGURATION_MISMATCH")
					}
				}
				return nil
			}},
			{Config: bound, Check: func(s *terraform.State) error { return live(s, true) }},
			{Config: bound, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: bound, Check: verifyState},
			{ResourceName: vpn, ImportState: true, ImportStateVerify: true},
			{Config: bound, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: bound, Check: verifyState},
		},
	})
}

func vpnAttachmentCenBackfillConfig(name string, bound, reverse bool) string {
	config := AlicloudVpnGatewayVpnAttachmentBasicDependenceRole(name) + `
resource "alicloud_cen_instance" "default" {
 cen_instance_name=var.name
}
resource "alicloud_cen_transit_router" "default" {
 cen_id=alicloud_cen_instance.default.id
 transit_router_name=var.name
}
resource "alicloud_cen_transit_router_cidr" "default" {
 transit_router_id=alicloud_cen_transit_router.default.transit_router_id
 cidr="192.168.0.0/16"
 publish_cidr_route=false
}
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
 vpn_attachment_name=var.name
 network_type="public"
 local_subnet="0.0.0.0/0"
 remote_subnet="0.0.0.0/0"
 enable_tunnels_bgp=false
 effect_immediately=false
`
	indices := []int{1, 2}
	if reverse {
		indices = []int{2, 1}
	}
	for _, index := range indices {
		config += fmt.Sprintf(`
 tunnel_options_specification {
  tunnel_index=%d
  customer_gateway_id=alicloud_vpn_customer_gateway.cgw%d.id
  role="master"
  enable_dpd=true
  enable_nat_traversal=true
  tunnel_ike_config {
   psk="tf-acc-backfill-%d"
   ike_version="ikev2"
   ike_mode="main"
   ike_auth_alg="sha1"
   ike_enc_alg="aes"
   ike_pfs="group2"
   ike_lifetime=86400
  }
 }
`, index, index, index)
	}
	config += "}\n"
	if bound {
		config += `
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
 cen_id=alicloud_cen_instance.default.id
 transit_router_id=alicloud_cen_transit_router_cidr.default.transit_router_id
 vpn_id=alicloud_vpn_gateway_vpn_attachment.default.id
 transit_router_vpn_attachment_name=var.name
 auto_publish_route_enabled=false
 order_type="PayByCenOwner"
}
`
	}
	return config
}

func vpnAttachmentCenReproAssignedIP(value string) bool {
	ip := net.ParseIP(value)
	return ip != nil && !ip.IsUnspecified()
}

func TestUnitVpnGatewayVpnAttachmentSetReorderNoDiff(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	client, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentDistinctConfig()))
	refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), old.State(), client)
	if diags.HasError() {
		t.Fatalf("%v", diags)
	}
	vpnAttachmentUnitAssertNoDiff(t, r, refreshed, config, client)
}

func TestUnitVpnGatewayVpnAttachmentSetExplicitWholeIke(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	original := vpnAttachmentDistinctConfig()
	old := schema.TestResourceDataRaw(t, r.Schema, original)
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
	tunnels[0].(map[string]interface{})["tunnel_ike_config"] = vpnAttachmentDistinctTunnel(1)["tunnel_ike_config"]
	diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	state, diags := r.Apply(context.Background(), old.State(), diff, client)
	if diags.HasError() {
		t.Fatalf("%v", diags)
	}
	calls := requests.modifyCalls()
	if len(calls) != 1 {
		t.Fatalf("explicit IKE replacement sent %d modify calls, want 1", len(calls))
	}
	if calls[0].Get("TunnelOptionsSpecification.1.TunnelIndex") != "2" || calls[0].Get("TunnelOptionsSpecification.1.TunnelIkeConfig.LocalId") != "tunnel-1-local" {
		t.Fatal("explicit IKE replacement was not routed to tunnel 2")
	}
	vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
}

func TestUnitVpnGatewayVpnAttachmentSetRejectsDuplicateIndex(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentDistinctConfig()
	tunnels := config["tunnel_options_specification"].([]interface{})
	tunnels[1].(map[string]interface{})["tunnel_index"] = 1
	if _, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil); err == nil {
		t.Fatal("two different configurations for tunnel_index 1 must fail planning")
	}
}

func TestUnitVpnGatewayVpnAttachmentSetUnknownValues(t *testing.T) {
	for _, field := range []string{"tunnel_index", "customer_gateway_id", "enable_dpd", "tunnel_ike_config", "psk"} {
		t.Run(field, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
			old.SetId("vco-unit-test")
			config := vpnAttachmentDistinctConfig()
			tunnel := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			if field == "psk" {
				tunnel = tunnel["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
			}
			tunnel[field] = "74D93920-ED26-11E3-AC10-0800200C9A66"
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			computed := false
			if diff != nil {
				for key, change := range diff.Attributes {
					suffix := "." + field
					if field == "tunnel_ike_config" {
						suffix += ".#"
					}
					if strings.HasPrefix(key, "tunnel_options_specification.~") && strings.HasSuffix(key, suffix) && change.NewComputed {
						computed = true
					}
				}
			}
			if !computed {
				t.Fatal(fmt.Sprintf("unknown %s was replaced with a known prior value", field))
			}
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentSetPartialUpdate(t *testing.T) {
	for _, omitBlock := range []bool{false, true} {
		t.Run(fmt.Sprint(omitBlock), func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			full := vpnAttachmentDistinctConfig()
			if omitBlock {
				full["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})["enable_dpd"] = true
			}
			old := schema.TestResourceDataRaw(t, r.Schema, full)
			old.SetId("vco-unit-test")
			config := vpnAttachmentDistinctConfig()
			tunnels := config["tunnel_options_specification"].([]interface{})
			tunnels[0], tunnels[1] = tunnels[1], tunnels[0]
			for _, raw := range tunnels {
				tunnel := raw.(map[string]interface{})
				for _, key := range []string{"enable_dpd", "enable_nat_traversal", "role", "tunnel_bgp_config", "tunnel_ipsec_config"} {
					delete(tunnel, key)
				}
				if omitBlock {
					delete(tunnel, "tunnel_ike_config")
				} else {
					tunnel["tunnel_ike_config"] = []interface{}{map[string]interface{}{"psk": fmt.Sprintf("distinct-psk-%d", tunnel["tunnel_index"])}}
				}
			}
			// Change only a top-level flag while all nested blocks are omitted, or a
			// single IKE leaf while its other server-populated fields are omitted.
			target := tunnels[0].(map[string]interface{})
			expected := vpnAttachmentDistinctConfig()
			expectedTarget := expected["tunnel_options_specification"].([]interface{})[1].(map[string]interface{})
			wire, want := "EnableDpd", "false"
			if omitBlock {
				target["enable_dpd"] = false
				expectedTarget["enable_dpd"] = false
			} else {
				target["tunnel_ike_config"] = []interface{}{map[string]interface{}{"local_id": "updated-local"}}
				expectedTarget["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["local_id"] = "updated-local"
				wire, want = "TunnelIkeConfig.LocalId", "updated-local"
			}
			// Keep the other tunnel wholly omitted so this update targets only index 2.
			delete(tunnels[1].(map[string]interface{}), "tunnel_ike_config")
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(expected))
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			calls := requests.modifyCalls()
			if len(calls) != 1 {
				t.Fatalf("modify calls=%d", len(calls))
			}
			if calls[0].Get("TunnelOptionsSpecification.1.TunnelIndex") != "2" || calls[0].Get("TunnelOptionsSpecification.1."+wire) != want {
				t.Fatal("partial update not routed by index")
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentSetLegacyRefreshAndImport(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	config := vpnAttachmentDistinctConfig()
	legacy := resourceAliCloudVpnGatewayVpnAttachment()
	legacy.Schema["tunnel_options_specification"].Set = func(v interface{}) int {
		m := v.(map[string]interface{})
		return int(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%d-%s", m["tunnel_index"].(int), m["customer_gateway_id"].(string)))))
	}
	old := schema.TestResourceDataRaw(t, legacy.Schema, config)
	old.SetId("vco-unit-test")
	client, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	for _, imported := range []bool{false, true} {
		t.Run(fmt.Sprint(imported), func(t *testing.T) {
			state := old.State()
			if imported {
				data := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
				data.SetId("vco-unit-test")
				results, err := r.Importer.State(data, client)
				if err != nil || len(results) != 1 {
					t.Fatal("import failed")
				}
				state = results[0].State()
			}
			refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), state, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			vpnAttachmentUnitAssertNoDiff(t, r, refreshed, config, client)
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentSetHashWritableOnly(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	setSchema := r.Schema["tunnel_options_specification"]
	hash := setSchema.Set
	fields := setSchema.Elem.(*schema.Resource).Schema
	if !fields["role"].Optional {
		t.Fatal("hash schema changed public role configuration compatibility")
	}
	var walk func(map[string]*schema.Schema, []string)
	walk = func(fields map[string]*schema.Schema, path []string) {
		for name, field := range fields {
			leaf := append(append([]string{}, path...), name)
			if child, ok := field.Elem.(*schema.Resource); ok {
				walk(child.Schema, leaf)
				continue
			}
			t.Run(strings.Join(leaf, "."), func(t *testing.T) {
				old, next := vpnAttachmentDistinctTunnel(1), vpnAttachmentDistinctTunnel(1)
				target := next
				for _, group := range path {
					target = target[group].([]interface{})[0].(map[string]interface{})
				}
				switch field.Type {
				case schema.TypeBool:
					value, _ := target[name].(bool)
					target[name] = !value
				case schema.TypeInt:
					target[name] = formatInt(target[name]) + 1
				case schema.TypeString:
					target[name] = fmt.Sprint(target[name]) + "-changed"
				default:
					t.Fatalf("uncovered hash field type for %s", name)
				}
				ignored := (!field.Optional && !field.Required) || (len(path) == 0 && name == "role")
				if (hash(old) == hash(next)) != ignored {
					t.Fatalf("hash inclusion incorrect for %s", strings.Join(leaf, "."))
				}
			})
		}
	}
	walk(fields, nil)
}

func TestUnitVpnGatewayVpnAttachmentRoleCompatibility(t *testing.T) {
	for _, mixed := range []bool{false, true} {
		t.Run(fmt.Sprint(mixed), func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
			old.SetId("vco-unit-test")
			baselineClient, _ := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(vpnAttachmentDistinctConfig()))
			baseline, diags := r.RefreshWithoutUpgrade(context.Background(), old.State(), baselineClient)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			config := vpnAttachmentDistinctConfig()
			tunnel := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			tunnel["role"] = "slave"
			if mixed {
				tunnel["tunnel_ipsec_config"].([]interface{})[0].(map[string]interface{})["ipsec_lifetime"] = 12345
			}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
			diff, err := r.Diff(context.Background(), baseline, terraform.NewResourceConfigRaw(config), client)
			if err != nil {
				t.Fatal(err)
			}
			state := baseline
			if !mixed {
				if diff != nil && !diff.Empty() {
					t.Fatalf("role-only configuration must have no diff: %#v", diff.Attributes)
				}
			} else {
				applied, applyDiags := r.Apply(context.Background(), baseline, diff, client)
				if applyDiags.HasError() {
					t.Fatalf("%v", applyDiags)
				}
				state = applied
			}
			calls := requests.modifyCalls()
			expectedCalls := 0
			if mixed {
				expectedCalls = 1
			}
			if len(calls) != expectedCalls {
				t.Fatalf("Modify calls=%d, want %d", len(calls), expectedCalls)
			}
			if mixed {
				if calls[0].Get("TunnelOptionsSpecification.1.TunnelIpsecConfig.IpsecLifetime") != "12345" {
					t.Fatal("writable change lost with ignored role")
				}
				for key := range calls[0] {
					if strings.HasSuffix(key, ".Role") || strings.Contains(key, ".TunnelIkeConfig.") {
						t.Fatal("read-only role or unchanged IKE included in update")
					}
				}
			}
			state, refreshDiags := r.RefreshWithoutUpgrade(context.Background(), state, client)
			if refreshDiags.HasError() {
				t.Fatalf("%v", refreshDiags)
			}
			if vpnAttachmentStateAttrByIndex(state, 1, "role") != "master" {
				t.Fatal("configured role overwrote service readback")
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
		})
	}
}

func TestUnitVpnGatewayVpnAttachmentSetUnknownResolvedUpdate(t *testing.T) {
	for _, tc := range []struct {
		field, wire string
		value       interface{}
	}{
		{"enable_dpd", "EnableDpd", false},
		{"customer_gateway_id", "CustomerGatewayId", "cgw-unit-updated"},
		{"psk", "TunnelIkeConfig.Psk", "resolved-test-psk"},
		{"local_asn", "TunnelBgpConfig.LocalAsn", 0},
	} {
		t.Run(tc.field, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
			old.SetId("vco-unit-test")
			config := vpnAttachmentDistinctConfig()
			target := config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})
			if tc.field == "psk" {
				target = target["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})
			}
			if tc.field == "local_asn" {
				target = target["tunnel_bgp_config"].([]interface{})[0].(map[string]interface{})
			}
			target[tc.field] = "74D93920-ED26-11E3-AC10-0800200C9A66"
			unknown, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for key, change := range unknown.Attributes {
				if strings.HasPrefix(key, "tunnel_options_specification.~") && strings.HasSuffix(key, "."+tc.field) && change.NewComputed {
					found = true
				}
			}
			if !found {
				t.Fatal("target unknown was overwritten")
			}
			target[tc.field] = tc.value
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			calls := requests.modifyCalls()
			if len(calls) != 1 {
				t.Fatalf("modify calls=%d", len(calls))
			}
			if calls[0].Get("TunnelOptionsSpecification.1.TunnelIndex") != "1" || calls[0].Get("TunnelOptionsSpecification.1."+tc.wire) != fmt.Sprint(tc.value) {
				t.Fatal("resolved write was not applied")
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, config, client)
		})
	}
}

func vpnAttachmentStateAttrByIndex(state *terraform.InstanceState, index int, field string) string {
	for key, value := range state.Attributes {
		if strings.HasPrefix(key, "tunnel_options_specification.") && strings.HasSuffix(key, ".tunnel_index") && value == fmt.Sprint(index) {
			return state.Attributes[strings.TrimSuffix(key, ".tunnel_index")+"."+field]
		}
	}
	return ""
}

func TestUnitVpnGatewayVpnAttachmentSetExplicitEmptyPsk(t *testing.T) {
	r := resourceAliCloudVpnGatewayVpnAttachment()
	old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
	old.SetId("vco-unit-test")
	config := vpnAttachmentDistinctConfig()
	config["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["psk"] = ""
	diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(config))
	if _, diags := r.Apply(context.Background(), old.State(), diff, client); !diags.HasError() || !strings.Contains(fmt.Sprintf("%v", diags), "pre-shared key") {
		t.Fatal("explicit empty PSK must reach the safety guard, not inherit the prior key")
	}
	if len(requests.modifyCalls()) != 0 {
		t.Fatal("explicit empty PSK reached API")
	}
}

func TestUnitVpnGatewayVpnAttachmentSetUnknownIdentityAndBlockResolved(t *testing.T) {
	for _, field := range []string{"tunnel_index", "tunnel_ike_config"} {
		t.Run(field, func(t *testing.T) {
			r := resourceAliCloudVpnGatewayVpnAttachment()
			old := schema.TestResourceDataRaw(t, r.Schema, vpnAttachmentDistinctConfig())
			old.SetId("vco-unit-test")
			pending := vpnAttachmentDistinctConfig()
			pending["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})[field] = "74D93920-ED26-11E3-AC10-0800200C9A66"
			if _, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(pending), nil); err != nil {
				t.Fatal(err)
			}
			resolved := vpnAttachmentDistinctConfig()
			resolved["tunnel_options_specification"].([]interface{})[0].(map[string]interface{})["tunnel_ike_config"].([]interface{})[0].(map[string]interface{})["local_id"] = "resolved-local"
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(resolved), nil)
			if err != nil {
				t.Fatal(err)
			}
			client, requests := vpnAttachmentUnitClient(t, vpnAttachmentUnitAPI(resolved))
			state, diags := r.Apply(context.Background(), old.State(), diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			calls := requests.modifyCalls()
			if len(calls) != 1 || calls[0].Get("TunnelOptionsSpecification.1.TunnelIndex") != "1" || calls[0].Get("TunnelOptionsSpecification.1.TunnelIkeConfig.LocalId") != "resolved-local" {
				t.Fatal("resolved tunnel configuration was not applied")
			}
			vpnAttachmentUnitAssertNoDiff(t, r, state, resolved, client)
		})
	}
}

// A complete IKE write equal to the previous list occupant's configuration
// must still be applied to the tunnel now occupying that position.
func TestAccAliCloudVPNGatewayVpnAttachment_reorderExplicitIkeBlock(t *testing.T) {
	var v map[string]interface{}
	var originalID string
	const resourceID = "alicloud_vpn_gateway_vpn_attachment.default"
	rc := resourceCheckInitWithDescribeMethod(resourceID, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVpnAttachment")
	name := fmt.Sprintf("tfaccvpnreorder%d", acctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		IDRefreshName:     resourceID,
		CheckDestroy: func(state *terraform.State) error {
			if err := rc.checkResourceDestroy()(state); err != nil {
				return err
			}
			t.Log("SAFE_ATTACHMENT_DESTROY_VERIFIED")
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config:             vpnAttachmentReorderExplicitIkeConfig(name, false),
				ExpectNonEmptyPlan: false,
				Check:              vpnAttachmentCheckExplicitIke(resourceID, &originalID, false),
			},
			{
				Config:             vpnAttachmentReorderExplicitIkeConfig(name, true),
				ExpectNonEmptyPlan: false,
				Check:              vpnAttachmentCheckExplicitIke(resourceID, &originalID, true),
			},
		},
	})
}

func vpnAttachmentReorderExplicitIkeConfig(name string, reorder bool) string {
	first := vpnAttachmentReorderExplicitIkeTunnel(1, 1)
	second := vpnAttachmentReorderExplicitIkeTunnel(2, 2)
	if reorder {
		// All IKE fields at position zero equal the previous occupant's
		// values, while tunnel 2 must receive tunnel 1's complete IKE block.
		first = vpnAttachmentReorderExplicitIkeTunnel(2, 1)
		second = vpnAttachmentReorderExplicitIkeTunnel(1, 1)
	}
	return AlicloudVpnGatewayVpnAttachmentBasicDependenceRole(name) + fmt.Sprintf(`
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  vpn_attachment_name = var.name
  network_type = "public"
  local_subnet = "0.0.0.0/0"
  remote_subnet = "0.0.0.0/0"
  enable_tunnels_bgp = false
  effect_immediately = false
%s
%s
}
`, first, second)
}

func vpnAttachmentReorderExplicitIkeTunnel(index, ikeIdentity int) string {
	return fmt.Sprintf(`
  tunnel_options_specification {
    tunnel_index = %d
    customer_gateway_id = alicloud_vpn_customer_gateway.cgw%d.id
    role = "master"
    enable_dpd = true
    enable_nat_traversal = true
    tunnel_ike_config {
      ike_auth_alg = "sha1"
      ike_enc_alg = "aes"
      ike_version = "ikev2"
      ike_mode = "main"
      ike_lifetime = 86400
      ike_pfs = "group2"
      psk = "tf-acc-reorder-test-%d"
      local_id = "198.51.100.%d"
      remote_id = "203.0.113.%d"
    }
  }
`, index, index, ikeIdentity, 10+ikeIdentity, 10+ikeIdentity)
}

func TestAccAliCloudVpnGatewayVpnAttachment_partialIkeUpdate(t *testing.T) {
	var v map[string]interface{}
	const resourceID = "alicloud_vpn_gateway_vpn_attachment.default"
	rc := resourceCheckInitWithDescribeMethod(resourceID, &v, func() interface{} { return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)} }, "DescribeVpnGatewayVpnAttachment")
	name := fmt.Sprintf("tfaccvpnpartial%d", acctest.RandIntRange(10000, 99999))
	config := func(step int) string {
		blocks := ""
		for _, index := range []int{2, 1} {
			localID := fmt.Sprintf("198.51.100.%d", 10+index)
			psk := fmt.Sprintf("tf-acc-partial-test-%d", index)
			if step > 0 && index == 2 {
				localID = "198.51.100.22"
				psk = "tf-acc-partial-updated-2"
			}
			pskLine := fmt.Sprintf("psk=%q", psk)
			if step == 2 && index == 2 {
				pskLine = ""
				localID = "198.51.100.23"
			}
			blocks += fmt.Sprintf(`
 tunnel_options_specification {
  tunnel_index=%d
  customer_gateway_id=alicloud_vpn_customer_gateway.cgw%d.id
  tunnel_ike_config {
   %s
   local_id="%s"
  }
 }
`, index, index, pskLine, localID)
		}
		return AlicloudVpnGatewayVpnAttachmentBasicDependenceRole(name) + fmt.Sprintf(`
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
 vpn_attachment_name=var.name
 network_type="public"
 local_subnet="0.0.0.0/0"
 remote_subnet="0.0.0.0/0"
 enable_tunnels_bgp=false
 effect_immediately=false
%s
}
`, blocks)
	}
	var originalID string
	baseline := map[int]map[string]string{}
	checkPhase := func(phase int) resource.TestCheckFunc {
		return vpnAttachmentCheckConfiguration(resourceID, &originalID, func(state *terraform.State) map[int]map[string]string {
			result := map[int]map[string]string{}
			for _, index := range []int{1, 2} {
				if phase == 0 {
					baseline[index] = map[string]string{}
					for _, field := range []string{"role", "enable_dpd", "enable_nat_traversal", "tunnel_ike_config.0.ike_auth_alg", "tunnel_ike_config.0.ike_enc_alg", "tunnel_ike_config.0.ike_version", "tunnel_ike_config.0.ike_mode", "tunnel_ike_config.0.ike_lifetime", "tunnel_ike_config.0.ike_pfs", "tunnel_ike_config.0.remote_id"} {
						baseline[index][field] = vpnAttachmentStateAttrByIndex(state.RootModule().Resources[resourceID].Primary, index, field)
					}
				}
				m := map[string]string{}
				for k, v := range baseline[index] {
					m[k] = v
				}
				m["customer_gateway_id"] = state.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", index)].Primary.ID
				m["tunnel_ike_config.0.psk"] = fmt.Sprintf("tf-acc-partial-test-%d", index)
				m["tunnel_ike_config.0.local_id"] = fmt.Sprintf("198.51.100.%d", 10+index)
				if phase > 0 && index == 2 {
					m["tunnel_ike_config.0.psk"] = "tf-acc-partial-updated-2"
					m["tunnel_ike_config.0.local_id"] = fmt.Sprintf("198.51.100.%d", 21+phase)
				}
				result[index] = m
			}
			return result
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory, IDRefreshName: resourceID,
		CheckDestroy: func(state *terraform.State) error {
			if err := rc.checkResourceDestroy()(state); err != nil {
				return err
			}
			t.Log("SAFE_ATTACHMENT_DESTROY_VERIFIED")
			return nil
		},
		Steps: []resource.TestStep{
			{Config: config(0), ExpectNonEmptyPlan: false, Check: checkPhase(0)},
			{Config: config(1), ExpectNonEmptyPlan: false, Check: checkPhase(1)},
			{Config: config(2), ExpectNonEmptyPlan: false, Check: checkPhase(2)},
			{Config: config(2), PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

// Check independently by the API index and the state's index, never Set position.
func vpnAttachmentCheckConfiguration(resourceID string, originalID *string, expected func(*terraform.State) map[int]map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		mismatch := fmt.Errorf("SAFE_TUNNEL_CONFIGURATION_MISMATCH")
		entry := state.RootModule().Resources[resourceID]
		if entry == nil || entry.Primary == nil {
			return mismatch
		}
		if *originalID == "" {
			*originalID = entry.Primary.ID
		}
		if *originalID != entry.Primary.ID {
			return mismatch
		}
		service := VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		obj, err := service.DescribeVpnGatewayVpnAttachment(entry.Primary.ID)
		if err != nil {
			return fmt.Errorf("SAFE_VPN_LIVE_READ_FAILED")
		}
		outer, ok := obj["TunnelOptionsSpecification"].(map[string]interface{})
		if !ok {
			return mismatch
		}
		tunnels, ok := outer["TunnelOptions"].([]interface{})
		if !ok || len(tunnels) != 2 {
			return mismatch
		}
		want := expected(state)
		seen := map[int]bool{}
		fields := map[string]string{
			"customer_gateway_id": "CustomerGatewayId", "role": "Role", "enable_dpd": "EnableDpd", "enable_nat_traversal": "EnableNatTraversal",
			"tunnel_ike_config.0.ike_auth_alg": "TunnelIkeConfig.IkeAuthAlg", "tunnel_ike_config.0.ike_enc_alg": "TunnelIkeConfig.IkeEncAlg", "tunnel_ike_config.0.ike_version": "TunnelIkeConfig.IkeVersion", "tunnel_ike_config.0.ike_mode": "TunnelIkeConfig.IkeMode", "tunnel_ike_config.0.ike_lifetime": "TunnelIkeConfig.IkeLifetime", "tunnel_ike_config.0.ike_pfs": "TunnelIkeConfig.IkePfs", "tunnel_ike_config.0.psk": "TunnelIkeConfig.Psk", "tunnel_ike_config.0.local_id": "TunnelIkeConfig.LocalId", "tunnel_ike_config.0.remote_id": "TunnelIkeConfig.RemoteId",
			"tunnel_ipsec_config.0.ipsec_auth_alg": "TunnelIpsecConfig.IpsecAuthAlg", "tunnel_ipsec_config.0.ipsec_enc_alg": "TunnelIpsecConfig.IpsecEncAlg", "tunnel_ipsec_config.0.ipsec_lifetime": "TunnelIpsecConfig.IpsecLifetime", "tunnel_ipsec_config.0.ipsec_pfs": "TunnelIpsecConfig.IpsecPfs",
		}
		for _, raw := range tunnels {
			tunnel, ok := raw.(map[string]interface{})
			if !ok {
				return mismatch
			}
			index := formatInt(tunnel["TunnelIndex"])
			if seen[index] || want[index] == nil {
				return mismatch
			}
			seen[index] = true
			for field, value := range want[index] {
				if vpnAttachmentStateAttrByIndex(entry.Primary, index, field) != value {
					return mismatch
				}
				var actual interface{} = tunnel
				for _, key := range strings.Split(fields[field], ".") {
					m, ok := actual.(map[string]interface{})
					if !ok {
						return mismatch
					}
					actual = m[key]
				}
				// A missing PSK is deliberately not disclosed by Describe; state must still retain the exact configured secret.
				if field == "tunnel_ike_config.0.psk" && (actual == nil || actual == "") {
					continue
				}
				// Read maps an absent optional remote ID to the schema string zero value.
				if field == "tunnel_ike_config.0.remote_id" && actual == nil && value == "" {
					continue
				}
				if fmt.Sprint(actual) != value {
					return mismatch
				}
			}
		}
		return nil
	}
}

func vpnAttachmentCheckRolePhase(resourceID string, originalID *string, phase int) resource.TestCheckFunc {
	return vpnAttachmentCheckConfiguration(resourceID, originalID, func(state *terraform.State) map[int]map[string]string {
		result := map[int]map[string]string{}
		for _, index := range []int{1, 2} {
			cgw := 1
			flag := "true"
			auth, enc, version, mode, lifetime, pfs := "md5", "aes", "ikev2", "main", "86400", "group2"
			suffix := fmt.Sprint(index)
			if phase == 2 {
				cgw = 2
				flag = "false"
				auth, enc, version, mode, lifetime, pfs = "sha1", "aes256", "ikev1", "aggressive", "43200", "group14"
				suffix = "upd-" + suffix
			}
			m := map[string]string{"customer_gateway_id": state.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", cgw)].Primary.ID, "role": "master", "enable_dpd": flag, "enable_nat_traversal": flag}
			ike := map[string]string{"ike_auth_alg": auth, "ike_enc_alg": enc, "ike_version": version, "ike_mode": mode, "ike_lifetime": lifetime, "ike_pfs": pfs, "psk": "tf-testvpn-role-" + suffix, "local_id": "role-local-" + suffix, "remote_id": "role-remote-" + suffix}
			for k, v := range ike {
				m["tunnel_ike_config.0."+k] = v
			}
			ipsec := map[string]string{"ipsec_auth_alg": "md5", "ipsec_enc_alg": "aes", "ipsec_lifetime": "86400", "ipsec_pfs": "group5"}
			if phase == 2 || (phase == 1 && index == 1) {
				ipsec = map[string]string{"ipsec_auth_alg": "sha1", "ipsec_enc_alg": "aes256", "ipsec_lifetime": "43200", "ipsec_pfs": "group14"}
			}
			for k, v := range ipsec {
				m["tunnel_ipsec_config.0."+k] = v
			}
			result[index] = m
		}
		return result
	})
}

func vpnAttachmentCheckExplicitIke(resourceID string, originalID *string, reordered bool) resource.TestCheckFunc {
	return vpnAttachmentCheckConfiguration(resourceID, originalID, func(state *terraform.State) map[int]map[string]string {
		result := map[int]map[string]string{}
		for _, index := range []int{1, 2} {
			identity := index
			if reordered {
				identity = 1
			}
			m := map[string]string{"customer_gateway_id": state.RootModule().Resources[fmt.Sprintf("alicloud_vpn_customer_gateway.cgw%d", index)].Primary.ID}
			for k, v := range map[string]string{"ike_auth_alg": "sha1", "ike_enc_alg": "aes", "ike_version": "ikev2", "ike_mode": "main", "ike_lifetime": "86400", "ike_pfs": "group2", "psk": fmt.Sprintf("tf-acc-reorder-test-%d", identity), "local_id": fmt.Sprintf("198.51.100.%d", 10+identity), "remote_id": fmt.Sprintf("203.0.113.%d", 10+identity)} {
				m["tunnel_ike_config.0."+k] = v
			}
			result[index] = m
		}
		return result
	})
}
