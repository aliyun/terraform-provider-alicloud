package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oos"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOOSDefaultPatchBaseline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_default_patch_baseline.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOOSDefaultPatchBaselineMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosDefaultPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccOOSDefaultPatchBaseline%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOOSDefaultPatchBaselineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"patch_baseline_name": "${alicloud_oos_patch_baseline.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"patch_baseline_name": acctest.CHECKSET,
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

var AlicloudOOSDefaultPatchBaselineMap0 = map[string]string{
	"patch_baseline_name": acctest.CHECKSET,
}

func AlicloudOOSDefaultPatchBaselineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_oos_patch_baseline" "default" {
  operation_system    = "Windows"
  patch_baseline_name = var.name
  description         = var.name
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
}
`, name)
}
