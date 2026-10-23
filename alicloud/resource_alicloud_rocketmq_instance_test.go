package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	strings "strings"
	"sync"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The three collection-order markers below record stock TypeList coverage debt:
// those lists already had no reorder coverage before the trace_on fix and are
// untouched by it, so they are bypassed instead of widening this change.
// collection-order: bypass acl_info.acl_types -- pre-existing list unchanged by the trace_on fix
// collection-order: bypass network_info.internet_info.ip_whitelist -- pre-existing list unchanged by the trace_on fix
// collection-order: bypass network_info.vpc_info.vswitches -- pre-existing list unchanged by the trace_on fix

// Test case to verify the fix for send_receive_ratio validation issue
func TestAccAliCloudRocketmqInstance_SendReceiveRatioValidation(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.validation_test"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	//testAccCheck is not being used
	//testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstancevalidation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4101)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							// This should be out of range [0.05, 0.5] to test validation
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.03",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "validation test",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected product_info\.0\.send_receive_ratio to be in the range \(0\.050000 - 0\.500000\), got 0\.030000`),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("alicloud_rocketmq_instance", &resource.Sweeper{
		Name: "alicloud_rocketmq_instance",
		F:    testSweepRocketMq,
	})
}

// Instances created by the documented examples carry the documented default
// name, so the regular sweeper has to treat that prefix as its own cleanup
// target; otherwise those pay-as-you-go instances are left behind whenever
// sweepAll is disabled.
var rocketmqSweepDocPrefixes = []string{
	"terraform-example",
}

var rocketmqSweepTestPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"testAcc",
}

func shouldSweepRocketmqInstance(name string) bool {
	if sweepAll() {
		return true
	}
	lowerName := strings.ToLower(name)
	for _, prefixes := range [][]string{rocketmqSweepTestPrefixes, rocketmqSweepDocPrefixes} {
		for _, prefix := range prefixes {
			if strings.HasPrefix(lowerName, strings.ToLower(prefix)) {
				return true
			}
		}
	}
	return false
}

// listRocketmqInstances walks every page before anything is deleted. Collecting
// the full result first keeps the deletion from shifting the pages that are
// still to be read, which would silently skip instances on later pages.
func listRocketmqInstances(fetchPage func(pageNumber, pageSize int) ([]interface{}, error)) ([]interface{}, error) {
	const pageSize = 200
	var instances []interface{}
	for pageNumber := 1; ; pageNumber++ {
		page, err := fetchPage(pageNumber, pageSize)
		if err != nil {
			return nil, WrapErrorf(err, "failed to list rocketmq instances on page %d.", pageNumber)
		}
		instances = append(instances, page...)
		if len(page) < pageSize {
			return instances, nil
		}
	}
}

func testSweepRocketMq(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "/instances"
	instances, err := listRocketmqInstances(func(pageNumber, pageSize int) ([]interface{}, error) {
		query := map[string]*string{
			"pageNumber": tea.String(strconv.Itoa(pageNumber)),
			"pageSize":   tea.String(strconv.Itoa(pageSize)),
		}
		response, err := client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.data.list", response)
		if err != nil {
			return nil, WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.list", response)
		}
		if resp == nil {
			return nil, nil
		}
		list, ok := resp.([]interface{})
		if !ok {
			return nil, WrapErrorf(fmt.Errorf("unexpected $.data.list type %T", resp), FailedGetAttributeMsg, action, "$.data.list", response)
		}
		return list, nil
	})
	if err != nil {
		return err
	}

	return sweepRocketmqInstances(instances, shouldSweepRocketmqInstance, func(instanceId string) error {
		action := fmt.Sprintf("/instances/%s", instanceId)
		body := map[string]interface{}{"instanceId": instanceId}
		_, err := client.RoaDelete("RocketMQ", "2022-08-01", action, make(map[string]*string), nil, body, true)
		return err
	})
}

// sweepRocketmqInstances deletes the selected instances and reports the failures
// after the remaining instances have been processed, so an undeletable instance
// can no longer make the whole sweep look successful. An instance the API
// reports as already absent is a successful cleanup.
func sweepRocketmqInstances(instances []interface{}, shouldDelete func(name string) bool, deleteInstance func(instanceId string) error) error {
	var failures []string
	for _, v := range instances {
		item := v.(map[string]interface{})
		name := item["instanceName"].(string)
		if !shouldDelete(name) {
			log.Printf("[INFO] Skipping rocketmq instance: %s ", name)
			continue
		}
		log.Printf("[INFO] delete rocketmq instance: %s ", name)
		instanceId := item["instanceId"].(string)
		if err := deleteInstance(instanceId); err != nil && !NotFoundError(err) {
			log.Printf("[ERROR] Failed to delete rocketmq instance (%s): %s", name, err)
			failures = append(failures, fmt.Sprintf("%s: %s", name, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("failed to delete %d rocketmq instance(s): %s", len(failures), strings.Join(failures, "; "))
	}
	return nil
}

func TestUnitSweepRocketmqInstances(t *testing.T) {
	instance := func(name, instanceId string) interface{} {
		return map[string]interface{}{"instanceName": name, "instanceId": instanceId}
	}
	all := func(string) bool { return true }

	t.Run("failed sole deletion is reported", func(t *testing.T) {
		err := sweepRocketmqInstances([]interface{}{instance("tf-testAcc-sweep-sole", "rmq-1")}, all, func(string) error {
			return fmt.Errorf("InvalidParameter")
		})
		if err == nil {
			t.Fatal("the failed deletion of the only instance must be reported")
		}
		if !strings.Contains(err.Error(), "tf-testAcc-sweep-sole") || !strings.Contains(err.Error(), "InvalidParameter") {
			t.Fatalf("the error must name the instance and the delete failure, got %s", err)
		}
	})

	t.Run("failed final deletion is reported after the remaining instances were deleted", func(t *testing.T) {
		var deleted []string
		err := sweepRocketmqInstances([]interface{}{
			instance("tf-testAcc-sweep-first", "rmq-1"),
			instance("tf-testAcc-sweep-final", "rmq-2"),
		}, all, func(instanceId string) error {
			deleted = append(deleted, instanceId)
			if instanceId == "rmq-2" {
				return fmt.Errorf("Throttling")
			}
			return nil
		})
		if err == nil {
			t.Fatal("the failed deletion of the final instance must be reported")
		}
		if len(deleted) != 2 {
			t.Fatalf("the remaining instances must still be processed, deleted %v", deleted)
		}
	})

	t.Run("successful sweep returns nil", func(t *testing.T) {
		if err := sweepRocketmqInstances([]interface{}{instance("tf-testAcc-sweep-ok", "rmq-1")}, all, func(string) error {
			return nil
		}); err != nil {
			t.Fatalf("a successful sweep must return nil, got %s", err)
		}
	})

	t.Run("verified already absent instance counts as cleaned", func(t *testing.T) {
		if err := sweepRocketmqInstances([]interface{}{instance("tf-testAcc-sweep-gone", "rmq-1")}, all, func(string) error {
			return errs.GetNotFoundErrorFromString("instance rmq-1 does not exist")
		}); err != nil {
			t.Fatalf("an already absent instance must not fail the sweep, got %s", err)
		}
	})

	t.Run("name filtering is preserved", func(t *testing.T) {
		var deleted []string
		matchesTestPrefix := func(name string) bool {
			return strings.HasPrefix(name, "tf-testAcc")
		}
		if err := sweepRocketmqInstances([]interface{}{
			instance("keep-me", "rmq-1"),
			instance("tf-testAcc-sweep", "rmq-2"),
		}, matchesTestPrefix, func(instanceId string) error {
			deleted = append(deleted, instanceId)
			return nil
		}); err != nil {
			t.Fatalf("unexpected sweep error: %s", err)
		}
		if len(deleted) != 1 || deleted[0] != "rmq-2" {
			t.Fatalf("only the instance matching the test prefixes may be deleted, got %v", deleted)
		}
	})
}

// TestUnitRocketmqSweeperSelection exercises the production predicate through
// the same deletion callback the sweeper uses, so the documented example names
// and the test names are both proven selectable while unrelated instances stay
// protected.
func TestUnitRocketmqSweeperSelection(t *testing.T) {
	instance := func(name, instanceId string) interface{} {
		return map[string]interface{}{"instanceName": name, "instanceId": instanceId}
	}
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")

	t.Run("documented and test names are selected, unrelated names are not", func(t *testing.T) {
		var deleted []string
		err := sweepRocketmqInstances([]interface{}{
			instance("terraform-example", "rmq-doc"),
			instance("terraform-example-custom", "rmq-doc-prefixed"),
			instance("tf-testAccRocketmqInstance12345", "rmq-testacc-camel"),
			instance("tf_testAccrocketmq", "rmq-testacc-underscore"),
			instance("tf_test_rocketmq", "rmq-test-underscore"),
			instance("tf-test-rocketmq", "rmq-test-dash"),
			instance("testAccRocketmq", "rmq-testacc-plain"),
			instance("tf-testaccrocketmq", "rmq-testacc-lower"),
			instance("prod-rocketmq", "rmq-keep-prod"),
			instance("my-instance", "rmq-keep-custom"),
			instance("terraform", "rmq-keep-terraform"),
			instance("terraform-other", "rmq-keep-other"),
		}, shouldSweepRocketmqInstance, func(instanceId string) error {
			deleted = append(deleted, instanceId)
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected sweep error: %s", err)
		}
		selected := map[string]bool{}
		for _, id := range deleted {
			selected[id] = true
		}
		for _, expected := range []string{
			"rmq-doc", "rmq-doc-prefixed",
			"rmq-testacc-camel", "rmq-testacc-underscore", "rmq-test-underscore",
			"rmq-test-dash", "rmq-testacc-plain", "rmq-testacc-lower",
		} {
			if !selected[expected] {
				t.Fatalf("the sweeper must select %s, deleted %v", expected, deleted)
			}
		}
		for _, protected := range []string{"rmq-keep-prod", "rmq-keep-custom", "rmq-keep-terraform", "rmq-keep-other"} {
			if selected[protected] {
				t.Fatalf("the sweeper must not delete the protected instance %s, deleted %v", protected, deleted)
			}
		}
		if len(deleted) != 8 {
			t.Fatalf("exactly the eight selected instances may be deleted, deleted %v", deleted)
		}
	})

	t.Run("sweepAll still bypasses the name filter", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
		var deleted []string
		err := sweepRocketmqInstances([]interface{}{
			instance("prod-rocketmq", "rmq-prod"),
		}, shouldSweepRocketmqInstance, func(instanceId string) error {
			deleted = append(deleted, instanceId)
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected sweep error: %s", err)
		}
		if len(deleted) != 1 || deleted[0] != "rmq-prod" {
			t.Fatalf("sweepAll must select every instance, deleted %v", deleted)
		}
	})
}

func TestUnitListRocketmqInstances(t *testing.T) {
	instance := func(name, instanceId string) interface{} {
		return map[string]interface{}{"instanceName": name, "instanceId": instanceId}
	}
	pageOf := func(size int, prefix string) []interface{} {
		page := make([]interface{}, 0, size)
		for i := 0; i < size; i++ {
			name := fmt.Sprintf("%s-%03d", prefix, i)
			page = append(page, instance(name, name))
		}
		return page
	}

	t.Run("matching instance on a later page is collected and swept", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
		var fetched []int
		instances, err := listRocketmqInstances(func(pageNumber, pageSize int) ([]interface{}, error) {
			fetched = append(fetched, pageNumber)
			if pageSize != 200 {
				t.Fatalf("the sweeper must use the maximum supported page size, got %d", pageSize)
			}
			switch pageNumber {
			case 1:
				return pageOf(200, "prod-rocketmq"), nil
			case 2:
				return []interface{}{instance("tf-testAccRocketmqPageTwo", "rmq-page-two")}, nil
			default:
				return nil, fmt.Errorf("unexpected request for page %d after a short page", pageNumber)
			}
		})
		if err != nil {
			t.Fatalf("unexpected list error: %s", err)
		}
		if len(fetched) != 2 || len(instances) != 201 {
			t.Fatalf("the walk must read page 1 and the short page 2, fetched %v with %d instances", fetched, len(instances))
		}
		var deleted []string
		if err := sweepRocketmqInstances(instances, shouldSweepRocketmqInstance, func(instanceId string) error {
			deleted = append(deleted, instanceId)
			return nil
		}); err != nil {
			t.Fatalf("unexpected sweep error: %s", err)
		}
		if len(deleted) != 1 || deleted[0] != "rmq-page-two" {
			t.Fatalf("the instance that only exists on the second page must be swept, deleted %v", deleted)
		}
	})

	t.Run("exact page boundary keeps walking until a short page", func(t *testing.T) {
		var fetched []int
		instances, err := listRocketmqInstances(func(pageNumber, pageSize int) ([]interface{}, error) {
			fetched = append(fetched, pageNumber)
			if pageNumber == 1 {
				return pageOf(200, "prod-rocketmq"), nil
			}
			return []interface{}{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected list error: %s", err)
		}
		if len(instances) != 200 || len(fetched) != 2 {
			t.Fatalf("a full page must be followed by one more request, fetched %v with %d instances", fetched, len(instances))
		}
	})

	t.Run("later page list failure is returned", func(t *testing.T) {
		instances, err := listRocketmqInstances(func(pageNumber, pageSize int) ([]interface{}, error) {
			if pageNumber == 1 {
				return pageOf(200, "prod-rocketmq"), nil
			}
			return nil, fmt.Errorf("Throttling.User")
		})
		if err == nil {
			t.Fatal("a failed later page must fail the list")
		}
		if instances != nil {
			t.Fatalf("a partial list must not be returned, got %d instances", len(instances))
		}
		if !strings.Contains(err.Error(), "page 2") || !strings.Contains(err.Error(), "Throttling.User") {
			t.Fatalf("the error must name the failing page and the cause, got %s", err)
		}
	})

	t.Run("deleting instances cannot shift the pages still to be read", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
		// The backing slice behaves like the service list: deleting an instance
		// removes it and shifts every later entry to a lower offset.
		backing := pageOf(200, "prod-rocketmq")
		backing[0] = instance("tf-testAccSweepPageOne", "rmq-page-one")
		backing = append(backing, instance("tf-testAccSweepPageTwo", "rmq-page-two"))

		var steps []string
		instances, err := listRocketmqInstances(func(pageNumber, pageSize int) ([]interface{}, error) {
			steps = append(steps, fmt.Sprintf("fetch-%d", pageNumber))
			start := (pageNumber - 1) * pageSize
			if start >= len(backing) {
				return nil, nil
			}
			end := start + pageSize
			if end > len(backing) {
				end = len(backing)
			}
			page := make([]interface{}, end-start)
			copy(page, backing[start:end])
			return page, nil
		})
		if err != nil {
			t.Fatalf("unexpected list error: %s", err)
		}
		if len(instances) != 201 {
			t.Fatalf("both pages must be collected before any deletion, got %d instances", len(instances))
		}

		if err := sweepRocketmqInstances(instances, shouldSweepRocketmqInstance, func(instanceId string) error {
			steps = append(steps, "delete-"+instanceId)
			for i, v := range backing {
				if v.(map[string]interface{})["instanceId"] == instanceId {
					backing = append(backing[:i], backing[i+1:]...)
					break
				}
			}
			return nil
		}); err != nil {
			t.Fatalf("unexpected sweep error: %s", err)
		}

		if len(steps) < 3 || steps[0] != "fetch-1" || steps[1] != "fetch-2" {
			t.Fatalf("every page must be read before the first deletion, steps %v", steps)
		}
		selected := map[string]bool{}
		for _, step := range steps {
			if strings.HasPrefix(step, "delete-") {
				selected[strings.TrimPrefix(step, "delete-")] = true
			}
		}
		if !selected["rmq-page-one"] || !selected["rmq-page-two"] {
			t.Fatalf("both selected instances must be deleted even though deleting shifts the list, steps %v", steps)
		}
		if len(backing) != 199 {
			t.Fatalf("only the two selected instances may be deleted, %d entries left", len(backing))
		}
	})
}

func TestAccAliCloudRocketmqInstance_bugfix(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 4665
func TestAccAliCloudRocketmqInstance_basic4665(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4665)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "cluster_ha",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "ultimate",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "cluster_ha",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "ultimate",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				// The software maintenance window is managed by the service for
				// provisioned instances and cannot be set standalone here.
				// Removing the required sub_series_code stops the plan before
				// the placeholder reaches the API.
				Config: testAccConfig(map[string]interface{}{
					"software": []map[string]interface{}{
						{"maintain_time": "05:00:00"},
					},
					"sub_series_code": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				// Restore the required argument and drop the negative-step
				// placeholder; the post-test destroy reuses the last config, so
				// leaving sub_series_code removed fails the destroy with
				// "Missing required argument" and leaks the instance.
				Config: testAccConfig(map[string]interface{}{
					"software":        REMOVEKEY,
					"sub_series_code": "cluster_ha",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap4665 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence4665(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-j"
	}

	resource "alicloud_kms_instance" "default" {
  		product_version = "3"
  		vpc_num         = "1"
  		key_num         = "1000"
  		secret_num      = "0"
  		spec            = "1000"
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		"cn-hangzhou-k",
    		"cn-hangzhou-j"
  		]
  		timeouts {
    		delete = "60m"
  		}
	}

	resource "alicloud_kms_key" "default" {
  		description            = "Hello KMS"
  		status                 = "Enabled"
  		pending_window_in_days = "7"
  		dkms_instance_id       = alicloud_kms_instance.default.id
	}
`, name)
}

// Case 4652 From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic4652(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4652)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "single_node",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "standard",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "single_node",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "single_node",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap4652 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence4652(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 4128
func TestAccAliCloudRocketmqInstance_basic4128(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4128)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil, // bypass linter, prepaid (Subscription) resource cannot destroy
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":           "Subscription",
					"auto_renew":             "true",
					"auto_renew_period":      "1",
					"auto_renew_period_unit": "Month",
					"period":                 "1",
					"period_unit":            "Month",
					"sub_series_code":        "cluster_ha",
					"instance_name":          name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":           "rmq",
						"series_code":            "standard",
						"payment_type":           "Subscription",
						"sub_series_code":        "cluster_ha",
						"instance_name":          name,
						"auto_renew":             "true",
						"auto_renew_period":      "1",
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "_update",
					"auto_renew":        "true",
					"auto_renew_period": "1",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code":      "rmq",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"series_code":       "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "Subscription",
					"period":          "1",
					"sub_series_code": "cluster_ha",
					"remark":          "自动化测试购买使用11",
					"period_unit":     "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"auto_renew":        "true",
						"auto_renew_period": "1",
						"service_code":      "rmq",
						"resource_group_id": CHECKSET,
						"series_code":       "standard",
						"payment_type":      "Subscription",
						"period":            "1",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"period_unit":       "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap4128 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence4128(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}


`, name)
}

// Case 4101
func TestAccAliCloudRocketmqInstance_basic4101(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_code": "rmq",
					"series_code":  "professional",
					"product_info": []map[string]interface{}{
						{
							"auto_scaling":           "false",
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
								},
							},
						},
					},
					"payment_type":    "PayAsYouGo",
					"sub_series_code": "cluster_ha",
					"instance_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_code":    "rmq",
						"series_code":     "professional",
						"payment_type":    "PayAsYouGo",
						"sub_series_code": "cluster_ha",
						"instance_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_info": []map[string]interface{}{
						{
							"acl_types": []string{
								"default", "apache_acl", "aliyun_ram"},
							"default_vpc_auth_free": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "自动化测试购买使用11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "自动化测试购买使用11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "80",
							"auto_scaling":           "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.2.2", "192.168.1.1", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.2.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.2.2", "192.168.1.1", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "_update",
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap4101 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence4101(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

  zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

// Case 4665  twin
func TestAccAliCloudRocketmqInstance_basic4665_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4665)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4665)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.u2.10xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
							"storage_encryption":     "true",
							"storage_secret_key":     "${alicloud_kms_key.default.id}",
						},
					},
					"service_code": "rmq",
					"series_code":  "ultimate",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
									"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "ultimate",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"product_info.#":    "1",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 4652  twin From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic4652_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4652)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "single_node",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 4128  twin
func TestAccAliCloudRocketmqInstance_basic4128_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4128)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil, // bypass linter, postpay resource cannot destroy
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":          name,
					"auto_renew":             "true",
					"auto_renew_period":      "1",
					"auto_renew_period_unit": "Year",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s2.2xlarge",
							"send_receive_ratio":     "0.4",
							"message_retention_time": "70",
						},
					},
					"service_code":      "rmq",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"series_code":       "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
									"vswitch_id": "${data.alicloud_vswitches.default.vswitches.0.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":    "Subscription",
					"period":          "1",
					"sub_series_code": "cluster_ha",
					"remark":          "自动化测试购买使用11",
					"period_unit":     "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":          name,
						"auto_renew":             "true",
						"auto_renew_period":      "1",
						"auto_renew_period_unit": "Year",
						"service_code":           "rmq",
						"resource_group_id":      CHECKSET,
						"series_code":            "standard",
						"payment_type":           "Subscription",
						"period":                 "1",
						"sub_series_code":        "cluster_ha",
						"remark":                 "自动化测试购买使用11",
						"period_unit":            "Year",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period_unit": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

// Case 4101  twin
func TestAccAliCloudRocketmqInstance_basic4101_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap4101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence4101)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.p2.4xlarge",
							"send_receive_ratio":     "0.3",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "professional",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "cluster_ha",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"acl_info": []map[string]interface{}{
						{
							"acl_types":             []string{"default", "apache_acl", "aliyun_ram"},
							"default_vpc_auth_free": "false",
						},
					},
					"ip_whitelists": []string{"0.0.0.0/0", "192.168.1.1", "192.168.2.2"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "professional",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "cluster_ha",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"acl_info.#":        "1",
						"ip_whitelists.#":   "3",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Rocketmq Instance. >>> Resource test cases, automatically generated.
// Case 创建serverless实例 6747
func TestAccAliCloudRocketmqInstance_basic6747(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap6747)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence6747)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"series_code":  "standard",
					"payment_type": "PayAsYouGo",
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s3.nxlarge",
							"auto_scaling":           "true",
							"message_retention_time": "72",
							// "support_auto_scaling":   "true",
							"trace_on": "false",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":             "${alicloud_vpc.createVPC.id}",
									"security_group_ids": "${alicloud_security_group.CreateSecurityGroup.id}",
									"vswitches": []map[string]interface{}{
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
										},
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch2.id}",
										},
									},
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
								},
							},
						},
					},
					"sub_series_code": "serverless",
					"instance_name":   name,
					"service_code":    "rmq",
					"commodity_code":  "ons_rmqsrvlesspost_public_cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"resource_group_id": CHECKSET,
						"sub_series_code":   "serverless",
						"instance_name":     name,
						"service_code":      "rmq",
						"commodity_code":    "ons_rmqsrvlesspost_public_cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s3.nxlarge",
							"auto_scaling":           "true",
							"message_retention_time": "72",
							"trace_on":               "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_info.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap6747 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence6747(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "createVPC" {
  description = "111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_security_group" "CreateSecurityGroup" {
  name   = var.name
  vpc_id = alicloud_vpc.createVPC.id
}

resource "alicloud_vswitch" "createVSwitch" {
  vpc_id       = alicloud_vpc.createVPC.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.17.0.0/16"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "createVSwitch2" {
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.18.0.0/16"
  vswitch_name = format("%%s4", var.name)
  zone_id      = "cn-hangzhou-j"
}


`, name)
}

func TestAccAliCloudRocketmqInstance_traceOn(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap6747)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence6747)
	productInfo := func(extra map[string]interface{}) []map[string]interface{} {
		item := map[string]interface{}{
			"msg_process_spec":       "rmq.s3.nxlarge",
			"auto_scaling":           "true",
			"message_retention_time": "72",
		}
		for k, val := range extra {
			item[k] = val
		}
		return []map[string]interface{}{item}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"series_code":       "standard",
					"payment_type":      "PayAsYouGo",
					"product_info":      productInfo(nil),
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":             "${alicloud_vpc.createVPC.id}",
									"security_group_ids": "${alicloud_security_group.CreateSecurityGroup.id}",
									"vswitches": []map[string]interface{}{
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
										},
										{
											"vswitch_id": "${alicloud_vswitch.createVSwitch2.id}",
										},
									},
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec": "disable",
									"flow_out_type": "uninvolved",
								},
							},
						},
					},
					"sub_series_code": "serverless",
					"instance_name":   name,
					"service_code":    "rmq",
					"commodity_code":  "ons_rmqsrvlesspost_public_cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"product_info.#": "1",
					}),
				),
			},
			{
				// An unrelated update must not backfill the computed trace_on value.
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": productInfo(map[string]interface{}{"trace_on": "true"}),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_info.0.trace_on": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_info": productInfo(map[string]interface{}{"trace_on": "false"}),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_info.0.trace_on": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

// Case 创建单节点实例用例_副本1719467789892 7144 From November 9th, 2024 Beijing time, the RocketMQ version 5.x cannot create single node instances.
func SkipTestAccAliCloudRocketmqInstance_basic7144(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqInstanceMap7144)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqInstanceBasicDependence7144)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"product_info": []map[string]interface{}{
						{
							"msg_process_spec":       "rmq.s1.micro",
							"send_receive_ratio":     "0.3",
							"auto_scaling":           "false",
							"message_retention_time": "70",
						},
					},
					"service_code": "rmq",
					"series_code":  "standard",
					"network_info": []map[string]interface{}{
						{
							"vpc_info": []map[string]interface{}{
								{
									"vpc_id":     "${alicloud_vpc.createVPC.id}",
									"vswitch_id": "${alicloud_vswitch.createVSwitch.id}",
								},
							},
							"internet_info": []map[string]interface{}{
								{
									"internet_spec":      "enable",
									"flow_out_type":      "payByBandwidth",
									"flow_out_bandwidth": "30",
									"ip_whitelist": []string{
										"192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"},
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"sub_series_code":   "single_node",
					"remark":            "自动化测试购买使用11",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"commodity_code":    "ons_rmqpost_public_cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"service_code":      "rmq",
						"series_code":       "standard",
						"payment_type":      "PayAsYouGo",
						"sub_series_code":   "single_node",
						"remark":            "自动化测试购买使用11",
						"resource_group_id": CHECKSET,
						"commodity_code":    "ons_rmqpost_public_cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AliCloudRocketmqInstanceMap7144 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqInstanceBasicDependence7144(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
  zone_id      = data.alicloud_zones.default.zones.0.id
}

`, name)
}

// Test Rocketmq Instance. <<< Resource test cases, automatically generated.

// rocketmqTraceOnFixtureInstanceId is the instance ID used by the trace_on
// request serialization fixture.
const rocketmqTraceOnFixtureInstanceId = "rmq-fixture-instance"

type rocketmqTraceOnCapturedRequest struct {
	method string
	path   string
	body   []byte
}

// rocketmqTraceOnFixture serves the RocketMQ and BSS endpoints used by the
// instance resource. Create and Update write requests are answered with a 400
// after their body is captured: the resource waits minutes for a RUNNING
// instance after a successful write, while a 4xx is classified as
// non-retryable, so the read/update portions of the chain still run and the
// test stays fast.
type rocketmqTraceOnFixture struct {
	mu       sync.Mutex
	requests []rocketmqTraceOnCapturedRequest
}

func (f *rocketmqTraceOnFixture) serveHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	f.mu.Lock()
	f.requests = append(f.requests, rocketmqTraceOnCapturedRequest{
		method: r.Method,
		path:   r.URL.Path,
		body:   body,
	})
	f.mu.Unlock()

	instancePath := "/instances/" + rocketmqTraceOnFixtureInstanceId
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/instances":
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"code":"InvalidParameter","message":"fixture captured the create request","requestId":"fixture"}`)
	case r.Method == http.MethodPatch && r.URL.Path == instancePath:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"code":"InvalidParameter","message":"fixture captured the update request","requestId":"fixture"}`)
	case r.Method == http.MethodGet && r.URL.Path == instancePath:
		fmt.Fprintf(w, `{"data":{"instanceId":%q,"instanceName":"rmq-trace-on-instance","status":"RUNNING","seriesCode":"standard","serviceCode":"rmq","subSeriesCode":"serverless","paymentType":"PayAsYouGo","productInfo":{"msgProcessSpec":"rmq.s3.nxlarge","autoScaling":true,"messageRetentionTime":72,"sendReceiveRatio":0.2,"traceOn":false},"networkInfo":{"vpcInfo":{"vpcId":"vpc-rmqtraceon"},"internetInfo":{"internetSpec":"disable","flowOutType":"uninvolved"}}}}`, rocketmqTraceOnFixtureInstanceId)
	case r.Method == http.MethodGet && r.URL.Path == instancePath+"/account":
		fmt.Fprint(w, `{"data":{}}`)
	case r.Method == http.MethodGet && r.URL.Path == instancePath+"/ip/whitelists":
		fmt.Fprint(w, `{"data":{"ipWhitelists":[]}}`)
	case r.URL.Path == "/":
		fmt.Fprintf(w, `{"Success":true,"RequestId":"fixture","Data":{"InstanceList":[{"InstanceID":%q,"RenewStatus":"ManualRenewal","RenewalDuration":0,"RenewalDurationUnit":"M"}]}}`, rocketmqTraceOnFixtureInstanceId)
	default:
		fmt.Fprint(w, `{}`)
	}
}

func (f *rocketmqTraceOnFixture) captured(method, path string) []rocketmqTraceOnCapturedRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	captured := make([]rocketmqTraceOnCapturedRequest, 0)
	for _, request := range f.requests {
		if request.method == method && request.path == path {
			captured = append(captured, request)
		}
	}
	return captured
}

func newRocketmqTraceOnFixture(t *testing.T) (*rocketmqTraceOnFixture, *connectivity.AliyunClient) {
	t.Helper()
	fixture := &rocketmqTraceOnFixture{}
	server := httptest.NewServer(http.HandlerFunc(fixture.serveHTTP))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to build the fixture credential: %s", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("failed to build the fixture client: %s", err)
	}
	endpoints.Store("rocketmq", endpoint)
	endpoints.Store("bssopenapi", endpoint)
	return fixture, client
}

func rocketmqTraceOnFixtureRawConfig(instanceName string, traceOn *bool) map[string]interface{} {
	productInfo := map[string]interface{}{
		"msg_process_spec":       "rmq.s3.nxlarge",
		"auto_scaling":           true,
		"message_retention_time": 72,
	}
	if traceOn != nil {
		productInfo["trace_on"] = *traceOn
	}
	return map[string]interface{}{
		"instance_name":   instanceName,
		"payment_type":    "PayAsYouGo",
		"series_code":     "standard",
		"service_code":    "rmq",
		"sub_series_code": "serverless",
		"product_info":    []interface{}{productInfo},
		"network_info": []interface{}{
			map[string]interface{}{
				"vpc_info": []interface{}{
					map[string]interface{}{"vpc_id": "vpc-rmqtraceon"},
				},
				"internet_info": []interface{}{
					map[string]interface{}{
						"internet_spec": "disable",
						"flow_out_type": "uninvolved",
					},
				},
			},
		},
	}
}

// rocketmqTraceOnFixtureRawConfigValue mirrors the protocol raw config: an
// unset optional+computed attribute is present as a null value, not missing.
// cty.GetAttr panics on a missing attribute, so trace_on is always present.
func rocketmqTraceOnFixtureRawConfigValue(traceOn cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"product_info": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"msg_process_spec":       cty.StringVal("rmq.s3.nxlarge"),
				"auto_scaling":           cty.True,
				"message_retention_time": cty.NumberIntVal(72),
				"trace_on":               traceOn,
			}),
		}),
	})
}

func rocketmqTraceOnFixtureNewResourceData(t *testing.T, resourceInstance *schema.Resource, raw map[string]interface{}, traceOn cty.Value) *schema.ResourceData {
	t.Helper()
	diff, err := resourceInstance.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(raw), nil)
	if err != nil {
		t.Fatalf("failed to build the create diff: %s", err)
	}
	diff.RawConfig = rocketmqTraceOnFixtureRawConfigValue(traceOn)
	d, err := schema.InternalMap(resourceInstance.Schema).Data(nil, diff)
	if err != nil {
		t.Fatalf("failed to build the create resource data: %s", err)
	}
	d.MarkNewResource()
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	return d
}

// rocketmqTraceOnFixtureState simulates a resource read back from the API:
// trace_on is Optional+Computed, so state carries a value even when the
// configuration never set it.
func rocketmqTraceOnFixtureState(traceOn bool) *terraform.InstanceState {
	traceOnState := "false"
	if traceOn {
		traceOnState = "true"
	}
	return &terraform.InstanceState{
		ID: rocketmqTraceOnFixtureInstanceId,
		Attributes: map[string]string{
			"instance_name":                                "rmq-trace-on-instance",
			"payment_type":                                 "PayAsYouGo",
			"series_code":                                  "standard",
			"service_code":                                 "rmq",
			"sub_series_code":                              "serverless",
			"product_info.#":                               "1",
			"product_info.0.msg_process_spec":              "rmq.s3.nxlarge",
			"product_info.0.auto_scaling":                  "true",
			"product_info.0.message_retention_time":        "72",
			"product_info.0.trace_on":                      traceOnState,
			"network_info.#":                               "1",
			"network_info.0.vpc_info.#":                    "1",
			"network_info.0.vpc_info.0.vpc_id":             "vpc-rmqtraceon",
			"network_info.0.internet_info.#":               "1",
			"network_info.0.internet_info.0.internet_spec": "disable",
			"network_info.0.internet_info.0.flow_out_type": "uninvolved",
			"ip_whitelists.#":                              "0",
		},
	}
}

func rocketmqTraceOnFixtureExistingResourceData(t *testing.T, resourceInstance *schema.Resource, stateTraceOn bool, raw map[string]interface{}, configTraceOn cty.Value) *schema.ResourceData {
	t.Helper()
	state := rocketmqTraceOnFixtureState(stateTraceOn)
	diff, err := resourceInstance.Diff(context.Background(), state, terraform.NewResourceConfigRaw(raw), nil)
	if err != nil {
		t.Fatalf("failed to build the update diff: %s", err)
	}
	diff.RawConfig = rocketmqTraceOnFixtureRawConfigValue(configTraceOn)
	d, err := schema.InternalMap(resourceInstance.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build the update resource data: %s", err)
	}
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	return d
}

func rocketmqTraceOnFixtureProductInfo(t *testing.T, request rocketmqTraceOnCapturedRequest) map[string]interface{} {
	t.Helper()
	body := make(map[string]interface{})
	if err := json.Unmarshal(request.body, &body); err != nil {
		t.Fatalf("failed to decode the %s %s body %q: %s", request.method, request.path, string(request.body), err)
	}
	productInfo, ok := body["productInfo"].(map[string]interface{})
	if !ok {
		t.Fatalf("the %s %s body has no productInfo object: %s", request.method, request.path, string(request.body))
	}
	return productInfo
}

func rocketmqTraceOnFixtureAssertOmitted(t *testing.T, productInfo map[string]interface{}, label string) {
	t.Helper()
	if value, ok := productInfo["traceOn"]; ok {
		t.Fatalf("%s must omit the traceOn key, got %v", label, value)
	}
}

func rocketmqTraceOnFixtureAssertSent(t *testing.T, productInfo map[string]interface{}, label string, want bool) {
	t.Helper()
	value, ok := productInfo["traceOn"]
	if !ok {
		t.Fatalf("%s must send the traceOn key, body productInfo is %v", label, productInfo)
	}
	sent, ok := value.(bool)
	if !ok {
		t.Fatalf("%s must send traceOn as a JSON boolean, got %T (%v)", label, value, value)
	}
	if sent != want {
		t.Fatalf("%s must send traceOn=%t, got %t", label, want, sent)
	}
}

// TestUnitAliCloudRocketmqInstanceTraceOnRequestSerialization drives the real
// resource and Alibaba Cloud client request chain against an isolated HTTP
// fixture and locks the wire contract of product_info.trace_on: the key is
// serialized only when the configuration sets it explicitly. The pre-fix
// implementation sent traceOn=false for an unconfigured create tail update and
// for updates that changed unrelated attributes, so both omission assertions
// fail against it.
func TestUnitAliCloudRocketmqInstanceTraceOnRequestSerialization(t *testing.T) {
	resourceInstance := resourceAliCloudRocketmqInstance()
	createPath := "/instances"
	updatePath := "/instances/" + rocketmqTraceOnFixtureInstanceId

	t.Run("unconfigured create omits traceOn", func(t *testing.T) {
		fixture, client := newRocketmqTraceOnFixture(t)
		d := rocketmqTraceOnFixtureNewResourceData(t, resourceInstance, rocketmqTraceOnFixtureRawConfig("rmq-trace-on-create", nil), cty.NullVal(cty.Bool))

		if err := resourceInstance.Create(d, client); err == nil {
			t.Fatal("the fixture answers the create call with 400, Create must return the API error")
		}
		creates := fixture.captured(http.MethodPost, createPath)
		if len(creates) != 1 {
			t.Fatalf("expected exactly 1 create request, got %d", len(creates))
		}
		rocketmqTraceOnFixtureAssertOmitted(t, rocketmqTraceOnFixtureProductInfo(t, creates[0]), "an unconfigured create")
	})

	t.Run("unconfigured create tail update sends no update request", func(t *testing.T) {
		fixture, client := newRocketmqTraceOnFixture(t)
		d := rocketmqTraceOnFixtureNewResourceData(t, resourceInstance, rocketmqTraceOnFixtureRawConfig("rmq-trace-on-create", nil), cty.NullVal(cty.Bool))

		if err := resourceInstance.Update(d, client); err != nil {
			t.Fatalf("the tail update of an unconfigured create must not call the API: %s", err)
		}
		updates := fixture.captured(http.MethodPatch, updatePath)
		if len(updates) != 0 {
			t.Fatalf("the tail update of an unconfigured create must issue no update request, got %d: %+v", len(updates), updates)
		}
	})

	t.Run("unrelated update omits traceOn", func(t *testing.T) {
		fixture, client := newRocketmqTraceOnFixture(t)
		d := rocketmqTraceOnFixtureExistingResourceData(t, resourceInstance, false, rocketmqTraceOnFixtureRawConfig("rmq-trace-on-renamed", nil), cty.NullVal(cty.Bool))

		if err := resourceInstance.Update(d, client); err == nil {
			t.Fatal("the fixture answers the update call with 400, Update must return the API error")
		}
		updates := fixture.captured(http.MethodPatch, updatePath)
		if len(updates) != 1 {
			t.Fatalf("expected exactly 1 update request for the instance_name change, got %d", len(updates))
		}
		rocketmqTraceOnFixtureAssertOmitted(t, rocketmqTraceOnFixtureProductInfo(t, updates[0]), "an update of unrelated attributes")
		body := make(map[string]interface{})
		if err := json.Unmarshal(updates[0].body, &body); err != nil {
			t.Fatalf("failed to decode the update body: %s", err)
		}
		if body["instanceName"] != "rmq-trace-on-renamed" {
			t.Fatalf("the unrelated update must still send the new instanceName, got %v", body["instanceName"])
		}
	})

	t.Run("explicit true is sent as true", func(t *testing.T) {
		fixture, client := newRocketmqTraceOnFixture(t)
		enabled := true
		d := rocketmqTraceOnFixtureNewResourceData(t, resourceInstance, rocketmqTraceOnFixtureRawConfig("rmq-trace-on-create", &enabled), cty.True)

		if err := resourceInstance.Update(d, client); err == nil {
			t.Fatal("the fixture answers the update call with 400, Update must return the API error")
		}
		updates := fixture.captured(http.MethodPatch, updatePath)
		if len(updates) != 1 {
			t.Fatalf("expected exactly 1 update request for trace_on=true, got %d", len(updates))
		}
		rocketmqTraceOnFixtureAssertSent(t, rocketmqTraceOnFixtureProductInfo(t, updates[0]), "an explicit trace_on=true", true)
	})

	t.Run("explicit false is sent as false", func(t *testing.T) {
		fixture, client := newRocketmqTraceOnFixture(t)
		disabled := false
		d := rocketmqTraceOnFixtureNewResourceData(t, resourceInstance, rocketmqTraceOnFixtureRawConfig("rmq-trace-on-create", &disabled), cty.False)

		if err := resourceInstance.Update(d, client); err == nil {
			t.Fatal("the fixture answers the update call with 400, Update must return the API error")
		}
		updates := fixture.captured(http.MethodPatch, updatePath)
		if len(updates) != 1 {
			t.Fatalf("expected exactly 1 update request for trace_on=false, got %d", len(updates))
		}
		rocketmqTraceOnFixtureAssertSent(t, rocketmqTraceOnFixtureProductInfo(t, updates[0]), "an explicit trace_on=false", false)
	})
}
