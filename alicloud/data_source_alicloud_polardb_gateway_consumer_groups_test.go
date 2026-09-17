package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayConsumerGroups_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: `data "alicloud_polardb_gateway_consumer_groups" "default" {
  gateway_id = "pg-test"
  ids        = ["cg-test"]
}`,
		ExpectNonEmptyPlan: true,
		PlanOnly:           true,
	}}})
}
