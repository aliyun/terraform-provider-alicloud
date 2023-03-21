package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudCloudConnectNetworkGrant_basic(t *testing.T) {
	var grantRule smartag.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	ra := resourceAttrInit(resourceId, ccnGrantMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkGrant-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnGrantBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCcnGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":  "${alicloud_cloud_connect_network.ccn.id}",
					"cen_id":  "${alicloud_cen_instance.cen.id}",
					"cen_uid": os.Getenv("ALICLOUD_ACCOUNT_ID_2"),
					"depends_on": []string{
						"alicloud_cloud_connect_network.ccn",
						"alicloud_cen_instance.cen"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnGrantExistsWithProviders(resourceId, &grantRule, &providers),
					testAccCheck(map[string]string{
						"ccn_id":  CHECKSET,
						"cen_id":  CHECKSET,
						"cen_uid": os.Getenv("ALICLOUD_ACCOUNT_ID_2"),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCloudConnectNetworkGrant_multi(t *testing.T) {
	var grantRule smartag.GrantRule
	resourceId := "alicloud_cloud_connect_network_grant.default.2"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	ra := resourceAttrInit(resourceId, ccnGrantMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkGrant-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnGrantBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCcnGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":   "3",
					"ccn_id":  "${alicloud_cloud_connect_network.ccn.id}",
					"cen_id":  "${alicloud_cen_instance.cen.id}",
					"cen_uid": os.Getenv("ALICLOUD_ACCOUNT_ID_2"),
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
	"ccn_id":  CHECKSET,
	"cen_id":  CHECKSET,
	"cen_uid": CHECKSET,
}

func resourceCcnGrantBasicDependence(name string) string {
	access2 := os.Getenv("ALICLOUD_ACCESS_KEY_2")
	secret2 := os.Getenv("ALICLOUD_SECRET_KEY_2")

	return fmt.Sprintf(`
	provider "alicloud" {
  		alias = "ccn_account"
	}

	provider "alicloud" {
  		region     = "cn-hangzhou"
  		access_key = "%s"
  		secret_key = "%s"
  		alias      = "cen_account"
	}

	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cen_instance" "cen" {
  		provider = "alicloud.cen_account"
  		name     = "${var.name}"
	}

	resource "alicloud_cloud_connect_network" "ccn" {
  		provider   = "alicloud.ccn_account"
  		name       = "${var.name}"
  		is_default = "true"
	}
`, access2, secret2, name)
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
	sagService := SagService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cloud_connect_network_grant" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		ccnId := parts[0]

		if err != nil {
			return err
		}

		rule, err := sagService.DescribeCloudConnectNetworkGrant(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Child instance %s still grant to CEN %s", ccnId, rule.CenInstanceId)
		}
	}

	return nil
}

func testAccCheckCcnGrantExistsWithProviders(n string, rule *smartag.GrantRule, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			sagService := SagService{client}

			resp, err := sagService.DescribeCloudConnectNetworkGrant(rs.Primary.ID)
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

		return fmt.Errorf("Cen Child Instance Grant not found")
	}
}
