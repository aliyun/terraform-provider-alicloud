package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_eci_container_group",
		&resource.Sweeper{
			Name: "alicloud_eci_container_group",
			F:    testSweepEciContainerGroup,
		})
}

func testSweepEciContainerGroup(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EciContainerGroupRegions) {
		log.Printf("[INFO] Skipping ECI unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc`",
	}
	var response map[string]interface{}
	action := "DescribeContainerGroups"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eci_container_group", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	resp, err := jsonpath.Get("$.ContainerGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ContainerGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["ContainerGroupName"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Eci containerGroup: %s ", name)
			continue
		}
		log.Printf("[INFO] Delete Eci containerGroup: %s ", name)
		action := "DeleteContainerGroup"
		request := map[string]interface{}{
			"ContainerGroupId": item["ContainerGroupId"],
		}
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*10, func() *retry.RetryError {
			response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Eci containerGroup (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudEciContainerGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":   strings.ToLower(name),
					"security_group_id":      "${alicloud_security_group.default.id}",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"auto_match_image_cache": "false",
					"ephemeral_storage":      "20",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"ephemeral_storage":    "20",
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"user_name": name + "_update",
							"password":  "tftestacc" + "_update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_registry_credential.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "/tmp/test",
									"read_only":  "false",
									"name":       "empty1",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"containers.#":      "1",
						"volumes.#":         "1",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"restart_policy": "OnFailure",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"restart_policy": "OnFailure",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"restart_policy":    "Always",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"containers": []map[string]interface{}{
						{
							"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/centos:7", defaultRegionToTest),
							"name":     "centos",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"restart_policy":    "Always",
						"containers.#":      "1",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"auto_match_image_cache": "false",
					"auto_create_eip":        "true",
					"eip_bandwidth":          "5",
					"cpu":                    "2",
					"memory":                 "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"resource_group_id":    CHECKSET,
						"internet_ip":          CHECKSET,
						"cpu":                  "2",
						"memory":               "4",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":   strings.ToLower(name),
					"security_group_id":      "${alicloud_security_group.default.id}",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"auto_match_image_cache": "false",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"plain_http_registry": "harbor.pre.com,192.168.1.1:5000,reg.test.com:80",
					"insecure_registry":   "harbor.pre.com,192.168.1.1:5000,reg.test.com:80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "insecure_registry", "plain_http_registry", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"auto_match_image_cache": "false",
					"eip_instance_id":        "${alicloud_eip_address.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"internet_ip":          CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "eip_instance_id", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacceci-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":   strings.ToLower(name),
					"security_group_id":      "${alicloud_security_group.default.id}",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"auto_match_image_cache": "false",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"exec": []map[string]interface{}{
										{
											"commands": []string{"cat /tmp/healthy"},
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"exec": []map[string]interface{}{
										{
											"commands": []string{"cat /tmp/healthy"},
										},
									},
								},
							},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"acr_registry_info": []map[string]interface{}{
						{
							"instance_name": "acr_test_name",
							"instance_id":   "acr-1",
							"domains":       []string{fmt.Sprintf("registry.%s.cr.aliyuncs.com", defaultRegionToTest)},
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"exec": []map[string]interface{}{
										{
											"commands": []string{"/bin/sh cat /tmp/healthy"},
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"exec": []map[string]interface{}{
										{
											"commands": []string{"/bin/sh cat /tmp/healthy"},
										},
									},
								},
							},
						},
					},
					"acr_registry_info": []map[string]interface{}{
						{
							"instance_name": "acr_test_name_2",
							"instance_id":   "acr-2",
							"domains":       []string{"test"},
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"containers.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "acr_registry_info", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacceci-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":   strings.ToLower(name),
					"security_group_id":      "${alicloud_security_group.default.id}",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"auto_match_image_cache": "false",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"http_get": []map[string]interface{}{
										{
											"scheme": "HTTP",
											"path":   "/healthyz",
											"port":   "8080",
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"http_get": []map[string]interface{}{
										{
											"scheme": "HTTP",
											"path":   "/healthyz",
											"port":   "8888",
										},
									},
								},
							},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"acr_registry_info": []map[string]interface{}{
						{
							"instance_name": "acr_test_name",
							"instance_id":   "acr-1",
							"domains":       []string{fmt.Sprintf("registry.%s.cr.aliyuncs.com", defaultRegionToTest)},
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"http_get": []map[string]interface{}{
										{
											"scheme": "HTTPS",
											"path":   "/usr/local/bin",
											"port":   "8081",
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"http_get": []map[string]interface{}{
										{
											"scheme": "HTTPS",
											"path":   "/usr/",
											"port":   "8081",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"containers.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "acr_registry_info", "auto_match_image_cache", "containers.0.restart_count"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":             strings.ToLower(name),
					"security_group_id":                "${alicloud_security_group.default.id}",
					"vswitch_id":                       "${alicloud_vswitch.default.id}",
					"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"termination_grace_period_seconds": "10",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "0.5",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh"},
							"args":              []string{"-c", "sleep 9999"},
							"cpu":               "1",
							"memory":            "2",
							"gpu":               "0",
							"working_dir":       "/home/1",
							"ports": []map[string]interface{}{
								{
									"protocol": "UDP",
									"port":     "800",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"capability": []map[string]interface{}{
										{
											"add": []string{"NET_ADMIN"},
										},
									},
									"run_as_user": "0",
									"privileged":  true,
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "container_n1",
									"value": "container_v1",
								},
								{
									"key": "container_n2",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "status.hostIP",
										},
									},
								},
							},
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"tcp_socket": []map[string]interface{}{
										{
											"port": "9090",
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "5",
									"initial_delay_seconds": "5",
									"success_threshold":     "1",
									"failure_threshold":     "3",
									"timeout_seconds":       "1",
									"tcp_socket": []map[string]interface{}{
										{
											"port": "8080",
										},
									},
								},
							},
							"lifecycle_pre_stop_handler_exec": []string{"/bin/sh", "-c", "sleep 10"}},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
							"cpu":               "1",
							"memory":            "2",
							"gpu":               "0",
							"working_dir":       "/home/2",
							"ports": []map[string]interface{}{
								{
									"protocol": "TCP",
									"port":     "80",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"capability": []map[string]interface{}{
										{
											"add": []string{"SYS_PTRACE"},
										},
									},
									"run_as_user": "0",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "init_container_n1",
									"value": "init_container_v1",
								},
								{
									"key": "init_container_n2",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "status.podIP",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "/tmp/test",
									"read_only":  "false",
									"name":       "empty1",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "1.1.1.1",
							"hostnames": []string{"hehe.com"},
						},
					},
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", defaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"auto_match_image_cache": "false",
					"auto_create_eip":        "true",
					"eip_bandwidth":          "5",
					"cpu":                    "2",
					"memory":                 "4",
					"dns_config": []map[string]interface{}{
						{
							"name_servers": []string{"1.1.1.1"},
							"searches":     []string{"aliyun.com"},
							"options": []map[string]interface{}{
								{
									"name":  "ndots",
									"value": "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"resource_group_id":    CHECKSET,
						"internet_ip":          CHECKSET,
						"cpu":                  "2",
						"memory":               "4",
						"spot_strategy":        "SpotWithPriceLimit",
						"spot_price_limit":     "0.5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"termination_grace_period_seconds": "11",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "0.5",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:perl", defaultRegionToTest),
							"image_pull_policy": "Always",
							"commands":          []string{"/bin/sh", "-c"},
							"args":              []string{"sleep 9999"},
							"cpu":               "1.5",
							"memory":            "2.5",
							"working_dir":       "/home/1-update",
							"ports": []map[string]interface{}{
								{
									"protocol": "TCP",
									"port":     "8000",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"capability": []map[string]interface{}{
										{
											"add": []string{"SYS_PTRACE"},
										},
									},
									"run_as_user": "1",
									"privileged":  false,
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "container_n1_update",
									"value": "container_v1_update",
								},
								{
									"key": "container_n2",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "status.podIP",
										},
									},
								},
							},
							"liveness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "6",
									"timeout_seconds":       "2",
									"tcp_socket": []map[string]interface{}{
										{
											"port": "1010",
										},
									},
								},
							},
							"readiness_probe": []map[string]interface{}{
								{
									"period_seconds":        "10",
									"initial_delay_seconds": "10",
									"success_threshold":     "1",
									"failure_threshold":     "6",
									"timeout_seconds":       "2",
									"tcp_socket": []map[string]interface{}{
										{
											"port": "2020",
										},
									},
								},
							},
							"lifecycle_pre_stop_handler_exec": []string{"/bin/sh", "-c", "sleep 11"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/ubuntu:18.04", defaultRegionToTest),
							"image_pull_policy": "Always",
							"commands":          []string{"/bin/sh", "-c", "echo"},
							"args":              []string{"hello initcontainer-update"},
							"cpu":               "0.5",
							"memory":            "1",
							"working_dir":       "/home/2-update",
							"ports": []map[string]interface{}{
								{
									"protocol": "UDP",
									"port":     "8000",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"capability": []map[string]interface{}{
										{
											"add": []string{"SYS_CHROOT"},
										},
									},
									"run_as_user": "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "init_container_n1_update",
									"value": "init_container_v1_update",
								},
								{
									"key": "init_container_n2_update",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "status.hostIP",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "/tmp/test-update",
									"read_only":  "true",
									"name":       "empty1",
								},
							},
						},
					},

					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
					},
					"dns_config": []map[string]interface{}{
						{
							"name_servers": []string{"2.2.2.2"},
							"searches":     []string{"taobao.com"},
							"options": []map[string]interface{}{
								{
									"name":  "com",
									"value": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"init_containers.#":    "1",
						"host_aliases.#":       "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"resource_group_id":    CHECKSET,
						"internet_ip":          CHECKSET,
						"cpu":                  "2",
						"memory":               "4",
						"dns_config.#":         "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "termination_grace_period_seconds", "containers.0.lifecycle_pre_stop_handler_exec", "auto_match_image_cache", "containers.0.security_context.0.privileged"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic7(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":   strings.ToLower(name),
					"security_group_id":      "${alicloud_security_group.default.id}",
					"auto_match_image_cache": "false",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"zone_id":                "${data.alicloud_zones.default.zones.0.id}",
					"instance_type":          "${data.alicloud_instance_types.default.instance_types.0.id}",
					"ram_role_name":          "${alicloud_ram_role.default.name}",
					"security_context": []map[string]interface{}{
						{
							"sysctl": []map[string]interface{}{
								{
									"name":  "net.ipv4.ping_group_range",
									"value": "1",
								},
							},
						},
					},
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", "${data.alicloud_regions.default.regions.0.id}"),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "/tmp/nas",
									"read_only":  "true",
									"name":       "nfsVolume",
								},
								{
									"mount_path": "/tmp/config",
									"read_only":  "true",
									"name":       "configFileVolume",
								},
								{
									"mount_path": "/tmp/flex",
									"read_only":  "true",
									"name":       "flexVolume",
								},
								{
									"mount_path": "/tmp/disk",
									"read_only":  "true",
									"name":       "diskVolume",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name":                 "nfsVolume",
							"type":                 "NFSVolume",
							"nfs_volume_server":    "${alicloud_nas_mount_target.default.mount_target_domain}",
							"nfs_volume_path":      "/",
							"nfs_volume_read_only": "true",
						},
						{
							"name": "configFileVolume",
							"type": "ConfigFileVolume",
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "aGVsbG8gd29ybGQ=",
									"path":    "config",
								},
							},
						},
						{
							"name":                "flexVolume",
							"type":                "FlexVolume",
							"flex_volume_driver":  "alicloud/disk",
							"flex_volume_options": "{\\\"volumeSize\\\":\\\"50\\\",\\\"tags\\\":\\\"test:eci\\\"}",
							"flex_volume_fs_type": "ext4",
						},
						{
							"name":                "diskVolume",
							"type":                "DiskVolume",
							"disk_volume_fs_type": "ext4",
							"disk_volume_disk_id": "${alicloud_ecs_disk.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_context": []map[string]interface{}{
						{
							"sysctl": []map[string]interface{}{
								{
									"name":  "net.ipv4.ping_group_range",
									"value": "0 10",
								},
							},
						},
					},
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", "${data.alicloud_regions.default.regions.0.id}"),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "/tmp/nas",
									"read_only":  "true",
									"name":       "nfsVolume",
								},
								{
									"mount_path": "/tmp/config",
									"read_only":  "true",
									"name":       "configFileVolume",
								},
								{
									"mount_path": "/tmp/flex",
									"read_only":  "true",
									"name":       "flexVolume",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name":                 "nfsVolume",
							"type":                 "NFSVolume",
							"nfs_volume_server":    "${alicloud_nas_mount_target.default.mount_target_domain}",
							"nfs_volume_path":      "/",
							"nfs_volume_read_only": "true",
						},
						{
							"name": "configFileVolume",
							"type": "ConfigFileVolume",
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "aGVsbG8gd29ybGQ=",
									"path":    "config",
								},
							},
						},
						{
							"name":                "flexVolume",
							"type":                "FlexVolume",
							"flex_volume_driver":  "alicloud/disk",
							"flex_volume_options": "{\\\"volumeSize\\\":\\\"50\\\",\\\"tags\\\":\\\"test:eci\\\"}",
							"flex_volume_fs_type": "ext4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "eip_instance_id", "volumes.1.config_file_volume_config_file_to_paths.0.content", "auto_match_image_cache", "auto_match_image_cache"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic8(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh"},
							"args":              []string{"-c", "sleep 9999"},
							"cpu":               "1",
							"memory":            "2",
							"gpu":               "0",
						},
					},
					"auto_match_image_cache": "false",
					"cpu":                    "2",
					"memory":                 "4",
					"dns_config": []map[string]interface{}{
						{
							"name_servers": []string{"1.1.1.1"},
							"searches":     []string{"aliyun.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"resource_group_id":    CHECKSET,
						"cpu":                  "2",
						"memory":               "4",
						"dns_policy":           "Default",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "termination_grace_period_seconds", "containers.0.lifecycle_pre_stop_handler_exec", "auto_match_image_cache", "containers.0.security_context.0.privileged"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic9(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh"},
							"args":              []string{"-c", "sleep 9999"},
							"cpu":               "1",
							"memory":            "2",
							"gpu":               "0",
						},
					},
					"auto_match_image_cache": "false",
					"cpu":                    "2",
					"memory":                 "4",
					"dns_policy":             "None",
					"dns_config": []map[string]interface{}{
						{
							"name_servers": []string{"1.1.1.1"},
							"searches":     []string{"aliyun.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
						"containers.#":         "1",
						"security_group_id":    CHECKSET,
						"vswitch_id":           CHECKSET,
						"resource_group_id":    CHECKSET,
						"cpu":                  "2",
						"memory":               "4",
						"dns_policy":           "None",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential", "eip_bandwidth", "auto_create_eip", "termination_grace_period_seconds", "containers.0.lifecycle_pre_stop_handler_exec", "auto_match_image_cache", "containers.0.security_context.0.privileged"},
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic10(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":         strings.ToLower(name),
					"security_group_id":            "${alicloud_security_group.default.id}",
					"vswitch_id":                   "${alicloud_vswitch.default.id}",
					"auto_match_image_cache":       "false",
					"active_deadline_seconds":      "3600",
					"compute_category_list":        []string{"economy"},
					"container_resource_view":      "true",
					"core_pattern":                 "/var/core/core-%e-%p",
					"cpu":                          "1",
					"cpu_architecture":             "AMD64",
					"cpu_options_core":             "1",
					"cpu_options_numa":             "None",
					"cpu_options_threads_per_core": "1",
					"egress_bandwidth":             "1048576",
					"ingress_bandwidth":            "1048576",
					"environment_var_hide":         "false",
					"fixed_ip":                     "false",
					"fixed_ip_retain_hour":         "48",
					"image_accelerate_mode":        "imc",
					"memory":                       "4",
					"ntp_server":                   []string{"ntp1.aliyun.com", "ntp2.aliyun.com"},
					"private_ip_address":           "172.16.0.100",
					"schedule_strategy":            "VSwitchOrdered",
					"share_process_namespace":      "true",
					"update_type":                  "RenewUpdate",
					"overhead_reservation_option": []map[string]interface{}{
						{
							"enable_overhead_reservation": "false",
						},
					},
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
							"stdin":             "true",
							"stdin_once":        "true",
							"tty":               "true",
							"environment_vars": []map[string]interface{}{
								{
									"key": "test",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "metadata.name",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "/tmp/test",
									"read_only":         "false",
									"name":              "empty1",
									"sub_path":          "sub",
									"mount_propagation": "None",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"run_as_user":               "1000",
									"read_only_root_filesystem": "false",
								},
							},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
							"environment_vars": []map[string]interface{}{
								{
									"key": "init-test",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "metadata.name",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "/tmp/init",
									"read_only":         "false",
									"name":              "empty1",
									"mount_propagation": "None",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"read_only_root_filesystem": "false",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name":          strings.ToLower(name),
						"containers.#":                  "1",
						"init_containers.#":             "1",
						"active_deadline_seconds":       "3600",
						"compute_category_list.#":       "1",
						"compute_category_list.0":       "economy",
						"container_resource_view":       "true",
						"cpu":                           "1",
						"cpu_architecture":              "AMD64",
						"cpu_options_core":              "1",
						"cpu_options_threads_per_core":  "1",
						"egress_bandwidth":              "1048576",
						"ingress_bandwidth":             "1048576",
						"environment_var_hide":          "false",
						"fixed_ip":                      "false",
						"fixed_ip_retain_hour":          "48",
						"image_accelerate_mode":         "imc",
						"memory":                        "4",
						"ntp_server.#":                  "2",
						"ntp_server.0":                  "ntp1.aliyun.com",
						"ntp_server.1":                  "ntp2.aliyun.com",
						"private_ip_address":            "172.16.0.100",
						"schedule_strategy":             "VSwitchOrdered",
						"share_process_namespace":       "true",
						"update_type":                   "RenewUpdate",
						"overhead_reservation_option.#": "1",
						"overhead_reservation_option.0.enable_overhead_reservation": "false",
						"containers.0.stdin":                                             "true",
						"containers.0.stdin_once":                                        "true",
						"containers.0.tty":                                               "true",
						"containers.0.environment_vars.#":                                "1",
						"containers.0.environment_vars.0.key":                            "test",
						"containers.0.environment_vars.0.field_ref.#":                    "1",
						"containers.0.environment_vars.0.field_ref.0.field_path":         "metadata.name",
						"containers.0.volume_mounts.0.sub_path":                          "sub",
						"containers.0.volume_mounts.0.mount_propagation":                 "None",
						"containers.0.security_context.0.read_only_root_filesystem":      "false",
						"init_containers.0.volume_mounts.0.mount_propagation":            "None",
						"init_containers.0.security_context.0.read_only_root_filesystem": "false",
						"security_group_id":                                              CHECKSET,
						"vswitch_id":                                                     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":         strings.ToLower(name),
					"security_group_id":            "${alicloud_security_group.default.id}",
					"vswitch_id":                   "${alicloud_vswitch.default.id}",
					"auto_match_image_cache":       "false",
					"active_deadline_seconds":      "3600",
					"compute_category_list":        []string{"economy"},
					"container_resource_view":      "true",
					"core_pattern":                 "/var/core/core-%e-%p",
					"cpu":                          "1",
					"cpu_architecture":             "AMD64",
					"cpu_options_core":             "1",
					"cpu_options_numa":             "None",
					"cpu_options_threads_per_core": "1",
					"egress_bandwidth":             "1048576",
					"ingress_bandwidth":            "1048576",
					"environment_var_hide":         "false",
					"fixed_ip":                     "false",
					"fixed_ip_retain_hour":         "48",
					"image_accelerate_mode":        "imc",
					"memory":                       "4",
					"ntp_server":                   []string{"ntp1.aliyun.com", "ntp2.aliyun.com"},
					"private_ip_address":           "172.16.0.100",
					"schedule_strategy":            "VSwitchOrdered",
					"share_process_namespace":      "true",
					"update_type":                  "IncrementalUpdate",
					"overhead_reservation_option": []map[string]interface{}{
						{
							"enable_overhead_reservation": "false",
						},
					},
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
							"stdin":             "false",
							"stdin_once":        "false",
							"tty":               "false",
							"environment_vars": []map[string]interface{}{
								{
									"key": "test",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "metadata.name",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "/tmp/test",
									"read_only":         "false",
									"name":              "empty1",
									"sub_path":          "sub2",
									"mount_propagation": "HostToContainer",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"run_as_user":               "1000",
									"read_only_root_filesystem": "true",
								},
							},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"echo"},
							"args":              []string{"hello initcontainer"},
							"environment_vars": []map[string]interface{}{
								{
									"key": "init-test",
									"field_ref": []map[string]interface{}{
										{
											"field_path": "metadata.name",
										},
									},
								},
							},
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "/tmp/init",
									"read_only":         "false",
									"name":              "empty1",
									"mount_propagation": "HostToContainer",
								},
							},
							"security_context": []map[string]interface{}{
								{
									"read_only_root_filesystem": "true",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":    "1",
						"memory": "4",
						"overhead_reservation_option.0.enable_overhead_reservation": "false",
						"update_type":                                                    "IncrementalUpdate",
						"containers.0.stdin":                                             "false",
						"containers.0.stdin_once":                                        "false",
						"containers.0.tty":                                               "false",
						"containers.0.volume_mounts.0.sub_path":                          "sub2",
						"containers.0.volume_mounts.0.mount_propagation":                 "HostToContainer",
						"containers.0.security_context.0.read_only_root_filesystem":      "true",
						"init_containers.0.volume_mounts.0.mount_propagation":            "HostToContainer",
						"init_containers.0.security_context.0.read_only_root_filesystem": "true",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEciContainerGroup_basic11(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":          strings.ToLower(name),
					"security_group_id":             "${alicloud_security_group.default.id}",
					"vswitch_id":                    "${alicloud_vswitch.default.id}",
					"auto_match_image_cache":        "false",
					"data_cache_bucket":             "${alicloud_oss_bucket.default.id}",
					"data_cache_bursting_enabled":   "false",
					"data_cache_pl":                 "PL1",
					"data_cache_provisioned_iops":   "3000",
					"image_snapshot_id":             "${alicloud_eci_image_cache.default.id}",
					"auto_create_eip":               "false",
					"eip_isp":                       "BGP",
					"eip_common_bandwidth_package":  "${alicloud_common_bandwidth_package.default.id}",
					"ipv6_address_count":            "1",
					"ipv6_gateway_bandwidth_enable": "true",
					"ipv6_gateway_bandwidth":        "200",
					"spot_strategy":                 "SpotWithPriceLimit",
					"spot_price_limit":              "0.5",
					"spot_duration":                 "1",
					"strict_spot":                   "false",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
							"cpu":               "2",
							"memory":            "4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name":        strings.ToLower(name),
						"containers.#":                "1",
						"ipv6_address":                CHECKSET,
						"spot_strategy":               "SpotWithPriceLimit",
						"spot_price_limit":            "0.5",
						"data_cache_bursting_enabled": "false",
						"data_cache_pl":               "PL1",
						"data_cache_provisioned_iops": "3000",
						"ipv6_address_count":          "1",
						"security_group_id":           CHECKSET,
						"vswitch_id":                  CHECKSET,
					}),
				),
			},
		},
	})
}

func AliCloudEciContainerGroupBasicDependence5(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    default = "%s"
	}

	variable "name_lower" {
    default = "%s"
	}

	data "alicloud_zones" "default" {
    available_disk_category     = "cloud_efficiency"
    available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
    cidr_block  = "172.16.0.0/16"
    enable_ipv6 = true
    ipv6_isp    = "BGP"
    vpc_name    = var.name
	}

	resource "alicloud_vswitch" "default" {
    vpc_id               = alicloud_vpc.default.id
    cidr_block           = "172.16.0.0/24"
    zone_id              = data.alicloud_zones.default.zones.0.id
    vswitch_name         = var.name
    ipv6_cidr_block_mask = "3"
	}

	resource "alicloud_security_group" "default" {
    name   = var.name
    vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_oss_bucket" "default" {
    bucket = "%s"
	}

	resource "alicloud_common_bandwidth_package" "default" {
    bandwidth = "1000"
    name      = var.name
	}

	resource "alicloud_eci_image_cache" "default" {
    image_cache_name  = var.name_lower
    images            = ["registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine"]
    security_group_id = alicloud_security_group.default.id
    vswitch_id        = alicloud_vswitch.default.id
	}
`, name, strings.ToLower(name), strings.ToLower(name), defaultRegionToTest)
}

var AliCloudEciContainerGroupMap = map[string]string{
	"resource_group_id": CHECKSET,
	"restart_policy":    "Always",
	"status":            CHECKSET,
}

func AliCloudEciContainerGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
%s
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, EcsInstanceCommonTestCase, name)
}

func AliCloudEciContainerGroupBasicDependence2(name string) string {
	return fmt.Sprintf(`
%s
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_eip_address" "default" {
  		address_name = var.name
	}
`, EcsInstanceCommonTestCase, name)
}

func AliCloudEciContainerGroupBasicDependence3(name string) string {
	return fmt.Sprintf(`
%s
	variable "name" {
  		default = "%s"
	}


`, EcsInstanceCommonTestCase, name)
}

func AliCloudEciContainerGroupBasicDependence4(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
  		default = "%s"
	}

	data "alicloud_regions" "default" {
  		current = true
	}

	data "alicloud_nas_protocols" "default" {
  		type = "Capacity"
	}

	resource "alicloud_ram_role" "default" {
  		name        = var.name
  		document    = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
  		description = "this is a test"
  		force       = true
	}

	resource "alicloud_nas_file_system" "default" {
  		description   = var.name
  		storage_type  = "Capacity"
  		protocol_type = data.alicloud_nas_protocols.default.protocols.1
	}

	resource "alicloud_nas_mount_target" "default" {
  		file_system_id    = alicloud_nas_file_system.default.id
  		access_group_name = "DEFAULT_VPC_GROUP_NAME"
  		vswitch_id        = alicloud_vswitch.default.id
	}

	resource "alicloud_ecs_disk" "default" {
  		zone_id              = data.alicloud_zones.default.zones.0.id
  		category             = "cloud_efficiency"
  		description          = "Test For Terraform"
  		disk_name            = var.name
  		enable_auto_snapshot = "true"
  		encrypted            = "true"
  		size                 = "500"
  		tags = {
    		Created     = "TF"
    		Environment = "Acceptance-test"
  		}
	}
`, EcsInstanceCommonTestCase, name)
}

// TypeList order coverage: the configurable multi-member TypeList attributes
// bypassed below pre-date this change and the diff does not touch their
// ordering semantics; see scripts/collection-order/README.md.
// collection-order: bypass acr_registry_info.domains -- pre-existing attribute outside this change
// collection-order: bypass containers.args -- pre-existing attribute outside this change
// collection-order: bypass containers.commands -- pre-existing attribute outside this change
// collection-order: bypass containers.environment_vars -- pre-existing attribute outside this change
// collection-order: bypass containers.environment_vars.field_ref -- pre-existing attribute outside this change
// collection-order: bypass containers.lifecycle_pre_stop_handler_exec -- pre-existing attribute outside this change
// collection-order: bypass containers.liveness_probe -- pre-existing attribute outside this change
// collection-order: bypass containers.liveness_probe.exec -- pre-existing attribute outside this change
// collection-order: bypass containers.liveness_probe.exec.commands -- pre-existing attribute outside this change
// collection-order: bypass containers.liveness_probe.http_get -- pre-existing attribute outside this change
// collection-order: bypass containers.liveness_probe.tcp_socket -- pre-existing attribute outside this change
// collection-order: bypass containers.ports -- pre-existing attribute outside this change
// collection-order: bypass containers.readiness_probe -- pre-existing attribute outside this change
// collection-order: bypass containers.readiness_probe.exec -- pre-existing attribute outside this change
// collection-order: bypass containers.readiness_probe.exec.commands -- pre-existing attribute outside this change
// collection-order: bypass containers.readiness_probe.http_get -- pre-existing attribute outside this change
// collection-order: bypass containers.readiness_probe.tcp_socket -- pre-existing attribute outside this change
// collection-order: bypass containers.security_context.capability -- pre-existing attribute outside this change
// collection-order: bypass containers.security_context.capability.add -- pre-existing attribute outside this change
// collection-order: bypass dns_config.name_servers -- pre-existing attribute outside this change
// collection-order: bypass dns_config.options -- pre-existing attribute outside this change
// collection-order: bypass dns_config.searches -- pre-existing attribute outside this change
// collection-order: bypass host_aliases -- pre-existing attribute outside this change
// collection-order: bypass host_aliases.hostnames -- pre-existing attribute outside this change
// collection-order: bypass init_containers.args -- pre-existing attribute outside this change
// collection-order: bypass init_containers.commands -- pre-existing attribute outside this change
// collection-order: bypass init_containers.environment_vars -- pre-existing attribute outside this change
// collection-order: bypass init_containers.environment_vars.field_ref -- pre-existing attribute outside this change
// collection-order: bypass init_containers.ports -- pre-existing attribute outside this change
// collection-order: bypass init_containers.security_context.capability -- pre-existing attribute outside this change
// collection-order: bypass init_containers.security_context.capability.add -- pre-existing attribute outside this change
// collection-order: bypass security_context.sysctl -- pre-existing attribute outside this change
// collection-order: bypass volumes -- pre-existing attribute outside this change
// collection-order: bypass volumes.config_file_volume_config_file_to_paths -- pre-existing attribute outside this change

// basic12 establishes TypeList order coverage for the attributes this change
// introduces or extends: each consecutive apply pair reorders exactly one
// collection and holds every other value identical, so the empty post-apply
// plan proves the resource converges on the reordered members.
func TestAccAliCloudEciContainerGroup_basic12(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name":  strings.ToLower(name),
					"security_group_id":     "${alicloud_security_group.default.id}",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"cpu":                   "1",
					"memory":                "4",
					"update_type":           "RenewUpdate",
					"compute_category_list": []string{"economy", "general"},
					"ntp_server":            []string{"ntp1.aliyun.com", "ntp2.aliyun.com"},
					"containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-ct-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-ct-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
							},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-init-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-init-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
							},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"name": "empty1",
							"type": "EmptyDirVolume",
						},
						{
							"name": "empty2",
							"type": "EmptyDirVolume",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_group_name": strings.ToLower(name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-ct-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-ct-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"init_containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-init-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-init-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-ct-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-ct-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/mount-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/mount-b",
									"read_only":  "false",
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"init_containers": []map[string]interface{}{
						{
							"name":     fmt.Sprintf("%s-init-2", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
							},
						},
						{
							"name":     fmt.Sprintf("%s-init-1", strings.ToLower(name)),
							"image":    "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/busybox:1.30",
							"commands": []string{"echo"},
							"volume_mounts": []map[string]interface{}{
								{
									"name":       "empty1",
									"mount_path": "/tmp/init-a",
									"read_only":  "false",
								},
								{
									"name":       "empty2",
									"mount_path": "/tmp/init-b",
									"read_only":  "false",
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_category_list": []string{"general", "economy"},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ntp_server": []string{"ntp2.aliyun.com", "ntp1.aliyun.com"},
				}),
			},
		},
	})
}
