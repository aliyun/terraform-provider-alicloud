package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudMnsTopicDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_mns_topics.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
		dataSourceMnsTopicConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_prefix": alicloud_mns_topic.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_prefix": "${alicloud_mns_topic.default.name}-fake",
		}),
	}

	var existMnsTopicMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                       "1",
			"names.0":                       fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"topics.#":                      "1",
			"topics.0.name":                 fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"topics.0.id":                   fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"topics.0.maximum_message_size": "12357",
			"topics.0.logging_enabled":      "true",
		}
	}

	var fakeMnsTopicMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"topics.#": "0",
		}
	}

	var mnsTopicCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMnsTopicMapFunc,
		fakeMapFunc:  fakeMnsTopicMapFunc,
	}

	mnsTopicCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceMnsTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_mns_topic" "default"{
		name="%s"
		maximum_message_size=12357
		logging_enabled=true
	}
	`, name)
}
