package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSSynchronizationJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSSynchronizationJobBasicDependence0)
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
					"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
					"dts_job_name":                       "tf-testAccCase",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${var.region_id}",
					"source_endpoint_database_name":      "test_database",
					"source_endpoint_user_name":          "${alicloud_rds_account.source_account.account_name}",
					"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_region":        "${var.region_id}",
					"destination_endpoint_database_name": "test_database",
					"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
					"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             os.Getenv("ALICLOUD_REGION"),
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        os.Getenv("ALICLOUD_REGION"),
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "Lazypeople123+",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "Lazypeople123+",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "N1cetest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "N1cetest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "Lazypeople123+",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "Lazypeople123+",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "N1cetest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "N1cetest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Suspending",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Suspending",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "Lazypeople123+",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "Lazypeople123+",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "N1cetest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "N1cetest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Synchronizing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Synchronizing",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

func TestAccAlicloudDTSSynchronizationJob_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSSynchronizationJobBasicDependence1)
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
					"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
					"dts_job_name":                       "tf-testAccCase",
					"source_endpoint_instance_type":      "PolarDB",
					"source_endpoint_instance_id":        "${alicloud_polardb_cluster.source.id}",
					"source_endpoint_engine_name":        "PolarDB",
					"source_endpoint_region":             "${var.region_id}",
					"source_endpoint_database_name":      "test_database",
					"source_endpoint_user_name":          "${alicloud_polardb_account.source_account.account_name}",
					"source_endpoint_password":           "${alicloud_polardb_account.source_account.account_password}",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_region":        "${var.region_id}",
					"destination_endpoint_database_name": "test_database",
					"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
					"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "PolarDB",
						"source_endpoint_engine_name":        "PolarDB",
						"source_endpoint_region":             os.Getenv("ALICLOUD_REGION"),
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        os.Getenv("ALICLOUD_REGION"),
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

func TestAccAlicloudDTSSynchronizationJob_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSSynchronizationJobBasicDependence0)
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
					"dts_instance_id":                    "${alicloud_dts_synchronization_instance.default.id}",
					"dts_job_name":                       "tf-testAccCase",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.source.id}",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${var.region_id}",
					"source_endpoint_database_name":      "test_database",
					"source_endpoint_user_name":          "${alicloud_rds_account.source_account.account_name}",
					"source_endpoint_password":           "${alicloud_rds_account.source_account.account_password}",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${alicloud_db_instance.target.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_region":        "${var.region_id}",
					"destination_endpoint_database_name": "test_database",
					"destination_endpoint_user_name":     "${alicloud_rds_account.target_account.account_name}",
					"destination_endpoint_password":      "${alicloud_rds_account.target_account.account_password}",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             os.Getenv("ALICLOUD_REGION"),
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        os.Getenv("ALICLOUD_REGION"),
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
					}),
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
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

var AlicloudDTSSynchronizationJobMap0 = map[string]string{
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

func AlicloudDTSSynchronizationJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region_id" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
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
  db_instance_storage_type = "local_ssd"
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
  source_endpoint_region           = var.region_id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = var.region_id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

`, name, os.Getenv("ALICLOUD_REGION"))
}

func AlicloudDTSSynchronizationJobBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region_id" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
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
  db_instance_storage_type = "local_ssd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "5.6"
  pay_type   = "PostPaid"
  zone_id    = data.alicloud_db_zones.default.zones.0.id
}

## PolarDB PolarDB Source
resource "alicloud_polardb_cluster" "source" {
  db_type       = "MySQL"
  db_version    = "5.6"
  pay_type      = "PostPaid"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  vswitch_id    = data.alicloud_vswitches.default.ids.0
  description   = "polardb_cluster_description"
}

resource "alicloud_polardb_database" "source_db" {
  db_cluster_id = alicloud_polardb_cluster.source.id
  db_name       = "test_database"
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
  source_endpoint_region           = var.region_id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = var.region_id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

`, name, os.Getenv("ALICLOUD_REGION"))
}
