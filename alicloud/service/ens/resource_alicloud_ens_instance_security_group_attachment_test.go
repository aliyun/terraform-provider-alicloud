package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ens InstanceSecurityGroupAttachment. >>> Resource test cases, automatically generated.
// Case 实例安全组绑定 5684
func TestAccAliCloudEnsInstanceSecurityGroupAttachment_basic5684(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance_security_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEnsInstanceSecurityGroupAttachmentMap5684)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsInstanceSecurityGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstancesecuritygroupattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEnsInstanceSecurityGroupAttachmentBasicDependence5684)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_ens_instance.创建实例.id}",
					"security_group_id": "${alicloud_ens_security_group.创建安全组.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsInstanceSecurityGroupAttachmentMap5684 = map[string]string{
	"instance_id": acctest.CHECKSET,
}

func AlicloudEnsInstanceSecurityGroupAttachmentBasicDependence5684(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_instance" "创建实例" {
  system_disk {
    size = "20"
  }
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678ABCabc"
  amount                     = "1"
  period                     = "1"
  internet_max_bandwidth_out = "10"
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit                = "Month"
}

resource "alicloud_ens_security_group" "创建安全组" {
  description         = "InstanceSecurityGroupAttachment_Description"
  security_group_name = var.name

}


`, name)
}

// Test Ens InstanceSecurityGroupAttachment. <<< Resource test cases, automatically generated.
