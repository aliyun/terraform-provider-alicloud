package alicloud

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alicloud_service_mesh_service_mesh",
		&resource.Sweeper{
			Name: "alicloud_alicloud_service_mesh_service_mesh",
			F:    testSweepServiceMeshServiceMesh,
		})
}

func testSweepServiceMeshServiceMesh(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-testacc",
		"test",
		"tf-example",
		"tf_example",
		"terraform-example",
	}
	action := "DescribeServiceMeshes"
	request := map[string]interface{}{}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("servicemesh", "2020-01-11", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.ServiceMeshes", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServiceMeshes", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["ServiceMeshInfo"].(map[string]interface{})["Name"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Service Mesh: %s", item["ServiceMeshInfo"].(map[string]interface{})["Name"].(string))
			continue
		}
		action := "DeleteServiceMesh"
		request := map[string]interface{}{
			"ServiceMeshId": item["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"],
		}
		_, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Service Mesh (%s): %s", item["ServiceMeshInfo"].(map[string]interface{})["Name"].(string), err)
		}
		log.Printf("[INFO] Delete Service Mesh success: %s ", item["ServiceMeshInfo"].(map[string]interface{})["Name"].(string))
	}
	return nil
}

func TestAccAliCloudServiceMeshServiceMesh_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemeshdefault%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Default",
					"version":           "${local.version_2}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							// In testing, `customized_zipkin` needs to be set to true for the following reasons:
							// 1. When creating with version 1.22.89 or higher, setting `customizedZipkin` to false will enable new trace configuration and set `customizedZipkin` to true, which is a product-side behavior.
							// 2. `tracing` controls whether it is enabled, while `customizedZipkin` controls where it is exported to; this field will gradually become deprecated.
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled": "false",
										},
									},
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${local.version_1}",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin": "false",
							"telemetry":         "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "warn",
									"request_cpu":    "2",
									"request_memory": "1024Mi",
									"limit_cpu":      "4",
									"limit_memory":   "2048Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_2}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "2048Mi",
									"request_cpu":    "200m",
									"limit_cpu":      "3000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "true",
									"request_memory":                "256Mi",
									"limit_memory":                  "2048Mi",
									"request_cpu":                   "400m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "3000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "",
										},
									},
								},
							},
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin": "true",
							"telemetry":         "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "false",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "warn",
									"request_cpu":    "2",
									"request_memory": "1024Mi",
									"limit_cpu":      "4",
									"limit_memory":   "2048Mi",
								},
							},
							//"audit": []map[string]interface{}{
							//	{
							//		"enabled": "false",
							//		"project": "",
							//	},
							//},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "2048Mi",
									"request_cpu":    "200m",
									"limit_cpu":      "3000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "true",
									"request_memory":                "256Mi",
									"limit_memory":                  "2048Mi",
									"request_cpu":                   "400m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "3000m",
								},
							},
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AlicloudServiceMeshServiceMeshMap0 = map[string]string{
	"cluster_spec": "standard",
}

func AlicloudServiceMeshServiceMeshBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_service_mesh_versions" "default" {
		edition = "Default"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_log_project" "default_1" {
		name        = "${var.name}-01"
		description = "created by terraform"
	}

	resource "alicloud_log_project" "default_2" {
		name        = "${var.name}-02"
		description = "created by terraform"
	}

	locals {
		vswitch_id    = data.alicloud_vswitches.default.ids.0
		vpc_id        = data.alicloud_vpcs.default.ids.0
		log_project_1 = alicloud_log_project.default_1.name
		log_project_2 = alicloud_log_project.default_2.name
		version_1     = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
		version_2     = reverse(data.alicloud_service_mesh_versions.default.versions).1.version
	}
`, name)
}

func TestAccAliCloudServiceMeshServiceMesh_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemeshstandard%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Default",
					"cluster_spec":      "standard",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled": "false",
										},
									},
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
					"cluster_ids": []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "standard",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids": REMOVEKEY,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudServiceMeshServiceMesh_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc-servicemesh%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Pro",
					"cluster_spec":      "enterprise",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
					"cluster_ids": []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids": REMOVEKEY,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}
func TestAccAliCloudServiceMeshServiceMesh_basic3(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemeshultimate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Pro",
					"cluster_spec":      "ultimate",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "false",
									"gateway_enabled": "false",
									"sidecar_enabled": "false",
								},
							},
						},
					},
					"cluster_ids": []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "ultimate",
						"edition":               "Pro",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids": REMOVEKEY,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}
func AlicloudServiceMeshServiceMeshBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_service_mesh_versions" "default" {
  		edition = "Default"
	}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {}
	
	resource "alicloud_vpc" "default" {
	  vpc_name   = var.name
	  cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
	  vswitch_name = var.name
	  cidr_block   = "10.4.0.0/24"
	  vpc_id       = alicloud_vpc.default.id
	  zone_id      = data.alicloud_zones.default.zones.3.id
	}
	
	resource "alicloud_cs_managed_kubernetes" "default" {
	  name_prefix          = "tf-test-service_mesh"
	  cluster_spec         = "ack.pro.small"
	  worker_vswitch_ids   = [alicloud_vswitch.default.id]
	  new_nat_gateway      = true
	  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
	  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
	  slb_internet_enabled = true
	}

	resource "alicloud_key_pair" "default" {
	  key_pair_name = var.name
	}

	data "alicloud_instance_types" "default" {
	  availability_zone    = alicloud_vswitch.default.zone_id
	  kubernetes_node_role = "Worker"
	}
	
	resource "alicloud_cs_kubernetes_node_pool" "default" {
	  name                 = "desired_size"
	  cluster_id           = alicloud_cs_managed_kubernetes.default.id
	  vswitch_ids          = [alicloud_vswitch.default.id]
	  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
	  system_disk_category = "cloud_efficiency"
	  system_disk_size     = 40
	  key_name             = alicloud_key_pair.default.key_name
	  desired_size         = 2
	}

	resource "alicloud_log_project" "default_1" {
		name        = "${var.name}-01"
		description = "created by terraform"
	}

	resource "alicloud_log_project" "default_2" {
		name        = "${var.name}-02"
		description = "created by terraform"
	}

	resource "alicloud_eip" "default" {
		bandwidth            = "10"
		internet_charge_type = "PayByBandwidth"
	}

	locals {
		vswitch_id    = alicloud_vswitch.default.id
		vpc_id        = alicloud_vpc.default.id
  		log_project_1 = alicloud_log_project.default_1.name
  		log_project_2 = alicloud_log_project.default_2.name
		version_1     = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
		version_2     = reverse(data.alicloud_service_mesh_versions.default.versions).1.version
	}
`, name)
}

func AlicloudServiceMeshServiceMeshBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_service_mesh_versions" "default" {
  		edition = "Pro"
	}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {}
	
	resource "alicloud_vpc" "default" {
	  vpc_name   = var.name
	  cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
	  vswitch_name = var.name
	  cidr_block   = "10.4.0.0/24"
	  vpc_id       = alicloud_vpc.default.id
	  zone_id      = data.alicloud_zones.default.zones.3.id
	}
	
	resource "alicloud_cs_managed_kubernetes" "default" {
	  name_prefix          = "tf-test-service_mesh"
	  cluster_spec         = "ack.pro.small"
	  worker_vswitch_ids   = [alicloud_vswitch.default.id]
	  new_nat_gateway      = true
	  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
	  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
	  slb_internet_enabled = true
	}

	resource "alicloud_key_pair" "default" {
	  key_pair_name = var.name
	}

	data "alicloud_instance_types" "default" {
	  availability_zone    = alicloud_vswitch.default.zone_id
	  kubernetes_node_role = "Worker"
	}
	
	resource "alicloud_cs_kubernetes_node_pool" "default" {
	  name                 = "desired_size"
	  cluster_id           = alicloud_cs_managed_kubernetes.default.id
	  vswitch_ids          = [alicloud_vswitch.default.id]
	  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
	  system_disk_category = "cloud_efficiency"
	  system_disk_size     = 40
	  key_name             = alicloud_key_pair.default.key_name
	  desired_size         = 2
	}

	resource "alicloud_log_project" "default_1" {
		name        = "${var.name}-01"
		description = "created by terraform"
	}

	resource "alicloud_log_project" "default_2" {
		name        = "${var.name}-02"
		description = "created by terraform"
	}

	locals {
		vswitch_id    = alicloud_vswitch.default.id
		vpc_id        = alicloud_vpc.default.id
  		log_project_1 = alicloud_log_project.default_1.name
  		log_project_2 = alicloud_log_project.default_2.name
		version_1     = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
		version_2     = reverse(data.alicloud_service_mesh_versions.default.versions).1.version
	}
`, name)
}

func TestAccAliCloudServiceMeshServiceMesh_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Default",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled": "false",
										},
									},
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"project":         "${local.log_project_1}",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
							"control_plane_log": []map[string]interface{}{
								{
									"enabled":        "true",
									"project":        "${local.log_project_1}",
									"log_ttl_in_day": "10",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "false",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
									"project": "${local.log_project_1}",
								},
							},
							"control_plane_log": []map[string]interface{}{
								{
									"enabled":        "true",
									"project":        "${local.log_project_2}",
									"log_ttl_in_day": "20",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "false",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
									"project":         "${local.log_project_1}",
								},
							},
							"control_plane_log": []map[string]interface{}{
								{
									"enabled": "false",
									"project": "${local.log_project_2}",
								},
							},
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func AlicloudServiceMeshServiceMeshBasicDependence3(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_service_mesh_versions" "default" {
		edition = "Default"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_log_project" "default_1" {
		name        = "${var.name}-01"
		description = "created by terraform"
	}

	resource "alicloud_log_project" "default_2" {
		name        = "${var.name}-02"
		description = "created by terraform"
	}

	locals {
		vswitch_id    = data.alicloud_vswitches.default.ids.0
		vpc_id        = data.alicloud_vpcs.default.ids.0
		log_project_1 = alicloud_log_project.default_1.name
		log_project_2 = alicloud_log_project.default_2.name
		version_1     = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
		version_2     = reverse(data.alicloud_service_mesh_versions.default.versions).1.version
	}
`, name)
}

func TestAccAliCloudServiceMeshServiceMesh_basic5(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemeshstandard%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Pro",
					"cluster_spec":      "enterprise",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "false",
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},

							"tracing": "false",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "0",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "false",
									"limit_cpu":                     "2000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "another,istio-system",
										},
									},
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"include_ip_ranges": "172.1.1.2/32",
						},
					},
					"cluster_ids": []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"customized_prometheus": "false",
					"prometheus_url":        "https://out.prometheus.url/",
					"force":                 "true",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "enterprise",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids":           REMOVEKEY,
					"customized_prometheus": "true",
					"prometheus_url":        "https://out.prometheus.url",
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "false",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},

							"tracing": "true",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "100",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "true",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "true",
											"exclude_namespaces": "another,istio-system,kube-system",
										},
									},
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":           "true",
									"project":           "${local.log_project_1}",
									"gateway_enabled":   "true",
									"gateway_lifecycle": "3",
									"sidecar_enabled":   "true",
									"sidecar_lifecycle": "3",
								},
							},
							"include_ip_ranges": "*",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_configuration.#": "1",
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
				ImportStateVerifyIgnore: []string{"force", "customized_prometheus", "prometheus_url"},
			},
		},
	})
}

func TestAccAliCloudServiceMeshServiceMesh_basic6(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemeshstandard%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence6)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"edition": "Default",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"network.#":         "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "customized_prometheus", "prometheus_url"},
			},
		},
	})
}

func AlicloudServiceMeshServiceMeshBasicDependence6(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_service_mesh_versions" "default" {
		edition = "Default"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_log_project" "default_1" {
		name        = "${var.name}-01"
		description = "created by terraform"
	}

	resource "alicloud_log_project" "default_2" {
		name        = "${var.name}-02"
		description = "created by terraform"
	}

	locals {
		vswitch_id    = data.alicloud_vswitches.default.ids.0
		vpc_id        = data.alicloud_vpcs.default.ids.0
		log_project_1 = alicloud_log_project.default_1.name
		log_project_2 = alicloud_log_project.default_2.name
		version_1     = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
		version_2     = reverse(data.alicloud_service_mesh_versions.default.versions).1.version
	}
`, name)
}

func TestAccAliCloudServiceMeshServiceMesh_basic7(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc-servicemesh%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Pro",
					"cluster_spec":      "enterprise",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"tracing":            "false",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "0",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"cluster_domain": "cluster.local",
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
					"customized_prometheus": "false",
					"cluster_ids":           []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"telemetry": "true",
							"kiali": []map[string]interface{}{
								{
									"enabled":                   "true",
									"integrate_clb":             "true",
									"kiali_service_annotations": "{\\\"${alicloud_cs_kubernetes_node_pool.default.cluster_id}\\\":{\\\"service.beta.kubernetes.io/alicloud-loadbalancer-address-type\\\":\\\"intranet\\\"}}",
									"kiali_arms_auth_tokens":    "{\\\"${alicloud_cs_kubernetes_node_pool.default.cluster_id}\\\":\\\"abcdefg\\\"}",
									"auth_strategy":             "ramoauth",
									"ram_oauth_config": []map[string]interface{}{
										{
											"redirect_uris": "http://www.terraformtest.com",
										},
									},
								},
							},
						},
					},
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"telemetry": "true",
							"kiali": []map[string]interface{}{
								{
									"enabled":               "true",
									"auth_strategy":         "openid",
									"custom_prometheus_url": "https://out.prometheus.url",
									"open_id_config": []map[string]interface{}{
										{
											"client_id":     "testid",
											"client_secret": "testsecret",
											"issuer_uri":    "www.terraformtest.com",
											"scopes":        []string{"openid"},
										},
									},
									"server_config": []map[string]interface{}{
										{
											"web_fqdn":   "www.terraformtest.com",
											"web_root":   "/",
											"web_port":   "80",
											"web_schema": "http",
										},
									},
								},
							},
						},
					},
					"customized_prometheus": "true",
					"prometheus_url":        "https://out.prometheus.url",
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name,
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids": REMOVEKEY,
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "customized_prometheus", "prometheus_url"},
			},
		},
	})
}

func TestAccAliCloudServiceMeshServiceMesh_basic8(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc-servicemesh%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Pro",
					"cluster_spec":      "enterprise",
					"version":           "${local.version_1}",
					"network": []map[string]interface{}{
						{
							"vpc_id":        "${local.vpc_id}",
							"vswitche_list": []string{"${local.vswitch_id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"customized_zipkin":  "true",
							"enable_locality_lb": "false",
							"telemetry":          "true",
							"tracing":            "false",
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "true",
									"trace_sampling": "0",
								},
							},
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"log_level":      "info",
									"request_cpu":    "1",
									"request_memory": "512Mi",
									"limit_cpu":      "2",
									"limit_memory":   "1024Mi",
								},
							},
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"audit": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_2}",
								},
							},
							"proxy": []map[string]interface{}{
								{
									"cluster_domain": "cluster.local",
									"request_memory": "128Mi",
									"limit_memory":   "1024Mi",
									"request_cpu":    "100m",
									"limit_cpu":      "2000m",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"request_memory":                "128Mi",
									"limit_memory":                  "1024Mi",
									"request_cpu":                   "100m",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "2000m",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"access_log": []map[string]interface{}{
								{
									"enabled":         "true",
									"gateway_enabled": "true",
									"sidecar_enabled": "true",
								},
							},
						},
					},
					"customized_prometheus": "false",
					"cluster_ids":           []string{"${alicloud_cs_kubernetes_node_pool.default.cluster_id}"},
					"extra_configuration": []map[string]interface{}{
						{
							"cr_aggregation_enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"mesh_config.#":         "1",
						"cluster_ids.#":         "1",
						"extra_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}-m",
					"cluster_spec":      "ultimate",
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip_id":   "${alicloud_eip.default.id}",
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
									"project": "",
								},
							},
						},
					},
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name + "-m",
						"cluster_spec":      "ultimate",
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}-m",
					"cluster_spec":      "ultimate",
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip_id":   "",
							"api_server_public_eip": "false",
						},
					},
					"mesh_config": []map[string]interface{}{
						{
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
									"project": "",
								},
							},
						},
					},
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name + "-m",
						"cluster_spec":      "ultimate",
						"mesh_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_ids": REMOVEKEY,
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name": name + "-m",
						"cluster_spec":      "ultimate",
						"cluster_ids.#":     "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "customized_prometheus", "prometheus_url"},
			},
		},
	})
}

type vpcASMAPI struct {
	region           string
	list             func() (map[string]interface{}, error)
	detail           func(string) (map[string]interface{}, error)
	clusterDetail    func(string) (map[string]interface{}, error)
	clusterWasScoped func(string) bool
	remove           func(string) (map[string]interface{}, error)
	now              func() time.Time
	sleep            func(time.Duration)
	maxPolls         int
	releaseTimeout   time.Duration
}

func nativeVpcASMAPI(client *connectivity.AliyunClient, region string) vpcASMAPI {
	ack := nativeVpcACKAPI(client, region)
	return vpcASMAPI{region: region,
		list: func() (map[string]interface{}, error) {
			return client.RpcGet("servicemesh", "2020-01-11", "DescribeServiceMeshes", map[string]interface{}{}, nil)
		},
		detail: func(id string) (map[string]interface{}, error) {
			return client.RpcPost("servicemesh", "2020-01-11", "DescribeServiceMeshDetail", nil, map[string]interface{}{"ServiceMeshId": id}, true)
		},
		clusterDetail: ack.detail,
		// Scope and associated-cluster checks run before this delete adapter.
		remove: func(id string) (map[string]interface{}, error) {
			return nativeVpcASMDelete(client.RpcPost, id)
		},
		now: time.Now, sleep: time.Sleep, maxPolls: 180, releaseTimeout: 30 * time.Minute}
}

// Force is the native ASM parameter, used only by explicitly authorized full
// cleanup after VPC/cluster ownership verification. Ordinary deletion is unchanged.
func nativeVpcASMDelete(rpc vpcSweepRPC, id string) (map[string]interface{}, error) {
	request := map[string]interface{}{"ServiceMeshId": id}
	if sweepAll() {
		request["Force"] = true
	}
	return rpc("servicemesh", "2020-01-11", "DeleteServiceMesh", nil, request, true)
}

// A regional not-found alone cannot prove that a cluster in another region
// is gone. Retain only positive, exact ownership observed by this invocation's
// ACK sweeper so its completed deletion can safely precede ASM cleanup.
func trackVpcACKOwnership(region, vpcID string, api vpcACKAPI) (vpcACKAPI, func(string) bool) {
	var mu sync.RWMutex
	verified := map[string]bool{}
	proof := func(id string) bool {
		mu.RLock()
		defer mu.RUnlock()
		return verified[id]
	}
	original := api.detail
	if original == nil {
		return api, proof
	}
	api.detail = func(id string) (map[string]interface{}, error) {
		row, err := original(id)
		if err == nil && region != "" && vpcID != "" && id != "" && row["cluster_id"] == id && row["region_id"] == region && row["vpc_id"] == vpcID {
			mu.Lock()
			verified[id] = true
			mu.Unlock()
		}
		return row, err
	}
	return api, proof
}

func asmExactCode(err error, code string) bool {
	for {
		if e, ok := err.(*ComplexError); ok {
			err = e.Cause
			continue
		}
		break
	}
	switch e := err.(type) {
	case *tea.SDKError:
		return tea.StringValue(e.Code) == code
	case *sdkerrors.ServerError:
		return e.ErrorCode() == code
	}
	return false
}
func asmList(api vpcASMAPI) ([]map[string]interface{}, error) {
	r, e := api.list()
	if e != nil {
		return nil, e
	}
	for _, k := range []string{"NextToken", "nextToken", "HasMore", "PageNumber", "PageSize"} {
		if v, ok := r[k]; ok && v != nil {
			return nil, fmt.Errorf("ASM unexpected pagination %s", k)
		}
	}
	raw, ok := r["ServiceMeshes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("ASM missing complete ServiceMeshes")
	}
	rows := []map[string]interface{}{}
	seen := map[string]bool{}
	for _, v := range raw {
		row, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ASM malformed list")
		}
		info, ok := row["ServiceMeshInfo"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ASM missing identity")
		}
		id, ok := info["ServiceMeshId"].(string)
		if !ok || id == "" || seen[id] {
			return nil, fmt.Errorf("ASM invalid/duplicate identity")
		}
		seen[id] = true
		rows = append(rows, row)
	}
	return rows, nil
}
func asmScope(row map[string]interface{}, region, vpc, id string) (bool, error) {
	info, ok := row["ServiceMeshInfo"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing info")
	}
	if info["ServiceMeshId"] != id {
		return false, fmt.Errorf("ASM identity changed")
	}
	spec, ok := row["Spec"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing spec")
	}
	network, ok := spec["Network"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing network")
	}
	owner, ok := network["VpcId"].(string)
	if !ok || owner == "" {
		return false, fmt.Errorf("ASM unknown VPC")
	}
	if owner != vpc {
		return false, nil
	}
	if info["RegionId"] != region {
		return false, fmt.Errorf("ASM region mismatch")
	}
	return true, nil
}
func asmClusterScope(row map[string]interface{}, api vpcASMAPI, region, vpc string) error {
	raw, ok := row["Clusters"].([]interface{})
	if !ok {
		return fmt.Errorf("ASM complete associated cluster list missing")
	}
	seen := map[string]bool{}
	for _, v := range raw {
		id, ok := v.(string)
		if !ok || id == "" || seen[id] {
			return fmt.Errorf("ASM malformed associated cluster")
		}
		seen[id] = true
		d, e := api.clusterDetail(id)
		if e != nil {
			if asmExactCode(e, "ErrorClusterNotFound") && api.clusterWasScoped != nil && api.clusterWasScoped(id) {
				continue
			}
			return fmt.Errorf("ASM associated ACK detail: %w", e)
		}
		if d["cluster_id"] != id || d["region_id"] != region || d["vpc_id"] != vpc {
			return fmt.Errorf("ASM shared/unknown associated ACK ownership; retained")
		}
	}
	return nil
}
func vpcASMSweeper(region, vpcID string, api vpcASMAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_service_mesh_service_mesh", F: func(actual string) error {
		if actual != region || region == "" || vpcID == "" || api.region != region || api.list == nil || api.detail == nil || api.clusterDetail == nil || api.remove == nil || api.now == nil || api.sleep == nil || api.maxPolls < 1 || api.releaseTimeout <= 0 {
			return fmt.Errorf("ASM requires exact scope and bounded API")
		}
		rows, e := asmList(api)
		if e != nil {
			return e
		}
		var failures []string
		for _, row := range rows {
			id := row["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"].(string)
			selected, e := asmScope(row, region, vpcID, id)
			if e != nil {
				failures = append(failures, e.Error())
				continue
			}
			if !selected {
				continue
			}
			r, e := api.detail(id)
			if e != nil {
				failures = append(failures, e.Error())
				continue
			}
			detail, ok := r["ServiceMesh"].(map[string]interface{})
			if !ok {
				failures = append(failures, "ASM missing fresh detail")
				continue
			}
			selected, e = asmScope(detail, region, vpcID, id)
			if e != nil || !selected {
				failures = append(failures, "ASM fresh ownership changed")
				continue
			}
			if e = asmClusterScope(detail, api, region, vpcID); e != nil {
				failures = append(failures, e.Error())
				continue
			}
			if _, e = api.remove(id); e != nil {
				failures = append(failures, e.Error())
				continue
			}
			deadline := api.now().Add(api.releaseTimeout)
			gone := false
			for poll := 0; poll < api.maxPolls && api.now().Before(deadline); poll++ {
				_, e = api.detail(id)
				if e != nil {
					if !asmExactCode(e, "ServiceMesh.NotFound") {
						break
					}
					var current []map[string]interface{}
					current, e = asmList(api)
					if e != nil {
						break
					}
					found := false
					for _, x := range current {
						if x["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"] == id {
							found = true
						}
					}
					if !found {
						gone = true
						break
					}
				}
				api.sleep(10 * time.Second)
			}
			if !gone {
				failures = append(failures, fmt.Sprintf("ASM release not confirmed: %v", e))
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("ASM owner cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}

func asmUnitMesh(clusters interface{}) map[string]interface{} {
	return map[string]interface{}{"ServiceMeshInfo": map[string]interface{}{"ServiceMeshId": "mesh", "RegionId": "region"}, "Spec": map[string]interface{}{"Network": map[string]interface{}{"VpcId": "target"}}, "Clusters": clusters}
}
func asmUnitAPI(mesh map[string]interface{}) (vpcASMAPI, *int) {
	deletes := 0
	gone := false
	now := time.Unix(0, 0)
	api := vpcASMAPI{region: "region", now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }, maxPolls: 3, releaseTimeout: time.Minute}
	api.list = func() (map[string]interface{}, error) {
		rows := []interface{}{}
		if !gone {
			rows = append(rows, mesh)
		}
		return map[string]interface{}{"ServiceMeshes": rows}, nil
	}
	api.detail = func(string) (map[string]interface{}, error) {
		if gone {
			return nil, &tea.SDKError{Code: tea.String("ServiceMesh.NotFound")}
		}
		return map[string]interface{}{"ServiceMesh": mesh}, nil
	}
	api.clusterDetail = func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}, nil
	}
	api.remove = func(string) (map[string]interface{}, error) {
		deletes++
		gone = true
		return map[string]interface{}{}, nil
	}
	return api, &deletes
}
func TestUnitVpcASMScopedOwnerAndExactAsyncConfirmation(t *testing.T) {
	for _, knownGone := range []bool{false, true} {
		t.Run(fmt.Sprint(knownGone), func(t *testing.T) {
			api, n := asmUnitAPI(asmUnitMesh([]interface{}{"cluster"}))
			if knownGone {
				api.clusterWasScoped = func(id string) bool { return id == "cluster" }
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
			}
			if e := vpcASMSweeper("region", "target", api).F("region"); e != nil || *n != 1 {
				t.Fatalf("delete=%d err=%v", *n, e)
			}
		})
	}
}
func TestUnitVpcASMSharedUnknownAndOpaqueErrorsNeverDelete(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"shared", "region", "identity", "missingClusters", "malformedClusters", "permission", "opaqueNotFound", "regionalNotFound", "freshVpc", "freshError"} {
		t.Run(mode, func(t *testing.T) {
			mesh := asmUnitMesh([]interface{}{"cluster"})
			api, n := asmUnitAPI(mesh)
			switch mode {
			case "shared", "region", "identity":
				api.clusterDetail = func(id string) (map[string]interface{}, error) {
					d := map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}
					if mode == "shared" {
						d["vpc_id"] = "protected"
					}
					if mode == "region" {
						d["region_id"] = "other"
					}
					if mode == "identity" {
						d["cluster_id"] = "other"
					}
					return d, nil
				}
			case "missingClusters":
				delete(mesh, "Clusters")
			case "malformedClusters":
				mesh["Clusters"] = []interface{}{1}
			case "regionalNotFound":
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
			case "permission", "opaqueNotFound":
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("Forbidden"), Message: tea.String("cluster not found or inaccessible")}
				}
			case "freshVpc":
				api.detail = func(string) (map[string]interface{}, error) {
					d := asmUnitMesh([]interface{}{})
					d["Spec"].(map[string]interface{})["Network"].(map[string]interface{})["VpcId"] = "protected"
					return map[string]interface{}{"ServiceMesh": d}, nil
				}
			case "freshError":
				api.detail = func(string) (map[string]interface{}, error) { return nil, fmt.Errorf("MiranaErr.Common") }
			}
			if e := vpcASMSweeper("region", "target", api).F("region"); e == nil || *n != 0 {
				t.Fatalf("delete=%d err=%v", *n, e)
			}
		})
	}
}
func TestUnitVpcASMCompleteListAndBoundedRelease(t *testing.T) {
	for _, mode := range []string{"missingList", "pagination", "acceptedStillPresent", "opaqueGone", "listStillPresent", "emptyClusters"} {
		t.Run(mode, func(t *testing.T) {
			api, n := asmUnitAPI(asmUnitMesh([]interface{}{}))
			switch mode {
			case "missingList":
				api.list = func() (map[string]interface{}, error) { return map[string]interface{}{}, nil }
			case "pagination":
				orig := api.list
				api.list = func() (map[string]interface{}, error) { d, e := orig(); d["NextToken"] = "more"; return d, e }
			case "acceptedStillPresent":
				api.remove = func(string) (map[string]interface{}, error) { *n++; return map[string]interface{}{}, nil }
			case "opaqueGone":
				api.remove = func(string) (map[string]interface{}, error) {
					*n++
					api.detail = func(string) (map[string]interface{}, error) { return nil, fmt.Errorf("not found or inaccessible") }
					return map[string]interface{}{}, nil
				}
			case "listStillPresent":
				api.remove = func(string) (map[string]interface{}, error) {
					*n++
					api.detail = func(string) (map[string]interface{}, error) {
						return nil, &tea.SDKError{Code: tea.String("ServiceMesh.NotFound")}
					}
					return map[string]interface{}{}, nil
				}
			}
			e := vpcASMSweeper("region", "target", api).F("region")
			if mode == "emptyClusters" {
				if e != nil || *n != 1 {
					t.Fatalf("delete=%d err=%v", *n, e)
				}
			} else if e == nil {
				t.Fatal("invalid/incomplete release accepted")
			}
		})
	}
}

func TestUnitVpcASMNativeForceOnlyForFullCleanup(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprint(full), func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", fmt.Sprint(full))
			calls := 0
			expected := fmt.Errorf("native ASM delete rejected")
			_, err := nativeVpcASMDelete(func(product, version, action string, query, request map[string]interface{}, auto bool) (map[string]interface{}, error) {
				calls++
				force, present := request["Force"]
				fields := 1
				if full {
					fields++
				}
				if product != "servicemesh" || version != "2020-01-11" || action != "DeleteServiceMesh" || query != nil || request["ServiceMeshId"] != "mesh-target" || present != full || (full && force != true) || len(request) != fields {
					t.Fatalf("wrong native full/ordinary ASM delete contract: %s %s %s %v", product, version, action, request)
				}
				return nil, expected
			}, "mesh-target")
			if err != expected || calls != 1 {
				t.Fatalf("native deletion result was lost: err=%v calls=%d", err, calls)
			}
		})
	}
}

func TestUnitVpcASMACKOwnershipProofIsLocalAndExact(t *testing.T) {
	for _, scenario := range []string{"owned", "other-vpc", "other-region", "other-id", "read-error", "missing-detail"} {
		t.Run(scenario, func(t *testing.T) {
			api := vpcACKAPI{detail: func(id string) (map[string]interface{}, error) {
				row := map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}
				switch scenario {
				case "other-vpc":
					row["vpc_id"] = "protected"
				case "other-region":
					row["region_id"] = "reserved-region"
				case "other-id":
					row["cluster_id"] = "different-cluster"
				case "read-error":
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
				return row, nil
			}}
			if scenario == "missing-detail" {
				api.detail = nil
			}
			tracked, proof := trackVpcACKOwnership("region", "target", api)
			if proof("cluster") {
				t.Fatal("ownership existed before a successful exact read")
			}
			if tracked.detail != nil {
				_, _ = tracked.detail("cluster")
			}
			if proof("cluster") != (scenario == "owned") || proof("another-cluster") {
				t.Fatal("unknown, protected or cross-region identity counted as owned")
			}
			_, freshProof := trackVpcACKOwnership("region", "target", api)
			if freshProof("cluster") {
				t.Fatal("ownership leaked into another invocation")
			}
		})
	}
}
