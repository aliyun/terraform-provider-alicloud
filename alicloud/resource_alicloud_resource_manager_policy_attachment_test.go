package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
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
	resourceManagerService := ResourcemanagerService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}

	request := resourcemanager.CreateListPolicyAttachmentsRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.PolicyType = "Custom"
	var attachmentIds []string
	for {
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListPolicyAttachments(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager policy attachment in service list: %s", err)
		}

		response, _ := raw.(*resourcemanager.ListPolicyAttachmentsResponse)

		for _, v := range response.PolicyAttachments.PolicyAttachment {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(v.PolicyName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager policy attachment with policy: %s ", v.PolicyName)
			} else {
				attachmentIds = append(attachmentIds, fmt.Sprintf("%v:%v:%v:%v:%v", v.PolicyName, v.PolicyType, v.PrincipalName, v.PrincipalType, v.ResourceGroupId))
			}
		}
		if len(response.PolicyAttachments.PolicyAttachment) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, attachmentId := range attachmentIds {
		log.Printf("[INFO] Delete resource manager policy attachment: %s ", attachmentId)

		ids := strings.Split(attachmentId, ":")
		request := resourcemanager.CreateDetachPolicyRequest()
		request.PolicyName = ids[0]
		request.PolicyType = ids[1]
		request.PrincipalName = ids[2]
		request.PrincipalType = ids[3]
		request.ResourceGroupId = ids[4]

		_, err = resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.DetachPolicy(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager policy attachment (%s): %s", attachmentId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerPolicyAttachment_basic(t *testing.T) {
	var v resourcemanager.PolicyAttachment
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
