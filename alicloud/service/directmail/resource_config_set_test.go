package directmail_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/directmail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := sweep.SharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "ConfigSetList"
	request := map[string]interface{}{
		"PageIndex": 1,
		"PageSize":  helper.PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			if errs.NotFoundError(err) {
				log.Printf("[INFO] DirectMail Config Set list reported not found; nothing to sweep: %s", err)
				return nil
			}
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
			if !directMailNameShouldSweep(name) {
				log.Printf("[INFO] Skipping DirectMail Config Set: %s", name)
				continue
			}
			action := "ConfigSetDelete"
			deleteRequest := map[string]interface{}{
				"Ids": item["Id"],
			}
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, deleteRequest, true)
			if err != nil {
				if errs.NotFoundError(err) {
					log.Printf("[INFO] DirectMail Config Set %s was already removed: %s", name, err)
				} else {
					log.Printf("[ERROR] Failed to delete DirectMail Config Set (%s): %s", name, err)
				}
			} else {
				log.Printf("[INFO] Delete DirectMail Config Set success: %s", name)
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageIndex"] = request["PageIndex"].(int) + 1
	}
	return nil
}

func TestAccAliCloudDirectMailConfigSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailConfigSetMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDirectMailServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailConfigSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailConfigSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				// A flat-field-only change (validation_option unchanged in
				// config) must still resend validation_option, otherwise the
				// full-replace update would silently clear it server-side.
				Config: testAccConfig(map[string]interface{}{
					"name":                      name + "-update",
					"description":               "tf-testAcc-description-kept",
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "tf-testAcc-description-kept",
						"validation_option.0.enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"description":       "tf-testAcc-description-update",
					"is_force":          "true",
					"validation_option": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"description":                 "tf-testAcc-description-update",
						"validation_option.#":         "0",
						"validation_option.0.enabled": acctest.REMOVEKEY,
						"validation_option.0.forbidden_status_list.#": acctest.REMOVEKEY,
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

// Associating a configuration set with a dedicated IP pool requires purchased
// dedicated IP instances in the account (pool creation is rejected otherwise);
// the shared tf acceptance test account owns none, so this case skips at
// runtime unless ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN is set. Removing
// ip_pool_id from the configuration exercises the full-replace update that
// clears the pool association.
func TestAccAliCloudDirectMailConfigSet_withIpPool(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs to create the associated pool; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailConfigSetMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDirectMailServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailConfigSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailConfigSetWithIpPoolDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		// No IDRefreshName here: the SDK v2 ID-only refresh rebuilds a minimal
		// provider config, which would drop the assume_role block and query the
		// base test account instead of the target role account.
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       name,
					"ip_pool_id": "${alicloud_direct_mail_dedicated_ip_pool.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"ip_pool_id":   acctest.CHECKSET,
						"ip_pool_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       name,
					"ip_pool_id": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"ip_pool_id":   "",
						"ip_pool_name": "",
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
%s

resource "alicloud_direct_mail_dedicated_ip_pool" "default" {
  name = "%s-pool"
}
`, testAccDirectMailAssumeRoleProviderConfig(), name)
}

var AlicloudDirectMailConfigSetMap0 = map[string]string{}

func AlicloudDirectMailConfigSetBasicDependence0(name string) string {
	return ""
}

// directMailConfigSetPropagationClient serves ConfigSetDetail with the stale
// detail first and the requested detail afterwards, counting the calls.
func directMailConfigSetPropagationClient(t *testing.T, stale, fresh map[string]interface{}, calls *int32) *connectivity.AliyunClient {
	return directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.Form.Get("Action"); action != "ConfigSetDetail" {
			t.Errorf("unexpected action %s", action)
		}
		detail := fresh
		if atomic.AddInt32(calls, 1) == 1 {
			detail = stale
		}
		directMailWriteJSON(t, w, map[string]interface{}{"Detail": detail})
	})
}

func TestUnitDirectMailConfigSetWaitForValidationOption(t *testing.T) {
	var calls int32
	stale := map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "desc"}
	fresh := map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "desc", "ValidationOption": map[string]interface{}{
		"Enabled":             true,
		"ForbiddenStatusList": []interface{}{"DO_NOT_MAIL", "INVALID"},
	}}
	client := directMailConfigSetPropagationClient(t, stale, fresh, &calls)
	d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{
		"name":        "cs",
		"description": "desc",
		"validation_option": []interface{}{map[string]interface{}{
			"enabled":               true,
			"forbidden_status_list": []interface{}{"DO_NOT_MAIL", "INVALID"},
		}},
	})
	d.SetId("cs-1")
	service := directmail.NewDirectMailServiceV2(client)
	if err := directmail.WaitForDirectMailConfigSetUpdatePropagated(d, service); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("the waiter must not accept a detail with stale validation settings; calls=%d", calls)
	}
}

func TestUnitDirectMailConfigSetWaitForValidationListOrder(t *testing.T) {
	for _, field := range []string{"ForbiddenStatusList", "ForbiddenSubStatusList"} {
		t.Run(field, func(t *testing.T) {
			var calls int32
			configField := "forbidden_status_list"
			if field == "ForbiddenSubStatusList" {
				configField = "forbidden_sub_status_list"
			}
			stale := map[string]interface{}{"Id": "cs-1", "Name": "cs", "ValidationOption": map[string]interface{}{
				"Enabled": true, field: []interface{}{"FIRST", "SECOND"},
			}}
			fresh := map[string]interface{}{"Id": "cs-1", "Name": "cs", "ValidationOption": map[string]interface{}{
				"Enabled": true, field: []interface{}{"SECOND", "FIRST"},
			}}
			client := directMailConfigSetPropagationClient(t, stale, fresh, &calls)
			d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{
				"name": "cs",
				"validation_option": []interface{}{map[string]interface{}{
					"enabled": true, configField: []interface{}{"SECOND", "FIRST"},
				}},
			})
			d.SetId("cs-1")
			if err := directmail.WaitForDirectMailConfigSetUpdatePropagated(d, directmail.NewDirectMailServiceV2(client)); err != nil {
				t.Fatal(err)
			}
			if got := atomic.LoadInt32(&calls); got < 2 {
				t.Fatalf("the waiter accepted stale %s order after %d call(s)", field, got)
			}
		})
	}
}

func TestUnitDirectMailConfigSetWaitForIpPoolAssociation(t *testing.T) {
	var calls int32
	stale := map[string]interface{}{"Id": "cs-1", "Name": "cs"}
	fresh := map[string]interface{}{"Id": "cs-1", "Name": "cs", "IpPool": map[string]interface{}{"IpPoolId": "pool-1", "IpPoolName": "pool"}}
	client := directMailConfigSetPropagationClient(t, stale, fresh, &calls)
	d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{
		"name":       "cs",
		"ip_pool_id": "pool-1",
	})
	d.SetId("cs-1")
	service := directmail.NewDirectMailServiceV2(client)
	if err := directmail.WaitForDirectMailConfigSetUpdatePropagated(d, service); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("the waiter must not accept a detail without the requested IP pool; calls=%d", calls)
	}
}

func TestUnitDirectMailConfigSetWaitForIpPoolRemoval(t *testing.T) {
	var calls int32
	stale := map[string]interface{}{"Id": "cs-1", "Name": "cs", "IpPool": map[string]interface{}{"IpPoolId": "pool-1"}}
	fresh := map[string]interface{}{"Id": "cs-1", "Name": "cs"}
	client := directMailConfigSetPropagationClient(t, stale, fresh, &calls)
	d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{"name": "cs"})
	d.SetId("cs-1")
	service := directmail.NewDirectMailServiceV2(client)
	if err := directmail.WaitForDirectMailConfigSetUpdatePropagated(d, service); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("the waiter must not accept a detail that still carries the removed IP pool; calls=%d", calls)
	}
}

func TestUnitDirectMailConfigSetReadDisabledEmptyValidation(t *testing.T) {
	for _, explicitBlock := range []bool{false, true} {
		client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			directMailWriteJSON(t, w, map[string]interface{}{"Detail": map[string]interface{}{
				"Id": "cs-1", "Name": "cs",
				"ValidationOption": map[string]interface{}{"Enabled": false, "ForbiddenStatusList": []interface{}{}, "ForbiddenSubStatusList": []interface{}{}},
			}})
		})
		config := map[string]interface{}{"name": "cs"}
		wantCount := 0
		if explicitBlock {
			config["validation_option"] = []interface{}{map[string]interface{}{"enabled": false}}
			wantCount = 1
		}
		d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, config)
		d.SetId("cs-1")
		if err := directmail.ResourceConfigSetRead(d, client); err != nil {
			t.Fatal(err)
		}
		if got := len(d.Get("validation_option").([]interface{})); got != wantCount {
			t.Fatalf("explicitBlock=%t: expected %d validation blocks, got %d", explicitBlock, wantCount, got)
		}
	}
}

func TestUnitDirectMailConfigSetUpdateClearsValidation(t *testing.T) {
	var updates atomic.Int32
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ConfigSetUpdate":
			updates.Add(1)
			var option map[string]interface{}
			if err := json.Unmarshal([]byte(r.Form.Get("ValidationOption")), &option); err != nil {
				t.Error(err)
			}
			if enabled, ok := option["Enabled"].(bool); !ok || enabled {
				t.Errorf("validation removal must explicitly disable validation: %v", option)
			}
			for _, field := range []string{"ForbiddenStatusList", "ForbiddenSubStatusList"} {
				if values, ok := option[field].([]interface{}); !ok || len(values) != 0 {
					t.Errorf("validation removal must send an empty %s: %v", field, option)
				}
			}
			directMailWriteJSON(t, w, map[string]interface{}{"Id": "cs-1"})
		case "ConfigSetDetail":
			directMailWriteJSON(t, w, map[string]interface{}{"Detail": map[string]interface{}{
				"Id": "cs-1", "Name": "cs",
				"ValidationOption": map[string]interface{}{"Enabled": false, "ForbiddenStatusList": []interface{}{}, "ForbiddenSubStatusList": []interface{}{}},
			}})
		default:
			t.Errorf("unexpected action %s", r.Form.Get("Action"))
		}
	})
	attrs := map[string]string{
		"name": "cs", "description": "", "ip_pool_id": "", "is_public_channel_backoff": "false",
		"validation_option.#": "1", "validation_option.0.enabled": "true",
		"validation_option.0.forbidden_status_list.#": "1", "validation_option.0.forbidden_status_list.0": "INVALID",
		"validation_option.0.forbidden_sub_status_list.#": "0",
	}
	diff := sdktf.NewInstanceDiff()
	diff.Attributes["validation_option.#"] = &sdktf.ResourceAttrDiff{Old: "1", New: "0"}
	d, err := schema.InternalMap(directmail.ResourceConfigSet().Schema).Data(&sdktf.InstanceState{ID: "cs-1", Attributes: attrs}, diff)
	if err != nil {
		t.Fatal(err)
	}
	if got := directmail.DirectMailConfiguredValidationOptionKey(d); got != "false||" {
		t.Fatalf("removed validation block must ignore stale SDK children, got %q", got)
	}
	if err := directmail.ResourceConfigSetUpdate(d, client); err != nil {
		t.Fatal(err)
	}
	if updates.Load() != 1 || len(d.Get("validation_option").([]interface{})) != 0 {
		t.Fatalf("expected one update and no validation block, got updates=%d option=%v", updates.Load(), d.Get("validation_option"))
	}
}

func TestUnitDirectMailConfigSetWaitForBackoffPropagation(t *testing.T) {
	var calls int32
	stale := map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "desc"}
	fresh := map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "desc", "IsPublicChannelBackoff": true}
	client := directMailConfigSetPropagationClient(t, stale, fresh, &calls)
	d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{
		"name":                      "cs",
		"description":               "desc",
		"is_public_channel_backoff": true,
	})
	d.SetId("cs-1")
	service := directmail.NewDirectMailServiceV2(client)
	if err := directmail.WaitForDirectMailConfigSetUpdatePropagated(d, service); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("the waiter must not accept a detail that omits IsPublicChannelBackoff while the target is true; calls=%d", calls)
	}
}

func TestUnitDirectMailConfigSetEmptyValidationOptionBlock(t *testing.T) {
	// An empty validation_option {} block can surface as a nil list element; the
	// configured key must normalize it as disabled-empty instead of panicking.
	d := schema.TestResourceDataRaw(t, directmail.ResourceConfigSet().Schema, map[string]interface{}{
		"name":              "cs",
		"validation_option": []interface{}{nil},
	})
	if got := directmail.DirectMailConfiguredValidationOptionKey(d); got != "false||" {
		t.Fatalf("an empty validation block must normalize to disabled-empty, got %q", got)
	}

	var updates atomic.Int32
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ConfigSetUpdate":
			updates.Add(1)
			var option map[string]interface{}
			if err := json.Unmarshal([]byte(r.Form.Get("ValidationOption")), &option); err != nil {
				t.Error(err)
			}
			if enabled, ok := option["Enabled"].(bool); !ok || enabled {
				t.Errorf("an empty validation block must send Enabled=false: %v", option)
			}
			for _, field := range []string{"ForbiddenStatusList", "ForbiddenSubStatusList"} {
				if values, ok := option[field].([]interface{}); !ok || len(values) != 0 {
					t.Errorf("an empty validation block must send an empty %s: %v", field, option)
				}
			}
			directMailWriteJSON(t, w, map[string]interface{}{"Id": "cs-1"})
		case "ConfigSetDetail":
			directMailWriteJSON(t, w, map[string]interface{}{"Detail": map[string]interface{}{
				"Id": "cs-1", "Name": "cs", "Description": "updated",
			}})
		default:
			t.Errorf("unexpected action %s", r.Form.Get("Action"))
		}
	})
	attrs := map[string]string{
		"name": "cs", "description": "cs", "ip_pool_id": "", "is_public_channel_backoff": "false",
		"validation_option.#": "1",
	}
	diff := sdktf.NewInstanceDiff()
	diff.Attributes["description"] = &sdktf.ResourceAttrDiff{Old: "cs", New: "updated"}
	state, err := schema.InternalMap(directmail.ResourceConfigSet().Schema).Data(&sdktf.InstanceState{ID: "cs-1", Attributes: attrs}, diff)
	if err != nil {
		t.Fatal(err)
	}
	if err := directmail.ResourceConfigSetUpdate(state, client); err != nil {
		t.Fatal(err)
	}
	if updates.Load() != 1 {
		t.Fatalf("expected exactly one ConfigSetUpdate, got %d", updates.Load())
	}
}
