package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testSweepRocketMq(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	action := "/instances"
	request := make(map[string]interface{})
	query := make(map[string]*string)
	query["pageNumber"] = tea.String("1")
	query["pageSize"] = tea.String("200")
	var response map[string]interface{}
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.body.data.list", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.body.data.list", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["instanceName"].(string)
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping rocketmq instance: %s ", name)
				continue
			}
		}
		log.Printf("[INFO] delete rocketmq instance: %s ", name)

		conn, err := client.NewRocketmqClient()
		if err != nil {
			return WrapError(err)
		}
		action = fmt.Sprintf("/instances/%s", item["instanceId"])
		query := make(map[string]*string)
		body := make(map[string]interface{})
		if err != nil {
			return WrapError(err)
		}
		request = make(map[string]interface{})
		request["instanceId"] = item["instanceId"]

		body = request
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to delete rocketmq instance (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudRocketmqInstance_bugfix(t *testing.T) {
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
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
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

// Case 4652 From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic4652(t *testing.T) {
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

// Case 4652  twin From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic4652_twin(t *testing.T) {
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

// Test Rocketmq Instance. >>> Resource test cases, automatically generated.
// Case 创建serverless实例 6747
func TestAccAliCloudRocketmqInstance_basic6747(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap6747)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence6747)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"series_code":  "standard",
					"payment_type": "PayAsYouGo",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s3.nxlarge",
							"auto_scaling":           "true",
							"message_retention_time": "72",
							// "support_auto_scaling":   "true",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":             "${alicloud_vpc.createVPC.id}",
									"security_group_ids": "${alicloud_security_group.CreateSecurityGroup.id}",
									"vswitches": []map[string]interface{}{
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
										},
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch2.id}",
										},
									},
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
								},
							},
						},
					},
					"sub_series_code": "serverless",
					"instance_name":   name,
					"service_code":    "rmq",
					"commodity_code":  "ons_rmqsrvlesspost_public_cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"resource_group_id": CHECKSET,
						"sub_series_code":   "serverless",
						"instance_name":     name,
						"service_code":      "rmq",
						"commodity_code":    "ons_rmqsrvlesspost_public_cn",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AlicloudRocketmqInstanceMap6747 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence6747(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "createVPC" {
  description = "111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_security_group" "CreateSecurityGroup" {
  name   = var.name
  vpc_id = alicloud_vpc.createVPC.id
}

resource "alicloud_vswitch" "createVSwitch" {
  vpc_id       = alicloud_vpc.createVPC.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.17.0.0/16"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "createVSwitch2" {
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.18.0.0/16"
  vswitch_name = format("%%s4", var.name)
  zone_id      = "cn-hangzhou-j"
}


`, name)
}

// Case 创建单节点实例用例_副本1719467789892 7144 From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic7144(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqInstanceMap7144)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqInstanceBasicDependence7144)
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
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"auto_scaling":           "false",
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
					"commodity_code":    "ons_rmqpost_public_cn",
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
						"commodity_code":    "ons_rmqpost_public_cn",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AlicloudRocketmqInstanceMap7144 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudRocketmqInstanceBasicDependence7144(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
  zone_id      = data.alicloud_zones.default.zones.0.id
}

`, name)
}

// Test Rocketmq Instance. <<< Resource test cases, automatically generated.
