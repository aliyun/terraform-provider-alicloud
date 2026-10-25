// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo NodeGroupAttachment. >>> Resource test cases, automatically generated.
// Case 本地盘场景Attachment完整生命周期验证 11092
func TestAccAliCloudEfloNodeGroupAttachment_basic11092(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMap11092)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNodeGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeGroupAttachmentBasicDependence11092)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":        "vsw-uf63gbmvwgreao66opmie",
					"hostname":          "attachment-test-e01-cn-smw4d1bzd0a",
					"login_password":    "G7f$2kL9@vQx3Zp5*",
					"cluster_id":        "i118976621753269898628",
					"node_group_id":     "i127582271753269898630",
					"node_id":           "e01-cn-smw4d1bzd0a",
					"vpc_id":            "vpc-uf6t73bb01dfprb2qvpqa",
					"security_group_id": "${alicloud_security_group.create_security_group.id}",
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
				ImportStateVerifyIgnore: []string{"data_disk", "login_password", "user_data", "security_group_id"},
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

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "ng-attachment-test-11092"
  security_group_type = "normal"
  vpc_id              = "vpc-uf6t73bb01dfprb2qvpqa"
}

`, name)
}

// Case 云盘系统盘Attachment完整生命周期验证 11072
func TestAccAliCloudEfloNodeGroupAttachment_basic11072(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMap11072)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNodeGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeGroupAttachmentBasicDependence11072)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"vswitch_id":        "vsw-2zewbiowk4qgd7cbcj3gj",
					"hostname":          "attachment-test-e01-cn-62l4ccwfx2i",
					"login_password":    "G7f$2kL9@vQx3Zp5*",
					"cluster_id":        "i116429021752634788999",
					"node_group_id":     "i123439231752634789001",
					"node_id":           "e01-cn-oo04e28330o",
					"vpc_id":            "vpc-2ze9znsls3gm3xr17mnrv",
					"user_data":         "ZWNobyAiaGVsbG8i",
					"security_group_id": "${alicloud_security_group.create_security_group.id}",
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
				ImportStateVerifyIgnore: []string{"data_disk", "login_password", "user_data", "security_group_id"},
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

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "ng-attachment-test-11072"
  security_group_type = "normal"
  vpc_id              = "vpc-2ze9znsls3gm3xr17mnrv"
}
`, name)
}

// Case data_disk TypeList 收敛 reorder 覆盖 (TypeListOrderCoverage)
// 见 scripts/collection-order/README.md: apply 两个 distinct members 后,
// 重排顺序再 apply, post-apply plan 必须为空, 不设 ExpectNonEmptyPlan.
func TestAccAliCloudEfloNodeGroupAttachment_dataDiskOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeGroupAttachmentMapDataDiskOrder)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNodeGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	_ = rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeGroupAttachmentDataDiskOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
							"delete_with_node":  "true",
							"category":          "cloud_essd",
							"size":              "100",
							"performance_level": "PL1",
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
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"delete_with_node":  "true",
							"category":          "cloud_essd",
							"size":              "100",
							"performance_level": "PL1",
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
			},
		},
	})
}

var AlicloudEfloNodeGroupAttachmentMapDataDiskOrder = map[string]string{}

func AlicloudEfloNodeGroupAttachmentDataDiskOrderDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Eflo NodeGroupAttachment. <<< Resource test cases, automatically generated.
