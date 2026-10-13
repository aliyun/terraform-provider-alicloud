package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudECSImageBasic(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_instance.default.id}",
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"name":        name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
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
					"name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
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
					"name":        name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"name":         name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudECSImageBasic1(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
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

func TestAccAliCloudECSImageBasic2(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageBasicConfigDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
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
							"size":        20,
							"snapshot_id": "${alicloud_ecs_snapshot.default.id}",
						},
						{
							"disk_type":   "data",
							"device":      "/dev/xvdb",
							"size":        20,
							"snapshot_id": "${alicloud_ecs_snapshot.data.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":                        name,
						"description":                       fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":                            "2",
						"tags.Created":                      "TF",
						"tags.For":                          "acceptance test123",
						"disk_device_mapping.#":             "2",
						"disk_device_mapping.0.disk_type":   "system",
						"disk_device_mapping.0.device":      "/dev/xvda",
						"disk_device_mapping.0.size":        "20",
						"disk_device_mapping.0.snapshot_id": CHECKSET,
						"disk_device_mapping.1.disk_type":   "data",
						"disk_device_mapping.1.device":      "/dev/xvdb",
						"disk_device_mapping.1.size":        "20",
						"disk_device_mapping.1.snapshot_id": CHECKSET,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "disk_device_mapping.0.snapshot_id", "alicloud_ecs_snapshot.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "disk_device_mapping.1.snapshot_id", "alicloud_ecs_snapshot.data", "id"),
				),
			},
		},
	})
}

var testAccImageCheckMap = map[string]string{
	"usable": "true",
}

func resourceImageBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g8i"
  system_disk_category = "cloud_essd"
}

locals {
  zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
  instance_type = data.alicloud_instance_types.default.ids.0
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id           = local.zone_id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id             = "${data.alicloud_images.default.ids[0]}"
  instance_type        = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups      = "${[alicloud_security_group.default.id]}"
  vswitch_id           = alicloud_vswitch.default.id
  system_disk_category = "cloud_essd"
  instance_name        = "${var.name}"
  status               = "Stopped"
}
`, name)
}

func resourceImageBasicConfigDependence1(name string) string {
	return resourceImageBasicConfigDependenceWithSnapshot(name, "")
}

func resourceImageBasicConfigDependence2(name string) string {
	return resourceImageBasicConfigDependenceWithSnapshot(name, "  system_disk_size     = 20\n") + `
resource "alicloud_disk" "data" {
  disk_name         = var.name
  availability_zone = local.zone_id
  category          = "cloud_essd"
  size              = 20
}

resource "alicloud_disk_attachment" "data" {
  disk_id     = alicloud_disk.data.id
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot" "data" {
  category      = "standard"
  description   = "Test data disk snapshot for Terraform"
  disk_id       = alicloud_disk_attachment.data.disk_id
  snapshot_name = "${var.name}-data"

  timeouts {
    create = "10m"
  }
}
`
}

func resourceImageBasicConfigDependenceWithSnapshot(name, systemDiskSize string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g8i"
  system_disk_category = "cloud_essd"
}

locals {
  zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
  instance_type = data.alicloud_instance_types.default.ids.0
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id           = local.zone_id
  vswitch_name      = var.name
}

locals {
  vswitch_id = alicloud_vswitch.default.id
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_instance" "default" {
  image_id             = "${data.alicloud_images.default.ids[0]}"
  instance_type        = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups      = "${[alicloud_security_group.default.id]}"
  vswitch_id           = local.vswitch_id
  system_disk_category = "cloud_essd"
%s  instance_name        = "${var.name}"
  status               = "Stopped"
}

resource "alicloud_ecs_snapshot" "default" {
	category = "standard"
	description = "Test For Terraform"
	disk_id = alicloud_instance.default.system_disk_id
	timeouts {
		create = "10m"
	}
	retention_days = "20"
	snapshot_name = var.name
	tags 				 = {
		Created = "TF"
		For 	= "Acceptance-test"
	}
}

`, name, systemDiskSize)
}

func TestAccAliCloudECSImageBasic7009(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_image.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsImageMap7009)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsImageBasicDependence7009)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name":           name,
					"instance_id":          "${alicloud_instance.default.id}",
					"platform":             "Ubuntu",
					"force":                "true",
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name,
						"platform":   "Ubuntu",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "create",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "create",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"boot_mode":    "BIOS",
					"license_type": "BYOL",
					"features": []map[string]interface{}{
						{
							"nvme_support": "supported",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"boot_mode":    "BIOS",
						"license_type": "BYOL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_family": "test-tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_family": "test-tf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-creat",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-creat",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"boot_mode": "UEFI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"boot_mode": "UEFI",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_family": "test-tf-123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_family": "test-tf-123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-aaaa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-aaaa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "create",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "create",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"boot_mode": "BIOS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"boot_mode": "BIOS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_family": "test-tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_family": "test-tf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":        "create",
					"instance_id":        "${alicloud_instance.default.id}",
					"image_name":         name + "_update",
					"detection_strategy": "Standard",
					"architecture":       "x86_64",
					"boot_mode":          "BIOS",
					"image_family":       "test-tf",
					"image_version":      "1",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "create",
						"instance_id":        CHECKSET,
						"image_name":         name + "_update",
						"detection_strategy": "Standard",
						"architecture":       "x86_64",
						"boot_mode":          "BIOS",
						"image_family":       "test-tf",
						"image_version":      "1",
						"resource_group_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detection_strategy", "features", "instance_id", "license_type", "snapshot_id", "delete_auto_snapshot", "force"},
			},
		},
	})
}

var AlicloudEcsImageMap7009 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"usable":      "true",
}

func AlicloudEcsImageBasicDependence7009(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g8i"
  system_disk_category = "cloud_essd"
}

locals {
  zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
  instance_type = data.alicloud_instance_types.default.ids.0
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id           = local.zone_id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_instance" "default" {
  image_id             = "${data.alicloud_images.default.ids[0]}"
  instance_type        = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups      = "${[alicloud_security_group.default.id]}"
  vswitch_id           = alicloud_vswitch.default.id
  system_disk_category = "cloud_essd"
  instance_name        = "${var.name}"
  status               = "Stopped"
}

`, name)
}
