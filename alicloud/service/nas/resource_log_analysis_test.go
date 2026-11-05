package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Nas LogAnalysis. >>> Resource test cases.
// Basic acceptance test exercising create/read/destroy of the NAS LogAnalysis
// resource. CreateLogAnalysis provisions log delivery on an existing NAS file
// system; the SLS project/logstore/role are managed service-side and exposed as
// computed attributes.

func TestAccAliCloudNasLogAnalysis_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_log_analysis.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasLogAnalysisMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasLogAnalysis")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasloganalysis%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasLogAnalysisBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_nas_file_system.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudNasLogAnalysisMap = map[string]string{
	"file_system_id": acctest.CHECKSET,
	"logstore":       acctest.CHECKSET,
	"project":        acctest.CHECKSET,
	"region":         acctest.CHECKSET,
	"role_arn":       acctest.CHECKSET,
}

func AliCloudNasLogAnalysisBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}
`, name)
}

// Test Nas LogAnalysis. <<< Resource test cases.
