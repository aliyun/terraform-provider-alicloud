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
		"alicloud_direct_mail_mail_address",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_mail_address",
			F:    testSweepDirectMailAddress,
		})
}

func testSweepDirectMailAddress(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "QueryMailAddressByParam"
	request := map[string]interface{}{
		"PageSize": PageSizeLarge,
	}
	var response map[string]interface{}
	conn, err := client.NewDmClient()
	if err != nil {
		log.Println(WrapError(err))
		return nil
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] Failed to fetch DirectMail Mail Address: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_direct_mail_mail_addresss", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		v, err := jsonpath.Get("$.data.mailAddress", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse DirectMail Mail Address: %s", WrapErrorf(err, FailedGetAttributeMsg, action, "$.Groups", response))
			return nil
		}
		if len(v.([]interface{})) < 1 {
			log.Printf("[ERROR] Failed to fetch DirectMail Mail Address: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_direct_mail_mail_addresss", action, AlibabaCloudSdkGoERROR))
			return nil
		}

		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AccountName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DirectMail Mail Address: %s", item["AccountName"].(string))
				continue
			}

			action := "DeleteMailAddress"
			request := map[string]interface{}{
				"MailAddressId": item["MailAddressId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DirectMail Mail Address (%s): %s", item["AccountName"].(string), err)
			}
			log.Printf("[INFO] Delete DirectMail Mail Address success: %s ", item["AccountName"].(string))
		}

		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil

}

// this resource depends on user's account which tf account not support
func SkipTestAccAlicloudDirectMailMailAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_mail_address.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailMailAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailMailAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdirectmailmailaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailMailAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sendtype":     "batch",
					"account_name": name + "s@xxx.changes.com.cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sendtype":     "batch",
						"account_name": name + "s@xxx.changes.com.cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reply_address": "r@xxx.changes.com.cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reply_address": "r@xxx.changes.com.cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Tf12345678password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Tf12345678password",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reply_address": "r2@xxx.changes.com.cn",
					"password":      "Tf987654321password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reply_address": "r2@xxx.changes.com.cn",
						"password":      "Tf987654321password",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudDirectMailMailAddressMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudDirectMailMailAddressBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
