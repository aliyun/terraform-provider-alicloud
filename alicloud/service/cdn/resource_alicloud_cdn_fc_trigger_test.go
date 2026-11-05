package cdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCDNFcTrigger_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_fc_trigger.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCDNFcTriggerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cdn.NewCdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCdnFcTrigger")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scdnfctrigger%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCDNFcTriggerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_meta_name":    "LogFileCreated",
					"event_meta_version": "1.0.0",
					"notes":              "${var.name}",
					"trigger_arn":        "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:services/FCTestService/functions/printEvent/triggers/testtrigger",
					"source_arn":         "acs:cdn:*:${data.alicloud_account.default.id}:domain/example.com",
					"role_arn":           "acs:ram::${data.alicloud_account.default.id}:role/aliyuncdneventnotificationrole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_meta_name":    "LogFileCreated",
						"event_meta_version": "1.0.0",
						"notes":              name,
						"trigger_arn":        acctest.CHECKSET,
						"source_arn":         acctest.CHECKSET,
						"role_arn":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notes":        "${var.name}_update",
					"function_arn": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:services/FCTestService/functions/printEvent",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notes": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_arn": "acs:cdn:*:${data.alicloud_account.default.id}:domain/example1.com",
					"role_arn":   "acs:ram::${data.alicloud_account.default.id}:role/aliyuncdneventnotificationrole2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_arn": acctest.CHECKSET,
						"role_arn":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"function_arn"},
			},
		},
	})
}

var AlicloudCDNFcTriggerMap0 = map[string]string{
	"notes":              acctest.CHECKSET,
	"role_arn":           acctest.CHECKSET,
	"source_arn":         acctest.CHECKSET,
	"trigger_arn":        acctest.CHECKSET,
	"event_meta_name":    acctest.CHECKSET,
	"event_meta_version": acctest.CHECKSET,
}

func AlicloudCDNFcTriggerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {}

data "alicloud_regions" "default" {
  current = true
}
`, name)
}
