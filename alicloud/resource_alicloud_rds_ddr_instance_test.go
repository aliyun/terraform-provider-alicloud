package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

/*
	Because the generation time of the disaster recovery set is uncertain, the query backup set may not have a value, so the 'Test' of the disaster recovery new instance automation cannot be run regularly online. The 'Test' of the disaster recovery new instance automation has been simulated offline.
*/

func init() {
	resource.AddTestSweepers("alicloud_rds_ddr_instance", &resource.Sweeper{
		Name: "alicloud_rds_ddr_instance",
		F:    testSweepDBDdrInstances,
	})
}

func testSweepDBDdrInstances(region string) error {
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
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

func SkipTestAccAlicloudRdsDdrInstanceMysql(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_rds_ddr_instance.default"
	ra := resourceAttrInit(resourceId, ddrinstanceBasicMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDdrDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdrDBInstanceConfigDependence)
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
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"payment_type":             "PayAsYouGo",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "local_ssd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"restore_type":             "BackupSet",
					"backup_set_id":            "${data.alicloud_rds_cross_region_backups.backups.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "local_ssd",
						"resource_group_id":          CHECKSET,
						"restore_type":               CHECKSET,
						"backup_set_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
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
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckDdrSecurityIpExists("alicloud_rds_ddr_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"db_instance_storage_type":   "local_ssd",
					"instance_name":              "tf-testAccDdrDBInstanceConfig",
					"monitoring_period":          "60",
					"payment_type":               "PayAsYouGo",
					"restore_type":               "BackupSet",
					"backup_set_id":              "${data.alicloud_rds_cross_region_backups.backups.ids.0}",
					"security_group_ids":         []string{},
					"auto_upgrade_minor_version": "Manual",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "delayed_insert_timeout",
							"value": "70",
						},
					},
					"encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"db_instance_storage_type":   "local_ssd",
						"instance_name":              "tf-testAccDdrDBInstanceConfig",
						"monitoring_period":          "60",
						"zone_id":                    CHECKSET,
						"payment_type":               "PayAsYouGo",
						"connection_string":          CHECKSET,
						"port":                       CHECKSET,
						"security_group_ids.#":       "0",
						"auto_upgrade_minor_version": "Manual",
						"parameters.#":               "1",
						"restore_type":               CHECKSET,
						"backup_set_id":              CHECKSET,
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
		},
	})
}

func SkipTestAccAlicloudRdsDdrInstanceMysqlTime(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_rds_ddr_instance.default"
	ra := resourceAttrInit(resourceId, ddrinstanceBasicMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDdrDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdrDBInstanceConfigDependence)
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
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"payment_type":             "PayAsYouGo",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "local_ssd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"restore_type":             "BackupTime",
					"restore_time":             "${data.alicloud_rds_cross_region_backups.backups.backups.0.recovery_end_time}",
					"source_region":            "${data.alicloud_rds_cross_region_backups.backups.backups.0.restore_regions.0}",
					"source_db_instance_name":  "${data.alicloud_rds_cross_region_backups.backups.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "local_ssd",
						"resource_group_id":          CHECKSET,
						"restore_type":               "BackupTime",
						"restore_time":               CHECKSET,
						"source_region":              CHECKSET,
						"source_db_instance_name":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
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
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckDdrSecurityIpExists("alicloud_rds_ddr_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"db_instance_storage_type":   "local_ssd",
					"instance_name":              "tf-testAccDdrDBInstanceConfig",
					"monitoring_period":          "60",
					"payment_type":               "PayAsYouGo",
					"restore_type":               "BackupTime",
					"restore_time":               "${data.alicloud_rds_cross_region_backups.backups.backups.0.recovery_end_time}",
					"source_region":              "${data.alicloud_rds_cross_region_backups.backups.backups.0.restore_regions.0}",
					"source_db_instance_name":    "${data.alicloud_rds_cross_region_backups.backups.db_instance_id}",
					"security_group_ids":         []string{},
					"auto_upgrade_minor_version": "Manual",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "delayed_insert_timeout",
							"value": "70",
						},
					},
					"encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "8.0",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"db_instance_storage_type":   "local_ssd",
						"instance_name":              "tf-testAccDdrDBInstanceConfig",
						"monitoring_period":          "60",
						"zone_id":                    CHECKSET,
						"payment_type":               "PayAsYouGo",
						"connection_string":          CHECKSET,
						"port":                       CHECKSET,
						"security_group_ids.#":       "0",
						"auto_upgrade_minor_version": "Manual",
						"parameters.#":               "1",
						"restore_type":               CHECKSET,
						"restore_time":               CHECKSET,
						"source_region":              CHECKSET,
						"source_db_instance_name":    CHECKSET,
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
		},
	})
}

func resourceDdrDBInstanceConfigDependence(name string) string {
	startTime := time.Now().AddDate(0, 0, -2).Format("2006-01-02T15:04:05Z")
	endTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "startTime" {
 default = "%v"
}

variable "endTime" {
 default = "%v"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.5.id
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
  zone_id = data.alicloud_db_zones.default.zones.5.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.5
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.5
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

data "alicloud_rds_cross_region_backups" "backups" {
  db_instance_id = "rm-xxx"
  start_time = var.startTime
  end_time = var.endTime
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status               = "Enabled"
}
`, name, startTime, endTime)
}

func SkipTestAccAlicloudRdsDdrInstanceSQLServer(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_rds_ddr_instance.default"
	ra := resourceAttrInit(resourceId, ddrinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBDdrInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdrDBInstanceSQLServerConfigDependence)
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
					"engine_version":           "2016_ent_ha",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.19.instance_class}",
					"instance_storage":         "100",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"restore_type":             "BackupSet",
					"backup_set_id":            "${data.alicloud_rds_cross_region_backups.backups.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2016_ent_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         "100",
						"db_instance_storage_type": "cloud_essd",
						"restore_type":             "BackupSet",
						"backup_set_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
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
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.20.instance_class}",
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
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInDdrMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2016_ent_ha",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.19.instance_class}",
					"instance_storage":         "100",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"restore_type":             "BackupSet",
					"backup_set_id":            "${data.alicloud_rds_cross_region_backups.backups.ids.0}",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_group_ids":       []string{"${alicloud_security_group.default.0.id}"},
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2016_ent_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         "100",
						"db_instance_storage_type": "cloud_essd",
						"instance_name":            "tf-testAccDBDdrInstanceConfig",
						"monitoring_period":        "60",
						"zone_id":                  CHECKSET,
						"payment_type":             "PayAsYouGo",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_ids.#":     "1",
						"restore_type":             CHECKSET,
						"backup_set_id":            CHECKSET,
					}),
				),
			},
		},
	})
}

func SkipTestAccAlicloudRdsDdrInstanceSQLServerTime(t *testing.T) {
	var instance map[string]interface{}
	var ips []map[string]interface{}

	resourceId := "alicloud_rds_ddr_instance.default"
	ra := resourceAttrInit(resourceId, ddrinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBDdrInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdrDBInstanceSQLServerConfigDependence)
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
					"engine_version":           "2016_ent_ha",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.19.instance_class}",
					"instance_storage":         "100",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"restore_type":             "BackupTime",
					"restore_time":             "${data.alicloud_rds_cross_region_backups.backups.backups.0.recovery_end_time}",
					"source_region":            "${data.alicloud_rds_cross_region_backups.backups.backups.0.restore_regions.0}",
					"source_db_instance_name":  "${data.alicloud_rds_cross_region_backups.backups.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2016_ent_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         "100",
						"db_instance_storage_type": "cloud_essd",
						"restore_type":             "BackupTime",
						"restore_time":             CHECKSET,
						"source_region":            CHECKSET,
						"source_db_instance_name":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
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
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes.20.instance_class}",
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
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInDdrMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2016_ent_ha",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.19.instance_class}",
					"instance_storage":         "100",
					"db_instance_storage_type": "cloud_essd",
					"payment_type":             "PayAsYouGo",
					"restore_type":             "BackupTime",
					"restore_time":             "${data.alicloud_rds_cross_region_backups.backups.backups.0.recovery_end_time}",
					"source_region":            "${data.alicloud_rds_cross_region_backups.backups.backups.0.restore_regions.0}",
					"source_db_instance_name":  "${data.alicloud_rds_cross_region_backups.backups.db_instance_id}",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"security_group_ids":       []string{"${alicloud_security_group.default.0.id}"},
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2016_ent_ha",
						"instance_type":            CHECKSET,
						"instance_storage":         "100",
						"db_instance_storage_type": "cloud_essd",
						"instance_name":            "tf-testAccDBDdrInstanceConfig",
						"monitoring_period":        "60",
						"zone_id":                  CHECKSET,
						"payment_type":             "PayAsYouGo",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_ids.#":     "1",
						"restore_type":             "BackupTime",
						"restore_time":             CHECKSET,
						"source_region":            CHECKSET,
						"source_db_instance_name":  CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceDdrDBInstanceSQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "SQLServer"
	engine_version = "2016_ent_ha"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.2.id
	engine = "SQLServer"
	engine_version = "2016_ent_ha"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
}

data "alicloud_vpcs" "default" {
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[length(data.alicloud_vpcs.default.ids)-1]
  zone_id = data.alicloud_db_zones.default.zones.2.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids[length(data.alicloud_vpcs.default.ids)-1]
 zone_id = data.alicloud_db_zones.default.ids.2
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_rds_cross_region_backups" "backups" {
  db_instance_id = "rm-xxx"
  start_time = var.startTime
  end_time = var.endTime
}

resource "alicloud_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func testAccCheckDdrSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
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

func testAccCheckKeyValueInDdrMaps(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

var ddrinstanceBasicMap = map[string]string{
	"engine":            "SQLServer",
	"engine_version":    "2016_ent_ha",
	"instance_type":     CHECKSET,
	"instance_storage":  "100",
	"instance_name":     "tf-testAccDBDdrInstanceConfig",
	"monitoring_period": "60",
	"zone_id":           CHECKSET,
	"payment_type":      "PayAsYouGo",
	"connection_string": CHECKSET,
	"port":              CHECKSET,
}

var ddrinstanceBasicMap2 = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "8.0",
	"instance_type":     CHECKSET,
	"instance_storage":  "5",
	"instance_name":     "tf-testAccDdrDBInstanceConfig",
	"monitoring_period": "60",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
}
