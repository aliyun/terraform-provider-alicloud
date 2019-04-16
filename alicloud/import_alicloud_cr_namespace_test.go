package alicloud

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCRNamespace_Import(t *testing.T) {
	resourceName := "alicloud_cr_namespace.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRNamespace_Basic(acctest.RandIntRange(100000, 999999)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
