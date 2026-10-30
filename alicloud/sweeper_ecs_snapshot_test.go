package alicloud

// The alicloud_ecs_snapshot sweeper stays in the root test binary: the make
// sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, sweepAll, PageSizeLarge) live in the root alicloud
// test package. The resource implementation itself lives in alicloud/service/ecs.

import (
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_snapshot", &resource.Sweeper{
		Name: "alicloud_ecs_snapshot",
		F:    testSweepEcsSnapshots,
	})
}

func testSweepEcsSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeSnapshots"

	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}

	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_snapshot", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Snapshots.Snapshot", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Snapshots.Snapshot", response)
		}

		result, _ := resp.([]interface{})

		for _, v := range result {
			item := v.(map[string]interface{})

			name := item["SnapshotName"]
			id := item["SnapshotId"]
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping snapshot: %s (%s)", name, id)
					continue
				}
			}
			log.Printf("[INFO] Deleting snapshot: %s (%s)", name, id)
			action = "DeleteSnapshot"
			request := map[string]interface{}{
				"SnapshotId": item["SnapshotId"],
			}

			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)

			if err != nil {
				log.Printf("[ERROR] Failed to delete snapshot(%s (%s)): %s", name, id, err)
			}

			log.Printf("[INFO] Delete snapshot success: %s ", item["SnapshotId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}
