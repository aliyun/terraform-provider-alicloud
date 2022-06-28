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
	resource.AddTestSweepers("alicloud_slb_load_balancer", &resource.Sweeper{
		Name: "alicloud_slb_load_balancer",
		F:    testSweepSLBs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepSLBs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	k8sPrefix := "kubernetes"

	action := "DescribeLoadBalancers"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	slbs := make([]map[string]interface{}, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error retrieving SLBs: %s", err)
		}
		resp, err := jsonpath.Get("$.LoadBalancers.LoadBalancer", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LoadBalancers.LoadBalancer", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			slbs = append(slbs, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	service := SlbService{client}
	vpcService := VpcService{client}
	csService := CsService{client}
	for _, loadBalancer := range slbs {
		name := fmt.Sprint(loadBalancer["LoadBalancerName"])
		id := fmt.Sprint(loadBalancer["LoadBalancerId"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(fmt.Sprint(loadBalancer["VpcId"]), fmt.Sprint(loadBalancer["VSwitchId"])); err == nil {
				skip = !need
			}

		}
		// If a slb tag key has prefix "kubernetes", this is a slb for k8s cluster and it should be deleted if cluster not exist.
		if skip {
			if resp, err := jsonpath.Get("$.Tags.Tag", loadBalancer); err == nil {
				tag, _ := resp.([]interface{})
				for _, v := range tag {
					t, _ := v.(map[string]interface{})
					if strings.HasPrefix(strings.ToLower(t["TagKey"].(string)), strings.ToLower(k8sPrefix)) {
						_, err := csService.DescribeCsKubernetes(name)
						if NotFoundError(err) {
							skip = false
						} else {
							skip = true
							break
						}
					}
				}
			}
		}
		if skip {
			log.Printf("[INFO] Skipping SLB: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting SLB: %s (%s)", name, id)
		if err := service.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSLBLoadBalancer_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence0)
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
					"load_balancer_name":   "${var.name}",
					"load_balancer_spec":   "slb.s3.small",
					"internet_charge_type": "PayByBandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"load_balancer_spec":   "slb.s3.small",
						"internet_charge_type": "PayByBandwidth",
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
					"bandwidth": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_protection": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_protection": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_protection": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_protection": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_status": "ConsoleProtection",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_status": "ConsoleProtection",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_reason": "tf-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s2.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s2.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
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
					"load_balancer_name":             "${var.name}",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_spec":             "slb.s1.small",
					"bandwidth":                      `1`,
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":             name,
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_spec":             "slb.s1.small",
						"bandwidth":                      `1`,
						"tags.%":                         "2",
						"tags.Created":                   "TF-update",
						"tags.For":                       "Test-update",
					}),
				),
			},
		},
	})
}

var AlicloudSlbLoadBalancerMap0 = map[string]string{
	"address":                        CHECKSET,
	"address_ip_version":             "ipv4",
	"address_type":                   "internet",
	"bandwidth":                      CHECKSET,
	"delete_protection":              "off",
	"internet_charge_type":           "PayByBandwidth",
	"load_balancer_name":             "",
	"master_zone_id":                 CHECKSET,
	"modification_protection_reason": "",
	"modification_protection_status": CHECKSET,
	"payment_type":                   "PayAsYouGo",
	"resource_group_id":              CHECKSET,
	"slave_zone_id":                  CHECKSET,
	"load_balancer_spec":             "slb.s1.small",
	"status":                         "active",
	"tags.#":                         "0",
	"vswitch_id":                     "",
}

func AlicloudSlbLoadBalancerBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}

func TestAccAlicloudSLBLoadBalancer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence1)
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
					"address_type":       "intranet",
					"load_balancer_name": "${var.name}",
					"load_balancer_spec": "slb.s1.small",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type":       "intranet",
						"load_balancer_name": name,
						"load_balancer_spec": "slb.s1.small",
						"vswitch_id":         CHECKSET,
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
					"delete_protection": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_protection": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_status": "NonProtection",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_status": "NonProtection",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s2.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s2.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
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
					"delete_protection":              "off",
					"load_balancer_name":             "${var.name}",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_spec":             "slb.s1.small",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_protection":              "off",
						"load_balancer_name":             name,
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_spec":             "slb.s1.small",
						"tags.%":                         "2",
						"tags.Created":                   "TF-update",
						"tags.For":                       "Test-update",
					}),
				),
			},
		},
	})
}

var AlicloudSlbLoadBalancerMap1 = map[string]string{
	"address":                        CHECKSET,
	"address_ip_version":             "ipv4",
	"address_type":                   "intranet",
	"bandwidth":                      CHECKSET,
	"delete_protection":              "off",
	"internet_charge_type":           "PayByTraffic",
	"load_balancer_name":             "",
	"master_zone_id":                 CHECKSET,
	"modification_protection_reason": "",
	"modification_protection_status": CHECKSET,
	"payment_type":                   "PayAsYouGo",
	"resource_group_id":              CHECKSET,
	"slave_zone_id":                  CHECKSET,
	"load_balancer_spec":             "slb.s1.small",
	"status":                         "active",
	"tags.#":                         "0",
	"vswitch_id":                     CHECKSET,
}

func AlicloudSlbLoadBalancerBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_slb_zones.default.zones.0.id
}
`, name)
}

func TestAccAlicloudSLBLoadBalancer_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence2)
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
					"address_type":                   "intranet",
					"name":                           name,
					"specification":                  "slb.s1.small",
					"vswitch_id":                     "${data.alicloud_vswitches.default.ids[0]}",
					"address":                        "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)}",
					"master_zone_id":                 "${data.alicloud_slb_zones.default.zones.0.master_zone_id}",
					"modification_protection_status": "ConsoleProtection",
					"modification_protection_reason": name,
					"payment_type":                   "PayAsYouGo",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"slave_zone_id":                  "${data.alicloud_slb_zones.default.zones.0.slave_zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type":                   "intranet",
						"name":                           name,
						"specification":                  "slb.s1.small",
						"vswitch_id":                     CHECKSET,
						"address":                        CHECKSET,
						"master_zone_id":                 CHECKSET,
						"modification_protection_status": "ConsoleProtection",
						"modification_protection_reason": name,
						"payment_type":                   "PayAsYouGo",
						"resource_group_id":              CHECKSET,
						"slave_zone_id":                  CHECKSET,
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

var AlicloudSlbLoadBalancerMap2 = map[string]string{
	"address_ip_version":   "ipv4",
	"address_type":         "intranet",
	"bandwidth":            CHECKSET,
	"delete_protection":    "off",
	"internet_charge_type": "PayByTraffic",
	"resource_group_id":    CHECKSET,
	"slave_zone_id":        CHECKSET,
	"load_balancer_spec":   "slb.s1.small",
	"status":               "active",
	"tags.#":               "0",
}

func AlicloudSlbLoadBalancerBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "^default$"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_slb_zones.default.zones.0.id
}

`, name)
}
