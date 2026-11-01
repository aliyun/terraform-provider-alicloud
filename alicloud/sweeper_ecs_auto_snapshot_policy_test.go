package alicloud

// The alicloud_ecs_auto_snapshot_policy sweeper stays in the root test binary:
// the make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, testSweepPreCheckWithRegions, addDebug,
// incrementalWait, NeedRetry, PageSizeLarge) live in the root alicloud test
// package. The resource implementation itself lives in alicloud/service/ecs.

import (
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_auto_snapshot_policy", &resource.Sweeper{
		Name: "alicloud_ecs_auto_snapshot_policy",
		F:    testSweepEcsAutoSnapshotPolicy,
	})
}

func testSweepEcsAutoSnapshotPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	for {
		action := "DescribeAutoSnapshotPolicyEx"
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AutoSnapshotPolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs SnapShot Policy: %s (%s)", item["AutoSnapshotPolicyName"], item["AutoSnapshotPolicyId"])
				continue
			}
			action = "DeleteAutoSnapshotPolicy"
			request := map[string]interface{}{
				"autoSnapshotPolicyId": item["AutoSnapshotPolicyId"],
				"regionId":             client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs SnapShot Policy (%s (%s)): %s", item["AutoSnapshotPolicyName"], item["InstanceId"], err)
				continue
			}
			log.Printf("[INFO] Delete Ecs SnapShot Policy success: %s ", item["AutoSnapshotPolicyId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
