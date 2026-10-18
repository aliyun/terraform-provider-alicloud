// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Operation. >>> Resource test cases.
func TestAccAliCloudApigOperation_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_operation.default"
	ra := resourceAttrInit(resourceId, AlicloudApigOperationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigOperation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigOperationBasicDependence)
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
					"http_api_id":    "${alicloud_apig_http_api.operation_httpapi.id}",
					"operation_name": name,
					"path":           "/op-basic",
					"method":         "GET",
					"description":    "basic operation description",
					"mock": []map[string]interface{}{
						{
							"enable":           "true",
							"response_code":    "200",
							"response_content": "hello world",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_name":          name,
						"path":                    "/op-basic",
						"method":                  "GET",
						"description":             "basic operation description",
						"http_api_id":             CHECKSET,
						"operation_id":            CHECKSET,
						"create_time":             CHECKSET,
						"mock.0.enable":           "true",
						"mock.0.response_code":    "200",
						"mock.0.response_content": "hello world",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_name": fmt.Sprintf("%s-updated", name),
					"description":    "updated description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_name": fmt.Sprintf("%s-updated", name),
						"description":    "updated description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path":   "/updated-path",
					"method": "POST",
					"mock": []map[string]interface{}{
						{
							"enable":           "false",
							"response_code":    "404",
							"response_content": "not found",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path":                    "/updated-path",
						"method":                  "POST",
						"mock.0.enable":           "false",
						"mock.0.response_code":    "404",
						"mock.0.response_content": "not found",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudApigOperationMap = map[string]string{
	"operation_id": CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudApigOperationBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_apig_http_api" "operation_httpapi" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "operation test httpapi"
  base_path     = "/${var.name}"
}

`, name)
}

// Test Apig Operation. <<< Resource test cases, automatically generated.
