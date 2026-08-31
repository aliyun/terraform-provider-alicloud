package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Hbr CrossAccount. >>> Resource test cases, automatically generated.
// Case CrossAccount验证 9417
func TestAccAliCloudHbrCrossAccount_basic9417(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_cross_account.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrCrossAccountMap9417)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrCrossAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrcrossaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrCrossAccountBasicDependence9417)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-guangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_account_user_id":   "1",
					"cross_account_role_name": "zhenyuan-时间引用：GetCurrentUnixTimeStamp(0,'ms','s')",
					"alias":                   "镇元测试用例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_account_user_id":   "1",
						"cross_account_role_name": "zhenyuan-时间引用：GetCurrentUnixTimeStamp(0,'ms','s')",
						"alias":                   "镇元测试用例",
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

var AlicloudHbrCrossAccountMap9417 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrCrossAccountBasicDependence9417(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Hbr CrossAccount. <<< Resource test cases, automatically generated.
