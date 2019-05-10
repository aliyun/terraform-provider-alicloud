package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenInstanceGrant_basic(t *testing.T) {
	var rule vpc.CbnGrantRule

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceGrantBasic(os.Getenv("ALICLOUD_ACCESS_KEY_2"), os.Getenv("ALICLOUD_SECRET_KEY_2"), os.Getenv("ALICLOUD_ACCOUNT_ID_1"), os.Getenv("ALICLOUD_ACCOUNT_ID_2")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceGrantExistsWithProviders("alicloud_cen_instance_grant.foo", &rule, &providers),
					resource.TestCheckResourceAttr("alicloud_cen_instance_grant.foo", "cen_owner_id", os.Getenv("ALICLOUD_ACCOUNT_ID_2")),
				),
			},
		},
	})
}

func testAccCheckCenInstanceGrantExistsWithProviders(n string, rule *vpc.CbnGrantRule, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			vpcService := VpcService{client}

			resp, err := vpcService.DescribeGrantRulesToCen(rs.Primary.ID)
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

func testAccCheckCenInstanceGrantDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInstanceGrantDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInstanceGrantDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_instance_grant" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 3)
		if err != nil {
			return WrapError(err)
		}
		instanceId := parts[1]

		if err != nil {
			return err
		}

		rule, err := vpcService.DescribeGrantRulesToCen(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Child instance %s still grant to CEN %s", instanceId, rule.CenInstanceId)
		}
	}

	return nil
}

func testAccCenInstanceGrantBasic(access, secret, uid1, uid2 string) string {
	return fmt.Sprintf(`
	provider "alicloud" {
		alias = "account1"
	}

	provider "alicloud" {
		access_key = "%s"
		secret_key = "%s"
		alias = "account2"
	}

	variable "name" {
		default = "tf-testAccCenInstanceGrantBasic"
	}

	resource "alicloud_cen_instance" "cen" {
		provider = "alicloud.account2"
		name = "${var.name}"
	}

	resource "alicloud_vpc" "vpc" {
		provider = "alicloud.account1"
		name = "${var.name}"
		cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_cen_instance_grant" "foo" {
		provider = "alicloud.account1"
		cen_id = "${alicloud_cen_instance.cen.id}"
		child_instance_id = "${alicloud_vpc.vpc.id}"
		cen_owner_id = "%s"
	}

    resource "alicloud_cen_instance_attachment" "foo" {
        provider = "alicloud.account2"
        instance_id = "${alicloud_cen_instance.cen.id}"
        child_instance_id = "${alicloud_vpc.vpc.id}"
        child_instance_region_id = "cn-qingdao"
        child_instance_owner_id = "%s"
        depends_on = [
            "alicloud_cen_instance_grant.foo"]
    }
	`, access, secret, uid2, uid1)
}
