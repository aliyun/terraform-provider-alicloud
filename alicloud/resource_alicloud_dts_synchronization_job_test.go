package alicloud

import (
	"fmt"
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
					"source_endpoint_instance_id":        "${data.alicloud_db_instances.db_instances_rs.instances.0.id}",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "cn-hangzhou",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestdts",
					"source_endpoint_password":           "Test12345",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${data.alicloud_db_instances.db_instances_ds.instances.0.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_region":        "cn-hangzhou",
					"destination_endpoint_database_name": "tfaccountpri_0",
					"destination_endpoint_user_name":     "tftestdts",
					"destination_endpoint_password":      "Test12345",
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
						"source_endpoint_region":             "cn-hangzhou",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestdts",
						"source_endpoint_password":           "Test12345",
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        "cn-hangzhou",
						"destination_endpoint_database_name": "tfaccountpri_0",
						"destination_endpoint_user_name":     "tftestdts",
						"destination_endpoint_password":      "Test12345",
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
					"source_endpoint_password": "Test12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "Test12345",
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
					"destination_endpoint_password": "Test12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "Test12345",
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
					"destination_endpoint_password": "Test12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "Test12345",
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
variable "creation" {
  default = "Rds"
}

data "alicloud_db_instances" "db_instances_ds" {
  name_regex = "dts_used_dest"
  status     = "Running"
}

data "alicloud_db_instances" "db_instances_rs" {
  name_regex = "dts_used_source"
  status     = "Running"
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = data.alicloud_db_instances.db_instances_ds.instances.0.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = data.alicloud_db_instances.db_instances_ds.instances.0.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = data.alicloud_db_instances.db_instances_ds.instances.0.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_db_database" "db_r" {
  count       = 2
  instance_id = data.alicloud_db_instances.db_instances_rs.instances.0.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account_r" {
  db_instance_id      = data.alicloud_db_instances.db_instances_rs.instances.0.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege_r" {
  instance_id  = data.alicloud_db_instances.db_instances_rs.instances.0.id
  account_name = alicloud_db_account.account_r.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db_r.*.name
}

resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                        = "PayAsYouGo"
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = "cn-hangzhou"
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = "cn-hangzhou"
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}

`, name)
}
