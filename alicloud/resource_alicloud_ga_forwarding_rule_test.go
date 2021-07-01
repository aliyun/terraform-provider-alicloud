package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaForwardingRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaForwardingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaForwardingRuleBasicDependence)
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
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Host",
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.test.com"},
								},
							},
						},
					},
					"priority": "1000",
					"rule_actions": []map[string]interface{}{
						{
							"order":            "20",
							"rule_action_type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"endpoint_group_id": "${alicloud_ga_endpoint_group.example.id}",
										},
									},
								},
							},
						},
					},
					"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
					"listener_id":    "${alicloud_ga_listener.example.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
						"priority":          "1000",
						"rule_actions.#":    "1",
						"accelerator_id":    CHECKSET,
						"listener_id":       CHECKSET,
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
					"forwarding_rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarding_rule_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarding_rule_name": name + "update",
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Host",
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.test3.com"},
								},
							},
						},
					},
					"priority": "2000",
					"rule_actions": []map[string]interface{}{
						{
							"order":            "30",
							"rule_action_type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"endpoint_group_id": "${alicloud_ga_endpoint_group.example.id}",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarding_rule_name": name + "update",
						"rule_conditions.#":    "1",
						"priority":             "2000",
						"rule_actions.#":       "1",
					}),
				),
			},
		},
	})
}

func AlicloudGaForwardingRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_ga_accelerators" "default"{
	}
	resource "alicloud_ga_listener" "example" {
	 accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	 port_ranges {
	   from_port = 70
	   to_port   = 70
	 }
	 protocol="HTTP"
	}
	resource "alicloud_eip_address" "example" {
	 bandwidth            = "10"
	 internet_charge_type = "PayByBandwidth"
	}
	resource "alicloud_ga_endpoint_group" "example" {
	 accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	 endpoint_configurations {
	   endpoint = alicloud_eip_address.example.ip_address
	   type     = "PublicIp"
	   weight   = "20"
	 }
	 endpoint_group_region = "cn-hangzhou"
	 listener_id           = alicloud_ga_listener.example.id
	}`, name)
}
