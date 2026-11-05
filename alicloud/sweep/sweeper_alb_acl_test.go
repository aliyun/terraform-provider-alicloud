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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_acl",
		&resource.Sweeper{
			Name: "alicloud_alb_acl",
			F:    testSweepAlbAcl,
		})
}

func testSweepAlbAcl(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListAcls"
	request := map[string]interface{}{}

	request["MaxResults"] = tfhelper.PageSizeXLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
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

		resp, err := jsonpath.Get("$.Acls", response)
		totalCount, formatErr := tfconvert.FormatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Acls", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["AclName"]; !ok {
				continue
			}
			skip := true
			if !tfsweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["AclName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Alb Acl: %s", item["AclName"].(string))
					continue
				}
			}
			action := "DeleteAcl"
			request := map[string]interface{}{
				"AclId": item["AclId"],
			}
			request["ClientToken"] = helper.BuildClientToken("DeleteAcl")
			_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Acl (%s): %s", item["AclId"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Acl success: %s ", item["AclId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}
