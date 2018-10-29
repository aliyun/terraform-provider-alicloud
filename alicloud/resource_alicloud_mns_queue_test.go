package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_mns_queue", &resource.Sweeper{
		Name: "alicloud_mns_queue",
		F:    testSweepMnsQueues,
	})
}

func testSweepMnsQueues(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	var queueAttrs []ali_mns.QueueAttribute
	for _, namePrefix := range prefixes {
		for {
			var nextMaker string
			raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
				return queueManager.ListQueueDetail(nextMaker, 1000, namePrefix)
			})
			if err != nil {
				return fmt.Errorf("get queueDetails  error: %#v", err)
			}
			queueDetails, _ := raw.(ali_mns.QueueDetails)
			for _, attr := range queueDetails.Attrs {
				queueAttrs = append(queueAttrs, attr)
			}
			nextMaker = queueDetails.NextMarker
			if nextMaker == "" {
				break
			}
		}
	}
	for _, queueAttr := range queueAttrs {
		name := queueAttr.QueueName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping mns queque: %s ", name)
			continue
		}
		log.Printf("[INFO] delete  mns queque: %s ", name)
		_, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return nil, queueManager.DeleteQueue(queueAttr.QueueName)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete mnsQueue (%s (%s)): %s", queueAttr.QueueName, queueAttr.QueueName, err)
		}
	}

	return nil
}

func TestAccAlicloudMnsQueue_basic(t *testing.T) {

	var attr ali_mns.QueueAttribute
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSQueueDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccMNSQueueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSQueueExist("alicloud_mns_queue.queue", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "name", "tf-testAccMNSQueueConfig"),
				),
			},
			{

				Config: testAccMNSQueueConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSQueueExist("alicloud_mns_queue.queue", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "name", "tf-testAccMNSQueueConfig"),
					resource.TestCheckResourceAttr("alicloud_mns_queue.queue", "delay_seconds", "60478"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return queueManager.GetQueueAttributes(rs.Primary.ID)
		})

		if err != nil {
			return err
		}
		instance, _ := raw.(ali_mns.QueueAttribute)
		if instance.QueueName != rs.Primary.ID {
			return fmt.Errorf("mns queue:%s not found", n)
		}

		*attr = instance
		return nil
	}

}

func testAccCheckMNSQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	mnsService := MnsService{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_queue" {
			continue
		}
		_, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return queueManager.GetQueueAttributes(rs.Primary.ID)
		})
		if err != nil {
			if mnsService.QueueNotExistFunc(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("MNS Queue %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccMNSQueueConfig = `
variable "name" {
	default = "tf-testAccMNSQueueConfig"
}
resource "alicloud_mns_queue" "queue"{
	name="${var.name}"
}`

const testAccMNSQueueConfigUpdate = `
variable "name" {
	default = "tf-testAccMNSQueueConfig"
}
resource "alicloud_mns_queue" "queue"{
	name="${var.name}"
	delay_seconds=60478
	maximum_message_size=12357
	message_retention_period=256000
	visibility_timeout=30
	polling_wait_seconds=3
}`
