package ecs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccAliCloudEcsStorageSet_basic exercises the full lifecycle of
// alicloud_ecs_storage_set: create with resource_group_id and tags, update
// name/description in place (ModifyStorageSetAttribute), and add/change/clear
// tags and add resource_group_id — all of which are ForceNew because ECS
// Storage Set is not a supported ResourceType of JoinResourceGroup,
// TagResources or UntagResources, so those attributes can only be set at
// creation and any change recreates the resource. Clearing tags requires an
// explicit tags = {} (omitting tags is a no-op that preserves the cloud tags).
//
// Note: CreateStorageSet may return 403 InvalidUser.NotSupported on
// accounts without Storage Set capability. This is an account-level
// limitation, not a code defect — see the route plan for context.
func TestAccAliCloudEcsStorageSet_basic(t *testing.T) {
	r := sdkacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsStorageSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsStorageSetConfig(name, "wTest", 2, name, false, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "description", "wTest"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "zone_id", "cn-hangzhou-i"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "max_partition_number", "2"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "storage_set_name", name),
					resource.TestCheckResourceAttrSet("alicloud_ecs_storage_set.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_ecs_storage_set.default", "storage_set_id"),
					resource.TestCheckResourceAttrSet("alicloud_ecs_storage_set.default", "region_id"),
				),
			},
			{
				Config: testAccEcsStorageSetConfig(name, "update", 2, name+"_update", false, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "description", "update"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "storage_set_name", name+"_update"),
				),
			},
			{
				Config: testAccEcsStorageSetConfig(name, "update", 2, name+"_update", true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ecs_storage_set.default", "resource_group_id"),
				),
			},
			{
				Config: testAccEcsStorageSetConfig(name, "update", 2, name+"_update", true, map[string]string{
					"Created": "TF",
					"For":     "Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.For", "Test"),
				),
			},
			{
				Config: testAccEcsStorageSetConfig(name, "update", 2, name+"_update", true, map[string]string{
					"Created": "TF-update",
					"For":     "Test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.Created", "TF-update"),
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.For", "Test-update"),
				),
			},
			{
				// Clearing tags is ForceNew: an explicit tags = {} (NOT omitting
				// tags, which is a no-op that preserves the cloud tags) differs
				// from the prior non-empty map, so the resource is recreated
				// with no tags instead of calling the unsupported UntagResources.
				Config: testAccEcsStorageSetConfig(name, "update", 2, name+"_update", true, map[string]string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_storage_set.default", "tags.%", "0"),
				),
			},
			{
				ResourceName:      "alicloud_ecs_storage_set.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccEcsStorageSetConfig builds the HCL for a Storage Set resource.
// When useResourceGroup is true, the config dynamically creates a real
// resource manager resource group (alicloud_resource_manager_resource_group)
// and wires resource_group_id to its ID, rather than relying on a pre-existing
// resource group fetched via a data source whose second id (ids.1) may not
// exist on the account and would surface as an invalid ResourceGroupId at
// CreateStorageSet. The managed resource group is created once and reused
// across steps; the test harness destroys it at teardown.
// When tags is non-nil, the config includes a tags block.
func testAccEcsStorageSetConfig(name, description string, maxPartition int, storageSetName string, useResourceGroup bool, tags map[string]string) string {
	rgResource := ""
	rgBlock := ""
	if useResourceGroup {
		rgResource = fmt.Sprintf(`
resource "alicloud_resource_manager_resource_group" "default" {
  resource_group_name = %q
  display_name        = %q
}
`, name+"rg", name+"rg-display")
		rgBlock = "  resource_group_id = alicloud_resource_manager_resource_group.default.id\n"
	}

	tagsBlock := ""
	if tags != nil {
		tagsBlock = "  tags = {\n"
		for k, v := range tags {
			tagsBlock += fmt.Sprintf("    %q = %q\n", k, v)
		}
		tagsBlock += "  }\n"
	}

	return fmt.Sprintf(`
%s
resource "alicloud_ecs_storage_set" "default" {
  description          = %q
  zone_id              = "cn-hangzhou-i"
  max_partition_number = %d
  storage_set_name     = %q
%s%s
}
`, rgResource, description, maxPartition, storageSetName, rgBlock, tagsBlock)
}

// testAccCheckEcsStorageSetDestroy verifies that every alicloud_ecs_storage_set
// left in state has been removed: after DeleteStorageSet the resource must no
// longer be visible via DescribeStorageSets. It delegates to
// ecs.CheckEcsStorageSetDestroy, which is target-first: it collects the
// alicloud_ecs_storage_set instances still in state, returns nil only when
// there are none to verify, and when targets remain requires a usable provider
// client (nil or wrong-type meta surfaces as an error rather than a silent
// pass that treats "cannot check" as "destroyed") and reuses the strict
// ecs.StorageSetExists not-found logic from Read so a transient Describe error
// surfaces rather than masking a leaked resource as absent.
func testAccCheckEcsStorageSetDestroy(s *terraform.State) error {
	return ecs.CheckEcsStorageSetDestroy(acctest.Provider.Meta(), s)
}
