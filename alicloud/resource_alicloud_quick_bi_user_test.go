package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_quick_bi_user",
		&resource.Sweeper{
			Name: "alicloud_quick_bi_user",
			F:    testSweepQuickBIUser,
		})
}

func testSweepQuickBIUser(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "QueryUserList"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 1

	var response map[string]interface{}
	conn, err := client.NewQuickbiClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-08-01"), StringPointer("AK"), request, nil, &runtime)
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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Result.Data", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Result.Data", action, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["NickName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["NickName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping QuickBI User: %s", item["NickName"].(string))
				continue
			}

			action := "DeleteUser"
			request := map[string]interface{}{
				"UserId": item["UserId"],
			}
			request["ClientToken"] = buildClientToken("DeleteUser")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-08-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete QuickBI User (%s): %s", item["UserId"].(string), err)
			}
			log.Printf("[INFO] Delete QuickBI User success: %s ", item["UserId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return nil
}

func TestAccAlicloudQuickBIUser_basic0(t *testing.T) {
	t.Skip()
	var v map[string]interface{}
	resourceId := "alicloud_quick_bi_user.default"
	ra := resourceAttrInit(resourceId, AlicloudQuickBIUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuickbiPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuickBiUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squickbiuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuickBIUserBasicDependence0)
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
					"nick_name":       name,
					"account_name":    "${local.account_name}",
					"admin_user":      "false",
					"auth_admin_user": "false",
					"user_type":       "Developer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":       name,
						"account_name":    CHECKSET,
						"admin_user":      "false",
						"auth_admin_user": "false",
						"user_type":       "Developer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_user": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"admin_user": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_type": "Analyst",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_type": "Analyst",
					}),
				),
			},
			// todo fix  product problem
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auth_admin_user": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auth_admin_user": "true",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_user":      "false",
					"auth_admin_user": "false",
					"user_type":       "Developer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"admin_user":      "false",
						"auth_admin_user": "false",
						"user_type":       "Developer",
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

var AlicloudQuickBIUserMap0 = map[string]string{}

func AlicloudQuickBIUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_ram_users" "default" {
  name_regex  = "terraform*"
}
locals{
  account_name  = join(":",["%s",data.alicloud_ram_users.default.users.0.name])
}
`, name, "openapiautomation@test.aliyunid.com")
}
