package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls OssExportSink. >>> Resource test cases, automatically generated.
// Case test-oss-example2 9137
func TestAccAliCloudSlsOssExportSink_basic9137(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_oss_export_sink.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsOssExportSinkMap9137)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsOssExportSink")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsossexportsink%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsOssExportSinkBasicDependence9137)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.defaulteyHJsO.name}",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"bucket":           "${alicloud_oss_bucket.defaultiwj0xO.bucket}",
									"role_arn":         "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
									"time_zone":        "+0700",
									"content_type":     "json",
									"compression_type": "none",
									"content_detail":   "{\\\"enableTag\\\": false} ",
									"buffer_interval":  "300",
									"buffer_size":      "256",
									"endpoint":         "https://oss-cn-shanghai-internal.aliyuncs.com",
								},
							},
							"from_time": "1732165733",
							"to_time":   "1732166733",
						},
					},
					"job_name":     "export-oss-1731404933-00001",
					"display_name": "testterraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":      CHECKSET,
						"job_name":     "export-oss-1731404933-00001",
						"display_name": "testterraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::1395891111111111111:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"bucket":           "gyhuolang",
									"role_arn":         "acs:ram::1395891111111111111:role/aliyunlogdefaultrole",
									"time_zone":        "+0800",
									"content_type":     "csv",
									"compression_type": "snappy",
									"content_detail":   "{\\\"null\\\": \\\"-\\\", \\\"header\\\": false, \\\"lineFeed\\\": \\\"\\\\n\\\", \\\"quote\\\": \\\"\\\", \\\"delimiter\\\": \\\",\\\", \\\"columns\\\": [\\\"a\\\", \\\"b\\\", \\\"c\\\", \\\"d\\\"]}",
									"buffer_interval":  "300",
									"buffer_size":      "256",
									"endpoint":         "https://oss-cn-beijing-internal.aliyuncs.com",
								},
							},
							"from_time": "1",
							"to_time":   "0",
						},
					},
					"display_name": "test-huolang-hcl3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "test-huolang-hcl3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::1311111111111111111:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"endpoint":         "https://oss-cn-hangzhou-internal.aliyuncs.com",
									"bucket":           "gyhuolang",
									"prefix":           "terraform22",
									"suffix":           ".csv",
									"role_arn":         "acs:ram::1311111111111111111:role/aliyunlogdefaultrole",
									"path_format":      "%Y/%m/%d/%H",
									"time_zone":        "+0800",
									"content_type":     "parquet",
									"compression_type": "gzip",
									"buffer_interval":  "300",
									"buffer_size":      "256",
									"delay_seconds":    "909",
									"content_detail":   "{\\\"columns\\\": [{\\\"name\\\": \\\"a\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"b\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"c\\\", \\\"type\\\": \\\"string\\\"}]}",
								},
							},
							"from_time": "1",
							"to_time":   "0",
						},
					},
					"display_name": "test-huolang-hcl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "test-huolang-hcl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"endpoint":         "https://oss-cn-beijing-internal.aliyuncs.com",
									"bucket":           "${alicloud_oss_bucket.defaultiwj0xO.bucket}",
									"prefix":           "terraform",
									"suffix":           ".gzip",
									"role_arn":         "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
									"path_format":      "%Y/%m/%d/%H",
									"path_format_type": "time",
									"time_zone":        "+0800",
									"content_type":     "orc",
									"compression_type": "none",
									"buffer_interval":  "301",
									"buffer_size":      "255",
									"delay_seconds":    "909",
									"content_detail":   "{\\\"columns\\\": [{\\\"name\\\": \\\"a\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"b\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"c\\\", \\\"type\\\": \\\"string\\\"}]}",
								},
							},
							"from_time": "1732165734",
							"to_time":   "1732186733",
						},
					},
					"display_name": "testterraform",
					"description":  "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "testterraform",
						"description":  "test",
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

var AlicloudSlsOssExportSinkMap9137 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudSlsOssExportSinkBasicDependence9137(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaulteyHJsO" {
  description = "terraform-oss-test-788"
  name        = var.name
}

resource "alicloud_log_store" "defaultxeHfXC" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaulteyHJsO.name
  name             = format("%%s1", var.name)
}

resource "alicloud_oss_bucket" "defaultiwj0xO" {
  storage_class = "Standard"
}


`, name)
}

// Case test-oss 9051
func TestAccAliCloudSlsOssExportSink_basic9051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_oss_export_sink.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsOssExportSinkMap9051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsOssExportSink")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsossexportsink%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsOssExportSinkBasicDependence9051)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":     "${alicloud_log_project.defaulteyHJsO.name}",
					"description": "test-huolang-hcl",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"endpoint":         "https://oss-cn-beijing-internal.aliyuncs.com",
									"bucket":           "${alicloud_oss_bucket.default2aGp8n.bucket}",
									"prefix":           "terraform",
									"role_arn":         "acs:ram::12345678901234567:role/aliyunlogdefaultrole",
									"path_format":      "%Y/%m/%d/%H",
									"time_zone":        "+0700",
									"content_type":     "json",
									"compression_type": "none",
									"buffer_interval":  "301",
									"buffer_size":      "255",
									"delay_seconds":    "900",
									"suffix":           ".gzip",
									"path_format_type": "time",
									"content_detail":   "{\\\"enableTag\\\": true}",
								},
							},
							"from_time": "1732165733",
							"to_time":   "1732166733",
						},
					},
					"job_name":     "export-oss-1731404933-00001",
					"display_name": "testterraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":      CHECKSET,
						"description":  "test-huolang-hcl",
						"job_name":     "export-oss-1731404933-00001",
						"display_name": "testterraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-huolang-hcl2",
					"configuration": []map[string]interface{}{
						{
							"logstore": "${alicloud_log_store.defaultxeHfXC.name}",
							"role_arn": "acs:ram::1395891111111111111:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"endpoint":         "https://oss-cn-hangzhou-internal.aliyuncs.com",
									"bucket":           "${alicloud_oss_bucket.default2aGp8n.bucket}",
									"prefix":           "terraform2",
									"suffix":           "json",
									"role_arn":         "acs:ram::1395891111111111111:role/aliyunlogdefaultrole",
									"path_format":      "%Y/%m/%d/%H/%M",
									"time_zone":        "+0800",
									"content_type":     "csv",
									"compression_type": "snappy",
									"buffer_interval":  "300",
									"buffer_size":      "256",
									"content_detail":   "{\\\"null\\\": \\\"-\\\", \\\"header\\\": false, \\\"lineFeed\\\": \\\"\\\\n\\\", \\\"quote\\\": \\\"\\\", \\\"delimiter\\\": \\\",\\\", \\\"columns\\\": [\\\"a\\\", \\\"b\\\", \\\"c\\\", \\\"d\\\"]}",
									"delay_seconds":    "901",
									"path_format_type": "time",
								},
							},
							"from_time": "1",
							"to_time":   "0",
						},
					},
					"display_name": "test-huolang-hcl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test-huolang-hcl2",
						"display_name": "test-huolang-hcl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"role_arn": "acs:ram::1311111111111111111:role/aliyunlogdefaultrole",
							"sink": []map[string]interface{}{
								{
									"endpoint":         "https://oss-cn-hangzhou-internal.aliyuncs.com",
									"bucket":           "${alicloud_oss_bucket.default2aGp8n.bucket}",
									"prefix":           "terraform22",
									"suffix":           ".csv",
									"role_arn":         "acs:ram::1311111111111111111:role/aliyunlogdefaultrole",
									"path_format":      "%Y/%m/%d/%H",
									"time_zone":        "+0800",
									"content_type":     "orc",
									"compression_type": "zstd",
									"buffer_interval":  "300",
									"buffer_size":      "256",
									"delay_seconds":    "909",
									"content_detail":   "{\\\"columns\\\": [{\\\"name\\\": \\\"a\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"b\\\", \\\"type\\\": \\\"string\\\"}, {\\\"name\\\": \\\"c\\\", \\\"type\\\": \\\"string\\\"}]}",
								},
							},
							"from_time": "1",
							"to_time":   "0",
							"logstore":  "${alicloud_log_store.defaultxeHfXC.name}",
						},
					},
				}),
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

var AlicloudSlsOssExportSinkMap9051 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudSlsOssExportSinkBasicDependence9051(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaulteyHJsO" {
  description = "terraform-oss-test-500"
  name        = var.name
}

resource "alicloud_log_store" "defaultxeHfXC" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaulteyHJsO.name
  name             = format("%%s1", var.name)
}

resource "alicloud_log_store" "defaultoVvRkC" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaulteyHJsO.name
  name             = format("%%s2", var.name)
}

resource "alicloud_oss_bucket" "default2aGp8n" {
  storage_class = "Standard"
}


`, name)
}

// Test Sls OssExportSink. <<< Resource test cases, automatically generated.
