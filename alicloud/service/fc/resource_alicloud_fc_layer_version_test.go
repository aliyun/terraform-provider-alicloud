package fc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	fcsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func SkipTestAccAlicloudFCLayerVersion_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fc_layer_version.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcLayerVersionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcsvc.NewFcOpenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFcLayerVersion")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sfclayerversion%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcLayerVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"layer_name":         name,
					"oss_bucket_name":    "tf-testacc-cdn-3392079",
					"oss_object_name":    "terraform.zip",
					"compatible_runtime": []string{"nodejs14"},
					"description":        name,
					"skip_destroy":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layer_name":      name,
						"oss_bucket_name": "tf-testacc-cdn-3392079",
						"oss_object_name": "terraform.zip",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oss_bucket_name", "oss_object_name", "zip_file", "skip_destroy"},
			},
		},
	})
}

var AlicloudFcLayerVersionMap0 = map[string]string{}

func AlicloudFcLayerVersionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
