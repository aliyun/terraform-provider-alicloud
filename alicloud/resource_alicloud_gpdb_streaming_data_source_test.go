package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb StreamingDataSource. >>> Resource test cases, automatically generated.
// Case StreamingDataSource_资源依赖_case_003 7178
func TestAccAliCloudGpdbStreamingDataSource_basic7178(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7178)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7178)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":     "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":   "test-kafka3",
					"data_source_config": "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
					"data_source_type":   "kafka",
					"service_id":         "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":     CHECKSET,
						"data_source_name":   "test-kafka3",
						"data_source_config": "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
						"data_source_type":   "kafka",
						"service_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_description": "test-kafka",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_description": "test-kafka",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_config": "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"#\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_config": "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"#\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_description": "test-kafka2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_description": "test-kafka2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
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

var AlicloudGpdbStreamingDataSourceMap7178 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudGpdbStreamingDataSourceBasicDependence7178(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "kafka-config-modify" {
  default = <<EOF
{
    "brokers": "alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092",
    "delimiter": "#",
    "format": "delimited",
    "topic": "ziyuan_test"
}
EOF
}

variable "kafka-config" {
  default = <<EOF
{
    "brokers": "alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092",
    "delimiter": "|",
    "format": "delimited",
    "topic": "ziyuan_test"
}
EOF
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultDfkYOR" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default59ZqyD" {
  vpc_id     = alicloud_vpc.defaultDfkYOR.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default7mX6ld" {
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
  vswitch_id                 = alicloud_vswitch.default59ZqyD.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultDfkYOR.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "defaultwruvdv" {
  service_name        = "test"
  db_instance_id      = alicloud_gpdb_instance.default7mX6ld.id
  service_description = "test"
  service_spec        = "8"
}


`, name)
}

// Case StreamingDataSource_资源依赖_case_002 7172
func TestAccAliCloudGpdbStreamingDataSource_basic7172(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7172)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":     "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":   "test-kafka3",
					"data_source_config": "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
					"data_source_type":   "kafka",
					"service_id":         "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":     CHECKSET,
						"data_source_name":   "test-kafka3",
						"data_source_config": "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
						"data_source_type":   "kafka",
						"service_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_description": "test-kafka",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_description": "test-kafka",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_config": "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"#\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_config": "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"#\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_description": "test-kafka2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_description": "test-kafka2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
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

var AlicloudGpdbStreamingDataSourceMap7172 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudGpdbStreamingDataSourceBasicDependence7172(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "kafka-config" {
  default = <<EOF
{
    "brokers": "alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092",
    "topic": "lineitem",
    "format": "delimited",
    "delimiter": "|"
}
EOF
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultDfkYOR" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default59ZqyD" {
  vpc_id     = alicloud_vpc.defaultDfkYOR.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default7mX6ld" {
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
  vswitch_id                 = alicloud_vswitch.default59ZqyD.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultDfkYOR.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "defaultwruvdv" {
  service_name        = "test"
  db_instance_id      = alicloud_gpdb_instance.default7mX6ld.id
  service_description = "test"
  service_spec        = "8"
}


`, name)
}

// Case StreamingDataSource_资源依赖_case_003 7178  twin
func TestAccAliCloudGpdbStreamingDataSource_basic7178_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7178)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7178)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
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

// Case StreamingDataSource_资源依赖_case_002 7172  twin
func TestAccAliCloudGpdbStreamingDataSource_basic7172_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7172)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
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

// Case StreamingDataSource_资源依赖_case_003 7178  raw
func TestAccAliCloudGpdbStreamingDataSource_basic7178_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7178)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7178)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"#\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"ziyuan_test\\\"}",
					"data_source_description": "test-kafka2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"#\",\"format\":\"delimited\",\"topic\":\"ziyuan_test\"}",
						"data_source_description": "test-kafka2",
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

// Case StreamingDataSource_资源依赖_case_002 7172  raw
func TestAccAliCloudGpdbStreamingDataSource_basic7172_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataSourceMap7172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataSourceBasicDependence7172)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_gpdb_instance.default7mX6ld.id}",
					"data_source_name":        "test-kafka3",
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"|\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
					"data_source_type":        "kafka",
					"data_source_description": "test-kafka",
					"service_id":              "${alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"data_source_name":        "test-kafka3",
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"|\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
						"data_source_type":        "kafka",
						"data_source_description": "test-kafka",
						"service_id":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_config":      "{\\\"brokers\\\":\\\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\\\",\\\"delimiter\\\":\\\"#\\\",\\\"format\\\":\\\"delimited\\\",\\\"topic\\\":\\\"lineitem\\\"}",
					"data_source_description": "test-kafka2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_config":      "{\"brokers\":\"alikafka-post-cn-uax3gim8q004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-uax3gim8q004-3-vpc.alikafka.aliyuncs.com:9092\",\"delimiter\":\"#\",\"format\":\"delimited\",\"topic\":\"lineitem\"}",
						"data_source_description": "test-kafka2",
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

// Test Gpdb StreamingDataSource. <<< Resource test cases, automatically generated.
