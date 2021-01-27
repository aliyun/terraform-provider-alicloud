package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// The quota product does not support deletion, so skip the test.
func SkipTestAccAlicloudQuotasApplicationInfo_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_application_info.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasApplicationInfoMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasApplicationInfo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudQuotasApplicationInfoBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"notice_type":       "0",
					"desire_value":      "60",
					"product_code":      "ess",
					"quota_action_code": "q_db_instance",
					"reason":            "For Terraform Test",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notice_type":       "0",
						"desire_value":      "60",
						"product_code":      "ess",
						"quota_action_code": "q_db_instance",
						"reason":            "For Terraform Test",
						"dimensions.#":      "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"quota_category"},
			},
		},
	})
}

var AlicloudQuotasApplicationInfoMap = map[string]string{
	"notice_type": "0",
	"status":      CHECKSET,
}

func AlicloudQuotasApplicationInfoBasicDependence(name string) string {
	return ""
}
