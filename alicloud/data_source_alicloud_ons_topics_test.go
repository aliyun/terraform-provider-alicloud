package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudOnsTopicsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ons_topics.default"
	name := fmt.Sprintf("tf-testacc%sonstopic%v", defaultRegionToTest, rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": alicloud_ons_instance.default.id,
			"name_regex":  alicloud_ons_topic.default.topic,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": alicloud_ons_instance.default.id,
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	var existOnsTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"topics.#":                    "1",
			"topics.0.topic":              fmt.Sprintf("tf-testacc%sonstopic%v", defaultRegionToTest, rand),
			"topics.0.relation":           "1",
			"topics.0.message_type":       "0",
			"topics.0.independent_naming": "true",
			"topics.0.remark":             "alicloud_ons_topic_remark",
		}
	}

	var fakeOnsTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#": "0",
			"names.#":  "0",
		}
	}

	var onsTopicsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsTopicsMapFunc,
		fakeMapFunc:  fakeOnsTopicsMapFunc,
	}

	onsTopicsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceOnsTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "topic" {
 default = "%v"
}

resource "alicloud_ons_instance" "default" {
name = "tf-testaccOnsInstanceTopicbasic"
}

resource "alicloud_ons_topic" "default" {
  instance_id = alicloud_ons_instance.default.id
  topic = var.topic
  message_type = "0"
  remark = "alicloud_ons_topic_remark"
  perm = "6"
}
`, name)
}
