package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSJobMonitorRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_job_monitor_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSJobMonitorRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsJobMonitorRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSJobMonitorRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_id": "${alicloud_dts_migration_job.default.id}",
					"type":       "delay",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_id": CHECKSET,
						"type":       "delay",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"state": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"state": "Y",
					}),
				),
			},
			// There needs a real phone number
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"phone": "12345678987",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"phone": "12345678987",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay_rule_time": "233",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_rule_time": "233",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":            "delay",
					"state":           "N",
					"delay_rule_time": "234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":            "delay",
						"state":           "N",
						"delay_rule_time": "234",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDTSJobMonitorRuleMap0 = map[string]string{
	"dts_job_id": CHECKSET,
	"state":      CHECKSET,
}

func AlicloudDTSJobMonitorRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  count            = 2
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    =  data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = var.region
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = var.region
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = alicloud_dts_migration_instance.default.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.default.0.id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = var.region
  source_endpoint_user_name          = alicloud_rds_account.default.0.name
  source_endpoint_password           = var.password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_instance.default.1.id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = var.region
  destination_endpoint_user_name     = alicloud_rds_account.default.1.name
  destination_endpoint_password      = var.password
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  structure_initialization           = true
  data_initialization                = true
  data_synchronization               = true
  status                             = "Migrating"
  depends_on                         = [alicloud_db_account_privilege.default]
}

`, name, os.Getenv("ALICLOUD_REGION"))
}

func TestAccAlicloudDTSJobMonitorRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_job_monitor_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSJobMonitorRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsJobMonitorRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSJobMonitorRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_id": "${alicloud_dts_synchronization_job.default.id}",
					"type":       "delay",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_id": CHECKSET,
						"type":       "delay",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"state": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"state": "Y",
					}),
				),
			},
			// There needs a real phone number
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"phone": "12345678987",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"phone": "12345678987",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay_rule_time": "233",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_rule_time": "233",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":            "delay",
					"state":           "N",
					"delay_rule_time": "234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":            "delay",
						"state":           "N",
						"delay_rule_time": "234",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AlicloudDTSJobMonitorRuleBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region_id" {
  default = "%s"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "source" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}
resource "alicloud_db_instance" "dest" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}

resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                        = "PayAsYouGo"
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = var.region_id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = var.region_id
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}


resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.dest.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.dest.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_account.account.instance_id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_db_database" "db_r" {
  count       = 2
  instance_id = alicloud_db_instance.source.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account_r" {
  db_instance_id      = alicloud_db_instance.source.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege_r" {
  instance_id  = alicloud_db_account.account_r.instance_id
  account_name = alicloud_db_account.account_r.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db_r.*.name
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                     = alicloud_dts_synchronization_instance.default.id
  dts_job_name                        = "tf-testAccCase1"
  source_endpoint_instance_type       = "RDS"
  source_endpoint_instance_id         = alicloud_db_instance.source.id
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = var.region_id
  source_endpoint_database_name       = "tfaccountpri_0"
  source_endpoint_user_name           = "tftestdts"
  source_endpoint_password            = "Test12345"
  destination_endpoint_instance_type  = "RDS"
  destination_endpoint_instance_id    = alicloud_db_instance.dest.id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = var.region_id
  destination_endpoint_database_name  = "tfaccountpri_0"
  destination_endpoint_user_name      = "tftestdts"
  destination_endpoint_password       = "Test12345"
  db_list                             = "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization            = "true"
  data_initialization                 = "true"
  data_synchronization                = "true"
}

`, name, os.Getenv("ALICLOUD_REGION"))
}

func TestAccAlicloudDTSJobMonitorRule_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_job_monitor_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSJobMonitorRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsJobMonitorRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSJobMonitorRuleBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_id": "${alicloud_dts_subscription_job.default.id}",
					"type":       "delay",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_id": CHECKSET,
						"type":       "delay",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"state": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"state": "Y",
					}),
				),
			},
			// There needs a real phone number
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"phone": "12345678987",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"phone": "12345678987",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay_rule_time": "233",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_rule_time": "233",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":            "delay",
					"state":           "N",
					"delay_rule_time": "234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":            "delay",
						"state":           "N",
						"delay_rule_time": "234",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AlicloudDTSJobMonitorRuleBasicDependence2(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region_id" {
  default = "%s"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "5.6"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "5.6"
	instance_charge_type = "PostPaid"
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.instance.id
  account_name        = "tftestprivilege"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_dts_subscription_job" "default" {
    dts_job_name                        = var.name
    payment_type                        = "PayAsYouGo"
    source_endpoint_engine_name         = "MySQL"
    source_endpoint_region              = var.region_id
    source_endpoint_instance_type       = "RDS"
    source_endpoint_instance_id         = alicloud_db_instance.instance.id
    source_endpoint_database_name       = "tfaccountpri_0"
    source_endpoint_user_name           = "tftestprivilege"
    source_endpoint_password            = "Test12345"
    db_list                             =  <<EOF
        {"dtstestdata": {"name": "tfaccountpri_0", "all": true}}
    EOF
    subscription_instance_network_type  = "vpc"
    subscription_instance_vpc_id        = data.alicloud_vpcs.default.ids[0]
    subscription_instance_vswitch_id    = data.alicloud_vswitches.default.ids[0]
    status                              = "Normal"
}

`, name, os.Getenv("ALICLOUD_REGION"))
}
