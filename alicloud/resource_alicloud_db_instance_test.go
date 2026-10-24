package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass template_id_list -- pre-existing attribute outside this change (PR adds force_encryption, does not modify template_id_list)

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
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

			if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
				if err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
						return retry.RetryableError(err)
					}
					if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidRwSplitNetType.NotFound"}) {
						return nil
					}
					return retry.NonRetryableError(err)
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
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})

		action = "DeleteDBInstance"
		request = map[string]interface{}{
			"DBInstanceId": id,
			"SourceIp":     client.SourceIp,
		}
		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			if err != nil && !NotFoundError(err) {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) || NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	// Mysql_8_0 exercises encryption_key in Step 18, so it uses the KMS
	// fixture (resourceDBInstanceConfigDependence) which declares
	// alicloud_kms_key.default. The KMS key omits dkms_instance_id so it
	// is a standard key that does not require a pre-existing DKMS instance.
	// Tests that do not touch encryption use resourceDBInstanceConfigDependenceNoKMS.
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.n1e.small.1",
					"instance_storage":         "20",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_ssd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"db_is_ignore_case":        "false",
					"auto_create_proxy":        "false",
					"auto_pay":                 "false",
					"connection_mode":          "Standard",
					"force":                    "No",
					"sql_collector_status":     "Disabled",
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
						"connection_mode":            "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sql_collector_status":     "Enabled",
					"ssl_connection_string":    "",
					"db_instance_storage_type": "cloud_essd",
					"io_acceleration_enabled":  "0",
					"compression_mode":         "off",
					"server_cert":              "",
					"server_key":               "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "25",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "mysql.n2e.medium.1",
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
				ImportStateVerifyIgnore: []string{"force_restart", "force", "db_is_ignore_case", "auto_create_proxy", "auto_pay", "business_info", "promotion_code", "user_backup_id", "allow_major_version_upgrade"},
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
					"instance_type":              "mysql.n1e.small.1",
					"instance_storage":           "60",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "force", "db_is_ignore_case", "auto_create_proxy", "auto_pay", "business_info", "promotion_code", "user_backup_id", "parameters", "encryption_key", "security_group_id", "storage_auto_scale", "storage_threshold", "storage_upper_bound", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

// Standard KMS key (dkms_instance_id omitted) so creation does not depend on
// a pre-existing DKMS instance in the test account. The alicloud_kms_service
// data source above calls OpenKmsService to enable KMS before the key is
// created, making this fixture self-contained for encryption_key tests.
resource "alicloud_kms_key" "default" {
  description             = var.name
  pending_window_in_days  = 7
  status                  = "Enabled"
}
`, name)
}

// resourceDBInstanceConfigDependenceNoKMS is identical to
// resourceDBInstanceConfigDependence but omits the KMS data sources and key
// resource. Tests that do not exercise encryption (encryption_key /
// tde_encryption_key) use this fixture to avoid requiring a KMS instance in
// the test account.
func resourceDBInstanceConfigDependenceNoKMS(name string) string {
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

data "alicloud_db_instance_classes" "default" {
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

// resourceDBInstanceConfigDependenceWhitelistTemplate extends the non-KMS
// fixture with two whitelist template resources used by the template_id_list
// reorder convergence test.
func resourceDBInstanceConfigDependenceWhitelistTemplate(name string) string {
	return resourceDBInstanceConfigDependenceNoKMS(name) + fmt.Sprintf(`
resource "alicloud_rds_whitelist_template" "first" {
  template_name = "%s-first"
  ip_white_list = "10.10.10.10/32"
}

resource "alicloud_rds_whitelist_template" "second" {
  template_name = "%s-second"
  ip_white_list = "10.10.20.20/32"
  depends_on = [alicloud_rds_whitelist_template.first]
}
`, name, name)
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                     "MySQL",
					"engine_version":             "5.7",
					"instance_type":              "rds.mysql.t1.small",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes[1].storage_range[0].min}",
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
						"sql_collector_config_value": "30",
						"sql_collector_status":       "Enabled",
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
				ImportStateVerifyIgnore: []string{"force", "force_restart", "db_is_ignore_case", "tde_status", "role_arn", "tde_encryption_key", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_instances" "default" {
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status            = "Enabled"
  dkms_instance_id        = data.alicloud_kms_instances.default.ids.0
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "5.7",
					"instance_type":            "rds.mysql.t1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[1].storage_range[0].min}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "effective_time", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                     "SQLServer",
					"engine_version":             "2012_std_ha",
					"instance_type":              "mssql.x4.medium.s2",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
					"db_instance_storage_type":   "cloud_essd",
					"instance_charge_type":       "Postpaid",
					"instance_name":              "${var.name}",
					"vswitch_id":                 "${local.vswitch_id}",
					"monitoring_period":          "60",
					"category":                   "HighAvailability",
					"time_zone":                  "China Standard Time",
					"collation":                  "Chinese_PRC_CI_AS",
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "SQLServer",
						"engine_version":             "2012_std_ha",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"db_instance_storage_type":   "cloud_essd",
						"category":                   "HighAvailability",
						"time_zone":                  "China Standard Time",
						"collation":                  "Chinese_PRC_CI_AS",
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "30",
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
					"instance_type": "mssql.x4.medium.s2",
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
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "180",
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
					"engine_version": "2016_std_ha",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "2016_std_ha",
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
	security_group_name   = var.name
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                     "PostgreSQL",
					"engine_version":             "12.0",
					"instance_type":              "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
					"db_instance_storage_type":   "cloud_essd",
					"zone_id":                    "${data.alicloud_db_instance_classes.default.instance_classes[0].zone_ids[0].id}",
					"instance_charge_type":       "Postpaid",
					"instance_name":              "${var.name}",
					"vswitch_id":                 "${local.vswitch_id}",
					"monitoring_period":          "60",
					"category":                   "HighAvailability",
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "PostgreSQL",
						"engine_version":             "12.0",
						"instance_storage":           CHECKSET,
						"instance_type":              CHECKSET,
						"db_instance_storage_type":   "cloud_essd",
						"category":                   "HighAvailability",
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "30",
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
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes[1].instance_class}",
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
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "180",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "12.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_instances" "default" {
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status            = "Enabled"
  dkms_instance_id        = data.alicloud_kms_instances.default.ids.0
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "13.0",
					"instance_type":            "pg.x2.medium.2c",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes[0].zone_ids[0].id}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "allow_major_version_upgrade"},
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
					"instance_storage":            "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min + data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].step}",
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
			{
				// node_id targets a specific HA node and the server
				// certificate pair must be generated offline per the suite
				// notes above; removing the required instance_type stops the
				// plan before the placeholders reach the API.
				Config: testAccConfig(map[string]interface{}{
					"force":         "No",
					"node_id":       "8fi5v8m2s38",
					"server_cert":   "-----BEGIN CERTIFICATE-----tfacc-placeholder-----END CERTIFICATE-----",
					"server_key":    "-----BEGIN PRIVATE KEY-----tfacc-placeholder-----END PRIVATE KEY-----",
					"instance_type": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				// Restore a destroy-valid config after the preceding ExpectError
				// step: re-add instance_type (Required) and drop the placeholder
				// force/node_id/server_cert/server_key it introduced, so the
				// post-test destroy does not fail with
				// "Missing required argument: instance_type".
				Config: testAccConfig(map[string]interface{}{
					"force":         REMOVEKEY,
					"node_id":       REMOVEKEY,
					"server_cert":   REMOVEKEY,
					"server_key":    REMOVEKEY,
					"instance_type": "pg.x8.medium.2c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
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
	security_group_name   = var.name
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "15.0",
					"instance_type":            "pg.x4.large.4c",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes[0].zone_ids[0].id}",
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
					"instance_storage":            "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min + data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].step}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "db_is_ignore_case", "fresh_white_list_readins", "released_keep_policy",
					"client_ca_enabled", "client_crl_enabled", "db_instance_ip_array_name", "encryption_key", "security_group_id", "modify_mode", "security_ip_type",
					"whitelist_network_type", "babelfish_config.2289427611.babelfish_enabled", "babelfish_config.2289427611.master_user_password", "babelfish_config.2289427611.master_username",
					"babelfish_config.2289427611.migration_mode", "allow_major_version_upgrade"},
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
  	engine_version       = "15.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "15.0"
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MariaDB",
					"engine_version":           "10.3",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
					"db_instance_storage_type": "cloud_essd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes[0].zone_ids[0].id}",
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
					"instance_type": "${data.alicloud_db_instance_classes.default.instance_classes[1].instance_class}",
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
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min + data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].step}",
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_instances" "default" {
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description            = var.name
  pending_window_in_days = 7
  status                 = "Enabled"
  dkms_instance_id       = data.alicloud_kms_instances.default.ids.0
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "rds.mysql.t1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[1].storage_range[0].min}",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes[0].zone_ids[0].id}",
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "rds.mysql.s1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[1].storage_range[0].min}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "encryption_key", "db_is_ignore_case", "allow_major_version_upgrade"},
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

func TestAccAliCloudRdsDBInstance_Mysql_8_0_PrePaidForceDelete(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfigForceDelete"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceHighAvailabilityConfigDependence1)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "rds.mysql.s1.small",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.1.storage_range.min}",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"force_delete":             "true",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
					"zone_id":                  "${local.zone_id}",
					"vswitch_id":               "${local.vswitch_id}",
					"monitoring_period":        "60",
					"security_group_ids":       "${alicloud_security_group.default.*.id}",
					"security_ips":             []string{"10.168.1.12", "100.69.7.112"},
					"db_time_zone":             "America/New_York",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"engine_version":       "8.0",
						"instance_charge_type": "Prepaid",
						"force_delete":         "true",
						"instance_name":        "tf-testAccDBInstanceConfigForceDelete",
						"instance_storage":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "period", "force_delete", "encryption_key", "db_is_ignore_case", "allow_major_version_upgrade"},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "mysql.x2.medium.xc",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[1].storage_range[0].min}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "encryption_key", "db_is_ignore_case", "allow_major_version_upgrade"},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
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
					"instance_type": "${data.alicloud_db_instance_classes.this.instance_classes[0].instance_class}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "14.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
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
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

`, name)
}

func TestAccAliCloudRdsDBInstance_SQLServer_2019_ServerlessHA(t *testing.T) {
	t.Skipf("Skipping test: ssl_certificate, ssl_password, tde_certificate, tde_private_key, tde_password, tde_db_name require customer-generated certificates offline")

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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "SQLServer",
					"engine_version":           "2019_std_sl",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes[0].instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min}",
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
					"ssl_action":      "Open",
					"ssl_certificate": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.pfx",
					"ssl_password":    "123456",
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
				Config: testAccConfig(map[string]interface{}{
					"tde_db_name":     "test01",
					"tde_password":    "cMr2O6BnU3XFMlHu86/DHbOlbGFgo5prUrGCbcARI/E=",
					"tde_private_key": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.pvk",
					"tde_certificate": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.cer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_db_name":     "test01",
						"tde_password":    "cMr2O6BnU3XFMlHu86/DHbOlbGFgo5prUrGCbcARI/E=",
						"tde_private_key": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.pvk",
						"tde_certificate": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.cer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_db_name":     "test02",
					"tde_password":    "cMr2O6BnU3XFMlHu86/DHbOlbGFgo5prUrGCbcARI/E=",
					"tde_private_key": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.pvk",
					"tde_certificate": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.cer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_db_name":     "test02",
						"tde_password":    "cMr2O6BnU3XFMlHu86/DHbOlbGFgo5prUrGCbcARI/E=",
						"tde_private_key": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.pvk",
						"tde_certificate": "oss-cn-beijing.aliyuncs.com:sjcbeijing:rdsTest.cer",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"target_minor_version":     "rds_20201031",
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
						"target_minor_version":       "rds_20201031",
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
			// UpgradeDBInstanceKernelVersion
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_time":         "Immediate",
					"target_minor_version": "rds_20260228",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "rds_20260228",
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period", "upgrade_time", "allow_major_version_upgrade"},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period", "allow_major_version_upgrade"},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min}",
					"category":                 "HighAvailability",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "cloud_essd",
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
						"db_instance_storage_type":   "cloud_essd",
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "direction", "auto_renew", "auto_renew_period", "allow_major_version_upgrade"},
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
	security_group_name   = var.name
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

// Standard KMS key (dkms_instance_id omitted) so creation does not depend on
// a pre-existing DKMS instance in the test account. The alicloud_kms_service
// data source above calls OpenKmsService to enable KMS before the key is
// created, making this fixture self-contained for encryption_key tests.
resource "alicloud_kms_key" "default" {
  description             = var.name
  pending_window_in_days  = 7
  status                  = "Enabled"
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
   security_group_name   = var.name
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
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_kms_instances" "default" {
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  key_state               = "Enabled"
  dkms_instance_id        = data.alicloud_kms_instances.default.ids.0
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

data "alicloud_db_instance_classes" "default" {
	zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
	category = "HighAvailability"
	db_instance_storage_type = "general_essd"
	instance_charge_type = "PostPaid"
}

locals {
	# Storage compression requires a dedicated (x4) instance family; the class is
	# resolved from the data source so a retired class or a sold-out zone does not
	# turn into a non-semantic test failure.
	compression_instance_class = [for c in data.alicloud_db_instance_classes.default.instance_classes : c.instance_class if can(regex("^mysql\\.x4\\.", c.instance_class))][0]
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "encryption_key", "direction", "auto_renew", "auto_renew_period", "allow_major_version_upgrade"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstanceMysql_general_essd_supplement(t *testing.T) {
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "encryption_key", "direction", "auto_renew", "auto_renew_period", "allow_major_version_upgrade"},
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"template_id_list":         []string{"${alicloud_rds_whitelist_template.default1.id}", "${alicloud_rds_whitelist_template.default.id}"},
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
				Config: testAccConfig(map[string]interface{}{
					"engine_version":    "18.0",
					"collect_stat_mode": "Before",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "18.0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "tde_encryption_key", "collect_stat_mode", "allow_major_version_upgrade"},
			},
		},
	})
}

func TestAccAliCloudRdsDBInstance_PG_MajorVersionUpgrade_Prepaid(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap7)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceConfig%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePGMajorVersionUpgradePrepaidDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "PostgreSQL",
					"engine_version":           "17.0",
					"instance_type":            "pg.n4.2c.2m",
					"instance_storage":         "30",
					"instance_charge_type":     "Prepaid",
					"period":                   "1",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "17.0",
						"instance_charge_type": "Prepaid",
						"instance_name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":    "18.0",
					"collect_stat_mode": "Before",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "18.0",
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
				ImportStateVerifyIgnore: []string{"force_restart", "period", "collect_stat_mode", "auto_renew", "auto_renew_period", "allow_major_version_upgrade"},
			},
		},
	})
}

func resourceDBInstancePGMajorVersionUpgradePrepaidDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_zones" "default" {
	engine         = "PostgreSQL"
	engine_version = "17.0"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_db_zones.default.zones.0.id
}
`, name)
}

func resourceDBInstanceConfigGeneralEssdPgSql(name string) string {
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

data "alicloud_kms_instances" "default" {
}

data "alicloud_kms_service" "default" {
  enable = "On"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  status            = "Enabled"
  dkms_instance_id        = data.alicloud_kms_instances.default.ids.0
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
			},
		},
	})
}

func resourceDBInstanceConfigGeneralEssdSqlService(name string) string {
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

// TestAccAliCloudRdsDBInstance_networkAttributes verifies the network related
// attributes that are read back from DescribeDBInstanceAttribute
// (network_type from InstanceNetworkType, db_instance_net_type from
// DBInstanceNetType) and the write-only allow_major_version_upgrade switch.
func TestAccAliCloudRdsDBInstance_networkAttributes(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"engine":               "MySQL",
		"engine_version":       "8.0",
		"instance_type":        CHECKSET,
		"instance_storage":     CHECKSET,
		"instance_name":        CHECKSET,
		"instance_charge_type": "Postpaid",
		"status":               CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceNet%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceNoKMS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min}",
					"category":                 "HighAvailability",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":         "VPC",
						"db_instance_net_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_major_version_upgrade": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_major_version_upgrade": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "allow_major_version_upgrade"},
			},
		},
	})
}

// TestUnitParseDBInstanceParamGroupId is a regression test for the provider
// panic reported when creating alicloud_db_instance.
//
// resourceAliCloudDBInstanceRead used to run, without a comma-ok guard:
//
//	dbParamGroupInfo := response["ParamGroupInfo"].(map[string]interface{})
//	...
//	d.Set("db_param_group_id", dbParamGroupInfo["ParamGroupId"].(string))
//
// A freshly created instance whose DescribeParameters response does not yet
// contain "ParamGroupInfo" therefore crashed the plugin with:
//
//	interface conversion: interface {} is nil, not map[string]interface {}
//
// The parsing is now isolated in parseDBInstanceParamGroupId, which must never
// panic and must report presence via its second return value. The first two
// cases below are the crash inputs (they panic against the old inline code).
func TestUnitParseDBInstanceParamGroupId(t *testing.T) {
	cases := []struct {
		name      string
		response  map[string]interface{}
		wantID    string
		wantFound bool
	}{
		{
			name:     "ParamGroupInfo absent (freshly created instance)",
			response: map[string]interface{}{},
		},
		{
			name:     "ParamGroupInfo present but nil",
			response: map[string]interface{}{"ParamGroupInfo": nil},
		},
		{
			name:     "ParamGroupInfo present but ParamGroupId missing",
			response: map[string]interface{}{"ParamGroupInfo": map[string]interface{}{}},
		},
		{
			name:     "ParamGroupId present but not a string",
			response: map[string]interface{}{"ParamGroupInfo": map[string]interface{}{"ParamGroupId": 123}},
		},
		{
			name:      "ParamGroupInfo and ParamGroupId both present",
			response:  map[string]interface{}{"ParamGroupInfo": map[string]interface{}{"ParamGroupId": "rpg-test123"}},
			wantID:    "rpg-test123",
			wantFound: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, ok := parseDBInstanceParamGroupId(c.response)
			if ok != c.wantFound {
				t.Fatalf("parseDBInstanceParamGroupId(%v) found = %v, want %v", c.response, ok, c.wantFound)
			}
			if id != c.wantID {
				t.Fatalf("parseDBInstanceParamGroupId(%v) id = %q, want %q", c.response, id, c.wantID)
			}
		})
	}
}

type fakeFieldChangeReporter map[string]bool

func (f fakeFieldChangeReporter) HasChange(field string) bool {
	return f[field]
}

func TestUnitDBInstanceCreateOnlyFieldsUnchanged(t *testing.T) {
	if err := validateDBInstanceCreateOnlyFieldsUnchanged(fakeFieldChangeReporter{}); err != nil {
		t.Fatalf("expected no error when no create-only field changes, got %v", err)
	}
	for _, field := range dbInstanceCreateOnlyFields {
		err := validateDBInstanceCreateOnlyFieldsUnchanged(fakeFieldChangeReporter{field: true})
		if err == nil {
			t.Fatalf("expected an error when %s changes on an existing instance", field)
		}
		if !strings.Contains(err.Error(), "`"+field+"`") {
			t.Fatalf("error %q does not name %s", err, field)
		}
	}
}

// TestUnitDBInstanceBoolFromConfig pins the replacement for the deprecated
// GetOkExists on auto_create_proxy/auto_pay: an explicit false must be sent while
// an omitted or unknown value must not be sent at all.
func TestUnitDBInstanceBoolFromConfig(t *testing.T) {
	cases := []struct {
		name      string
		config    cty.Value
		expected  bool
		present   bool
		attribute string
	}{
		{name: "explicit false", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.False}), attribute: "auto_pay", expected: false, present: true},
		{name: "explicit true", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.True}), attribute: "auto_pay", expected: true, present: true},
		{name: "null value", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.NullVal(cty.Bool)}), attribute: "auto_pay"},
		{name: "unknown value", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.UnknownVal(cty.Bool)}), attribute: "auto_pay"},
		{name: "absent attribute", config: cty.ObjectVal(map[string]cty.Value{"auto_create_proxy": cty.True}), attribute: "auto_pay"},
		{name: "nil config", config: cty.NilVal, attribute: "auto_pay"},
		{name: "null config", config: cty.NullVal(cty.EmptyObject), attribute: "auto_pay"},
	}
	for _, c := range cases {
		got, present := dbInstanceBoolFromConfig(c.config, c.attribute)
		if present != c.present || got != c.expected {
			t.Fatalf("case %s: dbInstanceBoolFromConfig() = (%v, %v), want (%v, %v)", c.name, got, present, c.expected, c.present)
		}
	}
}

// TestUnitNonEmptyStringAttr covers the Read guard used for
// Optional+Computed attributes: an omitted, null or empty API value must not
// overwrite the last known state.
func TestUnitNonEmptyStringAttr(t *testing.T) {
	cases := []struct {
		name     string
		instance map[string]interface{}
		expected string
		present  bool
	}{
		{name: "present", instance: map[string]interface{}{"ConnectionMode": "Safe"}, expected: "Safe", present: true},
		{name: "absent", instance: map[string]interface{}{}},
		{name: "null", instance: map[string]interface{}{"ConnectionMode": nil}},
		{name: "empty", instance: map[string]interface{}{"ConnectionMode": ""}},
	}
	for _, c := range cases {
		got, present := nonEmptyStringAttr(c.instance, "ConnectionMode")
		if got != c.expected || present != c.present {
			t.Fatalf("case %s: nonEmptyStringAttr() = (%q, %v), want (%q, %v)", c.name, got, present, c.expected, c.present)
		}
	}
}

// TestUnitDBInstancePgBouncerEnabledFromInstance covers the Read conversion of
// the PgBouncer flag: only an explicit "true"/"false" from the API is trusted,
// while an omitted or unparsable value is reported as false instead of keeping
// an unconfirmed configured true.
func TestUnitDBInstancePgBouncerEnabledFromInstance(t *testing.T) {
	cases := []struct {
		name     string
		instance map[string]interface{}
		want     bool
	}{
		{name: "enabled", instance: map[string]interface{}{"PGBouncerEnabled": "true"}, want: true},
		{name: "disabled", instance: map[string]interface{}{"PGBouncerEnabled": "false"}},
		{name: "boolean true", instance: map[string]interface{}{"PGBouncerEnabled": true}, want: true},
		{name: "omitted", instance: map[string]interface{}{}},
		{name: "null", instance: map[string]interface{}{"PGBouncerEnabled": nil}},
		{name: "unparsable", instance: map[string]interface{}{"PGBouncerEnabled": "unsupported"}},
	}
	for _, c := range cases {
		if got := pgBouncerEnabledFromInstance(c.instance); got != c.want {
			t.Fatalf("case %s: pgBouncerEnabledFromInstance() = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestUnitDBInstanceUserBackupIdValidation(t *testing.T) {
	legal := map[string]interface{}{
		"engine":               "MySQL",
		"engine_version":       "5.7",
		"instance_charge_type": "Postpaid",
		"category":             "Basic",
	}
	cases := []struct {
		name    string
		raw     map[string]interface{}
		wantErr bool
	}{
		{name: "legal combination", raw: legal},
		{name: "prepaid", raw: copyWith(legal, "instance_charge_type", "Prepaid"), wantErr: true},
		{name: "non MySQL", raw: copyWith(legal, "engine", "PostgreSQL"), wantErr: true},
		{name: "non 5.7", raw: copyWith(legal, "engine_version", "8.0"), wantErr: true},
		{name: "non Basic category", raw: copyWith(legal, "category", "HighAvailability"), wantErr: true},
		{name: "category unset", raw: copyWith(legal, "category", ""), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAliCloudDBInstance().Schema, c.raw)
			err := validateDBInstanceUserBackupId(d)
			if (err != nil) != c.wantErr {
				t.Fatalf("validateDBInstanceUserBackupId() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func copyWith(src map[string]interface{}, key string, value interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	dst[key] = value
	return dst
}

func TestUnitDBInstanceSQLServerMajorVersionUpgradeRequest(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou", SourceIp: "127.0.0.1"}
	base := map[string]interface{}{
		"engine":                      "SQLServer",
		"engine_version":              "2019_std_ha",
		"instance_type":               "mssql.x4.medium.s2",
		"instance_charge_type":        "Postpaid",
		"category":                    "Basic",
		"zone_id":                     "cn-hangzhou-k",
		"vswitch_id":                  "vsw-primary",
		"allow_major_version_upgrade": true,
	}

	build := func(raw map[string]interface{}) (map[string]interface{}, error) {
		d := schema.TestResourceDataRaw(t, resourceAliCloudDBInstance().Schema, raw)
		d.SetId("rm-test123")
		return buildDBInstanceSQLServerMajorVersionUpgradeRequest(d, client)
	}

	t.Run("basic category carries the required parameters", func(t *testing.T) {
		request, err := build(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for key, want := range map[string]interface{}{
			"Category":        "Basic",
			"ZoneId":          "cn-hangzhou-k",
			"VSwitchId":       "vsw-primary",
			"EngineVersion":   "2019_std_ha",
			"DBInstanceId":    "rm-test123",
			"DBInstanceClass": "mssql.x4.medium.s2",
		} {
			if got := request[key]; got != want {
				t.Fatalf("request[%s] = %v, want %v", key, got, want)
			}
		}
		if _, ok := request["ZoneIdSlave1"]; ok {
			t.Fatalf("ZoneIdSlave1 must not be sent for a Basic category instance")
		}
		if token, _ := request["ClientToken"].(string); token == "" {
			t.Fatalf("ClientToken must be set for an idempotent retry")
		}
	})

	t.Run("explicit false is forwarded", func(t *testing.T) {
		request, err := build(copyWith(base, "allow_major_version_upgrade", false))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, ok := request["AllowMajorVersionUpgrade"]
		if !ok || got != false {
			t.Fatalf("AllowMajorVersionUpgrade = %v (present %v), want false", got, ok)
		}
	})

	t.Run("high availability without a slave zone stays single AZ", func(t *testing.T) {
		request, err := build(copyWith(base, "category", "HighAvailability"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := request["ZoneIdSlave1"]; got != "cn-hangzhou-k" {
			t.Fatalf("ZoneIdSlave1 = %v, want the primary zone for a single-AZ deployment", got)
		}
	})

	t.Run("high availability with a different slave zone needs two vswitches", func(t *testing.T) {
		raw := copyWith(copyWith(base, "category", "HighAvailability"), "zone_id_slave_a", "cn-hangzhou-j")
		if _, err := build(raw); err == nil {
			t.Fatalf("expected an error when the slave zone differs but only one vSwitch is configured")
		}
		raw = copyWith(raw, "vswitch_id", "vsw-primary,vsw-standby")
		request, err := build(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := request["ZoneIdSlave1"]; got != "cn-hangzhou-j" {
			t.Fatalf("ZoneIdSlave1 = %v, want cn-hangzhou-j", got)
		}
		if got := request["VSwitchId"]; got != "vsw-primary,vsw-standby" {
			t.Fatalf("VSwitchId = %v, want both vSwitch IDs", got)
		}
	})

	for _, tc := range []struct {
		name  string
		raw   map[string]interface{}
		field string
	}{
		{name: "missing category", raw: copyWith(base, "category", ""), field: "category"},
		{name: "missing zone", raw: copyWith(base, "zone_id", ""), field: "zone_id"},
		{name: "missing vswitch", raw: copyWith(base, "vswitch_id", ""), field: "vswitch_id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := build(tc.raw); err == nil || !strings.Contains(err.Error(), "`"+tc.field+"` is required") {
				t.Fatalf("expected a %s required error, got %v", tc.field, err)
			}
		})
	}
}

func TestUnitDBInstanceConcurrentErrorClassification(t *testing.T) {
	// The codes are listed explicitly (not taken from the production slices) so
	// the test fails when a documented wait-and-retry code is dropped, and so the
	// create-only / modify-only split is asserted in both directions.
	shared := []string{
		"InvalidConcurrentOperate",
		"System.ConcurrentConflict",
		"OperationDenied.SystemConcurrent",
	}
	for _, code := range shared {
		err := &tea.SDKError{Code: tea.String(code), Message: tea.String("concurrent"), Data: tea.String("")}
		if !isDBInstanceConcurrentError(err) || !isDBInstanceCreateWaitError(err) || !isDBInstanceModifyWaitError(err) {
			t.Fatalf("expected %s to be classified on every write surface", code)
		}
	}

	createOnly := &tea.SDKError{Code: tea.String("Order.UnfinishedRefundOrder"), Message: tea.String("refund"), Data: tea.String("")}
	if !isDBInstanceCreateWaitError(createOnly) {
		t.Fatal("Order.UnfinishedRefundOrder must be retried on create")
	}
	if isDBInstanceModifyWaitError(createOnly) {
		t.Fatal("Order.UnfinishedRefundOrder must not be retried on modify")
	}

	modifyOnly := &tea.SDKError{Code: tea.String("OperationDenied.DurationLimit"), Message: tea.String("duration"), Data: tea.String("")}
	if !isDBInstanceModifyWaitError(modifyOnly) {
		t.Fatal("OperationDenied.DurationLimit must be retried on modify")
	}
	if isDBInstanceCreateWaitError(modifyOnly) {
		t.Fatal("OperationDenied.DurationLimit must not be retried on create")
	}

	permanent := &tea.SDKError{Code: tea.String("InvalidDBInstanceId.NotFound"), Message: tea.String("not found"), Data: tea.String("")}
	if isDBInstanceCreateWaitError(permanent) || isDBInstanceModifyWaitError(permanent) || isDBInstanceConcurrentError(permanent) {
		t.Fatal("a permanent error must not be classified as retryable")
	}
	if isDBInstanceCreateWaitError(nil) || isDBInstanceModifyWaitError(nil) {
		t.Fatal("nil must not be classified as retryable")
	}
}

// TestAccAliCloudRdsDBInstance_createOnlyFields verifies that the write-only
// creation parameters reach CreateDBInstance with real values, and that a later
// change to one of them is rejected at plan time instead of silently changing
// only the local state.
func TestAccAliCloudRdsDBInstance_createOnlyFields(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"engine":                   "MySQL",
		"engine_version":           "8.0",
		"instance_charge_type":     "Postpaid",
		"db_instance_storage_type": "cloud_essd",
		"status":                   CHECKSET,
		"connection_string":        CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceCreateOnly%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceNoKMS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min}",
					"category":                 "HighAvailability",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"auto_create_proxy":        "false",
					"auto_pay":                 "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage_type": "cloud_essd",
						"auto_create_proxy":        "false",
						"auto_pay":                 "false",
					}),
				),
			},
			{
				// business_info is exercised through the create-only guard only:
				// CreateDBInstance rejects the value documented for BusinessInfo
				// (InvalidBusinessInfo.Malformed, request id
				// 01A182E3-08A3-5AFF-9D9F-39471F56EBC1), so no valid value could be
				// sent from the test account.
				Config: testAccConfig(map[string]interface{}{
					"business_info": "121436975448952",
				}),
				ExpectError: regexp.MustCompile("`business_info` can only be specified when the DB instance is created"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// Drop business_info so this step fails on promotion_code only:
					// the step config map is accumulated across steps.
					"business_info":  REMOVEKEY,
					"promotion_code": "aliwood-1688-mobile-promotion",
				}),
				ExpectError: regexp.MustCompile("`promotion_code` can only be specified when the DB instance is created"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"promotion_code": REMOVEKEY,
					"user_backup_id": "backup-121436975448952",
				}),
				ExpectError: regexp.MustCompile("`user_backup_id` can only be specified when the DB instance is created"),
			},
			{
				// The boolean write-only fields must be rejected on an existing
				// instance through the same plan-time guard, so a later apply can
				// never silently rewrite only the local state.
				Config: testAccConfig(map[string]interface{}{
					"user_backup_id":    REMOVEKEY,
					"auto_create_proxy": "true",
				}),
				ExpectError: regexp.MustCompile("`auto_create_proxy` can only be specified when the DB instance is created"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// Restore auto_create_proxy so this step fails on auto_pay only:
					// the step config map is accumulated across steps and the
					// previous step was rejected without changing state.
					"auto_create_proxy": "false",
					"auto_pay":          "true",
				}),
				ExpectError: regexp.MustCompile("`auto_pay` can only be specified when the DB instance is created"),
			},
		},
	})
}

// TestUnitDBInstanceTemplateIdListOrderSuppress exercises the bounded-suppressor
// unit net for template_id_list. The AttachWhitelistTemplate/Detach operations
// use set-based comparison, so order-only changes are semantically meaningless.
// The cases run the same algorithm the SDK reaches through the DiffSuppressFunc
// (templateIdListEqualIgnoringOrder with the values from d.GetChange), so the
// production path is verified.
func TestUnitDBInstanceTemplateIdListOrderSuppress(t *testing.T) {
	listCases := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		{"ReorderOnly", []interface{}{1, 2}, []interface{}{2, 1}, true},
		{"ReorderOnlyLonger", []interface{}{1, 2, 3}, []interface{}{3, 1, 2}, true},
		{"MemberAdded", []interface{}{1, 2}, []interface{}{1, 2, 3}, false},
		{"MemberRemoved", []interface{}{1, 2, 3}, []interface{}{1, 2}, false},
		{"MemberReplaced", []interface{}{1, 2}, []interface{}{1, 3}, false},
		{"DuplicateCollapses", []interface{}{1, 1}, []interface{}{1}, false},
		{"DuplicateReorder", []interface{}{1, 1, 2}, []interface{}{2, 1, 1}, true},
		{"SameOrder", []interface{}{1, 2}, []interface{}{1, 2}, true},
		{"Empty", []interface{}{}, []interface{}{}, true},
	}
	for _, tc := range listCases {
		if got := templateIdListEqualIgnoringOrder(tc.old, tc.new); got != tc.expected {
			t.Fatalf("list case %s: got %v, expected %v", tc.name, got, tc.expected)
		}
	}
}

// TestAccAliCloudRdsDBInstance_templateIdListOrder verifies that reordering
// template_id_list members converges: apply members in one order, then apply
// the same members in a different order, and require the post-apply plan to be
// empty. The schema carries a DiffSuppressFunc because the RDS
// AttachWhitelistTemplate/DetachWhitelistTemplate operations use set-based
// comparison and member order is semantically meaningless.
func TestAccAliCloudRdsDBInstance_templateIdListOrder(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	name := fmt.Sprintf("tf-testAccDBInstanceTplOrder%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceWhitelistTemplate)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min}",
					"category":                 "HighAvailability",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${local.vswitch_id}",
					"db_instance_storage_type": "cloud_essd",
					"template_id_list": []interface{}{
						"${alicloud_rds_whitelist_template.first.id}",
						"${alicloud_rds_whitelist_template.second.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "engine", "MySQL"),
					resource.TestCheckResourceAttr(resourceId, "engine_version", "8.0"),
					resource.TestCheckResourceAttr(resourceId, "template_id_list.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_id_list": []interface{}{
						"${alicloud_rds_whitelist_template.second.id}",
						"${alicloud_rds_whitelist_template.first.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "template_id_list.#", "2"),
				),
			},
		},
	})
}

// TestAccAliCloudRdsDBInstance_compressionMode verifies that
// compression_mode="on" is applied during create via ModifyDBInstanceSpec
// (RDS-02: CreateDBInstance does not honour compression_mode). Uses a
// HighAvailability dedicated instance (mysql.x4.medium.2c) with general_essd
// and 1000 GB storage — the minimum prerequisites for RDS storage
// compression. Once enabled, compression cannot be disabled, so no on→off
// update is performed.
func TestAccAliCloudRdsDBInstance_compressionMode(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceCompression%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDBGeneralEssd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${local.compression_instance_class}",
					"instance_storage":         "1000",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
					"compression_mode":         "on",
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
						"compression_mode":         "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "force", "db_is_ignore_case", "auto_create_proxy", "auto_pay", "business_info", "promotion_code", "user_backup_id", "allow_major_version_upgrade"},
			},
		},
	})
}

// TestAccAliCloudRdsDBInstance_compressionModeUpdate verifies the in-place
// off -> on path: an existing instance created with compression disabled is
// updated to "on" and must converge (the post-apply plan is required to be
// empty). on -> off is not exercised because the server does not support
// disabling storage compression.
func TestAccAliCloudRdsDBInstance_compressionModeUpdate(t *testing.T) {
	var instance map[string]interface{}
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBInstanceCompressionUpd%d", rand.Intn(1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependenceDBGeneralEssd)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"instance_type":            "${local.compression_instance_class}",
					"instance_storage":         "1000",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"db_instance_storage_type": "general_essd",
					"compression_mode":         "off",
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
						"compression_mode":         "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compression_mode": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compression_mode": "on",
					}),
				),
			},
		},
	})
}

// dbInstanceMockClient builds an AliyunClient whose RDS endpoint is an isolated
// in-process HTTP server, so request-level tests exercise the same code path
// Terraform uses without touching the cloud.
func dbInstanceMockClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("rds", endpoint)
	return client
}

func dbInstanceRawConfig(values map[string]interface{}) cty.Value {
	attrs := map[string]cty.Value{}
	for key, value := range values {
		switch typed := value.(type) {
		case bool:
			attrs[key] = cty.BoolVal(typed)
		case int:
			attrs[key] = cty.NumberIntVal(int64(typed))
		case string:
			attrs[key] = cty.StringVal(typed)
		default:
			attrs[key] = cty.NullVal(cty.DynamicPseudoType)
		}
	}
	return cty.ObjectVal(attrs)
}

func dbInstanceMockCreateConfig(extra map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"engine":                   "MySQL",
		"engine_version":           "8.0",
		"instance_type":            "mysql.n2.small.2c",
		"instance_storage":         20,
		"category":                 "Basic",
		"instance_charge_type":     "Postpaid",
		"db_instance_storage_type": "cloud_ssd",
	}
	for key, value := range extra {
		config[key] = value
	}
	return config
}

// TestUnitDBInstanceCreateWireWriteOnlyBooleans drives the public Create entry
// point against an isolated HTTP endpoint and asserts the serialized
// CreateDBInstance form. An omitted value must not be sent (the API default is
// true), while an explicit false/true must reach the wire unchanged.
func TestUnitDBInstanceCreateWireWriteOnlyBooleans(t *testing.T) {
	cases := []struct {
		name        string
		extra       map[string]interface{}
		want        map[string]string
		wantMissing []string
	}{
		{
			name:        "omitted",
			extra:       map[string]interface{}{},
			wantMissing: []string{"AutoPay", "AutoCreateProxy"},
		},
		{
			name:  "explicit false",
			extra: map[string]interface{}{"auto_pay": false, "auto_create_proxy": false},
			want:  map[string]string{"AutoPay": "false", "AutoCreateProxy": "false"},
		},
		{
			name:  "explicit true",
			extra: map[string]interface{}{"auto_pay": true, "auto_create_proxy": true},
			want:  map[string]string{"AutoPay": "true", "AutoCreateProxy": "true"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := dbInstanceMockCreateConfig(tc.extra)
			var captured url.Values
			client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				if action := r.Form.Get("Action"); action != "CreateDBInstance" {
					t.Errorf("unexpected action %q", action)
				}
				captured = r.Form
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Code":    "InvalidParameter",
					"Message": "stop after request capture",
				})
			})

			res := resourceAliCloudDBInstance()
			diff, err := res.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("diff failed: %v", err)
			}
			if diff == nil {
				t.Fatal("expected a create diff")
			}
			diff.RawConfig = dbInstanceRawConfig(config)
			if _, diags := res.Apply(context.Background(), nil, diff, client); !diags.HasError() {
				t.Fatal("expected the mocked CreateDBInstance failure")
			}
			if captured == nil {
				t.Fatal("CreateDBInstance was not called")
			}
			for key, want := range tc.want {
				if got := captured.Get(key); got != want {
					t.Fatalf("%s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.wantMissing {
				if _, ok := captured[key]; ok {
					t.Fatalf("%s must not be sent when omitted, got %q", key, captured.Get(key))
				}
			}
		})
	}
}

// TestUnitDBInstanceCreateOnlyFieldsPlanGuard drives the real CustomizeDiff
// path with a prior state, so the create-only rejection is proven for the
// boolean write-only fields and not only through a fake HasChange reporter.
func TestUnitDBInstanceCreateOnlyFieldsPlanGuard(t *testing.T) {
	cases := []struct {
		field string
		from  interface{}
		to    interface{}
	}{
		{field: "auto_create_proxy", from: false, to: true},
		{field: "auto_pay", from: false, to: true},
		{field: "business_info", from: "", to: "info-1"},
		{field: "promotion_code", from: "", to: "CODE-1"},
		{field: "user_backup_id", from: "", to: "backup-1"},
	}

	baseAttributes := func(field string, value interface{}) map[string]string {
		return map[string]string{
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"instance_type":            "mysql.n2.small.2c",
			"instance_storage":         "20",
			"category":                 "Basic",
			"instance_charge_type":     "Postpaid",
			"db_instance_storage_type": "cloud_ssd",
			"connection_mode":          "Standard",
			field:                      fmt.Sprint(value),
		}
	}

	for _, tc := range cases {
		t.Run(tc.field+"/in_place_rejected", func(t *testing.T) {
			config := dbInstanceMockCreateConfig(map[string]interface{}{tc.field: tc.to})
			state := &sdkterraform.InstanceState{ID: "rm-existing", Attributes: baseAttributes(tc.field, tc.from)}

			res := resourceAliCloudDBInstance()
			_, err := res.Diff(context.Background(), state, sdkterraform.NewResourceConfigRaw(config), nil)
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("expected a plan-time rejection for %s, got %v", tc.field, err)
			}
		})

		t.Run(tc.field+"/replacement_allowed", func(t *testing.T) {
			// An independent ForceNew change replaces the instance, so the
			// create-only parameter is applied to the new instance and the plan
			// must not be rejected.
			config := dbInstanceMockCreateConfig(map[string]interface{}{
				tc.field:          tc.to,
				"connection_mode": "Safe",
			})
			state := &sdkterraform.InstanceState{ID: "rm-existing", Attributes: baseAttributes(tc.field, tc.from)}

			res := resourceAliCloudDBInstance()
			if _, err := res.Diff(context.Background(), state, sdkterraform.NewResourceConfigRaw(config), nil); err != nil {
				t.Fatalf("a replacement plan must accept the create-only change on %s, got %v", tc.field, err)
			}
		})
	}
}

// TestUnitDBInstanceForceNewFieldListMatchesSchema keeps the replacement list in
// sync with the schema: a ForceNew attribute missing from the list would bring
// back the create-only false rejection, and a listed field that lost ForceNew
// would silently bypass the guard.
func TestUnitDBInstanceForceNewFieldListMatchesSchema(t *testing.T) {
	schemaMap := resourceAliCloudDBInstance().Schema
	listed := make(map[string]bool, len(dbInstanceForceNewFields))
	for _, field := range dbInstanceForceNewFields {
		if _, ok := schemaMap[field]; !ok {
			t.Fatalf("listed ForceNew field %q does not exist in the schema", field)
		}
		listed[field] = true
	}
	for name, item := range schemaMap {
		if item.ForceNew != listed[name] {
			t.Fatalf("ForceNew mismatch for %q: schema=%v list=%v", name, item.ForceNew, listed[name])
		}
	}
}
