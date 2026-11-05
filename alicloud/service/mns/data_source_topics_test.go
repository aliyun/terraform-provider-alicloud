// Package mns_test.
package mns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMnsTopicDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_mns_topics.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
		dataSourceMnsTopicConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_prefix": "${alicloud_mns_topic.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var mnsTopicCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMnsTopicMapFunc,
		FakeMapFunc:  fakeMnsTopicMapFunc,
	}

	mnsTopicCheckInfo.DataSourceTestCheck(t, rand, allConf)
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
