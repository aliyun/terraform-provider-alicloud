package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudLogStoreIndex_fullText(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore
	var index sls.Index

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogStoreIndexFullText,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.foo", &store),
					testAccCheckAlicloudLogStoreIndexExists("alicloud_log_store_index.foo", &index),
					resource.TestCheckResourceAttr("alicloud_log_store_index.foo", "full_text.#", "1"),
					resource.TestCheckResourceAttr("alicloud_log_store_index.foo", "full_text.1.case_sensitive", "true"),
					resource.TestCheckResourceAttr("alicloud_log_store_index.foo", "full_text.1.token", " #$%^*\r\n\t"),
				),
			},
		},
	})
}

func TestAccAlicloudLogStoreIndex_field(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore
	var index sls.Index

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogStoreIndexField,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.bar", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.bar", &store),
					testAccCheckAlicloudLogStoreIndexExists("alicloud_log_store_index.bar", &index),
					resource.TestCheckResourceAttr("alicloud_log_store_index.bar", "field_search.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudLogStoreIndex_all(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore
	var index sls.Index

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogStoreIndexAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.all", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.all", &store),
					testAccCheckAlicloudLogStoreIndexExists("alicloud_log_store_index.all", &index),
					resource.TestCheckResourceAttr("alicloud_log_store_index.all", "full_text.#", "1"),
					resource.TestCheckResourceAttr("alicloud_log_store_index.all", "full_text.1.case_sensitive", "true"),
					resource.TestCheckResourceAttr("alicloud_log_store_index.all", "field_search.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogStoreIndexExists(name string, index *sls.Index) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Log store ID is set")
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		i, err := testAccProvider.Meta().(*AliyunClient).DescribeLogStoreIndex(split[0], split[1])
		if err != nil {
			return err
		}

		index = i
		return nil
	}
}

func testAccCheckAlicloudLogStoreIndexDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_store_index" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		i, err := client.DescribeLogStoreIndex(split[0], split[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check log store index got an error: %#v.", err)
		}

		if len(split) == 2 {
			if i.Line == nil {
				continue
			}
		} else {
			if _, ok := i.Keys[split[2]]; !ok {
				continue
			}
		}

		return fmt.Errorf("Log store index %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAlicloudLogStoreIndexFullText = `
variable "name" {
    default = "tf-test-log-store-index-full"
}
resource "alicloud_log_project" "foo" {
    name = "${var.name}"
    description = "tf unit test"
}
resource "alicloud_log_store" "foo" {
    project = "${alicloud_log_project.foo.name}"
    name = "${var.name}"
    retention_period = "3000"
    shard_count = 1
}
resource "alicloud_log_store_index" "foo" {
    project = "${alicloud_log_project.foo.name}"
    logstore = "${alicloud_log_store.foo.name}"
    full_text {
	case_sensitive = true
	token = " #$%^*\r\n\t"
    }
}
`
const testAlicloudLogStoreIndexField = `
variable "name" {
    default = "tf-test-log-store-index-field"
}
resource "alicloud_log_project" "bar" {
    name = "${var.name}"
    description = "tf unit test"
}
resource "alicloud_log_store" "bar" {
    project = "${alicloud_log_project.bar.name}"
    name = "${var.name}"
    retention_period = "3000"
    shard_count = 1
}
resource "alicloud_log_store_index" "bar" {
    project = "${alicloud_log_project.bar.name}"
    logstore = "${alicloud_log_store.bar.name}"
    field_search {
      name = "${var.name}"
      enable_analytics = true
      token = " #$%^*\r\n\t"
      name = "${var.name}-1"
      type = "text"
    }
}
`

const testAlicloudLogStoreIndexAll = `
variable "name" {
    default = "tf-log-store-index-all"
}
resource "alicloud_log_project" "all" {
    name = "${var.name}"
    description = "tf unit test"
}
resource "alicloud_log_store" "all" {
    project = "${alicloud_log_project.all.name}"
    name = "${var.name}"
    retention_period = "3000"
    shard_count = 1
}

resource "alicloud_log_store_index" "all" {
    project = "${alicloud_log_project.all.name}"
    logstore = "${alicloud_log_store.all.name}"
    full_text {
	case_sensitive = true
	token = " #$%^*\r\n\t"
    }
    field_search = [
    {
	name = "${var.name}-1"
	enable_analytics = true
    },
    {
	token = " #$%^*\r\n\t"
	name = "${var.name}-2"
	type = "text"
    }
    ]
}
`
