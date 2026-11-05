package wafopenapi_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/wafopenapi"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Skip this testcase because you can only have one instance.
func TestAccAliCloudWafv3Instance_basic2294(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	resourceId := "alicloud_wafv3_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudWafV3InstanceMap2294)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return wafopenapi.NewWafOpenapiService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeWafv3Instance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sWafV3Instance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudWafV3InstanceBasicDependence2294)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheckForCleanUpInstances(t, string(connectivity.Hangzhou), "waf", "waf", "waf", "waf")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": acctest.CHECKSET,
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

var AlicloudWafV3InstanceMap2294 = map[string]string{}

func AlicloudWafV3InstanceBasicDependence2294(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}`, name)
}
