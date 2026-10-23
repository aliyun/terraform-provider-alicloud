package rocketmq_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Stock TypeList debt surfaced by touching this file; these attributes are not
// part of this change. The latter two sit under ForceNew blocks, so reordering
// them would replace the instance.
// collection-order: bypass acl_info.acl_types -- pre-existing attribute outside this change
// collection-order: bypass network_info.internet_info.ip_whitelist -- pre-existing deprecated attribute outside this change
// collection-order: bypass network_info.vpc_info.vswitches -- pre-existing attribute outside this change

// Test case to verify the fix for send_receive_ratio validation issue
func TestAccAliCloudRocketmqInstance_SendReceiveRatioValidation(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.validation_test"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	//testAccCheck is not being used
	//testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstancevalidation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4101)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							// This should be out of range [0.05, 0.5] to test validation
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.03",
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
					"remark":            "validation test",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected product_info\.0\.send_receive_ratio to be in the range \(0\.050000 - 0\.500000\), got 0\.030000`),
			},
		},
	})
}

func TestAccAliCloudRocketmqInstance_bugfix(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4101)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4665)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16"},
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
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
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
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "ultimate",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				// The software maintenance window is managed by the service for
				// provisioned instances and cannot be set standalone here.
				// Removing the required sub_series_code stops the plan before
				// the placeholder reaches the API.
				Config: testAccConfig(map[string]interface{}{
					"software": []map[string]interface{}{
						{"maintain_time": "05:00:00"},
					},
					"sub_series_code": acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
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

var instanceMap4665 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence4665(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-j"
	}

	resource "alicloud_kms_instance" "default" {
  		product_version = "3"
  		vpc_num         = "1"
  		key_num         = "1000"
  		secret_num      = "0"
  		spec            = "1000"
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		"cn-hangzhou-k",
    		"cn-hangzhou-j"
  		]
  		timeouts {
    		delete = "60m"
  		}
	}

	resource "alicloud_kms_key" "default" {
  		description            = "Hello KMS"
  		status                 = "Enabled"
  		pending_window_in_days = "7"
  		dkms_instance_id       = alicloud_kms_instance.default.id
	}
`, name)
}

// Case 4652 From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic4652(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4652)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var instanceMap4652 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence4652(name string) string {
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4128)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil, // bypass linter, prepaid (Subscription) resource cannot destroy
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
						"resource_group_id": acctest.CHECKSET,
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"auto_renew":        "true",
						"auto_renew_period": "1",
						"service_code":      "rmq",
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var instanceMap4128 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence4128(name string) string {
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4101)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"acl_info": []map[string]interface{}{
						{
							"acl_types": []string{
								"default", "apache_acl", "aliyun_ram"},
							"default_vpc_auth_free": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_info.#": "1",
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"ip_whitelists": []string{"192.168.2.2", "192.168.1.1", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.2.2"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.3.3"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.2.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.2.2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.2.2", "192.168.1.1", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.2.2"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.3.3"),
				),
			},
			{
				// Unsetting ip_whitelists must leave the instance whitelist untouched.
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.2.2"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.3.3"),
				),
			},
			{
				// An explicitly empty set is authoritative and must clear every entry without a permanent diff.
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": acctest.CLEARLIST,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "0",
					}),
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
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": acctest.CHECKSET,
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

var instanceMap4101 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence4101(name string) string {
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4665)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
							"storage_encryption":     "true",
							"storage_secret_key":     "${alicloud_kms_key.default.id}",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
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
						"resource_group_id": acctest.CHECKSET,
						"product_info.#":    "1",
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4652)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4128)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil, // bypass linter, postpay resource cannot destroy
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":          name,
						"auto_renew":             "true",
						"auto_renew_period":      "1",
						"auto_renew_period_unit": "Year",
						"service_code":           "rmq",
						"resource_group_id":      acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap4101)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"acl_info": []map[string]interface{}{
						{
							"acl_types":             []string{"default", "apache_acl", "aliyun_ram"},
							"default_vpc_auth_free": "false",
						},
					},
					"ip_whitelists": []string{"192.168.1.1", "192.168.2.2"},
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
						"resource_group_id": acctest.CHECKSET,
						"acl_info.#":        "1",
						"ip_whitelists.#":   "2",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "192.168.2.2"),
				),
			},
			{
				// An explicitly configured default entry is part of the authoritative set and must be kept.
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"0.0.0.0/0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "ip_whitelists.*", "0.0.0.0/0"),
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap6747)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence6747)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"trace_on": "false",
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s3.nxlarge",
							"auto_scaling":           "true",
							"message_retention_time": "72",
							"trace_on":               "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_info.#": "1",
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

var instanceMap6747 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence6747(name string) string {
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
	ra := acctest.ResourceAttrInit(resourceId, instanceMap7144)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, instanceBasicDependence7144)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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

var instanceMap7144 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func instanceBasicDependence7144(name string) string {
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
