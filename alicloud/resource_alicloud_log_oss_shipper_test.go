package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAliCloudLogOssShipper_basic(t *testing.T) {
	var v *sls.Shipper
	resourceId := "alicloud_log_oss_shipper.default"
	ra := resourceAttrInit(resourceId, logOssShipperMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("test-log-oss-shipper-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogOssShipperConfigDependence)

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
					"project_name":    "${alicloud_log_store.default.project}",
					"logstore_name":   "${alicloud_log_store.default.name}",
					"shipper_name":    "test_shipper",
					"oss_bucket":      "test_bucket",
					"oss_prefix":      "",
					"buffer_interval": "300",
					"buffer_size":     "250",
					"compress_type":   "none",
					"path_format":     "%Y/%m/%d/%H/%M",
					"format":          "json",
					"json_enable_tag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name,
						"logstore_name":   name,
						"shipper_name":    "test_shipper",
						"oss_bucket":      "test_bucket",
						"oss_prefix":      "",
						"buffer_interval": "300",
						"buffer_size":     "250",
						"compress_type":   "none",
						"path_format":     "%Y/%m/%d/%H/%M",
						"format":          "json",
						"json_enable_tag": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket": "test_bucket_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket": "test_bucket_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"buffer_interval": "350",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buffer_interval": "350",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"buffer_size": "128",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buffer_size": "128",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path_format": "%Y/%m/%d/%H",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path_format": "%Y/%m/%d/%H",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "snappy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "snappy",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"format": "parquet",
					"parquet_config": []map[string]interface{}{
						{
							"name": "name1",
							"type": "string",
						},
						{
							"name": "name2",
							"type": "int64",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"format":           "parquet",
						"parquet_config.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"format":                    "csv",
					"csv_config_delimiter":      ",",
					"csv_config_header":         "false",
					"csv_config_linefeed":       "",
					"csv_config_quote":          ",",
					"csv_config_columns":        []string{"aini", "aliyun"},
					"csv_config_nullidentifier": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"format":                    "csv",
						"csv_config_delimiter":      ",",
						"csv_config_header":         "false",
						"csv_config_linefeed":       "",
						"csv_config_quote":          ",",
						"csv_config_columns.#":      "2",
						"csv_config_nullidentifier": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket":      "test_bucket",
					"oss_prefix":      "root",
					"buffer_interval": "300",
					"buffer_size":     "250",
					"compress_type":   "none",
					"path_format":     "%Y/%m/%d/%H/%M",
					"format":          "json",
					"json_enable_tag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket":      "test_bucket",
						"oss_prefix":      "root",
						"buffer_interval": "300",
						"buffer_size":     "250",
						"compress_type":   "none",
						"path_format":     "%Y/%m/%d/%H/%M",
						"format":          "json",
						"json_enable_tag": "true",
					}),
				),
			},
		},
	})
}

func resourceLogOssShipperConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_log_project" "default" {
  		name        = var.name
  		description = "tf unit test"
	}

	resource "alicloud_log_store" "default" {
  		project          = alicloud_log_project.default.name
  		name             = var.name
  		retention_period = "3000"
  		shard_count      = 1
	}
`, name)
}

var logOssShipperMap = map[string]string{
	"project_name":  CHECKSET,
	"logstore_name": CHECKSET,
	"shipper_name":  CHECKSET,
}
