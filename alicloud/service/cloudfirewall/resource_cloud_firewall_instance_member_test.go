package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudFirewallInstanceMember_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_member.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallInstanceMemberMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallInstanceMember")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCFInstanceMember%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallInstanceMemberBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// currently, international test account has not enabled RD
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"member_desc": "${var.name}",
					"member_uid":  "${alicloud_resource_manager_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"member_desc": name,
						"member_uid":  acctest.CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"member_desc": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"member_desc": name + "_update",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCloudFirewallInstanceMemberMap = map[string]string{}

func AlicloudCloudFirewallInstanceMemberBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_resource_manager_account" "default" {
  display_name = var.name
  abandon_able_check_id = ["SP_fc_fc"]
}
`, name)
}
