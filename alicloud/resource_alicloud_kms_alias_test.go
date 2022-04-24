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
	resource.AddTestSweepers("alicloud_kms_alias", &resource.Sweeper{
		Name: "alicloud_kms_alias",
		F:    testSweepKmsAlias,
	})
}

func testSweepKmsAlias(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"alias/tf-testacc",
		"alias/tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "ListAliases"

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	sweeped := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_alias", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Aliases.Alias", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Aliases.Alias", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["AliasName"]; !ok {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AliasName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Kms Key alias: %s", item["AliasName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteAlias"
			request := map[string]interface{}{
				"AliasName": item["AliasName"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Alias (%s): %s", item["AliasName"], err)
			}
			log.Printf("[INFO] Delete Kms Key ALias success: %s ", item["AliasName"])
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudKMSAlias_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_kms_alias.default"
	ra := resourceAttrInit(resourceId, kmsAliasBasicMap)

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("alias/tf_testaccKmsAlias_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsAliasConfigDependence)

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
					"alias_name": name,
					"key_id":     "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name,
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
					"key_id": "${alicloud_kms_key.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceKmsAliasConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
	description = "tf-testacckmskeyforaliasdefault"
	pending_window_in_days = 7
}

resource "alicloud_kms_key" "default1" {
	description = "tf-testacckmskeyforaliasdefault1"
	pending_window_in_days = 7
}
`)
}

var kmsAliasBasicMap = map[string]string{
	"key_id": CHECKSET,
}
