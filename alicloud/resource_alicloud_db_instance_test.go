package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_db_instance", &resource.Sweeper{
		Name: "alicloud_db_instance",
		F:    testSweepDBInstances,
	})
}

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

	var insts []rds.DBInstance
	req := rds.CreateDescribeDBInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDBInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving RDS Instances: %s", err)
		}
		resp, _ := raw.(*rds.DescribeDBInstancesResponse)
		if resp == nil || len(resp.Items.DBInstance) < 1 {
			break
		}
		insts = append(insts, resp.Items.DBInstance...)

		if len(resp.Items.DBInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.DBInstanceDescription
		id := v.DBInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}

		}

		if skip {
			log.Printf("[INFO] Skipping RDS Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting RDS Instance: %s (%s)", name, id)
		if len(v.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId) > 0 {
			request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
			request.DBInstanceId = id
			if _, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ReleaseReadWriteSplittingConnection(request)
			}); err != nil {
				log.Printf("[ERROR] ReleaseReadWriteSplittingConnection error: %#v", err)
			} else {
				time.Sleep(5 * time.Second)
			}
		}
		req := rds.CreateDeleteDBInstanceRequest()
		req.DBInstanceId = id
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete RDS Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudDBInstanceMysql(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
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
					"engine":                   "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":           "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"monitoring_period":        "60",
					"db_instance_storage_type": "local_ssd",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "5.6",
						"instance_type":              CHECKSET,
						"instance_storage":           CHECKSET,
						"auto_upgrade_minor_version": "Auto",
						"db_instance_storage_type":   "local_ssd",
						"resource_group_id":          CHECKSET,
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
						"instance_storage": "10",
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
					"sql_collector_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sql_collector_config_value": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sql_collector_config_value": "30",
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
					"engine":                     "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":             "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":              "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":           "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min * 3}",
					"db_instance_storage_type":   "local_ssd",
					"sql_collector_status":       "Enabled",
					"sql_collector_config_value": "30",
					"instance_name":              "tf-testAccDBInstanceConfig",
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
					"encryption_key": "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                     "MySQL",
						"engine_version":             "5.6",
						"instance_type":              CHECKSET,
						"instance_storage":           "15",
						"db_instance_storage_type":   "local_ssd",
						"sql_collector_status":       "Enabled",
						"sql_collector_config_value": "30",
						"instance_name":              "tf-testAccDBInstanceConfig",
						"monitoring_period":          "60",
						"zone_id":                    CHECKSET,
						"instance_charge_type":       "Postpaid",
						"connection_string":          CHECKSET,
						"port":                       CHECKSET,
						"security_group_id":          "",
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
		},
	})
}

func resourceDBInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  key_state               = "Enabled"
}

`, RdsCommonTestCase, name)
}

func TestAccAlicloudDBInstanceMultiInstance(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alicloud_db_instance.default.4"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
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
					"count":                "5",
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type": "Postpaid",
					"instance_name":        "${var.name}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"monitoring_period":    "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudDBInstanceHighAvailabilityInstance(t *testing.T) {
	var instance *rds.DBInstanceAttribute
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
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type": "Postpaid",
					"instance_name":        "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.availability_zone}",
					"zone_id_slave_a":      "${alicloud_vswitch.slave_a.availability_zone}",
					"vswitch_id":           "${local.vswitch_id}",
					"monitoring_period":    "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "MySQL",
						"engine_version": "8.0",
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

func resourceDBInstanceHighAvailabilityConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

locals {
  vswitch_id = format("%%s,%%s",alicloud_vswitch.default.id,alicloud_vswitch.slave_a.id)
}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "8.0"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "8.0"
  storage_type       = "local_ssd"
}

resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_vswitch" "slave_a" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.1.sub_zone_ids.0}"
  name              = "tf-testaccvswitchslave"
}
`, RdsCommonTestCase, name)
}

func TestAccAlicloudDBInstanceEnterpriseEditionInstance(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig_slave_zone"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceEnterpriseEditionConfigDependence)

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
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "mysql.n2.small.25",
					"instance_storage":     "5",
					"instance_charge_type": "Postpaid",
					"instance_name":        "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.availability_zone}",
					"zone_id_slave_a":      "${alicloud_vswitch.slave_a.availability_zone}",
					"zone_id_slave_b":      "${alicloud_vswitch.slave_b.availability_zone}",
					"vswitch_id":           "${local.vswitch_id}",
					"monitoring_period":    "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "MySQL",
						"engine_version": "8.0",
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

func resourceDBInstanceEnterpriseEditionConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}
 
locals {
  vswitch_id = format("%%s,%%s,%%s",alicloud_vswitch.default.id,alicloud_vswitch.slave_a.id,alicloud_vswitch.slave_b.id)
}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "8.0"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "8.0"
  storage_type       = "local_ssd"
}

resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_vswitch" "slave_a" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.1.sub_zone_ids.0}"
  name              = "tf-testaccvswitcha"
}
resource "alicloud_vswitch" "slave_b" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.2.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.2.sub_zone_ids.0}"
  name              = "tf-testaccvswitchb"
}

`, RdsCommonTestCase, name)
}

// Unknown current resource exists
func TestAccAlicloudDBInstanceSQLServer(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					"engine":                   "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":           "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "cloud_essd",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2012",
						"instance_type":            CHECKSET,
						"instance_storage":         "20",
						"db_instance_storage_type": "cloud_essd",
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
					"engine":                   "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":           "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"db_instance_storage_type": "cloud_essd",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"security_group_ids":       []string{},
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "SQLServer",
						"engine_version":           "2012",
						"instance_type":            CHECKSET,
						"instance_storage":         "25",
						"db_instance_storage_type": "cloud_essd",
						"instance_name":            "tf-testAccDBInstanceConfig",
						"monitoring_period":        "60",
						"zone_id":                  CHECKSET,
						"instance_charge_type":     "Postpaid",
						"connection_string":        CHECKSET,
						"port":                     CHECKSET,
						"security_group_id":        "",
						"security_group_ids.#":     "0",
					}),
				),
			},
		},
	})
}

func resourceDBInstanceSQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "SQLServer"
  engine_version       = "2012"
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "SQLServer"
  engine_version       = "2012"
}

resource "alicloud_security_group" "default" {
	count = 2
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccAlicloudDBInstancePostgreSQL(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					"engine":                   "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":           "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "local_ssd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "9.4",
						"instance_storage":         "20",
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "local_ssd",
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
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_charge_type": "Postpaid",
					"instance_name":        "${var.name}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"security_group_ids":   []string{},
					"monitoring_period":    "60",
					"encryption_key":       "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "9.4",
						"instance_type":        CHECKSET,
						"instance_storage":     "25",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    "",
						"security_group_ids.#": "0",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDBInstancePostgreSQLSSL(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					"engine_version":           "13.0",
					"instance_type":            "pg.n2.small.1",
					"instance_storage":         "30",
					"db_instance_storage_type": "cloud_ssd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "vsw-bp1choevojf5gasmazq3c",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PostgreSQL",
						"engine_version":           "13.0",
						"instance_storage":         "30",
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "cloud_ssd",
						"private_ip_address":       CHECKSET,
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
					"private_ip_address": "172.19.96.98",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip_address": "172.19.96.98",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"port":                     "3333",
					"connection_string_prefix": "rm-ccccccc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":                     "3333",
						"connection_string_prefix": "rm-ccccccc",
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
					"engine":                      "PostgreSQL",
					"engine_version":              "13.0",
					"instance_type":               "pg.n2.small.1",
					"instance_storage":            "30",
					"instance_charge_type":        "Postpaid",
					"instance_name":               "${var.name}",
					"vswitch_id":                  "vsw-bp1choevojf5gasmazq3c",
					"security_group_ids":          []string{},
					"monitoring_period":           "60",
					"encryption_key":              "${alicloud_kms_key.default.id}",
					"ssl_action":                  "Open",
					"ca_type":                     "aliyun",
					"client_ca_enabled":           "1",
					"client_ca_cert":              client_ca_cert,
					"client_crl_enabled":          "1",
					"client_cert_revocation_list": client_cert_revocation_list,
					"acl":                         "cert",
					"replication_acl":             "cert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                      "PostgreSQL",
						"engine_version":              "13.0",
						"instance_type":               CHECKSET,
						"instance_storage":            "30",
						"instance_name":               "tf-testAccDBInstanceConfig",
						"monitoring_period":           "60",
						"zone_id":                     CHECKSET,
						"instance_charge_type":        "Postpaid",
						"connection_string":           CHECKSET,
						"port":                        CHECKSET,
						"security_group_id":           "",
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
					}),
				),
			},
		},
	})
}

const client_ca_cert = `-----BEGIN CERTIFICATE-----\nMIIC+TCCAeGgAwIBAgIJAKfv52qIKAi7MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV\nBAMMCHJvb3QtY2ExMB4XDTIxMDQyMzA3Mjk1M1oXDTMxMDQyMTA3Mjk1M1owEzER\nMA8GA1UEAwwIcm9vdC1jYTEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCyCXrZgqdge6oSji+URDXN0pMWnq4D8doP8quz09shN9TU4iqtyX+Bw+uYOoNF\ndNL4W09p8ykca3RzZghXdbHvtSZy5oCe1rup0xaATAgejDZKBi32ogLXdlA5UMyi\nc0OqIQpOZ+OmeMEVEZP7wsbDy7jS2v59d5OI4tnH2V2SDoWlI/7F9QOq36ER0UqY\nnnjJGnOsTDVeSy4ZXHMT0pXvSSLHsMMhzSJa6t3CiOuAeAW43zIS9tag0yvJI1v7\nxKSJTLs9O5V/h+oD9xofQ4kb4kOdStB2KpDteNfJWJoJYdvRMO+g1u6c2ovlc7KR\nrJPX2ZMJh14q99gPt6Dd+beVAgMBAAGjUDBOMB0GA1UdDgQWBBTDGEb5Aj6SI7hM\nC+AJa3YTNLdDrTAfBgNVHSMEGDAWgBTDGEb5Aj6SI7hMC+AJa3YTNLdDrTAMBgNV\nHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAXWXp6H4bAMZZN6b/rmuxvn4XP\n8p/7NN7BgPQSvQ24U5n8Lo2X8yXYZ4Si/NfWBitAqHceTk6rYTFhODG8CykiduHh\nowfhSjlMj9MGVw3j6I7crBuQ8clUGpy0mUNWJ9ObIdEMaVT+S1Jwk88Byf5FEBxO\nZLg+hg4NQh9qspFAtnhprU9LbcpVtQFY6uyCPs6OEOpPWF1Vtcu+ibQdIQV/e1SQ\n3NJ54R3MCfgEb9errFPv/rXscgahSMxW0sDvObAYdeIeiVeBp3wYKKFHeRNFPGT1\njzei5hlUJzGHf9DlgAH/KODvWUY5cvpuMtJY2yLyJv9xHjjyMnZZAOtHZxfR\n-----END CERTIFICATE-----`
const client_ca_cert2 = "-----BEGIN CERTIFICATE-----\nMIIC+TCCAeGgAwIBAgIJAKfv52qIKAi7MA0GCSqGSIb3DQEBCwUAMBMxETAPBgNV\nBAMMCHJvb3QtY2ExMB4XDTIxMDQyMzA3Mjk1M1oXDTMxMDQyMTA3Mjk1M1owEzER\nMA8GA1UEAwwIcm9vdC1jYTEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB\nAQCyCXrZgqdge6oSji+URDXN0pMWnq4D8doP8quz09shN9TU4iqtyX+Bw+uYOoNF\ndNL4W09p8ykca3RzZghXdbHvtSZy5oCe1rup0xaATAgejDZKBi32ogLXdlA5UMyi\nc0OqIQpOZ+OmeMEVEZP7wsbDy7jS2v59d5OI4tnH2V2SDoWlI/7F9QOq36ER0UqY\nnnjJGnOsTDVeSy4ZXHMT0pXvSSLHsMMhzSJa6t3CiOuAeAW43zIS9tag0yvJI1v7\nxKSJTLs9O5V/h+oD9xofQ4kb4kOdStB2KpDteNfJWJoJYdvRMO+g1u6c2ovlc7KR\nrJPX2ZMJh14q99gPt6Dd+beVAgMBAAGjUDBOMB0GA1UdDgQWBBTDGEb5Aj6SI7hM\nC+AJa3YTNLdDrTAfBgNVHSMEGDAWgBTDGEb5Aj6SI7hMC+AJa3YTNLdDrTAMBgNV\nHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAXWXp6H4bAMZZN6b/rmuxvn4XP\n8p/7NN7BgPQSvQ24U5n8Lo2X8yXYZ4Si/NfWBitAqHceTk6rYTFhODG8CykiduHh\nowfhSjlMj9MGVw3j6I7crBuQ8clUGpy0mUNWJ9ObIdEMaVT+S1Jwk88Byf5FEBxO\nZLg+hg4NQh9qspFAtnhprU9LbcpVtQFY6uyCPs6OEOpPWF1Vtcu+ibQdIQV/e1SQ\n3NJ54R3MCfgEb9errFPv/rXscgahSMxW0sDvObAYdeIeiVeBp3wYKKFHeRNFPGT1\njzei5hlUJzGHf9DlgAH/KODvWUY5cvpuMtJY2yLyJv9xHjjyMnZZAOtHZxfR\n-----END CERTIFICATE-----"
const client_cert_revocation_list = `-----BEGIN X509 CRL-----\nMIIBpzCBkAIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290LWNhMRcN\nMjEwNDI5MDYwODMyWhcNMjEwNTI5MDYwODMyWjA4MBoCCQCG3wQwiFfYbRcNMjEw\nNDIzMTE0MTI4WjAaAgkAht8EMIhX2G8XDTIxMDQyOTA2MDc1N1qgDzANMAsGA1Ud\nFAQEAgIQATANBgkqhkiG9w0BAQsFAAOCAQEAq/M+t0zWLZzqw0T23rZsOhjd2/7+\nu1aHAW5jtjWU+lY4UxGqRsjUTJZnOiSq1w7CWhGxanyjtY/hmSeO6hGMuCmini8f\nNEq/jRvfeS7yJieFucnW4JFmz1HbqSr2S1uXRuHB1ziTRtGm3Epe0qynKm6O4L4q\nCIIqba1gye6H4BmEHaQIi4fplN7buWoeC5Ae9EdxRr3+59P4qJhHD4JGller8/QS\n3m1g75AHJO1dxvAEWy8DrrbP5SrqrsP8mmoNVIHXzCQPGEMnA1sG84365krwR+GC\noi1eBKozVqfnyLRA1C/ZY+dtt3I6zocA2Lt2+JX47VsbXApGgAPVIpKN6A==\n-----END X509 CRL-----`
const client_cert_revocation_list2 = "-----BEGIN X509 CRL-----\nMIIBpzCBkAIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290LWNhMRcN\nMjEwNDI5MDYwODMyWhcNMjEwNTI5MDYwODMyWjA4MBoCCQCG3wQwiFfYbRcNMjEw\nNDIzMTE0MTI4WjAaAgkAht8EMIhX2G8XDTIxMDQyOTA2MDc1N1qgDzANMAsGA1Ud\nFAQEAgIQATANBgkqhkiG9w0BAQsFAAOCAQEAq/M+t0zWLZzqw0T23rZsOhjd2/7+\nu1aHAW5jtjWU+lY4UxGqRsjUTJZnOiSq1w7CWhGxanyjtY/hmSeO6hGMuCmini8f\nNEq/jRvfeS7yJieFucnW4JFmz1HbqSr2S1uXRuHB1ziTRtGm3Epe0qynKm6O4L4q\nCIIqba1gye6H4BmEHaQIi4fplN7buWoeC5Ae9EdxRr3+59P4qJhHD4JGller8/QS\n3m1g75AHJO1dxvAEWy8DrrbP5SrqrsP8mmoNVIHXzCQPGEMnA1sG84365krwR+GC\noi1eBKozVqfnyLRA1C/ZY+dtt3I6zocA2Lt2+JX47VsbXApGgAPVIpKN6A==\n-----END X509 CRL-----"

func resourceDBInstancePostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "alicloud_db_instance_engines" "default" {
  	instance_charge_type = "PostPaid"
  	engine               = "PostgreSQL"
  	engine_version       = "9.4"
	multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
  	instance_charge_type = "PostPaid"
  	engine               = "PostgreSQL"
  	engine_version       = "9.4"
  	multi_zone           = true
}

resource "alicloud_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = alicloud_vpc.default.id
}

resource "alicloud_kms_key" "default" {
  description = var.name
  pending_window_in_days  = 7
  key_state               = "Enabled"
}

`, RdsCommonTestCase, name)
}

// Unknown current resource exists
func TestAccAlicloudDBInstancePPAS(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceAZConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsPPASNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":           "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":            "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":         "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"db_instance_storage_type": "local_ssd",
					"zone_id":                  "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_charge_type":     "Postpaid",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"monitoring_period":        "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                   "PPAS",
						"engine_version":           "9.3",
						"instance_storage":         "250",
						"instance_type":            CHECKSET,
						"db_instance_storage_type": "local_ssd",
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
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min + data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_charge_type": "Postpaid",
					"instance_name":        "${var.name}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"security_group_ids":   []string{},
					"monitoring_period":    "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PPAS",
						"engine_version":       "9.3",
						"instance_type":        CHECKSET,
						"instance_storage":     "500",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    "",
						"security_group_ids.#": "0",
					}),
				),
			},
		},
	})
}

func resourceDBInstanceAZConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PPAS"
  engine_version       = "9.3"
  multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PPAS"
  engine_version       = "9.3"
  multi_zone           = true
}

resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}
`, name)
}

// Unknown current resource exists
func TestAccAlicloudDBInstanceMultiAZ(t *testing.T) {
	var instance = &rds.DBInstanceAttribute{}
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
			testAccPreCheckWithRegions(t, false, connectivity.RdsMultiAzNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":            "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":    "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":     "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":  "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":           "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_name":     "${var.name}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"monitoring_period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       REGEXMATCH + ".*" + MULTI_IZ_SYMBOL + ".*",
						"instance_name": "tf-testAccDBInstance_multiAZ",
					}),
				),
			},
		},
	})

}

func resourceDBInstanceMysqlAZConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}
data "alicloud_db_instance_engines" "default" {
  	engine               = "MySQL"
  	engine_version       = "5.6"
	multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
  	engine               = "MySQL"
  	engine_version       = "5.6"
	multi_zone           = true
}
resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccAlicloudDBInstanceClassic(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceClassicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":               "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_charge_type": "Postpaid",
					"zone_id":              `${lookup(data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids[length(data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids)-1], "id")}`,
					"instance_name":        "${var.name}",
					"monitoring_period":    "60",
					"tde_status":           "Enabled",
					"ssl_action":           "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
						"ssl_status": "Yes",
					}),
				),
			},
		},
	})

}

func resourceDBInstanceClassicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_db_instance_engines" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_db_instance_classes" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_zones" "default" {
  	available_resource_creation= "Rds"
}`, name)
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
	"instance_storage":     "5",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
}

var instanceBasicMap2 = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "8.0",
	"instance_type":        CHECKSET,
	"instance_storage":     "5",
	"instance_name":        "tf-testAccDBInstanceConfig_slave_zone",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
}
