// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks Certificate. >>> Resource test cases, automatically generated.
// Case 测试Certificate1（有前置步骤，上海预发执行） 9054
func TestAccAliCloudDataWorksCertificate_basic9054(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksCertificateMap9054)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdataworks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksCertificateBasicDependence9054)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":       "${alicloud_data_works_project.defaultxWQGsP.id}",
					"name":             name,
					"certificate_file": "http://oxs-dataworks-openapi-console-cn-shanghai.oss-cn-shanghai.aliyuncs.com/d713_1107550004253538_8b2e08f49dd44414a8f27c4630b32c4b",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":       CHECKSET,
						"name":             CHECKSET,
						"certificate_file": "http://oxs-dataworks-openapi-console-cn-shanghai.oss-cn-shanghai.aliyuncs.com/d713_1107550004253538_8b2e08f49dd44414a8f27c4630b32c4b",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_file"},
			},
		},
	})
}

var AlicloudDataWorksCertificateMap9054 = map[string]string{
	"create_time":        CHECKSET,
	"create_user":        CHECKSET,
	"file_size_in_bytes": CHECKSET,
}

func AlicloudDataWorksCertificateBasicDependence9054(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}


`, name)
}

// Case 测试Certificate2（没有前置步骤，能够运行就行） 9055
func TestAccAliCloudDataWorksCertificate_basic9055(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksCertificateMap9055)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdataworks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksCertificateBasicDependence9055)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "tf-acc-cert",
					"project_id":       "${alicloud_data_works_project.defaultxWQGsP.id}",
					"name":             name,
					"certificate_file": "http://oxs-dataworks-openapi-console-cn-shanghai.oss-cn-shanghai.aliyuncs.com/d713_1107550004253538_8b2e08f49dd44414a8f27c4630b32c4b",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "tf-acc-cert",
						"project_id":       CHECKSET,
						"name":             CHECKSET,
						"certificate_file": "http://oxs-dataworks-openapi-console-cn-shanghai.oss-cn-shanghai.aliyuncs.com/d713_1107550004253538_8b2e08f49dd44414a8f27c4630b32c4b",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_file"},
			},
		},
	})
}

var AlicloudDataWorksCertificateMap9055 = map[string]string{
	"create_time":        CHECKSET,
	"create_user":        CHECKSET,
	"file_size_in_bytes": CHECKSET,
}

func AlicloudDataWorksCertificateBasicDependence9055(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}


`, name)
}

// Test DataWorks Certificate. <<< Resource test cases, automatically generated.
