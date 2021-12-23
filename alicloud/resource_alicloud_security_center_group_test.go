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
		"alicloud_security_center_group",
		&resource.Sweeper{
			Name: "alicloud_security_center_group",
			F:    testSweepSasGroup,
		})
}

func testSweepSasGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeAllGroups"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewSasClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &runtime)
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
		log.Printf("[ERROR] Failed to fetch Sas Groups: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_security_center_groups", action, AlibabaCloudSdkGoERROR))
		return nil
	}
	resp, err := jsonpath.Get("$.Groups", response)
	if err != nil {
		log.Printf("[ERROR] Failed to parse Sas Groups: %s", WrapErrorf(err, FailedGetAttributeMsg, action, "$.Groups", response))
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["GroupName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Sas Group: %s", item["GroupName"].(string))
			continue
		}

		action := "DeleteGroup"
		request := map[string]interface{}{
			"GroupId": item["GroupId"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Sas Group (%s): %s", item["GroupName"].(string), err)
		}
		log.Printf("[INFO] Delete Sas Group success: %s ", item["GroupName"].(string))

	}

	return nil
}

func TestAccAlicloudSASGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_center_group.default"
	ra := resourceAttrInit(resourceId, AlicloudSASGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAllGroups")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssasgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSASGroupBasicDependence0)
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
					"group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "_update",
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

var AlicloudSASGroupMap0 = map[string]string{}

func AlicloudSASGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
