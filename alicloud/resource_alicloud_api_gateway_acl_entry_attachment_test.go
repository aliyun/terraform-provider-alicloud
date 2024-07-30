package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAliCloudApiGatewayAclEntryAttachment(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_acl_entry_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayAclEntryAttachmentAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_example%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayAclEntryAttachmentBasicDependence)
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
					"acl_id":  "${alicloud_api_gateway_access_control_list.default.id}",
					"entry":   "128.0.0.1/32",
					"comment": "test comment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":  CHECKSET,
						"entry":   "128.0.0.1/32",
						"comment": "test comment",
					}),
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

func AlicloudApiGatewayAclEntryAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_api_gateway_instance" "default" {
  instance_name = var.name
  instance_spec = "api.s1.small"
  https_policy  = "HTTPS2_TLS1_0"
  zone_id       = "cn-hangzhou-MAZ6"
  payment_type  = "PayAsYouGo"
}

resource "alicloud_api_gateway_access_control_list" "default" {
  access_control_list_name = var.name
  address_ip_version 	   = "ipv4"
  depends_on               = ["alicloud_api_gateway_instance.default"]
}
`, name)
}
