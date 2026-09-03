package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeSqlFilesDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudRealtimeComputeSqlFilesSourceConfig(rand, map[string]string{
			"workspace":   `"${alicloud_realtime_compute_vvp_instance.default.resource_id}"`,
			"namespace":   `"${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"`,
			"ids":         `["${alicloud_realtime_compute_sql_file.default.id}"]`,
			"output_file": `"./test_output_file"`,
		}),
		fakeConfig: testAccCheckAliCloudRealtimeComputeSqlFilesSourceConfig(rand, map[string]string{
			"workspace": `"${alicloud_realtime_compute_vvp_instance.default.resource_id}"`,
			"namespace": `"${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"`,
			"ids":       `["${alicloud_realtime_compute_sql_file.default.id}_fake"]`,
		}),
	}

	RealtimeComputeSqlFilesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existRealtimeComputeSqlFilesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"files.#":                      "1",
		"files.0.id":                   CHECKSET,
		"files.0.workspace":            CHECKSET,
		"files.0.namespace":            CHECKSET,
		"files.0.sql_file_id":          CHECKSET,
		"files.0.name":                 fmt.Sprintf("tfaccrealtimecompute%d", rand),
		"files.0.sql_script":           "SELECT 1;",
		"files.0.description":          "tf acc test sql file",
		"files.0.batch_mode":           CHECKSET,
		"files.0.session_cluster_name": CHECKSET,
	}
}

var fakeRealtimeComputeSqlFilesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"files.#": "0",
	}
}

var RealtimeComputeSqlFilesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_realtime_compute_sql_files.default",
	existMapFunc: existRealtimeComputeSqlFilesMapFunc,
	fakeMapFunc:  fakeRealtimeComputeSqlFilesMapFunc,
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
