// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass data_disk -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo NodeGroupAttachment. >>> Resource test cases, automatically generated.
// Case 本地盘场景Attachment完整生命周期验证 11092
func TestAccAliCloudEfloNodeGroupAttachment_basic11092(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMap11092)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupAttachmentBasicDependence11092)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":     "vsw-uf63gbmvwgreao66opmie",
					"hostname":       "attachment-test-e01-cn-smw4d1bzd0a",
					"login_password": "G7f$2kL9@vQx3Zp5*",
					"cluster_id":     "i118976621753269898628",
					"node_group_id":  "i127582271753269898630",
					"node_id":        "e01-cn-smw4d1bzd0a",
					"vpc_id":         "vpc-uf6t73bb01dfprb2qvpqa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":     "vsw-uf63gbmvwgreao66opmie",
						"hostname":       "attachment-test-e01-cn-smw4d1bzd0a",
						"login_password": "G7f$2kL9@vQx3Zp5*",
						"cluster_id":     "i118976621753269898628",
						"node_group_id":  "i127582271753269898630",
						"node_id":        "e01-cn-smw4d1bzd0a",
						"vpc_id":         "vpc-uf6t73bb01dfprb2qvpqa",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_disk", "login_password", "user_data"},
			},
		},
	})
}

var AlicloudEfloNodeGroupAttachmentMap11092 = map[string]string{}

func AlicloudEfloNodeGroupAttachmentBasicDependence11092(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 云盘系统盘Attachment完整生命周期验证 11072
func TestAccAliCloudEfloNodeGroupAttachment_basic11072(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMap11072)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupAttachmentBasicDependence11072)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"delete_with_node":  "false",
							"category":          "cloud_essd",
							"size":              "40",
							"performance_level": "PL0",
						},
						{
							"delete_with_node":  "false",
							"category":          "cloud_essd",
							"size":              "40",
							"performance_level": "PL0",
						},
					},
					"vswitch_id":     "vsw-2zewbiowk4qgd7cbcj3gj",
					"hostname":       "attachment-test-e01-cn-62l4ccwfx2i",
					"login_password": "G7f$2kL9@vQx3Zp5*",
					"cluster_id":     "i116429021752634788999",
					"node_group_id":  "i123439231752634789001",
					"node_id":        "e01-cn-oo04e28330o",
					"vpc_id":         "vpc-2ze9znsls3gm3xr17mnrv",
					"user_data":      "ZWNobyAiaGVsbG8i",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":    "2",
						"vswitch_id":     "vsw-2zewbiowk4qgd7cbcj3gj",
						"hostname":       "attachment-test-e01-cn-62l4ccwfx2i",
						"login_password": "G7f$2kL9@vQx3Zp5*",
						"cluster_id":     "i116429021752634788999",
						"node_group_id":  "i123439231752634789001",
						"node_id":        "e01-cn-oo04e28330o",
						"vpc_id":         "vpc-2ze9znsls3gm3xr17mnrv",
						"user_data":      "ZWNobyAiaGVsbG8i",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_disk", "login_password", "user_data"},
			},
		},
	})
}

var AlicloudEfloNodeGroupAttachmentMap11072 = map[string]string{}

func AlicloudEfloNodeGroupAttachmentBasicDependence11072(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_eflo_node_group_attachment" "default0" {
  node_group_id = "i123439231752634789001"
  login_password = "G7f$2kL9@vQx3Zp5*"
  vpc_id = "vpc-2ze9znsls3gm3xr17mnrv"
  hostname = "attachment-test-e01-cn-62l4ccwfx2j"
  cluster_id = "i116429021752634788999"
  node_id = "e01-cn-62l4ccwfx2k"
  data_disk {
    delete_with_node = "false"
    category = "cloud_essd"
    size = "40"
    performance_level = "PL0"
  }
  
  vswitch_id = "vsw-2zewbiowk4qgd7cbcj3gj"
} 
`, name)
}

// Test Eflo NodeGroupAttachment. <<< Resource test cases, automatically generated.

// Case VPD 集群 Attachment 生命周期验证：节点通过 VpdSubnets/IpAllocationPolicy 从集群子网分配 IP，不配置 vpc_id/vswitch_id
//
// 前置条件（一次性环境，非本测试自建）：一个可用的 VPD 集群、其节点组与待挂载节点，以及该集群所属 VPD 的子网和节点机型，
// 分别通过 ALICLOUD_EFLO_VPD_CLUSTER_ID / _NODE_GROUP_ID / _NODE_ID / _SUBNET_ID / _MACHINE_TYPE 传入；
// ALICLOUD_EFLO_VPD_REGION 用于区域前置校验，须与 Provider 实际运行地域一致（缺省 cn-hangzhou）。
// 注意：Destroy 通过 ShrinkCluster 把该节点移出集群，一次执行即消耗该预置节点，重复执行前需由环境责任人补充可用节点。
func TestAccAliCloudEfloNodeGroupAttachment_vpd(t *testing.T) {
	clusterId := os.Getenv("ALICLOUD_EFLO_VPD_CLUSTER_ID")
	nodeGroupId := os.Getenv("ALICLOUD_EFLO_VPD_NODE_GROUP_ID")
	nodeId := os.Getenv("ALICLOUD_EFLO_VPD_NODE_ID")
	subnetId := os.Getenv("ALICLOUD_EFLO_VPD_SUBNET_ID")
	machineType := os.Getenv("ALICLOUD_EFLO_VPD_MACHINE_TYPE")
	region := os.Getenv("ALICLOUD_EFLO_VPD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	if clusterId == "" || nodeGroupId == "" || nodeId == "" || subnetId == "" || machineType == "" {
		t.Skip("Skipping VPD attachment acceptance test: ALICLOUD_EFLO_VPD_CLUSTER_ID, ALICLOUD_EFLO_VPD_NODE_GROUP_ID, ALICLOUD_EFLO_VPD_NODE_ID, ALICLOUD_EFLO_VPD_SUBNET_ID and ALICLOUD_EFLO_VPD_MACHINE_TYPE must point to a pre-provisioned VPD cluster, subnet and available node")
	}

	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMapVpd)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEfloAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupAttachmentBasicDependenceVpd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(region)})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":    clusterId,
					"node_group_id": nodeGroupId,
					"node_id":       nodeId,
					"hostname":      name,
					"vpd_subnets": []string{
						subnetId,
					},
					"ip_allocation_policy": []map[string]interface{}{
						{
							"bond_policy": []map[string]interface{}{
								{
									"bond_default_subnet": subnetId,
									"bonds": []map[string]interface{}{
										{
											"name":   "Bond0",
											"subnet": subnetId,
										},
									},
								},
							},
							"machine_type_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "Bond0",
											"subnet": subnetId,
										},
									},
									"machine_type": machineType,
								},
							},
							"node_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "Bond0",
											"subnet": subnetId,
										},
									},
									"hostname": name,
									"node_id":  nodeId,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":             clusterId,
						"node_group_id":          nodeGroupId,
						"node_id":                nodeId,
						"hostname":               name,
						"vpd_subnets.#":          "1",
						"ip_allocation_policy.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudEfloNodeGroupAttachmentMapVpd = map[string]string{}

func AlicloudEfloNodeGroupAttachmentBasicDependenceVpd(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
