package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Skip this testcase because you can only have one instance.
func TestAccAliCloudWafv3Instance_basic2294(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	resourceId := "alicloud_wafv3_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudWafV3InstanceMap2294)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Instance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sWafV3Instance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafV3InstanceBasicDependence2294)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWafv3Instance(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": CHECKSET,
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

// testAccPreCheckWafv3Instance ensures no WAFv3 instance already exists in the test
// region before running a lifecycle test. Each Alibaba Cloud account can only own one
// WAF instance per region, and CreatePostpaidInstance fails with TokenCreateLimited
// when one already exists. Unlike the BSS-based testAccPreCheckForCleanUpInstances
// (which only queries Subscription/prepaid instances and therefore misses postpaid
// PayAsYouGo ones), this uses the authoritative WAF DescribeInstance API so that
// existing instances of either billing type are detected. When an instance exists, it
// waits for it to be cleaned up; if it still exists after the deadline, the test is
// skipped (recording the instance ID and region) instead of proceeding and failing
// with the deterministic 400 TokenCreateLimited error.
func testAccPreCheckWafv3Instance(t *testing.T) {
	rawClient, err := sharedClientForRegion(defaultRegionToTest)
	if err != nil {
		t.Fatalf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}
	deadline := time.Now().Add(10 * time.Minute)
	for {
		object, err := wafOpenapiService.DescribeWafv3Instance("")
		if err != nil {
			if NotFoundError(err) {
				// No WAFv3 instance exists in the region; proceed with the lifecycle test.
				return
			}
			t.Fatalf("error describing WAFv3 instance in region %s: %s", client.RegionId, err)
		}
		instanceId := ""
		if v, ok := object["InstanceId"]; ok {
			instanceId = fmt.Sprint(v)
		}
		if instanceId == "" {
			return
		}
		if time.Now().After(deadline) {
			t.Skipf("Skip test: a WAFv3 instance (%s) already exists in region %s and was not cleaned up within 10m. "+
				"Each Alibaba Cloud account can only own one WAF instance per region; release the existing instance "+
				"(postpaid via ReleaseInstance, prepaid via the console) before re-running.", instanceId, client.RegionId)
		}
		t.Logf("waiting for cleaning up existed WAFv3 instance %s in region %s...", instanceId, client.RegionId)
		time.Sleep(30 * time.Second)
	}
}
