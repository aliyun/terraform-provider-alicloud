package vpngateway_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpngateway"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ssl_vpn_client_cert", &resource.Sweeper{
		Name: "alicloud_ssl_vpn_client_cert",
		F:    testSweepSslVpnClientCerts,
	})
}

func testSweepSslVpnClientCerts(region string) error {
	rawClient, err := vpngatewayTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var certs []vpc.SslVpnClientCertKey
	request := vpc.CreateDescribeSslVpnClientCertsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnClientCerts(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s", errs.WrapError(err))
			return errs.WrapError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSslVpnClientCertsResponse)
		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < 1 {
			break
		}
		certs = append(certs, response.SslVpnClientCertKeys.SslVpnClientCertKey...)

		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < helper.PageSizeLarge {
			break
		}

		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			log.Printf("[ERROR] %s", errs.WrapError(err))
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range certs {
		name := v.Name
		id := v.SslVpnClientCertId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ssl Client Cert: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ssl Client Cert: %s (%s)", name, id)
		request := vpc.CreateDeleteSslVpnClientCertRequest()
		request.SslVpnClientCertId = id
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnClientCert(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ssl Client Cert (%s (%s)): %s", name, id, errs.WrapError(err))
			return errs.WrapError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	if sweeped {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func TestAccAliCloudSslVpnClientCert_basic(t *testing.T) {
	var v vpc.DescribeSslVpnClientCertResponse

	resourceId := "alicloud_ssl_vpn_client_cert.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return vpngateway.NewVpnGatewayService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithTime(t, []int{1})
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnClientCertConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              fmt.Sprintf("tf-testAccSslVpnClientCertConfig%d", rand),
						"ssl_vpn_server_id": acctest.CHECKSET,
						"ca_cert":           acctest.CHECKSET,
						"client_cert":       acctest.CHECKSET,
						"client_key":        acctest.CHECKSET,
						"client_config":     acctest.CHECKSET,
						"status":            "normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSslVpnClientCertConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSslVpnClientCertConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccSslVpnClientCertConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSslVpnClientCertConfig%d", rand),
					}),
				),
			},
		},
	})

}

func testAccCheckSslVpnClientCertDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := vpngateway.NewVpnGatewayService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_client_cert" {
			continue
		}

		_, err := vpnGatewayService.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
		}
	}

	return nil
}

func testAccSslVpnClientCertConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_vswitches" "default0" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}


resource "alicloud_vpn_gateway" "default" {
  vpn_gateway_name             = var.name
  vpc_id                       = data.alicloud_vpcs.default.ids.0
  bandwidth                    = "10"
  enable_ssl                   = true
  description                  = var.name
  payment_type                 = "Subscription"
  vswitch_id                   = data.alicloud_vswitches.default0.ids.0
  disaster_recovery_vswitch_id = data.alicloud_vswitches.default1.ids.0
}


resource "alicloud_ssl_vpn_server" "default" {
  name           = var.name
  vpn_gateway_id = alicloud_vpn_gateway.default.id
  client_ip_pool = "192.168.0.0/16"
  local_subnet   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 8)
  protocol       = "UDP"
  cipher         = "AES-128-CBC"
  port           = "1194"
  compress       = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
  ssl_vpn_server_id = alicloud_ssl_vpn_server.default.id
  name              = var.name
}
`, rand)
}

func testAccSslVpnClientCertConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_vswitches" "default0" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}


resource "alicloud_vpn_gateway" "default" {
  vpn_gateway_name             = var.name
  vpc_id                       = data.alicloud_vpcs.default.ids.0
  bandwidth                    = "10"
  enable_ssl                   = true
  description                  = var.name
  payment_type                 = "Subscription"
  vswitch_id                   = data.alicloud_vswitches.default0.ids.0
  disaster_recovery_vswitch_id = data.alicloud_vswitches.default1.ids.0
}

resource "alicloud_ssl_vpn_server" "default" {
  name           = var.name
  vpn_gateway_id = alicloud_vpn_gateway.default.id
  client_ip_pool = "192.168.0.0/16"
  local_subnet   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 8)
  protocol       = "UDP"
  cipher         = "AES-128-CBC"
  port           = "1194"
  compress       = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
  ssl_vpn_server_id = alicloud_ssl_vpn_server.default.id
  name              = "${var.name}_change"
}
`, rand)
}

func testAccSslVpnClientCertConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_vswitches" "default0" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}


resource "alicloud_vpn_gateway" "default" {
  vpn_gateway_name             = var.name
  vpc_id                       = data.alicloud_vpcs.default.ids.0
  bandwidth                    = "10"
  enable_ssl                   = true
  description                  = var.name
  payment_type                 = "Subscription"
  vswitch_id                   = data.alicloud_vswitches.default0.ids.0
  disaster_recovery_vswitch_id = data.alicloud_vswitches.default1.ids.0
}


resource "alicloud_ssl_vpn_server" "default" {
  name           = var.name
  vpn_gateway_id = alicloud_vpn_gateway.default.id
  client_ip_pool = "192.168.0.0/16"
  local_subnet   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 8)
  protocol       = "UDP"
  cipher         = "AES-128-CBC"
  port           = "1194"
  compress       = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
  ssl_vpn_server_id = alicloud_ssl_vpn_server.default.id
  name              = var.name
}
`, rand)
}
