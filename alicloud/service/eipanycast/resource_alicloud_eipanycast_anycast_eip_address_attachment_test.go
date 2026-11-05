package eipanycast_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eipanycast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEipanycastAnycastEipAddressAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eipanycast_anycast_eip_address_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipanycastAnycastEipAddressAttachmentMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eipanycast.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipanycastAnycastEipAddressAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEipanycastAnycastEipAddressAttachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "bind_instance_region_id" must be consistent with the region of slb instance.
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address.default.id}",
					"bind_instance_id":        "${alicloud_slb_load_balancer.default.id}",
					"bind_instance_region_id": "${data.alicloud_regions.current_regions.regions.0.id}",
					"bind_instance_type":      "SlbInstance",
					// PrivateIpAddress only takes effect for BindInstanceType=NetworkInterface
					// (verified via AssociateAnycastEipAddress OpenAPI). For an SlbInstance
					// binding the field is not effective: Create does not send it and Read
					// reads back an empty value. Assert that actual read-back here.
					"private_ip_address": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anycast_id":              acctest.CHECKSET,
						"bind_instance_id":        acctest.CHECKSET,
						"bind_instance_region_id": acctest.CHECKSET,
						"bind_instance_type":      "SlbInstance",
						"private_ip_address":      "",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// pop_locations is cloud-populated with default pop locations for a
				// ChineseMainland anycast attachment when not configured; this basic
				// case does not manage pop_locations, so ignore it during import verify.
				ImportStateVerifyIgnore: []string{"pop_locations"},
			},
		},
	})
}

var AlicloudEipanycastAnycastEipAddressAttachmentMap = map[string]string{
	"bind_time": acctest.CHECKSET,
}

func AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
	address_type = "intranet"
	vswitch_id = data.alicloud_vswitches.default.ids[0]
	load_balancer_name = var.name
	load_balancer_spec = "slb.s1.small"
    master_zone_id = "${data.alicloud_slb_zones.default.zones.0.id}"
}

resource "alicloud_eipanycast_anycast_eip_address" "default" {
  anycast_eip_address_name = "${var.name}"
  service_location = "ChineseMainland"
}

data "alicloud_regions" "current_regions" {
  current = true
}

`, name)
}

// Test Eipanycast AnycastEipAddressAttachment. >>> Resource test cases, automatically generated.
// Case 3732
func TestAccAliCloudEipanycastAnycastEipAddressAttachment_basic3732(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_eipanycast_anycast_eip_address_attachment.default"

	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipanycastAnycastEipAddressAttachmentMap3732)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipanycastanycasteipaddressattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence3732)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckEipanycastEipAddressAttachmentDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_instance_id":        "${alicloud_instance.defaultEcs2.network_interface_id}",
					"bind_instance_type":      "NetworkInterface",
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address_attachment.defaultEfYBJY.anycast_id}",
					"bind_instance_region_id": "${data.alicloud_regions.current_regions.regions.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bind_instance_id":        acctest.CHECKSET,
						"bind_instance_type":      "NetworkInterface",
						"anycast_id":              acctest.CHECKSET,
						"bind_instance_region_id": acctest.CHECKSET,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"association_mode": "Default",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
			//		testAccCheck(map[string]string{
			//			"association_mode": "Default",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"pop_locations": []map[string]interface{}{
						{
							"pop_location": "cn-guangzhou-pop",
						},
						{
							"pop_location": "cn-shanghai-pop",
						},
						{
							"pop_location": "cn-beijing-pop",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"pop_locations.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pop_locations": []map[string]interface{}{
						{
							"pop_location": "cn-chengdu-pop",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"pop_locations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pop_locations": []map[string]interface{}{
						{
							"pop_location": "cn-guangzhou-pop",
						},
						{
							"pop_location": "cn-shanghai-pop",
						},
						{
							"pop_location": "cn-beijing-pop",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"pop_locations.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pop_locations": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"pop_locations.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pop_locations": []map[string]interface{}{
						{
							"pop_location": "cn-chengdu-pop",
						},
						{
							"pop_location": "cn-beijing-pop",
						},
					},
					"association_mode": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"pop_locations.#":  "2",
						"association_mode": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_instance_id":        "${alicloud_instance.defaultEcs2.network_interface_id}",
					"bind_instance_type":      "NetworkInterface",
					"association_mode":        "Normal",
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id}",
					"bind_instance_region_id": "${data.alicloud_regions.current_regions.regions.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bind_instance_id":        acctest.CHECKSET,
						"bind_instance_type":      "NetworkInterface",
						"association_mode":        "Normal",
						"anycast_id":              acctest.CHECKSET,
						"bind_instance_region_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudEipanycastAnycastEipAddressAttachmentMap3732 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence3732(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  alias  = "beijing"
  region = "cn-beijing"
}

provider "alicloud" {
  alias  = "hangzhou"
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  provider                    = "alicloud.beijing"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  provider    = "alicloud.beijing"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  provider          = "alicloud.beijing"
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "defaultVpc" {
  provider   = "alicloud.beijing"
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVsw" {
  provider   = "alicloud.beijing"
  vpc_id     = alicloud_vpc.defaultVpc.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "defaultuBsECI" {
  provider = "alicloud.beijing"
  vpc_id   = alicloud_vpc.defaultVpc.id
}

resource "alicloud_instance" "default9KDlN7" {
  provider                   = "alicloud.beijing"
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = ["${alicloud_security_group.defaultuBsECI.id}"]
  availability_zone          = alicloud_vswitch.defaultVsw.zone_id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.defaultVsw.id
}

resource "alicloud_eipanycast_anycast_eip_address" "defaultXkpFRs" {
  provider         = "alicloud.hangzhou"
  service_location = "ChineseMainland"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "defaultEfYBJY" {
  provider                = "alicloud.beijing"
  bind_instance_id        = alicloud_instance.default9KDlN7.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-beijing"
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id
}

resource "alicloud_vpc" "defaultVpc2" {
  provider   = "alicloud.hangzhou"
  vpc_name   = "${var.name}6"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_regions" "current_regions" {
  provider = "alicloud.hangzhou"
  current = true
}

data "alicloud_zones" "default2" {
  provider                    = "alicloud.hangzhou"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default2" {
  provider    = "alicloud.hangzhou"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default2" {
  provider          = "alicloud.hangzhou"
  availability_zone = data.alicloud_zones.default2.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vswitch" "defaultdsVsw2" {
  provider   = "alicloud.hangzhou"
  vpc_id     = alicloud_vpc.defaultVpc2.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default2.zones.1.id
}

resource "alicloud_security_group" "defaultuBsECI2" {
  provider = "alicloud.hangzhou"
  vpc_id   = alicloud_vpc.defaultVpc2.id
}

resource "alicloud_instance" "defaultEcs2" {
  provider                   = "alicloud.hangzhou"
  image_id                   = data.alicloud_images.default2.images[0].id
  instance_type              = data.alicloud_instance_types.default2.instance_types[0].id
  instance_name              = var.name
  security_groups            = ["${alicloud_security_group.defaultuBsECI2.id}"]
  availability_zone          = alicloud_vswitch.defaultdsVsw2.zone_id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.defaultdsVsw2.id
}


`, name)
}

func AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence3732_region(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  alias  = "hongkong"
  region = "cn-hongkong"
}

provider "alicloud" {
  alias  = "central"
  region = "eu-central-1"
}

data "alicloud_zones" "default" {
  provider                    = "alicloud.hongkong"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  provider    = "alicloud.hongkong"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  provider          = "alicloud.hongkong"
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "defaultVpc" {
  provider   = "alicloud.hongkong"
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVsw" {
  provider   = "alicloud.hongkong"
  vpc_id     = alicloud_vpc.defaultVpc.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "defaultuBsECI" {
  provider = "alicloud.hongkong"
  vpc_id   = alicloud_vpc.defaultVpc.id
}

resource "alicloud_instance" "default9KDlN7" {
  provider                   = "alicloud.hongkong"
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = ["${alicloud_security_group.defaultuBsECI.id}"]
  availability_zone          = alicloud_vswitch.defaultVsw.zone_id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.defaultVsw.id
}

data "alicloud_regions" "current_regions" {
  provider         = "alicloud.central"
  current = true
}

resource "alicloud_eipanycast_anycast_eip_address" "defaultXkpFRs" {
  provider         = "alicloud.central"
  service_location = "international"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "defaultEfYBJY" {
  provider                = "alicloud.hongkong"
  bind_instance_id        = alicloud_instance.default9KDlN7.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-hongkong"
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id
}

resource "alicloud_vpc" "defaultVpc2" {
  provider   = "alicloud.central"
  vpc_name   = "${var.name}6"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "default2" {
  provider                    = "alicloud.central"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default2" {
  provider    = "alicloud.central"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default2" {
  provider          = "alicloud.central"
  availability_zone = data.alicloud_zones.default2.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vswitch" "defaultdsVsw2" {
  provider   = "alicloud.central"
  vpc_id     = alicloud_vpc.defaultVpc2.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default2.zones.1.id
}

resource "alicloud_security_group" "defaultuBsECI2" {
  provider = "alicloud.central"
  vpc_id   = alicloud_vpc.defaultVpc2.id
}

resource "alicloud_instance" "defaultEcs2" {
  provider                   = "alicloud.central"
  image_id                   = data.alicloud_images.default2.images[0].id
  instance_type              = data.alicloud_instance_types.default2.instance_types[0].id
  instance_name              = var.name
  security_groups            = ["${alicloud_security_group.defaultuBsECI2.id}"]
  availability_zone          = alicloud_vswitch.defaultdsVsw2.zone_id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.defaultdsVsw2.id
}


`, name)
}

// Case 3732  twin
func TestAccAliCloudEipanycastAnycastEipAddressAttachment_basic3732_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_eipanycast_anycast_eip_address_attachment.default"

	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipanycastAnycastEipAddressAttachmentMap3732)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipanycastanycasteipaddressattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence3732)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckEipanycastEipAddressAttachmentDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_instance_id":        "${alicloud_instance.defaultEcs2.network_interface_id}",
					"bind_instance_type":      "NetworkInterface",
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address_attachment.defaultEfYBJY.anycast_id}",
					"bind_instance_region_id": "${data.alicloud_regions.current_regions.regions.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bind_instance_id":        acctest.CHECKSET,
						"bind_instance_type":      "NetworkInterface",
						"anycast_id":              acctest.CHECKSET,
						"bind_instance_region_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEipanycastAnycastEipAddressAttachment_basic3732_region(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_eipanycast_anycast_eip_address_attachment.default"

	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	ra := acctest.ResourceAttrInit(resourceId, AlicloudEipanycastAnycastEipAddressAttachmentMap3732)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipanycastanycasteipaddressattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence3732_region)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckEipanycastEipAddressAttachmentDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_instance_id":        "${alicloud_instance.defaultEcs2.network_interface_id}",
					"bind_instance_type":      "NetworkInterface",
					"association_mode":        "Normal",
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id}",
					"bind_instance_region_id": "eu-central-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bind_instance_id":        acctest.CHECKSET,
						"bind_instance_type":      "NetworkInterface",
						"association_mode":        "Normal",
						"anycast_id":              acctest.CHECKSET,
						"bind_instance_region_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckEipanycastEipAddressAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckEipanycastEipAddressAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckEipanycastEipAddressAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	eipanycastServiceV2 := eipanycast.NewService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eipanycast_anycast_eip_address_attachment" {
			continue
		}
		_, err := eipanycastServiceV2.DescribeEipanycastAnycastEipAddressAttachment(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("alicloud_eipanycast_anycast_eip_address_attachment still exist, ID %s ", fmt.Sprint(rs.Primary.ID))
		}
	}

	return nil
}

func testAccCheckEipanycastEipAddressAttachmentExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_eipanycast_anycast_eip_address_attachment ID is set. ")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			eipanycastServiceV2 := eipanycast.NewService(client)

			resp, err := eipanycastServiceV2.DescribeEipanycastAnycastEipAddressAttachment(rs.Primary.ID)
			if err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_inter_region_traffic_qos_queue not found")
	}
}

// Test Eipanycast AnycastEipAddressAttachment. <<< Resource test cases, automatically generated.
