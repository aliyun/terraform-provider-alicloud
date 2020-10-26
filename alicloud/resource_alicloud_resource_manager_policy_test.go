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
	resourceManagerService := ResourcemanagerService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}

	request := resourcemanager.CreateListPoliciesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.PolicyType = "Custom"
	var policyIds []string
	for {
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListPolicies(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager policy in service list: %s", err)
		}

		response, _ := raw.(*resourcemanager.ListPoliciesResponse)

		for _, v := range response.Policies.Policy {
			policyIds = append(policyIds, v.PolicyName)
		}
		if len(response.Policies.Policy) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
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
		versionReq := resourcemanager.CreateListPolicyVersionsRequest()
		versionReq.PolicyType = "Custom"
		versionReq.PolicyName = policyId
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListPolicyVersions(versionReq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resource manager policy version (%s): %s", policyId, err)
			continue
		}
		versions := raw.(*resourcemanager.ListPolicyVersionsResponse)
		if len(versions.PolicyVersions.PolicyVersion) > 0 {
			for _, version := range versions.PolicyVersions.PolicyVersion {
				// default version can not deleted, skip it.
				if version.IsDefaultVersion {
					continue
				}
				log.Printf("[INFO] Delete resource manager policy version: (%s:%s)", policyId, version.VersionId)

				delRequest := resourcemanager.CreateDeletePolicyVersionRequest()
				delRequest.PolicyName = policyId
				delRequest.VersionId = version.VersionId
				_, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
					return resourceManagerClient.DeletePolicyVersion(delRequest)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete resource manager policy version (%s:%s): %s", policyId, version.VersionId, err)
				}
			}
		}

		log.Printf("[INFO] Delete resource manager policy: %s ", policyId)

		request := resourcemanager.CreateDeletePolicyRequest()
		request.PolicyName = policyId

		_, err = resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.DeletePolicy(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager policy (%s): %s", policyId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerPolicy_basic(t *testing.T) {
	var v resourcemanager.Policy
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
	"default_version": CHECKSET,
	"policy_type":     CHECKSET,
}

func ResourceManagerPolicyBasicdependence(name string) string {
	return ""
}
