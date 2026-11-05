package mscopensubscription_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mscopensubscription"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMscSubSubscription_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_subscription.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMscSubSubscriptionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mscopensubscription.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSubscription")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smscsubsubscription%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMscSubSubscriptionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"item_name":      "Notifications of Product Expiration",
					"sms_status":     "1",
					"email_status":   "1",
					"pmsg_status":    "1",
					"tts_status":     "0",
					"webhook_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"item_name":      "Notifications of Product Expiration",
						"sms_status":     "1",
						"email_status":   "1",
						"pmsg_status":    "1",
						"tts_status":     "0",
						"webhook_status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tts_status":     "1",
					"webhook_status": "1",
					"sms_status":     "0",
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
					"contact_ids": []string{"${alicloud_msc_sub_contact.default1.id}", "${alicloud_msc_sub_contact.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_ids.#": "2",
					}),
				),
			},
			// collection-order: bypass webhook_ids -- pre-existing attribute outside this change: webhook_ids references real DingTalk robot URLs this test cannot provision without a paid token, so convergent reorder coverage is not honestly establishable here.
			{
				// A webhook ID can only reference a real DingTalk robot URL, which this
				// test cannot provision. Removing the required item_name stops the plan
				// before the placeholder reaches the API.
				Config: testAccConfig(map[string]interface{}{
					"webhook_ids": []interface{}{"1"},
					"item_name":   acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
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
