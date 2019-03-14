package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCRReposDataSource_Empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCRReposDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cr_repos.all_repos"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_repos.all_repos", "repos.#"),
				),
			},
		},
	})
}

const testAccAlicloudCRReposDataSourceEmpty = `
data "alicloud_cr_repos" "all_repos" {
}
`

func TestAccAlicloudCRReposDataSource_New(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCRReposDataSourceNew(acctest.RandIntRange(1000, 9999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cr_repos.my_repos"),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.my_repos", "ids.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_cr_repos.my_repos", "ids.0", regexp.MustCompile("^tf-testacc-cr-repo-basic-*")),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.my_repos", "repos.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_cr_repos.my_repos", "repos.0.namespace", regexp.MustCompile("^tf-testacc-cr-repo-basic-*")),
					resource.TestMatchResourceAttr("data.alicloud_cr_repos.my_repos", "repos.0.name", regexp.MustCompile("^tf-testacc-cr-repo-basic-*")),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.my_repos", "repos.0.summary", "OLD"),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.my_repos", "repos.0.repo_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_repos.my_repos", "repos.0.domain_list.vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_repos.my_repos", "repos.0.domain_list.public"),
					resource.TestCheckResourceAttrSet("data.alicloud_cr_repos.my_repos", "repos.0.domain_list.internal"),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.my_repos", "repos.0.tags.#", "0"),
				),
			},
			{
				Config: testAccAlicloudCRReposDataSourceNew_NonExists(acctest.RandIntRange(1000, 9999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cr_repos.no_repos"),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.no_repos", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cr_repos.no_repos", "repos.#", "0"),
				),
			},
		},
	})
}

func testAccAlicloudCRReposDataSourceNew(rand int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-cr-repo-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
    name = "${var.name}"
    auto_create	= false
    default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "default" {
    namespace = "${alicloud_cr_namespace.default.name}"
    name = "${var.name}"
    summary = "OLD"
    repo_type = "PUBLIC"
    detail  = "OLD"
}

data "alicloud_cr_repos" "my_repos" {
    name_regex = "${alicloud_cr_repo.default.name}"
    enable_details = true
}
`, rand)
}

func testAccAlicloudCRReposDataSourceNew_NonExists(rand int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-cr-repo-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
    name = "${var.name}"
    auto_create	= false
    default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "default" {
    namespace = "${alicloud_cr_namespace.default.name}"
    name = "${var.name}"
    summary = "OLD"
    repo_type = "PUBLIC"
    detail  = "OLD"
}

data "alicloud_cr_repos" "no_repos" {
    name_regex = "${alicloud_cr_repo.default.name}-nonexists"
    enable_details = true
}
`, rand)
}
