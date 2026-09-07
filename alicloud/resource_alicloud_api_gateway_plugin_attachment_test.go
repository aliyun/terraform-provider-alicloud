package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudApiGatewayPluginAttachment(t *testing.T) {
	var v *cloudapi.PluginAttribute

	resourceId := "alicloud_api_gateway_plugin_attachment.default"
	ra := resourceAttrInit(resourceId, apiGatewayPluginAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccPlugin_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApiGatewayPluginAttachmentConfigDependence)

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
					"stage_name": "RELEASE",
					"plugin_id":  "${alicloud_api_gateway_plugin.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func resourceApiGatewayPluginAttachmentConfigDependence(name string) string {
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

  stage_names = [
    "RELEASE",
    "TEST",
  ]
}

resource "alicloud_api_gateway_plugin" "default" {
  plugin_name 	= "${var.name}"
  plugin_data 	= jsonencode({"allowOrigins": "api.foo.com","allowMethods": "GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH","allowHeaders": "Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid","exposeHeaders": "Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message","maxAge": 172800,"allowCredentials": true})
  plugin_type 	= "cors"
  description 	= "tf_testAccApiPlugin Description"
}

 `, name)
}

var apiGatewayPluginAttachmentBasicMap = map[string]string{
	"api_id":     CHECKSET,
	"group_id":   CHECKSET,
	"stage_name": "RELEASE",
	"plugin_id":  CHECKSET,
}
