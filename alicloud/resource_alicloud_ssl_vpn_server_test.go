package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ssl_vpn_server", &resource.Sweeper{
		Name: "alicloud_ssl_vpn_server",
		F:    testSweepSslVpnServers,
	})
}

func testSweepSslVpnServers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var servers []vpc.SslVpnServer
	request := vpc.CreateDescribeSslVpnServersRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnServers(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSslVpnServersResponse)
		if len(response.SslVpnServers.SslVpnServer) < 1 {
			break
		}
		servers = append(servers, response.SslVpnServers.SslVpnServer...)

		if len(response.SslVpnServers.SslVpnServer) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range servers {
		name := v.Name
		id := v.SslVpnServerId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ssl Vpn Server: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ssl Vpn Server: %s (%s)", name, id)
		request := vpc.CreateDeleteSslVpnServerRequest()
		request.SslVpnServerId = id
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnServer(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ssl Vpn Server (%s (%s)): %s", name, id, WrapError(err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	if sweeped {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func testAccCheckSslVpnServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_server" {
			continue
		}

		instance, err := vpnGatewayService.DescribeSslVpnServer(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("SSL VPN server %s still exist", instance.SslVpnServerId)
	}

	return nil
}

func TestAccAlicloudSslVpnServerBasic(t *testing.T) {
	var v vpc.SslVpnServer

	resourceId := "alicloud_ssl_vpn_server.default"
	ra := resourceAttrInit(resourceId, testAccSslVpnServerCheckMap)
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
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSslVpnServerConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSslVpnServerConfig_client_ip_pool(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ip_pool": "192.168.0.0/17",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_local_subnet(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_subnet": "172.16.1.0/24",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSslVpnServerConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_protocol(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "TCP",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_cipher(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher": "AES-192-CBC",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_port(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "23",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfig_compress(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress": "false",
					}),
				),
			},
			{
				Config: testAccSslVpnServerConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            fmt.Sprintf("tf-testAccSslVpnServerConfig%d", rand),
						"vpn_gateway_id":  CHECKSET,
						"client_ip_pool":  "192.168.0.0/16",
						"local_subnet":    "172.16.1.0/24,172.16.2.0/24",
						"protocol":        "UDP",
						"port":            "1194",
						"cipher":          "AES-128-CBC",
						"compress":        "false",
						"connections":     "0",
						"max_connections": "5",
						"internet_ip":     CHECKSET,
					}),
				),
			},
		},
	})

}

func TestAccAlicloudSslVpnServerMulti(t *testing.T) {
	var v vpc.SslVpnServer

	resourceId := "alicloud_ssl_vpn_server.default.1"
	ra := resourceAttrInit(resourceId, testAccSslVpnServerCheckMap)
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
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfig_mulit(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            fmt.Sprintf("tf-testAccSslVpnServerConfig%d", rand),
						"vpn_gateway_id":  CHECKSET,
						"client_ip_pool":  "192.169.0.0/16",
						"local_subnet":    "172.16.1.0/24,172.16.2.0/24",
						"protocol":        "UDP",
						"port":            "1194",
						"cipher":          "AES-128-CBC",
						"compress":        "false",
						"connections":     "0",
						"max_connections": "5",
						"internet_ip":     CHECKSET,
					}),
				),
			},
		},
	})

}

func testAccSslVpnServerConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.1.0/24,172.16.2.0/24"
}
`, rand)
}

func testAccSslVpnServerConfig_client_ip_pool(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24,172.16.2.0/24"
}
`, rand)
}

func testAccSslVpnServerConfig_local_subnet(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
}
`, rand)
}

func testAccSslVpnServerConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}_change"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
}
`, rand)
}

func testAccSslVpnServerConfig_protocol(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}_change"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
	protocol = "TCP"
}
`, rand)
}

func testAccSslVpnServerConfig_cipher(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}_change"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
	protocol = "TCP"
	cipher = "AES-192-CBC"
}
`, rand)
}

func testAccSslVpnServerConfig_port(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}_change"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
	protocol = "TCP"
	cipher = "AES-192-CBC"
	port = "23"
}
`, rand)
}

func testAccSslVpnServerConfig_compress(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	name = "${var.name}_change"
	vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
	client_ip_pool = "192.168.0.0/17"
	local_subnet = "172.16.1.0/24"
	protocol = "TCP"
	cipher = "AES-192-CBC"
	port = "23"
	compress = "false"
}
`, rand)
}

func testAccSslVpnServerConfig_mulit(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServerConfig%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}
resource "alicloud_vpn_gateway" "default" {
	count = "2"
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
resource "alicloud_ssl_vpn_server" "default" {
	count = "2"
	name = "${var.name}"
	vpn_gateway_id = "${element(alicloud_vpn_gateway.default.*.id,count.index)}"
	client_ip_pool = "192.${count.index + 168}.0.0/16"
	local_subnet = "172.16.1.0/24,172.16.2.0/24"
}
`, rand)
}

var testAccSslVpnServerCheckMap = map[string]string{
	"vpn_gateway_id":  CHECKSET,
	"client_ip_pool":  "192.168.0.0/16",
	"local_subnet":    "172.16.1.0/24,172.16.2.0/24",
	"protocol":        "UDP",
	"port":            "1194",
	"cipher":          "AES-128-CBC",
	"compress":        "false",
	"connections":     "0",
	"max_connections": "5",
	"internet_ip":     CHECKSET,
}
