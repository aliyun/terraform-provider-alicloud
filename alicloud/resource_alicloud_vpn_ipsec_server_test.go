package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_vpn_ipsec_server",
		&resource.Sweeper{
			Name: "alicloud_vpn_ipsec_server",
			F:    testSweepVpnIpsecServer,
		})
}

func testSweepVpnIpsecServer(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListIpsecServers"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = aliyunClient.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.IpsecServers", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.XXXXX", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["IpsecServerName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpn Ipsec Server: %s", item["IpsecServerName"].(string))
				continue
			}
			action := "DeleteIpsecServer"
			request := map[string]interface{}{
				"IpsecServerId": item["IpsecServerId"],
				"RegionId":      aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpn Ipsec Server (%s): %s", item["IpsecServerName"].(string), err)
			}
			log.Printf("[INFO] Delete Vpn Ipsec Server success: %s ", item["IpsecServerName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudVPNGatewayIpsecServer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_ipsec_server.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnIpsecServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnIpsecServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpnipsecserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnIpsecServerBasicDependence0)
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
					"local_subnet":   "192.168.0.0/24",
					"client_ip_pool": "10.0.0.0/24",
					"vpn_gateway_id": "${local.vpn_gateway_id}",
					"psk_enabled":    "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":   "192.168.0.0/24",
						"client_ip_pool": "10.0.0.0/24",
						"vpn_gateway_id": CHECKSET,
						"psk_enabled":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_server_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipsec_server_name": name,
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
					"client_ip_pool": "10.0.1.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ip_pool": "10.0.1.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"psk": "tf-testpask",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"psk": "tf-testpask",
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
					"ike_config": []map[string]interface{}{
						{
							"ike_version":  "ikev2",
							"ike_mode":     "main",
							"ike_enc_alg":  "aes",
							"ike_auth_alg": "sha1",
							"ike_pfs":      "group2",
							"ike_lifetime": "86400",
							"local_id":     "127.0.0.2",
							"remote_id":    "127.0.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ike_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_pfs":      "group2",
							"ipsec_lifetime": "66400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipsec_config.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudVPNGatewayIpsecServer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_ipsec_server.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnIpsecServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnIpsecServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpnipsecserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnIpsecServerBasicDependence0)
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
					"local_subnet":       "192.168.0.0/24",
					"client_ip_pool":     "10.0.0.0/24",
					"ipsec_server_name":  "${var.name}",
					"effect_immediately": "true",
					"psk_enabled":        "true",
					"psk":                "tf-testpask",
					"vpn_gateway_id":     "${local.vpn_gateway_id}",
					"ike_config": []map[string]interface{}{
						{
							"ike_version":  "ikev2",
							"ike_mode":     "main",
							"ike_enc_alg":  "aes",
							"ike_auth_alg": "sha1",
							"ike_pfs":      "group2",
							"ike_lifetime": "86400",
							"local_id":     "127.0.0.2",
							"remote_id":    "127.0.0.1",
						},
					},
					"ipsec_config": []map[string]interface{}{
						{
							"ipsec_enc_alg":  "aes",
							"ipsec_auth_alg": "sha1",
							"ipsec_pfs":      "group2",
							"ipsec_lifetime": "86400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet":       "192.168.0.0/24",
						"client_ip_pool":     "10.0.0.0/24",
						"ipsec_server_name":  name,
						"effect_immediately": "true",
						"psk_enabled":        "true",
						"psk":                "tf-testpask",
						"vpn_gateway_id":     CHECKSET,
						"ike_config.#":       "1",
						"ipsec_config.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVpnIpsecServerMap0 = map[string]string{
	"dry_run": NOSET,
}

func AlicloudVpnIpsecServerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}
locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

data "alicloud_vpn_gateways" "default" {
  vpc_id       = data.alicloud_vpcs.default.ids.0
  ssl_vpn      = "enable"
}

locals {
  vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
}
`, name)
}

// lintignore: R001
