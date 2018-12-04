package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDatahubTopic_importBasic(t *testing.T) {
	resourceName := "alicloud_datahub_topic.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatahubTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic(acctest.RandIntRange(datahubProjectSuffixMin, datahubProjectSuffixMax)),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
