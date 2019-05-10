package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
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
		if skip {
			log.Printf("[INFO] Skipping RDS Instance: %s (%s)", name, id)
			continue
		}

		sweeped = true
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
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudDBInstance_mysql(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_storage(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_name(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_type(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "rds.mysql.t1.small",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_monitoring_period(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_IPs(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alicloud_db_instance.default", ips)),
				),
			},
			{
				Config: testAccDBInstanceConfig_securitygroup(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"engine_version":       "5.6",
						"instance_type":        "rds.mysql.s2.large",
						"instance_storage":     "30",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDBInstance_multi_instance(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alicloud_db_instance.default.4"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig_multi(RdsCommonTestCase, "MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

// Unknown current resource exists
func TestAccAlicloudDBInstance_SQLServer(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "SQLServer",
						"engine_version": "2012",
						"instance_type":  "rds.mssql.s2.large",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_storage(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_name(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_type(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "rds.mssql.s2.xlarge",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_monitoring_period(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_IPs(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccDBInstanceConfig_securitygroup(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_sqlserver(RdsCommonTestCase, "SQLServer", "2012", "rds.mssql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "SQLServer",
						"engine_version":       "2012",
						"instance_type":        "rds.mssql.s2.large",
						"instance_storage":     "60",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDBInstance_PostgreSQL(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.s1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "PostgreSQL",
						"engine_version": "9.4",
						"instance_type":  "rds.pg.s1.small",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_storage(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.s1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_name(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.s1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_instance_type(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "rds.pg.t1.small",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_monitoring_period(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_IPs(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccDBInstanceConfig_securitygroup(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.t1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig(RdsCommonTestCase, "PostgreSQL", "9.4", "rds.pg.s1.small"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "9.4",
						"instance_type":        "rds.pg.s1.small",
						"instance_storage":     "30",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

// Unknown current resource exists
func SkipTestAccAlicloudDBInstance_PPAS(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsPPASNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig_ppas(DBMultiAZCommonTestCase, "PPAS", "9.3", "ppas.x4.medium.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_instance_storage(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.small.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "50",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_instance_name(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.small.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_instance_type(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.medium.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "rds.mysql.t1.small",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_monitoring_period(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.medium.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_IPs(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.medium.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas_securitygroup(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.medium.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDBInstanceConfig_ppas(DBMultiAZCommonTestCase, "PPAS", "10", "ppas.x4.small.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PPAS",
						"engine_version":       "9.3",
						"instance_type":        "ppas.x4.small.2",
						"instance_storage":     "30",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"monitoring_period":    "60",
						"zone_id":              CHECKSET,
						"instance_charge_type": "Postpaid",
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_id":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

// Unknown current resource exists
func TestAccAlicloudDBInstance_multiAZ(t *testing.T) {
	var instance = &rds.DBInstanceAttribute{}
	resourceId := "alicloud_db_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsMultiAzNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase),
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

func TestAccAlicloudDBInstance_classic(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alicloud_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig_classic("MySQL", "5.6", "rds.mysql.s2.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
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
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func testAccCheckDBInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_instance" {
			continue
		}

		ins, err := rdsService.DescribeDBInstance(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

var instanceBasicMap = map[string]string{
	"engine":               "MySQL",
	"engine_version":       "5.6",
	"instance_type":        "rds.mysql.s2.large",
	"instance_storage":     "30",
	"instance_name":        "tf-testAccDBInstanceConfig",
	"monitoring_period":    "60",
	"zone_id":              CHECKSET,
	"instance_charge_type": "Postpaid",
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
}

func testAccDBInstanceConfig(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_instance_storage(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_instance_name(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_instance_type(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_monitoring_period(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "300"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_IPs(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccDBInstance_instance_name"
	}
	resource "alicloud_db_instance" "default" {
		engine = "%s"
		engine_version = "%s"
		instance_type = "%s"
		instance_storage = "50"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		monitoring_period = "300"
	}
	`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_securitygroup(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
				%s
		variable "creation" {
			default = "Rds"
		}
		variable "name" {
			default = "tf-testAccDBInstance_instance_name"
		}
		resource "alicloud_security_group" "foo-sg1" {
			name   = "${var.name}"
			vpc_id = "${alicloud_vpc.default.id}"
		}

		resource "alicloud_db_instance" "default" {
			engine = "%s"
			engine_version = "%s"
			instance_type = "%s"
			instance_storage = "50"
			instance_charge_type = "Postpaid"
			instance_name = "${var.name}"
			vswitch_id = "${alicloud_vswitch.default.id}"
			security_group_id = "${alicloud_security_group.foo-sg1.id}"
		}
`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_multi(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	count = 5
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstance_vpc_multiAZ(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccDBInstance_multiAZ"
	}
	resource "alicloud_db_instance" "default" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "30"
		zone_id = "${data.alicloud_zones.default.zones.0.id}"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		monitoring_period = "60"
	}
	`, common)
}

func testAccDBInstanceConfigTags(tags, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstanceConfigTags"
	}

	resource "alicloud_db_instance" "default" {
		engine = "%s"
		engine_version = "%s"
		instance_type = "%s"
		instance_storage = "50"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		tags {
			%s
		}
	}`, RdsCommonTestCase, tags, engine, engine_version, instance_type)
}

func testAccDBInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_instance_name"
	}

	resource "alicloud_db_instance" "default" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "50"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		monitoring_period = "300"
	}
	`, common)
}

func testAccDBInstanceConfig_classic(engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}

data "alicloud_zones" "default" {
  	available_resource_creation= "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	zone_id = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)], "id")}"
	instance_name = "${var.name}"
	monitoring_period = "60"
}`, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_sqlserver(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "60"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_instance_storage(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstanceConfig"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_instance_name(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_instance_type(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "60"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_monitoring_period(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDBInstance_instance_name"
}
variable "creation" {
		default = "Rds"
}
resource "alicloud_db_instance" "default" {
	engine = "%s"
	engine_version = "%s"
	instance_type = "%s"
	instance_storage = "50"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
	monitoring_period = "300"
}`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_IPs(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccDBInstance_instance_name"
	}
	resource "alicloud_db_instance" "default" {
		engine = "%s"
		engine_version = "%s"
		instance_type = "%s"
		instance_storage = "50"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		monitoring_period = "300"
	}
	`, common, engine, engine_version, instance_type)
}

func testAccDBInstanceConfig_ppas_securitygroup(common, engine, engine_version, instance_type string) string {
	return fmt.Sprintf(`
				%s
		variable "creation" {
			default = "Rds"
		}
		variable "name" {
			default = "tf-testAccDBInstance_instance_name"
		}
		resource "alicloud_security_group" "foo-sg1" {
			name   = "${var.name}"
			vpc_id = "${alicloud_vpc.default.id}"
		}

		resource "alicloud_db_instance" "default" {
			engine = "%s"
			engine_version = "%s"
			instance_type = "%s"
			instance_storage = "50"
			instance_charge_type = "Postpaid"
			instance_name = "${var.name}"
			vswitch_id = "${alicloud_vswitch.default.id}"
			security_group_id = "${alicloud_security_group.foo-sg1.id}"
		}
`, common, engine, engine_version, instance_type)
}
