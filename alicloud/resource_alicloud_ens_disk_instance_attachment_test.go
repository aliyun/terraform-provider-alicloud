package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens DiskInstanceAttachment. >>> Resource test cases, automatically generated.
// Case 磁盘挂载 5683
func TestAccAliCloudEnsDiskInstanceAttachment_basic5683(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_disk_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsDiskInstanceAttachmentMap5683)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsDiskInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensdiskinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsDiskInstanceAttachmentBasicDependence5683)
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
					"instance_id":          "${alicloud_ens_instance.创建实例.id}",
					"delete_with_instance": "false",
					"disk_id":              "${alicloud_ens_disk.创建磁盘.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":          CHECKSET,
						"delete_with_instance": "false",
						"disk_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_with_instance"},
			},
		},
	})
}

var AlicloudEnsDiskInstanceAttachmentMap5683 = map[string]string{}

func AlicloudEnsDiskInstanceAttachmentBasicDependence5683(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_disk" "创建磁盘" {
  size          = "20"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  payment_type  = "PayAsYouGo"
  category      = "cloud_efficiency"
}

resource "alicloud_ens_instance" "创建实例" {
  system_disk {
    size = "20"
  }
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678ABCabc"
  amount                     = "1"
  internet_max_bandwidth_out = "10"
  unique_suffix              = true
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  schedule_area_level        = "Region"
  period_unit                = "Month"
  period                     = "1"
}


`, name)
}

// Test Ens DiskInstanceAttachment. <<< Resource test cases, automatically generated.
