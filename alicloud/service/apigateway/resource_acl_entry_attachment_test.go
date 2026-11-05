package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudApiGatewayAclEntryAttachment(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_acl_entry_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayAclEntryAttachmentAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_example%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayAclEntryAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_api_gateway_access_control_list.default.id}",
					"entry":   "128.0.0.1/32",
					"comment": "test comment",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":  acctest.CHECKSET,
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
