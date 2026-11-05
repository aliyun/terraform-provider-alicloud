package kms_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_alias", &resource.Sweeper{
		Name: "alicloud_kms_alias",
		F:    testSweepKmsAlias,
	})
}

func testSweepKmsAlias(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"alias/tf-testacc",
		"alias/tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "ListAliases"

	var response map[string]interface{}
	sweeped := false
	for {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_kms_alias", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Aliases.Alias", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Aliases.Alias", response)
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
			_, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Alias (%s): %s", item["AliasName"], err)
			}
			log.Printf("[INFO] Delete Kms Key ALias success: %s ", item["AliasName"])
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAliCloudKmsAlias_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_kms_alias.default"
	ra := acctest.ResourceAttrInit(resourceId, kmsAliasBasicMap)

	serviceFunc := func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}

	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("alias/tf_testaccKmsAlias_%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceKmsAliasConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"key_id": acctest.CHECKSET,
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
	"key_id": acctest.CHECKSET,
}

// lintignore: R001
