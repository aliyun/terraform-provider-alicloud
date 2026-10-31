package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudMonitorService BasicPublic. >>> Resource test cases, automatically generated.
// Case 5537
func TestAccAliCloudCloudMonitorServiceBasicPublic_basic5537(t *testing.T) {
	t.Skipf("Skipping alicloud_cloud_monitor_service_basic_public testing.")
	t.Skipped()
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_basic_public.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudMonitorServiceBasicPublicMap5537)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceBasicPublic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicebasicpublic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudMonitorServiceBasicPublicBasicDependence5537)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

var AlicloudCloudMonitorServiceBasicPublicMap5537 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudCloudMonitorServiceBasicPublicBasicDependence5537(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudMonitorService BasicPublic. <<< Resource test cases, automatically generated.
