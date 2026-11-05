package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionInstancesDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionInstanceSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	ThreatDetectionInstanceCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existThreatDetectionInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                    "2",
		"instances.#":              "2",
		"instances.0.id":           acctest.CHECKSET,
		"instances.0.create_time":  acctest.CHECKSET,
		"instances.0.instance_id":  acctest.CHECKSET,
		"instances.0.payment_type": acctest.CHECKSET,
		"instances.0.status":       acctest.CHECKSET,
	}
}

var fakeThreatDetectionInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var ThreatDetectionInstanceCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_instances.default",
	ExistMapFunc: existThreatDetectionInstanceMapFunc,
	FakeMapFunc:  fakeThreatDetectionInstanceMapFunc,
}

func testAccCheckAlicloudThreatDetectionInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionInstance%d"
}

data "alicloud_threat_detection_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
