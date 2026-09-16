package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayConsumers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: `data "alicloud_polardb_gateway_consumers" "default" {
  gateway_id        = "pg-test"
  ids               = ["c-test"]
  consumer_group_id = "cg-test"
}`,
		ExpectNonEmptyPlan: true,
		PlanOnly:           true,
	}}})
}
