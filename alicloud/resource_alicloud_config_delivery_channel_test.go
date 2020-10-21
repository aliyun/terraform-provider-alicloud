package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudConfigDeliveryChannel_basic(t *testing.T) {
	var v config.DeliveryChannel
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.bucket}",
					"delivery_channel_type":            "OSS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "OSS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change the role arn must using resource manager master account.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"delivery_channel_assume_role_arn": "acs:ram::118272523xxxxxxx:role/aliyunserviceroleforconfig",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"delivery_channel_assume_role_arn": "acs:ram::118272523xxxxxxx:role/aliyunserviceroleforconfig",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.bucket_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.bucket}",
					"delivery_channel_type":            "OSS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "OSS",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigDeliveryChannel_MNS(t *testing.T) {
	var v config.DeliveryChannel
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence_MNS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.mns}",
					"delivery_channel_type":            "MNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "MNS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_condition": deliveryChannelCondition,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_condition": strings.Replace(strings.Replace(deliveryChannelCondition, `\n`, "\n", -1), `\"`, "\"", -1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.mns_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.mns}",
					"delivery_channel_type":            "MNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "MNS",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigDeliveryChannel_SLS(t *testing.T) {
	var v config.DeliveryChannel
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence_SLS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.sls}",
					"delivery_channel_type":            "SLS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "SLS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.sls_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.sls}",
					"delivery_channel_type":            "SLS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "SLS",
					}),
				),
			},
		},
	})
}

var ConfigDeliveryChannelMap = map[string]string{}

// Because the bucket cannot be deleted after being used by the delivery channel.
// Use pre-created Oss bucket in this test.
func ConfigDeliveryChannelBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  bucket	       = format("acs:oss:cn-beijing:%%s:ci-test-bucket-for-config",local.uid)
  bucket_change	   = format("acs:oss:cn-beijing:%%s:ci-test-bucket-for-config-update",local.uid)
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, name)
}

func ConfigDeliveryChannelBasicdependence_MNS(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  mns	       	   = format("acs:oss:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.default.name)
  mns_change	   = format("acs:oss:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.change.name)
}

resource "alicloud_mns_topic" "default" {
  name = var.name
}
resource "alicloud_mns_topic" "change" {
  name = format("%%s-change",var.name)
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, name, defaultRegionToTest)
}

func ConfigDeliveryChannelBasicdependence_SLS(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  sls	       	   = format("acs:oss:%[2]s:%%s:/project/%%s/logstore/%%s",local.uid,alicloud_log_project.this.name,alicloud_log_store.this.name)
  sls_change	   = format("acs:oss:%[2]s:%%s:/project/%%s/logstore/%%s",local.uid,alicloud_log_project.this.name,alicloud_log_store.change.name)
}

resource "alicloud_log_project" "this" {
  name = var.name
}
resource "alicloud_log_store" "this" {
  name = var.name
  project = alicloud_log_project.this.name
}
resource "alicloud_log_store" "change" {
  name = format("%%s-change",var.name)
  project = alicloud_log_project.this.name
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, strings.ToLower(name), defaultRegionToTest)
}

const deliveryChannelCondition = `[\n{\n\"filterType\":\"ResourceType\",\n\"values\":[\n\"ACS::CEN::CenInstance\",\n],\n\"multiple\":true\n}\n]\n`
