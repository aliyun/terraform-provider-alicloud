package alicloud

import (
	"fmt"
	"testing"

	"log"
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
	resource.AddTestSweepers("alicloud_ram_account_alias", &resource.Sweeper{
		Name: "alicloud_ram_account_alias",
		F:    testSweepAccountAliases,
	})
}

func testSweepAccountAliases(region string) error {
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
		return ramClient.GetAccountAlias()
	})
	if err != nil {
		return fmt.Errorf("Error retrieving Ram account alias: %s", err)
	}
	sweeped := false
	resp, _ := raw.(ram.AccountAliasResponse)
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

	_, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ClearAccountAlias()
	})
	if err != nil {
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
				Config: testAccRamAccountAliasConfig(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccountAliasExists(
						"alicloud_ram_account_alias.alias", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_account_alias.alias",
						"account_alias",
						regexp.MustCompile("^tf-testaccramaccountalias*")),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetAccountAlias()
		})

		if err == nil {
			response, _ := raw.(ram.AccountAliasResponse)
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetAccountAlias()
		})

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

func testAccRamAccountAliasConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_account_alias" "alias" {
	  account_alias = "tf-testaccramaccountalias%d"
	}`, rand)
}
