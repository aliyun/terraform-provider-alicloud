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

// Test Amqp OpenSourcePermission. >>> Resource test cases, automatically generated.
// Case 开源权限_副本1779693632554 12809
func TestAccAliCloudAmqpOpenSourcePermission_basic12809(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_open_source_permission.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAmqpOpenSourcePermissionMap12809)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return amqpopen.NewAmqpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAmqpOpenSourcePermission")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAmqpOpenSourcePermissionBasicDependence12809)
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
					"write":       ".*",
					"read":        ".*",
					"vhost":       "${var.vhost}",
					"user_name":   "${var.user_name}",
					"instance_id": "${data.alicloud_amqp_instances.default.ids.0}",
					"configure":   ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write":       ".*",
						"read":        ".*",
						"vhost":       acctest.CHECKSET,
						"user_name":   acctest.CHECKSET,
						"instance_id": acctest.CHECKSET,
						"configure":   ".*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"write":     "^$",
					"read":      "^$",
					"configure": "^$",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write":     "^$",
						"read":      "^$",
						"configure": "^$",
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

var AlicloudAmqpOpenSourcePermissionMap12809 = map[string]string{}

func AlicloudAmqpOpenSourcePermissionBasicDependence12809(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "instance_name" {
  default = "测试开源鉴权实例"
}

variable "vhost" {
  default = "/"
}

variable "user_name" {
  default = "Suhao123_WithPer"
}

data "alicloud_amqp_instances" "default" {
  status = "SERVING"
}


`, name)
}

// Case 开源权限 12794
func TestAccAliCloudAmqpOpenSourcePermission_basic12794(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_open_source_permission.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAmqpOpenSourcePermissionMap12794)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return amqpopen.NewAmqpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAmqpOpenSourcePermission")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAmqpOpenSourcePermissionBasicDependence12794)
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
					"write":       ".*",
					"read":        ".*",
					"vhost":       "${var.vhost}",
					"user_name":   "${var.user_name}",
					"instance_id": "${data.alicloud_amqp_instances.default.ids.0}",
					"configure":   ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write":       ".*",
						"read":        ".*",
						"vhost":       acctest.CHECKSET,
						"user_name":   acctest.CHECKSET,
						"instance_id": acctest.CHECKSET,
						"configure":   ".*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudAmqpOpenSourcePermissionMap12794 = map[string]string{}

func AlicloudAmqpOpenSourcePermissionBasicDependence12794(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "instance_name" {
  default = "测试开源鉴权实例"
}

variable "vhost" {
  default = "/"
}

variable "user_name" {
  default = "Suhao123_WithPer"
}

data "alicloud_amqp_instances" "default" {
  status = "SERVING"
}


`, name)
}

// Test Amqp OpenSourcePermission. <<< Resource test cases, automatically generated.
