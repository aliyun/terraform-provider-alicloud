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
	resource.AddTestSweepers("alicloud_dms_enterprise_user", &resource.Sweeper{
		Name: "alicloud_dms_enterprise_user",
		F:    testSweepDMSEnterpriseUsers,
	})
}

func testSweepDMSEnterpriseUsers(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "ListUsers"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_users", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.UserList.User", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.UserList.User", response)
		}

		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["NickName"]; !ok {
				skip = false
			} else {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["NickName"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DMS Enterprise User: %s", item["NickName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteUser"
			request := map[string]interface{}{
				"Uid": item["Uid"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise User (%s): %s", item["NickName"].(string), err)
				continue
			}
			if sweeped {
				// Waiting 30 seconds to ensure these DMS Enterprise User have been deleted.
				time.Sleep(30 * time.Second)
			}
			log.Printf("[INFO] Delete DMS Enterprise User Success: %s ", item["NickName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDMSEnterpriseUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_user.default"
	ra := resourceAttrInit(resourceId, DmsEnterpriseUserMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsEnterpriseUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseUser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DmsEnterpriseUserBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"uid":        "${alicloud_ram_user.user.id}",
					"nick_name":  name,
					"role_names": []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":    name,
						"role_names.#": "1",
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
					"max_execute_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_result_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_result_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nick_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_names": []string{"USER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "NORMAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_execute_count": "1000",
					"max_result_count":  "1000",
					"nick_name":         name + "change",
					"role_names":        []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "1000",
						"max_result_count":  "1000",
						"nick_name":         name + "change",
						"role_names.#":      "1",
					}),
				),
			},
		},
	})
}

var DmsEnterpriseUserMap = map[string]string{
	"status": CHECKSET,
}

func DmsEnterpriseUserBasicdependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name         = "%s"
	  display_name = "user_display_name"
	  mobile       = "86-18688888888"
	  email        = "hello.uuu@aaa.com"
	  comments     = "yoyoyo"
	}`, name)
}
