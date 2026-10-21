package eci_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEciCommitContainerTask_basic verifies the create-only
// CommitContainerTask resource: create via CommitContainer, read via
// DescribeCommitContainerTask, import, a two-task read identity check, and that
// destroy removes the resource from state without deleting the retained task
// record (the create-only contract has no delete API).
func TestAccAliCloudEciCommitContainerTask_basic(t *testing.T) {
	// ECI container group only supports cn-hangzhou
	// (connectivity.EciContainerGroupRegions = []Region{Hangzhou}). The test
	// depends on alicloud_eci_container_group reaching Running, which only
	// happens in cn-hangzhou; force the provider region so the dependency does
	// not stall in Pending (20m timeout) on other regions like cn-beijing.
	os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	resourceId := "alicloud_eci_commit_container_task.default"
	secondId := "alicloud_eci_commit_container_task.second"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION") // "cn-hangzhou"
	// ECI ContainerGroupName only accepts lowercase letters, digits and hyphens, so the
	// standard tf-testAcc prefix is rejected by the API (InvalidParameter); use the
	// lowercase variant and keep the random suffix for uniqueness.
	name := fmt.Sprintf("tf-testacc%d", rand)
	var taskId, secondTaskId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckEciCommitContainerTaskRetained(&taskId),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEciCommitContainerTaskBasic(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttr(resourceId, "container_name", name),
					resource.TestCheckResourceAttr(resourceId, "image.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "image.0.author", "terraform"),
					resource.TestCheckResourceAttr(resourceId, "image.0.message", "built by terraform"),
					testAccCaptureEciCommitContainerTaskIds(resourceId, secondId, &taskId, &secondTaskId),
				),
			},
			{
				// Read must select the task by the resource identity instead of the
				// first element of the Describe response: importing the second task
				// has to surface the second task id even though both tasks are
				// visible in the same region.
				ResourceName:            secondId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"container_group_id", "image", "arn"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					if got := states[0].Attributes["id"]; got != secondTaskId {
						return fmt.Errorf("imported id = %q, want %q", got, secondTaskId)
					}
					if got := states[0].Attributes["container_name"]; got != name {
						return fmt.Errorf("imported container_name = %q, want %q", got, name)
					}
					return nil
				},
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

// TestAccAliCloudEciCommitContainerTask_missingArn verifies that the live
// CommitContainer API rejects a create without the arn block. The OpenAPI
// metadata marks Arn as optional, so this negative case is the only evidence
// that the documented requirement is real.
func TestAccAliCloudEciCommitContainerTask_missingArn(t *testing.T) {
	os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	name := fmt.Sprintf("tf-testacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudEciCommitContainerTaskMissingArn(name, region),
				ExpectError: regexp.MustCompile(`(?i)(MissingParameter|Missing.*Arn|Arn.*required)`),
			},
		},
	})
}

func testAccCaptureEciCommitContainerTaskIds(resourceId, secondId string, taskId, secondTaskId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no id", resourceId)
		}
		rs2, ok := s.RootModule().Resources[secondId]
		if !ok || rs2.Primary.ID == "" {
			return fmt.Errorf("resource %s has no id", secondId)
		}
		if rs.Primary.ID == rs2.Primary.ID {
			return fmt.Errorf("both commit container tasks resolved to the same id %q", rs.Primary.ID)
		}
		*taskId, *secondTaskId = rs.Primary.ID, rs2.Primary.ID
		return nil
	}
}

// testAccCheckEciCommitContainerTaskRetained is the create-only CheckDestroy:
// it never relies on the resource Read path and asserts the opposite of the
// usual destroy check - the task record must still be visible after destroy,
// because CommitContainer has no delete API and the service retains the record.
func testAccCheckEciCommitContainerTaskRetained(taskId *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if taskId == nil || *taskId == "" {
			return fmt.Errorf("no commit container task id was captured before destroy")
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"RegionId": client.RegionId,
			"TaskId.1": *taskId,
		}
		raw, err := client.RpcPost("Eci", "2018-08-08", "DescribeCommitContainerTask", nil, request, true)
		if err != nil {
			return fmt.Errorf("describing retained task %s after destroy: %s", *taskId, err)
		}
		tasks, _ := raw["CommitTasks"].([]interface{})
		for _, item := range tasks {
			object, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(object["TaskId"]) == *taskId {
				return nil
			}
		}
		return fmt.Errorf("create-only task %s is no longer visible after destroy; the service is expected to retain the task record", *taskId)
	}
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
# via DescribeCommitContainerTask. The API description requires a RAM role
# trusted by the ECI service with access to the target ACR repository before
# the task can be submitted.
resource "alicloud_ram_role" "default" {
  role_name                   = var.name
  force                       = true
  assume_role_policy_document = <<EOF
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
  role_name   = alicloud_ram_role.default.role_name
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
    author     = "terraform"
    message    = "built by terraform"
  }
}

resource "alicloud_eci_commit_container_task" "second" {
  depends_on = [alicloud_eci_commit_container_task.default]

  container_group_id = alicloud_eci_container_group.default.id
  container_name     = var.name

  arn {
    role_type = "service"
    role_arn  = alicloud_ram_role.default.arn
  }

  image {
    repository = "registry.%s.aliyuncs.com/%s/commit-container-task"
    tag        = "v2.0.0"
  }
}
`, name, region, region, name, region, name)
}

func testAccAliCloudEciCommitContainerTaskMissingArn(name, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

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

resource "alicloud_eci_commit_container_task" "default" {
  container_group_id = alicloud_eci_container_group.default.id
  container_name     = var.name

  image {
    repository = "registry.%s.aliyuncs.com/%s/commit-container-task"
    tag        = "v1.0.0"
  }
}
`, name, region, region, name)
}
