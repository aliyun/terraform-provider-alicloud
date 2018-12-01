package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_network_interface", &resource.Sweeper{
		Name: "alicloud_network_interface",
		F:    testAlicloudNetworkInterface,
	})
}

func testAlicloudNetworkInterface(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %#v", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	req := ecs.CreateDescribeNetworkInterfacesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	var enis []ecs.NetworkInterfaceSet
	for {
		raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
			return client.DescribeNetworkInterfaces(req)
		})
		if err != nil {
			return fmt.Errorf("Describe NetworkInterfaces failed, %#v", err)
		}

		resp := raw.(*ecs.DescribeNetworkInterfacesResponse)
		if resp == nil || len(resp.NetworkInterfaceSets.NetworkInterfaceSet) == 0 {
			break
		}

		enis = append(enis, resp.NetworkInterfaceSets.NetworkInterfaceSet...)

		if len(resp.NetworkInterfaceSets.NetworkInterfaceSet) < PageSizeLarge {
			break
		}

		if pageNumber, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = pageNumber
		}
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	sweeped := false
	for _, eni := range enis {
		name := eni.NetworkInterfaceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping NetworkInterface %s", name)
			continue
		}
		sweeped = true
		if eni.InstanceId != "" {
			req := ecs.CreateDetachNetworkInterfaceRequest()
			req.InstanceId = eni.InstanceId
			req.NetworkInterfaceId = eni.NetworkInterfaceId
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DetachNetworkInterface(req)
			})

			if err != nil {
				log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
				continue
			}

			if err := ecsService.WaitForEcsNetworkInterface(eni.NetworkInterfaceId, Available, DefaultTimeout); err != nil {
				log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
				continue
			}
		}

		log.Printf("[INFO] Deleting NetworkInterface %s", name)
		req := ecs.CreateDeleteNetworkInterfaceRequest()
		req.NetworkInterfaceId = eni.NetworkInterfaceId
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteNetworkInterface(req)
		})

		if err != nil {
			log.Printf("[ERROR] Delete NetworkInterface failed, %#v", err)
			continue
		}
	}

	if sweeped {
		time.Sleep(30 * time.Second)
	}

	return nil
}

func TestAccAlicloudNetworkInterfaceBasic(t *testing.T) {
	var eni ecs.NetworkInterfaceSet
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface.eni",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "private_ip"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "0"),
					resource.TestCheckNoResourceAttr("alicloud_network_interface.eni", "private_ips_count"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Basic test"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-TAG", "0.11.3")),
			},
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "private_ip"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "0"),
					resource.TestCheckNoResourceAttr("alicloud_network_interface.eni", "private_ips_count"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Basic2 test")),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceWithPrivateIpList(t *testing.T) {
	var eni ecs.NetworkInterfaceSet
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface.eni",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigWithPrivateIpAddressList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "private_ip"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "2"),
					resource.TestCheckNoResourceAttr("alicloud_network_interface.eni", "private_ips_count"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Address list test"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigWithPrivateIpAddressList2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "private_ip"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "3"),
					resource.TestCheckNoResourceAttr("alicloud_network_interface.eni", "private_ips_count"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Address list test2"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceWithPrivateIpCount(t *testing.T) {
	var eni ecs.NetworkInterfaceSet
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface.eni",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigWithPrivateIpAddressCount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "5"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips_count", "5"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Address count test"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigWithPrivateIpAddressCount2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "3"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips_count", "3"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "Address count test2"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceWithoutPrimaryIpAddress(t *testing.T) {
	var eni ecs.NetworkInterfaceSet
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface.eni",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceConfigWithoutPrimaryIpAddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("alicloud_network_interface.eni", &eni),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "name", "tf-testAcc-eni"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface.eni", "private_ip"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "private_ips.#", "0"),
					resource.TestCheckNoResourceAttr("alicloud_network_interface.eni", "private_ips_count"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "description", "No primary private IP address test"),
					resource.TestCheckResourceAttr("alicloud_network_interface.eni", "tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

func testAccCheckEniExists(n string, eni *ecs.NetworkInterfaceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		d, err := ecsService.DescribeNetworkInterfaceById("", rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("While checking ENI existing, describing disk got an error: %#v.", err)
		}

		*eni = d

		return nil
	}
}

func testAccCheckNetworkInterfaceDestroy(t *terraform.State) error {
	for _, rs := range t.RootModule().Resources {
		if rs.Type != "alicloud_network_interface" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceById("", rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccNetworkInterfaceConfigBasic = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-TAG = "0.11.3"
	}
}
`

const testAccNetworkInterfaceConfigBasic2 = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Basic2 test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3"
	}
}
`

const testAccNetworkInterfaceConfigWithPrivateIpAddressList = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Address list test"
	private_ip = "192.168.0.2"
	private_ips = ["192.168.0.3", "192.168.0.4"]
	tags = {
		TF-VER = "0.11.3"
	}
}
`

const testAccNetworkInterfaceConfigWithPrivateIpAddressList2 = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Address list test2"
	private_ip = "192.168.0.2"
	private_ips = ["192.168.0.3", "192.168.0.5", "192.168.0.6"]
	tags = {
		TF-VER = "0.11.3"
	}
}
`

const testAccNetworkInterfaceConfigWithPrivateIpAddressCount = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Address count test"
	private_ip = "192.168.0.2"
	private_ips_count = 5
	tags = {
		TF-VER = "0.11.3"
	}
}
`
const testAccNetworkInterfaceConfigWithPrivateIpAddressCount2 = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Address count test2"
	private_ip = "192.168.0.2"
	private_ips_count = 3
	tags = {
		TF-VER = "0.11.3"
	}
}
`

const testAccNetworkInterfaceConfigWithoutPrimaryIpAddress = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "No primary private IP address test"
	tags = {
		TF-VER = "0.11.3"
	}
}
`
