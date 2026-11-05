package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudECSSnapshot_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id": "${alicloud_ecs_disk_attachment.default.disk_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id": CHECKSET,
						"status":  "accomplished",
					}),
					func(s *terraform.State) error {
						if value := s.RootModule().Resources[resourceId].Primary.Attributes["wait_until"]; value != "" {
							return fmt.Errorf("omitted wait_until must remain unset, got %q", value)
						}
						return nil
					},
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":           "${alicloud_ecs_disk_attachment.default.disk_id}",
					"category":          "standard",
					"wait_until":        "accomplished",
					"retention_days":    "50",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"snapshot_name":     name,
					"description":       name,
					"force":             "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":           CHECKSET,
						"category":          "standard",
						"wait_until":        "accomplished",
						"status":            "accomplished",
						"retention_days":    "50",
						"resource_group_id": CHECKSET,
						"snapshot_name":     name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":        "${alicloud_ecs_disk_attachment.default.disk_id}",
					"retention_days": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":        CHECKSET,
						"retention_days": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":           "${alicloud_ecs_disk_attachment.default.disk_id}",
					"category":          "standard",
					"retention_days":    "50",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"name":              name,
					"description":       name,
					"force":             "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":           CHECKSET,
						"category":          "standard",
						"retention_days":    "50",
						"resource_group_id": CHECKSET,
						"name":              name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

var AliCloudEcsSnapshotMap0 = map[string]string{
	"available":   "true",
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"category":    CHECKSET,
	"status":      CHECKSET,
}

func AliCloudEcsSnapshotBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_essd"
	available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
	available_resource_creation = "Instance"
}

data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
	architecture = "x86_64"
	os_type      = "linux"
}

data "alicloud_instance_types" "default" {
	instance_type_family = "ecs.g7"
	sorted_by            = "CPU"
  image_id             = data.alicloud_images.default.images.0.id
  system_disk_category = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
	availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_essd"
	system_disk_encrypted      = true
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
	data_disks {
		category = "cloud_essd"
		encrypted = true
		size     = 20
	}
}

resource "alicloud_ecs_disk" "default" {
  disk_name = var.name
	zone_id   = data.alicloud_zones.default.zones.0.id
  category  = "cloud_essd"
	encrypted = true
	size      = 20
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id     = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}
`, name)
}

func AliCloudEcsSnapshotBasicDependence1(name string) string {
	return AliCloudEcsSnapshotBasicDependence0(name)
}

func TestAccAliCloudECSSnapshot_availableDisk(t *testing.T) {
	var snapshot, disk map[string]interface{}
	snapshotID := "alicloud_ecs_snapshot.default"
	diskID := "alicloud_ecs_disk.restored"
	snapshotCheck := resourceCheckInitWithDescribeMethod(snapshotID, &snapshot, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsSnapshot")
	diskCheck := resourceCheckInitWithDescribeMethod(diskID, &disk, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsDisk")
	name := fmt.Sprintf("tf-testacc-snapshot-%d", tfacctest.RandIntRange(10000, 99999))
	config := AliCloudEcsSnapshotBasicDependence0(name) + `
resource "alicloud_ecs_snapshot" "default" {
  disk_id       = alicloud_ecs_disk_attachment.default.disk_id
  snapshot_name = var.name
  force         = true
  wait_until    = "available"
}

resource "alicloud_ecs_disk" "restored" {
  disk_name   = "${var.name}-restored"
  zone_id     = alicloud_ecs_disk.default.zone_id
  category    = "cloud_essd"
  snapshot_id = alicloud_ecs_snapshot.default.id
  encrypted   = true
  size        = 20
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			snapshotCheck.checkResourceDestroy(), diskCheck.checkResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					snapshotCheck.checkResourceExists(),
					resource.TestCheckResourceAttr(snapshotID, "wait_until", "available"),
					resource.TestCheckResourceAttr(snapshotID, "available", "true"),
					diskCheck.checkResourceExists(),
					resource.TestCheckResourceAttrPair(diskID, "snapshot_id", snapshotID, "id"),
					func(s *terraform.State) error {
						service := ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
						restored, err := service.DescribeEcsDisk(s.RootModule().Resources[diskID].Primary.ID)
						if err != nil {
							return err
						}
						if fmt.Sprint(restored["SourceSnapshotId"]) != s.RootModule().Resources[snapshotID].Primary.ID {
							return fmt.Errorf("restored disk must use the created snapshot")
						}
						return nil
					},
				),
			},
			{
				ResourceName: snapshotID, ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}
