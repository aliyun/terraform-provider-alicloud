package cms_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCmsMetricListResourceId = "data.alicloud_cms_metric_list.default"
	testCmsMetricListNamespace  = "acs_ecs_dashboard"
)

// TestAccAliCloudCmsMetricListDataSource exercises the data source against the
// read-only acs_ecs_dashboard CPUUtilization metric, so no resource is
// provisioned. The account is not guaranteed to have recent datapoints, so the
// exist steps assert only structural attributes (list size and request echo).
// The last step queries a metric name that never exists: the API answers
// 10002, which must clear the previous read's state.
func TestAccAliCloudCmsMetricListDataSource(t *testing.T) {
	existCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testCmsMetricListResourceId, "metric_name", "CPUUtilization"),
		resource.TestCheckResourceAttr(testCmsMetricListResourceId, "metric_namespace", testCmsMetricListNamespace),
		resource.TestCheckResourceAttrSet(testCmsMetricListResourceId, "id"),
		resource.TestMatchResourceAttr(testCmsMetricListResourceId, "datapoints.#", regexp.MustCompile(`^\d+$`)),
	)
	shortWindow := testAccCmsMetricListConfig("CPUUtilization", "60", 100, -10*time.Minute, -5*time.Minute)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: shortWindow,
				Check:  existCheck,
			},
			{
				// The identical config re-reads the data source in place.
				Config: shortWindow,
				Check:  existCheck,
			},
			{
				Config: testAccCmsMetricListConfig("CPUUtilization", "300", 100, -2*time.Hour, -1*time.Hour),
				Check: resource.ComposeTestCheckFunc(
					existCheck,
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "period", "300"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "length", "100"),
				),
			},
			{
				// A never-existing metric hits the 10002 branch, which must
				// clear the state of the preceding successful read.
				Config: testAccCmsMetricListConfig("FakeMetricNotExist", "300", 100, -2*time.Hour, -1*time.Hour),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "datapoints.#", "0"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "next_token", ""),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "actual_period", ""),
				),
			},
		},
	})
}

func testAccCmsMetricListConfig(metricName, period string, length int, startOffset, endOffset time.Duration) string {
	now := time.Now()
	// StartTime/EndTime go on the wire as ISO 8601 UTC (Z suffix); the CMS API
	// rejects the RFC 3339 local-offset form (parse param time error), which a
	// non-UTC test machine would otherwise emit.
	return fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
  metric_name      = "%s"
  metric_namespace = "%s"
  start_time       = "%s"
  end_time         = "%s"
  period           = "%s"
  length           = %d
}
`, metricName, testCmsMetricListNamespace, now.Add(startOffset).UTC().Format(time.RFC3339), now.Add(endOffset).UTC().Format(time.RFC3339), period, length)
}

// TestAccAliCloudCmsMetricListDataSourceOutputFile verifies the output_file
// contract through the Terraform read path:
//   - step 1 queries an existing metric and the file must be written as a JSON
//     array; its content is then replaced with a sentinel so step 2 can detect
//     an overwrite even when the account has no datapoints in the window;
//   - step 2 switches to a metric that never exists (10002): the state is
//     cleared and the file must be overwritten with an empty JSON array;
//   - step 3 points output_file at a directory that does not exist: the read
//     must fail instead of reporting success without exporting the result.
func TestAccAliCloudCmsMetricListDataSourceOutputFile(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "cms_metric_list.json")
	missingDirFile := filepath.Join(dir, "missing-dir", "cms_metric_list.json")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsMetricListOutputFileConfig("CPUUtilization", outputFile),
				Check: resource.ComposeTestCheckFunc(
					checkCmsMetricListOutputFileIsJSONArray(outputFile),
					seedCmsMetricListOutputFile(outputFile),
				),
			},
			{
				Config: testAccCmsMetricListOutputFileConfig("FakeMetricNotExist", outputFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "datapoints.#", "0"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "next_token", ""),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "actual_period", ""),
					checkCmsMetricListOutputFileContent(outputFile, "[]"),
				),
			},
			{
				Config:      testAccCmsMetricListOutputFileConfig("CPUUtilization", missingDirFile),
				ExpectError: regexp.MustCompile(`no such file or directory`),
			},
		},
	})
}

func testAccCmsMetricListOutputFileConfig(metricName, outputFile string) string {
	now := time.Now()
	return fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
  metric_name      = "%s"
  metric_namespace = "%s"
  start_time       = "%s"
  end_time         = "%s"
  period           = "300"
  length           = 100
  output_file      = %q
}
`, metricName, testCmsMetricListNamespace, now.Add(-2*time.Hour).UTC().Format(time.RFC3339), now.Add(-1*time.Hour).UTC().Format(time.RFC3339), outputFile)
}

func checkCmsMetricListOutputFileIsJSONArray(file string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("output_file %s was not written: %w", file, err)
		}
		var datapoints []map[string]interface{}
		if err := json.Unmarshal(content, &datapoints); err != nil {
			return fmt.Errorf("output_file %s is not a JSON array: %w", file, err)
		}
		return nil
	}
}

func seedCmsMetricListOutputFile(file string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		return os.WriteFile(file, []byte(`[{"seed":true}]`), 0600)
	}
}

func checkCmsMetricListOutputFileContent(file, expected string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("output_file %s was not written: %w", file, err)
		}
		if string(content) != expected {
			return fmt.Errorf("output_file %s contains %q, expected %q", file, string(content), expected)
		}
		return nil
	}
}
