package alicloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMNSTopic_importBasic(t *testing.T) {
	resourceName := "alicloud_mns_topic.topic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMNSTopicConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}