package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAlicloudMNSQueue_basic(t *testing.T) {

	var attr ali_mns.QueueAttribute

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// 配置 资源销毁结果检查函数
		CheckDestroy: testAccCheckMNSQueueDestroy,

		Steps: []resource.TestStep{
			{
				// 配置 配置内容
				Config: testAccMNSQueueConfig,
				// 配置 验证函数
				Check: resource.ComposeTestCheckFunc(
					// 验证资源ID
					testAccMNSQueueExist("alicloud_mns_queue.queue", &attr),
					// 验证资源属性，能匹配到，肯定就是创建成功了
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "name", "tf-testAccMNSQueueConfig"),
				),
			},
			{
				// 配置 配置内容
				Config: testAccMNSQueueConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSQueueExist("alicloud_mns_queue.queue", &attr),
					// 验证修改后的属性值，如果能匹配到，肯定就是修改成功了
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "name", "tf-testAccMNSQueueConfig"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "delay_seconds", "60482"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "maximum_message_size", "12357"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "message_retention_period", "256000"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "visibility_timeout", "30"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "polling_wait_seconds", "3"),
				),
			},
		},
	})
}

func testAccMNSQueueExist(n string, attr *ali_mns.QueueAttribute) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MNSQueue ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating alicoudMNSQueue  error: %#v", err)
		}
		queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
		instance, err := queueManager.GetQueueAttributes(rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.QueueName != rs.Primary.ID {
			return fmt.Errorf("mns queue:%snot found", n)
		}

		*attr = instance
		return nil
	}

}

func testAccCheckMNSQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_queue" {
			continue
		}
		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating MNS Queue  error: %#v", err)
		}
		queueManager := ali_mns.NewMNSQueueManager(*mnsClient)

		if _, err := queueManager.GetQueueAttributes(rs.Primary.ID); err != nil {
			if strings.Contains(err.Error(), "QueueNotExist") {
				continue
			}
			return err
		}

		return fmt.Errorf("MNS Queue %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccMNSQueueConfig = `variable "name" {
	default = "tf-testAccMNSQueueConfig"
}
resource "alicloud_mns_queue" "queue"{
	name="${var.name}"
}`

const testAccMNSQueueConfigUpdate = `variable "name" {
	default = "tf-testAccMNSQueueConfig"
}
resource "alicloud_mns_queue" "queue"{
	name="${var.name}"
	delay_seconds=60482
	maximum_message_size=12357
	message_retention_period=256000
	visibility_timeout=30
	polling_wait_seconds=3
}`
