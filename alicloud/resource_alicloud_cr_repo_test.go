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

func TestAccAlicloudCRRepo_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_repo.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRRepo_Basic(acctest.RandIntRange(1000, 9999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCRRepoExists("alicloud_cr_repo.default"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "namespace", regexp.MustCompile("tf-testacc-cr-repo-basic-*")),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "name", regexp.MustCompile("tf-testacc-cr-repo-basic-*")),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "summary", "summary"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "detail", "detail"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "repo_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.public"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.vpc"),
				),
			},
		},
	})
}

func TestAccAlicloudCRRepo_Update(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_repo.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRRepo_UpdateBefore(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCRRepoExists("alicloud_cr_repo.default"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "namespace", regexp.MustCompile("tf-testacc-cr-repo-update-*")),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "name", regexp.MustCompile("tf-testacc-cr-repo-update-*")),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "summary", "OLD"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "detail", "OLD"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "repo_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.public"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.vpc"),
				),
			},
			{
				Config: testAccCRRepo_UpdateAfter(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCRRepoExists("alicloud_cr_repo.default"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "namespace", regexp.MustCompile("tf-testacc-cr-repo-update-*")),
					resource.TestMatchResourceAttr("alicloud_cr_repo.default", "name", regexp.MustCompile("tf-testacc-cr-repo-update-*")),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "summary", "NEW"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "detail", "NEW"),
					resource.TestCheckResourceAttr("alicloud_cr_repo.default", "repo_type", "PRIVATE"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.public"),
					resource.TestCheckResourceAttrSet("alicloud_cr_repo.default", "domain_list.vpc"),
				),
			},
		},
	})
}

func testAccCheckCRRepoDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	crService := CrService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_repo" {
			continue
		}
		_, err := crService.DescribeRepo(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return fmt.Errorf("error namespace/repo %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCRRepoExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		crService := CrService{client}

		repo, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("resource not found: %s", n))
		}
		if repo.Primary.ID == "" {
			return WrapError(fmt.Errorf("resource id not set: %s", n))
		}

		_, err := crService.DescribeRepo(repo.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return WrapError(fmt.Errorf("resource not exists: %s %s", n, repo.Primary.ID))
			}
			return WrapError(err)
		}
		return nil
	}
}

func testAccCRRepo_Basic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-repo-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PRIVATE"
}

resource "alicloud_cr_repo" "default" {
	namespace = "${alicloud_cr_namespace.default.name}"
	name = "${var.name}"
	summary = "summary"
	repo_type = "PUBLIC"
	detail  = "detail"
}
`, rand)
}

func testAccCRRepo_UpdateBefore(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-repo-update-%d"
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
`, rand)
}

func testAccCRRepo_UpdateAfter(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-repo-update-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "default" {
	namespace = "${alicloud_cr_namespace.default.name}"
	name = "${var.name}"
	summary = "NEW"
	repo_type = "PRIVATE"
	detail  = "NEW"
}
`, rand)
}
