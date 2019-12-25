package alicloud

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
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
		return fmt.Errorf("Error retrieving FC services: %s", err)
	}

	swept := false
	services, _ := raw.(*fc.ListServicesOutput)
	for _, v := range services.Services {
		serviceName := v.ServiceName
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListFunctions(fc.NewListFunctionsInput(*serviceName))
		})
		if err != nil {
			fmt.Println(err.Error())
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

func TestAccAlicloudFCFunction_basic(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "python2.7",
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
					"service":     alicloud_fc_service.default.name,
					"name":        var.name,
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  alicloud_oss_bucket.default.id,
					"oss_key":     alicloud_oss_bucket_object.default.key,
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
					"runtime": "python3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "python3",
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
					"service":     alicloud_fc_service.default.name,
					"name":        var.name,
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  alicloud_oss_bucket.default.id,
					"oss_key":     alicloud_oss_bucket_object.default.key,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
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
		"runtime":     "python2.7",
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
					"service":     alicloud_fc_service.default.name,
					"name":        "${var.name}-${count.index}",
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  alicloud_oss_bucket.default.id,
					"oss_key":     alicloud_oss_bucket_object.default.key,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceFCFunctionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%v"
}
resource "alicloud_log_project" "default" {
  name = var.name
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name = var.name
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "default" {
    name = var.name
    description = "tf unit test"
    log_config {
	project = alicloud_log_project.default.name
	logstore = alicloud_log_store.default.name
    }
    role = alicloud_ram_role.default.arn
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_ram_role" "default" {
  name = var.name
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
`, name, testFCRoleTemplate)
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
		  name = var.name
		  description = "tf unit test"
		}
		
		resource "alicloud_log_store" "default" {
		  project = alicloud_log_project.default.name
		  name = var.name
		  retention_period = "3000"
		  shard_count = 1
		}
		resource "alicloud_fc_service" "default" {
		    name = var.name
		    description = "tf unit test"
		    log_config {
			project = alicloud_log_project.default.name
			logstore = alicloud_log_store.default.name
		    }
		    role = alicloud_ram_role.default.arn
		    depends_on = ["alicloud_ram_role_policy_attachment.default"]
		}
		
		resource "alicloud_ram_role" "default" {
		  name = var.name
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
		  role_name = alicloud_ram_role.default.name
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
					"service":       alicloud_fc_service.default.name,
					"name":          var.name,
					"runtime":       "python2.7",
					"description":   "tf",
					"handler":       "hello.handler",
					"filename":      path,
					"code_checksum": data.alicloud_file_crc64_checksum.default.checksum,
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
