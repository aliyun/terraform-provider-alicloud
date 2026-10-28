package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case alicloud_ram_access_key_policy basic
func TestAccAliCloudRamAccessKeyPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_access_key_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRamAccessKeyPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamAccessKeyPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccramakp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRamAccessKeyPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_key_id":  "${alicloud_ram_access_key.default.id}",
					"user_principal_name": "${alicloud_ram_user.default.name}@${data.alicloud_account.default.id}.onaliyun.com",
					"access_key_policy":   "{\\\"Version\\\":1,\\\"Status\\\":\\\"Active\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.1/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_key_id":  acctest.CHECKSET,
						"user_principal_name": acctest.CHECKSET,
						"access_key_policy":   acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_key_policy": "{\\\"Version\\\":1,\\\"Status\\\":\\\"Active\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.2/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_key_policy": acctest.CHECKSET,
					}),
				),
			},
			{
				// Disable the policy while retaining the rule (Status Inactive, non-empty Statements).
				// Note: a disabled policy with no statements is the reset baseline and is treated as
				// "not exist", so a statement is retained here to keep the resource present.
				Config: testAccConfig(map[string]interface{}{
					"access_key_policy": "{\\\"Version\\\":1,\\\"Status\\\":\\\"Inactive\\\",\\\"Statements\\\":[{\\\"Type\\\":\\\"ClassicWhiteList\\\",\\\"IPList\\\":[\\\"10.0.0.2/32\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_key_policy": acctest.CHECKSET,
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

var AlicloudRamAccessKeyPolicyMap0 = map[string]string{}

func AlicloudRamAccessKeyPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}
`, name)
}
