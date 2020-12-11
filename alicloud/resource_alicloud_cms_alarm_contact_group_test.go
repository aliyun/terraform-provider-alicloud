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
	resource.AddTestSweepers("alicloud_cms_alarm_contact_group", &resource.Sweeper{
		Name: "alicloud_cms_alarm_contact_group",
		F:    testSweepCmsAlarmContactGroup,
	})
}

func testSweepCmsAlarmContactGroup(region string) error {
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

	request := cms.CreateDescribeContactGroupListRequest()

	raw, err := cmsService.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactGroupList(request)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Cms Alarm Contact Group in service list: %s", err)
	}

	var response *cms.DescribeContactGroupListResponse
	response, _ = raw.(*cms.DescribeContactGroupListResponse)

	for _, v := range response.ContactGroupList.ContactGroup {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alarm contact group: %s ", name)
			continue
		}
		log.Printf("[INFO] delete alarm contact group: %s ", name)

		request := cms.CreateDeleteContactGroupRequest()
		request.ContactGroupName = v.Name
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteContactGroup(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alarm contact group (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlicloudCmsAlarmContactGroup_basic(t *testing.T) {
	var v cms.ContactGroup
	resourceId := "alicloud_cms_alarm_contact_group.default"
	ra := resourceAttrInit(resourceId, CmsAlarmContactGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlarmContactGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactGrouptf-test%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CmsAlarmContactGroupBasicdependence)
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
					"alarm_contact_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_contact_group_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"describe": "tf-test-describe",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"describe": "tf-test-describe",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_subscribed": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_subscribed": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contacts": []string{"${alicloud_cms_alarm_contact.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contacts.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"describe":          "tf-test-describe-update",
					"contacts":          []string{"${alicloud_cms_alarm_contact.default.id}", "${alicloud_cms_alarm_contact.default0.id}"},
					"enable_subscribed": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"describe":          "tf-test-describe-update",
						"contacts.#":        "2",
						"enable_subscribed": "false",
					}),
				),
			},
		},
	})
}

var CmsAlarmContactGroupMap = map[string]string{}

func CmsAlarmContactGroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cms_alarm_contact" "default" {
  alarm_contact_name = "${var.name}-1"
  describe           = "For Test 1234567"
  channels_mail      = "hello.uuuu@aaa.com"
  lifecycle {
    ignore_changes = [channels_mail]
  }
}
resource "alicloud_cms_alarm_contact" "default0" {
  alarm_contact_name = "${var.name}-0"
  describe           = "For Test 1234567"
  channels_mail      = "hello.uuuu@aaa.com"
  lifecycle {
    ignore_changes = [channels_mail]
  }
}
`, name)
}
