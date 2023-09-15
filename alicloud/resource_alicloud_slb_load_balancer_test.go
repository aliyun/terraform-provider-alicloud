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
		if !sweepAll() {
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
		}
		log.Printf("[INFO] Deleting SLB: %s (%s)", name, id)
		if err := service.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSLBLoadBalancer_basic0(t *testing.T) {
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
	"instance_charge_type":           "PayBySpec",
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
					"load_balancer_name":             name,
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
						"load_balancer_name":             name,
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

func TestAccAlicloudSLBLoadBalancer_basic3(t *testing.T) {
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
					"internet_charge_type": "PayByTraffic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"load_balancer_spec":   "slb.s3.small",
						"internet_charge_type": "PayByTraffic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PayByCLCU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PayByCLCU",
						"load_balancer_spec":   "slb.lcu.elastic",
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

func TestAccAlicloudSLBLoadBalancer_basic4(t *testing.T) {
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
					"address_type":         "intranet",
					"load_balancer_name":   name,
					"vswitch_id":           "${data.alicloud_vswitches.default.ids[0]}",
					"address":              "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)}",
					"master_zone_id":       "${data.alicloud_slb_zones.default.zones.0.master_zone_id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"slave_zone_id":        "${data.alicloud_slb_zones.default.zones.0.slave_zone_id}",
					"instance_charge_type": "PayByCLCU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type":       "intranet",
						"load_balancer_name": name,
						"specification":      "slb.lcu.elastic",
						"load_balancer_spec": "slb.lcu.elastic",
						"vswitch_id":         CHECKSET,
						"address":            CHECKSET,
						"master_zone_id":     CHECKSET,
						"payment_type":       "PayAsYouGo",
						"resource_group_id":  CHECKSET,
						"slave_zone_id":      CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":             name,
						"modification_protection_status": "ConsoleProtection",
						"tags.%":                         "2",
						"tags.Created":                   "TF-update",
						"tags.For":                       "Test-update",
					}),
				),
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
 
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
	slave_zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "^default$"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_slb_zones.default.zones.0.id
}

`, name)
}

// Test Slb LoadBalancer. >>> Resource test cases, automatically generated.
// Case 4052
func TestAccAlicloudSlbLoadBalancer_basic4052(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4052)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4052)
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
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
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
					"instance_charge_type": "PayByCLCU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PayByCLCU",
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
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_reason": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.rg1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
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
					"modification_protection_reason": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "test-update",
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
					"tags": []map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.#": "0",
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
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_reason": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
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
					"modification_protection_reason": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                         "active",
					"address_ip_version":             "ipv4",
					"address":                        "10.0.10.1",
					"instance_charge_type":           "PayByCLCU",
					"vswitch_id":                     "vsw-bp1fpj92chcwmdla73oxg",
					"slave_zone_id":                  "cn-hangzhou-h",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_name":             name + "_update",
					"delete_protection":              "off",
					"vpc_id":                         "vpc-bp18uccoyc62e4gk6033e",
					"payment_type":                   "PayAsYouGo",
					"modification_protection_reason": "test",
					"address_type":                   "intranet",
					"master_zone_id":                 "cn-hangzhou-j",
					"resource_group_id":              "${alicloud_resource_manager_resource_group.rg1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                         "active",
						"address_ip_version":             "ipv4",
						"address":                        "10.0.10.1",
						"instance_charge_type":           "PayByCLCU",
						"vswitch_id":                     "vsw-bp1fpj92chcwmdla73oxg",
						"slave_zone_id":                  "cn-hangzhou-h",
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_name":             name + "_update",
						"delete_protection":              "off",
						"vpc_id":                         "vpc-bp18uccoyc62e4gk6033e",
						"payment_type":                   "PayAsYouGo",
						"modification_protection_reason": "test",
						"address_type":                   "intranet",
						"master_zone_id":                 "cn-hangzhou-j",
						"resource_group_id":              CHECKSET,
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudSlbLoadBalancerMap4052 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudSlbLoadBalancerBasicDependence4052(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "${var.name}1"
}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "slb01"
  resource_group_name = "${var.name}2"
}


`, name)
}

// Case 4064
func TestAccAlicloudSlbLoadBalancer_basic4064(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4064)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4064)
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
					"load_balancer_name":   name,
					"instance_charge_type": "PayBySpec",
					"load_balancer_spec":   "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"instance_charge_type": "PayBySpec",
						"load_balancer_spec":   "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PayBySpec",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PayBySpec",
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
					"modification_protection_reason": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"modification_protection_reason": "test",
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
						"modification_protection_reason": "",
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
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
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
					"internet_charge_type": "paybybandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "paybybandwidth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1",
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
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
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
					"bandwidth": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":           "PayBySpec",
					"slave_zone_id":                  "cn-hangzhou-h",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_name":             name + "_update",
					"delete_protection":              "off",
					"payment_type":                   "PayAsYouGo",
					"modification_protection_reason": "test",
					"address_type":                   "internet",
					"master_zone_id":                 "cn-hangzhou-j",
					"address_ip_version":             "ipv4",
					"internet_charge_type":           "paybybandwidth",
					"load_balancer_spec":             "slb.s1.small",
					"bandwidth":                      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":           "PayBySpec",
						"slave_zone_id":                  "cn-hangzhou-h",
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_name":             name + "_update",
						"delete_protection":              "off",
						"payment_type":                   "PayAsYouGo",
						"modification_protection_reason": "test",
						"address_type":                   "internet",
						"master_zone_id":                 "cn-hangzhou-j",
						"address_ip_version":             "ipv4",
						"internet_charge_type":           "paybybandwidth",
						"load_balancer_spec":             "slb.s1.small",
						"bandwidth":                      "1",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudSlbLoadBalancerMap4064 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudSlbLoadBalancerBasicDependence4064(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4065
func TestAccAlicloudSlbLoadBalancer_basic4065(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4065)
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
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "paybybandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "paybybandwidth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
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
					"bandwidth": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_spec": "slb.s1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_spec": "slb.s1.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
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
					"auto_pay":             "true",
					"slave_zone_id":        "cn-hangzhou-h",
					"load_balancer_name":   name + "_update",
					"payment_type":         "Subscription",
					"address_type":         "internet",
					"master_zone_id":       "cn-hangzhou-j",
					"address_ip_version":   "ipv4",
					"internet_charge_type": "paybybandwidth",
					"load_balancer_spec":   "slb.s1.small",
					"bandwidth":            "1",
					"pricing_cycle":        "month",
					"duration":             "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_pay":             "true",
						"slave_zone_id":        "cn-hangzhou-h",
						"load_balancer_name":   name + "_update",
						"payment_type":         "Subscription",
						"address_type":         "internet",
						"master_zone_id":       "cn-hangzhou-j",
						"address_ip_version":   "ipv4",
						"internet_charge_type": "paybybandwidth",
						"load_balancer_spec":   "slb.s1.small",
						"bandwidth":            "1",
						"pricing_cycle":        "month",
						"duration":             "1",
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
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudSlbLoadBalancerMap4065 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudSlbLoadBalancerBasicDependence4065(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4052  twin
func TestAccAlicloudSlbLoadBalancer_basic4052_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4052)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4052)
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
					"status":                         "active",
					"address_ip_version":             "ipv4",
					"address":                        "10.0.10.5",
					"instance_charge_type":           "PayByCLCU",
					"vswitch_id":                     "vsw-bp1fpj92chcwmdla73oxg",
					"slave_zone_id":                  "cn-hangzhou-h",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_name":             name,
					"delete_protection":              "off",
					"vpc_id":                         "vpc-bp18uccoyc62e4gk6033e",
					"payment_type":                   "PayAsYouGo",
					"modification_protection_reason": "test-update",
					"address_type":                   "intranet",
					"master_zone_id":                 "cn-hangzhou-j",
					"resource_group_id":              "${alicloud_resource_manager_resource_group.rg1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                         "active",
						"address_ip_version":             "ipv4",
						"address":                        "10.0.10.5",
						"instance_charge_type":           "PayByCLCU",
						"vswitch_id":                     "vsw-bp1fpj92chcwmdla73oxg",
						"slave_zone_id":                  "cn-hangzhou-h",
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_name":             name,
						"delete_protection":              "off",
						"vpc_id":                         "vpc-bp18uccoyc62e4gk6033e",
						"payment_type":                   "PayAsYouGo",
						"modification_protection_reason": "test-update",
						"address_type":                   "intranet",
						"master_zone_id":                 "cn-hangzhou-j",
						"resource_group_id":              CHECKSET,
						"tags.%":                         "2",
						"tags.Created":                   "TF",
						"tags.For":                       "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

// Case 4064  twin
func TestAccAlicloudSlbLoadBalancer_basic4064_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4064)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4064)
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
					"instance_charge_type":           "PayBySpec",
					"slave_zone_id":                  "cn-hangzhou-h",
					"modification_protection_status": "ConsoleProtection",
					"load_balancer_name":             name,
					"delete_protection":              "off",
					"payment_type":                   "PayAsYouGo",
					"modification_protection_reason": "test-update",
					"address_type":                   "internet",
					"master_zone_id":                 "cn-hangzhou-j",
					"address_ip_version":             "ipv4",
					"internet_charge_type":           "paybybandwidth",
					"load_balancer_spec":             "slb.s1.small",
					"bandwidth":                      "1",
					"pricing_cycle":                  "month",
					"duration":                       "1",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":           "PayBySpec",
						"slave_zone_id":                  "cn-hangzhou-h",
						"modification_protection_status": "ConsoleProtection",
						"load_balancer_name":             name,
						"delete_protection":              "off",
						"payment_type":                   "PayAsYouGo",
						"modification_protection_reason": "test-update",
						"address_type":                   "internet",
						"master_zone_id":                 "cn-hangzhou-j",
						"address_ip_version":             "ipv4",
						"internet_charge_type":           "paybybandwidth",
						"load_balancer_spec":             "slb.s1.small",
						"bandwidth":                      "1",
						"pricing_cycle":                  "month",
						"duration":                       "1",
						"tags.%":                         "2",
						"tags.Created":                   "TF",
						"tags.For":                       "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

// Case 4065  twin
func TestAccAlicloudSlbLoadBalancer_basic4065_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudSlbLoadBalancerMap4065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbLoadBalancerBasicDependence4065)
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
					"auto_pay":             "true",
					"slave_zone_id":        "cn-hangzhou-h",
					"load_balancer_name":   name,
					"payment_type":         "Subscription",
					"address_type":         "internet",
					"master_zone_id":       "cn-hangzhou-j",
					"address_ip_version":   "ipv4",
					"internet_charge_type": "paybybandwidth",
					"load_balancer_spec":   "slb.s2.small",
					"bandwidth":            "1",
					"pricing_cycle":        "year",
					"duration":             "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_pay":             "true",
						"slave_zone_id":        "cn-hangzhou-h",
						"load_balancer_name":   name,
						"payment_type":         "Subscription",
						"address_type":         "internet",
						"master_zone_id":       "cn-hangzhou-j",
						"address_ip_version":   "ipv4",
						"internet_charge_type": "paybybandwidth",
						"load_balancer_spec":   "slb.s2.small",
						"bandwidth":            "1",
						"pricing_cycle":        "year",
						"duration":             "2",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle"},
			},
		},
	})
}

// Test Slb LoadBalancer. <<< Resource test cases, automatically generated.
