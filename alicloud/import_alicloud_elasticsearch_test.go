package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudElasticsearch_import(t *testing.T) {
	resourceName := "alicloud_elasticsearch.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticsearchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(EcsInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"es_admin_password"},
			},
		},
	})
}
