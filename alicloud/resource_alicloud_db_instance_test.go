package alicloud

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_db_instance", &resource.Sweeper{
		Name: "alicloud_db_instance",
		F:    testSweepDBInstances,
	})
}

/*
"ssl_connection_string" There may be issues with circular dependencies, which cannot be tested in the case
"server_key","server_cert"  These two parameters need to be generated offline and cannot be generated in online tests
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
			// If a rds name is set by other service, it should be fetched by vswitch name and deleted.
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

func TestAccAliCloudRdsDBInstance_Mysql_8_0(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
					"db_instance_storage_type": "cloud_ssd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"db_is_ignore_case":        "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"instance_name":              name,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "cloud_ssd",
						"resource_group_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case"},
			},
			// test default port and there should not changes
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "delayed_insert_timeout",
							"value": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "22:00Z-02:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "22:00Z-02:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_minor_version": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_minor_version": "Auto",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
					"ssl_action":                 "Close",
					"instance_name":              "${var.name}",
					"monitoring_period":          "60",
					"instance_charge_type":       "Postpaid",
					"security_group_ids":         []string{},
					"auto_upgrade_minor_version": "Manual",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "delayed_insert_timeout",
							"value": "70",
						},
					},
					"encryption_key":           "${alicloud_kms_key.default.id}",
					"port":                     "3306",
					"connection_string_prefix": connectionStringPrefixSecond,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"db_instance_storage_type":   "cloud_essd",
						"ssl_action":                 "Close",
						"instance_name":              name,
						"monitoring_period":          "60",
						"zone_id":                    CHECKSET,
						"instance_charge_type":       "Postpaid",
						"connection_string":          CHECKSET,
						"port":                       "3306",
						"connection_string_prefix":   connectionStringPrefixSecond,
						"security_group_id":          CHECKSET,
						"security_group_ids.#":       "0",
						"auto_upgrade_minor_version": "Manual",
						"parameters.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_mode": SafetyMode,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_mode": SafetyMode,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "parameters", "encryption_key", "security_group_id", "storage_auto_scale", "storage_threshold", "storage_upper_bound"},
			},
		},
	})
}

func resourceDBInstanceConfigDependence(name string) string {
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

func resourceDBInstanceConfigDependenceSwitchDBInstanceHA(name string) string {
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instances" "default" { 
	name_regex = var.name
}
`, name)
}

func TestAccAliCloudRdsDBInstance_VpcId(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	rand2 := acctest.RandIntRange(1, 255)
	name := fmt.Sprintf("tftestaccdbcreatemysql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceHighAvailabilityConfigDependenceVpcId)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                     "MySQL",
					"engine_version":             "5.7",
					"instance_type":              "rds.mysql.t1.small",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"instance_charge_type":       "Postpaid",
					"instance_name":              "${var.name}",
					"db_instance_storage_type":   "local_ssd",
					"target_minor_version":       "rds_20201031",
					"zone_id":                    "${local.zone_id}",
					"zone_id_slave_a":            "${local.zone_id}",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                 "${local.vswitch_id}",
					"monitoring_period":          "60",
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "30",
					"force_restart":              false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "5.7",
						"db_instance_storage_type":   "local_ssd",
						"instance_storage":           CHECKSET,
						"instance_name":              name,
						"sql_collector_config_value": CHECKSET,
						"tde_status":                 CHECKSET,
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
					"sql_collector_config_value": "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_config_value": "180",
					}),
				),
			},
			//UpgradeDBInstanceKernelVersion
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_time":         "Immediate",
					"switch_time":          "2020-01-15T00:00:00Z",
					"target_minor_version": "rds_20201031",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "rds_20201031",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${local.vswitch_id}",
					"private_ip_address": fmt.Sprintf("${cidrhost(data.alicloud_vswitches.default.vswitches[0].cidr_block, %d)}", rand2),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "force_restart", "db_is_ignore_case", "tde_status", "sql_collector_status", "role_arn", "tde_encryption_key"},
			},
		},
	})
}

func resourceDBInstanceHighAvailabilityConfigDependenceVpcId(name string) string {
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

func TestAccAliCloudRdsDBInstance_Mysql_5_7_HighAvailabilityInstance(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_slave_zone"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceHighAvailabilityConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"db_instance_storage_type": "local_ssd",
						"category":                 "HighAvailability",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "effective_time"},
			},
		},
	})
}

func resourceDBInstanceHighAvailabilityConfigDependence(name string) string {
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

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func TestAccAliCloudRdsDBInstance_SQLServer(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap4)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceSQLServerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2012_std_ha",
					"instance_type":            "mssql.x4.medium.s2",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"category":                 "HighAvailability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2012_std_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"category":                 "HighAvailability",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2012_std_ha",
					"instance_type":            "mssql.x4.large.s2",
					"instance_storage":         "50",
					"db_instance_storage_type": "cloud_essd",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_group_ids":       []string{"${alicloud_security_group.default.0.id}"},
					"monitoring_period":        "300",
					"category":                 "HighAvailability",
					"connection_string_prefix": connectionStringPrefixSecond,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2012_std_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         "50",
						"db_instance_storage_type": "cloud_essd",
						"instance_name":            "tf-testAccDBInstanceConfig",
						"monitoring_period":        "300",
						"zone_id":                  CHECKSET,
						"instance_charge_type":     "Postpaid",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_id":        CHECKSET,
						"security_group_ids.#":     "1",
						"category":                 "HighAvailability",
						"connection_string_prefix": connectionStringPrefixSecond,
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
		},
	})
}

func resourceDBInstanceSQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  engine = "SQLServer"
  engine_version = "2012_std_ha"
  instance_charge_type = "PostPaid"
  category = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id = data.alicloud_db_zones.default.zones.0.id
  engine = "SQLServer"
  engine_version = "2012_std_ha"
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
`, name)
}

func TestAccAliCloudRdsDBInstance_PostgreSQL_12_0(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instancePostgreSQLBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "12.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"category":                 "HighAvailability",
					"target_minor_version":     "rds_postgres_1200_20231030",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "12.0",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"category":                 "HighAvailability",
						"target_minor_version":     "rds_postgres_1200_20231030",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_minor_version": "rds_postgres_1200_20240229",
					"upgrade_time":         "Immediate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "rds_postgres_1200_20240229",
						"upgrade_time":         "Immediate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_essd2",
					"instance_storage":         "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_essd2",
						"instance_storage":         "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pg_hba_conf": []interface{}{
						map[string]interface{}{
							"type":        "host",
							"user":        "all",
							"address":     "192.168.126.3",
							"database":    "all",
							"method":      "md5",
							"priority_id": "0",
							"mask":        "0",
							"option":      "ldapbasedn=CN=Users",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pg_hba_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
					"babelfish_port":           "1433",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
						"babelfish_port":           "1433",
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
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "12.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "500",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_group_ids":       []string{},
					"monitoring_period":        "60",
					"category":                 "HighAvailability",
					"db_instance_storage_type": "cloud_essd2",
					"connection_string_prefix": connectionStringPrefixSecond,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "12.0",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_name":            "tf-testAccDBInstanceConfig",
						"monitoring_period":        "60",
						"zone_id":                  CHECKSET,
						"instance_charge_type":     "Postpaid",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_id":        CHECKSET,
						"security_group_ids.#":     "0",
						"category":                 "HighAvailability",
						"db_instance_storage_type": "cloud_essd2",
						"connection_string_prefix": connectionStringPrefixSecond,
					}),
				),
			},
		},
	})
}

func resourceDBInstancePostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  	engine               = "PostgreSQL"
  	engine_version       = "12.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "12.0"
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

func TestAccAliCloudRdsDBInstance_PostgreSQL_13_0_SSL(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	manualHATime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbinstanceconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePostgreSQLSSLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "13.0",
					"instance_type":            "pg.x2.medium.2c",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_time_zone":             "America/New_York",
					"connection_string_prefix": "${var.name}",
					"port":                     "5999",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "13.0",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"private_ip_address":       CHECKSET,
						"db_time_zone":             "America/New_York",
						"deletion_protection":      "false",
						"port":                     "5999",
						"connection_string_prefix": CHECKSET,
						"instance_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_minor_version": "Manual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_minor_version": "Manual",
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
					"tcp_connection_type": "SHORT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_connection_type": "SHORT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "pg.x4.medium.2c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":                   []string{"10.168.1.12", "100.69.7.112"},
					"db_instance_ip_array_name":      "default",
					"security_ip_type":               "IPv4",
					"db_instance_ip_array_attribute": "",
					"whitelist_network_type":         "MIX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
					testAccCheck(map[string]string{
						"db_instance_ip_array_name":      "default",
						"security_ip_type":               "IPv4",
						"whitelist_network_type":         "MIX",
						"db_instance_ip_array_attribute": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3333",
					"connection_string_prefix": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config":      "Manual",
					"manual_ha_time": manualHATime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config":      "Manual",
						"manual_ha_time": manualHATime,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config":      "Auto",
						"manual_ha_time": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":      "Open",
						"ca_type":         "aliyun",
						"acl":             "prefer",
						"replication_acl": "prefer",
						"server_cert":     CHECKSET,
						"server_key":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":  "Close",
						"ca_type":     "",
						"server_cert": "",
						"server_key":  "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                      "PostgreSQL",
					"engine_version":              "13.0",
					"instance_type":               "pg.x8.medium.2c",
					"instance_storage":            "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_charge_type":        "Postpaid",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "${local.vswitch_id}",
					"security_group_ids":          []string{},
					"monitoring_period":           "60",
					"ssl_action":                  "Open",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
					"deletion_protection":         "false",
					"auto_upgrade_minor_version":  "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                      "PostgreSQL",
						"engine_version":              "13.0",
						"instance_type":               CHECKSET,
						"instance_storage":            CHECKSET,
						"instance_name":               CHECKSET,
						"monitoring_period":           "60",
						"zone_id":                     CHECKSET,
						"instance_charge_type":        "Postpaid",
						"connection_string":           CHECKSET,
						"port":                        CHECKSET,
						"security_group_id":           CHECKSET,
						"security_group_ids.#":        "0",
						"ssl_action":                  "Open",
						"ca_type":                     "aliyun",
						"client_ca_enabled":           "1",
						"client_ca_cert":              client_ca_cert2,
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": client_cert_revocation_list2,
						"acl":                         "cert",
						"replication_acl":             "cert",
						"server_cert":                 CHECKSET,
						"server_key":                  CHECKSET,
						"deletion_protection":         "false",
						"auto_upgrade_minor_version":  "Auto",
					}),
				),
			},
		},
	})
}
func resourceDBInstancePostgreSQLUpdateDBInstanceSSLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
  	engine               = "PostgreSQL"
  	engine_version       = "13.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "13.0"
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
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}
resource "alicloud_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func TestAccAliCloudRdsDBInstance_PostgreSQL_15_0_Babelfish(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixTwo := acctest.RandString(8) + "rm"
	manualHATime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePostgreSQLSSLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "15.0",
					"instance_type":            "pg.x2.medium.2c",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_time_zone":             "America/New_York",
					"connection_string_prefix": connectionStringPrefix,
					"ssl_connection_string":    "",
					"port":                     "5999",
					"babelfish_config": []interface{}{
						map[string]interface{}{
							"babelfish_enabled":    "true",
							"migration_mode":       "single-db",
							"master_username":      "test01",
							"master_user_password": "test_123456",
						},
					},
					"deletion_protection":  "true",
					"released_keep_policy": "Lastest",
					"private_ip_address":   "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "15.0",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"db_time_zone":             "America/New_York",
						"deletion_protection":      "true",
						"port":                     "5999",
						"connection_string_prefix": connectionStringPrefix,
						"babelfish_config.#":       "1",
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
					"tcp_connection_type": "SHORT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_connection_type": "SHORT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":                   []string{"10.168.1.12", "100.69.7.112"},
					"db_instance_ip_array_name":      "default",
					"security_ip_type":               "IPv4",
					"db_instance_ip_array_attribute": "",
					"whitelist_network_type":         "MIX",
					"fresh_white_list_readins":       "",
					"modify_mode":                    "Cover",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
					testAccCheck(map[string]string{
						"db_instance_ip_array_name":      "default",
						"security_ip_type":               "IPv4",
						"whitelist_network_type":         "MIX",
						"db_instance_ip_array_attribute": "",
						"fresh_white_list_readins":       "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3333",
					"connection_string_prefix": connectionStringPrefixTwo,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": connectionStringPrefixTwo,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config":      "Manual",
					"manual_ha_time": manualHATime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config":      "Manual",
						"manual_ha_time": manualHATime,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_config": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_config":      "Auto",
						"manual_ha_time": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":      "Open",
						"ca_type":         "aliyun",
						"acl":             "prefer",
						"replication_acl": "prefer",
						"server_cert":     CHECKSET,
						"server_key":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action":  "Close",
						"ca_type":     "",
						"server_cert": "",
						"server_key":  "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                      "PostgreSQL",
					"engine_version":              "15.0",
					"instance_storage":            "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_charge_type":        "Postpaid",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "${local.vswitch_id}",
					"security_group_ids":          []string{},
					"monitoring_period":           "60",
					"ssl_action":                  "Open",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
					"deletion_protection":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                      "PostgreSQL",
						"engine_version":              "15.0",
						"instance_type":               CHECKSET,
						"instance_storage":            CHECKSET,
						"instance_name":               "tf-testAccDBInstanceConfig",
						"monitoring_period":           "60",
						"zone_id":                     CHECKSET,
						"instance_charge_type":        "Postpaid",
						"connection_string":           CHECKSET,
						"port":                        CHECKSET,
						"security_group_id":           CHECKSET,
						"security_group_ids.#":        "0",
						"ssl_action":                  "Open",
						"ca_type":                     "aliyun",
						"client_ca_enabled":           "1",
						"client_ca_cert":              client_ca_cert2,
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": client_cert_revocation_list2,
						"acl":                         "cert",
						"replication_acl":             "cert",
						"server_cert":                 CHECKSET,
						"server_key":                  CHECKSET,
						"deletion_protection":         "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"released_keep_policy": "None",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"released_keep_policy": "None",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "fresh_white_list_readins", "released_keep_policy", "babelfish_config.#",
					"client_ca_enabled", "client_crl_enabled", "db_instance_ip_array_name", "encryption_key", "security_group_id", "modify_mode", "security_ip_type",
					"whitelist_network_type", "babelfish_config.2289427611.babelfish_enabled", "babelfish_config.2289427611.master_user_password", "babelfish_config.2289427611.master_username",
					"babelfish_config.2289427611.migration_mode"},
			},
		},
	})
}

const client_ca_cert = `-----BEGIN CERTIFICATE-----\nMIIC+TCCAeGgAwIBAgIJAKfv52qIKAi7MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV\nBAMMCHJvb3QtY2ExMB4XDTIxMDQyMzA3Mjk1M1oXDTMxMDQyMTA3Mjk1M1owEzER\nMA8GA1UEAwwIcm9vdC1jYTEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCyCXrZgqdge6oSji+URDXN0pMWnq4D8doP8quz09shN9TU4iqtyX+Bw+uYOoNF\ndNL4W09p8ykca3RzZghXdbHvtSZy5oCe1rup0xaATAgejDZKBi32ogLXdlA5UMyi\nc0OqIQpOZ+OmeMEVEZP7wsbDy7jS2v59d5OI4tnH2V2SDoWlI/7F9QOq36ER0UqY\nnnjJGnOsTDVeSy4ZXHMT0pXvSSLHsMMhzSJa6t3CiOuAeAW43zIS9tag0yvJI1v7\nxKSJTLs9O5V/h+oD9xofQ4kb4kOdStB2KpDteNfJWJoJYdvRMO+g1u6c2ovlc7KR\nrJPX2ZMJh14q99gPt6Dd+beVAgMBAAGjUDBOMB0GA1UdDgQWBBTDGEb5Aj6SI7hM\nC+AJa3YTNLdDrTAfBgNVHSMEGDAWgBTDGEb5Aj6SI7hMC+AJa3YTNLdDrTAMBgNV\nHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAXWXp6H4bAMZZN6b/rmuxvn4XP\n8p/7NN7BgPQSvQ24U5n8Lo2X8yXYZ4Si/NfWBitAqHceTk6rYTFhODG8CykiduHh\nowfhSjlMj9MGVw3j6I7crBuQ8clUGpy0mUNWJ9ObIdEMaVT+S1Jwk88Byf5FEBxO\nZLg+hg4NQh9qspFAtnhprU9LbcpVtQFY6uyCPs6OEOpPWF1Vtcu+ibQdIQV/e1SQ\n3NJ54R3MCfgEb9errFPv/rXscgahSMxW0sDvObAYdeIeiVeBp3wYKKFHeRNFPGT1\njzei5hlUJzGHf9DlgAH/KODvWUY5cvpuMtJY2yLyJv9xHjjyMnZZAOtHZxfR\n-----END CERTIFICATE-----`
const client_ca_cert2 = "-----BEGIN CERTIFICATE-----\nMIIC+TCCAeGgAwIBAgIJAKfv52qIKAi7MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV\nBAMMCHJvb3QtY2ExMB4XDTIxMDQyMzA3Mjk1M1oXDTMxMDQyMTA3Mjk1M1owEzER\nMA8GA1UEAwwIcm9vdC1jYTEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCyCXrZgqdge6oSji+URDXN0pMWnq4D8doP8quz09shN9TU4iqtyX+Bw+uYOoNF\ndNL4W09p8ykca3RzZghXdbHvtSZy5oCe1rup0xaATAgejDZKBi32ogLXdlA5UMyi\nc0OqIQpOZ+OmeMEVEZP7wsbDy7jS2v59d5OI4tnH2V2SDoWlI/7F9QOq36ER0UqY\nnnjJGnOsTDVeSy4ZXHMT0pXvSSLHsMMhzSJa6t3CiOuAeAW43zIS9tag0yvJI1v7\nxKSJTLs9O5V/h+oD9xofQ4kb4kOdStB2KpDteNfJWJoJYdvRMO+g1u6c2ovlc7KR\nrJPX2ZMJh14q99gPt6Dd+beVAgMBAAGjUDBOMB0GA1UdDgQWBBTDGEb5Aj6SI7hM\nC+AJa3YTNLdDrTAfBgNVHSMEGDAWgBTDGEb5Aj6SI7hMC+AJa3YTNLdDrTAMBgNV\nHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAXWXp6H4bAMZZN6b/rmuxvn4XP\n8p/7NN7BgPQSvQ24U5n8Lo2X8yXYZ4Si/NfWBitAqHceTk6rYTFhODG8CykiduHh\nowfhSjlMj9MGVw3j6I7crBuQ8clUGpy0mUNWJ9ObIdEMaVT+S1Jwk88Byf5FEBxO\nZLg+hg4NQh9qspFAtnhprU9LbcpVtQFY6uyCPs6OEOpPWF1Vtcu+ibQdIQV/e1SQ\n3NJ54R3MCfgEb9errFPv/rXscgahSMxW0sDvObAYdeIeiVeBp3wYKKFHeRNFPGT1\njzei5hlUJzGHf9DlgAH/KODvWUY5cvpuMtJY2yLyJv9xHjjyMnZZAOtHZxfR\n-----END CERTIFICATE-----"
const client_cert_revocation_list = `-----BEGIN X509 CRL-----\nMIIBpzCBkAIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290LWNhMRcN\nMjEwNDI5MDYwODMyWhcNMjEwNTI5MDYwODMyWjA4MBoCCQCG3wQwiFfYbRcNMjEw\nNDIzMTE0MTI4WjAaAgkAht8EMIhX2G8XDTIxMDQyOTA2MDc1N1qgDzANMAsGA1Ud\nFAQEAgIQATANBgkqhkiG9w0BAQsFAAOCAQEAq/M+t0zWLZzqw0T23rZsOhjd2/7+\nu1aHAW5jtjWU+lY4UxGqRsjUTJZnOiSq1w7CWhGxanyjtY/hmSeO6hGMuCmini8f\nNEq/jRvfeS7yJieFucnW4JFmz1HbqSr2S1uXRuHB1ziTRtGm3Epe0qynKm6O4L4q\nCIIqba1gye6H4BmEHaQIi4fplN7buWoeC5Ae9EdxRr3+59P4qJhHD4JGller8/QS\n3m1g75AHJO1dxvAEWy8DrrbP5SrqrsP8mmoNVIHXzCQPGEMnA1sG84365krwR+GC\noi1eBKozVqfnyLRA1C/ZY+dtt3I6zocA2Lt2+JX47VsbXApGgAPVIpKN6A==\n-----END X509 CRL-----`
const client_cert_revocation_list2 = "-----BEGIN X509 CRL-----\nMIIBpzCBkAIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290LWNhMRcN\nMjEwNDI5MDYwODMyWhcNMjEwNTI5MDYwODMyWjA4MBoCCQCG3wQwiFfYbRcNMjEw\nNDIzMTE0MTI4WjAaAgkAht8EMIhX2G8XDTIxMDQyOTA2MDc1N1qgDzANMAsGA1Ud\nFAQEAgIQATANBgkqhkiG9w0BAQsFAAOCAQEAq/M+t0zWLZzqw0T23rZsOhjd2/7+\nu1aHAW5jtjWU+lY4UxGqRsjUTJZnOiSq1w7CWhGxanyjtY/hmSeO6hGMuCmini8f\nNEq/jRvfeS7yJieFucnW4JFmz1HbqSr2S1uXRuHB1ziTRtGm3Epe0qynKm6O4L4q\nCIIqba1gye6H4BmEHaQIi4fplN7buWoeC5Ae9EdxRr3+59P4qJhHD4JGller8/QS\n3m1g75AHJO1dxvAEWy8DrrbP5SrqrsP8mmoNVIHXzCQPGEMnA1sG84365krwR+GC\noi1eBKozVqfnyLRA1C/ZY+dtt3I6zocA2Lt2+JX47VsbXApGgAPVIpKN6A==\n-----END X509 CRL-----"

func resourceDBInstancePostgreSQLSSLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
  	engine               = "PostgreSQL"
  	engine_version       = "13.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "13.0"
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
  status               = "Enabled"
}
`, name)
}
func TestAccAliCloudRdsDBInstance_MariaDB_10_3(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap4)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMariaDBDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MariaDB",
					"engine_version":           "10.3",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MariaDB",
						"engine_version":           "10.3",
						"instance_storage":         CHECKSET,
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "cloud_essd",
						"monitoring_period":        "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
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
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"engine":                   "MariaDB",
					"engine_version":           "10.3",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_group_ids":       []string{},
					"connection_string_prefix": connectionStringPrefixSecond,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MariaDB",
						"engine_version":           "10.3",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_name":            "tf-testAccDBInstanceConfig",
						"monitoring_period":        "300",
						"zone_id":                  CHECKSET,
						"instance_charge_type":     "Postpaid",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_id":        CHECKSET,
						"security_group_ids.#":     "0",
						"connection_string_prefix": connectionStringPrefixSecond,
					}),
				),
			},
		},
	})
}

func resourceDBInstanceMariaDBDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  	engine               = "MariaDB"
  	engine_version       = "10.3"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "MariaDB"
  	engine_version       = "10.3"
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
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description            = var.name
  pending_window_in_days = 7
  status                 = "Enabled"
}
`, name)
}

// Unknown current resource exists
func TestAccAliCloudRdsDBInstance_Mysql_8_0_MultiAZ(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstance_multiAZ"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysqlAZConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "rds.mysql.t1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"engine_version":           "8.0",
						"db_instance_storage_type": "local_ssd",
						"instance_storage":         CHECKSET,
						"zone_id":                  REGEXMATCH + ".*" + MULTI_IZ_SYMBOL + ".*",
						"instance_name":            "tf-testAccDBInstance_multiAZ",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudRdsDBInstance_Mysql_8_0_PrePaid(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_slave_zone"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceHighAvailabilityConfigDependence1)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "rds.mysql.s1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
					"zone_id":                  "${local.zone_id}",
					//"zone_id_slave_a":          "${local.zone_id}",
					//"zone_id_slave_b":          "${local.zone_id}",
					"vswitch_id":         "${local.vswitch_id}",
					"monitoring_period":  "60",
					"security_group_ids": "${alicloud_security_group.default.*.id}",
					"encryption_key":     "${alicloud_kms_key.default.id}",
					"security_ips":       []string{"10.168.1.12", "100.69.7.112"},
					"db_time_zone":       "America/New_York",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"engine_version":           "8.0",
						"db_instance_storage_type": "local_ssd",
						"instance_storage":         CHECKSET,
						"db_time_zone":             "America/New_York",
						"resource_group_id":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "encryption_key", "db_is_ignore_case"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":        "true",
					"auto_renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "false",
						"auto_renew_period": "0",
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
		},
	})
}

func TestAccAliCloudRdsDBInstance_Mysql_8_0_Cluster(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap6)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_Cluster"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceHighAvailabilityConfigDependence2)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.x2.medium.xc",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${local.zone_id}",
					"zone_id_slave_a":          "${local.zone_id}",
					"zone_id_slave_b":          "${local.zone_id}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_ips":             []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"engine_version":           "8.0",
						"db_instance_storage_type": "cloud_essd",
						"instance_storage":         CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "Postpaid", "encryption_key", "db_is_ignore_case"},
			},
		},
	})
}
func TestAccAliCloudRdsDBInstance_MySQL_8_0_ServerlessBasic(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceServerlessMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_MysqlServerlessBasic_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysqlServerlessBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":                  "${data.alicloud_db_zones.default.ids.1}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "cloud_essd",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type":     "Serverless",
					"category":                 "serverless_basic",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "0.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                           "MySQL",
						"engine_version":                   "8.0",
						"db_instance_storage_type":         "cloud_essd",
						"zone_id":                          CHECKSET,
						"instance_name":                    CHECKSET,
						"instance_charge_type":             CHECKSET,
						"category":                         CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "0.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "7",
							"min_capacity": "1.5",
							"auto_pause":   true,
							"switch_force": true,
						},
					},
					"instance_type": "${data.alicloud_db_instance_classes.this.instance_classes.0.instance_class}",
					"category":      "serverless_standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "7",
						"serverless_config.0.min_capacity": "1.5",
						"serverless_config.0.auto_pause":   "true",
						"serverless_config.0.switch_force": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}

func resourceDBInstanceMysqlServerlessBasicConfigDependence(name string) string {
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

data "alicloud_db_instance_classes" "this" {
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

func TestAccAliCloudRdsDBInstance_PostgreSQL_14_0_Serverless(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceServerlessMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_PostgreSQLServerlessBasic_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePostgreSQLServerlessBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":                  "${data.alicloud_db_zones.default.ids.1}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "cloud_essd",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type":     "Serverless",
					"category":                 "serverless_basic",
					"deletion_protection":      "true",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "12",
							"min_capacity": "0.5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                           CHECKSET,
						"engine_version":                   CHECKSET,
						"db_instance_storage_type":         CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_name":                    CHECKSET,
						"instance_charge_type":             CHECKSET,
						"category":                         CHECKSET,
						"deletion_protection":              "true",
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "12",
						"serverless_config.0.min_capacity": "0.5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "10",
							"min_capacity": "3.5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "10",
						"serverless_config.0.min_capacity": "3.5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
					}),
				),
			},
			// ssl_action is ignored for pg serverless
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}

func resourceDBInstancePostgreSQLServerlessBasicConfigDependence(name string) string {
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

func TestAccAliCloudRdsDBInstanceMySQL_ServerlessStandard(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceServerlessMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_MysqlServerlessStandard_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysqlServerlessStandardConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":                  "${data.alicloud_db_zones.default.ids.0}",
					"zone_id_slave_a":          "${data.alicloud_db_zones.default.ids.1}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "cloud_essd",
					"vswitch_id":               "${join(\",\", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])}",
					"instance_charge_type":     "Serverless",
					"category":                 "serverless_standard",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "0.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                           "MySQL",
						"engine_version":                   "8.0",
						"db_instance_storage_type":         "cloud_essd",
						"zone_id":                          CHECKSET,
						"instance_name":                    CHECKSET,
						"instance_charge_type":             CHECKSET,
						"category":                         CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "0.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "7",
							"min_capacity": "1.5",
							"auto_pause":   false,
							"switch_force": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "7",
						"serverless_config.0.min_capacity": "1.5",
						"serverless_config.0.auto_pause":   "false",
						"serverless_config.0.switch_force": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}

func resourceDBInstanceMysqlServerlessStandardConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "MySQL"
    engine_version = "8.0"
    instance_charge_type = "Serverless"
    category = "serverless_standard"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
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
data "alicloud_vswitches" "vswitche1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}
data "alicloud_vswitches" "vswitche2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

`, name)
}

func TestAccAliCloudRdsDBInstance_SQLServer_2019_ServerlessHA(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceServerlessMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_MssqlServerlessHA_%d", rand)
	connectionStringPrefix := acctest.RandString(8) + "rm"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMssqlServerlessHAConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2019_std_sl",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":                  "${data.alicloud_db_zones.default.ids.0}",
					"zone_id_slave_a":          "${data.alicloud_db_zones.default.ids.1}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "cloud_essd",
					"vswitch_id":               "${join(\",\", [data.alicloud_vswitches.vswitche1.ids.0, data.alicloud_vswitches.vswitche2.ids.0])}",
					"instance_charge_type":     "Serverless",
					"category":                 "serverless_ha",
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "8",
							"min_capacity": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                           CHECKSET,
						"engine_version":                   CHECKSET,
						"db_instance_storage_type":         CHECKSET,
						"zone_id":                          CHECKSET,
						"instance_name":                    CHECKSET,
						"instance_charge_type":             CHECKSET,
						"category":                         CHECKSET,
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "8",
						"serverless_config.0.min_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []interface{}{
						map[string]interface{}{
							"max_capacity": "6",
							"min_capacity": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#":              "1",
						"serverless_config.0.max_capacity": "6",
						"serverless_config.0.min_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"port": "3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": connectionStringPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": connectionStringPrefix,
					}),
				),
			},
			// ssl_action is ignored for SQLServer
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}
func resourceDBInstanceConfigDependenceDowngrade(name string) string {
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

`, name)
}
func TestAccAliCloudRdsDBInstanceMysql_Downgrade(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDowngrade)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.n4.medium.2c",
					"instance_storage":         "100",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"instance_charge_type":       CHECKSET,
						"instance_name":              name,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "cloud_ssd",
						"resource_group_id":          CHECKSET,
						"node_id":                    CHECKSET,
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
						"instance_type": "mysql.n2.medium.2c",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period"},
			},
		},
	})
}
func resourceDBInstanceConfigDependenceDBEncryptionKey(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}
`, name)
}
func TestAccAliCloudRdsDBInstanceMysql_DBEncryptionKey(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDBEncryptionKey)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.x4.medium.2c",
					"instance_storage":         "30",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "cloud_essd",
					"optimized_writes":         "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"engine_version":           "8.0",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_charge_type":     CHECKSET,
						"instance_name":            name,
						"db_instance_storage_type": "cloud_essd",
						"monitoring_period":        CHECKSET,
						"optimized_writes":         "{\"optimized_writes\":true,\"init_optimized_writes\":true}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period"},
			},
		},
	})
}
func TestAccAliCloudRdsDBInstanceMysql_DBParamGroupId(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDBParamGroupId)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.n4.medium.2c",
					"instance_storage":         "100",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_ssd",
					"db_param_group_id":        "${local.db_param_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"instance_charge_type":       CHECKSET,
						"instance_name":              name,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "cloud_ssd",
						"resource_group_id":          CHECKSET,
						"node_id":                    CHECKSET,
						"db_param_group_id":          CHECKSET,
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
						"instance_type": "mysql.n2.medium.2c",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period"},
			},
		},
	})
}
func resourceDBInstanceConfigDependenceDBParamGroupId(name string) string {
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_rds_parameter_group" "default" {
  engine = "mysql"
  engine_version = "8.0"
  param_detail{
    param_name = "back_log"
    param_value = "4000"
  }
  param_detail{
    param_name = "wait_timeout"
    param_value = "86460"
  }
  parameter_group_desc = "terrarform_test"
  parameter_group_name = "terrarform_test"
}

locals {
  db_param_group_id = alicloud_rds_parameter_group.default.id
}

`, name)
}
func resourceDBInstanceMssqlServerlessHAConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
    engine = "SQLServer"
    engine_version = "2019_std_sl"
    instance_charge_type = "Serverless"
    category = "serverless_ha"
    db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.ids.1
    engine = "SQLServer"
    engine_version = "2019_std_sl"
    category = "serverless_ha"
    db_instance_storage_type = "cloud_essd"
    instance_charge_type = "Serverless"
    commodity_code = "rds_serverless_public_cn"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "vswitche1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}
data "alicloud_vswitches" "vswitche2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

`, name)
}
func resourceDBInstanceHighAvailabilityConfigDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
   default = "%s"
}
data "alicloud_db_zones" "default"{
   engine = "MySQL"
   engine_version = "8.0"
   instance_charge_type = "PrePaid"
   category = "HighAvailability"
   db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
   engine = "MySQL"
   engine_version = "8.0"
    category = "HighAvailability"
   db_instance_storage_type = "local_ssd"
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
   security_group_name   = var.name
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

func resourceDBInstanceHighAvailabilityConfigDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
   default = "%s"
}
data "alicloud_db_zones" "default"{
   engine = "MySQL"
   engine_version = "8.0"
   instance_charge_type = "PostPaid"
   category = "cluster"
   db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
   engine = "MySQL"
   engine_version = "8.0"
    category = "cluster"
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

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
   status = "OK"
}

resource "alicloud_security_group" "default" {
   name   = var.name
   vpc_id = data.alicloud_vpcs.default.ids.0
}

`, name)
}
func resourceDBInstanceMysqlAZConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
	db_instance_storage_type = "local_ssd"
	multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
	multi_zone           = true
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.multi_zone_ids.0
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.zones.0.multi_zone_ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.zones.0.multi_zone_ids.0
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
  key_state               = "Enabled"
}

`, name)
}
func resourceDBInstanceConfigDependenceDBGeneralEssd(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_auto"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

`, name)
}
func TestAccAliCloudRdsDBInstanceMysql_general_essd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDBGeneralEssd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.x4.medium.2c",
					"instance_storage":         "100",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
					"bursting_enabled":         "true",
					"optimized_writes":         "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "MySQL",
						"engine_version":           "8.0",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_charge_type":     CHECKSET,
						"instance_name":            name,
						"db_instance_storage_type": "general_essd",
						"monitoring_period":        CHECKSET,
						"bursting_enabled":         CHECKSET,
						"optimized_writes":         "{\"optimized_writes\":true,\"init_optimized_writes\":true}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"optimized_writes": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"optimized_writes": "{\"optimized_writes\":false,\"init_optimized_writes\":true}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"optimized_writes": "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"optimized_writes": "{\"optimized_writes\":true,\"init_optimized_writes\":true}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "encryption_key", "direction", "auto_renew", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstancePostgreSQL(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap7)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigGeneralEssdPgSql)
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
					"engine":                   "PostgreSQL",
					"engine_version":           "17.0",
					"instance_type":            "pg.n4.2c.2m",
					"instance_storage":         "30",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "17.0",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_charge_type":     CHECKSET,
						"instance_name":            name,
						"db_instance_storage_type": "general_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"pg_bouncer_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pg_bouncer_enabled": "false",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "tde_encryption_key"},
			},
		},
	})
}

func resourceDBInstanceConfigGeneralEssdPgSql(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "17.0"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
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

func TestAccAliCloudRdsDBInstanceSqlService_general_essd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap8)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigGeneralEssdSqlService)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2022_web",
					"instance_type":            "mssql.x2.medium.w1",
					"instance_storage":         "50",
					"instance_charge_type":     "Postpaid",
					"category":                 "Basic",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2022_web",
						"instance_type":            CHECKSET,
						"instance_storage":         CHECKSET,
						"instance_charge_type":     CHECKSET,
						"instance_name":            name,
						"db_instance_storage_type": "general_essd",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}

func resourceDBInstanceConfigGeneralEssdSqlService(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "SQLServer"
	engine_version = "2022_web"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

`, name)
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
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

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
