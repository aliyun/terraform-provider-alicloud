package alicloud

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogLogsResource_basic(t *testing.T) {
	if val := os.Getenv("TEST_LOG_LOGS"); val != "true" {
		return
	}
	var project sls.LogProject
	var store sls.LogStore

	config := fmt.Sprintf(testAlicloudLogLogsBasic, time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.foo", &store),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "retention_period", "3000"),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "shards.#", "1"),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "auto_split", "true"),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "max_split_shard_count", "60"),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "append_meta", "true"),
					resource.TestCheckResourceAttr("alicloud_log_store.foo", "enable_web_tracking", "false"),
				),
			},
		},
	})
}

const testAlicloudLogLogsBasic = `
variable "name" {
    default = "tf-testacc-log-store"
}
resource "alicloud_log_project" "foo" {
    name = "${var.name}"
    description = "tf unit test"
}
resource "alicloud_log_store" "foo" {
    project = "${alicloud_log_project.foo.name}"
    name = "${var.name}"
    retention_period = 3000
	shard_count = 1
	auto_split = true
	max_split_shard_count = 60
	append_meta = true
	enable_web_tracking = false
}
resource "alicloud_log_logs" "foo" {
    project = "${alicloud_log_project.foo.name}"
    logstore = "${alicloud_log_store.foo.name}"
	source = "10.1.2.3"
	topic = "test_topic"
	logs = [
		{
		  time = %d
		  contents = {
			key1 = "value1"
			key2 = "value2"
			key3 = "value3"
		   }
		}
	]
	tags = {
		tag1 = "value1"
		tag2 = "value2"
	   }
}
`
