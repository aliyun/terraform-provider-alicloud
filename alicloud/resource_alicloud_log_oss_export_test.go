package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogOssExport_basic(t *testing.T) {
	var v *sls.Shipper
	resourceId := "alicloud_log_oss_export.default"
	ra := resourceAttrInit(resourceId, logOssExportMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogOssExportConfigDependence)

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
					"project_name":    name,
					"logstore_name":   name,
					"export_name":     "test_oss_export",
					"display_name":    "test display",
					"bucket":          "test_bucket",
					"prefix":          "",
					"suffix":          "",
					"buffer_interval": "300",
					"buffer_size":     "250",
					"compress_type":   "none",
					"path_format":     "%Y/%m/%d/%H/%M",
					"content_type":    "json",
					"json_enable_tag": "true",
					"time_zone":       "+0800",
					"role_arn":        "test-role-arn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name,
						"logstore_name":   name,
						"export_name":     "test_oss_export",
						"display_name":    "test display",
						"bucket":          "test_bucket",
						"prefix":          "",
						"suffix":          "",
						"buffer_interval": "300",
						"buffer_size":     "250",
						"compress_type":   "none",
						"path_format":     "%Y/%m/%d/%H/%M",
						"content_type":    "json",
						"json_enable_tag": "true",
						"time_zone":       "+0800",
						"role_arn":        "test-role-arn",
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
					"bucket":          "test_bucket_1",
					"buffer_interval": "350",
					"buffer_size":     "128",
					"path_format":     "%Y/%m/%d/%H",
					"compress_type":   "snappy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":          "test_bucket_1",
						"buffer_interval": "350",
						"buffer_size":     "128",
						"path_format":     "%Y/%m/%d/%H",
						"compress_type":   "snappy",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogOssExport_parquest(t *testing.T) {
	var v *sls.Shipper
	resourceId := "alicloud_log_oss_export.default"
	ra := resourceAttrInit(resourceId, logOssExportMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogOssExportConfigDependence)

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
					"project_name":    name,
					"logstore_name":   name,
					"export_name":     "test_oss_export",
					"display_name":    "test display",
					"bucket":          "test_bucket",
					"prefix":          "",
					"suffix":          "",
					"buffer_interval": "300",
					"buffer_size":     "250",
					"compress_type":   "none",
					"path_format":     "%Y/%m/%d/%H/%M",
					"content_type":    "json",
					"json_enable_tag": "true",
					"time_zone":       "+0800",
					"role_arn":        "test-role-arn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name,
						"logstore_name":   name,
						"export_name":     "test_oss_export",
						"display_name":    "test display",
						"bucket":          "test_bucket",
						"prefix":          "",
						"suffix":          "",
						"buffer_interval": "300",
						"buffer_size":     "250",
						"compress_type":   "none",
						"path_format":     "%Y/%m/%d/%H/%M",
						"content_type":    "json",
						"json_enable_tag": "true",
						"time_zone":       "+0800",
						"role_arn":        "test-role-arn",
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
					"content_type": "parquet",
					"config_columns": []map[string]interface{}{
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
						"content_type":     "parquet",
						"config_columns.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogOssExport_csv(t *testing.T) {
	var v *sls.Shipper
	resourceId := "alicloud_log_oss_export.default"
	ra := resourceAttrInit(resourceId, logOssExportMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogOssExportConfigDependence)

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
					"project_name":    name,
					"logstore_name":   name,
					"export_name":     "test_oss_export",
					"display_name":    "test display",
					"bucket":          "test_bucket",
					"prefix":          "",
					"suffix":          "",
					"buffer_interval": "300",
					"buffer_size":     "250",
					"compress_type":   "none",
					"path_format":     "%Y/%m/%d/%H/%M",
					"content_type":    "json",
					"json_enable_tag": "true",
					"time_zone":       "+0800",
					"role_arn":        "test-role-arn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name,
						"logstore_name":   name,
						"export_name":     "test_oss_export",
						"display_name":    "test display",
						"bucket":          "test_bucket",
						"prefix":          "",
						"suffix":          "",
						"buffer_interval": "300",
						"buffer_size":     "250",
						"compress_type":   "none",
						"path_format":     "%Y/%m/%d/%H/%M",
						"content_type":    "json",
						"json_enable_tag": "true",
						"time_zone":       "+0800",
						"role_arn":        "test-role-arn",
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
					"content_type":         "csv",
					"csv_config_delimiter": ",",
					"csv_config_header":    "false",
					"csv_config_linefeed":  "",
					"csv_config_quote":     ",",
					"csv_config_columns":   []string{"aini", "aliyun"},
					"csv_config_null":      "",
					"csv_config_escape":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_type":         "csv",
						"csv_config_delimiter": ",",
						"csv_config_header":    "false",
						"csv_config_linefeed":  "",
						"csv_config_quote":     ",",
						"csv_config_columns.#": "2",
						"csv_config_null":      "",
						"csv_config_escape":    "",
					}),
				),
			},
		},
	})
}

func resourceLogOssExportConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	resource "alicloud_log_store" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    shard_count = 1
	}
	`, name)
}

var logOssExportMap = map[string]string{
	"project_name":  CHECKSET,
	"logstore_name": CHECKSET,
	"export_name":   CHECKSET,
}
