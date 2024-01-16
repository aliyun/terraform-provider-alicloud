package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudMonitorService NaamPublic. >>> Resource test cases, automatically generated.
// Case 5535
func TestAccAliCloudCloudMonitorServiceNaamPublic_basic5535(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_naam_public.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceNaamPublicMap5535)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceNaamPublic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicenaampublic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceNaamPublicBasicDependence5535)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudCloudMonitorServiceNaamPublicMap5535 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudCloudMonitorServiceNaamPublicBasicDependence5535(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudMonitorService NaamPublic. <<< Resource test cases, automatically generated.
