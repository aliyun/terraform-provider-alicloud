package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCRNamespace_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_namespace.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRNamespace_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "name", "tf-test-acc-cr-namespace-basic"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "false"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PUBLIC"),
				),
			},
		},
	})
}

func testAccCheckCRNamespaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_namespace" {
			continue
		}

		crService := CrService{client}
		resp, err := crService.GetNamespace(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
				continue
			}
			return err
		}

		if resp.Data.Namespace.Namespace != "" {
			return fmt.Errorf("error namespace %s still exists.", rs.Primary.ID)
		}
	}
	return nil
}

const testAccCRNamespace_Basic = `
variable "name" {
	default = "tf-test-acc-cr-namespace-basic"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}
`
