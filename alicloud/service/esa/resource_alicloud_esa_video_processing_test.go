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

// Test ESA VideoProcessing. >>> Resource test cases, automatically generated.
// Case videoprocessing_test
func TestAccAliCloudESAVideoProcessingvideoprocessing_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_video_processing.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAVideoProcessingvideoprocessing_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaVideoProcessing")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAVideoProcessing%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAVideoProcessingvideoprocessing_testBasicDependence)
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
					"video_seek_enable":   "on",
					"site_id":             "${alicloud_esa_site.default.id}",
					"rule_enable":         "on",
					"flv_video_seek_mode": "by_byte",
					"mp4_seek_end":        "end",
					"flv_seek_start":      "start",
					"rule":                "(http.host eq \\\"video.example.com\\\")",
					"sequence":            "1",
					"mp4_seek_start":      "start",
					"site_version":        "0",
					"flv_seek_end":        "end",
					"rule_name":           "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "test_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"video_seek_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flv_video_seek_mode": "by_time",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mp4_seek_start": "s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mp4_seek_end": "e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flv_seek_start": "s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flv_seek_end": "e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
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

var AliCloudESAVideoProcessingvideoprocessing_testMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAVideoProcessingvideoprocessing_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA VideoProcessing. <<< Resource test cases, automatically generated.
