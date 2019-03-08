package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDatahubProject_import(t *testing.T) {
	resourceName := "alicloud_datahub_project.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatahubProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatahubProject(acctest.RandIntRange(datahubProjectSuffixMin, datahubProjectSuffixMax)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
