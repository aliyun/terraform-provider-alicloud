package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cdt Service. >>> Resource test cases, automatically generated.
// Case 开通测试_副本1706495831058 5864
func TestAccAliCloudCdtService_basic5864(t *testing.T) {
	var v map[string]interface{}
	resourceId := "data.alicloud_cdt_service.default"
	ra := resourceAttrInit(resourceId, AlicloudCdtServiceMap5864)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CdtServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCdtInternetService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scdtinternetservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCdtServiceBasicDependence5864)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlicloudCdtServiceMap5864 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCdtServiceBasicDependence5864(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Cdt Service. <<< Resource test cases, automatically generated.
