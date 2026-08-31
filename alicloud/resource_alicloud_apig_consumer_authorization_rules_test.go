package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigConsumerAuthorizationRules_basic(t *testing.T) {
	resourceID := "alicloud_apig_consumer_authorization_rules.default"
	name := fmt.Sprintf("tfaccapigauth%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceID, name, alicloudApigAuthorizationBasicDependence)
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
					"consumer_id":        "${alicloud_apig_consumer.lifecycle.id}",
					"environment_id":     "${alicloud_apig_http_api_deployment.authorization.environment_id}",
					"parent_resource_id": "${alicloud_apig_http_api.lifecycle.id}",
					"resource_type":      "HttpApiRoute",
					"resource_ids":       []string{"${alicloud_apig_route.lifecycle.route_id}"},
					"principal_type":     "Consumer",
					"expire_mode":        "LongTerm",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "resource_type", "HttpApiRoute"),
					resource.TestCheckResourceAttr(resourceID, "resource_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceID, "authorization_rule_ids.%", "1"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateIdFunc: testAccApigAuthorizationRulesImportID,
				ImportStateVerify: true,
			},
		},
	})
}
