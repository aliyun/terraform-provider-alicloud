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
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
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
											"endpoint_group_id": "${alicloud_ga_endpoint_group.default.id}",
										},
									},
								},
							},
						},
					},
					"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
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
											"endpoint_group_id": "${alicloud_ga_endpoint_group.default.id}",
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
  default  = "%s"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="70"
    to_port="70"
  }
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity="SOURCE_IP"
  protocol="HTTP"
  name=var.name
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name = var.name
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  endpoint_configurations{
    endpoint=alicloud_eip_address.default.ip_address
    type="PublicIp"
    weight="20"
  }
  description=var.name
  name=var.name
  threshold_count=4
  endpoint_group_region="%s"
  health_check_interval_seconds="3"
  health_check_path="/healthcheck"
  health_check_port="9999"
  health_check_protocol="http"
  port_overrides{
    endpoint_port="10"
    listener_port="70"
  }
  traffic_percentage=20
  listener_id=alicloud_ga_listener.default.id
  endpoint_group_type = "virtual"
}
`, name, defaultRegionToTest)
}
