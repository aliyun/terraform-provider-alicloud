package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	req := vpc.CreateDescribeSslVpnServersRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnServers(req)
		})
		if err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		}
		resp, _ := raw.(*vpc.DescribeSslVpnServersResponse)
		if resp == nil || len(resp.SslVpnServers.SslVpnServer) < 1 {
			break
		}
		servers = append(servers, resp.SslVpnServers.SslVpnServer...)

		if len(resp.SslVpnServers.SslVpnServer) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		} else {
			req.PageNumber = page
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
		req := vpc.CreateDeleteSslVpnServerRequest()
		req.SslVpnServerId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnServer(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ssl Vpn Server (%s (%s)): %s", name, id, WrapError(err))
		}
	}
	if sweeped {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func TestAccAlicloudSslVpnServer_basic(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},

		IDRefreshName: "alicloud_ssl_vpn_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.1.0/24,172.16.2.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
		},
	})

}

func TestAccAlicloudSslVpnServer_update_single_subnet(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfigSingleSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
			{
				Config: testAccSslVpnServerConfigUpdateSingleSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_update"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.2.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "none"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudSslVpnServer_update(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.1.0/24,172.16.2.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
			{
				Config: testAccSslVpnServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_update"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.0.0/24,172.16.1.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "none"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudSslVpnServer_update_singleToMulti_subnet(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfigSingleSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
			{
				Config: testAccSslVpnServerMultiConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "tf-testAccSslVpnServerConfig_update"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.0.0/24,172.16.1.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1195"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "none"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "true"),
				),
			},
		},
	})
}

func testAccCheckSslVpnServerExists(n string, vpn *vpc.SslVpnServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpnGatewayService := VpnGatewayService{client}
		instance, err := vpnGatewayService.DescribeSslVpnServer(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
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

const testAccSslVpnServerConfigSingleSubnet = `
variable "name" {
	default = "tf-testAccSslVpnServerConfig_create"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.1.0/24"
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = 1194
	compress = "false"
}
`

const testAccSslVpnServerConfig = `
variable "name" {
	default = "tf-testAccSslVpnServerConfig_create"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.1.0/24,172.16.2.0/24"
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = 1194
	compress = "false"
}
`

const testAccSslVpnServerConfigUpdate = `
variable "name" {
	default = "tf-testAccSslVpnServerConfig_update"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "${var.name}"
    vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.0.0/24,172.16.1.0/24"
    protocol = "UDP"
    cipher = "none"
    port = 1194
    compress = "false"
}
`

const testAccSslVpnServerMultiConfigUpdate = `
variable "name" {
	default = "tf-testAccSslVpnServerConfig_update"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "${var.name}"
    vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.0.0/24,172.16.1.0/24"
    protocol = "TCP"
    cipher = "none"
    port = 1195
    compress = "true"
}
`

const testAccSslVpnServerConfigUpdateSingleSubnet = `
variable "name" {
	default = "tf-testAccSslVpnServerConfig_update"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "${var.name}"
    vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.2.0/24"
    protocol = "UDP"
    cipher = "none"
    port = 1194
    compress = "false"
}
`
