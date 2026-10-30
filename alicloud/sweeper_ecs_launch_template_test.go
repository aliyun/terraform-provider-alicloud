package alicloud

// The alicloud_ecs_launch_template sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, addDebug, PageSizeLarge) live in the root alicloud
// test package. The resource implementation itself lives in alicloud/service/ecs.

import (
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_launch_template", &resource.Sweeper{
		Name: "alicloud_ecs_launch_template",
		F:    testAliCloudEcsLaunchTemplate,
	})
}

func testAliCloudEcsLaunchTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "DescribeLaunchTemplates"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_launch_templates", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.LaunchTemplateSets.LaunchTemplateSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LaunchTemplateSets.LaunchTemplateSet", response)
		}

		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["LaunchTemplateName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Launch Template: %s", item["LaunchTemplateName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteLaunchTemplateVersion"
			request := map[string]interface{}{
				"LaunchTemplateId": item["LaunchTemplateId"],
				"RegionId":         client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Launch Template (%s): %s", item["LaunchTemplateName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros Template have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Launch Template success: %s ", item["LaunchTemplateName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
