package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_direct_mail_config_set",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_config_set",
			F:    testSweepDirectMailConfigSet,
		})
}

func testSweepDirectMailConfigSet(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ConfigSetList"
	request := map[string]interface{}{
		"PageIndex": "1",
		"PageSize":  fmt.Sprint(PageSizeLarge),
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
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
			log.Printf("[ERROR] Failed to fetch DirectMail Config Set: %s", err)
			return nil
		}
		v, err := jsonpath.Get("$.ConfigSets", response)
		if err != nil {
			return nil
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DirectMail Config Set: %s", name)
				continue
			}
			action := "ConfigSetDelete"
			deleteRequest := map[string]interface{}{
				"Ids": item["Id"],
			}
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, deleteRequest, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DirectMail Config Set (%s): %s", name, err)
			} else {
				log.Printf("[INFO] Delete DirectMail Config Set success: %s", name)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageIndex"] = fmt.Sprint(request["PageIndex"].(int) + 1)
	}
	return nil
}

func TestAccAliCloudDirectMailConfigSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailConfigSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailConfigSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailConfigSetBasicDependence0)
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
					"name":                      name,
					"description":               "tf-testAcc-description",
					"is_public_channel_backoff": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"description":               "tf-testAcc-description",
						"is_public_channel_backoff": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name + "-update",
						"description":               "tf-testAcc-description-update",
						"is_public_channel_backoff": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
					"validation_option": []map[string]interface{}{
						{
							"enabled": "true",
							"forbidden_status_list": []string{
								"INVALID",
								"DO_NOT_MAIL",
							},
							"forbidden_sub_status_list": []string{
								"DISPOSABLE",
								"MAILBOX_FULL",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"validation_option.#":                         "1",
						"validation_option.0.enabled":                 "true",
						"validation_option.0.forbidden_status_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
					"validation_option": []map[string]interface{}{
						{
							"enabled": "true",
							"forbidden_status_list": []string{
								"DO_NOT_MAIL",
								"INVALID",
							},
							"forbidden_sub_status_list": []string{
								"DISPOSABLE",
								"MAILBOX_FULL",
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
					"validation_option": []map[string]interface{}{
						{
							"enabled": "true",
							"forbidden_status_list": []string{
								"DO_NOT_MAIL",
								"INVALID",
							},
							"forbidden_sub_status_list": []string{
								"DISPOSABLE",
								"MAILBOX_FULL",
							},
						},
					},
				}),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
					"validation_option": []map[string]interface{}{
						{
							"enabled": "true",
							"forbidden_status_list": []string{
								"DO_NOT_MAIL",
								"INVALID",
							},
							"forbidden_sub_status_list": []string{
								"MAILBOX_FULL",
								"DISPOSABLE",
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-update",
					"is_public_channel_backoff": "false",
					"validation_option": []map[string]interface{}{
						{
							"enabled": "true",
							"forbidden_status_list": []string{
								"DO_NOT_MAIL",
								"INVALID",
							},
							"forbidden_sub_status_list": []string{
								"MAILBOX_FULL",
								"DISPOSABLE",
							},
						},
					},
				}),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "tf-testAcc-description-update",
					"is_force":    "true",
					"ip_pool_id":  "CHECKSET",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf-testAcc-description-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_force"},
			},
		},
	})
}

func TestAccAliCloudDirectMailConfigSetsDataSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailConfigSet")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.name", name),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_config_sets.default", "ids.#"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  keyword = alicloud_direct_mail_config_set.default.name
}
`, name)
}

// Associating a configuration set with a dedicated IP pool requires purchased
// dedicated IP instances in the account (pool creation is rejected otherwise);
// the shared tf acceptance test account owns none, so this case is skipped the
// same way as other DirectMail tests.
func SkipTestAccAlicloudDirectMailConfigSet_withIpPool(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailConfigSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailConfigSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailConfigSetWithIpPoolDependence0)
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
					"name":       name,
					"ip_pool_id": "${alicloud_direct_mail_dedicated_ip_pool.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"ip_pool_id":   CHECKSET,
						"ip_pool_name": CHECKSET,
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

func AlicloudDirectMailConfigSetWithIpPoolDependence0(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_dedicated_ip_pool" "default" {
  name = "%s-pool"
}
`, name)
}

var AlicloudDirectMailConfigSetMap0 = map[string]string{}

func AlicloudDirectMailConfigSetBasicDependence0(name string) string {
	return ""
}
