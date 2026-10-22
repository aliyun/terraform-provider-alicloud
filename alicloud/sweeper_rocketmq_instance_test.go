package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_rocketmq_instance",
		&resource.Sweeper{
			Name: "alicloud_rocketmq_instance",
			F:    testSweepRocketmqInstance,
		})
}

var sweepRocketmqInstanceNamePrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"testAcc",
}

// sweepRocketmqInstanceMatchesName reports whether the instance name carries one
// of the AccTest prefixes the sweeper treats as its own.
func sweepRocketmqInstanceMatchesName(name string) bool {
	for _, prefix := range sweepRocketmqInstanceNamePrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// sweepRocketmqInstanceSkip reports whether the sweeper has to spare the listed
// instance. The name selection applies unless SWEEP_ALL is set, and the shared
// age protection applies in both name modes. ListInstances returns the creation
// time under the lowercase createTime key as a "2006-01-02 15:04:05" string,
// which the shared contract does not look up, so it is mapped onto
// CreationTime first; the shared parser reads that layout as UTC, which can
// only make an age look smaller than it is and therefore keeps the protection
// conservative.
func sweepRocketmqInstanceSkip(item map[string]interface{}) (bool, string) {
	name, _ := item["instanceName"].(string)
	if !sweep.SweepAll() && !sweepRocketmqInstanceMatchesName(name) {
		return true, "name does not carry an AccTest prefix"
	}
	ageItem := make(map[string]interface{}, 1)
	if v, ok := item["createTime"]; ok && v != nil {
		ageItem["CreationTime"] = v
	}
	return sweep.SweepRecentCheck(ageItem)
}

// sweepRocketmqInstances applies the name selection and the shared age
// protection to a ListInstances payload before anything is removed. remove is
// only invoked for the instances the sweeper is allowed to delete.
func sweepRocketmqInstances(items []interface{}, remove func(item map[string]interface{}) error) {
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := item["instanceName"].(string)
		if skip, reason := sweepRocketmqInstanceSkip(item); skip {
			log.Printf("[INFO] Skipping rocketmq instance %s: %s ", name, reason)
			continue
		}
		log.Printf("[INFO] delete rocketmq instance: %s ", name)
		if err := remove(item); err != nil {
			log.Printf("[ERROR] Failed to delete rocketmq instance (%s): %s", name, err)
		}
	}
}

func testSweepRocketmqInstance(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "/instances"
	query := make(map[string]*string)
	query["pageNumber"] = tea.String("1")
	query["pageSize"] = tea.String("200")
	response, err := client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.data.list", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.list", response)
	}

	result, _ := resp.([]interface{})
	sweepRocketmqInstances(result, func(item map[string]interface{}) error {
		deleteAction := fmt.Sprintf("/instances/%s", item["instanceId"])
		body := map[string]interface{}{"instanceId": item["instanceId"]}
		_, err := client.RoaDelete("RocketMQ", "2022-08-01", deleteAction, make(map[string]*string), nil, body, true)
		return err
	})

	return nil
}

// sweepTestCreateTime renders a creation time the way the RocketMQ
// ListInstances API does. The shared parser reads that layout as UTC, so the
// value is formatted as UTC to keep the computed age exact in tests.
func sweepTestCreateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}

func sweepTestItem(id, name string, createTime interface{}) map[string]interface{} {
	item := map[string]interface{}{"instanceId": id, "instanceName": name}
	if createTime != nil {
		item["createTime"] = createTime
	}
	return item
}

func sweepTestRemovals(t *testing.T, items []interface{}, removeErr map[string]error) []string {
	t.Helper()
	var removed []string
	sweepRocketmqInstances(items, func(item map[string]interface{}) error {
		id := item["instanceId"].(string)
		removed = append(removed, id)
		return removeErr[id]
	})
	return removed
}

func TestUnitRocketmqSweeperNamePrefixes(t *testing.T) {
	for _, name := range []string{
		"tf-testAcc-instance",
		"tf_testAcc-instance",
		"tf_test_instance",
		"tf-test-instance",
		"testAcc-instance",
		"TF-TESTACC-INSTANCE",
		"tf_testacC-mixed",
	} {
		if !sweepRocketmqInstanceMatchesName(name) {
			t.Fatalf("expected %q to match an AccTest prefix", name)
		}
	}
	for _, name := range []string{"", "prod-instance", "xtf-testAcc-instance", "tf testAcc instance"} {
		if sweepRocketmqInstanceMatchesName(name) {
			t.Fatalf("expected %q not to match an AccTest prefix", name)
		}
	}
}

func TestUnitRocketmqSweeperAgeProtection(t *testing.T) {
	old := sweepTestCreateTime(time.Now().Add(-72 * time.Hour))
	recent := sweepTestCreateTime(time.Now().Add(-time.Hour))
	future := sweepTestCreateTime(time.Now().Add(48 * time.Hour))

	for _, tc := range []struct {
		name           string
		sweepAll       bool
		item           map[string]interface{}
		skip           bool
		reasonContains string
	}{
		{"recent_matching_name", false, sweepTestItem("id", "tf-testAcc-recent", recent), true, "newer than"},
		{"recent_sweep_all", true, sweepTestItem("id", "prod-recent", recent), true, "newer than"},
		{"future_matching_name", false, sweepTestItem("id", "tf_testAcc-future", future), true, "newer than"},
		{"missing_time_matching_name", false, sweepTestItem("id", "tf-test-1", nil), true, "no creation time field found"},
		{"missing_time_sweep_all", true, sweepTestItem("id", "prod-1", nil), true, "no creation time field found"},
		{"nil_time_sweep_all", true, map[string]interface{}{"instanceId": "id", "instanceName": "prod-2", "createTime": nil}, true, "no creation time field found"},
		{"empty_time_matching_name", false, sweepTestItem("id", "testAcc-1", ""), true, "no creation time field found"},
		{"invalid_time_matching_name", false, sweepTestItem("id", "tf_test_1", "not-a-time"), true, "cannot parse creation time"},
		{"invalid_time_sweep_all", true, sweepTestItem("id", "prod-3", "2021/04/18 15:02:37"), true, "cannot parse creation time"},
		{"old_matching_name", false, sweepTestItem("id", "tf-testAcc-old", old), false, ""},
		{"old_sweep_all", true, sweepTestItem("id", "prod-old", old), false, ""},
		{"old_non_matching_name_without_sweep_all", false, sweepTestItem("id", "prod-old", old), true, "name does not carry"},
		{"recent_non_matching_name_without_sweep_all", false, sweepTestItem("id", "prod-recent", recent), true, "name does not carry"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "24h")
			if tc.sweepAll {
				t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
			} else {
				t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
			}
			skip, reason := sweepRocketmqInstanceSkip(tc.item)
			if skip != tc.skip {
				t.Fatalf("got skip=%t reason=%q, want skip=%t", skip, reason, tc.skip)
			}
			if tc.reasonContains == "" {
				if reason != "" {
					t.Fatalf("got unexpected reason %q", reason)
				}
			} else if !strings.Contains(reason, tc.reasonContains) {
				t.Fatalf("got reason %q, want it to contain %q", reason, tc.reasonContains)
			}
		})
	}
}

func TestUnitRocketmqSweeperDeleteBoundary(t *testing.T) {
	old := sweepTestCreateTime(time.Now().Add(-72 * time.Hour))
	recent := sweepTestCreateTime(time.Now().Add(-time.Hour))
	future := sweepTestCreateTime(time.Now().Add(48 * time.Hour))

	t.Run("prefix_mode", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "24h")
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
		items := []interface{}{
			sweepTestItem("i-old-match", "tf-testAcc-old", old),
			sweepTestItem("i-old-other", "prod-old", old),
			sweepTestItem("i-recent-match", "tf-testAcc-recent", recent),
			sweepTestItem("i-missing-match", "tf_testAcc-missing", nil),
			sweepTestItem("i-bad-match", "tf_test_-bad", "bogus"),
			sweepTestItem("i-future-match", "tf-test--future", future),
			"unexpected-entry",
		}
		removed := sweepTestRemovals(t, items, nil)
		if len(removed) != 1 || removed[0] != "i-old-match" {
			t.Fatalf("got removals %v, want only [i-old-match]", removed)
		}
	})

	t.Run("sweep_all_mode", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "24h")
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
		items := []interface{}{
			sweepTestItem("i-old-other", "prod-old", old),
			sweepTestItem("i-recent-other", "prod-recent", recent),
			sweepTestItem("i-missing-other", "prod-missing", nil),
			sweepTestItem("i-bad-other", "prod-bad", "bogus"),
		}
		removed := sweepTestRemovals(t, items, nil)
		if len(removed) != 1 || removed[0] != "i-old-other" {
			t.Fatalf("got removals %v, want only [i-old-other]", removed)
		}
	})

	t.Run("remove_error_does_not_stop_others", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "24h")
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
		items := []interface{}{
			sweepTestItem("i-old-1", "prod-1", old),
			sweepTestItem("i-old-2", "prod-2", old),
		}
		removed := sweepTestRemovals(t, items, map[string]error{"i-old-1": fmt.Errorf("simulated delete failure")})
		if len(removed) != 2 {
			t.Fatalf("got removals %v, want both instances attempted", removed)
		}
	})
}

func TestUnitRocketmqSweeperLegacySemanticsWithoutOlderThan(t *testing.T) {
	recent := sweepTestCreateTime(time.Now().Add(-time.Hour))

	for _, tc := range []struct {
		name      string
		olderThan string
		sweepAll  bool
		item      map[string]interface{}
		skip      bool
	}{
		{"no_threshold_matching_without_time", "", false, sweepTestItem("id", "tf-testAcc-legacy", nil), false},
		{"no_threshold_matching_recent", "", false, sweepTestItem("id", "tf_testAcc-recent", recent), false},
		{"no_threshold_non_matching", "", false, sweepTestItem("id", "prod-legacy", nil), true},
		{"invalid_threshold_matching", "not-a-duration", false, sweepTestItem("id", "tf_test_-legacy", nil), false},
		{"no_threshold_sweep_all_non_matching", "", true, sweepTestItem("id", "prod-legacy", nil), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", tc.olderThan)
			if tc.sweepAll {
				t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
			} else {
				t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
			}
			if skip, reason := sweepRocketmqInstanceSkip(tc.item); skip != tc.skip {
				t.Fatalf("got skip=%t reason=%q, want skip=%t", skip, reason, tc.skip)
			}
		})
	}

	t.Run("removal_uses_name_selection_only", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
		items := []interface{}{
			sweepTestItem("i-match", "tf-testAcc-legacy", nil),
			sweepTestItem("i-other", "prod-legacy", nil),
		}
		removed := sweepTestRemovals(t, items, nil)
		if len(removed) != 1 || removed[0] != "i-match" {
			t.Fatalf("got removals %v, want only [i-match]", removed)
		}
	})
}
