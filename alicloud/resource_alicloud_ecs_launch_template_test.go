package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_launch_template", &resource.Sweeper{
		Name: "alicloud_ecs_launch_template",
		F:    testAliCloudEcsLaunchTemplate,
	})
}

func testAliCloudEcsLaunchTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "DescribeLaunchTemplates"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_launch_templates", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.LaunchTemplateSets.LaunchTemplateSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LaunchTemplateSets.LaunchTemplateSet", response)
		}

		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["LaunchTemplateName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Launch Template: %s", item["LaunchTemplateName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteLaunchTemplateVersion"
			request := map[string]interface{}{
				"LaunchTemplateId": item["LaunchTemplateId"],
				"RegionId":         client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Launch Template (%s): %s", item["LaunchTemplateName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros Template have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Launch Template success: %s ", item["LaunchTemplateName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudECSLaunchTemplateBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_launch_template.default"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLaunchTemplateBasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence)
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
					"name":                          name,
					"description":                   name,
					"image_id":                      "${data.alicloud_images.default.images.0.id}",
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 "${data.alicloud_instance_types.default.instance_types.0.id}",
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "optimized",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_ids":            []string{"${alicloud_security_group.default.id}"},
					"system_disk": []map[string]interface{}{
						{
							"category":             "cloud_ssd",
							"description":          name,
							"name":                 name,
							"size":                 "40",
							"delete_with_instance": "false",
							"encrypted":            "false",
						},
					},
					"resource_group_id": "rg-zkdfjahg9zxncv0",
					"user_data":         "xxxxxxx",
					"vswitch_id":        "${alicloud_vswitch.shareVswitch1.id}",
					"vpc_id":            "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":           "cn-hangzhou-i",

					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"template_tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"data_disks": []map[string]string{
						{
							"name":                 "disk1",
							"description":          "test1",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
						{
							"name":                 "disk2",
							"description":          "test2",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"description":            name,
						"host_name":              name,
						"instance_name":          name,
						"key_pair_name":          name,
						"ram_role_name":          name,
						"default_version_number": "1",
						"latest_version_number":  "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk": []map[string]interface{}{
						{
							"category":             "cloud_ssd",
							"description":          name + "Update",
							"name":                 name + "Update",
							"size":                 "50",
							"delete_with_instance": "true",
							"encrypted":            "true",
							"kms_key_id":           "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk.#":         "1",
						"latest_version_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk": []map[string]interface{}{
						{
							"encrypted": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk.#":         "1",
						"latest_version_number": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"io_optimized": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"io_optimized":          "none",
						"latest_version_number": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + ":_.-,][[" + name + "]" + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":         name + ":_.-,][[" + name + "]" + "_change",
						"latest_version_number": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "PayByTraffic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type":  "PayByTraffic",
						"latest_version_number": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]string{
						{
							"name":                 "disk1update",
							"description":          "test1 update",
							"delete_with_instance": "true",
							"category":             "cloud_ssd",
							"encrypted":            "true",
							"performance_level":    "PL1",
							"size":                 "25",
						},
						{
							"name":                 "disk2update",
							"description":          "test2 update",
							"delete_with_instance": "true",
							"category":             "cloud_ssd",
							"encrypted":            "true",
							"performance_level":    "PL1",
							"size":                 "25",
							"kms_key_id":           "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#":          "2",
						"latest_version_number": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":         name + "_change",
						"latest_version_number": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#":  "1",
						"latest_version_number": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "ecs.g6.xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":         "ecs.g6.xlarge",
						"latest_version_number": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":         "NoSpot",
						"latest_version_number": "11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": "dGhpcyBpcyBhIGV4YW1wbGU=",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":             CHECKSET,
						"latest_version_number": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id": "${alicloud_vswitch.shareVswitch2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":            CHECKSET,
						"latest_version_number": "13",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type": "classic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":          "classic",
						"latest_version_number": "14",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_price_limit": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_price_limit":      "1",
						"latest_version_number": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ram_role_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_role_name":         name + "_change",
						"latest_version_number": "16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name":             name + "_change",
						"latest_version_number": "17",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "rg-zkdfjahg9zxxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":     CHECKSET,
						"latest_version_number": "18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name + "_change",
						"latest_version_number": "19",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "cn-hangzhou-f",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":               "cn-hangzhou-f",
						"latest_version_number": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_enhancement_strategy": "Deactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_enhancement_strategy": "Deactive",
						"latest_version_number":         "21",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "5",
						"latest_version_number":      "22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "6",
						"latest_version_number":     "23",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tag1": "tag1",
						"tag2": "tag2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":                "2",
						"tags.tag1":             "tag1",
						"tags.tag2":             "tag2",
						"latest_version_number": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello",
							"primary_ip":        "10.0.0.6",
							"security_group_id": "xxxxx",
							"vswitch_id":        "xxxxx",
							"delete_on_release": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#":  "1",
						"latest_version_number": "25",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello",
							"primary_ip":        "10.0.0.6",
							"security_group_id": "xxxxx",
							"vswitch_id":        "xxxxx",
							"delete_on_release": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#":  "1",
						"latest_version_number": "26",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_version_number": "26",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_version_number": "26",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_version_number":        REMOVEKEY,
					"update_default_version_number": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_version_number": "26",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_options": []map[string]string{
						{
							"login_as_non_root": "true",
						},
					}}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_options.#":        "1",
						"default_version_number": "27",
						"latest_version_number":  "27",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"description":                   name,
					"host_name":                     name,
					"instance_name":                 name,
					"instance_type":                 "ecs.g6.large",
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "none",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_ids":            []string{"${alicloud_security_group.default.id}"},
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"user_data":                     "xxxxxxx",
					"vswitch_id":                    "${alicloud_vswitch.shareVswitch1.id}",
					"vpc_id":                        "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":                       "cn-hangzhou-i",

					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": CHECKSET,
							"vswitch_id":        CHECKSET,
							"delete_on_release": "true",
						},
					},
					"system_disk": []map[string]interface{}{
						{
							"category":             "cloud_ssd",
							"description":          name,
							"name":                 name,
							"size":                 "40",
							"delete_with_instance": "false",
							"encrypted":            "false",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":                 "disk1",
							"description":          "test1",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
						{
							"name":                 "disk2",
							"description":          "test2",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"description":                   name,
						"host_name":                     name,
						"instance_name":                 name,
						"key_pair_name":                 name,
						"ram_role_name":                 name,
						"image_id":                      CHECKSET,
						"instance_type":                 CHECKSET,
						"internet_charge_type":          "PayByBandwidth",
						"internet_max_bandwidth_in":     "5",
						"internet_max_bandwidth_out":    "0",
						"io_optimized":                  "none",
						"network_type":                  "vpc",
						"security_enhancement_strategy": "Active",
						"spot_price_limit":              "5",
						"spot_strategy":                 "SpotWithPriceLimit",
						"security_group_ids.#":          "1",
						"resource_group_id":             CHECKSET,
						"user_data":                     "xxxxxxx",
						"vswitch_id":                    CHECKSET,
						"vpc_id":                        CHECKSET,
						"zone_id":                       CHECKSET,
						"tags.%":                        "2",
						"tags.tag1":                     "hello",
						"tags.tag2":                     "world",
						"system_disk.#":                 "1",
						"network_interfaces.#":          "1",
						"data_disks.#":                  "2",
						"latest_version_number":         "28",
						"default_version_number":        "28",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudECSLaunchTemplateBasic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_launch_template.default"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap1)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccLaunchTemplateBasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence1)
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
					"launch_template_name":          name,
					"description":                   name,
					"image_id":                      "${data.alicloud_images.default.images.0.id}",
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name + ":_.-,][[" + name + "]",
					"instance_type":                 "${data.alicloud_instance_types.default.instance_types.0.id}",
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "5",
					"io_optimized":                  "optimized",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_ids":            []string{"${alicloud_security_group.default.id}"},
					"auto_release_time":             time.Now().Add(10 * time.Hour).Format("2021-12-30T12:05:05Z"),
					"deployment_set_id":             "${alicloud_ecs_deployment_set.default.id}",
					"enable_vm_os_config":           "false",
					"image_owner_alias":             "system",
					"password_inherit":              "false",
					"period":                        "1",
					"private_ip_address":            "172.16.0.10",
					"template_resource_group_id":    "rg-zkdfjahg9zxncv0",
					"version_description":           name,
					"system_disk_category":          "cloud_ssd",
					"system_disk_description":       name,
					"system_disk_name":              name,
					"system_disk_size":              "40",
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"userdata":                      "xxxxxxx",
					"vswitch_id":                    "${alicloud_vswitch.shareVswitch1.id}",
					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"template_tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
							"delete_on_release": "true",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":                 "disk1",
							"description":          "test1",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
						{
							"name":                 "disk2",
							"description":          "test2",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
						},
						{
							"name":              "disk2",
							"description":       "test2",
							"category":          "cloud",
							"performance_level": "PL0",
						},
						{
							"delete_with_instance": "true",
							"encrypted":            "false",
							"size":                 "20",
						},
					},
					"image_options": []map[string]string{
						{
							"login_as_non_root": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_name":   name,
						"description":            name,
						"host_name":              name,
						"instance_name":          name + ":_.-,][[" + name + "]",
						"key_pair_name":          name,
						"ram_role_name":          name,
						"auto_release_time":      CHECKSET,
						"deployment_set_id":      CHECKSET,
						"enable_vm_os_config":    CHECKSET,
						"image_owner_alias":      CHECKSET,
						"period":                 "1",
						"private_ip_address":     CHECKSET,
						"version_description":    name,
						"data_disks.#":           "4",
						"network_interfaces.#":   "1",
						"latest_version_number":  "1",
						"default_version_number": "1",
						"image_options.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]string{
						{
							"name":                 "disk1_update",
							"description":          "test1",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
							"device":               "/dev/xvdb",
						},
						{
							"name":                 "disk2",
							"description":          "test2_update",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
							"device":               "/dev/xvdc",
						},
						{
							"name":              "disk2",
							"description":       "test2",
							"category":          "cloud",
							"performance_level": "PL1",
							"device":            "/dev/xvdd",
						},
						{
							"delete_with_instance": "false",
							"encrypted":            "true",
							"size":                 "30",
							"device":               "/dev/xvde",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#":           "4",
						"latest_version_number":  "2",
						"default_version_number": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_resource_group_id"},
			},
		},
	})
}

func TestAccAliCloudECSLaunchTemplateBasic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_launch_template.default"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap1)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccLaunchTemplateBasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence1)
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
					"auto_renew":                    "true",
					"auto_renew_period":             "1",
					"launch_template_name":          name,
					"description":                   name,
					"image_id":                      "${data.alicloud_images.default.images.0.id}",
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 "${data.alicloud_instance_types.default.instance_types.0.id}",
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "5",
					"io_optimized":                  "optimized",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_ids":            []string{"${alicloud_security_group.default.id}"},
					"auto_release_time":             time.Now().Add(10 * time.Hour).Format("2021-12-30T12:05:05Z"),
					"deployment_set_id":             "${alicloud_ecs_deployment_set.default.id}",
					"enable_vm_os_config":           "false",
					"image_owner_alias":             "system",
					"password_inherit":              "false",
					"period":                        "1",
					"period_unit":                   "Month",
					"private_ip_address":            "172.16.0.10",
					"template_resource_group_id":    "rg-zkdfjahg9zxncv0",
					"version_description":           name,
					"system_disk_category":          "cloud_ssd",
					"system_disk_description":       name,
					"system_disk_name":              name,
					"system_disk_size":              "40",
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"userdata":                      "xxxxxxx",
					"http_endpoint":                 "enabled",
					"http_tokens":                   "optional",
					"http_put_response_hop_limit":   "3",
					"vswitch_id":                    "${alicloud_vswitch.shareVswitch1.id}",
					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"template_tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
							"delete_on_release": "true",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":                 "disk1",
							"description":          "test1",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
							"device":               "/dev/xvdb",
						},
						{
							"name":                 "disk2",
							"description":          "test2",
							"delete_with_instance": "true",
							"category":             "cloud",
							"encrypted":            "false",
							"performance_level":    "PL0",
							"size":                 "20",
							"device":               "/dev/xvdc",
						},
						{
							"name":              "disk2",
							"description":       "test2",
							"category":          "cloud",
							"performance_level": "PL0",
							"device":            "/dev/xvdd",
						},
						{
							"delete_with_instance": "true",
							"encrypted":            "false",
							"size":                 "20",
							"device":               "/dev/xvde",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_name":   name,
						"description":            name,
						"host_name":              name,
						"instance_name":          name,
						"key_pair_name":          name,
						"ram_role_name":          name,
						"auto_release_time":      CHECKSET,
						"deployment_set_id":      CHECKSET,
						"enable_vm_os_config":    CHECKSET,
						"image_owner_alias":      CHECKSET,
						"period":                 "1",
						"private_ip_address":     CHECKSET,
						"version_description":    name,
						"data_disks.#":           "4",
						"network_interfaces.#":   "1",
						"latest_version_number":  "1",
						"default_version_number": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_resource_group_id"},
			},
		},
	})
}

func TestAccAliCloudECSLaunchTemplateMulti(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_launch_template.default.4"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccLaunchTemplateMulti%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence)
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
					"count":                         "5",
					"name":                          name + "${count.index}",
					"description":                   name,
					"image_id":                      "${data.alicloud_images.default.images.0.id}",
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 "${data.alicloud_instance_types.default.instance_types.0.id}",
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "optimized",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_ids":            []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.update.id}"},
					"system_disk": []map[string]interface{}{
						{
							"category":             "cloud_ssd",
							"description":          name,
							"name":                 name,
							"size":                 "40",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"kms_key_id":           "${alicloud_kms_key.default.id}",
						},
					},
					"resource_group_id": "rg-zkdfjahg9zxncv0",
					"user_data":         "xxxxxxxxxxxxxx",
					"vswitch_id":        "${alicloud_vswitch.shareVswitch1.id}",
					"vpc_id":            "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":           "cn-beijing-a",
					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"description": "test1",
						},
						{
							"name":        "disk2",
							"description": "test2",
							"kms_key_id":  "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name + "4",
						"description":            name,
						"host_name":              name,
						"instance_name":          name,
						"key_pair_name":          name,
						"ram_role_name":          name,
						"security_group_ids.#":   "2",
						"network_interfaces.#":   "1",
						"latest_version_number":  "1",
						"default_version_number": "1",
					}),
				),
			},
		},
	})
}

var testAccLaunchTemplateCheckMap = map[string]string{
	"image_id":                      CHECKSET,
	"instance_charge_type":          "PrePaid",
	"instance_type":                 CHECKSET,
	"internet_charge_type":          "PayByBandwidth",
	"internet_max_bandwidth_in":     "5",
	"internet_max_bandwidth_out":    "0",
	"io_optimized":                  "optimized",
	"network_type":                  "vpc",
	"security_enhancement_strategy": "Active",
	"spot_price_limit":              "5",
	"spot_strategy":                 "SpotWithPriceLimit",
	"security_group_ids.#":          "1",
	"system_disk.#":                 "1",
	"resource_group_id":             CHECKSET,
	"userdata":                      CHECKSET,
	"vswitch_id":                    CHECKSET,
	"vpc_id":                        CHECKSET,
	"zone_id":                       CHECKSET,
	"data_disks.#":                  "2",
}

var testAccLaunchTemplateCheckMap1 = map[string]string{
	"image_id":                      CHECKSET,
	"instance_charge_type":          "PrePaid",
	"instance_type":                 CHECKSET,
	"internet_charge_type":          "PayByBandwidth",
	"internet_max_bandwidth_in":     "5",
	"internet_max_bandwidth_out":    "5",
	"io_optimized":                  "optimized",
	"network_type":                  "vpc",
	"security_enhancement_strategy": "Active",
	"spot_price_limit":              "5",
	"spot_strategy":                 "SpotWithPriceLimit",
	"security_group_ids.#":          "1",
	"system_disk.#":                 "1",
	"resource_group_id":             CHECKSET,
	"userdata":                      CHECKSET,
	"vswitch_id":                    CHECKSET,
	"network_interfaces.#":          "1",
	"data_disks.#":                  "2",
}

func resourceLaunchTemplateConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

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
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
 vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "shareVswitch1" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "shareVswitch2" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_security_group" "update" {
  name   = "${var.name}1"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_kms_key" "default" {
	description             = var.name
	pending_window_in_days  = "7"
	status               = "Enabled"
}
`, name)
}

func resourceLaunchTemplateConfigDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

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
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
 vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "shareVswitch1" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "shareVswitch2" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_security_group" "update" {
  name   = "${var.name}1"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
  description         = var.name
}
`, name)
}
