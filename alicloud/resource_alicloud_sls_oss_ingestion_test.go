package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls OssIngestion. >>> Resource test cases, automatically generated.
// Case 标准测试用例 5983
func TestAccAliCloudSlsOssIngestion_basic5983(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_oss_ingestion.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsOssIngestionMap5983)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsOssIngestion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsossingestion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsOssIngestionBasicDependence5983)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.defaultProject.name}",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultILogStore.name}",
							"source": []map[string]interface{}{
								{
									"endpoint":          "oss-cn-chengdu.aliyuncs.com",
									"bucket":            "${alicloud_oss_bucket.defaultBucket.bucket_name}",
									"encoding":          "utf-16",
									"format":            "{\\\"type\\\": \\\"text\\\"}",
									"interval":          "3m",
									"pattern":           "1day*",
									"prefix":            "1day",
									"start_time":        "1706792358",
									"end_time":          "1706792558",
									"time_field":        "__time__",
									"time_format":       "epoch",
									"time_pattern":      "\\\\d+:\\\\d+:\\\\d+",
									"time_zone":         "GMT-09:00",
									"use_meta_index":    "false",
									"compression_codec": "none",
									"role_arn":          "acs:log:asdasd:asdas:dasdasdasd",
								},
							},
						},
					},
					"oss_ingestion_name": name,
					"display_name":       "tf-testacc-875",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":            CHECKSET,
						"oss_ingestion_name": name,
						"display_name":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testacc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project":     "${alicloud_log_project.defaultProject.name}",
					"description": "tf-testacc",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultILogStore.name}",
							"source": []map[string]interface{}{
								{
									"endpoint":          "oss-cn-chengdu.aliyuncs.com",
									"bucket":            "${alicloud_oss_bucket.defaultBucket.bucket_name}",
									"encoding":          "utf-16",
									"format":            "{\\\"type\\\": \\\"text\\\"}",
									"interval":          "3m",
									"pattern":           "1day*",
									"prefix":            "1day",
									"start_time":        "1706792358",
									"end_time":          "1706792558",
									"time_field":        "__time__",
									"time_format":       "epoch",
									"time_pattern":      "\\\\d+:\\\\d+:\\\\d+",
									"time_zone":         "GMT-09:00",
									"use_meta_index":    "false",
									"compression_codec": "none",
									"role_arn":          "acs:log:asdasd:asdas:dasdasdasd",
								},
							},
						},
					},
					"oss_ingestion_name": name + "_update",
					"schedule": []map[string]interface{}{
						{
							"type":            "Resident",
							"interval":        "5m",
							"cron_expression": "*/3 * * * * *",
							"run_immediately": "false",
							"time_zone":       "GMT+8:00",
							"delay":           "5m",
						},
					},
					"display_name": "tf-testacc-12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":            CHECKSET,
						"description":        "tf-testacc",
						"oss_ingestion_name": name + "_update",
						"display_name":       CHECKSET,
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

var AlicloudSlsOssIngestionMap5983 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudSlsOssIngestionBasicDependence5983(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultBucket" {
  bucket = var.name

  storage_class = "Standard"
}

resource "alicloud_log_project" "defaultProject" {
  name = var.name

}

resource "alicloud_log_store" "defaultILogStore" {
  retention_period = "30"
  project          = alicloud_log_project.defaultProject.name
  name             = var.name

}


`, name)
}

// Case 标准测试用例 5983  twin
func TestAccAliCloudSlsOssIngestion_basic5983_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_oss_ingestion.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsOssIngestionMap5983)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsOssIngestion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsossingestion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsOssIngestionBasicDependence5983)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":     "${alicloud_log_project.defaultProject.name}",
					"description": "tf-testacc",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultILogStore.name}",
							"source": []map[string]interface{}{
								{
									"endpoint": "oss-cn-hangzhou.aliyuncs.com",
									"bucket":   "${alicloud_oss_bucket.defaultBucket.bucket}",
									"encoding": "utf-16",
									"format": map[string]interface{}{
										"type": "line",
									},
									"interval":          "3m",
									"pattern":           "1day*",
									"prefix":            "1day",
									"start_time":        "1706792358",
									"end_time":          "1706792558",
									"time_field":        "__time__",
									"time_format":       "epoch",
									"time_pattern":      "\\\\d+:\\\\d+:\\\\d+",
									"time_zone":         "GMT-09:00",
									"use_meta_index":    "false",
									"compression_codec": "none",
									"role_arn":          "acs:log:asdasd:asdas:dasdasdasd",
								},
							},
						},
					},
					"oss_ingestion_name": name,
					"schedule": []map[string]interface{}{
						{
							"type":            "Resident",
							"interval":        "5m",
							"cron_expression": "*/3 * * * * *",
							"run_immediately": "false",
							"time_zone":       "GMT+8:00",
							"delay":           "5",
						},
					},
					"display_name": "tf-testacc-854",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":            CHECKSET,
						"description":        "tf-testacc",
						"oss_ingestion_name": name,
						"display_name":       CHECKSET,
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

// Test Sls OssIngestion. <<< Resource test cases, automatically generated.
