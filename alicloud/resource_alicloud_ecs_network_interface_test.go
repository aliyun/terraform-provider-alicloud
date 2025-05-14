package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_network_interface", &resource.Sweeper{
		Name: "alicloud_ecs_network_interface",
		F:    testAliCloudEcsNetworkInterface,
	})
}

func testAliCloudEcsNetworkInterface(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"fc-eni", // Clean up the eni which created by fc service
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "DescribeNetworkInterfaces"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Describe NetworkInterface failed, error: %s. Return!", err.Error())
			return nil
		}

		resp, err := jsonpath.Get("$.NetworkInterfaceSets.NetworkInterfaceSet", response)
		if err != nil {
			log.Printf("[ERROR] jsonpath Get NetworkInterface failed, %#v", err)
			continue
		}

		result, _ := resp.([]interface{})
		service := VpcService{client}
		ecsService := EcsService{client}
		for _, v := range result {
			item := v.(map[string]interface{})
			if _, ok := item["NetworkInterfaceName"]; !ok {
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["NetworkInterfaceName"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				// If a nat gateway name is not set successfully, it should be fetched by vpc name and deleted.
				if skip {
					if need, err := service.needSweepVpc(item["VpcId"].(string), ""); err == nil {
						skip = !need
					}
				}
				if skip {
					log.Printf("[INFO] Skipping NetworkInterface: %s", item["NetworkInterfaceName"].(string))
					continue
				}
			}
			if item["InstanceId"] != "" {
				requestDetach := map[string]interface{}{
					"InstanceId":         item["InstanceId"],
					"NetworkInterfaceId": item["NetworkInterfaceId"],
					"RegionId":           client.RegionId,
				}
				actionDetach := "DetachNetworkInterface"

				response, err = client.RpcPost("Ecs", "2014-05-26", actionDetach, nil, requestDetach, true)
				if err != nil {
					log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
					continue
				}
				stateConf := BuildStateConf([]string{}, []string{"Available"}, DefaultTimeout, 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(item["NetworkInterfaceId"].(string), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
					continue
				}
			}
			action = "DeleteNetworkInterface"
			request := map[string]interface{}{
				"NetworkInterfaceId": item["NetworkInterfaceId"],
				"RegionId":           client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete NetworkInterface (%s): %s", item["NetworkInterfaceName"].(string), err)
				continue
			}
			log.Printf("[INFO] Delete NetworkInterface success: %s ", item["NetworkInterfaceName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

func TestAccAliCloudECSNetworkInterface_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 255)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"network_interface_name": name,
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_group_ids.#":   "1",
						"resource_group_id":      CHECKSET,
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
					"private_ip_addresses": []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+2)},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_addresses.#": "3",
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
					"queue_number": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_number": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Test For Terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Test For Terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interface_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip_addresses": []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+2)},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_addresses.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip_addresses":   []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+2)},
					"network_interface_name": name,
					"description":            "Test For Terraform Update",
					"queue_number":           "2",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_addresses.#": "3",
						"network_interface_name": name,
						"description":            "Test For Terraform Update",
						"queue_number":           "2",
						"tags.%":                 "2",
						"tags.Created":           "TF-update",
						"tags.For":               "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}", "${alicloud_security_group.default.2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "3",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudECSNetworkInterface_primary_ip_address(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceIpv6Dependence)
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
					"network_interface_name": name,
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"primary_ip_address":     "${cidrhost(alicloud_vswitch.default.cidr_block, 26)}",
					"ipv6_addresses":         []string{"${cidrhost(alicloud_vswitch.default.ipv6_cidr_block, 64)}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_group_ids.#":   "1",
						"resource_group_id":      CHECKSET,
						"primary_ip_address":     CHECKSET,
						"ipv6_addresses.#":       "1",
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
					"ipv6_addresses": []string{
						"${cidrhost(alicloud_vswitch.default.ipv6_cidr_block, 64)}",
						"${cidrhost(alicloud_vswitch.default.ipv6_cidr_block, 32)}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_addresses.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_addresses": []string{
						"${cidrhost(alicloud_vswitch.default.ipv6_cidr_block, 32)}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_addresses.#": "1",
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

func TestAccAliCloudECSNetworkInterface_secondary_private_ip_address_count(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceIpv6Dependence)
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
					"network_interface_name": name,
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"ipv6_address_count":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_group_ids.#":   "1",
						"resource_group_id":      CHECKSET,
						"ipv6_address_count":     "1",
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
					"secondary_private_ip_address_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_private_ip_address_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_private_ip_address_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_private_ip_address_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_count": "1",
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

func TestAccAliCloudECSNetworkInterface_ipv4_prefix_address(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	checkoutSupportedRegions(t, true, connectivity.EcsActivationsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceIpv4PrefixDependence)
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
					"network_interface_name": name,
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"primary_ip_address":     "${cidrhost(alicloud_vswitch.default.cidr_block, 26)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_group_ids.#":   "1",
						"resource_group_id":      CHECKSET,
						"primary_ip_address":     CHECKSET,
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
					"ipv4_prefixes": []string{
						"172.16.10.16/28",
						"172.16.10.32/28",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv4_prefixes.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv4_prefixes": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv4_prefixes.#": "0",
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

func TestAccAliCloudECSNetworkInterface_ipv4_prefix_count(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	checkoutSupportedRegions(t, true, connectivity.EcsActivationsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceIpv4PrefixDependence)
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
					"network_interface_name":         name,
					"vswitch_id":                     "${alicloud_vswitch.default.id}",
					"security_group_ids":             []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"network_interface_traffic_mode": "HighPerformance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name":         CHECKSET,
						"vswitch_id":                     CHECKSET,
						"security_group_ids.#":           "1",
						"resource_group_id":              CHECKSET,
						"network_interface_traffic_mode": "HighPerformance",
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
					"ipv4_prefix_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv4_prefix_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv4_prefix_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv4_prefix_count": "1",
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

var AlicloudEcsNetworkInterfaceMap = map[string]string{
	"mac":                    CHECKSET,
	"network_interface_name": CHECKSET,
	"primary_ip_address":     CHECKSET,
	"queue_number":           CHECKSET,
	"resource_group_id":      CHECKSET,
	"status":                 CHECKSET,
	"vswitch_id":             CHECKSET,
}

func TestAccAliCloudECSNetworkInterface_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":                               name,
					"vswitch_id":                         "${alicloud_vswitch.default.id}",
					"security_groups":                    []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}", "${alicloud_security_group.default.2.id}"},
					"resource_group_id":                  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":                        name,
					"private_ip":                         fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand),
					"queue_number":                       "1",
					"secondary_private_ip_address_count": "1",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                               CHECKSET,
						"vswitch_id":                         CHECKSET,
						"security_groups.#":                  "3",
						"resource_group_id":                  CHECKSET,
						"description":                        name,
						"private_ip":                         CHECKSET,
						"queue_number":                       "1",
						"secondary_private_ip_address_count": "1",
						"tags.%":                             "2",
						"tags.Created":                       "TF-update",
						"tags.For":                           "Test-update",
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

func TestAccAliCloudECSNetworkInterface_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":              name,
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_groups":   []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":       name,
					"private_ip":        fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand),
					"queue_number":      "1",
					"private_ips_count": "1",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_groups.#": "1",
						"resource_group_id": CHECKSET,
						"description":       name,
						"private_ip":        CHECKSET,
						"queue_number":      "1",
						"private_ips_count": "1",
						"tags.%":            "2",
						"tags.Created":      "TF-update",
						"tags.For":          "Test-update",
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

func TestAccAliCloudECSNetworkInterface_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":                 name,
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"security_groups":      []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":          name,
					"instance_type":        "Trunk",
					"private_ip_addresses": []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1)},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_groups.#":      "1",
						"resource_group_id":      CHECKSET,
						"description":            name,
						"instance_type":          "Trunk",
						"private_ip_addresses.#": "2",
						"tags.%":                 "2",
						"tags.Created":           "TF-update",
						"tags.For":               "Test-update",
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

func TestAccAliCloudECSNetworkInterface_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":                 name,
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"security_groups":      []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":          name,
					"private_ip_addresses": []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1)},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_groups.#": "1",
						"resource_group_id": CHECKSET,
						"description":       name,
						"private_ips.#":     "2",
						"tags.%":            "2",
						"tags.Created":      "TF-update",
						"tags.For":          "Test-update",
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

func TestAccAliCloudECSNetworkInterface_name_deprecated(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":               name,
					"vswitch_id":         "${alicloud_vswitch.default.id}",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 CHECKSET,
						"vswitch_id":           CHECKSET,
						"security_group_ids.#": "1",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "Update",
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

func TestAccAliCloudECSNetworkInterface_private_ips_deprecated(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":               name,
					"vswitch_id":         "${alicloud_vswitch.default.id}",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 CHECKSET,
						"vswitch_id":           CHECKSET,
						"security_group_ids.#": "1",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ips": []string{fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand), fmt.Sprintf("${cidrhost(alicloud_vswitch.default.cidr_block, %d)}", rand+1)},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips.#": "2",
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

func TestAccAliCloudECSNetworkInterface_private_ips_count_deprecated(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":               name,
					"vswitch_id":         "${alicloud_vswitch.default.id}",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}", "${alicloud_security_group.default.2.id}"},
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 CHECKSET,
						"vswitch_id":           CHECKSET,
						"security_group_ids.#": "3",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ips_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips_count": "2",
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

func TestAccAliCloudECSNetworkInterface_security_groups_deprecated(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsNetworkInterfaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(2, 253)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsNetworkInterfaceBasicDependence)
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
					"name":              name,
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_groups":   []string{"${alicloud_security_group.default.0.id}"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_groups.#": "1",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}", "${alicloud_security_group.default.2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "3",
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

func AliCloudEcsNetworkInterfaceBasicDependence(name string) string {
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

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}
	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		count  = 3
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

func AliCloudEcsNetworkInterfaceIpv6Dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = "true"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id                = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name           = "${var.name}"
  ipv6_cidr_block_mask   =  64
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = alicloud_vpc.default.id
}
data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}

func AliCloudEcsNetworkInterfaceIpv4PrefixDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  zone_id                = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name           = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_vpc_vswitch_cidr_reservation" "default" {
	vswitch_id = "${alicloud_vswitch.default.id}"
	cidr_reservation_cidr = "172.16.10.0/24"
}

resource "alicloud_ecs_network_interface" "example" {
    name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
	depends_on = ["alicloud_vpc_vswitch_cidr_reservation.default"]
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}
