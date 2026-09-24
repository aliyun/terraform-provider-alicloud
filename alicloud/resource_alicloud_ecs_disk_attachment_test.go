package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudECSDiskAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudECSDiskAttachmentMulti(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default.1"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependenceMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudECSDiskAttachmentBasic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default.1"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependence1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var AlicloudEcsDiskAttachmentMap = map[string]string{
	"disk_id":     CHECKSET,
	"instance_id": CHECKSET,
}

func AlicloudEcsDiskAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
   name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
 	most_recent = true
	owners = "system"
}
resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	zone_id = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}
`, name)
}
func AlicloudEcsDiskAttachmentBasicDependenceMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

variable "number" {
	default = "2"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_images" "default" {
    name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  	most_recent = true
	owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	count = "${var.number}"
	zone_id = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  count = "${var.number}"
  disk_id = "${element(alicloud_ecs_disk.default.*.id, count.index)}"
  instance_id = alicloud_instance.default.id
}
`, name)
}

func AlicloudEcsDiskAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "number" {
  default = "2"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6"
}

data "alicloud_images" "default" {
  name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent   = true
  owners        = "system"
  instance_type = data.alicloud_instance_types.default.instance_types.0.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_name     = var.name
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_ecs_disk" "default" {
  count                = var.number
  zone_id              = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category             = "cloud_efficiency"
  delete_auto_snapshot = "true"
  description          = "Test For Terraform"
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created     = "TF"
    Environment = "Acceptance-test"
  }
}

resource "alicloud_ecs_disk_attachment" "default" {
  count         = var.number
  disk_id       = element(alicloud_ecs_disk.default.*.id, count.index)
  instance_id   = alicloud_instance.default.id
  password      = "YouPassword123"
  key_pair_name = var.name
}
`, name)
}
