package sweep_test

// Legacy product sweepers run in this dedicated test binary.

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

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
	listErrors := make([]string, 0)

	for {
		var response map[string]interface{}
		wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
		listErr := retry.Retry(1*time.Minute, func() *retry.RetryError {
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
		if listErr != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, listErr)
			// The pipelines of the pages already read are still deleted below;
			// a listing failure must not hide those deletion results.
			listErrors = append(listErrors, fmt.Sprintf("%s failed to list image pipelines: %s", action, listErr))
			break
		}

		result, err := imagePipelineRecords(response)
		if err != nil {
			log.Printf("[ERROR] Reading the image pipeline list failed!!! Body: %v.", response)
			listErrors = append(listErrors, fmt.Sprintf("%s returned an unexpected image pipeline list: %s", action, err))
			break
		}
		for _, item := range result {
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
	return imagePipelineSweepError(region, deleteErrors, listErrors)
}

// imagePipelineRecords decodes the pipeline list of one page. A response without
// the list is a valid empty page, because the service omits it when nothing
// matches. A member or entry of another shape is a malformed response and is
// reported instead of being dropped, so a broken response can never look like
// "nothing left to sweep".
func imagePipelineRecords(response map[string]interface{}) ([]map[string]interface{}, error) {
	rawContainer, ok := response["ImagePipeline"]
	if !ok || rawContainer == nil {
		return nil, nil
	}
	container, ok := rawContainer.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("ImagePipeline is %T instead of an object", rawContainer)
	}
	rawList, ok := container["ImagePipelineSet"]
	if !ok || rawList == nil {
		return nil, nil
	}
	entries, ok := rawList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ImagePipelineSet is %T instead of a list", rawList)
	}
	result := make([]map[string]interface{}, 0, len(entries))
	for index, entry := range entries {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ImagePipelineSet[%d] is %T instead of an object", index, entry)
		}
		result = append(result, item)
	}
	return result, nil
}

// imagePipelineSweepError reports every problem the sweep met. Deletion failures
// are always part of the message, so a listing error that stops later pages
// cannot hide an object that was already left behind.
func imagePipelineSweepError(region string, deleteErrors, listErrors []string) error {
	parts := make([]string, 0, len(deleteErrors)+len(listErrors))
	if len(deleteErrors) > 0 {
		parts = append(parts, fmt.Sprintf("failed to delete %d Ecs Image Pipeline(s): %s", len(deleteErrors), strings.Join(deleteErrors, "; ")))
	}
	if len(listErrors) > 0 {
		parts = append(parts, fmt.Sprintf("failed to list Ecs Image Pipelines in %s: %s", region, strings.Join(listErrors, "; ")))
	}
	if len(parts) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(parts, "; "))
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

func TestUnitImagePipelineRecords(t *testing.T) {
	records, err := imagePipelineRecords(map[string]interface{}{
		"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{
				map[string]interface{}{"ImagePipelineId": "ip-1", "Name": "tf-testAcc-1"},
			},
		},
	})
	if err != nil {
		t.Fatalf("a well formed list must decode, got %v", err)
	}
	if len(records) != 1 || records[0]["ImagePipelineId"] != "ip-1" {
		t.Fatalf("expected the pipeline record, got %v", records)
	}

	records, err = imagePipelineRecords(map[string]interface{}{
		"ImagePipeline": map[string]interface{}{"ImagePipelineSet": []interface{}{}},
	})
	if err != nil || len(records) != 0 {
		t.Fatalf("an explicit empty list is a valid empty page, got %v, %v", records, err)
	}

	// The service omits the list when nothing matches, which is a valid empty
	// page and must not fail the sweep.
	for name, response := range map[string]map[string]interface{}{
		"missing container": {},
		"missing list":      {"ImagePipeline": map[string]interface{}{}},
		"nil list":          {"ImagePipeline": map[string]interface{}{"ImagePipelineSet": nil}},
	} {
		records, err := imagePipelineRecords(response)
		if err != nil || len(records) != 0 {
			t.Errorf("%s must be a valid empty page, got %v, %v", name, records, err)
		}
	}

	for name, response := range map[string]map[string]interface{}{
		"container wrong type": {"ImagePipeline": "not-an-object"},
		"list wrong type":      {"ImagePipeline": map[string]interface{}{"ImagePipelineSet": "not-a-list"}},
		"entry wrong type": {"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{"not-an-object"},
		}},
	} {
		if _, err := imagePipelineRecords(response); err == nil {
			t.Errorf("%s must be reported instead of being treated as an empty list", name)
		}
	}
}

// A listing error must never hide the objects that failed to delete earlier in
// the sweep: both failures are reported together.
func TestUnitImagePipelineSweepError(t *testing.T) {
	if err := imagePipelineSweepError("cn-hangzhou", nil, nil); err != nil {
		t.Fatalf("no failures must not produce an error, got %v", err)
	}

	deleteErrors := []string{"terraform-example (ip-1): Forbidden"}
	if err := imagePipelineSweepError("cn-hangzhou", deleteErrors, nil); err == nil || !strings.Contains(err.Error(), "failed to delete 1") {
		t.Fatalf("a deletion failure must be reported, got %v", err)
	}

	listErrors := []string{"DescribeImagePipelines failed to list image pipelines: timeout"}
	if err := imagePipelineSweepError("cn-hangzhou", nil, listErrors); err == nil || !strings.Contains(err.Error(), "failed to list") {
		t.Fatalf("a listing failure must be reported, got %v", err)
	}

	combined := imagePipelineSweepError("cn-hangzhou", deleteErrors, listErrors)
	if combined == nil {
		t.Fatal("a combined failure must produce an error")
	}
	for _, expected := range []string{"failed to delete 1", "terraform-example (ip-1): Forbidden", "failed to list"} {
		if !strings.Contains(combined.Error(), expected) {
			t.Fatalf("the combined error must keep %q, got %v", expected, combined)
		}
	}
}
