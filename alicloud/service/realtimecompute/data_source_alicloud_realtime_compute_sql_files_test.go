package realtimecompute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRealtimeComputeSqlFilesDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudRealtimeComputeSqlFilesSourceConfig(rand, map[string]string{
			"workspace":   `"${alicloud_realtime_compute_vvp_instance.default.resource_id}"`,
			"namespace":   `"${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"`,
			"ids":         `["${alicloud_realtime_compute_sql_file.default.id}"]`,
			"output_file": `"./test_output_file"`,
		}),
		FakeConfig: testAccCheckAliCloudRealtimeComputeSqlFilesSourceConfig(rand, map[string]string{
			"workspace": `"${alicloud_realtime_compute_vvp_instance.default.resource_id}"`,
			"namespace": `"${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"`,
			"ids":       `["${alicloud_realtime_compute_sql_file.default.id}_fake"]`,
		}),
	}

	RealtimeComputeSqlFilesCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

var existRealtimeComputeSqlFilesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"files.#":                      "1",
		"files.0.id":                   acctest.CHECKSET,
		"files.0.workspace":            acctest.CHECKSET,
		"files.0.namespace":            acctest.CHECKSET,
		"files.0.sql_file_id":          acctest.CHECKSET,
		"files.0.name":                 fmt.Sprintf("tfaccrealtimecompute%d", rand),
		"files.0.sql_script":           "SELECT 1;",
		"files.0.description":          "tf acc test sql file",
		"files.0.batch_mode":           acctest.CHECKSET,
		"files.0.session_cluster_name": acctest.CHECKSET,
	}
}

var fakeRealtimeComputeSqlFilesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"files.#": "0",
	}
}

var RealtimeComputeSqlFilesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_realtime_compute_sql_files.default",
	ExistMapFunc: existRealtimeComputeSqlFilesMapFunc,
	FakeMapFunc:  fakeRealtimeComputeSqlFilesMapFunc,
}

func testAccCheckAliCloudRealtimeComputeSqlFilesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	return AlicloudRealtimeComputeSqlFileBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_realtime_compute_sql_file" "default" {
  workspace            = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace            = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name                 = "%s"
  sql_script           = "SELECT 1;"
  description          = "tf acc test sql file"
  batch_mode           = "true"
  session_cluster_name = "%s"
}

data "alicloud_realtime_compute_sql_files" "default" {
  %s
}
`, name, name, strings.Join(pairs, "\n  "))
}
