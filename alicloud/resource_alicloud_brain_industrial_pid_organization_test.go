package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_brain_industrial_pid_organization",
		&resource.Sweeper{
			Name: "alicloud_brain_industrial_pid_organization",
			F:    testSweepBrainIndustrialPidOrganization,
		})
}

func testSweepBrainIndustrialPidOrganization(region string) error {
	if !testSweepPreCheckWithRegions(region, false, connectivity.BrainIndustrialRegions) {
		log.Printf("[INFO] Skipping Brain Industrial unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{}
	var response map[string]interface{}
	action := "ListPidOrganizations"
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_brain_industrial_pid_organization", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.OrganizationList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.OrganizationList", response)
	}
	sweeped := false
	for _, v := range resp.([]interface{}) {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["OrganizationName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Brain Industrial Organization: %s", item["OrganizationName"].(string))
			continue
		}
		sweeped = true
		action = "DeletePidOrganization"
		request := map[string]interface{}{
			"OrganizationId": item["OrganizationId"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Brain Industrial Organization (%s): %s", item["OrganizationName"].(string), err)
		}
		if sweeped {
			// Waiting 5 seconds to ensure  Brain Industrial Organization have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete Brain Industrial Organization success: %s ", item["OrganizationName"].(string))
	}
	return nil

}

func TestAccAlicloudBrainIndustrialPidOrganization_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_organization.default"
	ra := resourceAttrInit(resourceId, AlicloudBrainIndustrialPidOrganizationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Brain_industrialService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBrainIndustrialPidOrganization")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBrainIndustrialPidOrganizationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_pid_organization_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_name": "tf-testAccUp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": "tf-testAccUp",
					}),
				),
			},
		},
	})
}

var AlicloudBrainIndustrialPidOrganizationMap = map[string]string{}

func AlicloudBrainIndustrialPidOrganizationBasicDependence(name string) string {
	return ""
}
