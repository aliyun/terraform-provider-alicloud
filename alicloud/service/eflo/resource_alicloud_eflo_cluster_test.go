package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass nimiz_vswitches -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass node_groups -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass node_groups.nodes -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass components -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass components.component_config.node_units -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy.bond_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy.machine_type_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy.machine_type_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy.node_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.ip_allocation_policy.node_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.new_vpd_info.vpd_subnets -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass networks.vpd_info.vpd_subnets -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEfloCluster_basic10311(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloClusterMap10311)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloClusterBasicDependence10311)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"hpn_zone": "B1",
					"networks": []map[string]interface{}{
						{
							"new_vpd_info": []map[string]interface{}{
								{
									"cloud_link_cidr":    "169.254.128.0/23",
									"monitor_vpc_id":     "${alicloud_vpc.create_vpc.id}",
									"monitor_vswitch_id": "${alicloud_vswitch.create_vswitch.id}",
									"cen_id":             "11111",
									"cloud_link_id":      "1111",
									"vpd_cidr":           "111",
									"vpd_subnets": []map[string]interface{}{
										{
											"subnet_cidr": "111",
											"subnet_type": "111",
											"zone_id":     "1111",
										},
									},
								},
							},
							"security_group_id": "${alicloud_security_group.create_security_group.id}",
							"vswitch_zone_id":   "cn-wulanchabu-b",
							"vpc_id":            "${alicloud_vpc.create_vpc.id}",
							"vswitch_id":        "${alicloud_vswitch.create_vswitch.id}",
							"vpd_info": []map[string]interface{}{
								{
									"vpd_id": "111",
									"vpd_subnets": []string{
										"111"},
								},
							},
							"ip_allocation_policy": []map[string]interface{}{
								{
									"bond_policy": []map[string]interface{}{
										{
											"bond_default_subnet": "111",
											"bonds": []map[string]interface{}{
												{
													"name":   "111",
													"subnet": "111",
												},
											},
										},
									},
									"machine_type_policy": []map[string]interface{}{
										{
											"bonds": []map[string]interface{}{
												{
													"name":   "111",
													"subnet": "111",
												},
											},
											"machine_type": "111",
										},
									},
									"node_policy": []map[string]interface{}{
										{
											"bonds": []map[string]interface{}{
												{
													"name":   "111",
													"subnet": "111",
												},
											},
											"node_id": "111",
										},
									},
								},
							},
							"tail_ip_version": "ipv4",
						},
					},
					"ignore_failed_node_tasks": "true",
					"cluster_type":             "Lite",
					"cluster_name":             name,
					"cluster_description":      "cluster-resource-test",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"node_groups": []map[string]interface{}{
						{
							"node_group_name":        "cluster-resource-test",
							"node_group_description": "cluster-resource-test",
							"machine_type":           "efg1.nvga1n",
							"image_id":               "i190982651690986913088",
							"zone_id":                "cn-wulanchabu-b",
						},
					},
					"nimiz_vswitches": []string{
						"1111"},
					"open_eni_jumbo_frame": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hpn_zone":                 "B1",
						"ignore_failed_node_tasks": "true",
						"cluster_type":             "Lite",
						"cluster_name":             name,
						"cluster_description":      "cluster-resource-test",
						"resource_group_id":        acctest.CHECKSET,
						"node_groups.#":            "1",
						"nimiz_vswitches.#":        "1",
						"open_eni_jumbo_frame":     "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"components", "hpn_zone", "ignore_failed_node_tasks", "networks", "nimiz_vswitches", "node_groups", "open_eni_jumbo_frame"},
			},
		},
	})
}

var AlicloudEfloClusterMap10311 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEfloClusterBasicDependence10311(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "cluster-resoure-test"
}

resource "alicloud_vswitch" "create_vswitch" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = "cn-wulanchabu-b"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "cluster-resoure-test"
}

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "cluster-resoure-test"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.create_vpc.id
}


`, name)
}

// Case 创建集群 1874
func TestAccAliCloudEfloCluster_basic1874(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloClusterMap1874)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloClusterBasicDependence1874)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"components", "hpn_zone", "ignore_failed_node_tasks", "networks", "nimiz_vswitches", "node_groups", "open_eni_jumbo_frame"},
			},
		},
	})
}

var AlicloudEfloClusterMap1874 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEfloClusterBasicDependence1874(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
