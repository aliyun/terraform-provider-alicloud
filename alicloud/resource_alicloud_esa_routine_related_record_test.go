package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test ESA RoutineRelatedRecord. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaRoutineRelatedRecord_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_related_record.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaRoutineRelatedRecordMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRelatedRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serrr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaRoutineRelatedRecordBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"name":        "${alicloud_esa_routine.default.id}",
					"record_name": name + ".${data.alicloud_esa_sites.default.sites.0.site_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     CHECKSET,
						"name":        CHECKSET,
						"record_name": CHECKSET,
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

var AliCloudEsaRoutineRelatedRecordMap0 = map[string]string{
	"record_id": CHECKSET,
	"site_name": CHECKSET,
}

func AliCloudEsaRoutineRelatedRecordBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "tftestacc.com"
}

resource "alicloud_esa_routine" "default" {
  name = var.name
}
`, name)
}

// Test ESA RoutineRelatedRecord. <<< Resource test cases, automatically generated.
