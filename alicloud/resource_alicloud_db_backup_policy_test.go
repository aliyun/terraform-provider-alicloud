package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func testAccCheckDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_backup_policy" {
			continue
		}
		request := rds.CreateDescribeBackupPolicyRequest()
		request.DBInstanceId = rs.Primary.ID
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func TestAccAlicloudDBBackupPolicy_mysql(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBBackupPolicy_mysql_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_backup_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_backup_time,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_log_backup_false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_log_backup_true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_log_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_mysql_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func TestAccAlicloudDBBackupPolicy_pgdb_high_edition(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_backup_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_backup_time,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_log_backup_false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_log_backup_true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_log_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_high_edition_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func TestAccAlicloudDBBackupPolicy_pgdb_basic_edition(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBBackupPolicy_pgdb_basic_edition_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_basic_edition_backup_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_basic_edition_backup_time,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_basic_edition_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_pgdb_basic_edition_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func TestAccAlicloudDBBackupPolicy_SQLServer(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBBackupPolicy_SQLServer_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_SQLServer_backup_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_SQLServer_backup_time,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_SQLServer_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_SQLServer_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

// Unknown current resource exists
func SkipTestAccAlicloudDBBackupPolicy_PPAS(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBBackupPolicy_PPAS_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_backup_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_backup_time,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_log_backup_false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_log_backup_true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_log_retention_period,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccDBBackupPolicy_PPAS_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

const testAccDBBackupPolicy_mysql_base = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id = "${alicloud_db_instance.default.id}"
}`

const testAccDBBackupPolicy_mysql_backup_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
}`

const testAccDBBackupPolicy_mysql_backup_time = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}`

const testAccDBBackupPolicy_mysql_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_mysql_log_backup_false = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = false
}`

const testAccDBBackupPolicy_mysql_log_backup_true = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = true
}`

const testAccDBBackupPolicy_mysql_log_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`

const testAccDBBackupPolicy_mysql_all = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Tuesday", "Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`

const testAccDBBackupPolicy_pgdb_high_edition_base = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id = "${alicloud_db_instance.default.id}"
}`

const testAccDBBackupPolicy_pgdb_high_edition_backup_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
}`

const testAccDBBackupPolicy_pgdb_high_edition_backup_time = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}`

const testAccDBBackupPolicy_pgdb_high_edition_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_pgdb_high_edition_log_backup_false = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = false
}`

const testAccDBBackupPolicy_pgdb_high_edition_log_backup_true = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = true
}`

const testAccDBBackupPolicy_pgdb_high_edition_log_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`

const testAccDBBackupPolicy_pgdb_high_edition_all = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Tuesday", "Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`

const testAccDBBackupPolicy_pgdb_basic_edition_base = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id = "${alicloud_db_instance.default.id}"
}`

const testAccDBBackupPolicy_pgdb_basic_edition_backup_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
}`

const testAccDBBackupPolicy_pgdb_basic_edition_backup_time = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}`

const testAccDBBackupPolicy_pgdb_basic_edition_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_pgdb_basic_edition_all = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PostgreSQL"
  engine_version   = "10.0"
  instance_type    = "rds.pg.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_SQLServer_base = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "SQLServer"
  engine_version   = "2008r2"
  instance_type    = "rds.mssql.s2.large"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id = "${alicloud_db_instance.default.id}"
}`

const testAccDBBackupPolicy_SQLServer_backup_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "SQLServer"
  engine_version   = "2008r2"
  instance_type    = "rds.mssql.s2.large"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
}`

const testAccDBBackupPolicy_SQLServer_backup_time = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "SQLServer"
  engine_version   = "2008r2"
  instance_type    = "rds.mssql.s2.large"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}`

const testAccDBBackupPolicy_SQLServer_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "SQLServer"
  engine_version   = "2008r2"
  instance_type    = "rds.mssql.s2.large"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_SQLServer_all = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "SQLServer"
  engine_version   = "2008r2"
  instance_type    = "rds.mssql.s2.large"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_PPAS_base = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id = "${alicloud_db_instance.default.id}"
}`

const testAccDBBackupPolicy_PPAS_backup_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
}`

const testAccDBBackupPolicy_PPAS_backup_time = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id   = "${alicloud_db_instance.default.id}"
  backup_period = ["Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}`

const testAccDBBackupPolicy_PPAS_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
}`

const testAccDBBackupPolicy_PPAS_log_backup_false = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = false
}`

const testAccDBBackupPolicy_PPAS_log_backup_true = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id      = "${alicloud_db_instance.default.id}"
  backup_period    = ["Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  retention_period = 10
  log_backup       = true
}`

const testAccDBBackupPolicy_PPAS_log_retention_period = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`

const testAccDBBackupPolicy_PPAS_all = `
variable "name" {
  default = "tf-testAccDBbackuppolicy"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
  engine           = "PPAS"
  engine_version   = "10.0"
  instance_type    = "rds.ppas.t1.small"
  instance_storage = "20"
}
resource "alicloud_db_backup_policy" "default" {
  instance_id          = "${alicloud_db_instance.default.id}"
  backup_period        = ["Tuesday", "Wednesday"]
  backup_time          = "10:00Z-11:00Z"
  retention_period     = 10
  log_backup           = true
  log_retention_period = 7
}`
