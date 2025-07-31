package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudFc3FunctionMap6916 = map[string]string{
	"cpu":             "0.5",
	"function_name":   CHECKSET,
	"disk_size":       "512",
	"memory_size":     "512",
	"create_time":     CHECKSET,
	"internet_access": "true",
}

func AlicloudFc3FunctionBasicDependence6916(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
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

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
  depends_on = [alicloud_log_store.default]
}

resource "alicloud_oss_bucket" "default1" {
  bucket = format("%%supdate", var.name)
}

resource "alicloud_oss_bucket_object" "default1" {
  bucket  = alicloud_oss_bucket.default1.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
  depends_on = [alicloud_log_store.default1]
}

resource "alicloud_log_project" "default1" {
  project_name = format("%%supdate", var.name)
  description  = format("%%supdate", var.name)
}

resource "alicloud_log_store" "default1" {
  project_name          = alicloud_log_project.default1.name
  logstore_name         = format("%%supdate", var.name)
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
`, name)
}

var AlicloudFc3FunctionMap6950 = map[string]string{
	"cpu":                  "0.5",
	"function_name":        CHECKSET,
	"disk_size":            "512",
	"instance_concurrency": CHECKSET,
	"memory_size":          "512",
	"create_time":          CHECKSET,
	"internet_access":      "true",
}

func AlicloudFc3FunctionBasicDependence6950(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestCustomContainer_GPU"
}

variable "image1" {
  default = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/nginx:alpine"
}

variable "image2" {
  default = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30"
}


`, name)
}

var AlicloudFc3FunctionMap6917 = map[string]string{
	"cpu":                  "0.5",
	"function_name":        CHECKSET,
	"disk_size":            "512",
	"instance_concurrency": CHECKSET,
	"memory_size":          "512",
	"create_time":          CHECKSET,
	"internet_access":      "true",
}

func AlicloudFc3FunctionBasicDependence6917(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestCustomRuntime_Full"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}


`, name)
}

var AlicloudFc3FunctionMap6927 = map[string]string{
	"cpu":             "0.5",
	"function_name":   CHECKSET,
	"disk_size":       "512",
	"memory_size":     "512",
	"create_time":     CHECKSET,
	"internet_access": "true",
}

func AlicloudFc3FunctionBasicDependence6927(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestNativeRuntimePython_OSSMount"
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

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket" "default1" {
  bucket = format("%%supdate", var.name)
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}

resource "alicloud_oss_bucket_object" "default1" {
  bucket  = alicloud_oss_bucket.default1.bucket
  key     = "fc3Py39.zip"
  content = "print('hello world')"
}

`, name)
}

var AlicloudFc3FunctionMap6936 = map[string]string{
	"cpu":                  "0.5",
	"function_name":        CHECKSET,
	"disk_size":            "512",
	"instance_concurrency": CHECKSET,
	"memory_size":          "512",
	"create_time":          CHECKSET,
	"internet_access":      "true",
}

func AlicloudFc3FunctionBasicDependence6936(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "image1" {
  default = "registry-vpc.cn-shanghai.aliyuncs.com/fc-demo2/fc-demo:luoni-py39-http-v0.1"
}

variable "image2" {
  default = "registry-vpc.cn-shanghai.aliyuncs.com/fc-demo2/fc-demo:luoni-py39-http-v0.2"
}


`, name)
}

var AlicloudFc3FunctionMap6895 = map[string]string{
	"function_name":   CHECKSET,
	"memory_size":     "512",
	"create_time":     CHECKSET,
	"internet_access": "true",
}

func AlicloudFc3FunctionBasicDependence6895(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

var AlicloudFc3FunctionMap6938 = map[string]string{
	"cpu":             "0.5",
	"function_name":   CHECKSET,
	"disk_size":       "512",
	"memory_size":     "512",
	"create_time":     CHECKSET,
	"internet_access": "true",
}

func AlicloudFc3FunctionBasicDependence6938(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestNativeRuntimePython_VPC"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc-a" {
  description = "Alibaba-Fc-V3-Component-Generated"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vsw-a-1" {
  description  = "Alibaba-Fc-V3-Component-Generated"
  vpc_id       = alicloud_vpc.vpc-a.id
  zone_id      = "cn-hangzhou-h"
  cidr_block   = "10.20.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "vsw-a-2" {
  vpc_id       = alicloud_vpc.vpc-a.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "vsw-a-3" {
  vpc_id       = alicloud_vpc.vpc-a.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "10.2.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "sg-a" {
  name = var.name
  vpc_id              = alicloud_vpc.vpc-a.id
  security_group_type = "normal"
  inner_access_policy = "Accept"
}

resource "alicloud_vpc" "vpc-b" {
  cidr_block  = "192.168.0.0/16"
  vpc_name    = var.name
  description = "luoni.fz"
}

resource "alicloud_vswitch" "vsw-b-1" {
  vpc_id       = alicloud_vpc.vpc-b.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "192.168.1.0/24"
  vswitch_name = var.name
  description  = "luoni.fz"
}

resource "alicloud_security_group" "sg-b" {
  name = var.name
  vpc_id              = alicloud_vpc.vpc-b.id
  security_group_type = "normal"
  inner_access_policy = "Accept"
  description         = "luoni.fz"
}


`, name)
}

var AlicloudFc3FunctionMap7025 = map[string]string{
	"cpu":             "0.5",
	"function_name":   CHECKSET,
	"disk_size":       "512",
	"memory_size":     "512",
	"create_time":     CHECKSET,
	"internet_access": "true",
}

func AlicloudFc3FunctionBasicDependence7025(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestNativeRuntimePython_Nas_Pre2"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc-a" {
  description = "Alibaba-Fc-V3-Component-Generated"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vsw-a-1" {
  description  = "Alibaba-Fc-V3-Component-Generated"
  vpc_id       = alicloud_vpc.vpc-a.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "10.20.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "sg-a" {
  name = var.name
  vpc_id              = alicloud_vpc.vpc-a.id
  security_group_type = "normal"
  inner_access_policy = "Accept"
}

resource "alicloud_nas_file_system" "fs1" {
  storage_type     = "Performance"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
  description      = "Alibaba-Fc-V3-Component-Generated"
}

resource "alicloud_nas_access_group" "default" {
  access_group_name = var.name
  access_group_type = "Vpc"
  description       = "test_access_group"
  file_system_type  = "standard"
}

resource "alicloud_nas_mount_target" "mg1" {
  vpc_id         = alicloud_vpc.vpc-a.id
  vswitch_id     = alicloud_vswitch.vsw-a-1.id
  network_type   = "Vpc"
  file_system_id = alicloud_nas_file_system.fs1.id
  access_group_name = alicloud_nas_access_group.default.access_group_name
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

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}

`, name)
}

// Case TestNativeRuntimePython_Full 6916  raw
func TestAccAliCloudFcv3Function_basic6916_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6916)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6916)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"timeout":       "6",
					"handler":       "index.handler",
					"invocation_restriction": []map[string]interface{}{
						{
							"disable": "true",
							"reason":  "test",
						},
					},
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.key}",
							"checksum":        "4270285996107335518",
						},
					},
					"description": "Create",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "1",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "1",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":             "0.5",
					"disk_size":       "512",
					"internet_access": "true",
					"log_config": []map[string]interface{}{
						{
							"project":                 "${alicloud_log_project.default.project_name}",
							"logstore":                "${alicloud_log_store.default.logstore_name}",
							"log_begin_rule":          "None",
							"enable_instance_metrics": "false",
							"enable_request_metrics":  "false",
						},
					},
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"8.8.8.8", "100.100.2.136", "100.100.2.138"},
							"searches": []string{
								"ns1.svc.cluster-domain.example", "example.com", "mydomain.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  "timeout",
									"value": "2",
								},
								{
									"name":  "attempts",
									"value": "2",
								},
								{
									"name":  "ndots",
									"value": "2",
								},
							},
						},
					},
					"role": "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   name,
						"memory_size":     "512",
						"runtime":         "python3.9",
						"timeout":         "6",
						"handler":         "index.handler",
						"description":     "Create",
						"cpu":             "0.5",
						"disk_size":       "512",
						"internet_access": "true",
						"role":            "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"runtime":     "python3.10",
					"handler":     "index.HandlerX",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default1.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default1.key}",
							"checksum":        "4270285996107335518",
						},
					},
					"invocation_restriction": []map[string]interface{}{
						{
							"disable": "false",
						},
					},
					"description": "update",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal2",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.initializer",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.prestop",
								},
							},
						},
					},
					"cpu":       "1",
					"disk_size": "10240",
					"log_config": []map[string]interface{}{
						{
							"enable_instance_metrics": "true",
							"enable_request_metrics":  "true",
							"project":                 "${alicloud_log_project.default1.project_name}",
							"logstore":                "${alicloud_log_store.default1.logstore_name}",
							"log_begin_rule":          "DefaultRegex",
						},
					},
					"custom_dns": []map[string]interface{}{
						{
							"dns_options": []map[string]interface{}{
								{
									"name":  "timeout",
									"value": "1",
								},
								{
									"name":  "attempts",
									"value": "1",
								},
							},
							"name_servers": []string{
								"8.8.8.8"},
							"searches": []string{
								"example.com"},
						},
					},
					"role": "acs:ram::1511928242963727:role/FC3APIResourceTest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
						"runtime":     "python3.10",
						"handler":     "index.HandlerX",
						"description": "update",
						"cpu":         "1",
						"disk_size":   "10240",
						"role":        "acs:ram::1511928242963727:role/FC3APIResourceTest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"timeout":     "3",
					"invocation_restriction": []map[string]interface{}{
						{
							"disable": "true",
							"reason":  "test-change",
						},
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.start",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.prestop",
								},
							},
						},
					},
					"cpu":             "0.5",
					"disk_size":       "512",
					"internet_access": "false",
					//"custom_dns": []map[string]interface{}{
					//	{
					//		"dns_options": []map[string]interface{}{},
					//	},
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size":     "512",
						"timeout":         "3",
						"cpu":             "0.5",
						"disk_size":       "512",
						"internet_access": "false",
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
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestCustomContainer_GPU 6950  raw
func TestAccAliCloudFcv3Function_basic6950_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6950)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6950)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "custom-container",
					"timeout":       "3",
					"handler":       "index.handler",
					"description":   "Create",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":       "0.5",
					"disk_size": "512",
					"custom_container_config": []map[string]interface{}{
						{
							"image": "${var.image1}",
							"port":  "9000",
							"entrypoint": []string{
								"python", "-c", "xx"},
							"command": []string{
								"xx", "x", "a"},
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "1",
									"http_get_url":          "/ready",
									"initial_delay_seconds": "1",
									"period_seconds":        "1",
									"success_threshold":     "1",
									"timeout_seconds":       "1",
								},
							},
						},
					},
					"gpu_config": []map[string]interface{}{
						{
							"gpu_type":        "fc.gpu.tesla.1",
							"gpu_memory_size": "16384",
						},
					},
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "custom-container",
						"timeout":       "3",
						"handler":       "index.handler",
						"description":   "Create",
						"cpu":           "0.5",
						"disk_size":     "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"timeout":     "6",
					"handler":     "index.Handler",
					"description": "update",
					"cpu":         "1",
					"custom_container_config": []map[string]interface{}{
						{
							"image": "${var.image2}",
							"port":  "9001",
							"entrypoint": []string{
								"python3"},
							"command": []string{
								"app.py"},
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "2",
									"http_get_url":          "/readyz",
									"initial_delay_seconds": "2",
									"period_seconds":        "2",
									"success_threshold":     "2",
									"timeout_seconds":       "2",
								},
							},
						},
					},
					"gpu_config": []map[string]interface{}{
						{
							"gpu_type":        "fc.gpu.ampere.1",
							"gpu_memory_size": "24576",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
						"timeout":     "6",
						"handler":     "index.Handler",
						"description": "update",
						"cpu":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"cpu":         "0.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
						"cpu":         "0.5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestCustomRuntime_Full 6917  raw
func TestAccAliCloudFcv3Function_basic6917_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6917)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6917)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "custom.debian10",
					"timeout":       "3",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.key}",
							"checksum":        "4270285996107335518",
						},
					},
					"description": "Create",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "1",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "1",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":             "0.5",
					"disk_size":       "512",
					"internet_access": "true",
					"layers": []string{
						"acs:fc:cn-shanghai:official:layers/Python39-Aliyun-SDK/versions/3"},
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{
								"python", "-c", "test"},
							"args": []string{
								"app.py", "xx", "x"},
							"port": "9000",
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "3",
									"http_get_url":          "/ready",
									"initial_delay_seconds": "1",
									"period_seconds":        "10",
									"success_threshold":     "1",
									"timeout_seconds":       "1",
								},
							},
						},
					},
					"instance_concurrency": "2",
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        name,
						"memory_size":          "512",
						"runtime":              "custom.debian10",
						"timeout":              "3",
						"handler":              "index.handler",
						"description":          "Create",
						"cpu":                  "0.5",
						"disk_size":            "512",
						"internet_access":      "true",
						"layers.#":             "1",
						"instance_concurrency": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"timeout":     "6",
					"handler":     "index.HandlerX",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.key}",
							"checksum":        "4270285996107335518",
						},
					},
					"description": "update",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal2",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.initializer",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.prestop",
								},
							},
						},
					},
					"cpu":       "1",
					"disk_size": "10240",
					"layers": []string{
						"acs:fc:cn-shanghai:official:layers/Python39-Aliyun-SDK/versions/3"},
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{
								"python3"},
							"args": []string{
								"server.py"},
							"port": "9001",
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "5",
									"http_get_url":          "/readyx",
									"initial_delay_seconds": "2",
									"period_seconds":        "2",
									"success_threshold":     "2",
									"timeout_seconds":       "2",
								},
							},
						},
					},
					"instance_concurrency": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size":          "1024",
						"timeout":              "6",
						"handler":              "index.HandlerX",
						"description":          "update",
						"cpu":                  "1",
						"disk_size":            "10240",
						"layers.#":             "1",
						"instance_concurrency": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"runtime":     "custom",
					"timeout":     "60",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":       "0.5",
					"disk_size": "512",
					"layers":    []string{},
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{},
							"args":    []string{},
							"port":    "9000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
						"runtime":     "custom",
						"timeout":     "60",
						"cpu":         "0.5",
						"disk_size":   "512",
						"layers.#":    "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestNativeRuntimePython_OSSMount 6927  raw
func TestAccAliCloudFcv3Function_basic6927_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6927)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6927)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"timeout":       "6",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.key}",
							"checksum":        "4270285996107335518",
						},
					},
					"description": "Create",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
						},
					},
					"cpu":             "0.5",
					"disk_size":       "512",
					"internet_access": "true",
					"log_config": []map[string]interface{}{
						{
							"project":                 "${alicloud_log_project.default.name}",
							"logstore":                "${alicloud_log_store.default.logstore_name}",
							"log_begin_rule":          "None",
							"enable_instance_metrics": "false",
							"enable_request_metrics":  "false",
						},
					},
					"oss_mount_config": []map[string]interface{}{
						{
							"mount_points": []map[string]interface{}{
								{
									"bucket_name": "${alicloud_oss_bucket.default.bucket}",
									"bucket_path": "/test",
									"endpoint":    "http://oss-" + defaultRegionToTest + "-internal.aliyuncs.com",
									"mount_dir":   "/mnt1",
									"read_only":   "false",
								},
								{
									"bucket_name": "${alicloud_oss_bucket.default1.bucket}",
									"bucket_path": "/test2",
									"endpoint":    "http://oss-" + defaultRegionToTest + "-internal.aliyuncs.com",
									"mount_dir":   "/mnt2",
									"read_only":   "false",
								},
							},
						},
					},
					"role": "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   name,
						"memory_size":     "512",
						"runtime":         "python3.9",
						"timeout":         "6",
						"handler":         "index.handler",
						"description":     "Create",
						"cpu":             "0.5",
						"disk_size":       "512",
						"internet_access": "true",
						"role":            "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "3",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default1.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default1.key}",
							"checksum":        "11178804849154173829",
						},
					},
					"oss_mount_config": []map[string]interface{}{
						{
							"mount_points": []map[string]interface{}{
								{
									"bucket_name": "${alicloud_oss_bucket.default.bucket}",
									"bucket_path": "/test3",
									"endpoint":    "http://oss-" + defaultRegionToTest + ".aliyuncs.com",
									"mount_dir":   "/mnt4",
									"read_only":   "true",
								},
								{
									"bucket_name": "${alicloud_oss_bucket.default1.bucket}",
									"bucket_path": "/test2",
									"endpoint":    "http://oss-" + defaultRegionToTest + "-internal.aliyuncs.com",
									"mount_dir":   "/mnt2",
									"read_only":   "false",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestCustomContainer_Base 6936  raw
func TestAccAliCloudFcv3Function_basic6936_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6936)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6936)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "custom-container",
					"timeout":       "3",
					"handler":       "index.handler",
					"description":   "Create",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":       "0.5",
					"disk_size": "512",
					"custom_container_config": []map[string]interface{}{
						{
							"image": "${var.image1}",
							"port":  "9000",
							"entrypoint": []string{
								"python", "-c", "xx"},
							"command": []string{
								"xx", "x", "a"},
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "1",
									"http_get_url":          "/ready",
									"initial_delay_seconds": "1",
									"period_seconds":        "1",
									"success_threshold":     "1",
									"timeout_seconds":       "1",
								},
							},
						},
					},
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "custom-container",
						"timeout":       "3",
						"handler":       "index.handler",
						"description":   "Create",
						"cpu":           "0.5",
						"disk_size":     "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"timeout":     "6",
					"handler":     "index.Handler",
					"description": "update",
					"cpu":         "1",
					"custom_container_config": []map[string]interface{}{
						{
							"image": "${var.image2}",
							"port":  "9001",
							"entrypoint": []string{
								"python3"},
							"command": []string{
								"app.py"},
							"health_check_config": []map[string]interface{}{
								{
									"failure_threshold":     "2",
									"http_get_url":          "/readyz",
									"initial_delay_seconds": "2",
									"period_seconds":        "2",
									"success_threshold":     "2",
									"timeout_seconds":       "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
						"timeout":     "6",
						"handler":     "index.Handler",
						"description": "update",
						"cpu":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"cpu":         "0.5",
					//"custom_container_config": []map[string]interface{}{
					//	{
					//		"entrypoint": []string{},
					//		"command":    []string{},
					//	},
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
						"cpu":         "0.5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestNativeRuntimePython_Base 6895  raw
func TestAccAliCloudFcv3Function_basic6895_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6895)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(0, 9999999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6895)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "python3.9",
						"handler":       "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
					"description": "change_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "python3.9",
						"handler":       "index.handler",
						"description":   "change_description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role": "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": "acs:ram::1511928242963727:role/AliyunFCDefaultRole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestNativeRuntimePython_VPC 6938  raw
func TestAccAliCloudFcv3Function_basic6938_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap6938)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence6938)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"timeout":       "3",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
					"description": "Create",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":       "0.5",
					"disk_size": "512",
					"vpc_config": []map[string]interface{}{
						{
							"security_group_id": "${alicloud_security_group.sg-a.id}",
							"vpc_id":            "${alicloud_vpc.vpc-a.id}",
							"vswitch_ids": []string{
								"${alicloud_vswitch.vsw-a-1.id}", "${alicloud_vswitch.vsw-a-2.id}", "${alicloud_vswitch.vsw-a-3.id}"},
						},
					},
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "python3.9",
						"timeout":       "3",
						"handler":       "index.handler",
						"description":   "Create",
						"cpu":           "0.5",
						"disk_size":     "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"runtime":     "python3.10",
					"timeout":     "6",
					"handler":     "index.Handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZJBEncQRhGJgbuHQvPelK0lKpCp3KjHoUz+KdvIJkZHTQhcv+/fiPXxlWVys0GqJ0W8zWrjZL4uIwajKwdl2U7vx8mlScy/CgIPF7JlhPiBIteo6vlNCStzkRit5snLZ13ROPlef4MkvV6FAHbeaBxB4DmY9cr+82tzebqreBS5dhZPITIQ4j04L9FczSWFSBn7An1uPH+5vLEKi9xIpGxejZyq3LgNMMStid91Qd2f0pK8oLoIrSapF/90Tp8tK5pbv3EphSQQcSu8ZPPZCBDobd6TgVqw/TF1W6f8S/tD0xK46aOOTLZyKbk+Ayxzr/DAAA//9QSwcIBlYFIgIBAACxAQAAUEsBAhQAFAAIAAgAAAAAAAZWBSICAQAAsQEAAAgAAAAAAAAAAAAAAAAAAAAAAGluZGV4LnB5UEsFBgAAAAABAAEANgAAADgBAAAAAA==",
						},
					},
					"description": "update",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal2",
					},
					"cpu": "1",
					"vpc_config": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.vpc-b.id}",
							"vswitch_ids": []string{
								"${alicloud_vswitch.vsw-b-1.id}"},
							"security_group_id": "${alicloud_security_group.sg-b.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
						"runtime":     "python3.10",
						"timeout":     "6",
						"handler":     "index.Handler",
						"description": "update",
						"cpu":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"cpu":         "0.5",
					"vpc_config": []map[string]interface{}{
						{
							"vswitch_ids": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
						"cpu":         "0.5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Case TestNativeRuntimePython_Nas 7025  raw
func TestAccAliCloudFcv3Function_basic7025_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3FunctionMap7025)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfc3function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3FunctionBasicDependence7025)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"timeout":       "3",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.key}",
						},
					},
					"description": "Create",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":       "0.5",
					"disk_size": "512",
					"vpc_config": []map[string]interface{}{
						{
							"security_group_id": "${alicloud_security_group.sg-a.id}",
							"vpc_id":            "${alicloud_vpc.vpc-a.id}",
							"vswitch_ids": []string{
								"${alicloud_vswitch.vsw-a-1.id}"},
						},
					},
					"nas_config": []map[string]interface{}{
						{
							"group_id": "0",
							"user_id":  "0",
							"mount_points": []map[string]interface{}{
								{
									"server_addr": "${alicloud_nas_mount_target.mg1.mount_target_domain}:/luoni-test",
									"mount_dir":   "/luoni-test",
									"enable_tls":  "false",
								},
								{
									"server_addr": "${alicloud_nas_mount_target.mg1.mount_target_domain}:/luoni-test1",
									"mount_dir":   "/luoni-test1",
									"enable_tls":  "false",
								},
								{
									"server_addr": "${alicloud_nas_mount_target.mg1.mount_target_domain}:/luoni-test2",
									"mount_dir":   "/luoni-test2",
									"enable_tls":  "false",
								},
							},
						},
					},
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"memory_size":   "512",
						"runtime":       "python3.9",
						"timeout":       "3",
						"handler":       "index.handler",
						"description":   "Create",
						"cpu":           "0.5",
						"disk_size":     "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
					"runtime":     "python3.10",
					"timeout":     "6",
					"handler":     "index.Handler",
					"description": "update",
					"cpu":         "1",
					"vpc_config": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.vpc-a.id}",
							"vswitch_ids": []string{
								"${alicloud_vswitch.vsw-a-1.id}"},
							"security_group_id": "${alicloud_security_group.sg-a.id}",
						},
					},
					"nas_config": []map[string]interface{}{
						{
							"group_id": "1",
							"user_id":  "1",
							"mount_points": []map[string]interface{}{
								{
									"server_addr": "${alicloud_nas_mount_target.mg1.mount_target_domain}:/luoni-test3",
									"mount_dir":   "/luoni-test3",
									"enable_tls":  "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
						"runtime":     "python3.10",
						"timeout":     "6",
						"handler":     "index.Handler",
						"description": "update",
						"cpu":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
					"cpu":         "0.5",
					"vpc_config": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.vpc-a.id}",
							"vswitch_ids": []string{
								"${alicloud_vswitch.vsw-a-1.id}"},
							"security_group_id": "${alicloud_security_group.sg-a.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
						"cpu":         "0.5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "layers"},
			},
		},
	})
}

// Test Fc3 Function. <<< Resource test cases, automatically generated.
// Case TestNativeRuntimePython_Base_FC_Session_Feature 6895
func TestAccAliCloudFcv3Function_basic6895(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3FunctionMap6895)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3FunctionBasicDependence6895)
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
					"function_name": name,
					"memory_size":   "512",
					"runtime":       "python3.9",
					"timeout":       "3",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
					"description": "Create",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop",
								},
							},
						},
					},
					"cpu":                     "0.5",
					"disk_size":               "512",
					"session_affinity_config": "{\\\"sessionConcurrencyPerInstance\\\":20,\\\"sseEndpointPath\\\":\\\"sse\\\"}",
					"instance_isolation_mode": "SHARE",
					"session_affinity":        "MCP_SSE",
					// there is a bug in api,
					"log_config": []map[string]interface{}{
						{
							"log_begin_rule": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":           name,
						"memory_size":             "512",
						"runtime":                 "python3.9",
						"timeout":                 "3",
						"handler":                 "index.handler",
						"description":             "Create",
						"cpu":                     CHECKSET,
						"disk_size":               "512",
						"session_affinity_config": CHECKSET,
						"instance_isolation_mode": "SHARE",
						"session_affinity":        "MCP_SSE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "2048",
					"runtime":     "python3.10",
					"timeout":     "6",
					"handler":     "index.HandlerX",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZJBEncQRhGJgbuHQvPelK0lKpCp3KjHoUz+KdvIJkZHTQhcv+/fiPXxlWVys0GqJ0W8zWrjZL4uIwajKwdl2U7vx8mlScy/CgIPF7JlhPiBIteo6vlNCStzkRit5snLZ13ROPlef4MkvV6FAHbeaBxB4DmY9cr+82tzebqreBS5dhZPITIQ4j04L9FczSWFSBn7An1uPH+5vLEKi9xIpGxejZyq3LgNMMStid91Qd2f0pK8oLoIrSapF/90Tp8tK5pbv3EphSQQcSu8ZPPZCBDobd6TgVqw/TF1W6f8S/tD0xK46aOOTLZyKbk+Ayxzr/DAAA//9QSwcIBlYFIgIBAACxAQAAUEsBAhQAFAAIAAgAAAAAAAZWBSICAQAAsQEAAAgAAAAAAAAAAAAAAAAAAAAAAGluZGV4LnB5UEsFBgAAAAABAAEANgAAADgBAAAAAA==",
						},
					},
					"description": "update",
					"environment_variables": map[string]interface{}{
						"\"TestEnvKey\"": "TestEnvVal2",
					},
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"initializer": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.init2",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"timeout": "3",
									"handler": "index.stop2",
								},
							},
						},
					},
					"cpu":                     "2",
					"session_affinity_config": "{\\\"sessionConcurrencyPerInstance\\\":1,\\\"sseEndpointPath\\\":\\\"sse\\\"}",
					"instance_isolation_mode": "SESSION_EXCLUSIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size":             "2048",
						"runtime":                 "python3.10",
						"timeout":                 "6",
						"handler":                 "index.HandlerX",
						"description":             "update",
						"cpu":                     "2",
						"session_affinity_config": CHECKSET,
						"instance_isolation_mode": "SESSION_EXCLUSIVE",
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
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

var AlicloudFcv3FunctionMap6895 = map[string]string{
	"tracing_config.#":   CHECKSET,
	"function_id":        CHECKSET,
	"function_arn":       CHECKSET,
	"create_time":        CHECKSET,
	"code_size":          CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFcv3FunctionBasicDependence6895(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Fcv3 Function. <<< Resource test cases, automatically generated.
