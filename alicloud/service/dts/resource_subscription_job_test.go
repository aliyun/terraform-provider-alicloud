package dts_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dtsSubscriptionJobResourceId = "alicloud_dts_subscription_job.default"

func init() {
	resource.AddTestSweepers(
		"alicloud_dts_subscription_job",
		&resource.Sweeper{
			Name: "alicloud_dts_subscription_job",
			F:    testSweepDtsSubscriptionJob,
		})
}

func testSweepDtsSubscriptionJob(region string) error {
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDtsJobs"
	request := make(map[string]interface{})
	request["JobType"] = "SUBSCRIBE"
	request["Region"] = region
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
			log.Printf("[ERROR] Failed to fetch Dts SubscriptionJobs: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse Dts SubscriptionJobs: %s", err)
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
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

// lintignore: AT001
func TestAccAliCloudDTSSubscriptionJob_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", dtsTestRegion(), rand)

	// attrs mirrors the legacy cumulative resourceTestAccConfigFunc builder: each
	// step mutates this map and re-renders the full config. The placeholder step
	// below uses a FRESH map so its omitted source_endpoint_region (the legacy
	// REMOVEKEY) does not leak into the restore step.
	attrs := map[string]interface{}{
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
		"db_list":                            `{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}`,
		"subscription_instance_network_type": "vpc",
		"subscription_instance_vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
		"subscription_instance_vswitch_id":   "${data.alicloud_vswitches.default.ids.0}",
	}

	// placeholderAttrs is rendered on its own for the ExpectError step. It carries
	// every placeholder coverage field and omits the required source_endpoint_region
	// (the legacy REMOVEKEY), so the plan fails before any value reaches the API.
	placeholderAttrs := map[string]interface{}{
		"checkpoint":                       "",
		"compute_unit":                     "1",
		"database_count":                   "1",
		"db_list":                          `{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"},\"second_database\":{\"name\":\"second_database\",\"all\":true,\"state\":\"normal\"}}`,
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
		"source_endpoint_role":             "AliyunDTSDefaultRole",
		"source_endpoint_user_name":        "tftestprivilege",
		"subscription_data_type_ddl":       "true",
		"subscription_data_type_dml":       "true",
		"subscription_instance_vpc_id":     "vpc-tfacc-placeholder",
		"subscription_instance_vswitch_id": "vsw-tfacc-placeholder",
		"sync_architecture":                "oneway",
		"synchronization_direction":        "Forward",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite gates (checkoutSupportedRegions + testAccPreCheckWithRegions)
			// converged to cn-hangzhou for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// TODO: there is an api bug that the API DescribeDtsJobDetail can get resource even if it has been deleted.
		CheckDestroy: nil,
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: dtsSubscriptionJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(dtsSubscriptionJobResourceId, "dts_job_name"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "payment_type", "PayAsYouGo"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSubscriptionJobResourceId, "source_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_database_name", "N1cetest"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_user_name", "test_mysql"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_password", "N1cetest"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_ssl", "1"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "db_list", `{"test_database":{"name":"test_database","all":true,"state":"normal"}}`),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "subscription_instance_network_type", "vpc"),
					),
				},
			}
			// Turn SSL off on the source endpoint in place. The field is not ForceNew, so this
			// must go through ConfigureSubscription with the srcSSL reserve key rather than
			// recreating the job.
			attrs["source_endpoint_ssl"] = "0"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_ssl", "0"),
				),
			})
			// Turn SSL back on, covering the enable-by-update direction as well as disable.
			attrs["source_endpoint_ssl"] = "1"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_ssl", "1"),
				),
			})
			attrs["status"] = "Normal"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "status", "Normal"),
				),
			})
			attrs["subscription_instance_network_type"] = "vpc"
			attrs["subscription_instance_vpc_id"] = "${data.alicloud_vpcs.default.ids.0}"
			attrs["subscription_instance_vswitch_id"] = "${data.alicloud_vswitches.default.ids.0}"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "subscription_instance_network_type", "vpc"),
				),
			})
			attrs["status"] = "Normal"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "status", "Normal"),
				),
			})
			attrs["dts_job_name"] = "tf-testAccCase2"
			attrs["source_endpoint_instance_id"] = "${alicloud_db_instance.source.id}"
			attrs["subscription_instance_network_type"] = "vpc"
			attrs["subscription_instance_vpc_id"] = "${data.alicloud_vpcs.default.ids.0}"
			attrs["subscription_instance_vswitch_id"] = "${data.alicloud_vswitches.default.ids.0}"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "dts_job_name", "tf-testAccCase2"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "subscription_instance_network_type", "vpc"),
				),
			})
			// The expected-error step removes a required argument (source_endpoint_region) and
			// fills every other field with placeholders. It uses a FRESH attrs map so the
			// REMOVEKEY does not drop source_endpoint_region from the cumulative state used by
			// the restore step below.
			steps = append(steps, resource.TestStep{
				Config:      dtsSubscriptionJobConfig(name, placeholderAttrs),
				ExpectError: regexp.MustCompile("Missing required argument"),
			})
			// The expected-error step left an incomplete config in the working directory, and the
			// SDK reuses that directory for the post-test destroy. Re-apply the last valid
			// cumulative config so the import verification and the destroy that follows it see
			// a complete configuration.
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "dts_job_name", "tf-testAccCase2"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:      dtsSubscriptionJobResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region",
					"database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class",
					"destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone",
					"compute_unit", "error_notice", "sync_architecture",
				},
			})
			return steps
		}(),
	})
}

// lintignore: AT001
func TestAccAliCloudDTSSubscriptionJob_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"dts_job_name":                       "tf-testAccCase",
		"payment_type":                       "PayAsYouGo",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_database_name":      "tfaccountpri_0",
		"source_endpoint_user_name":          "tftestprivilege",
		"source_endpoint_password":           "Test12345",
		"db_list":                            `{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}`,
		"subscription_instance_network_type": "classic",
		"tags": map[string]string{
			"Created": "TF",
			"For":     "acceptance test",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// TODO: there is an api bug that the API DescribeDtsJobDetail can get resource even if it has been deleted.
		CheckDestroy: nil,
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: dtsSubscriptionJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "dts_job_name", "tf-testAccCase"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "payment_type", "PayAsYouGo"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSubscriptionJobResourceId, "source_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_database_name", "tfaccountpri_0"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_user_name", "tftestprivilege"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_password", "Test12345"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "db_list", `{"tfaccountpri_0":{"name":"tfaccountpri_0","all":true,"state":"normal"}}`),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "subscription_instance_network_type", "classic"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "2"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.Created", "TF"),
						resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.For", "acceptance test"),
					),
				},
			}
			attrs["tags"] = map[string]string{"Created": "TF", "For": "subscribeJob"}
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.For", "subscribeJob"),
				),
			})
			attrs["tags"] = map[string]string{"Created": "TF", "For": "Test"}
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.For", "Test"),
				),
			})
			attrs["tags"] = map[string]string{"Created": "TF-update", "For": "Test-update"}
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.Created", "TF-update"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.For", "Test-update"),
				),
			})
			// REMOVEKEY on tags: the cumulative builder drops the tags key, so the
			// rendered config has no tags block and the state ends with tags.% = 0.
			delete(attrs, "tags")
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "0"),
				),
			})
			attrs["dts_job_name"] = "tf-testAccCase1"
			steps = append(steps, resource.TestStep{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "dts_job_name", "tf-testAccCase1"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:      dtsSubscriptionJobResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region",
					"database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class",
					"destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone",
					"compute_unit", "error_notice", "sync_architecture",
				},
			})
			return steps
		}(),
	})
}

func TestAccAliCloudDTSSubscriptionJob_basic2(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
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
		"db_list":                            `{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}`,
		"subscription_instance_network_type": "classic",
		"tags": map[string]string{
			"Created": "TF",
			"For":     "acceptance test",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: dtsSubscriptionJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "dts_job_name", "tf-testAccCase"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "payment_type", "Subscription"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSubscriptionJobResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_instance_type", "RDS"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_database_name", "tfaccountpri_0"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_user_name", "tftestprivilege"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "source_endpoint_password", "Test12345"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "db_list", `{"tfaccountpri_0":{"name":"tfaccountpri_0","all":true,"state":"normal"}}`),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "subscription_instance_network_type", "classic"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(dtsSubscriptionJobResourceId, "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      dtsSubscriptionJobResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region",
					"database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class",
					"destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone",
					"compute_unit", "error_notice", "sync_architecture",
				},
			},
		},
	})
}

// dtsSubscriptionJobConfig renders the dependence plus a resource block whose
// attributes come from attrs. Callers mutate a single attrs map across steps to
// replicate the legacy cumulative configBuild closure; for REMOVEKEY entries the
// caller deletes the key from the map before re-rendering.
func dtsSubscriptionJobConfig(name string, attrs map[string]interface{}) string {
	return dtsSubscriptionJobDependence(name) + fmt.Sprintf(`
resource "alicloud_dts_subscription_job" "default" {
%s}
`, dtsRenderAttrs(attrs))
}

func dtsSubscriptionJobDependence(name string) string {
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
