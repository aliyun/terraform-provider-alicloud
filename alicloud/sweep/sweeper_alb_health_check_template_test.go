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
		"alicloud_alb_health_check_template",
		&resource.Sweeper{
			Name: "alicloud_alb_health_check_template",
			F:    testSweepAlbHealthCheckTemplate,
		})
}

func testSweepAlbHealthCheckTemplate(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListHealthCheckTemplates"
	request := map[string]interface{}{
		"MaxResults": tfhelper.PageSizeLarge,
	}
	var response map[string]interface{}
	for {
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

		resp, err := jsonpath.Get("$.HealthCheckTemplates", response)

		totalCount, formatErr := tfconvert.FormatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.HealthCheckTemplates", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["HealthCheckTemplateName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ALB HealthCheckTemplate: %s", item["HealthCheckTemplateName"].(string))
				continue
			}

			sweeped = true
			action := "DeleteHealthCheckTemplates"
			request := map[string]interface{}{
				"HealthCheckTemplateIds.1": item["HealthCheckTemplateId"],
			}
			request["ClientToken"] = tfhelper.BuildClientToken("DeleteHealthCheckTemplate")
			_, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ALB HealthCheckTemplate (%s): %s", item["HealthCheckTemplateName"].(string), err)
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete ALB HealthCheckTemplate success: %s ", item["HealthCheckTemplateName"].(string))
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}
