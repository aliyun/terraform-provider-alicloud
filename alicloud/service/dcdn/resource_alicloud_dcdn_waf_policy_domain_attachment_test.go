package dcdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDcdnWafPolicyDomainAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_policy_domain_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNWafPolicyDomainAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnWafPolicyDomainAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccdcdnwafpolicydomainattachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNWafPolicyDomainAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${alicloud_dcdn_waf_domain.default.domain_name}",
					"policy_id":   "${alicloud_dcdn_waf_policy.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": acctest.CHECKSET,
						"policy_id":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDCDNWafPolicyDomainAttachmentMap0 = map[string]string{
	"policy_id": acctest.CHECKSET,
}

func AlicloudDCDNWafPolicyDomainAttachmentBasicDependence0(name string) string {
	rand := tfacctest.RandIntRange(10000, 99999)
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "domain_name" {	
	default = "tf-testacc%d.alicloud-provider.cn"
}
resource "alicloud_dcdn_domain" "default" {
  domain_name = var.domain_name
  sources {
    content = "1.1.1.1"
    port = "80"
    priority = "20"
    type = "ipaddr"
  }
}
resource "alicloud_dcdn_waf_domain" "default" {
	domain_name = alicloud_dcdn_domain.default.domain_name
	client_ip_tag = "X-Forwarded-For"
}
resource "alicloud_dcdn_waf_policy" "default" {
	policy_type  =   "custom"
	policy_name  =    var.name
	defense_scene = "waf_group"
	status =    "on"
}
`, name, rand)
}
