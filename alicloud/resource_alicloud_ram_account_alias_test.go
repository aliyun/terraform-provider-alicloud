package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_account_alias", &resource.Sweeper{
		Name: "alicloud_ram_account_alias",
		F:    testSweepAccountAliases,
	})
}

func testSweepAccountAliases(region string) error {
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

	resp, err := conn.ramconn.GetAccountAlias()
	if err != nil {
		return fmt.Errorf("Error retrieving Ram account alias: %s", err)
	}
	sweeped := false

	name := resp.AccountAlias
	skip := true
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			skip = false
			break
		}
	}
	if skip {
		log.Printf("[INFO] Skipping Ram account alias: %s", name)
		return nil
	}
	sweeped = true
	log.Printf("[INFO] Deleting Ram account alias: %s", name)

	if _, err := conn.ramconn.ClearAccountAlias(); err != nil {
		log.Printf("[ERROR] Failed to delete Ram account alias (%s): %s", name, err)
	}

	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRamAccountAlias_basic(t *testing.T) {
	var v string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_account_alias.alias",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAccountAliasDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamAccountAliasConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccountAliasExists(
						"alicloud_ram_account_alias.alias", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_account_alias.alias",
						"account_alias",
						"testaccramaccountaliasconfig"),
				),
			},
		},
	})

}

func testAccCheckRamAccountAliasExists(n string, alias *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Alias ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		response, err := conn.GetAccountAlias()

		if err == nil {
			*alias = response.AccountAlias
			return nil
		}
		return fmt.Errorf("Error finding alias %s.", rs.Primary.ID)
	}
}

func testAccCheckRamAccountAliasDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_account_alias" {
			continue
		}

		// Try to find the alias
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		_, err := conn.GetAccountAlias()

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamAccountAliasConfig = `
resource "alicloud_ram_account_alias" "alias" {
  account_alias = "testaccramaccountaliasconfig"
}`
