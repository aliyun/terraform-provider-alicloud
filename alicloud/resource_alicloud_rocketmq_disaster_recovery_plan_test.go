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
							"instance_role":     "SOURCE",
							"instance_type":     "rmq",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.security_group_ids}",
							"auth_type":         "USER",
							"network_type":      "VPC",
							"username":          "tf-testacc-user",
							"password":          "Test123456!",
							"consumer_group_id": "tf-testacc-cg",
							"message_property": []map[string]interface{}{
								{
									"property_key":   "tag",
									"property_value": "*",
								},
							},
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
						"instances.#":                                   "1",
						"instances.0.instance_id":                       CHECKSET,
						"instances.0.instance_role":                     "SOURCE",
						"instances.0.instance_type":                     "rmq",
						"instances.0.region_id":                         CHECKSET,
						"instances.0.endpoint_url":                      CHECKSET,
						"instances.0.vpc_id":                            CHECKSET,
						"instances.0.vswitch_id":                        CHECKSET,
						"instances.0.security_group_id":                 CHECKSET,
						"instances.0.auth_type":                         "USER",
						"instances.0.network_type":                      "VPC",
						"instances.0.username":                          "tf-testacc-user",
						"instances.0.consumer_group_id":                 "tf-testacc-cg",
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
							"instance_role":     "SOURCE",
							"instance_type":     "rmq",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.security_group_ids}",
							"auth_type":         "USER",
							"network_type":      "VPC",
							"username":          "tf-testacc-user-updated",
							"password":          "Test123456!",
							"consumer_group_id": "tf-testacc-cg-updated",
							"message_property": []map[string]interface{}{
								{
									"property_key":   "tag",
									"property_value": "important",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":                                     "${var.name}-updated",
						"plan_desc":                                     "tf-testacc-desc-updated",
						"auto_sync_checkpoint":                          "false",
						"sync_checkpoint_enabled":                       "false",
						"instances.0.username":                          "tf-testacc-user-updated",
						"instances.0.consumer_group_id":                 "tf-testacc-cg-updated",
						"instances.0.message_property.0.property_value": "important",
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
							"instance_role":     "SOURCE",
							"instance_type":     "rmq",
							"region_id":         "${alicloud_rocketmq_instance.default.region_id}",
							"endpoint_url":      "${alicloud_rocketmq_instance.default.network_info.0.endpoints.0.endpoint_url}",
							"vpc_id":            "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vpc_id}",
							"vswitch_id":        "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.vswitches.0.vswitch_id}",
							"security_group_id": "${alicloud_rocketmq_instance.default.network_info.0.vpc_info.0.security_group_ids}",
							"auth_type":         "USER",
							"network_type":      "VPC",
							"username":          "tf-testacc-ds-user",
							"password":          "Test123456!",
							"consumer_group_id": "tf-testacc-ds-cg",
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
