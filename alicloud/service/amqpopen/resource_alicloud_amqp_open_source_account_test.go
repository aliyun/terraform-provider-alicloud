// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package amqpopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/amqpopen"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Amqp OpenSourceAccount. >>> Resource test cases, automatically generated.
// Case 开源用户 12792
func TestAccAliCloudAmqpOpenSourceAccount_basic12792(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_open_source_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAmqpOpenSourceAccountMap12792)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return amqpopen.NewAmqpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAmqpOpenSourceAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAmqpOpenSourceAccountBasicDependence12792)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_name":   "${var.user_name}",
					"description": "${var.user_name}",
					"password":    "${var.user_name}",
					"instance_id": "${data.alicloud_amqp_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   acctest.CHECKSET,
						"description": acctest.CHECKSET,
						"password":    acctest.CHECKSET,
						"instance_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.user_name_update}",
					"password":    "${var.user_name_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
						"password":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudAmqpOpenSourceAccountMap12792 = map[string]string{}

func AlicloudAmqpOpenSourceAccountBasicDependence12792(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "instance_name" {
  default = "测试开源鉴权实例"
}

variable "user_name" {
  default = "Suhao123_"
}

variable "user_name_update" {
  default = "Suhao456_"
}

data "alicloud_amqp_instances" "default" {
  status = "SERVING"
}


`, name)
}

// Test Amqp OpenSourceAccount. <<< Resource test cases, automatically generated.
