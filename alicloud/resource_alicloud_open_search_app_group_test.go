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
	resource.AddTestSweepers("alicloud_open_search_app_group", &resource.Sweeper{
		Name: "alicloud_open_search_app_group",
		F:    testSweepOpenSearchAppGroup,
	})
}

func testSweepOpenSearchAppGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	conn, err := client.NewOpensearchClient()
	if err != nil {
		return WrapError(err)
	}

	action := "/v4/openapi/app-groups"
	request := make(map[string]*string)
	ids := make([]string, 0)
	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List AppGroup Failed!", err)
			return nil
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
		}
		resp, err := jsonpath.Get("$.result", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AppGroup.result", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if convertOpenSearchAppGroupPaymentTypeResponse(item["chargeType"].(string)) == "Subscription" {
				log.Printf("[INFO] Skipping AppGroup: %v (%v)", item["name"], item["id"])
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["name"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping AppGroup: %v (%v)", item["name"], item["id"])
				continue
			}
			ids = append(ids, fmt.Sprint(item["name"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
	}
	for _, id := range ids {
		log.Printf("[INFO] Deleting App Group: (%s)", id)
		action := "/v4/openapi/app-groups/" + id
		conn, err := client.NewOpensearchClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]*string{
			"appGroupIdentity": StringPointer(id),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*9, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed To Delete AppGroup : %s", err)
		}
		log.Printf("[INFO] Delete AppGroup Success : %s", id)
	}
	return nil
}

func TestAccAlicloudOpenSearchAppGroup_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.OpenSearchSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_open_search_app_group.default"
	ra := resourceAttrInit(resourceId, AlicloudOpenSearchAppGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OpenSearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOpenSearchAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testaccosappgroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOpenSearchAppGroupBasicDependence0)
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
					"app_group_name": name,
					"payment_type":   "PayAsYouGo",
					"type":           "standard",
					"quota": []map[string]interface{}{
						{
							"doc_size":         "10",
							"compute_resource": "20",
							"spec":             "opensearch.share.common",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_group_name": name,
						"type":           "standard",
						"payment_type":   "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"charge_way": "compute_resource",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"charge_way": "compute_resource",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_type": "UPGRADE",
					"quota": []map[string]interface{}{
						{
							"doc_size":         "20",
							"compute_resource": "1000",
							"spec":             "opensearch.share.compute",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "quota"},
			},
		},
	})
}

// There is an api error： InternalError. Reopen it after the error has been fixed.
func SkipTestAccAlicloudOpenSearchAppGroup_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.OpenSearchSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_open_search_app_group.default"
	ra := resourceAttrInit(resourceId, AlicloudOpenSearchAppGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OpenSearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOpenSearchAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testaccosappgroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOpenSearchAppGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_group_name": name,
					"payment_type":   "Subscription",
					"type":           "standard",
					"quota": []map[string]interface{}{
						{
							"doc_size":         "10",
							"compute_resource": "20",
							"spec":             "opensearch.share.common",
						},
					},
					"order": []map[string]interface{}{
						{
							"duration":      "1",
							"pricing_cycle": "Month",
							"auto_renew":    "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_group_name": name,
						"type":           "standard",
						"payment_type":   "Subscription",
						"quota.#":        "1",
						"order.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"charge_way": "compute_resource",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"charge_way": "compute_resource",
					}),
				),
			},
			// todo： Reopen after OpenSearch fixing the issue
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"order_type": "UPGRADE",
			//		"quota": []map[string]interface{}{
			//			{
			//				"doc_size":         "2",
			//				"compute_resource": "1000",
			//				"spec":             "opensearch.share.compute",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"quota.#": "1",
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "quota", "order"},
			},
		},
	})
}

var AlicloudOpenSearchAppGroupMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AlicloudOpenSearchAppGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
