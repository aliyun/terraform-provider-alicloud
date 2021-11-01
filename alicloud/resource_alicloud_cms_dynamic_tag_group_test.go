package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_cms_dynamic_tag_group",
		&resource.Sweeper{
			Name: "alicloud_cms_dynamic_tag_group",
			F:    testSweepCmsDynamicTagGroup,
		})
}
func testSweepCmsDynamicTagGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDynamicTagRuleList"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.TagGroupList.TagGroup", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TagKey"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DynamicTagGroup: %s", item["TagKey"].(string))
				continue
			}

			action := "DeleteDynamicTagGroup"
			request := map[string]interface{}{
				"DynamicTagRuleId": item["DynamicTagRuleId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DynamicTagGroup (%s): %s", item["TagKey"].(string), err)
			}
			log.Printf("[INFO] Delete DynamicTagGroup success: %s ", item["TagKey"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudCloudMonitorServiceDynamicTagGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CmsDynamicTagGroupSupportRegions)
	resourceId := "alicloud_cms_dynamic_tag_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceDynamicTagGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDynamicTagGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceDynamicTagGroupBasicDependence0)
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
					"contact_group_list": []string{"${alicloud_cms_alarm_contact_group.default.id}", "${alicloud_cms_alarm_contact_group.default0.id}"},
					"tag_key":            "appgroup",
					"match_express": []map[string]interface{}{
						{
							"tag_value":                "landingzone",
							"tag_value_match_function": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_group_list.#": "2",
						"tag_key":              "appgroup",
						"match_express.#":      "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"contact_group_list"},
			},
		},
	})
}

var AlicloudCloudMonitorServiceDynamicTagGroupMap0 = map[string]string{
	"match_express.#":               CHECKSET,
	"status":                        CHECKSET,
	"contact_group_list.#":          CHECKSET,
	"tag_key":                       CHECKSET,
	"match_express_filter_relation": CHECKSET,
	"template_id_list.#":            CHECKSET,
}

func AlicloudCloudMonitorServiceDynamicTagGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_cms_alarm_contact_group" "default" {
	alarm_contact_group_name = "${var.name}-1"
	describe                 = "For Test"
	enable_subscribed        = true
}

resource "alicloud_cms_alarm_contact_group" "default0" {
	alarm_contact_group_name = "${var.name}-0"
	describe                 = "For Test"
	enable_subscribed        = true
}

`, name)
}
