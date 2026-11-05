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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	return sweepEcsImagePipelines(region, aliyunClient.RegionId,
		func(request map[string]interface{}) (map[string]interface{}, error) {
			return aliyunClient.RpcPost("Ecs", "2014-05-26", "DescribeImagePipelines", nil, request, true)
		},
		func(regionId, id string) error {
			_, err := aliyunClient.RpcPost("Ecs", "2014-05-26", "DeleteImagePipeline", nil, map[string]interface{}{
				"ImagePipelineId": id,
				"RegionId":        regionId,
			}, false)
			return err
		})
}

// sweepEcsImagePipelines is the sweep itself: it follows every page, deletes the
// pipelines the sweeper owns and reports every problem it met. The transport is
// passed in so the orchestration can be exercised without a cloud client.
func sweepEcsImagePipelines(region, regionId string, list func(request map[string]interface{}) (map[string]interface{}, error), deletePipeline func(regionId, id string) error) error {
	action := "DescribeImagePipelines"
	request := map[string]interface{}{}
	request["RegionId"] = regionId
	request["MaxResults"] = tfhelper.PageSizeLarge

	deleteErrors := make([]string, 0)
	listErrors := make([]string, 0)
	recordErrors := make([]string, 0)
	seenTokens := map[string]struct{}{}
	pages := 0

	for {
		// The token guard below catches an echoed cursor; this cap catches a
		// server that keeps returning distinct tokens, so the listing always
		// terminates like the data source helpers do.
		if imagePipelinePageBudgetExceeded(pages) {
			listErrors = append(listErrors, fmt.Sprintf("%s kept returning a continuation token for %d pages; the sweep stopped instead of listing forever", action, pages))
			break
		}
		pages++
		var response map[string]interface{}
		wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
		listErr := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var err error
			response, err = list(request)
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

		result, pageRecordErrors, err := imagePipelineRecords(response)
		if err != nil {
			log.Printf("[ERROR] Reading the image pipeline list failed!!! Body: %v.", response)
			listErrors = append(listErrors, fmt.Sprintf("%s returned an unexpected image pipeline list: %s", action, err))
			break
		}
		// A malformed member is reported on its own: the records that could be
		// read are still swept, so one bad member cannot keep an owned pipeline
		// from being deleted.
		recordErrors = append(recordErrors, pageRecordErrors...)
		outcome := sweepImagePipelineRecords(result, func(name, id string) error {
			err := deletePipeline(regionId, id)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Image Pipeline (%s): %s", name, err)
			} else {
				log.Printf("[INFO] Delete Ecs Image Pipeline success: %s ", name)
			}
			return err
		})
		deleteErrors = append(deleteErrors, outcome.deleteErrors...)
		recordErrors = append(recordErrors, outcome.recordErrors...)
		nextToken, tokenErr := imagePipelineNextToken(response, seenTokens)
		if tokenErr != nil {
			listErrors = append(listErrors, fmt.Sprintf("%s pagination failed: %s", action, tokenErr))
			break
		}
		if nextToken == "" {
			break
		}
		request["NextToken"] = nextToken
	}
	return imagePipelineSweepError(region, deleteErrors, recordErrors, listErrors)
}

// imagePipelinePageBudgetExceeded reports whether the sweep has already read
// the pages one listing may use and must stop instead of following another
// token.
func imagePipelinePageBudgetExceeded(pages int) bool {
	return pages >= imagePipelineSweepMaxPages
}

// imagePipelineNextToken returns the continuation token the next request must
// send, or "" when the page was the last one. A token that was already followed
// is reported: the same page would be requested again, forever. A member of the
// wrong type is reported for the same reason.
func imagePipelineNextToken(response map[string]interface{}, seen map[string]struct{}) (string, error) {
	raw, ok := response["NextToken"]
	if !ok || raw == nil {
		return "", nil
	}
	token, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("NextToken is %T instead of a string", raw)
	}
	if token == "" {
		return "", nil
	}
	if _, repeated := seen[token]; repeated {
		return "", fmt.Errorf("the server repeated the pagination token %q", token)
	}
	seen[token] = struct{}{}
	return token, nil
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
// matches. A malformed container or list is reported as an error, which stops
// the listing. A single malformed member is reported on its own and the readable
// records are still returned, so one bad entry cannot make the whole page look
// empty and leave every owned pipeline behind.
func imagePipelineRecords(response map[string]interface{}) ([]map[string]interface{}, []string, error) {
	rawContainer, ok := response["ImagePipeline"]
	if !ok || rawContainer == nil {
		return nil, nil, nil
	}
	container, ok := rawContainer.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("ImagePipeline is %T instead of an object", rawContainer)
	}
	rawList, ok := container["ImagePipelineSet"]
	if !ok || rawList == nil {
		return nil, nil, nil
	}
	entries, ok := rawList.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("ImagePipelineSet is %T instead of a list", rawList)
	}
	result := make([]map[string]interface{}, 0, len(entries))
	recordErrors := make([]string, 0)
	for index, entry := range entries {
		item, ok := entry.(map[string]interface{})
		if !ok {
			recordErrors = append(recordErrors, fmt.Sprintf("ImagePipelineSet[%d] is %T instead of an object", index, entry))
			continue
		}
		result = append(result, item)
	}
	return result, recordErrors, nil
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

// imagePipelineSweepMaxPages bounds one sweeper listing, matching the page
// budget the read helpers use.
const imagePipelineSweepMaxPages = 1000

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

// imagePipelineSweepOutcome collects the decisions of one page: the pipelines
// that were deleted, the deletions that failed and the records whose identity or
// ownership could not be read.
type imagePipelineSweepOutcome struct {
	deleted      []string
	deleteErrors []string
	recordErrors []string
}

// sweepImagePipelineRecords applies the ownership rule to the records of one
// page and deletes every pipeline it owns. Ownership is decided from the name
// first: a fixture name prefix identifies the pipeline on its own, so an
// unrelated malformed tag member cannot keep an owned pipeline from being
// deleted. Only the documentation example's generic name needs the ownership
// tag, and a record whose ownership cannot be established there is left alone
// and reported. A deletion that fails is reported and the remaining records are
// still processed, so one failure cannot hide another.
func sweepImagePipelineRecords(records []map[string]interface{}, deletePipeline func(name, id string) error) imagePipelineSweepOutcome {
	outcome := imagePipelineSweepOutcome{}
	for index, item := range records {
		name, err := imagePipelineRecordName(item)
		if err != nil {
			outcome.recordErrors = append(outcome.recordErrors, fmt.Sprintf("DescribeImagePipelines record %d: %s", index, err))
			continue
		}
		owned, err := imagePipelineOwnedBySweeper(name, item["Tags"])
		if err != nil {
			outcome.recordErrors = append(outcome.recordErrors, fmt.Sprintf("image pipeline %s: %s", name, err))
			continue
		}
		if !owned {
			log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", name)
			continue
		}
		id, err := imagePipelineRecordId(item)
		if err != nil {
			outcome.recordErrors = append(outcome.recordErrors, fmt.Sprintf("image pipeline %s: %s", name, err))
			continue
		}
		if err := deletePipeline(name, id); err != nil {
			outcome.deleteErrors = append(outcome.deleteErrors, fmt.Sprintf("%s (%s): %s", name, id, err))
			continue
		}
		outcome.deleted = append(outcome.deleted, id)
	}
	return outcome
}

// imagePipelineOwnedBySweeper reports whether a pipeline belongs to a fixture
// this sweeper owns. Acceptance-test fixtures are matched by name prefix, which
// identifies them without looking at their tags. The documentation example's
// generic name is matched only together with its ownership tag, so an unrelated
// pipeline with the same name is left alone and a malformed tag member makes the
// ownership unknown instead of owned.
func imagePipelineOwnedBySweeper(name string, tags interface{}) (bool, error) {
	lower := strings.ToLower(name)
	for _, prefix := range []string{"tf-testacc", "tf_testacc"} {
		if strings.HasPrefix(lower, prefix) {
			return true, nil
		}
	}
	if lower != "terraform-example" {
		return false, nil
	}
	tagMap, err := imagePipelineTagMap(tags)
	if err != nil {
		return false, err
	}
	return tagMap[imagePipelineOwnerTagKey] == imagePipelineOwnerTagValue, nil
}

func TestUnitImagePipelineOwnedBySweeper(t *testing.T) {
	ownedTags := map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": imagePipelineOwnerTagKey, "TagValue": imagePipelineOwnerTagValue},
		},
	}
	notOwnedTags := map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": imagePipelineOwnerTagKey, "TagValue": "false"},
		},
	}
	cases := []struct {
		name string
		tags interface{}
		want bool
	}{
		{"tf-testAccimagepipeline1234", nil, true},
		{"tf_testAccimage", "not-an-object", true},
		{"terraform-example", ownedTags, true},
		{"terraform-example", nil, false},
		{"terraform-example", notOwnedTags, false},
		{"terraform-example-2", ownedTags, false},
		{"terraform-examples-backup", ownedTags, false},
		{"production-image-pipeline", ownedTags, false},
	}
	for _, testCase := range cases {
		got, err := imagePipelineOwnedBySweeper(testCase.name, testCase.tags)
		if err != nil {
			t.Errorf("imagePipelineOwnedBySweeper(%q, %v) returned %v", testCase.name, testCase.tags, err)
			continue
		}
		if got != testCase.want {
			t.Errorf("imagePipelineOwnedBySweeper(%q, %v) = %v, want %v", testCase.name, testCase.tags, got, testCase.want)
		}
	}

	// The generic name depends on its tags, so unreadable tags leave the
	// ownership unknown instead of guessing either way.
	if _, err := imagePipelineOwnedBySweeper("terraform-example", "not-an-object"); err == nil {
		t.Fatal("unreadable ownership tags must be reported instead of deciding ownership")
	}
}

func TestUnitImagePipelinePageBudget(t *testing.T) {
	for pages := 0; pages < imagePipelineSweepMaxPages; pages++ {
		if imagePipelinePageBudgetExceeded(pages) {
			t.Fatalf("page %d of the %d page budget must be allowed", pages, imagePipelineSweepMaxPages)
		}
	}
	if !imagePipelinePageBudgetExceeded(imagePipelineSweepMaxPages) {
		t.Fatalf("page %d must end the listing", imagePipelineSweepMaxPages)
	}
	if !imagePipelinePageBudgetExceeded(imagePipelineSweepMaxPages + 1) {
		t.Fatal("a page beyond the budget must end the listing")
	}
}

func TestUnitImagePipelineNextToken(t *testing.T) {
	for name, response := range map[string]map[string]interface{}{
		"missing token": {},
		"nil token":     {"NextToken": nil},
		"empty token":   {"NextToken": ""},
	} {
		token, err := imagePipelineNextToken(response, map[string]struct{}{})
		if err != nil || token != "" {
			t.Errorf("%s must end the sweep, got %q, %v", name, token, err)
		}
	}

	seen := map[string]struct{}{}
	token, err := imagePipelineNextToken(map[string]interface{}{"NextToken": "page-2"}, seen)
	if err != nil || token != "page-2" {
		t.Fatalf("a new token must be followed, got %q, %v", token, err)
	}
	if _, err := imagePipelineNextToken(map[string]interface{}{"NextToken": "page-2"}, seen); err == nil {
		t.Fatal("a repeated token must be reported instead of requesting the same page again")
	}
	if _, err := imagePipelineNextToken(map[string]interface{}{"NextToken": 2}, map[string]struct{}{}); err == nil {
		t.Fatal("a non-string token must be reported")
	}
}

// A pipeline identified by name is deleted even when it carries an unrelated
// malformed tag, an object whose ownership cannot be read is protected and
// reported, and one failed deletion does not hide the others.
func TestUnitSweepImagePipelineRecords(t *testing.T) {
	t.Run("name identified fixture ignores unrelated malformed tags", func(t *testing.T) {
		deleted := make([]string, 0)
		outcome := sweepImagePipelineRecords([]map[string]interface{}{
			{"ImagePipelineId": "ip-fixture", "Name": "tf-testAcc-image", "Tags": "not-an-object"},
		}, func(name, id string) error {
			deleted = append(deleted, id)
			return nil
		})
		assert.Equal(t, []string{"ip-fixture"}, deleted, "the fixture must be deleted")
		assert.Empty(t, outcome.deleteErrors)
		assert.Empty(t, outcome.recordErrors)
	})

	t.Run("unknown ownership is protected and reported", func(t *testing.T) {
		deleted := make([]string, 0)
		outcome := sweepImagePipelineRecords([]map[string]interface{}{
			{"ImagePipelineId": "ip-generic", "Name": "terraform-example", "Tags": "not-an-object"},
		}, func(name, id string) error {
			deleted = append(deleted, id)
			return nil
		})
		assert.Empty(t, deleted, "an object whose ownership is unknown must not be deleted")
		require.Len(t, outcome.recordErrors, 1)
		assert.Contains(t, outcome.recordErrors[0], "terraform-example")
	})

	t.Run("unrelated pipeline is left alone", func(t *testing.T) {
		deleted := make([]string, 0)
		outcome := sweepImagePipelineRecords([]map[string]interface{}{
			{"ImagePipelineId": "ip-unrelated", "Name": "production-image-pipeline"},
		}, func(name, id string) error {
			deleted = append(deleted, id)
			return nil
		})
		assert.Empty(t, deleted)
		assert.Empty(t, outcome.recordErrors)
		assert.Empty(t, outcome.deleteErrors)
	})

	t.Run("independent deletion failures are preserved", func(t *testing.T) {
		deleted := make([]string, 0)
		outcome := sweepImagePipelineRecords([]map[string]interface{}{
			{"ImagePipelineId": "ip-1", "Name": "tf-testAcc-1"},
			{"ImagePipelineId": "ip-2", "Name": "tf-testAcc-2"},
		}, func(name, id string) error {
			if id == "ip-1" {
				return fmt.Errorf("Forbidden")
			}
			deleted = append(deleted, id)
			return nil
		})
		assert.Equal(t, []string{"ip-2"}, deleted, "the remaining records must still be deleted")
		require.Len(t, outcome.deleteErrors, 1)
		assert.Contains(t, outcome.deleteErrors[0], "ip-1")
		assert.Contains(t, outcome.deleteErrors[0], "Forbidden")
		assert.Empty(t, outcome.recordErrors)
	})
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
	records, recordErrors, err := imagePipelineRecords(map[string]interface{}{
		"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{
				map[string]interface{}{"ImagePipelineId": "ip-1", "Name": "tf-testAcc-1"},
			},
		},
	})
	if err != nil || len(recordErrors) != 0 {
		t.Fatalf("a well formed list must decode, got %v, %v", recordErrors, err)
	}
	if len(records) != 1 || records[0]["ImagePipelineId"] != "ip-1" {
		t.Fatalf("expected the pipeline record, got %v", records)
	}

	records, recordErrors, err = imagePipelineRecords(map[string]interface{}{
		"ImagePipeline": map[string]interface{}{"ImagePipelineSet": []interface{}{}},
	})
	if err != nil || len(recordErrors) != 0 || len(records) != 0 {
		t.Fatalf("an explicit empty list is a valid empty page, got %v, %v, %v", records, recordErrors, err)
	}

	// The service omits the list when nothing matches, which is a valid empty
	// page and must not fail the sweep.
	for name, response := range map[string]map[string]interface{}{
		"missing container": {},
		"missing list":      {"ImagePipeline": map[string]interface{}{}},
		"nil list":          {"ImagePipeline": map[string]interface{}{"ImagePipelineSet": nil}},
	} {
		records, recordErrors, err := imagePipelineRecords(response)
		if err != nil || len(recordErrors) != 0 || len(records) != 0 {
			t.Errorf("%s must be a valid empty page, got %v, %v, %v", name, records, recordErrors, err)
		}
	}

	for name, response := range map[string]map[string]interface{}{
		"container wrong type": {"ImagePipeline": "not-an-object"},
		"list wrong type":      {"ImagePipeline": map[string]interface{}{"ImagePipelineSet": "not-a-list"}},
	} {
		if _, _, err := imagePipelineRecords(response); err == nil {
			t.Errorf("%s must be reported instead of being treated as an empty list", name)
		}
	}

	// A single malformed member is isolated: the readable records are returned
	// and the bad entry is reported, so the page is not thrown away as a whole.
	records, recordErrors, err = imagePipelineRecords(map[string]interface{}{
		"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{
				map[string]interface{}{"ImagePipelineId": "ip-1", "Name": "tf-testAcc-1"},
				"not-an-object",
				map[string]interface{}{"ImagePipelineId": "ip-2", "Name": "tf-testAcc-2"},
			},
		},
	})
	if err != nil {
		t.Fatalf("a mixed page must not fail as a whole, got %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("the readable records must be kept, got %v", records)
	}
	require.Len(t, recordErrors, 1)
	assert.Contains(t, recordErrors[0], "ImagePipelineSet[1]")
}

// The sweep orchestration keeps the readable records of a mixed page, follows
// the next page, deletes every owned pipeline and reports the malformed member
// and the independent deletion failure in one summary.
func TestUnitSweepEcsImagePipelinesMixedPages(t *testing.T) {
	deleted := make([]string, 0)
	listCalls := 0
	list := func(request map[string]interface{}) (map[string]interface{}, error) {
		listCalls++
		if token, _ := request["NextToken"].(string); token == "page-2" {
			return map[string]interface{}{
				"ImagePipeline": map[string]interface{}{
					"ImagePipelineSet": []interface{}{
						map[string]interface{}{
							"ImagePipelineId": "ip-page-2",
							"Name":            "terraform-example",
							"Tags": map[string]interface{}{
								"Tag": []interface{}{
									map[string]interface{}{"TagKey": imagePipelineOwnerTagKey, "TagValue": imagePipelineOwnerTagValue},
								},
							},
						},
					},
				},
			}, nil
		}
		return map[string]interface{}{
			"NextToken": "page-2",
			"ImagePipeline": map[string]interface{}{
				"ImagePipelineSet": []interface{}{
					map[string]interface{}{"ImagePipelineId": "ip-page-1", "Name": "tf-testAcc-1"},
					"not-an-object",
					map[string]interface{}{"ImagePipelineId": "ip-fails", "Name": "tf-testAcc-2"},
				},
			},
		}, nil
	}
	deletePipeline := func(regionId, id string) error {
		if id == "ip-fails" {
			return fmt.Errorf("Forbidden")
		}
		deleted = append(deleted, id)
		return nil
	}

	err := sweepEcsImagePipelines("cn-hangzhou", "cn-hangzhou", list, deletePipeline)
	require.Error(t, err, "the sweep must report the failures it met")
	text := err.Error()
	assert.Equal(t, 2, listCalls, "the sweep must follow the continuation token")
	assert.ElementsMatch(t, []string{"ip-page-1", "ip-page-2"}, deleted, "both readable owned pipelines must be deleted")
	assert.Contains(t, text, "ImagePipelineSet[1]", "the malformed member must be reported")
	assert.Contains(t, text, "ip-fails", "the independent deletion failure must be reported")
	assert.Contains(t, text, "Forbidden")
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
