package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudGaEndpointGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence0)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"listener_id":               CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_request_protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_request_protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/healthCheck",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/healthCheck",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval_seconds": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval_seconds": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_percentage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_percentage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":                     "${alicloud_eip_address.default.1.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_ecs_network_interface.default.id}",
							"type":     "ENI",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_ecs_network_interface.default.id}",
							"type":                         "ENI",
							"weight":                       "30",
							"sub_address":                  "${tolist(alicloud_ecs_network_interface.default.private_ip_addresses).0}",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
						{
							"endpoint":                     "${alicloud_ecs_network_interface.update.id}",
							"type":                         "ENI",
							"weight":                       "30",
							"sub_address":                  "${tolist(alicloud_ecs_network_interface.update.private_ip_addresses).0}",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "www.alicloud-provider.cn",
							"type":                         "Domain",
							"weight":                       "50",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_overrides.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence0)
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
					"accelerator_id":                "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":                   "${alicloud_ga_listener.default.id}",
					"endpoint_group_region":         defaultRegionToTest,
					"endpoint_group_type":           "virtual",
					"endpoint_request_protocol":     "HTTP",
					"health_check_enabled":          "true",
					"health_check_path":             "/healthCheck",
					"health_check_port":             "30",
					"health_check_protocol":         "HTTP",
					"health_check_interval_seconds": "5",
					"threshold_count":               "5",
					"traffic_percentage":            "30",
					"name":                          name,
					"description":                   name,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
						{
							"endpoint":                     "${alicloud_ecs_network_interface.default.id}",
							"type":                         "ENI",
							"weight":                       "30",
							"sub_address":                  "${tolist(alicloud_ecs_network_interface.default.private_ip_addresses).0}",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
						{
							"endpoint":                     "www.alicloud-provider.cn",
							"type":                         "Domain",
							"weight":                       "50",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
					},
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "60",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":                CHECKSET,
						"listener_id":                   CHECKSET,
						"endpoint_group_region":         defaultRegionToTest,
						"endpoint_group_type":           "virtual",
						"endpoint_request_protocol":     "HTTP",
						"health_check_enabled":          "true",
						"health_check_path":             "/healthCheck",
						"health_check_port":             "30",
						"health_check_protocol":         "HTTP",
						"health_check_interval_seconds": "5",
						"threshold_count":               "5",
						"traffic_percentage":            "30",
						"name":                          name,
						"description":                   name,
						"endpoint_configurations.#":     "3",
						"port_overrides.#":              "1",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence1)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"listener_id":               CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_request_protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_request_protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_protocol_version": "HTTP2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_protocol_version": "HTTP2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/healthCheck",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/healthCheck",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "tcp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "https",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "https",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval_seconds": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval_seconds": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_percentage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_percentage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":                     "${alicloud_eip_address.default.1.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "www.alicloud-provider.cn",
							"type":                         "Domain",
							"weight":                       "30",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "8080",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_overrides.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence1)
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
					"accelerator_id":                "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":                   "${alicloud_ga_listener.default.id}",
					"endpoint_group_region":         defaultRegionToTest,
					"endpoint_group_type":           "virtual",
					"endpoint_request_protocol":     "HTTPS",
					"endpoint_protocol_version":     "HTTP2",
					"health_check_enabled":          "true",
					"health_check_path":             "/healthCheck",
					"health_check_port":             "30",
					"health_check_protocol":         "http",
					"health_check_interval_seconds": "5",
					"threshold_count":               "5",
					"traffic_percentage":            "30",
					"name":                          name,
					"description":                   name,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "8080",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":                CHECKSET,
						"listener_id":                   CHECKSET,
						"endpoint_group_region":         defaultRegionToTest,
						"endpoint_group_type":           "virtual",
						"endpoint_request_protocol":     "HTTPS",
						"endpoint_protocol_version":     "HTTP2",
						"health_check_enabled":          "true",
						"health_check_path":             "/healthCheck",
						"health_check_port":             "30",
						"health_check_protocol":         "http",
						"health_check_interval_seconds": "5",
						"threshold_count":               "5",
						"traffic_percentage":            "30",
						"name":                          name,
						"description":                   name,
						"endpoint_configurations.#":     "1",
						"port_overrides.#":              "1",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence2)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint": "${alicloud_eip_address.default.0.ip_address}",
							"type":     "PublicIp",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"listener_id":               CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_request_protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_request_protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_protocol_version": "HTTP2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_protocol_version": "HTTP2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/healthCheck",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/healthCheck",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "tcp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "https",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "https",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval_seconds": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval_seconds": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_percentage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_percentage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":                     "${alicloud_eip_address.default.1.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "30",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":    "1.1.1.2",
							"type":        "IpTarget",
							"weight":      "50",
							"vpc_id":      "${alicloud_vpc.default.id}",
							"vswitch_ids": []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.update.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "8080",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_overrides.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence2)
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
					"accelerator_id":                "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":                   "${alicloud_ga_listener.default.id}",
					"endpoint_group_region":         defaultRegionToTest,
					"endpoint_group_type":           "virtual",
					"endpoint_request_protocol":     "HTTPS",
					"endpoint_protocol_version":     "HTTP2",
					"health_check_enabled":          "true",
					"health_check_path":             "/healthCheck",
					"health_check_port":             "30",
					"health_check_protocol":         "http",
					"health_check_interval_seconds": "5",
					"threshold_count":               "5",
					"traffic_percentage":            "30",
					"name":                          name,
					"description":                   name,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":                     "${alicloud_eip_address.default.1.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "30",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":    "1.1.1.2",
							"type":        "IpTarget",
							"weight":      "50",
							"vpc_id":      "${alicloud_vpc.default.id}",
							"vswitch_ids": []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.update.id}"},
						},
					},
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "8080",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":                CHECKSET,
						"listener_id":                   CHECKSET,
						"endpoint_group_region":         defaultRegionToTest,
						"endpoint_group_type":           "virtual",
						"endpoint_request_protocol":     "HTTPS",
						"endpoint_protocol_version":     "HTTP2",
						"health_check_enabled":          "true",
						"health_check_path":             "/healthCheck",
						"health_check_port":             "30",
						"health_check_protocol":         "http",
						"health_check_interval_seconds": "5",
						"threshold_count":               "5",
						"traffic_percentage":            "30",
						"name":                          name,
						"description":                   name,
						"endpoint_configurations.#":     "3",
						"port_overrides.#":              "1",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "EndpointGroup",
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

var AliCloudGaEndpointGroupMap0 = map[string]string{
	"endpoint_group_type":       CHECKSET,
	"endpoint_request_protocol": CHECKSET,
	"threshold_count":           CHECKSET,
	"endpoint_group_ip_list.#":  CHECKSET,
	"status":                    CHECKSET,
}

func AliCloudGaEndpointGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_ga_accelerators" "default" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_eip_address" "default" {
  		count                = 2
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}

	resource "alicloud_ecs_network_interface" "default" {
  		vswitch_id           = alicloud_vswitch.default.id
  		security_group_ids   = [alicloud_security_group.default.id]
  		private_ip_addresses = [cidrhost(alicloud_vswitch.default.cidr_block, 26)]
	}

	resource "alicloud_ecs_network_interface" "update" {
  		vswitch_id           = alicloud_vswitch.default.id
  		security_group_ids   = [alicloud_security_group.default.id]
  		private_ip_addresses = [cidrhost(alicloud_vswitch.default.cidr_block, 28)]
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Enhanced"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
		accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
		client_affinity = "SOURCE_IP"
		protocol        = "HTTP"
		name            = var.name
  		port_ranges {
    		from_port = "60"
    		to_port   = "60"
  		}
	}
`, name)
}

func AliCloudGaEndpointGroupBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_eip_address" "default" {
  		count                = 2
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQUNnSVa/sQNeb9pBN9NhkwTANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMzA4MDkwMzM4MThaFw0yODA4MDcwMzM4MThaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOgskr8dEfZYdjr0xaIqlCkmE802vABoj3SQNn3rLWnUj+1v
Wqbpsj6Bu61Scb8mtl/OZOOM7sgq0Q1hpdO8xvMGxTMuZ2bjX0EqCMqh4AvFofHL
a/iVD07hfoM1Jo8CEidh1uvcOuXP1TlaqU020x1TX3a3niJu4JVkmCkCOwAbWYuj
O8IsgBCsFaF9d4+C1JRYOtRbIHCNhd0sxG8AGovUDLvlkePeH5NF7DNvFXgGJ4iv
EQcY9pP08RBFUkaznOw/r64Up7zhLb+Ie4SyAvs1FulhMAmIXOcbsND39hJ+/WIP
8beWvIN1eCS8zcvgAvDgMkV8oqqVbQu1dqx5WuMCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALd0hFZAd2XHJgETbHQs
h4YUBNKxrIy6JiWfxffhIL1ZK5pI443DC4VRGfxVi3zWqs01WbNtJ2b1KdfSoovH
Zwi3hdMF1IwoAB/Y2sS4zjqS0H1od7MN9KKHes6bl3yCgpmaYs5cHbyg0IJHmeq3
rCgbKsvHfUwtzBNNPHlpANakAYd/5O1pztmUskWMUVaExfpMoQLo/AX9Lqm8pVjw
xs921I703l/E5zEnd3PVSYagy/KQJrwVt+wQZS11HsAryfO9kct/9f+c85VDo6Ht
iRirW/EnNPQRSno4z0V2x1Rn5+ZaoJo8cWzPvKrdfCG9TUozt4AR/LIudNLb6NNW
n7g=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6CySvx0R9lh2OvTFoiqUKSYTzTa8AGiPdJA2festadSP7W9a
pumyPoG7rVJxvya2X85k44zuyCrRDWGl07zG8wbFMy5nZuNfQSoIyqHgC8Wh8ctr
+JUPTuF+gzUmjwISJ2HW69w65c/VOVqpTTbTHVNfdreeIm7glWSYKQI7ABtZi6M7
wiyAEKwVoX13j4LUlFg61FsgcI2F3SzEbwAai9QMu+WR494fk0XsM28VeAYniK8R
Bxj2k/TxEEVSRrOc7D+vrhSnvOEtv4h7hLIC+zUW6WEwCYhc5xuw0Pf2En79Yg/x
t5a8g3V4JLzNy+AC8OAyRXyiqpVtC7V2rHla4wIDAQABAoIBABKGQ+sluaIrKrvH
feFTfmDOHfRYsqVhslh9jSt80THJePZb1SLOMJ+WIFBS7Kpwv0pjoF8bho3IBMgJ
i36aaFFJsABGao+mApqjbPIl+kdWLHarYWEDG6aSjVKQshPk+WfVAZ3uA3EEpSGf
XzS+9Bc56LsDKYXbzOV+kjlraSO35AMec3CpISdx4K1caEAhKX6it9bvPq4pSYXi
PQspba0Jv46VV7MaabVjLzsinz5/md4vxyYHNIJAukHUfwJIsVC9ZNxukwSw+CzE
MMO64ylq2DGokNerGsLetuViV8UWi7qmUmms2fAmchodW16olgNkYTz27+V/A42S
eex63pkCgYEA+CqKhqp3qPe2E9KVrycrwjoycxmhOn3Iz1xiN7uAEv+DzfKtfZVf
mcOIiqw4Z82RkgjHb9vJuTigKdDkB1zE2gSDnep44sDWJM/5nPjGlMgnkiJWJhci
CnD0P4d6cT5wyDt7Q0/tS6ql2UrCpW4ktw1AP0Rm/z/VBD8jGkVenjcCgYEA74DM
Z2Qmh3bPt1TykpOlw+H+sEuvlkYxqMlbtn3Rv3WgEPIBekOFrgP7n/uLW1Aizn8w
EhNBBAE8w5jvklqZWYbpFMJQc09eqUkI8aTbLooZbzYj1f3CrzBRKn1GoTPmN9V0
j9r+TbH3/5CEoqlsJdmeQPofuv5Qid2oEutZcrUCgYBuZ16hco0xmqJiRzlYZvDM
w99V3X0g7Hy947e+W6gqy4nzwZb1W9LgMWE5cEzXwViVw1oWpY0k3dBDSi9oJxlc
dM2pH3sQRgH+9pdyAis2XaVdGfGBmKEITCAdc0RBxSmfqva3h4NmOlD2TpAx0MJ8
vWRrwR6hR+CYtw4CzgG+GQKBgQDGmi5lugW9JUe/xeBUrbyyv0+cT1auLUz2ouq7
XIA23Mo74wJYqW9Lyp+4nTWFJeGHDK8G/hJWyNPjeomG+jvZombbQPrHc9SSWi7h
eowKfpfywZlb1M7AyTc1HacY+9l3CTlcJQPl16NHuEZUQFue02NIjGENhd+xQy4h
ainFVQKBgAoPs9ebtWbBOaqGpOnquzb7WibvW/ifOfzv5/aOhkbpp+KCjcSON6CB
QF3BEXMcNMGWlpPrd8PaxCAzR4MyU///ekJri2icS9lrQhGSz2TtYhdED4pv1Aag
7eTPl5L7xAwphCSwy8nfCKmvlqcX/MSJ7A+LHB/2hdbuuEOyhpbu
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Enhanced"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		name           = var.name
  		protocol       = "HTTPS"
  		port_ranges {
    		from_port = 8080
    		to_port   = 8080
  		}
  		certificates {
    		id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
  		}
	}
`, name, defaultRegionToTest)
}

func AliCloudGaEndpointGroupBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	resource "alicloud_eip_address" "default" {
  		count                = 2
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQUNnSVa/sQNeb9pBN9NhkwTANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMzA4MDkwMzM4MThaFw0yODA4MDcwMzM4MThaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOgskr8dEfZYdjr0xaIqlCkmE802vABoj3SQNn3rLWnUj+1v
Wqbpsj6Bu61Scb8mtl/OZOOM7sgq0Q1hpdO8xvMGxTMuZ2bjX0EqCMqh4AvFofHL
a/iVD07hfoM1Jo8CEidh1uvcOuXP1TlaqU020x1TX3a3niJu4JVkmCkCOwAbWYuj
O8IsgBCsFaF9d4+C1JRYOtRbIHCNhd0sxG8AGovUDLvlkePeH5NF7DNvFXgGJ4iv
EQcY9pP08RBFUkaznOw/r64Up7zhLb+Ie4SyAvs1FulhMAmIXOcbsND39hJ+/WIP
8beWvIN1eCS8zcvgAvDgMkV8oqqVbQu1dqx5WuMCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALd0hFZAd2XHJgETbHQs
h4YUBNKxrIy6JiWfxffhIL1ZK5pI443DC4VRGfxVi3zWqs01WbNtJ2b1KdfSoovH
Zwi3hdMF1IwoAB/Y2sS4zjqS0H1od7MN9KKHes6bl3yCgpmaYs5cHbyg0IJHmeq3
rCgbKsvHfUwtzBNNPHlpANakAYd/5O1pztmUskWMUVaExfpMoQLo/AX9Lqm8pVjw
xs921I703l/E5zEnd3PVSYagy/KQJrwVt+wQZS11HsAryfO9kct/9f+c85VDo6Ht
iRirW/EnNPQRSno4z0V2x1Rn5+ZaoJo8cWzPvKrdfCG9TUozt4AR/LIudNLb6NNW
n7g=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6CySvx0R9lh2OvTFoiqUKSYTzTa8AGiPdJA2festadSP7W9a
pumyPoG7rVJxvya2X85k44zuyCrRDWGl07zG8wbFMy5nZuNfQSoIyqHgC8Wh8ctr
+JUPTuF+gzUmjwISJ2HW69w65c/VOVqpTTbTHVNfdreeIm7glWSYKQI7ABtZi6M7
wiyAEKwVoX13j4LUlFg61FsgcI2F3SzEbwAai9QMu+WR494fk0XsM28VeAYniK8R
Bxj2k/TxEEVSRrOc7D+vrhSnvOEtv4h7hLIC+zUW6WEwCYhc5xuw0Pf2En79Yg/x
t5a8g3V4JLzNy+AC8OAyRXyiqpVtC7V2rHla4wIDAQABAoIBABKGQ+sluaIrKrvH
feFTfmDOHfRYsqVhslh9jSt80THJePZb1SLOMJ+WIFBS7Kpwv0pjoF8bho3IBMgJ
i36aaFFJsABGao+mApqjbPIl+kdWLHarYWEDG6aSjVKQshPk+WfVAZ3uA3EEpSGf
XzS+9Bc56LsDKYXbzOV+kjlraSO35AMec3CpISdx4K1caEAhKX6it9bvPq4pSYXi
PQspba0Jv46VV7MaabVjLzsinz5/md4vxyYHNIJAukHUfwJIsVC9ZNxukwSw+CzE
MMO64ylq2DGokNerGsLetuViV8UWi7qmUmms2fAmchodW16olgNkYTz27+V/A42S
eex63pkCgYEA+CqKhqp3qPe2E9KVrycrwjoycxmhOn3Iz1xiN7uAEv+DzfKtfZVf
mcOIiqw4Z82RkgjHb9vJuTigKdDkB1zE2gSDnep44sDWJM/5nPjGlMgnkiJWJhci
CnD0P4d6cT5wyDt7Q0/tS6ql2UrCpW4ktw1AP0Rm/z/VBD8jGkVenjcCgYEA74DM
Z2Qmh3bPt1TykpOlw+H+sEuvlkYxqMlbtn3Rv3WgEPIBekOFrgP7n/uLW1Aizn8w
EhNBBAE8w5jvklqZWYbpFMJQc09eqUkI8aTbLooZbzYj1f3CrzBRKn1GoTPmN9V0
j9r+TbH3/5CEoqlsJdmeQPofuv5Qid2oEutZcrUCgYBuZ16hco0xmqJiRzlYZvDM
w99V3X0g7Hy947e+W6gqy4nzwZb1W9LgMWE5cEzXwViVw1oWpY0k3dBDSi9oJxlc
dM2pH3sQRgH+9pdyAis2XaVdGfGBmKEITCAdc0RBxSmfqva3h4NmOlD2TpAx0MJ8
vWRrwR6hR+CYtw4CzgG+GQKBgQDGmi5lugW9JUe/xeBUrbyyv0+cT1auLUz2ouq7
XIA23Mo74wJYqW9Lyp+4nTWFJeGHDK8G/hJWyNPjeomG+jvZombbQPrHc9SSWi7h
eowKfpfywZlb1M7AyTc1HacY+9l3CTlcJQPl16NHuEZUQFue02NIjGENhd+xQy4h
ainFVQKBgAoPs9ebtWbBOaqGpOnquzb7WibvW/ifOfzv5/aOhkbpp+KCjcSON6CB
QF3BEXMcNMGWlpPrd8PaxCAzR4MyU///ekJri2icS9lrQhGSz2TtYhdED4pv1Aag
7eTPl5L7xAwphCSwy8nfCKmvlqcX/MSJ7A+LHB/2hdbuuEOyhpbu
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ga_accelerator" "default" {
  		bandwidth_billing_type = "CDT"
  		payment_type           = "PayAsYouGo"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.1.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_vswitch" "update" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.2.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_accelerator.default.id
  		name           = var.name
  		protocol       = "HTTPS"
  		port_ranges {
    		from_port = 8080
    		to_port   = 8080
  		}
  		certificates {
    		id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
  		}
	}
`, name, defaultRegionToTest)
}

func TestUnitAliCloudGaEndpointGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id":        "CreateEndpointGroupValue",
		"description":           "CreateEndpointGroupValue",
		"endpoint_group_region": "CreateEndpointGroupValue",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": true,
				"endpoint":                     "CreateEndpointGroupValue",
				"type":                         "PublicIp",
				"weight":                       20,
			},
		},
		"endpoint_group_type":           "CreateEndpointGroupValue",
		"listener_id":                   "CreateEndpointGroupValue",
		"endpoint_request_protocol":     "CreateEndpointGroupValue",
		"health_check_interval_seconds": 3,
		"health_check_path":             "CreateEndpointGroupValue",
		"health_check_port":             20,
		"health_check_protocol":         "CreateEndpointGroupValue",
		"name":                          "CreateEndpointGroupValue",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 10,
				"listener_port": 60,
			},
		},
		"threshold_count":    3,
		"traffic_percentage": 20,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeEndpointGroup
		"AcceleratorId": "CreateEndpointGroupValue",
		"Description":   "CreateEndpointGroupValue",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": true,
				"Endpoint":                   "CreateEndpointGroupValue",
				"Type":                       "PublicIp",
				"Weight":                     20,
			},
		},
		"EndpointGroupRegion":        "CreateEndpointGroupValue",
		"EndpointGroupType":          "CreateEndpointGroupValue",
		"HealthCheckIntervalSeconds": 3,
		"HealthCheckPath":            "CreateEndpointGroupValue",
		"HealthCheckPort":            20,
		"HealthCheckProtocol":        "CreateEndpointGroupValue",
		"ListenerId":                 "CreateEndpointGroupValue",
		"EndpointRequestProtocol":    "CreateEndpointGroupValue",
		"Name":                       "CreateEndpointGroupValue",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 10,
				"ListenerPort": 60,
			},
		},
		"State":             "active",
		"ThresholdCount":    3,
		"TrafficPercentage": 20,
	}
	CreateMockResponse := map[string]interface{}{
		// CreateEndpointGroup
		"EndpointGroupId": "CreateEndpointGroupValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_endpoint_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeEndpointGroup Response
		"EndpointGroupId": "CreateEndpointGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEndpointGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateEndpointGroup
	attributesDiff := map[string]interface{}{
		"description": "UpdateEndpointGroup",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": false,
				"endpoint":                     "UpdateEndpointGroup",
				"type":                         "UpdateEndpointGroup",
				"weight":                       30,
			},
		},
		"endpoint_request_protocol":     "UpdateEndpointGroup",
		"health_check_interval_seconds": 4,
		"health_check_path":             "UpdateEndpointGroup",
		"health_check_port":             30,
		"health_check_protocol":         "UpdateEndpointGroup",
		"name":                          "UpdateEndpointGroup",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 20,
				"listener_port": 70,
			},
		},
		"threshold_count":    4,
		"traffic_percentage": 30,
	}
	diff, err := newInstanceDiff("alicloud_ga_endpoint_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEndpointGroup Response
		"Description": "UpdateEndpointGroup",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": false,
				"Endpoint":                   "UpdateEndpointGroup",
				"Type":                       "UpdateEndpointGroup",
				"Weight":                     30,
			},
		},
		"HealthCheckIntervalSeconds": 4,
		"HealthCheckPath":            "UpdateEndpointGroup",
		"HealthCheckPort":            30,
		"HealthCheckProtocol":        "UpdateEndpointGroup",
		"EndpointRequestProtocol":    "UpdateEndpointGroup",
		"Name":                       "UpdateEndpointGroup",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 20,
				"ListenerPort": 70,
			},
		},
		"State":             "active",
		"ThresholdCount":    4,
		"TrafficPercentage": 30,
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEndpointGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeEndpointGroup" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEndpointGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DeleteEndpointGroup" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
