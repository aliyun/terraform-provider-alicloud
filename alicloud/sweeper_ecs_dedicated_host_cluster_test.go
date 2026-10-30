package alicloud

// The alicloud_ecs_dedicated_host_cluster sweeper stays in the root test
// binary: the make sweep target runs only ./alicloud, and the SDKv2 sweeper
// helpers (sharedClientForRegion, addDebug, incrementalWait, NeedRetry,
// PageSizeLarge) live in the root alicloud test package. The resource
// implementation itself lives in alicloud/service/ecs.

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ecs_dedicated_host_cluster",
		&resource.Sweeper{
			Name: "alicloud_ecs_dedicated_host_cluster",
			F:    testSweepEcsDedicatedHostCluster,
		})
}

func testSweepEcsDedicatedHostCluster(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDedicatedHostClusters"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DedicatedHostClusters.DedicatedHostCluster", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DedicatedHostClusters.DedicatedHostCluster", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DedicatedHostClusterName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs Dedicated Host Cluster: %s", item["DedicatedHostClusterName"].(string))
				continue
			}
			action := "DeleteDedicatedHostCluster"
			request := map[string]interface{}{
				"DedicatedHostClusterId": item["DedicatedHostClusterId"],
				"RegionId":               client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Dedicated Host Cluster (%s): %s", item["DedicatedHostClusterName"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs Dedicated Host Cluster success: %s ", item["DedicatedHostClusterName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
