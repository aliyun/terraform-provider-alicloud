package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagGrant_basic(t *testing.T) {
	var grantRule smartag.GrantRule
	resourceId := "alicloud_sag_grant.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	ra := resourceAttrInit(resourceId, SagGrantMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSagGrant-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagGrantDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithSmartAccessGatewaySetting(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSagGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"ccn_id":     "${alicloud_cloud_connect_network.ccn.id}",
					"ccn_uid":    os.Getenv("ALICLOUD_ACCOUNT_ID_2"),
					"depends_on": []string{"alicloud_cloud_connect_network.ccn"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSagGrantExistsWithProviders(resourceId, &grantRule, &providers),
					testAccCheck(map[string]string{
						"sag_id":  os.Getenv("SAG_INSTANCE_ID"),
						"ccn_id":  CHECKSET,
						"ccn_uid": os.Getenv("ALICLOUD_ACCOUNT_ID_2"),
					}),
				),
			},
		},
	})
}

var SagGrantMap = map[string]string{
	"ccn_id": CHECKSET,
}

func resourceSagGrantDependence(name string) string {
	region2 := os.Getenv("ALICLOUD_REGION_2")
	access2 := os.Getenv("ALICLOUD_ACCESS_KEY_2")
	secret2 := os.Getenv("ALICLOUD_SECRET_KEY_2")

	return fmt.Sprintf(`
	provider "alicloud" {
		alias = "sag_account"
	}
	provider "alicloud" {
		region = "%s"
		access_key = "%s"
		secret_key = "%s"
		alias = "ccn_account"
	}
	variable "name" {
		default = "%s"
	}
	resource "alicloud_cloud_connect_network" "ccn" {
		provider = "alicloud.ccn_account"
	  	name = "${var.name}"
	  	is_default = "true"
	}
`, region2, access2, secret2, name)
}

func testAccCheckSagGrantDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckSagGrantDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckSagGrantDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	sagService := SagService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_sag_grant" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		sagId := parts[0]

		if err != nil {
			return err
		}

		rule, err := sagService.DescribeSagGrant(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Child instance %s still grant to CCN %s", sagId, rule.CcnInstanceId)
		}
	}

	return nil
}

func testAccCheckSagGrantExistsWithProviders(n string, rule *smartag.GrantRule, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SAG Child Instance Grant ID is set")
		}

		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			sagService := SagService{client}

			resp, err := sagService.DescribeSagGrant(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					// only one provider can get the rule
					continue
				}
				return err
			}

			*rule = resp
			return nil
		}

		return fmt.Errorf("SAG Child Instance Grant not found")
	}
}
