package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigHttpApiDeployment_basic(t *testing.T) {
	resourceID := "alicloud_apig_http_api_deployment.default"
	name := fmt.Sprintf("tfaccapigdeploy%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceID, name, alicloudApigLifecycleBasicDependence)
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
					"http_api_id":    "${alicloud_apig_http_api.lifecycle.id}",
					"route_id":       "${alicloud_apig_route.lifecycle.route_id}",
					"environment_id": "${alicloud_apig_gateway.lifecycle.environments.0.environment_id}",
					"gateway_id":     "${alicloud_apig_gateway.lifecycle.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "status", "Deployed"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
