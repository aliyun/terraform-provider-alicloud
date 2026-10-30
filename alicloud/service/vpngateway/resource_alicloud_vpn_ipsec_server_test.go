package vpngateway_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpngateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := vpngatewayTestClientForRegion(region)
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
	request["MaxResults"] = helper.PageSizeLarge

	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpnIpsecServerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpngateway.NewVpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpnIpsecServer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpnipsecserver%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpnIpsecServerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vpn_gateway_id": acctest.CHECKSET,
						"psk_enabled":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"psk_enabled":       "false",
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpnIpsecServerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpngateway.NewVpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpnIpsecServer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpnipsecserver%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpnIpsecServerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vpn_gateway_id":     acctest.CHECKSET,
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
	"dry_run": acctest.NOSET,
}

func AlicloudVpnIpsecServerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones.0.id
  name              = var.name
}

resource "alicloud_vswitch" "default1" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.17.0.0/22"
  availability_zone = data.alicloud_zones.default.zones.1.id
  name              = var.name
}

resource "alicloud_vpn_gateway" "default" {
  name                         = var.name
  instance_charge_type         = "PostPaid"
  vpc_id                       = alicloud_vpc.default.id
  vswitch_id                   = alicloud_vswitch.default.id
  disaster_recovery_vswitch_id = alicloud_vswitch.default1.id
  bandwidth                    = "10"
  description                  = "test"
}

locals {
  vpn_gateway_id = alicloud_vpn_gateway.default.id
}
`, name)
}

// lintignore: R001
