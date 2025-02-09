package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudFcv3TriggerMap6980 = map[string]string{

	"create_time":  CHECKSET,
	"trigger_name": CHECKSET,
}

func AlicloudFcv3TriggerBasicDependence6980(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_HTTP"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}
`, name)
}

var AlicloudFcv3TriggerMap6982 = map[string]string{

	"create_time":  CHECKSET,
	"trigger_name": CHECKSET,
}

func AlicloudFcv3TriggerBasicDependence6982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_Timer"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

`, name)
}

var AlicloudFcv3TriggerMap6981 = map[string]string{
	"create_time":  CHECKSET,
	"trigger_name": CHECKSET,
}

func AlicloudFcv3TriggerBasicDependence6981(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_OSS"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

`, name)
}

var AlicloudFcv3TriggerMap6983 = map[string]string{

	"create_time":  CHECKSET,
	"trigger_name": CHECKSET,
}

func AlicloudFcv3TriggerBasicDependence6983(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_Log"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = var.name
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_project" "default1" {
  project_name = format("%%supdate", var.name)
  description  = format("%%supdate", var.name)
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

`, name)
}

var AlicloudFcv3TriggerMap6985 = map[string]string{
	"create_time":  CHECKSET,
	"trigger_name": CHECKSET,
}

func AlicloudFcv3TriggerBasicDependence6985(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_CDN"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}


`, name)
}

// Case TestTrigger_HTTP 6980  raw
func TestAccAliCloudFcv3Trigger_basic6980_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_trigger.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3TriggerMap6980)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Trigger")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3trigger%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3TriggerBasicDependence6980)
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
					"function_name":  "${alicloud_fcv3_function.function.function_name}",
					"trigger_type":   "http",
					"trigger_name":   name,
					"description":    "create",
					"qualifier":      "LATEST",
					"trigger_config": "{\\\"methods\\\":[\\\"GET\\\",\\\"POST\\\",\\\"PUT\\\",\\\"DELETE\\\"],\\\"authType\\\":\\\"anonymous\\\",\\\"disableURLInternet\\\":false}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":  CHECKSET,
						"trigger_type":   "http",
						"trigger_name":   name,
						"description":    "create",
						"qualifier":      "LATEST",
						"trigger_config": "{\"methods\":[\"GET\",\"POST\",\"PUT\",\"DELETE\"],\"authType\":\"anonymous\",\"disableURLInternet\":false}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "update1",
					"trigger_config": "{\\\"methods\\\":[\\\"GET\\\",\\\"POST\\\",\\\"PUT\\\",\\\"DELETE\\\"],\\\"authType\\\":\\\"function\\\",\\\"disableURLInternet\\\":false}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "update1",
						"trigger_config": "{\"methods\":[\"GET\",\"POST\",\"PUT\",\"DELETE\"],\"authType\":\"function\",\"disableURLInternet\":false}",
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

// Case TestTrigger_Timer 6982  raw
func TestAccAliCloudFcv3Trigger_basic6982_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_trigger.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3TriggerMap6982)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Trigger")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3trigger%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3TriggerBasicDependence6982)
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
					"function_name":  "${alicloud_fcv3_function.function.function_name}",
					"trigger_type":   "timer",
					"trigger_name":   name,
					"description":    "create",
					"qualifier":      "LATEST",
					"trigger_config": "{\\\"payload\\\":\\\"hello\\\",\\\"cronExpression\\\":\\\"@every 1m\\\",\\\"enable\\\":true}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":  CHECKSET,
						"trigger_type":   "timer",
						"trigger_name":   name,
						"description":    "create",
						"qualifier":      "LATEST",
						"trigger_config": "{\"payload\":\"hello\",\"cronExpression\":\"@every 1m\",\"enable\":true}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "update",
					"trigger_config": "{\\\"payload\\\":\\\"hello1\\\",\\\"cronExpression\\\":\\\"@every 1m\\\",\\\"enable\\\":true}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "update",
						"trigger_config": "{\"payload\":\"hello1\",\"cronExpression\":\"@every 1m\",\"enable\":true}",
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

// Case TestTrigger_OSS 6981  raw
func TestAccAliCloudFcv3Trigger_basic6981_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_trigger.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3TriggerMap6981)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Trigger")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3trigger%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3TriggerBasicDependence6981)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":   "${alicloud_fcv3_function.function.function_name}",
					"trigger_type":    "oss",
					"trigger_name":    name,
					"description":     "create",
					"qualifier":       "LATEST",
					"trigger_config":  "{\\\"events\\\":[\\\"oss:ObjectCreated:PutObject\\\",\\\"oss:ObjectCreated:PostObject\\\",\\\"oss:ObjectCreated:CompleteMultipartUpload\\\",\\\"oss:ObjectCreated:PutSymlink\\\"],\\\"filter\\\":{\\\"key\\\":{\\\"prefix\\\":\\\"fc_oss_trigger_api_test\\\",\\\"suffix\\\":\\\".zip1\\\"}}}",
					"source_arn":      "acs:oss:cn-shanghai:1511928242963727:${alicloud_oss_bucket.default.bucket}",
					"invocation_role": "acs:ram::1511928242963727:role/aliyunosseventnotificationrole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   CHECKSET,
						"trigger_type":    "oss",
						"trigger_name":    name,
						"description":     "create",
						"qualifier":       "LATEST",
						"trigger_config":  "{\"events\":[\"oss:ObjectCreated:PutObject\",\"oss:ObjectCreated:PostObject\",\"oss:ObjectCreated:CompleteMultipartUpload\",\"oss:ObjectCreated:PutSymlink\"],\"filter\":{\"key\":{\"prefix\":\"fc_oss_trigger_api_test\",\"suffix\":\".zip1\"}}}",
						"source_arn":      CHECKSET,
						"invocation_role": "acs:ram::1511928242963727:role/aliyunosseventnotificationrole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":     "update",
					"trigger_config":  "{\\\"events\\\":[\\\"oss:ObjectCreated:PutObject\\\",\\\"oss:ObjectCreated:PostObject\\\",\\\"oss:ObjectCreated:CompleteMultipartUpload\\\",\\\"oss:ObjectCreated:PutSymlink\\\"],\\\"filter\\\":{\\\"key\\\":{\\\"prefix\\\":\\\"fc_oss_trigger_api_test\\\",\\\"suffix\\\":\\\".text1\\\"}}}",
					"invocation_role": "acs:ram::1511928242963727:role/aliyunosseventnotificationrole2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":     "update",
						"trigger_config":  "{\"events\":[\"oss:ObjectCreated:PutObject\",\"oss:ObjectCreated:PostObject\",\"oss:ObjectCreated:CompleteMultipartUpload\",\"oss:ObjectCreated:PutSymlink\"],\"filter\":{\"key\":{\"prefix\":\"fc_oss_trigger_api_test\",\"suffix\":\".text1\"}}}",
						"invocation_role": "acs:ram::1511928242963727:role/aliyunosseventnotificationrole2",
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

// Case TestTrigger_Log 6983  raw
func TestAccAliCloudFcv3Trigger_basic6983_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_trigger.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3TriggerMap6983)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Trigger")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3trigger%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3TriggerBasicDependence6983)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":   "${alicloud_fcv3_function.function.function_name}",
					"trigger_type":    "log",
					"trigger_name":    name,
					"description":     "create",
					"qualifier":       "LATEST",
					"trigger_config":  "{\\\"sourceConfig\\\":{\\\"logstore\\\":\\\"${alicloud_log_project.default1.project_name}\\\",\\\"startTime\\\":null},\\\"jobConfig\\\":{\\\"maxRetryTime\\\":3,\\\"triggerInterval\\\":60},\\\"functionParameter\\\":{},\\\"logConfig\\\":{\\\"project\\\":\\\"${alicloud_log_project.default.project_name}\\\",\\\"logstore\\\":\\\"${alicloud_log_store.default.logstore_name}\\\"},\\\"enable\\\":true}",
					"source_arn":      "acs:log:cn-shanghai:1511928242963727:project/${alicloud_log_project.default1.project_name}-cn-shanghai",
					"invocation_role": "acs:ram::1511928242963727:role/aliyunlogetlrole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   CHECKSET,
						"trigger_type":    "log",
						"trigger_name":    name,
						"description":     "create",
						"qualifier":       "LATEST",
						"trigger_config":  CHECKSET,
						"source_arn":      CHECKSET,
						"invocation_role": "acs:ram::1511928242963727:role/aliyunlogetlrole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "update",
					"trigger_config": "{\\\"sourceConfig\\\":{\\\"logstore\\\":\\\"function-log\\\",\\\"startTime\\\":null},\\\"jobConfig\\\":{\\\"maxRetryTime\\\":3,\\\"triggerInterval\\\":120},\\\"functionParameter\\\":{},\\\"logConfig\\\":{\\\"project\\\":\\\"fc3-api-1511928242963727-cn-shanghai\\\",\\\"logstore\\\":\\\"fc-trigger-log\\\"},\\\"enable\\\":true}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "update",
						"trigger_config": "{\"sourceConfig\":{\"logstore\":\"function-log\",\"startTime\":null},\"jobConfig\":{\"maxRetryTime\":3,\"triggerInterval\":120},\"functionParameter\":{},\"logConfig\":{\"project\":\"fc3-api-1511928242963727-cn-shanghai\",\"logstore\":\"fc-trigger-log\"},\"enable\":true}",
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

// Case TestTrigger_CDN 6985  raw
func TestAccAliCloudFcv3Trigger_basic6985_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_trigger.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3TriggerMap6985)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Trigger")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3trigger%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3TriggerBasicDependence6985)
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
					"function_name":   "${alicloud_fcv3_function.function.function_name}",
					"trigger_type":    "cdn_events",
					"trigger_name":    name,
					"description":     "create",
					"qualifier":       "LATEST",
					"trigger_config":  "{\\\"eventName\\\":\\\"CachedObjectsPushed\\\",\\\"eventVersion\\\":\\\"1.0.0\\\",\\\"notes\\\":\\\"test\\\",\\\"filter\\\":{\\\"domain\\\":[\\\"example.com\\\"]}}",
					"source_arn":      "acs:cdn:*:1511928242963727",
					"invocation_role": "acs:ram::1511928242963727:role/aliyuncdneventnotificationrole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   CHECKSET,
						"trigger_type":    "cdn_events",
						"trigger_name":    name,
						"description":     "create",
						"qualifier":       "LATEST",
						"trigger_config":  "{\"eventName\":\"CachedObjectsPushed\",\"eventVersion\":\"1.0.0\",\"notes\":\"test\",\"filter\":{\"domain\":[\"example.com\"]}}",
						"source_arn":      "acs:cdn:*:1511928242963727",
						"invocation_role": "acs:ram::1511928242963727:role/aliyuncdneventnotificationrole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "update",
					"trigger_config": "{\\\"eventName\\\":\\\"CachedObjectsPushed\\\",\\\"eventVersion\\\":\\\"1.0.0\\\",\\\"notes\\\":\\\"test\\\",\\\"filter\\\":{\\\"domain\\\":[\\\"example1.com\\\"]}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "update",
						"trigger_config": "{\"eventName\":\"CachedObjectsPushed\",\"eventVersion\":\"1.0.0\",\"notes\":\"test\",\"filter\":{\"domain\":[\"example1.com\"]}}",
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

// Test Fcv3 Trigger. <<< Resource test cases, automatically generated.
