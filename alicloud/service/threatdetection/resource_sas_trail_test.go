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

// Test ThreatDetection SasTrail. >>> Resource test cases, automatically generated.
// Case 4975
func TestAccAliCloudThreatDetectionSasTrail_basic4975(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_sas_trail.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionSasTrailMap4975)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionSasTrail")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionsastrail%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionSasTrailBasicDependence4975)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

var AlicloudThreatDetectionSasTrailMap4975 = map[string]string{
	"service_trail.#": "1",
	"create_time":     acctest.CHECKSET,
}

func AlicloudThreatDetectionSasTrailBasicDependence4975(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection SasTrail. <<< Resource test cases, automatically generated.
