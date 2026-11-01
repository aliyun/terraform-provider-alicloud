package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA RoutineRelatedRecord. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaRoutineRelatedRecord_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_related_record.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRoutineRelatedRecordMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutineRelatedRecord")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serrr%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRoutineRelatedRecordBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"name":        "${alicloud_esa_routine.default.id}",
					"record_name": name + ".${data.alicloud_esa_sites.default.sites.0.site_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     acctest.CHECKSET,
						"name":        acctest.CHECKSET,
						"record_name": acctest.CHECKSET,
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
	"record_id": acctest.CHECKSET,
	"site_name": acctest.CHECKSET,
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
