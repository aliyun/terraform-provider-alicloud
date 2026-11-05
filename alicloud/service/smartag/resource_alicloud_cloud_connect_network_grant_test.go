package smartag_test

import (
	"fmt"
	"os"
	"testing"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudCloudConnectNetworkGrant_basic(t *testing.T) {
	var grantRule smartagsdk.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default"
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := alicloud.Provider()
			providers = append(providers, p)
			return p, nil
		},
	}

	ra := acctest.ResourceAttrInit(resourceId, ccnGrantMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkGrant-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCcnGrantBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCcnGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":  "${alicloud_cloud_connect_network.ccn.id}",
					"cen_id":  "${alicloud_cen_instance.cen.id}",
					"cen_uid": testAccPrimaryAccountUid(t),
					"depends_on": []string{
						"alicloud_cloud_connect_network.ccn",
						"alicloud_cen_instance.cen"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnGrantExistsWithProviders(resourceId, &grantRule, &providers),
					testAccCheck(map[string]string{
						"ccn_id":  acctest.CHECKSET,
						"cen_id":  acctest.CHECKSET,
						"cen_uid": testAccPrimaryAccountUid(t),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCloudConnectNetworkGrant_multi(t *testing.T) {
	var grantRule smartagsdk.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default.2"
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := alicloud.Provider()
			providers = append(providers, p)
			return p, nil
		},
	}

	ra := acctest.ResourceAttrInit(resourceId, ccnGrantMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkGrant-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCcnGrantBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCcnGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":   "3",
					"ccn_id":  "${alicloud_cloud_connect_network.ccn.id}",
					"cen_id":  "${alicloud_cen_instance.cen.id}",
					"cen_uid": testAccPrimaryAccountUid(t),
					"depends_on": []string{
						"alicloud_cloud_connect_network.ccn",
						"alicloud_cen_instance.cen"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnGrantExistsWithProviders(resourceId, &grantRule, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

var ccnGrantMap = map[string]string{
	"ccn_id":  acctest.CHECKSET,
	"cen_id":  acctest.CHECKSET,
	"cen_uid": acctest.CHECKSET,
}

func resourceCcnGrantBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cen_instance" "cen" {
  		name = "${var.name}"
	}

	resource "alicloud_cloud_connect_network" "ccn" {
  		name       = "${var.name}"
  		is_default = "true"
	}
`, name)
}

// testAccPrimaryAccountUid resolves the default provider's (account_1) account
// UID at runtime via GetCallerIdentity. The provider-run v2 ACC environment does
// not expose ALICLOUD_ACCOUNT_ID, so the same-account grant test cannot read it
// from the environment; this reuses the shared sweeper client (no new
// credentials, no product code change) to obtain the caller's account id.
func testAccPrimaryAccountUid(t *testing.T) string {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	raw, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Fatalf("failed to build shared client for account id: %s", err)
	}
	client, ok := raw.(*connectivity.AliyunClient)
	if !ok {
		t.Fatalf("unexpected shared client type %T", raw)
	}
	uid, err := client.AccountId()
	if err != nil {
		t.Fatalf("failed to resolve primary account id: %s", err)
	}
	return uid
}

func testAccCheckCcnGrantDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCcnGrantDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCcnGrantDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	sagService := smartag.NewService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cloud_connect_network_grant" {
			continue
		}

		parts, err := helper.ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return errs.WrapError(err)
		}
		ccnId := parts[0]

		if err != nil {
			return err
		}

		rule, err := sagService.DescribeCloudConnectNetworkGrant(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Child instance %s still grant to CEN %s", ccnId, rule.CenInstanceId)
		}
	}

	return nil
}

func testAccCheckCcnGrantExistsWithProviders(n string, rule *smartagsdk.GrantRule, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cen Child Instance Grant ID is set")
		}

		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			sagService := smartag.NewService(client)

			resp, err := sagService.DescribeCloudConnectNetworkGrant(rs.Primary.ID)
			if err != nil {
				if errs.NotFoundError(err) {
					// only one provider can get the rule
					continue
				}
				return err
			}

			*rule = resp
			return nil
		}

		return fmt.Errorf("Cen Child Instance Grant not found")
	}
}
