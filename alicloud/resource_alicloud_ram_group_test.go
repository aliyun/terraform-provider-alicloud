package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_group", &resource.Sweeper{
		Name: "alicloud_ram_group",
		F:    testSweepRamGroups,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_ram_user",
		},
	})
}

func testSweepRamGroups(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var groups []ram.Group
	args := ram.GroupListRequest{}
	for {
		resp, err := conn.ramconn.ListGroup(args)
		if err != nil {
			return fmt.Errorf("Error retrieving Ram groups: %s", err)
		}
		if len(resp.Groups.Group) < 1 {
			break
		}
		groups = append(groups, resp.Groups.Group...)

		if !resp.IsTruncated {
			break
		}
		args.Marker = resp.Marker
	}
	sweeped := false

	for _, v := range groups {
		name := v.GroupName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Group: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ram Group: %s", name)
		req := ram.GroupQueryRequest{
			GroupName: name,
		}
		if _, err := conn.ramconn.DeleteGroup(req); err != nil {
			log.Printf("[ERROR] Failed to delete Ram User (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRamGroup_basic(t *testing.T) {
	var v ram.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_group.group",
						"name",
						"tf-testAccRamGroupConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_group.group",
						"comments",
						"group comments"),
				),
			},
		},
	})

}

func testAccCheckRamGroupExists(n string, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Group ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		response, err := conn.GetGroup(ram.GroupQueryRequest{
			GroupName: rs.Primary.ID,
		})

		if err == nil {
			*group = response.Group
			return nil
		}
		return fmt.Errorf("Error finding group %#v", err)
	}
}

func testAccCheckRamGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group" {
			continue
		}

		// Try to find the group
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.ID,
		}

		_, err := conn.GetGroup(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamGroupConfig = `
resource "alicloud_ram_group" "group" {
  name = "tf-testAccRamGroupConfig"
  comments = "group comments"
  force=true
}`
