package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigPolicy_basic(t *testing.T) {
	resourceID := "alicloud_apig_policy.default"
	name := fmt.Sprintf("tfaccapigpolicy%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceID, name, alicloudApigPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name":          name,
					"class_name":           "RateLimit",
					"config":               "{\\\"behaviorType\\\":0,\\\"bodyEncoding\\\":1,\\\"enable\\\":true,\\\"responseContentBody\\\":\\\"rate limited\\\",\\\"responseRedirectUrl\\\":\\\"\\\",\\\"responseStatusCode\\\":429,\\\"threshold\\\":10}",
					"description":          "Terraform acceptance policy",
					"attach_resource_ids":  []string{"${alicloud_apig_route.policy_primary.route_id}"},
					"attach_resource_type": "GatewayRoute",
					"environment_id":       "${alicloud_apig_gateway.policy_primary.environments.0.environment_id}",
					"gateway_id":           "${alicloud_apig_gateway.policy_primary.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "policy_name", name),
					resource.TestCheckResourceAttr(resourceID, "description", "Terraform acceptance policy"),
					resource.TestCheckResourceAttr(resourceID, "attach_resource_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceID, "policy_attachment_id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name":         name + "-updated",
					"config":              "{\\\"behaviorType\\\":0,\\\"bodyEncoding\\\":1,\\\"enable\\\":true,\\\"responseContentBody\\\":\\\"rate limited\\\",\\\"responseRedirectUrl\\\":\\\"\\\",\\\"responseStatusCode\\\":429,\\\"threshold\\\":20}",
					"description":         "Terraform acceptance policy updated",
					"attach_resource_ids": []string{"${alicloud_apig_route.policy_secondary.route_id}"},
					"environment_id":      "${alicloud_apig_gateway.policy_secondary.environments.0.environment_id}",
					"gateway_id":          "${alicloud_apig_gateway.policy_secondary.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "policy_name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceID, "description", "Terraform acceptance policy updated"),
					resource.TestCheckResourceAttr(resourceID, "attach_resource_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateIdFunc: testAccApigPolicyImportID,
				ImportStateVerify: true,
			},
		},
	})
}
