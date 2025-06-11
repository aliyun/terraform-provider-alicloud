// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MessageService Service. >>> Resource test cases, automatically generated.
// Case Service测试用例 10620
func TestAccAliCloudMessageServiceService_basic10620(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_service.default"
	ra := resourceAttrInit(resourceId, AliCloudMessageServiceServiceMap10620)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMessageServiceServiceBasicDependence10620)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AliCloudMessageServiceServiceMap10620 = map[string]string{
	"status": CHECKSET,
}

func AliCloudMessageServiceServiceBasicDependence10620(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test MessageService Service. <<< Resource test cases, automatically generated.
