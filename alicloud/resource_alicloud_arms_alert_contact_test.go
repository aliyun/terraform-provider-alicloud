package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_alert_contact", &resource.Sweeper{
		Name: "alicloud_arms_alert_contact",
		F:    testSweepArmsAlertContact,
	})
}

func testSweepArmsAlertContact(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := cms.CreateDescribeContactListRequest()

	raw, err := cmsService.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
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

func TestAccAlicloudArmsAlertContact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_alert_contact.default"
	ra := resourceAttrInit(resourceId, ArmsAlertContactMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsAlertContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertContact%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsAlertContactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}",
					"email":              "hello.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name,
						"email":              "hello.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "helloupdate.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "helloupdate.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"phone_num": "12345678900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"phone_num": "12345678900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_noc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_noc": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name":     "${var.name}",
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
					"email":                  "hello.uuuu@aaa.com",
					"phone_num":              "12345678901",
					"system_noc":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name":     name,
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
						"email":                  "hello.uuuu@aaa.com",
						"phone_num":              "12345678901",
						"system_noc":             "false",
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

var ArmsAlertContactMap = map[string]string{
	"ding_robot_webhook_url": "",
	"email":                  "",
	"phone_num":              "",
	"system_noc":             "false",
}

func ArmsAlertContactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
