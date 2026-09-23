// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa LoadBalancer. >>> Resource test cases, automatically generated.
// Case test_1 11682
func TestAccAliCloudEsaLoadBalancer_basic11682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaLoadBalancerMap11682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaLoadBalancerBasicDependence11682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fallback_pool":      "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}",
					"default_pools":      []string{"${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}"},
					"load_balancer_name": "lb." + name + ".com",
					"site_id":            "${alicloud_esa_site.resource_Site_OriginPool_test.id}",
					"steering_policy":    "geo",
					"monitor": []map[string]interface{}{
						{
							"path":              "/healthcheck",
							"type":              "HTTP",
							"header":            "{\\\"Host\\\":[\\\"example1.com\\\"]}",
							"expected_codes":    "2xx",
							"follow_redirects":  "true",
							"timeout":           "5",
							"port":              "80",
							"consecutive_up":    "3",
							"consecutive_down":  "5",
							"method":            "GET",
							"interval":          "60",
							"monitoring_region": "Global",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fallback_pool":      CHECKSET,
						"default_pools.#":    "1",
						"load_balancer_name": "lb." + name + ".com",
						"site_id":            CHECKSET,
						"steering_policy":    "geo",
						"monitor.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"adaptive_routing": []map[string]interface{}{
						{
							"failover_across_pools": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adaptive_routing.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_pools": []string{"${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}", "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}", "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_3.origin_pool_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_pools.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
					"fallback_pool": "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fallback_pool": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor": []map[string]interface{}{
						{
							"path":              "/healthcheck1",
							"type":              "HTTPS",
							"header":            "{\\\"Host\\\":[\\\"example2.com\\\"]}",
							"expected_codes":    "200",
							"follow_redirects":  "false",
							"timeout":           "6",
							"port":              "443",
							"monitoring_region": "OutsideChineseMainland",
							"consecutive_up":    "2",
							"consecutive_down":  "3",
							"method":            "HEAD",
							"interval":          "50",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"random_steering": []map[string]interface{}{
						{
							"pool_weights": map[string]interface{}{
								"\"${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}\"": "50",
							},
							"default_weight": "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"random_steering.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_pools": "{\\\"CNM,NAM\\\":[${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_pools": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_enable": "on",
							"rule":        "(http.host eq \\\"video.example.com\\\")",
							"sequence":    "1",
							"rule_name":   "rule1",
							"fixed_response": []map[string]interface{}{
								{
									"content_type": "text/plain",
									"status_code":  "200",
									"location":     "http://www.example.com/index.html",
									"message_body": "Hello World!",
								},
							},
						},
						{
							"rule_enable": "on",
							"overrides":   "{\\\"session_affinity\\\":\\\"cookie\\\"}",
							"rule":        "http.request.method eq \\\"GET\\\"",
							"sequence":    "2",
							"terminates":  "false",
							"rule_name":   "rule2",
						},
						{
							"rule_enable": "off",
							"overrides":   "{\\\"steering_policy\\\":\\\"geo\\\"}",
							"rule":        "http.request.uri eq \\\"/t\\\"",
							"sequence":    "3",
							"terminates":  "false",
							"rule_name":   "rule3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"session_affinity": "ip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"session_affinity": "ip",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"steering_policy": "order",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"steering_policy": "order",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_region_pools": "{\\\"CN-AH,CN-BJ\\\":[${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_region_pools": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "300",
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

func TestAccAliCloudEsaLoadBalancer_basic11682_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_load_balancer.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaLoadBalancerMap11682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaLoadBalancerBasicDependence11682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fallback_pool":      "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}",
					"default_pools":      []string{"${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}", "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}", "${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_3.origin_pool_id}"},
					"load_balancer_name": "lb." + name + ".com",
					"site_id":            "${alicloud_esa_site.resource_Site_OriginPool_test.id}",
					"steering_policy":    "geo",
					"description":        name,
					"enabled":            "true",
					"session_affinity":   "ip",
					"ttl":                "300",
					"region_pools":       "{\\\"CNM,NAM\\\":[${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}]}",
					"sub_region_pools":   "{\\\"CN-AH,CN-BJ\\\":[${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id}]}",
					"adaptive_routing": []map[string]interface{}{
						{
							"failover_across_pools": "true",
						},
					},
					"monitor": []map[string]interface{}{
						{
							"path":              "/healthcheck",
							"type":              "HTTP",
							"header":            "{\\\"Host\\\":[\\\"example1.com\\\"]}",
							"expected_codes":    "2xx",
							"follow_redirects":  "true",
							"timeout":           "5",
							"port":              "80",
							"consecutive_up":    "3",
							"consecutive_down":  "5",
							"method":            "GET",
							"interval":          "60",
							"monitoring_region": "Global",
						},
					},
					"random_steering": []map[string]interface{}{
						{
							"pool_weights": map[string]interface{}{
								"\"${alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id}\"": "50",
							},
							"default_weight": "30",
						},
					},
					"rules": []map[string]interface{}{
						{
							"rule_enable": "on",
							"rule":        "(http.host eq \\\"video.example.com\\\")",
							"sequence":    "1",
							"rule_name":   "rule1",
							"fixed_response": []map[string]interface{}{
								{
									"content_type": "text/plain",
									"status_code":  "200",
									"location":     "http://www.example.com/index.html",
									"message_body": "Hello World!",
								},
							},
						},
						{
							"rule_enable": "on",
							"overrides":   "{\\\"session_affinity\\\":\\\"cookie\\\"}",
							"rule":        "http.request.method eq \\\"GET\\\"",
							"sequence":    "2",
							"terminates":  "false",
							"rule_name":   "rule2",
						},
						{
							"rule_enable": "off",
							"overrides":   "{\\\"steering_policy\\\":\\\"geo\\\"}",
							"rule":        "http.request.uri eq \\\"/t\\\"",
							"sequence":    "3",
							"terminates":  "false",
							"rule_name":   "rule3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fallback_pool":      CHECKSET,
						"default_pools.#":    "3",
						"load_balancer_name": "lb." + name + ".com",
						"site_id":            CHECKSET,
						"steering_policy":    "geo",
						"description":        name,
						"enabled":            "true",
						"session_affinity":   "ip",
						"ttl":                "300",
						"region_pools":       CHECKSET,
						"sub_region_pools":   CHECKSET,
						"monitor.#":          "1",
						"adaptive_routing.#": "1",
						"random_steering.#":  "1",
						"rules.#":            "3",
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

var AliCloudEsaLoadBalancerMap11682 = map[string]string{
	"status":           CHECKSET,
	"load_balancer_id": CHECKSET,
	"ttl":              CHECKSET,
}

func AliCloudEsaLoadBalancerBasicDependence11682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_OriginPool_test" {
  site_name   = "${var.name}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_origin_pool" "resource_OriginPool_LoadBalancer_test_1_3" {
  origins {
    type    = "ip_domain"
    address = "www.example3.com"
	header  = "{\"Host\":[\"www.example3.com\"]}"
    enabled = true
    weight  = "30"
    name    = "origin3"
  }
  site_id          = alicloud_esa_site.resource_Site_OriginPool_test.id
  origin_pool_name = "testoriginpool3"
  enabled          = true
}

resource "alicloud_esa_origin_pool" "resource_OriginPool_LoadBalancer_test_1_2" {
  origins {
    type    = "ip_domain"
    address = "www.example1.com"
	header  = "{\"Host\":[\"www.example1.com\"]}"
    enabled = true
    weight  = "30"
    name    = "origin2"
  }
  site_id          = alicloud_esa_site.resource_Site_OriginPool_test.id
  origin_pool_name = "testoriginpool2"
  enabled          = true
}

resource "alicloud_esa_origin_pool" "resource_OriginPool_LoadBalancer_test_1_1" {
  origins {
    type    = "ip_domain"
    address = "www.example.com"
	header  = "{\"Host\":[\"www.example.com\"]}"
    enabled = true
    weight  = "30"
    name    = "origin1"
  }
  site_id          = alicloud_esa_site.resource_Site_OriginPool_test.id
  origin_pool_name = "testoriginpool1"
  enabled          = true
}
`, name)
}

// Test Esa LoadBalancer. <<< Resource test cases, automatically generated.
