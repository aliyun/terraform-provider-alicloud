package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// These cases mirror the SDKv2 PortRangeList acceptance tests by name
// (basic12358 / 12361 / 12362 / 12363) so the remote ACC harness can match
// them, but are rewritten as Framework-native, self-contained HCL using
// ProtoV5ProviderFactories. Each step carries the full configuration forward
// (the SDKv2 harness merged per-step deltas; Framework configs do not), so a
// step that only changes tags still repeats the entry/max_entries/name fields.

const (
	prlResourceID     = "alicloud_ecs_port_range_list.default"
	prlSecondResource = "alicloud_ecs_port_range_list.second"
)

func testAccEcsPortRangeListDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

// testAccEcsPortRangeListCheckDestroy is the fail-closed destroy verification.
// Delete only returns once DeletePortRangeList succeeds, so a green Delete step
// proves the destroy call was accepted — not that the cloud resource is gone.
// This re-describes each managed id through the resource's own Read call chain
// (the strict DescribePortRangeLists shape guards): a still-present object, a
// malformed/truncated response, a transport or permission failure, or a missing
// client all fail the check; only a structurally-correct empty array (or the
// absence of the id) confirms the resource is gone. No polling is added.
func testAccEcsPortRangeListCheckDestroy(s *terraform.State) error {
	client, ok := acctest.Provider.Meta().(*connectivity.AliyunClient)
	if !ok || client == nil {
		return fmt.Errorf("CheckDestroy: provider client is unavailable; cannot verify ECS PortRangeList deletion (fail-closed)")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_port_range_list" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		object, err := ecs.DescribePortRangeListForTest(client, rs.Primary.ID)
		if err != nil {
			// describePortRangeList returns (nil, nil) only for a clean empty
			// array. Any error means the resource's absence is not confirmed,
			// so the destroy check fails closed rather than guessing it is gone.
			return fmt.Errorf("CheckDestroy: DescribePortRangeLists for %s returned an error; the resource may still exist: %w", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("CheckDestroy: ECS PortRangeList %s still exists after destroy", rs.Primary.ID)
		}
		// object == nil, err == nil: structurally-correct empty array → gone.
	}
	return nil
}

// Case PortRangeList_生命周期 12358: create with one entry, update entry and
// force-recreate on resource_group_id change, then add/update/remove tags and
// import.
func TestAccAliCloudEcsPortRangeList_basic12358(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccEcsPortRangeListCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  entry {
    port_range  = "80/80"
    description = "resource_test_80"
  }
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "10"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "1"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
					resource.TestCheckResourceAttrSet(prlResourceID, "association_count"),
					resource.TestCheckResourceAttrSet(prlResourceID, "create_time"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 15
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "22/22"
    description = "Entry_test_22"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "15"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "1"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test_update"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name+"_update"),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 15
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "22/22"
    description = "Entry_test_22"
  }

  tags = {
    Created = "TF"
    For     = "Test"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 15
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "22/22"
    description = "Entry_test_22"
  }

  tags = {
    Created = "TF-update"
    For     = "Test-update"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF-update"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test-update"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 15
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "22/22"
    description = "Entry_test_22"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
			{
				ResourceName:      prlResourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case PortRangeList_生命周期_ Entry和Tag_置空元素 12361: create with no entry,
// update scalar fields, then exercise tags add/update/remove and import. The
// empty-entry path keeps the entry set null so it matches a config without an
// entry block (a null set has no entry.# key; ImportStateVerify covers it).
func TestAccAliCloudEcsPortRangeList_basic12361(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccEcsPortRangeListCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "10"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
					resource.TestCheckResourceAttrSet(prlResourceID, "association_count"),
					resource.TestCheckResourceAttrSet(prlResourceID, "create_time"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "20"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test_update"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name+"_update"),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  tags = {
    Created = "TF"
    For     = "Test"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  tags = {
    Created = "TF-update"
    For     = "Test-update"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF-update"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test-update"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
			{
				ResourceName:      prlResourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case PortRangeList_生命周期_ Entry和Tag_3个元素 12362: create with three
// entries, update all three to different ranges, then tags add/update/remove
// and import. Exercises diffEntries with multiple add/remove pairs.
func TestAccAliCloudEcsPortRangeList_basic12362(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccEcsPortRangeListCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  entry {
    port_range  = "5/5"
    description = "resource_test"
  }
  entry {
    port_range  = "2/2"
    description = "resource_test_2"
  }
  entry {
    port_range  = "1/1"
    description = "resource_test_1"
  }
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "10"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "3"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
					resource.TestCheckResourceAttrSet(prlResourceID, "association_count"),
					resource.TestCheckResourceAttrSet(prlResourceID, "create_time"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "10/10"
    description = "resource_test_update"
  }
  entry {
    port_range  = "5/5"
    description = "resource_test_update_5"
  }
  entry {
    port_range  = "2/2"
    description = "resource_test_update_2"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "20"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "3"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test_update"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name+"_update"),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "10/10"
    description = "resource_test_update"
  }
  entry {
    port_range  = "5/5"
    description = "resource_test_update_5"
  }
  entry {
    port_range  = "2/2"
    description = "resource_test_update_2"
  }

  tags = {
    Created = "TF"
    For     = "Test"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "10/10"
    description = "resource_test_update"
  }
  entry {
    port_range  = "5/5"
    description = "resource_test_update_5"
  }
  entry {
    port_range  = "2/2"
    description = "resource_test_update_2"
  }

  tags = {
    Created = "TF-update"
    For     = "Test-update"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "2"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF-update"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.For", "Test-update"),
				),
			},
			{
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 20
  port_range_list_name = "%s_update"
  description          = "resource_test_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "10/10"
    description = "resource_test_update"
  }
  entry {
    port_range  = "5/5"
    description = "resource_test_update_5"
  }
  entry {
    port_range  = "2/2"
    description = "resource_test_update_2"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
			{
				ResourceName:      prlResourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case PortRangeList_多实例Read隔离 12363: a second PortRangeList resource is
// kept alive alongside the primary one so every Read of the primary must return
// the list identified by its own PortRangeListId (the PortRangeListId.1
// server-side filter), not drift onto the sibling list. The second resource's
// attributes are re-asserted after each primary update.
func TestAccAliCloudEcsPortRangeList_basic12363(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	secondName := name + "_second"
	dependence := testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "second" {
  max_entries          = 10
  port_range_list_name = "%s"
  description          = "second_isolation"

  entry {
    port_range  = "443/443"
    description = "second_entry"
  }
}
`, secondName)

	secondChecks := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prlSecondResource, "port_range_list_name", secondName),
		resource.TestCheckResourceAttr(prlSecondResource, "description", "second_isolation"),
		resource.TestCheckResourceAttr(prlSecondResource, "entry.#", "1"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccEcsPortRangeListCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: dependence + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "first_isolation"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  entry {
    port_range  = "80/80"
    description = "first_entry"
  }
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "10"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "1"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "first_isolation"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
					resource.TestCheckResourceAttrSet(prlResourceID, "association_count"),
					resource.TestCheckResourceAttrSet(prlResourceID, "create_time"),
					secondChecks,
				),
			},
			{
				Config: dependence + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 15
  port_range_list_name = "%s_update"
  description          = "first_isolation_update"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[1]

  entry {
    port_range  = "22/22"
    description = "first_entry_update"
  }
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "15"),
					resource.TestCheckResourceAttr(prlResourceID, "entry.#", "1"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "first_isolation_update"),
					resource.TestCheckResourceAttr(prlResourceID, "port_range_list_name", name+"_update"),
					resource.TestCheckResourceAttrSet(prlResourceID, "resource_group_id"),
					secondChecks,
				),
			},
			{
				ResourceName:      prlResourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case PortRangeList_生命周期_ explicit-empty-tags 12364: the D1 null-vs-empty
// consistency path the historical cases never configured. The cloud API reports
// no tags for every step, so each re-apply would drift an explicit tags = {} to
// null under the old unconditional MapNull readback. The steps exercise: create
// with an explicit empty map; add a non-empty set; return to an explicit empty
// map (the non-empty -> {} path); then drop the block to null (the {} -> null
// path). ImportStateVerify is intentionally omitted here because the import step
// would read a null/empty plan back over an explicit empty and conflate the two
// semantics; import is already covered by the lifecycle cases above, and this
// case isolates the empty-tag transitions.
func TestAccAliCloudEcsPortRangeList_basic12364(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccEcsPortRangeListCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Explicit tags = {} on create: the configured empty map must
				// survive the readback as a non-null, zero-element map.
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  tags = {}
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "max_entries", "10"),
					resource.TestCheckResourceAttr(prlResourceID, "description", "resource_test"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
			{
				// Non-empty tags applied over an explicit-empty prior state.
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  tags = {
    Created = "TF"
  }
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "1"),
					resource.TestCheckResourceAttr(prlResourceID, "tags.Created", "TF"),
				),
			},
			{
				// D1 non-empty -> {} transition: remove every tag but keep the
				// block explicit-empty; the readback must stay a zero-element
				// map, not drift to null.
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  tags = {}
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
			{
				// {} -> null: omit the block so the plan is null and the API
				// reports no tags; the readback stays null (tags.% == 0).
				Config: testAccEcsPortRangeListDependence(name) + fmt.Sprintf(`
resource "alicloud_ecs_port_range_list" "default" {
  max_entries          = 10
  port_range_list_name = var.name
  description          = "resource_test"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prlResourceID, "tags.%", "0"),
				),
			},
		},
	})
}
