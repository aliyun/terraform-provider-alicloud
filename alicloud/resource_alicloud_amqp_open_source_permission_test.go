// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Amqp OpenSourcePermission. >>> Resource test cases, automatically generated.
// Case 开源权限_副本1779693632554 12809
func TestAccAliCloudAmqpOpenSourcePermission_basic12809(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_open_source_permission.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpOpenSourcePermissionMap12809)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpOpenSourcePermission")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpOpenSourcePermissionBasicDependence12809)
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
					"write":       ".*",
					"read":        ".*",
					"vhost":       "${var.vhost}",
					"user_name":   "${var.user_name}",
					"instance_id": "${alicloud_amqp_instance.CreateInstance.id}",
					"configure":   ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write":       ".*",
						"read":        ".*",
						"vhost":       CHECKSET,
						"user_name":   CHECKSET,
						"instance_id": CHECKSET,
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

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = var.instance_name
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  vpc_id                = alicloud_vpc.default.id
  vswitch_ids           = [alicloud_vswitch.default_b.id, alicloud_vswitch.default_g.id]
  security_group_id     = alicloud_security_group.default.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default_b" {
  vswitch_name = "${var.name}-b"
  cidr_block   = "172.16.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-b"
}

resource "alicloud_vswitch" "default_g" {
  vswitch_name = "${var.name}-g"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-g"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}


`, name)
}

// Case 开源权限 12794
func TestAccAliCloudAmqpOpenSourcePermission_basic12794(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_open_source_permission.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpOpenSourcePermissionMap12794)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpOpenSourcePermission")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpOpenSourcePermissionBasicDependence12794)
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
					"write":       ".*",
					"read":        ".*",
					"vhost":       "${var.vhost}",
					"user_name":   "${var.user_name}",
					"instance_id": "${alicloud_amqp_instance.CreateInstance.id}",
					"configure":   ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write":       ".*",
						"read":        ".*",
						"vhost":       CHECKSET,
						"user_name":   CHECKSET,
						"instance_id": CHECKSET,
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

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = var.instance_name
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  vpc_id                = alicloud_vpc.default.id
  vswitch_ids           = [alicloud_vswitch.default_b.id, alicloud_vswitch.default_g.id]
  security_group_id     = alicloud_security_group.default.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default_b" {
  vswitch_name = "${var.name}-b"
  cidr_block   = "172.16.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-b"
}

resource "alicloud_vswitch" "default_g" {
  vswitch_name = "${var.name}-g"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-g"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}


`, name)
}

// Test Amqp OpenSourcePermission. <<< Resource test cases, automatically generated.
