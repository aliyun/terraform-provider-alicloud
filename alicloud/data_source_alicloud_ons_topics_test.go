package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOnsTopicsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ons_topics.default"
	name := fmt.Sprintf("tf-testacc%sonstopic%v", defaultRegionToTest, rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"name_regex":  "${alicloud_ons_topic.default.topic}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"ids":         []string{"${alicloud_ons_topic.default.topic}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"ids":         []string{"fake_tf-testacc*"},
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"ids":         []string{"${alicloud_ons_topic.default.topic}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_ons_instance.default.id}",
			"ids":         []string{"${alicloud_ons_topic.default.topic}"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	var existOnsTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"topics.#":                    "1",
			"topics.0.topic_name":         fmt.Sprintf("tf-testacc%sonstopic%v", defaultRegionToTest, rand),
			"topics.0.relation":           "6",
			"topics.0.relation_name":      "授权发布/订阅",
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

	preCheck := func() {
		testAccPreCheck(t)
	}
	onsTopicsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf)
}

func dataSourceOnsTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "topic" {
 default = "%v"
}

resource "alicloud_ons_instance" "default" {
instance_name = var.topic
}

resource "alicloud_ons_topic" "default" {
  instance_id = "${alicloud_ons_instance.default.id}"
  topic_name = "${var.topic}"
  message_type = "0"
  remark = "alicloud_ons_topic_remark"
  tags 	 = {
	Created = "TF"
	For 	= "acceptance test"
  }
}
`, name)
}
