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
	recordErrors := make([]string, 0)

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
		for index, item := range result {
			// A record whose identity or ownership cannot be read is not
			// silently skipped: skipping it would leave the pipeline behind
			// while the sweep still reports success, so it is reported next to
			// the listing and deletion failures.
			name, nameErr := imagePipelineRecordName(item)
			if nameErr != nil {
				recordErrors = append(recordErrors, fmt.Sprintf("%s record %d: %s", action, index, nameErr))
				continue
			}
			tags, tagsErr := imagePipelineTagMap(item["Tags"])
			if tagsErr != nil {
				recordErrors = append(recordErrors, fmt.Sprintf("image pipeline %s: %s", name, tagsErr))
				continue
			}
			if !shouldSweepEcsImagePipeline(name, tags) {
				log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", name)
				continue
			}
			id, idErr := imagePipelineRecordId(item)
			if idErr != nil {
				recordErrors = append(recordErrors, fmt.Sprintf("image pipeline %s: %s", name, idErr))
				continue
			}
			action := "DeleteImagePipeline"
			request := map[string]interface{}{
				"ImagePipelineId": id,
				"RegionId":        aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
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
	return imagePipelineSweepError(region, deleteErrors, recordErrors, listErrors)
}

// imagePipelineRecordName returns the name the sweep logs and matches on. A
// missing or non-string name cannot be matched against the owned fixture
// patterns, so it is reported instead of being read as an empty name.
func imagePipelineRecordName(item map[string]interface{}) (string, error) {
	raw, ok := item["Name"]
	if !ok || raw == nil {
		return "", fmt.Errorf("record has no Name")
	}
	name, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("Name is %T instead of a string", raw)
	}
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("Name is blank")
	}
	return name, nil
}

// imagePipelineRecordId returns the identity DeleteImagePipeline needs. A
// missing or non-string identity cannot be deleted and must not be dropped.
func imagePipelineRecordId(item map[string]interface{}) (string, error) {
	raw, ok := item["ImagePipelineId"]
	if !ok || raw == nil {
		return "", fmt.Errorf("record has no ImagePipelineId")
	}
	id, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("ImagePipelineId is %T instead of a string", raw)
	}
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("ImagePipelineId is blank")
	}
	return id, nil
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
// and unreadable records are always part of the message, so a listing error that
// stops later pages cannot hide an object that was already left behind, and a
// record whose ownership could not be read cannot be mistaken for a clean sweep.
func imagePipelineSweepError(region string, deleteErrors, recordErrors, listErrors []string) error {
	parts := make([]string, 0, len(deleteErrors)+len(recordErrors)+len(listErrors))
	if len(deleteErrors) > 0 {
		parts = append(parts, fmt.Sprintf("failed to delete %d Ecs Image Pipeline(s): %s", len(deleteErrors), strings.Join(deleteErrors, "; ")))
	}
	if len(recordErrors) > 0 {
		parts = append(parts, fmt.Sprintf("could not read %d Ecs Image Pipeline record(s): %s", len(recordErrors), strings.Join(recordErrors, "; ")))
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

// imagePipelineTagMap decodes the Tags member into a plain map. An absent
// member means "no tags", which is a legitimate "not owned" answer; a member of
// the wrong shape is reported instead, because reading it as "no tags" would
// silently leave an owned pipeline behind. Entries without a string key are
// reported for the same reason.
func imagePipelineTagMap(tags interface{}) (map[string]string, error) {
	result := map[string]string{}
	if tags == nil {
		return result, nil
	}
	container, ok := tags.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Tags is %T instead of an object", tags)
	}
	raw, ok := container["Tag"]
	if !ok || raw == nil {
		return result, nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Tags.Tag is %T instead of a list", raw)
	}
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Tags.Tag[%d] is %T instead of an object", index, entry)
		}
		key, ok := item["TagKey"].(string)
		if !ok || key == "" {
			return nil, fmt.Errorf("Tags.Tag[%d] has no string key", index)
		}
		value := ""
		if item["TagValue"] != nil {
			text, ok := item["TagValue"].(string)
			if !ok {
				return nil, fmt.Errorf("Tags.Tag[%d].TagValue is %T instead of a string", index, item["TagValue"])
			}
			value = text
		}
		result[key] = value
	}
	return result, nil
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
	tags, err := imagePipelineTagMap(map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": imagePipelineOwnerTagKey, "TagValue": imagePipelineOwnerTagValue},
			map[string]interface{}{"TagKey": "no-value"},
		},
	})
	if err != nil {
		t.Fatalf("a well formed tag list must decode, got %v", err)
	}
	if tags[imagePipelineOwnerTagKey] != imagePipelineOwnerTagValue {
		t.Fatalf("expected the ownership tag, got %v", tags)
	}
	if len(tags) != 2 {
		t.Fatalf("expected both tags, got %v", tags)
	}

	// An absent member means the pipeline has no tags, which is a valid "not
	// owned" answer rather than a malformed record.
	for name, value := range map[string]interface{}{
		"nil member":     nil,
		"missing Tag":    map[string]interface{}{},
		"nil Tag":        map[string]interface{}{"Tag": nil},
		"empty Tag list": map[string]interface{}{"Tag": []interface{}{}},
	} {
		tags, err := imagePipelineTagMap(value)
		if err != nil || len(tags) != 0 {
			t.Errorf("%s must be an empty tag set, got %v, %v", name, tags, err)
		}
	}

	// A malformed ownership member must be reported: reading it as "no tags"
	// would skip an owned pipeline while the sweep reports success.
	for name, value := range map[string]interface{}{
		"Tags wrong type":   "not-an-object",
		"Tag wrong type":    map[string]interface{}{"Tag": "not-a-list"},
		"entry wrong type":  map[string]interface{}{"Tag": []interface{}{"malformed"}},
		"entry without key": map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagValue": "no-key"}}},
		"blank key":         map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": ""}}},
		"non-string key":    map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": 1}}},
		"non-string value":  map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": "k", "TagValue": 1}}},
	} {
		if _, err := imagePipelineTagMap(value); err == nil {
			t.Errorf("%s must be reported instead of being read as an empty tag set", name)
		}
	}
}

// A record whose name or identity cannot be read must be reported: skipping it
// silently would leave the pipeline behind while the sweep reports success.
func TestUnitImagePipelineRecordIdentity(t *testing.T) {
	if name, err := imagePipelineRecordName(map[string]interface{}{"Name": "terraform-example"}); err != nil || name != "terraform-example" {
		t.Fatalf("a well formed name must decode, got %q, %v", name, err)
	}
	if id, err := imagePipelineRecordId(map[string]interface{}{"ImagePipelineId": "ip-1"}); err != nil || id != "ip-1" {
		t.Fatalf("a well formed identity must decode, got %q, %v", id, err)
	}

	for name, item := range map[string]map[string]interface{}{
		"missing name":        {"ImagePipelineId": "ip-1"},
		"non-string name":     {"Name": 7, "ImagePipelineId": "ip-1"},
		"blank name":          {"Name": "  ", "ImagePipelineId": "ip-1"},
		"missing identity":    {"Name": "tf-testAcc-1"},
		"non-string identity": {"Name": "tf-testAcc-1", "ImagePipelineId": 7},
		"blank identity":      {"Name": "tf-testAcc-1", "ImagePipelineId": " "},
	} {
		_, nameErr := imagePipelineRecordName(item)
		_, idErr := imagePipelineRecordId(item)
		if nameErr == nil && idErr == nil {
			t.Errorf("%s must be reported as unreadable", name)
		}
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
// the sweep, and an unreadable record must never look like a clean sweep: every
// failure is reported together.
func TestUnitImagePipelineSweepError(t *testing.T) {
	if err := imagePipelineSweepError("cn-hangzhou", nil, nil, nil); err != nil {
		t.Fatalf("no failures must not produce an error, got %v", err)
	}

	deleteErrors := []string{"terraform-example (ip-1): Forbidden"}
	if err := imagePipelineSweepError("cn-hangzhou", deleteErrors, nil, nil); err == nil || !strings.Contains(err.Error(), "failed to delete 1") {
		t.Fatalf("a deletion failure must be reported, got %v", err)
	}

	recordErrors := []string{"DescribeImagePipelines record 2: Name is blank"}
	if err := imagePipelineSweepError("cn-hangzhou", nil, recordErrors, nil); err == nil || !strings.Contains(err.Error(), "could not read 1") {
		t.Fatalf("an unreadable record must be reported, got %v", err)
	}

	listErrors := []string{"DescribeImagePipelines failed to list image pipelines: timeout"}
	if err := imagePipelineSweepError("cn-hangzhou", nil, nil, listErrors); err == nil || !strings.Contains(err.Error(), "failed to list") {
		t.Fatalf("a listing failure must be reported, got %v", err)
	}

	combined := imagePipelineSweepError("cn-hangzhou", deleteErrors, recordErrors, listErrors)
	if combined == nil {
		t.Fatal("a combined failure must produce an error")
	}
	for _, expected := range []string{"failed to delete 1", "terraform-example (ip-1): Forbidden", "could not read 1", "Name is blank", "failed to list"} {
		if !strings.Contains(combined.Error(), expected) {
			t.Fatalf("the combined error must keep %q, got %v", expected, combined)
		}
	}
}
