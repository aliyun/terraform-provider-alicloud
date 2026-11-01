package alicloud

// The alicloud_alb_server_group sweeper stays in the root test binary: the make sweep target runs
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
		"alicloud_alb_server_group",
		&resource.Sweeper{
			Name: "alicloud_alb_server_group",
			F:    testSweepAlbServerGroup,
		})
}

func testSweepAlbServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListServerGroups"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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

		resp, err := jsonpath.Get("$.ServerGroups", response)
		totalCount, formatErr := formatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServerGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["ServerGroupName"]; !ok {
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["ServerGroupName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Alb Server Group: %s", item["ServerGroupName"].(string))
					continue
				}
			}
			action := "DeleteServerGroup"
			request := map[string]interface{}{
				"ServerGroupId": item["ServerGroupId"],
			}
			_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Server Group (%s): %s", item["ServerGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Server Group success: %s ", item["ServerGroupName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}
