package alicloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_fc_function",
		&resource.Sweeper{
			Name: "alicloud_fc_function",
			F:    testSweepFcFunction,
			Dependencies: []string{
				"alicloud_fc_trigger",
			},
		})
}

func testSweepFcFunction(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping FC unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.ListServices(fc.NewListServicesInput())
	})
	if err != nil {
		log.Printf("Error retrieving FC services: %s", err)
		return nil
	}

	swept := false
	services, _ := raw.(*fc.ListServicesOutput)
	for _, v := range services.Services {
		serviceName := v.ServiceName
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListFunctions(fc.NewListFunctionsInput(*serviceName))
		})
		if err != nil {
			log.Println(err.Error())
		} else {
			functions := raw.(*fc.ListFunctionsOutput)
			for _, v := range functions.Functions {
				functionName := v.FunctionName
				skip := true
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(*functionName), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					continue
				}
				swept = true

				_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
					return fcClient.DeleteFunction(&fc.DeleteFunctionInput{
						ServiceName:  serviceName,
						FunctionName: functionName,
					})
				})

				if err != nil {
					log.Printf("[ERROR] Failed to delete Api (%s): %s", *functionName, err)
				}
			}
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

var filePath string

func TestAccAlicloudFCFunction_basic(t *testing.T) {
	defer os.Remove(filePath)
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "nodejs12",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFCFunctionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "nodejs12",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "filename", "oss_bucket", "oss_key"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": map[string]string{
						"test":   "terraform",
						"prefix": "tfAcc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.test":   "terraform",
						"environment_variables.prefix": "tfAcc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "nodejs12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "nodejs12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "hello.initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "hello.initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "e1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "e1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "nodejs12",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
					// Check the function invocation result.
					func(*terraform.State) error {
						return checkInvocation(name, name, "hello world")
					},
				),
			},
		},
	})
}

func TestAccAlicloudFCFunctionMulti(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name + "-9",
		"runtime":     "nodejs12",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default.9"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFCFunctionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "10",
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}-${count.index}",
					"runtime":     "nodejs12",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudFCFunction_custom_container(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	basicMap := map[string]string{
		"service": CHECKSET,
		"name":    name,
		"runtime": CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFCFunctionConfigDependence)

	// REQUIREMENT: the image must be in the repo already.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service": "${alicloud_fc_service.default.name}",
					"name":    "${var.name}",
					"handler": "fake",
					"runtime": "custom-container",
					"custom_container_config": []map[string]string{
						{
							"image": fmt.Sprintf("registry.%s.aliyuncs.com/eci_open/nginx:alpine", os.Getenv("ALICLOUD_REGION")),
						},
					},
					"ca_port": "9527",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_container_config.0.image": fmt.Sprintf("registry.%s.aliyuncs.com/eci_open/nginx:alpine", os.Getenv("ALICLOUD_REGION")),
						"ca_port":                         "9527",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service": "${alicloud_fc_service.default.name}",
					"name":    "${var.name}",
					"handler": "fake",
					"runtime": "custom-container",
					"custom_container_config": []map[string]string{
						{
							"image":   fmt.Sprintf("registry.%s.aliyuncs.com/eci_open/nginx:alpine", os.Getenv("ALICLOUD_REGION")),
							"command": "${local.container_command}",
							"args":    "${local.container_args}",
						},
					},
					"ca_port": "9900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_container_config.0.image":   fmt.Sprintf("registry.%s.aliyuncs.com/eci_open/nginx:alpine", os.Getenv("ALICLOUD_REGION")),
						"custom_container_config.0.command": `["python", "server.py"]`,
						"custom_container_config.0.args":    `["a1", "a2"]`,
						"ca_port":                           "9900",
					}),
				),
			},
			// TODO: Check the invocation when supporing public image access.
		},
	})
}

func resourceFCFunctionConfigDependence(name string) string {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Printf("Failed to create temp directory: %s. Error: %v", dir, err)
		return ""
	}
	result := `'use strict';
	           exports.initializer = (context, callback) => {
		           console.log('hello init');
		           callback(null, 'hello init');
			   }
	           exports.handler = (event, context, callback) => {
		           console.log(event.toString())
		           callback(null, 'hello world');
			   }`
	filePath := filepath.Join(dir, "hello.js")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s. Error: %v", filePath, err)
		return ""
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		return ""
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
		return ""
	}

	return fmt.Sprintf(`
variable "name" {
    default = "%v"
}

// After serveral hours of investigation, finally figure out how to escape the double quotes.
// https://github.com/hashicorp/terraform/issues/17144
// https://discuss.hashicorp.com/t/how-can-i-escape-double-quotes-in-a-variable-value/4697
locals {
	container_command = "[\"python\", \"server.py\"]"
	container_args = "[\"a1\", \"a2\"]"
}

output "container_command" {
	value = "${local.container_command}"
}

output "container_args" {
	value = "${local.container_args}"
}

resource "alicloud_log_project" "default" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project = "${alicloud_log_project.default.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    description = "tf unit test"
    log_config {
	project = "${alicloud_log_project.default.name}"
	logstore = "${alicloud_log_store.default.name}"
    }
    role = "${alicloud_ram_role.default.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = "${alicloud_oss_bucket.default.id}"
  key = "fc/hello.zip"
  source = "%s"
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
resource "alicloud_ram_role_policy_attachment" "acr" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}
`, name, zipFilePath, testFCRoleTemplate)
}

func TestAccAlicloudFCFunction_custom_runtime(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	caPort := "9527"
	var basicMap = map[string]string{
		"service": CHECKSET,
		"name":    name,
		"runtime": "custom",
		"ca_port": caPort,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s. Error: %v", dir, err)
	}
	// Create the server file.
	result := fmt.Sprintf(`
<?php
$http = new swoole_http_server("0.0.0.0", %s);
$http->on("request", function ($request, $response) {
    $response->header("Content-Type", "text/plain");
    $response->end("hello custom runtime");
});
$http->start();
	`, caPort)
	filePath := filepath.Join(dir, "server.php")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %s. Error: %v", filePath, err)
	}
	// Create the bootstrap file.
	result = fmt.Sprintf(
		`#!/bin/bash
php server.php`)
	filePath = filepath.Join(dir, "bootstrap")
	err = ioutil.WriteFile(filePath, []byte(result), 0744)
	if err != nil {
		t.Fatalf("Failed to write file: %s. Error: %v", filePath, err)
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		t.Fatalf("Failed to zip directory: %s. Error: %v", dir, err)
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
	}

	dependence := func(name string) string {
		return fmt.Sprintf(`
		variable "name" {
		    default = "%v"
		}
		resource "alicloud_log_project" "default" {
		  name = "${var.name}"
		  description = "tf unit test"
		}
		
		resource "alicloud_log_store" "default" {
		  project = "${alicloud_log_project.default.name}"
		  name = "${var.name}"
		  retention_period = "3000"
		  shard_count = 1
		}
		resource "alicloud_fc_service" "default" {
		    name = "${var.name}"
		    description = "tf unit test"
		    log_config {
			project = "${alicloud_log_project.default.name}"
			logstore = "${alicloud_log_store.default.name}"
		    }
		    role = "${alicloud_ram_role.default.arn}"
		    depends_on = ["alicloud_ram_role_policy_attachment.default"]
		}
		
		resource "alicloud_ram_role" "default" {
		  name = "${var.name}"
		  document = <<EOF
		  %s
		  EOF
		  description = "this is a test"
		  force = true
		}
		
		data "alicloud_file_crc64_checksum" "default" {
			filename = "%s"
		}
		
		resource "alicloud_ram_role_policy_attachment" "default" {
		  role_name = "${alicloud_ram_role.default.name}"
		  policy_name = "AliyunLogFullAccess"
		  policy_type = "System"
		}
`, name, testFCRoleTemplate, zipFilePath)
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, dependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":  "${alicloud_fc_service.default.name}",
					"name":     "${var.name}",
					"runtime":  "custom",
					"handler":  "fake",
					"filename": zipFilePath,
					"ca_port":  caPort,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
					// Check the function invocation result.
					func(*terraform.State) error {
						return checkInvocation(name, name, "hello custom runtime")
					},
				),
			},
		},
	})
}

func TestAccAlicloudFCFunction_code_checksum(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "python2.7",
		"description": "tf",
		"handler":     "hello.handler",
		"filename":    CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	path, file, err := createTempFile(name)
	file.WriteString(`    # -*- coding: utf-8 -*-`)
	file.WriteString("\n")
	file.WriteString(`	def handler(event, context):`)
	file.WriteString("\n")
	file.WriteString(`	    print "hello world"`)
	file.WriteString("\n")

	if err != nil {
		t.Fatal(WrapError(err))
	}
	defer func() {
		file.Close()
		os.Remove(path)
	}()

	dependence := func(name string) string {
		return fmt.Sprintf(`
		variable "name" {
		    default = "%v"
		}
		resource "alicloud_log_project" "default" {
		  name = "${var.name}"
		  description = "tf unit test"
		}
		
		resource "alicloud_log_store" "default" {
		  project = "${alicloud_log_project.default.name}"
		  name = "${var.name}"
		  retention_period = "3000"
		  shard_count = 1
		}
		resource "alicloud_fc_service" "default" {
		    name = "${var.name}"
		    description = "tf unit test"
		    log_config {
			project = "${alicloud_log_project.default.name}"
			logstore = "${alicloud_log_store.default.name}"
		    }
		    role = "${alicloud_ram_role.default.arn}"
		    depends_on = ["alicloud_ram_role_policy_attachment.default"]
		}
		
		resource "alicloud_ram_role" "default" {
		  name = "${var.name}"
		  document = <<EOF
		  %s
		  EOF
		  description = "this is a test"
		  force = true
		}
		
		data "alicloud_file_crc64_checksum" "default" {
			filename = "%s"
		}
		
		resource "alicloud_ram_role_policy_attachment" "default" {
		  role_name = "${alicloud_ram_role.default.name}"
		  policy_name = "AliyunLogFullAccess"
		  policy_type = "System"
		}
`, name, testFCRoleTemplate, path)
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, dependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":       "${alicloud_fc_service.default.name}",
					"name":          "${var.name}",
					"runtime":       "python2.7",
					"description":   "tf",
					"handler":       "hello.handler",
					"filename":      path,
					"code_checksum": "${data.alicloud_file_crc64_checksum.default.checksum}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code_checksum": CHECKSET,
					}),
				),
			},
			{
				PreConfig: func() {
					file.WriteString(`	    return "hello world"`)
				},
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code_checksum": CHECKSET,
					}),
				),
			},
		},
	})
}

func createTempFile(prefix string) (string, *os.File, error) {
	f, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return "", nil, err
	}

	pathToFile, err := filepath.Abs(f.Name())
	if err != nil {
		return "", nil, err
	}
	return pathToFile, f, nil
}

func checkInvocation(service string, function string, gold string) (err error) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	res, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.InvokeFunction(fc.NewInvokeFunctionInput(service, function))
	})
	if err != nil {
		return err
	}
	str := string(res.(*fc.InvokeFunctionOutput).Payload[:])
	if gold != str {
		return Error(fmt.Sprintf("Expect invocation result: %s, but got: %s", gold, str))
	}
	return nil
}
