package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudRealtimeComputeMembersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-flinkmember-ds%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudRealtimeComputeMembersDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_realtime_compute_members.default", "members.0.member"),
					resource.TestCheckResourceAttrSet("data.alicloud_realtime_compute_members.default", "members.0.role"),
					resource.TestCheckResourceAttrSet("data.alicloud_realtime_compute_members.default", "total_size"),
					testAccCheckAlicloudRealtimeComputeMembersContainsViewer("data.alicloud_realtime_compute_members.default", "alicloud_ram_user.user"),
				),
			},
		},
	})
}

// testAccCheckAlicloudRealtimeComputeMembersContainsViewer verifies that the member created
// by the test (role=VIEWER) is present in the data source's members list, without
// relying on a fixed member count: the Flink workspace is created with default
// owner members in addition to the viewer added by the test, so the returned
// total size is not deterministic. role is compared case-insensitively because
// the Flink member API normalizes role to lowercase on read.
func testAccCheckAlicloudRealtimeComputeMembersContainsViewer(membersResource, ramUserResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		mrs, ok := s.RootModule().Resources[membersResource]
		if !ok || mrs.Primary == nil {
			return fmt.Errorf("resource not found in state: %s", membersResource)
		}
		urs, ok := s.RootModule().Resources[ramUserResource]
		if !ok || urs.Primary == nil {
			return fmt.Errorf("resource not found in state: %s", ramUserResource)
		}
		userID := urs.Primary.Attributes["id"]
		if userID == "" {
			return fmt.Errorf("ram user id is empty for resource %s", ramUserResource)
		}
		count, err := strconv.Atoi(mrs.Primary.Attributes["members.#"])
		if err != nil {
			return fmt.Errorf("failed to parse members count for %s: %v", membersResource, err)
		}
		for i := 0; i < count; i++ {
			if mrs.Primary.Attributes[fmt.Sprintf("members.%d.member", i)] == userID &&
				strings.EqualFold(mrs.Primary.Attributes[fmt.Sprintf("members.%d.role", i)], "VIEWER") {
				return nil
			}
		}
		return fmt.Errorf("viewer member (member=%s, role=VIEWER) not found in %s (members count=%d)", userID, membersResource, count)
	}
}

// testAccAlicloudRealtimeComputeMembersDataSourceConfig provisions the backing resources
// (ram user / vpc / vswitch / oss bucket / vvp instance / flink member) and the data
// source in a single config. The data source's resource_id and namespace inputs
// reference the flink member resource's computed attributes instead of using
// depends_on: this creates a natural apply-graph dependency edge (data source Read
// runs only after the member resource's Create+Read have completed), which both
// guarantees the viewer member is already visible to ListRealtimeComputeFlinkMembers
// (the member Read/Describe retries until the member is visible) and avoids the
// SDK v1 (terraform 0.12.7) data-source state pruning ("DestroyEdgeTransformer:
// pruning unused resource node") observed when depends_on is used to introduce a
// data source in a later step — that pruning left the data source in state with
// empty attributes and produced a non-empty refresh plan. The single-config
// reference pattern matches every other alicloud list data source test (none use
// depends_on or ExpectNonEmptyPlan).
func testAccAlicloudRealtimeComputeMembersDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ram_user" "user" {
  name = var.name
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-flink-member-ds"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-flink-member-ds"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch.zone_id
}

resource "alicloud_realtime_compute_member" "default" {
  resource_id = alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id
  namespace    = "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default"
  member       = alicloud_ram_user.user.id
  role         = "VIEWER"
}

data "alicloud_realtime_compute_members" "default" {
  resource_id = alicloud_realtime_compute_member.default.resource_id
  namespace   = alicloud_realtime_compute_member.default.namespace
}
`, name)
}
