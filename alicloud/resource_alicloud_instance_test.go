package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudInstance_basic(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(*terraform.State) error {
		log.Printf("[WARN] instances: %#v", instance)
		if instance.ZoneId == "" {
			return fmt.Errorf("bad availability zone")
		}
		if len(instance.SecurityGroupIds.SecurityGroupId) == 0 {
			return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"testAccInstanceConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
					testAccCheckSystemDiskSize("alicloud_instance.foo", 80),
				),
			},

			// test for multi steps
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"testAccInstanceConfig"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_vpc(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfigVPC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"system_disk_category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_userData(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfigUserData,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"system_disk_category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"user_data",
						"echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_multipleRegions(t *testing.T) {
	var instance ecs.Instance

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckInstanceDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfigMultipleRegions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExistsWithProviders(
						"alicloud_instance.foo", &instance, &providers),
					testAccCheckInstanceExistsWithProviders(
						"alicloud_instance.bar", &instance, &providers),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_multiSecurityGroup(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(sgCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(instance.SecurityGroupIds.SecurityGroupId) < 0 {
				return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			if len(instance.SecurityGroupIds.SecurityGroupId) < sgCount {
				return fmt.Errorf("less security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig_multiSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(2),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"2"),
				),
			},
			resource.TestStep{
				Config: testAccInstanceConfig_multiSecurityGroup_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(3),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"3"),
				),
			},
			resource.TestStep{
				Config: testAccInstanceConfig_multiSecurityGroup_remove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(1),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"1"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_multiSecurityGroupByCount(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(sgCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(instance.SecurityGroupIds.SecurityGroupId) < 0 {
				return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			if len(instance.SecurityGroupIds.SecurityGroupId) < sgCount {
				return fmt.Errorf("less security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig_multiSecurityGroupByCount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(2),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"2"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_NetworkInstanceSecurityGroups(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceNetworkInstanceSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_tags(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"tags.foo",
						"bar"),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceConfigTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"tags.bar",
						"zzz"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceConfigOrigin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"testAccCheckInstanceConfigOrigin-foo"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"host_name",
						"host-foo"),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceConfigOriginUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"testAccCheckInstanceConfigOrigin-bar"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"host_name",
						"host-bar"),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceImage_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceImageOrigin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.update_image", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.update_image",
						"system_disk_size",
						"50"),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceImageUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.update_image", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.update_image",
						"system_disk_size",
						"60"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_associatePublicIP(t *testing.T) {
	var instance ecs.Instance

	testCheckPrivateIP := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			privateIP := instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
			if privateIP == "" {
				return fmt.Errorf("can't get private IP")
			}

			return nil
		}
	}

	testCheckPublicIP := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			publicIP := instance.PublicIpAddress.IpAddress[0]
			if publicIP == "" {
				return fmt.Errorf("can't get public IP")
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfigAssociatePublicIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					testCheckPrivateIP(),
					testCheckPublicIP(),
				),
			},
		},
	})
}

func TestAccAlicloudInstancePrivateIp_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstancePrivateIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"private_ip",
						"10.1.1.3"),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstancePrivateIpUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"private_ip",
						"10.1.2.3"),
				),
			},
		},
	})
}

// At present, One account only support at most 16 cpu core modify in one month.
func SkipTestAccAlicloudInstanceChargeType_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceChargeType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceChargeTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_spot(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.spot",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSpotInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.spot", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.spot",
						"spot_strategy", "SpotWithPriceLimit"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.spot",
						"spot_price_limit", "1.002"),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceType_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceNetworkSpec_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceNetworkSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.network", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_in", "0"),
				),
			},

			resource.TestStep{
				Config: testAccCheckInstanceNetworkSpecUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.network", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_out", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_in", "50"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_ramrole(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.role",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckInstanceRamRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.role", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.role",
						"role_name",
						"testAccCheckInstanceRamRole"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, i *ecs.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckInstanceExistsWithProviders(n, i, &providers)
}

func testAccCheckInstanceExistsWithProviders(n string, i *ecs.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*AliyunClient)
			instance, err := client.DescribeInstanceById(rs.Primary.ID)
			log.Printf("[WARN]get ecs instance %#v", instance)
			// Verify the error is what we want
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err

			}

			*i = instance
			return nil
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckInstanceDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance" {
			continue
		}

		// Try to find the resource
		instance, err := client.DescribeInstanceById(rs.Primary.ID)
		if err == nil {
			if instance.Status != "" && instance.Status != string(Stopped) {
				return fmt.Errorf("Found unstopped instance: %s", instance.InstanceId)
			}
		}

		// Verify the error is what we want
		if NotFoundError(err) {
			continue
		}

		return err
	}

	return nil
}

func testAccCheckSystemDiskSize(n string, size int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providers := []*schema.Provider{testAccProvider}
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		for _, provider := range providers {
			if provider.Meta() == nil {
				continue
			}
			client := provider.Meta().(*AliyunClient)
			systemDisk, err := client.QueryInstanceSystemDisk(rs.Primary.ID)
			if err != nil {
				log.Printf("[ERROR]get system disk size error: %#v", err)
				return err
			}

			if systemDisk.Size != size {
				return fmt.Errorf("system disk size not equal %d, the instance system size is %d",
					size, systemDisk.Size)
			}
		}

		return nil
	}
}

const testAccInstanceConfig = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_ssd"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "testAccInstanceConfig"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}"
	description = "bar"
}

resource "alicloud_instance" "foo" {
	image_id = "${data.alicloud_images.default.images.0.id}"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	system_disk_category = "cloud_ssd"
	system_disk_size = 80

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"

	tags {
		foo = "bar"
		work = "test"
	}
}
`
const testAccInstanceConfigVPC = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "testAccInstanceConfigVPC"
}

resource "alicloud_vpc" "foo" {
 	name = "${var.name}"
 	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"

	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
}
`

const testAccInstanceConfigUserData = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfigUserData"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/21"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
	user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}
`

const testAccInstanceConfigMultipleRegions = `
provider "alicloud" {
	alias = "beijing"
	region = "cn-beijing"
}
provider "alicloud" {
	alias = "shanghai"
	region = "cn-shanghai"
}
data "alicloud_zones" "default" {
	provider = "alicloud.beijing"
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
	provider = "alicloud.beijing"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
	provider = "alicloud.beijing"
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
data "alicloud_zones" "sh" {
	provider = "alicloud.shanghai"
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "sh" {
	provider = "alicloud.shanghai"
 	availability_zone = "${data.alicloud_zones.sh.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "sh" {
	provider = "alicloud.shanghai"
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfigMultipleRegions"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	provider = "alicloud.beijing"
	description = "foo"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}"
	provider = "alicloud.shanghai"
	description = "bar"
}

resource "alicloud_instance" "foo" {
  	# cn-beijing
  	provider = "alicloud.beijing"
  	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

  	internet_charge_type = "PayByBandwidth"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	system_disk_category = "cloud_efficiency"
  	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  	instance_name = "${var.name}"
}

resource "alicloud_instance" "bar" {
	# cn-shanghai
	provider = "alicloud.shanghai"
  	image_id = "${data.alicloud_images.sh.images.0.id}"
	availability_zone = "${data.alicloud_zones.sh.zones.0.id}"

	internet_charge_type = "PayByBandwidth"

	instance_type = "${data.alicloud_instance_types.sh.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.tf_test_bar.id}"]
	instance_name = "${var.name}"
}
`

const testAccInstanceConfig_multiSecurityGroup = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfig_multiSecurityGroup"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}-foo"
	description = "foo"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}-bar"
	description = "bar"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}", "${alicloud_security_group.tf_test_bar.id}"]
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
}`

const testAccInstanceConfig_multiSecurityGroup_add = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfig_multiSecurityGroup"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}-foo"
	description = "foo"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}-bar"
	description = "bar"
}

resource "alicloud_security_group" "tf_test_add_sg" {
	name = "${var.name}-add-sg"
	description = "sg"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}", "${alicloud_security_group.tf_test_bar.id}",
				"${alicloud_security_group.tf_test_add_sg.id}"]
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
}
`

const testAccInstanceConfig_multiSecurityGroup_remove = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfig_multiSecurityGroup"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}-foo"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
}
`

const testAccInstanceConfig_multiSecurityGroupByCount = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfig_multiSecurityGroupByCount"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	count = 2
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	security_groups = ["${alicloud_security_group.tf_test_foo.*.id}"]
	instance_name = "test_foo"
	system_disk_category = "cloud_efficiency"
}
`

const testAccInstanceNetworkInstanceSecurityGroups = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceNetworkInstanceSecurityGroups"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/21"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"

	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"

	internet_max_bandwidth_out = 5
	internet_charge_type = "PayByBandwidth"
}
`
const testAccCheckInstanceConfigTags = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceConfigTags"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"

	tags {
		foo = "bar"
	}
}
`

const testAccCheckInstanceConfigTagsUpdate = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceConfigTags"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "tf_test_foo"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"

	tags {
		bar = "zzz"
	}
}
`
const testAccCheckInstanceConfigOrigin = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceConfigOrigin"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "http-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "80/80"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	instance_name = "${var.name}-foo"
	host_name = "host-foo"
}
`

const testAccCheckInstanceConfigOriginUpdate = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceConfigOrigin"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "http-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "80/80"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	instance_name = "${var.name}-bar"
	host_name = "host-bar"
}
`
const testAccInstanceConfigAssociatePublicIP = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccInstanceConfigAssociatePublicIP"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	internet_max_bandwidth_out = 5
	internet_charge_type = "PayByBandwidth"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	instance_name = "${var.name}"
}
`
const testAccCheckInstanceImageOrigin = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "centos" {
	most_recent = true
	owners = "system"
	name_regex = "^centos_6\\w{1,5}[64]{1}.*"
}
variable "name" {
	default = "testAccCheckInstanceImageOrigin"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "update_image" {
	image_id = "${data.alicloud_images.centos.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 50
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	password = "Test12345"
  	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
}
`
const testAccCheckInstanceImageUpdate = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
variable "name" {
	default = "testAccCheckInstanceImageOrigin"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "update_image" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 60
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	password = "Test12345"
  	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccCheckInstancePrivateIp = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckInstancePrivateIp"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.2.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}-bar"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "private_ip" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
}
`

const testAccCheckInstancePrivateIpUpdate = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckInstancePrivateIp"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.2.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}-bar"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "private_ip" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.bar.id}"
	private_ip = "10.1.2.3"
}
`

const testAccCheckInstanceChargeType = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckInstanceChargeType"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "charge_type" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
	instance_charge_type = "PrePaid"
	period_unit = "Week"
}
`

const testAccCheckInstanceChargeTypeUpdate = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckInstanceChargeType"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "charge_type" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
}
`
const testAccCheckSpotInstance = `
data "alicloud_zones" "default" {
  available_disk_category= "cloud_efficiency"
  available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckSpotInstance"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "spot" {
  # cn-beijing
  vswitch_id = "${alicloud_vswitch.foo.id}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  instance_name = "${var.name}"
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit = "1.002"
}
`
const testAccCheckInstanceType = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
  	available_disk_category= "cloud_efficiency"
  	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceType"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "type" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
}
`

const testAccCheckInstanceTypeUpdate = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
  	available_disk_category= "cloud_efficiency"
  	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceType"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "type" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
}
`

const testAccCheckInstanceNetworkSpec = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
  	available_disk_category= "cloud_efficiency"
  	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceNetworkSpec"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "network" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
}
`

const testAccCheckInstanceNetworkSpecUpdate = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
  	available_disk_category= "cloud_efficiency"
  	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceType"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "network" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 40
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = ["${alicloud_security_group.group.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "10.1.1.3"
	internet_charge_type = "PayByBandwidth"
  	internet_max_bandwidth_out = 5
  	internet_max_bandwidth_in = 50
}
`
const testAccCheckInstanceRamRole = `
data "alicloud_images" "ubuntu" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
}
data "alicloud_zones" "default" {
  	available_disk_category= "cloud_efficiency"
  	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccCheckInstanceRamRole"
}
resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "role" {
	image_id = "${data.alicloud_images.ubuntu.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	system_disk_category = "cloud_efficiency"
  	system_disk_size = 60
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	password = "Test12345"
  	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	vswitch_id = "${alicloud_vswitch.foo.id}"
	role_name = "${alicloud_ram_role.role.name}"
}

resource "alicloud_ram_role" "role" {
  	name = "${var.name}"
  	services = ["ecs.aliyuncs.com"]
  	force = "true"
}

resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Allow"
      action = ["CreateInstance"]
      resource = ["*"]
    }
  ]
  force = "true"
}

resource "alicloud_ram_role_policy_attachment" "role-policy" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
`
