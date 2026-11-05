package ram_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_account_alias", &resource.Sweeper{
		Name: "alicloud_ram_account_alias",
		F:    testSweepAccountAliases,
	})
}

func testSweepAccountAliases(region string) error {
	rawClient, err := ramTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := ramsdk.CreateGetAccountAliasRequest()
	raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		return errs.WrapError(err)
	}
	sweeped := false
	resp, _ := raw.(*ramsdk.GetAccountAliasResponse)
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

	_, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		request := ramsdk.CreateClearAccountAliasRequest()
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

func TestAccAliCloudRAMAccountAlias_basic(t *testing.T) {
	randInt := tfacctest.RandIntRange(1000, 9999)
	var v *ramsdk.GetAccountAliasResponse
	resourceId := "alicloud_ram_account_alias.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.RamNoSkipRegions)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRamAccountAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccountAliasConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": fmt.Sprintf("tf-testacc%s%d", acctest.DefaultRegionToTest, randInt),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRamAccountAliasConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_account_alias" "default" {
	  account_alias = "tf-testacc%s%d"
	}`, acctest.DefaultRegionToTest, rand)
}

func testAccCheckRamAccountAliasDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_account_alias" {
			continue
		}

		// Try to find the alias
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := ramsdk.CreateGetAccountAliasRequest()
		_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.GetAccountAlias(request)
		})

		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return errs.WrapError(err)
		}
	}
	return nil
}

// Test Ram AccountAlias. >>> Resource test cases, automatically generated.
// Case  AccountAlias测试 9000
// lintignore: AT001
func TestAccAliCloudRamAccountAlias_basic9000(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_account_alias.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRamAccountAliasMap9000)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamAccountAlias")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	aliasCreate := fmt.Sprintf("tf-testacccreate%d", rand)
	aliasUpdate := fmt.Sprintf("tf-testaccupdate%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRamAccountAliasBasicDependence9000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_alias": aliasUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_alias": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRamAccountAliasMap9000 = map[string]string{}

func AlicloudRamAccountAliasBasicDependence9000(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ram AccountAlias. <<< Resource test cases, automatically generated.
