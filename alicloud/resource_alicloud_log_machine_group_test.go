package alicloud

import (
	"fmt"
	"testing"

	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogMachineGroup_ip(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogMachineGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogMachineGroupIp(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.foo", &group),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.foo", "identify_type", "ip"),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.foo", "topic", "terraform"),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.foo", "identify_list.#", "3"),
				),
			},
		},
	})
}

func TestAccAlicloudLogMachineGroup_userdefined(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogMachineGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogMachineGroupUserDefined(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.bar", &project),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.bar", &group),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.bar", "identify_type", "userdefined"),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.bar", "topic", "terraform"),
					resource.TestCheckResourceAttr("alicloud_log_machine_group.bar", "identify_list.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogMachineGroupExists(name string, group *sls.MachineGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Log machine group ID is set")
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}
		g, err := logService.DescribeLogMachineGroup(split[0], split[1])
		if err != nil {
			return err
		}

		group = g
		return nil
	}
}

func testAccCheckAlicloudLogMachineGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	logService := LogService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_machine_group" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := logService.DescribeLogMachineGroup(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check log machine group got an error: %#v.", err)
		}
		return fmt.Errorf("Log machine group %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAlicloudLogMachineGroupIp(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "tf-testacclogmachinegroupip-%d"
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}

	resource "alicloud_log_machine_group" "foo" {
	    project = "${alicloud_log_project.foo.name}"
	    name = "${var.name}"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
	}
	`, rand)
}

func testAlicloudLogMachineGroupUserDefined(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "tf-testacclogmachinegroup-%d"
	}
	resource "alicloud_log_project" "bar" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	resource "alicloud_log_machine_group" "bar" {
	    project = "${alicloud_log_project.bar.name}"
	    name = "${var.name}"
	    identify_type = "userdefined"
	    topic = "terraform"
	    identify_list = ["terraform", "abc1234"]
	}
	`, rand)
}
