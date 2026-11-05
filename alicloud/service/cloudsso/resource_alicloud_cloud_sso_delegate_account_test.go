// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudsso_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudsso "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudSSO DelegateAccount. >>> Resource test cases, automatically generated.
// Case DelegateAccount 11151
func TestAccAliCloudCloudSSODelegateAccount_basic11151(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_delegate_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudSSODelegateAccountMap11151)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSODelegateAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudSSODelegateAccountBasicDependence11151)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id": "${alicloud_resource_manager_delegated_administrator.default.account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id": acctest.CHECKSET,
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

var AlicloudCloudSSODelegateAccountMap11151 = map[string]string{}

func AlicloudCloudSSODelegateAccountBasicDependence11151(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
resource "alicloud_resource_manager_delegated_administrator" "default" {
	account_id = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
	service_principal = "cloudsso.aliyuncs.com"
}

`, name)
}

// Test CloudSSO DelegateAccount. <<< Resource test cases, automatically generated.
