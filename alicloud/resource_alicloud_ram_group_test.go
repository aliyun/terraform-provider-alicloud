package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var groups []ram.Group
	request := ram.CreateListGroupsRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroups(request)
		})
		if err != nil {
			return WrapError(err)
		}
		resp, _ := raw.(*ram.ListGroupsResponse)
		if len(resp.Groups.Group) < 1 {
			break
		}
		groups = append(groups, resp.Groups.Group...)

		if !resp.IsTruncated {
			break
		}
		request.Marker = resp.Marker
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
		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = name

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Group (%s): %s", name, err)
		}
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		for _, p := range response.Policies.Policy {
			request := ram.CreateDetachPolicyFromGroupRequest()
			request.PolicyType = p.PolicyType
			request.GroupName = name
			request.PolicyName = p.PolicyName
			log.Printf("[INFO] Detaching Ram policy %s from group: %s", p.PolicyName, name)
			_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.DetachPolicyFromGroup(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to detach policy from Group (%s): %s", name, err)
			}
		}
		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateDeleteGroupRequest()
			request.GroupName = name
			return ramClient.DeleteGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Group (%s): %s", name, err)
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
			{
				Config: testAccRamGroupConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_group.group",
						"name",
						regexp.MustCompile("^tf-testAccRamGroupConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ram_group.group",
						"comments",
						"group comments"),
				),
			},
		},
	})

}

func TestAccAlicloudRamGroup_rename(t *testing.T) {
	var v ram.Group
    
	randInt := acctest.RandIntRange(1000000, 99999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists("alicloud_ram_group.group", &v),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","name",fmt.Sprintf("tf-testAccRamGroupConfig-%d", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","comments","group comments"),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","force","true"),
				),
			},
			{
				Config: testAccRamGroupConfig_rename(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists("alicloud_ram_group.group", &v),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","name",fmt.Sprintf("tf-testAccRamGroupConfig-new-%d", randInt)),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","comments","group comments"),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","force","true"),
				),
			},
		},
	})

}

func TestAccAlicloudRamGroup_recomments(t *testing.T) {
	var v ram.Group
    
	randInt := acctest.RandIntRange(1000000, 99999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists("alicloud_ram_group.group", &v),
					resource.TestMatchResourceAttr("alicloud_ram_group.group","name",regexp.MustCompile("^tf-testAccRamGroupConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","comments","group comments"),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","force","true"),
				),
			},
			{
				Config: testAccRamGroupConfig_recomments(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists("alicloud_ram_group.group", &v),
					resource.TestMatchResourceAttr("alicloud_ram_group.group","name",regexp.MustCompile("^tf-testAccRamGroupConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","comments","group comments new"),
					resource.TestCheckResourceAttr("alicloud_ram_group.group","force","true"),
				),
			},
		},
	})

}

func testAccCheckRamGroupExists(n string, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Group ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetGroupRequest()
		request.GroupName = rs.Primary.ID
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetGroup(request)
		})

		if err == nil {
			response, _ := raw.(*ram.GetGroupResponse)
			*group = response.Group
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group" {
			continue
		}

		// Try to find the group
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetGroupRequest()
		request.GroupName = rs.Primary.ID

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamGroupConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "group" {
	  name = "tf-testAccRamGroupConfig-%d"
	  comments = "group comments"
	  force=true
	}`, rand)
}

func testAccRamGroupConfig_rename(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "group" {
	  name = "tf-testAccRamGroupConfig-new-%d"
	  comments = "group comments"
	  force=true
	}`, rand)
}

func testAccRamGroupConfig_recomments(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_group" "group" {
	  name = "tf-testAccRamGroupConfig-%d"
	  comments = "group comments new"
	  force=true
	}`, rand)
}
