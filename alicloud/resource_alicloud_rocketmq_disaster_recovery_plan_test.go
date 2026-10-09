package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq DisasterRecoveryPlan. >>> Resource test cases.
func TestAccAliCloudRocketmqDisasterRecoveryPlan_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_disaster_recovery_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqDisasterRecoveryPlanMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqDisasterRecoveryPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqdrplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqDisasterRecoveryPlanBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name":               "${var.name}",
					"plan_desc":               "tf-testacc-desc",
					"plan_type":               "ACTIVE_PASSIVE",
					"auto_sync_checkpoint":    true,
					"sync_checkpoint_enabled": true,
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "${var.auth_type}",
							"network_type":      "TCP_VPC",
							"username":          "tf-testacc-dr-user-2",
							"password":          "Test654321!",
							"consumer_group_id": "${alicloud_rocketmq_consumer_group.default.consumer_group_id}",
							"message_property": []map[string]interface{}{
								{
									"property_key":   "tag",
									"property_value": "*",
								},
							},
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":                                     "${var.name}",
						"plan_desc":                                     "tf-testacc-desc",
						"plan_type":                                     "ACTIVE_PASSIVE",
						"auto_sync_checkpoint":                          "true",
						"sync_checkpoint_enabled":                       "true",
						"instances.#":                                   "2",
						"instances.0.instance_id":                       CHECKSET,
						"instances.0.instance_role":                     "ACTIVE",
						"instances.0.instance_type":                     "ALIYUN_ROCKETMQ",
						"instances.0.region_id":                         CHECKSET,
						"instances.0.endpoint_url":                      CHECKSET,
						"instances.0.vpc_id":                            CHECKSET,
						"instances.0.vswitch_id":                        CHECKSET,
						"instances.0.security_group_id":                 CHECKSET,
						"instances.0.auth_type":                         "ACL_AUTH",
						"instances.0.username":                          CHECKSET,
						"instances.0.network_type":                      "TCP_VPC",
						"instances.0.consumer_group_id":                 CHECKSET,
						"instances.0.message_property.#":                "1",
						"instances.0.message_property.0.property_key":   "tag",
						"instances.0.message_property.0.property_value": "*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name":               "${var.name}-updated",
					"plan_desc":               "tf-testacc-desc-updated",
					"plan_type":               "ACTIVE_PASSIVE",
					"auto_sync_checkpoint":    false,
					"sync_checkpoint_enabled": false,
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "ACL_AUTH",
							"network_type":      "TCP_VPC",
							"username":          "${alicloud_rocketmq_account.default.username}",
							"password":          "${alicloud_rocketmq_account.default.password}",
							"consumer_group_id": "${alicloud_rocketmq_consumer_group.default.consumer_group_id}",
							"message_property": []map[string]interface{}{
								{
									"property_key":   "tag2",
									"property_value": "important",
								},
							},
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":               "${var.name}-updated",
						"plan_desc":               "tf-testacc-desc-updated",
						"auto_sync_checkpoint":    "false",
						"sync_checkpoint_enabled": "false",
						"instances.0.auth_type":   "ACL_AUTH",
						"instances.0.username":    CHECKSET,
						"instances.0.message_property.0.property_key":   "tag2",
						"instances.0.message_property.0.property_value": "important",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instances.0.password"},
			},
		},
	})
}

var AliCloudRocketmqDisasterRecoveryPlanMap = map[string]string{
	"plan_id":     CHECKSET,
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqDisasterRecoveryPlanBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "auth_type" {
  default = "ACL_AUTH"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "example"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.createVPC.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"]
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  acl_info {
    default_vpc_auth_free = true
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_instance" "second" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = "${var.name}-second"
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"]
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  acl_info {
    default_vpc_auth_free = true
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_account" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  username    = "tf-testacc-dr-user"
  password    = "Test123456!"
}

resource "alicloud_rocketmq_consumer_group" "default" {
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = "tf-dr-cg-${var.name}"
  delivery_order_type = "Concurrently"
  remark              = "example"
  consume_retry_policy {
    retry_policy    = "DefaultRetryPolicy"
    max_retry_times = 3
  }
}
`, name)
}

func TestAccAliCloudRocketmqDisasterRecoveryPlan_instanceOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_disaster_recovery_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqDisasterRecoveryPlanMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqDisasterRecoveryPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqdrplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqDisasterRecoveryPlanOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "${var.name}",
					"plan_desc": "tf-testacc-order",
					"plan_type": "ACTIVE_PASSIVE",
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "ACL_AUTH",
							"username":          "${alicloud_rocketmq_account.default.username}",
							"password":          "${alicloud_rocketmq_account.default.password}",
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":                 "${var.name}",
						"instances.#":               "2",
						"instances.0.instance_role": "ACTIVE",
						"instances.1.instance_role": "PASSIVE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "${var.name}",
					"plan_desc": "tf-testacc-order",
					"plan_type": "ACTIVE_PASSIVE",
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "ACL_AUTH",
							"username":          "${alicloud_rocketmq_account.default.username}",
							"password":          "${alicloud_rocketmq_account.default.password}",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "${var.name}",
					"plan_desc": "tf-testacc-order",
					"plan_type": "ACTIVE_PASSIVE",
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "ACL_AUTH",
							"username":          "${alicloud_rocketmq_account.default.username}",
							"password":          "${alicloud_rocketmq_account.default.password}",
						},
					},
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func AliCloudRocketmqDisasterRecoveryPlanOrderDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "example"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.createVPC.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = "${var.name}-first"
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  series_code       = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  acl_info {
    default_vpc_auth_free = true
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_instance" "second" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = "${var.name}-second"
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  series_code       = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  acl_info {
    default_vpc_auth_free = true
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_account" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  username    = "tf-testacc-dr-user"
  password    = "Test123456!"
}
`, name)
}

func TestAccAliCloudRocketmqDisasterRecoveryPlansDataSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_disaster_recovery_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqDisasterRecoveryPlanMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqDisasterRecoveryPlan")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqdrplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqDisasterRecoveryPlanBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "data.alicloud_rocketmq_disaster_recovery_plans.default",
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name":               "${var.name}",
					"plan_desc":               "tf-testacc-ds-desc",
					"plan_type":               "ACTIVE_PASSIVE",
					"auto_sync_checkpoint":    true,
					"sync_checkpoint_enabled": true,
					"instances": []map[string]interface{}{
						{
							"instance_id":       "${alicloud_rocketmq_instance.default.id}",
							"instance_role":     "ACTIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "ACL_AUTH",
							"network_type":      "TCP_VPC",
							"username":          "${alicloud_rocketmq_account.default.username}",
							"password":          "${alicloud_rocketmq_account.default.password}",
							"consumer_group_id": "${alicloud_rocketmq_consumer_group.default.consumer_group_id}",
						},
						{
							"instance_id":       "${alicloud_rocketmq_instance.second.id}",
							"instance_role":     "PASSIVE",
							"instance_type":     "ALIYUN_ROCKETMQ",
							"region_id":         "${alicloud_rocketmq_instance.second.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.second.network_info.0.endpoints.0.endpoint_url}",
							"network_type":      "TCP_VPC",
							"vpc_id":            "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.second.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"auth_type":         "NO_AUTH",
						},
					},
				}) + `
data "alicloud_rocketmq_disaster_recovery_plans" "default" {
  ids        = [alicloud_rocketmq_disaster_recovery_plan.default.plan_id]
  name_regex = alicloud_rocketmq_disaster_recovery_plan.default.plan_name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_rocketmq_disaster_recovery_plans.default", "plans.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_rocketmq_disaster_recovery_plans.default", "plans.0.plan_id"),
				),
			},
		},
	})
}

// Test Rocketmq DisasterRecoveryPlan. <<< Resource test cases.
