package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"regexp"
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
				Config: testAccCRNamespace_Basic(acctest.RandIntRange(100000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-test-acc-cr-ns-basic*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "false"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PUBLIC"),
				),
			},
		},
	})
}

func TestAccAlicloudCRNamespace_Update(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_namespace.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRNamespace_UpdateBefore(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-test-acc-cr-ns*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "false"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PUBLIC"),
				),
			},
			{
				Config: testAccCRNamespace_UpdateAfter(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-test-acc-cr-ns*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "true"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PRIVATE"),
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
		resp, err := crService.DescribeNamespace(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
				continue
			}
			return err
		}

		if resp.Data.Namespace.Namespace != "" {
			return fmt.Errorf("error namespace %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCRNamespace_Basic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-test-acc-cr-ns-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}
`, rand)
}

func testAccCRNamespace_UpdateBefore(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-test-acc-cr-ns-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}
`, rand)
}

func testAccCRNamespace_UpdateAfter(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-test-acc-cr-ns-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}
`, rand)
}
