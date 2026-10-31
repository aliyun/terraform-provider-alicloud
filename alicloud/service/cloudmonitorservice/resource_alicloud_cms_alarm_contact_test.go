package cloudmonitorservice_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_alarm_contact", &resource.Sweeper{
		Name: "alicloud_cms_alarm_contact",
		F:    testSweepCmsAlarmContact,
	})
}

func testSweepCmsAlarmContact(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := cms.CreateDescribeContactListRequest()

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Cms Alarm in service list: %s", err)
	}

	var response *cms.DescribeContactListResponse
	response, _ = raw.(*cms.DescribeContactListResponse)

	for _, v := range response.Contacts.Contact {
		name := v.Name
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alarm contact: %s ", name)
				continue
			}
		}
		log.Printf("[INFO] delete alarm contact: %s ", name)

		request := cms.CreateDeleteContactRequest()
		request.ContactName = v.Name
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteContact(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alarm contact (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudCmsAlarmContact_basic(t *testing.T) {
	var v cms.Contact
	resourceId := "alicloud_cms_alarm_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, CmsAlarmContactMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlarmContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCmsAlarmContact%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, CmsAlarmContactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_contact_name": "${var.name}",
					"describe":           "For-test",
					"channels_mail":      "hello.uuuu@aaa.com",
					"lang":               "en",
					"lifecycle": []map[string]interface{}{
						{
							"ignore_changes": []string{"channels_mail"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_contact_name": name,
						"describe":           "For-test",
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

func TestAccAliCloudCmsAlarmContact_modify(t *testing.T) {
	var v cms.Contact
	resourceId := "alicloud_cms_alarm_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, CmsAlarmContactMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlarmContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCmsAlarmContact%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, CmsAlarmContactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_contact_name":     "${var.name}",
					"describe":               "For-test-mod",
					"channels_mail":          "tf-mod-mail1@aaa.com",
					"channels_aliim":         "1000111111",
					"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=tf-test-token-1",
					"channels_sms":           "13800000001",
					"lang":                   "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_contact_name":     name,
						"describe":               "For-test-mod",
						"channels_mail":          "tf-mod-mail1@aaa.com",
						"channels_aliim":         "1000111111",
						"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=tf-test-token-1",
						"channels_sms":           "13800000001",
						"lang":                   "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_contact_name":     "${var.name}",
					"describe":               "For-test-mod-updated",
					"channels_mail":          "tf-mod-mail2@aaa.com",
					"channels_aliim":         "1000222222",
					"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=tf-test-token-2",
					"channels_sms":           "13800000002",
					"lang":                   "zh-cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_contact_name":     name,
						"describe":               "For-test-mod-updated",
						"channels_mail":          "tf-mod-mail2@aaa.com",
						"channels_aliim":         "1000222222",
						"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=tf-test-token-2",
						"channels_sms":           "13800000002",
						"lang":                   "zh-cn",
					}),
				),
			},
		},
	})
}

var CmsAlarmContactMap = map[string]string{}

func CmsAlarmContactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
