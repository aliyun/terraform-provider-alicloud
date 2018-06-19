package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudLogProject_basic(t *testing.T) {
	var project sls.LogProject

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogProjectBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					resource.TestCheckResourceAttr("alicloud_log_project.foo", "description", "tf unit test"),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogProjectExists(name string, project *sls.LogProject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Log project ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		p, err := client.DescribeLogProject(rs.Primary.ID)
		if err != nil {
			return err
		}
		if p == nil || p.Name == "" {
			return fmt.Errorf("Log project %s is not exist.", rs.Primary.ID)
		}
		project = p

		return nil
	}
}

func testAccCheckAlicloudLogProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AliyunClient).logconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_project" {
			continue
		}

		exist, err := conn.CheckProjectExist(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Check log project got an error: %#v.", err)
		}

		if !exist {
			return nil
		}

		return fmt.Errorf("Log project %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAlicloudLogProjectBasic = `
resource "alicloud_log_project" "foo" {
    name = "for-tf-test"
    description = "tf unit test"
}`
