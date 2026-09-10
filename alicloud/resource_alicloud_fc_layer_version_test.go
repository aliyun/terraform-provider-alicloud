package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudFCLayerVersion_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fc_layer_version.default"
	ra := resourceAttrInit(resourceId, AlicloudFcLayerVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &FcOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcLayerVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sfclayerversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcLayerVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
