package alicloud

// The alicloud_ecs_deployment_set sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
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
		"alicloud_ecs_deployment_set",
		&resource.Sweeper{
			Name: "alicloud_ecs_deployment_set",
			F:    testSweepEcsDeploymentSet,
		})
}

func testSweepEcsDeploymentSet(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeDeploymentSets"
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

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
		resp, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DeploymentSets.DeploymentSet", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["DeploymentSetName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DeploymentSetName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs DeploymentSet: %s", item["DeploymentSetName"].(string))
				continue
			}
			action := "DeleteDeploymentSet"
			request := map[string]interface{}{
				"DeploymentSetId": item["DeploymentSetId"],
				"RegionId":        client.RegionId,
			}
			request["ClientToken"] = buildClientToken("DeleteDeploymentSet")
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs DeploymentSet (%s): %s", item["DeploymentSetId"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs DeploymentSet success: %s ", item["DeploymentSetId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}
