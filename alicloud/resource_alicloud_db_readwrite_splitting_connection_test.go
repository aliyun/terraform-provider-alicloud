package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBReadWriteSplittingConnection_basic(t *testing.T) {
	var connection rds.DBInstanceNetInfo
	var primary rds.DBInstanceAttribute
	var readonly rds.DBInstanceAttribute

	randomPrefix := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_read_write_splitting_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_basic(testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttrPtr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_string", &connection.ConnectionString),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "port", "3306"),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_prefix", regexp.MustCompile("^t-con-*")),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "distribution_type", "Standard"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "weight"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "max_delay_time", "30"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_read_write_splitting_connection.foo", "instance_id"),
				),
			},
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_update(testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"max_delay_time", "300"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"weight.%", "2"),
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &primary),
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &readonly),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", primary.DBInstanceId), "0")(s)
					},
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", readonly.DBInstanceId), "500")(s)
					},
				),
			},
		},
	})

}

func TestAccAlicloudDBReadWriteSplittingConnection_multiAZ_classic(t *testing.T) {
	var connection rds.DBInstanceNetInfo
	var primary rds.DBInstanceAttribute
	var readonly rds.DBInstanceAttribute

	randomPrefix := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_read_write_splitting_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_basic(testAccDBReadonlyInstance_multiAZ(testAccDBInstance_multiAZ), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttrPtr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_string", &connection.ConnectionString),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "port", "3306"),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_prefix", regexp.MustCompile("^t-con-*")),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "distribution_type", "Standard"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "weight"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "max_delay_time", "30"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_read_write_splitting_connection.foo", "instance_id"),
				),
			},
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_update(testAccDBReadonlyInstance_multiAZ(testAccDBInstance_multiAZ), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"max_delay_time", "300"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"weight.%", "2"),
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &primary),
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &readonly),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", primary.DBInstanceId), "0")(s)
					},
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", readonly.DBInstanceId), "500")(s)
					},
				),
			},
		},
	})

}

func TestAccAlicloudDBReadWriteSplittingConnection_multiAZ_vpc(t *testing.T) {
	var connection rds.DBInstanceNetInfo
	var primary rds.DBInstanceAttribute
	var readonly rds.DBInstanceAttribute

	randomPrefix := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_read_write_splitting_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_basic(testAccDBReadonlyInstance_multiAZ_vpc(testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttrPtr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_string", &connection.ConnectionString),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "port", "3306"),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "connection_prefix", regexp.MustCompile("^t-con-*")),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "distribution_type", "Standard"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "weight"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo", "max_delay_time", "30"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_read_write_splitting_connection.foo", "instance_id"),
				),
			},
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_update(testAccDBReadonlyInstance_multiAZ_vpc(testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"max_delay_time", "300"),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"weight.%", "2"),
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &primary),
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &readonly),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", primary.DBInstanceId), "0")(s)
					},
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							"alicloud_db_read_write_splitting_connection.foo", fmt.Sprintf("weight.%s", readonly.DBInstanceId), "500")(s)
					},
				),
			},
		},
	})

}

func testAccCheckDBReadWriteSplittingConnectionExists(n string, d *rds.DBInstanceNetInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB connection ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rdsService := RdsService{client}

		conn, err := rdsService.DescribeReadWriteSplittingConnection(rs.Primary.ID)
		if err != nil {
			return err
		}

		*d = *conn
		return nil
	}
}

func testAccCheckDBReadWriteSplittingConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_read_write_splitting_connection" {
			continue
		}

		conn, err := rdsService.DescribeReadWriteSplittingConnection(rs.Primary.ID)
		if conn != nil {
			return fmt.Errorf("Error db connection string %s is still existing.", conn.ConnectionString)
		}

		if NotFoundError(err) {
			continue
		}

		return err
	}

	return nil
}

func testAccDBReadWriteSplittingConnection_basic(common string, rand int) string {
	return fmt.Sprintf(`
	%s

	resource "alicloud_db_read_write_splitting_connection" "foo" {
	    instance_id = "${alicloud_db_instance.foo.id}"
		connection_prefix = "t-con-%d"
		distribution_type = "Standard"
		
		depends_on = ["alicloud_db_readonly_instance.foo"]
	}
	`, common, rand)
}

func testAccDBReadWriteSplittingConnection_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s

	resource "alicloud_db_read_write_splitting_connection" "foo" {
		instance_id = "${alicloud_db_instance.foo.id}"
		connection_prefix = "t-con-%d"
		distribution_type = "Custom"
		max_delay_time = 300
		weight = "${map(
			"${alicloud_db_instance.foo.id}", "0",
			"${alicloud_db_readonly_instance.foo.id}", "500"
		)}"
		
		depends_on = ["alicloud_db_readonly_instance.foo"]
	}
	`, common, rand)
}
