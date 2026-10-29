package imm_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers(
		"alicloud_imm_project",
		&resource.Sweeper{
			Name: "alicloud_imm_project",
			F:    testSweepIMMProject,
		})
}

func testSweepIMMProject(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "ListProjects"
	request := map[string]interface{}{
		"MaxKeys":  helper.PageSizeLarge,
		"RegionId": client.RegionId,
	}

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("imm", "2017-09-06", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.Projects", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Projects", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		sweeped := false
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Project"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping IMM Project: %s", item["Project"].(string))
				continue
			}

			sweeped = true
			action := "DeleteProject"
			var response map[string]interface{}
			request := map[string]interface{}{
				"Project": item["Project"],
			}

			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("imm", "2017-09-06", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				log.Printf("[ERROR] %s get an error: %#v", action, err)
				return nil
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete IMM Project success: %s ", item["Project"].(string))

			if marker, ok := response["Marker"].(string); ok && marker != "" {
				request["Marker"] = marker
			} else {
				break
			}
		}

		return nil
	}
}
