package alicloud

// Cloud Config (配置审计) sweeper gap fill.
//
// Product code: Config, API version: 2020-09-07 (RPC style for all gap
// resources below). Cloud Config is supported in cn-shanghai and
// cn-hangzhou only (connectivity.CloudConfigSupportedRegions); it is NOT
// available in cn-beijing. When the sweep fleet targets cn-beijing the list
// calls below will fail/empty-list, so every sweeper returns nil (log +
// continue) and baselines-0. That graceful no-op is intentional.
//
// Existing config sweepers (inline in their own *_test.go files) already
// cover: alicloud_config_rule, alicloud_config_aggregator,
// alicloud_config_compliance_pack, alicloud_config_aggregate_compliance_pack
// and alicloud_config_aggregate_config_rule. This file fills the remaining
// gaps: aggregate_delivery, aggregate_remediation, delivery, remediation and
// report_template.
//
// Two gap resources are intentionally DEFERRED (see report):
//   - alicloud_config_configuration_recorder: singleton per account/region,
//     TF Delete is a no-op ("Cannot destroy"), no delete API exists (only
//     Start/Stop ConfigurationRecorder). Sweeping does not apply.
//   - alicloud_config_delivery_channel: legacy 2019-01-08 delivery channel,
//     TF Delete is a no-op ("Cannot destroy"), no DeleteDeliveryChannel API
//     exists for the legacy version. The new alicloud_config_delivery
//     sweeper below (ListConfigDeliveryChannels + DeleteConfigDeliveryChannel
//     on 2020-09-07) sweeps the same underlying delivery channels regardless
//     of which API version created them, so legacy channels are covered.

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

const (
	configProduct = "Config"
	configVersion = "2020-09-07"
)

// configSweepPrefixes is the shared name-prefix allowlist used by the
// fleet-wipe gate when SWEEP_ALL is not active.
var configSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-test",
	"tf_test",
	"testAcc",
}

func init() {
	resource.AddTestSweepers("alicloud_config_aggregate_delivery", &resource.Sweeper{
		Name: "alicloud_config_aggregate_delivery",
		F:    testSweepConfigAggregateDelivery,
	})
	resource.AddTestSweepers("alicloud_config_aggregate_remediation", &resource.Sweeper{
		Name: "alicloud_config_aggregate_remediation",
		F:    testSweepConfigAggregateRemediation,
	})
	resource.AddTestSweepers("alicloud_config_delivery", &resource.Sweeper{
		Name: "alicloud_config_delivery",
		F:    testSweepConfigDelivery,
	})
	resource.AddTestSweepers("alicloud_config_remediation", &resource.Sweeper{
		Name: "alicloud_config_remediation",
		F:    testSweepConfigRemediation,
	})
	resource.AddTestSweepers("alicloud_config_report_template", &resource.Sweeper{
		Name: "alicloud_config_report_template",
		F:    testSweepConfigReportTemplate,
	})
}

// configSweepNameOk is the fleet-wipe gate shared by every config sweeper.
// The NODELETING guard (protectedSweepVpcName) is ALWAYS active in both
// modes and runs first; SWEEP_ALL bypasses the prefix check and cleans all
// non-NODELETING resources; otherwise a case-insensitive prefix match on the
// resource name/id is required.
func configSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range configSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// configRpcGet wraps RpcGet with retry.Retry + NeedRetry + incrementalWait.
func configRpcGet(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet(configProduct, configVersion, action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// configRpcPost wraps RpcPost with retry.Retry + NeedRetry + incrementalWait.
// autoRetry is forwarded to RpcPost to match each resource's delete convention.
func configRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost(configProduct, configVersion, action, query, body, autoRetry)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// configListAggregatorIds lists all config aggregators matching the fleet-wipe
// gate (parent scope for aggregate_delivery / aggregate_remediation). It is
// shared by the aggregate child sweepers to avoid duplicating the
// ListAggregators NextToken-pagination logic.
func configListAggregatorIds(client *connectivity.AliyunClient) []string {
	aggregatorIds := make([]string, 0)
	action := "ListAggregators"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	for {
		response, err := configRpcGet(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list config aggregators: %s", err)
			return aggregatorIds
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			log.Printf("[ERROR] Failed to get config aggregators path: %s", err)
			return aggregatorIds
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["AggregatorName"])
			if !configSweepNameOk(name) {
				log.Printf("[INFO] Skipping config aggregator: %s (%v)", name, item["AggregatorId"])
				continue
			}
			aggregatorIds = append(aggregatorIds, fmt.Sprint(item["AggregatorId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return aggregatorIds
}

// testSweepConfigDelivery sweeps alicloud_config_delivery (account/region
// scoped delivery channels on the 2020-09-07 API).
func testSweepConfigDelivery(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListConfigDeliveryChannels"
	response, err := configRpcPost(client, action, nil, map[string]interface{}{}, true)
	if err != nil {
		log.Printf("[ERROR] Failed to list config deliveries: %s", err)
		return nil
	}
	resp, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		log.Printf("[ERROR] Failed to get config delivery channels path: %s", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name := fmt.Sprint(item["DeliveryChannelName"])
		if !configSweepNameOk(name) {
			log.Printf("[INFO] Skipping config delivery: %s (%v)", name, item["DeliveryChannelId"])
			continue
		}
		deliveryChannelId := fmt.Sprint(item["DeliveryChannelId"])
		log.Printf("[INFO] Deleting config delivery: %s", deliveryChannelId)
		delAction := "DeleteConfigDeliveryChannel"
		delQuery := map[string]interface{}{
			"DeliveryChannelId": deliveryChannelId,
		}
		_, err = configRpcPost(client, delAction, delQuery, map[string]interface{}{}, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"DeliveryChannelNotExists"}) {
				log.Printf("[INFO] config delivery already deleted: %s", deliveryChannelId)
				continue
			}
			log.Printf("[ERROR] Failed to delete config delivery (%s): %s", deliveryChannelId, err)
			continue
		}
		log.Printf("[INFO] Delete config delivery success: %s", deliveryChannelId)
	}
	return nil
}

// testSweepConfigRemediation sweeps alicloud_config_remediation
// (account/region scoped remediation settings). Remediation has no
// user-settable name, so the gate matches on the auto-generated
// RemediationId (crr-xxx) which only deletes under SWEEP_ALL.
func testSweepConfigRemediation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListRemediations"
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := configRpcGet(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list config remediations: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Remediations", response)
		if err != nil {
			log.Printf("[ERROR] Failed to get config remediations path: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			remediationId := fmt.Sprint(item["RemediationId"])
			if !configSweepNameOk(remediationId) {
				log.Printf("[INFO] Skipping config remediation: %s", remediationId)
				continue
			}
			log.Printf("[INFO] Deleting config remediation: %s", remediationId)
			delAction := "DeleteRemediations"
			delBody := map[string]interface{}{
				"RemediationIds": remediationId,
			}
			_, err = configRpcPost(client, delAction, nil, delBody, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"RemediationConfigNotExist", "RemediationNotExist"}) {
					log.Printf("[INFO] config remediation already deleted: %s", remediationId)
					continue
				}
				log.Printf("[ERROR] Failed to delete config remediation (%s): %s", remediationId, err)
				continue
			}
			log.Printf("[INFO] Delete config remediation success: %s", remediationId)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// testSweepConfigReportTemplate sweeps alicloud_config_report_template
// (account/region scoped compliance report templates).
func testSweepConfigReportTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListReportTemplates"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	for {
		response, err := configRpcGet(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list config report templates: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.ReportTemplateList", response)
		if err != nil {
			log.Printf("[ERROR] Failed to get config report templates path: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["ReportTemplateName"])
			if !configSweepNameOk(name) {
				log.Printf("[INFO] Skipping config report template: %s (%v)", name, item["ReportTemplateId"])
				continue
			}
			reportTemplateId := fmt.Sprint(item["ReportTemplateId"])
			log.Printf("[INFO] Deleting config report template: %s", reportTemplateId)
			delAction := "DeleteReportTemplate"
			delBody := map[string]interface{}{
				"ReportTemplateId": reportTemplateId,
			}
			_, err = configRpcPost(client, delAction, map[string]interface{}{}, delBody, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ReportTemplateNotExist"}) {
					log.Printf("[INFO] config report template already deleted: %s", reportTemplateId)
					continue
				}
				log.Printf("[ERROR] Failed to delete config report template (%s): %s", reportTemplateId, err)
				continue
			}
			log.Printf("[INFO] Delete config report template success: %s", reportTemplateId)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// testSweepConfigAggregateDelivery sweeps alicloud_config_aggregate_delivery
// (parent-scoped: aggregator). It lists aggregators first, then per aggregator
// lists + deletes its aggregate delivery channels.
func testSweepConfigAggregateDelivery(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	aggregatorIds := configListAggregatorIds(client)
	for _, aggregatorId := range aggregatorIds {
		action := "ListAggregateConfigDeliveryChannels"
		request := map[string]interface{}{
			"AggregatorId": aggregatorId,
		}
		response, err := configRpcPost(client, action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to list config aggregate deliveries for aggregator %s: %s", aggregatorId, err)
			continue
		}
		resp, err := jsonpath.Get("$.DeliveryChannels", response)
		if err != nil {
			log.Printf("[ERROR] Failed to get config aggregate delivery channels path: %s", err)
			continue
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["DeliveryChannelName"])
			if !configSweepNameOk(name) {
				log.Printf("[INFO] Skipping config aggregate delivery: %s (%v)", name, item["DeliveryChannelId"])
				continue
			}
			deliveryChannelId := fmt.Sprint(item["DeliveryChannelId"])
			log.Printf("[INFO] Deleting config aggregate delivery: %s", deliveryChannelId)
			delAction := "DeleteAggregateConfigDeliveryChannel"
			delQuery := map[string]interface{}{
				"AggregatorId":      aggregatorId,
				"DeliveryChannelId": deliveryChannelId,
			}
			_, err = configRpcPost(client, delAction, delQuery, map[string]interface{}{}, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"DeliveryChannelNotExists", "Invalid.AggregatorId.Value"}) {
					log.Printf("[INFO] config aggregate delivery already deleted: %s", deliveryChannelId)
					continue
				}
				log.Printf("[ERROR] Failed to delete config aggregate delivery (%s): %s", deliveryChannelId, err)
				continue
			}
			log.Printf("[INFO] Delete config aggregate delivery success: %s", deliveryChannelId)
		}
	}
	return nil
}

// testSweepConfigAggregateRemediation sweeps
// alicloud_config_aggregate_remediation (parent-scoped: aggregator). It lists
// aggregators first, then per aggregator lists + deletes its aggregate
// remediation settings. Remediation has no user-settable name, so the gate
// matches on the auto-generated RemediationId and only deletes under
// SWEEP_ALL.
func testSweepConfigAggregateRemediation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	aggregatorIds := configListAggregatorIds(client)
	for _, aggregatorId := range aggregatorIds {
		action := "ListAggregateRemediations"
		request := map[string]interface{}{
			"AggregatorId": aggregatorId,
		}
		response, err := configRpcGet(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list config aggregate remediations for aggregator %s: %s", aggregatorId, err)
			continue
		}
		resp, err := jsonpath.Get("$.Remediations", response)
		if err != nil {
			log.Printf("[ERROR] Failed to get config aggregate remediations path: %s", err)
			continue
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			remediationId := fmt.Sprint(item["RemediationId"])
			if !configSweepNameOk(remediationId) {
				log.Printf("[INFO] Skipping config aggregate remediation: %s", remediationId)
				continue
			}
			log.Printf("[INFO] Deleting config aggregate remediation: %s", remediationId)
			delAction := "DeleteAggregateRemediations"
			delBody := map[string]interface{}{
				"RemediationIds": remediationId,
				"AggregatorId":   aggregatorId,
			}
			_, err = configRpcPost(client, delAction, map[string]interface{}{}, delBody, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"RemediationNotExist", "Invalid.AggregatorId.Value"}) {
					log.Printf("[INFO] config aggregate remediation already deleted: %s", remediationId)
					continue
				}
				log.Printf("[ERROR] Failed to delete config aggregate remediation (%s): %s", remediationId, err)
				continue
			}
			log.Printf("[INFO] Delete config aggregate remediation success: %s", remediationId)
		}
	}
	return nil
}
