package hbr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbr"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Hbr CrossAccount. >>> Resource test cases, automatically generated.
// Case CrossAccount验证 9417
func TestAccAliCloudHbrCrossAccount_basic9417(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_cross_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHbrCrossAccountMap9417)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrCrossAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrcrossaccount%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHbrCrossAccountBasicDependence9417)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-guangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"create_time": acctest.CHECKSET,
}

func AlicloudHbrCrossAccountBasicDependence9417(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Hbr CrossAccount. <<< Resource test cases, automatically generated.
