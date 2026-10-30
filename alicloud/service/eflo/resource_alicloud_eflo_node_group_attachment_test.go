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
