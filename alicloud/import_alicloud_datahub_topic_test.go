package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDatahubTopic_importBasic(t *testing.T) {
	resourceName := "alicloud_datahub_topic.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatahubTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
