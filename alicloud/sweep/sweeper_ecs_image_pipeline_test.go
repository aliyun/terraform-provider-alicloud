package sweep_test

// Legacy product sweepers run in this dedicated test binary.

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ecs_image_pipeline",
		&resource.Sweeper{
			Name: "alicloud_ecs_image_pipeline",
			F:    testSweepEcsImagePipeline,
		})
}

func testSweepEcsImagePipeline(region string) error {

	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	action := "DescribeImagePipelines"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = tfhelper.PageSizeLarge

	var response map[string]interface{}
	for {
		wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ImagePipeline.ImagePipelineSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ImagePipeline.ImagePipelineSet", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			name, _ := item["Name"].(string)
			if !shouldSweepEcsImagePipeline(name) {
				log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", item["Name"])
				continue
			}
			action := "DeleteImagePipeline"
			request := map[string]interface{}{
				"ImagePipelineId": item["ImagePipelineId"],
				"RegionId":        aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Image Pipeline (%s): %s", item["Name"], err)
			}
			log.Printf("[INFO] Delete Ecs Image Pipeline success: %s ", item["Name"])
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// shouldSweepEcsImagePipeline reports whether a pipeline name belongs to the
// fixtures this sweeper owns. The documentation example uses the exact name
// `terraform-example`; a prefix match on it would also delete unrelated
// pipelines that merely share that prefix, so it is matched exactly. The
// acceptance-test prefixes stay prefix matches.
func shouldSweepEcsImagePipeline(name string) bool {
	lower := strings.ToLower(name)
	if lower == "terraform-example" {
		return true
	}
	for _, prefix := range []string{"tf-testacc", "tf_testacc"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func TestUnitShouldSweepEcsImagePipeline(t *testing.T) {
	cases := map[string]bool{
		"terraform-example":           true,
		"tf-testAccimagepipeline1234": true,
		"tf_testAccimage":             true,
		"terraform-example-2":         false,
		"terraform-examples-backup":   false,
		"production-image-pipeline":   false,
	}
	for name, want := range cases {
		if got := shouldSweepEcsImagePipeline(name); got != want {
			t.Errorf("shouldSweepEcsImagePipeline(%q) = %v, want %v", name, got, want)
		}
	}
}
