package sweep_test

// Legacy product sweepers run in this dedicated test binary.

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_load_balancer",
		&resource.Sweeper{
			Name: "alicloud_alb_load_balancer",
			F:    testSweepAlbLoadBalancer,
		})
}

func testSweepAlbLoadBalancer(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListLoadBalancers"
	request := map[string]interface{}{
		"MaxResults": tfhelper.PageSizeXLarge,
	}
	var response map[string]interface{}
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.LoadBalancers", response)

	totalCount, formatErr := tfconvert.FormatInt(response["TotalCount"])
	if formatErr != nil {
		return formatErr
	}
	if totalCount != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.LoadBalancers", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if _, ok := item["LoadBalancerName"]; !ok {
			continue
		}
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["LoadBalancerName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ALB LoadBalancer: %s", item["LoadBalancerName"].(string))
			continue
		}

		action := "DeleteLoadBalancer"
		request := map[string]interface{}{
			"LoadBalancerId": item["LoadBalancerId"],
		}
		request["ClientToken"] = tfhelper.BuildClientToken("DeleteLoadBalancer")
		_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete ALB LoadBalancer (%s): %s", item["LoadBalancerId"].(string), err)
		}
		log.Printf("[INFO] Delete ALB LoadBalancer success: %s ", item["LoadBalancerId"].(string))
	}
	return nil
}
