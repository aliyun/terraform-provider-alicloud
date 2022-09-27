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

func TestAccAlicloudServiceMeshServiceMesh_basic0(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshServiceMeshBasicDependence0)
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
					"edition":           "Default",
					"version":           "${data.alicloud_service_mesh_versions.default.versions.0.version}",
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
									"enabled": "true",
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
									"trace_sampling": "80",
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
								},
							},
							"outbound_traffic_policy": "REGISTRY_ONLY",
							"access_log": []map[string]interface{}{
								{
									"enabled": "true",
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

var AlicloudServiceMeshServiceMeshMap0 = map[string]string{}

func AlicloudServiceMeshServiceMeshBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_service_mesh_versions" "default" {
	edition = "Default"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id     	= data.alicloud_zones.default.zones.0.id
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
  vswitch_id = data.alicloud_vswitches.default.ids.0
  vpc_id = data.alicloud_vpcs.default.ids.0
  log_project_1 = alicloud_log_project.default_1.name
  log_project_2 = alicloud_log_project.default_2.name
}

`, name)
}

func TestAccAlicloudServiceMeshServiceMesh_basic1(t *testing.T) {
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
					"edition":           "Default",
					"cluster_spec":      "standard",
					"version":           "${data.alicloud_service_mesh_versions.default.versions.0.version}",
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
									"enabled": "true",
								},
							},
						},
					},
					"cluster_ids": []string{"${data.alicloud_cs_managed_kubernetes_clusters.default.clusters[0].id}"},
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

func AlicloudServiceMeshServiceMeshBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_service_mesh_versions" "default" {
	edition = "Default"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
  tags = {
    For = "Terraform-CI-Test"
  }
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id     	= data.alicloud_zones.default.zones.0.id
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "default-NODELETING"
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
  vswitch_id = data.alicloud_vswitches.default.ids.0
  vpc_id = data.alicloud_vpcs.default.ids.0
  log_project_1 = alicloud_log_project.default_1.name
  log_project_2 = alicloud_log_project.default_2.name
}

`, name)
}

func TestAccAlicloudServiceMeshServiceMesh_basic3(t *testing.T) {
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_mesh_name": "${var.name}",
					"edition":           "Default",
					"version":           "${data.alicloud_service_mesh_versions.default.versions.0.version}",
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
									"enabled": "true",
									"project": "${local.log_project_1}",
								},
							},
							"control_plane_log": []map[string]interface{}{
								{
									"enabled": "true",
									"project": "${local.log_project_1}",
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
	available_resource_creation= "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id     	= data.alicloud_zones.default.zones.0.id
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
  vswitch_id = data.alicloud_vswitches.default.ids.0
  vpc_id = data.alicloud_vpcs.default.ids.0
  log_project_1 = alicloud_log_project.default_1.name
  log_project_2 = alicloud_log_project.default_2.name
}

`, name)
}
