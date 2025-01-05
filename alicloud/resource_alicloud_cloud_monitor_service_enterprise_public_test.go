package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudMonitorService EnterprisePublic. >>> Resource test cases, automatically generated.
// Case 5536
func TestAccAliCloudCloudMonitorServiceEnterprisePublic_basic5536(t *testing.T) {
	t.Skipf("Skipping alicloud_cloud_monitor_service_enterprise_public testing because of the service limitation that only run once per day.")
	t.Skipped()
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_enterprise_public.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceEnterprisePublicMap5536)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceEnterprisePublic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorserviceenterprisepublic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceEnterprisePublicBasicDependence5536)
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

var AlicloudCloudMonitorServiceEnterprisePublicMap5536 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudCloudMonitorServiceEnterprisePublicBasicDependence5536(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudMonitorService EnterprisePublic. <<< Resource test cases, automatically generated.
