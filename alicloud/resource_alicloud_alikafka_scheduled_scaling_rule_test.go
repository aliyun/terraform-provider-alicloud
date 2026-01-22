package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Alikafka ScheduledScalingRule. >>> Resource test cases, automatically generated.
// Case 定时策略全生命周期-网段-195.0.0.0/25 12412
func TestAccAliCloudAlikafkaScheduledScalingRule_basic12412(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_scheduled_scaling_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaScheduledScalingRuleMap12412)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaScheduledScalingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	firstScheduledTime := time.Now().AddDate(0, 0, 6).UnixMilli()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaScheduledScalingRuleBasicDependence12412)
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
					"duration_minutes":     "100",
					"first_scheduled_time": firstScheduledTime,
					"instance_id":          "${alicloud_alikafka_instance.default.id}",
					"reserved_pub_flow":    "200",
					"reserved_sub_flow":    "200",
					"rule_name":            name,
					"schedule_type":        "at",
					"time_zone":            "GMT+8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration_minutes":     "100",
						"first_scheduled_time": CHECKSET,
						"instance_id":          CHECKSET,
						"reserved_pub_flow":    "200",
						"reserved_sub_flow":    "200",
						"rule_name":            name,
						"schedule_type":        "at",
						"time_zone":            "GMT+8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
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

func TestAccAliCloudAlikafkaScheduledScalingRule_basic12412_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_scheduled_scaling_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaScheduledScalingRuleMap12412)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaScheduledScalingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	firstScheduledTime := time.Now().AddDate(0, 0, 6).UnixMilli()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaScheduledScalingRuleBasicDependence12412)
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
					"duration_minutes":     "100",
					"first_scheduled_time": firstScheduledTime,
					"instance_id":          "${alicloud_alikafka_instance.default.id}",
					"reserved_pub_flow":    "200",
					"reserved_sub_flow":    "200",
					"rule_name":            name,
					"schedule_type":        "repeat",
					"repeat_type":          "Weekly",
					"time_zone":            "GMT+8",
					"enable":               "false",
					"weekly_types":         []string{"Tuesday", "Monday", "Friday", "Wednesday", "Thursday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration_minutes":     "100",
						"first_scheduled_time": CHECKSET,
						"instance_id":          CHECKSET,
						"reserved_pub_flow":    "200",
						"reserved_sub_flow":    "200",
						"rule_name":            name,
						"schedule_type":        "repeat",
						"repeat_type":          "Weekly",
						"time_zone":            "GMT+8",
						"enable":               "false",
						"weekly_types.#":       "5",
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

var AliCloudAlikafkaScheduledScalingRuleMap12412 = map[string]string{}

func AliCloudAlikafkaScheduledScalingRuleBasicDependence12412(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "10.4.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_alikafka_instance" "default" {
  		deploy_type     = "4"
  		instance_type   = "alikafka_serverless"
  		vswitch_id      = alicloud_vswitch.default.id
  		spec_type       = "normal"
  		service_version = "3.3.1"
  		security_group  = alicloud_security_group.default.id
  		config          = "{\"enable.acl\":\"true\"}"
  		serverless_config {
    		reserved_publish_capacity   = 60
    		reserved_subscribe_capacity = 60
  		}
	}
`, name)
}

// Test Alikafka ScheduledScalingRule. <<< Resource test cases, automatically generated.
