package ga_test

// collection-order: bypass rule_conditions.host_config.values -- pre-existing multi-member TypeList attribute left untouched by this migration; convergent reorder coverage is pre-existing debt outside this change
// collection-order: bypass rule_conditions.path_config.values -- pre-existing multi-member TypeList attribute left untouched by this migration; convergent reorder coverage is pre-existing debt outside this change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaForwardingRule_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaForwardingRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaForwardingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaForwardingRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaForwardingRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Path",
							"path_config": []map[string]interface{}{
								{
									"values": []string{"/tf-testAcc"},
								},
							},
						},
					},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":    acctest.CHECKSET,
						"listener_id":       acctest.CHECKSET,
						"rule_conditions.#": "1",
						"rule_actions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2000",
					}),
				),
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
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Path",
							"path_config": []map[string]interface{}{
								{
									"values": []string{"/tf-testAcc-update"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Host",
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.alicloud.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"order":            "25",
							"rule_action_type": "Drop",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"order":            "30",
							"rule_action_type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"endpoint_group_id": "${alicloud_ga_endpoint_group.update.id}",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaForwardingRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaForwardingRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaForwardingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaForwardingRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaForwardingRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":          "${alicloud_ga_listener.default.id}",
					"priority":             "2000",
					"forwarding_rule_name": name,
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Path",
							"path_config": []map[string]interface{}{
								{
									"values": []string{"/tf-testAcc"},
								},
							},
						},
					},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"listener_id":          acctest.CHECKSET,
						"priority":             "2000",
						"forwarding_rule_name": name,
						"rule_conditions.#":    "1",
						"rule_actions.#":       "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaForwardingRule_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaForwardingRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaForwardingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaForwardingRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaForwardingRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type":  "Path",
							"rule_condition_value": `[\"/tf-testAcc\"]`,
						},
					},
					"rule_actions": []map[string]interface{}{
						{
							"order":             "20",
							"rule_action_type":  "ForwardGroup",
							"rule_action_value": `[{\"type\":\"endpointgroup\", \"value\":\"` + "${alicloud_ga_endpoint_group.default.id}" + `\"}]`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":    acctest.CHECKSET,
						"listener_id":       acctest.CHECKSET,
						"rule_conditions.#": "1",
						"rule_actions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2000",
					}),
				),
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
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type":  "Path",
							"rule_condition_value": `[\"/tf-testAcc-update\"]`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type":  "Host",
							"rule_condition_value": `[\"www.alicloud.com\"]`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"order":            "25",
							"rule_action_type": "Drop",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"order":             "30",
							"rule_action_type":  "ForwardGroup",
							"rule_action_value": `[{\"type\":\"endpointgroup\", \"value\":\"` + "${alicloud_ga_endpoint_group.update.id}" + `\"}]`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudGaForwardingRule_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaForwardingRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaForwardingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaForwardingRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaForwardingRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":          "${alicloud_ga_listener.default.id}",
					"priority":             "2000",
					"forwarding_rule_name": name,
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type":  "Host",
							"rule_condition_value": `[\"www.alicloud.com\"]`,
						},
					},
					"rule_actions": []map[string]interface{}{
						{
							"order":             "20",
							"rule_action_type":  "ForwardGroup",
							"rule_action_value": `[{\"type\":\"endpointgroup\", \"value\":\"` + "${alicloud_ga_endpoint_group.default.id}" + `\"}]`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"listener_id":          acctest.CHECKSET,
						"priority":             "2000",
						"forwarding_rule_name": name,
						"rule_conditions.#":    "1",
						"rule_actions.#":       "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudGaForwardingRuleMap0 = map[string]string{
	"priority":               acctest.CHECKSET,
	"forwarding_rule_id":     acctest.CHECKSET,
	"forwarding_rule_status": acctest.CHECKSET,
}

func AliCloudGaForwardingRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
        bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Basic"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
  		auto_pay               = true
  		auto_use_coupon        = true
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		port_ranges {
    		from_port = "70"
    		to_port   = "70"
  		}
  		accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		client_affinity = "SOURCE_IP"
  		protocol        = "HTTP"
  		name            = var.name
	}

	resource "alicloud_eip_address" "default" {
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}

	resource "alicloud_ga_endpoint_group" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		endpoint_configurations {
    		endpoint = alicloud_eip_address.default.ip_address
    		type     = "PublicIp"
    		weight   = "20"
  		}
  		description                   = var.name
  		name                          = var.name
  		threshold_count               = 4
  		endpoint_group_region         = "%s"
  		health_check_interval_seconds = "3"
  		health_check_path             = "/healthcheck"
  		health_check_port             = "9999"
  		health_check_protocol         = "http"
  		port_overrides {
    		endpoint_port = "10"
    		listener_port = "70"
  		}
  		traffic_percentage        = 20
  		listener_id               = alicloud_ga_listener.default.id
  		endpoint_group_type       = "virtual"
  		endpoint_request_protocol = "HTTP"
	}

	resource "alicloud_ga_endpoint_group" "update" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		endpoint_configurations {
    		endpoint = alicloud_eip_address.default.ip_address
    		type     = "PublicIp"
    		weight   = "20"
  		}
  		description                   = "${var.name}-update"
  		name                          = "${var.name}-update"
  		threshold_count               = 4
  		endpoint_group_region         = "%s"
  		health_check_interval_seconds = "3"
  		health_check_path             = "/healthcheck"
  		health_check_port             = "9999"
  		health_check_protocol         = "http"
  		port_overrides {
    		endpoint_port = "10"
    		listener_port = "70"
  		}
  		traffic_percentage        = 20
  		listener_id               = alicloud_ga_listener.default.id
  		endpoint_group_type       = "virtual"
  		endpoint_request_protocol = "HTTP"
	}
`, name, acctest.DefaultRegionToTest, acctest.DefaultRegionToTest)
}
