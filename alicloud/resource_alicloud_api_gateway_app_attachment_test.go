package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func SkipTestAccAlicloudApigatewayAppAttachment(t *testing.T) {
	var v *cloudapi.AuthorizedApp

	resourceId := "alicloud_api_gateway_app_attachment.default"
	ra := resourceAttrInit(resourceId, apigatewayAppAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApp_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayAppAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"api_id":     "${alicloud_api_gateway_api.default.api_id}",
					"group_id":   "${alicloud_api_gateway_group.default.id}",
					"stage_name": "PRE",
					"app_id":     "${alicloud_api_gateway_app.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceApigatewayAppAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_api_gateway_group" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiGroup Description"
}
resource "alicloud_api_gateway_api" "default" {
  name        = "${var.name}"
  group_id    = "${alicloud_api_gateway_group.default.id}"
  description = "description"
  auth_type   = "APP"

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 22
    aone_name = "cloudapi-openapi"
  }

  request_parameters {
      name         = "aa"
      type         = "STRING"
      required     = "OPTIONAL"
      in           = "QUERY"
      in_service   = "QUERY"
      name_service = "testparams"
    }
}

resource "alicloud_api_gateway_app" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiAPP Description"
}

 `, name)
}

var apigatewayAppAttachmentBasicMap = map[string]string{
	"api_id":     CHECKSET,
	"group_id":   CHECKSET,
	"stage_name": "PRE",
	"app_id":     CHECKSET,
}
