package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq Instance. >>> Resource test cases, automatically generated.
// Case 4665
func TestAccAliCloudRocketmqInstance_basic4665(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4665)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "cluster_ha",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "ultimate",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "cluster_ha",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "ultimate",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudRocketmqInstanceMap4665 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence4665(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 4652
func TestAccAliCloudRocketmqInstance_basic4652(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4652)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "single_node",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "standard",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "single_node",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "single_node",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudRocketmqInstanceMap4652 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence4652(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 4128
func TestAccAliCloudRocketmqInstance_basic4128(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4128)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":           "Subscription",
					"auto_renew":             "true",
					"auto_renew_period":      "1",
					"auto_renew_period_unit": "Month",
					"period":                 "1",
					"period_unit":            "Month",
					"sub_series_code":        "cluster_ha",
					"instance_name":          name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":           "rmq",
						"series_code":            "standard",
						"payment_type":           "Subscription",
						"sub_series_code":        "cluster_ha",
						"instance_name":          name,
						"auto_renew":             "true",
						"auto_renew_period":      "1",
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "_update",
					"auto_renew":        "true",
					"auto_renew_period": "1",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code":      "rmq",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"series_code":       "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "Subscription",
					"period":          "1",
					"sub_series_code": "cluster_ha",
					"remark":          "自动化测试购买使用11",
					"period_unit":     "Month",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"auto_renew":        "true",
						"auto_renew_period": "1",
						"service_code":      "rmq",
						"resource_group_id": CHECKSET,
						"series_code":       "standard",
						"payment_type":      "Subscription",
						"period":            "1",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"period_unit":       "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AlicloudRocketmqInstanceMap4128 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence4128(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 4101
func TestAccAliCloudRocketmqInstance_basic4101(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_code": "rmq",
					"series_code":  "professional",
					"product_info": []map[string]interface{}{
						{
							"auto_scaling":           "false",
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "cluster_ha",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "professional",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "cluster_ha",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"10.10.0.0/16", "172.168.0.0/16", "192.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "80",
							"auto_scaling":           "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
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

var AlicloudRocketmqInstanceMap4101 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence4101(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

// Case 4665  twin
func TestAccAliCloudRocketmqInstance_basic4665_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4665)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "ultimate",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Case 4652  twin
func TestAccAliCloudRocketmqInstance_basic4652_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4652)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "single_node",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
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

// Case 4128  twin
func TestAccAliCloudRocketmqInstance_basic4128_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4128)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":          name,
					"auto_renew":             "true",
					"auto_renew_period":      "1",
					"auto_renew_period_unit": "Year",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code":      "rmq",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"series_code":       "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "Subscription",
					"period":          "1",
					"sub_series_code": "cluster_ha",
					"remark":          "自动化测试购买使用11",
					"period_unit":     "Year",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":          name,
						"auto_renew":             "true",
						"auto_renew_period":      "1",
						"auto_renew_period_unit": "Year",
						"service_code":           "rmq",
						"resource_group_id":      CHECKSET,
						"series_code":            "standard",
						"payment_type":           "Subscription",
						"period":                 "1",
						"sub_series_code":        "cluster_ha",
						"remark":                 "自动化测试购买使用11",
						"period_unit":            "Year",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

// Case 4101  twin
func TestAccAliCloudRocketmqInstance_basic4101_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap4101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"software": []map[string]interface{}{
						{
							"maintain_time": "02:00-06:00",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Test Rocketmq Instance. <<< Resource test cases, automatically generated.
