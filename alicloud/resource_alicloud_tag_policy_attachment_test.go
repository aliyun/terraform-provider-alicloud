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
	resource.AddTestSweepers("alicloud_tag_policy_attachment", &resource.Sweeper{
		Name: "alicloud_tag_policy_attachment",
		F:    testSweepTagPolicyAttachment,
	})
}

func testSweepTagPolicyAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListPoliciesForTarget"
	request := make(map[string]interface{})
	request1 := make(map[string]interface{})

	request["MaxResult"] = PageSizeLarge
	var response map[string]interface{}
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	var attachmentIds []string

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.Target", "EntityNotExist.Policy"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve tag policy attachment in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			request1["PolicyId"] = item["PolicyId"]
			actionListTargetsForPolicy := "ListTargetsForPolicy"
			runtime1 := util.RuntimeOptions{}
			runtime1.SetAutoretry(true)
			response, err = conn.DoRequest(StringPointer(actionListTargetsForPolicy), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request1, &runtime1)
			if err != nil {
				if IsExpectedErrors(err, []string{"EntityNotExists.Target", "EntityNotExist.Policy"}) {
					return nil
				}
				log.Printf("[ERROR] Failed to retrieve tag policy attachment in service list: %s", err)
				return nil
			}
			resp1, err := jsonpath.Get("$.Targets", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
			}
			result1, _ := resp1.([]interface{})
			for _, w := range result1 {
				item1 := w.(map[string]interface{})
				attachmentIds = append(attachmentIds, fmt.Sprintf("%v:%v:%v", item["PolicyId"], item1["TargetId"], item1["TargetType"]))
			}

		}
		if len(result) < PageSizeLarge {
			break
		}
	}

	for _, attachmentId := range attachmentIds {
		log.Printf("[INFO] Delete tag policy attachment: %s ", attachmentId)
		action := "DetachPolicy"
		ids := strings.Split(attachmentId, ":")
		request := map[string]interface{}{
			"PolicyId":   ids[0],
			"TargetId":   ids[1],
			"TargetType": ids[2],
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete tag policy attachment (%s): %s", attachmentId, err)
		}
	}
	return nil
}

func TestAccAlicloudTagPolicyAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_policy_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.TagPolicySupportRegions)
	ra := resourceAttrInit(resourceId, TagPolicyAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicyAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, TagPolicyAttachmentBasicdependence)
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
					"policy_id":   "${alicloud_tag_policy.this.id}",
					"target_id":   "${data.alicloud_account.default.id}",
					"target_type": "USER",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_id":   CHECKSET,
						"target_id":   CHECKSET,
						"target_type": "USER",
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

var TagPolicyAttachmentMap = map[string]string{}

func TagPolicyAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_account" "default" {}

resource "alicloud_tag_policy" "this" {
  policy_name     = "${var.name}"
  policy_desc 	  = "policy_attachment"
  policy_content = <<EOF
        {
    "tags":{
        "CostCenter":{
            "tag_value":{
                "@@assign":[
                    "Beijing",
                    "Shanghai"
                ]
            },
            "tag_key":{
                "@@assign":"CostCenter"
            }
        }
    }
}
    EOF
}
`, name)
}
