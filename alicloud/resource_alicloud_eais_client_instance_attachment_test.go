package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eais ClientInstanceAttachment. >>> Resource test cases, automatically generated.
// Case ca_ei_pro 10133
func TestAccAliCloudEaisClientInstanceAttachment_basic10133(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_client_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisClientInstanceAttachmentMap10133)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisClientInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceais%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisClientInstanceAttachmentBasicDependence10133)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_eais_instance.eais.id}",
					"client_instance_id": "${alicloud_instance.example.id}",
					"category":           "ei",
					"status":             "Bound",
					"ei_instance_type":   "eais.ei-a6.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"client_instance_id": CHECKSET,
						"status":             "Bound",
						"ei_instance_type":   "eais.ei-a6.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "InUse",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "InUse",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Bound",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Bound",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category"},
			},
		},
	})
}

var AlicloudEaisClientInstanceAttachmentMap10133 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEaisClientInstanceAttachmentBasicDependence10133(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "ecs_image" {
  default = "ubuntu_20_04_x64_20G_alibase_20230316.vhd"
}

variable "ecs_type" {
  default = "ecs.g7.large"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "category" {
  default = "ei"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "example" {
  availability_zone = "cn-hangzhou-i"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  description         = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = "cn-hangzhou-i"
  vswitch_id                 = alicloud_vswitch.example.id
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.example.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

resource "alicloud_eais_instance" "eais" {
  instance_name     = var.name
  vswitch_id        = alicloud_vswitch.example.id
  security_group_id = alicloud_security_group.example.id
  instance_type     = "eais.ei-a6.2xlarge"
  category = "ei"
}


`, name)
}

// Case ca_eai_pro 10134
func TestAccAliCloudEaisClientInstanceAttachment_basic10134(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_client_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisClientInstanceAttachmentMap10134)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisClientInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceais%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisClientInstanceAttachmentBasicDependence10134)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_eais_instance.eais.id}",
					"client_instance_id": "${alicloud_instance.example.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"client_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEaisClientInstanceAttachmentMap10134 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEaisClientInstanceAttachmentBasicDependence10134(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "ecs_image" {
  default = "ubuntu_20_04_x64_20G_alibase_20230316.vhd"
}

variable "ecs_type" {
  default = "ecs.g7.large"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "category" {
  default = "eais"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "example" {
  availability_zone = "cn-hangzhou-i"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  description         = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = "cn-hangzhou-i"
  vswitch_id                 = alicloud_vswitch.example.id
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.example.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

resource "alicloud_eais_instance" "eais" {
  instance_name     = var.name
  vswitch_id        = alicloud_vswitch.example.id
  security_group_id = alicloud_security_group.example.id
  instance_type     = "eais.ei-a6.2xlarge"
}


`, name)
}

// Case ca_pre_eai 10092
func TestAccAliCloudEaisClientInstanceAttachment_basic10092(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_client_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisClientInstanceAttachmentMap10092)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisClientInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceais%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisClientInstanceAttachmentBasicDependence10092)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_eais_instance.eais.id}",
					"client_instance_id": "${alicloud_instance.example.id}",
					"category":           "eais",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"client_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category"},
			},
		},
	})
}

var AlicloudEaisClientInstanceAttachmentMap10092 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEaisClientInstanceAttachmentBasicDependence10092(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "ecs_image" {
  default = "ubuntu_20_04_x64_20G_alibase_20230316.vhd"
}

variable "ecs_type" {
  default = "ecs.g7.large"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "category" {
  default = "eais"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "example" {
  availability_zone = "cn-hangzhou-i"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  description         = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = "cn-hangzhou-i"
  vswitch_id                 = alicloud_vswitch.example.id
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.example.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

resource "alicloud_eais_instance" "eais" {
  instance_name     = var.name
  vswitch_id        = alicloud_vswitch.example.id
  security_group_id = alicloud_security_group.example.id
  instance_type     = "eais.ei-a6.2xlarge"
}


`, name)
}

// Test Eais ClientInstanceAttachment. <<< Resource test cases, automatically generated.
