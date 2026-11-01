package alicloud

// The alicloud_alb_rule sweeper stays in the root test binary: the make sweep target runs
// only ./alicloud, so a sweeper registered from a product package is never
// loaded. The resource implementation and acceptance tests live in
// alicloud/service/alb; keeping only the sweeper here mirrors
// alicloud/sweeper_alb_acl_test.go.

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
		"alicloud_alb_rule",
		&resource.Sweeper{
			Name: "alicloud_alb_rule",
			F:    testSweepAlbRule,
		})
}

func testSweepAlbRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListRules"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Rules", response)

		totalCount, formatErr := formatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Rules", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["RuleName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ALB Rule: %s", item["RuleName"].(string))
				continue
			}

			sweeped = true
			action := "DeleteRule"
			request := map[string]interface{}{
				"RuleId": item["RuleId"],
			}
			request["ClientToken"] = buildClientToken("DeleteRule")
			_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ALB Rule (%s): %s", item["RuleName"].(string), err)
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete ALB Rule success: %s ", item["RuleName"].(string))
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}
