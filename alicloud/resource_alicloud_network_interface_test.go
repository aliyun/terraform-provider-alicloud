package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"

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
	service := VpcService{client}
	for _, eni := range enis {
		name := eni.NetworkInterfaceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a nat gateway name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(eni.VpcId, ""); err == nil {
				skip = !need
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

			if err := ecsService.WaitForNetworkInterface(eni.NetworkInterfaceId, Available, DefaultTimeout); err != nil {
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
		_, err := ecsService.DescribeNetworkInterface(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func TestAccAlicloudNetworkInterfaceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	var v ecs.NetworkInterfaceSet
	resourceId := "alicloud_network_interface.default"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNetworkInterface%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkInterfaceConfig_privateIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_private_ips(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips.#":     "3",
						"private_ips_count": "3",
					}),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_private_ips_count(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// There is a bug when d.Set a set parameter. The new values can not overwrite the state
						// when a parameter is a TypeSet and Computed. https://github.com/hashicorp/terraform/issues/20504
						// "private_ips.#": "4",
						"private_ips_count": "4",
					}),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNetworkInterfaceChange%d", rand),
						// Same with last step
						"private_ips.#": "4",
					}),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-eni-description",
					}),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_tags(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "1",
					}),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// There is a bug when d.Set a set parameter. The new values can not overwrite the state
						// when a parameter is a TypeSet and Computed. https://github.com/hashicorp/terraform/issues/20504
						// "private_ips.#":     "1",
						"private_ips_count": "1",
						"description":       "tf-testAcc-eni-description_all",
						"tags.%":            "0",
						"name":              fmt.Sprintf("tf-testAccNetworkInterface%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceMulti(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alicloud_network_interface.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNetworkInterface%d", rand),
					}),
				),
			},
		},
	})
}

func testAccNetworkInterfaceConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}
`, rand)
}

func testAccNetworkInterfaceConfig_privateIp(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNetworkInterface"
}
resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
    private_ip = "192.168.0.2"
}
`, rand)
}

func testAccNetworkInterfaceConfig_private_ips(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips = ["192.168.0.3", "192.168.0.5", "192.168.0.6"]
}
`, rand)
}

func testAccNetworkInterfaceConfig_private_ips_count(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips_count = 4
}
`, rand)
}

func testAccNetworkInterfaceConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips_count = 4
    name = "${var.name}Change%d"
}
`, rand)
}

func testAccNetworkInterfaceConfig_description(rand int) string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips_count = 4
    name = "${var.name}Change%d"
    description = "tf-testAcc-eni-description"
}
`, rand)
}

func testAccNetworkInterfaceConfig_tags(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips_count = 4
    name = "${var.name}Change%d"
    description = "tf-testAcc-eni-description"
    tags = {
		TF-VER = "0.11.3"
	}
}
`, rand)
}

func testAccNetworkInterfaceConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips_count = 1
    name = "${var.name}%d"
    description = "tf-testAcc-eni-description_all"
}
`, rand)
}

func testAccNetworkInterfaceConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}%d"
    count = 3
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}
`, rand)
}

var testAccCheckNetworkInterfaceCheckMap = map[string]string{
	"vswitch_id":        CHECKSET,
	"security_groups.#": "1",
	"private_ip":        CHECKSET,
	"private_ips.#":     "0",
	"private_ips_count": "0",
	"description":       "",
	"tags.%":            NOSET,
}
