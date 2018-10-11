package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDatahubTopic_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"name", "tf_testacc_datahub_topic_basic"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubTopic_Tuple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopicTuple,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"name", "tf_testacc_datahub_topic_tuple"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubTopic_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"life_cycle", "7"),
				),
			},

			resource.TestStep{
				Config: testAccDatahubTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"life_cycle", "1"),
				),
			},
		},
	})
}

func testAccCheckDatahubTopicExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub topic: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub topic ID is set")
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		projectName := split[0]
		topicName := split[1]
		_, err := dh.GetTopic(projectName, topicName)

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckDatahubTopicDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_datahub_topic" {
			continue
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		projectName := split[0]
		topicName := split[1]
		_, err := dh.GetTopic(projectName, topicName)

		if err != nil && isDatahubNotExistError(err) {
			continue
		}

		return fmt.Errorf("Datahub topic %s may still exist", rs.Primary.ID)
	}

	return nil
}

const testAccDatahubTopic = `
variable "project_name" {
  default = "tf_testacc_datahub_project"
}
variable "topic_name" {
  default = "tf_testacc_datahub_topic_basic"
}
variable "record_type" {
  default = "BLOB"
}
resource "alicloud_datahub_project" "basic" {
  name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  name = "${var.topic_name}"
  project_name = "${alicloud_datahub_project.basic.name}"
  record_type = "${var.record_type}"
  shard_count = 3
  life_cycle = 7
  comment = "topic for basic."
}
`

const testAccDatahubTopicUpdate = `
variable "project_name" {
  default = "tf_testacc_datahub_project"
}
variable "topic_name" {
  default = "tf_testacc_datahub_topic_basic"
}
variable "record_type" {
  default = "BLOB"
}
resource "alicloud_datahub_project" "basic" {
  name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  name = "${var.topic_name}"
  project_name = "${alicloud_datahub_project.basic.name}"
  record_type = "${var.record_type}"
  shard_count = 3
  life_cycle = 1
  comment = "topic for update."
}
`

const testAccDatahubTopicTuple = `
variable "project_name" {
  default = "tf_testacc_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  name = "tf_testacc_datahub_topic_tuple"
  project_name = "${alicloud_datahub_project.basic.name}"
  record_type = "TUPLE"
  record_schema = {
    bigint_field = "BIGINT"
    timestamp_field = "TIMESTAMP"
    string_field = "STRING"
    double_field = "DOUBLE"
    boolean_field = "BOOLEAN"
  }
  shard_count = 3
  life_cycle = 7
  comment = "a tuple topic."
}
`
