package mnsopen_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	mnsopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/mnsopen"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test MessageService EndpointAcl. >>> Resource test cases, automatically generated.
// Case EndpointAcl资源测试v1.0_对接TF 10161
func TestAccAliCloudMessageServiceEndpointAcl_basic10161(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_endpoint_acl.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudMessageServiceEndpointAclMap10161)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEndpointAcl")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smessageserviceendpointacl%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudMessageServiceEndpointAclBasicDependence10161)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"endpoint_type": tfacctest.CHECKSET,
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
