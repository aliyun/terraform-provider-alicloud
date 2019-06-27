package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

func TestAccAlicloudSnapshotPolicyBasic(t *testing.T) {

	resourceId := "alicloud_snapshot_policy.default"
	randInt := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSnapshotPolicyBasic%d", randInt)
	basicMap := map[string]string{
		"name":              name,
		"repeat_weekdays.#": "1",
		"retention_days":    "-1",
		"time_points.#":     "1",
	}
	var v *ecs.AutoSnapshotPolicy
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_weekdays": []string{"1"},
					"retention_days":  "-1",
					"time_points":     []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_weekdays": []string{"1"},
					"retention_days":  "-1",
					"time_points":     []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"repeat_weekdays.#": "1",
						"retention_days":    "-1",
						"time_points.#":     "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSnapshotPolicyMulti(t *testing.T) {

	resourceId := "alicloud_snapshot_policy.default.4"
	randInt := acctest.RandIntRange(10000, 99999)
	var v *ecs.AutoSnapshotPolicy
	name := fmt.Sprintf("tf-testAccSnapshotPolicyMulti%d", randInt)
	basicMap := map[string]string{
		"name":              name,
		"repeat_weekdays.#": "1",
		"retention_days":    "-1",
		"time_points.#":     "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":           "5",
					"name":            name,
					"repeat_weekdays": []string{"1"},
					"retention_days":  "-1",
					"time_points":     []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
