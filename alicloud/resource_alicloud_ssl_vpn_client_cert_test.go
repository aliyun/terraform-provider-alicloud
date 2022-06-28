package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ssl_vpn_client_cert", &resource.Sweeper{
		Name: "alicloud_ssl_vpn_client_cert",
		F:    testSweepSslVpnClientCerts,
	})
}

func testSweepSslVpnClientCerts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var certs []vpc.SslVpnClientCertKey
	request := vpc.CreateDescribeSslVpnClientCertsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnClientCerts(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSslVpnClientCertsResponse)
		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < 1 {
			break
		}
		certs = append(certs, response.SslVpnClientCertKeys.SslVpnClientCertKey...)

		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
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
			log.Printf("[ERROR] Failed to delete Ssl Client Cert (%s (%s)): %s", name, id, WrapError(err))
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	if sweeped {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func TestAccAlicloudSslVpnClientCert_basic(t *testing.T) {
	var v vpc.DescribeSslVpnClientCertResponse

	resourceId := "alicloud_ssl_vpn_client_cert.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnClientCertConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              fmt.Sprintf("tf-testAccSslVpnClientCertConfig%d", rand),
						"ssl_vpn_server_id": CHECKSET,
						"ca_cert":           CHECKSET,
						"client_cert":       CHECKSET,
						"client_key":        CHECKSET,
						"client_config":     CHECKSET,
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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_client_cert" {
			continue
		}

		_, err := vpnGatewayService.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccSslVpnClientCertConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PrePaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	local_subnet = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = "1194"
	compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
	ssl_vpn_server_id = "${alicloud_ssl_vpn_server.default.id}"
	name = "${var.name}"
}
`, rand)
}

func testAccSslVpnClientCertConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PrePaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	local_subnet = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = "1194"
	compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
	ssl_vpn_server_id = "${alicloud_ssl_vpn_server.default.id}"
	name = "${var.name}_change"
}
`, rand)
}

func testAccSslVpnClientCertConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PrePaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	local_subnet = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = "1194"
	compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
	ssl_vpn_server_id = "${alicloud_ssl_vpn_server.default.id}"
	name = "${var.name}"
}
`, rand)
}
