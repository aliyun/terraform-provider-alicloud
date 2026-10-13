package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Nas LogAnalysis. >>> Resource test cases.
// Basic acceptance test exercising create/read/destroy of the NAS LogAnalysis
// resource. CreateLogAnalysis provisions log delivery on an existing NAS file
// system; the SLS project/logstore/role are managed service-side and exposed as
// computed attributes.

func TestAccAliCloudNasLogAnalysis_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_log_analysis.default"
	ra := resourceAttrInit(resourceId, AliCloudNasLogAnalysisMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasLogAnalysis")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasloganalysis%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasLogAnalysisBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_nas_file_system.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudNasLogAnalysisMap = map[string]string{
	"file_system_id": CHECKSET,
	"logstore":       CHECKSET,
	"project":        CHECKSET,
	"region":         CHECKSET,
	"role_arn":       CHECKSET,
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
