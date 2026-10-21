package alicloud

import (
	"fmt"
	"log"
	"regexp"
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
		"alicloud_dts_subscription_job",
		&resource.Sweeper{
			Name: "alicloud_dts_subscription_job",
			F:    testSweepDtsSubscriptionJob,
		})
}

func testSweepDtsSubscriptionJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDtsJobs"
	request := make(map[string]interface{})
	request["JobType"] = "SUBSCRIBE"
	request["Region"] = region
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
			log.Printf("[ERROR] Failed to fetch Dts SubscriptionJobs: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dts_subscription_jobs", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse Dts SubscriptionJobs: %s", WrapErrorf(err, FailedGetAttributeMsg, action, "$.DtsJobList", response))
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DtsJobName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dts SubscriptionJob: %s", item["DtsJobName"].(string))
				continue
			}

			action := "DeleteDtsJob"
			request := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			request["DtsInstanceId"] = item["DtsInstanceID"]
			request["RegionId"] = client.RegionId
			_, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dts SubscriptionJob (%s): %s", item["DtsJobName"].(string), err)
			}
			log.Printf("[INFO] Delete Dts SubscriptionJob success: %s ", item["DtsJobName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

// lintignore: AT001
func TestAccAliCloudDTSSubscriptionJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSubscriptionJobBasicDependence0)
	// The expected-error step removes a required argument through REMOVEKEY and
	// fills every other field with placeholders. It needs its own builder:
	// feeding that map into the cumulative builder would drop
	// source_endpoint_region from all later steps and leak the placeholders
	// into the restore step and the post-test destroy.
	testAccConfigPlaceholder := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSubscriptionJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		// TODO: there is an api bug that the API DescribeDtsJobDetail can get resource even if it has been deleted.
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "${var.name}",
					"payment_type":                       "PayAsYouGo",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"source_endpoint_database_name":      "${alicloud_rds_account.source_account.account_password}",
					"source_endpoint_user_name":          "${alicloud_rds_account.source_account.account_name}",
					"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
					"source_endpoint_ssl":                "1",
					"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
					"subscription_instance_vswitch_id":   "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       CHECKSET,
						"payment_type":                       "PayAsYouGo",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "N1cetest",
						"source_endpoint_user_name":          "test_mysql",
						"source_endpoint_password":           "N1cetest",
						"source_endpoint_ssl":                "1",
						"db_list":                            "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			// Turn SSL off on the source endpoint in place. The field is not ForceNew, so this
			// must go through ConfigureSubscription with the srcSSL reserve key rather than
			// recreating the job.
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_ssl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl": "0",
					}),
				),
			},
			// Turn SSL back on, covering the enable-by-update direction as well as disable.
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_ssl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl": "1",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"dts_job_name": "tf-testAccCase1",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"dts_job_name": "tf-testAccCase1",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "Lazypeople123+",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": "Lazypeople123+",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "Test12345",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": "Test12345",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Normal",
					}),
				),
			},
			// TODO: There is an api bug that does not support to update db_list even if its status is Normal
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"db_list": "{\\\"dtstestdata\\\": {   \\\"name\\\": \\\"tfaccountpri_0\\\",   \\\"all\\\": true }}",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"db_list": "{\"dtstestdata\": {   \"name\": \"tfaccountpri_0\",   \"all\": true }}",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Abnormal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Abnormal",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
					"subscription_instance_vswitch_id":   "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase2",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
					"subscription_instance_vswitch_id":   "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase2",
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			{
				// Cross-account owner/role, Oracle SID, the incremental
				// checkpoint and the phone/Reserve placeholders cannot be
				// provisioned in a single-account test environment, and the
				// remaining endpoint parameters are only exercised as
				// placeholder coverage here. Removing the required
				// source_endpoint_region stops the plan before any of these
				// values reach the API.
				Config: testAccConfigPlaceholder(map[string]interface{}{
					"checkpoint":                       "",
					"compute_unit":                     "1",
					"database_count":                   "1",
					"db_list":                          "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"},\\\"second_database\\\":{\\\"name\\\":\\\"second_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"delay_notice":                     "true",
					"delay_phone":                      "",
					"delay_rule_time":                  "10",
					"destination_endpoint_engine_name": "MySQL",
					"destination_region":               "cn-hangzhou",
					"dts_instance_id":                  "dts-tfacc-placeholder",
					"error_notice":                     "true",
					"error_phone":                      "",
					"instance_class":                   "large",
					"payment_type":                     "PayAsYouGo",
					"reserve":                          "",
					"source_endpoint_database_name":    "tfacc_db",
					"source_endpoint_engine_name":      "Oracle",
					"source_endpoint_instance_id":      "rm-tfacc-placeholder",
					"source_endpoint_instance_type":    "ECS",
					"source_endpoint_ip":               "172.16.0.1",
					"source_endpoint_oracle_sid":       "orcl",
					"source_endpoint_owner_id":         "",
					"source_endpoint_port":             "3306",
					"source_endpoint_region":           REMOVEKEY,
					"source_endpoint_role":             "AliyunDTSDefaultRole",
					"source_endpoint_user_name":        "tftestprivilege",
					"subscription_data_type_ddl":       "true",
					"subscription_data_type_dml":       "true",
					"subscription_instance_vpc_id":     "vpc-tfacc-placeholder",
					"subscription_instance_vswitch_id": "vsw-tfacc-placeholder",
					"sync_architecture":                "oneway",
					"synchronization_direction":        "Forward",
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// The expected-error step left an incomplete config in the working
			// directory, and the SDK reuses that directory for the post-test
			// destroy. Re-apply the last valid config so the import verification
			// and the destroy that follows it see a complete configuration.
			{
				Config: testAccConfig(nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccCase2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudDTSSubscriptionJob_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSubscriptionJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		// TODO: there is an api bug that the API DescribeDtsJobDetail can get resource even if it has been deleted.
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase",
					"payment_type":                       "PayAsYouGo",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestprivilege",
					"source_endpoint_password":           "Test12345",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "classic",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"payment_type":                       "PayAsYouGo",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestprivilege",
						"source_endpoint_password":           "Test12345",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "classic",
						"tags.%":                             "2",
						"tags.Created":                       "TF",
						"tags.For":                           "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "subscribeJob",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "subscribeJob",
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
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name": "tf-testAccCase1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccCase1",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "Lazypeople123+",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": "Lazypeople123+",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "Test12345",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": "Test12345",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Normal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Normal",
			//		}),
			//	),
			//},
			// TODO: There is an api bug that does not support to update db_list even if its status is Normal
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"db_list": "{\\\"dtstestdata\\\": {   \\\"name\\\": \\\"tfaccountpri_0\\\",   \\\"all\\\": true }}",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"db_list": "{\"dtstestdata\": {   \"name\": \"tfaccountpri_0\",   \"all\": true }}",
			//		}),
			//	),
			//},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

func TestAccAliCloudDTSSubscriptionJob_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSubscriptionJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase",
					"payment_type":                       "Subscription",
					"payment_duration_unit":              "Month",
					"payment_duration":                   "1",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestprivilege",
					"source_endpoint_password":           "Test12345",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "classic",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"payment_type":                       "Subscription",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestprivilege",
						"source_endpoint_password":           "Test12345",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "classic",
						"tags.%":                             "2",
						"tags.Created":                       "TF",
						"tags.For":                           "acceptance test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

var AliCloudDTSSubscriptionJobMap0 = map[string]string{
	"auto_pay":                         NOSET,
	"delay_rule_time":                  NOSET,
	"compute_unit":                     NOSET,
	"delay_phone":                      NOSET,
	"subscription_data_type_dml":       CHECKSET,
	"error_notice":                     NOSET,
	"sync_architecture":                NOSET,
	"quantity":                         NOSET,
	"error_phone":                      NOSET,
	"period":                           NOSET,
	"destination_region":               NOSET,
	"delay_notice":                     NOSET,
	"reserve":                          NOSET,
	"synchronization_direction":        NOSET,
	"auto_start":                       NOSET,
	"database_count":                   NOSET,
	"instance_class":                   NOSET,
	"subscription_data_type_ddl":       CHECKSET,
	"destination_endpoint_engine_name": NOSET,
	"used_time":                        NOSET,
	"source_endpoint_ssl":              NOSET,
	"status":                           CHECKSET,
}

func AliCloudDTSSubscriptionJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "tf-testaccdts%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

## RDS MySQL Source
resource "alicloud_db_instance" "source" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-source"
  ssl_action       = "Open"
}

resource "alicloud_db_database" "source_db" {
  instance_id = alicloud_db_instance.source.id
  name        = "test_database"
}

resource "alicloud_rds_account" "source_account" {
  db_instance_id   = alicloud_db_instance.source.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_db_account_privilege" "source_privilege" {
  instance_id  = alicloud_db_instance.source.id
  account_name = alicloud_rds_account.source_account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.source_db.*.name
}

## RDS MySQL Target
resource "alicloud_db_instance" "target" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-target"
}

resource "alicloud_rds_account" "target_account" {
  db_instance_id   = alicloud_db_instance.target.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}
`, name)
}
