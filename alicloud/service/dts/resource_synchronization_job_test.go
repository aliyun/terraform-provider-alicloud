package dts_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dtsSynchronizationJobResourceId = "alicloud_dts_synchronization_job.default"

// dtsSynchronizationJobConfig renders the cumulative attrs map for the DTS
// synchronization job resource, mirroring the legacy resourceTestAccConfigFunc
// builder: each call returns the dependence block plus the resource block with
// every attr currently in the map rendered via dtsRenderAttrs (defined in
// resource_instance_test.go). The caller mutates the shared attrs map between
// steps to accumulate state, matching the legacy closure semantics.
func dtsSynchronizationJobConfig(name string, attrs map[string]interface{}) string {
	return fmt.Sprintf(`
resource "alicloud_dts_synchronization_job" "default" {
%s}
`, dtsRenderAttrs(attrs))
}

func TestAccAliCloudDTSSynchronizationJob_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", dtsTestRegion(), rand)
	synchronizationConfigValue := `[{\"module\":\"03\",\"name\":\"sink.batch.size.minimum\",\"value\":\"64\"},{\"module\":\"03\",\"name\":\"sink.task.number\",\"value\":\"4\"}]`
	expectedSynchronizationConfigValue := `[{"module":"03","name":"sink.batch.size.minimum","value":"64"},{"module":"03","name":"sink.task.number","value":"4"}]`

	attrs := map[string]interface{}{
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccCase",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_db_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		"structure_initialization":           "true",
		"data_initialization":                "true",
		"data_synchronization":               "true",
		"instance_class":                     "4xlarge",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite testAccPreCheckWithRegions converged the configured
			// region to cn-hangzhou for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccCase"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
						// source_endpoint_region is interpolated from a data source; pin it via a Pair
						// check against that source rather than asserting a concrete region string.
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "source_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "destination_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "db_list", "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"),
						// status is server-driven and Computed; the legacy Map0 marked it CHECKSET.
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
						// The two region attributes are pinned to the configured region via the data source.
						resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
						resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					),
				},
			}
			attrs["dts_job_name"] = "tf-testAccCase1"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccCase1"),
					// Carried-over Map0 CHECKSET defaults remain in force.
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			attrs["job_parameters"] = synchronizationConfigValue
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "job_parameters", expectedSynchronizationConfigValue),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:            dtsSynchronizationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			})
			return steps
		}(),
	})
}

// The negative plan must run as its own test: the rejected config stays in the
// working directory and the framework reuses the last config for the post-test
// destroy, which cannot succeed on a config that is missing a required argument.
// Nothing is created here because the plan fails validation first, so the run
// ends with an empty state and no destroy.
func TestAccAliCloudDTSSynchronizationJob_createOnlyNegative(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", dtsTestRegion(), rand)

	// A fresh placeholder map is used for the ExpectError step so the cumulative
	// state is not disturbed; the legacy REMOVEKEY on dts_job_name here is a
	// no-op against an empty map, so the rendered config simply omits the
	// required dts_job_name argument, producing the expected validation error.
	// The key is therefore deliberately absent from this map.
	placeholderAttrs := map[string]interface{}{
		"checkpoint":                    "",
		"data_check_configure":          "",
		"dedicated_cluster_id":          "",
		"delay_phone":                   "",
		"destination_endpoint_owner_id": "",
		"destination_endpoint_role":     "",
		"source_endpoint_oracle_sid":    "",
		"source_endpoint_role":          "",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// Nothing is created: the plan fails before any API call, so there is
		// nothing to destroy. Keep the legacy nil-destroy semantics verbatim.
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				// Negative coverage for the create-only attributes that cannot be
				// exercised against a single-account test environment: cross-account
				// owner/role, Oracle SID, dedicated cluster, alert phones, data check
				// and the incremental checkpoint. The missing required dts_job_name
				// stops the plan before any of these values reach the API.
				Config:      AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, placeholderAttrs),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccAliCloudDTSSynchronizationJob_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccCase",
		"source_endpoint_instance_type":      "PolarDB",
		"source_endpoint_instance_id":        "${alicloud_polardb_cluster.source.id}",
		"source_endpoint_engine_name":        "PolarDB",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_polardb_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_polardb_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		"structure_initialization":           "true",
		"data_initialization":                "true",
		"data_synchronization":               "true",
		"instance_class":                     "4xlarge",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		// The job keeps advancing from prechecking to synchronizing after the
		// create wait releases, so two reads of this server-driven lifecycle value
		// can disagree even though nothing in the configuration changed.
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: []resource.TestStep{
			{
				Config: AliCloudDTSSynchronizationJobBasicDependence1(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccCase"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_instance_type", "PolarDB"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "PolarDB"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_instance_type", "RDS"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "destination_endpoint_region"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "db_list", "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
				),
			},
			{
				// The job keeps advancing from prechecking to synchronizing after the
				// create wait releases, so status differs between the apply read and the
				// import read and cannot be compared here.
				ResourceName:            dtsSynchronizationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			},
		},
	})
}

func TestAccAliCloudDTSSynchronizationJob_basic2(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccCase",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_db_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		"structure_initialization":           "true",
		"data_initialization":                "true",
		"data_synchronization":               "true",
		"instance_class":                     "4xlarge",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccCase"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "source_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_instance_type", "RDS"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "destination_endpoint_region"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "db_list", "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"),
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
						resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
						resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					),
				},
			}
			attrs["db_list"] = "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "db_list", "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:            dtsSynchronizationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			})
			return steps
		}(),
	})
}

func TestAccAliCloudDTSSynchronizationJob_ssl(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdtssyncjobssl%d", rand)

	attrs := map[string]interface{}{
		"instance_class":                     "4xlarge",
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccCaseSsl",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_db_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		// reserve is set alongside the SSL attributes so that the srcSSL/destSSL keys
		// have to be merged into an existing reserve rather than composed from scratch.
		// A merge that dropped targetTableMode would change how the job treats existing
		// destination tables, and one that produced invalid JSON would fail to configure.
		"reserve":                  "{\\\"targetTableMode\\\":\\\"2\\\"}",
		"source_endpoint_ssl":      "1",
		"destination_endpoint_ssl": "1",
		"structure_initialization": "true",
		"data_initialization":      "true",
		"data_synchronization":     "true",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// Preserve the legacy cn-beijing pin: SSL coverage runs against the
			// default VPC whose vswitches exist in cn-beijing, and the account VPC
			// quota in cn-hangzhou is frequently exhausted.
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: AliCloudDTSSynchronizationJobSslDependence(name) + dtsSynchronizationJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_ssl", "1"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_ssl", "1"),
						// The merge happens on the request only; state must still hold the
						// reserve string exactly as it was written, with no SSL keys folded in.
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "reserve", "{\"targetTableMode\":\"2\"}"),
						// Map0 CHECKSET defaults: engine names are concrete from the config, status is Computed.
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
					),
				},
			}
			// Turn SSL off on both endpoints in place. The fields are not ForceNew, so this must
			// go through ModifyDtsJob with ModifyTypeEnum=UPDATE_RESERVED rather than recreating
			// the job.
			attrs["source_endpoint_ssl"] = "0"
			attrs["destination_endpoint_ssl"] = "0"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobSslDependence(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_ssl", "0"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_ssl", "0"),
					// reserve carries over from step 1 (cumulative checkMap).
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "reserve", "{\"targetTableMode\":\"2\"}"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			// Turn SSL back on, covering the enable-by-update direction as well as disable.
			attrs["source_endpoint_ssl"] = "1"
			attrs["destination_endpoint_ssl"] = "1"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobSslDependence(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "reserve", "{\"targetTableMode\":\"2\"}"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			// Change only one endpoint, so the Reserved payload carries a single changed key.
			attrs["source_endpoint_ssl"] = "0"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobSslDependence(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_ssl", "0"),
					// destination_endpoint_ssl carries over as "1" from step 3 (cumulative checkMap).
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "reserve", "{\"targetTableMode\":\"2\"}"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:            dtsSynchronizationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			})
			return steps
		}(),
	})
}

func TestAccAliCloudDTSSynchronizationJob_vpcNat(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdtssyncjobvpcnat%d", rand)

	attrs := map[string]interface{}{
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccCaseVpcNat",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_db_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		"src_primary_vswitch_id":             "${alicloud_vswitch.primary.id}",
		"src_secondary_vswitch_id":           "${alicloud_vswitch.secondary.id}",
		"dest_primary_vswitch_id":            "${alicloud_vswitch.primary.id}",
		"dest_secondary_vswitch_id":          "${alicloud_vswitch.secondary.id}",
		"structure_initialization":           "true",
		"data_initialization":                "true",
		"data_synchronization":               "true",
		"instance_class":                     "4xlarge",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: []resource.TestStep{
			// The four VPC NAT vswitch attributes only exist on ConfigureDtsJob, so this step is
			// what proves the request is accepted with them set — a rejected or misspelled
			// parameter fails the create outright. The primary and secondary vswitches are in
			// different zones, which is the arrangement the paired parameters exist for.
			{
				Config: AliCloudDTSSynchronizationJobVpcNatDependence(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccCaseVpcNat"),
					// Compared against the vswitch resources rather than asserted with CHECKSET,
					// so that a primary/secondary or src/dest swap in the create request mapping
					// is caught instead of passing on any non-empty value.
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "src_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "src_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "dest_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(dtsSynchronizationJobResourceId, "dest_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
					// Map0 CHECKSET defaults.
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsSynchronizationJobResourceId, "status"),
				),
			},
			// DTS does not return the vswitch ids on DescribeDtsJobDetail, so Read cannot set them
			// and an imported job has them empty. They are ignored here for that reason, not
			// because the values are expected to drift. The status is server driven and can still
			// be in a precheck phase when the import refresh runs.
			{
				ResourceName:            dtsSynchronizationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "src_primary_vswitch_id", "src_secondary_vswitch_id", "dest_primary_vswitch_id", "dest_secondary_vswitch_id", "instance_class"},
			},
		},
	})
}

func AliCloudDTSSynchronizationJobVpcNatDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
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

	data "alicloud_db_instance_classes" "default" {
 		zone_id                  = data.alicloud_db_zones.default.zones.0.id
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		category                 = "HighAvailability"
 		db_instance_storage_type = "cloud_essd"
 		instance_charge_type     = "PostPaid"
	}

	resource "alicloud_vpc" "default" {
 		vpc_name   = var.name
 		cidr_block = "172.16.0.0/16"
	}

	## The RDS instances live in the primary vswitch; the secondary is created only to be handed
	## to DTS as the standby side of the VPC NAT link, so it holds no resources of its own.
	resource "alicloud_vswitch" "primary" {
 		vpc_id       = alicloud_vpc.default.id
 		cidr_block   = "172.16.0.0/24"
 		zone_id      = data.alicloud_db_zones.default.zones.0.id
 		vswitch_name = "${var.name}-primary"
	}

	resource "alicloud_vswitch" "secondary" {
 		vpc_id       = alicloud_vpc.default.id
 		cidr_block   = "172.16.1.0/24"
 		zone_id      = data.alicloud_db_zones.default.zones.1.id
 		vswitch_name = "${var.name}-secondary"
	}

	## RDS MySQL Source
	resource "alicloud_db_instance" "source" {
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
 		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
 		db_instance_storage_type = "cloud_essd"
 		vswitch_id               = alicloud_vswitch.primary.id
 		instance_name            = "${var.name}-source"
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
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
 		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
 		db_instance_storage_type = "cloud_essd"
 		vswitch_id               = alicloud_vswitch.primary.id
 		instance_name            = "${var.name}-target"
	}

	resource "alicloud_rds_account" "target_account" {
 		db_instance_id   = alicloud_db_instance.target.id
 		account_name     = "test_mysql"
 		account_password = "N1cetest"
	}

	## DTS Data Synchronization
	resource "alicloud_dts_synchronization_instance" "default" {
 		payment_type                     = "PayAsYouGo"
 		source_endpoint_engine_name      = "MySQL"
 		source_endpoint_region           = data.alicloud_regions.default.regions.0.id
 		destination_endpoint_engine_name = "MySQL"
 		destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
 		instance_class                   = "4xlarge"
 		sync_architecture                = "oneway"
	}
`, name)
}

func AliCloudDTSSynchronizationJobSslDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
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

	data "alicloud_db_instance_classes" "default" {
 		zone_id                  = data.alicloud_db_zones.default.zones.0.id
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		category                 = "HighAvailability"
 		db_instance_storage_type = "cloud_essd"
 		instance_charge_type     = "PostPaid"
	}

	// Reuse the shared default VPC instead of creating a new one: the test runs
	// in cn-beijing where the account VPC quota is frequently exhausted, and a
	// new alicloud_vpc would fail with QuotaExceeded.Vpc before the job under
	// test is reached. The default VPC already has vswitches in the db_zones zone.
	data "alicloud_vpcs" "default" {
 		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
 		vpc_id  = data.alicloud_vpcs.default.ids.0
 		zone_id = data.alicloud_db_zones.default.zones.0.id
	}

	## RDS MySQL Source, with SSL opened
	resource "alicloud_db_instance" "source" {
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
 		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
 		db_instance_storage_type = "cloud_essd"
 		vswitch_id               = data.alicloud_vswitches.default.ids.0
 		instance_name            = "rds-mysql-source-ssl"
 		ssl_action               = "Open"
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

	## RDS MySQL Target, with SSL opened
	resource "alicloud_db_instance" "target" {
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
 		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
 		db_instance_storage_type = "cloud_essd"
 		vswitch_id               = data.alicloud_vswitches.default.ids.0
 		instance_name            = "rds-mysql-target-ssl"
 		ssl_action               = "Open"
	}

	resource "alicloud_rds_account" "target_account" {
 		db_instance_id   = alicloud_db_instance.target.id
 		account_name     = "test_mysql"
 		account_password = "N1cetest"
	}

	## DTS Data Synchronization
	resource "alicloud_dts_synchronization_instance" "default" {
 		payment_type                     = "PayAsYouGo"
 		source_endpoint_engine_name      = "MySQL"
 		source_endpoint_region           = data.alicloud_regions.default.regions.0.id
 		destination_endpoint_engine_name = "MySQL"
 		destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
 		instance_class                   = "4xlarge"
 		sync_architecture                = "oneway"
	}
`, name)
}

func AliCloudDTSSynchronizationJobBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
	}

	// The region has to come from the provider rather than from ALICLOUD_REGION: that variable is not
	// always set, and an empty region_id makes CreateDtsInstance fail with MissingDestinationRegion
	// before the job under test is ever reached.
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
 		name_regex =  "^default-NODELETING$"
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

	## DTS Data Synchronization
	resource "alicloud_dts_synchronization_instance" "default" {
 		payment_type                     = "PayAsYouGo"
 		source_endpoint_engine_name      = "MySQL"
 		source_endpoint_region           = data.alicloud_regions.default.regions.0.id
 		destination_endpoint_engine_name = "MySQL"
 		destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
 		instance_class                   = "4xlarge"
 		sync_architecture                = "oneway"
	}
`, name)
}

func AliCloudDTSSynchronizationJobBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
	}

	// Both endpoints below are created in whatever region the provider is configured for, so the DTS
	// instance has to be told that same region instead of a hardcoded one.
	data "alicloud_regions" "default" {
 		current = true
	}

	data "alicloud_polardb_zones" "default" {}

	// PolarDB sales zones in the configured region do not always overlap with the
	// first zone returned by alicloud_db_zones; pick a zone that actually sells
	// PolarDB so CreateDBCluster does not fail with InvalidZoneID.NotFound. RDS
	// MySQL HA is available in the same zone, so the db_instance_classes lookup
	// and the RDS target keep working.
	locals {
 		zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
	}

	data "alicloud_vpcs" "default" {
 		name_regex =  "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
 		vpc_id  = data.alicloud_vpcs.default.ids.0
 		zone_id = local.zone_id
	}

	data "alicloud_db_instance_classes" "default" {
 		zone_id                  = local.zone_id
 		engine                   = "MySQL"
 		engine_version           = "8.0"
 		category                 = "HighAvailability"
 		db_instance_storage_type = "cloud_essd"
 		instance_charge_type     = "PostPaid"
	}

	data "alicloud_polardb_node_classes" "default" {
 		db_type    = "MySQL"
 		db_version = "8.0"
 		pay_type   = "PostPaid"
 		zone_id    = local.zone_id
	}

	## PolarDB PolarDB Source
	resource "alicloud_polardb_cluster" "source" {
 		db_type       = "MySQL"
 		db_version    = "8.0"
 		pay_type      = "PostPaid"
 		//db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
 		db_node_class = "polar.mysql.x4.medium.c"
 		vswitch_id    = data.alicloud_vswitches.default.ids.0
 		description   = "polardb_cluster_description"
		encrypt_new_tables = "OFF"
		// PolarDB Binlog is off by default, and DTS fails precheck for a PolarDB
		// source without it ("loose_polar_log_bin" must be ON).
		loose_polar_log_bin = "ON"
	}

	resource "alicloud_polardb_database" "source_db" {
 		db_cluster_id = alicloud_polardb_cluster.source.id
 		db_name       = "test_database"
 		account_name  = "test_polardb"
 		depends_on = [alicloud_polardb_account.source_account]
	}

	resource "alicloud_polardb_account" "source_account" {
 		db_cluster_id    = alicloud_polardb_cluster.source.id
 		account_name     = "test_polardb"
 		account_password = "N1cetest"
	}

	resource "alicloud_polardb_account_privilege" "source_privilege" {
 		db_cluster_id     = alicloud_polardb_cluster.source.id
 		account_name      = alicloud_polardb_account.source_account.account_name
 		account_privilege = "ReadWrite"
 		db_names          = alicloud_polardb_database.source_db.*.db_name
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

	## DTS Data Synchronization
	resource "alicloud_dts_synchronization_instance" "default" {
 		payment_type                     = "PayAsYouGo"
 		source_endpoint_engine_name      = "PolarDB"
 		source_endpoint_region           = data.alicloud_regions.default.regions.0.id
 		destination_endpoint_engine_name = "MySQL"
 		destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
 		instance_class                   = "4xlarge"
 		sync_architecture                = "oneway"
	}
`, name)
}

// TestAccAliCloudDTSSynchronizationJob_instanceClassDowngrade is an opt-in probe
// of a class change to a smaller spec. TransferInstanceClass documents only
// UPGRADE, and any downgrade path is limited to allowlisted synchronization
// instances of class large or above with a minimum target of medium, so the case
// is gated on TF_ACC_DTS_DOWNGRADE and skipped by default. It records the
// reserved class-change path; it is not evidence that a downgrade is available.
func TestAccAliCloudDTSSynchronizationJob_instanceClassDowngrade(t *testing.T) {
	// Downgrade requires an allowlisted account: TransferInstanceClass returns
	// NoPermission for accounts without downgrade permission. Skip by default;
	// set TF_ACC_DTS_DOWNGRADE=1 to opt in once the account is allowlisted.
	if os.Getenv("TF_ACC_DTS_DOWNGRADE") == "" {
		t.Skip("Skipping instance_class downgrade test: requires allowlisted account; set TF_ACC_DTS_DOWNGRADE=1 to run")
	}
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssyncjobdowngrade%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
		"dts_job_name":                       "tf-testAccDowngrade",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_region":             "${data.alicloud_regions.default.regions.0.id}",
		"source_endpoint_database_name":      "test_database",
		"source_endpoint_user_name":          "${alicloud_db_account_privilege.source_privilege.account_name}",
		"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_region":        "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_database_name": "test_database",
		"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
		"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
		"db_list":                            "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
		"structure_initialization":           "true",
		"data_initialization":                "true",
		"data_synchronization":               "true",
		"instance_class":                     "4xlarge",
		"synchronization_direction":          "Forward",
		"source_endpoint_port":               "3306",
		"source_endpoint_role":               "Source",
		"source_endpoint_vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
		"destination_endpoint_port":          "3306",
		"destination_endpoint_role":          "Destination",
		"dts_bis_label":                      "tf-testacc",
		"data_check_configure":               "{}",
		"checkpoint":                         "",
		"status":                             "Synchronizing",
		"dedicated_cluster_id":               "",
		"delay_notice":                       "false",
		"delay_phone":                        "",
		"delay_rule_time":                    "",
		"destination_endpoint_ip":            "",
		"destination_endpoint_oracle_sid":    "",
		"destination_endpoint_owner_id":      "",
		"error_notice":                       "false",
		"error_phone":                        "",
		"source_endpoint_ip":                 "",
		"source_endpoint_oracle_sid":         "",
		"source_endpoint_owner_id":           "",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsSynchronizationJobDestroy(dtsSynchronizationJobResourceId),
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "dts_job_name", "tf-testAccDowngrade"),
						// Map0 marked status CHECKSET, but the config explicitly sets status="Synchronizing",
						// so the cumulative checkMap carries the concrete value rather than CHECKSET.
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "status", "Synchronizing"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "instance_class", "4xlarge"),
						// Map0 CHECKSET defaults for engine names (concrete from config).
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					),
				},
			}
			attrs["instance_class"] = "large"
			attrs["db_list"] = "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"},\\\"another_db\\\":{\\\"name\\\":\\\"another_db\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "instance_class", "large"),
					// status carries over as "Synchronizing" from step 1 (cumulative checkMap).
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "status", "Synchronizing"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
				),
			})
			attrs["instance_class"] = "4xlarge"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "instance_class", "4xlarge"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "status", "Synchronizing"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
				),
			})
			attrs["instance_class"] = "4xlarge"
			attrs["status"] = "Suspending"
			steps = append(steps, resource.TestStep{
				Config: AliCloudDTSSynchronizationJobBasicDependence0(name) + dtsSynchronizationJobConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "instance_class", "4xlarge"),
					// status overridden to "Suspending" in this step's changeMap.
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "status", "Suspending"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsSynchronizationJobResourceId, "destination_endpoint_engine_name", "MySQL"),
				),
			})
			return steps
		}(),
	})
}

func checkDtsSynchronizationJobDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_synchronization_job" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeDtsSynchronizationJob(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("the resource %s %s was not destroyed", resourceId, rs.Primary.ID)
		}
		return nil
	}
}
