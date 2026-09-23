package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Live Caster. >>> Resource test cases, automatically generated.
// Case 导播测试_Terraform 9293
func TestAccAliCloudLiveCaster_basic9293(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_live_caster.default"
	ra := resourceAttrInit(resourceId, AlicloudLiveCasterMap9293)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LiveServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLiveCaster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLiveCasterBasicDependence9293)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"caster_name":       name,
					"payment_type":      "PayAsYouGo",
					"norm_type":         "1",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"caster_name":       name,
						"payment_type":      "PayAsYouGo",
						"norm_type":         "1",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"caster_name":               name + "_update",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"resource_type":             "caster",
					"domain_name":               "guoxinplay.alivecdn.com",
					"urgent_image_url":          "http://learn.aliyundoc.com/AppName/image.jpg",
					"auto_switch_urgent_config": "{\\\"eofThres\\\":3}",
					"urgent_material_id":        "a2b8e671",
					"transcode_config":          "{\\\"casterTemplate\\\": \\\"lp_ld\\\"}",
					"program_name":              "program_name",
					"urgent_image_id":           "a089175eb5f4427684fc0715159a",
					"delay":                     "0",
					"program_effect":            "1",
					"auto_switch_urgent_on":     "false",
					"side_output_url_list":      "[\\\"rtmp://antang-test-del04.alivecdn.com/app1/stream1\\\",\\\"rtmp://antang-test-del04.alivecdn.com/app2/stream2\\\"]",
					"side_output_url":           "rtmp://antang-test-del04.alivecdn.com/app1/stream1",
					"callback_url":              "http://aliyundoc.com:8000/caster/4a82a3d1b7f0462ea37348366201",
					"record_config":             "{\\\"endpoint\\\":\\\"http://oss-cn-shanghai.aliyuncs.com/api\\\",\\\"ossBucket\\\":\\\"liveBucketabcd\\\",\\\"VideoFormat\\\":[{\\\"OssObjectPrefix\\\":\\\"record/apptest/streamtest/1733396815_1733396875\\\",\\\"Format\\\":\\\"m3u8\\\",\\\"CycleDuration\\\":21600,\\\"SliceOssObjectPrefix\\\":\\\"record/apptest/streamtest-s/1733396815\\\"},{\\\"OssObjectPrefix\\\":\\\"record/apptest2/streamtest2/1733396815_1733396875\\\",\\\"Format\\\":\\\"flv\\\",\\\"CycleDuration\\\":21600}],\\\"interval\\\":5}",
					"sync_groups_config":        "[{\\\"mode\\\":1,\\\"resourceIds\\\":[],\\\"hostResourceId\\\":\\\"3aa2b39a-fd0e-4b8c-be73-b7af31c4****\\\"}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"caster_name":               name + "_update",
						"resource_group_id":         CHECKSET,
						"resource_type":             "caster",
						"domain_name":               "guoxinplay.alivecdn.com",
						"urgent_image_url":          "http://learn.aliyundoc.com/AppName/image.jpg",
						"auto_switch_urgent_config": "{\"eofThres\":3}",
						"urgent_material_id":        "a2b8e671",
						"transcode_config":          "{\"casterTemplate\": \"lp_ld\"}",
						"program_name":              "program_name",
						"urgent_image_id":           "a089175eb5f4427684fc0715159a",
						"delay":                     "0",
						"program_effect":            "1",
						"auto_switch_urgent_on":     "false",
						"side_output_url_list":      "[\"rtmp://antang-test-del04.alivecdn.com/app1/stream1\",\"rtmp://antang-test-del04.alivecdn.com/app2/stream2\"]",
						"side_output_url":           "rtmp://antang-test-del04.alivecdn.com/app1/stream1",
						"callback_url":              "http://aliyundoc.com:8000/caster/4a82a3d1b7f0462ea37348366201",
						"record_config":             "{\"endpoint\":\"http://oss-cn-shanghai.aliyuncs.com/api\",\"ossBucket\":\"liveBucketabcd\",\"VideoFormat\":[{\"OssObjectPrefix\":\"record/apptest/streamtest/1733396815_1733396875\",\"Format\":\"m3u8\",\"CycleDuration\":21600,\"SliceOssObjectPrefix\":\"record/apptest/streamtest-s/1733396815\"},{\"OssObjectPrefix\":\"record/apptest2/streamtest2/1733396815_1733396875\",\"Format\":\"flv\",\"CycleDuration\":21600}],\"interval\":5}",
						"sync_groups_config":        "[{\"mode\":1,\"resourceIds\":[],\"hostResourceId\":\"3aa2b39a-fd0e-4b8c-be73-b7af31c4****\"}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"urgent_live_stream_url": "rtmp://demo.aliyundoc.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"urgent_live_stream_url": "rtmp://demo.aliyundoc.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_switch_urgent_config", "auto_switch_urgent_on", "callback_url", "delay", "domain_name", "program_effect", "program_name", "record_config", "resource_type", "side_output_url", "side_output_url_list", "sync_groups_config", "transcode_config", "urgent_image_id", "urgent_image_url", "urgent_live_stream_url", "urgent_material_id"},
			},
		},
	})
}

var AlicloudLiveCasterMap9293 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudLiveCasterBasicDependence9293(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Live Caster. <<< Resource test cases, automatically generated.
