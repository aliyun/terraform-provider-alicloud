package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCRNamespacesDataSource_Empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCRNamespacesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cr_namespaces.all_namespaces"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_namespaces.all_namespaces", "namespaces.#"),
				),
			},
		},
	})
}

const testAccAlicloudCRNamespacesDataSourceEmpty = `
data "alicloud_cr_namespaces" "all_namespaces" {
    name_regex = "kubernetes-ingress"
    enable_details = false
}
`

func TestAccAlicloudCRNamespacesDataSource_New(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCRNamespacesDataSourceNew(acctest.RandIntRange(1000, 9999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cr_namespaces.my_namespaces"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_namespaces.my_namespaces", "namespaces.#"),
					resource.TestMatchResourceAttr("data.alicloud_cr_namespaces.my_namespaces", "namespaces.0.name", regexp.MustCompile("^tf-testacc-cr-ns-basic-*")),
					resource.TestCheckResourceAttr("data.alicloud_cr_namespaces.my_namespaces", "namespaces.0.default_visibility", "PUBLIC"),
					resource.TestCheckResourceAttr("data.alicloud_cr_namespaces.my_namespaces", "namespaces.0.auto_create", "false"),
				),
			},
		},
	})
}

func testAccAlicloudCRNamespacesDataSourceNew(rand int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-cr-ns-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
    name = "${var.name}"
    auto_create	= false
    default_visibility = "PUBLIC"
}

data "alicloud_cr_namespaces" "m_namespaces" {
    name_regex = "${alicloud_cr_namespace.default.name}"
    enable_details = true
}
`, rand)
}
