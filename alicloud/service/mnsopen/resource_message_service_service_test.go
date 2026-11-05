// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

// Test MessageService Service. >>> Resource test cases, automatically generated.
// Case Service测试用例 10620
// lintignore: AT001
func TestAccAliCloudMessageServiceService_basic10620(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_service.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudMessageServiceServiceMap10620)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceService")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudMessageServiceServiceBasicDependence10620)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
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
	"status": tfacctest.CHECKSET,
}

func AliCloudMessageServiceServiceBasicDependence10620(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test MessageService Service. <<< Resource test cases, automatically generated.
