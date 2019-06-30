package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_launch_template", &resource.Sweeper{
		Name: "alicloud_launch_template",
		F:    testAlicloudLaunchTemplate,
	})
}

func testAccCheckLaunchTemplateDestroy(t *terraform.State) error {
	for _, rs := range t.RootModule().Resources {
		if rs.Type != "alicloud_launch_template" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No launch template is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeLaunchTemplate(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func testAlicloudLaunchTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %#v", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	req := ecs.CreateDescribeLaunchTemplatesRequest()
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplates(req)
	})
	if err != nil {
		return fmt.Errorf("describe launch template failed in sweepers, %v", err)
	}
	resp := raw.(*ecs.DescribeLaunchTemplatesResponse)
	if resp == nil {
		return fmt.Errorf("describe launch template version got nil")
	}

	var ids []string
	for _, tpl := range resp.LaunchTemplateSets.LaunchTemplateSet {
		if strings.HasPrefix(tpl.LaunchTemplateName, "tf-testAcc") {
			ids = append(ids, tpl.LaunchTemplateId)
		}
	}

	for i := range ids {
		args := ecs.CreateDeleteLaunchTemplateRequest()
		args.LaunchTemplateId = ids[i]
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteLaunchTemplate(args)
		})
		if err != nil || !NotFoundError(err) {
			log.Printf("delete template failed in sweepers, %v", err)
		}
	}

	return nil
}

func TestAccAlicloudLaunchTemplate_basic(t *testing.T) {
	var tpl ecs.LaunchTemplateVersionSet
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_launch_template.template",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLaunchTemplateDestroy,

		Steps: []resource.TestStep{
			{
				Config: lauchTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckLaunchTemplateExists("alicloud_launch_template.template", &tpl),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "name", "tf-testAcc-template"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "description", "test1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "image_id", "m-xlhga1245bbngasdfgd"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "host_name", "tf-test-host"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_name", "tf-instance-name"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_type", "ecs.g5.8xlarge"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_max_bandwidth_in", "5"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "io_optimized", "none"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "key_pair_name", "test-key-pair"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "ram_role_name", "xxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_type", "vpc"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "spot_price_limit", "5"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "spot_strategy", "SpotWithPriceLimit"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "security_group_id", "sg-zxcvj0lasdf102350asdf9a"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_description", "test disk"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_name", "hello"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "resource_group_id", "rg-zkdfjahg9zxncv0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "userdata", "xxxxxxxxxxxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "vswitch_id", "sw-ljkngaksdjfj0nnasdf"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "vpc_id", "vpc-asdfnbg0as8dfk1nb2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "zone_id", "beijing-a"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.tag1", "hello"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.tag2", "world"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.#", "1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.name", "eth0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.description", "hello1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.primary_ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.security_group_id", "xxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.vswitch_id", "xxxxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.#", "2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.0.name", "disk1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.0.description", "test1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.1.name", "disk2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.1.description", "test2")),
			},
			{
				Config: lauchTemplateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckLaunchTemplateExists("alicloud_launch_template.template", &tpl),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "name", "tf-testAcc-template"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "description", "test2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "image_id", "m-xlhga1245bbngasdfgd"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "host_name", "tf-test-host"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_name", "tf-instance-name"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "instance_type", "ecs.g5.8xlarge"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_max_bandwidth_in", "5"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "io_optimized", "none"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "key_pair_name", "test-key-pair"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "ram_role_name", "xxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_type", "vpc"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "spot_price_limit", "5"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "spot_strategy", "SpotWithPriceLimit"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "security_group_id", "sg-zxcvj0lasdf102350asdf9a"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_description", "test disk"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_name", "hello"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "system_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "resource_group_id", "rg-zkdfjahg9zxncv0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "userdata", "xxxxxxxxxxxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "vswitch_id", "sw-ljkngaksdjfj0nnasdf"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "vpc_id", "vpc-asdfnbg0as8dfk1nb2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "zone_id", "beijing-a"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.tag1", "bye"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "tags.tag2", "world"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.#", "1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.name", "eth0"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.description", "hello1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.primary_ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.security_group_id", "xxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "network_interfaces.0.vswitch_id", "xxxxxxx"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.#", "2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.0.name", "disk1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.0.description", "test1"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.1.name", "disk2"),
					resource.TestCheckResourceAttr("alicloud_launch_template.template", "data_disks.1.description", "test2")),
			},
		},
	})
}

func testCheckLaunchTemplateExists(n string, tpl *ecs.LaunchTemplateVersionSet) resource.TestCheckFunc {
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

		d, err := ecsService.DescribeLaunchTemplate(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("While checking LauchTemplate existing, describing LauchTemplate got an error: %#v.", err)
		}

		*tpl = d

		return nil
	}
}

const lauchTemplateConfig = `
resource "alicloud_launch_template" "template" {
    name = "tf-testAcc-template"
    description = "test1"
    image_id = "m-xlhga1245bbngasdfgd"
    host_name = "tf-test-host"
    instance_charge_type = "PrePaid"
    instance_name = "tf-instance-name"
    instance_type = "ecs.g5.8xlarge"
    internet_charge_type = "PayByBandwidth"
    internet_max_bandwidth_in = 5
    internet_max_bandwidth_out = 0
    io_optimized = "none"
    key_pair_name = "test-key-pair"
    ram_role_name = "xxxxx"
    network_type = "vpc"
    security_enhancement_strategy = "Active"
    spot_price_limit = 5
    spot_strategy = "SpotWithPriceLimit"
    security_group_id = "sg-zxcvj0lasdf102350asdf9a"
    system_disk_category = "cloud_ssd"
    system_disk_description = "test disk"
    system_disk_name = "hello"
    system_disk_size = 40
    resource_group_id = "rg-zkdfjahg9zxncv0"
    userdata = "xxxxxxxxxxxxxx"
    vswitch_id = "sw-ljkngaksdjfj0nnasdf"
    vpc_id = "vpc-asdfnbg0as8dfk1nb2"
    zone_id = "beijing-a"
    
    tags = {
        tag1 = "hello"
        tag2 = "world"
    }
    network_interfaces {
            name = "eth0"
            description = "hello1"
            primary_ip = "10.0.0.2"
            security_group_id = "xxxx"
            vswitch_id = "xxxxxxx"
        }
    data_disks {
            name = "disk1"
            description = "test1"
        }
    data_disks {
            name = "disk2"
            description = "test2"
        }
}
`
const lauchTemplateConfigUpdate = `
resource "alicloud_launch_template" "template" {
    name = "tf-testAcc-template"
    description = "test2"
    image_id = "m-xlhga1245bbngasdfgd"
    host_name = "tf-test-host"
    instance_charge_type = "PrePaid"
    instance_name = "tf-instance-name"
    instance_type = "ecs.g5.8xlarge"
    internet_charge_type = "PayByBandwidth"
    internet_max_bandwidth_in = 5
    internet_max_bandwidth_out = 0
    io_optimized = "none"
    key_pair_name = "test-key-pair"
    ram_role_name = "xxxxx"
    network_type = "vpc"
    security_enhancement_strategy = "Active"
    spot_price_limit = 5
    spot_strategy = "SpotWithPriceLimit"
    security_group_id = "sg-zxcvj0lasdf102350asdf9a"
    system_disk_category = "cloud_ssd"
    system_disk_description = "test disk"
    system_disk_name = "hello"
    system_disk_size = 40
    resource_group_id = "rg-zkdfjahg9zxncv0"
    userdata = "xxxxxxxxxxxxxx"
    vswitch_id = "sw-ljkngaksdjfj0nnasdf"
    vpc_id = "vpc-asdfnbg0as8dfk1nb2"
    zone_id = "beijing-a"
    
    tags = {
        tag1 = "bye"
        tag2 = "world"
    }
    network_interfaces {
            name = "eth0"
            description = "hello1"
            primary_ip = "10.0.0.2"
            security_group_id = "xxxx"
            vswitch_id = "xxxxxxx"
        }
    data_disks {
            name = "disk1"
            description = "test1"
        }
    data_disks {
            name = "disk2"
            description = "test2"
        }
}
`
