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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete ALB LoadBalancer (%s): %s", item["LoadBalancerId"].(string), err)
		}
		log.Printf("[INFO] Delete ALB LoadBalancer success: %s ", item["LoadBalancerId"].(string))
	}
	return nil
}

func TestAccAliCloudAlbLoadBalancer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlbLoadBalancerBasicDependence0)
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
					"load_balancer_edition":  "Basic",
					"address_type":           "Internet",
					"vpc_id":                 "${alicloud_vpc.default.id}",
					"address_allocated_mode": "Fixed",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.zone_a.id}",
							"zone_id":    "${alicloud_vswitch.zone_a.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_b.id}",
							"zone_id":    "${alicloud_vswitch.zone_b.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_edition":          "Basic",
						"address_type":                   "Internet",
						"vpc_id":                         CHECKSET,
						"address_allocated_mode":         "Fixed",
						"load_balancer_billing_config.#": "1",
						"zone_mappings.#":                "2",
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
					"deletion_protection_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled": "false",
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
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
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
					"deletion_protection_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": name,
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
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
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
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
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
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.update.project}",
							"log_store":   "${alicloud_log_store.update.name}",
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
					"access_log_config": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_log_config.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
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
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.zone_a.id}",
							"zone_id":    "${alicloud_vswitch.zone_a.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_c.id}",
							"zone_id":    "${alicloud_vswitch.zone_c.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "LoadBalancer",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "LoadBalancer",
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

func TestAccAliCloudAlbLoadBalancer_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlbLoadBalancerBasicDependence0)
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
					"load_balancer_edition":       "Basic",
					"address_type":                "Internet",
					"vpc_id":                      "${alicloud_vpc.default.id}",
					"address_allocated_mode":      "Fixed",
					"address_ip_version":          "IPv4",
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"load_balancer_name":          name,
					"deletion_protection_enabled": "false",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": name,
						},
					},
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.zone_a.id}",
							"zone_id":    "${alicloud_vswitch.zone_a.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_b.id}",
							"zone_id":    "${alicloud_vswitch.zone_b.zone_id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "LoadBalancer",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_edition":            "Basic",
						"address_type":                     "Internet",
						"vpc_id":                           CHECKSET,
						"address_allocated_mode":           "Fixed",
						"address_ip_version":               "IPv4",
						"bandwidth_package_id":             CHECKSET,
						"resource_group_id":                CHECKSET,
						"load_balancer_name":               name,
						"deletion_protection_enabled":      "false",
						"load_balancer_billing_config.#":   "1",
						"modification_protection_config.#": "1",
						"access_log_config.#":              "1",
						"zone_mappings.#":                  "2",
						"tags.%":                           "2",
						"tags.Created":                     "TF",
						"tags.For":                         "LoadBalancer",
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

func TestAccAliCloudAlbLoadBalancer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlbLoadBalancerBasicDependence1)
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
					"load_balancer_edition":  "Basic",
					"address_type":           "Internet",
					"vpc_id":                 "${alicloud_vpc_ipv6_gateway.default.vpc_id}",
					"address_allocated_mode": "Fixed",
					"address_ip_version":     "DualStack",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_a.id}",
							"intranet_address": "192.168.10.1",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_b.id}",
							"zone_id":    "${alicloud_vswitch.zone_b.zone_id}",
						},
					},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_edition":          "Basic",
						"address_type":                   "Internet",
						"vpc_id":                         CHECKSET,
						"address_allocated_mode":         "Fixed",
						"address_ip_version":             "DualStack",
						"load_balancer_billing_config.#": "1",
						"zone_mappings.#":                "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_a.id}",
							"intranet_address": "192.168.10.1",
						},
						{
							"vswitch_id":       "${alicloud_vswitch.zone_c.id}",
							"zone_id":          "${alicloud_vswitch.zone_c.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_b.id}",
							"intranet_address": "192.168.192.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
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
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_type": "Internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_type": "Internet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
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
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": name,
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
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
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
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
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
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.update.project}",
							"log_store":   "${alicloud_log_store.update.name}",
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
					"access_log_config": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_log_config.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "LoadBalancer",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "LoadBalancer",
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

func TestAccAliCloudAlbLoadBalancer_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlbLoadBalancerBasicDependence1)
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
					"load_balancer_edition":       "Basic",
					"address_type":                "Internet",
					"vpc_id":                      "${alicloud_vpc_ipv6_gateway.default.vpc_id}",
					"address_allocated_mode":      "Fixed",
					"address_ip_version":          "DualStack",
					"ipv6_address_type":           "Internet",
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"load_balancer_name":          name,
					"deletion_protection_enabled": "false",
					"load_balancer_billing_config": []map[string]interface{}{
						{
							"pay_type": "PayAsYouGo",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": name,
						},
					},
					"access_log_config": []map[string]interface{}{
						{
							"log_project": "${alicloud_log_store.default.project}",
							"log_store":   "${alicloud_log_store.default.name}",
						},
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_a.id}",
							"intranet_address": "192.168.10.1",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_b.id}",
							"zone_id":    "${alicloud_vswitch.zone_b.zone_id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "LoadBalancer",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_edition":            "Basic",
						"address_type":                     "Internet",
						"vpc_id":                           CHECKSET,
						"address_allocated_mode":           "Fixed",
						"address_ip_version":               "DualStack",
						"ipv6_address_type":                "Internet",
						"bandwidth_package_id":             CHECKSET,
						"resource_group_id":                CHECKSET,
						"load_balancer_name":               name,
						"deletion_protection_enabled":      "false",
						"load_balancer_billing_config.#":   "1",
						"modification_protection_config.#": "1",
						"access_log_config.#":              "1",
						"zone_mappings.#":                  "2",
						"tags.%":                           "2",
						"tags.Created":                     "TF",
						"tags.For":                         "LoadBalancer",
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

var AliCloudAlbLoadBalancerMap0 = map[string]string{
	"address_ip_version": CHECKSET,
	"resource_group_id":  CHECKSET,
	"dns_name":           CHECKSET,
	"status":             CHECKSET,
	"create_time":        CHECKSET,
}

var AliCloudAlbLoadBalancerMap1 = map[string]string{
	"address_ip_version": CHECKSET,
	"ipv6_address_type":  CHECKSET,
	"resource_group_id":  CHECKSET,
	"dns_name":           CHECKSET,
	"status":             CHECKSET,
	"create_time":        CHECKSET,
}

func AliCloudAlbLoadBalancerBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "zone_a" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.0.0/18"
  		zone_id      = data.alicloud_alb_zones.default.zones.0.id
	}

	resource "alicloud_vswitch" "zone_b" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.128.0/18"
  		zone_id      = data.alicloud_alb_zones.default.zones.1.id
	}

	resource "alicloud_vswitch" "zone_c" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/18"
  		zone_id      = data.alicloud_alb_zones.default.zones.2.id
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 1000
  		internet_charge_type = "PayByBandwidth"
	}

	resource "alicloud_log_project" "default" {
  		name = var.name
	}

	resource "alicloud_log_store" "default" {
  		project = alicloud_log_project.default.name
  		name    = var.name
	}

	resource "alicloud_log_project" "update" {
  		name = "${var.name}-update"
	}

	resource "alicloud_log_store" "update" {
  		project = alicloud_log_project.update.name
  		name    = "${var.name}-update"
	}
`, name)
}

func AliCloudAlbLoadBalancerBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name    = var.name
  		cidr_block  = "192.168.0.0/16"
  		enable_ipv6 = "true"
	}

	resource "alicloud_eip" "zone_a" {
	  bandwidth            = "10"
	  internet_charge_type = "PayByTraffic"
	}

	resource "alicloud_eip" "zone_b" {
	  bandwidth            = "10"
	  internet_charge_type = "PayByTraffic"
	}

	resource "alicloud_eipanycast_anycast_eip_address" "zone_c" {
	  bandwidth                = 200
	  service_location         = "international"
	  internet_charge_type     = "PayByTraffic"
	  payment_type             = "PayAsYouGo"
	}

	resource "alicloud_vswitch" "zone_a" {
  		vswitch_name         = var.name
  		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.0.0/18"
  		zone_id              = data.alicloud_alb_zones.default.zones.0.id
  		ipv6_cidr_block_mask = "6"
	}

	resource "alicloud_vswitch" "zone_b" {
  		vswitch_name         = var.name
  		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.128.0/18"
  		zone_id              = data.alicloud_alb_zones.default.zones.1.id
  		ipv6_cidr_block_mask = "8"
	}

	resource "alicloud_vswitch" "zone_c" {
  		vswitch_name         = var.name
  		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.192.0/24"
  		zone_id              = data.alicloud_alb_zones.default.zones.2.id
  		ipv6_cidr_block_mask = "18"
	}

	resource "alicloud_vpc_ipv6_gateway" "default" {
  		ipv6_gateway_name = var.name
  		vpc_id            = alicloud_vpc.default.id
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 1000
  		internet_charge_type = "PayByBandwidth"
	}

	resource "alicloud_log_project" "default" {
  		name = var.name
	}

	resource "alicloud_log_store" "default" {
  		project = alicloud_log_project.default.name
  		name    = var.name
	}

	resource "alicloud_log_project" "update" {
  		name = "${var.name}-update"
	}

	resource "alicloud_log_store" "update" {
  		project = alicloud_log_project.update.name
  		name    = "${var.name}-update"
	}
`, name)
}
