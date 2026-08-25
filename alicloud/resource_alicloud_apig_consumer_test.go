package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigConsumer_basic(t *testing.T) {
	resourceID := "alicloud_apig_consumer.default"
	name := fmt.Sprintf("tfaccapigconsumer%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceID, name, func(name string) string {
		return fmt.Sprintf("variable \"name\" { default = %q }\n", name)
	})
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
					"consumer_name":            name,
					"description":              "Terraform acceptance consumer",
					"enable":                   true,
					"gateway_type":             "API",
					"credential_generate_mode": "System",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "consumer_name", name),
					resource.TestCheckResourceAttr(resourceID, "description", "Terraform acceptance consumer"),
					resource.TestCheckResourceAttr(resourceID, "enable", "true"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Terraform acceptance consumer updated",
					"enable":      false,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "description", "Terraform acceptance consumer updated"),
					resource.TestCheckResourceAttr(resourceID, "enable", "false"),
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
