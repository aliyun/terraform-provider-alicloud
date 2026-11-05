package alb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlbLoadBalancer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbLoadBalancerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vpc_id":                         acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"access_log_config": acctest.REMOVEKEY,
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
							"eip_type":   "Common",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_c.id}",
							"zone_id":    "${alicloud_vswitch.zone_c.zone_id}",
							"eip_type":   "Common",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbLoadBalancerBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vpc_id":                           acctest.CHECKSET,
						"address_allocated_mode":           "Fixed",
						"address_ip_version":               "IPv4",
						"bandwidth_package_id":             acctest.CHECKSET,
						"resource_group_id":                acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbLoadBalancerBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vpc_id":                         acctest.CHECKSET,
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
							"allocation_id":    "${alicloud_eip.zone_c.id}",
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
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_b.id}",
							"zone_id":          "${alicloud_vswitch.zone_b.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_b.id}",
							"intranet_address": "192.168.176.169",
						},
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"eip_type":         "Common",
							"allocation_id":    "${alicloud_eip.zone_a.id}",
							"intranet_address": "192.168.10.1",
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
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"intranet_address": "192.168.10.1",
						},
						{
							"vswitch_id":       "${alicloud_vswitch.zone_c.id}",
							"zone_id":          "${alicloud_vswitch.zone_c.zone_id}",
							"intranet_address": "192.168.192.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Intranet",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":       "${alicloud_vswitch.zone_a.id}",
							"zone_id":          "${alicloud_vswitch.zone_a.zone_id}",
							"intranet_address": "192.168.10.1",
						},
						{
							"vswitch_id": "${alicloud_vswitch.zone_b.id}",
							"zone_id":    "${alicloud_vswitch.zone_b.zone_id}",
						},
					},
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
						"resource_group_id": acctest.CHECKSET,
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
					"access_log_config": acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbLoadBalancerMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbLoadBalancerBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vpc_id":                           acctest.CHECKSET,
						"address_allocated_mode":           "Fixed",
						"address_ip_version":               "DualStack",
						"ipv6_address_type":                "Internet",
						"bandwidth_package_id":             acctest.CHECKSET,
						"resource_group_id":                acctest.CHECKSET,
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
	"address_ip_version": acctest.CHECKSET,
	"resource_group_id":  acctest.CHECKSET,
	"dns_name":           acctest.CHECKSET,
	"status":             acctest.CHECKSET,
	"create_time":        acctest.CHECKSET,
}

var AliCloudAlbLoadBalancerMap1 = map[string]string{
	"address_ip_version": acctest.CHECKSET,
	"ipv6_address_type":  acctest.CHECKSET,
	"resource_group_id":  acctest.CHECKSET,
	"dns_name":           acctest.CHECKSET,
	"status":             acctest.CHECKSET,
	"create_time":        acctest.CHECKSET,
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

	resource "alicloud_eip" "zone_c" {
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
