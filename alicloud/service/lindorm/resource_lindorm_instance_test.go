package lindorm_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/lindorm"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lindormTestClientForRegion delegates the package-level sweeper client to the
// shared sweep.SharedClientForRegion, which isolates the per-region endpoint
// cache and keeps signVersion, AK/STS, AccountId and AccountType.
func lindormTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return rawClient.(*connectivity.AliyunClient), nil
}

func init() {
	resource.AddTestSweepers("alicloud_lindorm_instance", &resource.Sweeper{
		Name: "alicloud_lindorm_instance",
		F:    testSweepLindormInstances,
	})
}

func testSweepLindormInstances(region string) error {
	rawClient, err := lindormTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "GetLindormInstanceList"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	lindormInstanceIds := make([]string, 0)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_lindorm_instances", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") != "true" {
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
		if len(result) < helper.PageSizeLarge {
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
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
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return lindorm.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLindormInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   lindormHangzhouSaleZoneId,
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
					"file_engine_node_count": "3",
					"instance_storage":       "160",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_engine_node_count": "3",
						"instance_storage":       "160",
					}),
				),
			},
			// The product team no longer recommends the time-series engine, but
			// the repository attribute-coverage gate still requires an update for
			// every schema attribute, so it is opened and resized exactly once.
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
				// The deprecated time_series alias must still drive a real
				// upgrade, not merely exist in the schema. The two fields
				// conflict, so the new field is removed from the config first.
				Config: testAccConfig(map[string]interface{}{
					"time_series_engine_specification":  acctest.REMOVEKEY,
					"time_serires_engine_specification": "lindorm.g.4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_serires_engine_specification": "lindorm.g.4xlarge",
						// The lookup refreshes both fields, so the accumulated
						// assertion on the new field has to follow the upgrade.
						"time_series_engine_specification": "lindorm.g.4xlarge",
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
					"search_engine_specification": "lindorm.g.xlarge",
					"search_engine_node_count":    "2",
					"instance_storage":            "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.xlarge",
						"search_engine_node_count":    "2",
						"instance_storage":            "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lts_node_specification": "lindorm.g.xlarge",
					"lts_node_count":         "1",
					"instance_storage":       "880",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lts_node_specification": "lindorm.g.xlarge",
						"lts_node_count":         "1",
						"instance_storage":       "880",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_engine_specification": "lindorm.g.xlarge",
					"stream_engine_node_count":    "2",
					"instance_storage":            "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_engine_specification": "lindorm.g.xlarge",
						"stream_engine_node_count":    "2",
						"instance_storage":            "1200",
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
					"ip_white_list": []string{"115.115.115.115", "116.116.116.116", "117.117.117.117"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "3",
					}),
				),
			},
			{
				// Replacing the set (rather than only appending) keeps the
				// removal semantics of the original case covered.
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// file_engine_specification only has one valid value, so no step can modify it; skip it here
				// instead of mutating the attribute to an unsupported value. The deprecated
				// time_serires_engine_specification alias is not written back by Read, so it is
				// ignored here while the regression still configures it in the step above.
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration", "auto_renew", "auto_renew_period", "file_engine_specification", "time_serires_engine_specification"},
			},
		},
	})
}

func TestAccAliCloudLindormInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1})
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return lindorm.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLindormInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudLindormInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"arch_version":               "2.0",
					"disk_category":              "cloud_efficiency",
					"payment_type":               "PayAsYouGo",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                    lindormMultiZonePrimaryZoneId,
					"vswitch_id":                 "${data.alicloud_vswitches.primary.ids[0]}",
					"primary_zone_id":            lindormMultiZonePrimaryZoneId,
					"primary_vswitch_id":         "${data.alicloud_vswitches.primary.ids[0]}",
					"instance_name":              name,
					"table_engine_node_count":    "4",
					"table_engine_specification": "lindorm.g.2xlarge",
					//"core_spec":               "lindorm.i2.xlarge",
					"log_num":                "4",
					"log_single_storage":     "400",
					"arbiter_zone_id":        lindormMultiZoneArbiterZoneId,
					"arbiter_vswitch_id":     "${data.alicloud_vswitches.arbiter.ids[0]}",
					"standby_zone_id":        lindormMultiZoneStandbyZoneId,
					"log_spec":               "lindorm.sn1.large",
					"log_disk_category":      "cloud_efficiency",
					"core_single_storage":    "400",
					"standby_vswitch_id":     "${data.alicloud_vswitches.standby.ids[0]}",
					"multi_zone_combination": lindormMultiZoneCombination,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":           "cloud_efficiency",
						"payment_type":            "PayAsYouGo",
						"zone_id":                 acctest.CHECKSET,
						"vswitch_id":              acctest.CHECKSET,
						"primary_zone_id":         acctest.CHECKSET,
						"primary_vswitch_id":      acctest.CHECKSET,
						"instance_name":           name,
						"table_engine_node_count": "4",
						"instance_storage":        acctest.CHECKSET,
						"core_spec":               acctest.CHECKSET,
						"log_num":                 "4",
						"log_single_storage":      "400",
						"arbiter_zone_id":         acctest.CHECKSET,
						"multi_zone_combination":  lindormMultiZoneCombination,
						"arbiter_vswitch_id":      acctest.CHECKSET,
						"standby_zone_id":         acctest.CHECKSET,
						"log_spec":                "lindorm.sn1.large",
						"log_disk_category":       "cloud_efficiency",
						"core_single_storage":     "400",
						"standby_vswitch_id":      acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return lindorm.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLindormInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Only the smallest local-SSD shape is reliably in stock in
				// Hangzhou; the case covers the local_ssd_pro create/update
				// lifecycle, not a particular cluster size.
				Config: testAccConfig(map[string]interface{}{
					"disk_category":           "local_ssd_pro",
					"payment_type":            "PayAsYouGo",
					"zone_id":                 lindormHangzhouSaleZoneId,
					"vswitch_id":              "${data.alicloud_vswitches.default.ids[0]}",
					"instance_name":           "${var.name}",
					"core_spec":               "lindorm.i2.xlarge",
					"table_engine_node_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":    "local_ssd_pro",
						"payment_type":     "PayAsYouGo",
						"instance_name":    name,
						"core_spec":        "lindorm.i2.xlarge",
						"instance_storage": acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudLindormInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return lindorm.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLindormInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"vpc_id":                           "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"zone_id":                          lindormHangzhouSaleZoneId,
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
						"vswitch_id":                       acctest.CHECKSET,
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
						"vpc_id":                           acctest.CHECKSET,
						"resource_group_id":                acctest.CHECKSET,
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
	"cold_storage":                  acctest.CHECKSET,
	"deletion_proection":            acctest.CHECKSET,
	"status":                        acctest.CHECKSET,
	"time_series_engine_node_count": acctest.CHECKSET,
	"file_engine_node_count":        acctest.CHECKSET,
	"lts_node_count":                acctest.CHECKSET,
	"search_engine_node_count":      acctest.CHECKSET,
	"core_spec":                     "",
	"table_engine_node_count":       acctest.CHECKSET,
	"instance_storage":              acctest.CHECKSET,
	"zone_id":                       acctest.CHECKSET,
	"disk_category":                 "cloud_efficiency",
	"payment_type":                  "PayAsYouGo",
	"vswitch_id":                    acctest.CHECKSET,
	"instance_name":                 acctest.CHECKSET,
	//"lts_node_specification":        acctest.CHECKSET,
	//"stream_engine_specification":   acctest.CHECKSET,
	//"file_engine_specification":     acctest.CHECKSET,
	//"search_engine_specification":   acctest.CHECKSET,
	//"table_engine_specification":    acctest.CHECKSET,
	"stream_engine_node_count":    acctest.CHECKSET,
	"service_type":                acctest.CHECKSET,
	"enabled_file_engine":         acctest.CHECKSET,
	"enabled_time_serires_engine": acctest.CHECKSET,
	"enabled_table_engine":        acctest.CHECKSET,
	"enabled_search_engine":       acctest.CHECKSET,
	"enabled_lts_engine":          acctest.CHECKSET,
	"enabled_stream_engine":       acctest.CHECKSET,
	"arch_version":                acctest.CHECKSET,
}

// Lindorm is only on sale in these Hangzhou zones; the zone order returned by
// the alicloud_zones data source starts with zones where it cannot be purchased.
const lindormHangzhouSaleZoneId = "cn-hangzhou-h"

// Multi-zone instances must use one of the preset combinations published by
// Lindorm. The cn-hangzhou-ehi-aliyun combination is rejected by the sale check
// (SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED) and a dedicated VPC in Singapore
// cannot be torn down before the instance's security group is reclaimed, so the
// multi-zone case reuses the default VPC and its existing switches; each role
// switch is placed in the zone of the ap-southeast-1-abc-aliyun combination
// instead of reusing the first zones returned by the alicloud_zones data source.
const (
	lindormMultiZoneCombination   = "ap-southeast-1-abc-aliyun"
	lindormMultiZonePrimaryZoneId = "ap-southeast-1a"
	lindormMultiZoneStandbyZoneId = "ap-southeast-1c"
	lindormMultiZoneArbiterZoneId = "ap-southeast-1b"
)

func AliCloudLindormInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vswitches" "default" {
		zone_id    = "%s"
		name_regex = "^default-NODELETING-vsw-h$"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name, lindormHangzhouSaleZoneId)
}

func AliCloudLindormInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "primary" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "%s"
	}

	data "alicloud_vswitches" "standby" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "%s"
	}

	data "alicloud_vswitches" "arbiter" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name, lindormMultiZonePrimaryZoneId, lindormMultiZoneStandbyZoneId, lindormMultiZoneArbiterZoneId)
}

const lindormInstanceTestId = "ld-test-instance"

type lindormInstanceTestAPI struct {
	t             *testing.T
	engineType    int
	logNum        int
	rejectUpgrade bool

	payType             string
	rejectPayTypeChange bool
	released            bool
	goneAfterDescribe   int

	describeFailures      int
	describeFailureStatus int
	describeFailureCode   string

	mu            sync.Mutex
	upgrades      []map[string]string
	payTypeCalls  []map[string]string
	releaseCalls  []map[string]string
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
		for _, key := range []string{"UpgradeType", "LogNum", "LindormNum", "LindormSpec", "ClusterStorage", "InstanceId",
			"TsdbSpec", "TsdbNum", "SolrSpec", "SolrNum", "LtsCoreSpec", "LtsCoreNum", "StreamSpec", "StreamNum"} {
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
		gone := a.goneAfterDescribe > 0 && a.describeCalls > a.goneAfterDescribe
		payType := a.payType
		a.mu.Unlock()
		if gone {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{
				"Code":      "Instance.IsDeleted",
				"Message":   "The instance is deleted.",
				"RequestId": "test-request-id",
			}
			break
		}
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
		if payType == "" {
			payType = "POSTPAY"
		}
		response = map[string]interface{}{
			"InstanceId":       lindormInstanceTestId,
			"InstanceStatus":   "ACTIVATION",
			"InstanceAlias":    "tf-test-lindorm",
			"PayType":          payType,
			"ZoneId":           "cn-hangzhou-h",
			"ServiceType":      "lindorm",
			"LogNum":           a.logNum,
			"LogSingleStorage": 400,
			"InstanceStorage":  1000,
			"EngineType":       a.engineType,
		}
	case "ModifyInstancePayType":
		call := map[string]string{}
		for _, key := range []string{"InstanceId", "PayType", "Duration", "PricingCycle"} {
			if value := r.Form.Get(key); value != "" {
				call[key] = value
			}
		}
		a.mu.Lock()
		a.payTypeCalls = append(a.payTypeCalls, call)
		reject := a.rejectPayTypeChange
		if !reject {
			a.payType = "POSTPAY"
		}
		a.mu.Unlock()
		if reject {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{
				"Code":      "InvalidParameter",
				"Message":   "mock ModifyInstancePayType rejected",
				"RequestId": "test-request-id",
			}
		}
	case "ReleaseLindormInstance":
		a.mu.Lock()
		a.releaseCalls = append(a.releaseCalls, map[string]string{"InstanceId": r.Form.Get("InstanceId")})
		a.released = true
		a.mu.Unlock()
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

func (a *lindormInstanceTestAPI) payTypeChangeCalls() []map[string]string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]map[string]string(nil), a.payTypeCalls...)
}

func (a *lindormInstanceTestAPI) releaseInstanceCalls() []map[string]string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]map[string]string(nil), a.releaseCalls...)
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
	d, err := schema.InternalMap(lindorm.ResourceLindormInstance().Schema).Data(
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

// lindormInstanceEngineState returns the stored state of an instance that
// already has the given engine enabled, so a request-level update test can
// exercise the post-enable specification and node count branches that the ACC
// cases only hit once at enable time.
func lindormInstanceEngineState(engine string) map[string]string {
	state := map[string]string{
		"instance_storage": "1000",
		"zone_id":          "cn-hangzhou-h",
	}
	switch engine {
	case "table":
		state["enabled_table_engine"] = "true"
		state["table_engine_node_count"] = "4"
		state["table_engine_specification"] = "lindorm.g.xlarge"
		state["log_num"] = "4"
	case "search":
		state["enabled_search_engine"] = "true"
		state["search_engine_node_count"] = "2"
		state["search_engine_specification"] = "lindorm.g.xlarge"
	case "lts":
		state["enabled_lts_engine"] = "true"
		state["lts_node_count"] = "1"
		state["lts_node_specification"] = "lindorm.g.xlarge"
	case "stream":
		state["enabled_stream_engine"] = "true"
		state["stream_engine_node_count"] = "2"
		state["stream_engine_specification"] = "lindorm.g.xlarge"
	case "ts":
		state["enabled_time_serires_engine"] = "true"
		state["time_series_engine_node_count"] = "2"
		state["time_series_engine_specification"] = "lindorm.g.xlarge"
	}
	return state
}

func shortenLindormUpgradeWait(t *testing.T) {
	t.Helper()
	t.Cleanup(lindorm.ShortenLindormStateWaitForTest(time.Second))
}

func shortenLindormPayTypeWait(t *testing.T) {
	t.Helper()
	t.Cleanup(lindorm.ShortenLindormStateWaitForTest(time.Second))
}

func TestUnitLindormInstanceDeletePayType(t *testing.T) {
	shortenLindormPayTypeWait(t)

	t.Run("prepaid instance is converted to pay-as-you-go before release", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, payType: "PREPAY"}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, map[string]string{"payment_type": "Subscription", "zone_id": "cn-hangzhou-h"}, nil)

		if err := lindorm.ResourceLindormInstanceDeleteForTest(d, client); err != nil {
			t.Fatalf("deleting a prepaid instance failed: %v", err)
		}

		// Before the fix Delete only logged a warning and kept the prepaid
		// instance in the cloud, so this assertion is the counter-example.
		changes := api.payTypeChangeCalls()
		if len(changes) != 1 {
			t.Fatalf("a prepaid instance must be converted before release, got %d ModifyInstancePayType call(s): %#v", len(changes), changes)
		}
		if got := changes[0]["PayType"]; got != "POSTPAY" {
			t.Fatalf("expected PayType=POSTPAY, got %q", got)
		}
		if got := changes[0]["InstanceId"]; got != lindormInstanceTestId {
			t.Fatalf("expected InstanceId=%s, got %q", lindormInstanceTestId, got)
		}
		releases := api.releaseInstanceCalls()
		if len(releases) != 1 {
			t.Fatalf("expected exactly one ReleaseLindormInstance call after the conversion, got %d: %#v", len(releases), releases)
		}
	})

	t.Run("pay-as-you-go instance is released without a pay type change", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, payType: "POSTPAY"}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, map[string]string{"payment_type": "PayAsYouGo", "zone_id": "cn-hangzhou-h"}, nil)

		if err := lindorm.ResourceLindormInstanceDeleteForTest(d, client); err != nil {
			t.Fatalf("deleting a pay-as-you-go instance failed: %v", err)
		}
		if calls := api.payTypeChangeCalls(); len(calls) != 0 {
			t.Fatalf("a pay-as-you-go instance must not be converted, got %#v", calls)
		}
		if releases := api.releaseInstanceCalls(); len(releases) != 1 {
			t.Fatalf("expected exactly one ReleaseLindormInstance call, got %d: %#v", len(releases), releases)
		}
	})

	t.Run("a rejected pay type change fails the delete without releasing", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, payType: "PREPAY", rejectPayTypeChange: true}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, map[string]string{"payment_type": "Subscription", "zone_id": "cn-hangzhou-h"}, nil)

		err := lindorm.ResourceLindormInstanceDeleteForTest(d, client)
		if err == nil {
			t.Fatal("a rejected ModifyInstancePayType must fail the delete")
		}
		if !strings.Contains(err.Error(), "mock ModifyInstancePayType rejected") {
			t.Fatalf("the cloud error must be returned as-is, got: %v", err)
		}
		if releases := api.releaseInstanceCalls(); len(releases) != 0 {
			t.Fatalf("a failed conversion must not release the instance, got %#v", releases)
		}
	})

	t.Run("an instance released while the conversion is in progress is treated as deleted", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, payType: "PREPAY", goneAfterDescribe: 1}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, map[string]string{"payment_type": "Subscription", "zone_id": "cn-hangzhou-h"}, nil)

		if err := lindorm.ResourceLindormInstanceDeleteForTest(d, client); err != nil {
			t.Fatalf("an instance that disappeared during the conversion must not fail the delete: %v", err)
		}
		if changes := api.payTypeChangeCalls(); len(changes) != 1 {
			t.Fatalf("expected one ModifyInstancePayType call, got %#v", changes)
		}
	})
}

func TestUnitLindormInstanceTimeoutDefaults(t *testing.T) {
	r := lindorm.ResourceLindormInstance()
	cases := []struct {
		name string
		got  *time.Duration
		want time.Duration
	}{
		{"create", r.Timeouts.Create, 180 * time.Minute},
		{"update", r.Timeouts.Update, 720 * time.Minute},
		{"delete", r.Timeouts.Delete, 30 * time.Minute},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.got == nil {
				t.Fatalf("%s timeout is not set", c.name)
			}
			if *c.got != c.want {
				t.Fatalf("%s timeout: expected %s, got %s", c.name, c.want, *c.got)
			}
		})
	}
}

// TestUnitLindormInstanceUpdateEnginePaths keeps the post-enable update paths
// that the consolidated ACC cases no longer walk: after an engine is enabled,
// a specification or node count change must still reach the matching
// UpgradeLindormInstance branch for every engine, including the deprecated
// time_serires_engine_specification alias.
func TestUnitLindormInstanceUpdateEnginePaths(t *testing.T) {
	shortenLindormUpgradeWait(t)

	cases := []struct {
		name        string
		engine      string
		changes     map[string]string
		upgradeType string
		want        map[string]string
	}{
		{
			name:        "table specification change after enable is upgraded",
			engine:      "table",
			changes:     map[string]string{"table_engine_specification": "lindorm.c.8xlarge"},
			upgradeType: "upgrade-lindorm-engine",
			want:        map[string]string{"LindormSpec": "lindorm.c.8xlarge"},
		},
		{
			name:        "table node count change after enable is upgraded",
			engine:      "table",
			changes:     map[string]string{"table_engine_node_count": "8"},
			upgradeType: "upgrade-lindorm-core-num",
			want:        map[string]string{"LindormNum": "8"},
		},
		{
			name:        "search specification change after enable is upgraded",
			engine:      "search",
			changes:     map[string]string{"search_engine_specification": "lindorm.g.2xlarge"},
			upgradeType: "upgrade-search-engine",
			want:        map[string]string{"SolrSpec": "lindorm.g.2xlarge"},
		},
		{
			name:        "search node count change after enable is upgraded",
			engine:      "search",
			changes:     map[string]string{"search_engine_node_count": "3"},
			upgradeType: "upgrade-search-core-num",
			want:        map[string]string{"SolrNum": "3"},
		},
		{
			name:        "lts specification change after enable is upgraded",
			engine:      "lts",
			changes:     map[string]string{"lts_node_specification": "lindorm.g.2xlarge"},
			upgradeType: "upgrade-bds-transfer",
			want:        map[string]string{"LtsCoreSpec": "lindorm.g.2xlarge"},
		},
		{
			name:        "lts node count change after enable is upgraded",
			engine:      "lts",
			changes:     map[string]string{"lts_node_count": "2"},
			upgradeType: "upgrade-bds-core-num",
			want:        map[string]string{"LtsCoreNum": "2"},
		},
		{
			name:        "stream specification change after enable is upgraded",
			engine:      "stream",
			changes:     map[string]string{"stream_engine_specification": "lindorm.c.2xlarge"},
			upgradeType: "upgrade-stream-engine",
			want:        map[string]string{"StreamSpec": "lindorm.c.2xlarge"},
		},
		{
			name:        "stream node count change after enable is upgraded",
			engine:      "stream",
			changes:     map[string]string{"stream_engine_node_count": "3"},
			upgradeType: "upgrade-stream-core-num",
			want:        map[string]string{"StreamNum": "3"},
		},
		{
			name:        "time series node count change after enable is upgraded",
			engine:      "ts",
			changes:     map[string]string{"time_series_engine_node_count": "3"},
			upgradeType: "upgrade-tsdb-core-num",
			want:        map[string]string{"TsdbNum": "3"},
		},
		{
			name:        "deprecated time series specification alias is upgraded",
			engine:      "ts",
			changes:     map[string]string{"time_serires_engine_specification": "lindorm.g.2xlarge"},
			upgradeType: "upgrade-tsdb-engine",
			want:        map[string]string{"TsdbSpec": "lindorm.g.2xlarge"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			api := &lindormInstanceTestAPI{t: t, engineType: 4}
			client := lindormInstanceTestClient(t, api)
			d := lindormInstanceTestData(t, lindormInstanceEngineState(c.engine), c.changes)

			if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
				t.Fatalf("update failed: %v", err)
			}

			calls := api.upgradeCalls()
			if len(calls) != 1 {
				t.Fatalf("expected exactly one UpgradeLindormInstance call, got %d: %#v", len(calls), calls)
			}
			if got := calls[0]["UpgradeType"]; got != c.upgradeType {
				t.Fatalf("expected UpgradeType=%s, got %q (%#v)", c.upgradeType, got, calls[0])
			}
			for key, want := range c.want {
				if got := calls[0][key]; got != want {
					t.Fatalf("expected %s=%s, got %q (%#v)", key, want, got, calls[0])
				}
			}
		})
	}
}

func TestUnitLindormInstanceUpdateLogNum(t *testing.T) {
	shortenLindormUpgradeWait(t)
	t.Run("log_num only is upgraded independently", func(t *testing.T) {
		api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 8}
		client := lindormInstanceTestClient(t, api)
		d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"log_num": "8"})

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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

		if err := lindorm.ResourceLindormInstanceUpdateForTest(d, client); err != nil {
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
	shortenLindormUpgradeWait(t)
	api := &lindormInstanceTestAPI{t: t, engineType: 4, logNum: 4, rejectUpgrade: true}
	client := lindormInstanceTestClient(t, api)
	d := lindormInstanceTestData(t, lindormInstanceTableEngineState(), map[string]string{"log_num": "8"})

	err := lindorm.ResourceLindormInstanceUpdateForTest(d, client)
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
		service := lindorm.NewService(client)
		var tolerance lindorm.LindormPostCreateReadTolerance
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
		if lastTolerated := lindorm.LindormPostCreateReadToleranceLastForTest(&tolerance); lastTolerated == nil || !strings.Contains(lastTolerated.Error(), "OperationFailed") {
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
		service := lindorm.NewService(client)
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
		service := lindorm.NewService(client)
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

// lindormServerErrorStub mirrors the legacy SDK error shape (HttpStatus /
// ErrorCode) that the shared client can return for a failed read.
type lindormServerErrorStub struct {
	status int
	code   string
}

func (e lindormServerErrorStub) Error() string     { return e.code }
func (e lindormServerErrorStub) HttpStatus() int   { return e.status }
func (e lindormServerErrorStub) ErrorCode() string { return e.code }

func TestUnitLindormTransientPostCreateReadErrorShapes(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "legacy sdk error with 451 operation failed is transient",
			err:  lindormServerErrorStub{status: http.StatusUnavailableForLegalReasons, code: "OperationFailed"},
			want: true,
		},
		{
			name: "wrapped legacy sdk error keeps the transient classification",
			err:  errs.WrapError(lindormServerErrorStub{status: http.StatusUnavailableForLegalReasons, code: "OperationFailed"}),
			want: true,
		},
		{
			name: "legacy sdk error with another code is not transient",
			err:  lindormServerErrorStub{status: http.StatusUnavailableForLegalReasons, code: "SomeOtherCode"},
			want: false,
		},
		{
			name: "legacy sdk error with 400 is not transient",
			err:  lindormServerErrorStub{status: http.StatusBadRequest, code: "OperationFailed"},
			want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lindorm.IsLindormTransientPostCreateReadErrorForTest(c.err); got != c.want {
				t.Fatalf("expected transient=%v, got %v for %v", c.want, got, c.err)
			}
		})
	}
}

func TestUnitLindormPostCreateReadToleranceForCreateTimeout(t *testing.T) {
	cases := []struct {
		name          string
		createTimeout time.Duration
		expected      time.Duration
	}{
		{"default create timeout keeps the upper bound", 180 * time.Minute, 20 * time.Minute},
		{"timeout above the upper bound keeps the upper bound", 30 * time.Minute, 20 * time.Minute},
		{"timeout just above the upper bound keeps room for the original error", 21 * time.Minute, 15*time.Minute + 45*time.Second},
		{"timeout equal to the upper bound leaves room for the original error", 20 * time.Minute, 15 * time.Minute},
		{"short timeout is shortened to three quarters", 12 * time.Minute, 9 * time.Minute},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lindorm.LindormPostCreateReadToleranceForCreateTimeoutForTest(c.createTimeout); got != c.expected {
				t.Fatalf("create timeout %s: expected tolerance %s, got %s", c.createTimeout, c.expected, got)
			}
		})
	}
}
