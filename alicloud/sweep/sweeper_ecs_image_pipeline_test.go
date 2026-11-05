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

	deleteErrors := make([]string, 0)

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
			return fmt.Errorf("%s failed to list image pipelines in %s: %s", action, region, err)
		}

		resp, err := jsonpath.Get("$.ImagePipeline.ImagePipelineSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ImagePipeline.ImagePipelineSet", action, err)
			return fmt.Errorf("%s returned an unexpected image pipeline list in %s: %s", action, region, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				deleteErrors = append(deleteErrors, fmt.Sprintf("malformed pipeline record %T", v))
				continue
			}

			name, _ := item["Name"].(string)
			if !shouldSweepEcsImagePipeline(name, imagePipelineTagMap(item["Tags"])) {
				log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", name)
				continue
			}
			action := "DeleteImagePipeline"
			request := map[string]interface{}{
				"ImagePipelineId": item["ImagePipelineId"],
				"RegionId":        aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				id, _ := item["ImagePipelineId"].(string)
				log.Printf("[ERROR] Failed to delete Ecs Image Pipeline (%s): %s", name, err)
				deleteErrors = append(deleteErrors, fmt.Sprintf("%s (%s): %s", name, id, err))
				continue
			}
			log.Printf("[INFO] Delete Ecs Image Pipeline success: %s ", name)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if len(deleteErrors) > 0 {
		return fmt.Errorf("failed to delete %d Ecs Image Pipeline(s): %s", len(deleteErrors), strings.Join(deleteErrors, "; "))
	}
	return nil
}

// The documentation example deliberately uses a generic name, so the sweeper
// only treats an exactly named pipeline as its own when the example's ownership
// tag is present. Acceptance-test fixtures keep their name prefixes.
const (
	imagePipelineOwnerTagKey   = "terraform-example"
	imagePipelineOwnerTagValue = "true"
)

// imagePipelineTagMap decodes the Tags member into a plain map. A malformed
// member yields an empty map, which means "not owned" and never matches.
func imagePipelineTagMap(tags interface{}) map[string]string {
	result := map[string]string{}
	container, ok := tags.(map[string]interface{})
	if !ok {
		return result
	}
	raw, ok := container["Tag"]
	if !ok {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := item["TagKey"].(string)
		value, _ := item["TagValue"].(string)
		if key != "" {
			result[key] = value
		}
	}
	return result
}

// shouldSweepEcsImagePipeline reports whether a pipeline belongs to a fixture
// this sweeper owns. Acceptance-test fixtures are matched by name prefix. The
// documentation example's generic name is matched only together with its
// ownership tag, so an unrelated pipeline with the same name is left alone.
func shouldSweepEcsImagePipeline(name string, tags map[string]string) bool {
	lower := strings.ToLower(name)
	for _, prefix := range []string{"tf-testacc", "tf_testacc"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	if lower == "terraform-example" {
		return tags[imagePipelineOwnerTagKey] == imagePipelineOwnerTagValue
	}
	return false
}

func TestUnitShouldSweepEcsImagePipeline(t *testing.T) {
	owned := map[string]string{imagePipelineOwnerTagKey: imagePipelineOwnerTagValue}
	cases := []struct {
		name string
		tags map[string]string
		want bool
	}{
		{"tf-testAccimagepipeline1234", nil, true},
		{"tf_testAccimage", nil, true},
		{"terraform-example", owned, true},
		{"terraform-example", nil, false},
		{"terraform-example", map[string]string{imagePipelineOwnerTagKey: "false"}, false},
		{"terraform-example-2", owned, false},
		{"terraform-examples-backup", owned, false},
		{"production-image-pipeline", owned, false},
	}
	for _, testCase := range cases {
		if got := shouldSweepEcsImagePipeline(testCase.name, testCase.tags); got != testCase.want {
			t.Errorf("shouldSweepEcsImagePipeline(%q, %v) = %v, want %v", testCase.name, testCase.tags, got, testCase.want)
		}
	}
}

func TestUnitImagePipelineTagMap(t *testing.T) {
	tags := imagePipelineTagMap(map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": imagePipelineOwnerTagKey, "TagValue": imagePipelineOwnerTagValue},
			"malformed",
			map[string]interface{}{"TagValue": "no-key"},
		},
	})
	if tags[imagePipelineOwnerTagKey] != imagePipelineOwnerTagValue {
		t.Fatalf("expected the ownership tag, got %v", tags)
	}
	if len(tags) != 1 {
		t.Fatalf("expected malformed entries to be skipped, got %v", tags)
	}
	if len(imagePipelineTagMap("not-an-object")) != 0 {
		t.Fatal("expected an empty map for a malformed Tags member")
	}
}
