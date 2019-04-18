package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var DBReadWriteMap = map[string]string{
	"port":              "3306",
	"distribution_type": "Standard",
	"weight":            NOSET,
	"max_delay_time":    "30",
	"instance_id":       CHECKSET,
}

func TestAccAlicloudDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection = &rds.DBInstanceNetInfo{}
	var primary = &rds.DBInstanceAttribute{}
	var readonly = &rds.DBInstanceAttribute{}

	resourceId := "alicloud_db_read_write_splitting_connection.default"
	ra := resourceAttrInit(resourceId, DBReadWriteMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rc_connection := resourceCheckInit(resourceId, &connection, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rc_primary := resourceCheckInit("alicloud_db_instance.default", &primary, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rc_readonly := resourceCheckInit("alicloud_db_readonly_instance.default", &readonly, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	randomPrefix := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_read_write_splitting_connection.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadWriteSplittingConnection_basic(testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					rc_connection.checkResourceExists(),
					testAccCheck(map[string]string{
						"connection_string": CHECKSET,
						"connection_prefix": fmt.Sprintf("t-con-%d", randomPrefix),
					}),
				),
			},
			{
				Config: testAccDBReadWriteSplittingConnection_update(testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					rc_connection.checkResourceExists(),
					rc_primary.checkResourceExists(),
					rc_readonly.checkResourceExists(),
					testAccCheck(map[string]string{
						"distribution_type": "Custom",
						"max_delay_time":    "300",
						"weight.%":          "2",
					}),
				),
			},
			{
				Config: testAccDBReadWriteSplittingConnection_update(testAccDBReadonlyInstance_multiAZ(testAccDBInstance_multiAZ), randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					rc_connection.checkResourceExists(),
					rc_primary.checkResourceExists(),
					rc_readonly.checkResourceExists(),
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckDBReadWriteSplittingConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_read_write_splitting_connection" {
			continue
		}

		conn, err := rdsService.DescribeDBReadWriteSplittingConnection(rs.Primary.ID)
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

	resource "alicloud_db_read_write_splitting_connection" "default" {
	    instance_id = "${alicloud_db_instance.default.id}"
		connection_prefix = "t-con-%d"
		distribution_type = "Standard"
		
		depends_on = ["alicloud_db_readonly_instance.default"]
	}
	`, common, rand)
}

func testAccDBReadWriteSplittingConnection_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s

	resource "alicloud_db_read_write_splitting_connection" "default" {
		instance_id = "${alicloud_db_instance.default.id}"
		connection_prefix = "t-con-%d"
		distribution_type = "Custom"
		max_delay_time = 300
		weight = "${map(
			"${alicloud_db_instance.default.id}", "0",
			"${alicloud_db_readonly_instance.default.id}", "500"
		)}"
		
		depends_on = ["alicloud_db_readonly_instance.default"]
	}
	`, common, rand)
}
