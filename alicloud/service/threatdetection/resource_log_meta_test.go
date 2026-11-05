package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudThreatDetectionLogMeta_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_log_meta.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionLogMetaMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionLogMeta")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionLogMetaBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "disabled",
					"log_meta_name": "aegis-log-client",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":        "disabled",
						"log_meta_name": "aegis-log-client",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "enabled",
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

var AlicloudThreatDetectionLogMetaMap0 = map[string]string{}

func AlicloudThreatDetectionLogMetaBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
