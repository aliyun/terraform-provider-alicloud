package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApiGatewayGroupPluginAttachment(t *testing.T) {
	var v *cloudapi.PluginAttribute

	resourceId := "alicloud_api_gateway_group_plugin_attachment.default"
	ra := resourceAttrInit(resourceId, apiGatewayGroupPluginAttachmentBasicMap)

	//serviceFunc := func() interface{} {
	//	return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	//}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayGroupPluginAttachment")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccGroupPlugin_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApiGatewayGroupPluginAttachmentConfigDependence)

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

func resourceApiGatewayGroupPluginAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_api_gateway_group" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiGroup Description"
}

resource "alicloud_api_gateway_plugin" "default" {
  plugin_name 	= "${var.name}"
  plugin_data 	= jsonencode({"allowOrigins": "api.foo.com","allowMethods": "GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH","allowHeaders": "Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid","exposeHeaders": "Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message","maxAge": 172800,"allowCredentials": true})
  plugin_type 	= "cors"
  description 	= "tf_testAccApiPlugin Description"
}

 `, name)
}

var apiGatewayGroupPluginAttachmentBasicMap = map[string]string{
	"group_id":   CHECKSET,
	"stage_name": "RELEASE",
	"plugin_id":  CHECKSET,
}
