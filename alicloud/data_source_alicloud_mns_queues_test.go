package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudMnsQueueDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_mns_queues.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMNSQueue-%d", rand),
		dataSourceMnsQueueConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_prefix": alicloud_mns_queue.queue.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_prefix": "${alicloud_mns_queue.queue.name}-fake",
		}),
	}

	var existMnsQueueMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                           "1",
			"names.0":                           fmt.Sprintf("tf-testAccMNSQueue-%d", rand),
			"queues.#":                          "1",
			"queues.0.id":                       fmt.Sprintf("tf-testAccMNSQueue-%d", rand),
			"queues.0.name":                     fmt.Sprintf("tf-testAccMNSQueue-%d", rand),
			"queues.0.delay_seconds":            "60478",
			"queues.0.maximum_message_size":     "12357",
			"queues.0.message_retention_period": "256000",
			"queues.0.visibility_timeouts":      "30",
			"queues.0.polling_wait_seconds":     "3",
		}
	}

	var fakeMnsQueueMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"queues.#": "0",
		}
	}

	var mnsQueueCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMnsQueueMapFunc,
		fakeMapFunc:  fakeMnsQueueMapFunc,
	}

	mnsQueueCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceMnsQueueConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_mns_queue" "queue"{
		name="%s"
		delay_seconds=60478
		maximum_message_size=12357
		message_retention_period=256000
		visibility_timeout=30
		polling_wait_seconds=3
	}
	`, name)
}
