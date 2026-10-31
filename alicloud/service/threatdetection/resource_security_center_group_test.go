package threatdetection_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers(
		"alicloud_security_center_group",
		&resource.Sweeper{
			Name: "alicloud_security_center_group",
			F:    testSweepSasGroup,
		})
}

func testSweepSasGroup(region string) error {
	rawClient, err := threatdetectionTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeAllGroups"
	request := make(map[string]interface{})
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch Sas Groups: %s", errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_security_center_groups", action, errs.AlibabaCloudSdkGoERROR))
		return nil
	}
	resp, err := jsonpath.Get("$.Groups", response)
	if err != nil {
		log.Printf("[ERROR] Failed to parse Sas Groups: %s", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Groups", response))
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
		_, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Sas Group (%s): %s", item["GroupName"].(string), err)
		}
		log.Printf("[INFO] Delete Sas Group success: %s ", item["GroupName"].(string))

	}

	return nil
}

func TestAccAliCloudSASGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_center_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSASGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAllGroups")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssasgroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSASGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

// lintignore: R001
