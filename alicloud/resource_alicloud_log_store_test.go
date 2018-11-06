package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogStore_basic(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogStoreBasic,
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

func testAccCheckAlicloudLogStoreExists(name string, store *sls.LogStore) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Log store ID is set")
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}

		logstore, err := logService.DescribeLogStore(split[0], split[1])
		if err != nil {
			return err
		}
		if logstore == nil || logstore.Name == "" {
			return fmt.Errorf("Log store %s is not exist.", split[1])
		}
		store = logstore

		return nil
	}
}

func testAccCheckAlicloudLogStoreDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	logService := LogService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_store" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		if _, err := logService.DescribeLogStore(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check log store got an error: %#v.", err)
		}
		return fmt.Errorf("Log store %s still exists.", split[0])
	}

	return nil
}

const testAlicloudLogStoreBasic = `
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
`
