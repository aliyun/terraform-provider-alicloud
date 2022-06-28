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
		"alicloud_alb_load_balancer",
		&resource.Sweeper{
			Name: "alicloud_alb_load_balancer",
			F:    testSweepAlbLoadBalancer,
		})
}

func testSweepAlbLoadBalancer(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListLoadBalancers"
	request := map[string]interface{}{
		"MaxResults": PageSizeXLarge,
	}
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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

	resp, err := jsonpath.Get("$.LoadBalancers", response)

	if formatInt(response["TotalCount"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.LoadBalancers", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if _, ok := item["LoadBalancerName"]; !ok {
			continue
		}
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["LoadBalancerName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ALB LoadBalancer: %s", item["LoadBalancerName"].(string))
			continue
		}

		action := "DeleteLoadBalancer"
		request := map[string]interface{}{
			"LoadBalancerId": item["LoadBalancerId"],
		}
		request["ClientToken"] = buildClientToken("DeleteLoadBalancer")
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ALB LoadBalancer (%s): %s", item["LoadBalancerId"].(string), err)
		}
		log.Printf("[INFO] Delete ALB LoadBalancer success: %s ", item["LoadBalancerId"].(string))
	}
	return nil
}

func TestAccAlicloudALBLoadBalancer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudALBLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBLoadBalancerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":                 "${local.vpc_id}",
					"address_type":           "Internet",
					"address_allocated_mode": "Fixed",
					"load_balancer_name":     "${var.name}",
					"load_balancer_edition":  "Basic",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                         CHECKSET,
						"address_type":                   "Internet",
						"address_allocated_mode":         "Fixed",
						"load_balancer_name":             name,
						"load_balancer_edition":          "Basic",
						"load_balancer_billing_config.#": "1",
						"zone_mappings.#":                "2",
						"dns_name":                       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_edition": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_edition": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${local.log_project}",
							"log_store":   "${local.log_store}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_log_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "Test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.For":     "Test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": "TF_Test123.-",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name":          name,
					"deletion_protection_enabled": "false",
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "Test2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":               name,
						"deletion_protection_enabled":      "false",
						"modification_protection_config.#": "1",
						"tags.%":                           "2",
						"tags.Created":                     "TF2",
						"tags.For":                         "Test2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "deletion_protection_enabled"},
			},
		},
	})
}

func TestAccAlicloudALBLoadBalancer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudALBLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBLoadBalancerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":                 "${local.vpc_id}",
					"address_type":           "Internet",
					"address_allocated_mode": "Fixed",
					"load_balancer_name":     "${var.name}",
					"load_balancer_edition":  "Basic",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
					"deletion_protection_enabled": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": "TF_Test123.-",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                           CHECKSET,
						"address_type":                     "Internet",
						"address_allocated_mode":           "Fixed",
						"load_balancer_name":               name,
						"load_balancer_edition":            "Basic",
						"load_balancer_billing_config.#":   "1",
						"zone_mappings.#":                  "2",
						"deletion_protection_enabled":      "false",
						"resource_group_id":                CHECKSET,
						"modification_protection_config.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "deletion_protection_enabled"},
			},
		},
	})
}

var AlicloudALBLoadBalancerMap0 = map[string]string{}

func AlicloudALBLoadBalancerBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id =  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name              = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

locals {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id_1 =  data.alicloud_alb_zones.default.zones.0.id
 vswitch_id_1 =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
 zone_id_2 =  data.alicloud_alb_zones.default.zones.1.id
 vswitch_id_2 =  length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
 log_project = alicloud_log_project.default.name
 log_store =   alicloud_log_store.default.name
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
