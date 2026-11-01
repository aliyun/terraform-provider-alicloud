package alicloud

// The alicloud_ecs_activations sweeper stays in the root test binary: the make
// sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, testSweepPreCheckWithRegions, addDebug,
// incrementalWait, NeedRetry, PageSizeLarge) live in the root alicloud test
// package. The resource implementation itself lives in alicloud/service/ecs.

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
		"alicloud_ecs_activations",
		&resource.Sweeper{
			Name: "alicloud_ecs_activations",
			F:    testSweepEcsActivations,
		})
}

func testSweepEcsActivations(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EcsActivationsSupportRegions) {
		log.Printf("[INFO] Skipping Ecs Activations unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeActivations"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ActivationList", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ActivationList", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["InstanceName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs Activations: %s", item["InstanceName"].(string))
				continue
			}
			action := "DeleteActivation"
			request := map[string]interface{}{
				"ActivationId": item["ActivationId"],
				"RegionId":     aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Activations (%s): %s", item["InstanceName"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs Activations success: %s ", item["InstanceName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
