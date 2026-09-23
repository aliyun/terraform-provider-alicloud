// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ResourceManager MessageContact. >>> Resource test cases, automatically generated.
// Case MessageContact-资源用例new 11390
func TestAccAliCloudResourceManagerMessageContact_basic11390(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_message_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerMessageContactMap11390)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerMessageContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerMessageContactBasicDependence11390)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"message_types": []string{
						"AccountExpenses"},
					"phone_number":         "86-18626811111",
					"title":                "TechnicalDirector",
					"email_address":        "resourcetest@126.com",
					"message_contact_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message_types.#":      "1",
						"phone_number":         "86-18626811111",
						"title":                "TechnicalDirector",
						"email_address":        "resourcetest@126.com",
						"message_contact_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"message_types": []string{
						"ProductMessage", "SecurityMessage", "ActivityMessage"},
					"phone_number":         "86-18626811112",
					"title":                "MaintenanceDirector",
					"email_address":        "resourcetestnew@126.com",
					"message_contact_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message_types.#":      "3",
						"phone_number":         "86-18626811112",
						"title":                "MaintenanceDirector",
						"email_address":        "resourcetestnew@126.com",
						"message_contact_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"message_types": []string{
						"ServiceMessage"},
					"message_contact_name": name + "new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message_types.#":      "1",
						"message_contact_name": name + "new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudResourceManagerMessageContactMap11390 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudResourceManagerMessageContactBasicDependence11390(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager MessageContact. <<< Resource test cases, automatically generated.
