package sweep_test

// Legacy product sweepers run in this dedicated test binary.

import (
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_auto_snapshot_policy", &resource.Sweeper{
		Name: "alicloud_ecs_auto_snapshot_policy",
		F:    testSweepEcsAutoSnapshotPolicy,
	})
}

func testSweepEcsAutoSnapshotPolicy(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return tferrs.WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	request["PageSize"] = tfhelper.PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	for {
		action := "DescribeAutoSnapshotPolicyEx"
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DataDefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, tferrs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, action, "$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
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
		if len(result) < tfhelper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
