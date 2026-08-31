package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MessageService EndpointAcl. >>> Resource test cases, automatically generated.
// Case EndpointAcl资源测试v1.0_对接TF 10161
func TestAccAliCloudMessageServiceEndpointAcl_basic10161(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_endpoint_acl.default"
	ra := resourceAttrInit(resourceId, AliCloudMessageServiceEndpointAclMap10161)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceEndpointAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smessageserviceendpointacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMessageServiceEndpointAclBasicDependence10161)
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
					"cidr":          "192.168.1.1/23",
					"endpoint_type": "${alicloud_message_service_endpoint.default.id}",
					"acl_strategy":  "allow",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr":          "192.168.1.1/23",
						"endpoint_type": CHECKSET,
						"acl_strategy":  "allow",
					}),
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

var AliCloudMessageServiceEndpointAclMap10161 = map[string]string{}

func AliCloudMessageServiceEndpointAclBasicDependence10161(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	resource "alicloud_message_service_endpoint" "default" {
  		endpoint_enabled = true
  		endpoint_type    = "public"
	}
`, name)
}

// Test MessageService EndpointAcl. <<< Resource test cases, automatically generated.
