package alicloud

import (
	"fmt"
	"log"
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
			log.Printf("[ERROR] Describe NetworkInterface failed, %#v", err)
			continue
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
					"vswitch_id":             "${alicloud_vswitch.default.id}",
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
					"private_ip_addresses": []string{"192.168.0.3", "192.168.0.5", "192.168.0.6"},
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
					"private_ip_addresses": []string{"192.168.0.3", "192.168.0.5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_addresses.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip_addresses":   []string{"192.168.0.3", "192.168.0.5", "192.168.0.6", "192.168.0.8"},
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
						"private_ip_addresses.#": "4",
						"network_interface_name": name,
						"description":            "Test For Terraform Update",
						"queue_number":           "2",
						"tags.%":                 "2",
						"tags.Created":           "TF-update",
						"tags.For":               "Test-update",
					}),
				),
			},
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
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"security_group_ids":     []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"primary_ip_address":     "192.168.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": CHECKSET,
						"vswitch_id":             CHECKSET,
						"security_group_ids.#":   "1",
						"resource_group_id":      CHECKSET,
						"primary_ip_address":     "192.168.0.2",
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
					"primary_ip_address": "192.168.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"primary_ip_address": "192.168.0.4",
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
					"vswitch_id":             "${alicloud_vswitch.default.id}",
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

func AlicloudEcsNetworkInterfaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    vswitch_name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}
data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}
