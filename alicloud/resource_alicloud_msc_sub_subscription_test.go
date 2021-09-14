package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMscSubSubscription_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_subscription.default"
	ra := resourceAttrInit(resourceId, AlicloudMscSubSubscriptionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MscOpenSubscriptionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMscSubSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smscsubsubscription%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMscSubSubscriptionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"item_name":      "Notifications of Product Expiration",
					"sms_status":     "1",
					"email_status":   "1",
					"pmsg_status":    "1",
					"tts_status":     "1",
					"webhook_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"item_name":      "Notifications of Product Expiration",
						"sms_status":     "1",
						"email_status":   "1",
						"pmsg_status":    "1",
						"tts_status":     "1",
						"webhook_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sms_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sms_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pmsg_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pmsg_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tts_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tts_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_ids": []string{"${alicloud_msc_sub_contact.default1.id}", "${alicloud_msc_sub_contact.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_ids.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudMscSubSubscriptionMap0 = map[string]string{}

func AlicloudMscSubSubscriptionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_msc_sub_contact" "default1" {
  contact_name = "tfceo"
  position = "CEO"
  email = "123@163.com"
  mobile = "12312345906"
}
resource "alicloud_msc_sub_contact" "default2" {
  contact_name = "tfdirector"
  position = "Technical Director"
  email = "123@163.com"
  mobile = "12312345906"
}
`, name)
}
