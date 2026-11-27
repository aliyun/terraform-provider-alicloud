package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass ip_allocation_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.bond_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass vpd_subnets -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEfloNodeGroup_basic10344(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupMap10344)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupBasicDependence10344)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_group_name":           name,
					"node_group_description":    "resource-test1",
					"cluster_id":                "${alicloud_eflo_cluster.default.id}",
					"machine_type":              "efg2.GN9De.cn8",
					"az":                        "cn-wulanchabu-d",
					"image_id":                  "i195757131775635040498",
					"login_password":            "YourPassword123!",
					"key_pair_name":             "tfacc-key-pair",
					"file_system_mount_enabled": "false",
					"zone_id":                   "cn-wulanchabu-d",
					"vswitch_zone_id":           "cn-wulanchabu-d",
					"user_data":                 "YWxpLGFsaSxhbGliYWJh",
					"ignore_failed_node_tasks":  "true",
					"vpd_subnets":               []string{"test"},
					"system_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "180",
							"performance_level": "PL1",
						},
					},
					"ip_allocation_policy": []map[string]interface{}{
						{
							"bond_policy": []map[string]interface{}{
								{
									"bond_default_subnet": "test",
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
								},
							},
							"machine_type_policy": []map[string]interface{}{
								{
									"machine_type": "test",
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
								},
							},
							"node_policy": []map[string]interface{}{
								{
									"node_id": "e01-cn-rno46i6rdfn",
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_name":                 name,
						"node_group_description":          "resource-test1",
						"cluster_id":                      acctest.CHECKSET,
						"machine_type":                    "efg2.GN9De.cn8",
						"az":                              "cn-wulanchabu-d",
						"image_id":                        "i195757131775635040498",
						"key_pair_name":                   "tfacc-key-pair",
						"file_system_mount_enabled":       "false",
						"zone_id":                         "cn-wulanchabu-d",
						"vswitch_zone_id":                 "cn-wulanchabu-d",
						"user_data":                       "YWxpLGFsaSxhbGliYWJh",
						"ignore_failed_node_tasks":        "true",
						"vpd_subnets.#":                   "1",
						"ip_allocation_policy.#":          "1",
						"system_disk.#":                   "1",
						"system_disk.0.category":          "cloud_essd",
						"system_disk.0.size":              "180",
						"system_disk.0.performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_group_name":           name + "_update",
					"image_id":                  "i195048661660874208657",
					"file_system_mount_enabled": "true",
					"system_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "200",
							"performance_level": "PL0",
						},
					},
					"nodes": []map[string]interface{}{
						{
							"node_id":        "e01-cn-rno46i6rdfn",
							"hostname":       "tfacc-hostname",
							"login_password": "YourPassword123!",
							"vpc_id":         "${alicloud_vpc.create_vpc.id}",
							"vswitch_id":     "${alicloud_vswitch.create_vswitch.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_name":                 name + "_update",
						"image_id":                        "i195048661660874208657",
						"file_system_mount_enabled":       "true",
						"system_disk.0.size":              "200",
						"system_disk.0.performance_level": "PL0",
						"nodes.#":                         "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_failed_node_tasks", "ip_allocation_policy", "user_data", "vswitch_zone_id", "vpd_subnets", "zone_id", "login_password", "key_pair_name", "nodes.0.login_password"},
			},
		},
	})
}

// TestAccAliCloudEfloNodeGroup_basic10344_twin covers the node-group lifecycle
// only: it never sets ip_allocation_policy.node_policy.node_id nor nodes, so no
// reserved physical Lingjun node is attached (no ExtendCluster). That keeps it
// re-entrant and far cheaper to run than _basic10344, which carries those
// hardware-bound attributes for schema coverage. The update step changes ONLY
// system_disk.size 180 -> 200 (plus performance_level and name); it deliberately
// keeps image_id and file_system_mount_enabled at their step-1 values because
// UpdateNodeGroup currently rejects both: switching to a non-system image
// surfaces as a 500, and enabling file_system_mount_enabled is gated off in the
// backend (417 "not open to current user"). Read writes the real remote disk size
// back into state, so the size assertion below is the empirical check that
// UpdateNodeGroup honors an in-place Size expansion.
func TestAccAliCloudEfloNodeGroup_basic10344_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupMap10344)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupBasicDependence10344)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_group_name":           name,
					"node_group_description":    "resource-test1",
					"cluster_id":                "${alicloud_eflo_cluster.default.id}",
					"machine_type":              "efg2.GN9De.cn8",
					"az":                        "cn-wulanchabu-d",
					"image_id":                  "i195757131775635040498",
					"login_password":            "YourPassword123!",
					"file_system_mount_enabled": "false",
					"zone_id":                   "cn-wulanchabu-d",
					"vswitch_zone_id":           "cn-wulanchabu-d",
					"user_data":                 "YWxpLGFsaSxhbGliYWJh",
					"ignore_failed_node_tasks":  "true",
					"vpd_subnets":               []string{"test"},
					"system_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "180",
							"performance_level": "PL1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_name":                 name,
						"node_group_description":          "resource-test1",
						"cluster_id":                      acctest.CHECKSET,
						"machine_type":                    "efg2.GN9De.cn8",
						"az":                              "cn-wulanchabu-d",
						"image_id":                        "i195757131775635040498",
						"file_system_mount_enabled":       "false",
						"zone_id":                         "cn-wulanchabu-d",
						"vswitch_zone_id":                 "cn-wulanchabu-d",
						"user_data":                       "YWxpLGFsaSxhbGliYWJh",
						"ignore_failed_node_tasks":        "true",
						"vpd_subnets.#":                   "1",
						"system_disk.#":                   "1",
						"system_disk.0.category":          "cloud_essd",
						"system_disk.0.size":              "180",
						"system_disk.0.performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_group_name": name + "_update",
					"system_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "200",
							"performance_level": "PL0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_name":                 name + "_update",
						"system_disk.0.size":              "200",
						"system_disk.0.performance_level": "PL0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_failed_node_tasks", "user_data", "vswitch_zone_id", "vpd_subnets", "zone_id", "login_password"},
			},
		},
	})
}

var AlicloudEfloNodeGroupMap10344 = map[string]string{
	"create_time":   acctest.CHECKSET,
	"node_group_id": acctest.CHECKSET,
}

func AlicloudEfloNodeGroupBasicDependence10344(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "cluster-resoure-test"
}

resource "alicloud_vswitch" "create_vswitch" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = "cn-wulanchabu-d"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "cluster-resoure-test"
}

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "cluster-resoure-test"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.create_vpc.id
}

resource "alicloud_eflo_cluster" "default" {
  cluster_description      = "cluster-resource-test"
  open_eni_jumbo_frame     = "false"
  hpn_zone                 = "D3"
  nimiz_vswitches          = ["1111"]
  ignore_failed_node_tasks = "true"
  cluster_name             = var.name
  cluster_type             = "Lite"

  node_groups {
    node_group_name        = "cluster-resource-test"
    node_group_description = "cluster-resource-test"
    machine_type           = "efg2.GN9De.cn8"
    image_id               = "i195757131775635040498"
    zone_id                = "cn-wulanchabu-d"
    login_password         = "YourPassword123!"
    system_disk {
      category          = "cloud_essd"
      size              = 180
      performance_level = "PL1"
    }
  }

  networks {
    tail_ip_version   = "ipv4"
    security_group_id = alicloud_security_group.create_security_group.id
    vswitch_zone_id   = "cn-wulanchabu-d"
    vpc_id            = alicloud_vpc.create_vpc.id
    vswitch_id        = alicloud_vswitch.create_vswitch.id

    new_vpd_info {
      monitor_vpc_id     = alicloud_vpc.create_vpc.id
      monitor_vswitch_id = alicloud_vswitch.create_vswitch.id
      vpd_cidr           = "10.0.0.0/16"
      cloud_link_cidr    = "169.254.128.0/23"
      vpd_subnets {
        zone_id     = "cn-wulanchabu-d"
        subnet_cidr = "10.0.1.0/24"
      }
    }
  }
}
`, name)
}
