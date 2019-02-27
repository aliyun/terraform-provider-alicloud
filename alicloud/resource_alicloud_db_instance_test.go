package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/hashcode"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
		"tf_test_",
		"tf-test-",
		"testAcc",
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

func TestAccAlicloudDBInstance_classic(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"2012"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"SQLServer"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_name",
						"tf-testAccDBInstanceConfig"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_vpc(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_vpc(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"20"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"MySQL"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_parameter(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_vpc(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"20"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"MySQL"),
				),
			},
			// update parameter
			resource.TestStep{
				Config: testAccDBInstance_parameter(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
				),
			},
			// update multi parameter
			resource.TestStep{
				Config: testAccDBInstance_parameterMulti(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("connect_timeout")), "50"),
				),
			},
			// remove parameter definition, parameter value not change
			resource.TestStep{
				Config: testAccDBInstance_parameter(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
					testAccCheckDBParameterExpects(
						"alicloud_db_instance.foo", "connect_timeout", "50"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_CreateWithParameter(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_parameter(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"20"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_classic_multiAZ(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_multiAZ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					testAccCheckDBInstanceMultiIZ(&instance),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_Vpc_multiAZ(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsMultiAzNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					testAccCheckDBInstanceMultiIZ(&instance),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_securityIps(t *testing.T) {
	var ips []map[string]interface{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_securityIps(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityIpExists(
						"alicloud_db_instance.foo", ips),
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "127.0.0.1"),
				),
			},

			{
				Config: testAccDBInstance_securityIpsUpdate(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityIpExists(
						"alicloud_db_instance.foo", ips),
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_upgradeClass(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_class(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_type", "rds.pg.t1.small"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_storage", "20"),
				),
			},

			{
				Config: testAccDBInstance_classUpgrade(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_type", "rds.pg.s1.small"),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_storage", "30"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_tags(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfigTags(
					`foo = "bar"
					bar = "foo"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.foo", "bar"),
				),
			},

			{
				Config: testAccDBInstanceConfigTags(
					`bar1 = "zzz"
					bar2 = "bar"
					bar3 = "bar"
					bar4 = "bar"
					bar5 = "zzz"
					bar6 = "bar"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.%", "6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.bar5", "zzz"),
				),
			},

			{
				Config: testAccDBInstanceConfigTags(
					`bar1 = "zzz"
					bar2 = "bar"
					bar3 = "bar"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.bar3", "bar"),
				),
			},

			{
				Config: testAccDBInstanceConfigTags(
					`bar1 = "zzz"
					bar2 = "bar"
					bar3 = "bar_update"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo", "tags.bar3", "bar_update"),
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

func testAccCheckDBInstanceMultiIZ(i *rds.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(i.ZoneId, MULTI_IZ_SYMBOL) {
			return fmt.Errorf("Current region does not support multiIZ.")
		}
		return nil
	}
}

func testAccCheckDBInstanceExists(n string, d *rds.DBInstanceAttribute) resource.TestCheckFunc {
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
		attr, err := rdsService.DescribeDBInstanceById(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = *attr
		return nil
	}
}

// check instance parameter value using SDK API, make sure that
// real parameter value is expected
func testAccCheckDBParameterExpects(n string, key string, value string) resource.TestCheckFunc {
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
		response, err := rdsService.DescribeParameters(rs.Primary.ID)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] check instance %s attribute %#v", rs.Primary.ID, response)
		for _, i := range response.ConfigParameters.DBInstanceParameter {
			if i.ParameterName == key {
				if i.ParameterValue != value {
					return fmt.Errorf(
						"%s: Parameter '%s' expected %#v, got %#v",
						rs.Primary.ID, key, value, i.ParameterValue)
				} else {
					return nil
				}
			}
		}
		for _, i := range response.RunningParameters.DBInstanceParameter {
			if i.ParameterName == key {
				if i.ParameterValue != value {
					return fmt.Errorf(
						"%s: Parameter '%s' expected %#v, got %#v",
						rs.Primary.ID, key, value, i.ParameterValue)
				}
			}
		}

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

		ins, err := rdsService.DescribeDBInstanceById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccDBInstanceConfig = `
data "alicloud_zones" "default" {
  available_resource_creation= "Rds"
}
resource "alicloud_db_instance" "foo" {
	engine = "SQLServer"
	engine_version = "2012"
	instance_type = "rds.mssql.s2.large"
	instance_storage = "30"
	instance_charge_type = "Postpaid"
	instance_name = "tf-testAccDBInstanceConfig"
	zone_id = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%length(data.alicloud_zones.default.zones)], "id")}"
}
`

func testAccDBInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
	`, common)
}

func testAccDBInstance_parameter(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		parameters = [{
			name = "innodb_large_prefix"
			value = "ON"
		}]
	}
	`, common)
}

func testAccDBInstance_parameterMulti(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		parameters = [{
			name = "innodb_large_prefix"
			value = "ON"
		},{
			name = "connect_timeout"
			value = "50"
		}]
	}
	`, common)
}

const testAccDBInstance_multiAZ = `
data "alicloud_zones" "default" {
  available_resource_creation= "Rds"
  multi = true
}
variable "name" {
	default = "tf-testAccDBInstance_multiAZ"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	instance_name = "${var.name}"
}
`

func testAccDBInstance_vpc_multiAZ(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccDBInstance_multiAZ"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		zone_id = "${data.alicloud_zones.default.zones.0.id}"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBInstance_securityIps(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_securityIps"
	}
	resource "alicloud_db_instance" "foo" {
		engine = "SQLServer"
		engine_version = "2012"
		instance_type = "rds.mssql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}
func testAccDBInstance_securityIpsUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_securityIpsUpdate"
	}
	resource "alicloud_db_instance" "foo" {
		engine = "SQLServer"
		engine_version = "2012"
		instance_type = "rds.mssql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBInstance_class(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_class"
	}
	resource "alicloud_db_instance" "foo" {
		engine = "PostgreSQL"
		engine_version = "9.4"
		instance_type = "rds.pg.t1.small"
		instance_storage = "20"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}
func testAccDBInstance_classUpgrade(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBInstance_class"
	}
	resource "alicloud_db_instance" "foo" {
		engine = "PostgreSQL"
		engine_version = "9.4"
		instance_type = "rds.pg.s1.small"
		instance_storage = "30"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBInstanceConfigTags(tags string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
  		available_resource_creation= "Rds"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		instance_name = "tf-testAccDBInstanceConfigTags"
		zone_id = "${data.alicloud_zones.default.zones.0.id}"
		tags {
			%s
		}
	}`, tags)
}
