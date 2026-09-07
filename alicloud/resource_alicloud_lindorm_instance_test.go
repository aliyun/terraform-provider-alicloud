package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
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
		request := map[string]interface{}{
			"InstanceId": id,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete lindorm Instance (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAliCloudLindormInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"eu-central-1"})
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"instance_storage":          "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_engine_specification": "lindorm.c.xlarge-upd",
					"file_engine_node_count":    "3",
					"instance_storage":          "160",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_engine_node_count": "3",
						"instance_storage":       "160",
					}),
				),
			},
			// Field `file_engine_specification` only support lindorm.c.xlarge, so disable this upgrade test case
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"file_engine_specification": "lindorm.c.2xlarge",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"file_engine_specification": "lindorm.c.2xlarge",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_specification": "lindorm.g.2xlarge",
					"time_series_engine_node_count":    "2",
					"instance_storage":                 "320",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_specification": "lindorm.g.2xlarge",
						"time_series_engine_node_count":    "2",
						"instance_storage":                 "320",
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
					"instance_storage":              "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_node_count": "3",
						"instance_storage":              "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"table_engine_specification": "lindorm.g.4xlarge",
					"table_engine_node_count":    "2",
					"instance_storage":           "560",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_specification": "lindorm.g.4xlarge",
						"table_engine_node_count":    "2",
						"instance_storage":           "560",
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
					"instance_storage":        "640",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_node_count": "3",
						"instance_storage":        "640",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_specification": "lindorm.g.2xlarge",
					"search_engine_node_count":    "2",
					"instance_storage":            "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.2xlarge",
						"search_engine_node_count":    "2",
						"instance_storage":            "800",
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
					"instance_storage":         "880",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_node_count": "3",
						"instance_storage":         "880",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lts_node_specification": "lindorm.g.xlarge",
					"lts_node_count":         "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lts_node_specification": "lindorm.g.xlarge",
						"lts_node_count":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lts_node_specification": "lindorm.g.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lts_node_specification": "lindorm.g.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lts_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lts_node_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_engine_specification": "lindorm.g.xlarge",
					"stream_engine_node_count":    "2",
					"instance_storage":            "1040",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_engine_specification": "lindorm.g.xlarge",
						"stream_engine_node_count":    "2",
						"instance_storage":            "1040",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_engine_specification": "lindorm.c.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_engine_specification": "lindorm.c.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_engine_node_count": "3",
					"instance_storage":         "1120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_engine_node_count": "3",
						"instance_storage":         "1120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage": "800",
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
					"instance_storage": "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "1200",
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
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudLindormInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"eu-central-1"})
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudLindormInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.LindormInstanceRegions)
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLindormInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"table_engine_specification": "lindorm.g.2xlarge",
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
					"log_num": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_node_count": "8",
						"log_num":                 "12",
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
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudLindormInstance_basic3(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"eu-central-1"})
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"disk_category":    "local_ssd_pro",
						"payment_type":     "PayAsYouGo",
						"instance_name":    name,
						"core_spec":        "lindorm.i2.4xlarge",
						"instance_storage": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudLindormInstance_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"eu-central-1"})
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":                       "${data.alicloud_vswitches.default.ids[0]}",
					"disk_category":                    "cloud_efficiency",
					"payment_type":                     "Subscription",
					"cold_storage":                     "800",
					"deletion_proection":               "false",
					"duration":                         "1",
					"file_engine_node_count":           "2",
					"file_engine_specification":        "lindorm.c.xlarge",
					"instance_name":                    "${var.name}",
					"instance_storage":                 "800",
					"ip_white_list":                    []string{"10.0.0.0/8", "11.0.0.0/8", "33.0.0.0/8"},
					"lts_node_count":                   "2",
					"lts_node_specification":           "lindorm.g.xlarge",
					"pricing_cycle":                    "Month",
					"search_engine_node_count":         "2",
					"search_engine_specification":      "lindorm.g.xlarge",
					"table_engine_node_count":          "2",
					"table_engine_specification":       "lindorm.g.4xlarge",
					"time_series_engine_node_count":    "2",
					"time_series_engine_specification": "lindorm.g.2xlarge",
					"stream_engine_node_count":         "2",
					"stream_engine_specification":      "lindorm.g.xlarge",
					"vpc_id":                           "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                          "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"auto_renew":                       "true",
					"auto_renew_period":                "1",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "HITS",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":                       CHECKSET,
						"disk_category":                    "cloud_efficiency",
						"payment_type":                     "Subscription",
						"cold_storage":                     "800",
						"deletion_proection":               "false",
						"file_engine_node_count":           "2",
						"file_engine_specification":        "lindorm.c.xlarge",
						"instance_name":                    name,
						"instance_storage":                 "800",
						"ip_white_list.#":                  "3",
						"lts_node_count":                   "2",
						"lts_node_specification":           "lindorm.g.xlarge",
						"pricing_cycle":                    "Month",
						"search_engine_node_count":         "2",
						"search_engine_specification":      "lindorm.g.xlarge",
						"table_engine_node_count":          "2",
						"table_engine_specification":       "lindorm.g.4xlarge",
						"time_series_engine_node_count":    "2",
						"time_series_engine_specification": "lindorm.g.2xlarge",
						"stream_engine_node_count":         "2",
						"stream_engine_specification":      "lindorm.g.xlarge",
						"vpc_id":                           CHECKSET,
						"resource_group_id":                CHECKSET,
						"tags.%":                           "2",
						"tags.Created":                     "TF",
						"tags.For":                         "HITS",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

var AliCloudLindormInstanceMap0 = map[string]string{
	"cold_storage":                  CHECKSET,
	"deletion_proection":            CHECKSET,
	"status":                        CHECKSET,
	"time_series_engine_node_count": CHECKSET,
	"file_engine_node_count":        CHECKSET,
	"lts_node_count":                CHECKSET,
	"search_engine_node_count":      CHECKSET,
	"core_spec":                     "",
	"table_engine_node_count":       CHECKSET,
	"instance_storage":              CHECKSET,
	"zone_id":                       CHECKSET,
	"disk_category":                 "cloud_efficiency",
	"payment_type":                  "PayAsYouGo",
	"vswitch_id":                    CHECKSET,
	"instance_name":                 CHECKSET,
	//"lts_node_specification":        CHECKSET,
	//"stream_engine_specification":   CHECKSET,
	//"file_engine_specification":     CHECKSET,
	//"search_engine_specification":   CHECKSET,
	//"table_engine_specification":    CHECKSET,
	"stream_engine_node_count":    CHECKSET,
	"service_type":                CHECKSET,
	"enabled_file_engine":         CHECKSET,
	"enabled_time_serires_engine": CHECKSET,
	"enabled_table_engine":        CHECKSET,
	"enabled_search_engine":       CHECKSET,
	"enabled_lts_engine":          CHECKSET,
	"enabled_stream_engine":       CHECKSET,
	"arch_version":                CHECKSET,
}

func AliCloudLindormInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

func AliCloudLindormInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
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

const lindormInstanceTestId = "ld-test-instance"

type lindormInstanceTestAPI struct {
	t             *testing.T
	engineType    int
	logNum        int
	rejectUpgrade bool

	describeFailures      int
	describeFailureStatus int
	describeFailureCode   string

	mu            sync.Mutex
	upgrades      []map[string]string
	describeCalls int
}

func (a *lindormInstanceTestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.t.Errorf("parse hitsdb request: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"RequestId": "test-request-id"}
	switch action := r.Form.Get("Action"); action {
	case "UpgradeLindormInstance":
		call := map[string]string{}
		for _, key := range []string{"UpgradeType", "LogNum", "LindormNum", "LindormSpec", "ClusterStorage", "InstanceId"} {
			if value := r.Form.Get(key); value != "" {
				call[key] = value
			}
		}
		a.mu.Lock()
		a.upgrades = append(a.upgrades, call)
		a.mu.Unlock()
		if a.rejectUpgrade {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{
				"Code":      "InvalidParameter",
				"Message":   "mock UpgradeLindormInstance rejected",
				"RequestId": "test-request-id",
			}
		}
	case "GetLindormInstance":
		a.mu.Lock()
		a.describeCalls++
		failRead := a.describeCalls <= a.describeFailures
		a.mu.Unlock()
		if failRead {
			status := a.describeFailureStatus
			if status == 0 {
				status = http.StatusUnavailableForLegalReasons
			}
			code := a.describeFailureCode
			if code == "" {
				code = "OperationFailed"
			}
			w.WriteHeader(status)
			response = map[string]interface{}{
				"Code":      code,
				"Message":   "mock GetLindormInstance failure",
				"RequestId": "test-request-id",
			}
			break
		}
		response = map[string]interface{}{
			"InstanceId":       lindormInstanceTestId,
			"InstanceStatus":   "ACTIVATION",
			"InstanceAlias":    "tf-test-lindorm",
			"PayType":          "POSTPAY",
			"ZoneId":           "cn-hangzhou-h",
			"ServiceType":      "lindorm",
			"LogNum":           a.logNum,
			"LogSingleStorage": 400,
			"InstanceStorage":  1000,
			"EngineType":       a.engineType,
		}
	case "GetInstanceIpWhiteList":
		response = map[string]interface{}{"InstanceId": lindormInstanceTestId, "IpList": []interface{}{"127.0.0.1"}}
	case "GetLindormInstanceEngineInfo":
		response = map[string]interface{}{
			"InstanceId":               lindormInstanceTestId,
			"TableEngineNodeCount":     4,
			"TableEngineSpecification": "lindorm.g.xlarge",
		}
	case "ListTagResources":
		response = map[string]interface{}{"TagResources": []interface{}{}}
	default:
		a.t.Errorf("unexpected hitsdb action %q", action)
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]interface{}{"Code": "UnexpectedAction", "Message": action}
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.t.Errorf("encode hitsdb response: %v", err)
	}
}

func (a *lindormInstanceTestAPI) upgradeCalls() []map[string]string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]map[string]string(nil), a.upgrades...)
}

func (a *lindormInstanceTestAPI) describeCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.describeCalls
}

func lindormInstanceTestClient(t *testing.T, handler http.Handler) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("hitsdb", endpoint)
	return client
}

func lindormInstanceTestData(t *testing.T, state, changes map[string]string) *schema.ResourceData {
	t.Helper()
	attributes := map[string]string{"id": lindormInstanceTestId}
	for key, value := range state {
		attributes[key] = value
	}
	diff := sdktf.NewInstanceDiff()
	for key, value := range changes {
		diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: attributes[key], New: value}
	}
	d, err := schema.InternalMap(resourceAliCloudLindormInstance().Schema).Data(
		&sdktf.InstanceState{ID: lindormInstanceTestId, Attributes: attributes}, diff)
	if err != nil {
		t.Fatalf("build resource data: %v", err)
	}
	return d
}

func lindormInstanceTableEngineState() map[string]string {
	return map[string]string{
		"enabled_table_engine":       "true",
		"table_engine_node_count":    "4",
		"table_engine_specification": "lindorm.g.xlarge",
		"log_num":                    "4",
		"instance_storage":           "1000",
		"zone_id":                    "cn-hangzhou-h",
	}
}

func TestUnitLindormInstanceUpdateLogNum(t *testing.T) {
	t.Run("log_num only is upgraded independently", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 8}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"log_num": "8"})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("log_num-only update failed: %v", err)
		}

		// The pre-fix gate keyed only off table_engine_node_count and
		// table_engine_specification, so it sent no request here at all and this
		// assertion is the counter-example that fails without the log_num gate.
		calls := api.upgradeCalls()
		if len(calls) != 1 {
			t.Fatalf("expected exactly one UpgradeLindormInstance call for a log_num-only change, got %d: %#v", len(calls), calls)
		}
		if got := calls[0]["UpgradeType"]; got != "upgrade-lindorm-core-num" {
			t.Fatalf("expected UpgradeType=upgrade-lindorm-core-num, got %q", got)
		}
		if got := calls[0]["LogNum"]; got != "8" {
			t.Fatalf("expected LogNum=8 (the new value), got %q", got)
		}
		if got := calls[0]["LindormNum"]; got != "4" {
			t.Fatalf("expected LindormNum=4 (the unchanged table node count), got %q", got)
		}
		if got, ok := calls[0]["ClusterStorage"]; ok {
			t.Fatalf("a log_num-only change must not send ClusterStorage, got %q", got)
		}
	})

	t.Run("log_num and table node count change together are sent as one request", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 8}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{
			"log_num":                 "8",
			"table_engine_node_count": "8",
		})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("combined log_num and table node count update failed: %v", err)
		}

		calls := api.upgradeCalls()
		if len(calls) != 1 {
			t.Fatalf("expected exactly one UpgradeLindormInstance call for a combined change, got %d: %#v", len(calls), calls)
		}
		if got := calls[0]["UpgradeType"]; got != "upgrade-lindorm-core-num" {
			t.Fatalf("expected UpgradeType=upgrade-lindorm-core-num, got %q", got)
		}
		if got := calls[0]["LogNum"]; got != "8" {
			t.Fatalf("expected LogNum=8 (the new value), got %q", got)
		}
		if got := calls[0]["LindormNum"]; got != "8" {
			t.Fatalf("expected LindormNum=8 (the new node count), got %q", got)
		}
		if got, ok := calls[0]["ClusterStorage"]; ok {
			t.Fatalf("a combined log_num and node count change must not send ClusterStorage, got %q", got)
		}
	})

	t.Run("log_num and specification change together upgrade both", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 8}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{
			"log_num":                    "8",
			"table_engine_specification": "lindorm.g.2xlarge",
		})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("combined log_num and specification update failed: %v", err)
		}

		calls := api.upgradeCalls()
		if len(calls) != 2 {
			t.Fatalf("expected a spec upgrade followed by a log_num upgrade, got %d call(s): %#v", len(calls), calls)
		}
		if got := calls[0]["UpgradeType"]; got != "upgrade-lindorm-engine" {
			t.Fatalf("expected the first call to be upgrade-lindorm-engine, got %q", got)
		}
		if got := calls[0]["LindormSpec"]; got != "lindorm.g.2xlarge" {
			t.Fatalf("expected LindormSpec=lindorm.g.2xlarge on the spec upgrade, got %q", got)
		}
		if got := calls[0]["ClusterStorage"]; got != "1000" {
			t.Fatalf("expected ClusterStorage=1000 on the spec upgrade, got %q", got)
		}
		for _, key := range []string{"LogNum", "LindormNum"} {
			if got, ok := calls[0][key]; ok {
				t.Fatalf("the spec upgrade must not send %s, got %q", key, got)
			}
		}
		if got := calls[1]["UpgradeType"]; got != "upgrade-lindorm-core-num" {
			t.Fatalf("expected the second call to be upgrade-lindorm-core-num, got %q", got)
		}
		if got := calls[1]["LogNum"]; got != "8" {
			t.Fatalf("expected LogNum=8 on the log_num upgrade, got %q", got)
		}
		if got := calls[1]["LindormNum"]; got != "4" {
			t.Fatalf("expected LindormNum=4 (the unchanged table node count) on the log_num upgrade, got %q", got)
		}
		if got, ok := calls[1]["ClusterStorage"]; ok {
			t.Fatalf("the log_num upgrade must not send ClusterStorage, got %q", got)
		}
	})

	t.Run("log_num only does not open a disabled table engine", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 0, logNum: 8}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, map[string]string{
			"enabled_table_engine":    "false",
			"table_engine_node_count": "0",
			"log_num":                 "4",
			"instance_storage":        "1000",
			"zone_id":                 "cn-hangzhou-h",
		}, map[string]string{"log_num": "8"})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("log_num-only update on a disabled table engine failed: %v", err)
		}

		if calls := api.upgradeCalls(); len(calls) != 0 {
			t.Fatalf("a log_num-only change must not trigger any upgrade while the table engine is disabled, got %#v", calls)
		}
	})

	t.Run("table node count change keeps its request shape", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"table_engine_node_count": "8"})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("table node count update failed: %v", err)
		}

		calls := api.upgradeCalls()
		if len(calls) != 1 {
			t.Fatalf("expected exactly one UpgradeLindormInstance call, got %d: %#v", len(calls), calls)
		}
		if got := calls[0]["UpgradeType"]; got != "upgrade-lindorm-core-num" {
			t.Fatalf("expected UpgradeType=upgrade-lindorm-core-num, got %q", got)
		}
		if got := calls[0]["LindormNum"]; got != "8" {
			t.Fatalf("expected LindormNum=8 (the new node count), got %q", got)
		}
		if got := calls[0]["LogNum"]; got != "" {
			t.Fatalf("a node count change must not send LogNum, got %q", got)
		}
		if got := calls[0]["ClusterStorage"]; got != "1000" {
			t.Fatalf("expected ClusterStorage=1000, got %q", got)
		}
	})

	t.Run("table engine spec change keeps its request shape", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"table_engine_specification": "lindorm.g.2xlarge"})

		if err := resourceAliCloudLindormInstanceUpdate(d, client); err != nil {
			t.Fatalf("table engine spec update failed: %v", err)
		}

		calls := api.upgradeCalls()
		if len(calls) != 1 {
			t.Fatalf("expected exactly one UpgradeLindormInstance call, got %d: %#v", len(calls), calls)
		}
		if got := calls[0]["UpgradeType"]; got != "upgrade-lindorm-engine" {
			t.Fatalf("expected UpgradeType=upgrade-lindorm-engine, got %q", got)
		}
		if got := calls[0]["LindormSpec"]; got != "lindorm.g.2xlarge" {
			t.Fatalf("expected LindormSpec=lindorm.g.2xlarge, got %q", got)
		}
		if got := calls[0]["LindormNum"]; got != "" {
			t.Fatalf("a spec change must not send LindormNum, got %q", got)
		}
		if got := calls[0]["LogNum"]; got != "" {
			t.Fatalf("a spec change must not send LogNum, got %q", got)
		}
	})
}

func TestUnitLindormInstanceUpdateLogNumFailure(t *testing.T) {
	api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4, rejectUpgrade: true}
	client := lindormInstanceTestClient(t, api)
	d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"log_num": "8"})

	err := resourceAliCloudLindormInstanceUpdate(d, client)
	if err == nil {
		t.Fatal("a rejected UpgradeLindormInstance must fail the update")
	}
	if !strings.Contains(err.Error(), "mock UpgradeLindormInstance rejected") {
		t.Fatalf("the cloud error must be returned as-is, got: %v", err)
	}
	if calls := api.upgradeCalls(); len(calls) != 1 {
		t.Fatalf("expected exactly one UpgradeLindormInstance attempt, got %d: %#v", len(calls), calls)
	}
	if got := api.describeCount(); got != 0 {
		t.Fatalf("a rejected upgrade must not poll or re-read the instance, got %d GetLindormInstance calls", got)
	}
}

func TestUnitLindormInstanceStateRefreshPostCreateReadTolerance(t *testing.T) {
	t.Run("freshly created instance read failure is treated as provisioning", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4, describeFailures: 2}
		client := lindormInstanceTestClient(t, api)
		service := HitsdbService{client}
		var tolerance LindormPostCreateReadTolerance
		refresh := service.LindormInstanceStateRefreshFuncWithPostCreateReadTolerance(lindormInstanceTestId, []string{}, time.Now().Add(time.Minute), &tolerance)

		for i := 0; i < 2; i++ {
			_, state, err := refresh()
			if err != nil {
				t.Fatalf("call %d: a transient post-create read failure must not fail the wait, got: %v", i+1, err)
			}
			if state == "ACTIVATION" {
				t.Fatalf("call %d: an unreadable instance must not report ACTIVATION", i+1)
			}
		}
		if lastTolerated := tolerance.last(); lastTolerated == nil || !strings.Contains(lastTolerated.Error(), "OperationFailed") {
			t.Fatalf("the tolerated error must be recorded for diagnostics, got: %v", lastTolerated)
		}

		_, state, err := refresh()
		if err != nil {
			t.Fatalf("call 3: the refresh must recover once the instance is readable, got: %v", err)
		}
		if state != "ACTIVATION" {
			t.Fatalf("call 3: expected state ACTIVATION once readable, got %q", state)
		}
	})

	t.Run("the same failure is returned after the tolerance window", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4, describeFailures: 5}
		client := lindormInstanceTestClient(t, api)
		service := HitsdbService{client}
		refresh := service.LindormInstanceStateRefreshFuncWithPostCreateReadTolerance(lindormInstanceTestId, []string{}, time.Now().Add(-time.Second), nil)

		_, _, err := refresh()
		if err == nil {
			t.Fatal("a read failure after the tolerance window must be returned")
		}
		if !strings.Contains(err.Error(), "OperationFailed") {
			t.Fatalf("the original cloud error must be returned after the window, got: %v", err)
		}
	})

	t.Run("other errors are not tolerated", func(t *testing.T) {
		api := &lindormInstanceTestAPI{
			t:                     t,
			engineType:            4,
			logNum:                4,
			describeFailures:      5,
			describeFailureStatus: http.StatusBadRequest,
			describeFailureCode:   "InvalidParameter",
		}
		client := lindormInstanceTestClient(t, api)
		service := HitsdbService{client}
		refresh := service.LindormInstanceStateRefreshFuncWithPostCreateReadTolerance(lindormInstanceTestId, []string{}, time.Now().Add(time.Minute), nil)

		_, _, err := refresh()
		if err == nil {
			t.Fatal("a non-transient read error must propagate even inside the tolerance window")
		}
		if !strings.Contains(err.Error(), "InvalidParameter") {
			t.Fatalf("expected the original error, got: %v", err)
		}
	})
}

func TestUnitLindormPostCreateReadToleranceForCreateTimeout(t *testing.T) {
	cases := []struct {
		name          string
		createTimeout time.Duration
		expected      time.Duration
	}{
		{"default create timeout keeps the upper bound", 90 * time.Minute, 20 * time.Minute},
		{"timeout above the upper bound keeps the upper bound", 30 * time.Minute, 20 * time.Minute},
		{"timeout just above the upper bound keeps room for the original error", 21 * time.Minute, 15*time.Minute + 45*time.Second},
		{"timeout equal to the upper bound leaves room for the original error", 20 * time.Minute, 15 * time.Minute},
		{"short timeout is shortened to three quarters", 12 * time.Minute, 9 * time.Minute},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lindormPostCreateReadToleranceFor(c.createTimeout); got != c.expected {
				t.Fatalf("create timeout %s: expected tolerance %s, got %s", c.createTimeout, c.expected, got)
			}
		})
	}
}
