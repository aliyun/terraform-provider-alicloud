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
	}
	action := "DescribeServiceMeshes"
	request := map[string]interface{}{}

	var response map[string]interface{}
	conn, err := client.NewServicemeshClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-01-11"), StringPointer("AK"), request, nil, &runtime)
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
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
							"customized_zipkin":  "false",
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

// Test ServiceMesh ServiceMesh. >>> Resource test cases, automatically generated.
// Case Terraform_online_new 4755
func TestAccAliCloudServiceMeshServiceMesh_basic4755(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap4755)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceMeshServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicemeshservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence4755)
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
					"service_mesh_name": name,
					"version":           "${var.istio-version}",
					"cluster_spec":      "enterprise",
					"mesh_config": []map[string]interface{}{
						{
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"request_memory": "256Mi",
									"log_level":      "info",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "false",
									"trace_sampling": "0",
								},
							},
							"customized_zipkin": "false",
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"auto_injection_policy_enabled": "true",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "excluded,istio-system,kube-system",
										},
									},
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"request_memory": "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"telemetry": "false",
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
									"url":     "https://invalid-url",
								},
							},
							"tracing":            "false",
							"enable_locality_lb": "false",
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"include_ip_ranges": "*",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
						},
					},
					"network": []map[string]interface{}{
						{
							"vpc_id": "${var.vpc-id}",
							"vswitche_list": []string{
								"${var.vsw-id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
					"edition":               "Pro",
					"customized_prometheus": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"version":               CHECKSET,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"customized_prometheus": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"pilot": []map[string]interface{}{
								{
									"trace_sampling": "50",
									"http10_enabled": "true",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"limit_memory":                  "512Mi",
									"request_cpu":                   "512Mi",
									"request_memory":                "512Mi",
									"enable_namespaces_by_default":  "true",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "512Mi",
								},
							},
							"include_ip_ranges": "172.1.1.2/32",
							"opa": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"log_level":      "warn",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
									"enabled":        "true",
								},
							},
							"telemetry":               "true",
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
									"url":     "https://some-kiali.url",
								},
							},
							"tracing":            "true",
							"enable_locality_lb": "true",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"cluster_domain": "512Mi",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
								},
							},
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"customized_zipkin": "true",
						},
					},
					"customized_prometheus": "true",
					"prometheus_url":        "https://out.prometheus.url",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customized_prometheus": "true",
						"prometheus_url":        "https://out.prometheus.url",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"sidecar_injector": []map[string]interface{}{
								{
									"init_cni_configuration": []map[string]interface{}{
										{
											"exclude_namespaces": "another,istio-system,kube-system",
											"enabled":            "true",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${var.istio-version-2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${var.istio-version}",
					"mesh_config": []map[string]interface{}{
						{
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"request_memory": "256Mi",
									"log_level":      "info",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "false",
									"trace_sampling": "0",
								},
							},
							"customized_zipkin": "false",
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"auto_injection_policy_enabled": "true",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "excluded,istio-system,kube-system",
										},
									},
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"request_memory": "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"telemetry": "false",
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
									"url":     "https://invalid-url",
								},
							},
							"tracing":            "false",
							"enable_locality_lb": "false",
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"include_ip_ranges": "*",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
						},
					},
					"customized_prometheus": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":               CHECKSET,
						"customized_prometheus": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"pilot": []map[string]interface{}{
								{
									"trace_sampling": "50",
									"http10_enabled": "true",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"limit_memory":                  "512Mi",
									"request_cpu":                   "512Mi",
									"request_memory":                "512Mi",
									"enable_namespaces_by_default":  "true",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "512Mi",
								},
							},
							"include_ip_ranges": "172.1.1.2/32",
							"opa": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"log_level":      "warn",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
									"enabled":        "true",
								},
							},
							"telemetry":               "true",
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
									"url":     "https://some-kiali.url",
								},
							},
							"tracing":            "true",
							"enable_locality_lb": "true",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"cluster_domain": "512Mi",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
								},
							},
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"customized_zipkin": "true",
						},
					},
					"customized_prometheus": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customized_prometheus": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"customized_prometheus", "prometheus_url"},
			},
		},
	})
}

var AlicloudServiceMeshServiceMeshMap4755 = map[string]string{
	"mesh_config.#":         CHECKSET,
	"kubeconfig":            CHECKSET,
	"status":                CHECKSET,
	"create_time":           CHECKSET,
	"load_balancer.#":       CHECKSET,
	"cluster_spec":          CHECKSET,
	"extra_configuration.#": CHECKSET,
}

func AlicloudServiceMeshServiceMeshBasicDependence4755(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster-1" {
  default = "c9962195a41814399ba2b022283a9f2b2"
}

variable "cluster-2" {
  default = "c503428fdd3c6400585c35c4eaa1a6dba"
}

variable "vsw-id" {
  default = "vsw-bp1k8z9x1z9m5bx4p7tcc"
}

variable "istio-version" {
  default = "v1.17.2.27-g854e3b73-aliyun"
}

variable "vpc-id" {
  default = "vpc-bp1imspdu1ikr5j4jslkv"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "cluster-3" {
  default = "cb926fe87f31d4fefa577dbf4108202ae"
}

variable "istio-version-2" {
  default = "v1.18.0.117-g0febec19-aliyun("
}


`, name)
}

// Case Terraform_online 4586
func TestAccAliCloudServiceMeshServiceMesh_basic4586(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap4586)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceMeshServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicemeshservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence4586)
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
					"service_mesh_name": name,
					"version":           "${var.istio-version}",
					"cluster_spec":      "enterprise",
					"mesh_config": []map[string]interface{}{
						{
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"request_memory": "256Mi",
									"log_level":      "info",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "false",
									"trace_sampling": "0",
								},
							},
							"customized_zipkin": "false",
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"auto_injection_policy_enabled": "true",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "excluded,istio-system,kube-system",
										},
									},
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"request_memory": "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"telemetry": "false",
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
									"url":     "https://invalid-url",
								},
							},
							"tracing":            "false",
							"enable_locality_lb": "false",
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"include_ip_ranges": "*",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
						},
					},
					"network": []map[string]interface{}{
						{
							"vpc_id": "${var.vpc-id}",
							"vswitche_list": []string{
								"${var.vsw-id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
					"edition":               "Pro",
					"customized_prometheus": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"version":               CHECKSET,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"customized_prometheus": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"pilot": []map[string]interface{}{
								{
									"trace_sampling": "50",
									"http10_enabled": "true",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"limit_memory":                  "512Mi",
									"request_cpu":                   "512Mi",
									"request_memory":                "512Mi",
									"enable_namespaces_by_default":  "true",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "512Mi",
								},
							},
							"include_ip_ranges": "172.1.1.2/32",
							"opa": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"log_level":      "warn",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
									"enabled":        "true",
								},
							},
							"telemetry":               "true",
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
									"url":     "https://some-kiali.url",
								},
							},
							"tracing":            "true",
							"enable_locality_lb": "true",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"cluster_domain": "512Mi",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
								},
							},
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"customized_zipkin": "true",
						},
					},
					"customized_prometheus": "true",
					"prometheus_url":        "https://out.prometheus.url",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customized_prometheus": "true",
						"prometheus_url":        "https://out.prometheus.url",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"sidecar_injector": []map[string]interface{}{
								{
									"init_cni_configuration": []map[string]interface{}{
										{
											"exclude_namespaces": "another,istio-system,kube-system",
											"enabled":            "true",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${var.istio-version-2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"customized_prometheus", "prometheus_url"},
			},
		},
	})
}

var AlicloudServiceMeshServiceMeshMap4586 = map[string]string{
	"mesh_config.#":         CHECKSET,
	"kubeconfig":            CHECKSET,
	"status":                CHECKSET,
	"create_time":           CHECKSET,
	"load_balancer.#":       CHECKSET,
	"cluster_spec":          CHECKSET,
	"extra_configuration.#": CHECKSET,
}

func AlicloudServiceMeshServiceMeshBasicDependence4586(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster-1" {
  default = "c812c50970fb34e92b1b557a33e205e5d"
}

variable "cluster-2" {
  default = "ced1f5de824f74967b472813616916eb5"
}

variable "vsw-id" {
  default = "vsw-bp1k8z9x1z9m5bx4p7tcc"
}

variable "istio-version" {
  default = "v1.17.2.27-g854e3b73-aliyun"
}

variable "vpc-id" {
  default = "vpc-bp1imspdu1ikr5j4jslkv"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "cluster-3" {
  default = "cb926fe87f31d4fefa577dbf4108202ae"
}

variable "istio-version-2" {
  default = "v1.18.0.114-g105fcfd0-aliyun"
}


`, name)
}

// Case Tag_pre 4521
func TestAccAliCloudServiceMeshServiceMesh_basic4521(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap4521)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceMeshServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicemeshservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence4521)
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
					"mesh_config": []map[string]interface{}{
						{
							"audit": []map[string]interface{}{
								{
									"project": "mesh-log-{meshId}",
								},
							},
							"telemetry":          "false",
							"tracing":            "false",
							"enable_locality_lb": "false",
						},
					},
					"network": []map[string]interface{}{
						{
							"vpc_id": "${var.vpc-id}",
							"vswitche_list": []string{
								"${var.vsw-id}"},
						},
					},
					"version":           "v1.17.2.40-gca3a4011-aliyun",
					"service_mesh_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":           "v1.17.2.40-gca3a4011-aliyun",
						"service_mesh_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"audit": []map[string]interface{}{
								{
									"project": "mesh-log-{meshId}",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"audit": []map[string]interface{}{
								{
									"project": "mesh-log-{meshId}",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"customized_prometheus", "prometheus_url"},
			},
		},
	})
}

var AlicloudServiceMeshServiceMeshMap4521 = map[string]string{
	"mesh_config.#":         CHECKSET,
	"kubeconfig":            CHECKSET,
	"status":                CHECKSET,
	"create_time":           CHECKSET,
	"load_balancer.#":       CHECKSET,
	"cluster_spec":          CHECKSET,
	"extra_configuration.#": CHECKSET,
}

func AlicloudServiceMeshServiceMeshBasicDependence4521(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw-id" {
  default = "vsw-bp1k8z9x1z9m5bx4p7tcc"
}

variable "vpc-id" {
  default = "vpc-bp1imspdu1ikr5j4jslkv"
}


`, name)
}

// Case Terraform_pre 4560
func TestAccAliCloudServiceMeshServiceMesh_basic4560(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_service_mesh.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap4560)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceMeshServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshServiceMesh")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicemeshservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence4560)
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
					"service_mesh_name": name,
					"version":           "${var.istio-version}",
					"cluster_spec":      "enterprise",
					"mesh_config": []map[string]interface{}{
						{
							"opa": []map[string]interface{}{
								{
									"enabled":        "false",
									"request_memory": "256Mi",
									"log_level":      "info",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"pilot": []map[string]interface{}{
								{
									"http10_enabled": "false",
									"trace_sampling": "0",
								},
							},
							"customized_zipkin": "false",
							"sidecar_injector": []map[string]interface{}{
								{
									"enable_namespaces_by_default":  "false",
									"auto_injection_policy_enabled": "true",
									"init_cni_configuration": []map[string]interface{}{
										{
											"enabled":            "false",
											"exclude_namespaces": "excluded,istio-system,kube-system",
										},
									},
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"request_memory": "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
							"telemetry": "false",
							"kiali": []map[string]interface{}{
								{
									"enabled": "false",
									"url":     "https://invalid-url",
								},
							},
							"tracing":            "false",
							"enable_locality_lb": "false",
							"audit": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"outbound_traffic_policy": "ALLOW_ANY",
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"include_ip_ranges": "*",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "256Mi",
									"limit_memory":   "256Mi",
									"request_cpu":    "256Mi",
									"limit_cpu":      "256Mi",
								},
							},
						},
					},
					"network": []map[string]interface{}{
						{
							"vpc_id": "${var.vpc-id}",
							"vswitche_list": []string{
								"${var.vsw-id}"},
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"pilot_public_eip":      "false",
							"api_server_public_eip": "false",
						},
					},
					"edition":               "Pro",
					"customized_prometheus": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_name":     name,
						"version":               CHECKSET,
						"cluster_spec":          "enterprise",
						"edition":               "Pro",
						"customized_prometheus": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"pilot": []map[string]interface{}{
								{
									"trace_sampling": "50",
									"http10_enabled": "true",
								},
							},
							"sidecar_injector": []map[string]interface{}{
								{
									"limit_memory":                  "512Mi",
									"request_cpu":                   "512Mi",
									"request_memory":                "512Mi",
									"enable_namespaces_by_default":  "true",
									"auto_injection_policy_enabled": "true",
									"limit_cpu":                     "512Mi",
								},
							},
							"include_ip_ranges": "172.1.1.2/32",
							"opa": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"log_level":      "warn",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
									"enabled":        "true",
								},
							},
							"telemetry":               "true",
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"kiali": []map[string]interface{}{
								{
									"enabled": "true",
									"url":     "https://some-kiali.url",
								},
							},
							"tracing":            "true",
							"enable_locality_lb": "true",
							"proxy": []map[string]interface{}{
								{
									"request_memory": "512Mi",
									"cluster_domain": "512Mi",
									"limit_memory":   "512Mi",
									"request_cpu":    "512Mi",
									"limit_cpu":      "512Mi",
								},
							},
							"prometheus": []map[string]interface{}{
								{
									"use_external": "true",
								},
							},
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"customized_zipkin": "true",
						},
					},
					"customized_prometheus": "true",
					"prometheus_url":        "https://out.prometheus.url",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"customized_prometheus": "true",
						"prometheus_url":        "https://out.prometheus.url",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mesh_config": []map[string]interface{}{
						{
							"sidecar_injector": []map[string]interface{}{
								{
									"init_cni_configuration": []map[string]interface{}{
										{
											"exclude_namespaces": "another,istio-system,kube-system",
											"enabled":            "true",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${var.istio-version-2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"customized_prometheus", "prometheus_url"},
			},
		},
	})
}

var AlicloudServiceMeshServiceMeshMap4560 = map[string]string{
	"mesh_config.#":         CHECKSET,
	"kubeconfig":            CHECKSET,
	"status":                CHECKSET,
	"create_time":           CHECKSET,
	"load_balancer.#":       CHECKSET,
	"cluster_spec":          CHECKSET,
	"extra_configuration.#": CHECKSET,
}

func AlicloudServiceMeshServiceMeshBasicDependence4560(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "cluster-1" {
  default = "c812c50970fb34e92b1b557a33e205e5d"
}

variable "cluster-2" {
  default = "ced1f5de824f74967b472813616916eb5"
}

variable "vsw-id" {
  default = "vsw-bp1k8z9x1z9m5bx4p7tcc"
}

variable "istio-version" {
  default = "v1.17.2.27-g854e3b73-aliyun"
}

variable "vpc-id" {
  default = "vpc-bp1imspdu1ikr5j4jslkv"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "cluster-3" {
  default = "cb926fe87f31d4fefa577dbf4108202ae"
}

variable "istio-version-2" {
  default = "v1.18.0.114-g105fcfd0-aliyun"
}


`, name)
}

// Test ServiceMesh ServiceMesh. <<< Resource test cases, automatically generated.
