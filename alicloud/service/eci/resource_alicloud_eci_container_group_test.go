package eci_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eci"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	if eciTestSweepPreCheckWithRegions(region, true, connectivity.EciContainerGroupRegions) {
		log.Printf("[INFO] Skipping ECI unsupported region: %s", region)
		return nil
	}
	rawClient, err := eciTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_eci_container_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	resp, err := jsonpath.Get("$.ContainerGroups", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ContainerGroups", response)
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*10, func() *retry.RetryError {
			response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_registry_credential": []map[string]interface{}{
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
						{
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
						"resource_group_id": acctest.CHECKSET,
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
							"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/centos:7", acctest.DefaultRegionToTest),
							"name":     "centos",
							"commands": []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"internet_ip":          acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
							"image_pull_policy": "IfNotPresent",
							"commands":          []string{"/bin/sh", "-c", "sleep 9999"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"name":              "init-busybox",
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"internet_ip":          acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacceci-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"acr_registry_info": []map[string]interface{}{
						{
							"instance_name": "acr_test_name",
							"instance_id":   "acr-1",
							"domains":       []string{fmt.Sprintf("registry.%s.cr.aliyuncs.com", acctest.DefaultRegionToTest)},
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacceci-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
							"user_name": name,
							"password":  "tftestacc",
						},
					},
					"acr_registry_info": []map[string]interface{}{
						{
							"instance_name": "acr_test_name",
							"instance_id":   "acr-1",
							"domains":       []string{fmt.Sprintf("registry.%s.cr.aliyuncs.com", acctest.DefaultRegionToTest)},
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":              strings.ToLower(name),
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
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
							"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/google_containers/etcd:3.4.3-0", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"internet_ip":          acctest.CHECKSET,
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:perl", acctest.DefaultRegionToTest),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/ubuntu:18.04", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"internet_ip":          acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"image":             fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
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
						"security_group_id":    acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
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

// EcsInstanceCommonTestCase is the shared ECS common test-case fixture (zones /
// vpc / vswitch / security group scaffolding) retained locally in the eci test
// package because the original const lives in the root alicloud test package
// (service_alicloud_common_test.go) and test-only symbols cannot be imported
// across packages. The root copy remains for unmigrated consumers (ess,
// cloud_monitor). Behavior is unchanged.
const EcsInstanceCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`

var AliCloudEciContainerGroupMap = map[string]string{
	"resource_group_id": acctest.CHECKSET,
	"restart_policy":    "Always",
	"status":            acctest.CHECKSET,
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

// collection-order coverage for alicloud_eci_container_group TypeList candidates.
// Each step pair applies the original member order, then a legitimate convergent
// reorder of a single TypeList candidate (>=2 distinct members, clean apply,
// no StaticSkip, no ExpectNonEmptyPlan). Single-value (loop-overwrite) and
// object-list candidates whose API order-preservation cannot be verified are
// declared as stock debt via bypass markers below.
func TestAccAliCloudEciContainerGroup_collectionOrderContainers(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// containers.commands: original order [/bin/sh, -c] -> [-c, /bin/sh]
			{Config: testAccConfig(map[string]interface{}{
				"container_group_name":   strings.ToLower(name),
				"security_group_id":      "${alicloud_security_group.default.id}",
				"vswitch_id":             "${alicloud_vswitch.default.id}",
				"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				"auto_match_image_cache": "false",
				"containers": []map[string]interface{}{{
					"name":     strings.ToLower(name),
					"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"commands": []string{"/bin/sh", "-c"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":     strings.ToLower(name),
					"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"commands": []string{"-c", "/bin/sh"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.args: [-c, sleep 1] -> [sleep 1, -c]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"args":  []string{"-c", "sleep 1"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"args":  []string{"sleep 1", "-c"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.ports: [{80,TCP},{443,TCP}] -> [{443,TCP},{80,TCP}]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"ports": []map[string]interface{}{
						{"port": "80", "protocol": "TCP"},
						{"port": "443", "protocol": "TCP"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"ports": []map[string]interface{}{
						{"port": "443", "protocol": "TCP"},
						{"port": "80", "protocol": "TCP"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.environment_vars: [{k1,v1},{k2,v2}] -> [{k2,v2},{k1,v1}]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"environment_vars": []map[string]interface{}{
						{"key": "k1", "value": "v1"},
						{"key": "k2", "value": "v2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"environment_vars": []map[string]interface{}{
						{"key": "k2", "value": "v2"},
						{"key": "k1", "value": "v1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.volume_mounts: [{/tmp/a,v1},{/tmp/b,v2}] -> [{/tmp/b,v2},{/tmp/a,v1}]
			{Config: testAccConfig(map[string]interface{}{
				"volumes": []map[string]interface{}{
					{"name": "v1", "type": "EmptyDirVolume"},
					{"name": "v2", "type": "EmptyDirVolume"},
				},
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"volume_mounts": []map[string]interface{}{
						{"mount_path": "/tmp/a", "name": "v1"},
						{"mount_path": "/tmp/b", "name": "v2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"volume_mounts": []map[string]interface{}{
						{"mount_path": "/tmp/b", "name": "v2"},
						{"mount_path": "/tmp/a", "name": "v1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.lifecycle_pre_stop_handler_exec: [/bin/sh,-c,sleep 10] -> [sleep 10,-c,/bin/sh]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":                            strings.ToLower(name),
					"image":                           fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"lifecycle_pre_stop_handler_exec": []string{"/bin/sh", "-c", "sleep 10"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":                            strings.ToLower(name),
					"image":                           fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"lifecycle_pre_stop_handler_exec": []string{"sleep 10", "-c", "/bin/sh"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.liveness_probe.exec.commands: [cmd1,cmd2] -> [cmd2,cmd1]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"liveness_probe": []map[string]interface{}{{
						"exec": []map[string]interface{}{{
							"commands": []string{"cmd1", "cmd2"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"liveness_probe": []map[string]interface{}{{
						"exec": []map[string]interface{}{{
							"commands": []string{"cmd2", "cmd1"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.readiness_probe.exec.commands: [r1,r2] -> [r2,r1]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"readiness_probe": []map[string]interface{}{{
						"exec": []map[string]interface{}{{
							"commands": []string{"r1", "r2"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"readiness_probe": []map[string]interface{}{{
						"exec": []map[string]interface{}{{
							"commands": []string{"r2", "r1"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// containers.security_context.capability.add: [NET_ADMIN,SYS_TIME] -> [SYS_TIME,NET_ADMIN]
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"security_context": []map[string]interface{}{{
						"capability": []map[string]interface{}{{
							"add": []string{"NET_ADMIN", "SYS_TIME"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
					"security_context": []map[string]interface{}{{
						"capability": []map[string]interface{}{{
							"add": []string{"SYS_TIME", "NET_ADMIN"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
		},
	})
}

func TestAccAliCloudEciContainerGroup_collectionOrderInitContainers(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// init_containers.commands: [echo, hi] -> [hi, echo]
			{Config: testAccConfig(map[string]interface{}{
				"container_group_name":   strings.ToLower(name),
				"security_group_id":      "${alicloud_security_group.default.id}",
				"vswitch_id":             "${alicloud_vswitch.default.id}",
				"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				"auto_match_image_cache": "false",
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
				}},
				"init_containers": []map[string]interface{}{{
					"name":     "init-one",
					"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"commands": []string{"echo", "hi"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":     "init-one",
					"image":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"commands": []string{"hi", "echo"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// init_containers.args: [-c, hello] -> [hello, -c]
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"args":  []string{"-c", "hello"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"args":  []string{"hello", "-c"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// init_containers.ports: [{80,TCP},{443,TCP}] -> [{443,TCP},{80,TCP}]
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"ports": []map[string]interface{}{
						{"port": "80", "protocol": "TCP"},
						{"port": "443", "protocol": "TCP"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"ports": []map[string]interface{}{
						{"port": "443", "protocol": "TCP"},
						{"port": "80", "protocol": "TCP"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// init_containers.environment_vars: [{k1,v1},{k2,v2}] -> [{k2,v2},{k1,v1}]
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"environment_vars": []map[string]interface{}{
						{"key": "k1", "value": "v1"},
						{"key": "k2", "value": "v2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"environment_vars": []map[string]interface{}{
						{"key": "k2", "value": "v2"},
						{"key": "k1", "value": "v1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// init_containers.volume_mounts: [{/tmp/a,v1},{/tmp/b,v2}] -> [{/tmp/b,v2},{/tmp/a,v1}]
			{Config: testAccConfig(map[string]interface{}{
				"volumes": []map[string]interface{}{
					{"name": "v1", "type": "EmptyDirVolume"},
					{"name": "v2", "type": "EmptyDirVolume"},
				},
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"volume_mounts": []map[string]interface{}{
						{"mount_path": "/tmp/a", "name": "v1"},
						{"mount_path": "/tmp/b", "name": "v2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"volume_mounts": []map[string]interface{}{
						{"mount_path": "/tmp/b", "name": "v2"},
						{"mount_path": "/tmp/a", "name": "v1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// init_containers.security_context.capability.add: [NET_ADMIN,SYS_TIME] -> [SYS_TIME,NET_ADMIN]
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"security_context": []map[string]interface{}{{
						"capability": []map[string]interface{}{{
							"add": []string{"NET_ADMIN", "SYS_TIME"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"init_containers": []map[string]interface{}{{
					"name":  "init-one",
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/busybox:1.30", acctest.DefaultRegionToTest),
					"security_context": []map[string]interface{}{{
						"capability": []map[string]interface{}{{
							"add": []string{"SYS_TIME", "NET_ADMIN"},
						}},
					}},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
		},
	})
}

func TestAccAliCloudEciContainerGroup_collectionOrderToplevel(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEciContainerGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciContainerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEciContainerGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEciContainerGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// host_aliases.hostnames: [a.com,b.com] -> [b.com,a.com] (ForceNew recreate)
			{Config: testAccConfig(map[string]interface{}{
				"container_group_name":   strings.ToLower(name),
				"security_group_id":      "${alicloud_security_group.default.id}",
				"vswitch_id":             "${alicloud_vswitch.default.id}",
				"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				"auto_match_image_cache": "false",
				"containers": []map[string]interface{}{{
					"name":  strings.ToLower(name),
					"image": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest),
				}},
				"host_aliases": []map[string]interface{}{{
					"ip":        "1.1.1.1",
					"hostnames": []string{"a.com", "b.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"host_aliases": []map[string]interface{}{{
					"ip":        "1.1.1.1",
					"hostnames": []string{"b.com", "a.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// dns_config.name_servers: [1.1.1.1,8.8.8.8] -> [8.8.8.8,1.1.1.1]
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"name_servers": []string{"1.1.1.1", "8.8.8.8"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"name_servers": []string{"8.8.8.8", "1.1.1.1"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// dns_config.searches: [a.com,b.com] -> [b.com,a.com]
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"searches": []string{"a.com", "b.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"searches": []string{"b.com", "a.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// dns_config.options: [{ndots,2},{timeout,1}] -> [{timeout,1},{ndots,2}]
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"options": []map[string]interface{}{
						{"name": "ndots", "value": "2"},
						{"name": "timeout", "value": "1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"dns_config": []map[string]interface{}{{
					"options": []map[string]interface{}{
						{"name": "timeout", "value": "1"},
						{"name": "ndots", "value": "2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// security_context.sysctl: [{n1,v1},{n2,v2}] -> [{n2,v2},{n1,v1}] (ForceNew recreate)
			{Config: testAccConfig(map[string]interface{}{
				"security_context": []map[string]interface{}{{
					"sysctl": []map[string]interface{}{
						{"name": "net.ipv4.ping_group_range", "value": "1"},
						{"name": "kernel.shm_rmid_forced", "value": "1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"security_context": []map[string]interface{}{{
					"sysctl": []map[string]interface{}{
						{"name": "kernel.shm_rmid_forced", "value": "1"},
						{"name": "net.ipv4.ping_group_range", "value": "1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// acr_registry_info.domains: [d1.com,d2.com] -> [d2.com,d1.com]
			{Config: testAccConfig(map[string]interface{}{
				"acr_registry_info": []map[string]interface{}{{
					"domains": []string{"d1.com", "d2.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"acr_registry_info": []map[string]interface{}{{
					"domains": []string{"d2.com", "d1.com"},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			// volumes.config_file_volume_config_file_to_paths: [{a,p1},{b,p2}] -> [{b,p2},{a,p1}] (ForceNew recreate)
			{Config: testAccConfig(map[string]interface{}{
				"volumes": []map[string]interface{}{{
					"name": "cfv",
					"type": "ConfigFileVolume",
					"config_file_volume_config_file_to_paths": []map[string]interface{}{
						{"content": "a", "path": "p1"},
						{"content": "b", "path": "p2"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
			{Config: testAccConfig(map[string]interface{}{
				"volumes": []map[string]interface{}{{
					"name": "cfv",
					"type": "ConfigFileVolume",
					"config_file_volume_config_file_to_paths": []map[string]interface{}{
						{"content": "b", "path": "p2"},
						{"content": "a", "path": "p1"},
					},
				}},
			}), Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{"container_group_name": strings.ToLower(name)}))},
		},
	})
}

// collection-order bypass declarations for alicloud_eci_container_group
// TypeList candidates not covered by a convergent member reorder above.
// All are pre-existing attributes surfaced by the ECI subpackage migration
// (file relocation, schema unchanged) and are stock debt carried forward.

// collection-order: bypass containers -- list-of-objects; ECI DescribeContainerGroups does not guarantee config-order preservation for complex object members, so a convergent member reorder cannot be verified.
// collection-order: bypass init_containers -- list-of-objects; ECI DescribeContainerGroups does not guarantee config-order preservation for complex object members, so a convergent member reorder cannot be verified.
// collection-order: bypass host_aliases -- list-of-objects; ECI DescribeContainerGroups does not guarantee config-order preservation for complex object members, so a convergent member reorder cannot be verified.
// collection-order: bypass volumes -- list-of-objects; ECI DescribeContainerGroups does not guarantee config-order preservation for complex object members, so a convergent member reorder cannot be verified.
// collection-order: bypass containers.environment_vars.field_ref -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.liveness_probe -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.liveness_probe.exec -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.liveness_probe.http_get -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.liveness_probe.tcp_socket -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.readiness_probe -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.readiness_probe.exec -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.readiness_probe.http_get -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.readiness_probe.tcp_socket -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.security_context -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass containers.security_context.capability -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass init_containers.environment_vars.field_ref -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass init_containers.security_context -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
// collection-order: bypass init_containers.security_context.capability -- pre-existing attribute surfaced by the ECI subpackage migration; ECI Create builds this single-value field from the last list member (loop-overwrite), so no convergent member reorder exists.
