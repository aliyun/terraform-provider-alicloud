package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// This file 补齐 (fills in) sweepers for the remaining alicloud_dts_* resources
// that had no sweeper registration: alicloud_dts_synchronization_job,
// alicloud_dts_instance, alicloud_dts_migration_instance,
// alicloud_dts_synchronization_instance and alicloud_dts_consumer_channel.
//
// The two pre-existing sweepers (alicloud_dts_migration_job and
// alicloud_dts_subscription_job) live in their own test files and were patched
// in place to add the same name-based NODELETING protection and sweepAll()
// bypass used here.
//
// All sweepers mirror the Delete API call shape used by the provider
// resource_alicloud_dts_*.go Delete handlers. Child resources (jobs and
// consumer channels) declare Dependencies on their parents so the SDK runs
// them before the instance-level sweepers.

func init() {
	resource.AddTestSweepers(
		"alicloud_dts_synchronization_job",
		&resource.Sweeper{
			Name: "alicloud_dts_synchronization_job",
			F:    testSweepDtsSynchronizationJob,
		})
	resource.AddTestSweepers(
		"alicloud_dts_instance",
		&resource.Sweeper{
			Name:         "alicloud_dts_instance",
			Dependencies: []string{"alicloud_dts_migration_job", "alicloud_dts_subscription_job", "alicloud_dts_synchronization_job"},
			F:            testSweepDtsInstance,
		})
	resource.AddTestSweepers(
		"alicloud_dts_migration_instance",
		&resource.Sweeper{
			Name:         "alicloud_dts_migration_instance",
			Dependencies: []string{"alicloud_dts_migration_job"},
			F:            testSweepDtsMigrationInstance,
		})
	resource.AddTestSweepers(
		"alicloud_dts_synchronization_instance",
		&resource.Sweeper{
			Name:         "alicloud_dts_synchronization_instance",
			Dependencies: []string{"alicloud_dts_synchronization_job"},
			F:            testSweepDtsSynchronizationInstance,
		})
	resource.AddTestSweepers(
		"alicloud_dts_consumer_channel",
		&resource.Sweeper{
			Name:         "alicloud_dts_consumer_channel",
			Dependencies: []string{"alicloud_dts_instance"},
			F:            testSweepDtsConsumerChannel,
		})
}

// dtsSweepPrefixes is the canonical tf-test prefix list shared by every DTS
// sweeper. When ALICLOUD_SWEEP_ALL_RESOURCES is not set, only resources whose
// name begins with one of these prefixes are candidates for deletion.
var dtsSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
}

// dtsJobNameFromItem returns the DtsJobName field of a DescribeDtsJobs list
// item, or "" when the field is missing. Centralising the lookup keeps the
// per-sweeper loops free of repeated type-assertion boilerplate.
func dtsJobNameFromItem(item map[string]interface{}) string {
	if v, ok := item["DtsJobName"]; ok {
		return fmt.Sprint(v)
	}
	return ""
}

// dtsShouldDelete decides whether a listed resource should be deleted by the
// sweeper. It returns (true, "") when the resource should be deleted and
// (false, reason) when it must be skipped. The decision combines the
// name-based NODELETING protection (always enforced, even under SWEEP_ALL),
// the prefix match (enforced only when sweepAll() is false), and the
// age-based sweepRecentCheck filter.
func dtsShouldDelete(name string, item map[string]interface{}) (bool, string) {
	if protectedSweepVpcName(name) {
		return false, "NODELETING name marker"
	}
	if !sweepAll() {
		matches := false
		for _, prefix := range dtsSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				matches = true
				break
			}
		}
		if !matches {
			return false, "name does not match test prefixes"
		}
	}
	if skip, reason := sweepRecentCheck(item); skip {
		return false, reason
	}
	return true, ""
}

// listDtsJobs lists every DTS job of the given JobType (MIGRATION, SYNC or
// SUBSCRIBE) via DescribeDtsJobs. It paginates through the full result set and
// returns the raw list items for the caller to filter and delete.
func listDtsJobs(client *connectivity.AliyunClient, jobType string) ([]map[string]interface{}, error) {
	action := "DescribeDtsJobs"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"Region":     client.RegionId,
		"JobType":    jobType,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response = resp
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			return nil, fmt.Errorf("get $.DtsJobList for %s failed: %v", action, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				objects = append(objects, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}

// deleteDtsJob mirrors resourceAlicloudDtsMigrationJobDelete /
// resourceAliCloudDtsSubscriptionJobDelete / resourceAlicloudDtsSynchronizationJobDelete:
// all three call DeleteDtsJob with DtsJobId (+ DtsInstanceId when known).
func deleteDtsJob(client *connectivity.AliyunClient, dtsJobId, dtsInstanceId string) error {
	action := "DeleteDtsJob"
	request := map[string]interface{}{
		"DtsJobId": dtsJobId,
		"RegionId": client.RegionId,
	}
	if dtsInstanceId != "" {
		request["DtsInstanceId"] = dtsInstanceId
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil && !IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
		return err
	}
	return nil
}

// listDtsInstances lists every DTS instance via DescribeDtsInstances and
// returns the raw DtsInstanceStatusList items. alicloud_dts_instance,
// alicloud_dts_migration_instance and alicloud_dts_synchronization_instance
// are all created via CreateDtsInstance (with different Type values), so a
// single list call enumerates all three; each sweeper filters by Type.
func listDtsInstances(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeDtsInstances"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"PageNum":  1,
		"PageSize": PageSizeLarge,
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response = resp
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.DtsInstanceStatusList", response)
		if err != nil {
			return nil, fmt.Errorf("get $.DtsInstanceStatusList for %s failed: %v", action, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				objects = append(objects, item)
			}
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return objects, nil
}

// isDtsPrepayInstance reports whether a DescribeDtsInstances item is a
// Subscription (prepaid) instance. The provider's Delete handlers no-op on
// Subscription instances (see resourceAliCloudDtsInstanceDelete and
// resourceAlicloudDtsSynchronizationInstanceDelete), so sweepers preserve
// them too: releasing a prepaid instance via API has billing consequences the
// sweeper must not trigger silently.
func isDtsPrepayInstance(item map[string]interface{}) bool {
	switch fmt.Sprint(item["PayType"]) {
	case "PREPAY", "PrePaid", "Subscription":
		return true
	}
	return false
}

// testSweepDtsSynchronizationJob sweeps alicloud_dts_synchronization_job. It
// lists SYNC jobs via DescribeDtsJobs and deletes each candidate via
// DeleteDtsJob, mirroring resourceAlicloudDtsSynchronizationJobDelete.
func testSweepDtsSynchronizationJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := listDtsJobs(client, "SYNC")
	if err != nil {
		log.Printf("[ERROR] Failed to list DTS SYNC jobs: %s", err)
		return nil
	}
	for _, item := range items {
		name := dtsJobNameFromItem(item)
		if name == "" {
			continue
		}
		ok, reason := dtsShouldDelete(name, item)
		if !ok {
			log.Printf("[INFO] Skipping DTS SYNC Job (%s): %s", name, reason)
			continue
		}
		dtsInstanceId := fmt.Sprint(item["DtsInstanceID"])
		if err := deleteDtsJob(client, fmt.Sprint(item["DtsJobId"]), dtsInstanceId); err != nil {
			log.Printf("[ERROR] Failed to delete DTS SYNC Job (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete DTS SYNC Job success: %s", name)
	}
	return nil
}

// testSweepDtsInstance sweeps alicloud_dts_instance. It lists every DTS
// instance via DescribeDtsInstances and deletes each candidate via
// DeleteDtsJob with DtsInstanceId, mirroring resourceAliCloudDtsInstanceDelete.
// Subscription (PREPAY) instances are preserved because the provider Delete
// handler no-ops on them.
func testSweepDtsInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := listDtsInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list DTS instances: %s", err)
		return nil
	}
	for _, item := range items {
		name := fmt.Sprint(item["InstanceName"])
		if isDtsPrepayInstance(item) {
			log.Printf("[INFO] Skipping DTS Instance (Subscription/PREPAY): %s (%s)", name, fmt.Sprint(item["DtsInstanceId"]))
			continue
		}
		ok, reason := dtsShouldDelete(name, item)
		if !ok {
			log.Printf("[INFO] Skipping DTS Instance (%s): %s", name, reason)
			continue
		}
		action := "DeleteDtsJob"
		request := map[string]interface{}{
			"DtsInstanceId": item["DtsInstanceId"],
			"RegionId":      client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
			log.Printf("[ERROR] Failed to delete DTS Instance (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete DTS Instance success: %s", name)
	}
	return nil
}

// testSweepDtsMigrationInstance sweeps alicloud_dts_migration_instance. It
// lists DTS instances of Type=MIGRATION via DescribeDtsInstances and deletes
// each candidate via DeleteMigrationJob with MigrationJobId, mirroring
// resourceAlicloudDtsMigrationInstanceDelete.
func testSweepDtsMigrationInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := listDtsInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list DTS instances: %s", err)
		return nil
	}
	for _, item := range items {
		if fmt.Sprint(item["Type"]) != "MIGRATION" {
			continue
		}
		name := fmt.Sprint(item["InstanceName"])
		if isDtsPrepayInstance(item) {
			log.Printf("[INFO] Skipping DTS Migration Instance (Subscription/PREPAY): %s (%s)", name, fmt.Sprint(item["DtsInstanceId"]))
			continue
		}
		ok, reason := dtsShouldDelete(name, item)
		if !ok {
			log.Printf("[INFO] Skipping DTS Migration Instance (%s): %s", name, reason)
			continue
		}
		action := "DeleteMigrationJob"
		request := map[string]interface{}{
			"MigrationJobId": item["DtsInstanceId"],
			"RegionId":       client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
			log.Printf("[ERROR] Failed to delete DTS Migration Instance (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete DTS Migration Instance success: %s", name)
	}
	return nil
}

// testSweepDtsSynchronizationInstance sweeps
// alicloud_dts_synchronization_instance. It lists DTS instances of Type=SYNC
// via DescribeDtsInstances and deletes each candidate via
// DeleteSynchronizationJob with SynchronizationJobId, mirroring
// resourceAlicloudDtsSynchronizationInstanceDelete.
func testSweepDtsSynchronizationInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := listDtsInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list DTS instances: %s", err)
		return nil
	}
	for _, item := range items {
		if fmt.Sprint(item["Type"]) != "SYNC" {
			continue
		}
		name := fmt.Sprint(item["InstanceName"])
		if isDtsPrepayInstance(item) {
			log.Printf("[INFO] Skipping DTS Synchronization Instance (Subscription/PREPAY): %s (%s)", name, fmt.Sprint(item["DtsInstanceId"]))
			continue
		}
		ok, reason := dtsShouldDelete(name, item)
		if !ok {
			log.Printf("[INFO] Skipping DTS Synchronization Instance (%s): %s", name, reason)
			continue
		}
		action := "DeleteSynchronizationJob"
		request := map[string]interface{}{
			"SynchronizationJobId": item["DtsInstanceId"],
			"RegionId":             client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{"InvalidJobId", "Forbidden.InstanceNotFound"}) {
			log.Printf("[ERROR] Failed to delete DTS Synchronization Instance (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete DTS Synchronization Instance success: %s", name)
	}
	return nil
}

// testSweepDtsConsumerChannel sweeps alicloud_dts_consumer_channel. The
// DescribeConsumerChannel list API requires a DtsInstanceId, so the sweeper
// enumerates every DTS instance first and queries channels per instance.
// Each candidate is deleted via DeleteConsumerChannel with ConsumerGroupId
// and DtsInstanceId, mirroring resourceAlicloudDtsConsumerChannelDelete.
func testSweepDtsConsumerChannel(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instances, err := listDtsInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list DTS instances for consumer channel sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		dtsInstanceId := fmt.Sprint(inst["DtsInstanceId"])
		if dtsInstanceId == "" {
			continue
		}
		channels, err := listDtsConsumerChannels(client, dtsInstanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list consumer channels for DTS instance %s: %s", dtsInstanceId, err)
			continue
		}
		for _, ch := range channels {
			name := fmt.Sprint(ch["ConsumerGroupName"])
			ok, reason := dtsShouldDelete(name, ch)
			if !ok {
				log.Printf("[INFO] Skipping DTS Consumer Channel (%s): %s", name, reason)
				continue
			}
			consumerGroupId := fmt.Sprint(ch["ConsumerGroupId"])
			action := "DeleteConsumerChannel"
			request := map[string]interface{}{
				"ConsumerGroupId": consumerGroupId,
				"DtsInstanceId":   dtsInstanceId,
				"RegionId":        client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DTS Consumer Channel (%s): %s", name, err)
				continue
			}
			log.Printf("[INFO] Delete DTS Consumer Channel success: %s", name)
		}
	}
	return nil
}

// listDtsConsumerChannels lists every consumer channel attached to the given
// DTS instance via DescribeConsumerChannel. The API shape mirrors
// DescribeDtsConsumerChannel in service_alicloud_dts.go and the
// alicloud_dts_consumer_channels data source.
func listDtsConsumerChannels(client *connectivity.AliyunClient, dtsInstanceId string) ([]map[string]interface{}, error) {
	action := "DescribeConsumerChannel"
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"DtsInstanceId": dtsInstanceId,
		"PageSize":      PageSizeMedium,
		"PageNumber":    1,
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
					return retry.NonRetryableError(fmt.Errorf("not found: %s", dtsInstanceId))
				}
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response = resp
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.ConsumerChannels", response)
		if err != nil {
			return nil, fmt.Errorf("get $.ConsumerChannels for %s failed: %v", action, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				objects = append(objects, item)
			}
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}
