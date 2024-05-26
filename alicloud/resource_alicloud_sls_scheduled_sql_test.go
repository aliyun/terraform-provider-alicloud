package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls ScheduledSQL. >>> Resource test cases, automatically generated.
// Case test-_schedule 6614
func TestAccAliCloudSlsScheduledSQL_basic6614(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6614)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6614)
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
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-tf-scheduled-sql-0006",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":       "test-tf-scheduled-sql-0006",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-tf-scheduled-sql-0006",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-tf-scheduled-sql-0006",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-tf-scheduled-sql-0006",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-tf-scheduled-sql-0006",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name + "_update",
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-tf-scheduled-sql-0006",
						"display_name":       "test-tf-scheduled-sql-0006",
						"scheduled_sql_name": name + "_update",
						"project":            CHECKSET,
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

var AlicloudSlsScheduledSQLMap6614 = map[string]string{}

func AlicloudSlsScheduledSQLBasicDependence6614(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultKIe4KV" {
  description = "terraform-scheduledsql-test-941"
  name        = var.name
}

resource "alicloud_log_store" "default1LI9we" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultKIe4KV.name
  name             = var.name
}


`, name)
}

// Case scheduled-common 6613
func TestAccAliCloudSlsScheduledSQL_basic6613(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6613)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6613)
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
					"schedule": []map[string]interface{}{
						{
							"type":            "FixedRate",
							"interval":        "5m",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
						},
					},
					"display_name": "test-scheduled-sql-00001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":       "test-scheduled-sql-00001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-00001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-scheduled-sql-00001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-00001-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-scheduled-sql-00001-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "test-scheduled-sql-00001-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "test-scheduled-sql-00001-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-00001",
					"schedule": []map[string]interface{}{
						{
							"type":            "FixedRate",
							"interval":        "5m",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
						},
					},
					"display_name": "test-scheduled-sql-00001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name + "_update",
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-00001",
						"display_name":       "test-scheduled-sql-00001",
						"scheduled_sql_name": name + "_update",
						"project":            CHECKSET,
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

var AlicloudSlsScheduledSQLMap6613 = map[string]string{}

func AlicloudSlsScheduledSQLBasicDependence6613(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultKIe4KV" {
  description = "terraform-scheduledsql-test-690"
  name        = var.name
}

resource "alicloud_log_store" "default1LI9we" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultKIe4KV.name
  name             = var.name
}


`, name)
}

// Case scheduledSQL-cron 6609
func TestAccAliCloudSlsScheduledSQL_basic6609(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6609)
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
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-scheduled-sql-log2metric-000001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-hcl",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-hcl",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/test-hcl",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713192800",
							"to_time":                 "1713197800",
							"data_format":             "log2metric",
							"parameters": map[string]interface{}{
								"\"timeKey\"":    "timeKey",
								"\"metricKeys\"": "[\\\"df\\\",\\\"ef\\\",\\\"ff\\\"]",
								"\"labelKeys\"":  "[\\\"gf\\\",\\\"hf\\\",\\\"if\\\"]",
								"\"hashLabels\"": "[\\\"kf\\\",\\\"lf\\\",\\\"mf\\\"]",
								"\"addLabels\"":  "{\\\"sa\\\": \\\"asdfdfasdfasdfasdfasdf\\\",\\\"ef\\\": \\\"f\\\"}",
							},
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":       "test-scheduled-sql-log2metric-000001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-log2metric-000001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-scheduled-sql-log2metric-000001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-log2metric-000001",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-scheduled-sql-log2metric-000001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-hcl",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-hcl",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/test-hcl",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713192800",
							"to_time":                 "1713197800",
							"data_format":             "log2metric",
							"parameters": map[string]interface{}{
								"\"timeKey\"":    "timeKey",
								"\"metricKeys\"": "[\\\"df\\\",\\\"ef\\\",\\\"ff\\\"]",
								"\"labelKeys\"":  "[\\\"gf\\\",\\\"hf\\\",\\\"if\\\"]",
								"\"hashLabels\"": "[\\\"kf\\\",\\\"lf\\\",\\\"mf\\\"]",
								"\"addLabels\"":  "{\\\"sa\\\": \\\"asdfdfasdfasdfasdfasdf\\\",\\\"ef\\\": \\\"f\\\"}",
							},
						},
					},
					"scheduled_sql_name": name + "_update",
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-log2metric-000001",
						"display_name":       "test-scheduled-sql-log2metric-000001",
						"scheduled_sql_name": name + "_update",
						"project":            CHECKSET,
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

var AlicloudSlsScheduledSQLMap6609 = map[string]string{}

func AlicloudSlsScheduledSQLBasicDependence6609(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultKIe4KV" {
  description = "terraform-scheduledsql-test-245"
  name        = var.name
}

resource "alicloud_log_store" "default1LI9we" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultKIe4KV.name
  name             = var.name
}


`, name)
}

// Case test-_schedule 6614  twin
func TestAccAliCloudSlsScheduledSQL_basic6614_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6614)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6614)
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
					"description": "test-tf-scheduled-sql-0006",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-tf-scheduled-sql-0006",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-tf-scheduled-sql-0006",
						"display_name":       "test-tf-scheduled-sql-0006",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
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

// Case scheduled-common 6613  twin
func TestAccAliCloudSlsScheduledSQL_basic6613_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6613)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6613)
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
					"description": "test-scheduled-sql-00001",
					"schedule": []map[string]interface{}{
						{
							"type":            "FixedRate",
							"interval":        "5m",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
						},
					},
					"display_name": "test-scheduled-sql-00001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-00001",
						"display_name":       "test-scheduled-sql-00001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
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

// Case scheduledSQL-cron 6609  twin
func TestAccAliCloudSlsScheduledSQL_basic6609_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6609)
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
					"description": "test-scheduled-sql-log2metric-000001",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-scheduled-sql-log2metric-000001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-hcl",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-hcl",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/test-hcl",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713192800",
							"to_time":                 "1713197800",
							"data_format":             "log2metric",
							"parameters": map[string]interface{}{
								"\"timeKey\"":    "timeKey",
								"\"metricKeys\"": "[\\\"df\\\",\\\"ef\\\",\\\"ff\\\"]",
								"\"labelKeys\"":  "[\\\"gf\\\",\\\"hf\\\",\\\"if\\\"]",
								"\"hashLabels\"": "[\\\"kf\\\",\\\"lf\\\",\\\"mf\\\"]",
								"\"addLabels\"":  "{\\\"sa\\\": \\\"asdfdfasdfasdfasdfasdf\\\",\\\"ef\\\": \\\"f\\\"}",
							},
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-log2metric-000001",
						"display_name":       "test-scheduled-sql-log2metric-000001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
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

// Case test-_schedule 6614  raw
func TestAccAliCloudSlsScheduledSQL_basic6614_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6614)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6614)
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
					"description": "test-tf-scheduled-sql-0006",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-tf-scheduled-sql-0006",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-tf-scheduled-sql-0006",
						"display_name":       "test-tf-scheduled-sql-0006",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immediately": "true",
							"time_zone":       "+0800",
							"delay":           "20",
							"cron_expression": "0 0/3 * * *",
						},
					},
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "standard",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogdefaultrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_retries":             "10",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "default",
						},
					},
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

// Case scheduled-common 6613  raw
func TestAccAliCloudSlsScheduledSQL_basic6613_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6613)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6613)
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
					"description": "test-scheduled-sql-00001",
					"schedule": []map[string]interface{}{
						{
							"type":            "FixedRate",
							"interval":        "5m",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
						},
					},
					"display_name": "test-scheduled-sql-00001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api02",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-00001",
						"display_name":       "test-scheduled-sql-00001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-scheduled-sql-00001-update",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immediately": "true",
							"time_zone":       "+0800",
							"delay":           "20",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-scheduled-sql-00001-update",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-open-api",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_retries":             "10",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test-scheduled-sql-00001-update",
						"display_name": "test-scheduled-sql-00001-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":            "FixedRate",
							"interval":        "10m",
							"run_immediately": "false",
							"time_zone":       "+0800",
							"delay":           "20",
						},
					},
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "cn-hangzhou.log.aliyuncs.com",
							"dest_project":            "gy-hangzhou-huolang-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-3",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"from_time_expr":          "@m-5m",
							"to_time_expr":            "@m+1m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "10",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
							"parameters": map[string]interface{}{
								"\"a\"": "c",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":      "FixedRate",
							"interval":  "10m",
							"time_zone": "+0800",
							"delay":     "22",
						},
					},
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select bucket,end_time from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "cn-hangzhou.log.aliyuncs.com",
							"dest_project":            "gy-hangzhou-huolang-1",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-3",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/test-hcl",
							"from_time_expr":          "@m-5m",
							"to_time_expr":            "@m+1m",
							"max_run_time_in_seconds": "1600",
							"resource_pool":           "enhanced",
							"max_retries":             "25",
							"from_time":               "1713196800",
							"to_time":                 "0",
							"data_format":             "log2log",
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

// Case scheduledSQL-cron 6609  raw
func TestAccAliCloudSlsScheduledSQL_basic6609_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_scheduled_sql.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsScheduledSQLMap6609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsScheduledSQL")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsscheduledsql%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsScheduledSQLBasicDependence6609)
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
					"description": "test-scheduled-sql-log2metric-000001",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"time_zone":       "+0700",
							"delay":           "20",
							"run_immediately": "false",
							"cron_expression": "0 0/1 * * *",
						},
					},
					"display_name": "test-scheduled-sql-log2metric-000001",
					"scheduled_sql_configuration": []map[string]interface{}{
						{
							"script":                  "* | select * from log",
							"sql_type":                "searchQuery",
							"dest_endpoint":           "ap-northeast-1.log.aliyuncs.com",
							"dest_project":            "job-e2e-project-jj78kur-ap-southeast-hcl",
							"source_logstore":         "${alicloud_log_store.default1LI9we.name}",
							"dest_logstore":           "test-hcl",
							"role_arn":                "acs:ram::1395894005868720:role/aliyunlogetlrole",
							"dest_role_arn":           "acs:ram::1395894005868720:role/test-hcl",
							"from_time_expr":          "@m-1m",
							"to_time_expr":            "@m",
							"max_run_time_in_seconds": "1800",
							"resource_pool":           "enhanced",
							"max_retries":             "5",
							"from_time":               "1713192800",
							"to_time":                 "1713197800",
							"data_format":             "log2metric",
							"parameters": map[string]interface{}{
								"\"timeKey\"":    "timeKey",
								"\"metricKeys\"": "[\\\"df\\\",\\\"ef\\\",\\\"ff\\\"]",
								"\"labelKeys\"":  "[\\\"gf\\\",\\\"hf\\\",\\\"if\\\"]",
								"\"hashLabels\"": "[\\\"kf\\\",\\\"lf\\\",\\\"mf\\\"]",
								"\"addLabels\"":  "{\\\"sa\\\": \\\"asdfdfasdfasdfasdfasdf\\\",\\\"ef\\\": \\\"f\\\"}",
							},
						},
					},
					"scheduled_sql_name": name,
					"project":            "${alicloud_log_project.defaultKIe4KV.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "test-scheduled-sql-log2metric-000001",
						"display_name":       "test-scheduled-sql-log2metric-000001",
						"scheduled_sql_name": name,
						"project":            CHECKSET,
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

// Test Sls ScheduledSQL. <<< Resource test cases, automatically generated.
