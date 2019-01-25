package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudElasticsearch_import(t *testing.T) {
	resourceName := "alicloud_elasticsearch_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticsearchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
