package alicloud

import (
	"testing"
)

func TestAccAliCloudThreatDetectionIncidentsDataSource_basic(t *testing.T) {
	allConf := dataSourceTestAccConfig{
		existConfig: `
data "alicloud_threat_detection_incidents" "default" {
  page_size = 10
}
`,
		fakeConfig: `
data "alicloud_threat_detection_incidents" "default" {
  page_size = 10
}
`,
	}

	ThreatDetectionIncidentsCheckInfo.dataSourceTestCheck(t, 0, allConf)
}

var existThreatDetectionIncidentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       CHECKSET,
		"incidents.#": CHECKSET,
	}
}

var fakeThreatDetectionIncidentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       CHECKSET,
		"incidents.#": CHECKSET,
	}
}

var ThreatDetectionIncidentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_incidents.default",
	existMapFunc: existThreatDetectionIncidentsMapFunc,
	fakeMapFunc:  fakeThreatDetectionIncidentsMapFunc,
}
