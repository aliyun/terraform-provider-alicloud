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
		"alicloud_eci_container_group",
		&resource.Sweeper{
			Name: "alicloud_eci_container_group",
			F:    testSweepEciContainerGroup,
		})
}

func testSweepEciContainerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc`",
	}
	var response map[string]interface{}
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DescribeContainerGroups"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
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
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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

func TestAccAlicloudEciContainerGroup_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEciContainerGroupBasicDependence)
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
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.group.id}",
					"vswitch_id":           "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_registry_credential"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
						"containers.#": "1",
						"volumes.#":    "1",
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
					"restart_policy": "Always",
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
						"containers.#":   "1",
						"restart_policy": "Always",
						"tags.%":         "2",
						"tags.Created":   "TF",
						"tags.For":       "Test",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEciContainerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEciContainerGroupBasicDependence)
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
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.group.id}",
					"vswitch_id":           "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
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
					"auto_match_image_cache": "true",
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
				ImportStateVerifyIgnore: []string{"image_registry_credential", "auto_match_image_cache", "eip_bandwidth", "auto_create_eip"},
			},
		},
	})
}

func TestAccAlicloudEciContainerGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEciContainerGroupBasicDependence)
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
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.group.id}",
					"vswitch_id":           "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
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
				ImportStateVerifyIgnore: []string{"image_registry_credential", "insecure_registry", "plain_http_registry"},
			},
		},
	})
}

func TestAccAlicloudEciContainerGroup_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_container_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEciContainerGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciContainerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEciContainerGroupBasicDependence3)
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
					"container_group_name": strings.ToLower(name),
					"security_group_id":    "${alicloud_security_group.group.id}",
					"vswitch_id":           "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
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
					"auto_match_image_cache": "true",
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
				ImportStateVerifyIgnore: []string{"image_registry_credential", "auto_match_image_cache", "eip_bandwidth", "auto_create_eip", "eip_instance_id"},
			},
		},
	})
}

var AlicloudEciContainerGroupMap = map[string]string{
	"resource_group_id": CHECKSET,
	"restart_policy":    "Always",
	"status":            CHECKSET,
}

func AlicloudEciContainerGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "tf-eci-image-test"
  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
}
`, name)
}

func AlicloudEciContainerGroupBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "tf-eci-image-test"
  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}
`, name)
}
