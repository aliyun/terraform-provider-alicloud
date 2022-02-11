package alicloud

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_network_interface", &resource.Sweeper{
		Name: "alicloud_ecs_network_interface",
		F:    testAlicloudEcsNetworkInterface,
	})
}

func testAlicloudEcsNetworkInterface(region string) error {
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
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	sweeped := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
			sweeped = true
			if item["InstanceId"] != "" {
				requestDetach := map[string]interface{}{
					"InstanceId":         item["InstanceId"],
					"NetworkInterfaceId": item["NetworkInterfaceId"],
					"RegionId":           client.RegionId,
				}
				actionDetach := "DetachNetworkInterface"

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				response, err = conn.DoRequest(StringPointer(actionDetach), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, requestDetach, &runtime)
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
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudEcsNetworkInterface_basic(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":             "${local.vswitch_id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
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
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids":     []string{"${alicloud_security_group.default.id}","security"},
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
			//{
			//	ResourceName:      resourceId,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"private_ip_addresses": []string{"${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 7)}"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"private_ip_addresses.#": "3",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"tags": map[string]string{
			//			"Created": "TF",
			//			"For":     "Test",
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"tags.%":       "2",
			//			"tags.Created": "TF",
			//			"tags.For":     "Test",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"queue_number": "1",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"queue_number": "1",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"description": "Test For Terraform",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"description": "Test For Terraform",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"network_interface_name": name + "Update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"network_interface_name": name + "Update",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"private_ip_addresses": []string{"${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 7)}"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"private_ip_addresses.#": "2",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"private_ip_addresses":   []string{"${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 5)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 6)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 8)}"},
			//		"network_interface_name": name,
			//		"description":            "Test For Terraform Update",
			//		"queue_number":           "2",
			//		"tags": map[string]string{
			//			"Created": "TF-update",
			//			"For":     "Test-update",
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"private_ip_addresses.#": "4",
			//			"network_interface_name": name,
			//			"description":            "Test For Terraform Update",
			//			"queue_number":           "2",
			//			"tags.%":                 "2",
			//			"tags.Created":           "TF-update",
			//			"tags.For":               "Test-update",
			//		}),
			//	),
			//},
		},
	})
}

func TestAccAlicloudEcsNetworkInterface_primary_ip_address(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":             "${local.vswitch_id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"primary_ip_address":     "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)}",
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
					"primary_ip_address": "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 7)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"primary_ip_address": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEcsNetworkInterface_secondary_private_ip_address_count(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":             "${local.vswitch_id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
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

func TestAccAlicloudEcsNetworkInterface_basic1(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":                         "${local.vswitch_id}",
					"security_groups":                    []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":                  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":                        name,
					"private_ip":                         "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}",
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
						"security_groups.#":                  "1",
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

func TestAccAlicloudEcsNetworkInterface_basic2(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":        "${local.vswitch_id}",
					"security_groups":   []string{"${alicloud_security_group.default.id}"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":       name,
					"private_ip":        "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)}",
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

func TestAccAlicloudEcsNetworkInterface_basic3(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":           "${local.vswitch_id}",
					"security_groups":      []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":          name,
					"private_ip_addresses": []string{"${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 7)}"},
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

func TestAccAlicloudEcsNetworkInterface_basic4(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsNetworkInterfaceBasicDependence)
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
					"vswitch_id":           "${local.vswitch_id}",
					"security_groups":      []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":          name,
					"private_ip_addresses": []string{"${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 3)}", "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 7)}"},
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

func AlicloudEcsNetworkInterfaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}

func TestAccAlicloudEcsNetworkInterface_unit(t *testing.T) {
	resourceName := "alicloud_ecs_network_interface"
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	attributes := map[string]interface{}{
		"network_interface_name": "network_interface_name",
		"vswitch_id":             "vswitch_id",
		"security_group_ids":     []string{"security_group_ids"},
		"resource_group_id":"resource_group_id",
		"description":"description",
	}
	for key, value := range attributes {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"NetworkInterfaceSets":map[string]interface{}{
			"NetworkInterfaceSet":[]interface{}{
				map[string]interface{}{
					"NetworkInterfaceName":"network_interface_name",
					"VSwitchId":"vswitch_id",
					"Description":"description",
					"ResourceGroupId":"resource_group_id",
					"Status":"Available",
					"NetworkInterfaceId":"MockId",
					"PrivateIpSets": map[string]interface{}{
						"PrivateIpSet": []interface{}{},
					},
					"SecurityGroupIds":map[string]interface{}{
						"SecurityGroupId": []interface{}{},
					},
				},

			},
		},
		"TagResources":map[string]interface{}{
			"TagResource":[]interface{}{},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["NetworkInterfaceId"] = "MockId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
	}
	t.Run("Create", func(t *testing.T) {
		// Client Error
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err = resourceAlicloudEcsNetworkInterfaceCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)


		retryFlag, noRetryFlag, abnormal := true, true, true
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		for {
			err = resourceAlicloudEcsNetworkInterfaceCreate(dCreate, rawClient)
			if abnormal {
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for key, _ := range attributes {
					assert.False(t, dCreate.HasChange(key))
				}
				break
			}
		}
		patches.Reset()
	})

	t.Run("Update", func(t *testing.T) {
		rand1 := acctest.RandIntRange(1000, 9999)
		// Client Error
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})

		err := resourceAlicloudEcsNetworkInterfaceUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		// DoRequest
		retryFlag,noRetryFlag,abnormal := false,false,false
		targetFunc := "ModifyNetworkInterfaceAttribute"

		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if targetFunc == *action && retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if targetFunc == *action && noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})

		// UpdateMoveResourceGroup
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"security_group_ids"} {
			switch p[resourceName].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(2)})
			case schema.TypeMap:
				diff.SetAttribute(fmt.Sprintf("%s.%",key), &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute(fmt.Sprintf("%s.%sKey1",key,key), &terraform.ResourceAttrDiff{Old: "", New: fmt.Sprintf("%sValue1",key)})
				diff.SetAttribute(fmt.Sprintf("%s.%sKey2",key,key), &terraform.ResourceAttrDiff{Old: "", New: fmt.Sprintf("%sValue2",key)})
			case schema.TypeList:
				diff.SetAttribute("security_group_ids.#",&terraform.ResourceAttrDiff{Old: "1", New: "2"})
				diff.SetAttribute("security_group_ids.0",&terraform.ResourceAttrDiff{Old: "prev_value0", New: "current_value0"})
				diff.SetAttribute("security_group_ids.1",&terraform.ResourceAttrDiff{Old: "", New: "current_value1"})
			case schema.TypeSet:
				diff.SetAttribute(fmt.Sprintf("%s.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field, _ := range p[resourceName].Schema[key].Elem.(*schema.Resource).Schema {
						diff.SetAttribute(fmt.Sprintf("%s.%d.%s", key, rand1, field), &terraform.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"})
					}
				}
			}
		}
		d, _ = schema.InternalMap(p[resourceName].Schema).Data(dCreate.State(), diff)
		d.SetId(dCreate.Id())

		for {
			err = resourceAlicloudEcsNetworkInterfaceUpdate(d, rawClient)
			if abnormal {
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for key, _ := range attributes {
					assert.False(t, dCreate.HasChange(key))
				}
				break
			}
		}
		patches.Reset()

	})
}
