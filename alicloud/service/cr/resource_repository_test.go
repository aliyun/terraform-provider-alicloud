// Package cr. This file is generated automatically. Please do not modify it manually, thank you!
package cr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// testAccCheckAliCloudCrRepositoryDestroy proves the repository no longer exists after destroy.
// Only an explicit repository/parent-instance NotFound passes; a still-present repository, an
// unusable success flag or any other failure keeps the check failing with the target ID.
func testAccCheckAliCloudCrRepositoryDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_repository" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 2 {
			return fmt.Errorf("unexpected alicloud_cr_repository id %q", rs.Primary.ID)
		}
		request := map[string]interface{}{"RegionId": client.RegionId, "InstanceId": parts[0], "RepoId": parts[1]}
		response, err := client.RpcPost("cr", "2018-12-01", "GetRepository", nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST", "INSTANCE_NOT_EXIST", "ResourceNotFound.Repository", "ResourceNotFound"}) {
				continue
			}
			return fmt.Errorf("checking alicloud_cr_repository %s: %s", rs.Primary.ID, err)
		}
		if response == nil {
			return fmt.Errorf("checking alicloud_cr_repository %s: empty GetRepository response", rs.Primary.ID)
		}
		code, _ := response["Code"].(string)
		if code == "REPO_NOT_EXIST" || code == "INSTANCE_NOT_EXIST" {
			continue
		}
		success, ok := response["IsSuccess"].(bool)
		if !ok {
			return fmt.Errorf("checking alicloud_cr_repository %s: response has no usable IsSuccess flag", rs.Primary.ID)
		}
		if success {
			return fmt.Errorf("alicloud_cr_repository %s still exists", rs.Primary.ID)
		}
		return fmt.Errorf("checking alicloud_cr_repository %s: IsSuccess=false, Code=%s", rs.Primary.ID, code)
	}
	return nil
}

// Test Cr Repository. >>> Resource test cases, automatically generated.
// Case repo 生命周期 5773
func TestAccAliCloudCrRepository_basic5773(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudCrRepositoryDestroy,
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
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "resource_group_id"),
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
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudCrRepositoryDestroy,
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
					resource.TestCheckResourceAttrSet("alicloud_cr_repository.default", "resource_group_id"),
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
