// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudSSO DelegateAccount. >>> Resource test cases, automatically generated.
// Case DelegateAccount 11151
func TestAccAliCloudCloudSSODelegateAccount_basic11151(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_delegate_account.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudSSODelegateAccountMap11151)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSODelegateAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSODelegateAccountBasicDependence11151)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id": "${alicloud_resource_manager_delegated_administrator.default.account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id": CHECKSET,
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
