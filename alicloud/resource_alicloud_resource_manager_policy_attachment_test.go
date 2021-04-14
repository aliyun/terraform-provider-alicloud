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
	resource.AddTestSweepers("alicloud_resource_manager_policy_attachment", &resource.Sweeper{
		Name: "alicloud_resource_manager_policy_attachment",
		F:    testSweepResourceManagerPolicyAttachment,
	})
}

func testSweepResourceManagerPolicyAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}

	action := "ListPolicyAttachments"
	request := make(map[string]interface{})

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["PolicyType"] = "Custom"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	var attachmentIds []string

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory", "EntityNotExist.Policy"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve resoure manager policy attachment in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.PolicyAttachments.PolicyAttachment", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PolicyAttachments.PolicyAttachment", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager policy attachment with policy: %s ", item["PolicyName"].(string))
			} else {
				attachmentIds = append(attachmentIds, fmt.Sprintf("%v:%v:%v:%v:%v", item["PolicyName"], item["PolicyType"], item["PrincipalName"], item["PrincipalType"], item["ResourceGroupId"]))
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, attachmentId := range attachmentIds {
		log.Printf("[INFO] Delete resource manager policy attachment: %s ", attachmentId)

		action := "DetachPolicy"
		ids := strings.Split(attachmentId, ":")
		request := map[string]interface{}{
			"PolicyName":      ids[0],
			"PolicyType":      ids[1],
			"PrincipalName":   ids[2],
			"PrincipalType":   ids[3],
			"ResourceGroupId": ids[4],
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			log.Printf("[ERROR] Failed to delete resource manager policy attachment (%s): %s", attachmentId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerPolicyAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_policy_attachment.default"
	ra := resourceAttrInit(resourceId, ResourceManagerPolicyAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicyAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerPolicyAttachmentBasicdependence)
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
					"policy_name":       "${alicloud_resource_manager_policy.this.policy_name}",
					"policy_type":       "Custom",
					"principal_name":    "${local.principal_name}",
					"principal_type":    "IMSUser",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.this.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":       name,
						"policy_type":       "Custom",
						"principal_name":    CHECKSET,
						"principal_type":    "IMSUser",
						"resource_group_id": CHECKSET,
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

var ResourceManagerPolicyAttachmentMap = map[string]string{}

func ResourceManagerPolicyAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_ram_user" "this" {
  name = "${var.name}"
}

resource "alicloud_resource_manager_policy" "this" {
  policy_name     = "${var.name}"
  description 	  = "policy_attachment"
  policy_document = <<EOF
        {
            "Statement": [{
                "Action": ["oss:*"],
                "Effect": "Allow",
                "Resource": ["acs:oss:*:*:*"]
            }],
            "Version": "1"
        }
    EOF
}

data "alicloud_account" "this" {}

data "alicloud_resource_manager_resource_groups" "this" {
  name_regex = "default"
}

locals{
	principal_name = format("%%s@%%s.onaliyun.com", alicloud_ram_user.this.name, data.alicloud_account.this.id)	
}
`, name)
}
