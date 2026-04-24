package alicloud

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_db_instance", &resource.Sweeper{
		Name: "alicloud_db_instance",
		F:    testSweepDBInstances,
	})
}

/*
"ssl_connection_string" There may be issues with circular dependencies, which cannot be tested in the case
"server_key","server_cert","ssl_certificate","ssl_password","tde_certificate","tde_private_key","tde_password","tde_db_name" These two parameters need to be generated offline and cannot be generated in online tests
*/

func testSweepDBInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	objects := make([]interface{}, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			continue
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		objects = append(objects, result...)
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	vpcService := VpcService{client}
	for _, v := range objects {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["DBInstanceDescription"])
		id := fmt.Sprint(item["DBInstanceId"])
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				if need, err := vpcService.needSweepVpc(fmt.Sprint(item["VpcId"]), fmt.Sprint(item["VSwitchId"])); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping RDS Instance: %s (%s)", name, id)
				continue
			}
		}

		log.Printf("[INFO] Deleting RDS Instance: %s (%s)", name, id)
		if len(item["ReadOnlyDBInstanceIds"].(map[string]interface{})["ReadOnlyDBInstanceId"].([]interface{})) > 0 {
			action := "ReleaseReadWriteSplittingConnection"
			request := map[string]interface{}{
				"RegionId":     client.RegionId,
				"DBInstanceId": id,
				"SourceIp":     client.SourceIp,
			}

			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
				if err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
						return resource.RetryableError(err)
					}
					if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidRwSplitNetType.NotFound"}) {
						return nil
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			}); err != nil {
				log.Printf("[ERROR] ReleaseReadWriteSplittingConnection error: %#v", err)
			}
		}

		action = "ModifyDBInstanceDeletionProtection"
		request = map[string]interface{}{
			"RegionId":           client.RegionId,
			"DBInstanceId":       id,
			"SourceIp":           client.SourceIp,
			"DeletionProtection": false,
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		var response map[string]interface{}
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})

		action = "DeleteDBInstance"
		request = map[string]interface{}{
			"DBInstanceId": id,
			"SourceIp":     client.SourceIp,
		}
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil && !NotFoundError(err) {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete RDS Instance (%s (%s)): %s", name, id, err)
		}
	}

	return nil
}

func initDBInstanceTest(resourceId string, baseMap map[string]string, instance *map[string]interface{}, configDep func(name string) string) (
	testAccCheck resourceAttrMapUpdate,
	rac *resourceAttrCheck,
	name string,
	testAccConfig ResourceTestAccConfigFunc,
) {
	name = fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	ra := resourceAttrInit(resourceId, baseMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac = resourceAttrCheckInit(rc, ra)
	testAccCheck = rac.resourceAttrMapUpdateSet()
	testAccConfig = resourceTestAccConfigFunc(resourceId, name, configDep)
	return
}

func resourceDBInstanceMySQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "Basic"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "Basic"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_db_instance_classes" "ha" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
enable = "On"
}

resource "alicloud_kms_key" "default" {
description = var.name
pending_window_in_days  = 7
status               = "Enabled"
}
`, name)
}

func resourceDBInstanceMySQL57HAConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "5.7"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}
data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "5.7"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}
resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}
resource "alicloud_ram_policy" "default" {
  policy_name = "${var.name}"
  policy_document = <<EOF
	{
	  "Statement": [
		{
          "Action": [
              "kms:List*",
              "kms:DescribeKey",
              "kms:TagResource",
              "kms:UntagResource"
          ],
          "Resource": [
              "acs:kms:*:*:*"
          ],
          "Effect": "Allow"
      	},
      	{
          "Action": [
              "kms:Encrypt",
              "kms:Decrypt",
              "kms:GenerateDataKey"
          ],
          "Resource": [
              "acs:kms:*:*:*"
          ],
          "Effect": "Allow",
          "Condition": {
              "StringEqualsIgnoreCase": {
                  "kms:tag/acs:rds:instance-encryption": "true"
              }
          }
      	}
	  ],
		"Version": "1"
	}
  EOF
  description = "this is a policy test"
  force = true
}
resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<EOF
	{
	  "Statement": [
		{
		  "Action": "sts:AssumeRole",
		  "Effect": "Allow",
		  "Principal": {
			"Service": [
			  "rds.aliyuncs.com"
			]
		  }
		}
	  ],
	  "Version": "1"
	}
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  policy_name = "${alicloud_ram_policy.default.policy_name}"
  role_name = "${alicloud_ram_role.default.name}"
  policy_type = "${alicloud_ram_policy.default.type}"
}
data "alicloud_ram_roles" "default" {
  name_regex = "${alicloud_ram_role_policy_attachment.default.policy_name}"
}
resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status            = "Enabled"
}
`, name)
}

func resourceDBInstanceMySQLPrepaidConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PrePaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PrePaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status               = "Enabled"
}
`, name)
}

func resourceDBInstanceMySQLServerlessConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "MySQL"
    engine_version = "8.0"
    instance_charge_type = "Serverless"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "MySQL"
    engine_version = "8.0"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_db_instance_classes" "standard" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "MySQL"
    engine_version = "8.0"
    category = "serverless_standard"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func resourceDBInstancePostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  	engine               = "PostgreSQL"
  	engine_version       = "14.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "14.0"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status            = "Enabled"
}
`, name)
}

func resourceDBInstancePostgreSQLServerlessConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "PostgreSQL"
    engine_version = "14.0"
    instance_charge_type = "Serverless"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "PostgreSQL"
    engine_version = "14.0"
    category = "serverless_basic"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func resourceDBInstanceSQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  engine = "SQLServer"
  engine_version = "2022_web"
  instance_charge_type = "PostPaid"
  category = "Basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id = data.alicloud_db_zones.default.zones.0.id
  engine = "SQLServer"
  engine_version = "2022_web"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type = "PostPaid"
  category = "Basic"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func resourceDBInstanceMariaDBConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MariaDB"
	engine_version = "10.3"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MariaDB"
	engine_version = "10.3"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status               = "Enabled"
}
`, name)
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Postpaid_Basic(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, dbInstanceMySQLBasicMap, &instance, resourceDBInstanceMySQLConfigDependence)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_essd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"db_is_ignore_case":        "false",
					"collect_stat_mode":        "Before",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"instance_name":              name,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "cloud_essd",
						"resource_group_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_ssd",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_ssd",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3306",
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3306",
						"connection_string_prefix": connectionStringPrefix,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":       []string{"10.168.1.12", "100.69.7.112"},
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":                   []string{"10.168.1.13", "100.69.7.113"},
					"db_instance_ip_array_attribute": "test_attr_update",
					"whitelist_network_type":         "MIX",
					"modify_mode":                    "Append",
					"fresh_white_list_readins":       "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_attribute": "test_attr_update",
						"whitelist_network_type":         "MIX",
						"modify_mode":                    "Append",
						"fresh_white_list_readins":       "all",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{"Created": "TF", "For": "Test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2", "tags.Created": "TF", "tags.For": "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": REMOVEKEY, "tags.Created": REMOVEKEY, "tags.For": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters":                 []map[string]interface{}{{"name": "delayed_insert_timeout", "value": "70"}},
					"maintain_time":              "22:00Z-02:00Z",
					"auto_upgrade_minor_version": "Auto",
					"instance_name":              "${var.name}" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":               "1",
						"maintain_time":              "22:00Z-02:00Z",
						"auto_upgrade_minor_version": "Auto",
						"instance_name":              name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"collect_stat_mode": "After",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collect_stat_mode": "After",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "parameters", "encryption_key", "security_group_id"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Postpaid_Advanced(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, dbInstanceMySQLBasicMap, &instance, resourceDBInstanceMySQLConfigDependence)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MySQL",
						"engine_version":   "8.0",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
						"instance_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status":         "Enabled",
					"role_arn":           "${data.alicloud_ram_roles.default.roles.0.arn}",
					"tde_encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
						"role_arn":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_auto_scale":  "Enable",
					"storage_threshold":   "40",
					"storage_upper_bound": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_auto_scale":  "Enable",
						"storage_threshold":   "40",
						"storage_upper_bound": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action":            "Open",
					"ssl_connection_string": "true",
					"ca_type":               "aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Close",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                     "MySQL",
					"engine_version":             "8.0",
					"instance_type":              "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min * 3}",
					"db_instance_storage_type":   "cloud_essd",
					"instance_name":              "${var.name}",
					"monitoring_period":          "60",
					"instance_charge_type":       "Postpaid",
					"security_group_ids":         []string{},
					"auto_upgrade_minor_version": "Manual",
					"parameters": []map[string]interface{}{
						{"name": "delayed_insert_timeout", "value": "70"},
					},
					"encryption_key":            "${alicloud_kms_key.default.id}",
					"port":                      "3306",
					"connection_string_prefix":  connectionStringPrefix,
					"db_instance_ip_array_name": "default",
					"security_ip_type":          "IPv4",
					"security_ip_mode":          SafetyMode,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"db_instance_storage_type":   "cloud_essd",
						"instance_name":              name,
						"monitoring_period":          "60",
						"instance_charge_type":       "Postpaid",
						"port":                       "3306",
						"connection_string_prefix":   connectionStringPrefix,
						"security_group_ids.#":       "0",
						"auto_upgrade_minor_version": "Manual",
						"parameters.#":               "1",
						"security_ip_mode":           SafetyMode,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "parameters", "encryption_key", "security_group_id", "storage_auto_scale", "storage_threshold", "storage_upper_bound", "tde_status", "role_arn", "tde_encryption_key", "sql_collector_status", "sql_collector_config_value"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_57_Postpaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap2, &instance, resourceDBInstanceMySQL57HAConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "5.7",
					"instance_type":            "rds.mysql.t1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
					"zone_id":                  "${local.zone_id}",
					"zone_id_slave_a":          "${local.zone_id}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
					"category":                 "HighAvailability",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"force_restart":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"db_instance_storage_type": "local_ssd",
						"category":                 "HighAvailability",
						"vpc_id":                   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "8.0",
					"effective_time": "Immediate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "8.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status":         "Enabled",
					"role_arn":           "${data.alicloud_ram_roles.default.roles.0.arn}",
					"tde_encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
						"role_arn":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_minor_version": "rds_20201031",
					"upgrade_time":         "Immediate",
					"switch_time":          "2020-01-15T00:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "rds_20201031",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_group_id": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"private_ip_address": "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_address": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_id": "Slave",
					"force":   "Yes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "effective_time", "tde_status", "role_arn", "tde_encryption_key"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Prepaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap3, &instance, resourceDBInstanceMySQLPrepaidConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.n4.medium.2c",
					"instance_storage":         "100",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"released_keep_policy":     "None",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MySQL",
						"engine_version":   "8.0",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":        "true",
					"auto_renew_period": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "Postpaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "mysql.n2.medium.2c",
					"direction":     "Down",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.n2.medium.2c",
					"instance_storage":         "100",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"instance_charge_type": "Postpaid",
						"instance_type":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "direction", "auto_renew", "auto_renew_period", "period", "released_keep_policy"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Serverless(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceServerlessMap, &instance, resourceDBInstanceMySQLServerlessConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Serverless",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "serverless_basic",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 2,
							"min_capacity": 1,
							"auto_pause":   true,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"engine_version":       "8.0",
						"instance_charge_type": "Serverless",
						"category":             "serverless_basic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.standard.instance_classes.0.instance_class}",
					"category":      "serverless_standard",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 4,
							"min_capacity": 1,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":      "serverless_standard",
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3306",
					"connection_string_prefix": acctest.RandString(8) + "rm",
					"security_ips":             []string{"10.168.1.12"},
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                 "3306",
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
					"ssl_action":          "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
						"ssl_action":          "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 4,
							"min_capacity": 1,
							"auto_pause":   false,
							"switch_force": true,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_auto_scale":  "Enable",
					"storage_threshold":   "30",
					"storage_upper_bound": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_auto_scale":  "Enable",
						"storage_threshold":   "30",
						"storage_upper_bound": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
					"optimized_writes": "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
						"optimized_writes": "optimized",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.standard.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Serverless",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "serverless_standard",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 4,
							"min_capacity": 1,
						},
					},
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "serverless_standard",
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "serverless_config", "security_group_id", "storage_auto_scale", "storage_threshold", "storage_upper_bound"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Postpaid_GeneralEssd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap9, &instance, resourceDBInstanceMySQLConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.ha.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.ha.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
					"bursting_enabled":         "true",
					"optimized_writes":         "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MySQL",
						"bursting_enabled": "true",
						"optimized_writes": "optimized",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_data_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_data_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_data_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_data_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"optimized_writes": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"optimized_writes": "none",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "bursting_enabled", "optimized_writes", "cold_data_enabled"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MySQL_80_Postpaid_Cluster(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap6, &instance, resourceDBInstanceMySQLConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.ha.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.ha.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "cluster",
					"zone_id_slave_a":          "${local.zone_id}",
					"zone_id_slave_b":          "${local.zone_id}",
					"security_ips":             []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":   "MySQL",
						"category": "cluster",
						"zone_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config":      "Manual",
					"manual_ha_time": "2025-01-15T02:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config":      "Manual",
						"manual_ha_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config": "Auto",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "zone_id_slave_a", "zone_id_slave_b", "manual_ha_time"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_PostgreSQL_Postpaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instancePostgreSQLBasicMap, &instance, resourceDBInstancePostgreSQLConfigDependence)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
					"target_minor_version":     "rds_20230430",
					"db_time_zone":             "America/New_York",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "14.0",
						"target_minor_version": "rds_20230430",
						"db_time_zone":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pg_hba_conf":        "host all all 0.0.0.0/0 md5\nhost all all ::0/0 md5\nlocal all all md5",
					"pg_bouncer_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pg_bouncer_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_ssd",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_ssd",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
					"babelfish_port":           "1334",
					"ssl_action":               "Open",
					"ca_type":                  "aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
						"ssl_action":               "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action":                  "Close",
					"client_ca_enabled":           "true",
					"client_ca_cert":              "",
					"client_crl_enabled":          "false",
					"client_cert_revocation_list": "",
					"acl":                         "host all all 0.0.0.0/0 md5",
					"replication_acl":             "host replication all 0.0.0.0/0 md5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":        "Close",
						"client_ca_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "5432",
					"connection_string_prefix": acctest.RandString(8) + "rm",
					"security_ips":             []string{"10.168.1.12"},
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
					"deletion_protection":      "true",
					"monitoring_period":        "300",
					"instance_name":            "${var.name}" + "update",
					"tcp_connection_type":      "SHORT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                 "5432",
						"deletion_protection":  "true",
						"monitoring_period":    "300",
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status":         "Enabled",
					"tde_encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
					"security_group_ids":       []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "pg_hba_conf", "tde_status", "tde_encryption_key", "client_ca_enabled", "client_ca_cert", "client_crl_enabled", "client_cert_revocation_list", "acl", "replication_acl", "security_group_id", "babelfish_port"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_PostgreSQL_Serverless(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceServerlessMap, &instance, resourceDBInstancePostgreSQLServerlessConfigDependence)
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Serverless",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "serverless_basic",
					"deletion_protection":      "false",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 2,
							"min_capacity": 1,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"instance_charge_type": "Serverless",
						"deletion_protection":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 4,
							"min_capacity": 2,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "5432",
					"connection_string_prefix": acctest.RandString(8) + "rm",
					"security_ips":             []string{"10.168.1.12"},
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "5432",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
					"ssl_action":          "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
						"ssl_action":          "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Close",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Serverless",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "serverless_basic",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": 4,
							"min_capacity": 2,
						},
					},
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "serverless_basic",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "serverless_config", "security_group_id"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_SQLServer_Postpaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap8, &instance, resourceDBInstanceSQLServerConfigDependence)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2022_web",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "Basic",
					"time_zone":                "China Standard Time",
					"collation":                "Chinese_PRC_CI_AS",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "SQLServer",
						"engine_version": "2022_web",
						"time_zone":      "China Standard Time",
						"collation":      "Chinese_PRC_CI_AS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "50",
					"instance_type":    "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3306",
					"connection_string_prefix": connectionStringPrefix,
					"security_ips":             []string{"10.168.1.12"},
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                 "3306",
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"collation": "SQL_Latin1_General_CP1_CI_AS",
					"time_zone": "UTC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collation": "SQL_Latin1_General_CP1_CI_AS",
						"time_zone": "UTC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recovery_model": "simple",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recovery_model": "simple",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action":       "Open",
					"force_encryption": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":       "Open",
						"force_encryption": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action":       "Close",
					"force_encryption": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":       "Close",
						"force_encryption": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2022_web",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "50",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "Basic",
					"connection_string_prefix": connectionStringPrefix,
					"security_group_ids":       []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "SQLServer",
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "security_group_id", "recovery_model"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_MariaDB_Postpaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	testAccCheck, rac, name, testAccConfig := initDBInstanceTest(resourceId, instanceBasicMap, &instance, resourceDBInstanceMariaDBConfigDependence)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	_ = name
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MariaDB",
					"engine_version":           "10.3",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MariaDB",
						"engine_version":   "10.3",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":            "${var.name}" + "update",
					"port":                     "3306",
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
						"port":          "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":       []string{"10.168.1.12"},
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters":    []map[string]interface{}{{"name": "delayed_insert_timeout", "value": "70"}},
					"maintain_time": "02:00Z-06:00Z",
					"tags":          map[string]string{"Env": "Test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":  "1",
						"maintain_time": "02:00Z-06:00Z",
						"tags.%":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MariaDB",
					"engine_version":           "10.3",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"category":                 "HighAvailability",
					"monitoring_period":        "60",
					"security_group_ids":       []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MariaDB",
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "parameters", "security_group_id"},
			},
		},
	})
}
func testAccCheckSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rdsService := RdsService{client}
		resp, err := rdsService.DescribeDBSecurityIps(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s security ip %#v", rs.Primary.ID, resp)

		if err != nil {
			return err
		}

		if len(resp) < 1 {
			return fmt.Errorf("DB security ip not found")
		}

		ips = rdsService.flattenDBSecurityIPs(resp)
		return nil
	}
}

func testAccCheckKeyValueInMaps(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s', expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

const client_ca_cert = `-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIJAKfv52qIKAi7MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV
BAMMCHJvb3QtY2ExMB4XDTIxMDQyMzA3Mjk1M1oXDTMxMDQyMTA3Mjk1M1owEzER
MA8GA1UEAwwIcm9vdC1jYTEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQCyCXrZgqdge6oSji+URDXN0pMWnq4D8doP8quz09shN9TU4iqtyX+Bw+uYOoNF
dNL4W09p8ykca3RzZghXdbHvtSZy5oCe1rup0xaATAgejDZKBi32ogLXdlA5UMyi
c0OqIQpOZ+OmeMEVEZP7wsbDy7jS2v59d5OI4tnH2V2SDoWlI/7F9QOq36ER0UqY
nnjJGnOsTDVeSy4ZXHMT0pXvSSLHsMMhzSJa6t3CiOuAeAW43zIS9tag0yvJI1v7
xKSJTLs9O5V/h+oD9xofQ4kb4kOdStB2KpDteNfJWJoJYdvRMO+g1u6c2ovlc7KR
rJPX2ZMJh14q99gPt6Dd+beVAgMBAAGjUDBOMB0GA1UdDgQWBBTDGEb5Aj6SI7hM
C+AJa3YTNLdDrTAfBgNVHSMEGDAWgBTDGEb5Aj6SI7hMC+AJa3YTNLdDrTAMBgNV
HRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAXWXp6H4bAMZZN6b/rmuxvn4XP
8p/7NN7BgPQSvQ24U5n8Lo2X8yXYZ4Si/NfWBitAqHceTk6rYTFhODG8CykiduHh
owfhSjlMj9MGVw3j6I7crBuQ8clUGpy0mUNWJ9ObIdEMaVT+S1Jwk88Byf5FEBxO
ZLg+hg4NQh9qspFAtnhprU9LbcpVtQFY6uyCPs6OEOpPWF1Vtcu+ibQdIQV/e1SQ
3NJ54R3MCfgEb9errFPv/rXscgahSMxW0sDvObAYdeIeiVeBp3wYKKFHeRNFPGT1
jzei5hlUJzGHf9DlgAH/KODvWUY5cvpuMtJY2yLyJv9xHjjyMnZZAOtHZxfR
-----END CERTIFICATE-----`

const client_ca_cert2 = client_ca_cert

const client_cert_revocation_list = `-----BEGIN X509 CRL-----
MIIBpzCBkAIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290LWNhMRcN
MjEwNDI5MDYwODMyWhcNMjEwNTI5MDYwODMyWjA4MBoCCQCG3wQwiFfYbRcNMjEw
NDIzMTE0MTI4WjAaAgkAht8EMIhX2G8XDTIxMDQyOTA2MDc1N1qgDzANMAsGA1Ud
FAQEAgIQATANBgkqhkiG9w0BAQsFAAOCAQEAq/M+t0zWLZzqw0T23rZsOhjd2/7+
u1aHAW5jtjWU+lY4UxGqRsjUTJZnOiSq1w7CWhGxanyjtY/hmSeO6hGMuCmini8f
NEq/jRvfeS7yJieFucnW4JFmz1HbqSr2S1uXRuHB1ziTRtGm3Epe0qynKm6O4L4q
CIIqba1gye6H4BmEHaQIi4fplN7buWoeC5Ae9EdxRr3+59P4qJhHD4JGller8/QS
3m1g75AHJO1dxvAEWy8DrrbP5SrqrsP8mmoNVIHXzCQPGEMnA1sG84365krwR+GC
oi1eBKozVqfnyLRA1C/ZY+dtt3I6zocA2Lt2+JX47VsbXApGgAPVIpKN6A==
-----END X509 CRL-----`

const client_cert_revocation_list2 = client_cert_revocation_list

var instanceBasicMap = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "5.6",
	"instance_type":        CHECKSET,
	"db_instance_type":     "Primary",
	"instance_storage":     "5",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 "3306",
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"ssl_action":           "Close",
}

var dbInstanceMySQLBasicMap = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "5.6",
	"instance_type":        CHECKSET,
	"db_instance_type":     "Primary",
	"instance_storage":     "5",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 "3306",
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"ssl_action":           "Close",
}

var instanceBasicMap2 = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "5.7",
	"instance_type":        CHECKSET,
	"instance_name":        "tf-testAccDBInstanceConfig_slave_zone",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
	"ssl_action":           "Close",
}

var instanceBasicMap3 = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "8.0",
	"instance_type":     CHECKSET,
	"instance_storage":  "5",
	"instance_name":     "tf-testAccDBInstanceConfig_slave_zone",
	"monitoring_period": "60",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
	"ssl_action":        "Close",
}

var instanceBasicMap4 = map[string]string{
	"ssl_action": "Close",
}

var instanceBasicMap5 = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "8.0",
	"instance_type":     CHECKSET,
	"instance_storage":  "5",
	"instance_name":     "tf-testAccDBInstanceSwitchDBInstanceHA",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
	"ssl_action":        "Close",
}

var instanceBasicMap6 = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "8.0",
	"instance_type":     CHECKSET,
	"instance_storage":  "20",
	"instance_name":     "tf-testAccDBInstanceConfig_Cluster",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
	"ssl_action":        "Close",
}

var instanceServerlessMap = map[string]string{
	"ssl_action": "Close",
}

var instancePostgreSQLBasicMap = map[string]string{
	"engine":               "PostgreSQL",
	"engine_version":       "12.0",
	"instance_type":        CHECKSET,
	"db_instance_type":     "Primary",
	"instance_storage":     "5",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 "5432",
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"ssl_action":           "Close",
}
var instanceBasicMap7 = map[string]string{
	"engine":               "PostgreSQL",
	"engine_version":       "17.0",
	"instance_type":        CHECKSET,
	"db_instance_type":     "Primary",
	"instance_storage":     "30",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"ssl_action":           "Close",
}
var instanceBasicMap8 = map[string]string{
	"engine":               "SQLServer",
	"engine_version":       "2022_web",
	"instance_type":        CHECKSET,
	"instance_storage":     "50",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"status":               CHECKSET,
	"create_time":          CHECKSET,
}
var instanceBasicMap9 = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "8.0",
	"instance_type":        CHECKSET,
	"instance_name":        "tf-testAccDBInstanceConfig",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"status":               CHECKSET,
	"create_time":          CHECKSET,
}
