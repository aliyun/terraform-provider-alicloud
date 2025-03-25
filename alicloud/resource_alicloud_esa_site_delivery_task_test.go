package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA SiteDeliveryTask. >>> Resource test cases, automatically generated.
// Case resource_SiteDeliveryTask_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   "dcdn-test20240417",
							"sls_region":    "cn-hongkong",
							"sls_log_store": "accesslog-test",
						},
					},
					"site_id":       "${alicloud_esa_site.resource_Site_task_test.id}",
					"data_center":   "cn",
					"discard_rate":  "0.0",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"discard_rate":  "0.1",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"field_name":    "OriginIP,ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"discard_rate":  "0.2",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "online",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sls_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_task_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_SiteDeliveryTask_http_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id": "${alicloud_esa_site.resource_Site_http_test.id}",
					"http_delivery": []map[string]interface{}{
						{
							"compress":        "gzip",
							"log_body_suffix": "cdnVersion:1.0",
							"standard_auth_param": []map[string]interface{}{
								{
									"private_key":  "***",
									"url_path":     "v1/log/upload",
									"expired_time": "300",
								},
							},
							"standard_auth_on":  "false",
							"log_body_prefix":   "cdnVersion:1.0",
							"dest_url":          "http://11.177.129.13:8081",
							"max_batch_size":    "1000",
							"max_retry":         "3",
							"transform_timeout": "10",
							"max_batch_mb":      "5",
						},
					},
					"data_center":   "oversea",
					"discard_rate":  "0.0",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_http_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_SiteDeliveryTask_oss_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${alicloud_esa_site.resource_Site_http_test.id}",
					"data_center":   "oversea",
					"discard_rate":  "0.0",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"oss_delivery": []map[string]interface{}{
						{
							"bucket_name": "test-log3",
							"region":      "cn-hongkong",
							"prefix_path": "",
							"aliuid":      "1097011697834102",
						},
					},
					"delivery_type": "oss",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oss_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_http_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_SiteDeliveryTask_kafka_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${alicloud_esa_site.resource_Site_http_test.id}",
					"data_center":   "oversea",
					"discard_rate":  "0.0",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
					"kafka_delivery": []map[string]interface{}{
						{
							"compress":       "lz4",
							"user_name":      "kafkaadmin",
							"machanism_type": "plain",
							"brokers": []string{
								"11.177.129.13:9092",
							},
							"balancer":  "kafka.LeastBytes",
							"topic":     "access",
							"user_auth": "true",
							"password":  "kafkaadmin-cdnlog",
						},
					},
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "kafka",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kafka_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_http_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_SiteDeliveryTask_aws3_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.resource_Site_http_test.id}",
					"data_center":  "oversea",
					"discard_rate": "0.0",
					"task_name":    "dcdn-test-task",
					"s3_delivery": []map[string]interface{}{
						{
							"secret_key":             "qYl*****PJeFh",
							"endpoint":               "https://s3.oss-cn-hangzhou.aliyuncs.com",
							"vertify_type":           "",
							"region":                 "us-east-2",
							"bucket_path":            "openapi-test-esa/test",
							"server_side_encryption": "false",
							"access_key":             "AK****YMC",
							"prefix_path":            "logriver-test/log",
							"s3_cmpt":                "true",
						},
					},
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "aws3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     "dcdn-test-task",
					"business_type": "dcdn_log_access_l1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_http_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA SiteDeliveryTask. <<< Resource test cases, automatically generated.
