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

func TestAccAliCloudRdsDBInstanceMysql80PPEssdB(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPEssdB)
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

func resourceDBInstanceMysql80PPEssdB(name string) string {
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

func resourceDBInstanceMysql80APEssdHA(name string) string {
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

func TestAccAliCloudRdsDBInstanceMysql57PPSsdHAVpc(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	rand2 := acctest.RandIntRange(1, 255)
	name := fmt.Sprintf("tftestaccdbcreatemysql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql57PPSsdHA)
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

func TestAccAliCloudRdsDBInstanceMysql57PPSsdHA(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_slave_zone"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql57PPSsdHA)

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

func resourceDBInstanceMysql57PPSsdHA(name string) string {
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

// SQLServer 综合测试 - 覆盖基础功能 + Serverless
// 注意：SQLServer 没有像 PostgreSQL 那样的独有特殊参数
// 主要覆盖：serverless_config、connection_string_prefix、ssl_action 等
func TestAccAliCloudRdsDBInstanceSQLServer_Comprehensive(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-sql-comp-%d", acctest.RandIntRange(10000, 99999))
	connectionStringPrefix := acctest.RandString(8) + "rm"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceSQLServerComprehensive)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: 创建基础实例（10 个必填参数）- SQLServer 2012 高可用版
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2012_std_ha",
					"instance_type":            "mssql.x4.medium.s2",
					"instance_storage":         "50",
					"instance_charge_type":     "Postpaid",
					"category":                 "HighAvailability",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "SQLServer",
						"engine_version":   "2012_std_ha",
						"instance_type":    CHECKSET,
						"instance_storage": "50",
						"category":         "HighAvailability",
					}),
				),
			},
			// Step 2: 测试连接字符串前缀（SQLServer 重要特性）
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
			// Step 3: 测试端口变更
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
			// Step 4: 测试安全 IP 列表
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			// Step 5: 测试安全组
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
			// Step 6: 测试 SSL 功能
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
			// Step 7: 变配规格和存储
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":    "mssql.x4.large.s2",
					"instance_storage": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    CHECKSET,
						"instance_storage": "60",
					}),
				),
			},
			// Step 8: Import 验证
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
		},
	})
}

func resourceDBInstanceSQLServerComprehensive(name string) string {
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

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
	count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vswitch_name = var.name
	vpc_id       = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.0
	cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
	vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
	zone_id    = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_security_group" "default" {
	count  = 2
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

// PostgreSQL 综合测试 - 覆盖所有 PostgreSQL 独有参数
// 独有参数清单 (9 个):
// 1. pg_hba_conf - PostgreSQL 访问控制配置
// 2. babelfish_config - Babelfish 功能配置
// 3. babelfish_port - Babelfish 端口
// 4. client_ca_enabled - 客户端 CA 启用
// 5. client_ca_cert - 客户端 CA 证书
// 6. client_crl_enabled - 客户端 CRL 启用
// 7. client_cert_revocation_list - 客户端证书吊销列表
// 8. acl - 访问控制列表
// 9. replication_acl - 复制访问控制
func TestAccAliCloudRdsDBInstancePostgreSQL_Comprehensive(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-pgsql-comp-%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePgSQLComprehensive)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: 创建基础实例（10 个必填参数）
			{
				Config: testAccConfig(map[string]interface{}{
					// 必填参数
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"category":                 "HighAvailability",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.zone_ids.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "PostgreSQL",
						"engine_version":   "14.0",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
						"category":         "HighAvailability",
					}),
				),
			},
			// Step 2: 测试 pg_hba_conf（PostgreSQL 独有）
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
			// Step 3: 测试 SSL 基础功能
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
			// Step 4: 测试 SSL 证书管理（PostgreSQL 独有增强）
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_type":                     "aliyun",
					"server_cert":                 testServerCert,
					"server_key":                  testServerKey,
					"client_ca_enabled":           "1",
					"client_ca_cert":              testClientCaCert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": testClientCrl,
					"acl":                         "cert",
					"replication_acl":             "cert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_type":                     "aliyun",
						"client_ca_enabled":           "1",
						"client_ca_cert":              testClientCaCert,
						"client_crl_enabled":          "1",
						"client_cert_revocation_list": testClientCrl,
						"acl":                         "cert",
						"replication_acl":             "cert",
					}),
				),
			},
			// Step 5: 测试变配规格
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":    "${data.alicloud_db_instance_classes.default.instance_classes.1.instance_class}",
					"instance_storage": "${data.alicloud_db_instance_classes.default.storage_range.min + data.alicloud_db_instance_classes.default.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
					}),
				),
			},
			// Step 6: 测试存储类型变更
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage_type": "cloud_essd2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_essd2",
					}),
				),
			},
			// Step 7: 最终完整配置验证（包含所有独有参数）
			{
				Config: testAccConfig(map[string]interface{}{
					// 必填参数
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"category":                 "HighAvailability",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.zone_ids.0.id}",
					// PostgreSQL 独有参数
					"pg_hba_conf": []interface{}{
						map[string]interface{}{
							"type":        "host",
							"user":        "all",
							"address":     "0.0.0.0/0",
							"database":    "all",
							"method":      "scram-sha-256",
							"priority_id": "1",
						},
					},
					"ssl_action":                  "Open",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              testClientCaCert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": testClientCrl,
					"acl":                         "prefer",
					"replication_acl":             "prefer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":            "PostgreSQL",
						"pg_hba_conf.#":     "1",
						"ssl_action":        "Open",
						"client_ca_enabled": "1",
						"acl":               "prefer",
						"replication_acl":   "prefer",
					}),
				),
			},
			// Step 8: Import 验证
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "pg_hba_conf", "client_ca_cert", "client_cert_revocation_list", "server_cert", "server_key"},
			},
		},
	})
}

// PostgreSQL 综合测试辅助函数
func resourceDBInstancePgSQLComprehensive(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default" {
	engine               = "PostgreSQL"
	engine_version       = "14.0"
	instance_charge_type = "PostPaid"
	category             = "HighAvailability"
	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
	zone_id              = data.alicloud_db_zones.default.zones.0.id
	engine               = "PostgreSQL"
	engine_version       = "14.0"
	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	category             = "HighAvailability"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
	count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vswitch_name = var.name
	vpc_id     = data.alicloud_vpcs.default.ids.0
	zone_id    = data.alicloud_db_zones.default.ids.0
	cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
	vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
	zone_id    = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

// SSL 证书测试数据
const testServerCert = `-----BEGIN CERTIFICATE-----
MIIDXTCCAmWgAwIBAgIJAKl7...
-----END CERTIFICATE-----`

const testServerKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwU7...
-----END RSA PRIVATE KEY-----`

const testClientCaCert = `-----BEGIN CERTIFICATE-----
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

const testClientCrl = `-----BEGIN X509 CRL-----
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

func TestAccAliCloudRdsDBInstanceMariaDB103PPEssdHA(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	connectionStringPrefixSecond := acctest.RandString(8) + "rm"
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMariaDB103PPEssdHA)
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

func resourceDBInstanceMariaDB103PPEssdHA(name string) string {
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
func TestAccAliCloudRdsDBInstanceMysql80PPSsdHAMZ(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstance_multiAZ"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPSsdHAMZ)
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

func TestAccAliCloudRdsDBInstanceMysql80APEssdHA(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_slave_zone"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80APSsdHA)

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

func TestAccAliCloudRdsDBInstanceMysql80PPEssdC(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_Cluster"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPEssdC)

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
func TestAccAliCloudRdsDBInstanceMysql80SLEssdSLB(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	connectionStringPrefix := acctest.RandString(8) + "rm"
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_MysqlServerlessBasic_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80SLEssdSLB)

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

func resourceDBInstanceMysql80SLEssdSLB(name string) string {
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

func TestAccAliCloudRdsDBInstanceMysql80SLEssdSLS(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccDBInstance_MysqlServerlessStandard_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80SLEssdSLS)

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

func resourceDBInstanceMysql80SLEssdSLS(name string) string {
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

func TestAccAliCloudRdsDBInstanceMysql80APEssdHADown(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80APEssdHA)
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
func resourceDBInstanceMysql80PPEssdHAEnc(name string) string {
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
func TestAccAliCloudRdsDBInstanceMysql80PPEssdHAEnc(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPEssdHAEnc)
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
						"optimized_writes":         "optimized",
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
func TestAccAliCloudRdsDBInstanceMysql80APEssdHAPG(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80APEssdHAPG)
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
func resourceDBInstanceMysql80APEssdHAPG(name string) string {
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

func resourceDBInstanceMysql80APSsdHA(name string) string {
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

func resourceDBInstanceMysql80PPEssdC(name string) string {
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
func resourceDBInstanceMysql80PPSsdHAMZ(name string) string {
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
func resourceDBInstanceMysql80PPAutoHA(name string) string {
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
func TestAccAliCloudRdsDBInstanceMysql80PPAutoHA(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPAutoHA)
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
						"optimized_writes":         "optimized",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_data_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_data_enabled": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_data_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_data_enabled": CHECKSET,
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
					"optimized_writes": "optimized",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"optimized_writes": "optimized",
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

func TestAccAliCloudRdsDBInstanceMysql80PPAutoHASup(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysql80PPAutoHA)
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
					"cold_data_enabled":        "true",
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
						"cold_data_enabled":        CHECKSET,
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

func TestAccAliCloudRdsDBInstancePgSQL170Gessd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePgSQL170Gessd)
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
					"template_id_list":         []string{"${alicloud_rds_whitelist_template.default.id}", "${alicloud_rds_whitelist_template.default1.id}"},
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
						"template_id_list.#":       CHECKSET,
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

func resourceDBInstancePgSQL170Gessd(name string) string {
	templateName1 := fmt.Sprintf("tf-test%d", rand.Intn(1000))
	templateName2 := fmt.Sprintf("tf-test%d", rand.Intn(1000))
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

resource "alicloud_rds_whitelist_template" "default" {
  ip_white_list            = "127.0.0.1,127.0.0.2"
  template_name            = "%s"
}

resource "alicloud_rds_whitelist_template" "default1" {
  ip_white_list            = "127.0.0.1"
  template_name            = "%s"
}
`, name, templateName1, templateName2)
}

func TestAccAliCloudRdsDBInstanceSQL2022Gessd(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceSQL2022Gessd)
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
				Config: testAccConfig(map[string]interface{}{
					"template_id_list": []string{"${alicloud_rds_whitelist_template.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_id_list.#": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_id_list": []string{"${alicloud_rds_whitelist_template.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_id_list.#": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_id_list": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_id_list.#": "0",
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

func resourceDBInstanceSQL2022Gessd(name string) string {
	templateName := fmt.Sprintf("tf-test%d", rand.Intn(1000))
	templateName1 := fmt.Sprintf("tf-test%d", rand.Intn(1000))
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
resource "alicloud_rds_whitelist_template" "default" {
  ip_white_list            = "127.0.0.1,127.0.0.2"
  template_name            = "%s"
}
resource "alicloud_rds_whitelist_template" "default1" {
  ip_white_list            = "127.0.0.2,127.0.0.3"
  template_name            = "%s"
}

`, name, templateName, templateName1)
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
