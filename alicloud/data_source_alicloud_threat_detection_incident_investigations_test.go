package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudThreatDetectionIncidentInvestigationsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionIncidentInvestigationsDataSourceName(rand, map[string]string{
			"output_file": `"out.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionIncidentInvestigationsDataSourceName(rand, map[string]string{
			"incident_uuid": `"fake-incident-uuid-not-exist"`,
		}),
	}

	var existAlicloudThreatDetectionIncidentInvestigationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            CHECKSET,
			"investigations.#": CHECKSET,
		}
	}
	var fakeAlicloudThreatDetectionIncidentInvestigationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"investigations.#": "0",
		}
	}
	var alicloudThreatDetectionIncidentInvestigationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_threat_detection_incident_investigations.default",
		existMapFunc: existAlicloudThreatDetectionIncidentInvestigationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudThreatDetectionIncidentInvestigationsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudThreatDetectionIncidentInvestigationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudThreatDetectionIncidentInvestigationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccThreatDetectionIncidentInvestigation-%d"
	}

	data "alicloud_threat_detection_incident_investigations" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
