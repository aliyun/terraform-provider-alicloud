// Package tag_test.
package tag_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var tagTestEndpoints sync.Map
var tagTestSignVersion sync.Map

// tagTestClientForRegion rebuilds the package-level sweeper client on the
// exported connectivity API: the flat-suite sharedClientForRegion is not
// reachable from this package.
func tagTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &tagTestEndpoints,
		SignVersion: &tagTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	return conf.Client()
}

func init() {
	resource.AddTestSweepers("alicloud_tag_policy_attachment", &resource.Sweeper{
		Name: "alicloud_tag_policy_attachment",
		F:    testSweepTagPolicyAttachment,
	})
}

func testSweepTagPolicyAttachment(region string) error {
	rawClient, err := tagTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient

	action := "ListPoliciesForTarget"
	request := make(map[string]interface{})
	request1 := make(map[string]interface{})

	request["MaxResult"] = helper.PageSizeLarge
	var response map[string]interface{}
	var attachmentIds []string

	for {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExists.Target", "EntityNotExist.Policy"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve tag policy attachment in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			request1["PolicyId"] = item["PolicyId"]
			actionListTargetsForPolicy := "ListTargetsForPolicy"
			response, err = client.RpcPost("Tag", "2018-08-28", actionListTargetsForPolicy, nil, request1, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"EntityNotExists.Target", "EntityNotExist.Policy"}) {
					return nil
				}
				log.Printf("[ERROR] Failed to retrieve tag policy attachment in service list: %s", err)
				return nil
			}
			resp1, err := jsonpath.Get("$.Targets", response)
			if err != nil {
				return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data", response)
			}
			result1, _ := resp1.([]interface{})
			for _, w := range result1 {
				item1 := w.(map[string]interface{})
				attachmentIds = append(attachmentIds, fmt.Sprintf("%v:%v:%v", item["PolicyId"], item1["TargetId"], item1["TargetType"]))
			}

		}
		if len(result) < helper.PageSizeLarge {
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

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete tag policy attachment (%s): %s", attachmentId, err)
		}
	}
	return nil
}

func TestAccAliCloudTagPolicyAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_policy_attachment.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.TagPolicySupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, TagPolicyAttachmentMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return tag.NewTagService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTagPolicyAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicyAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, TagPolicyAttachmentBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_id":   "${alicloud_tag_policy.this.id}",
					"target_id":   "${data.alicloud_resource_manager_resource_directories.default.directories[0].root_folder_id}",
					"target_type": "ROOT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_id":   acctest.CHECKSET,
						"target_id":   acctest.CHECKSET,
						"target_type": "ROOT",
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

data "alicloud_resource_manager_resource_directories" "default" {}

resource "alicloud_tag_policy" "this" {
  policy_name     = "${var.name}"
  policy_desc 	  = "policy_attachment"
  user_type       = "RD"
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
