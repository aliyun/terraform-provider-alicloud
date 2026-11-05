// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order stock-debt declarations (pure rename migration; see scripts/collection-order/README.md)
// collection-order: bypass message_types -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259

// Test ResourceManager MessageContact. >>> Resource test cases, automatically generated.
// Case MessageContact-资源用例new 11390
func TestAccAliCloudResourceManagerMessageContact_basic11390(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_message_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerMessageContactMap11390)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerMessageContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := tfacctest.RandStringFromCharSet(5, tfacctest.CharSetAlpha)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerMessageContactBasicDependence11390)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudResourceManagerMessageContactBasicDependence11390(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager MessageContact. <<< Resource test cases, automatically generated.
