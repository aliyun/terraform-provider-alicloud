package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"log"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("alicloud_snapshot_policy", &resource.Sweeper{
		Name: "alicloud_snapshot_policy",
		F:    testSweepSnapshotPolicy,
	})
}

func testSweepSnapshotPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var snapshots []ecs.AutoSnapshotPolicy
	req := ecs.CreateDescribeAutoSnapshotPolicyExRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeAutoSnapshotPolicyEx(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving snapshots: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeAutoSnapshotPolicyExResponse)
		if resp == nil || len(resp.AutoSnapshotPolicies.AutoSnapshotPolicy) < 1 {
			break
		}
		snapshots = append(snapshots, resp.AutoSnapshotPolicies.AutoSnapshotPolicy...)

		if len(resp.AutoSnapshotPolicies.AutoSnapshotPolicy) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range snapshots {
		name := v.AutoSnapshotPolicyName
		id := v.AutoSnapshotPolicyId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping snapshot: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting snapshot: %s (%s)", name, id)
		req := ecs.CreateDeleteAutoSnapshotPolicyRequest()
		req.AutoSnapshotPolicyId = id
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteAutoSnapshotPolicy(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete snapshot(%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSnapshotPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snapshot_policy.sp",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSnapshotPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotPolicyExists("alicloud_snapshot_policy.sp"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "name", "tf-testAcc-sp"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "repeat_weekdays.#", "1"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "retention_days", "-1"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "time_points.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccSnapshotPolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotPolicyExists("alicloud_snapshot_policy.sp"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "name", "tf-testAcc-sp"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "repeat_weekdays.#", "1"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "retention_days", "-1"),
					resource.TestCheckResourceAttr("alicloud_snapshot_policy.sp", "time_points.#", "2"),
				),
			},
		},
	})
}

func testAccCheckSnapshotPolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snapshot" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		_, err := ecsService.DescribeSnapshotPolicy(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describing snapshot policy(%s) failed while destoring, error: %#v.", rs.Primary.ID, err)
		}
		return fmt.Errorf("Error ECS Snapshot Policy (%s) still exist", rs.Primary.ID)
	}

	return nil
}

func testAccCheckSnapshotPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Snapshot ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeSnapshotPolicy(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("While checking snapshot existing, describing snapshot got an error: %#v.", err)
		}

		return nil
	}
}

const testAccSnapshotPolicyConfig = `
resource "alicloud_snapshot_policy" "sp" {
    name = "tf-testAcc-sp"
    repeat_weekdays = [ "1" ]
    retention_days = -1
    time_points = ["1"]
}
`

const testAccSnapshotPolicyConfigUpdate = `
resource "alicloud_snapshot_policy" "sp" {
    name = "tf-testAcc-sp"
    repeat_weekdays = [ "2" ]
    retention_days = -1
    time_points = ["1", "2"]
}
`
