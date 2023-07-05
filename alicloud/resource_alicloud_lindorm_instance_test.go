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
	resource.AddTestSweepers("alicloud_lindorm_instance", &resource.Sweeper{
		Name: "alicloud_lindorm_instance",
		F:    testSweepLindormInstances,
	})
}

func testSweepLindormInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "GetLindormInstanceList"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	lindormInstanceIds := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_lindorm_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["InstanceAlias"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Lindorm Instance: %v (%v)", item["InstanceAlias"], item["InstanceId"])
					continue
				}
			}
			lindormInstanceIds = append(lindormInstanceIds, fmt.Sprint(item["InstanceId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, id := range lindormInstanceIds {
		log.Printf("[INFO] Deleting Lindorm Instance: %s", id)
		action := "ReleaseLindormInstance"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"InstanceId": id,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete lindorm Instance (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAlicloudLindormInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"118.118.118.118"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"117.117.117.117", "116.116.116.116"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "2400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "2400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_proection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_proection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":      name,
					"deletion_proection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":      name,
						"deletion_proection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_serires_engine_specification": "lindorm.g.2xlarge",
					"time_series_engine_node_count":     "2",
					"instance_storage":                  "4320",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_serires_engine_specification": "lindorm.g.2xlarge",
						"time_series_engine_node_count":     "2",
						"instance_storage":                  "4320",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_serires_engine_specification": "lindorm.g.4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_serires_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_node_count": "3",
					"instance_storage":              "5440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_node_count": "3",
						"instance_storage":              "5440",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_2"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":              "cloud_efficiency",
					"payment_type":               "PayAsYouGo",
					"zone_id":                    "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":              "${var.name}",
					"table_engine_specification": "lindorm.g.4xlarge",
					"table_engine_node_count":    "2",
					"instance_storage":           "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":              "cloud_efficiency",
						"payment_type":               "PayAsYouGo",
						"instance_name":              name,
						"table_engine_specification": "lindorm.g.4xlarge",
						"table_engine_node_count":    "2",
						"instance_storage":           "1920",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"table_engine_specification": "lindorm.c.8xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_specification": "lindorm.c.8xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"table_engine_node_count": "3",
					"instance_storage":        "3200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_node_count": "3",
						"instance_storage":        "3200",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_3"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_specification": "lindorm.g.2xlarge",
					"search_engine_node_count":    "2",
					"instance_storage":            "4320",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.2xlarge",
						"search_engine_node_count":    "2",
						"instance_storage":            "4320",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_specification": "lindorm.g.4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_node_count": "3",
					"instance_storage":         "5440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_node_count": "3",
						"instance_storage":         "5440",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_essd_pl0",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "HITS",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_essd_pl0",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
						"resource_group_id":         CHECKSET,
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "HITS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "HITS Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "HITS Update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_5"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_specification": "lindorm.g.2xlarge",
					"time_series_engine_node_count":    "2",
					"instance_storage":                 "4320",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_specification": "lindorm.g.2xlarge",
						"time_series_engine_node_count":    "2",
						"instance_storage":                 "4320",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_specification": "lindorm.g.4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_node_count": "3",
					"instance_storage":              "5440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_node_count": "3",
						"instance_storage":              "5440",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic6(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.LindormInstanceRegions)
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence1)
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
					"arch_version":               "2.0",
					"disk_category":              "cloud_efficiency",
					"payment_type":               "PayAsYouGo",
					"zone_id":                    "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids[0]}",
					"primary_zone_id":            "${data.alicloud_zones.default.zones.0.id}",
					"primary_vswitch_id":         "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":              name,
					"table_engine_node_count":    "4",
					"table_engine_specification": "lindorm.g.4xlarge",
					//"core_spec":               "lindorm.i2.xlarge",
					"log_num":                "4",
					"log_single_storage":     "400",
					"arbiter_zone_id":        "${data.alicloud_zones.default.zones.1.id}",
					"arbiter_vswitch_id":     "${data.alicloud_vswitches.arbitervswitchid.ids[0]}",
					"standby_zone_id":        "${data.alicloud_zones.default.zones.2.id}",
					"log_spec":               "lindorm.sn1.large",
					"log_disk_category":      "cloud_efficiency",
					"core_single_storage":    "400",
					"standby_vswitch_id":     "${data.alicloud_vswitches.standbyvswitchid.ids[0]}",
					"multi_zone_combination": "ap-southeast-1-abc-aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":           "cloud_efficiency",
						"payment_type":            "PayAsYouGo",
						"zone_id":                 CHECKSET,
						"vswitch_id":              CHECKSET,
						"primary_zone_id":         CHECKSET,
						"primary_vswitch_id":      CHECKSET,
						"instance_name":           name,
						"table_engine_node_count": "4",
						"instance_storage":        CHECKSET,
						"core_spec":               CHECKSET,
						"log_num":                 "4",
						"log_single_storage":      "400",
						"arbiter_zone_id":         CHECKSET,
						"multi_zone_combination":  "ap-southeast-1-abc-aliyun",
						"arbiter_vswitch_id":      CHECKSET,
						"standby_zone_id":         CHECKSET,
						"log_spec":                "lindorm.sn1.large",
						"log_disk_category":       "cloud_efficiency",
						"core_single_storage":     "400",
						"standby_vswitch_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_single_storage": "440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_single_storage": "440",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_single_storage": "440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_single_storage": "440",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"table_engine_node_count": "8",
					"log_num":                 "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_node_count": "8",
						"log_num":                 "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_spec": "lindorm.sn1.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_spec": "lindorm.sn1.2xlarge",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration", "arch_version"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic7(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":           "local_ssd_pro",
					"payment_type":            "PayAsYouGo",
					"zone_id":                 "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":           "${var.name}",
					"core_spec":               "lindorm.i2.4xlarge",
					"table_engine_node_count": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "local_ssd_pro",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"core_spec":                 "lindorm.i2.4xlarge",
						"file_engine_specification": CHECKSET,
						"instance_storage":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"10.0.0.0/8", "11.0.0.0/8", "33.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"117.117.117.117", "116.116.116.116"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_proection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_proection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":      name,
					"deletion_proection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":      name,
						"deletion_proection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_VpcId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
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
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids[0]}",
					"vpc_id":                    "${data.alicloud_vpcs.default.ids.0}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "1920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "1920",
						"vpc_id":                    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "pricing_cycle", "duration"},
			},
		},
	})
}

var AlicloudLindormInstanceMap0 = map[string]string{
	"cold_storage":                  CHECKSET,
	"search_engine_specification":   CHECKSET,
	"duration":                      NOSET,
	"deletion_proection":            CHECKSET,
	"file_engine_specification":     CHECKSET,
	"status":                        CHECKSET,
	"core_num":                      NOSET,
	"phoenix_node_count":            CHECKSET,
	"phoenix_node_specification":    CHECKSET,
	"group_name":                    NOSET,
	"lts_node_specification":        CHECKSET,
	"time_series_engine_node_count": CHECKSET,
	"file_engine_node_count":        CHECKSET,
	"lts_node_count":                CHECKSET,
	"search_engine_node_count":      CHECKSET,
	"core_spec":                     "",
	"pricing_cycle":                 NOSET,
	"table_engine_node_count":       CHECKSET,
	"instance_storage":              CHECKSET,
	"zone_id":                       CHECKSET,
	"disk_category":                 "cloud_efficiency",
	"payment_type":                  "PayAsYouGo",
	"vswitch_id":                    CHECKSET,
	"instance_name":                 CHECKSET,
	"table_engine_specification":    CHECKSET,
	"service_type":                  CHECKSET,
	"enabled_file_engine":           CHECKSET,
	"enabled_time_serires_engine":   CHECKSET,
	"enabled_table_engine":          CHECKSET,
	"enabled_search_engine":         CHECKSET,
	"enabled_lts_engine":            CHECKSET,
	"arch_version":                  CHECKSET,
}

func AlicloudLindormInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	
	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

func AlicloudLindormInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	
	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

    data "alicloud_vswitches" "arbitervswitchid" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.1.id
	}

    data "alicloud_vswitches" "standbyvswitchid" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.2.id
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}
