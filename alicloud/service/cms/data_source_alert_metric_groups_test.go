package cms_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCmsAlertMetricGroupsDataSource exercises every Optional filter
// of the data source: ids, datasource_type, include_details and output_file.
//
// Alert metric groups are predefined, read-only system entries shared by all
// accounts, so the test creates no resource: an unfiltered data source ("all")
// provides a real alert metric group id that each filtered config then pins on,
// keeping the exist assertions deterministic. The fake configs apply a value no
// group matches, exercising the empty-results path.
func TestAccAliCloudCmsAlertMetricGroupsDataSource(t *testing.T) {
	// deterministic per-run suffix for the fake ids and output file names.
	rand := helper.RandIntRange(10000, 99999)

	realID := "data.alicloud_cms_alert_metric_groups.all.groups.0.id"
	fakeID := fmt.Sprintf(`"${data.alicloud_cms_alert_metric_groups.all.groups.0.id}_fake_%d"`, rand)
	firstDatasourceType := "jsondecode(data.alicloud_cms_alert_metric_groups.all.groups.0.datasource_types)[0]"

	existChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.alicloud_cms_alert_metric_groups.default", "groups.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_cms_alert_metric_groups.default", "ids.#", "1"),
		resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_metric_groups.default", "groups.0.id"),
		resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_metric_groups.default", "groups.0.alert_metric_group_id"),
		resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_metric_groups.default", "groups.0.datasource_types"),
		resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_metric_groups.default", "groups.0.display_name_cn"),
		resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_metric_groups.default", "groups.0.display_name_en"),
	}
	fakeChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.alicloud_cms_alert_metric_groups.default", "groups.#", "0"),
		resource.TestCheckResourceAttr("data.alicloud_cms_alert_metric_groups.default", "ids.#", "0"),
	}

	cases := []struct {
		name      string
		existAttr map[string]string
		fakeAttr  map[string]string
	}{
		{
			name:      "ids",
			existAttr: map[string]string{"ids": "[" + realID + "]"},
			fakeAttr:  map[string]string{"ids": "[" + fakeID + "]"},
		},
		{
			// datasource_type is a server-side filter; the value is derived from
			// the first group's own datasource_types (a JSON array string such as
			// ["arms_metrics"]) so the pinned id always matches the filter. The
			// ids filter narrows the result to that single group, matching the
			// original SDK v2 contract.
			name: "datasource_type",
			existAttr: map[string]string{
				"ids":             "[" + realID + "]",
				"datasource_type": firstDatasourceType,
			},
			fakeAttr: map[string]string{
				"ids":             "[" + fakeID + "]",
				"datasource_type": `"nonexistent_datasource_type"`,
			},
		},
		{
			name:      "include_details",
			existAttr: map[string]string{"ids": "[" + realID + "]", "include_details": "true"},
			fakeAttr:  map[string]string{"ids": "[" + fakeID + "]", "include_details": "true"},
		},
		{
			name:      "output_file",
			existAttr: map[string]string{"ids": "[" + realID + "]", "output_file": fmt.Sprintf(`"/tmp/cms_alert_metric_groups_ds_output_%d.txt"`, rand)},
			fakeAttr:  map[string]string{"ids": "[" + fakeID + "]", "output_file": fmt.Sprintf(`"/tmp/cms_alert_metric_groups_ds_output_fake_%d.txt"`, rand)},
		},
	}

	steps := make([]resource.TestStep, 0, len(cases)*2)
	for _, c := range cases {
		steps = append(steps,
			resource.TestStep{
				Config: testAccCmsAlertMetricGroupsConfig(c.existAttr),
				Check:  resource.ComposeTestCheckFunc(existChecks...),
			},
			resource.TestStep{
				Config: testAccCmsAlertMetricGroupsConfig(c.fakeAttr),
				Check:  resource.ComposeTestCheckFunc(fakeChecks...),
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// testAccCmsAlertMetricGroupsConfig renders the two-instance data source
// config: an unfiltered "all" instance exposing the real group id, plus the
// "default" instance carrying the filter under test.
func testAccCmsAlertMetricGroupsConfig(attrMap map[string]string) string {
	pairs := make([]string, 0, len(attrMap))
	for k, v := range attrMap {
		pairs = append(pairs, fmt.Sprintf("  %s = %s", k, v))
	}
	return fmt.Sprintf(`
data "alicloud_cms_alert_metric_groups" "all" {}
data "alicloud_cms_alert_metric_groups" "default" {
%s
}
`, strings.Join(pairs, "\n"))
}
