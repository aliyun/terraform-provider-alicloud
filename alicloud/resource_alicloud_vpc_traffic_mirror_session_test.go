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
	resource.AddTestSweepers(
		"alicloud_vpc_traffic_mirror_session",
		&resource.Sweeper{
			Name: "alicloud_vpc_traffic_mirror_session",
			F:    testSweepVpcTrafficMirrorSession,
		})
}

func testSweepVpcTrafficMirrorSession(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.VpcTrafficMirrorSupportRegions) {
		log.Printf("[INFO] Skipping Vpc Traffic Mirror Session unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListTrafficMirrorSessions"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.TrafficMirrorSessions", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.TrafficMirrorSessions", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TrafficMirrorSessionName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpc Traffic Mirror Session: %s", item["TrafficMirrorSessionName"].(string))
				continue
			}
			action := "DeleteTrafficMirrorSession"
			request := map[string]interface{}{
				"TrafficMirrorSessionId": item["TrafficMirrorSessionId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpc Traffic Mirror Session (%s): %s", item["TrafficMirrorSessionName"].(string), err)
			}
			log.Printf("[INFO] Delete Vpc Traffic Mirror Session success: %s ", item["TrafficMirrorSessionName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudVPCTrafficMirrorSession_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_session.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorSessionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorSession")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorsession-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorSessionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_USE_HOLOGRAPHIC_ACCOUNT")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":                           "1",
					"traffic_mirror_session_description": "${var.name}",
					"traffic_mirror_session_name":        "${var.name}",
					"traffic_mirror_target_id":           "${data.alicloud_ecs_network_interfaces.default.ids.0}",
					"traffic_mirror_source_ids":          []string{"${data.alicloud_ecs_network_interfaces.default.ids.1}"},
					"traffic_mirror_filter_id":           "${alicloud_vpc_traffic_mirror_filter.default.0.id}",
					"traffic_mirror_target_type":         "NetworkInterface",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":                           "1",
						"traffic_mirror_session_description": name,
						"traffic_mirror_session_name":        name,
						"traffic_mirror_target_id":           CHECKSET,
						"traffic_mirror_source_ids.#":        "1",
						"traffic_mirror_filter_id":           CHECKSET,
						"traffic_mirror_target_type":         "NetworkInterface",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_session_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_session_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_session_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_session_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_target_id": "${data.alicloud_ecs_network_interfaces.default.ids.2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_target_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_target_type": "SLB",
					"traffic_mirror_target_id":   "${alicloud_slb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_target_type": "SLB",
						"traffic_mirror_target_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_network_id": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_network_id": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_source_ids": []string{"${data.alicloud_ecs_network_interfaces.default.ids.0}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_source_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_filter_id":           "${alicloud_vpc_traffic_mirror_filter.default.0.id}",
					"traffic_mirror_target_id":           "${data.alicloud_ecs_network_interfaces.default.ids.3}",
					"traffic_mirror_source_ids":          []string{"${data.alicloud_ecs_network_interfaces.default.ids.1}"},
					"traffic_mirror_target_type":         "NetworkInterface",
					"traffic_mirror_session_description": "${var.name}",
					"traffic_mirror_session_name":        "${var.name}",
					"enabled":                            "false",
					"virtual_network_id":                 "20",
					"priority":                           "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_id":           CHECKSET,
						"traffic_mirror_source_ids.#":        "1",
						"traffic_mirror_target_id":           CHECKSET,
						"traffic_mirror_target_type":         "NetworkInterface",
						"traffic_mirror_session_description": name,
						"traffic_mirror_session_name":        name,
						"enabled":                            "false",
						"virtual_network_id":                 "20",
						"priority":                           "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudVPCTrafficMirrorSession_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_session.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorSessionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorSession")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorsession-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorSessionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_USE_HOLOGRAPHIC_ACCOUNT")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":                           "1",
					"traffic_mirror_session_description": "${var.name}",
					"traffic_mirror_session_name":        "${var.name}",
					"traffic_mirror_target_id":           "${data.alicloud_ecs_network_interfaces.default.ids.0}",
					"traffic_mirror_source_ids":          []string{"${data.alicloud_ecs_network_interfaces.default.ids.1}"},
					"traffic_mirror_filter_id":           "${alicloud_vpc_traffic_mirror_filter.default.0.id}",
					"traffic_mirror_target_type":         "NetworkInterface",
					"dry_run":                            "false",
					"enabled":                            "true",
					"virtual_network_id":                 "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":                           "1",
						"traffic_mirror_session_description": name,
						"traffic_mirror_session_name":        name,
						"traffic_mirror_target_id":           CHECKSET,
						"traffic_mirror_source_ids.#":        "1",
						"traffic_mirror_filter_id":           CHECKSET,
						"traffic_mirror_target_type":         "NetworkInterface",
						"dry_run":                            "false",
						"enabled":                            "true",
						"virtual_network_id":                 "10",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVPCTrafficMirrorSessionMap0 = map[string]string{
	"dry_run": NOSET,
	"status":  CHECKSET,
}

func AlicloudVPCTrafficMirrorSessionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "intranet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = alicloud_vswitch.default.id
}

resource "alicloud_vpc_traffic_mirror_filter" "default" {
  count                      = 2
  traffic_mirror_filter_name = var.name
}

data "alicloud_ecs_network_interfaces" "default" {
  tags = {
    tf-testacc = "vpctrafficmirrorsession"
  }
}`, name)
}
