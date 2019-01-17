package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_role", &resource.Sweeper{
		Name: "alicloud_ram_role",
		F:    testSweepRamRoles,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_fc_service",
		},
	})
}

func testSweepRamRoles(region string) error {
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
		"tftest",
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ListRoles()
	})
	if err != nil {
		return fmt.Errorf("Error retrieving Ram roles: %s", err)
	}
	resp, _ := raw.(ram.ListRoleResponse)
	if len(resp.Roles.Role) < 1 {
		return nil
	}

	sweeped := false

	for _, v := range resp.Roles.Role {
		name := v.RoleName
		id := v.RoleId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Role: %s (%s)", name, id)
			continue
		}
		sweeped = true

		log.Printf("[INFO] Detaching Ram Role: %s (%s) policies.", name, id)
		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListPoliciesForRole(ram.RoleQueryRequest{
				RoleName: name,
			})
		})
		resp, _ := raw.(ram.PolicyListResponse)
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Role (%s (%s)) policies: %s", name, id, err)
		} else if len(resp.Policies.Policy) > 0 {
			for _, v := range resp.Policies.Policy {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.DetachPolicyFromRole(ram.AttachPolicyToRoleRequest{
						PolicyRequest: ram.PolicyRequest{
							PolicyName: v.PolicyName,
							PolicyType: ram.Type(v.PolicyType),
						},
						RoleName: name,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					log.Printf("[ERROR] Failed detach Policy %s: %#v", v.PolicyName, err)
				}
			}
		}

		log.Printf("[INFO] Deleting Ram Role: %s (%s)", name, id)
		req := ram.RoleQueryRequest{
			RoleName: name,
		}
		_, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteRole(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Role (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRamRole_basic(t *testing.T) {
	var v ram.Role

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_role.role",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists(
						"alicloud_ram_role.role", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_role.role",
						"name",
						regexp.MustCompile("^tf-testAccRamRoleConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ram_role.role",
						"description",
						"this is a test"),
				),
			},
		},
	})

}

func testAccCheckRamRoleExists(n string, role *ram.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.RoleQueryRequest{
			RoleName: rs.Primary.Attributes["name"],
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		log.Printf("[WARN] Role id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(ram.RoleResponse)
			*role = response.Role
			return nil
		}
		return fmt.Errorf("Error finding role %#v", rs.Primary.ID)
	}
}

func testAccCheckRamRoleDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role" {
			continue
		}

		// Try to find the role
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.RoleQueryRequest{
			RoleName: rs.Primary.Attributes["name"],
		}

		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetRole(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

func testAccRamRoleConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRoleConfig-%d"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}`, rand)
}
