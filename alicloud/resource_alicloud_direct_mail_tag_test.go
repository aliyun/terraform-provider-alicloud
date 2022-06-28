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
		"alicloud_direct_mail_tag",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_tag",
			F:    testSweepDirectMailTag,
		})
}

func testSweepDirectMailTag(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DmSupportRegions) {
		log.Printf("[INFO] Skipping Direct Mail Tag unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tfTestAcc",
	}
	action := "QueryTagByParam"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNo"] = 1

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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.data.tag", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.data.tag", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TagName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Direct Mail Tag: %s", item["TagName"].(string))
				continue
			}
			action := "DeleteTag"
			request := map[string]interface{}{
				"TagId": item["TagId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Direct Mail Tag (%s): %s", item["TagName"].(string), err)
			}
			log.Printf("[INFO] Delete Direct Mail Tag success: %s ", item["TagName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDirectMailTag_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.DmSupportRegions)
	resourceId := "alicloud_direct_mail_tag.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailTagMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailTag")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccdirectmailtag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailTagBasicDependence0)
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
					"tag_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_name": "${var.name}update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_name": name,
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

var AlicloudDirectMailTagMap0 = map[string]string{}

func AlicloudDirectMailTagBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
