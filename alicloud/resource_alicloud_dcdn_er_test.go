package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_er",
		&resource.Sweeper{
			Name: "alicloud_dcdn_er",
			F:    testSweepDcdnEr,
		})
}

func testSweepDcdnEr(region string) error {
	// DCDN ER (Edge Routine) is a global/account-scoped service — no region gate
	// (cn-beijing must run).
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// RED LINE #5: prefix guard always applied (even under SWEEP_ALL=1).
	prefixes := []string{
		"tf-",
		"tf-testacc",
		"tf_testacc",
		"test-",
		"testacc",
		"test_acc",
	}

	action := "DescribeRoutineUserInfo"
	request := map[string]interface{}{}

	var response map[string]interface{}
	listWait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				listWait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		// ER service may not be activated in the sweep account — non-fatal.
		log.Printf("[ERROR] %s got an error: %v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.Content.Routines", response)
	if err != nil {
		log.Printf("[ERROR] Getting %s attribute by path %s failed: %v", action, "$.Content.Routines", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["RoutineName"])

		// RED LINE #2: NODELETING / do-not-delete marker protection.
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping Dcdn Er (NODELETING): %s", name)
			continue
		}
		// RED LINE #5: prefix guard always applied.
		matches := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), prefix) {
				matches = true
				break
			}
		}
		if !matches {
			log.Printf("[INFO] Skipping Dcdn Er (no test prefix): %s", name)
			continue
		}

		deleteAction := "DeleteRoutine"
		deleteRequest := map[string]interface{}{
			"Name": name,
		}
		log.Printf("[INFO] Deleting Dcdn Er: %s", name)
		deleteWait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", deleteAction, nil, deleteRequest, false)
			if err != nil {
				if NeedRetry(err) {
					deleteWait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(deleteAction, response, deleteRequest)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Dcdn Er (%s): %v", name, err)
			continue
		}
		log.Printf("[INFO] Deleted Dcdn Er: %s", name)
	}

	return nil
}

func TestAccAliCloudDcdnEr_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
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

var AlicloudDcdnErMap = map[string]string{}

func AlicloudDcdnErBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
`, name)
}
