package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudOosOpsItem covers the create/update/import lifecycle of the
// Ops Item. It exercises the governance-mapped fields: status (applied via a
// post-create update since CreateOpsItem has no Status parameter), tags,
// resources and solutions, alongside the core title/severity/category/source/
// priority/description. Destroying the resource transitions the Ops Item to
// Resolved because OOS exposes no DeleteOpsItem API.
func TestAccAliCloudOosOpsItem(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosOpsItemConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item.default", "ops_item_id"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "title", fmt.Sprintf("tf-testacc-opsitem-%d", r)),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "severity", "High"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "category", "Security"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "source", "tf-testacc"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "priority", "2"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "status", "InProgress"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "tags.Env", "tf-test"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "resources.#", "1"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "solutions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "solutions.0.priority", "3"),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item.default", "create_date"),
				),
			},
			{
				Config: testAccAliCloudOosOpsItemConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "severity", "Low"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "category", "Cost"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "status", "Resolved"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "tags.Env", "tf-test-2"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "tags.Owner", "tf"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "resources.#", "1"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "solutions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "solutions.0.description", fmt.Sprintf("tf-testacc solution-upd-%d", r)),
				),
			},
			{
				ResourceName:      "alicloud_oos_ops_item.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudOosOpsItemDedupString covers the create-only, write-only
// dedup_string field. OOS does not return dedup_string on GetOpsItem, so the
// configured value is preserved in state and cannot be verified by import;
// this is a standalone create-only case with no import step, per the
// write-only ForceNew coverage guidance.
func TestAccAliCloudOosOpsItemDedupString(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosOpsItemDedupStringConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item.default", "ops_item_id"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "severity", "Medium"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item.default", "dedup_string", fmt.Sprintf("tf-testacc-dedup-%d", r)),
				),
			},
		},
	})
}

func testAccAliCloudOosOpsItemConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_oos_ops_item" "default" {
  title       = "tf-testacc-opsitem-%[1]d"
  severity    = "High"
  category    = "Security"
  source      = "tf-testacc"
  priority    = 2
  description = "tf-testacc ops item %[1]d"
  status      = "InProgress"
  tags = {
    Env  = "tf-test"
    Name = "ops-item-%[1]d"
  }
  resources = [
    "acs:oos:cn-hangzhou:1563457855438522:application/tf-testacc-%[1]d/applicationgroup/default",
  ]
  solutions {
    priority    = 3
    type        = "SingleAZEcs"
    url         = "https://example.com/%[1]d"
    description = "tf-testacc solution %[1]d"
  }
}
`, r)
}

func testAccAliCloudOosOpsItemConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_oos_ops_item" "default" {
  title       = "tf-testacc-opsitem-%[1]d"
  severity    = "Low"
  category    = "Cost"
  source      = "tf-testacc"
  priority    = 5
  description = "tf-testacc ops item %[1]d updated"
  status      = "Resolved"
  tags = {
    Env   = "tf-test-2"
    Owner = "tf"
  }
  resources = [
    "acs:oos:cn-hangzhou:1563457855438522:application/tf-testacc-upd-%[1]d/applicationgroup/default",
  ]
  solutions {
    priority    = 1
    type        = "SingleAZEcs"
    url         = "https://example.com/%[1]d-upd"
    description = "tf-testacc solution-upd-%[1]d"
  }
}
`, r)
}

func testAccAliCloudOosOpsItemDedupStringConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_oos_ops_item" "default" {
  title       = "tf-testacc-opsitem-dedup-%[1]d"
  severity    = "Medium"
  source      = "tf-testacc"
  priority    = 2
  dedup_string = "tf-testacc-dedup-%[1]d"
}
`, r)
}
