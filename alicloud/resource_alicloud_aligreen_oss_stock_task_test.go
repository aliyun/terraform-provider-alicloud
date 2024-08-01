package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Aligreen OssStockTask. >>> Resource test cases, automatically generated.
var AlicloudAligreenOssStockTaskMap7310 = map[string]string{}

func AlicloudAligreenOssStockTaskBasicDependence7310(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultPyhXOV" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_aligreen_callback" "defaultJnW8Na" {
  callback_url         = "https://www.aliyun.com/"
  crypt_type           = "0"
  callback_name        = var.name
  callback_types       = ["machineScan"]
  callback_suggestions = ["block"]
}


`, name)
}

// Case oss1.0存量检测 7310  raw
func TestAccAliCloudAligreenOssStockTask_basic7310_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_oss_stock_task.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenOssStockTaskMap7310)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenOssStockTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenossstocktask%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenOssStockTaskBasicDependence7310)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_opened":             "true",
					"scan_image_no_file_type":  "true",
					"audio_opened":             "false",
					"image_auto_freeze_opened": "false",
					"auto_freeze_type":         "acl",
					"audio_max_size":           "200",
					"video_opened":             "false",
					"image_scan_limit":         "1",
					"video_frame_interval":     "1",
					"audio_auto_freeze_opened": "false",
					"video_scan_limit":         "1000",
					"video_auto_freeze_opened": "false",
					"audio_scan_limit":         "1000",
					"video_max_frames":         "200",
					"video_max_size":           "500",
					"start_date":               "2024-08-30 00:00:00 +0800",
					"end_date":                 "2024-08-30 20:42:29 +0800",
					"buckets":                  "[{\\\"Bucket\\\":\\\"${alicloud_oss_bucket.defaultPyhXOV.bucket}\\\",\\\"Prefixes\\\":[],\\\"Selected\\\":true}]",
					"image_scenes": []string{
						"porn"},
					"audio_antispam_freeze_config":       "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"image_live_freeze_config":           "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"video_terrorism_freeze_config":      "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"image_terrorism_freeze_config":      "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"callback_id":                        "${alicloud_aligreen_callback.defaultJnW8Na.id}",
					"image_ad_freeze_config":             "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"biz_type":                           "recommend_massmedia_template_01",
					"audio_scenes":                       "[\\\"antispam\\\"]",
					"image_porn_freeze_config":           "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"video_live_freeze_config":           "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"video_porn_freeze_config":           "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"video_voice_antispam_freeze_config": "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
					"video_scenes":                       "[\\\"ad\\\",\\\"terrorism\\\",\\\"live\\\",\\\"porn\\\",\\\"antispam\\\"]",
					"video_ad_freeze_config":             "{\\\"Type\\\":\\\"suggestion\\\",\\\"Value\\\":\\\"block\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_opened":                       "true",
						"scan_image_no_file_type":            "true",
						"audio_opened":                       "false",
						"image_auto_freeze_opened":           "false",
						"auto_freeze_type":                   "acl",
						"audio_max_size":                     "200",
						"video_opened":                       "false",
						"image_scan_limit":                   "1",
						"video_frame_interval":               "1",
						"audio_auto_freeze_opened":           "false",
						"video_scan_limit":                   "1000",
						"video_auto_freeze_opened":           "false",
						"audio_scan_limit":                   "1000",
						"video_max_frames":                   "200",
						"video_max_size":                     "500",
						"start_date":                         CHECKSET,
						"end_date":                           CHECKSET,
						"buckets":                            CHECKSET,
						"image_scenes.#":                     "1",
						"audio_antispam_freeze_config":       "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"image_live_freeze_config":           "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_terrorism_freeze_config":      "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"image_terrorism_freeze_config":      "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"callback_id":                        CHECKSET,
						"image_ad_freeze_config":             "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"biz_type":                           "recommend_massmedia_template_01",
						"audio_scenes":                       "[\"antispam\"]",
						"image_porn_freeze_config":           "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_live_freeze_config":           "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_porn_freeze_config":           "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_voice_antispam_freeze_config": "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_scenes":                       "[\"ad\",\"terrorism\",\"live\",\"porn\",\"antispam\"]",
						"video_ad_freeze_config":             "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"end_date"},
			},
		},
	})
}

// Test Aligreen OssStockTask. <<< Resource test cases, automatically generated.
