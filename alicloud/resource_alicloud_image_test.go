package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSImageBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "alicloud_image.default"
	ra := resourceAttrInit(resourceId, testAccImageCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_instance.default.id}",
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":   name,
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "acceptance test1232",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.For":     "acceptance test1232",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"image_name":   name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudECSImageBasic1(t *testing.T) {
	var v ecs.Image

	resourceId := "alicloud_image.default"
	ra := resourceAttrInit(resourceId, testAccImageCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageBasicConfigDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
					"snapshot_id": "${alicloud_ecs_snapshot.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":   name,
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
						"snapshot_id":  CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudECSImageBasic2(t *testing.T) {
	var v ecs.Image

	resourceId := "alicloud_image.default"
	ra := resourceAttrInit(resourceId, testAccImageCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageBasicConfigDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
					"disk_device_mapping": []map[string]interface{}{
						{
							"disk_type":   "system",
							"device":      "/dev/xvda",
							"size":        "2000",
							"snapshot_id": "${alicloud_ecs_snapshot.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":            name,
						"description":           fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test123",
						"disk_device_mapping.#": "1",
					}),
				),
			},
		},
	})
}

var testAccImageCheckMap = map[string]string{}

func resourceImageBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instance_types" "default" {
 	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.ids[0]}"
  instance_type = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups = "${[alicloud_security_group.default.id]}"
  vswitch_id = local.vswitch_id
  instance_name = "${var.name}"
}
`, name)
}

func resourceImageBasicConfigDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}


data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.ids[0]}"
  instance_type = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups = "${[alicloud_security_group.default.id]}"
  vswitch_id = local.vswitch_id
  instance_name = "${var.name}"
}

resource "alicloud_disk" "default" {
  count = "2"
  disk_name = "${var.name}"
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category          = "cloud_efficiency"
  size              = "20"
}

resource "alicloud_disk_attachment" "default" {
  count = "2"
  disk_id     = "${element(alicloud_disk.default.*.id,count.index)}"
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot" "default" {
	category = "standard"
	description = "Test For Terraform"
	disk_id = alicloud_disk_attachment.default.0.disk_id
	retention_days = "20"
	snapshot_name = var.name
	tags 				 = {
		Created = "TF"
		For 	= "Acceptance-test"
	}
}

`, name)
}
