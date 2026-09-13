package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDTSSynchronizationJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobBasicDependence0)
	synchronizationConfigValue := `[{\"module\":\"03\",\"name\":\"sink.batch.size.minimum\",\"value\":\"64\"},{\"module\":\"03\",\"name\":\"sink.task.number\",\"value\":\"4\"}]`

	expectedSynchronizationConfigValue := `[{"module":"03","name":"sink.batch.size.minimum","value":"64"},{"module":"03","name":"sink.task.number","value":"4"}]`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        CHECKSET,
						"db_list":                            "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}",
					}),
					// The two region attributes are only CHECKSET above because the expected value is
					// whatever region the provider is configured for; these pairs pin them to it.
					resource.TestCheckResourceAttrPair(resourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					resource.TestCheckResourceAttrPair(resourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
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
			{
				Config: testAccConfig(map[string]interface{}{
					"job_parameters": synchronizationConfigValue,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_parameters": expectedSynchronizationConfigValue,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Suspending",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Suspending",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Synchronizing",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Synchronizing",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "Lazypeople123+",
			//		"status":                   "Suspending",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": "Lazypeople123+",
			//			"status":                   "Suspending",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"source_endpoint_password": "${alicloud_rds_account.source_account.account_password}",
			//		"status":                   "Synchronizing",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"source_endpoint_password": CHECKSET,
			//			"status":                   "Synchronizing",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"destination_endpoint_password": "Lazypeople123+",
			//		"status":                        "Retrying",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"destination_endpoint_password": "Lazypeople123+",
			//			"status":                        "Retrying",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"destination_endpoint_password": "${alicloud_rds_account.target_account.account_password}",
			//		"status":                        "Synchronizing",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"destination_endpoint_password": CHECKSET,
			//			"status":                        "Synchronizing",
			//		}),
			//	),
			//},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			},
		},
	})
}

func TestAccAliCloudDTSSynchronizationJob_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
					"instance_class":                     "4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "PolarDB",
						"source_endpoint_engine_name":        "PolarDB",
						"source_endpoint_region":             CHECKSET,
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        CHECKSET,
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					resource.TestCheckResourceAttrPair(resourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			},
		},
	})
}

func TestAccAliCloudDTSSynchronizationJob_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        CHECKSET,
						"db_list":                            "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					resource.TestCheckResourceAttrPair(resourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_list": "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_list": "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			},
		},
	})
}

var AliCloudDTSSynchronizationJobMap0 = map[string]string{
	"error_phone":                      NOSET,
	"error_notice":                     NOSET,
	"delay_rule_time":                  NOSET,
	"delay_phone":                      NOSET,
	"source_endpoint_engine_name":      CHECKSET,
	"reserve":                          NOSET,
	"delay_notice":                     NOSET,
	"destination_endpoint_engine_name": CHECKSET,
	"status":                           CHECKSET,
}

func TestAccAliCloudDTSSynchronizationJob_ssl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdtssyncjobssl%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobSslDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl":      "1",
						"destination_endpoint_ssl": "1",
						// The merge happens on the request only; state must still hold the
						// reserve string exactly as it was written, with no SSL keys folded in.
						"reserve": "{\"targetTableMode\":\"2\"}",
					}),
				),
			},
			// Turn SSL off on both endpoints in place. The fields are not ForceNew, so this must
			// go through ModifyDtsJob with ModifyTypeEnum=UPDATE_RESERVED rather than recreating
			// the job.
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_ssl":      "0",
					"destination_endpoint_ssl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl":      "0",
						"destination_endpoint_ssl": "0",
					}),
				),
			},
			// Turn SSL back on, covering the enable-by-update direction as well as disable.
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_ssl":      "1",
					"destination_endpoint_ssl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl":      "1",
						"destination_endpoint_ssl": "1",
					}),
				),
			},
			// Change only one endpoint, so the Reserved payload carries a single changed key.
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_ssl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_ssl":      "0",
						"destination_endpoint_ssl": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "instance_class"},
			},
		},
	})
}

func TestAccAliCloudDTSSynchronizationJob_vpcNat(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdtssyncjobvpcnat%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobVpcNatDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// The four VPC NAT vswitch attributes only exist on ConfigureDtsJob, so this step is
			// what proves the request is accepted with them set — a rejected or misspelled
			// parameter fails the create outright. The primary and secondary vswitches are in
			// different zones, which is the arrangement the paired parameters exist for.
			{
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccCaseVpcNat",
					}),
					// Compared against the vswitch resources rather than asserted with CHECKSET,
					// so that a primary/secondary or src/dest swap in the create request mapping
					// is caught instead of passing on any non-empty value.
					resource.TestCheckResourceAttrPair(resourceId, "src_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "src_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "dest_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "dest_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
				),
			},
			// DTS does not return the vswitch ids on DescribeDtsJobDetail, so Read cannot set them
			// and an imported job has them empty. They are ignored here for that reason, not
			// because the values are expected to drift.
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password", "src_primary_vswitch_id", "src_secondary_vswitch_id", "dest_primary_vswitch_id", "dest_secondary_vswitch_id", "instance_class"},
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
  		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
  		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
  		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
  		instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
  		instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
  		instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
		storage_space = 20
		storage_type  = "ESSDPL0" 
	}

	resource "alicloud_polardb_database" "source_db" {
  		db_cluster_id = alicloud_polardb_cluster.source.id
  		db_name       = "test_database"
  		account_name  = "test_polardb"
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
  		instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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

// TestAccAliCloudDTSSynchronizationJob_instanceClassDowngrade exercises the
// TransferInstanceClass path in both directions: after the job is created with
// a high spec, it downgrades to a smaller spec and then upgrades back to the
// original. The provider must forward the transfer request whenever the
// configured class differs from the actual one and let the server decide the
// direction, instead of silently skipping downgrades.
func TestAccAliCloudDTSSynchronizationJob_instanceClassDowngrade(t *testing.T) {
	// Downgrade requires an allowlisted account: TransferInstanceClass returns
	// NoPermission for accounts without downgrade permission. Skip by default;
	// set TF_ACC_DTS_DOWNGRADE=1 to opt in once the account is allowlisted.
	if os.Getenv("TF_ACC_DTS_DOWNGRADE") == "" {
		t.Skip("Skipping instance_class downgrade test: requires allowlisted account; set TF_ACC_DTS_DOWNGRADE=1 to run")
	}
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AliCloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssyncjobdowngrade%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDTSSynchronizationJobBasicDependence0)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccDowngrade",
					}),
					resource.TestCheckResourceAttr(resourceId, "instance_class", "4xlarge"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "large",
					"db_list":        "{\\\"test_database\\\":{\\\"name\\\":\\\"test_database\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"},\\\"another_db\\\":{\\\"name\\\":\\\"another_db\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "large",
					}),
					resource.TestCheckResourceAttr(resourceId, "instance_class", "large"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "4xlarge",
					}),
					resource.TestCheckResourceAttr(resourceId, "instance_class", "4xlarge"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "4xlarge",
					"status":         "Suspending",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "4xlarge",
						"status":         "Suspending",
					}),
					resource.TestCheckResourceAttr(resourceId, "status", "Suspending"),
				),
			},
		},
	})
}
