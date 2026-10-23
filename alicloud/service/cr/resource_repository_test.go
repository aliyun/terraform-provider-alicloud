// Package cr. This file is generated automatically. Please do not modify it manually, thank you!
package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Test Cr Repository. >>> Resource test cases, automatically generated.
// Case repo 生命周期 5773
func TestAccAliCloudCrRepository_basic5773(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCrRepository_basic5773_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "repo_namespace_name"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "tag_immutability", "false"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repo_type", "PUBLIC"),

					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "summary", "sdad"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repository_name", name),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "detail", "ds"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("create_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("modified_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("repo_build_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCrRepository_basic5773_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "repo_namespace_name"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "tag_immutability", "true"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repo_type", "PUBLIC"),

					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "summary", "asfdsf"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repository_name", name),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "detail", "aesfafe"),
				),
			},
			{
				ResourceName:            "alicloud_cr_repository.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCrRepository_basic5773_base(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_cr_ee_instance" "default1TiZyc" {
  instance_name  = var.name
  payment_type   = "Subscription"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}

resource "alicloud_cr_ee_namespace" "defaultRZzh9M" {
  instance_id        = alicloud_cr_ee_instance.default1TiZyc.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}


`, name)
}

func testAccAliCloudCrRepository_basic5773_step0(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_cr_repository" "default" {
        repo_namespace_name = alicloud_cr_ee_namespace.defaultRZzh9M.name

        tag_immutability = false

        repo_type = "PUBLIC"

        instance_id = "${alicloud_cr_ee_instance.default1TiZyc.id}"

        summary = "sdad"

        repository_name = "%[1]s"
        detail = "ds"

}
`, name, testAccAliCloudCrRepository_basic5773_base(name))
}

func testAccAliCloudCrRepository_basic5773_step1(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_cr_repository" "default" {
        repo_namespace_name = alicloud_cr_ee_namespace.defaultRZzh9M.name

        tag_immutability = true

        repo_type = "PUBLIC"

        instance_id = "${alicloud_cr_ee_instance.default1TiZyc.id}"

        summary = "asfdsf"

        repository_name = "%[1]s"
        detail = "aesfafe"

}
`, name, testAccAliCloudCrRepository_basic5773_base(name))
}

// Case repo 生命周期_简化 5874
func TestAccAliCloudCrRepository_basic5874(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCrRepository_basic5874_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "repo_namespace_name"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "tag_immutability", "false"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repo_type", "PUBLIC"),

					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "summary", "sdad"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repository_name", name),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "detail", "ds"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("create_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("modified_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("repo_build_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cr_repository.default", tfjsonpath.New("repository_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCrRepository_basic5874_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "repo_namespace_name"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "tag_immutability", "true"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repo_type", "PUBLIC"),

					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "summary", "asfdsf"),

					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "repository_name", name),
					resource.TestCheckResourceAttr("alicloud_cr_repository.default", "detail", "aesfafe"),
				),
			},
			{
				ResourceName:            "alicloud_cr_repository.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCrRepository_basic5874_base(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_cr_ee_instance" "default1TiZyc" {
  instance_name  = var.name
  payment_type   = "Subscription"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Standard"
}

resource "alicloud_cr_ee_namespace" "defaultRZzh9M" {
  instance_id        = alicloud_cr_ee_instance.default1TiZyc.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}


`, name)
}

func testAccAliCloudCrRepository_basic5874_step0(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_cr_repository" "default" {
        repo_namespace_name = alicloud_cr_ee_namespace.defaultRZzh9M.name

        tag_immutability = false

        repo_type = "PUBLIC"

        instance_id = "${alicloud_cr_ee_instance.default1TiZyc.id}"

        summary = "sdad"

        repository_name = "%[1]s"
        detail = "ds"

}
`, name, testAccAliCloudCrRepository_basic5874_base(name))
}

func testAccAliCloudCrRepository_basic5874_step1(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_cr_repository" "default" {
        repo_namespace_name = alicloud_cr_ee_namespace.defaultRZzh9M.name

        tag_immutability = true

        repo_type = "PUBLIC"

        instance_id = "${alicloud_cr_ee_instance.default1TiZyc.id}"

        summary = "asfdsf"

        repository_name = "%[1]s"
        detail = "aesfafe"

}
`, name, testAccAliCloudCrRepository_basic5874_base(name))
}

// Test Cr Repository. <<< Resource test cases, automatically generated.
