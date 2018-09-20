package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAlicloudMNSTopic_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// 配置 资源销毁结果检查函数
		CheckDestroy: testAccCheckMNSTopicDestroy,

		Steps: []resource.TestStep{
			{
				// 配置 配置内容
				Config: testAccMNSTopicConfig,
				// 配置 验证函数
				Check: resource.ComposeTestCheckFunc(
					// 验证资源ID
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					// 验证资源属性，能匹配到，肯定就是创建成功了
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", "tf-testAccMNSTopicConfig"),
				),
			},
			{
				// 配置 配置内容
				Config: testAccMNSTopicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					// 验证修改后的属性值，如果能匹配到，肯定就是修改成功了
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", "tf-testAccMNSTopicConfig"),

					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "maximum_message_size", "12357"),

					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "logging_enabled", "true"),
				),
			},
		},
	})
}

func testAccMNSTopicExist(n string, attr *ali_mns.TopicAttribute) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MNSTopic ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
		}
		topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
		instance, err := topicManager.GetTopicAttributes(rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.TopicName != rs.Primary.ID {
			return fmt.Errorf("mns Topic:%snot found", n)
		}

		*attr = instance
		return nil
	}

}

func testAccCheckMNSTopicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_topic" {
			continue
		}
		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating MNS Topic  error: %#v", err)
		}
		topicManager := ali_mns.NewMNSTopicManager(*mnsClient)

		if _, err := topicManager.GetTopicAttributes(rs.Primary.ID); err != nil {
			if strings.Contains(err.Error(), "TopicNotExist") {
				continue
			}
			return err
		}

		return fmt.Errorf("MNS Topic %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccMNSTopicConfig = `variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
}`

const testAccMNSTopicConfigUpdate = `variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
	maximum_message_size=12357
	logging_enabled=true
}`
