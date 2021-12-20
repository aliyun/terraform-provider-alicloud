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
		"alicloud_ecd_user",
		&resource.Sweeper{
			Name: "alicloud_ecd_user",
			F:    testSweepEcdUser,
		})
}

func testSweepEcdUser(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeUsers"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := client.NewEdsuserClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &runtime)
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
		resp, err := jsonpath.Get("$.Users", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["EndUserId"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping User: %s", item["EndUserId"].(string))
				continue
			}
			action := "RemoveUsers"
			request := map[string]interface{}{
				"Users": []string{item["EndUserId"].(string)},
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete User (%s): %s", item["EndUserId"].(string), err)
			}
			log.Printf("[INFO] Delete User success: %s ", item["EndUserId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudECDUser_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_user.default"
	ra := resourceAttrInit(resourceId, AlicloudECDUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdsUserService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccecduser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDUserBasicDependence0)
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
					"end_user_id": name,
					"email":       "hello.uuu@aaa.com",
					"phone":       fmt.Sprintf("158016%d", rand),
					"password":    fmt.Sprintf("%d", rand),
					"status":      "Unlocked",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_user_id": name,
						"email":       "hello.uuu@aaa.com",
						"phone":       fmt.Sprintf("158016%d", rand),
						"password":    fmt.Sprintf("%d", rand),
						"status":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Locked",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": CHECKSET,
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

var AlicloudECDUserMap0 = map[string]string{
	"status":      CHECKSET,
	"end_user_id": CHECKSET,
}

func AlicloudECDUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
