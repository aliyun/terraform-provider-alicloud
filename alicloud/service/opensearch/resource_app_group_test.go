package opensearch_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/opensearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_open_search_app_group", &resource.Sweeper{
		Name: "alicloud_open_search_app_group",
		F:    testSweepOpenSearchAppGroup,
	})
}

func testSweepOpenSearchAppGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "/v4/openapi/app-groups"
	request := make(map[string]*string)
	ids := make([]string, 0)
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("OpenSearch", "2017-12-25", action, request, nil, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List AppGroup Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.result", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AppGroup.result", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			// chargeType "PREPAY" maps to payment_type "Subscription"; skip subscription app groups.
			if item["chargeType"].(string) == "PREPAY" {
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
		if len(result) < helper.PageSizeLarge {
			break
		}
	}
	for _, id := range ids {
		log.Printf("[INFO] Deleting App Group: (%s)", id)
		action := "/v4/openapi/app-groups/" + id
		request := map[string]*string{
			"appGroupIdentity": helper.StringPointer(id),
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*9, func() *retry.RetryError {
			response, err = client.RoaDelete("OpenSearch", "2017-12-25", action, request, nil, nil, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

func TestAccAliCloudOpenSearchAppGroup_basic0(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.OpenSearchSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_open_search_app_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudOpenSearchAppGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return opensearch.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAppGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testaccosappgroup%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudOpenSearchAppGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"compute_resource": "40",
							"spec":             "opensearch.share.common",
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
func TestAccAliCloudOpenSearchAppGroup_basic1(t *testing.T) {
	t.Skip("InternalError. Reopen it after the error has been fixed.")
	acctest.CheckoutSupportedRegions(t, true, connectivity.OpenSearchSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_open_search_app_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudOpenSearchAppGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return opensearch.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAppGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testaccosappgroup%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudOpenSearchAppGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

var AliCloudOpenSearchAppGroupMap0 = map[string]string{
	"instance_id": acctest.CHECKSET,
}

func AliCloudOpenSearchAppGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
