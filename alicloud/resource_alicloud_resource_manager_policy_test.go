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
	resource.AddTestSweepers("alicloud_resource_manager_policy", &resource.Sweeper{
		Name: "alicloud_resource_manager_policy",
		F:    testSweepResourceManagerPolicy,
		Dependencies: []string{
			"alicloud_resource_manager_policy_attachment",
		},
	})
}

func testSweepResourceManagerPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}
	action := "ListPolicies"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["PolicyType"] = "Custom"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	var policyIds []string
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve resoure manager policy in service list: %s", err)
			return nil
		}

		resp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			policyIds = append(policyIds, item["PolicyName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, policyId := range policyIds {
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(policyId), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping resource manager policy: %s ", policyId)
			continue
		}

		// Delete policy version before delete the policy
		action := "ListPolicyVersions"
		versionReq := make(map[string]interface{})
		versionReq["PolicyType"] = "Custom"
		versionReq["PolicyName"] = policyId
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, versionReq, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resource manager policy version (%s): %s", policyId, err)
			continue
		}
		resp, err := jsonpath.Get("$.PolicyVersions.PolicyVersion", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		if len(resp.([]interface{})) > 0 {
			for _, version := range resp.([]interface{}) {
				item := version.(map[string]interface{})
				// default version can not deleted, skip it.
				if v, ok := item["IsDefaultVersion"].(bool); ok && v {
					continue
				}
				log.Printf("[INFO] Delete resource manager policy version: (%s:%s)", policyId, item["VersionId"].(string))

				action := "DeletePolicyVersion"
				delRequest := make(map[string]interface{})
				delRequest["PolicyName"] = policyId
				delRequest["VersionId"] = item["VersionId"]
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(5*time.Minute, func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, delRequest, &util.RuntimeOptions{})
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
					log.Printf("[ERROR] Failed to delete resource manager policy version (%s:%s): %s", policyId, item["VersionId"].(string), err)
				}
			}
		}

		log.Printf("[INFO] Delete resource manager policy: %s ", policyId)

		deleteAction := "DeletePolicy"
		delRequest := map[string]interface{}{
			"PolicyName": policyId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(deleteAction), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, delRequest, &util.RuntimeOptions{})
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
			log.Printf("[ERROR] Failed to delete resource manager policy (%s): %s", policyId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_policy.default"
	ra := resourceAttrInit(resourceId, ResourceManagerPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerPolicyBasicdependence)
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
					"policy_document": `{\n\"Statement\": [{\n\"Action\": [\"oss:*\"],\n\"Effect\": \"Allow\",\n\"Resource\": [\"acs:oss:*:*:*\"]\n}],\n\"Version\": \"1\"\n}`,
					"policy_name":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": CHECKSET,
						"policy_name":     name,
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
					"default_version": "v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_version": "v1",
					}),
				),
			},
		},
	})
}

var ResourceManagerPolicyMap = map[string]string{
	"default_version":  CHECKSET,
	"policy_type":      CHECKSET,
	"attachment_count": CHECKSET,
}

func ResourceManagerPolicyBasicdependence(name string) string {
	return ""
}
