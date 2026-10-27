package aligreen_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	aligreen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/aligreen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
// lintignore: AT001
func TestAccAliCloudAligreenOssStockTask_basic7310_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_oss_stock_task.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenOssStockTaskMap7310)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenOssStockTask")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenossstocktask%d", acctest.DefaultRegionToTest, rand)
	cnLoc := time.FixedZone("CST", 8*3600)
	startDate := time.Now().In(cnLoc).AddDate(-1, 0, 0).Format("2006-01-02 00:00:00 +0800")
	endDate := time.Now().In(cnLoc).AddDate(1, 0, 0).Format("2006-01-02 00:00:00 +0800")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenOssStockTaskBasicDependence7310)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
					"start_date":               startDate,
					"end_date":                 endDate,
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
						"start_date":                         acctest.CHECKSET,
						"end_date":                           acctest.CHECKSET,
						"buckets":                            acctest.CHECKSET,
						"image_scenes.#":                     "1",
						"audio_antispam_freeze_config":       "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"image_live_freeze_config":           "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"video_terrorism_freeze_config":      "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"image_terrorism_freeze_config":      "{\"Type\":\"suggestion\",\"Value\":\"block\"}",
						"callback_id":                        acctest.CHECKSET,
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

// Case oss1.0存量检测 collection-order coverage
// Covers the configurable multi-member TypeList attribute image_scenes (ForceNew):
// apply the original member order, then a convergent apply with the members
// reordered. Schema ForceNew triggers a replacement, which the collection-order
// checker allows.
// lintignore: AT001
func TestAccAliCloudAligreenOssStockTask_collectionOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_oss_stock_task.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenOssStockTaskMap7310)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenOssStockTask")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenossstocktask%d", acctest.DefaultRegionToTest, rand)
	cnLoc := time.FixedZone("CST", 8*3600)
	startDate := time.Now().In(cnLoc).AddDate(-1, 0, 0).Format("2006-01-02 00:00:00 +0800")
	endDate := time.Now().In(cnLoc).AddDate(1, 0, 0).Format("2006-01-02 00:00:00 +0800")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenOssStockTaskBasicDependence7310)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
					"start_date":               startDate,
					"end_date":                 endDate,
					"buckets":                  "[{\\\"Bucket\\\":\\\"${alicloud_oss_bucket.defaultPyhXOV.bucket}\\\",\\\"Prefixes\\\":[],\\\"Selected\\\":true}]",
					"image_scenes": []string{
						"porn", "ad"},
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
						"image_scenes.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_scenes": []string{
						"ad", "porn"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_scenes.#": "2",
					}),
				),
			},
		},
	})
}

// Test Aligreen OssStockTask. <<< Resource test cases, automatically generated.
