package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudDmsProxyAccess_basic2036(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_proxy_access.default"
	ra := resourceAttrInit(resourceId, AlicloudDmsProxyAccessMap2036)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmsEnterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsEnterpriseProxyAccess")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.DMSEnterpriseProxyAccessSupportRegions)
	name := fmt.Sprintf("tf-testacc%sDmsProxyAccess%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDmsProxyAccessBasicDependence2036)
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
					"proxy_id": "${data.alicloud_dms_enterprise_proxies.ids.proxies.0.id}",
					"user_id":  "${data.alicloud_dms_enterprise_users.dms_enterprise_users_ds.users.0.user_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_id":        CHECKSET,
						"user_id":         CHECKSET,
						"instance_id":     CHECKSET,
						"proxy_access_id": CHECKSET,
						"origin_info":     CHECKSET,
						"access_id":       CHECKSET,
						"access_secret":   CHECKSET,
						"create_time":     CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"indep_password"},
			},
		},
	})
}

var AlicloudDmsProxyAccessMap2036 = map[string]string{}

func AlicloudDmsProxyAccessBasicDependence2036(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_dms_enterprise_users" "dms_enterprise_users_ds" {
  role   = "USER"
  status = "NORMAL"
}
data "alicloud_dms_enterprise_proxies" "ids" {}

`, name)
}
