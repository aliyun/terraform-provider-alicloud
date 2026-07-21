package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb StreamingJob. >>> Resource test cases, automatically generated.
// Case StreamingJob_资源依赖_case_002 7220
func TestAccAliCloudGpdbStreamingJob_basic7220(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_job.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingJobMap7220)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingJobBasicDependence7220)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account":         "test_001",
					"dest_schema":     "public",
					"mode":            "professional",
					"job_name":        "test-kafka",
					"job_description": "test-kafka",
					"dest_database":   "adb_sampledata_tpch",
					"db_instance_id":  "${alicloud_gpdb_instance.defaulth2ghc1.id}",
					"dest_table":      "customer",
					"data_source_id":  "${alicloud_gpdb_streaming_data_source.defaultcDQItu.data_source_id}",
					"password":        "test_001",
					"try_run":         "false",
					"job_config":      "DATABASE: adb_sampledata_tpch\\nUSER: test_001\\nPASSWORD: test_001\\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\\nPORT: 5432\\nKAFKA:\\n  INPUT:\\n    SOURCE:\\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\\n      TOPIC: ziyuan_test\\n      FALLBACK_OFFSET: LATEST\\n    KEY:\\n      COLUMNS:\\n      - NAME: c_custkey\\n        TYPE: int\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    VALUE:\\n      COLUMNS:\\n      - NAME: c_comment\\n        TYPE: varchar\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    ERROR_LIMIT: 10\\n  OUTPUT:\\n    SCHEMA: public\\n    TABLE: customer\\n    MODE: MERGE\\n    MATCH_COLUMNS:\\n    - c_custkey\\n    ORDER_COLUMNS:\\n    - c_custkey\\n    UPDATE_COLUMNS:\\n    - c_custkey\\n    MAPPING:\\n    - NAME: c_custkey\\n      EXPRESSION: c_custkey\\n  COMMIT:\\n    MAX_ROW: 1000\\n    MINIMAL_INTERVAL: 1000\\n    CONSISTENCY: ATLEAST\\n  POLL:\\n    BATCHSIZE: 1000\\n    TIMEOUT: 1000\\n  PROPERTIES:\\n    group.id: ziyuan_test_01\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account":         "test_001",
						"dest_schema":     "public",
						"mode":            "professional",
						"job_name":        "test-kafka",
						"job_description": "test-kafka",
						"dest_database":   "adb_sampledata_tpch",
						"db_instance_id":  CHECKSET,
						"dest_table":      "customer",
						"data_source_id":  CHECKSET,
						"password":        "test_001",
						"try_run":         "false",
						"job_config":      "DATABASE: adb_sampledata_tpch\nUSER: test_001\nPASSWORD: test_001\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\nPORT: 5432\nKAFKA:\n  INPUT:\n    SOURCE:\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\n      TOPIC: ziyuan_test\n      FALLBACK_OFFSET: LATEST\n    KEY:\n      COLUMNS:\n      - NAME: c_custkey\n        TYPE: int\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    VALUE:\n      COLUMNS:\n      - NAME: c_comment\n        TYPE: varchar\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    ERROR_LIMIT: 10\n  OUTPUT:\n    SCHEMA: public\n    TABLE: customer\n    MODE: MERGE\n    MATCH_COLUMNS:\n    - c_custkey\n    ORDER_COLUMNS:\n    - c_custkey\n    UPDATE_COLUMNS:\n    - c_custkey\n    MAPPING:\n    - NAME: c_custkey\n      EXPRESSION: c_custkey\n  COMMIT:\n    MAX_ROW: 1000\n    MINIMAL_INTERVAL: 1000\n    CONSISTENCY: ATLEAST\n  POLL:\n    BATCHSIZE: 1000\n    TIMEOUT: 1000\n  PROPERTIES:\n    group.id: ziyuan_test_01\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_description": "test-kafka2",
					"job_config":      "DATABASE: adb_sampledata_tpch\\nUSER: test_001\\nPASSWORD: test_001\\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\\nPORT: 5432\\nKAFKA:\\n  INPUT:\\n    SOURCE:\\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\\n      TOPIC: ziyuan_test\\n      FALLBACK_OFFSET: LATEST\\n    KEY:\\n      COLUMNS:\\n      - NAME: c_custkey\\n        TYPE: int\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    VALUE:\\n      COLUMNS:\\n      - NAME: c_comment\\n        TYPE: varchar\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    ERROR_LIMIT: 10\\n  OUTPUT:\\n    SCHEMA: public\\n    TABLE: customer\\n    MODE: MERGE\\n    MATCH_COLUMNS:\\n    - c_custkey\\n    ORDER_COLUMNS:\\n    - c_custkey\\n    UPDATE_COLUMNS:\\n    - c_custkey\\n    MAPPING:\\n    - NAME: c_custkey\\n      EXPRESSION: c_custkey\\n  COMMIT:\\n    MAX_ROW: 1000\\n    MINIMAL_INTERVAL: 1000\\n    CONSISTENCY: ATLEAST\\n  POLL:\\n    BATCHSIZE: 1000\\n    TIMEOUT: 1000\\n  PROPERTIES:\\n    group.id: ziyuan_test\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_description":  "test-kafka2",
						"job_config":       "DATABASE: adb_sampledata_tpch\nUSER: test_001\nPASSWORD: test_001\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\nPORT: 5432\nKAFKA:\n  INPUT:\n    SOURCE:\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\n      TOPIC: ziyuan_test\n      FALLBACK_OFFSET: LATEST\n    KEY:\n      COLUMNS:\n      - NAME: c_custkey\n        TYPE: int\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    VALUE:\n      COLUMNS:\n      - NAME: c_comment\n        TYPE: varchar\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    ERROR_LIMIT: 10\n  OUTPUT:\n    SCHEMA: public\n    TABLE: customer\n    MODE: MERGE\n    MATCH_COLUMNS:\n    - c_custkey\n    ORDER_COLUMNS:\n    - c_custkey\n    UPDATE_COLUMNS:\n    - c_custkey\n    MAPPING:\n    - NAME: c_custkey\n      EXPRESSION: c_custkey\n  COMMIT:\n    MAX_ROW: 1000\n    MINIMAL_INTERVAL: 1000\n    CONSISTENCY: ATLEAST\n  POLL:\n    BATCHSIZE: 1000\n    TIMEOUT: 1000\n  PROPERTIES:\n    group.id: ziyuan_test\n",
						"src_columns.#":    "0",
						"dest_columns.#":   "0",
						"update_columns.#": "0",
						"match_columns.#":  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_config": "DATABASE: adb_sampledata_tpch\\nUSER: test_001\\nPASSWORD: test_001\\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\\nPORT: 5432\\nKAFKA:\\n  INPUT:\\n    SOURCE:\\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\\n      TOPIC: ziyuan_test\\n      FALLBACK_OFFSET: LATEST\\n    KEY:\\n      COLUMNS:\\n      - NAME: c_custkey\\n        TYPE: int\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    VALUE:\\n      COLUMNS:\\n      - NAME: c_comment\\n        TYPE: varchar\\n      FORMAT: delimited\\n      DELIMITED_OPTION:\\n        DELIMITER: '|'\\n    ERROR_LIMIT: 10\\n  OUTPUT:\\n    SCHEMA: public\\n    TABLE: customer\\n    MODE: MERGE\\n    MATCH_COLUMNS:\\n    - c_custkey\\n    ORDER_COLUMNS:\\n    - c_custkey\\n    UPDATE_COLUMNS:\\n    - c_custkey\\n    MAPPING:\\n    - NAME: c_custkey\\n      EXPRESSION: c_custkey\\n  COMMIT:\\n    MAX_ROW: 1000\\n    MINIMAL_INTERVAL: 1000\\n    CONSISTENCY: ATLEAST\\n  POLL:\\n    BATCHSIZE: 1000\\n    TIMEOUT: 1000\\n  PROPERTIES:\\n    group.id: ziyuan_test_01\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_config": "DATABASE: adb_sampledata_tpch\nUSER: test_001\nPASSWORD: test_001\nHOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com\nPORT: 5432\nKAFKA:\n  INPUT:\n    SOURCE:\n      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092\n      TOPIC: ziyuan_test\n      FALLBACK_OFFSET: LATEST\n    KEY:\n      COLUMNS:\n      - NAME: c_custkey\n        TYPE: int\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    VALUE:\n      COLUMNS:\n      - NAME: c_comment\n        TYPE: varchar\n      FORMAT: delimited\n      DELIMITED_OPTION:\n        DELIMITER: '|'\n    ERROR_LIMIT: 10\n  OUTPUT:\n    SCHEMA: public\n    TABLE: customer\n    MODE: MERGE\n    MATCH_COLUMNS:\n    - c_custkey\n    ORDER_COLUMNS:\n    - c_custkey\n    UPDATE_COLUMNS:\n    - c_custkey\n    MAPPING:\n    - NAME: c_custkey\n      EXPRESSION: c_custkey\n  COMMIT:\n    MAX_ROW: 1000\n    MINIMAL_INTERVAL: 1000\n    CONSISTENCY: ATLEAST\n  POLL:\n    BATCHSIZE: 1000\n    TIMEOUT: 1000\n  PROPERTIES:\n    group.id: ziyuan_test_01\n",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"try_run"},
			},
		},
	})
}

var AlicloudGpdbStreamingJobMap7220 = map[string]string{
	"job_id": CHECKSET,
	"status": CHECKSET,
}

func AlicloudGpdbStreamingJobBasicDependence7220(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultTXqb15" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultaSWhbT" {
  vpc_id     = alicloud_vpc.defaultTXqb15.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaulth2ghc1" {
  instance_spec              = "2C8G"
  description                = "创建流数据源需要的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = alicloud_vswitch.defaultaSWhbT.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultTXqb15.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "default2dUszY" {
  service_name        = "test"
  db_instance_id      = alicloud_gpdb_instance.defaulth2ghc1.id
  service_description = "test"
  service_spec        = "8"
}

resource "alicloud_gpdb_streaming_data_source" "defaultcDQItu" {
  db_instance_id          = alicloud_gpdb_instance.defaulth2ghc1.id
  data_source_name        = "test"
  data_source_config      = "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}"
  data_source_type        = "kafka"
  data_source_description = "test"
  service_id              = alicloud_gpdb_streaming_data_service.default2dUszY.service_id
}


`, name)
}

// Case StreamingJob_资源依赖_case_001 7063
func TestAccAliCloudGpdbStreamingJob_basic7063(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_job.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingJobMap7063)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingJobBasicDependence7063)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":  "ziyuan_test",
					"account":     "test_001",
					"dest_schema": "public",
					"mode":        "basic",
					"job_name":    "test-kafka",
					"src_columns": []string{
						"l_orderkey", "l_partkey", "l_suppkey"},
					"dest_columns": []string{
						"l_orderkey", "l_partkey", "l_suppkey"},
					"job_description":   "test-kafka",
					"dest_database":     "adb_sampledata_tpch",
					"db_instance_id":    "${alicloud_gpdb_instance.defaulth2ghc1.id}",
					"fallback_offset":   "earliest",
					"dest_table":        "lineitem",
					"data_source_id":    "${alicloud_gpdb_streaming_data_source.defaultcDQItu.data_source_id}",
					"error_limit_count": "100",
					"write_mode":        "update",
					"password":          "test_001",
					"consistency":       "ATLEAST",
					"update_columns": []string{
						"l_linenumber", "l_quantity", "l_extendedprice"},
					"match_columns": []string{
						"l_orderkey", "l_partkey", "l_suppkey"},
					"try_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":        "ziyuan_test",
						"account":           "test_001",
						"dest_schema":       "public",
						"mode":              "basic",
						"job_name":          "test-kafka",
						"src_columns.#":     "3",
						"dest_columns.#":    "3",
						"job_description":   "test-kafka",
						"dest_database":     "adb_sampledata_tpch",
						"db_instance_id":    CHECKSET,
						"fallback_offset":   "earliest",
						"dest_table":        "lineitem",
						"data_source_id":    CHECKSET,
						"error_limit_count": "100",
						"write_mode":        "update",
						"password":          "test_001",
						"consistency":       "ATLEAST",
						"update_columns.#":  "3",
						"match_columns.#":   "3",
						"try_run":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":  "ziyuan_test_01",
					"account":     "test_002",
					"dest_schema": "public_01",
					"src_columns": []string{
						"l_linenumber", "l_quantity"},
					"dest_columns": []string{
						"l_linenumber", "l_quantity"},
					"job_description":   "test-job2",
					"dest_database":     "adb_sampledata_tpch_01",
					"fallback_offset":   "latest",
					"dest_table":        "lineitem_01",
					"error_limit_count": "50",
					"write_mode":        "merge",
					"password":          "test_002",
					"consistency":       "EXACTLY",
					"update_columns": []string{
						"l_discount", "l_tax"},
					"match_columns": []string{
						"l_quantity", "l_extendedprice"},
					"try_run": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":        "ziyuan_test_01",
						"account":           "test_002",
						"dest_schema":       "public_01",
						"src_columns.#":     "2",
						"dest_columns.#":    "2",
						"job_description":   "test-job2",
						"dest_database":     "adb_sampledata_tpch_01",
						"fallback_offset":   "latest",
						"dest_table":        "lineitem_01",
						"error_limit_count": "50",
						"write_mode":        "merge",
						"password":          "test_002",
						"consistency":       "EXACTLY",
						"update_columns.#":  "2",
						"match_columns.#":   "2",
						"try_run":           "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"try_run"},
			},
		},
	})
}

var AlicloudGpdbStreamingJobMap7063 = map[string]string{
	"job_id": CHECKSET,
	"status": CHECKSET,
}

func AlicloudGpdbStreamingJobBasicDependence7063(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultTXqb15" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultaSWhbT" {
  vpc_id     = alicloud_vpc.defaultTXqb15.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaulth2ghc1" {
  instance_spec              = "2C8G"
  description                = "创建流数据源需要的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = alicloud_vswitch.defaultaSWhbT.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultTXqb15.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "default2dUszY" {
  service_name        = "test"
  db_instance_id      = alicloud_gpdb_instance.defaulth2ghc1.id
  service_description = "test"
  service_spec        = "8"
}

resource "alicloud_gpdb_streaming_data_source" "defaultcDQItu" {
  db_instance_id          = alicloud_gpdb_instance.defaulth2ghc1.id
  data_source_name        = "test"
  data_source_config      = "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}"
  data_source_type        = "kafka"
  data_source_description = "test"
  service_id              = alicloud_gpdb_streaming_data_service.default2dUszY.service_id
}


`, name)
}

// Test Gpdb StreamingJob. <<< Resource test cases, automatically generated.
