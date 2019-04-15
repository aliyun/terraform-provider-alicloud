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
	resource.AddTestSweepers("alicloud_ram_account_alias", &resource.Sweeper{
		Name: "alicloud_ram_account_alias",
		F:    testSweepAccountAliases,
	})
}

func testSweepAccountAliases(region string) error {
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

	request := ram.CreateGetAccountAliasRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		return WrapError(err)
	}
	sweeped := false
	resp, _ := raw.(*ram.GetAccountAliasResponse)
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

	_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		request := ram.CreateClearAccountAliasRequest()
		return ramClient.ClearAccountAlias(request)
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
			{
				Config: testAccRamAccountAliasConfig(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(testAccCheckRamAccountAliasExists("alicloud_ram_account_alias.alias", &v),
					resource.TestMatchResourceAttr("alicloud_ram_account_alias.alias", "account_alias", regexp.MustCompile("^tf-testaccramaccountalias*")),
				),
			},
		},
	})
}

func testAccCheckRamAccountAliasExists(n string, alias *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Alias ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetAccountAliasRequest()
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetAccountAlias(request)
		})

		if err == nil {
			response, _ := raw.(*ram.GetAccountAliasResponse)
			*alias = response.AccountAlias
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamAccountAliasDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_account_alias" {
			continue
		}

		// Try to find the alias
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetAccountAliasRequest()
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetAccountAlias(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
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
