package eci_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudEciCommitContainerTask_basic verifies the create-only
// CommitContainerTask resource: create via CommitContainer, read via
// DescribeCommitContainerTask, import, and that destroy is a no-op that
// removes the resource from state without deleting the task record (no
// false destroy-success) consistent with the create-only contract.
func TestAccAliCloudEciCommitContainerTask_basic(t *testing.T) {
	// ECI container group only supports cn-hangzhou
	// (connectivity.EciContainerGroupRegions = []Region{Hangzhou}). The test
	// depends on alicloud_eci_container_group reaching Running, which only
	// happens in cn-hangzhou; force the provider region so the dependency does
	// not stall in Pending (20m timeout) on other regions like cn-beijing.
	os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	resourceId := "alicloud_eci_commit_container_task.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION") // "cn-hangzhou"
	name := fmt.Sprintf("tf-testacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEciCommitContainerTaskBasic(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttr(resourceId, "container_name", name),
					resource.TestCheckResourceAttr(resourceId, "image.#", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"container_group_id", "container_name", "image", "arn"},
			},
			{
				// CommitContainer has no delete operation; destroy is a no-op that
				// removes the resource from state without deleting the task record
				// (managed by the service), so terraform destroy completes.
				Config:  testAccAliCloudEciCommitContainerTaskBasic(name, region),
				Destroy: true,
			},
		},
	})
}

func testAccAliCloudEciCommitContainerTaskBasic(name, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

# alicloud_eci_zones returns only zones where ECI container groups can be
# created, so the vswitch lands in an ECI-supported zone regardless of the
# provider region (the generic alicloud_zones data source may return zones
# where ECI is unavailable).
data "alicloud_eci_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_eci_zones.default.zones.0.zone_ids.0
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_eci_container_group" "default" {
  container_group_name = var.name
  security_group_id    = alicloud_security_group.default.id
  vswitch_id           = alicloud_vswitch.default.id

  containers {
    name  = var.name
    image = "registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine"
  }
}

# CommitContainer (Eci/2018-08-08) is an asynchronous task API: it accepts the
# image reference and returns a CommitContainerTaskId, then describes the task
# via DescribeCommitContainerTask. The live API requires the Arn parameter
# authorizing the ECI service to access the target ACR registry on behalf of the
# caller: OpenAPI marks Arn optional, but the API returns MissingParameter when
# it is absent, and the API description states you must prepare the ACR
# repository and create a RAM role with the authorization before calling
# CommitContainer. The role below trusts the ECI service principal and is
# granted ACR access so the create path can return a CommitContainerTaskId.
resource "alicloud_ram_role" "default" {
  name     = var.name
  force    = true
  document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["eci.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
EOF
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_type = "System"
  policy_name = "AliyunContainerRegistryFullAccess"
}

resource "alicloud_eci_commit_container_task" "default" {
  container_group_id = alicloud_eci_container_group.default.id
  container_name     = var.name

  arn {
    role_type = "service"
    role_arn  = alicloud_ram_role.default.arn
  }

  image {
    repository = "registry.%s.aliyuncs.com/%s/commit-container-task"
    tag        = "v1.0.0"
  }
}
`, name, region, region, name)
}
