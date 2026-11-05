package dmsenterprise_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	dmsenterprise "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudDmsProxyAccess_basic2036(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_proxy_access.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudDmsProxyAccessMap2036)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDmsEnterpriseService(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDmsEnterpriseProxyAccess")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseProxyAccessSupportRegions)
	name := fmt.Sprintf("tf-testacc%sDmsProxyAccess%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudDmsProxyAccessBasicDependence2036)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_id": "${data.alicloud_dms_enterprise_proxies.ids.proxies.0.id}",
					"user_id":  "${data.alicloud_dms_enterprise_users.dms_enterprise_users_ds.users.0.user_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_id":        tfacctest.CHECKSET,
						"user_id":         tfacctest.CHECKSET,
						"instance_id":     tfacctest.CHECKSET,
						"proxy_access_id": tfacctest.CHECKSET,
						"origin_info":     tfacctest.CHECKSET,
						"access_id":       tfacctest.CHECKSET,
						"access_secret":   tfacctest.CHECKSET,
						"create_time":     tfacctest.CHECKSET,
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
